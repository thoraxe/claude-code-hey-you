package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration
type Config struct {
	Topic    string `toml:"topic"`
	Server   string `toml:"server"`
	Priority string `toml:"priority"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Server:   "https://ntfy.sh",
		Priority: "default",
	}
}

// LoadConfig loads configuration from CLI flags, environment variables, and config file.
// Priority: CLI flags > env vars > config file > defaults
func LoadConfig() Config {
	cfg := DefaultConfig()

	// Load from config file first (lowest priority)
	loadConfigFile(&cfg)

	// Override with environment variables
	loadEnvVars(&cfg)

	// Override with CLI flags (highest priority)
	loadFlags(&cfg)

	return cfg
}

// configFilePath returns the platform-appropriate config file path
func configFilePath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "claude-ntfy", "config.toml")
		}
	}

	// Unix-like: ~/.config/claude-ntfy/config.toml
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "claude-ntfy", "config.toml")
}

// loadConfigFile loads configuration from the config file if it exists
func loadConfigFile(cfg *Config) {
	path := configFilePath()
	if path == "" {
		return
	}

	var fileCfg Config
	if _, err := toml.DecodeFile(path, &fileCfg); err != nil {
		// File doesn't exist or can't be parsed - that's fine, it's optional
		return
	}

	// Apply non-empty values from file
	if fileCfg.Topic != "" {
		cfg.Topic = fileCfg.Topic
	}
	if fileCfg.Server != "" {
		cfg.Server = fileCfg.Server
	}
	if fileCfg.Priority != "" {
		cfg.Priority = fileCfg.Priority
	}
}

// loadEnvVars loads configuration from environment variables
func loadEnvVars(cfg *Config) {
	if topic := os.Getenv("NTFY_TOPIC"); topic != "" {
		cfg.Topic = topic
	}
	if server := os.Getenv("NTFY_SERVER"); server != "" {
		cfg.Server = server
	}
	if priority := os.Getenv("NTFY_PRIORITY"); priority != "" {
		cfg.Priority = priority
	}
}

// loadFlags loads configuration from CLI flags
func loadFlags(cfg *Config) {
	topic := flag.String("topic", "", "ntfy topic to publish to")
	server := flag.String("server", "", "ntfy server URL (default: https://ntfy.sh)")
	priority := flag.String("priority", "", "notification priority (min, low, default, high, urgent)")

	flag.Parse()

	if *topic != "" {
		cfg.Topic = *topic
	}
	if *server != "" {
		cfg.Server = *server
	}
	if *priority != "" {
		cfg.Priority = *priority
	}
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	if c.Topic == "" {
		return &ConfigError{"topic is required (use --topic, NTFY_TOPIC env var, or config file)"}
	}
	return nil
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
