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
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tasklists by name."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtasklist.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple)
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
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tasklists by name."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtasklist.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&multiple.Request.Path.ProjectID, "project-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple)
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
				mcp.Description("The ID of the tasklist."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twtasklist.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&single.ID, "tasklist-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &single); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(single)
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
			mcp.WithNumber("milestone-id",
				mcp.Description("The ID of the milestone to associate with the tasklist."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Create

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&tasklist.Name, "name"),
				twmcp.RequiredNumericParam(&tasklist.ProjectID, "project-id"),
				twmcp.OptionalPointerParam(&tasklist.Description, "description"),
				twmcp.OptionalNumericPointerParam(&tasklist.MilestoneID, "milestone-id"),
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

	mcpServer.AddTool(
		mcp.NewTool("update-tasklist",
			mcp.WithDescription("Update an existing tasklist in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithNumber("tasklist-id",
				mcp.Required(),
				mcp.Description("The ID of the tasklist to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the tasklist."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the tasklist."),
			),
			mcp.WithNumber("milestone-id",
				mcp.Description("The ID of the milestone to associate with the tasklist."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&tasklist.ID, "tasklist-id"),
				twmcp.OptionalPointerParam(&tasklist.Name, "name"),
				twmcp.OptionalPointerParam(&tasklist.Description, "description"),
				twmcp.OptionalNumericPointerParam(&tasklist.MilestoneID, "milestone-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Tasklist updated successfully"), nil
		},
	)
}
