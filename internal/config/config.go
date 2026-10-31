// Package config reads the service's settings from the environment.
package config

import (
	"errors"
	"fmt"
	"time"
)

// Config holds the settings the service runs with.
type Config struct {
	Port         string
	DrainTimeout time.Duration
}

// Load reads Config through getenv, using a default for each unset value.
func Load(getenv func(string) string) (Config, error) {
	port := getenv("PORT")
	if port == "" {
		port = "8080"
	}

	drain, err := drainTimeout(getenv("DRAIN_TIMEOUT"))
	if err != nil {
		return Config{}, err
	}

	return Config{Port: port, DrainTimeout: drain}, nil
}

func drainTimeout(raw string) (time.Duration, error) {
	if raw == "" {
		return 15 * time.Second, nil
	}

	drain, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("parse DRAIN_TIMEOUT: %w", err)
	}

	if drain <= 0 {
		return 0, errors.New("DRAIN_TIMEOUT must be positive")
	}

	return drain, nil
}
