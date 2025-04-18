package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	twtask "github.com/rafaeljusto/teamwork-ai/internal/teamwork/task"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerToolsRetrieve(mcpServer, configResources)
	registerToolsCreate(mcpServer, configResources)
	registerToolsUpdate(mcpServer, configResources)
}

func registerToolsRetrieve(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasks",
			mcp.WithDescription("Retrieve multiple tasks in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasks twtask.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&tasks.ProjectID, "projectId"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid project ID: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&tasks.TasklistID, "tasklistId"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid tasklist ID: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &tasks); err != nil {
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

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&task.ID, "taskId"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &task); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(task)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)
}

func registerToolsCreate(mcpServer *server.MCPServer, configResources *config.Resources) {
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

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&task.Name, "name"),
				twmcp.RequiredNumericParam(&task.TasklistID, "tasklistId"),
				twmcp.OptionalParam(&task.Description, "description"),
				twmcp.OptionalPointerParam(&task.Priority, "priority",
					twmcp.RestrictValues("low", "medium", "high"),
				),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			assignees, ok := request.Params.Arguments["assignees"]
			if ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assignees != nil {
					task.Assignees = new(teamwork.UserGroups)

					err = twmcp.ParamGroup(assigneesMap,
						twmcp.OptionalNumericListParam(&task.Assignees.UserIDs, "userIds"),
						twmcp.OptionalNumericListParam(&task.Assignees.CompanyIDs, "companyIds"),
						twmcp.OptionalNumericListParam(&task.Assignees.TeamIDs, "teamIds"),
					)
					if err != nil {
						return nil, fmt.Errorf("invalid assignees: %w", err)
					}
				}
			}

			if err := configResources.TeamworkEngine.Do(ctx, &task); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)
}

func registerToolsUpdate(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("update-task",
			mcp.WithDescription("Update an existing task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("taskId",
				mcp.Required(),
				mcp.Description("The ID of the task to update."),
			),
			mcp.WithString("name",
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

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&taskUpdate.ID, "taskId"),
				twmcp.OptionalParam(&taskUpdate.Task.Name, "name"),
				twmcp.OptionalPointerParam(&taskUpdate.Task.Description, "description"),
				twmcp.OptionalPointerParam(&taskUpdate.Task.Priority, "priority",
					twmcp.RestrictValues("low", "medium", "high"),
				),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			assignees, ok := request.Params.Arguments["assignees"]
			if ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assignees != nil {
					taskUpdate.Task.Assignees = new(teamwork.UserGroups)

					err = twmcp.ParamGroup(assigneesMap,
						twmcp.OptionalNumericListParam(&taskUpdate.Task.Assignees.UserIDs, "userIds"),
						twmcp.OptionalNumericListParam(&taskUpdate.Task.Assignees.CompanyIDs, "companyIds"),
						twmcp.OptionalNumericListParam(&taskUpdate.Task.Assignees.TeamIDs, "teamIds"),
					)
					if err != nil {
						return nil, fmt.Errorf("invalid assignees: %w", err)
					}
				}
			}

			if err := configResources.TeamworkEngine.Do(ctx, &taskUpdate); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Task created successfully"), nil
		},
	)
}
