package milestone

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the milestone resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage milestones in a
// customer site of Teamwork.com. Milestones are target dates representing a
// point of progress or goal within a project, which can be tracked using task
// lists. It also provides a list of all milestones and allows for the retrieval
// of a specific milestone by its ID, as well as the creation and updating of
// milestones.
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerResources(mcpServer, configResources)
	registerTools(mcpServer, configResources)
}
