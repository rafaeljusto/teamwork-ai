package timer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twtimer "github.com/rafaeljusto/teamwork-ai/internal/teamwork/timer"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-timers",
			mcp.WithDescription("Retrieve multiple timers in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithNumber("user-id",
				mcp.Description("The ID of the user to filter timers by. "+
					"Only timers associated with this user will be returned."),
			),
			mcp.WithNumber("task-id",
				mcp.Description("The ID of the task to filter timers by. "+
					"Only timers associated with this task will be returned."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to filter timers by. "+
					"Only timers associated with this project will be returned."),
			),
			mcp.WithBoolean("running-timers-only",
				mcp.Description("If true, only running timers will be returned. "+
					"Defaults to false, which returns all timers."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twtimer.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.UserID, "user-id"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.TaskID, "task-id"),
				twmcp.OptionalNumericParam(&multiple.Request.Filters.ProjectID, "project-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.RunningTimersOnly, "running-timers-only"),
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
		mcp.NewTool("retrieve-timer",
			mcp.WithDescription("Retrieve a specific timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithNumber("timer-id",
				mcp.Required(),
				mcp.Description("The ID of the timer."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twtimer.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "timer-id"),
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

	mcpServer.AddTool(
		mcp.NewTool("create-timer",
			mcp.WithDescription("Create a new timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithString("description",
				mcp.Description("A description of the timer."),
			),
			mcp.WithBoolean("billable",
				mcp.Description("If true, the timer is billable. Defaults to false."),
			),
			mcp.WithBoolean("running",
				mcp.Description("If true, the timer will start running immediately."),
			),
			mcp.WithNumber("seconds",
				mcp.Description("The number of seconds to set the timer for."),
			),
			mcp.WithBoolean("stop-running-timers",
				mcp.Description("If true, any other running timers will be stopped when this timer is created."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to associate the timer with."),
			),
			mcp.WithNumber("task-id",
				mcp.Description("The ID of the task to associate the timer with."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timer twtimer.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalPointerParam(&timer.Description, "description"),
				twmcp.OptionalPointerParam(&timer.Billable, "billable"),
				twmcp.OptionalPointerParam(&timer.Running, "running"),
				twmcp.OptionalNumericPointerParam(&timer.Seconds, "seconds"),
				twmcp.OptionalPointerParam(&timer.StopRunningTimers, "stop-running-timers"),
				twmcp.OptionalNumericPointerParam(&timer.ProjectID, "project-id"),
				twmcp.OptionalNumericPointerParam(&timer.TaskID, "task-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timer created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-timer",
			mcp.WithDescription("Update a timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithNumber("timer-id",
				mcp.Required(),
				mcp.Description("The ID of the timer to update."),
			),
			mcp.WithString("description",
				mcp.Description("A description of the timer."),
			),
			mcp.WithBoolean("billable",
				mcp.Description("If true, the timer is billable."),
			),
			mcp.WithBoolean("running",
				mcp.Description("If true, the timer will start running immediately."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to associate the timer with."),
			),
			mcp.WithNumber("task-id",
				mcp.Description("The ID of the task to associate the timer with."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timer twtimer.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&timer.ID, "timer-id"),
				twmcp.OptionalPointerParam(&timer.Description, "description"),
				twmcp.OptionalPointerParam(&timer.Billable, "billable"),
				twmcp.OptionalPointerParam(&timer.Running, "running"),
				twmcp.OptionalNumericPointerParam(&timer.ProjectID, "project-id"),
				twmcp.OptionalNumericPointerParam(&timer.TaskID, "task-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timer updated successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("pause-timer",
			mcp.WithDescription("Pause a running timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithNumber("timer-id",
				mcp.Required(),
				mcp.Description("The ID of the timer to update."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timer twtimer.Pause

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&timer.Request.Path.ID, "timer-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timer paused successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("complete-timer",
			mcp.WithDescription("Complete a running timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs. "+
				"A timer must have a project ID associated with it to be completed. "+
				"The user should be a member of the project to log the time."),
			mcp.WithNumber("timer-id",
				mcp.Required(),
				mcp.Description("The ID of the timer to update."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timer twtimer.Complete

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&timer.Request.Path.ID, "timer-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timer completed successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("resume-timer",
			mcp.WithDescription("Resume a running timer in a customer site of Teamwork.com. "+
				"Timer is used to track ongoing work that will generate timelogs."),
			mcp.WithNumber("timer-id",
				mcp.Required(),
				mcp.Description("The ID of the timer to update."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timer twtimer.Resume

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&timer.Request.Path.ID, "timer-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timer); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timer resumed successfully"), nil
		},
	)
}
