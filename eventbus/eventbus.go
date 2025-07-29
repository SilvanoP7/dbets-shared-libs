package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// EventHandler represents a function that handles events
type EventHandler func(event interface{}) error

// EventBus interface for publishing and subscribing to events
type EventBus interface {
	Publish(topic string, event interface{}) error
	Subscribe(topic string, handler EventHandler) error
	Unsubscribe(topic string) error
	Close() error
}

// RedisEventBus implements EventBus using Redis Pub/Sub
type RedisEventBus struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisEventBus creates a new Redis event bus
func NewRedisEventBus(redisURL string) (*RedisEventBus, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisEventBus{
		client: client,
		ctx:    ctx,
	}, nil
}

// Publish publishes an event to a topic
func (r *RedisEventBus) Publish(topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.client.Publish(r.ctx, topic, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event to topic %s: %+v", topic, event)
	return nil
}

// Subscribe subscribes to a topic and handles events
func (r *RedisEventBus) Subscribe(topic string, handler EventHandler) error {
	pubsub := r.client.Subscribe(r.ctx, topic)
	defer pubsub.Close()

	log.Printf("Subscribed to topic: %s", topic)

	for {
		msg, err := pubsub.ReceiveMessage(r.ctx)
		if err != nil {
			log.Printf("Error receiving message from topic %s: %v", topic, err)
			continue
		}

		// Parse the event (you might want to add event type information)
		var event interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Error unmarshaling event from topic %s: %v", topic, err)
			continue
		}

		// Handle the event
		if err := handler(event); err != nil {
			log.Printf("Error handling event from topic %s: %v", topic, err)
		}
	}
}

// Unsubscribe unsubscribes from a topic
func (r *RedisEventBus) Unsubscribe(topic string) error {
	pubsub := r.client.Subscribe(r.ctx, topic)
	return pubsub.Unsubscribe(r.ctx, topic)
}

// Close closes the Redis connection
func (r *RedisEventBus) Close() error {
	return r.client.Close()
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
