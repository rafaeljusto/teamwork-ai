package anthropic

import (
	"context"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
)

// SummarizeActivities summarizes the provided activities.
func (a *anthropic) SummarizeActivities(context.Context, []activity.Activity) (string, error) {
	// TODO(rafaeljusto): Figure out how to integrate the MCP server here, or
	// provide all tools to load the different activity item types.
	//
	// https://github.com/ollama/ollama/blob/main/docs/api.md#chat-request-with-tools
	return "", nil
}
