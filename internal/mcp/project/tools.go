package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twproject "github.com/rafaeljusto/teamwork-ai/internal/teamwork/project"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-projects",
			mcp.WithDescription("Retrieve multiple projects in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var projects twproject.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &projects); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(projects)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project",
			mcp.WithDescription("Retrieve a specific project in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&project.ID, "project-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(project)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-project",
			mcp.WithDescription("Create a new project in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the project."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the project."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Creation

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&project.Name, "name"),
				twmcp.OptionalParam(&project.Description, "description"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Project created successfully"), nil
		},
	)
}
