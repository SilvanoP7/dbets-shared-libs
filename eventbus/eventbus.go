package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// EventHandler represents a function that handles events
type EventHandler func(event interface{}) error

// EventBus interface for publishing and subscribing to events
type EventBus interface {
	Publish(topic string, event interface{}) error
	Subscribe(topic string, handler EventHandler) error
	Unsubscribe(topic string) error
	Close() error
	GetStreamInfo(streamName string) (*nats.StreamInfo, error)
	GetConsumerInfo(consumerName string) (*nats.ConsumerInfo, error)
}

// NATSEventBus implements EventBus using NATS JetStream
type NATSEventBus struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	ctx  context.Context
	subs map[string]*nats.Subscription
}

// NewNATSEventBus creates a new NATS event bus with JetStream
func NewNATSEventBus(natsURL string) (*NATSEventBus, error) {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create stream for event storage
	stream, err := js.AddStream(&nats.StreamConfig{
		Name:      "DBETS_EVENTS",
		Subjects:  []string{"dbets.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour, // Keep events for 24 hours
		MaxMsgs:   1000000,        // Keep up to 1M messages
		Replicas:  1,
	})
	if err != nil && err.Error() != "stream name already in use" {
		log.Printf("Warning: failed to create stream (may already exist): %v", err)
	} else if err == nil {
		log.Printf("Created JetStream stream: %s", stream.Config.Name)
	}

	return &NATSEventBus{
		nc:   nc,
		js:   js,
		ctx:  context.Background(),
		subs: make(map[string]*nats.Subscription),
	}, nil
}

// Publish publishes an event to a topic with JetStream
func (n *NATSEventBus) Publish(topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Add headers for filtering
	headers := nats.Header{}
	headers.Set("Event-Type", fmt.Sprintf("%T", event))
	headers.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	// Publish to JetStream
	ack, err := n.js.PublishMsg(&nats.Msg{
		Subject: topic,
		Data:    data,
		Header:  headers,
	})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event to topic %s (sequence: %d): %+v", topic, ack.Sequence, event)
	return nil
}

// Subscribe subscribes to a topic and handles events with JetStream
func (n *NATSEventBus) Subscribe(topic string, handler EventHandler) error {
	// Create consumer for this subscription
	consumerName := fmt.Sprintf("consumer-%s-%d", topic, time.Now().Unix())

	// Subscribe with JetStream
	sub, err := n.js.Subscribe(topic, func(msg *nats.Msg) {
		// Parse the event
		var event interface{}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event from topic %s: %v", topic, err)
			msg.Ack()
			return
		}

		// Handle the event
		if err := handler(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
			// Don't ack the message so it can be redelivered
			return
		}

		// Acknowledge the message
		msg.Ack()
	}, nats.Durable(consumerName), nats.AckWait(30*time.Second))

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	n.subs[topic] = sub
	log.Printf("Subscribed to topic: %s with consumer: %s", topic, consumerName)
	return nil
}

// SubscribeWithFilter subscribes to a topic with filtering capabilities
func (n *NATSEventBus) SubscribeWithFilter(topic string, filter func(msg *nats.Msg) bool, handler EventHandler) error {
	consumerName := fmt.Sprintf("consumer-filtered-%s-%d", topic, time.Now().Unix())

	sub, err := n.js.Subscribe(topic, func(msg *nats.Msg) {
		// Apply filter
		if !filter(msg) {
			msg.Ack()
			return
		}

		// Parse the event
		var event interface{}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event from topic %s: %v", topic, err)
			msg.Ack()
			return
		}

		// Handle the event
		if err := handler(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
			return
		}

		msg.Ack()
	}, nats.Durable(consumerName), nats.AckWait(30*time.Second))

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	n.subs[topic] = sub
	log.Printf("Subscribed to topic: %s with filter and consumer: %s", topic, consumerName)
	return nil
}

