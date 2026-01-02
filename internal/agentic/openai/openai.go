package openai

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

var _ agentic.Agentic = (*openai)(nil)

func init() {
	agentic.Register("openai", &openai{})
}

// openai is an american company that provides a suite of AI tools and services.
// This specific instance implements all required functions to allow the
// Teamwork AI agentic implementation. The OpenAI API is a cloud-based service
// that provides access to OpenAI's language models, including GPT-3.5 and
// GPT-4. It allows developers to integrate AI capabilities into their
// applications, enabling tasks such as natural language understanding, text
// generation, and conversation simulation.
//
// The API reference is available at:
// https://platform.openai.com/docs/api-reference/introduction
type openai struct {
	client *http.Client
	logger *slog.Logger
	model  string
	token  string
}

// Init initializes the OpenAI instance with the provided DSN. The DSN must have
// the format:
//
//	`model:token`.
//
// The model name should be the name of the model to be used (e.g.
// "gpt-3.5-turbo"). The token should be the OpenAI API key.
//
// TODO(rafaeljusto): Add support for custom HTTP client.
func (o *openai) Init(dsn string, logger *slog.Logger) error {
	o.client = http.DefaultClient
	o.logger = logger

	dsnParts := strings.Split(dsn, ":")
	if len(dsnParts) != 2 {
		return fmt.Errorf("invalid DSN format: %s", dsn)
	}
	o.model = dsnParts[0]
	o.token = dsnParts[1]
	return nil
}

func (o *openai) do(ctx context.Context, aiRequest request) (response, error) {
	body, err := json.Marshal(aiRequest)
	if err != nil {
		return response{}, fmt.Errorf("failed to encode request: %w", err)
	}

	url := "https://api.openai.com/v1/responses"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return response{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+o.token)
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
	Messages []requestMessage `json:"input"`
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
	Output []output `json:"output"`
}

func (r *response) decode(target any) error {
	if len(r.Output) == 0 {
		return fmt.Errorf("no outputs in response")
	}
	if len(r.Output) > 1 {
		return fmt.Errorf("multiple outputs in response")
	}
	if len(r.Output[0].Content) == 0 {
		return fmt.Errorf("no content in output")
	}
	if len(r.Output[0].Content) > 1 {
		return fmt.Errorf("multiple contents in output")
	}
	return json.Unmarshal([]byte(r.Output[0].Content[0].Text), target)
}

type output struct {
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
