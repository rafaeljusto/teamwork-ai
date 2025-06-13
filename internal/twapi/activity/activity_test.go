package activity_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
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

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	tests := []struct {
		name     string
		multiple activity.Multiple
	}{{
		name: "all activities",
	}, {
		name: "activities for project",
		multiple: func() activity.Multiple {
			var multiple activity.Multiple
			multiple.Request.Path.ProjectID = resourceIDs.projectID
			return multiple
		}(),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get activities: %v", err)

			} else if len(tt.multiple.Response.Activities) == 0 {
				t.Error("expected at least one activity, got none")
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
