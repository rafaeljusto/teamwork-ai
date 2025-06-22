package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

// Resources stores the resources used by different applications in the Teamwork
// AI ecosystem.
type Resources struct {
	Logger         *slog.Logger
	Agentic        agentic.Agentic
	TeamworkEngine interface {
		Do(context.Context, twapi.Entity, ...twapi.Option) error
	}
}

// InitResources creates a new set of resources for the many applications in the
// Teamwork AI ecosystem.
func InitResources(ctx context.Context, config *Config) (*Resources, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LoggerLevel,
	}))

	var mcpClient *agentic.MCPClient
	var mcpClientOptions []agentic.MCPOption
	if config.Agentic.MCPClient.SSEAddress != "" {
		mcpClientOptions = append(mcpClientOptions, agentic.WithMCPSSE(config.Agentic.MCPClient.SSEAddress))
	} else if config.Agentic.MCPClient.StdioPath != "" {
		mcpClientOptions = append(mcpClientOptions, agentic.WithMCPStdio(
			config.Agentic.MCPClient.StdioPath,
			config.Agentic.MCPClient.StdioEnvs,
			config.Agentic.MCPClient.StdioArgs...,
		))
	}
	if len(mcpClientOptions) > 0 {
		var err error
		if mcpClient, err = agentic.ConnectToMCP(ctx, logger, mcpClientOptions...); err != nil {
			return nil, fmt.Errorf("failed to connect to MCP: %w", err)
		}
	}

	resources := &Resources{
		Logger:         logger,
		Agentic:        agentic.Init(config.Agentic.Name, config.Agentic.DSN, mcpClient, logger),
		TeamworkEngine: twapi.NewEngine(config.TeamworkServer, config.TeamworkAPIToken, logger),
	}

	return resources, nil
}
