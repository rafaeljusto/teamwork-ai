package team

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	twteam "github.com/rafaeljusto/teamwork-ai/internal/twapi/team"
)

var resourceList = mcp.NewResource("twapi://teams", "teams",
	mcp.WithResourceDescription("Teams replicate your organization's structure and group people on your site based "+
		"on their position or contribution"),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://teams/{id}", "team",
	mcp.WithTemplateDescription("Team replicates your organization's structure and group people on your site based "+
		"on their position or contribution"),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twteam.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, team := range multiple.Response.Teams {
				encoded, err := json.Marshal(team)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://teams/%d", team.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTeamID := regexp.MustCompile(`twapi://teams/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTeamID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid team ID")
			}
			teamID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid team ID")
			}

			var team twteam.Single
			team.ID = twapi.LegacyNumber(teamID)
			if err := configResources.TeamworkEngine.Do(ctx, &team); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(team)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://teams/%d", team.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
