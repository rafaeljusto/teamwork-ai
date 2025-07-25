// Package main is a microservice to expose Teamwork.com operations to LLMs
// using the Model Context Protocol.
package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	mcpactivity "github.com/rafaeljusto/teamwork-ai/internal/mcp/activity"
	mcpcomment "github.com/rafaeljusto/teamwork-ai/internal/mcp/comment"
	mcpcompany "github.com/rafaeljusto/teamwork-ai/internal/mcp/company"
	mcpindustry "github.com/rafaeljusto/teamwork-ai/internal/mcp/industry"
	mcpjobrole "github.com/rafaeljusto/teamwork-ai/internal/mcp/jobrole"
	mcpmilestone "github.com/rafaeljusto/teamwork-ai/internal/mcp/milestone"
	mcpproject "github.com/rafaeljusto/teamwork-ai/internal/mcp/project"
	mcpskill "github.com/rafaeljusto/teamwork-ai/internal/mcp/skill"
	mcptag "github.com/rafaeljusto/teamwork-ai/internal/mcp/tag"
	mcptask "github.com/rafaeljusto/teamwork-ai/internal/mcp/task"
	mcptasklist "github.com/rafaeljusto/teamwork-ai/internal/mcp/tasklist"
	mcptimelog "github.com/rafaeljusto/teamwork-ai/internal/mcp/timelog"
	mcptimer "github.com/rafaeljusto/teamwork-ai/internal/mcp/timer"
	mcpuser "github.com/rafaeljusto/teamwork-ai/internal/mcp/user"
)

const (
	mcpName    = "Teamwork AI"
	mcpVersion = "1.0.0"
)

func main() {
	defer handleExit()

	serverMode := flag.String("mode", "sse", "server mode")
	flag.Parse()

	c, errs := config.ParseFromEnvs()
	if errs != nil {
		// We are using a logger to print the errors because we don't have a
		// logger yet. We could use the standard logger, but it's better to
		// create a logger with the correct configuration.
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
		for _, err := range multierr(errs) {
			logger.Error("failed to parse configuration",
				slog.String("error", err.Error()),
			)
		}
		exit(exitCodeInvalidInput)
	}
	resources := config.NewResources(c)

	mcpServer := newMCPServer(resources)
	switch *serverMode {
	case "stdio":
		stdioServer := server.NewStdioServer(mcpServer)
		if err := stdioServer.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
			resources.Logger.Error("failed to serve",
				slog.String("error", err.Error()),
			)
			exit(exitCodeSetupFailure)
		}

	case "sse":
		sseServerAddress := ":" + strconv.FormatInt(c.Port, 10)
		resources.Logger.Info("starting http server",
			slog.String("address", sseServerAddress),
		)

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		sseServer := server.NewSSEServer(mcpServer)
		go func() {
			if err := sseServer.Start(sseServerAddress); err != nil {
				if err != http.ErrServerClosed {
					resources.Logger.Error("failed to serve",
						slog.String("error", err.Error()),
					)
					select {
					case <-done:
					default:
						close(done)
					}
				}
			}
		}()

		<-done
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()
		if err := sseServer.Shutdown(ctx); err != nil {
			resources.Logger.Error("server shutdown failed",
				slog.String("error", err.Error()),
			)
		}
		resources.Logger.Info("server stopped")
	}
}

func newMCPServer(resources *config.Resources) *server.MCPServer {
	mcpServer := server.NewMCPServer(mcpName, mcpVersion,
		server.WithLogging(),
	)

	mcptask.Register(mcpServer, resources)
	mcptasklist.Register(mcpServer, resources)
	mcpproject.Register(mcpServer, resources)
	mcpuser.Register(mcpServer, resources)
	mcpskill.Register(mcpServer, resources)
	mcpcompany.Register(mcpServer, resources)
	mcpindustry.Register(mcpServer, resources)
	mcptag.Register(mcpServer, resources)
	mcpmilestone.Register(mcpServer, resources)
	mcpjobrole.Register(mcpServer, resources)
	mcpcomment.Register(mcpServer, resources)
	mcptimelog.Register(mcpServer, resources)
	mcptimer.Register(mcpServer, resources)
	mcpactivity.Register(mcpServer, resources)

	return mcpServer
}

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeInvalidInput
	exitCodeSetupFailure
)

type exitData struct {
	code exitCode
}

// exit allows to abort the program while still executing all defer statements.
func exit(code exitCode) {
	panic(exitData{code: code})
}

// handleExit exit code handler.
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exitData); ok {
			os.Exit(int(exit.code))
		}
		panic(e)
	}
}

// multierr unwraps multiple errors from a single error.
//
// https://pkg.go.dev/errors#Join
func multierr(errs error) []error {
	if errs == nil {
		return nil
	}
	if multierr, ok := errs.(interface{ Unwrap() []error }); ok {
		return multierr.Unwrap()
	}
	return []error{errs}
}
