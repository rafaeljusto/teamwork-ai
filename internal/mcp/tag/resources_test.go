package tag_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/tag"
)

func TestResources_tags(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
		TeamworkEngine: engineMock{},
	})

	request := &resourceRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		ReadResourceRequest: mcp.ReadResourceRequest{
			Request: mcp.Request{
				Method: string(mcp.MethodResourcesRead),
			},
		},
	}
	request.Params.URI = "twapi://tags"

	encodedRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	ctx := context.Background()
	message := mcpServer.HandleMessage(ctx, encodedRequest)
	if err, ok := message.(mcp.JSONRPCError); ok {
		t.Errorf("tool failed to execute: %v", err.Error)
	}
}

func TestResources_tag(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
		TeamworkEngine: engineMock{},
	})

	request := &resourceRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		ReadResourceRequest: mcp.ReadResourceRequest{
			Request: mcp.Request{
				Method: string(mcp.MethodResourcesRead),
			},
		},
	}
	request.Params.URI = "twapi://tags/123"

	encodedRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	ctx := context.Background()
	message := mcpServer.HandleMessage(ctx, encodedRequest)
	if err, ok := message.(mcp.JSONRPCError); ok {
		t.Errorf("tool failed to execute: %v", err.Error)
	}
}

type resourceRequest struct {
	mcp.ReadResourceRequest

	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
}
