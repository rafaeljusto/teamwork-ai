package team

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the team resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage teams in a customer
// site of Teamwork.com. Team replicates your organization's structure and group
// people on your site based on their position or contribution.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
