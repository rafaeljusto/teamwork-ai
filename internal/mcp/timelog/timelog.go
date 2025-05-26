package timelog

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the timelog resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage timelogs in a customer
// site of Teamwork.com. Timelogs are records of the amount that users spent
// working on a task or project.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
