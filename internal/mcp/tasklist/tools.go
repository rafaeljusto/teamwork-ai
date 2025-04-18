package tasklist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twtasklist "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasklist"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasklists",
			mcp.WithDescription("Retrieve multiple tasklists in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklists twtasklist.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &tasklists); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklists)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project-tasklists",
			mcp.WithDescription("Retrieve multiple tasklists from a specific project in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasklists."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklists twtasklist.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&tasklists.ProjectID, "project-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasklists); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklists)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasklist",
			mcp.WithDescription("Retrieve a specific tasklist in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithNumber("tasklist-id",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&tasklist.ID, "tasklist-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklist)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-tasklist",
			mcp.WithDescription("Create a new tasklist in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the tasklist."),
			),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the tasklist."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Creation

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&tasklist.Name, "name"),
				twmcp.RequiredNumericParam(&tasklist.ProjectID, "project-id"),
				twmcp.OptionalParam(&tasklist.Description, "description"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Tasklist created successfully"), nil
		},
	)
}
