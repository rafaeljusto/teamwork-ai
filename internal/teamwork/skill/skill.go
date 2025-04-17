package skill

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Skill represents a skill in Teamwork.com. It contains information about the
// skill such as its ID, name, creation and update timestamps, and the users who
// created, updated, or deleted the skill. Skills are knowledge or abilities
// that can be assigned to users, allowing for better task management and
// organization within projects.
type Skill struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
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

// Request creates an HTTP request to retrieve a single skill by its ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, t.ID)
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
		Skill Skill `json:"skill"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.Skill)
	return nil
}

// Multiple represents a request to retrieve multiple skills.
type Multiple []Skill

// Request creates an HTTP request to retrieve multiple skills.
//
// No public documentation available yet.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	url := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (t *Multiple) UnmarshalJSON(data []byte) error {
	var raw struct {
		Skills []Skill `json:"skills"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Skills
	return nil
}

// Creation represents the payload for creating a new skill in Teamwork.com.
//
// No public documentation available yet.
type Creation struct {
	Name    string  `json:"name"`
	UserIDs []int64 `json:"userIds"`
}

// Request creates an HTTP request to create a new skill in Teamwork.com.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
	paylaod := struct {
		Skill Creation `json:"skill"`
	}{Skill: t}
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

// Update represents the payload for updating an existing skill in Teamwork.com.
//
// No public documentation available yet.
type Update struct {
	ID    int64
	Skill struct {
		Name    *string `json:"name,omitempty"`
		UserIDs []int64 `json:"userIds,omitempty"`
	}
}

// Request creates an HTTP request to update an existing skill in Teamwork.com.
func (t Update) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, t.ID)
	paylaod := struct {
		Skill any `json:"skill"`
	}{Skill: t.Skill}
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
