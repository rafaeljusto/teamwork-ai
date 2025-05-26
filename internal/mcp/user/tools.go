package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twuser "github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
	twworkload "github.com/rafaeljusto/teamwork-ai/internal/teamwork/workload"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-users",
			mcp.WithDescription("Retrieve multiple users, also know as people, in a customer site of Teamwork.com. "+
				"Users, also known as people, are the individuals who can be assigned to tasks."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter users by first or last names, or e-mail. "+
					"The user will be selected if each word of the term matches the first or last name, or e-mail, not "+
					"requiring that the word matches are in the same field."),
			),
			mcp.WithNumber("type",
				mcp.Description("Type of user to filter by. The available options are account, collaborator or contact."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twuser.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalParam(&multiple.Request.Filters.Type, "type",
					twmcp.RestrictValues("account", "collaborator", "contact"),
				),
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
		mcp.NewTool("retrieve-project-users",
			mcp.WithDescription("Retrieve users, also known as people, from a specific project."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve users."),
			),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter users by first or last names, or e-mail. "+
					"The user will be selected if each word of the term matches the first or last name, or e-mail, not "+
					"requiring that the word matches are in the same field."),
			),
			mcp.WithNumber("type",
				mcp.Description("Type of user to filter by. The available options are account, collaborator or contact."),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twuser.Multiple

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&multiple.Request.Path.ProjectID, "project-id"),
				twmcp.OptionalParam(&multiple.Request.Filters.SearchTerm, "search-term"),
				twmcp.OptionalParam(&multiple.Request.Filters.Type, "type",
					twmcp.RestrictValues("account", "collaborator", "contact"),
				),
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
		mcp.NewTool("retrieve-user",
			mcp.WithDescription("Retrieve a specific user, also know as person, in a customer site of Teamwork.com. "+
				"Users, also known as person, is an individual who can be assigned to tasks."),
			mcp.WithNumber("user-id",
				mcp.Required(),
				mcp.Description("The ID of the user."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var user twuser.Single

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&user.ID, "user-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &user); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(user)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-user",
			mcp.WithDescription("Create a new user, who can be assigned to tasks. "+
				"User, also known as person, is an individual who can be assigned to tasks."),
			mcp.WithString("first-name",
				mcp.Required(),
				mcp.Description("The first name of the user."),
			),
			mcp.WithString("last-name",
				mcp.Required(),
				mcp.Description("The last name of the user."),
			),
			mcp.WithString("title",
				mcp.Description("The job title of the user, such as 'Project Manager' or 'Senior Software Developer'."),
			),
			mcp.WithString("email",
				mcp.Required(),
				mcp.Description("The email address of the user."),
			),
			mcp.WithBoolean("admin",
				mcp.Description("Indicates whether the user is an administrator."),
			),
			mcp.WithString("type",
				mcp.Description("The type of user, such as 'account', 'collaborator', or 'contact'."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company to which the user belongs."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var user twuser.Create

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredParam(&user.FirstName, "first-name"),
				twmcp.RequiredParam(&user.LastName, "last-name"),
				twmcp.OptionalPointerParam(&user.Title, "title"),
				twmcp.RequiredParam(&user.Email, "email"),
				twmcp.OptionalPointerParam(&user.Admin, "admin"),
				twmcp.OptionalPointerParam(&user.Type, "type",
					twmcp.RestrictValues("account", "collaborator", "contact"),
				),
				twmcp.OptionalNumericPointerParam(&user.CompanyID, "company-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &user); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("User created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-user",
			mcp.WithDescription("Update an existing user, who can be assigned to tasks. "+
				"User, also known as person, is an individual who can be assigned to tasks."),
			mcp.WithNumber("user-id",
				mcp.Required(),
				mcp.Description("The ID of the user to update."),
			),
			mcp.WithString("first-name",
				mcp.Description("The first name of the user."),
			),
			mcp.WithString("last-name",
				mcp.Description("The last name of the user."),
			),
			mcp.WithString("title",
				mcp.Description("The job title of the user, such as 'Project Manager' or 'Senior Software Developer'."),
			),
			mcp.WithString("email",
				mcp.Description("The email address of the user."),
			),
			mcp.WithBoolean("admin",
				mcp.Description("Indicates whether the user is an administrator."),
			),
			mcp.WithString("type",
				mcp.Description("The type of user, such as 'account', 'collaborator', or 'contact'."),
			),
			mcp.WithNumber("company-id",
				mcp.Description("The ID of the company to which the user belongs."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var user twuser.Update

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.OptionalPointerParam(&user.FirstName, "first-name"),
				twmcp.OptionalPointerParam(&user.LastName, "last-name"),
				twmcp.OptionalPointerParam(&user.Title, "title"),
				twmcp.OptionalPointerParam(&user.Email, "email"),
				twmcp.OptionalPointerParam(&user.Admin, "admin"),
				twmcp.OptionalPointerParam(&user.Type, "type",
					twmcp.RestrictValues("account", "collaborator", "contact"),
				),
				twmcp.OptionalNumericPointerParam(&user.CompanyID, "company-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &user); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("User updated successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("project-users",
			mcp.WithDescription("Assign users to a specific project."),
			mcp.WithNumber("project-id",
				mcp.Required(),
				mcp.Description("The ID of the project to which users will be assigned."),
			),
			mcp.WithArray("user-ids",
				mcp.Required(),
				mcp.Description("An array of user IDs to assign to the project."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var projectUsers twuser.AddProject

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredNumericParam(&projectUsers.Request.Path.ProjectID, "project-id"),
				twmcp.OptionalNumericListParam(&projectUsers.Request.Users.IDs, "user-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &projectUsers); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Users assigned to project successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-users-workload",
			mcp.WithDescription("Retrieve the workload of users, also known as people, in a customer site of Teamwork.com. "+
				"The workload allows you to see the users' overall workload on a short term, day-to-day basis, allowing for "+
				"a more granular view of each person's capacity. An individual's capacity is based on their working hours, "+
				"returned in the workload response, versus the total estimated time on their assigned tasks (minus any "+
				"unavailable time assigned to them) in the selected time frame. A user is considered over capacity when "+
				"their capacity exceeds their working hours. Missing dates in the response should be interpreted as the user "+
				"not having any tasks assigned to them on that date and being available."),
			mcp.WithString("start-date",
				mcp.Required(),
				mcp.Description("The start date of the workload period. The date must be in the format YYYY-MM-DD."),
			),
			mcp.WithString("end-date",
				mcp.Required(),
				mcp.Description("The end date of the workload period. The date must be in the format YYYY-MM-DD."),
			),
			mcp.WithArray("user-ids",
				mcp.Description("An array of user IDs to assign to the project."),
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
			var single twworkload.Single
			single.Request.Filters.Include = []string{"users.workingHours.workingHoursEntry"}

			err := twmcp.ParamGroup(request.GetArguments(),
				twmcp.RequiredDateParam(&single.Request.Filters.StartDate, "start-date"),
				twmcp.RequiredDateParam(&single.Request.Filters.EndDate, "end-date"),
				twmcp.OptionalNumericListParam(&single.Request.Filters.UserIDs, "user-ids"),
				twmcp.OptionalNumericParam(&single.Request.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&single.Request.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &single); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(single.Response)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)
}
