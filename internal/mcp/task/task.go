package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtask "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasks"
)

var resourceList = mcp.NewResource("twapi://tasks", "tasks",
	mcp.WithResourceDescription("Tasks are activities that need to be carried out by one or multiple project members."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://tasks/{id}", "task",
	mcp.WithTemplateDescription("Task is an activity that need to be carried out by one or multiple project members."),
	mcp.WithTemplateMIMEType("application/json"),
)

func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var tasks twtask.MultipleTasks
			if err := resources.TeamworkEngine.Do(&tasks); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, task := range tasks {
				encoded, err := json.Marshal(task)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasks/%d", task.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTaskID := regexp.MustCompile(`twapi://tasks/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTaskID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid task ID")
			}
			taskID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid task ID")
			}

			var task twtask.SingleTask
			task.ID = taskID
			if err := resources.TeamworkEngine.Do(&task); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(task)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasks/%d", task.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("task",
			mcp.WithDescription("Task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the task."),
			),
			mcp.WithNumber("tasklistId",
				mcp.Required(),
				mcp.Description("The ID of the tasklist."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var task twtask.TaskCreation
			var ok bool

			task.Name, ok = request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if task.Name == "" {
				return nil, fmt.Errorf("name is required")
			}

			tasklistID, ok := request.Params.Arguments["tasklistId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid tasklistId")
			} else if tasklistID == 0 {
				return nil, fmt.Errorf("tasklistId is required")
			}
			task.TasklistID = int64(tasklistID)

			task.Description, ok = request.Params.Arguments["description"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid description")
			}

			if err := resources.TeamworkEngine.Do(&task); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)
}
