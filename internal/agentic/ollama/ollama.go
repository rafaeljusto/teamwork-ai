package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rafaeljusto/teamwork-ai/internal/agentic"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
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
// https://github.com/ollama/ollama/blob/65bff664cb39ed16a1fa814b0228e4e48d7234ba/api/types.go
type ollama struct {
	server    string
	mcpClient *agentic.MCPClient
	client    *http.Client
	model     string
	logger    *slog.Logger
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
func (o *ollama) Init(dsn string, mcpClient *agentic.MCPClient, logger *slog.Logger) error {
	o.mcpClient = mcpClient
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

// do sends a request to the Ollama server. It injects MCP tools into the
// request, optionally filtering them by the provided methods. If the methods
// slice is empty, all tools are included. If the methods slice contains
// `mcp.MethodNone`, no tools are included in the request. It will handle the
// multiple roundtrips excuting MCP callbacks until no more tool calls are
// present in the response.
func (o *ollama) do(ctx context.Context, aiRequest request, methods ...twmcp.Method) (response, error) {
	if !slices.Contains(methods, twmcp.MethodNone) {
		mcpTools, err := o.mcpClient.Tools(ctx, methods...)
		if err != nil {
			return response{}, fmt.Errorf("failed to load tools: %w", err)
		}
		if aiRequest.Tools == nil {
			aiRequest.Tools = make([]requestTool, 0, len(mcpTools))
		}
		for _, tool := range mcpTools {
			aiRequest.Tools = append(aiRequest.Tools, newRequestTool(tool))
		}
	}

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
	aiResponse.adjustToolCalls(aiRequest.Tools)

	var rework bool
	for _, toolCall := range aiResponse.Message.ToolCalls {
		o.logger.Debug("executing tool",
			slog.String("name", toolCall.Function.Name),
			slog.Any("arguments", toolCall.Function.Arguments),
		)
		toolResult, err := o.mcpClient.ExecuteTool(ctx, toolCall.Function.Name, mcp.CallToolParams{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
		if err != nil {
			return response{}, fmt.Errorf("failed to execute tool %q: %w", toolCall.Function.Name, err)
		}

		aiRequest.addResponseMessage(aiResponse.Message)
		if toolResult.IsError {
			o.logger.Debug("tool returned an error",
				slog.String("name", toolCall.Function.Name),
				slog.Any("error", toolResult.Content),
			)
		}
		if len(toolResult.Content) > 0 {
			// https://github.com/ollama/ollama-python/blob/63ca74762284100b2f0ad207bc00fa3d32720fbd/examples/tools.py
			for _, content := range toolResult.Content {
				if t, ok := content.(mcp.TextContent); ok {
					aiRequest.addToolMessage(toolCall.Function.Name, t.Text)
				}
			}
			rework = true
		}
	}

	if rework {
		aiResponse, err = o.do(ctx, aiRequest)
		if err != nil {
			return response{}, fmt.Errorf("failed to iterate with the LLM: %w", err)
		}
	}

	return aiResponse, nil
}

type request struct {
	Model    string           `json:"model"`
	Messages []requestMessage `json:"messages"`
	Stream   bool             `json:"stream"`
	Tools    []requestTool    `json:"tools"`
}

func (r *request) addUserMessage(content string) {
	r.Messages = append(r.Messages, requestMessage{
		Role:    "user",
		Content: content,
	})
}

func (r *request) addResponseMessage(response responseMessage) {
	r.Messages = append(r.Messages, requestMessage{
		Role:      response.Role,
		Content:   response.Content,
		ToolCalls: response.ToolCalls,
	})
}

func (r *request) addToolMessage(name, content string) {
	r.Messages = append(r.Messages, requestMessage{
		Role:    "tool",
		Name:    name,
		Content: content,
	})
}

type requestMessage struct {
	Role      string `json:"role"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content"`
	ToolCalls []struct {
		Function struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"function"`
	} `json:"tool_calls"`
}

type requestTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string         `json:"type"`
			Properties map[string]any `json:"properties"`
			Required   []string       `json:"required"`
		} `json:"parameters"`
	} `json:"function"`
}

func newRequestTool(mcpTool mcp.Tool) requestTool {
	var requestTool requestTool
	requestTool.Type = "function"
	requestTool.Function.Name = mcpTool.Name
	requestTool.Function.Description = mcpTool.Description
	requestTool.Function.Parameters.Type = "object"
	requestTool.Function.Parameters.Properties = make(map[string]any)
	maps.Copy(requestTool.Function.Parameters.Properties, mcpTool.InputSchema.Properties)
	requestTool.Function.Parameters.Required = mcpTool.InputSchema.Required
	return requestTool
}

type response struct {
	Message struct {
		Role      string `json:"role"`
		Content   string `json:"content"`
		ToolCalls []struct {
			Function struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			} `json:"function"`
		} `json:"tool_calls"`
	} `json:"message"`
}

func (r *response) decode(target any) error {
	return json.Unmarshal([]byte(r.Message.Content), target)
}

// adjustToolCalls adjust the parameters of the tool calls in the response to
// match the expected format for the MCP client. This is necessary because the
// Ollama API returns tool calls in a specific format that may not directly
// match the MCP client's expectations.
func (r *response) adjustToolCalls(mcpTools []requestTool) {
	for i, toolCall := range r.Message.ToolCalls {
		for _, mcpTool := range mcpTools {
			if toolCall.Function.Name != mcpTool.Function.Name {
				continue
			}
			for llmArgumentName, llmArgumentValue := range toolCall.Function.Arguments {
				if mcpToolProperty, ok := mcpTool.Function.Parameters.Properties[llmArgumentName]; ok {
					// adjust the arguments to match the MCP client's expectations. LLMs
					// generate text tokens, and numbers like "30" and 30 are very similar
					// in context.
					mcpToolPropertyType := mcpToolProperty.(map[string]any)["type"].(string)
					if mcpToolPropertyType == "number" {
						switch v := llmArgumentValue.(type) {
						case string:
							if n, err := strconv.ParseFloat(v, 64); err == nil {
								llmArgumentValue = n
							}
						case []byte:
							if n, err := strconv.ParseFloat(string(v), 64); err == nil {
								llmArgumentValue = n
							}
						}
						r.Message.ToolCalls[i].Function.Arguments[llmArgumentName] = llmArgumentValue
					}
				}
			}
		}
	}
}

type responseMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	ToolCalls []struct {
		Function struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"function"`
	} `json:"tool_calls"`
}
