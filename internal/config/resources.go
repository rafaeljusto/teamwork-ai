package config

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
	twapi "github.com/teamwork/twapi-go-sdk"
	"github.com/teamwork/twapi-go-sdk/session"
)

// Resources stores the resources for the web server.
type Resources struct {
	Logger         *slog.Logger
	Agentic        agentic.Agentic
	TeamworkEngine *twapi.Engine
	MCPClient      *MCPClient
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
			session.NewBearerToken(config.TeamworkAPIToken, config.TeamworkServer),
			twapi.WithLogger(logger),
		),
		MCPClient: NewMCPClient(&mcp.StreamableClientTransport{
			Endpoint: config.MCPEndpoint,
			HTTPClient: &http.Client{
				Transport: &authTransport{token: config.TeamworkAPIToken},
			},
		}),
	}

	return resources
}

// MCPClient wraps an MCP client with a specific endpoint.
type MCPClient struct {
	client    *mcp.Client
	transport mcp.Transport
}

// NewMCPClient creates a new MCP client with the given endpoint.
func NewMCPClient(transport mcp.Transport) *MCPClient {
	return &MCPClient{
		client: mcp.NewClient(&mcp.Implementation{
			Name:    "teamwork-ai",
			Title:   "Teamwork AI",
			Version: "1.0.0",
		}, &mcp.ClientOptions{}),
		transport: transport,
	}
}

// Connect connects to the MCP server with the given options.
func (m *MCPClient) Connect(ctx context.Context) (*mcp.ClientSession, error) {
	return m.client.Connect(ctx, m.transport, &mcp.ClientSessionOptions{})
}

type authTransport struct {
	token string
}

// RoundTrip adds the basic auth header to the request.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
