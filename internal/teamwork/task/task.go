package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Task represents a task in Teamwork.com. It contains information about the
// task such as its ID, name, status, description, priority, progress,
// associated tasklist, assignees, tags, start and due dates, estimate minutes,
// creation and update timestamps, and completion details. Tasks are used to
// manage work within projects, allowing teams to track progress, assign
// responsibilities, and organize tasks effectively.
//
// It is a fundamental component of project management in Teamwork.com, enabling
// teams to break down work into manageable units, assign them to team members,
// and monitor their progress towards completion. Tasks can be organized within
// tasklists, which group related tasks together for better organization within
// a project. Each task can have various attributes such as priority, status,
// and progress, which help in tracking and managing the work effectively.
type Task struct {
	ID                     int64                   `json:"id"`
	Name                   string                  `json:"name"`
	Status                 string                  `json:"status"`
	Description            *string                 `json:"description"`
	DescriptionContentType *string                 `json:"descriptionContentType"`
	Priority               *string                 `json:"priority"`
	Progress               int64                   `json:"progress"`
	Tasklist               teamwork.Relationship   `json:"tasklist"`
	Assignees              []teamwork.Relationship `json:"assignees"`
	Tags                   []teamwork.Relationship `json:"tags"`
	StartDate              *time.Time              `json:"startDate"`
	DueDate                *time.Time              `json:"dueDate"`
	EstimateMinutes        int64                   `json:"estimateMinutes"`
	CreatedBy              *int64                  `json:"createdBy"`
	CreatedAt              *time.Time              `json:"createdAt"`
	UpdatedBy              *int64                  `json:"updatedBy"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	DeletedBy              *int64                  `json:"deletedBy"`
	DeletedAt              *time.Time              `json:"deletedAt"`
	CompletedBy            *int64                  `json:"completedBy,omitempty"`
	CompletedDate          *time.Time              `json:"completedAt,omitempty"`
	Meta                   map[string]any          `json:"meta,omitempty"`
}

// Single represents a request to retrieve a single task by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasks-task-id-json
type Single Task

// Request creates an HTTP request to retrieve a single task by its ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, t.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Single instance.
func (t *Single) UnmarshalJSON(data []byte) error {
	var raw struct {
		Task Task `json:"task"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.Task)
	return nil
}

// Multiple represents a request to retrieve multiple tasks.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasks-json
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-projects-project-id-tasks-json
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasklists-tasklist-id-tasks-json
type Multiple struct {
	Tasks      []Task
	ProjectID  int64
	TasklistID int64
}

// Request creates an HTTP request to retrieve multiple tasks.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case t.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/tasks.json", server, t.ProjectID)
	case t.TasklistID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, t.TasklistID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/tasks.json", server)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (t *Multiple) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tasks []Task `json:"tasks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Tasks = raw.Tasks
	return nil
}

// Creation represents the payload for creating a new task in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/post-projects-api-v3-tasklists-tasklist-id-tasks-json
type Creation struct {
	Name        string               `json:"name"`
	TasklistID  int64                `json:"tasklistId"`
	Description string               `json:"description"`
	Assignees   *teamwork.UserGroups `json:"assignees,omitempty"`
	Priority    *string              `json:"priority,omitempty"`
}

// Request creates an HTTP request to create a new task in a specific tasklist.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, t.TasklistID)
	paylaod := struct {
		Task Creation `json:"task"`
	}{Task: t}
	body, err := json.Marshal(paylaod)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Update represents the payload for updating an existing task in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/patch-projects-api-v3-tasks-task-id-json
type Update struct {
	ID   int64
	Task struct {
		Name        *string              `json:"name,omitempty"`
		Description *string              `json:"description,omitempty"`
		Assignees   *teamwork.UserGroups `json:"assignees"`
		Priority    *string              `json:"priority,omitempty"`
	}
}

// Request creates an HTTP request to update an existing task in Teamwork.com.
func (t Update) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, t.ID)
	paylaod := struct {
		Task any `json:"task"`
	}{Task: t.Task}
	body, err := json.Marshal(paylaod)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
