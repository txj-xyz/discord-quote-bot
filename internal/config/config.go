package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Bot struct {
		Token   string `yaml:"token"`
		GuildID string `yaml:"guild_id,omitempty"` // Optional guild ID for dev commands
	} `yaml:"bot"`
	Quote struct {
		Channel string `yaml:"channel"` // Channel ID to post quotes
	} `yaml:"quote"`
}

// LoadConfig loads the configuration from the provided YAML file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate required fields
	if config.Bot.Token == "" {
		return nil, fmt.Errorf("bot token is required in config")
	}

	if config.Quote.Channel == "" {
		return nil, fmt.Errorf("quote channel ID is required in config")
	}

	return &config, nil
} 