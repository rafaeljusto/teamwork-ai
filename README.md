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

The MCP server was moved to the official Teamwork.com repository:
https://github.com/teamwork/mcp

## Assigner

The assigner is a webservice that integrates with Teamwork.com webhooks,
handling task creation and updates. It behave as an Agentic AI, extracting
skills and job roles from the task and assigning the best user to fulfill it.

[![Assigner](https://img.youtube.com/vi/syeb50mia_M/0.jpg)](https://www.youtube.com/watch?v=syeb50mia_M)

**For more information check [our documentation](cmd/assigner/README.md).**