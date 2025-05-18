package team_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/team"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

func TestTools_retrieveTeams(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	team.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-teams"
	request.Params.Arguments = map[string]any{
		"search-term": "test",
		"page":        float64(1),
		"page-size":   float64(10),
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

func TestTools_retrieveTeam(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	team.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-team"
	request.Params.Arguments = map[string]any{
		"team-id": float64(123),
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

func TestTools_createTeam(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	team.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-team"
	request.Params.Arguments = map[string]any{
		"name":           "Example",
		"handle":         "example",
		"description":    "Example description",
		"parent-team-id": float64(123),
		"company-id":     float64(456),
		"project-id":     float64(789),
		"user-ids": []any{
			float64(1),
			float64(2),
			float64(3),
		},
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

func TestTools_updateTeam(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	team.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-team"
	request.Params.Arguments = map[string]any{
		"team-id":        float64(123),
		"name":           "Example",
		"handle":         "example",
		"description":    "Example description",
		"parent-team-id": float64(123),
		"company-id":     float64(456),
		"project-id":     float64(789),
		"user-ids": []any{
			float64(1),
			float64(2),
			float64(3),
		},
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
