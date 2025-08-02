package eventbus

import (
	"fmt"
	"log"
)

// NewEventBus creates a new NATS event bus
func NewEventBus(config *Config) (EventBus, error) {
	if config == nil {
		config = LoadConfigFromEnv()
	}

	// Only support NATS now
	log.Println("Creating NATS JetStream event bus")
	return NewNATSEventBus(config.GetNATSUrl())
}

// NewEventBusFromConfigFile creates a new event bus from a configuration file
func NewEventBusFromConfigFile(configPath string) (EventBus, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return NewEventBus(config)
}

// NewEventBusFromEnv creates a new event bus from environment variables
func NewEventBusFromEnv() (EventBus, error) {
	config := LoadConfigFromEnv()
	return NewEventBus(config)
}
