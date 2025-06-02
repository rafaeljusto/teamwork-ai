package tag_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/tag"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

func TestTools_retrieveTags(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-tags"
	request.Params.Arguments = map[string]any{
		"search-term": "test",
		"item-type":   "tasks",
		"project-ids": []int64{1, 2, 3},
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

func TestTools_retrieveTag(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-tag"
	request.Params.Arguments = map[string]any{
		"tag-id": float64(123),
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

func TestTools_createTag(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-tag"
	request.Params.Arguments = map[string]any{
		"name":       "Example",
		"project-id": float64(456),
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

func TestTools_updateTag(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	tag.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-tag"
	request.Params.Arguments = map[string]any{
		"tag-id":     float64(123),
		"name":       "Example",
		"project-id": float64(456),
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
