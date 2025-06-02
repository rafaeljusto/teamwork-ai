package tag_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/project"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/tag"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		projectID int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := tag.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tagID int64
	tagIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		tagID = i
	})
	if err := engine.Do(ctx, &create, tagIDSetter); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tagDelete tag.Delete
		tagDelete.Request.Path.ID = tagID
		if err := engine.Do(ctx, &tagDelete); err != nil {
			t.Logf("⚠️  failed to delete tag: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single tag.Single
	single.ID = tagID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get tag: %v", err)
	}
	if single.ID != tagID {
		t.Errorf("expected tag ID %d, got %d", tagID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := tag.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tagID int64
	tagIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		tagID = i
	})
	if err := engine.Do(ctx, &create, tagIDSetter); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tagDelete tag.Delete
		tagDelete.Request.Path.ID = tagID
		if err := engine.Do(ctx, &tagDelete); err != nil {
			t.Logf("⚠️  failed to delete tag: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple tag.Multiple
	}{{
		name: "all tags",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get tags: %v", err)

			} else if len(tt.multiple.Response.Tags) == 0 {
				t.Error("expected at least one tag, got none")
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
		create tag.Create
	}{{
		name: "only required fields",
		create: tag.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields",
		create: tag.Create{
			Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			ProjectID: &resourceIDs.projectID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var tagID int64
			tagIDSetter := teamwork.WithIDCallback("id", func(id int64) {
				tagID = id
			})

			if err := engine.Do(ctx, &tt.create, tagIDSetter); err != nil {
				t.Errorf("failed to create tag: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var tagDelete tag.Delete
					tagDelete.Request.Path.ID = tagID
					if err := engine.Do(ctx, &tagDelete); err != nil {
						t.Logf("⚠️  failed to delete tag: %v", err)
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

	create := tag.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tagID int64
	tagIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		tagID = i
	})
	if err := engine.Do(ctx, &create, tagIDSetter); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tagDelete tag.Delete
		tagDelete.Request.Path.ID = tagID
		if err := engine.Do(ctx, &tagDelete); err != nil {
			t.Logf("⚠️  failed to delete tag: %v", err)
		}
	})

	tests := []struct {
		name   string
		create tag.Update
	}{{
		name: "all fields",
		create: tag.Update{
			ID:        tagID,
			Name:      teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			ProjectID: &resourceIDs.projectID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update tag: %v", err)
			}
		})
	}
}

func createProject(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	projectCreate := project.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	projectIDSetter := teamwork.WithIDCallback("id", func(id int64) {
		resourceIDs.projectID = id
	})

	logger.Info("⚙️  Creating project")
	if err := engine.Do(ctx, &projectCreate, projectIDSetter); err != nil {
		logger.Error("failed to create project",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created project",
		slog.Int64("id", resourceIDs.projectID),
		slog.String("name", projectCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up project",
			slog.Int64("id", resourceIDs.projectID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var projectDelete project.Delete
		projectDelete.Request.Path.ID = resourceIDs.projectID
		if err := engine.Do(ctx, &projectDelete); err != nil {
			logger.Warn("⚠️  failed to delete project",
				slog.Int64("id", resourceIDs.projectID),
				slog.String("error", err.Error()),
			)
		}
	}
}

func startEngine() *teamwork.Engine {
	server, token := os.Getenv("TWAI_TEAMWORK_SERVER"), os.Getenv("TWAI_TEAMWORK_API_TOKEN")
	if server == "" || token == "" {
		return nil
	}
	return teamwork.NewEngine(server, token, nil)
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

	deleteProject := createProject(logger)
	if resourceIDs.projectID == 0 {
		exitCode = 1
		return
	}
	defer deleteProject()

	reference := time.Now()
	defer func() {
		if diff := time.Since(reference); diff < 200*time.Millisecond {
			time.Sleep(200*time.Millisecond - diff) // ensure tests have enough time to sync
		}
	}()

	exitCode = m.Run()
}
