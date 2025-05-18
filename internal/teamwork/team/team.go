// Package team provides functionality to manage teams in Teamwork.com. It
// includes operations for retrieving, creating, and updating team information.
// Teams are used to group people based on their position or contribution within
// an organization.
package team

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Team represents a team in Teamwork.com. It contains information about the
// team such as its name, description, handle, logo, and associated members.
type Team struct {
	ID          teamwork.LegacyNumber `json:"id"`
	Name        string                `json:"name"`
	Description *string               `json:"description"`
	Handle      string                `json:"handle"`
	LogoURL     *string               `json:"logoUrl"`
	LogoIcon    *string               `json:"logoIcon"`
	LogoColor   *string               `json:"logoColor"`

	ProjectID teamwork.LegacyNumber `json:"projectId"`
	Company   *struct {
		ID   teamwork.LegacyNumber `json:"id"`
		Name string                `json:"name"`
	} `json:"company"`
	ParentTeam *struct {
		ID     teamwork.LegacyNumber `json:"id"`
		Name   string                `json:"name"`
		Handle string                `json:"handle"`
	} `json:"parentTeam"`
	RootTeam *struct {
		ID     teamwork.LegacyNumber `json:"id"`
		Name   string                `json:"name"`
		Handle string                `json:"handle"`
	} `json:"rootTeam"`
	Members []teamwork.LegacyRelationship `json:"members"`

	CreatedBy teamwork.LegacyNumber      `json:"createdByUserId"`
	CreatedAt time.Time                  `json:"dateCreated"`
	UpdatedBy teamwork.LegacyNumber      `json:"updatedByUserId"`
	UpdatedAt time.Time                  `json:"dateUpdated"`
	Deleted   bool                       `json:"deleted"`
	DeletedAt *teamwork.OptionalDateTime `json:"deletedDate"`
	WebLink   *string                    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (t *Team) PopulateResourceWebLink(server string) {
	if t.ID == 0 {
		return
	}
	t.WebLink = teamwork.Ref(fmt.Sprintf("%s/app/teams/%d", server, t.ID))
}

// Single represents a request to retrieve a single team by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/teams/get-teams-id-json
type Single Team

// HTTPRequest creates an HTTP request to retrieve a single team by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/teams/%d.json", server, s.ID)
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
		Team Team `json:"team"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Team)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Team)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple teams.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/teams/get-teams-json
type Multiple struct {
	Request struct {
		Filters struct {
			SearchTerm string
			Page       int64
			PageSize   int64
		}
	}
	Response struct {
		Teams []Team `json:"teams"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple teams.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/teams.json", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
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
	for i := range m.Response.Teams {
		m.Response.Teams[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new team in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/teams/post-teams-json
type Create struct {
	Name        string  `json:"name"`
	Handle      *string `json:"handle,omitempty"`
	Description *string `json:"description,omitempty"`

	ParentTeamID *int64                     `json:"parentTeamId,omitempty"`
	CompanyID    *int64                     `json:"companyId,omitempty"`
	ProjectID    *int64                     `json:"projectId,omitempty"`
	UserIDs      teamwork.LegacyNumericList `json:"userIds,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new team.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/teams.json", server)
	payload := struct {
		Team Create `json:"team"`
	}{Team: c}
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

// Update represents the payload for updating an existing team in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/teams/put-teams-id-json
type Update struct {
	ID          int64   `json:"-"`
	Name        *string `json:"name,omitempty"`
	Handle      *string `json:"handle,omitempty"`
	Description *string `json:"description,omitempty"`

	CompanyID *int64                     `json:"companyId,omitempty"`
	ProjectID *int64                     `json:"projectId,omitempty"`
	UserIDs   teamwork.LegacyNumericList `json:"userIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update a team.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/teams/%d.json", server, u.ID)
	payload := struct {
		Team Update `json:"team"`
	}{Team: u}
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

// Delete represents the payload for deleting an existing team in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v1/teams/delete-teams-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to delete a team.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/teams/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}
