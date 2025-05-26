package timelog_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/timelog"
)

func TestResources_timelogs(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.URI = "twapi://timelogs"

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

func TestResources_timelog(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.URI = "twapi://timelogs/123"

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
