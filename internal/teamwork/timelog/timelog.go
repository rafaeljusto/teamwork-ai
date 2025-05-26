// Package timelog provides functionality to manage timelogs in Teamwork.com. It
// includes operations for retrieving, creating, and updating timelog
// information. It is part of the Teamwork AI project, which integrates with
// Teamwork.com to provide AI-driven insights and operations.
package timelog

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

// Timelog represents a timelog in Teamwork.com.
type Timelog struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Billable    bool      `json:"billable"`
	Minutes     int64     `json:"minutes"`
	LoggedAt    time.Time `json:"timeLogged"`

	User    teamwork.Relationship   `json:"user"`
	Task    *teamwork.Relationship  `json:"task"`
	Project teamwork.Relationship   `json:"project"`
	Tags    []teamwork.Relationship `json:"tags,omitempty"`

	CreatedAt time.Time  `json:"createdAt"`
	LoggedBy  int64      `json:"loggedBy"`
	UpdatedAt *time.Time `json:"updatedAt"`
	UpdatedBy *int64     `json:"updatedBy"`
	DeletedAt *time.Time `json:"deletedAt"`
	DeletedBy *int64     `json:"deletedBy"`
	Deleted   bool       `json:"deleted"`
}

// Single represents a request to retrieve a single timelog by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/get-projects-api-v3-time-timelog-id-json
type Single Timelog

// HTTPRequest creates an HTTP request to retrieve a single timelog by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/time/%d.json", server, s.ID)
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
		Timelog Timelog `json:"timelog"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Timelog)
	return nil
}

// Multiple represents a request to retrieve multiple timelogs.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/get-projects-api-v3-time-json
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/post-projects-api-v3-projects-project-id-time-json
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/get-projects-api-v3-tasks-task-id-time-json
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID int64
			TaskID    int64
		}
		Filters struct {
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
		Timelogs []Timelog `json:"timelogs"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple timelogs.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var uri string
	switch {
	case m.Request.Path.ProjectID > 0:
		uri = fmt.Sprintf("%s/projects/api/v3/projects/%d/time.json", server, m.Request.Path.ProjectID)
	case m.Request.Path.TaskID > 0:
		uri = fmt.Sprintf("%s/projects/api/v3/tasks/%d/time.json", server, m.Request.Path.TaskID)
	default:
		uri = fmt.Sprintf("%s/projects/api/v3/time.json", server)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
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

// Create represents the payload for creating a new timelog in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/post-projects-api-v3-tasks-task-id-time-json
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/post-projects-api-v3-projects-project-id-time-json
type Create struct {
	Description *string       `json:"description"`
	Date        teamwork.Date `json:"date"`
	Time        teamwork.Time `json:"time"`
	IsUTC       bool          `json:"isUTC"`
	Hours       int64         `json:"hours"`
	Minutes     int64         `json:"minutes"`
	Billable    bool          `json:"isBillable"`

	ProjectID int64   `json:"-"`
	TaskID    int64   `json:"-"`
	UserID    *int64  `json:"userId"`
	TagIDs    []int64 `json:"tagIds"`
}

// HTTPRequest creates an HTTP request to create a new timelog.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var uri string
	if c.TaskID > 0 {
		uri = fmt.Sprintf("%s/projects/api/v3/tasks/%d/time.json", server, c.TaskID)
	} else {
		uri = fmt.Sprintf("%s/projects/api/v3/projects/%d/time.json", server, c.ProjectID)
	}
	payload := struct {
		Timelog Create `json:"timelog"`
	}{Timelog: c}
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

// Update represents the payload for updating an existing timelog in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/patch-projects-api-v3-time-timelog-id-json
type Update struct {
	ID          int64          `json:"-"`
	Description *string        `json:"description,omitempty"`
	Date        *teamwork.Date `json:"date,omitempty"`
	Time        *teamwork.Time `json:"time,omitempty"`
	IsUTC       *bool          `json:"isUTC,omitempty"`
	Hours       *int64         `json:"hours,omitempty"`
	Minutes     *int64         `json:"minutes,omitempty"`
	Billable    *bool          `json:"isBillable,omitempty"`

	UserID *int64  `json:"userId,omitempty"`
	TagIDs []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update a timelog.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/time/%d.json", server, u.ID)
	payload := struct {
		Timelog Update `json:"timelog"`
	}{Timelog: u}
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

// Delete represents the payload for deleting an existing timelog in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/delete-projects-api-v3-time-timelog-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a timelog.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/time/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
