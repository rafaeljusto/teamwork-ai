package actions_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic/actions"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/comment"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/jobrole"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/skill"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/task"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/workload"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
)

func Test_AutoAssignTask(t *testing.T) {
	tests := []struct {
		name      string
		resources *config.Resources
		taskData  webhook.TaskData
		options   []actions.AutoAssignTaskOption
	}{{
		name: "it should assign a task and comment without checking rates or workload",
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
		taskData: func() webhook.TaskData {
			var taskData webhook.TaskData
			taskData.Task.ID = 1
			taskData.Task.Name = "task-1"
			return taskData
		}(),
		options: []actions.AutoAssignTaskOption{
			actions.WithAutoAssignTaskSkipRates(),
			actions.WithAutoAssignTaskSkipWorkload(),
		},
	}, {
		name: "it should assign a task and comment checking rates and not workload",
		resources: &config.Resources{
			TeamworkEngine: engineMock{
				do: teamworkEngine([]user.User{
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, true, false),
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
		taskData: func() webhook.TaskData {
			var taskData webhook.TaskData
			taskData.Task.ID = 1
			taskData.Task.Name = "task-1"
			return taskData
		}(),
		options: []actions.AutoAssignTaskOption{
			actions.WithAutoAssignTaskSkipWorkload(),
		},
	}, {
		name: "it should assign a task and comment checking workload and not rates",
		resources: &config.Resources{
			TeamworkEngine: engineMock{
				do: teamworkEngine([]user.User{
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, false, true),
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
		taskData: func() webhook.TaskData {
			var taskData webhook.TaskData
			taskData.Task.ID = 1
			taskData.Task.Name = "task-1"
			taskData.Task.StartDate = pointerTo(teamwork.Date(time.Now().AddDate(0, 0, 1)))
			taskData.Task.DueDate = pointerTo(teamwork.Date(time.Now().AddDate(0, 0, 2)))
			taskData.Task.EstimatedMinutes = 120
			return taskData
		}(),
		options: []actions.AutoAssignTaskOption{
			actions.WithAutoAssignTaskSkipRates(),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := actions.AutoAssignTask(
				context.Background(),
				tt.resources,
				tt.taskData,
				tt.options...,
			); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func pointerTo[T any](t T) *T {
	return &t
}

type engineMock struct {
	do func(context.Context, teamwork.Entity, ...teamwork.Option) error
}

func (e engineMock) Do(ctx context.Context, entity teamwork.Entity, optFuncs ...teamwork.Option) error {
	return e.do(ctx, entity, optFuncs...)
}

type agenticMock struct {
	findTaskSkillsAndJobRoles func(
		context.Context,
		webhook.TaskData,
		[]skill.Skill,
		[]jobrole.JobRole,
	) ([]int64, []int64, string, error)
}

func (a agenticMock) Init(string, *slog.Logger) error {
	return nil
}

func (a agenticMock) FindTaskSkillsAndJobRoles(
	ctx context.Context,
	taskData webhook.TaskData,
	availableSkills []skill.Skill,
	availableJobRoles []jobrole.JobRole,
) ([]int64, []int64, string, error) {
	return a.findTaskSkillsAndJobRoles(ctx, taskData, availableSkills, availableJobRoles)
}

func teamworkEngine(
	expectedAssignees []user.User,
	useRate, useWorkload bool,
) func(context.Context, teamwork.Entity, ...teamwork.Option) error {
	return func(_ context.Context, entity teamwork.Entity, _ ...teamwork.Option) error {
		switch t := entity.(type) {
		case *skill.Multiple:
			t.Response.Skills = []skill.Skill{
				{
					ID:   1,
					Name: "skill-1",
					Users: []teamwork.Relationship{
						{ID: 1, Type: "users"},
						{ID: 2, Type: "users"},
					},
				},
				{
					ID:   2,
					Name: "skill-2",
					Users: []teamwork.Relationship{
						{ID: 2, Type: "users"},
					},
				},
			}
		case *jobrole.Multiple:
			t.Response.JobRoles = []jobrole.JobRole{
				{
					ID:   1,
					Name: "jobrole-1",
					Users: []teamwork.Relationship{
						{ID: 1, Type: "users"},
						{ID: 2, Type: "users"},
					},
				},
				{
					ID:   2,
					Name: "jobrole-2",
					Users: []teamwork.Relationship{
						{ID: 2, Type: "users"},
					},
				},
			}
		case *user.Multiple:
			t.Response.Users = []user.User{
				{ID: 1, FirstName: "James", LastName: "Smith", Cost: pointerTo(teamwork.Money(20000))},
				{ID: 2, FirstName: "Michael", LastName: "Williams", Cost: pointerTo(teamwork.Money(10000))},
			}
		case *workload.Single:
			t.Response.Workload.Users = []workload.User{
				{
					ID: 1,
					Dates: map[teamwork.Date]workload.UserDate{
						teamwork.Date(time.Now().AddDate(0, 0, 1)): {
							Capacity:        87.5,
							CapacityMinutes: 420,
							UnavailableDay:  false,
						},
						teamwork.Date(time.Now().AddDate(0, 0, 2)): {
							UnavailableDay: true,
						},
					},
				},
				{
					ID: 2,
					Dates: map[teamwork.Date]workload.UserDate{
						teamwork.Date(time.Now().AddDate(0, 0, 1)): {
							Capacity:        10,
							CapacityMinutes: 48,
							UnavailableDay:  false,
						},
						teamwork.Date(time.Now().AddDate(0, 0, 2)): {
							Capacity:        80,
							CapacityMinutes: 384,
							UnavailableDay:  false,
						},
					},
				},
			}
		case *task.Update:
			if t.ID != 1 {
				return fmt.Errorf("unexpected task ID: %d", t.ID)
			}
			if len(t.Assignees.UserIDs) != len(expectedAssignees) {
				return fmt.Errorf("unexpected number of assigned users: %d", len(t.Assignees.UserIDs))
			}
			for i, expectedAssignee := range expectedAssignees {
				if t.Assignees.UserIDs[i] != expectedAssignee.ID {
					return fmt.Errorf("unexpected assigned user ID at index %d: %d", i, t.Assignees.UserIDs[i])
				}
			}
		case *comment.Create:
			if t.Object.Type != "tasks" {
				return fmt.Errorf("unexpected comment object type: %s", t.Object.Type)
			}
			if t.Object.ID != 1 {
				return fmt.Errorf("unexpected comment object ID: %d", t.Object.ID)
			}
			expectedBody := "🤖 Assignment of this task was performed by artificial intelligence.\n"
			for _, user := range expectedAssignees {
				expectedBody += fmt.Sprintf("\n  • %s %s", user.FirstName, user.LastName)
			}
			expectedBody += "\n\nSome interesting explanation."
			if useRate {
				expectedBody += " Concerns over user cost significantly impacted the decision."
			}
			if useWorkload {
				expectedBody += " Workload was a key consideration in the decision-making process."
			}
			if t.Body != expectedBody {
				return fmt.Errorf("unexpected comment body: %s", t.Body)
			}
		default:
			return fmt.Errorf("unexpected entity type: %T", t)
		}
		return nil
	}
}
