package milestone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	twmilestone "github.com/rafaeljusto/teamwork-ai/internal/teamwork/milestone"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-milestones",
			mcp.WithDescription("Retrieve multiple milestones in a customer site of Teamwork.com. "+
				"Milestone is a target date representing a point of progress, or goal within a "+
				"project, that you can use task lists to track progress towards."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter milestones by name. "+
					"Each word from the search term is used to match against the milestone name and description. "+
					"The milestone will be selected if each word of the term matches the milestone name or description, not "+
					"requiring that the word matches are in the same field."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter milestones by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match milestones that have all the specified tags. "+
					"If false, the search will match milestones that have any of the specified tags. "+
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
			var multiple twmilestone.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
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
		mcp.NewTool("retrieve-project-milestones",
			mcp.WithDescription("Retrieve multiple milestones from a specific project in a customer site of Teamwork.com. "+
				"Milestone is a target date representing a point of progress, or goal within a "+
				"project, that you can use task lists to track progress towards."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project to retrieve milestones from."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter milestones by name. "+
					"Each word from the search term is used to match against the milestone name and description. "+
					"The milestone will be selected if each word of the term matches the milestone name or description, not "+
					"requiring that the word matches are in the same field."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter milestones by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithBoolean("match-all-tags",
				mcp.Description("If true, the search will match milestones that have all the specified tags. "+
					"If false, the search will match milestones that have any of the specified tags. "+
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
			var multiple twmilestone.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
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
		mcp.NewTool("retrieve-milestone",
			mcp.WithDescription("Retrieve a specific milestone in a customer site of Teamwork.com. "+
				"Milestone is a target date representing a point of progress, or goal within a "+
				"project, that you can use task lists to track progress towards."),
			mcp.WithNumber("milestone-id",
				mcp.Required(),
				mcp.Description("The ID of the milestone."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twmilestone.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&single.ID, "milestone-id"),
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
		mcp.NewTool("create-milestone",
			mcp.WithDescription("Create a new milestone in a customer site of Teamwork.com. "+
				"Milestone is a target date representing a point of progress, or goal within a "+
				"project, that you can use task lists to track progress towards."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the milestone."),
			),
			mcp.WithString("desciption",
				mcp.Description("A description of the milestone."),
			),
			mcp.WithString("due-date",
				mcp.Required(),
				mcp.Description("The due date of the milestone in the format YYYYMMDD."),
			),
			mcp.WithObject("assignees",
				mcp.Required(),
				mcp.Description("A list of assignees for the milestone. At least one assignee must be provided."),
				mcp.Properties(map[string]any{
					"user-ids": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the milestone.",
					},
					"company-ids": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the milestone.",
					},
					"team-ids": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the milestone.",
					},
				}),
			),
			mcp.WithArray("tasklist-ids",
				mcp.Description("A list of tasklist IDs to associate with the milestone."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the milestone."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var milestone twmilestone.Create

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&milestone.Name, "name"),
				twmcp.OptionalPointerParam(&milestone.Description, "description"),
				twmcp.RequiredLegacyDateParam(&milestone.DueDate, "due-date"),
				twmcp.OptionalNumericListParam(&milestone.TasklistIDs, "tasklist-ids"),
				twmcp.OptionalNumericListParam(&milestone.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			assignees, ok := request.Params.Arguments["assignees"]
			if !ok {
				return nil, fmt.Errorf("missing required parameter: assignees")
			}
			assigneesMap, ok := assignees.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid assignees: expected an object, got %T", assignees)
			} else if assigneesMap == nil {
				return nil, fmt.Errorf("assignees cannot be null")
			}
			err = twmcp.ParamGroup(assigneesMap,
				twmcp.OptionalNumericListParam(&milestone.Assignees.UserIDs, "user-ids"),
				twmcp.OptionalNumericListParam(&milestone.Assignees.CompanyIDs, "company-ids"),
				twmcp.OptionalNumericListParam(&milestone.Assignees.TeamIDs, "team-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid assignees: %w", err)
			}
			if milestone.Assignees.IsEmpty() {
				return nil, fmt.Errorf("at least one assignee must be provided")
			}

			if err := configResources.TeamworkEngine.Do(ctx, &milestone); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Milestone created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-milestone",
			mcp.WithDescription("Update a milestone in a customer site of Teamwork.com. "+
				"Milestone is a target date representing a point of progress, or goal within a "+
				"project, that you can use task lists to track progress towards."),
			mcp.WithNumber("milestone-id",
				mcp.Required(),
				mcp.Description("The ID of the milestone to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the milestone."),
			),
			mcp.WithString("desciption",
				mcp.Description("A description of the milestone."),
			),
			mcp.WithString("due-date",
				mcp.Description("The due date of the milestone in the format YYYYMMDD."),
			),
			mcp.WithObject("assignees",
				mcp.Description("A list of assignees for the milestone."),
				mcp.Properties(map[string]any{
					"user-ids": map[string]any{
						"type":        "array",
						"description": "List of user IDs assigned to the milestone.",
					},
					"company-ids": map[string]any{
						"type":        "array",
						"description": "List of company IDs assigned to the milestone.",
					},
					"team-ids": map[string]any{
						"type":        "array",
						"description": "List of team IDs assigned to the milestone.",
					},
				}),
			),
			mcp.WithArray("tasklist-ids",
				mcp.Description("A list of tasklist IDs to associate with the milestone."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the milestone."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var milestone twmilestone.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&milestone.ID, "milestone-id"),
				twmcp.OptionalPointerParam(&milestone.Name, "name"),
				twmcp.OptionalPointerParam(&milestone.Description, "description"),
				twmcp.OptionalLegacyDatePointerParam(&milestone.DueDate, "due-date"),
				twmcp.OptionalNumericListParam(&milestone.TasklistIDs, "tasklist-ids"),
				twmcp.OptionalNumericListParam(&milestone.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if assignees, ok := request.Params.Arguments["assignees"]; ok {
				assigneesMap, ok := assignees.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid assignees")
				} else if assigneesMap != nil {
					milestone.Assignees = new(teamwork.LegacyUserGroups)

					err = twmcp.ParamGroup(assigneesMap,
						twmcp.OptionalNumericListParam(&milestone.Assignees.UserIDs, "user-ids"),
						twmcp.OptionalNumericListParam(&milestone.Assignees.CompanyIDs, "company-ids"),
						twmcp.OptionalNumericListParam(&milestone.Assignees.TeamIDs, "team-ids"),
					)
					if err != nil {
						return nil, fmt.Errorf("invalid assignees: %w", err)
					}
				}
			}

			if err := configResources.TeamworkEngine.Do(ctx, &milestone); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Milestone updated successfully"), nil
		},
	)
}
