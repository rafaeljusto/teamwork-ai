// Package workload provides functionality to interact with the Teamwork API for
// managing workloads. It includes structures and methods to retrieve workload
// information for users, including their capacity and availability.
package workload

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
)

// Workload represents the response from the API when retrieving a workload. It
// contains information about users, their capacity, and availability.
//
// More information can be found in the Teamwork API documentation:
// https://support.teamwork.com/projects/workload
type Workload struct {
	Users []User `json:"users"`
}

// User represents a user in the workload response. It contains the user's ID
// and a map of dates with their corresponding workload information.
type User struct {
	ID    int64                      `json:"userId"`
	Dates map[teamwork.Date]UserDate `json:"dates"`
}

// UserDate represents the workload information for a specific user on a
// specific date. It includes the user's capacity, capacity in minutes, and
// whether the user is unavailable on that date.
type UserDate struct {
	Capacity        float64 `json:"capacity"`
	CapacityMinutes int64   `json:"capacityMinutes"`
	UnavailableDay  bool    `json:"unavailableDay"`
}

// Single represents a request to retrieve a workload. You must provide the
// start and end dates.
//
// https://apidocs.teamwork.com/docs/teamwork/endpoints-by-object/workflows/get-projects-api-v3-workflows-json
type Single struct {
	Request struct {
		Filters struct {
			StartDate teamwork.Date
			EndDate   teamwork.Date
			UserIDs   []int64
			Page      int64
			PageSize  int64
			Include   []string
		}
	}
	Response struct {
		Meta struct {
			Page struct {
				HasMore bool `json:"hasMore"`
			} `json:"page"`
		} `json:"meta"`
		Workload Workload `json:"workload"`
		Included struct {
			Users map[string]struct {
				ID          int64                  `json:"id"`
				LengthOfDay float64                `json:"lengthOfDay"`
				WorkingHour *teamwork.Relationship `json:"workingHour"`
			} `json:"users,omitempty"`
			WorkingHours map[string]struct {
				ID      int64                   `json:"id"`
				Object  teamwork.Relationship   `json:"object"`
				Entries []teamwork.Relationship `json:"entries"`
			} `json:"workingHours,omitempty"`
			WorkingHoursEntries map[string]struct {
				ID          int64                 `json:"id"`
				WorkingHour teamwork.Relationship `json:"workingHour"`
				Weekday     string                `json:"weekday"`
				TaskHours   float64               `json:"taskHours"`
			} `json:"workingHourEntries,omitempty"`
		} `json:"included"`
	}
}

// HTTPRequest creates an HTTP request to retrieve a single user by their ID.
func (s Single) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/workload", server)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if !time.Time(s.Request.Filters.StartDate).IsZero() {
		query.Set("startDate", s.Request.Filters.StartDate.String())
	}
	if !time.Time(s.Request.Filters.EndDate).IsZero() {
		query.Set("endDate", s.Request.Filters.EndDate.String())
	}
	if len(s.Request.Filters.UserIDs) > 0 {
		var ids []string
		for _, id := range s.Request.Filters.UserIDs {
			ids = append(ids, strconv.FormatInt(id, 10))
		}
		query.Set("userIds", strings.Join(ids, ","))
	}
	if s.Request.Filters.Page > 0 {
		query.Set("page", strconv.FormatInt(s.Request.Filters.Page, 10))
	}
	if s.Request.Filters.PageSize > 0 {
		query.Set("pageSize", strconv.FormatInt(s.Request.Filters.PageSize, 10))
	}
	if len(s.Request.Filters.Include) > 0 {
		query.Set("include", strings.Join(s.Request.Filters.Include, ","))
	}

	// to reduce the size of the response, we omit empty date entries where the
	// user has no capacity and is not unavailable.
	query.Set("omitEmptyDateEntries", "true")

	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Single instance.
func (s *Single) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.Response)
}
