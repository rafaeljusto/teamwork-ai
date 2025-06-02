package jobrole

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twjobrole "github.com/rafaeljusto/teamwork-ai/internal/twapi/jobrole"
)

var resourceList = mcp.NewResource("twapi://jobroles", "jobroles",
	mcp.WithResourceDescription("Job roles are roles that can be assigned to users."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://jobroles/{id}", "jobrole",
	mcp.WithTemplateDescription("Job role is a role that can be assigned to users."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twjobrole.Multiple
			multiple.Request.Filters.Include = []string{"users"}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, jobrole := range multiple.Response.JobRoles {
				encoded, err := json.Marshal(jobrole)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://jobroles/%d", jobrole.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reJobRoleID := regexp.MustCompile(`twapi://jobroles/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reJobRoleID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid jobrole ID")
			}
			jobroleID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid jobrole ID")
			}

			var jobrole twjobrole.Single
			jobrole.ID = jobroleID
			if err := configResources.TeamworkEngine.Do(ctx, &jobrole); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(jobrole)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://jobroles/%d", jobrole.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
