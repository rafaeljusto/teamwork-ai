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
	resources := &Resources{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.LoggerLevel,
		})),
		TeamworkEngine: teamwork.NewEngine(config.TeamworkServer, config.TeamworkAPIToken),
	}

	return resources
}
