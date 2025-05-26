package timelog_test

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
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/task"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasklist"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/timelog"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		projectID  int64
		tasklistID int64
		taskID     int64
		tagID      int64
		userID     int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := timelog.Create{
		Date:      teamwork.Date(time.Now()),
		Hours:     1,
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timelogID int64
	timelogIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		timelogID = i
	})
	if err := engine.Do(ctx, &create, timelogIDSetter); err != nil {
		t.Fatalf("failed to create timelog: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timelogDelete timelog.Delete
		timelogDelete.Request.Path.ID = timelogID
		if err := engine.Do(ctx, &timelogDelete); err != nil {
			t.Logf("⚠️  failed to delete timelog: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single timelog.Single
	single.ID = timelogID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get timelog: %v", err)
	}
	if single.ID != timelogID {
		t.Errorf("expected timelog ID %d, got %d", timelogID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	createWithinProject := timelog.Create{
		Date:      teamwork.Date(time.Now()),
		Hours:     1,
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var projectTimelogID int64
	timelogIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		projectTimelogID = i
	})
	if err := engine.Do(ctx, &createWithinProject, timelogIDSetter); err != nil {
		t.Fatalf("failed to create timelog: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timelogDelete timelog.Delete
		timelogDelete.Request.Path.ID = projectTimelogID
		if err := engine.Do(ctx, &timelogDelete); err != nil {
			t.Logf("⚠️  failed to delete timelog: %v", err)
		}
	})

	createWithinTask := timelog.Create{
		Date:   teamwork.Date(time.Now()),
		Hours:  1,
		TaskID: resourceIDs.taskID,
	}

	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var taskTimelogID int64
	timelogIDSetter = teamwork.WithIDCallback("id", func(i int64) {
		taskTimelogID = i
	})
	if err := engine.Do(ctx, &createWithinTask, timelogIDSetter); err != nil {
		t.Fatalf("failed to create timelog: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timelogDelete timelog.Delete
		timelogDelete.Request.Path.ID = taskTimelogID
		if err := engine.Do(ctx, &timelogDelete); err != nil {
			t.Logf("⚠️  failed to delete timelog: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple timelog.Multiple
	}{{
		name: "all timelogs",
	}, {
		name: "timelogs for project",
		multiple: func() timelog.Multiple {
			var multiple timelog.Multiple
			multiple.Request.Path.ProjectID = resourceIDs.projectID
			return multiple
		}(),
	}, {
		name: "timelogs for task",
		multiple: func() timelog.Multiple {
			var multiple timelog.Multiple
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
				t.Errorf("failed to get timelogs: %v", err)

			} else if len(tt.multiple.Response.Timelogs) == 0 {
				t.Error("expected at least one timelog, got none")
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
		create timelog.Create
	}{{
		name: "only required fields for project",
		create: timelog.Create{
			Date:      teamwork.Date(time.Now()),
			Hours:     1,
			ProjectID: resourceIDs.projectID,
		},
	}, {
		name: "all fields for project",
		create: timelog.Create{
			Description: teamwork.Ref("This is a test timelog"),
			Date:        teamwork.Date(time.Now().UTC()),
			Time:        teamwork.Time(time.Now().UTC()),
			IsUTC:       true,
			Hours:       1,
			Minutes:     30,
			Billable:    true,
			ProjectID:   resourceIDs.projectID,
			UserID:      &resourceIDs.userID,
			TagIDs:      []int64{resourceIDs.tagID},
		},
	}, {
		name: "only required fields for task",
		create: timelog.Create{
			Date:   teamwork.Date(time.Now()),
			Hours:  1,
			TaskID: resourceIDs.taskID,
		},
	}, {
		name: "all fields for task",
		create: timelog.Create{
			Description: teamwork.Ref("This is a test timelog"),
			Date:        teamwork.Date(time.Now().UTC()),
			Time:        teamwork.Time(time.Now().UTC()),
			IsUTC:       true,
			Hours:       1,
			Minutes:     30,
			Billable:    true,
			TaskID:      resourceIDs.taskID,
			UserID:      &resourceIDs.userID,
			TagIDs:      []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var timelogID int64
			timelogIDSetter := teamwork.WithIDCallback("id", func(id int64) {
				timelogID = id
			})

			if err := engine.Do(ctx, &tt.create, timelogIDSetter); err != nil {
				t.Errorf("failed to create timelog: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var timelogDelete timelog.Delete
					timelogDelete.Request.Path.ID = timelogID
					if err := engine.Do(ctx, &timelogDelete); err != nil {
						t.Logf("⚠️  failed to delete timelog: %v", err)
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

	create := timelog.Create{
		Date:      teamwork.Date(time.Now()),
		Hours:     1,
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var timelogID int64
	timelogIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		timelogID = i
	})
	if err := engine.Do(ctx, &create, timelogIDSetter); err != nil {
		t.Fatalf("failed to create timelog: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var timelogDelete timelog.Delete
		timelogDelete.Request.Path.ID = timelogID
		if err := engine.Do(ctx, &timelogDelete); err != nil {
			t.Logf("⚠️  failed to delete timelog: %v", err)
		}
	})

	tests := []struct {
		name   string
		create timelog.Update
	}{{
		name: "all fields",
		create: timelog.Update{
			ID:          timelogID,
			Description: teamwork.Ref("Updated description"),
			Date:        teamwork.Ref(teamwork.Date(time.Now().UTC())),
			Time:        teamwork.Ref(teamwork.Time(time.Now().UTC())),
			IsUTC:       teamwork.Ref(true),
			Hours:       teamwork.Ref(int64(2)),
			Minutes:     teamwork.Ref(int64(15)),
			Billable:    teamwork.Ref(true),
			UserID:      teamwork.Ref(resourceIDs.userID),
			TagIDs:      []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update timelog: %v", err)
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

	logger.Info("⚙️ Creating project")
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

	tasklistIDSetter := teamwork.WithIDCallback("tasklistId", func(id int64) {
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

	taskIDSetter := teamwork.WithIDCallback("id", func(id int64) {
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

func createTag(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tagCreate := tag.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	tagIDSetter := teamwork.WithIDCallback("id", func(id int64) {
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

	userIDSetter := teamwork.WithIDCallback("id", func(id int64) {
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

	logger.Info("⚙️ Adding user to project")
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
