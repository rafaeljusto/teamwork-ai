package activity_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/activity"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

func TestTools_retrieveActivities(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	activity.Register(mcpServer, &config.Resources{
		TeamworkEngine: engineMock{},
	})

	request := &toolRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		CallToolRequest: mcp.CallToolRequest{
			Request: mcp.Request{
				Method: string(mcp.MethodToolsCall),
			},
		},
	}
	request.Params.Name = "retrieve-activities"
	request.Params.Arguments = map[string]any{
		"start-date": "2023-10-01T00:00:00Z",
		"end-date":   "2023-10-31T23:59:59Z",
		"log-item-types": []any{
			"message",
			"comment",
			"task",
			"tasklist",
			"taskgroup",
			"milestone",
			"file",
		},
		"page":      float64(1),
		"page-size": float64(10),
	}

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

func TestTools_retrieveProjectActivities(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	activity.Register(mcpServer, &config.Resources{
		TeamworkEngine: engineMock{},
	})

	request := &toolRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      1,
		CallToolRequest: mcp.CallToolRequest{
			Request: mcp.Request{
				Method: string(mcp.MethodToolsCall),
			},
		},
	}
	request.Params.Name = "retrieve-project-activities"
	request.Params.Arguments = map[string]any{
		"project-id": float64(123),
		"start-date": "2023-10-01T00:00:00Z",
		"end-date":   "2023-10-31T23:59:59Z",
		"log-item-types": []any{
			"message",
			"comment",
			"task",
			"tasklist",
			"taskgroup",
			"milestone",
			"file",
		},
		"page":      float64(1),
		"page-size": float64(10),
	}

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

type toolRequest struct {
	mcp.CallToolRequest

	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
}

type engineMock struct {
}

func (e engineMock) Do(context.Context, twapi.Entity, ...twapi.Option) error {
	return nil
}
