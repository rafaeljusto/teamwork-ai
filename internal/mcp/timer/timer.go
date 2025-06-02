package timer

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the timer resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage timers in a customer
// site of Teamwork.com. Timelogs are records of the amount that users spent
// working on a task or project.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