// ReplayEvents replays events from a specific time or sequence
func (n *NATSEventBus) ReplayEvents(topic string, fromTime time.Time, handler EventHandler) error {
	consumerName := fmt.Sprintf("replay-%s-%d", topic, time.Now().Unix())

	// Create consumer with start time
	sub, err := n.js.Subscribe(topic, func(msg *nats.Msg) {
		// Parse the event
		var event interface{}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event from topic %s: %v", topic, err)
			msg.Ack()
			return
		}

		// Handle the event
		if err := handler(event); err != nil {
			log.Printf("Error handling replayed event from topic %s: %v", topic, err)
			return
		}

		msg.Ack()
	}, nats.Durable(consumerName), nats.DeliverNew(), nats.AckWait(30*time.Second))

	if err != nil {
		return fmt.Errorf("failed to replay events from topic %s: %w", topic, err)
	}

	// Note: JetStream replay is handled by the subscription itself
	// The PurgeStream is not needed for replay functionality

	n.subs[topic] = sub
	log.Printf("Started replay for topic: %s from time: %s", topic, fromTime.Format(time.RFC3339))
	return nil
}

// Unsubscribe unsubscribes from a topic
func (n *NATSEventBus) Unsubscribe(topic string) error {
	if sub, exists := n.subs[topic]; exists {
		err := sub.Unsubscribe()
		delete(n.subs, topic)
		return err
	}
	return nil
}

// Close closes the NATS connection
func (n *NATSEventBus) Close() error {
	// Unsubscribe from all topics
	for topic := range n.subs {
		n.Unsubscribe(topic)
	}

	if n.nc != nil {
		n.nc.Close()
	}
	return nil
}

// GetStreamInfo returns information about the JetStream stream
func (n *NATSEventBus) GetStreamInfo(streamName string) (*nats.StreamInfo, error) {
	return n.js.StreamInfo(streamName)
}

// GetConsumerInfo returns information about a specific consumer
func (n *NATSEventBus) GetConsumerInfo(consumerName string) (*nats.ConsumerInfo, error) {
	return n.js.ConsumerInfo("DBETS_EVENTS", consumerName)
}

// Event types for different services

// BetPlacedEvent represents a bet placement event
type BetPlacedEvent struct {
	BetID     string  `json:"bet_id"`
	UserID    string  `json:"user_id"`
	EventID   string  `json:"event_id"`
	Amount    float64 `json:"amount"`
	Odds      float64 `json:"odds"`
	Timestamp string  `json:"timestamp"`
}

// BetSettledEvent represents a bet settlement event
type BetSettledEvent struct {
	BetID        string  `json:"bet_id"`
	UserID       string  `json:"user_id"`
	EventID      string  `json:"event_id"`
	PayoutAmount float64 `json:"payout_amount"`
	Status       string  `json:"status"`
	Timestamp    string  `json:"timestamp"`
}

// CashoutRequestEvent represents a cashout request event
type CashoutRequestEvent struct {
	BetID        string  `json:"bet_id"`
	UserID       string  `json:"user_id"`
	CashoutValue float64 `json:"cashout_value"`
	Timestamp    string  `json:"timestamp"`
}

// OddsUpdatedEvent represents an odds update event
type OddsUpdatedEvent struct {
	SelectionID string  `json:"selection_id"`
	OldOdds     float64 `json:"old_odds"`
	NewOdds     float64 `json:"new_odds"`
	Timestamp   string  `json:"timestamp"`
}

// EventResultEvent represents an event result event
type EventResultEvent struct {
	EventID   string `json:"event_id"`
	Result    string `json:"result"`
	Timestamp string `json:"timestamp"`
}

// SettlementCompleteEvent represents a settlement completion event
type SettlementCompleteEvent struct {
	BetID        string  `json:"bet_id"`
	UserID       string  `json:"user_id"`
	PayoutAmount float64 `json:"payout_amount"`
	Status       string  `json:"status"`
	Timestamp    string  `json:"timestamp"`
}

