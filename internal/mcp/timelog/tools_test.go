package timelog_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/timelog"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

func TestTools_retrieveTimelogs(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-timelogs"
	request.Params.Arguments = map[string]any{
		"tag-ids":        []float64{1, 2, 3},
		"match-all-tags": true,
		"page":           float64(1),
		"page-size":      float64(10),
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

func TestTools_retrieveProjectTimelogs(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-project-timelogs"
	request.Params.Arguments = map[string]any{
		"project-id":     float64(123),
		"tag-ids":        []float64{1, 2, 3},
		"match-all-tags": true,
		"page":           float64(1),
		"page-size":      float64(10),
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

func TestTools_retrieveTaskTimelogs(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-task-timelogs"
	request.Params.Arguments = map[string]any{
		"task-id":        float64(123),
		"tag-ids":        []float64{1, 2, 3},
		"match-all-tags": true,
		"page":           float64(1),
		"page-size":      float64(10),
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

func TestTools_retrieveTimelog(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-timelog"
	request.Params.Arguments = map[string]any{
		"timelog-id": float64(123),
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

func TestTools_createTimelog(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-timelog"
	request.Params.Arguments = map[string]any{
		"description": "Example timelog description",
		"date":        "2023-12-31",
		"time":        "12:00:00",
		"is-utc":      true,
		"hours":       float64(1),
		"minutes":     float64(30),
		"billable":    true,
		"project-id":  float64(123),
		"task-id":     float64(456),
		"user-id":     float64(789),
		"tag-ids":     []float64{10, 11, 12},
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

func TestTools_updateTimelog(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timelog.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-timelog"
	request.Params.Arguments = map[string]any{
		"timelog-id":  float64(123),
		"description": "Example timelog description",
		"date":        "2023-12-31",
		"time":        "12:00:00",
		"is-utc":      true,
		"hours":       float64(1),
		"minutes":     float64(30),
		"billable":    true,
		"project-id":  float64(123),
		"task-id":     float64(456),
		"user-id":     float64(789),
		"tag-ids":     []float64{10, 11, 12},
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
