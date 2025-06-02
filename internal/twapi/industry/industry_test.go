package industry_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/industry"
)

const timeout = 5 * time.Second

var engine *twapi.Engine

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var multiple industry.Multiple
	if err := engine.Do(ctx, &multiple); err != nil {
		t.Errorf("failed to get industries: %v", err)

	} else if len(multiple.Response.Industries) == 0 {
		t.Error("expected at least one industry, got none")
	}
}

func startEngine() *twapi.Engine {
	server, token := os.Getenv("TWAI_TEAMWORK_SERVER"), os.Getenv("TWAI_TEAMWORK_API_TOKEN")
	if server == "" || token == "" {
		return nil
	}
	return twapi.NewEngine(server, token, nil)
}

func TestMain(m *testing.M) {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	engine = startEngine()
	if engine == nil {
		logger.Info("Missing setup environment variables, skipping tests")
		return
	}

	exitCode = m.Run()
}
