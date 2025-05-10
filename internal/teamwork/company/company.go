// Package company provides functionality to manage companies in Teamwork.com.
// It includes operations for retrieving, creating, and updating company
// information. It defines structures for representing company data, including
// details like address, contact information, and relationships with industries
// and tags. It is part of the Teamwork AI project, which integrates with
// Teamwork.com to provide AI-driven insights and operations.
package company

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

// Company represents a company or client in Teamwork.com. It contains
// information about the company such as its name, address, contact details, and
// relationships with other entities like industry and tags.
type Company struct {
	ID          int64   `json:"id"`
	AddressOne  string  `json:"addressOne"`
	AddressTwo  string  `json:"addressTwo"`
	City        string  `json:"city"`
	CountryCode string  `json:"countryCode"`
	EmailOne    string  `json:"emailOne"`
	EmailThree  string  `json:"emailThree"`
	EmailTwo    string  `json:"emailTwo"`
	Fax         string  `json:"fax"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Profile     *string `json:"profileText"`
	State       string  `json:"state"`
	Website     string  `json:"website"`
	Zip         string  `json:"zip"`

	ManagedBy *teamwork.Relationship  `json:"clientManagedBy"`
	Industry  *teamwork.Relationship  `json:"industry"`
	Tags      []teamwork.Relationship `json:"tags"`

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Status    string     `json:"status"`
	WebLink   *string    `json:"webLink,omitempty"`
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (c *Company) PopulateResourceWebLink(server string) {
	if c.ID == 0 {
		return
	}
	c.WebLink = teamwork.Ref(fmt.Sprintf("%s/app/clients/%d", server, c.ID))
}

// Single represents a request to retrieve a single company by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/get-projects-api-v3-companies-company-id-json
type Single Company

// HTTPRequest creates an HTTP request to retrieve a single company by its ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies/%d.json", server, s.ID)
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
		Company Company `json:"company"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Single(raw.Company)
	return nil
}

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (s *Single) PopulateResourceWebLink(server string) {
	(*Company)(s).PopulateResourceWebLink(server)
}

// Multiple represents a request to retrieve multiple companies.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/get-projects-api-v3-companies-json
type Multiple struct {
	Request struct {
		Filters struct {
			SearchTerm   string
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
		Companies []Company `json:"companies"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple companies.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/projects/api/v3/companies.json", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if m.Request.Filters.SearchTerm != "" {
		query.Set("searchTerm", m.Request.Filters.SearchTerm)
	}
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

// PopulateResourceWebLink sets the website URL for the specific resource. It
// should be called after the object is loaded (the ID is set).
func (m *Multiple) PopulateResourceWebLink(server string) {
	for i := range m.Response.Companies {
		m.Response.Companies[i].PopulateResourceWebLink(server)
	}
}

// Create represents the payload for creating a new company in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/post-projects-api-v3-companies-json
type Create struct {
	AddressOne  *string `json:"addressOne,omitempty"`
	AddressTwo  *string `json:"addressTwo,omitempty"`
	City        *string `json:"city,omitempty"`
	CountryCode *string `json:"countrycode,omitempty"`
	EmailOne    *string `json:"emailOne,omitempty"`
	EmailTwo    *string `json:"emailTwo,omitempty"`
	EmailThree  *string `json:"emailThree,omitempty"`
	Fax         *string `json:"fax,omitempty"`
	Name        string  `json:"name"`
	Phone       *string `json:"phone,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	State       *string `json:"state,omitempty"`
	Website     *string `json:"website,omitempty"`
	Zip         *string `json:"zip,omitempty"`

	ManagerID  *int64  `json:"clientManagedBy"`
	IndustryID *int64  `json:"industryCatId,omitempty"`
	TagIDs     []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to create a new company.
func (c Create) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies.json", server)
	payload := struct {
		Company Create `json:"company"`
	}{Company: c}
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

// Update represents the payload for updating an existing company in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/patch-projects-api-v3-companies-company-id-json
type Update struct {
	ID          int64   `json:"-"`
	AddressOne  *string `json:"addressOne,omitempty"`
	AddressTwo  *string `json:"addressTwo,omitempty"`
	City        *string `json:"city,omitempty"`
	CountryCode *string `json:"countrycode,omitempty"`
	EmailOne    *string `json:"emailOne,omitempty"`
	EmailTwo    *string `json:"emailTwo,omitempty"`
	EmailThree  *string `json:"emailThree,omitempty"`
	Fax         *string `json:"fax,omitempty"`
	Name        *string `json:"name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	State       *string `json:"state,omitempty"`
	Website     *string `json:"website,omitempty"`
	Zip         *string `json:"zip,omitempty"`

	ManagerID  *int64  `json:"clientManagedBy"`
	IndustryID *int64  `json:"industryCatId,omitempty"`
	TagIDs     []int64 `json:"tagIds,omitempty"`
}

// HTTPRequest creates an HTTP request to update a company.
func (u Update) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies/%d.json", server, u.ID)
	payload := struct {
		Company Update `json:"company"`
	}{Company: u}
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

// Delete represents the payload for deleting an existing company in
// Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/delete-projects-api-v3-companies-company-id-json
type Delete struct {
	Request struct {
		Path struct {
			ID int64 `json:"-"`
		}
	}
}

// HTTPRequest creates an HTTP request to update a milestone.
func (d Delete) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies/%d.json", server, d.Request.Path.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
