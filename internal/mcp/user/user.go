package user

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twuser "github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

var resourceList = mcp.NewResource("twapi://users", "users",
	mcp.WithResourceDescription("Users, also known as people, are the individuals who can be assigned to tasks."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://users/{id}", "user",
	mcp.WithTemplateDescription("User, also known as person, is an individual who can be assigned to tasks."),
	mcp.WithTemplateMIMEType("application/json"),
)

func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var users twuser.MultipleUsers
			if err := resources.TeamworkEngine.Do(&users); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, user := range users.Users {
				encoded, err := json.Marshal(user)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://users/%d", user.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reUserID := regexp.MustCompile(`twapi://users/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reUserID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid user ID")
			}
			userID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid user ID")
			}

			var user twuser.SingleUser
			user.ID = userID
			if err := resources.TeamworkEngine.Do(&user); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(user)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://users/%d", user.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-users",
			mcp.WithDescription("Users, also known as people, are the individuals who can be assigned to tasks."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var users twuser.MultipleUsers
			if err := resources.TeamworkEngine.Do(&users); err != nil {
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
			var users twuser.MultipleUsers

			projectID, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if projectID == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			users.ProjectID = int64(projectID)

			if err := resources.TeamworkEngine.Do(&users); err != nil {
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
			var user twuser.SingleUser

			id, ok := request.Params.Arguments["userId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid userId")
			} else if id == 0 {
				return nil, fmt.Errorf("userId is required")
			}
			user.ID = int64(id)

			if err := resources.TeamworkEngine.Do(&user); err != nil {
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
			var user twuser.UserCreation
			var ok bool

			user.FirstName, ok = request.Params.Arguments["firstName"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid first name")
			} else if user.FirstName == "" {
				return nil, fmt.Errorf("first name is required")
			}

			user.LastName, ok = request.Params.Arguments["lastName"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid last name")
			} else if user.LastName == "" {
				return nil, fmt.Errorf("last name is required")
			}

			user.Email, ok = request.Params.Arguments["email"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid email")
			} else if user.Email == "" {
				return nil, fmt.Errorf("email is required")
			}

			user.Password, ok = request.Params.Arguments["password"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid password")
			} else if user.Password == "" {
				return nil, fmt.Errorf("password is required")
			}

			if err := resources.TeamworkEngine.Do(&user); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("User created successfully"), nil
		},
	)
}
