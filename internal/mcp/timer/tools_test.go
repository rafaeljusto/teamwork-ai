package timer_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/timer"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

func TestTools_retrieveTimers(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-timers"
	request.Params.Arguments = map[string]any{
		"user-id":             float64(123),
		"task-id":             float64(456),
		"project-id":          float64(789),
		"running-timers-only": true,
		"page":                float64(1),
		"page-size":           float64(10),
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

func TestTools_retrieveTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-timer"
	request.Params.Arguments = map[string]any{
		"timer-id": float64(123),
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

func TestTools_createTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-timer"
	request.Params.Arguments = map[string]any{
		"description":         "Example timer description",
		"billable":            true,
		"running":             true,
		"seconds":             float64(3600), // 1 hour
		"stop-running-timers": true,
		"project-id":          float64(123),
		"task-id":             float64(456),
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

func TestTools_updateTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-timer"
	request.Params.Arguments = map[string]any{
		"timer-id":    float64(123),
		"description": "Example timer description",
		"billable":    true,
		"running":     true,
		"project-id":  float64(123),
		"task-id":     float64(456),
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

func TestTools_pauseTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "pause-timer"
	request.Params.Arguments = map[string]any{
		"timer-id": float64(123),
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

func TestTools_completeTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "complete-timer"
	request.Params.Arguments = map[string]any{
		"timer-id": float64(123),
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

func TestTools_resumeTimer(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	timer.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "resume-timer"
	request.Params.Arguments = map[string]any{
		"timer-id": float64(123),
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

func (e engineMock) Do(context.Context, teamwork.Entity, ...teamwork.Option) error {
	return nil
}
