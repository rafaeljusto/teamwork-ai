package anthropic

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindTaskSkillsAndJobRoles finds the skills and job roles for a given task. It
// uses the task data, available skills, and available job roles to determine
// the most relevant skills and job roles IDs for the task.
func (a *anthropic) FindTaskSkillsAndJobRoles(
	ctx context.Context,
	promptMessages []*mcp.PromptMessage,
) ([]int64, []int64, string, error) {
	var aiRequest request
	aiRequest.Model = a.model
	aiRequest.MaxTokens = 1024

	for _, msg := range promptMessages {
		textContent, ok := msg.Content.(*mcp.TextContent)
		if !ok {
			return nil, nil, "", fmt.Errorf("unsupported prompt message content type: %T", msg)
		}
		if textContent == nil {
			return nil, nil, "", fmt.Errorf("nil text content in prompt message")
		}
		switch msg.Role {
		case "system":
			aiRequest.addSystemMessage(textContent.Text)
		case "user":
			aiRequest.addUserMessage(textContent.Text)
		default:
			return nil, nil, "", fmt.Errorf("unknown prompt message role: %s", msg.Role)
		}
	}

	aiResponse, err := a.do(ctx, aiRequest)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to find task skills and job roles: %w", err)
	}

	var skillAndJobRoles struct {
		SkillIDs   []int64 `json:"skillIds"`
		JobRoleIDs []int64 `json:"jobRoleIds"`
		Reasoning  string  `json:"reasoning"`
	}
	if err := aiResponse.decode(&skillAndJobRoles); err != nil {
		return nil, nil, "", fmt.Errorf("failed to decode task skills and job roles: %w", err)
	}
	return skillAndJobRoles.SkillIDs, skillAndJobRoles.JobRoleIDs, skillAndJobRoles.Reasoning, nil
}
