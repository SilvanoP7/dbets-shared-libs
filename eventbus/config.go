package eventbus

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the event bus configuration
type Config struct {
	EventBus struct {
		NATS struct {
			URL string `yaml:"url"`
			JetStream struct {
				StreamName string   `yaml:"stream_name"`
				Subjects   []string `yaml:"subjects"`
				Storage    string   `yaml:"storage"`
				Retention  string   `yaml:"retention"`
				MaxAge     string   `yaml:"max_age"`
				MaxMsgs    int      `yaml:"max_msgs"`
				Replicas   int      `yaml:"replicas"`
			} `yaml:"jetstream"`
			Consumer struct {
				AckWait       string `yaml:"ack_wait"`
				MaxDeliver    int    `yaml:"max_deliver"`
				FilterDurable bool   `yaml:"filter_durable"`
			} `yaml:"consumer"`
		} `yaml:"nats"`
		Events struct {
			Retention map[string]string `yaml:"retention"`
			Filters   struct {
				Enabled       bool     `yaml:"enabled"`
				IncludeHeaders []string `yaml:"include_headers"`
			} `yaml:"filters"`
			Replay struct {
				Enabled           bool   `yaml:"enabled"`
				MaxReplayDuration string `yaml:"max_replay_duration"`
				BatchSize         int    `yaml:"batch_size"`
			} `yaml:"replay"`
		} `yaml:"events"`
		Logging struct {
			Level            string `yaml:"level"`
			Format           string `yaml:"format"`
			IncludeSequence  bool   `yaml:"include_sequence"`
			IncludeTimestamp bool   `yaml:"include_timestamp"`
		} `yaml:"logging"`
	} `yaml:"eventbus"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{}
	
	// Set defaults
	config.EventBus.NATS.URL = getEnv("NATS_URL", "nats://localhost:4222")
	
	return config
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetNATSUrl returns the NATS URL from config or environment
func (c *Config) GetNATSUrl() string {
	if c.EventBus.NATS.URL != "" {
		return c.EventBus.NATS.URL
	}
	return getEnv("NATS_URL", "nats://localhost:4222")
} 