package config

import (
	"log/slog"
	"os"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
	twapi "github.com/teamwork/twapi-go-sdk"
	"github.com/teamwork/twapi-go-sdk/session"
)

// Resources stores the resources for the web server.
type Resources struct {
	Logger         *slog.Logger
	Agentic        agentic.Agentic
	TeamworkEngine *twapi.Engine
}

// NewResources creates a new set of resources for the web server.
func NewResources(config *Config) *Resources {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LoggerLevel,
	}))
	resources := &Resources{
		Logger:  logger,
		Agentic: agentic.Init(config.Agentic.Name, config.Agentic.DSN, logger),
		TeamworkEngine: twapi.NewEngine(
			session.NewBasicAuth(config.TeamworkAPIToken, "", config.TeamworkServer),
			twapi.WithLogger(logger),
		),
	}

	return resources
}
