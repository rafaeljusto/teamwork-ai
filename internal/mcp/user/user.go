package user

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the user resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage users in a customer
// site of Teamwork.com. A user, also known as a person, is an individual who
// can be assigned to tasks. It also provides a list of all users and allows for
// the retrieval of a specific user by their ID. Additionally, it provides tools
// to retrieve multiple users, a specific user, create a new user, and update an
// existing user.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
