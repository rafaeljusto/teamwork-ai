package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
)

var _ agentic.Agentic = (*anthropic)(nil)

func init() {
	agentic.Register("anthropic", &anthropic{})
}

// anthropic is an american company that provides a suite of AI tools and
// services. This specific instance implements all required functions to allow
// the Teamwork AI agentic implementation. The Anthropic API is a cloud-based
// service that provides access to Anthropic's language models, including Claude
// 1 and Claude 2. It allows developers to integrate AI capabilities into their
// applications, enabling tasks such as natural language understanding, text
// generation, and conversation simulation.
//
// The API reference is available at:
// https://docs.anthropic.com/en/api
type anthropic struct {
	client *http.Client
	logger *slog.Logger
	model  string
	token  string
}

// Init initializes the anthropic instance with the provided DSN. The DSN must
// have the format:
//
//	`model:token`.
//
// The model name should be the name of the model to be used (e.g.
// "claude-1"). The token should be the Anthropic API key.
//
// TODO(rafaeljusto): Add support for custom HTTP client.
func (a *anthropic) Init(dsn string, logger *slog.Logger) error {
	a.client = http.DefaultClient
	a.logger = logger

	dsnParts := strings.Split(dsn, ":")
	if len(dsnParts) != 2 {
		return fmt.Errorf("invalid DSN format: %s", dsn)
	}
	a.model = dsnParts[0]
	a.token = dsnParts[1]
	return nil
}

func (a *anthropic) do(ctx context.Context, aiRequest request) (response, error) {
	body, err := json.Marshal(aiRequest)
	if err != nil {
		return response{}, fmt.Errorf("failed to encode request: %w", err)
	}

	url := "https://api.anthropic.com/v1/messages"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return response{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpRequest.Header.Set("x-api-key", a.token)
	httpRequest.Header.Set("anthropic-version", "2023-06-01")
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := a.client.Do(httpRequest)
	if err != nil {
		return response{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			a.logger.Error("failed to close response body",
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
	Model     string           `json:"model"`
	Messages  []requestMessage `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
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
	Contents []content `json:"content"`
}

func (r *response) decode(target any) error {
	if len(r.Contents) == 0 {
		return fmt.Errorf("no content in response")
	}
	if len(r.Contents) > 1 {
		return fmt.Errorf("multiple contents in response")
	}
	return json.Unmarshal([]byte(r.Contents[0].Text), target)
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
