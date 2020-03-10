package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config represents service configuration for github-bots
type Config struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                ":8085",
		GracefulShutdownTimeout: 5 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
