package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtimer "github.com/rafaeljusto/teamwork-ai/internal/twapi/timer"
)

var resourceList = mcp.NewResource("twapi://timers", "timers",
	mcp.WithResourceDescription("Timers are used to track ongoing work that will generate timelogs."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://timers/{id}", "timer",
	mcp.WithTemplateDescription("Timer is used to track ongoing work that will generate timelogs."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twtimer.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, timer := range multiple.Response.Timers {
				encoded, err := json.Marshal(timer)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://timers/%d", timer.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTimerID := regexp.MustCompile(`twapi://timers/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTimerID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid timer ID")
			}
			timerID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid timer ID")
			}

			var timer twtimer.Single
			timer.ID = timerID
			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(timer)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://timers/%d", timer.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
