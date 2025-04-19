package teamwork

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// Engine is the main structure that handles communication with the Teamwork
// API. It is responsible for sending requests and processing responses for
// various entities such as projects, companies, and skills. The Engine uses an
// HTTP client to make requests to the Teamwork API and requires a server URL
// and an API token for authentication. It also accepts a logger for logging
// purposes.
type Engine struct {
	server     string
	apiToken   string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewEngine creates a new instance of the Engine with the provided server
// URL, API token, and logger.
//
// TODO(rafaeljusto): Add support for custom HTTP client.
func NewEngine(server, apiToken string, logger *slog.Logger) *Engine {
	return &Engine{
		server:     server,
		apiToken:   apiToken,
		httpClient: http.DefaultClient,
		logger:     logger,
	}
}

// Do executes the request for the given entity. It constructs an HTTP request
// using the entity's HTTPRequest method, sets the necessary authentication
// headers, and sends the request using the Engine's HTTP client. If the request
// is successful, it decodes the response body into the entity. If the request
// fails or the response status code indicates an error, it returns an error
// with a descriptive message. The method also ensures that the response body is
// closed after processing to prevent resource leaks.
func (e *Engine) Do(ctx context.Context, entity Entity) error {
	req, err := entity.HTTPRequest(ctx, e.server)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(e.apiToken, "")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			e.logger.Error("failed to close response body",
				slog.String("error", err.Error()),
			)
		}
	}()

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

// Entity is an interface that defines the methods required for an entity to be
// used with the Teamwork Engine. An entity must implement the Request method,
// which constructs an HTTP request for the entity. The HTTPRequest method takes
// a context and a server URL as parameters and returns an HTTP request and an
// error if any occurs during the request creation. This interface allows the
// Engine to handle different types of entities (like projects, companies,
// skills, etc.) in a uniform way, enabling the Engine to send requests and
// process responses without needing to know the specific details of each entity
// type. This abstraction allows for flexibility and extensibility in the
// Teamwork API client implementation, as new entity types can be added without
// modifying the Engine's core logic.
type Entity interface {
	HTTPRequest(ctx context.Context, server string) (*http.Request, error)
}
