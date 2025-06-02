package twapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// EngineOptions defines options for the Teamwork Engine.
type EngineOptions struct {
	idField    string
	idCallback func(id int64)
}

// Option is a function that modifies the EngineOptions.
type Option func(*EngineOptions)

// WithIDCallback sets a callback function that is called with the ID of an
// entity after it is created.
func WithIDCallback(idField string, callback func(id int64)) Option {
	return func(opts *EngineOptions) {
		if idField == "" {
			idField = "id"
		}
		if callback != nil {
			opts.idField = idField
			opts.idCallback = callback
		}
	}
}

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
func (e *Engine) Do(ctx context.Context, entity Entity, optFuncs ...Option) error {
	options := &EngineOptions{
		idField:    "id",
		idCallback: func(int64) {},
	}
	for _, optFunc := range optFuncs {
		optFunc(options)
	}
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

	switch req.Method {
	case http.MethodGet:
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(entity); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
		if resource, ok := entity.(interface{ PopulateResourceWebLink(server string) }); ok {
			resource.PopulateResourceWebLink(e.server)
		}
	case http.MethodPost:
		var body map[string]any
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&body); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
		if id, ok := idSearch(options.idField, body); ok {
			options.idCallback(id)
		}
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

// idSearch is a helper function that recursively searches for an "id" field in
// a map. It returns the first found ID as an int64 and a boolean indicating
// whether an ID was found. It uses a BFS approach to traverse nested maps,
// allowing it to find IDs even in complex JSON structures.
func idSearch(idField string, body map[string]any) (int64, bool) {
	var nestedMaps []map[string]any
	for key, value := range body {
		if strings.EqualFold(key, idField) {
			switch v := value.(type) {
			case int64:
				return v, true
			case float64:
				return int64(v), true
			case string:
				id, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("failed to parse %q as number: %v", v, err))
				}
				return id, true
			default:
				panic(fmt.Sprintf("unexpected type for %q: %T", idField, value))
			}
		} else if nestedMap, ok := value.(map[string]any); ok {
			nestedMaps = append(nestedMaps, nestedMap)
		}
	}
	for _, nestedMap := range nestedMaps {
		if id, found := idSearch(idField, nestedMap); found && id > 0 {
			return id, true
		}
	}
	return 0, false
}
