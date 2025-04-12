package skill_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/skill"
)

func TestResources_skills(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	skill.Register(mcpServer, &config.Resources{
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
	request.Params.URI = "twapi://skills"

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

func TestResources_skill(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	skill.Register(mcpServer, &config.Resources{
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
	request.Params.URI = "twapi://skills/123"

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
