// Package comment implements the API layer for managing comments in
// Teamwork.com. It provides structures and methods for creating, updating,
// retrieving, and listing comments.
package comment

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

// Comment represents a comment in Teamwork.com. It contains information about
// the comment's content, the user who posted it, and the associated project and
// object. The object can be a task, file, milestone, or notebook. The comment
// also includes metadata such as the date it was posted, edited, or deleted,
// and the user who performed those actions.
//
// More information can be found at:
// https://support.teamwork.com/projects/getting-started/comments-overview
type Comment struct {
	ID          int64  `json:"id"`
	Body        string `json:"body"`
	HTMLBody    string `json:"htmlBody"`
	ContentType string `json:"contentType"`

	Object  *twapi.Relationship `json:"object"`
	Project twapi.Relationship  `json:"project"`

	PostedBy     *int64     `json:"postedBy"`
	PostedAt     *time.Time `json:"postedDateTime"`
	LastEditedBy *int64     `json:"lastEditedBy"`
	EditedAt     *time.Time `json:"dateLastEdited"`
	Deleted      bool       `json:"deleted"`
	DeletedBy    *int64     `json:"deletedBy"`
	DeletedAt    *time.Time `json:"dateDeleted"`
	WebLink      *string    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (c *Comment) PopulateResourceWebLink(server string) {
	if c.Object == nil || c.ID == 0 {
		return
	}
	c.WebLink = twapi.Ref(fmt.Sprintf("%s/#%s/%d?c=%d", server, c.Object.Type, c.Object.ID, c.ID))
}

// Single represents a request to retrieve a single comment by its ID.
//
// No public documentation available yet.
type Single Comment

// HTTPRequest creates an HTTP request to retrieve a single comment by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/comments/%d.json", server, s.ID)
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
		Comment Comment `json:"comments"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Comment)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Comment)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple comments.
//
// No public documentation available yet.
type Multiple struct {
	Request struct {
		Path struct {
			FileID        int64
			FileVersionID int64
			MilestoneID   int64
			NotebookID    int64
			TaskID        int64
		}
		Filters struct {
			SearchTerm string
			UserIDs    []int64
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
		Comments []Comment `json:"comments"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple comments.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case m.Request.Path.FileID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/files/%d/comments.json", server, m.Request.Path.FileID)
	case m.Request.Path.FileVersionID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/fileversions/%d/comments.json", server, m.Request.Path.FileVersionID)
	case m.Request.Path.MilestoneID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/milestones/%d/comments.json", server, m.Request.Path.MilestoneID)
	case m.Request.Path.NotebookID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/notebooks/%d/comments.json", server, m.Request.Path.NotebookID)
	case m.Request.Path.TaskID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/tasks/%d/comments.json", server, m.Request.Path.TaskID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/comments.json", server)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
	}
	if len(m.Request.Filters.UserIDs) > 0 {
		userIDs := make([]string, len(m.Request.Filters.UserIDs))
		for i, id := range m.Request.Filters.UserIDs {
			userIDs[i] = strconv.FormatInt(id, 10)
		}
		query.Set("userIds", strings.Join(userIDs, ","))
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
	for i := range m.Response.Comments {
		m.Response.Comments[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new comment in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/comments/post-resource-resource-id-comments-json
type Create struct {
	Object      twapi.Relationship `json:"-"`
	Body        string             `json:"body"`
	ContentType *string            `json:"contentType,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new comment in a specific
// entity.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/%s/%d/comments.json", server, c.Object.Type, c.Object.ID)
	payload := struct {
		Comment Create `json:"comment"`
	}{Comment: c}
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

// Update represents the payload for updating an existing comment in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/comments/put-comments-id-json
type Update struct {
	ID          int64   `json:"-"`
	Body        string  `json:"body"`
	ContentType *string `json:"content-type,omitempty"`
}

// HTTPRequest creates an HTTP request to update an existing comment in
// Teamwork.com.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/comments/%d.json", server, u.ID)
	payload := struct {
		Comment Update `json:"comment"`
	}{Comment: u}
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

// Delete represents the payload for deleting an existing comment in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/comments/delete-comments-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a comment.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/comments/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
