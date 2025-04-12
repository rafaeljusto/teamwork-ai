package milestone_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/milestone"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

func TestTools_retrieveMilestones(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	milestone.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-milestones"
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

func TestTools_retrieveProjectMilestones(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	milestone.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-project-milestones"
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

func TestTools_retrievemilestone(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	milestone.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-milestone"
	request.Params.Arguments = map[string]any{
		"milestone-id": float64(123),
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

func TestTools_createMilestone(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	milestone.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-milestone"
	request.Params.Arguments = map[string]any{
		"name":        "Example",
		"description": "Example milestone description",
		"due-date":    "20231231",
		"assignees": map[string]any{
			"user-ids":    []float64{1, 2, 3},
			"company-ids": []float64{4, 5},
			"team-ids":    []float64{6, 7},
		},
		"tasklist-ids": []float64{8, 9},
		"tag-ids":      []float64{10, 11, 12},
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

func TestTools_updateMilestone(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	milestone.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-milestone"
	request.Params.Arguments = map[string]any{
		"milestone-id": float64(123),
		"name":         "Example",
		"description":  "Example milestone description",
		"due-date":     "20231231",
		"assignees": map[string]any{
			"user-ids":    []float64{1, 2, 3},
			"company-ids": []float64{4, 5},
			"team-ids":    []float64{6, 7},
		},
		"tasklist-ids": []float64{8, 9},
		"tag-ids":      []float64{10, 11, 12},
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
