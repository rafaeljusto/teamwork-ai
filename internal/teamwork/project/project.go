package project

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Project represents a project in Teamwork.com. It contains information about
// the project such as its ID, name, description, type, status, associated
// company, tags, start and end dates, category, ownership, updates, creation
// and update timestamps, completion details, project owner, custom field
// values, and whether it is starred or billable. It is used to manage and
// organize work, providing a central hub for all components related to what the
// team is working on.
type Project struct {
	ID                int64                   `json:"id"`
	Name              string                  `json:"name"`
	Description       *string                 `json:"description"`
	Type              string                  `json:"type"`
	Status            string                  `json:"status"`
	Company           teamwork.Relationship   `json:"company"`
	Tags              []teamwork.Relationship `json:"tags"`
	StartAt           *time.Time              `json:"startAt"`
	EndAt             *time.Time              `json:"endAt"`
	Category          *teamwork.Relationship  `json:"category"`
	OwnedBy           *int64                  `json:"ownedBy"`
	Update            *teamwork.Relationship  `json:"update"`
	CreatedBy         *int64                  `json:"createdBy"`
	CreatedAt         *time.Time              `json:"createdAt"`
	UpdatedAt         *time.Time              `json:"updatedAt"`
	UpdatedBy         *int64                  `json:"updatedBy"`
	CompletedAt       *time.Time              `json:"completedAt"`
	CompletedBy       *int64                  `json:"completedBy"`
	ProjectOwner      *teamwork.Relationship  `json:"projectOwner"`
	CustomFieldValues []teamwork.Relationship `json:"customfieldValues"`
	IsStarred         *bool                   `json:"isStarred,omitempty"`
	IsBillable        bool                    `json:"isBillable"`
}

// Single represents a request to retrieve a single project by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/projects/get-projects-api-v3-projects-project-id-json
type Single Project

// Request creates an HTTP request to retrieve a single project by its ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/projects/%d.json", server, t.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Single instance.
func (t *Single) UnmarshalJSON(data []byte) error {
	var raw struct {
		Project Project `json:"project"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.Project)
	return nil
}

// Multiple represents a request to retrieve multiple projects.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/projects/get-projects-api-v3-projects-json
type Multiple []Project

// Request creates an HTTP request to retrieve multiple projects.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/projects/api/v3/projects.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (t *Multiple) UnmarshalJSON(data []byte) error {
	var raw struct {
		Projects []Project `json:"projects"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Projects
	return nil
}

// Creation represents the payload for creating a new project in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/projects/post-projects-json
type Creation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Request creates an HTTP request to create a new project in Teamwork.com.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects.json", server)
	paylaod := struct {
		Project Creation `json:"project"`
	}{Project: t}
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
