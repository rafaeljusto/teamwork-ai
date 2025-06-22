package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twactivity "github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool(twmcp.MethodRetrieveActivities.String(),
			mcp.WithDescription("Retrieve multiple activities in a customer site of Teamwork.com. "+
				"Feed of all activity across your projects, including updates to various project items."),
			mcp.WithString("start-date",
				mcp.Description("Start date to filter activities. The date format follows RFC3339 - YYYY-MM-DDTHH:MM:SSZ."),
			),
			mcp.WithString("end-date",
				mcp.Description("End date to filter activities. The date format follows RFC3339 - YYYY-MM-DDTHH:MM:SSZ."),
			),
			mcp.WithArray("log-item-types",
				mcp.Description("Filter activities by item types."),
				mcp.Items(map[string]any{
					"type": "string",
					"enum": []any{
						"message",
						"comment",
						"task",
						"tasklist",
						"taskgroup",
						"milestone",
						"file",
						"form",
						"notebook",
						"timelog",
						"task_comment",
						"notebook_comment",
						"file_comment",
						"link_comment",
						"milestone_comment",
						"project",
						"link",
						"billingInvoice",
						"risk",
						"projectUpdate",
						"reacted",
						"budget",
					},
				}),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twactivity.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalTimeParam(&multiple.Request.Filters.StartDate, "start-date"),
				twmcp.OptionalTimeParam(&multiple.Request.Filters.EndDate, "end-date"),
				twmcp.OptionalListParam(&multiple.Request.Filters.LogItemTypes, "log-item-types"),
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
		mcp.NewTool(twmcp.MethodRetrieveProjectActivities.String(),
			mcp.WithDescription("Retrieve multiple activities from a project in a customer site of Teamwork.com. "+
				"Feed of all activity within a project."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project to retrieve activities from."),
			),
			mcp.WithString("start-date",
				mcp.Description("Start date to filter activities. The date format follows RFC3339 - YYYY-MM-DDTHH:MM:SSZ."),
			),
			mcp.WithString("end-date",
				mcp.Description("End date to filter activities. The date format follows RFC3339 - YYYY-MM-DDTHH:MM:SSZ."),
			),
			mcp.WithArray("log-item-types",
				mcp.Description("Filter activities by item types."),
				mcp.Items(map[string]any{
					"type": "string",
					"enum": []any{
						"message",
						"comment",
						"task",
						"tasklist",
						"taskgroup",
						"milestone",
						"file",
						"form",
						"notebook",
						"timelog",
						"task_comment",
						"notebook_comment",
						"file_comment",
						"link_comment",
						"milestone_comment",
						"project",
						"link",
						"billingInvoice",
						"risk",
						"projectUpdate",
						"reacted",
						"budget",
					},
				}),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twactivity.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.ProjectID, "project-id"),
				twmcp.OptionalTimeParam(&multiple.Request.Filters.StartDate, "start-date"),
				twmcp.OptionalTimeParam(&multiple.Request.Filters.EndDate, "end-date"),
				twmcp.OptionalListParam(&multiple.Request.Filters.LogItemTypes, "log-item-types"),
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
}
