// Package project provides the API implementation for managing projects in
// Teamwork.com. It includes functionalities for retrieving single or multiple
// projects, creating new projects, and updating existing ones.
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
	Name        string     `json:"name"`
	StartAt     *time.Time `json:"startAt"`
	EndAt       *time.Time `json:"endAt"`

	Company teamwork.Relationship   `json:"company"`
	Owner   *teamwork.Relationship  `json:"projectOwner"`
	Tags    []teamwork.Relationship `json:"tags"`

	CreatedAt   *time.Time `json:"createdAt"`
	CreatedBy   *int64     `json:"createdBy"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	UpdatedBy   *int64     `json:"updatedBy"`
	CompletedAt *time.Time `json:"completedAt"`
	CompletedBy *int64     `json:"completedBy"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	WebLink     *string    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (p *Project) PopulateResourceWebLink(server string) {
	if p.ID == 0 {
		return
	}
	p.WebLink = teamwork.Ref(fmt.Sprintf("%s/app/projects/%d", server, p.ID))
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

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Project)(s).PopulateResourceWebLink(server)
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

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (m *Multiple) PopulateResourceWebLink(server string) {
	for i := range m.Response.Projects {
		m.Response.Projects[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new project in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/post-projects-json
type Create struct {
	Name        string               `json:"name"`
	Description *string              `json:"description,omitempty"`
	StartAt     *teamwork.LegacyDate `json:"start-date,omitempty"`
	EndAt       *teamwork.LegacyDate `json:"end-date,omitempty"`

	CompanyID int64   `json:"companyId"`
	OwnerID   *int64  `json:"projectOwnerId,omitempty"`
	Tags      []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new project.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects.json", server)
	payload := struct {
		Project Create `json:"project"`
	}{Project: c}
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

// Update represents the payload for updating an existing project in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/put-projects-id-json
type Update struct {
	ID          int64                `json:"-"`
	Name        *string              `json:"name,omitempty"`
	Description *string              `json:"description,omitempty"`
	StartAt     *teamwork.LegacyDate `json:"start-date,omitempty"`
	EndAt       *teamwork.LegacyDate `json:"end-date,omitempty"`

	CompanyID *int64  `json:"companyId"`
	OwnerID   *int64  `json:"projectOwnerId,omitempty"`
	Tags      []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update a project.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d.json", server, u.ID)
	payload := struct {
		Project Update `json:"project"`
	}{Project: u}
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

// Delete represents the payload for deleting a project in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/delete-projects-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a project.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
