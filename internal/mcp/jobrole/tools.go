package jobrole

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twjobrole "github.com/rafaeljusto/teamwork-ai/internal/teamwork/jobrole"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-jobroles",
			mcp.WithDescription("Retrieve multiple job roles in a customer site of Teamwork.com. "+
				"Job role is a role that can be assigned to users."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter job roles by name. "+
					"Each word from the search term is used to match against the job role name or any assigned user name. "+
					"The job role will be selected if each word of the term matches the job role name or assigned user name, "+
					"not requiring that the word matches are in the same field."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twjobrole.Multiple
			multiple.Request.Filters.Include = []string{"users"}

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
			encoded, err := json.Marshal(multiple.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-jobrole",
			mcp.WithDescription("Retrieve a specific job role in a customer site of Teamwork.com. "+
				"Job role is a role that can be assigned to users."),
			mcp.WithNumber("jobrole-id",
				mcp.Required(),
				mcp.Description("The ID of the job role."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twjobrole.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&single.ID, "jobrole-id"),
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
		mcp.NewTool("create-jobrole",
			mcp.WithDescription("Create a new job role in a customer site of Teamwork.com. "+
				"Job role is a role that can be assigned to users."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the job role."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var jobrole twjobrole.Create

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&jobrole.Name, "name"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &jobrole); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Job role created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-jobrole",
			mcp.WithDescription("Update a job role in a customer site of Teamwork.com. "+
				"Job role is a role that can be assigned to users."),
			mcp.WithNumber("jobrole-id",
				mcp.Required(),
				mcp.Description("The ID of the jobrole to update."),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the jobrole."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var jobrole twjobrole.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&jobrole.ID, "jobrole-id"),
				twmcp.RequiredParam(&jobrole.Name, "name"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &jobrole); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Job role updated successfully"), nil
		},
	)
}
