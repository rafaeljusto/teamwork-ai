package skill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Skill struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	CreatedByUserID int64      `json:"createdByUser"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedByUserID *int64     `json:"updatedByUser"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedByUserID *int64     `json:"deletedByUser"`
	DeletedAt       *time.Time `json:"deletedAt"`
}

type SingleSkill Skill

func (t SingleSkill) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, t.ID)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *SingleSkill) UnmarshalJSON(data []byte) error {
	var raw struct {
		Skill Skill `json:"skill"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = SingleSkill(raw.Skill)
	return nil
}

type MultipleSkills []Skill

func (t MultipleSkills) Request(server string) (*http.Request, error) {
	url := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t *MultipleSkills) UnmarshalJSON(data []byte) error {
	var raw struct {
		Skills []Skill `json:"skills"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = raw.Skills
	return nil
}

type SkillCreation struct {
	Name    string  `json:"name"`
	UserIDs []int64 `json:"userIds"`
}

func (t SkillCreation) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills.json", server)
	paylaod := struct {
		Skill SkillCreation `json:"skill"`
	}{Skill: t}
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

type SkillUpdate struct {
	ID    int64
	Skill struct {
		Name    *string `json:"name,omitempty"`
		UserIDs []int64 `json:"userIds,omitempty"`
	}
}

func (t SkillUpdate) Request(server string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/projects/api/v3/skills/%d.json", server, t.ID)
	paylaod := struct {
		Skill any `json:"skill"`
	}{Skill: t.Skill}
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
