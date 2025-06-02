package timelog

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtimelog "github.com/rafaeljusto/teamwork-ai/internal/twapi/timelog"
)

var resourceList = mcp.NewResource("twapi://timelogs", "timelogs",
	mcp.WithResourceDescription("Timelogs are records of the amount that users spent working on a task or project."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://timelogs/{id}", "timelog",
	mcp.WithTemplateDescription("Timelog is record of the amount a user spent working on a task or project."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twtimelog.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, timelog := range multiple.Response.Timelogs {
				encoded, err := json.Marshal(timelog)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://timelogs/%d", timelog.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTimelogID := regexp.MustCompile(`twapi://timelogs/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTimelogID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid timelog ID")
			}
			timelogID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid timelog ID")
			}

			var timelog twtimelog.Single
			timelog.ID = timelogID
			if err := configResources.TeamworkEngine.Do(ctx, &timelog); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(timelog)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://timelogs/%d", timelog.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
