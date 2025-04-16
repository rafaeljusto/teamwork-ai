package task

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twtask "github.com/rafaeljusto/teamwork-ai/internal/teamwork/task"
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
			for _, task := range tasks.Tasks {
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
		mcp.NewTool("retrieve-tasks",
			mcp.WithDescription("Retrieve multiple tasks in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasks twtask.MultipleTasks
			if err := resources.TeamworkEngine.Do(&tasks); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasks)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project-tasks",
			mcp.WithDescription("Retrieve multiple tasks from a specific project in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("projectId",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasks."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasks twtask.MultipleTasks

			projectID, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if projectID == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			tasks.ProjectID = int64(projectID)

			if err := resources.TeamworkEngine.Do(&tasks); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasks)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasklist-tasks",
			mcp.WithDescription("Retrieve multiple tasks from a specific tasklist in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("tasklistId",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasks."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasks twtask.MultipleTasks

			tasklistID, ok := request.Params.Arguments["tasklistId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid tasklistId")
			} else if tasklistID == 0 {
				return nil, fmt.Errorf("tasklistId is required")
			}
			tasks.TasklistID = int64(tasklistID)

			if err := resources.TeamworkEngine.Do(&tasks); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasks)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-task",
			mcp.WithDescription("Retrieve a specific task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("taskId",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var task twtask.SingleTask

			id, ok := request.Params.Arguments["taskId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid taskId")
			} else if id == 0 {
				return nil, fmt.Errorf("taskId is required")
			}
			task.ID = int64(id)

			if err := resources.TeamworkEngine.Do(&task); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(task)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-task",
			mcp.WithDescription("Create a new task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
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
			mcp.WithObject("assignees",
				mcp.Description("The assignees of the task. This is a JSON object with user IDs, company IDs, and team IDs."),
				mcp.Properties(map[string]any{
					"userIds": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the task.",
					},
					"companyIds": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the task.",
					},
					"teamIds": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the task.",
					},
				}),
			),
			mcp.WithString("priority",
				mcp.Description("The priority of the task. Possible values are: low, medium, high."),
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

			assignees, ok := request.Params.Arguments["assignees"].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid assignees")
			} else if assignees != nil {
				if userIDs, ok := assignees["userIds"].([]any); ok {
					for _, userID := range userIDs {
						if id, ok := userID.(float64); ok {
							task.Assignees.UserIDs = append(task.Assignees.UserIDs, int64(id))
						}
					}
				}
				if companyIDs, ok := assignees["companyIds"].([]any); ok {
					for _, companyID := range companyIDs {
						if id, ok := companyID.(float64); ok {
							task.Assignees.CompanyIDs = append(task.Assignees.CompanyIDs, int64(id))
						}
					}
				}
				if teamIDs, ok := assignees["teamIds"].([]any); ok {
					for _, teamID := range teamIDs {
						if id, ok := teamID.(float64); ok {
							task.Assignees.TeamIDs = append(task.Assignees.TeamIDs, int64(id))
						}
					}
				}
			}

			priority, ok := request.Params.Arguments["priority"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid priority")
			} else if priority != "" {
				switch priority {
				case "low", "medium", "high":
					task.Priority = &priority
				default:
					return nil, fmt.Errorf("invalid priority: %s", priority)
				}
			}

			if err := resources.TeamworkEngine.Do(&task); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-task",
			mcp.WithDescription("Update an existing task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("The ID of the task to update."),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the task."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the task."),
			),
			mcp.WithObject("assignees",
				mcp.Description("The assignees of the task. This is a JSON object with user IDs, company IDs, and team IDs."),
				mcp.Properties(map[string]any{
					"userIds": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the task.",
					},
					"companyIds": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the task.",
					},
					"teamIds": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the task.",
					},
				}),
			),
			mcp.WithString("priority",
				mcp.Description("The priority of the task. Possible values are: low, medium, high."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var task twtask.TaskUpdate
			var ok bool

			id, ok := request.Params.Arguments["id"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid id")
			} else if id == 0 {
				return nil, fmt.Errorf("id is required")
			}
			task.ID = int64(id)

			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if name != "" {
				task.Name = &name
			}

			description, ok := request.Params.Arguments["description"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid description")
			} else if description != "" {
				task.Description = &description
			}

			assignees, ok := request.Params.Arguments["assignees"].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid assignees")
			} else if assignees != nil {
				if userIDs, ok := assignees["userIds"].([]any); ok {
					for _, userID := range userIDs {
						if id, ok := userID.(float64); ok {
							task.Assignees.UserIDs = append(task.Assignees.UserIDs, int64(id))
						}
					}
				}
				if companyIDs, ok := assignees["companyIds"].([]any); ok {
					for _, companyID := range companyIDs {
						if id, ok := companyID.(float64); ok {
							task.Assignees.CompanyIDs = append(task.Assignees.CompanyIDs, int64(id))
						}
					}
				}
				if teamIDs, ok := assignees["teamIds"].([]any); ok {
					for _, teamID := range teamIDs {
						if id, ok := teamID.(float64); ok {
							task.Assignees.TeamIDs = append(task.Assignees.TeamIDs, int64(id))
						}
					}
				}
			}

			priority, ok := request.Params.Arguments["priority"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid priority")
			} else if priority != "" {
				switch priority {
				case "low", "medium", "high":
					task.Priority = &priority
				default:
					return nil, fmt.Errorf("invalid priority: %s", priority)
				}
			}

			if err := resources.TeamworkEngine.Do(&task); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)
}
