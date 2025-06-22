package agentic

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
)

const (
	mcpClientName    = "Teamwork Agent"
	mcpClientVersion = "1.0.0"
)

// MCPOptions contains the options for the MCP client.
type MCPOptions struct {
	stdioPath string
	stdioEnvs []string
	stdioArgs []string
	sseURL    string
}

// MCPOption is a function that modifies the MCPOptions struct. It allows for
// optional configuration of the MCP client.
type MCPOption func(*MCPOptions)

// WithMCPStdio sets the path to the local MCP server binary to connect via
// stdio mode.
func WithMCPStdio(path string, envs []string, args ...string) MCPOption {
	return func(o *MCPOptions) {
		o.stdioPath = path
		o.stdioEnvs = envs
		o.stdioArgs = args
	}
}

// WithMCPSSE sets the URL to connect to the MCP server via SSE mode.
func WithMCPSSE(url string) MCPOption {
	return func(o *MCPOptions) {
		o.sseURL = url
	}
}

// MCPClient is a wrapper around the MCP client. It stores the client and the
// server information.
type MCPClient struct {
	client     *client.Client
	serverInfo *mcp.InitializeResult
}

// ConnectToMCP connects to the MCP server and returns the client. By default it
// will attempt to connect to a stdios MCP server using the path "teamwork-mcp".
func ConnectToMCP(ctx context.Context, logger *slog.Logger, optFunc ...MCPOption) (*MCPClient, error) {
	options := MCPOptions{
		stdioPath: "teamwork-mcp",
	}
	for _, opt := range optFunc {
		opt(&options)
	}

	var mcpTransport transport.Interface
	switch {
	case options.stdioPath != "":
		mcpTransport = transport.NewStdio(options.stdioPath, options.stdioEnvs, options.stdioArgs...)

	case options.sseURL != "":
		var err error
		if mcpTransport, err = transport.NewSSE(options.sseURL); err != nil {
			return nil, fmt.Errorf("failed to create SSE transport for URL %q: %w", options.sseURL, err)
		}

	default:
		return nil, fmt.Errorf("no transport method specified")
	}

	mcpClient := client.NewClient(mcpTransport)
	if err := mcpClient.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start MCP client: %w", err)
	}

	if options.stdioPath != "" {
		if stderr, ok := client.GetStderr(mcpClient); ok {
			// TODO(rafaeljusto): Add a better control to stop the goroutine. Maybe we
			// don't need to care about this since it could be sending an io.EOF when
			// the transport is closed.
			go func(stderr io.Reader) {
				buffer := make([]byte, 4096)
				for {
					n, err := stderr.Read(buffer)
					if err != nil {
						if err != io.EOF {
							logger.Error("failed to read from stderr",
								slog.String("error", err.Error()),
							)
						}
						return
					}
					if n > 0 {
						logger.Error("mcp server error",
							slog.String("error", string(buffer[:n])),
						)
					}
				}
			}(stderr)
		}
	}

	mcpClient.OnNotification(func(notification mcp.JSONRPCNotification) {
		logger.Info("MCP notification",
			slog.String("method", notification.Method),
		)
	})

	mcpServerInfo, err := mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    mcpClientName,
				Version: mcpClientVersion,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}

	logger.Info("MCP server info",
		slog.String("name", mcpServerInfo.ServerInfo.Name),
		slog.String("version", mcpServerInfo.ServerInfo.Version),
		slog.String("protocolVersion", mcpServerInfo.ProtocolVersion),
	)

	return &MCPClient{
		client:     mcpClient,
		serverInfo: mcpServerInfo,
	}, nil
}

// Tools returns the list of tools available in the MCP server. It's possible to
// filter the tools by methods.
func (m *MCPClient) Tools(ctx context.Context, methods ...twmcp.Method) ([]mcp.Tool, error) {
	if m.client == nil {
		return nil, fmt.Errorf("MCP client is not initialized")
	}

	if m.serverInfo == nil || m.serverInfo.Capabilities.Tools == nil {
		return nil, fmt.Errorf("MCP server does not support tools")
	}

	var toolsRequest mcp.ListToolsRequest
	result, err := m.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	if len(methods) > 0 {
		var i int
		for _, tool := range result.Tools {
			if slices.Contains(methods, twmcp.Method(tool.Name)) {
				result.Tools[i] = tool
				i++
			}
		}
		result.Tools = result.Tools[:i]
	}

	return result.Tools, nil
}

// ExecuteTool executes a tool with the given parameters.
func (m *MCPClient) ExecuteTool(
	ctx context.Context,
	method string,
	params mcp.CallToolParams,
) (*mcp.CallToolResult, error) {
	if m.client == nil {
		return nil, fmt.Errorf("MCP client is not initialized")
	}

	if m.serverInfo == nil || m.serverInfo.Capabilities.Tools == nil {
		return nil, fmt.Errorf("MCP server does not support tools")
	}

	toolResult, err := m.client.CallTool(ctx, mcp.CallToolRequest{
		Request: mcp.Request{
			Method: method,
		},
		Params: params,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute tool %q: %w", method, err)
	}

	return toolResult, nil
}

// Resources returns the list of resources available in the MCP server.
func (m *MCPClient) Resources(ctx context.Context) ([]mcp.Resource, error) {
	if m.client == nil {
		return nil, fmt.Errorf("MCP client is not initialized")
	}

	if m.serverInfo == nil || m.serverInfo.Capabilities.Resources == nil {
		return nil, fmt.Errorf("MCP server does not support resources")
	}

	var resourcesRequest mcp.ListResourcesRequest
	result, err := m.client.ListResources(ctx, resourcesRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return result.Resources, nil
}

// Close closes the MCP client connection.
func (m *MCPClient) Close() error {
	if m.client == nil {
		return fmt.Errorf("MCP client is not initialized")
	}

	if err := m.client.Close(); err != nil {
		return fmt.Errorf("failed to close MCP client: %w", err)
	}
	m.client = nil
	m.serverInfo = nil
	return nil
}
