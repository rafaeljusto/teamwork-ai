package skill

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the skill resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage skills in a customer
// site of Teamwork.com. A skill is a knowledge or ability that can be assigned
// to users. It also provides a list of all skills and allows for the retrieval
// of a specific skill by its ID. Additionally, it provides tools to retrieve
// multiple skills, a specific skill, create a new skill, and update an existing
// skill.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
