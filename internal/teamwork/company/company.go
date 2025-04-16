package company

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

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

type SingleCompany Company

func (t SingleCompany) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/companies/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleCompany) UnmarshalJSON(data []byte) error {
	var raw struct {
		Company Company `json:"company"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleCompany(raw.Company)
	return nil
}

type MultipleCompanies []Company

func (t MultipleCompanies) Request(server string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, server+"/projects/api/v3/companies.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleCompanies) UnmarshalJSON(data []byte) error {
	var raw struct {
		Companies []Company `json:"companies"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Companies
	return nil
}

type CompanyCreation struct {
	Name string `json:"name"`
}

func (t CompanyCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/companies.json", server)
	paylaod := struct {
		Company CompanyCreation `json:"company"`
	}{Company: t}
	body, err := json.Marshal(paylaod)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
