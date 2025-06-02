package jobrole_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/jobrole"
)

const timeout = 5 * time.Second

var (
	engine *twapi.Engine
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := jobrole.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var jobroleID int64
	jobroleIDSetter := twapi.WithIDCallback("id", func(i int64) {
		jobroleID = i
	})
	if err := engine.Do(ctx, &create, jobroleIDSetter); err != nil {
		t.Fatalf("failed to create job role: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var jobroleDelete jobrole.Delete
		jobroleDelete.Request.Path.ID = jobroleID
		if err := engine.Do(ctx, &jobroleDelete); err != nil {
			t.Logf("⚠️  failed to delete job role: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single jobrole.Single
	single.ID = jobroleID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get job role: %v", err)
	}
	if single.ID != jobroleID {
		t.Errorf("expected job role ID %d, got %d", jobroleID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := jobrole.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var jobroleID int64
	jobroleIDSetter := twapi.WithIDCallback("id", func(i int64) {
		jobroleID = i
	})
	if err := engine.Do(ctx, &create, jobroleIDSetter); err != nil {
		t.Fatalf("failed to create job role: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var jobroleDelete jobrole.Delete
		jobroleDelete.Request.Path.ID = jobroleID
		if err := engine.Do(ctx, &jobroleDelete); err != nil {
			t.Logf("⚠️  failed to delete job role: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple jobrole.Multiple
	}{{
		name: "all job roles",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get job roles: %v", err)

			} else if len(tt.multiple.Response.JobRoles) == 0 {
				t.Error("expected at least one job role, got none")
			}
		})
	}
}

func TestCreate(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	tests := []struct {
		name   string
		create jobrole.Create
	}{{
		name: "only required fields",
		create: jobrole.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var jobroleID int64
			jobroleIDSetter := twapi.WithIDCallback("id", func(id int64) {
				jobroleID = id
			})

			if err := engine.Do(ctx, &tt.create, jobroleIDSetter); err != nil {
				t.Errorf("failed to create job role: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var jobroleDelete jobrole.Delete
					jobroleDelete.Request.Path.ID = jobroleID
					if err := engine.Do(ctx, &jobroleDelete); err != nil {
						t.Logf("⚠️  failed to delete job role: %v", err)
					}
				})
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := jobrole.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var jobroleID int64
	jobroleIDSetter := twapi.WithIDCallback("id", func(i int64) {
		jobroleID = i
	})
	if err := engine.Do(ctx, &create, jobroleIDSetter); err != nil {
		t.Fatalf("failed to create job role: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var jobroleDelete jobrole.Delete
		jobroleDelete.Request.Path.ID = jobroleID
		if err := engine.Do(ctx, &jobroleDelete); err != nil {
			t.Logf("⚠️  failed to delete job role: %v", err)
		}
	})

	tests := []struct {
		name   string
		create jobrole.Update
	}{{
		name: "all fields",
		create: jobrole.Update{
			ID:   jobroleID,
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update job role: %v", err)
			}
		})
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
