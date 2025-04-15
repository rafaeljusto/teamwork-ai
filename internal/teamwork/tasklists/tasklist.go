package tasklists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

type Tasklist struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	DisplayOrder  int64                  `json:"displayOrder"`
	ProjectID     int64                  `json:"projectId"`
	Project       teamwork.Relationship  `json:"project"`
	MilestoneID   *int64                 `json:"milestoneId"`
	Milestone     *teamwork.Relationship `json:"milestone"`
	IsPinned      bool                   `json:"isPinned"`
	IsPrivate     bool                   `json:"isPrivate"`
	LockdownID    *int64                 `json:"lockdownId"`
	Status        string                 `json:"status"`
	DefaultTaskID *int64                 `json:"defaultTaskId"`
	DefaultTask   *teamwork.Relationship `json:"defaultTask"`
	IsBillable    *bool                  `json:"isBillable"`
	Budget        *teamwork.Relationship `json:"tasklistBudget"`
	CreatedAt     *time.Time             `json:"createdAt"`
	UpdatedAt     *time.Time             `json:"updatedAt"`
	Icon          *string                `json:"icon"`
}

type SingleTasklist Tasklist

func (t SingleTasklist) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleTasklist) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tasklist Tasklist `json:"tasklist"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleTasklist(raw.Tasklist)
	return nil
}

type MultipleTasklists []Tasklist

func (t MultipleTasklists) Request(server string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, server+"/projects/api/v3/tasklists.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleTasklists) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tasklists []Tasklist `json:"tasklists"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Tasklists
	return nil
}

type TasklistCreation struct {
	Name        string `json:"name"`
	ProjectID   int64  `json:"projectId"`
	Description string `json:"description"`
}

func (t TasklistCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/%d/tasklists.json", server, t.ProjectID)
	paylaod := struct {
		Tasklist TasklistCreation `json:"todo-list"`
	}{Tasklist: t}
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
