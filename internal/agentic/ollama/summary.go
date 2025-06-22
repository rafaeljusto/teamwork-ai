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
You are an expert project manager AI assistant. You will be given a list of
project activities. Your task is to interpret the activities and generate an
informative summary in one or more paragraphs.

Important behavioral rules:
- DO NOT return a JSON response or expose raw data.
- DO NOT return internal IDs (e.g., project ID, milestone ID).
- DO NOT ask the user for permission to retrieve data.
- Whenever an ID is present and entity details are missing (e.g., name,
description, date), you MUST use tool calls immediately to retrieve those
details. This is mandatory.
- You MUST automatically fetch all relevant data via tools when needed, without
prompting or requesting clarification.
- After retrieving the data (returned in JSON), interpret and summarize it
naturally in text, without revealing the underlying JSON structure.

If no activities are provided, return an empty string.
`
