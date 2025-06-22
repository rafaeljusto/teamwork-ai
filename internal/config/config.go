package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Config stores the configuration of the application.
type Config struct {
	// Port is the port of the MCP server.
	Port int64

	// LoggerLevel is the level of the logger.
	LoggerLevel slog.Level

	// TeamworkServer is the server of the Teamwork API.
	TeamworkServer string

	// TeamworkAPIToken is the API token of the Teamwork API.
	TeamworkAPIToken string

	// Agentic is the agentic configuration.
	Agentic struct {
		// Name is the name of the agentic implementation.
		Name string

		// DSN is the data source name for the agentic model. The format depends on
		// the chosen implementation.
		DSN string

		// MCPClient is the configuration for the MCP client used by the agentic
		// implementation.
		MCPClient struct {
			// StdioPath is the path to the stdio executable. It is used when the mode
			// is "stdio".
			StdioPath string

			// StdioArgs is the list of arguments to be passed to the stdio
			// executable. It is used when the mode is "stdio".
			StdioArgs []string

			// StdioEnvs is the list of environment variables to be passed to the
			// stdio executable. It is used when the mode is "stdio".
			StdioEnvs []string

			// SSEAddress is the address of the SSE server. It is used when the mode
			// is "sse".
			SSEAddress string
		}
	}
}

// DisableMCPClient disables the MCP client by clearing its configuration.
func (c *Config) DisableMCPClient() {
	c.Agentic.MCPClient.StdioPath = ""
	c.Agentic.MCPClient.StdioArgs = nil
	c.Agentic.MCPClient.StdioEnvs = nil
	c.Agentic.MCPClient.SSEAddress = ""
}

// ParseFromEnvs parses the configuration from environment variables.
func ParseFromEnvs() (*Config, error) {
	var config Config
	var errs error
	var err error

	portStr := os.Getenv("TWAI_PORT")
	if portStr != "" {
		config.Port, err = strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to parse TWAI_PORT: %w", err))
		}
	}

	loggerLevel := slog.LevelInfo
	if loggerLevelStr := os.Getenv("TWAI_LOG_LEVEL"); loggerLevelStr != "" {
		if err = loggerLevel.UnmarshalText([]byte(loggerLevelStr)); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to parse TWAI_LOG_LEVEL: %w", err))
		}
	}
	config.LoggerLevel = loggerLevel

	config.TeamworkServer = os.Getenv("TWAI_TEAMWORK_SERVER")
	config.TeamworkAPIToken = os.Getenv("TWAI_TEAMWORK_API_TOKEN")

	config.Agentic.Name = os.Getenv("TWAI_AGENTIC_NAME")
	config.Agentic.DSN = os.Getenv("TWAI_AGENTIC_DSN")

	config.Agentic.MCPClient.StdioPath = os.Getenv("TWAI_AGENTIC_MCP_CLIENT_STDIO_PATH")

	if mcpClientStdioArgs := os.Getenv("TWAI_AGENTIC_MCP_CLIENT_STDIO_ARGS"); mcpClientStdioArgs != "" {
		for arg := range strings.SplitSeq(mcpClientStdioArgs, ",") {
			config.Agentic.MCPClient.StdioArgs = append(config.Agentic.MCPClient.StdioArgs, strings.TrimSpace(arg))
		}
	}

	if mcpClientStdioEnvs := os.Getenv("TWAI_AGENTIC_MCP_CLIENT_STDIO_ENVS"); mcpClientStdioEnvs != "" {
		for env := range strings.SplitSeq(mcpClientStdioEnvs, ",") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				errs = errors.Join(errs, fmt.Errorf("invalid environment variable format: %q", env))
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "" || value == "" {
				errs = errors.Join(errs, fmt.Errorf("invalid environment variable format: %q", env))
				continue
			}
			config.Agentic.MCPClient.StdioEnvs = append(config.Agentic.MCPClient.StdioEnvs, fmt.Sprintf("%s=%s", key, value))
		}
	}

	config.Agentic.MCPClient.SSEAddress = os.Getenv("TWAI_AGENTIC_MCP_CLIENT_SSE_ADDRESS")

	if errs != nil {
		return nil, errs
	}
	return &config, nil
}
