package industry

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twindustry "github.com/rafaeljusto/teamwork-ai/internal/twapi/industry"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveIndustries.String(),
			mcp.WithDescription("Retrieve multiple industries in a customer site of Teamwork.com."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twindustry.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)
}
