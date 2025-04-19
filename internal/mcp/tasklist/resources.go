package tasklist

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtasklist "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasklist"
)

var resourceList = mcp.NewResource("twapi://tasklists", "tasklists",
	mcp.WithResourceDescription("Tasklists group tasks together in a project for better organization."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://tasklists/{id}", "task",
	mcp.WithTemplateDescription("Tasklist group tasks together in a project for better organization."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twtasklist.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, tasklist := range multiple.Response.Tasklists {
				encoded, err := json.Marshal(tasklist)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasklists/%d", tasklist.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTasklistID := regexp.MustCompile(`twapi://tasklists/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTasklistID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid tasklist ID")
			}
			tasklistID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid tasklist ID")
			}

			var tasklist twtasklist.Single
			tasklist.ID = tasklistID
			if err := configResources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(tasklist)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasklists/%d", tasklist.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
