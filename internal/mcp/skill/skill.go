package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twskill "github.com/rafaeljusto/teamwork-ai/internal/teamwork/skill"
)

var resourceList = mcp.NewResource("twapi://skills", "skills",
	mcp.WithResourceDescription("Skills are knowledge or abilities that can be assigned to users."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://skills/{id}", "skill",
	mcp.WithTemplateDescription("Skill is a knowledge or ability that can be assigned to users."),
	mcp.WithTemplateMIMEType("application/json"),
)

// Register registers the skill resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage skills in a customer
// site of Teamwork.com. A skill is a knowledge or ability that can be assigned
// to users. It also provides a list of all skills and allows for the retrieval
// of a specific skill by its ID. Additionally, it provides tools to retrieve
// multiple skills, a specific skill, create a new skill, and update an existing
// skill.
func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var skills twskill.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &skills); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, skill := range skills {
				encoded, err := json.Marshal(skill)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://skills/%d", skill.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reSkillID := regexp.MustCompile(`twapi://skills/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reSkillID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid skill ID")
			}
			skillID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid skill ID")
			}

			var skill twskill.Single
			skill.ID = skillID
			if err := resources.TeamworkEngine.Do(ctx, &skill); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(skill)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://skills/%d", skill.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-skills",
			mcp.WithDescription("Retrieve multiple skills in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skills twskill.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &skills); err != nil {
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
			mcp.WithNumber("skillId",
				mcp.Required(),
				mcp.Description("The ID of the skill."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Single

			id, ok := request.Params.Arguments["skillId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid skillId")
			} else if id == 0 {
				return nil, fmt.Errorf("skillId is required")
			}
			skill.ID = int64(id)

			if err := resources.TeamworkEngine.Do(ctx, &skill); err != nil {
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
			mcp.WithArray("userIds",
				mcp.Description("List of user IDs assigned to the skill. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skill twskill.Creation
			var ok bool

			skill.Name, ok = request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if skill.Name == "" {
				return nil, fmt.Errorf("name is required")
			}

			err := twmcp.OptionalNumericListParam(request.Params.Arguments, &skill.UserIDs, "userIds")
			if err != nil {
				return nil, fmt.Errorf("invalid userIds: %w", err)
			}

			if err := resources.TeamworkEngine.Do(ctx, &skill); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Skill created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-skill",
			mcp.WithDescription("Update an existing skill in a customer site of Teamwork.com. "+
				"Skill is a knowledge or ability that can be assigned to users."),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("The ID of the skill to update."),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the skill."),
			),
			mcp.WithArray("userIds",
				mcp.Description("List of user IDs assigned to the skill. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var skillUpdate twskill.Update
			var ok bool

			id, ok := request.Params.Arguments["id"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid id")
			} else if id == 0 {
				return nil, fmt.Errorf("id is required")
			}
			skillUpdate.ID = int64(id)

			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if name != "" {
				skillUpdate.Skill.Name = &name
			}

			err := twmcp.OptionalNumericListParam(request.Params.Arguments, &skillUpdate.Skill.UserIDs, "userIds")
			if err != nil {
				return nil, fmt.Errorf("invalid userIds: %w", err)
			}

			if err := resources.TeamworkEngine.Do(ctx, &skillUpdate); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Skill created successfully"), nil
		},
	)
}
