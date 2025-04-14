package teamwork

import (
	"encoding/json"
	"fmt"
	"io"
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if body, err := io.ReadAll(resp.Body); err == nil {
			return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if req.Method == http.MethodGet {
		decoder := json.NewDecoder(resp.Body)
		return decoder.Decode(entity)
	}
	return nil
}

type Entity interface {
	Request(server string) (*http.Request, error)
}
