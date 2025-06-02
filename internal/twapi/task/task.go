// Package task implements the API layer for managing tasks in Teamwork.com. It
// provides structures and methods for creating, updating, retrieving, and
// listing tasks.
package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
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
	ID                     int64      `json:"id"`
	Name                   string     `json:"name"`
	Description            *string    `json:"description"`
	DescriptionContentType *string    `json:"descriptionContentType"`
	Priority               *string    `json:"priority"`
	Progress               int64      `json:"progress"`
	StartAt                *time.Time `json:"startDate"`
	DueAt                  *time.Time `json:"dueDate"`
	EstimatedMinutes       int64      `json:"estimateMinutes"`

	Tasklist  twapi.Relationship   `json:"tasklist"`
	Assignees []twapi.Relationship `json:"assignees"`
	Tags      []twapi.Relationship `json:"tags"`

	CreatedBy     *int64     `json:"createdBy"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedBy     *int64     `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedBy     *int64     `json:"deletedBy"`
	DeletedAt     *time.Time `json:"deletedAt"`
	CompletedBy   *int64     `json:"completedBy,omitempty"`
	CompletedDate *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	WebLink       *string    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (t *Task) PopulateResourceWebLink(server string) {
	if t.ID == 0 {
		return
	}
	t.WebLink = twapi.Ref(fmt.Sprintf("%s/app/tasks/%d", server, t.ID))
}

// Single represents a request to retrieve a single task by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasks-task-id-json
type Single Task

// HTTPRequest creates an HTTP request to retrieve a single task by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, s.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Single instance.
func (s *Single) UnmarshalJSON(data []byte) error {
	var raw struct {
		Task Task `json:"task"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Task)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Task)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple tasks.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasks-json
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-projects-project-id-tasks-json
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/get-projects-api-v3-tasklists-tasklist-id-tasks-json
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID  int64
			TasklistID int64
		}
		Filters struct {
			SearchTerm   string
			TagIDs       []int64
			MatchAllTags *bool
			Page         int64
			PageSize     int64
		}
	}
	Response struct {
		Meta struct {
			Page struct {
				HasMore bool `json:"hasMore"`
			} `json:"page"`
		} `json:"meta"`
		Tasks []Task `json:"tasks"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple tasks.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case m.Request.Path.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/tasks.json", server, m.Request.Path.ProjectID)
	case m.Request.Path.TasklistID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, m.Request.Path.TasklistID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/tasks.json", server)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
	}
	if len(m.Request.Filters.TagIDs) > 0 {
		tagIDs := make([]string, len(m.Request.Filters.TagIDs))
		for i, id := range m.Request.Filters.TagIDs {
			tagIDs[i] = strconv.FormatInt(id, 10)
		}
		query.Set("tagIds", strings.Join(tagIDs, ","))
	}
	if m.Request.Filters.MatchAllTags != nil {
		query.Set("matchAllTags", strconv.FormatBool(*m.Request.Filters.MatchAllTags))
	}
	if m.Request.Filters.Page > 0 {
		query.Set("page", strconv.FormatInt(m.Request.Filters.Page, 10))
	}
	if m.Request.Filters.PageSize > 0 {
		query.Set("pageSize", strconv.FormatInt(m.Request.Filters.PageSize, 10))
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (m *Multiple) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Response)
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (m *Multiple) PopulateResourceWebLink(server string) {
	for i := range m.Response.Tasks {
		m.Response.Tasks[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new task in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/post-projects-api-v3-tasklists-tasklist-id-tasks-json
type Create struct {
	Name             string      `json:"name"`
	Description      *string     `json:"description,omitempty"`
	Priority         *string     `json:"priority,omitempty"`
	Progress         *int64      `json:"progress,omitempty"`
	StartAt          *twapi.Date `json:"startAt,omitempty"`
	DueAt            *twapi.Date `json:"dueAt,omitempty"`
	EstimatedMinutes *int64      `json:"estimatedMinutes,omitempty"`

	TasklistID int64             `json:"-"`
	Assignees  *twapi.UserGroups `json:"assignees,omitempty"`
	TagIDs     []int64           `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new task in a specific
// tasklist.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, c.TasklistID)
	payload := struct {
		Task Create `json:"task"`
	}{Task: c}
	body, err := json.Marshal(payload)
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
	ID               int64       `json:"-"`
	Name             *string     `json:"name,omitempty"`
	Description      *string     `json:"description,omitempty"`
	Priority         *string     `json:"priority,omitempty"`
	Progress         *int64      `json:"progress,omitempty"`
	StartAt          *twapi.Date `json:"startAt,omitempty"`
	DueAt            *twapi.Date `json:"dueAt,omitempty"`
	EstimatedMinutes *int64      `json:"estimatedMinutes,omitempty"`

	TasklistID *int64            `json:"tasklistId,omitempty"`
	Assignees  *twapi.UserGroups `json:"assignees,omitempty"`
	TagIDs     []int64           `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update an existing task in
// Teamwork.com.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, u.ID)
	payload := struct {
		Task Update `json:"task"`
	}{Task: u}
	body, err := json.Marshal(payload)
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

// Delete represents the payload for deleting an existing task in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tasks/delete-projects-api-v3-tasks-task-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a task.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
