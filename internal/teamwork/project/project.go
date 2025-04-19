package project

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

// Project represents a project in Teamwork.com. It contains information about
// the project such as its ID, name, description, type, status, associated
// company, tags, start and end dates, category, ownership, updates, creation
// and update timestamps, completion details and project owner. It is used to
// manage and organize work, providing a central hub for all components related
// to what the team is working on.
type Project struct {
	ID          int64      `json:"id"`
	Description *string    `json:"description"`
	EndAt       *time.Time `json:"endAt"`
	Name        string     `json:"name"`
	StartAt     *time.Time `json:"startAt"`

	Category *teamwork.Relationship  `json:"category"`
	Company  teamwork.Relationship   `json:"company"`
	Owner    *teamwork.Relationship  `json:"projectOwner"`
	Tags     []teamwork.Relationship `json:"tags"`

	CreatedAt   *time.Time `json:"createdAt"`
	CreatedBy   *int64     `json:"createdBy"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	UpdatedBy   *int64     `json:"updatedBy"`
	CompletedAt *time.Time `json:"completedAt"`
	CompletedBy *int64     `json:"completedBy"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
}

// Single represents a request to retrieve a single project by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/projects/get-projects-api-v3-projects-project-id-json
type Single Project

// HTTPRequest creates an HTTP request to retrieve a single project by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/projects/%d.json", server, s.ID)
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
		Project Project `json:"project"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Project)
	return nil
}

// Multiple represents a request to retrieve multiple projects.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/projects/get-projects-api-v3-projects-json
type Multiple struct {
	Request struct {
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
		Projects []Project `json:"projects"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple projects.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/projects/api/v3/projects.json", nil)
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
		query.Set("projectTagIds", strings.Join(tagIDs, ","))
	}
	if m.Request.Filters.MatchAllTags != nil {
		query.Set("matchAllProjectTags", strconv.FormatBool(*m.Request.Filters.MatchAllTags))
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

// Creation represents the payload for creating a new project in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/post-projects-json
type Creation struct {
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	StartAt     *time.Time `json:"start-date,omitempty"`
	EndAt       *time.Time `json:"end-date,omitempty"`

	CategoryID *int64  `json:"category-id,omitempty"`
	CompanyID  int64   `json:"companyId"`
	OwnerID    *int64  `json:"projectOwnerId,omitempty"`
	Tags       []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new project.
func (c Creation) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects.json", server)
	paylaod := struct {
		Project Creation `json:"project"`
	}{Project: c}
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

// Update represents the payload for updating an existing project in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/put-projects-id-json
type Update struct {
	ID          int64      `json:"-"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartAt     *time.Time `json:"start-date,omitempty"`
	EndAt       *time.Time `json:"end-date,omitempty"`

	CategoryID *int64  `json:"category-id,omitempty"`
	CompanyID  *int64  `json:"companyId"`
	OwnerID    *int64  `json:"projectOwnerId,omitempty"`
	Tags       []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update a project.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects.json", server)
	paylaod := struct {
		Project Update `json:"project"`
	}{Project: u}
	body, err := json.Marshal(paylaod)
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
