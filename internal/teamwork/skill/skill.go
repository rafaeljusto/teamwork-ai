// Package skill provides functionality to manage skills in Teamwork.com. It
// allows for the retrieval, creation, and updating of skills, which are
// knowledge or abilities that can be assigned to users. Skills help in better
// task management and organization within projects.
package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Skill represents a skill in Teamwork.com. It contains information about the
// skill such as its ID, name, creation and update timestamps, and the users who
// created, updated, or deleted the skill. Skills are knowledge or abilities
// that can be assigned to users, allowing for better task management and
// organization within projects.
type Skill struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	CreatedByUserID int64      `json:"createdByUser"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedByUserID *int64     `json:"updatedByUser"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedByUserID *int64     `json:"deletedByUser"`
	DeletedAt       *time.Time `json:"deletedAt"`
}

// Single represents a request to retrieve a single skill by its ID.
//
// No public documentation available yet.
type Single Skill

// HTTPRequest creates an HTTP request to retrieve a single skill by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, s.ID)
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
		Skill Skill `json:"skill"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Skill)
	return nil
}

// Multiple represents a request to retrieve multiple skills.
//
// No public documentation available yet.
type Multiple struct {
	Request struct {
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
		Skills []Skill `json:"skills"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple skills.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	url := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
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

// Create represents the payload for creating a new skill in Teamwork.com.
//
// No public documentation available yet.
type Create struct {
	Name    string  `json:"name"`
	UserIDs []int64 `json:"userIds"`
}

// HTTPRequest creates an HTTP request to create a new skill in Teamwork.com.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
	payload := struct {
		Skill Create `json:"skill"`
	}{Skill: c}
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

// Update represents the payload for updating an existing skill in Teamwork.com.
//
// No public documentation available yet.
type Update struct {
	ID      int64   `json:"-"`
	Name    *string `json:"name,omitempty"`
	UserIDs []int64 `json:"userIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update an existing skill in
// Teamwork.com.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, u.ID)
	payload := struct {
		Skill Update `json:"skill"`
	}{Skill: u}
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

// Delete represents the payload for deleting an existing skill in
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

// HTTPRequest creates an HTTP request to update a milestone.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
