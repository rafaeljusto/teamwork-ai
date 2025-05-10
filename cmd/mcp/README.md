# MCP server

Implements the Model Context Protocol (MCP) to allow AI agents to interact with
Teamwork.com. This server acts as a bridge between AI clients and Teamwork.com,
adding tools to create tasks, projects, and more.

### 📦 Installing

You can install the MCP server using [`go`](https://go.dev/doc/install):

```bash
go install -o teamwork-mcp github.com/rafaeljusto/teamwork-ai/cmd/mcp@latest
```

The binary will be installed in your `GOPATH/bin` directory, which is usually
`$HOME/go/bin`. Make sure to add this directory to your `PATH` environment
variable to run the `teamwork-mcp` command from anywhere.

Alternatively, you can use the pre-built binaries available in the
[releases](https://github.com/rafaeljusto/teamwork-ai/releases/latest) page.
Download the appropriate binary for your operating system, extract it, and place
it in a directory included in your `PATH`. For example, on Linux (amd64) using
`curl`:

```bash
# detect the latest release
twai_mcp_url=$(curl -s https://api.github.com/repos/rafaeljusto/teamwork-ai/releases/latest | \
  jq '.assets[] | select(.name | contains ("teamwork-mcp-linux-amd64")) | .browser_download_url')

# download the binary and place it in /usr/local/bin
sudo curl -s -O /usr/local/bin/teamwork-mcp ${twai_mcp_url}
```

The example above uses `jq` to parse the JSON response from the GitHub API,
which is a command-line JSON processor. You can find installation instructions
for it [here](https://jqlang.org/download/).

### ⚙️ Configuring

The following environment variables are required to run the MCP server:
- `TWAI_TEAMWORK_SERVER`: The URL of your Teamwork.com installation. For
  example, `https://<installation>.teamwork.com`.
- `TWAI_TEAMWORK_API_TOKEN`: The API token for your Teamwork.com account. You can
  generate a new API token in your Teamwork.com profile. For more information,
  check the [API documentation](https://apidocs.teamwork.com/guides/teamwork/authentication#basic-authentication).

Other optional environment variables are:
- `TWAI_PORT`: The port to run the MCP server. By default it will use a random
  available port in the machine. This will be used only when using the `sse`
  mode.
- `TWAI_LOG_LEVEL`: The log level for the MCP server. By default it will use
  `info`. Available log levels are `debug`, `info`, `warn` and `error`.

There are also some optional flags that you can use when running the MCP server:
- `-mode`: The mode to run the MCP server. It can be `stdio` or `sse`. By
  default it will use `sse`. For more information about the modes, check the
  [Usage section](#️usage), or the documentation
  [here](https://modelcontextprotocol.io/docs/concepts/transports#built-in-transport-types).

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

It assumes that [the binary](main.go) is compiled and installed in one of the
directories in your `PATH`. The `<installation>` is the name of your
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
| Comments          | ✅     | ✅       | ✅      | ❌     | Retrieve by task, milestone, notebook or file                   |

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