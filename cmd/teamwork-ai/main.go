// Package main is a microservice to expose Teamwork.com operations to LLMs
// using the Model Context Protocol.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	mcptask "github.com/rafaeljusto/teamwork-ai/internal/mcp/task"
)

const (
	mcpName    = "Teamwork AI"
	mcpVersion = "1.0.0"
)

func main() {
	defer handleExit()

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

	mcpServerAddress := ":" + strconv.FormatInt(c.Port, 10)
	resources.Logger.Info("starting web server",
		slog.String("address", mcpServerAddress),
	)

	mcpServer := server.NewMCPServer(mcpName, mcpVersion,
		server.WithLogging(),
	)
	mcptask.Register(mcpServer, resources)

	sseServer := server.NewSSEServer(mcpServer)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := sseServer.Start(mcpServerAddress); err != nil {
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
