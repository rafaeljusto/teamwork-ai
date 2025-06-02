package comment_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/comment"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/project"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/task"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/tasklist"
)

const timeout = 5 * time.Second

var (
	engine      *twapi.Engine
	resourceIDs struct {
		projectID  int64
		tasklistID int64
		taskID     int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := comment.Create{
		Object:      twapi.Relationship{ID: resourceIDs.taskID, Type: "tasks"},
		Body:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ContentType: twapi.Ref("TEXT"),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var commentID int64
	commentIDSetter := twapi.WithIDCallback("id", func(i int64) {
		commentID = i
	})
	if err := engine.Do(ctx, &create, commentIDSetter); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var commentDelete comment.Delete
		commentDelete.Request.Path.ID = commentID
		if err := engine.Do(ctx, &commentDelete); err != nil {
			t.Logf("⚠️  failed to delete comment: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single comment.Single
	single.ID = commentID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get comment: %v", err)
	}
	if single.ID != commentID {
		t.Errorf("expected comment ID %d, got %d", commentID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := comment.Create{
		Object:      twapi.Relationship{ID: resourceIDs.taskID, Type: "tasks"},
		Body:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ContentType: twapi.Ref("TEXT"),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var commentID int64
	commentIDSetter := twapi.WithIDCallback("id", func(i int64) {
		commentID = i
	})
	if err := engine.Do(ctx, &create, commentIDSetter); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var commentDelete comment.Delete
		commentDelete.Request.Path.ID = commentID
		if err := engine.Do(ctx, &commentDelete); err != nil {
			t.Logf("⚠️  failed to delete comment: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple comment.Multiple
	}{{
		name: "all comments",
	}, {
		name: "comments for task",
		multiple: func() comment.Multiple {
			var multiple comment.Multiple
			multiple.Request.Path.TaskID = resourceIDs.taskID
			return multiple
		}(),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get comments: %v", err)

			} else if len(tt.multiple.Response.Comments) == 0 {
				t.Error("expected at least one comment, got none")
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
		create comment.Create
	}{{
		name: "only required fields",
		create: comment.Create{
			Object:      twapi.Relationship{ID: resourceIDs.taskID, Type: "tasks"},
			Body:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			ContentType: twapi.Ref("TEXT"),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var commentID int64
			commentIDSetter := twapi.WithIDCallback("id", func(id int64) {
				commentID = id
			})

			if err := engine.Do(ctx, &tt.create, commentIDSetter); err != nil {
				t.Errorf("failed to create comment: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var commentDelete comment.Delete
					commentDelete.Request.Path.ID = commentID
					if err := engine.Do(ctx, &commentDelete); err != nil {
						t.Logf("⚠️  failed to delete comment: %v", err)
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

	create := comment.Create{
		Object:      twapi.Relationship{ID: resourceIDs.taskID, Type: "tasks"},
		Body:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ContentType: twapi.Ref("TEXT"),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var commentID int64
	commentIDSetter := twapi.WithIDCallback("id", func(i int64) {
		commentID = i
	})
	if err := engine.Do(ctx, &create, commentIDSetter); err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var commentDelete comment.Delete
		commentDelete.Request.Path.ID = commentID
		if err := engine.Do(ctx, &commentDelete); err != nil {
			t.Logf("⚠️  failed to delete comment: %v", err)
		}
	})

	tests := []struct {
		name   string
		create comment.Update
	}{{
		name: "all fields",
		create: comment.Update{
			ID:          commentID,
			Body:        "<h1>test</h1>",
			ContentType: twapi.Ref("HTML"),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update comment: %v", err)
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

	projectIDSetter := twapi.WithIDCallback("id", func(id int64) {
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

func createTasklist(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tasklistCreate := tasklist.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ProjectID: resourceIDs.projectID,
	}

	tasklistIDSetter := twapi.WithIDCallback("tasklistId", func(id int64) {
		resourceIDs.tasklistID = id
	})

	logger.Info("⚙️  Creating tasklist")
	if err := engine.Do(ctx, &tasklistCreate, tasklistIDSetter); err != nil {
		logger.Error("failed to create tasklist",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created tasklist",
		slog.Int64("id", resourceIDs.tasklistID),
		slog.String("name", tasklistCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up tasklist",
			slog.Int64("id", resourceIDs.tasklistID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tasklistDelete tasklist.Delete
		tasklistDelete.Request.Path.ID = resourceIDs.tasklistID
		if err := engine.Do(ctx, &tasklistDelete); err != nil {
			logger.Warn("⚠️  failed to delete tasklist",
				slog.Int64("id", resourceIDs.tasklistID),
				slog.String("error", err.Error()),
			)
		}
	}
}

func createTask(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	taskCreate := task.Create{
		Name:       fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		TasklistID: resourceIDs.tasklistID,
	}

	taskIDSetter := twapi.WithIDCallback("id", func(id int64) {
		resourceIDs.taskID = id
	})

	logger.Info("⚙️  Creating task")
	if err := engine.Do(ctx, &taskCreate, taskIDSetter); err != nil {
		logger.Error("failed to create task",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created task",
		slog.Int64("id", resourceIDs.taskID),
		slog.String("name", taskCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up task",
			slog.Int64("id", resourceIDs.taskID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var taskDelete task.Delete
		taskDelete.Request.Path.ID = resourceIDs.taskID
		if err := engine.Do(ctx, &taskDelete); err != nil {
			logger.Warn("⚠️  failed to delete task",
				slog.Int64("id", resourceIDs.taskID),
				slog.String("error", err.Error()),
			)
		}
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

	deleteProject := createProject(logger)
	if resourceIDs.projectID == 0 {
		exitCode = 1
		return
	}
	defer deleteProject()

	deleteTasklist := createTasklist(logger)
	if resourceIDs.tasklistID == 0 {
		exitCode = 1
		return
	}
	defer deleteTasklist()

	deleteTask := createTask(logger)
	if resourceIDs.taskID == 0 {
		exitCode = 1
		return
	}
	defer deleteTask()

	reference := time.Now()
	defer func() {
		if diff := time.Since(reference); diff < 200*time.Millisecond {
			time.Sleep(200*time.Millisecond - diff) // ensure tests have enough time to sync
		}
	}()

	exitCode = m.Run()
}
