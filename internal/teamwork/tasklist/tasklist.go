package tasklist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	DisplayOrder  int64                  `json:"displayOrder"`
	ProjectID     int64                  `json:"projectId"`
	Project       teamwork.Relationship  `json:"project"`
	MilestoneID   *int64                 `json:"milestoneId"`
	Milestone     *teamwork.Relationship `json:"milestone"`
	IsPinned      bool                   `json:"isPinned"`
	IsPrivate     bool                   `json:"isPrivate"`
	LockdownID    *int64                 `json:"lockdownId"`
	Status        string                 `json:"status"`
	DefaultTaskID *int64                 `json:"defaultTaskId"`
	DefaultTask   *teamwork.Relationship `json:"defaultTask"`
	IsBillable    *bool                  `json:"isBillable"`
	Budget        *teamwork.Relationship `json:"tasklistBudget"`
	CreatedAt     *time.Time             `json:"createdAt"`
	UpdatedAt     *time.Time             `json:"updatedAt"`
	Icon          *string                `json:"icon"`
}

// Single represents a request to retrieve a single tasklist by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-tasklists-tasklist-id
type Single Tasklist

// Request creates an HTTP request to retrieve a single tasklist by its ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d.json", server, t.ID)
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
		Tasklist Tasklist `json:"tasklist"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.Tasklist)
	return nil
}

// Multiple represents a request to retrieve multiple tasklists.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-tasklists
// https://apidocs.teamwork.com/docs/teamwork/v3/task-lists/get-projects-api-v3-projects-project-id-tasklists
type Multiple struct {
	Tasklists []Tasklist
	ProjectID int64
}

// Request creates an HTTP request to retrieve multiple tasklists.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case t.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/tasklists.json", server, t.ProjectID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/tasklists.json", server)
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
		Tasklists []Tasklist `json:"tasklists"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Tasklists = raw.Tasklists
	return nil
}

// Creation represents the payload for creating a new tasklist in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/task-lists/post-projects-id-tasklists-json
type Creation struct {
	Name        string `json:"name"`
	ProjectID   int64  `json:"projectId"`
	Description string `json:"description"`
}

// Request creates an HTTP request to create a new tasklist.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d/tasklists.json", server, t.ProjectID)
	paylaod := struct {
		Tasklist Creation `json:"todo-list"`
	}{Tasklist: t}
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
