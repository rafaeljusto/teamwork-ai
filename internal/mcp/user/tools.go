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
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-users",
			mcp.WithDescription("Users, also known as people, are the individuals who can be assigned to tasks."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var users twuser.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &users); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(users)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project-users",
			mcp.WithDescription("Retrieve users, also known as people, from a specific project."),
			mcp.WithNumber("projectId",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve users."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var users twuser.Multiple

			err := twmcp.NumericParam(request.Params.Arguments, &users.ProjectID, "projectId")
			if err != nil {
				return nil, fmt.Errorf("invalid project ID: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &users); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(users)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-user",
			mcp.WithDescription("Users, also known as person, is an individual who can be assigned to tasks."),
			mcp.WithNumber("userId",
				mcp.Required(),
				mcp.Description("The ID of the user."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var user twuser.Single

			err := twmcp.NumericParam(request.Params.Arguments, &user.ID, "userId")
			if err != nil {
				return nil, fmt.Errorf("invalid user ID: %w", err)
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
			mcp.WithDescription("User, also known as person, is an individual who can be assigned to tasks."),
			mcp.WithString("firstName",
				mcp.Required(),
				mcp.Description("The first name of the user."),
			),
			mcp.WithString("lastName",
				mcp.Required(),
				mcp.Description("The last name of the user."),
			),
			mcp.WithString("email",
				mcp.Description("The email address of the user."),
				mcp.Required(),
			),
			mcp.WithString("password",
				mcp.Description("The password for the user."),
				mcp.Required(),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var user twuser.Creation

			err := twmcp.Param(request.Params.Arguments, &user.FirstName, "firstName")
			if err != nil {
				return nil, fmt.Errorf("invalid first name: %w", err)
			}
			err = twmcp.Param(request.Params.Arguments, &user.LastName, "lastName")
			if err != nil {
				return nil, fmt.Errorf("invalid last name: %w", err)
			}
			err = twmcp.Param(request.Params.Arguments, &user.Email, "email")
			if err != nil {
				return nil, fmt.Errorf("invalid email: %w", err)
			}
			err = twmcp.Param(request.Params.Arguments, &user.Password, "password")
			if err != nil {
				return nil, fmt.Errorf("invalid password: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &user); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("User created successfully"), nil
		},
	)
}
