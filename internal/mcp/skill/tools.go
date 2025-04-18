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
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skills twskill.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &skills); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(skills)
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
				mcp.Description("List of user IDs assigned to the skill. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Creation

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
				mcp.Required(),
				mcp.Description("The name of the skill."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("List of user IDs assigned to the skill. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skillUpdate twskill.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&skillUpdate.ID, "skill-id"),
				twmcp.RequiredParam(&skillUpdate.Skill.Name, "name"),
				twmcp.OptionalNumericListParam(&skillUpdate.Skill.UserIDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &skillUpdate); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Skill created successfully"), nil
		},
	)
}
