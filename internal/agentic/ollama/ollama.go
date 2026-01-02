package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
)

var _ agentic.Agentic = (*ollama)(nil)

func init() {
	agentic.Register("ollama", &ollama{})
}

// ollama is an open-source, cross-platform framework that simplifies running
// large language models (LLMs) locally on your computer. This specific instance
// implements all required functions to allow the Teamwork AI agentic
// implementation.
//
// The API reference is available at:
// https://github.com/ollama/ollama/blob/6a74bba7e7e19bf5f5aeacb039a1537afa3522a5/docs/api.md
type ollama struct {
	server string
	client *http.Client
	model  string
	logger *slog.Logger
}

// Init initializes the Ollama instance with the provided DSN. The DSN must have
// the format:
//
//	`http[s]://[username[:password]@]host[:port]/model`.
//
// The server URL should point to the Ollama base URL, and the model name should
// be the name of the model to be used (e.g. "llama3.2").
//
// TODO(rafaeljusto): Add support for custom HTTP client.
func (o *ollama) Init(dsn string, logger *slog.Logger) error {
	o.client = http.DefaultClient
	o.logger = logger

	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s", parsedURL.Scheme)
	}
	if parsedURL.Path == "" {
		return fmt.Errorf("missing model name in DSN")
	}
	o.model = parsedURL.Path[1:]

	parsedURL.Path = ""
	o.server = parsedURL.String()
	return nil
}

func (o *ollama) do(ctx context.Context, aiRequest request) (response, error) {
	body, err := json.Marshal(aiRequest)
	if err != nil {
		return response{}, fmt.Errorf("failed to encode request: %w", err)
	}

	url, err := url.JoinPath(o.server, "/api/chat")
	if err != nil {
		return response{}, fmt.Errorf("failed to build url: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return response{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := o.client.Do(httpRequest)
	if err != nil {
		return response{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			o.logger.Error("failed to close response body",
				slog.String("error", err.Error()),
			)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		if body, err := io.ReadAll(httpResponse.Body); err == nil {
			return response{}, fmt.Errorf("unexpected status code: %d, body: %s", httpResponse.StatusCode, string(body))
		}
		return response{}, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	var aiResponse response
	if err = json.NewDecoder(httpResponse.Body).Decode(&aiResponse); err != nil {
		return response{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return aiResponse, nil
}

type request struct {
	Model    string           `json:"model"`
	Messages []requestMessage `json:"messages"`
	Stream   bool             `json:"stream"`
}

func (r *request) addSystemMessage(content string) {
	r.Messages = append(r.Messages, requestMessage{
		Role:    "system",
		Content: content,
	})
}

func (r *request) addUserMessage(content string) {
	r.Messages = append(r.Messages, requestMessage{
		Role:    "user",
		Content: content,
	})
}

type requestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func (r *response) decode(target any) error {
	return json.Unmarshal([]byte(r.Message.Content), target)
}
