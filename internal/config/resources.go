package config

import (
	"log/slog"
	"os"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Resources stores the resources for the web server.
type Resources struct {
	Logger         *slog.Logger
	TeamworkEngine *teamwork.Engine
}

// NewResources creates a new set of resources for the web server.
func NewResources(config *Config) *Resources {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LoggerLevel,
	}))
	resources := &Resources{
		Logger:         logger,
		TeamworkEngine: teamwork.NewEngine(config.TeamworkServer, config.TeamworkAPIToken, logger),
	}

	return resources
}
