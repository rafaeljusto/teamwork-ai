// Package main is an HTTP server that reacts to Teamwork.com webhooks and
// assigns users to tasks based on AI and workload decisions.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic/actions"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/anthropic"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/ollama"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/openai"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
)

var (
	skipRates      bool
	skipWorkload   bool
	skipAssignment bool
	skipComment    bool
)

func main() {
	defer handleExit()

	flag.BoolVar(&skipRates, "skip-rates", false, "Skip rate analysis when assigning a task")
	flag.BoolVar(&skipWorkload, "skip-workload", false, "Skip workload analysis when assigning a task")
	flag.BoolVar(&skipAssignment, "skip-assignment", false, "Skip task assignment (only comment)")
	flag.BoolVar(&skipComment, "skip-comment", false, "Skip task comment (only assign)")
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

	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(c.Port, 10))
	if err != nil {
		resources.Logger.Error("failed to listen",
			slog.String("error", err.Error()),
		)
		exit(exitCodeSetupFailure)
	}

	resources.Logger.Info("starting web server",
		slog.String("address", listener.Addr().String()),
	)

	router := http.NewServeMux()
	router.HandleFunc("POST /teamwork-ai/webhooks/task", handleTask(resources))

	server := http.Server{
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				resources.Logger.Error("failed to serve",
					slog.String("error", err.Error()),
				)
			}
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		resources.Logger.Error("server shutdown failed",
			slog.String("error", err.Error()),
		)
	}
	resources.Logger.Info("server stopped")
}

func handleTask(resources *config.Resources) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var taskData webhook.TaskData
		if err := decoder.Decode(&taskData); err != nil {
			resources.Logger.Error("failed to decode request body",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		var options []actions.AutoAssignTaskOption
		if skipRates {
			options = append(options, actions.WithAutoAssignTaskSkipRates())
		}
		if skipWorkload {
			options = append(options, actions.WithAutoAssignTaskSkipWorkload())
		}
		if skipAssignment {
			options = append(options, actions.WithAutoAssignTaskSkipAssignment())
		}
		if skipComment {
			options = append(options, actions.WithAutoAssignTaskSkipComment())
		}

		if err := actions.AutoAssignTask(r.Context(), resources, taskData, options...); err != nil {
			resources.Logger.Error("failed to auto assign task",
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to auto assign task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
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
