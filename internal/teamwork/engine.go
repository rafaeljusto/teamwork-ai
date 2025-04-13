package teamwork

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Engine struct {
	server     string
	apiToken   string
	httpClient *http.Client
}

func NewEngine(server, apiToken string) *Engine {
	return &Engine{
		server:     server,
		apiToken:   apiToken,
		httpClient: http.DefaultClient,
	}
}

func (e *Engine) Do(entity Entity) error {
	req, err := entity.Request(e.server)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(e.apiToken, "")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(entity)
}

type Entity interface {
	Request(server string) (*http.Request, error)
}
