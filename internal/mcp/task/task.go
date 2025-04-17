package task

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

// Register registers the task resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage tasks in a customer
// site of Teamwork.com. A task is an activity that needs to be carried out by
// one or multiple project members. It also provides a list of all tasks and
// allows for the retrieval of a specific task by its ID. Additionally, it
// provides tools to retrieve multiple tasks, a specific task, create a new
// task, and update an existing task. It also allows for the retrieval of tasks
// from a specific project or tasklist. Tasks can be assigned to users,
// companies, or teams, and can have a priority level (low, medium, high).
func Register(mcpServer *server.MCPServer, configResources *config.Resources) {
	resources(mcpServer, configResources)
	tools(mcpServer, configResources)
}
