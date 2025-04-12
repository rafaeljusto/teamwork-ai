package project

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the project resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage projects in a customer
// site of Teamwork.com. A project is a central hub to manage all of the
// components relating to what your team is working on. It also provides a list
// of all projects and allows for the retrieval of a specific project by its ID.
// It also provides tools to retrieve multiple projects, a specific project, and
// to create a new project.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
