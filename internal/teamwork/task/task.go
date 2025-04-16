package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

type Task struct {
	ID                     int64                   `json:"id"`
	Name                   string                  `json:"name"`
	Status                 string                  `json:"status"`
	Description            *string                 `json:"description"`
	DescriptionContentType *string                 `json:"descriptionContentType"`
	Priority               *string                 `json:"priority"`
	Progress               int64                   `json:"progress"`
	Tasklist               teamwork.Relationship   `json:"tasklist"`
	Assignees              []teamwork.Relationship `json:"assignees"`
	Tags                   []teamwork.Relationship `json:"tags"`
	StartDate              *time.Time              `json:"startDate"`
	DueDate                *time.Time              `json:"dueDate"`
	EstimateMinutes        int64                   `json:"estimateMinutes"`
	CreatedBy              *int64                  `json:"createdBy"`
	CreatedAt              *time.Time              `json:"createdAt"`
	UpdatedBy              *int64                  `json:"updatedBy"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	DeletedBy              *int64                  `json:"deletedBy"`
	DeletedAt              *time.Time              `json:"deletedAt"`
	CompletedBy            *int64                  `json:"completedBy,omitempty"`
	CompletedDate          *time.Time              `json:"completedAt,omitempty"`
	Meta                   map[string]any          `json:"meta,omitempty"`
}

type SingleTask Task

func (t SingleTask) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleTask) UnmarshalJSON(data []byte) error {
	var raw struct {
		Task Task `json:"task"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleTask(raw.Task)
	return nil
}

type MultipleTasks struct {
	Tasks      []Task
	ProjectID  int64
	TasklistID int64
}

func (t MultipleTasks) Request(server string) (*http.Request, error) {
	var url string
	switch {
	case t.ProjectID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/projects/%d/tasks.json", server, t.ProjectID)
	case t.TasklistID > 0:
		url = fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, t.TasklistID)
	default:
		url = fmt.Sprintf("%s/projects/api/v3/tasks.json", server)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleTasks) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tasks []Task `json:"tasks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Tasks = raw.Tasks
	return nil
}

type TaskCreation struct {
	Name        string               `json:"name"`
	TasklistID  int64                `json:"tasklistId"`
	Description string               `json:"description"`
	Assignees   *teamwork.UserGroups `json:"assignees,omitempty"`
	Priority    *string              `json:"priority,omitempty"`
}

func (t TaskCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasklists/%d/tasks.json", server, t.TasklistID)
	paylaod := struct {
		Task TaskCreation `json:"task"`
	}{Task: t}
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

type TaskUpdate struct {
	ID   int64
	Task struct {
		Name        *string              `json:"name,omitempty"`
		Description *string              `json:"description,omitempty"`
		Assignees   *teamwork.UserGroups `json:"assignees"`
		Priority    *string              `json:"priority,omitempty"`
	}
}

func (t TaskUpdate) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/tasks/%d.json", server, t.ID)
	paylaod := struct {
		Task any `json:"task"`
	}{Task: t.Task}
	body, err := json.Marshal(paylaod)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPatch, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
