package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddProject represents a request to add users to a project in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/people/put-projects-api-v3-projects-project-id-people-json
type AddProject struct {
	Request struct {
		Path struct {
			ProjectID int64
		}
		Users struct {
			IDs []int64 `json:"userIds"`
		}
	}
}

// HTTPRequest creates an HTTP request to add users to a project.
func (a AddProject) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/projects/%d/people.json", server, a.Request.Path.ProjectID)
	body, err := json.Marshal(a.Request.Users)
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
