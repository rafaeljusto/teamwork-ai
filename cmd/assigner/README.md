# Assigner

The Assigner works as a webhook server that listens for incoming requests from
Teamwork.com. The server will analyze the incoming request and assign the tasks
to the users based on skill, job roles, user costs and workload. The server will
also generate a comment explaining the assignment.

For the user costs and workload analysis, it will be used a scoring approach.
Where cheaper and more available users will have a higher score. The top-scoring
user (or users if scores are equal) will be assigned to the task.

Keep in mind that the workload analysis has bigger weight than the user cost
analysis. This means that if a user has a high workload, they will be less
likely to be assigned to the task, even if they are cheaper than other users.

### 📦 Installing

You can install the Assigner server using [`go`](https://go.dev/doc/install):

```bash
go install -o teamwork-assigner github.com/rafaeljusto/teamwork-ai/cmd/assigner@latest
```

The binary will be installed in your `GOPATH/bin` directory, which is usually
`$HOME/go/bin`. Make sure to add this directory to your `PATH` environment
variable to run the `teamwork-assigner` command from anywhere.

Alternatively, you can use the pre-built binaries available in the
[releases](https://github.com/rafaeljusto/teamwork-ai/releases/latest) page.
Download the appropriate binary for your operating system, extract it, and place
it in a directory included in your `PATH`. For example, on Linux (amd64) using
`curl`:

```bash
# detect the latest release
twai_assigner_url=$(curl -s https://api.github.com/repos/rafaeljusto/teamwork-ai/releases/latest | \
  jq '.assets[] | select(.name | contains ("teamwork-assigner-linux-amd64")) | .browser_download_url')

# download the binary and place it in /usr/local/bin
sudo curl -s -O /usr/local/bin/teamwork-assigner ${twai_assigner_url}
```

The example above uses `jq` to parse the JSON response from the GitHub API,
which is a command-line JSON processor. You can find installation instructions
for it [here](https://jqlang.org/download/).

### ⚙️  Configuring

The following environment variables are required to run the Assigner server:
- `TWAI_TEAMWORK_SERVER`: The URL of your Teamwork.com installation. For
  example, `https://<installation>.teamwork.com`.
- `TWAI_TEAMWORK_API_TOKEN`: The Bearer token for your Teamwork.com account.For more information,
  check the [documentation](https://apidocs.teamwork.com/guides/teamwork/authentication#o-auth-2-0).
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
- `TWAI_MCP_ENDPOINT`: The endpoint of the MCP server to use for retrieving
  the prompt used to extract skills and job roles from the task information.

Other optional environment variables are:
- `TWAI_PORT`: The port to run the Assigner server. By default it will use a
  random available port in the machine.
- `TWAI_LOG_LEVEL`: The log level for the Assigner server. By default it will
  use `info`. Available log levels are `debug`, `info`, `warn` and `error`.

There are also some optional flags that you can use when running the Assigner
server:
- `skip-rates`: Skip user cost analysis when assigning the tasks. By default,
  the server will analyze the user rates and assign the tasks to the users with
  the lowest cost. If multiple users have the same cost, the server will assign
  to all selected users.
- `skip-workload`: Skip workload analysis when assigning the tasks. By default,
  the server will analyze the user workload and assign the tasks to the users
  with the lowest workload. If multiple users have the same workload, the server
  will assign to all selected users.
- `skip-assignment`: Skip the assignment of tasks to users. This is useful when
  you only need a suggestion from the AI as a comment instead of proactively
  assigning the tasks to users. By default, the server will assign the task.
- `skip-comment`: Skip the comment generation. This is useful when you only need
  to assign the tasks to users without any explanation. By default, the server
  will generate a comment explaining the assignment.

### ⚡️ Usage

After installing and configuring the Assigner server, you will need to run the
server and register it as a webhook in your Teamwork.com account. All the steps
to configure the webhook can be found
[here](https://apidocs.teamwork.com/guides/teamwork/setting-up-webhooks).

The webhhok URL can be associated with `TASK.CREATED` and `TASK.UPDATED` events.
At the moment there's no token or checksum check implemented in the server and
only version 2 of Teamwork.com webhooks is supported.

> [!IMPORTANT]
> Do not forget to add the URL path `/teamwork-ai/webhooks/task` to the webhook
> URL.

> [!TIP]
> For testing purposes we recommend using `ngrok` to expose a local running
> server to the Internet. Follow more information about `ngrok`
> [here](https://ngrok.com/docs/getting-started/).

### 📜 API

The Assigner server exposes a single endpoint to receive the incoming requests
from Teamwork.com. The endpoint is `/teamwork-ai/webhooks/task` and it accepts
`POST` requests. An example of a JSON payload that the server will receive:

```json
{
  "eventCreator": {
    "id": 160342,
    "firstName": "Rafael",
    "lastName": "Dantas Justo",
    "avatar": "https://s3.amazonaws.com/TWFiles/488712/userAvatar/tf_67324674-9202-4b8d-9957-454392c49faa.avatar.gif"
  },
  "project": {
    "id": 581677,
    "name": "Game Development Project",
    "description": "A comprehensive project to develop a new game with defined milestones, task lists, and deadlines to ensure completion within 2025.",
    "status": "active",
    "startDate": "2025-04-22",
    "endDate": "2025-12-31",
    "tags": [],
    "ownerId": 0,
    "companyId": 796953,
    "categoryId": 0,
    "dateCreated": "2025-04-22T10:55:38Z"
  },
  "task": {
    "id": 16367318,
    "name": "Write a Go program for the Game menu",
    "description": "The game needs a menu and it would be important to be written in the Go language.\n",
    "priority": null,
    "status": "new",
    "assignedUserIds": [],
    "parentId": 0,
    "taskListId": 1404538,
    "startDate": null,
    "dueDate": null,
    "progress": 0,
    "estimatedMinutes": 0,
    "tags": [],
    "projectId": 581677,
    "dateCreated": "2025-05-01T17:34:45Z",
    "dateUpdated": "2025-05-01T17:34:45Z",
    "hasCustomFields": false
  },
  "taskList": {
    "id": 1404538,
    "name": "General tasks",
    "description": "",
    "status": "new",
    "milestoneId": 0,
    "projectId": 581677,
    "templateId": null,
    "tags": []
  },
  "users": []
}
```