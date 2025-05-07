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

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twskill.Multiple
			multiple.Request.Filters.Include = []string{"users"}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, skill := range multiple.Response.Skills {
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
			if err := configResources.TeamworkEngine.Do(ctx, &skill); err != nil {
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
}
