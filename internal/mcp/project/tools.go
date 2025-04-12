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
			mcp.WithString("searchTerm",
				mcp.Description("A search term to filter projects by name or description."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter projects by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match projects that have all the specified tags. "+
					"If false, the search will match projects that have any of the specified tags. "+
					"Defaults to false."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twproject.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.TagIDs, "tag-ids"),
				twmcp.OptionalPointerParam(&multiple.Request.Filters.MatchAllTags, "match-all-tags"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Response)
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
				mcp.Description("The ID of the project."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twproject.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&single.ID, "project-id"),
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
			mcp.WithString("start-at",
				mcp.Description("The start date of the project in the format YYYYMMDD."),
			),
			mcp.WithString("end-at",
				mcp.Description("The end date of the project in the format YYYYMMDD."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company associated with the project."),
			),
			mcp.WithNumber("owned-id",
				mcp.Description("The ID of the user who owns the project."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the project."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Create

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&project.Name, "name"),
				twmcp.OptionalPointerParam(&project.Description, "description"),
				twmcp.OptionalLegacyDatePointerParam(&project.StartAt, "start-at"),
				twmcp.OptionalLegacyDatePointerParam(&project.EndAt, "end-at"),
				twmcp.OptionalNumericParam(&project.CompanyID, "company-id"),
				twmcp.OptionalNumericPointerParam(&project.OwnerID, "owned-id"),
				twmcp.OptionalNumericListParam(&project.Tags, "tag-ids"),
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

	mcpServer.AddTool(
		mcp.NewTool("update-project",
			mcp.WithDescription("Update an existing project in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the project."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the project."),
			),
			mcp.WithString("start-at",
				mcp.Description("The start date of the project in the format YYYYMMDD."),
			),
			mcp.WithString("end-at",
				mcp.Description("The end date of the project in the format YYYYMMDD."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company associated with the project."),
			),
			mcp.WithNumber("owned-id",
				mcp.Description("The ID of the user who owns the project."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the project."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&project.ID, "project-id"),
				twmcp.OptionalPointerParam(&project.Name, "name"),
				twmcp.OptionalPointerParam(&project.Description, "description"),
				twmcp.OptionalLegacyDatePointerParam(&project.StartAt, "start-at"),
				twmcp.OptionalLegacyDatePointerParam(&project.EndAt, "end-at"),
				twmcp.OptionalNumericPointerParam(&project.CompanyID, "company-id"),
				twmcp.OptionalNumericPointerParam(&project.OwnerID, "owned-id"),
				twmcp.OptionalNumericListParam(&project.Tags, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Project updated successfully"), nil
		},
	)
}
