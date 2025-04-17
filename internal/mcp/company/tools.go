package company

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twcompany "github.com/rafaeljusto/teamwork-ai/internal/teamwork/company"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-companies",
			mcp.WithDescription("Retrieve multiple companies, also know as clients, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var companies twcompany.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &companies); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(companies)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-company",
			mcp.WithDescription("Retrieve a specific company, also know as client, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithNumber("companyId",
				mcp.Required(),
				mcp.Description("The ID of the company."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var company twcompany.Single

			err := twmcp.NumericParam(request.Params.Arguments, &company.ID, "companyId")
			if err != nil {
				return nil, fmt.Errorf("invalid company ID: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &company); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(company)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-company",
			mcp.WithDescription("Create a new company, also know as client, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the company."),
			),
			mcp.WithArray("userIds",
				mcp.Description("List of user IDs assigned to the company. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var company twcompany.Creation

			err := twmcp.Param(request.Params.Arguments, &company.Name, "name")
			if err != nil {
				return nil, fmt.Errorf("invalid name: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &company); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Company created successfully"), nil
		},
	)
}
