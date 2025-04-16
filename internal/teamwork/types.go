package teamwork

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
