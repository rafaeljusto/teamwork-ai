package webhook

import "github.com/rafaeljusto/teamwork-ai/internal/twapi"

// TaskData represents the payload for the task related webhook events in
// Teamwork.com.
type TaskData struct {
	Project struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"project"`
	Task struct {
		ID               int64       `json:"id"`
		Name             string      `json:"name"`
		Description      string      `json:"description"`
		AssignedUserIDs  []int64     `json:"assignedUserIds"`
		Status           string      `json:"status"`
		StartDate        *twapi.Date `json:"startDate"`
		DueDate          *twapi.Date `json:"dueDate"`
		EstimatedMinutes int64       `json:"estimatedMinutes"`
	} `json:"task"`
	Tasklist struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"taskList"`
}
