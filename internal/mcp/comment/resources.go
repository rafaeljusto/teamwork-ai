package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twcomment "github.com/rafaeljusto/teamwork-ai/internal/twapi/comment"
)

var resourceList = mcp.NewResource("twapi://comments", "comments",
	mcp.WithResourceDescription("Comments are messages or notes that can be added to various "+
		"objects in Teamwork, such as tasks, files, milestones, and notebooks."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://comments/{id}", "comment",
	mcp.WithTemplateDescription("Comment is a message or note that can be added to various "+
		"objects in Teamwork, such as tasks, files, milestones, and notebooks."),
	mcp.WithTemplateMIMEType("application/json"),
)

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twcomment.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, comment := range multiple.Response.Comments {
				encoded, err := json.Marshal(comment)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://comments/%d", comment.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reCommentID := regexp.MustCompile(`twapi://comments/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reCommentID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid comment ID")
			}
			commentID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid comment ID")
			}

			var comment twcomment.Single
			comment.ID = commentID
			if err := configResources.TeamworkEngine.Do(ctx, &comment); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(comment)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://comments/%d", comment.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)
}
