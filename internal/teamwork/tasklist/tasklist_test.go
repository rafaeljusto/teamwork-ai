package tasklist_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/milestone"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/project"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/tag"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/tasklist"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		projectID   int64
		milestoneID int64
		tagID       int64
		userID      int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := tasklist.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tasklistID int64
	tasklistIDSetter := teamwork.WithIDCallback("tasklistId", func(i int64) {
		tasklistID = i
	})
	if err := engine.Do(ctx, &create, tasklistIDSetter); err != nil {
		t.Fatalf("failed to create tasklist: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tasklistDelete tasklist.Delete
		tasklistDelete.Request.Path.ID = tasklistID
		if err := engine.Do(ctx, &tasklistDelete); err != nil {
			t.Logf("⚠️  failed to delete tasklist: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single tasklist.Single
	single.ID = tasklistID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get tasklist: %v", err)
	}
	if single.ID != tasklistID {
		t.Errorf("expected tasklist ID %d, got %d", tasklistID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := tasklist.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tasklistID int64
	tasklistIDSetter := teamwork.WithIDCallback("tasklistId", func(i int64) {
		tasklistID = i
	})
	if err := engine.Do(ctx, &create, tasklistIDSetter); err != nil {
		t.Fatalf("failed to create tasklist: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tasklistDelete tasklist.Delete
		tasklistDelete.Request.Path.ID = tasklistID
		if err := engine.Do(ctx, &tasklistDelete); err != nil {
			t.Logf("⚠️  failed to delete tasklist: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple tasklist.Multiple
	}{{
		name: "all tasklists",
	}, {
		name: "tasklists for project",
		multiple: func() tasklist.Multiple {
			var multiple tasklist.Multiple
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
				t.Errorf("failed to get tasklists: %v", err)

			} else if len(tt.multiple.Response.Tasklists) == 0 {
				t.Error("expected at least one tasklist, got none")
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
		create tasklist.Create
	}{{
		name: "only required fields",
		create: tasklist.Create{
			Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			ProjectID: resourceIDs.projectID,
		},
	}, {
		name: "all fields",
		create: tasklist.Create{
			Name:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Description: teamwork.Ref("This is a test tasklist"),
			ProjectID:   resourceIDs.projectID,
			MilestoneID: &resourceIDs.milestoneID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var tasklistID int64
			tasklistIDSetter := teamwork.WithIDCallback("tasklistId", func(id int64) {
				tasklistID = id
			})

			if err := engine.Do(ctx, &tt.create, tasklistIDSetter); err != nil {
				t.Errorf("failed to create tasklist: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var tasklistDelete tasklist.Delete
					tasklistDelete.Request.Path.ID = tasklistID
					if err := engine.Do(ctx, &tasklistDelete); err != nil {
						t.Logf("⚠️  failed to delete tasklist: %v", err)
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

	create := tasklist.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		ProjectID: resourceIDs.projectID,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tasklistID int64
	tasklistIDSetter := teamwork.WithIDCallback("tasklistId", func(i int64) {
		tasklistID = i
	})
	if err := engine.Do(ctx, &create, tasklistIDSetter); err != nil {
		t.Fatalf("failed to create tasklist: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var tasklistDelete tasklist.Delete
		tasklistDelete.Request.Path.ID = tasklistID
		if err := engine.Do(ctx, &tasklistDelete); err != nil {
			t.Logf("⚠️  failed to delete tasklist: %v", err)
		}
	})

	tests := []struct {
		name   string
		create tasklist.Update
	}{{
		name: "all fields",
		create: tasklist.Update{
			ID:          tasklistID,
			Name:        teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description: teamwork.Ref("This is a test tasklist"),
			MilestoneID: &resourceIDs.milestoneID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update tasklist: %v", err)
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

func createMilestone(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	milestoneCreate := milestone.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		DueDate:   teamwork.LegacyDate(time.Now().Add(24 * time.Hour)),
		ProjectID: resourceIDs.projectID,
		Assignees: teamwork.LegacyUserGroups{
			UserIDs: []int64{resourceIDs.userID},
		},
	}

	milestoneIDSetter := teamwork.WithIDCallback("milestoneId", func(id int64) {
		resourceIDs.milestoneID = id
	})

	logger.Info("⚙️  Creating milestone")
	if err := engine.Do(ctx, &milestoneCreate, milestoneIDSetter); err != nil {
		logger.Error("failed to create milestone",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created milestone",
		slog.Int64("id", resourceIDs.milestoneID),
		slog.String("name", milestoneCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up milestone",
			slog.Int64("id", resourceIDs.milestoneID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var milestoneDelete milestone.Delete
		milestoneDelete.Request.Path.ID = resourceIDs.milestoneID
		if err := engine.Do(ctx, &milestoneDelete); err != nil {
			logger.Warn("⚠️  failed to delete milestone",
				slog.Int64("id", resourceIDs.milestoneID),
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

	deleteMilestone := createMilestone(logger)
	if resourceIDs.milestoneID == 0 {
		exitCode = 1
		return
	}
	defer deleteMilestone()

	reference := time.Now()
	defer func() {
		if diff := time.Since(reference); diff < 200*time.Millisecond {
			time.Sleep(200*time.Millisecond - diff) // ensure tests have enough time to sync
		}
	}()

	exitCode = m.Run()
}
