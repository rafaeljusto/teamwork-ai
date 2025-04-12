package milestone

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmilestone "github.com/rafaeljusto/teamwork-ai/internal/teamwork/milestone"
)

var resourceList = mcp.NewResource("twapi://milestones", "milestones",
	mcp.WithResourceDescription("Milestones are a target date representing a point of progress, or goal within a "+
		"project, that you can use task lists to track progress towards."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://milestones/{id}", "milestone",
	mcp.WithTemplateDescription("Milestone is a target date representing a point of progress, or goal within a "+
		"project, that you can use task lists to track progress towards."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twmilestone.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, milestone := range multiple.Response.Milestones {
				encoded, err := json.Marshal(milestone)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://milestones/%d", milestone.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reMilestoneID := regexp.MustCompile(`twapi://milestones/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reMilestoneID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid milestone ID")
			}
			milestoneID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid milestone ID")
			}

			var milestone twmilestone.Single
			milestone.ID = milestoneID
			if err := configResources.TeamworkEngine.Do(ctx, &milestone); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(milestone)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://milestones/%d", milestone.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
