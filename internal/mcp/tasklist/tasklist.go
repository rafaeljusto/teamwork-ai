package tasklist

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
	twtasklist "github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasklist"
)

var resourceList = mcp.NewResource("twapi://tasklists", "tasklists",
	mcp.WithResourceDescription("Tasklists group tasks together in a project for better organization."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://tasklists/{id}", "task",
	mcp.WithTemplateDescription("Tasklist group tasks together in a project for better organization."),
	mcp.WithTemplateMIMEType("application/json"),
)

// Register registers the tasklist resources and tools with the MCP server. It
// provides functionality to retrieve, create, and manage tasklists in a
// customer site of Teamwork.com. A tasklist groups tasks together in a project
// for better organization. It also provides a list of all tasklists and allows
// for the retrieval of a specific tasklist by its ID. Additionally, it provides
// tools to retrieve multiple tasklists, a specific tasklist, create a new
// tasklist, and update an existing tasklist.
func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var tasklists twtasklist.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &tasklists); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, tasklist := range tasklists.Tasklists {
				encoded, err := json.Marshal(tasklist)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasklists/%d", tasklist.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reTasklistID := regexp.MustCompile(`twapi://tasklists/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reTasklistID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid tasklist ID")
			}
			tasklistID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid tasklist ID")
			}

			var tasklist twtasklist.Single
			tasklist.ID = tasklistID
			if err := resources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(tasklist)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://tasklists/%d", tasklist.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasklists",
			mcp.WithDescription("Retrieve multiple tasklists in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklists twtasklist.Multiple
			if err := resources.TeamworkEngine.Do(ctx, &tasklists); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklists)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-project-tasklists",
			mcp.WithDescription("Retrieve multiple tasklists from a specific project in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithNumber("projectId",
				mcp.Required(),
				mcp.Description("The ID of the project from which to retrieve tasklists."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklists twtasklist.Multiple

			projectID, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if projectID == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			tasklists.ProjectID = int64(projectID)

			if err := resources.TeamworkEngine.Do(ctx, &tasklists); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklists)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-tasklist",
			mcp.WithDescription("Retrieve a specific tasklist in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithNumber("tasklistId",
				mcp.Required(),
				mcp.Description("The ID of the task."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Single

			id, ok := request.Params.Arguments["tasklistId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid tasklistId")
			} else if id == 0 {
				return nil, fmt.Errorf("tasklistId is required")
			}
			tasklist.ID = int64(id)

			if err := resources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(tasklist)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-tasklist",
			mcp.WithDescription("Create a new tasklist in a customer site of Teamwork.com. "+
				"A tasklist group tasks together in a project for better organization."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the tasklist."),
			),
			mcp.WithNumber("projectId",
				mcp.Required(),
				mcp.Description("The ID of the project."),
			),
			mcp.WithString("description",
				mcp.Description("The description of the tasklist."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var tasklist twtasklist.Creation
			var ok bool

			tasklist.Name, ok = request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if tasklist.Name == "" {
				return nil, fmt.Errorf("name is required")
			}

			projectID, ok := request.Params.Arguments["projectId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid projectId")
			} else if projectID == 0 {
				return nil, fmt.Errorf("projectId is required")
			}
			tasklist.ProjectID = int64(projectID)

			description, ok, err := twmcp.OptionalParam[string](request.Params.Arguments, "description")
			if err != nil {
				return nil, fmt.Errorf("invalid description: %w", err)
			} else if ok {
				tasklist.Description = description
			}

			if err := resources.TeamworkEngine.Do(ctx, &tasklist); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Tasklist created successfully"), nil
		},
	)
}
