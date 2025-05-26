package tag

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twtag "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tag"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-tags",
			mcp.WithDescription("Retrieve multiple tags in a customer site of Teamwork.com. "+
				"Tags are a way to mark items so that you can use a filter to see just those items."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tags by name. "+
					"Each word from the search term is used to match against the tag name."),
			),
			mcp.WithString("item-type",
				mcp.Description("The type of item to filter tags by. Valid values are 'project', 'task', 'tasklist', "+
					"'milestone', 'message', 'timelog', 'notebook', 'file', 'company' and 'link'. "),
			),
			mcp.WithArray("project-ids",
				mcp.Description("A list of project IDs to filter tags by projects"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtag.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalParam(&multiple.Request.Filters.ItemType, "item-type"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.ProjectIDs, "project-ids"),
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
		mcp.NewTool("retrieve-tag",
			mcp.WithDescription("Retrieve a specific tag in a customer site of Teamwork.com. "+
				"Tags are a way to mark items so that you can use a filter to see just those items."),
			mcp.WithNumber("tag-id",
				mcp.Required(),
				mcp.Description("The ID of the tag."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twtag.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "tag-id"),
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
		mcp.NewTool("create-tag",
			mcp.WithDescription("Create a new tag in a customer site of Teamwork.com. "+
				"Tags are a way to mark items so that you can use a filter to see just those items."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the tag. It must have less than 50 characters."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to associate the tag with. This is for when you project-scoped tag."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tag twtag.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredParam(&tag.Name, "name"),
				twmcp.OptionalNumericPointerParam(&tag.ProjectID, "project-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if len(tag.Name) > 50 {
				return nil, fmt.Errorf("tag name must have less than 50 characters")
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tag); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Tag created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-tag",
			mcp.WithDescription("Update a tag in a customer site of Teamwork.com. "+
				"Tags are a way to mark items so that you can use a filter to see just those items."),
			mcp.WithNumber("tag-id",
				mcp.Required(),
				mcp.Description("The ID of the tag to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the tag. It must have less than 50 characters."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to associate the tag with. This is for when you project-scoped tag."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tag twtag.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&tag.ID, "tag-id"),
				twmcp.OptionalPointerParam(&tag.Name, "name"),
				twmcp.OptionalNumericPointerParam(&tag.ProjectID, "project-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if tag.Name != nil && len(*tag.Name) > 50 {
				return nil, fmt.Errorf("tag name must have less than 50 characters")
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tag); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Tag updated successfully"), nil
		},
	)
}
