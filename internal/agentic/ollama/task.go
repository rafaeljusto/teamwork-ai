package ollama

import (
	"context"
	"fmt"

	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/jobrole"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/skill"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
)

// FindTaskSkillsAndJobRoles finds the skills and job roles for a given task. It
// uses the task data, available skills, and available job roles to determine
// the most relevant skills and job roles IDs for the task.
func (o *ollama) FindTaskSkillsAndJobRoles(
	ctx context.Context,
	taskData webhook.TaskData,
	availableSkills []skill.Skill,
	availableJobRoles []jobrole.JobRole,
) ([]int64, []int64, string, error) {
	var encodedSkills string
	for i, skill := range availableSkills {
		if i > 0 {
			encodedSkills += ", "
		}
		encodedSkills += fmt.Sprintf(`{"id": %d, "name": "%s"}`, skill.ID, skill.Name)
	}

	var encodedJobRoles string
	for i, jobRole := range availableJobRoles {
		if i > 0 {
			encodedJobRoles += ", "
		}
		encodedJobRoles += fmt.Sprintf(`{"id": %d, "name": "%s"}`, jobRole.ID, jobRole.Name)
	}

	var aiRequest request
	aiRequest.Model = o.model
	aiRequest.addUserMessage(findTaskSkillsAndJobRolesPrompt)
	aiRequest.addUserMessage("Project name: " + taskData.Project.Name)
	aiRequest.addUserMessage("Project description: " + taskData.Project.Description)
	aiRequest.addUserMessage("Tasklist name: " + taskData.Tasklist.Name)
	aiRequest.addUserMessage("Tasklist description: " + taskData.Tasklist.Description)
	aiRequest.addUserMessage("Task name: " + taskData.Task.Name)
	aiRequest.addUserMessage("Task description: " + taskData.Task.Description)
	aiRequest.addUserMessage("Available skills: " + encodedSkills)
	aiRequest.addUserMessage("Available job roles: " + encodedJobRoles)

	aiResponse, err := o.do(ctx, aiRequest, twmcp.MethodNone)
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

const findTaskSkillsAndJobRolesPrompt = `
You are an project manager expert. You have access to a list of skills and job
roles that can be used to complete a task. You are given a task with its name,
description, and the project it belongs to. You need to analyze the task and
suggest the best skills and job roles to complete it.

Please send back a JSON object with the skills and job role IDs. The format
MUST be:

{
  "skillIds": [1, 2],
  "jobRoleIds": [3, 4]
  "reasoning": "The reasoning behind the suggestions"
}

You MUST NOT send anything else, just the JSON object. If there are no skills or
job roles, send an empty array. Do not allucinate or make up any skills or job
roles.
`
