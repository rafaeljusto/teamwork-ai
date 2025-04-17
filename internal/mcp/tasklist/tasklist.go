package tasklist

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the tasklist resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage tasklists in a
// customer site of Teamwork.com. A tasklist groups tasks together in a project
// for better organization. It also provides a list of all tasklists and allows
// for the retrieval of a specific tasklist by its ID. Additionally, it provides
// tools to retrieve multiple tasklists, a specific tasklist, create a new
// tasklist, and update an existing tasklist.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
