# Teamwork.com AI

[![Go Reference](https://pkg.go.dev/badge/github.com/rafaeljusto/teamwork-ai.svg)](https://pkg.go.dev/github.com/rafaeljusto/teamwork-ai)
![Test](https://github.com/rafaeljusto/teamwork-ai/actions/workflows/test.yml/badge.svg)

Unofficial extension for Teamwork.com to integrate AI capabilities.

## MCP server

Implements the Model Context Protocol (MCP) to allow AI agents to interact with
Teamwork.com. This server acts as a bridge between AI clients and Teamwork.com,
adding tools to create tasks, projects, and more.

[![MCP example](https://img.youtube.com/vi/QTGM7cQT7Ew/0.jpg)](https://www.youtube.com/watch?v=QTGM7cQT7Ew)

### Usage

The server works using 2 different modes:

1. **Local mode**: Runs as a local server (`stdio`), allowing AI clients to
   connect directly. You can, for example, install [Claude
   Desktop](https://claude.ai/download), which supports MCP, and connect
   configure it like:

```json
{
  "mcpServers": {
    "Teamwork AI": {
      "command": "teamwork-ai",
      "args": [
        "-mode=stdio"
      ],
      "env": {
        "TWAI_TEAMWORK_SERVER": "https://<installation>.teamwork.com",
        "TWAI_TEAMWORK_API_TOKEN": "<api-token>"
      }
    }
  }
}
```

It assumes that [the binary](cmd/mcp/main.go) is compiled and installed in one
of the directories in your `PATH`. The `<installation>` is the name of your
Teamwork.com installation, and `<api-token>` is your API token from Teamwork.com
profile (no support for OAuth2 yet).

2. **Remote mode**: Runs as a remote server (`sse`). This can also be used with
   [Claude Desktop](https://claude.ai/download) with the following
   configuration:

```json
{
  "mcpServers": {
    "math": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "https://<server>/sse"
      ]
    }
  }
}
```

Where `<server>` is the URL of the remote MCP server.

### Debug

For debugging purposes, you can run the `MCP Inspector tool`:

```bash
npx @modelcontextprotocol/inspector node build/index.js
```