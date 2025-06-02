package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AssignJobRole represents a request to assign users to a specific job role in
// Teamwork.com.
//
// No public documentation available yet.
type AssignJobRole struct {
	Request struct {
		Path struct {
			JobRoleID int64
		}
		Users struct {
			IDs       []int64 `json:"users"`
			IsPrimary bool    `json:"isPrimary"`
		}
	}
}

// HTTPRequest creates an HTTP request to assign users to a job role.
func (a AssignJobRole) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles/%d/people.json", server, a.Request.Path.JobRoleID)
	body, err := json.Marshal(a.Request.Users)
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

// UnassignJobRole represents a request to unassign users from a job role in
// Teamwork.com.
//
// No public documentation available yet.
type UnassignJobRole struct {
	Request struct {
		Path struct {
			JobRoleID int64
		}
		Users struct {
			IDs       []int64 `json:"users"`
			IsPrimary bool    `json:"isPrimary"`
		}
	}
}

// HTTPRequest creates an HTTP request to unassign users from a job role.
func (a UnassignJobRole) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/jobroles/%d/people.json", server, a.Request.Path.JobRoleID)
	body, err := json.Marshal(a.Request.Users)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
