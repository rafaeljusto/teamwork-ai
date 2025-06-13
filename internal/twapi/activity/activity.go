// Package activity provides the API implementation to retrieve activity of an
// installation or project.
package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
)

// Activity represents an activity log entry in Teamwork.com.
type Activity struct {
	ID               int64      `json:"id"`
	Action           Action     `json:"activityType"`
	LatestAction     Action     `json:"latestActivityType"`
	At               time.Time  `json:"dateTime"`
	Description      *string    `json:"description"`
	ExtraDescription *string    `json:"extraDescription"`
	PublicInfo       *string    `json:"publicInfo"`
	DueAt            *time.Time `json:"dueDate"`
	ForUserName      *string    `json:"forUserName"`
	ItemLink         *string    `json:"itemLink"`
	Link             *string    `json:"link"`

	User    twapi.Relationship  `json:"user"`
	ForUser *twapi.Relationship `json:"forUser"`
	Project twapi.Relationship  `json:"project"`
	Company twapi.Relationship  `json:"company"`
	Item    twapi.Relationship  `json:"item"`
}

// Action contains all possible activity types.
type Action string

// List of activity types.
const (
	LogTypeNew       Action = "new"
	LogTypeEdited    Action = "edited"
	LogTypeCompleted Action = "completed"
	LogTypeReopened  Action = "reopened"
	LogTypeDeleted   Action = "deleted"
	LogTypeUndeleted Action = "undeleted"
	LogTypeLiked     Action = "liked"
	LogTypeReacted   Action = "reacted"
	LogTypeViewed    Action = "viewed"
)

// LogItemType contains all possible activity item types.
type LogItemType string

// List of activity item types.
const (
	LogItemTypeMessage          LogItemType = "message"
	LogItemTypeComment          LogItemType = "comment"
	LogItemTypeTask             LogItemType = "task"
	LogItemTypeTasklist         LogItemType = "tasklist"
	LogItemTypeTaskgroup        LogItemType = "taskgroup"
	LogItemTypeMilestone        LogItemType = "milestone"
	LogItemTypeFile             LogItemType = "file"
	LogItemTypeForm             LogItemType = "form"
	LogItemTypeNotebook         LogItemType = "notebook"
	LogItemTypeTimelog          LogItemType = "timelog"
	LogItemTypeTaskComment      LogItemType = "task_comment"
	LogItemTypeNotebookComment  LogItemType = "notebook_comment"
	LogItemTypeFileComment      LogItemType = "file_comment"
	LogItemTypeLinkComment      LogItemType = "link_comment"
	LogItemTypeMilestoneComment LogItemType = "milestone_comment"
	LogItemTypeProject          LogItemType = "project"
	LogItemTypeLink             LogItemType = "link"
	LogItemTypeBillingInvoice   LogItemType = "billingInvoice"
	LogItemTypeRisk             LogItemType = "risk"
	LogItemTypeProjectUpdate    LogItemType = "projectUpdate"
	LogItemTypeReacted          LogItemType = "reacted"
	LogItemTypeBudget           LogItemType = "budget"
)

// UnmarshalText decodes the text into a LogItemType.
func (l *LogItemType) UnmarshalText(text []byte) error {
	if l == nil {
		panic("unmarshal LogItemType: nil pointer")
	}
	logItemType := LogItemType(strings.ToLower(string(text)))
	switch logItemType {
	case LogItemTypeMessage,
		LogItemTypeComment,
		LogItemTypeTask,
		LogItemTypeTasklist,
		LogItemTypeTaskgroup,
		LogItemTypeMilestone,
		LogItemTypeFile,
		LogItemTypeForm,
		LogItemTypeNotebook,
		LogItemTypeTimelog,
		LogItemTypeTaskComment,
		LogItemTypeNotebookComment,
		LogItemTypeFileComment,
		LogItemTypeLinkComment,
		LogItemTypeMilestoneComment,
		LogItemTypeProject,
		LogItemTypeLink,
		LogItemTypeBillingInvoice,
		LogItemTypeRisk,
		LogItemTypeProjectUpdate,
		LogItemTypeReacted,
		LogItemTypeBudget:
		*l = logItemType
	default:
		return fmt.Errorf("invalid log item type: %q", text)
	}
	return nil
}

// Multiple represents a request to retrieve multiple activities.
type Multiple struct {
	Request struct {
		Path struct {
			ProjectID int64
		}
		Filters struct {
			StartDate    time.Time
			EndDate      time.Time
			LogItemTypes []LogItemType
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
		Activities []Activity `json:"activities"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple activities.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	var uri string
	if m.Request.Path.ProjectID > 0 {
		uri = fmt.Sprintf("%s/projects/api/v3/projects/%d/latestactivity.json", server, m.Request.Path.ProjectID)
	} else {
		uri = fmt.Sprintf("%s/projects/api/v3/latestactivity.json", server)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if !m.Request.Filters.StartDate.IsZero() {
		query.Set("startDate", m.Request.Filters.StartDate.Format(time.RFC3339))
	}
	if !m.Request.Filters.EndDate.IsZero() {
		query.Set("endDate", m.Request.Filters.EndDate.Format(time.RFC3339))
	}
	if len(m.Request.Filters.LogItemTypes) > 0 {
		logItemTypes := make([]string, len(m.Request.Filters.LogItemTypes))
		for i, logType := range m.Request.Filters.LogItemTypes {
			logItemTypes[i] = string(logType)
		}
		query.Set("activityTypes", strings.Join(logItemTypes, ","))
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
