package industry

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the industry resources and tools with the MCP server. It
// provides functionality to retrieve industries in a customer site of
// Teamwork.com. Industries are categories that companies can belong to, each
// with an ID and a name.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
