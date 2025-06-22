package ollama

import (
	"context"
	"fmt"
	"time"

	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
)

// SummarizeActivities summarizes the provided activities.
func (o *ollama) SummarizeActivities(ctx context.Context, activities []activity.Activity) (string, error) {
	var aiRequest request
	aiRequest.Model = o.model
	aiRequest.addUserMessage(summaryPrompt)

	for _, activity := range activities {
		encodedActivity := fmt.Sprintf("Activity: Action %q at %s to %s with ID %d",
			activity.Action,
			activity.At.Format(time.RFC3339),
			activity.Item.Type,
			activity.Item.ID,
		)
		if activity.Description != nil {
			encodedActivity += fmt.Sprintf(" - Description: ```%s```", *activity.Description)
		}
		if activity.ExtraDescription != nil {
			encodedActivity += fmt.Sprintf(" - Extra Description: ```%s```", *activity.ExtraDescription)
		}
		aiRequest.addUserMessage(encodedActivity)
	}

	aiResponse, err := o.do(ctx, aiRequest,
		// allow the LLM to request additional information from the objects related
		// to the activities via MCP tools
		twmcp.MethodRetrieveComment,
		twmcp.MethodRetrieveTask,
		twmcp.MethodRetrieveTasklist,
		twmcp.MethodRetrieveMilestone,
		twmcp.MethodRetrieveTimelog,
		twmcp.MethodRetrieveProject,
	)
	if err != nil {
		return "", fmt.Errorf("failed to summarize activities: %w", err)
	}
	return aiResponse.Message.Content, nil
}

const summaryPrompt = `
You are an expert project manager AI assistant. You will receive a list of
project activities. Your task is to interpret the provided activities and
generate a concise, informative summary in one or more paragraphs.

If no activities are provided, return an empty string.

Use any available tools or functions to retrieve additional context or details
about objects related to the activities when necessary.
`
