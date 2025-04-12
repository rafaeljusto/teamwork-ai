package company_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/mcp/company"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

func TestTools_retrieveCompanies(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	company.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-companies"
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

func TestTools_retrieveCompany(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	company.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "retrieve-company"
	request.Params.Arguments = map[string]any{
		"company-id": float64(123),
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

func TestTools_createCompany(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	company.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "create-company"
	request.Params.Arguments = map[string]any{
		"name":         "Example",
		"address-one":  "123 Example St",
		"address-two":  "Suite 456",
		"city":         "Example City",
		"state":        "EX",
		"zip":          "12345",
		"country-code": "US",
		"phone":        "123-456-7890",
		"fax":          "098-765-4321",
		"email-one":    "example1@test.com",
		"email-two":    "example2@test.com",
		"email-three":  "example3@test.com",
		"website":      "https://www.example.com",
		"profile":      "Example Company Profile",
		"manager-id":   float64(456),
		"industry-id":  float64(789),
		"tag-ids":      []float64{1, 2, 3},
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

func TestTools_updateCompany(t *testing.T) {
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	company.Register(mcpServer, &config.Resources{
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
	request.Params.Name = "update-company"
	request.Params.Arguments = map[string]any{
		"company-id":   float64(123),
		"name":         "Example",
		"address-one":  "123 Example St",
		"address-two":  "Suite 456",
		"city":         "Example City",
		"state":        "EX",
		"zip":          "12345",
		"country-code": "US",
		"phone":        "123-456-7890",
		"fax":          "098-765-4321",
		"email-one":    "example1@test.com",
		"email-two":    "example2@test.com",
		"email-three":  "example3@test.com",
		"website":      "https://www.example.com",
		"profile":      "Example Company Profile",
		"manager-id":   float64(456),
		"industry-id":  float64(789),
		"tag-ids":      []float64{1, 2, 3},
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
