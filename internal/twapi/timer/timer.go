// Package timer provides functionality to manage timers in Teamwork.com. It
// includes operations for retrieving, creating, and updating timer information.
// It is part of the Teamwork AI project, which integrates with Teamwork.com to
// provide AI-driven insights and operations.
package timer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

// Timer represents a timer in Teamwork.com.
type Timer struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Running     bool   `json:"running"`
	Billable    bool   `json:"billable"`

	User    twapi.Relationship  `json:"user"`
	Task    *twapi.Relationship `json:"task"`
	Project twapi.Relationship  `json:"project"`
	Timelog *twapi.Relationship `json:"timelog,omitempty"`

	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	Deleted        bool       `json:"deleted"`
	Duration       int64      `json:"duration"`
	LastStartedAt  time.Time  `json:"lastStartedAt"`
	LastIntervalAt *time.Time `json:"timerLastIntervalEnd,omitempty"`
	Intervals      []struct {
		ID       int64     `json:"id"`
		From     time.Time `json:"from"`
		To       time.Time `json:"to"`
		Duration int64     `json:"duration"`
	} `json:"intervals"`
}

// Single represents a request to retrieve a single timer by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/get-projects-api-v3-timers-timer-id-json
type Single Timer

// HTTPRequest creates an HTTP request to retrieve a single timer by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/timers/%d.json", server, s.ID)
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
		Timer Timer `json:"timer"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Timer)
	return nil
}

// Multiple represents a request to retrieve multiple timers.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/get-projects-api-v3-timers-json
type Multiple struct {
	Request struct {
		Filters struct {
			UserID            int64
			TaskID            int64
			ProjectID         int64
			RunningTimersOnly bool
			Page              int64
			PageSize          int64
		}
	}
	Response struct {
		Meta struct {
			Page struct {
				HasMore bool `json:"hasMore"`
			} `json:"page"`
		} `json:"meta"`
		Timers []Timer `json:"timers"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple timers.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/timers.json", server)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.UserID > 0 {
		query.Set("userId", strconv.FormatInt(m.Request.Filters.UserID, 10))
	}
	if m.Request.Filters.TaskID > 0 {
		query.Set("taskId", strconv.FormatInt(m.Request.Filters.TaskID, 10))
	}
	if m.Request.Filters.ProjectID > 0 {
		query.Set("projectId", strconv.FormatInt(m.Request.Filters.ProjectID, 10))
	}
	if m.Request.Filters.RunningTimersOnly {
		query.Set("runningTimersOnly", strconv.FormatBool(m.Request.Filters.RunningTimersOnly))
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

// Create represents the payload for creating a new timer in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/post-projects-api-v3-me-timers-json
type Create struct {
	Description       *string `json:"description"`
	Billable          *bool   `json:"isBillable"`
	Running           *bool   `json:"isRunning"`
	Seconds           *int64  `json:"seconds"`
	StopRunningTimers *bool   `json:"stopRunningTimers"`

	ProjectID *int64 `json:"projectId"`
	TaskID    *int64 `json:"taskId"`
}

// HTTPRequest creates an HTTP request to create a new timer.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers.json", server)
	payload := struct {
		Timer Create `json:"timer"`
	}{Timer: c}
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

// Update represents the payload for updating an existing timer in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/patch-projects-api-v3-time-timer-id-json
type Update struct {
	ID          int64   `json:"-"`
	Description *string `json:"description"`
	Billable    *bool   `json:"isBillable"`
	Running     *bool   `json:"isRunning"`

	ProjectID *int64 `json:"projectId"`
	TaskID    *int64 `json:"taskId"`
}

// HTTPRequest creates an HTTP request to update a timer.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers/%d.json", server, u.ID)
	payload := struct {
		Timer Update `json:"timer"`
	}{Timer: u}
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

// Pause represents the payload for pausing an existing timer in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/put-projects-api-v3-me-timers-timer-id-pause-json
type Pause struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to pause a timer.
func (p Pause) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers/%d/pause.json", server, p.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Complete represents the payload for completing an existing timer in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/put-projects-api-v3-me-timers-timer-id-complete-json
type Complete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to complete a timer.
func (c Complete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers/%d/complete.json", server, c.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Resume represents the payload for resuming an existing timer in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/put-projects-api-v3-me-timers-timer-id-resume-json
type Resume struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to resume a timer.
func (r Resume) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers/%d/resume.json", server, r.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Delete represents the payload for deleting an existing timer in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/time-tracking/delete-projects-api-v3-time-timer-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a timer.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me/timers/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
