package actions_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic/actions"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/jobrole"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/skill"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/user"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
)

func Test_SummarizeActivities(t *testing.T) {
	tests := []struct {
		name      string
		resources *config.Resources
		options   []actions.SummarizeActivitiesOption
	}{{
		name: "it should summarize activities for the last 365 days",
		resources: &config.Resources{
			TeamworkEngine: engineMock{
				do: teamworkEngine([]user.User{
					{ID: 1, FirstName: "James", LastName: "Smith"},
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, false, false),
			},
			Agentic: agenticMock{
				findTaskSkillsAndJobRoles: func(
					_ context.Context,
					taskData webhook.TaskData,
					availableSkills []skill.Skill,
					availableJobRoles []jobrole.JobRole,
				) ([]int64, []int64, string, error) {
					if taskData.Task.ID != 1 {
						return nil, nil, "", fmt.Errorf("unexpected task ID: %d", taskData.Task.ID)
					}
					if len(availableSkills) != 2 {
						return nil, nil, "", fmt.Errorf("unexpected number of skills: %d", len(availableSkills))
					}
					if len(availableJobRoles) != 2 {
						return nil, nil, "", fmt.Errorf("unexpected number of job roles: %d", len(availableJobRoles))
					}
					return []int64{1}, []int64{}, "Some interesting explanation.", nil
				},
			},
			Logger: slog.New(slog.DiscardHandler),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := actions.SummarizeActivities(
				context.Background(),
				tt.resources,
				tt.options...,
			); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
