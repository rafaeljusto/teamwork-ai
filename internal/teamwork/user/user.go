package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

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

type SingleUser User

func (t SingleUser) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/people/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleUser) UnmarshalJSON(data []byte) error {
	var raw struct {
		User User `json:"person"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleUser(raw.User)
	return nil
}

type MultipleUsers struct {
	Users     []User
	ProjectID int64
}

func (t MultipleUsers) Request(server string) (*http.Request, error) {
	var url string
	switch {
	case t.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/people.json", server, t.ProjectID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/people.json", server)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleUsers) UnmarshalJSON(data []byte) error {
	var raw struct {
		Users []User `json:"people"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Users = raw.Users
	return nil
}

type UserCreation struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email-address"`
	Password  string `json:"password"`
}

func (t UserCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/people.json", server)
	paylaod := struct {
		User UserCreation `json:"person"`
	}{User: t}
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
