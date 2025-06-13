package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twactivity "github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
)

var resourceList = mcp.NewResource("twapi://activities", "activities",
	mcp.WithResourceDescription("Activities are logs of actions taken in Teamwork.com, such as "+
		"creating, editing, or deleting items. They provide a history of changes made to projects, tasks, "+
		"and other objects."),
	mcp.WithMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twactivity.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, activity := range multiple.Response.Activities {
				encoded, err := json.Marshal(activity)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://activities/%d", activity.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)
}
