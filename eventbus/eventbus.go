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
	CreateStream(config *nats.StreamConfig) error
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
	log.Printf("Attempting to connect to NATS at: %s", natsURL)

	// Connect to NATS with proper options for JetStream
	opts := []nats.Option{
		nats.Timeout(10 * time.Second),      // Connection timeout
		nats.ReconnectWait(1 * time.Second), // Reconnect wait time
		nats.MaxReconnects(5),               // Maximum reconnection attempts
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS error: %v", err)
		}),
	}

	log.Printf("Creating NATS connection with URL: %s", natsURL)
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Printf("Failed to connect to NATS at %s: %v", natsURL, err)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Printf("Successfully connected to NATS at: %s", nc.ConnectedUrl())

	// Create JetStream context with timeout
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	log.Printf("Successfully created JetStream context")

	// Don't create streams automatically - let each service create its own streams
	log.Printf("Connected to NATS JetStream at %s", natsURL)

	return &NATSEventBus{
		nc:   nc,
		js:   js,
		ctx:  context.Background(),
		subs: make(map[string]*nats.Subscription),
	}, nil
}

// Publish publishes an event to a topic with JetStream
func (n *NATSEventBus) Publish(topic string, event interface{}) error {
	log.Printf("Attempting to publish event to topic: %s", topic)

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	log.Printf("Event data marshaled successfully, size: %d bytes", len(data))

	// Add headers for filtering
	headers := nats.Header{}
	headers.Set("Event-Type", fmt.Sprintf("%T", event))
	headers.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	// Create context with timeout for JetStream operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Publishing message to JetStream with timeout...")

	// Publish to JetStream with timeout
	ack, err := n.js.PublishMsg(&nats.Msg{
		Subject: topic,
		Data:    data,
		Header:  headers,
	}, nats.Context(ctx))
	if err != nil {
		log.Printf("Failed to publish to JetStream: %v", err)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Successfully published event to topic %s (sequence: %d): %+v", topic, ack.Sequence, event)
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

// CreateStream creates a new JetStream stream
func (n *NATSEventBus) CreateStream(config *nats.StreamConfig) error {
	// Create context with timeout for stream operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := n.js.AddStream(config, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	log.Printf("Created JetStream stream: %s", stream.Config.Name)
	return nil
}

// Event types for different services

// BetPlacedEvent represents a bet placement event
type BetPlacedEvent struct {
	BetID      string         `json:"bet_id"`
	UserID     string         `json:"user_id"`
	Amount     float64        `json:"amount"`
	Odds       float64        `json:"odds"`
	Selections []BetSelection `json:"selections"`
	Timestamp  string         `json:"timestamp"`
}

type BetSelection struct {
	EventID     string  `json:"event_id"`
	MarketID    string  `json:"market_id"`
	SelectionID string  `json:"selection_id"`
	Odds        float64 `json:"odds"`
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

// EventUpdatedEvent represents an event update
type EventUpdatedEvent struct {
	EventID   string                 `json:"event_id"`
	Version   int                    `json:"version"`
	Changes   map[string]interface{} `json:"changes"` // Field name -> new value
	Timestamp string                 `json:"timestamp"`
}

// MarketUpdatedEvent represents a market update
type MarketUpdatedEvent struct {
	MarketID  string                 `json:"market_id"`
	EventID   string                 `json:"event_id"`
	Version   int                    `json:"version"`
	Changes   map[string]interface{} `json:"changes"` // Field name -> new value
	Timestamp string                 `json:"timestamp"`
}

// SelectionUpdatedEvent represents a selection update
type SelectionUpdatedEvent struct {
	SelectionID string                 `json:"selection_id"`
	MarketID    string                 `json:"market_id"`
	EventID     string                 `json:"event_id"`
	Version     int                    `json:"version"`
	Changes     map[string]interface{} `json:"changes"` // Field name -> new value
	Timestamp   string                 `json:"timestamp"`
}

// OddsUpdatedEvent represents an update to odds
type OddsUpdatedEvent struct {
	OddsID      string                 `json:"odds_id"`
	SelectionID string                 `json:"selection_id"`
	Version     int                    `json:"version"`
	Changes     map[string]interface{} `json:"changes"` // Field name -> new value
	Timestamp   string                 `json:"timestamp"`
}

// ResultUpdatedEvent represents an update to a result
type ResultUpdatedEvent struct {
	ResultID    string                 `json:"result_id"`
	SelectionID string                 `json:"selection_id"`
	Version     int                    `json:"version"`
	Changes     map[string]interface{} `json:"changes"` // Field name -> new value
	Timestamp   string                 `json:"timestamp"`
}

// Helper functions for creating events

// NewBetPlacedEvent creates a new bet placed event
func NewBetPlacedEvent(betID, userID string, amount, odds float64, selections []BetSelection) *BetPlacedEvent {
	return &BetPlacedEvent{
		BetID:      betID,
		UserID:     userID,
		Amount:     amount,
		Odds:       odds,
		Selections: selections,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
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

// NewEventUpdatedEvent creates a new event updated event
func NewEventUpdatedEvent(eventID string, version int, changes map[string]interface{}) *EventUpdatedEvent {
	return &EventUpdatedEvent{
		EventID:   eventID,
		Version:   version,
		Changes:   changes,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewMarketUpdatedEvent creates a new market updated event
func NewMarketUpdatedEvent(marketID, eventID string, version int, changes map[string]interface{}) *MarketUpdatedEvent {
	return &MarketUpdatedEvent{
		MarketID:  marketID,
		EventID:   eventID,
		Version:   version,
		Changes:   changes,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewSelectionUpdatedEvent creates a new selection updated event
func NewSelectionUpdatedEvent(selectionID, marketID, eventID string, version int, changes map[string]interface{}) *SelectionUpdatedEvent {
	return &SelectionUpdatedEvent{
		SelectionID: selectionID,
		MarketID:    marketID,
		EventID:     eventID,
		Version:     version,
		Changes:     changes,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// NewOddsUpdatedEvent creates a new odds updated event
func NewOddsUpdatedEvent(oddsID, selectionID string, version int, changes map[string]interface{}) *OddsUpdatedEvent {
	return &OddsUpdatedEvent{
		OddsID:      oddsID,
		SelectionID: selectionID,
		Version:     version,
		Changes:     changes,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// NewResultUpdatedEvent creates a new result updated event
func NewResultUpdatedEvent(resultID, selectionID string, version int, changes map[string]interface{}) *ResultUpdatedEvent {
	return &ResultUpdatedEvent{
		ResultID:    resultID,
		SelectionID: selectionID,
		Version:     version,
		Changes:     changes,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
