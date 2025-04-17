package company

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Company represents a company or client in Teamwork.com. It contains
// information about the company such as its name, address, contact details, and
// relationships with other entities like industry, tags, and currency.
type Company struct {
	ID               int64                   `json:"id"`
	Name             string                  `json:"name"`
	CreatedAt        *time.Time              `json:"createdAt"`
	UpdatedAt        *time.Time              `json:"updatedAt"`
	AddressOne       string                  `json:"addressOne"`
	AddressTwo       string                  `json:"addressTwo"`
	City             string                  `json:"city"`
	State            string                  `json:"state"`
	Zip              string                  `json:"zip"`
	CountryCode      string                  `json:"countryCode"`
	EmailOne         string                  `json:"emailOne"`
	EmailTwo         string                  `json:"emailTwo"`
	EmailThree       string                  `json:"emailThree"`
	Website          string                  `json:"website"`
	CID              string                  `json:"cid"`
	Phone            string                  `json:"phone"`
	Fax              string                  `json:"fax"`
	ProfileText      *string                 `json:"profileText,omitempty"`
	PrivateNotesText *string                 `json:"privateNotesText,omitempty"`
	PrivateNotes     *string                 `json:"privateNotes,omitempty"`
	CanSeePrivate    bool                    `json:"canSeePrivate"`
	IsOwner          bool                    `json:"isOwner"`
	Industry         *teamwork.Relationship  `json:"industry"`
	NameURL          string                  `json:"companyNameUrl"`
	ClientManagedBy  *teamwork.Relationship  `json:"clientManagedBy"`
	Tags             []teamwork.Relationship `json:"tags"`
	Status           string                  `json:"status"`
	Currency         *teamwork.Relationship  `json:"currency"`
}

// Single represents a request to retrieve a single company by its ID.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/get-projects-api-v3-companies-company-id-json
type Single Company

// Request creates an HTTP request to retrieve a single company by its ID.
func (t Single) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies/%d.json", server, t.ID)
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
		Company Company `json:"company"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = Single(raw.Company)
	return nil
}

// Multiple represents a request to retrieve multiple companies.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/get-projects-api-v3-companies-json
type Multiple []Company

// Request creates an HTTP request to retrieve multiple companies.
func (t Multiple) Request(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/projects/api/v3/companies.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (t *Multiple) UnmarshalJSON(data []byte) error {
	var raw struct {
		Companies []Company `json:"companies"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Companies
	return nil
}

// Creation represents the payload for creating a new company in Teamwork.com.
//
// https://apidocs.teamwork.com/docs/teamwork/v3/companies/post-projects-api-v3-companies-json
type Creation struct {
	Name string `json:"name"`
}

// Request creates an HTTP request to create a new company.
func (t Creation) Request(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/companies.json", server)
	paylaod := struct {
		Company Creation `json:"company"`
	}{Company: t}
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
