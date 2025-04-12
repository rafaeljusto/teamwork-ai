package teamwork

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// LegacyUserGroups represents a collection of users, companies, and teams
// in a legacy format, where IDs are represented as strings.
type LegacyUserGroups struct {
	UserIDs    []int64
	CompanyIDs []int64
	TeamIDs    []int64
}

// MarshalJSON encodes the LegacyUserGroups as a JSON object.
func (m LegacyUserGroups) MarshalJSON() ([]byte, error) {
	var result string
	for _, id := range m.UserIDs {
		if result != "" {
			result += ","
		}
		result += strconv.FormatInt(id, 10)
	}
	for _, id := range m.CompanyIDs {
		if result != "" {
			result += ","
		}
		result += "c" + strconv.FormatInt(id, 10)
	}
	for _, id := range m.TeamIDs {
		if result != "" {
			result += ","
		}
		result += "t" + strconv.FormatInt(id, 10)
	}
	return []byte(`"` + result + `"`), nil
}

// UnmarshalJSON decodes a JSON string into a LegacyUserGroups type.
func (m *LegacyUserGroups) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	for part := range strings.SplitSeq(str, ",") {
		if len(part) == 0 {
			continue
		}
		switch part[0] {
		case 'c':
			if len(part) < 2 {
				return fmt.Errorf("invalid company ID format: %s", part)
			}
			id, err := strconv.ParseInt(part[1:], 10, 64)
			if err != nil {
				return err
			}
			m.CompanyIDs = append(m.CompanyIDs, id)
		case 't':
			if len(part) < 2 {
				return fmt.Errorf("invalid team ID format: %s", part)
			}
			id, err := strconv.ParseInt(part[1:], 10, 64)
			if err != nil {
				return err
			}
			m.TeamIDs = append(m.TeamIDs, id)
		default:
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return err
			}
			m.UserIDs = append(m.UserIDs, id)
		}
	}
	return nil
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

// LegacyDate is a type alias for time.Time, used to represent date values in
// the API.
type LegacyDate time.Time

// MarshalJSON encodes the LegacyDate as a string in the format "20060102".
func (d LegacyDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("20060102") + `"`), nil
}

// UnmarshalJSON decodes a JSON string into a LegacyDate type.
func (d *LegacyDate) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsedTime, err := time.Parse("20060102", str)
	if err != nil {
		return err
	}
	*d = LegacyDate(parsedTime)
	return nil
}

// LegacyNumber is a type alias for int64, used to represent numeric values in
// the API.
type LegacyNumber int64

// MarshalJSON encodes the LegacyNumber as a string.
func (n LegacyNumber) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(n), 10) + `"`), nil
}

// UnmarshalJSON decodes a JSON string into a LegacyNumber type.
func (n *LegacyNumber) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsedInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*n = LegacyNumber(parsedInt)
	return nil
}
