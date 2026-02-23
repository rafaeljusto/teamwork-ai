package actions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/agentic/actions"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
	twapi "github.com/teamwork/twapi-go-sdk"
	"github.com/teamwork/twapi-go-sdk/projects"
	"github.com/teamwork/twapi-go-sdk/session"
)

var reTypeID = regexp.MustCompile(`/([a-z]+)/([0-9]+)/comments\.json`)

func Test_AutoAssignTask(t *testing.T) {
	tests := []struct {
		name      string
		resources *config.Resources
		taskData  webhook.TaskData
		options   []actions.AutoAssignTaskOption
	}{{
		name: "it should assign a task and comment without checking rates or workload",
		resources: &config.Resources{
			TeamworkEngine: twapi.NewEngine(session.NewBasicAuth("john", "abc123", "example.com"),
				twapi.WithHTTPClient(teamworkEngine([]projects.User{
					{ID: 1, FirstName: "James", LastName: "Smith"},
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, false, false)),
			),
			Agentic: agenticMock{
				findTaskSkillsAndJobRoles: func(
					_ context.Context,
					promptMessages []*mcp.PromptMessage,
				) ([]int64, []int64, string, error) {
					if len(promptMessages) != 2 {
						return nil, nil, "", fmt.Errorf("unexpected number of prompts: %d", len(promptMessages))
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
			TeamworkEngine: twapi.NewEngine(session.NewBasicAuth("john", "abc123", "example.com"),
				twapi.WithHTTPClient(teamworkEngine([]projects.User{
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, true, false)),
			),
			Agentic: agenticMock{
				findTaskSkillsAndJobRoles: func(
					_ context.Context,
					promptMessages []*mcp.PromptMessage,
				) ([]int64, []int64, string, error) {
					if len(promptMessages) != 2 {
						return nil, nil, "", fmt.Errorf("unexpected number of prompts: %d", len(promptMessages))
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
			TeamworkEngine: twapi.NewEngine(session.NewBasicAuth("john", "abc123", "example.com"),
				twapi.WithHTTPClient(teamworkEngine([]projects.User{
					{ID: 2, FirstName: "Michael", LastName: "Williams"},
				}, false, true)),
			),
			Agentic: agenticMock{
				findTaskSkillsAndJobRoles: func(
					_ context.Context,
					promptMessages []*mcp.PromptMessage,
				) ([]int64, []int64, string, error) {
					if len(promptMessages) != 2 {
						return nil, nil, "", fmt.Errorf("unexpected number of prompts: %d", len(promptMessages))
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
			taskData.Task.StartDate = new(twapi.Date(time.Now().AddDate(0, 0, 1)))
			taskData.Task.DueDate = new(twapi.Date(time.Now().AddDate(0, 0, 2)))
			taskData.Task.EstimatedMinutes = 120
			return taskData
		}(),
		options: []actions.AutoAssignTaskOption{
			actions.WithAutoAssignTaskSkipRates(),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory transports are not goroutine safe, so we need a new MCP mock
			// per test case
			tt.resources.MCPClient = config.NewMCPClient(mockMCP(t, func(srv *mcp.Server) {
				srv.AddPrompt(&mcp.Prompt{
					Name: "twprojects_task_skills_and_roles",
					Arguments: []*mcp.PromptArgument{
						{Name: "task_id", Required: true},
					},
				}, mcp.PromptHandler(func(context.Context, *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
					return &mcp.GetPromptResult{
						Messages: []*mcp.PromptMessage{
							{
								Role: "system",
								Content: &mcp.TextContent{
									Text: "You are an expert at identifying skills and job roles required for tasks based on their details.",
								},
							},
							{
								Role: "user",
								Content: &mcp.TextContent{
									Text: "Given the details of the task, identify the relevant skills and job roles required to complete it.",
								},
							},
						},
					}, nil
				}))
			}))

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

type agenticMock struct {
	findTaskSkillsAndJobRoles func(
		context.Context,
		[]*mcp.PromptMessage,
	) ([]int64, []int64, string, error)
}

func (a agenticMock) Init(string, *slog.Logger) error {
	return nil
}

func (a agenticMock) FindTaskSkillsAndJobRoles(
	ctx context.Context,
	promptMessages []*mcp.PromptMessage,
) ([]int64, []int64, string, error) {
	return a.findTaskSkillsAndJobRoles(ctx, promptMessages)
}

func teamworkEngine(expectedAssignees []projects.User, useRate, useWorkload bool) twapi.HTTPClientFunc {
	return func(req *http.Request) (*http.Response, error) {
		var entity any
		status := http.StatusOK

		switch {
		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "example.com/projects/api/v3/skills.json"):
			entity = projects.SkillListResponse{
				Skills: []projects.Skill{
					{
						ID:   1,
						Name: "skill-1",
						Users: []twapi.Relationship{
							{ID: 1, Type: "users"},
							{ID: 2, Type: "users"},
						},
					},
					{
						ID:   2,
						Name: "skill-2",
						Users: []twapi.Relationship{
							{ID: 2, Type: "users"},
						},
					},
				},
			}

		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "example.com/projects/api/v3/jobroles.json"):
			entity = projects.JobRoleListResponse{
				JobRoles: []projects.JobRole{
					{
						ID:   1,
						Name: "jobrole-1",
						Users: []twapi.Relationship{
							{ID: 1, Type: "users"},
							{ID: 2, Type: "users"},
						},
					},
					{
						ID:   2,
						Name: "jobrole-2",
						Users: []twapi.Relationship{
							{ID: 2, Type: "users"},
						},
					},
				},
			}

		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "example.com/projects/api/v3/people.json"):
			entity = projects.UserListResponse{
				Users: []projects.User{
					{ID: 1, FirstName: "James", LastName: "Smith", Cost: new(twapi.Money(20000))},
					{ID: 2, FirstName: "Michael", LastName: "Williams", Cost: new(twapi.Money(10000))},
				},
			}

		case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "example.com/projects/api/v3/workload.json"):
			entity = projects.WorkloadResponse{
				Workload: projects.Workload{
					Users: []projects.WorkloadUser{
						{
							ID: 1,
							Dates: map[twapi.Date]projects.WorkloadUserDate{
								twapi.Date(time.Now().AddDate(0, 0, 1)): {
									Capacity:        87.5,
									CapacityMinutes: 420,
									UnavailableDay:  false,
								},
								twapi.Date(time.Now().AddDate(0, 0, 2)): {
									UnavailableDay: true,
								},
							},
						},
						{
							ID: 2,
							Dates: map[twapi.Date]projects.WorkloadUserDate{
								twapi.Date(time.Now().AddDate(0, 0, 1)): {
									Capacity:        10,
									CapacityMinutes: 48,
									UnavailableDay:  false,
								},
								twapi.Date(time.Now().AddDate(0, 0, 2)): {
									Capacity:        80,
									CapacityMinutes: 384,
									UnavailableDay:  false,
								},
							},
						},
					},
				},
			}

		case req.Method == http.MethodPut && strings.HasPrefix(req.URL.Path, "example.com/projects/api/v3/tasks/"):
			id := strings.TrimPrefix(req.URL.Path, "example.com/projects/api/v3/tasks/")
			id = strings.TrimSuffix(id, ".json")
			if id != "1" {
				return nil, fmt.Errorf("unexpected task ID: %s", id)
			}

			var t struct {
				Task projects.TaskUpdateRequest `json:"task"`
			}
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&t); err != nil {
				return nil, fmt.Errorf("failed to decode task create request: %w", err)
			}
			if t.Task.Assignees == nil && len(expectedAssignees) > 0 {
				return nil, fmt.Errorf("expected assignees but none were provided")
			}
			if len(t.Task.Assignees.UserIDs) != len(expectedAssignees) {
				return nil, fmt.Errorf("unexpected number of assigned users: %d", len(t.Task.Assignees.UserIDs))
			}
			for i, expectedAssignee := range expectedAssignees {
				if t.Task.Assignees.UserIDs[i] != expectedAssignee.ID {
					return nil, fmt.Errorf("unexpected assigned user ID at index %d: %d", i, t.Task.Assignees.UserIDs[i])
				}
			}

			entity = projects.TaskUpdateResponse{
				Task: projects.Task{
					ID: 1,
				},
			}

		case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/comments.json"):
			matches := reTypeID.FindStringSubmatch(req.URL.Path)
			if len(matches) != 3 {
				return nil, fmt.Errorf("failed to extract comment object type and ID from URL path: %s", req.URL.Path)
			}
			if matches[1] != "tasks" {
				return nil, fmt.Errorf("unexpected comment object type: %s", matches[1])
			}
			if matches[2] != "1" {
				return nil, fmt.Errorf("unexpected comment object ID: %s", matches[2])
			}

			var t struct {
				Comment projects.CommentCreateRequest `json:"comment"`
			}
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&t); err != nil {
				return nil, fmt.Errorf("failed to decode comment create request: %w", err)
			}

			var expectedBody strings.Builder
			expectedBody.WriteString("🤖 Assignment of this task was performed by artificial intelligence.\n")
			for _, user := range expectedAssignees {
				fmt.Fprintf(&expectedBody, "\n  • %s %s", user.FirstName, user.LastName)
			}
			expectedBody.WriteString("\n\nSome interesting explanation.")
			if useRate {
				expectedBody.WriteString(" Concerns over user cost significantly impacted the decision.")
			}
			if useWorkload {
				expectedBody.WriteString(" Workload was a key consideration in the decision-making process.")
			}
			if t.Comment.Body != expectedBody.String() {
				return nil, fmt.Errorf("unexpected comment body: %s", t.Comment.Body)
			}

			status = http.StatusCreated
			entity = projects.CommentCreateResponse{
				ID: 1,
			}

		default:
			return nil, fmt.Errorf("unexpected method %q and URL path: %q", req.Method, req.URL.Path)
		}

		var body io.ReadCloser
		if entity != nil {
			encoded, err := json.Marshal(entity)
			if err != nil {
				return nil, err
			}
			body = io.NopCloser(strings.NewReader(string(encoded)))
		}

		return &http.Response{
			StatusCode: status,
			Body:       body,
			Header:     make(http.Header),
		}, nil
	}
}

func mockMCP(t *testing.T, register func(*mcp.Server)) mcp.Transport {
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, &mcp.ServerOptions{})

	register(server)

	serverSession, err := server.Connect(t.Context(), serverTransport, nil)
	if err != nil {
		t.Fatalf("failed to initialize a MCP server session: %v", err)
	}
	t.Cleanup(func() {
		if err := serverSession.Close(); err != nil {
			t.Logf("failed to close MCP server session: %v", err)
		}
	})

	return clientTransport
}
