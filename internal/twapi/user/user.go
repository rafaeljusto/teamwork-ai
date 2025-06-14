package user

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

// User represents a user in Teamwork.com. It contains information about the
// user such as their ID, name, email, company affiliation, admin status, client
// status, service account status, user type, deletion status, working hours,
// rate, cost, job roles, skills, placeholder status, timezone, creation and
// update details. Users can be administrators, clients, or service accounts,
// and they can have various roles and skills associated with them. This struct
// is used to manage user information within Teamwork.com, allowing teams to
// organize and manage their members effectively.
type User struct {
	ID        int64        `json:"id"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Title     *string      `json:"title"`
	Email     string       `json:"email"`
	Admin     bool         `json:"isAdmin"`
	Type      string       `json:"type"`
	Cost      *twapi.Money `json:"userCost"`
	Rate      *twapi.Money `json:"userRate"`

	Company  twapi.Relationship   `json:"company"`
	JobRoles []twapi.Relationship `json:"jobRoles,omitempty"`
	Skills   []twapi.Relationship `json:"skills,omitempty"`

	Deleted   bool                `json:"deleted"`
	CreatedBy *twapi.Relationship `json:"createdBy"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedBy *twapi.Relationship `json:"updatedBy"`
	UpdatedAt *time.Time          `json:"updatedAt"`
	WebLink   *string             `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (u *User) PopulateResourceWebLink(server string) {
	if u.ID == 0 {
		return
	}
	u.WebLink = twapi.Ref(fmt.Sprintf("%s/app/people/%d", server, u.ID))
}

// Single represents a request to retrieve a single user by their ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/person/get-projects-api-v3-people-person-id-json
type Single User

// HTTPRequest creates an HTTP request to retrieve a single user by their ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/people/%d.json", server, s.ID)
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
		User User `json:"person"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.User)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*User)(s).PopulateResourceWebLink(server)
}

// Me represents a request to retrieve the logged user information.
//
// No public documentation available yet.
type Me User

// HTTPRequest creates an HTTP request to retrieve the logged user information.
func (m Me) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/me.json", server)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Me instance.
func (m *Me) UnmarshalJSON(data []byte) error {
	var raw struct {
		User User `json:"person"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*m = Me(raw.User)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (m *Me) PopulateResourceWebLink(server string) {
	(*User)(m).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple users.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/people/get-projects-api-v3-people-json
// https://apidocs.teamwork.com/docs/teamwork/v3/people/get-projects-api-v3-projects-project-id-people-json
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID int64
		}
		Filters struct {
			SearchTerm string
			Type       string
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
		Users []User `json:"people"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple users.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var url string
	switch {
	case m.Request.Path.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/people.json", server, m.Request.Path.ProjectID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/people.json", server)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
	}
	if m.Request.Filters.Type != "" {
		query.Set("userType", m.Request.Filters.Type)
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
	for i := range m.Response.Users {
		m.Response.Users[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new user in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/people/post-people-json
type Create struct {
	FirstName string  `json:"first-name"`
	LastName  string  `json:"last-name"`
	Title     *string `json:"title,omitempty"`
	Email     string  `json:"email-address"`
	Admin     *bool   `json:"administrator,omitempty"`
	Type      *string `json:"user-type,omitempty"`

	CompanyID *int64 `json:"company-id,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new user.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/people.json", server)
	payload := struct {
		User Create `json:"person"`
	}{User: c}
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

// Update represents the payload for updating an existing user in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/people/put-people-id-json
type Update struct {
	ID        int64   `json:"-"`
	FirstName *string `json:"first-name,omitempty"`
	LastName  *string `json:"last-name,omitempty"`
	Title     *string `json:"title,omitempty"`
	Email     *string `json:"email-address,omitempty"`
	Password  *string `json:"password,omitempty"`
	Admin     *bool   `json:"administrator,omitempty"`
	Type      *string `json:"user-type,omitempty"`

	CompanyID *int64 `json:"company-id,omitempty"`
}

// HTTPRequest creates an HTTP request to update a new user.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/people/%d.json", server, u.ID)
	payload := struct {
		User Update `json:"person"`
	}{User: u}
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

// Delete represents the payload for deleting a user in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/people/delete-people-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a user.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/people/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
