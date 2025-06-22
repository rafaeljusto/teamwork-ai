# Summarizer

The summrizer is a command-line interface (CLI) tool that summarizes activities
from a Teamwork.com account in a date period. It can optionally target a
specific project.

### 📦 Installing

You can install the Summarizer using [`go`](https://go.dev/doc/install):

```bash
go install -o teamwork-summarizer github.com/rafaeljusto/teamwork-ai/cmd/summarizer@latest
```

The binary will be installed in your `GOPATH/bin` directory, which is usually
`$HOME/go/bin`. Make sure to add this directory to your `PATH` environment
variable to run the `teamwork-summarizer` command from anywhere.

Alternatively, you can use the pre-built binaries available in the
[releases](https://github.com/rafaeljusto/teamwork-ai/releases/latest) page.
Download the appropriate binary for your operating system, extract it, and place
it in a directory included in your `PATH`. For example, on Linux (amd64) using
`curl`:

```bash
# detect the latest release
twai_assigner_url=$(curl -s https://api.github.com/repos/rafaeljusto/teamwork-ai/releases/latest | \
  jq '.assets[] | select(.name | contains ("teamwork-summarizer-linux-amd64")) | .browser_download_url')

# download the binary and place it in /usr/local/bin
sudo curl -s -O /usr/local/bin/teamwork-summarizer ${twai_assigner_url}
```

The example above uses `jq` to parse the JSON response from the GitHub API,
which is a command-line JSON processor. You can find installation instructions
for it [here](https://jqlang.org/download/).

### ⚙️  Configuring

The following environment variables are required to run the Summarizer tool:
- `TWAI_TEAMWORK_SERVER`: The URL of your Teamwork.com installation. For
  example, `https://<installation>.teamwork.com`.
- `TWAI_TEAMWORK_API_TOKEN`: The API token for your Teamwork.com account. You can
  generate a new API token in your Teamwork.com profile. For more information,
  check the [API documentation](https://apidocs.teamwork.com/guides/teamwork/authentication#basic-authentication).
- `TWAI_AGENTIC_NAME`: The name of the agent that will be used to extract
  information from the task. The possible values are `anthropic`, `openai` and
  `ollama`.
- `TWAI_AGENTIC_DSN`: The connection string for the agentic model. The format of
  the connection string depends on the agentic name:
  * `anthropic`: `model:token`. Where `model` is the name of the model (e.g.,
    `claude-2`) and `token` is the API key for the Anthropic account. All
    available models can be found [here](https://docs.anthropic.com/en/docs/about-claude/models/all-models).
  * `openai`: `model:token`. Where `model` is the name of the model (e.g.,
    `gpt-3.5-turbo`) and `token` is the API key for the OpenAI account. All
    available models can be found [here](https://platform.openai.com/docs/models).
  * `ollama`: `http[s]://[username[:password]@]host[:port]/model`. Where
    `username` and `password` are the credentials for the Ollama account, `host`
    is the host name or IP address of the Ollama server, and `port` is the port
    number of the Ollama server. The `model` is the name of the model to use
    (e.g. `llama2`), all available models can be found
    [here](https://ollama.com/search).
- `TWAI_AGENTIC_MCP_CLIENT_STDIO_PATH` or `TWAI_AGENTIC_MCP_CLIENT_SSE_ADDRESS`:
  Depending on the MCP server mode one environment variable or the other must be
  set. While the `TWAI_AGENTIC_MCP_CLIENT_STDIO_PATH` should be the path to the
  MCP binary, the `TWAI_AGENTIC_MCP_CLIENT_SSE_ADDRESS` should be the address of
  the MCP server (please check the MCP server [here](../mcp/)).
  * `TWAI_AGENTIC_MCP_CLIENT_STDIO_ARGS`: Optional arguments to pass to the MCP
    client binary. The format is `ARG1,ARG2` and multiple arguments can be
    separated by a comma.
  * `TWAI_AGENTIC_MCP_CLIENT_STDIO_ENVS`: Optional environment variables to pass
    to the MCP client binary. The format is `KEY=VALUE` and multiple variables
    can be separated by a comma.


Other optional environment variables are:
- `TWAI_LOG_LEVEL`: The log level for the Summarizer server. By default it will
  use `info`. Available log levels are `debug`, `info`, `warn` and `error`.

There are also some command-line arguments that you can use when running the
Summarizer tool, some of them are required:
- `start-date`: The start date of the period to summarize the activities. The
  format is `YYYY-MM-DD`. This argument is required.
- `end-date`: The end date of the period to summarize the activities. The format
  is `YYYY-MM-DD`. This argument is required.
- `project-id`: The ID of the project to summarize the activities. If not
  provided, all activities from all projects will be summarized.