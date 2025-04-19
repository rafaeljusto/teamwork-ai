package teamwork

import (
	"encoding/json"
	"time"
)

// Relationship describes the relation between the main entity and a sideload type.
type Relationship struct {
	ID   int64          `json:"id"`
	Type string         `json:"type"`
	Meta map[string]any `json:"meta,omitempty"`
}

// UserGroups represents a collection of users, companies, and teams.
type UserGroups struct {
	UserIDs    []int64 `json:"userIds"`
	CompanyIDs []int64 `json:"companyIds"`
	TeamIDs    []int64 `json:"teamIds"`
}

// Date is a type alias for time.Time, used to represent date values in the API.
type Date time.Time

// MarshalJSON encodes the Date as a string in the format "2006-01-02".
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}

// UnmarshalJSON decodes a JSON string into a Date type.
func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = Date(parsedTime)
	return nil
}
