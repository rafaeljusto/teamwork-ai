package tag

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the tag resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage tags in a customer
// site of Teamwork.com. Tags are a way to mark items so that you can use a
// filter to see just those items. Tags can be added to projects, tasks,
// milestones, messages, time logs, notebooks, files, and links.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
