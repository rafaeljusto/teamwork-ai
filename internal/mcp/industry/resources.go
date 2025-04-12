package industry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twindustry "github.com/rafaeljusto/teamwork-ai/internal/teamwork/industry"
)

var resourceList = mcp.NewResource("twapi://industries", "industries",
	mcp.WithResourceDescription("Industries are categories that companies can belong to in Teamwork.com. "+
		"Each industry has an ID and a name."),
	mcp.WithMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twindustry.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, industry := range multiple.Response.Industries {
				encoded, err := json.Marshal(industry)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://industries/%d", industry.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)
}
