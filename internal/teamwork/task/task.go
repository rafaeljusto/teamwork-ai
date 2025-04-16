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

type MultipleTasks []Task

func (t MultipleTasks) Request(server string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, server+"/projects/api/v3/tasks.json", nil)
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
	*t = raw.Tasks
	return nil
}

type TaskCreation struct {
	Name        string `json:"name"`
	TasklistID  int64  `json:"tasklistId"`
	Description string `json:"description"`
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
