package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

// Config stores the configuration of the application.
type Config struct {
	// Port is the port of the server.
	Port int64

	// LoggerLevel is the level of the logger.
	LoggerLevel slog.Level

	// TeamworkServer is the server of the Teamwork API.
	TeamworkServer string

	// TeamworkAPIToken is the API token of the Teamwork API.
	TeamworkAPIToken string

	// Agentic is the agentic configuration.
	Agentic struct {
		// Name is the name of the agentic implementation.
		Name string

		// DSN is the data source name for the agentic model. The format depends on
		// the chosen implementation.
		DSN string
	}
}

// ParseFromEnvs parses the configuration from environment variables.
func ParseFromEnvs() (*Config, error) {
	var config Config
	var errs error
	var err error

	portStr := os.Getenv("TWAI_PORT")
	if portStr != "" {
		config.Port, err = strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to parse TWAI_PORT: %w", err))
		}
	}

	loggerLevel := slog.LevelInfo
	if loggerLevelStr := os.Getenv("TWAI_LOG_LEVEL"); loggerLevelStr != "" {
		if err = loggerLevel.UnmarshalText([]byte(loggerLevelStr)); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to parse TWAI_LOG_LEVEL: %w", err))
		}
	}
	config.LoggerLevel = loggerLevel

	config.TeamworkServer = os.Getenv("TWAI_TEAMWORK_SERVER")
	config.TeamworkAPIToken = os.Getenv("TWAI_TEAMWORK_API_TOKEN")

	config.Agentic.Name = os.Getenv("TWAI_AGENTIC_NAME")
	config.Agentic.DSN = os.Getenv("TWAI_AGENTIC_DSN")

	if errs != nil {
		return nil, errs
	}
	return &config, nil
}
