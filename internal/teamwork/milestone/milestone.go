// Package milestone provides functionality to manage milestones in
// Teamwork.com. It includes operations for retrieving, creating, and updating
// milestone information. It is part of the Teamwork AI project, which
// integrates with Teamwork.com to provide AI-driven insights and operations.
package milestone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Milestone represents a milestone in Teamwork.com.
type Milestone struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"deadline"`

	Project            teamwork.Relationship   `json:"project"`
	Tasklists          []teamwork.Relationship `json:"tasklists"`
	Tags               []teamwork.Relationship `json:"tags"`
	ResponsibleParties []teamwork.Relationship `json:"responsibleParties"`

	CreatedAt   *time.Time `json:"createdOn"`
	UpdatedAt   *time.Time `json:"lastChangedOn"`
	DeletedAt   *time.Time `json:"deletedOn"`
	CompletedAt *time.Time `json:"completedOn"`
	CompletedBy *int64     `json:"completedBy"`
	Completed   bool       `json:"completed"`
	Status      string     `json:"status"`
}

// Single represents a request to retrieve a single milestone by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/milestones/get-projects-api-v3-milestones-mileston-id-json
type Single Milestone

// HTTPRequest creates an HTTP request to retrieve a single milestone by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/milestones/%d.json", server, s.ID)
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
		Milestone Milestone `json:"milestone"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Milestone)
	return nil
}

// Multiple represents a request to retrieve multiple milestones.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/milestones/get-projects-api-v3-milestones-json
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID int64
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
		Milestones []Milestone `json:"milestones"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple milestones.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var uri string
	switch {
	case m.Request.Path.ProjectID > 0:
		uri = fmt.Sprintf("%s/projects/api/v3/projects/%d/milestones.json", server, m.Request.Path.ProjectID)
	default:
		uri = fmt.Sprintf("%s/projects/api/v3/milestones.json", server)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
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

// Create represents the payload for creating a new milestone in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/milestones/post-projects-id-milestones-json
type Create struct {
	Name        string              `json:"title"`
	Description *string             `json:"description,omitempty"`
	DueDate     teamwork.LegacyDate `json:"deadline"`

	ProjectID   int64                     `json:"-"`
	TasklistIDs []int64                   `json:"tasklistIds,omitempty"`
	TagIDs      []int64                   `json:"tagIds,omitempty"`
	Assignees   teamwork.LegacyUserGroups `json:"responsible-party-ids"`
}

// HTTPRequest creates an HTTP request to create a new milestone.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d/milestones.json", server, c.ProjectID)
	payload := struct {
		Milestone Create `json:"milestone"`
	}{Milestone: c}
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

// Update represents the payload for updating an existing milestone in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/milestones/put-milestones-id-json
type Update struct {
	ID          int64                `json:"-"`
	Name        *string              `json:"title,omitempty"`
	Description *string              `json:"description,omitempty"`
	DueDate     *teamwork.LegacyDate `json:"deadline,omitempty"`

	TasklistIDs []int64                    `json:"tasklistIds,omitempty"`
	TagIDs      []int64                    `json:"tagIds,omitempty"`
	Assignees   *teamwork.LegacyUserGroups `json:"responsible-party-ids,omitempty"`
}

// HTTPRequest creates an HTTP request to update a milestone.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/milestones/%d.json", server, u.ID)
	payload := struct {
		Milestone Update `json:"milestone"`
	}{Milestone: u}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Delete represents the payload for deleting an existing milestone in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/milestones/delete-milestones-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to update a milestone.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/milestones/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
