package openai

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
	"github.com/teamwork/twapi-go-sdk/projects"
)

var findTaskSkillsAndJobRolesCompiled = template.Must(template.New("prompt").Parse(findTaskSkillsAndJobRolesPrompt))

// FindTaskSkillsAndJobRoles finds the skills and job roles for a given task. It
// uses the task data, available skills, and available job roles to determine
// the most relevant skills and job roles IDs for the task.
func (o *openai) FindTaskSkillsAndJobRoles(
	ctx context.Context,
	taskData webhook.TaskData,
	availableSkills []projects.Skill,
	availableJobRoles []projects.JobRole,
) ([]int64, []int64, string, error) {
	var promptBuffer bytes.Buffer
	templateData := newFindTaskSkillsAndJobRolesData(taskData, availableSkills, availableJobRoles)
	if err := findTaskSkillsAndJobRolesCompiled.Execute(&promptBuffer, templateData); err != nil {
		return nil, nil, "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	var aiRequest request
	aiRequest.Model = o.model
	aiRequest.Input = promptBuffer.String()

	aiResponse, err := o.do(ctx, aiRequest)
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

type findTaskSkillsAndJobRolesData struct {
	Project struct {
		Name        string
		Description string
	}
	Tasklist struct {
		Name        string
		Description string
	}
	Task struct {
		Name        string
		Description string
	}
	Skills   []idName
	JobRoles []idName
}

type idName struct {
	ID   int64
	Name string
}

// Encode encodes the idName struct into a JSON string.
func (i idName) Encode() string {
	return fmt.Sprintf(`{"id":%d,"name":"%s"}`, i.ID, i.Name)
}

func newFindTaskSkillsAndJobRolesData(
	taskData webhook.TaskData,
	skills []projects.Skill,
	jobRoles []projects.JobRole,
) findTaskSkillsAndJobRolesData {
	var data findTaskSkillsAndJobRolesData
	data.Project.Name = taskData.Project.Name
	data.Project.Description = taskData.Project.Description
	data.Tasklist.Name = taskData.Tasklist.Name
	data.Tasklist.Description = taskData.Tasklist.Description
	data.Task.Name = taskData.Task.Name
	data.Task.Description = taskData.Task.Description

	for _, skill := range skills {
		data.Skills = append(data.Skills, idName{ID: skill.ID, Name: skill.Name})
	}

	for _, jobRole := range jobRoles {
		data.JobRoles = append(data.JobRoles, idName{ID: jobRole.ID, Name: jobRole.Name})
	}

	return data
}

//noling:lll
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

---
Project name: {{.Project.Name}}
---
Project description: {{.Project.Description}}
---
Tasklist name: {{.Tasklist.Name}}
---
Tasklist description: {{.Tasklist.Description}}
---
Task name: {{.Task.Name}}
---
Task description: {{.Task.Description}}
---
Available skills: {{range $i, $skill := .Skills}}{{if gt $i 0}},{{end}}{{$skill.Encode}}{{end}}
---
Available job roles: {{range $i, $jobRole := .JobRoles}}{{if gt $i 0}},{{end}}{{$jobRole.Encode}}{{end}}
`
