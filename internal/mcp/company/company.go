package company

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the company resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage companies in a
// customer site of Teamwork.com. Companies, also known as clients, are
// organizations that the customer offers services to. It also provides a list
// of all companies and allows for the retrieval of a specific company by its
// ID.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
