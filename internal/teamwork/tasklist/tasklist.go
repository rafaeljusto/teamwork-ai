// Package tasklist implements the API layer for managing task lists in
// Teamwork.com. It provides functionality to create, update, retrieve, and
// delete task lists, as well as to retrieve multiple task lists associated with
// a project. Task lists are used to organize tasks within a project, allowing
// teams to group related tasks together for better organization and management.
package tasklist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Tasklist represents a tasklist in Teamwork.com. It contains information about
// the tasklist such as its ID, name, description, display order, associated
// project, milestone, status, whether it is pinned or private, lockdown ID,
// default task ID, billable status, budget, creation and update timestamps, and
// icon. Tasklists are used to organize tasks within a project, allowing teams
// to group related tasks together for better organization and management. They
// can be associated with a project and can also have milestones, which
// represent significant points or events within the project timeline.
type Tasklist struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Project   teamwork.Relationship  `json:"project"`
	Milestone *teamwork.Relationship `json:"milestone"`

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Status    string     `json:"status"`
}

// Single represents a request to retrieve a single tasklist by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-tasklists-tasklist-id
type Single Tasklist

// HTTPRequest creates an HTTP request to retrieve a single tasklist by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d.json", server, s.ID)
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
		Tasklist Tasklist `json:"tasklist"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Tasklist)
	return nil
}

// Multiple represents a request to retrieve multiple tasklists.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-tasklists
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-projects-project-id-tasklists
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID int64
		}
		Filters struct {
			SearchTerm string
			Page       int64
			PageSize   int64
		}
	}
	Response struct {
		Meta struct {
			Page struct {
				HasMore bool `json:"hasMore"`
			} `json:"page"`
		} `json:"meta"`
		Tasklists []Tasklist `json:"tasklists"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple tasklists.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case m.Request.Path.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/tasklists.json", server, m.Request.Path.ProjectID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/tasklists.json", server)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
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

// Creation represents the payload for creating a new tasklist in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/task-lists/post-projects-id-tasklists-json
type Creation struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	ProjectID   int64  `json:"-"`
	MilestoneID *int64 `json:"milestone-Id,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new tasklist.
func (c Creation) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d/tasklists.json", server, c.ProjectID)
	paylaod := struct {
		Tasklist Creation `json:"todo-list"`
	}{Tasklist: c}
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

// Update represents the payload for updating an existing tasklist in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/task-lists/put-tasklists-id-json
type Update struct {
	ID          int64   `json:"-"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`

	ProjectID   *int64 `json:"projectId,omitempty"`
	MilestoneID *int64 `json:"milestone-Id,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new tasklist.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/tasklists/%d.json", server, u.ID)
	paylaod := struct {
		Tasklist Update `json:"todo-list"`
	}{Tasklist: u}
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