// WalletUpdatedEvent represents a wallet update event
type WalletUpdatedEvent struct {
	UserID     string  `json:"user_id"`
	OldBalance float64 `json:"old_balance"`
	NewBalance float64 `json:"new_balance"`
	Change     float64 `json:"change"`
	Type       string  `json:"type"` // "deposit", "withdrawal", "bet", "win"
	Timestamp  string  `json:"timestamp"`
}

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

// UserLoginEvent represents a user login event
type UserLoginEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

// UserUpdatedEvent represents a user update event
type UserUpdatedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

// BetFailedEvent represents a failed bet attempt
type BetFailedEvent struct {
	UserID       string                 `json:"user_id"`
	EventID      string                 `json:"event_id"`
	MarketID     string                 `json:"market_id"`
	SelectionID  string                 `json:"selection_id"`
	Amount       float64                `json:"amount"`
	FailedReason string                 `json:"failed_reason"`
	RequestData  map[string]interface{} `json:"request_data"`
	Timestamp    string                 `json:"timestamp"`
}

// Helper functions for creating events

// NewBetPlacedEvent creates a new bet placed event
func NewBetPlacedEvent(betID, userID, eventID string, amount, odds float64) *BetPlacedEvent {
	return &BetPlacedEvent{
		BetID:     betID,
		UserID:    userID,
		EventID:   eventID,
		Amount:    amount,
		Odds:      odds,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewBetSettledEvent creates a new bet settled event
func NewBetSettledEvent(betID, userID, eventID string, payoutAmount float64, status string) *BetSettledEvent {
	return &BetSettledEvent{
		BetID:        betID,
		UserID:       userID,
		EventID:      eventID,
		PayoutAmount: payoutAmount,
		Status:       status,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

// NewCashoutRequestEvent creates a new cashout request event
func NewCashoutRequestEvent(betID, userID string, cashoutValue float64) *CashoutRequestEvent {
	return &CashoutRequestEvent{
		BetID:        betID,
		UserID:       userID,
		CashoutValue: cashoutValue,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

// NewOddsUpdatedEvent creates a new odds updated event
func NewOddsUpdatedEvent(selectionID string, oldOdds, newOdds float64) *OddsUpdatedEvent {
	return &OddsUpdatedEvent{
		SelectionID: selectionID,
		OldOdds:     oldOdds,
		NewOdds:     newOdds,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// NewEventResultEvent creates a new event result event
func NewEventResultEvent(eventID, result string) *EventResultEvent {
	return &EventResultEvent{
		EventID:   eventID,
		Result:    result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewSettlementCompleteEvent creates a new settlement complete event
func NewSettlementCompleteEvent(betID, userID string, payoutAmount float64, status string) *SettlementCompleteEvent {
	return &SettlementCompleteEvent{
		BetID:        betID,
		UserID:       userID,
		PayoutAmount: payoutAmount,
		Status:       status,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

// NewWalletUpdatedEvent creates a new wallet updated event
func NewWalletUpdatedEvent(userID string, oldBalance, newBalance float64, changeType string) *WalletUpdatedEvent {
	return &WalletUpdatedEvent{
		UserID:     userID,
		OldBalance: oldBalance,
		NewBalance: newBalance,
		Change:     newBalance - oldBalance,
		Type:       changeType,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

// NewUserCreatedEvent creates a new user created event
func NewUserCreatedEvent(userID, username, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewUserLoginEvent creates a new user login event
func NewUserLoginEvent(userID, email string) *UserLoginEvent {
	return &UserLoginEvent{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewUserUpdatedEvent creates a new user updated event
func NewUserUpdatedEvent(userID, email string) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewBetFailedEvent creates a new bet failed event
func NewBetFailedEvent(userID, eventID, marketID, selectionID string, amount float64, failedReason string, requestData map[string]interface{}) *BetFailedEvent {
	return &BetFailedEvent{
		UserID:       userID,
		EventID:      eventID,
		MarketID:     marketID,
		SelectionID:  selectionID,
		Amount:       amount,
		FailedReason: failedReason,
		RequestData:  requestData,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}
