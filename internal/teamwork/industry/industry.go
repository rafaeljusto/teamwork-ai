// Package industry provides a functionality to manage industries in
// Teamwork.com. It allows for the retrieval of multiple industries, each
// represented by an Industry struct that contains an ID and a name.
package industry

import (
	"context"
	"encoding/json"
	"net/http"
)

// Industry represents an industry that a company can belong to in Teamwork.com.
// It contains an ID and a name for the industry.
type Industry struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Multiple represents a request to retrieve multiple industries.
//
// Not documented.
type Multiple struct {
	Response struct {
		Industries []Industry `json:"industries"`
	}
}

// HTTPRequest creates an HTTP request to retrieve multiple industries.
func (m Multiple) HTTPRequest(ctx context.Context, server string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+"/industries.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// UnmarshalJSON decodes the JSON data into a Multiple instance.
func (m *Multiple) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Response)
}
