package comment

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the comment resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage comments in a customer
// site of Teamwork.com. Comments are messages or notes that can be added to
// various objects in Teamwork, such as tasks, files, milestones, and notebooks.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
