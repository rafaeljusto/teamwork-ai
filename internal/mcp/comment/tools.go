package comment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twcomment "github.com/rafaeljusto/teamwork-ai/internal/teamwork/comment"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-comments",
			mcp.WithDescription("Retrieve multiple comments in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter comments by the content, also know as body in the response. "+
					"Each word from the search term is used to match against the comment content."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to filter comments by who posted them."),
				mcp.Items(map[string]any{
					"type": "number",
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
			var multiple twcomment.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.UserIDs, "user-ids"),
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
		mcp.NewTool("retrieve-file-comments",
			mcp.WithDescription("Retrieve multiple comments from a file in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithNumber("file-id",
				mcp.Required(),
				mcp.Description("The ID of the file to retrieve comments from."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter comments by the content, also know as body in the response. "+
					"Each word from the search term is used to match against the comment content."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to filter comments by who posted them."),
				mcp.Items(map[string]any{
					"type": "number",
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
			var multiple twcomment.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.FileID, "file-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.UserIDs, "user-ids"),
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
		mcp.NewTool("retrieve-milestone-comments",
			mcp.WithDescription("Retrieve multiple comments from a milestone in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithNumber("milestone-id",
				mcp.Required(),
				mcp.Description("The ID of the milestone to retrieve comments from."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter comments by the content, also know as body in the response. "+
					"Each word from the search term is used to match against the comment content."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to filter comments by who posted them."),
				mcp.Items(map[string]any{
					"type": "number",
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
			var multiple twcomment.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.MilestoneID, "milestone-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.UserIDs, "user-ids"),
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
		mcp.NewTool("retrieve-notebook-comments",
			mcp.WithDescription("Retrieve multiple comments from a notebook in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, notebooks, files and notebooks."),
			mcp.WithNumber("notebook-id",
				mcp.Required(),
				mcp.Description("The ID of the notebook to retrieve comments from."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter comments by the content, also know as body in the response. "+
					"Each word from the search term is used to match against the comment content."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to filter comments by who posted them."),
				mcp.Items(map[string]any{
					"type": "number",
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
			var multiple twcomment.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.NotebookID, "notebook-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.UserIDs, "user-ids"),
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
		mcp.NewTool("retrieve-task-comments",
			mcp.WithDescription("Retrieve multiple comments from a task in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithNumber("task-id",
				mcp.Required(),
				mcp.Description("The ID of the task to retrieve comments from."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter comments by the content, also know as body in the response. "+
					"Each word from the search term is used to match against the comment content."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("A list of user IDs to filter comments by who posted them."),
				mcp.Items(map[string]any{
					"type": "number",
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
			var multiple twcomment.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.TaskID, "task-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Request.Filters.UserIDs, "user-ids"),
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
		mcp.NewTool("retrieve-comment",
			mcp.WithDescription("Retrieve a specific comment in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithNumber("comment-id",
				mcp.Required(),
				mcp.Description("The ID of the comment."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twcomment.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&single.ID, "comment-id"),
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
		mcp.NewTool("create-comment",
			mcp.WithDescription("Create a new comment in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithObject("object",
				mcp.Required(),
				mcp.Description("The object to create the comment for. "+
					"It can be a tasks, messages, milestones, files or notebooks."),
				mcp.Properties(map[string]any{
					"type": map[string]any{
						"type": "string",
						"enum": []string{"tasks", "messages", "milestones", "files", "notebooks"},
						"description": "The type of object to create the comment for. " +
							"It can be a tasks, messages, milestones, files or notebooks.",
					},
					"id": map[string]any{
						"type":        "number",
						"description": "The ID of the object to create the comment for.",
					},
				}),
			),
			mcp.WithString("body",
				mcp.Required(),
				mcp.Description("The content of the comment. The content can be added as text or HTML."),
			),
			mcp.WithString("content-type",
				mcp.Description("The content type of the comment. It can be either 'TEXT' or 'HTML'."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var comment twcomment.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredParam(&comment.Body, "body"),
				twmcp.OptionalPointerParam(&comment.ContentType, "content-type"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			object, ok := request.GetArguments()["object"]
			if !ok {
				return nil, fmt.Errorf("missing required parameter: object")
			}
			objectMap, ok := object.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid object: expected an object, got %T", object)
			} else if objectMap == nil {
				return nil, fmt.Errorf("object cannot be nil")
			}
			err = twmcp.ParamGroup(objectMap,
				twmcp.RequiredParam(&comment.Object.Type, "type"),
				twmcp.RequiredNumericParam(&comment.Object.ID, "id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid object: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &comment); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Comment created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-comment",
			mcp.WithDescription("Update a comment in a customer site of Teamwork.com. "+
				"Within Teamwork.com, you can comment on project items such as tasks, milestones, files and notebooks."),
			mcp.WithNumber("comment-id",
				mcp.Required(),
				mcp.Description("The ID of the comment to update."),
			),
			mcp.WithString("body",
				mcp.Required(),
				mcp.Description("The content of the comment. The content can be added as text or HTML."),
			),
			mcp.WithString("content-type",
				mcp.Description("The content type of the comment. It can be either 'TEXT' or 'HTML'."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var comment twcomment.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&comment.ID, "comment-id"),
				twmcp.RequiredParam(&comment.Body, "body"),
				twmcp.OptionalPointerParam(&comment.ContentType, "content-type"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &comment); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Comment updated successfully"), nil
		},
	)
}
