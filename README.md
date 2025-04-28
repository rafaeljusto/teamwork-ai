# Teamwork.com AI

[![Go Reference](https://pkg.go.dev/badge/github.com/rafaeljusto/teamwork-ai.svg)](https://pkg.go.dev/github.com/rafaeljusto/teamwork-ai)
![Test](https://github.com/rafaeljusto/teamwork-ai/actions/workflows/test.yml/badge.svg)

![Logo](teamwork-ai.gif)

**Unofficial** extension for [Teamwork.com](https://teamwork.com) to integrate
AI capabilities.

> [!WARNING]
> When interacting with LLMs, be aware that the data you provide may be used to
> train and improve AI models. This may include sharing your data with
> third-party providers, which could lead to potential privacy and security
> risks. Always review the terms of service and privacy policies of the AI
> providers you use to understand how your data will be handled.

## MCP server

Implements the Model Context Protocol (MCP) to allow AI agents to interact with
Teamwork.com. This server acts as a bridge between AI clients and Teamwork.com,
adding tools to create tasks, projects, and more.

[![MCP example](https://img.youtube.com/vi/QTGM7cQT7Ew/0.jpg)](https://www.youtube.com/watch?v=QTGM7cQT7Ew)

Some interesting things you can do with this server:

```
> Could you please create a projects with the steps to create a new house?

The AI client will create a project named "New House" with tasklist and tasks
with the specific steps to create a new house.
```

```
> Could you assign the tasks from the "New House" project to users that have
> the available skills to fulfill them? Leave the tasks unassigned if no
> user has the required skills.

The AI client will automatically query the projects, project's members, 
tasklists, tasks and skills to correctly assign the tasks. It analyzes the
tasklist name, the task title and description to find the best match for the
users' skills.
```

### ⚡️ Usage

The server works using 2 different modes:

1. **Local mode**: Runs as a local server (`stdio`), allowing AI clients to
   connect directly. You can, for example, install [Claude
   Desktop](https://claude.ai/download), which supports MCP, and configure it
   like:

```json
{
  "mcpServers": {
    "Teamwork AI": {
      "command": "teamwork-mcp",
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

> [!TIP]
> For more information regarding the Claude Desktop configuration, refer to the
> [MCP documentation](https://modelcontextprotocol.io/quickstart/user). To
> further learn about Claude and MCP, you can also check [this
> content](https://www.claudemcp.com/).
>
> Be aware of the daily usage limits and how to follow the best practices
> [here](https://support.anthropic.com/en/articles/9797557-usage-limit-best-practices).

It assumes that [the binary](cmd/mcp/main.go) is compiled and installed in one
of the directories in your `PATH`. The `<installation>` is the name of your
Teamwork.com installation, and `<api-token>` is your API token from Teamwork.com
profile (no support for OAuth2 yet).

> [!IMPORTANT]
> The API token MUST have enough permissions to execute the desired actions when
> interacting with the AI client.

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

> [!NOTE]
> When using Claude Desktop, for every MCP tool execution the AI client will ask
> for confirmation before executing a tool. This is a safety feature to prevent
> unintended actions.

### 🔌 Supported entities

Below is a table summarizing the supported entities and their operations in the
MCP server.

| Entity            | Create | Retrieve | Update | Delete | Extra                                                           |
|-------------------|--------|----------|--------|--------|-----------------------------------------------------------------|
| Projects          | ✅     | ✅       | ✅      | ❌     |                                                                 |
| Tasklists         | ✅     | ✅       | ✅      | ❌     | Retrieve by project                                             |
| Tasks             | ✅     | ✅       | ✅      | ❌     | Retrieve by project; retrieve by tasklist                       |
| Companies/Clients | ✅     | ✅       | ✅      | ❌     |                                                                 |
| Users/People      | ✅     | ✅       | ✅      | ❌     | Retrieve by project; add to a project; assign/unassign job role |
| Skills            | ✅     | ✅       | ✅      | ❌     |                                                                 |
| Industries        | ❌     | ✅       | ❌      | ❌     |                                                                 |
| Tags              | ✅     | ✅       | ✅      | ❌     |                                                                 |
| Milestones        | ✅     | ✅       | ✅      | ❌     | Retrieve by project                                             |
| Job roles         | ✅     | ✅       | ✅      | ❌     |                                                                 |

> [!NOTE]
> Not all properties are supported for each entity. And, for now, delete actions
> are not implemented for safety.

### 🤓 Debug

For debugging purposes, you can run the [MCP Inspector
tool](https://github.com/modelcontextprotocol/inspector):

```bash
npx @modelcontextprotocol/inspector node build/index.js
```

> [!TIP]
> For more information regarding the MCP Inspector tool and how to use it, check
> it [here](https://modelcontextprotocol.io/docs/tools/inspector).