package team

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twteam "github.com/rafaeljusto/teamwork-ai/internal/twapi/team"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-teams",
			mcp.WithDescription("Retrieve multiple teams in a customer site of Teamwork.com. "+
				"Teams replicate your organization's structure and group people on your site based on their "+
				"position or contribution"),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter teams by name or handle."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twteam.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
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
		mcp.NewTool("retrieve-team",
			mcp.WithDescription("Retrieve a specific team in a customer site of Teamwork.com. "+
				"Teams replicate your organization's structure and group people on your site based on their "+
				"position or contribution"),
			mcp.WithNumber("team-id",
				mcp.Required(),
				mcp.Description("The ID of the team."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twteam.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "team-id"),
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
		mcp.NewTool("create-team",
			mcp.WithDescription("Create a new team in a customer site of Teamwork.com. "+
				"Teams replicate your organization's structure and group people on your site based on their "+
				"position or contribution"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the team."),
			),
			mcp.WithString("handle",
				mcp.Description("The handle of the team. It is a unique identifier for the team. It must not have spaces "+
					"or special characters."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the team."),
			),
			mcp.WithNumber("parent-team-id",
				mcp.Description("The ID of the parent team. This is used to create a hierarchy of teams."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company. This is used to create a team scoped for a specific company."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project. This is used to create a team scoped for a specific project."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to add to the team."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var team twteam.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredParam(&team.Name, "name"),
				twmcp.OptionalPointerParam(&team.Handle, "handle"),
				twmcp.OptionalPointerParam(&team.Description, "description"),
				twmcp.OptionalNumericPointerParam(&team.ParentTeamID, "parent-team-id"),
				twmcp.OptionalNumericPointerParam(&team.CompanyID, "company-id"),
				twmcp.OptionalNumericPointerParam(&team.ProjectID, "project-id"),
				twmcp.OptionalCustomNumericListParam(&team.UserIDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &team); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Team created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-team",
			mcp.WithDescription("Update a team in a customer site of Teamwork.com. "+
				"Teams replicate your organization's structure and group people on your site based on their "+
				"position or contribution"),
			mcp.WithNumber("team-id",
				mcp.Required(),
				mcp.Description("The ID of the team to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the team."),
			),
			mcp.WithString("handle",
				mcp.Description("The handle of the team. It is a unique identifier for the team. It must not have spaces "+
					"or special characters."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the team."),
			),
			mcp.WithNumber("parent-team-id",
				mcp.Description("The ID of the parent team. This is used to create a hierarchy of teams."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company. This is used to create a team scoped for a specific company."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project. This is used to create a team scoped for a specific project."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to add to the team."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var team twteam.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&team.ID, "team-id"),
				twmcp.OptionalPointerParam(&team.Name, "name"),
				twmcp.OptionalPointerParam(&team.Handle, "handle"),
				twmcp.OptionalPointerParam(&team.Description, "description"),
				twmcp.OptionalNumericPointerParam(&team.CompanyID, "company-id"),
				twmcp.OptionalNumericPointerParam(&team.ProjectID, "project-id"),
				twmcp.OptionalCustomNumericListParam(&team.UserIDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &team); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Team updated successfully"), nil
		},
	)
}
