package activity

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the activity resources and tools with the MCP server. It
// provides functionality to retrieve activities in a customer site of
// Teamwork.com. Activities are logs of actions taken in Teamwork, such as
// creating, editing, or deleting items. They provide a history of changes made
// to projects, tasks, and other objects.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
