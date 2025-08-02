package eventbus

import (
	"testing"
	"time"
)

func TestNewEventBusFromEnv(t *testing.T) {
	// Test creating event bus from environment
	eventBus, err := NewEventBusFromEnv()
	if err != nil {
		t.Skipf("Skipping test - failed to create event bus: %v", err)
	}
	defer eventBus.Close()

	// Test that we can create an event bus
	if eventBus == nil {
		t.Fatal("Expected event bus to be created")
	}
}

func TestEventCreation(t *testing.T) {
	// Test event creation functions
	betEvent := NewBetPlacedEvent("test-bet", "test-user", "test-event", 100.0, 2.5)
	if betEvent.BetID != "test-bet" {
		t.Errorf("Expected BetID to be 'test-bet', got %s", betEvent.BetID)
	}

	userEvent := NewUserCreatedEvent("test-user", "testuser", "test@example.com")
	if userEvent.UserID != "test-user" {
		t.Errorf("Expected UserID to be 'test-user', got %s", userEvent.UserID)
	}

	walletEvent := NewWalletUpdatedEvent("test-user", 100.0, 150.0, "deposit")
	if walletEvent.UserID != "test-user" {
		t.Errorf("Expected UserID to be 'test-user', got %s", walletEvent.UserID)
	}
}

func TestConfigLoading(t *testing.T) {
	// Test configuration loading
	config := LoadConfigFromEnv()
	if config == nil {
		t.Fatal("Expected config to be created")
	}

	// Test default values
	if config.GetNATSUrl() == "" {
		t.Error("Expected NATS URL to have a default value")
	}
}

func TestEventBusFactory(t *testing.T) {
	// Test factory function
	eventBus, err := NewEventBus(nil) // Use default config
	if err != nil {
		t.Skipf("Skipping test - failed to create event bus: %v", err)
	}
	defer eventBus.Close()

	if eventBus == nil {
		t.Fatal("Expected event bus to be created")
	}
}

func TestEventTimestamps(t *testing.T) {
	// Test that events have proper timestamps
	event := NewBetPlacedEvent("test-bet", "test-user", "test-event", 100.0, 2.5)

	// Parse the timestamp
	_, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		t.Errorf("Expected valid RFC3339 timestamp, got error: %v", err)
	}
}

func TestEventTypes(t *testing.T) {
	// Test all event types
	events := []interface{}{
		NewBetPlacedEvent("bet-1", "user-1", "event-1", 100.0, 2.0),
		NewBetSettledEvent("bet-1", "user-1", "event-1", 200.0, "won"),
		NewCashoutRequestEvent("bet-1", "user-1", 150.0),
		NewOddsUpdatedEvent("selection-1", 2.0, 2.5),
		NewEventResultEvent("event-1", "home_win"),
		NewSettlementCompleteEvent("bet-1", "user-1", 200.0, "completed"),
		NewWalletUpdatedEvent("user-1", 100.0, 150.0, "deposit"),
		NewUserCreatedEvent("user-1", "username", "user@example.com"),
		NewUserLoginEvent("user-1", "user@example.com"),
		NewUserUpdatedEvent("user-1", "user@example.com"),
	}

	for i, event := range events {
		if event == nil {
			t.Errorf("Event %d is nil", i)
		}
	}
}
