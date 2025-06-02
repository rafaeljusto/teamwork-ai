package timer_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/project"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/task"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/tasklist"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/timer"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/user"
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

	var create timer.Create

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single timer.Single
	single.ID = timerID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get timer: %v", err)
	}
	if single.ID != timerID {
		t.Errorf("expected timer ID %d, got %d", timerID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	var create timer.Create

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple timer.Multiple
	}{{
		name: "all timers",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get timers: %v", err)

			} else if len(tt.multiple.Response.Timers) == 0 {
				t.Error("expected at least one timer, got none")
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
		create timer.Create
	}{{
		name: "only required fields",
	}, {
		name: "all fields",
		create: timer.Create{
			Description:       twapi.Ref("This is a test timer"),
			Billable:          twapi.Ref(true),
			Running:           twapi.Ref(true),
			Seconds:           twapi.Ref(int64(3600)), // 1 hour in seconds
			StopRunningTimers: twapi.Ref(true),
			ProjectID:         &resourceIDs.projectID,
			TaskID:            &resourceIDs.taskID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var timerID int64
			timerIDSetter := twapi.WithIDCallback("id", func(id int64) {
				timerID = id
			})

			if err := engine.Do(ctx, &tt.create, timerIDSetter); err != nil {
				t.Errorf("failed to create timer: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var timerDelete timer.Delete
					timerDelete.Request.Path.ID = timerID
					if err := engine.Do(ctx, &timerDelete); err != nil {
						t.Logf("⚠️  failed to delete timer: %v", err)
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

	var create timer.Create

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	tests := []struct {
		name   string
		create timer.Update
	}{{
		name: "all fields",
		create: timer.Update{
			ID:          timerID,
			Description: twapi.Ref("Updated description"),
			Billable:    twapi.Ref(true),
			Running:     twapi.Ref(true),
			ProjectID:   &resourceIDs.projectID,
			TaskID:      &resourceIDs.taskID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update timer: %v", err)
			}
		})
	}
}

func TestPause(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	var create timer.Create

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var pause timer.Pause
	pause.Request.Path.ID = timerID
	if err := engine.Do(ctx, &pause); err != nil {
		t.Errorf("failed to pause timer: %v", err)
	}
}

func TestComplete(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	var create timer.Create
	create.Seconds = twapi.Ref(int64(3600)) // 1 hour in seconds
	create.ProjectID = &resourceIDs.projectID
	create.TaskID = &resourceIDs.taskID

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var complete timer.Complete
	complete.Request.Path.ID = timerID
	if err := engine.Do(ctx, &complete); err != nil {
		t.Errorf("failed to complete timer: %v", err)
	}
}

func TestResume(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	var create timer.Create

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timerID int64
	timerIDSetter := twapi.WithIDCallback("id", func(i int64) {
		timerID = i
	})
	if err := engine.Do(ctx, &create, timerIDSetter); err != nil {
		t.Fatalf("failed to create timer: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timerDelete timer.Delete
		timerDelete.Request.Path.ID = timerID
		if err := engine.Do(ctx, &timerDelete); err != nil {
			t.Logf("⚠️  failed to delete timer: %v", err)
		}
	})

	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var pause timer.Pause
	pause.Request.Path.ID = timerID
	if err := engine.Do(ctx, &pause); err != nil {
		t.Errorf("failed to pause timer: %v", err)
	}

	var resume timer.Resume
	resume.Request.Path.ID = timerID
	if err := engine.Do(ctx, &resume); err != nil {
		t.Errorf("failed to resume timer: %v", err)
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

func addLoggedUserAsProjectMember(logger *slog.Logger) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var me user.Me
	if err := engine.Do(ctx, &me); err != nil {
		logger.Error("failed to get current user",
			slog.String("error", err.Error()),
		)
		return
	}

	var addProject user.AddProject
	addProject.Request.Path.ProjectID = resourceIDs.projectID
	addProject.Request.Users.IDs = []int64{me.ID}

	logger.Info("⚙️  Adding user to project")
	if err := engine.Do(ctx, &addProject); err != nil {
		logger.Error("failed to add user to project",
			slog.Int64("userID", me.ID),
			slog.Int64("projectID", resourceIDs.projectID),
			slog.String("error", err.Error()),
		)
	}
	logger.Info("✅ Added user to project")
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

	addLoggedUserAsProjectMember(logger)

	reference := time.Now()
	defer func() {
		if diff := time.Since(reference); diff < 200*time.Millisecond {
			time.Sleep(200*time.Millisecond - diff) // ensure tests have enough time to sync
		}
	}()

	exitCode = m.Run()
}
