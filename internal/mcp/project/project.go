package project

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twproject "github.com/rafaeljusto/teamwork-ai/internal/teamwork/project"
)

var resourceList = mcp.NewResource("twapi://projects", "projects",
	mcp.WithResourceDescription("Projects are central hubs to manage all of the components relating to what your team "+
		"are working on."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://projects/{id}", "task",
	mcp.WithTemplateDescription("Project is central hubs to manage all of the components relating to what your team "+
		"is working on."),
	mcp.WithTemplateMIMEType("application/json"),
)

// Register registers the project resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage projects in a customer
// site of Teamwork.com. A project is a central hub to manage all of the
// components relating to what your team is working on. It also provides a list
// of all projects and allows for the retrieval of a specific project by its ID.
// It also provides tools to retrieve multiple projects, a specific project, and
// to create a new project.
func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var projects twproject.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &projects); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, project := range projects {
				encoded, err := json.Marshal(project)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://projects/%d", project.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reProjectID := regexp.MustCompile(`twapi://projects/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reProjectID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid project ID")
			}
			projectID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid project ID")
			}

			var project twproject.Single
			project.ID = projectID
			if err := resources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(project)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://projects/%d", project.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-projects",
			mcp.WithDescription("Retrieve multiple projects in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var projects twproject.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &projects); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(projects)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project",
			mcp.WithDescription("Retrieve a specific project in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
			mcp.WithNumber("projectId",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Single

			id, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if id == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			project.ID = int64(id)

			if err := resources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(project)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-project",
			mcp.WithDescription("Create a new project in a customer site of Teamwork.com. "+
				"A project is central hubs to manage all of the components relating to what your team is working on."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the project."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the project."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var project twproject.Creation
			var ok bool

			project.Name, ok = request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if project.Name == "" {
				return nil, fmt.Errorf("name is required")
			}

			err := twmcp.OptionalParam(request.Params.Arguments, &project.Description, "description")
			if err != nil {
				return nil, fmt.Errorf("invalid description: %w", err)
			}

			if err := resources.TeamworkEngine.Do(ctx, &project); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Project created successfully"), nil
		},
	)
}
