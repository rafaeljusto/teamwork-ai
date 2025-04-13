package tasks

import (
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
	uri := fmt.Sprintf(server+"/projects/api/v3/tasks/%d.json", t.ID)
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

type TaskList []Task

func (t TaskList) Request(server string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, server+"/projects/api/v3/tasks.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *TaskList) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tasks []Task `json:"tasks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Tasks
	return nil
}
