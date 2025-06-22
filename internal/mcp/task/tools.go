package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	twtask "github.com/rafaeljusto/teamwork-ai/internal/twapi/task"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	registerToolsRetrieve(mcpServer, configResources)
	registerToolsCreate(mcpServer, configResources)
	registerToolsUpdate(mcpServer, configResources)
}

func registerToolsRetrieve(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveTasks.String(),
			mcp.WithDescription("Retrieve multiple tasks in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tasks by name, description or the related tasklist's name. "+
					"The task will be selected if each word of the term matches the task name, task description, or the "+
					"tasklist name, not requiring that the word matches are in the same field."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter tasks by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match tasks that have all the specified tags. "+
					"If false, the search will match tasks that have any of the specified tags. "+
					"Defaults to false."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtask.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.TagIDs, "tag-ids"),
				twmcp.OptionalPointerParam(&multiple.Request.Filters.MatchAllTags, "match-all-tags"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveProjectTasks.String(),
			mcp.WithDescription("Retrieve multiple tasks from a specific project in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasks."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tasks by name, description or the related tasklist's name. "+
					"The task will be selected if each word of the term matches the task name, task description, or the "+
					"tasklist name, not requiring that the word matches are in the same field."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter tasks by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match tasks that have all the specified tags. "+
					"If false, the search will match tasks that have any of the specified tags. "+
					"Defaults to false."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtask.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.ProjectID, "project-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.TagIDs, "tag-ids"),
				twmcp.OptionalPointerParam(&multiple.Request.Filters.MatchAllTags, "match-all-tags"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveTasklistTasks.String(),
			mcp.WithDescription("Retrieve multiple tasks from a specific tasklist in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("tasklist-id",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasks."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter tasks by name, description or the related tasklist's name. "+
					"The task will be selected if each word of the term matches the task name, task description, or the "+
					"tasklist name, not requiring that the word matches are in the same field."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter tasks by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match tasks that have all the specified tags. "+
					"If false, the search will match tasks that have any of the specified tags. "+
					"Defaults to false."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtask.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.TasklistID, "tasklist-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.TagIDs, "tag-ids"),
				twmcp.OptionalPointerParam(&multiple.Request.Filters.MatchAllTags, "match-all-tags"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveTask.String(),
			mcp.WithDescription("Retrieve a specific task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("task-id",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twtask.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "task-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &single); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(single)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)
}

func registerToolsCreate(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodCreateTask.String(),
			mcp.WithDescription("Create a new task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the task."),
			),
			mcp.WithNumber("tasklist-id",
				mcp.Required(),
				mcp.Description("The ID of the tasklist."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the task."),
			),
			mcp.WithString("priority",
				mcp.Description("The priority of the task. Possible values are: low, medium, high."),
			),
			mcp.WithNumber("progress",
				mcp.Description("The progress of the task, as a percentage (0-100). Only whole numbers are allowed."),
			),
			mcp.WithString("start-date",
				mcp.Description("The start date of the task in ISO 8601 format (YYYY-MM-DD)."),
			),
			mcp.WithString("due-date",
				mcp.Description("The due date of the task in ISO 8601 format (YYYY-MM-DD)."),
			),
			mcp.WithNumber("estimated-minutes",
				mcp.Description("The estimated time to complete the task in minutes."),
			),
			mcp.WithObject("assignees",
				mcp.Description("The assignees of the task. This is a JSON object with user IDs, company IDs, and team IDs."),
				mcp.Properties(map[string]any{
					"user-ids": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the task.",
					},
					"company-ids": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the task.",
					},
					"team-ids": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the task.",
					},
				}),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to assign to the task."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var task twtask.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredParam(&task.Name, "name"),
				twmcp.RequiredNumericParam(&task.TasklistID, "tasklist-id"),
				twmcp.OptionalPointerParam(&task.Description, "description"),
				twmcp.OptionalPointerParam(&task.Priority, "priority",
					twmcp.RestrictValues("low", "medium", "high"),
				),
				twmcp.OptionalNumericPointerParam(&task.Progress, "progress"),
				twmcp.OptionalDatePointerParam(&task.StartAt, "start-date"),
				twmcp.OptionalDatePointerParam(&task.DueAt, "due-date"),
				twmcp.OptionalNumericPointerParam(&task.EstimatedMinutes, "estimated-minutes"),
				twmcp.OptionalNumericListParam(&task.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if assignees, ok := request.GetArguments()["assignees"]; ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assigneesMap != nil {
					task.Assignees = new(twapi.UserGroups)

					err = twmcp.ParamGroup(assigneesMap,
						twmcp.OptionalNumericListParam(&task.Assignees.UserIDs, "user-ids"),
						twmcp.OptionalNumericListParam(&task.Assignees.CompanyIDs, "company-ids"),
						twmcp.OptionalNumericListParam(&task.Assignees.TeamIDs, "team-ids"),
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
		mcp.NewTool(twmcp.MethodUpdateTask.String(),
			mcp.WithDescription("Update an existing task in a customer site of Teamwork.com. "+
				"A task is an activity that need to be carried out by one or multiple project members."),
			mcp.WithNumber("task-id",
				mcp.Required(),
				mcp.Description("The ID of the task to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the task."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the task."),
			),
			mcp.WithString("priority",
				mcp.Description("The priority of the task. Possible values are: low, medium, high."),
			),
			mcp.WithNumber("progress",
				mcp.Description("The progress of the task, as a percentage (0-100). Only whole numbers are allowed."),
			),
			mcp.WithString("start-date",
				mcp.Description("The start date of the task in ISO 8601 format (YYYY-MM-DD)."),
			),
			mcp.WithString("due-date",
				mcp.Description("The due date of the task in ISO 8601 format (YYYY-MM-DD)."),
			),
			mcp.WithNumber("estimated-minutes",
				mcp.Description("The estimated time to complete the task in minutes."),
			),
			mcp.WithObject("assignees",
				mcp.Description("The assignees of the task. This is a JSON object with user IDs, company IDs, and team IDs."),
				mcp.Properties(map[string]any{
					"user-ids": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the task.",
					},
					"company-ids": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the task.",
					},
					"team-ids": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the task.",
					},
				}),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to assign to the task."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var task twtask.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&task.ID, "task-id"),
				twmcp.OptionalPointerParam(&task.Name, "name"),
				twmcp.OptionalPointerParam(&task.Description, "description"),
				twmcp.OptionalPointerParam(&task.Priority, "priority",
					twmcp.RestrictValues("low", "medium", "high"),
				),
				twmcp.OptionalNumericPointerParam(&task.Progress, "progress"),
				twmcp.OptionalDatePointerParam(&task.StartAt, "start-date"),
				twmcp.OptionalDatePointerParam(&task.DueAt, "due-date"),
				twmcp.OptionalNumericPointerParam(&task.EstimatedMinutes, "estimated-minutes"),
				twmcp.OptionalNumericListParam(&task.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if assignees, ok := request.GetArguments()["assignees"]; ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assigneesMap != nil {
					task.Assignees = new(twapi.UserGroups)

					err = twmcp.ParamGroup(assigneesMap,
						twmcp.OptionalNumericListParam(&task.Assignees.UserIDs, "user-ids"),
						twmcp.OptionalNumericListParam(&task.Assignees.CompanyIDs, "company-ids"),
						twmcp.OptionalNumericListParam(&task.Assignees.TeamIDs, "team-ids"),
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
