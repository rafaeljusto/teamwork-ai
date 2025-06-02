package task_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/task"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

func TestTools_retrieveTasks(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-tasks"
	request.Params.Arguments = map[string]any{
		"search-term":    "test",
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

func TestTools_retrieveProjectTasks(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-project-tasks"
	request.Params.Arguments = map[string]any{
		"project-id":     float64(123),
		"search-term":    "test",
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

func TestTools_retrieveTasklistTasks(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-tasklist-tasks"
	request.Params.Arguments = map[string]any{
		"tasklist-id":    float64(123),
		"search-term":    "test",
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

func TestTools_retrieveTask(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-task"
	request.Params.Arguments = map[string]any{
		"task-id": float64(123),
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

func TestTools_createTask(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-task"
	request.Params.Arguments = map[string]any{
		"name":              "Example",
		"tasklist-id":       float64(123),
		"description":       "This is an example task.",
		"priority":          "high",
		"progress":          float64(50),
		"start-date":        "2023-10-01",
		"due-date":          "2023-10-15",
		"estimated-minutes": float64(120),
		"assignees": map[string]any{
			"user-ids":    []float64{1, 2, 3},
			"team-ids":    []float64{4, 5},
			"company-ids": []float64{6, 7},
		},
		"tag-ids": []float64{1, 2, 3},
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

func TestTools_updateTask(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	task.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-task"
	request.Params.Arguments = map[string]any{
		"task-id":           float64(123),
		"name":              "Example",
		"tasklist-id":       float64(123),
		"description":       "This is an example task.",
		"priority":          "high",
		"progress":          float64(50),
		"start-date":        "2023-10-01",
		"due-date":          "2023-10-15",
		"estimated-minutes": float64(120),
		"assignees": map[string]any{
			"user-ids":    []float64{1, 2, 3},
			"team-ids":    []float64{4, 5},
			"company-ids": []float64{6, 7},
		},
		"tag-ids": []float64{1, 2, 3},
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
