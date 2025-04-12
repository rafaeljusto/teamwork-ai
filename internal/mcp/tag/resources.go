package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtag "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tag"
)

var resourceList = mcp.NewResource("twapi://tags", "tags",
	mcp.WithResourceDescription("Tags are a way to mark items so that you can use a filter to see just those "+
		"items. Tags can be added to projects, tasks, milestones, messages, time logs, "+
		"notebooks, files and links."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://tags/{id}", "tag",
	mcp.WithTemplateDescription("Tag is a way to mark items so that you can use a filter to see just those "+
		"items. Tags can be added to projects, tasks, milestones, messages, time logs, "+
		"notebooks, files and links."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twtag.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, tag := range multiple.Response.Tags {
				encoded, err := json.Marshal(tag)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tags/%d", tag.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reCompanyID := regexp.MustCompile(`twapi://tags/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reCompanyID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid tag ID")
			}
			tagID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid tag ID")
			}

			var tag twtag.Single
			tag.ID = tagID
			if err := configResources.TeamworkEngine.Do(ctx, &tag); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(tag)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tags/%d", tag.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
