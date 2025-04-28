package jobrole

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the job role resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage job roles in a
// customer site of Teamwork.com.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
