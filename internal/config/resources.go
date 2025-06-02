package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

// Resources stores the resources for the web server.
type Resources struct {
	Logger         *slog.Logger
	Agentic        agentic.Agentic
	TeamworkEngine interface {
		Do(context.Context, twapi.Entity, ...twapi.Option) error
	}
}

// NewResources creates a new set of resources for the web server.
func NewResources(config *Config) *Resources {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LoggerLevel,
	}))
	resources := &Resources{
		Logger:         logger,
		Agentic:        agentic.Init(config.Agentic.Name, config.Agentic.DSN, logger),
		TeamworkEngine: twapi.NewEngine(config.TeamworkServer, config.TeamworkAPIToken, logger),
	}

	return resources
}
