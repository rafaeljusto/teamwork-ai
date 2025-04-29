// Package tag provides functionality to manage tags in Teamwork.com. Tags are
// used to mark items for filtering and organization across various resources
// such as projects, tasks, milestones, messages, time logs, notebooks, files,
// and links. It includes operations for retrieving, creating, and updating
// tags, as well as handling HTTP requests to the Teamwork.com API for tag
// management.
package tag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Tag are a way to mark items so that you can use a filter to see just those
// items. Tags can be added to projects, tasks, milestones, messages, time logs,
// notebooks, files and links.
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Project *teamwork.Relationship `json:"project"`

	WebLink *string `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (t *Tag) PopulateResourceWebLink(server string) {
	if t.ID == 0 {
		return
	}
	t.WebLink = teamwork.Ref(fmt.Sprintf("https://%s/app/settings/tags", server))
}

// Single represents a request to retrieve a single tag by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tags/get-projects-api-v3-tags-tag-id-json
type Single Tag

// HTTPRequest creates an HTTP request to retrieve a single tag by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tags/%d.json", server, s.ID)
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
		Tag Tag `json:"tag"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Tag)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Tag)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple tags.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tags/get-projects-api-v3-tags-json
type Multiple struct {
	Request struct {
		Filters struct {
			SearchTerm string
			ItemType   string
			ProjectIDs []int64
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
		Tags []Tag `json:"tags"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple tags.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/projects/api/v3/tags.json", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
	}
	if m.Request.Filters.ItemType != "" {
		query.Set("itemType", m.Request.Filters.ItemType)
	}
	if len(m.Request.Filters.ProjectIDs) > 0 {
		projectIDs := make([]string, len(m.Request.Filters.ProjectIDs))
		for i, id := range m.Request.Filters.ProjectIDs {
			projectIDs[i] = strconv.FormatInt(id, 10)
		}
		query.Set("projectIds", strings.Join(projectIDs, ","))
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
	for i := range m.Response.Tags {
		m.Response.Tags[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new tag in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tags/post-projects-api-v3-tags-json
type Create struct {
	Name string `json:"name"`

	ProjectID *int64 `json:"projectId,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new tag.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tags.json", server)
	payload := struct {
		Tag Create `json:"tag"`
	}{Tag: c}
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

// Update represents the payload for updating an existing tag in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/tags/patch-projects-api-v3-tags-tag-id-json
type Update struct {
	ID   int64   `json:"-"`
	Name *string `json:"name,omitempty"`

	ProjectID *int64 `json:"projectId"`
}

// HTTPRequest creates an HTTP request to update a tag.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tags/%d.json", server, u.ID)
	payload := struct {
		Tag Update `json:"tag"`
	}{Tag: u}
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

// Delete represents the payload for deleting a tag in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/tags/delete-tags-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a tag by its ID.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tags/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
