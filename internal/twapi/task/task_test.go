package task_test

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
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/tag"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/task"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/tasklist"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/user"
)

const timeout = 5 * time.Second

var (
	engine      *twapi.Engine
	resourceIDs struct {
		projectID  int64
		tasklistID int64
		tagID      int64
		userID     int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := task.Create{
		Name:       fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		TasklistID: resourceIDs.tasklistID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var taskID int64
	taskIDSetter := twapi.WithIDCallback("id", func(i int64) {
		taskID = i
	})
	if err := engine.Do(ctx, &create, taskIDSetter); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var taskDelete task.Delete
		taskDelete.Request.Path.ID = taskID
		if err := engine.Do(ctx, &taskDelete); err != nil {
			t.Logf("⚠️  failed to delete task: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single task.Single
	single.ID = taskID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get task: %v", err)
	}
	if single.ID != taskID {
		t.Errorf("expected task ID %d, got %d", taskID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := task.Create{
		Name:       fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		TasklistID: resourceIDs.tasklistID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var taskID int64
	taskIDSetter := twapi.WithIDCallback("id", func(i int64) {
		taskID = i
	})
	if err := engine.Do(ctx, &create, taskIDSetter); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var taskDelete task.Delete
		taskDelete.Request.Path.ID = taskID
		if err := engine.Do(ctx, &taskDelete); err != nil {
			t.Logf("⚠️  failed to delete task: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple task.Multiple
	}{{
		name: "all tasks",
	}, {
		name: "tasks for project",
		multiple: func() task.Multiple {
			var multiple task.Multiple
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
				t.Errorf("failed to get tasks: %v", err)

			} else if len(tt.multiple.Response.Tasks) == 0 {
				t.Error("expected at least one task, got none")
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
		create task.Create
	}{{
		name: "only required fields",
		create: task.Create{
			Name:       fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			TasklistID: resourceIDs.tasklistID,
		},
	}, {
		name: "all fields",
		create: task.Create{
			Name:             fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Description:      twapi.Ref("This is a test task"),
			Priority:         twapi.Ref("high"),
			Progress:         twapi.Ref(int64(50)),
			StartAt:          twapi.Ref(twapi.Date(time.Now().Add(24 * time.Hour))),
			DueAt:            twapi.Ref(twapi.Date(time.Now().Add(48 * time.Hour))),
			EstimatedMinutes: twapi.Ref(int64(120)),
			TasklistID:       resourceIDs.tasklistID,
			Assignees: &twapi.UserGroups{
				UserIDs: []int64{resourceIDs.userID},
			},
			TagIDs: []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var taskID int64
			taskIDSetter := twapi.WithIDCallback("id", func(id int64) {
				taskID = id
			})

			if err := engine.Do(ctx, &tt.create, taskIDSetter); err != nil {
				t.Errorf("failed to create task: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var taskDelete task.Delete
					taskDelete.Request.Path.ID = taskID
					if err := engine.Do(ctx, &taskDelete); err != nil {
						t.Logf("⚠️  failed to delete task: %v", err)
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

	create := task.Create{
		Name:       fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		TasklistID: resourceIDs.tasklistID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var taskID int64
	taskIDSetter := twapi.WithIDCallback("id", func(i int64) {
		taskID = i
	})
	if err := engine.Do(ctx, &create, taskIDSetter); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var taskDelete task.Delete
		taskDelete.Request.Path.ID = taskID
		if err := engine.Do(ctx, &taskDelete); err != nil {
			t.Logf("⚠️  failed to delete task: %v", err)
		}
	})

	tests := []struct {
		name   string
		create task.Update
	}{{
		name: "all fields",
		create: task.Update{
			ID:               taskID,
			Name:             twapi.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description:      twapi.Ref("This is a test task"),
			Priority:         twapi.Ref("high"),
			Progress:         twapi.Ref(int64(50)),
			StartAt:          twapi.Ref(twapi.Date(time.Now().Add(24 * time.Hour))),
			DueAt:            twapi.Ref(twapi.Date(time.Now().Add(48 * time.Hour))),
			EstimatedMinutes: twapi.Ref(int64(120)),
			TasklistID:       &resourceIDs.tasklistID,
			Assignees: &twapi.UserGroups{
				UserIDs: []int64{resourceIDs.userID},
			},
			TagIDs: []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update task: %v", err)
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

func createTag(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tagCreate := tag.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	tagIDSetter := twapi.WithIDCallback("id", func(id int64) {
		resourceIDs.tagID = id
	})

	logger.Info("⚙️  Creating tag")
	if err := engine.Do(ctx, &tagCreate, tagIDSetter); err != nil {
		logger.Error("failed to create tag",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created tag",
		slog.Int64("id", resourceIDs.tagID),
		slog.String("name", tagCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up tag",
			slog.Int64("id", resourceIDs.tagID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tagDelete tag.Delete
		tagDelete.Request.Path.ID = resourceIDs.tagID
		if err := engine.Do(ctx, &tagDelete); err != nil {
			logger.Warn("⚠️  failed to delete tag",
				slog.Int64("id", resourceIDs.tagID),
				slog.String("error", err.Error()),
			)
		}
	}
}

func createUser(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	userCreate := user.Create{
		FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
		Email:     fmt.Sprintf("test@test%d%d.com", time.Now().UnixNano(), rand.Intn(100)),
	}

	userIDSetter := twapi.WithIDCallback("id", func(id int64) {
		resourceIDs.userID = id
	})

	logger.Info("⚙️  Creating user")
	if err := engine.Do(ctx, &userCreate, userIDSetter); err != nil {
		logger.Error("failed to create user",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created user",
		slog.Int64("id", resourceIDs.userID),
		slog.String("name", fmt.Sprintf("%s %s", userCreate.FirstName, userCreate.LastName)),
	)

	var addProject user.AddProject
	addProject.Request.Path.ProjectID = resourceIDs.projectID
	addProject.Request.Users.IDs = []int64{resourceIDs.userID}

	logger.Info("⚙️  Adding user to project")
	if err := engine.Do(ctx, &addProject); err != nil {
		logger.Error("failed to add user to project",
			slog.Int64("userID", resourceIDs.userID),
			slog.Int64("projectID", resourceIDs.projectID),
			slog.String("error", err.Error()),
		)
	}
	logger.Info("✅ Added user to project")

	return func() {
		logger.Info("🗑️  Cleaning up user",
			slog.Int64("id", resourceIDs.userID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var userDelete user.Delete
		userDelete.Request.Path.ID = resourceIDs.userID
		if err := engine.Do(ctx, &userDelete); err != nil {
			logger.Warn("⚠️  failed to delete user",
				slog.Int64("id", resourceIDs.userID),
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

	deleteTag := createTag(logger)
	if resourceIDs.tagID == 0 {
		exitCode = 1
		return
	}
	defer deleteTag()

	deleteUser := createUser(logger)
	if resourceIDs.userID == 0 {
		exitCode = 1
		return
	}
	defer deleteUser()

	reference := time.Now()
	defer func() {
		if diff := time.Since(reference); diff < 200*time.Millisecond {
			time.Sleep(200*time.Millisecond - diff) // ensure tests have enough time to sync
		}
	}()

	exitCode = m.Run()
}
