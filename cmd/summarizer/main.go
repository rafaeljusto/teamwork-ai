// Package main is a command line tool to summarize installation or project activities during a period of time.
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/agentic/actions"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/anthropic"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/ollama"
	_ "github.com/rafaeljusto/teamwork-ai/internal/agentic/openai"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
)

func main() {
	defer handleExit()

	startDateStr := flag.String("start-date", "", "start date in YYYY-MM-DD format")
	endDateStr := flag.String("end-date", "", "end date in YYYY-MM-DD format")
	projectID := flag.Int64("project-id", 0, "project ID to summarize")
	flag.Parse()

	// We are using a logger to print the errors because we don't have a logger
	// yet. We could use the standard logger, but it's better to create a logger
	// with the correct configuration.
	preLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	var setupFailed bool
	if *startDateStr == "" || *endDateStr == "" {
		preLogger.Error("start-date and end-date are required")
		setupFailed = true
	}
	startDate, err := time.Parse("2006-01-02", *startDateStr)
	if err != nil {
		preLogger.Error("failed to parse start-date",
			slog.String("start-date", *startDateStr),
			slog.String("error", err.Error()),
		)
		setupFailed = true
	}
	endDate, err := time.Parse("2006-01-02", *endDateStr)
	if err != nil {
		preLogger.Error("failed to parse end-date",
			slog.String("end-date", *endDateStr),
			slog.String("error", err.Error()),
		)
		setupFailed = true
	}
	if *projectID < 0 {
		preLogger.Error("project-id should be a non-negative integer")
		setupFailed = true
	}
	if setupFailed {
		exit(exitCodeInvalidInput)
	}

	ctx := context.Background()

	c, errs := config.ParseFromEnvs()
	if errs != nil {
		for _, err := range multierr(errs) {
			preLogger.Error("failed to parse configuration",
				slog.String("error", err.Error()),
			)
		}
		exit(exitCodeInvalidInput)
	}

	resources, err := config.InitResources(ctx, c)
	if err != nil {
		resources.Logger.Error("failed to initialize resources",
			slog.String("error", err.Error()),
		)
		exit(exitCodeSetupFailure)
	}

	summary, err := actions.SummarizeActivities(ctx, resources,
		actions.WithSummarizeActivitiesPeriod(startDate, endDate),
		actions.WithSummarizeActivitiesProjectID(*projectID),
	)
	if err != nil {
		resources.Logger.Error("failed to summarize activities",
			slog.String("error", err.Error()),
		)
		exit(exitCodeInternalError)
	}

	if summary == "" {
		resources.Logger.Info("no activities found for the specified period",
			slog.String("start-date", startDate.Format("2006-01-02")),
			slog.String("end-date", endDate.Format("2006-01-02")),
			slog.Int64("project-id", *projectID),
		)
	} else {
		resources.Logger.Info("activities summary",
			slog.String("summary", summary),
			slog.String("start-date", startDate.Format("2006-01-02")),
			slog.String("end-date", endDate.Format("2006-01-02")),
			slog.Int64("project-id", *projectID),
		)
	}
}

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeInvalidInput
	exitCodeSetupFailure
	exitCodeInternalError
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
