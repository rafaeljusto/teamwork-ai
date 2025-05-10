// Package jobrole provides functionality to manage job roles in Teamwork.com.
// It includes operations for retrieving, creating, and updating a job role. It
// is part of the Teamwork AI project, which integrates with Teamwork.com to
// provide AI-driven insights and operations.
package jobrole

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

// JobRole represents a job role in Teamwork.com.
type JobRole struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Users        []teamwork.Relationship `json:"users"`
	PrimaryUsers []teamwork.Relationship `json:"primaryUsers"`

	CreatedByUserID int64      `json:"createdByUser"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedByUserID *int64     `json:"updatedByUser"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedByUserID *int64     `json:"deletedByUser"`
	DeletedAt       *time.Time `json:"deletedAt"`
	IsActive        bool       `json:"isActive"`
	WebLink         *string    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (j *JobRole) PopulateResourceWebLink(server string) {
	if j.ID == 0 {
		return
	}
	j.WebLink = teamwork.Ref(fmt.Sprintf("%s/people/roles", server))
}

// Single represents a request to retrieve a single job role by its ID.
//
// No public documentation available yet.
type Single JobRole

// HTTPRequest creates an HTTP request to retrieve a single job role by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles/%d.json", server, s.ID)
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
		JobRole JobRole `json:"jobRole"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.JobRole)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*JobRole)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple job roles.
//
// No public documentation available yet.
type Multiple struct {
	Request struct {
		Filters struct {
			SearchTerm string
			Page       int64
			PageSize   int64
			Include    []string
		}
	}
	Response struct {
		Meta struct {
			Page struct {
				HasMore bool `json:"hasMore"`
			} `json:"page"`
		} `json:"meta"`
		JobRoles []JobRole `json:"jobRoles"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple job roles.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles.json", server)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
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
	if len(m.Request.Filters.Include) > 0 {
		query.Set("include", strings.Join(m.Request.Filters.Include, ","))
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
	for i := range m.Response.JobRoles {
		m.Response.JobRoles[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new job role in Teamwork.com.
//
// No public documentation available yet.
type Create struct {
	Name string `json:"name"`
}

// HTTPRequest creates an HTTP request to create a new job role.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles.json", server)
	payload := struct {
		JobRole Create `json:"jobRole"`
	}{JobRole: c}
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

// Update represents the payload for updating an existing job role in
// Teamwork.com.
//
// No public documentation available yet.
type Update struct {
	ID   int64  `json:"-"`
	Name string `json:"name"`
}

// HTTPRequest creates an HTTP request to update a job role.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles/%d.json", server, u.ID)
	payload := struct {
		JobRole Update `json:"jobrole"`
	}{JobRole: u}
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

// Delete represents the payload for deleting an existing job role in
// Teamwork.com.
//
// No public documentation available yet.
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to update a job role.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
