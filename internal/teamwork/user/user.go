package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// User represents a user in Teamwork.com. It contains information about the
// user such as their ID, name, email, company affiliation, admin status, client
// status, service account status, user type, deletion status, working hours,
// rate, cost, job roles, skills, placeholder status, timezone, creation and
// update details, and any additional metadata. Users can be administrators,
// clients, or service accounts, and they can have various roles and skills
// associated with them. This struct is used to manage user information within
// Teamwork.com, allowing teams to organize and manage their members
// effectively.
type User struct {
	ID            int64                   `json:"id"`
	FirstName     string                  `json:"firstName"`
	LastName      string                  `json:"lastName"`
	Title         *string                 `json:"title"`
	Email         string                  `json:"email"`
	Company       teamwork.Relationship   `json:"company"`
	IsAdmin       bool                    `json:"isAdmin"`
	IsClient      bool                    `json:"isClientUser"`
	IsService     bool                    `json:"isServiceAccount"`
	Type          string                  `json:"type"`
	Deleted       bool                    `json:"deleted"`
	WorkingHour   *teamwork.Relationship  `json:"workingHour"`
	Rate          *int64                  `json:"userRate,omitempty"`
	Cost          *int64                  `json:"userCost,omitempty"`
	JobRoles      []teamwork.Relationship `json:"jobRoles,omitempty"`
	Skills        []teamwork.Relationship `json:"skills,omitempty"`
	IsPlaceholder bool                    `json:"isPlaceholderResource"`
	Timezone      *string                 `json:"timezone,omitempty"`
	CreatedBy     *teamwork.Relationship  `json:"createdBy"`
	CreatedAt     time.Time               `json:"createdAt"`
	UpdatedBy     *teamwork.Relationship  `json:"updatedBy"`
	UpdatedAt     *time.Time              `json:"updatedAt"`
	Meta          map[string]any          `json:"meta,omitempty"`
}

// Single represents a request to retrieve a single user by their ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/person/get-projects-api-v3-people-person-id-json
type Single User

// Request creates an HTTP request to retrieve a single user by their ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/people/%d.json", server, t.ID)
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
		User User `json:"person"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.User)
	return nil
}

// Multiple represents a request to retrieve multiple users.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/people/get-projects-api-v3-people-json
// https://apidocs.teamwork.com/docs/teamwork/v3/people/get-projects-api-v3-projects-project-id-people-json
type Multiple struct {
	Users     []User
	ProjectID int64
}

// Request creates an HTTP request to retrieve multiple users.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case t.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/people.json", server, t.ProjectID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/people.json", server)
	}
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
		Users []User `json:"people"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Users = raw.Users
	return nil
}

// Creation represents the payload for creating a new user in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/people/post-people-json
type Creation struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email-address"`
	Password  string `json:"password"`
}

// Request creates an HTTP request to create a new user.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/people.json", server)
	paylaod := struct {
		User Creation `json:"person"`
	}{User: t}
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
