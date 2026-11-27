package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Environment variable names read by the gateway.
// Names live in constants so a typo is a compile error, not a silent empty value.
const (
	envAddr            = "GATEWAY_ADDR"
	envAPIKey          = "GATEWAY_API_KEY"
	envLogLevel        = "GATEWAY_LOG_LEVEL"
	envShutdownTimeout = "GATEWAY_SHUTDOWN_TIMEOUT"
)

// Defaults used when an optional variable is empty.
// There is no default for the API key: a secret must never have a fallback.
const (
	defaultAddr            = ":8080"
	defaultLogLevel        = "info"
	defaultShutdownTimeout = "10s"
)

// config holds every setting the gateway reads from the environment.
type config struct {
	Addr            string
	APIKey          string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// loadConfig reads the settings through getenv, applies defaults and validates them.
// It collects every problem and returns them together with errors.Join, so a
// misconfigured deploy shows all mistakes at once instead of one per restart.
// getenv is os.Getenv in production and a map lookup in tests.
func loadConfig(getenv func(string) string) (config, error) {
	var errs []error

	cfg := config{
		Addr:   valueOr(getenv(envAddr), defaultAddr),
		APIKey: getenv(envAPIKey),
	}

	if cfg.APIKey == "" {
		errs = append(errs, fmt.Errorf("%s is required", envAPIKey))
	}

	level := valueOr(getenv(envLogLevel), defaultLogLevel)
	if err := cfg.LogLevel.UnmarshalText([]byte(level)); err != nil {
		errs = append(errs, fmt.Errorf("%s: %w", envLogLevel, err))
	}

	timeout, err := time.ParseDuration(valueOr(getenv(envShutdownTimeout), defaultShutdownTimeout))
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("%s: %w", envShutdownTimeout, err))
	case timeout <= 0:
		errs = append(errs, fmt.Errorf("%s must be positive, got %s", envShutdownTimeout, timeout))
	default:
		cfg.ShutdownTimeout = timeout
	}

	if err := errors.Join(errs...); err != nil {
		return config{}, err
	}
	return cfg, nil

}

// valueOr returns value, or fallback when value is empty.
func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
