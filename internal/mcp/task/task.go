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
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
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

// Register registers the task resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage tasks in a customer
// site of Teamwork.com. A task is an activity that needs to be carried out by
// one or multiple project members. It also provides a list of all tasks and
// allows for the retrieval of a specific task by its ID. Additionally, it
// provides tools to retrieve multiple tasks, a specific task, create a new
// task, and update an existing task. It also allows for the retrieval of tasks
// from a specific project or tasklist. Tasks can be assigned to users,
// companies, or teams, and can have a priority level (low, medium, high).
func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var tasks twtask.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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

			var task twtask.Single
			task.ID = taskID
			if err := resources.TeamworkEngine.Do(ctx, &task); err != nil {
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
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasks twtask.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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
			var tasks twtask.Multiple

			projectID, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if projectID == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			tasks.ProjectID = int64(projectID)

			if err := resources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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
			var tasks twtask.Multiple

			tasklistID, ok := request.Params.Arguments["tasklistId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid tasklistId")
			} else if tasklistID == 0 {
				return nil, fmt.Errorf("tasklistId is required")
			}
			tasks.TasklistID = int64(tasklistID)

			if err := resources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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
			var task twtask.Single

			id, ok := request.Params.Arguments["taskId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid taskId")
			} else if id == 0 {
				return nil, fmt.Errorf("taskId is required")
			}
			task.ID = int64(id)

			if err := resources.TeamworkEngine.Do(ctx, &task); err != nil {
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
			var task twtask.Creation
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

			description, ok, err := twmcp.OptionalParam[string](request.Params.Arguments, "description")
			if err != nil {
				return nil, fmt.Errorf("invalid description: %w", err)
			} else if ok {
				task.Description = description
			}

			assignees, ok := request.Params.Arguments["assignees"]
			if ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assignees != nil {
					task.Assignees = new(teamwork.UserGroups)

					userIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "userIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						task.Assignees.UserIDs = userIDs
					}
					companyIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "companyIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						task.Assignees.UserIDs = companyIDs
					}
					teamIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "teamIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						task.Assignees.UserIDs = teamIDs
					}
				}
			}

			priority, ok, err := twmcp.OptionalParam[string](request.Params.Arguments, "priority")
			if err != nil {
				return nil, fmt.Errorf("invalid priority: %w", err)
			} else if ok {
				if priority != "" {
					switch priority {
					case "low", "medium", "high":
						task.Priority = &priority
					default:
						return nil, fmt.Errorf("invalid priority: %s", priority)
					}
				}
			}

			if err := resources.TeamworkEngine.Do(ctx, &task); err != nil {
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
			var taskUpdate twtask.Update
			var ok bool

			id, ok := request.Params.Arguments["id"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid id")
			} else if id == 0 {
				return nil, fmt.Errorf("id is required")
			}
			taskUpdate.ID = int64(id)

			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if name != "" {
				taskUpdate.Task.Name = &name
			}

			description, ok, err := twmcp.OptionalParam[string](request.Params.Arguments, "description")
			if err != nil {
				return nil, fmt.Errorf("invalid description: %w", err)
			} else if ok {
				taskUpdate.Task.Description = &description
			}

			assignees, ok := request.Params.Arguments["assignees"]
			if ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assignees != nil {
					taskUpdate.Task.Assignees = new(teamwork.UserGroups)

					userIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "userIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						taskUpdate.Task.Assignees.UserIDs = userIDs
					}
					companyIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "companyIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						taskUpdate.Task.Assignees.CompanyIDs = companyIDs
					}
					teamIDs, ok, err := twmcp.OptionalNumericListParam[int64](assigneesMap, "teamIds")
					if err != nil {
						return nil, fmt.Errorf("invalid userIds: %w", err)
					} else if ok {
						taskUpdate.Task.Assignees.TeamIDs = teamIDs
					}
				}
			}

			priority, ok, err := twmcp.OptionalParam[string](request.Params.Arguments, "priority")
			if err != nil {
				return nil, fmt.Errorf("invalid priority: %w", err)
			} else if ok {
				if priority != "" {
					switch priority {
					case "low", "medium", "high":
						taskUpdate.Task.Priority = &priority
					default:
						return nil, fmt.Errorf("invalid priority: %s", priority)
					}
				}
			}

			if err := resources.TeamworkEngine.Do(ctx, &taskUpdate); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)
}
