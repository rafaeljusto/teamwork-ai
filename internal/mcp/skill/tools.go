package skill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twskill "github.com/rafaeljusto/teamwork-ai/internal/teamwork/skill"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-skills",
			mcp.WithDescription("Retrieve multiple skills in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter skills by name, or by the first or last names of "+
					"the user associated with the skill. The skill will be selected if each word of the term matches "+
					"the skill name or the user first or last name, not requiring that the word matches are in the same field."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twskill.Multiple
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
			encoded, err := json.Marshal(multiple)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-skill",
			mcp.WithDescription("Retrieve a specific skill in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
			mcp.WithNumber("skill-id",
				mcp.Required(),
				mcp.Description("The ID of the skill."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&skill.ID, "skill-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &skill); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(skill)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-skill",
			mcp.WithDescription("Create a new skill in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the skill."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs assigned to the skill."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Create

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&skill.Name, "name"),
				twmcp.OptionalNumericListParam(&skill.UserIDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &skill); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Skill created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-skill",
			mcp.WithDescription("Update an existing skill in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
			mcp.WithNumber("skill-id",
				mcp.Required(),
				mcp.Description("The ID of the skill to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the skill."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs assigned to the skill."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&skill.ID, "skill-id"),
				twmcp.OptionalPointerParam(&skill.Name, "name"),
				twmcp.OptionalNumericListParam(&skill.UserIDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &skill); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Skill created successfully"), nil
		},
	)
}
