package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

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

type SingleProject Project

func (t SingleProject) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/projects/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleProject) UnmarshalJSON(data []byte) error {
	var raw struct {
		Project Project `json:"project"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleProject(raw.Project)
	return nil
}

type MultipleProjects []Project

func (t MultipleProjects) Request(server string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, server+"/projects/api/v3/projects.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleProjects) UnmarshalJSON(data []byte) error {
	var raw struct {
		Projects []Project `json:"projects"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Projects
	return nil
}

type ProjectCreation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t ProjectCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects.json", server)
	paylaod := struct {
		Project ProjectCreation `json:"project"`
	}{Project: t}
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
