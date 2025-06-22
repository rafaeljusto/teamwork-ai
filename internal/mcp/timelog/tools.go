package timelog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twtimelog "github.com/rafaeljusto/teamwork-ai/internal/twapi/timelog"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveTimelogs.String(),
			mcp.WithDescription("Retrieve multiple timelogs in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter timelogs by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match timelogs that have all the specified tags. "+
					"If false, the search will match timelogs that have any of the specified tags. "+
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
			var multiple twtimelog.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
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
		mcp.NewTool(twmcp.MethodRetrieveProjectTimelogs.String(),
			mcp.WithDescription("Retrieve multiple timelogs from a specific project in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project to retrieve timelogs from."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter timelogs by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match timelogs that have all the specified tags. "+
					"If false, the search will match timelogs that have any of the specified tags. "+
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
			var multiple twtimelog.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.ProjectID, "project-id"),
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
		mcp.NewTool(twmcp.MethodRetrieveTaskTimelogs.String(),
			mcp.WithDescription("Retrieve multiple timelogs from a specific task in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithNumber("task-id",
				mcp.Required(),
				mcp.Description("The ID of the task to retrieve timelogs from."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter timelogs by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match timelogs that have all the specified tags. "+
					"If false, the search will match timelogs that have any of the specified tags. "+
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
			var multiple twtimelog.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.TaskID, "task-id"),
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
		mcp.NewTool(twmcp.MethodRetrieveTimelog.String(),
			mcp.WithDescription("Retrieve a specific timelog in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithNumber("timelog-id",
				mcp.Required(),
				mcp.Description("The ID of the timelog."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twtimelog.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "timelog-id"),
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
		mcp.NewTool(twmcp.MethodCreateTimelog.String(),
			mcp.WithDescription("Create a new timelog in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithString("description",
				mcp.Description("A description of the timelog."),
			),
			mcp.WithString("date",
				mcp.Required(),
				mcp.Description("The date of the timelog in the format YYYY-MM-DD."),
			),
			mcp.WithString("time",
				mcp.Required(),
				mcp.Description("The time of the timelog in the format HH:MM:SS."),
			),
			mcp.WithBoolean("is-utc",
				mcp.Description("If true, the time is in UTC. Defaults to false."),
			),
			mcp.WithNumber("hours",
				mcp.Required(),
				mcp.Description("The number of hours spent on the timelog. Must be a positive integer."),
			),
			mcp.WithNumber("minutes",
				mcp.Required(),
				mcp.Description("The number of minutes spent on the timelog. Must be a positive integer less than 60, "+
					"otherwise the hours attribute should be incremented."),
			),
			mcp.WithBoolean("billable",
				mcp.Description("If true, the timelog is billable. Defaults to false."),
			),
			mcp.WithNumber("project-id",
				mcp.Description("The ID of the project to associate the timelog with. "+
					"Either project-id or task-id must be provided, but not both."),
			),
			mcp.WithNumber("task-id",
				mcp.Description("The ID of the task to associate the timelog with. "+
					"Either project-id or task-id must be provided, but not both."),
			),
			mcp.WithNumber("user-id",
				mcp.Description("The ID of the user to associate the timelog with. "+
					"Defaults to the authenticated user if not provided."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the timelog."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timelog twtimelog.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalPointerParam(&timelog.Description, "description"),
				twmcp.RequiredDateParam(&timelog.Date, "date"),
				twmcp.RequiredTimeOnlyParam(&timelog.Time, "time"),
				twmcp.OptionalParam(&timelog.IsUTC, "is-utc"),
				twmcp.RequiredNumericParam(&timelog.Hours, "hours"),
				twmcp.RequiredNumericParam(&timelog.Minutes, "minutes"),
				twmcp.OptionalParam(&timelog.Billable, "billable"),
				twmcp.OptionalNumericParam(&timelog.ProjectID, "project-id"),
				twmcp.OptionalNumericParam(&timelog.TaskID, "task-id"),
				twmcp.OptionalNumericPointerParam(&timelog.UserID, "user-id"),
				twmcp.OptionalNumericListParam(&timelog.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timelog); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timelog created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodUpdateTimelog.String(),
			mcp.WithDescription("Update a timelog in a customer site of Teamwork.com. "+
				"Timelog is record of the amount a user spent working on a task or project."),
			mcp.WithNumber("timelog-id",
				mcp.Required(),
				mcp.Description("The ID of the timelog to update."),
			),
			mcp.WithString("description",
				mcp.Description("A description of the timelog."),
			),
			mcp.WithString("date",
				mcp.Description("The date of the timelog in the format YYYY-MM-DD."),
			),
			mcp.WithString("time",
				mcp.Description("The time of the timelog in the format HH:MM:SS."),
			),
			mcp.WithBoolean("is-utc",
				mcp.Description("If true, the time is in UTC."),
			),
			mcp.WithNumber("hours",
				mcp.Description("The number of hours spent on the timelog. Must be a positive integer."),
			),
			mcp.WithNumber("minutes",
				mcp.Description("The number of minutes spent on the timelog. Must be a positive integer less than 60, "+
					"otherwise the hours attribute should be incremented."),
			),
			mcp.WithBoolean("billable",
				mcp.Description("If true, the timelog is billable."),
			),
			mcp.WithNumber("user-id",
				mcp.Description("The ID of the user to associate the timelog with."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the timelog."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var timelog twtimelog.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&timelog.ID, "timelog-id"),
				twmcp.OptionalPointerParam(&timelog.Description, "description"),
				twmcp.OptionalDatePointerParam(&timelog.Date, "date"),
				twmcp.OptionalTimeOnlyPointerParam(&timelog.Time, "time"),
				twmcp.OptionalPointerParam(&timelog.IsUTC, "is-utc"),
				twmcp.OptionalNumericPointerParam(&timelog.Hours, "hours"),
				twmcp.OptionalNumericPointerParam(&timelog.Minutes, "minutes"),
				twmcp.OptionalPointerParam(&timelog.Billable, "billable"),
				twmcp.OptionalNumericPointerParam(&timelog.UserID, "user-id"),
				twmcp.OptionalNumericListParam(&timelog.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &timelog); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Timelog updated successfully"), nil
		},
	)
}
