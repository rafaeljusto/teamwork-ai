package milestone_test

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

	create := milestone.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		DueDate:   teamwork.LegacyDate(time.Now().Add(24 * time.Hour)),
		ProjectID: resourceIDs.projectID,
		Assignees: teamwork.LegacyUserGroups{
			UserIDs: []int64{resourceIDs.userID},
		},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var milestoneID int64
	milestoneIDSetter := teamwork.WithIDCallback("milestoneId", func(i int64) {
		milestoneID = i
	})
	if err := engine.Do(ctx, &create, milestoneIDSetter); err != nil {
		t.Fatalf("failed to create milestone: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var milestoneDelete milestone.Delete
		milestoneDelete.Request.Path.ID = milestoneID
		if err := engine.Do(ctx, &milestoneDelete); err != nil {
			t.Logf("⚠️  failed to delete milestone: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single milestone.Single
	single.ID = milestoneID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get milestone: %v", err)
	}
	if single.ID != milestoneID {
		t.Errorf("expected milestone ID %d, got %d", milestoneID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := milestone.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		DueDate:   teamwork.LegacyDate(time.Now().Add(24 * time.Hour)),
		ProjectID: resourceIDs.projectID,
		Assignees: teamwork.LegacyUserGroups{
			UserIDs: []int64{resourceIDs.userID},
		},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var milestoneID int64
	milestoneIDSetter := teamwork.WithIDCallback("milestoneId", func(i int64) {
		milestoneID = i
	})
	if err := engine.Do(ctx, &create, milestoneIDSetter); err != nil {
		t.Fatalf("failed to create milestone: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var milestoneDelete milestone.Delete
		milestoneDelete.Request.Path.ID = milestoneID
		if err := engine.Do(ctx, &milestoneDelete); err != nil {
			t.Logf("⚠️  failed to delete milestone: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple milestone.Multiple
	}{{
		name: "all milestones",
	}, {
		name: "milestones for project",
		multiple: func() milestone.Multiple {
			var multiple milestone.Multiple
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
				t.Errorf("failed to get milestones: %v", err)

			} else if len(tt.multiple.Response.Milestones) == 0 {
				t.Error("expected at least one milestone, got none")
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
		create milestone.Create
	}{{
		name: "only required fields",
		create: milestone.Create{
			Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			DueDate:   teamwork.LegacyDate(time.Now().Add(24 * time.Hour)),
			ProjectID: resourceIDs.projectID,
			Assignees: teamwork.LegacyUserGroups{
				UserIDs: []int64{resourceIDs.userID},
			},
		},
	}, {
		name: "all fields",
		create: milestone.Create{
			Name:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Description: pointerTo("This is a test milestone"),
			DueDate:     teamwork.LegacyDate(time.Now().Add(48 * time.Hour)),
			ProjectID:   resourceIDs.projectID,
			TasklistIDs: []int64{resourceIDs.tasklistID},
			TagIDs:      []int64{resourceIDs.tagID},
			Assignees: teamwork.LegacyUserGroups{
				UserIDs: []int64{resourceIDs.userID},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var milestoneID int64
			milestoneIDSetter := teamwork.WithIDCallback("milestoneId", func(id int64) {
				milestoneID = id
			})

			if err := engine.Do(ctx, &tt.create, milestoneIDSetter); err != nil {
				t.Errorf("failed to create milestone: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var milestoneDelete milestone.Delete
					milestoneDelete.Request.Path.ID = milestoneID
					if err := engine.Do(ctx, &milestoneDelete); err != nil {
						t.Logf("⚠️  failed to delete milestone: %v", err)
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

	create := milestone.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		DueDate:   teamwork.LegacyDate(time.Now().Add(24 * time.Hour)),
		ProjectID: resourceIDs.projectID,
		Assignees: teamwork.LegacyUserGroups{
			UserIDs: []int64{resourceIDs.userID},
		},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var milestoneID int64
	milestoneIDSetter := teamwork.WithIDCallback("milestoneId", func(i int64) {
		milestoneID = i
	})
	if err := engine.Do(ctx, &create, milestoneIDSetter); err != nil {
		t.Fatalf("failed to create milestone: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var milestoneDelete milestone.Delete
		milestoneDelete.Request.Path.ID = milestoneID
		if err := engine.Do(ctx, &milestoneDelete); err != nil {
			t.Logf("⚠️  failed to delete milestone: %v", err)
		}
	})

	tests := []struct {
		name   string
		create milestone.Update
	}{{
		name: "all fields",
		create: milestone.Update{
			ID:          milestoneID,
			Name:        pointerTo(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description: pointerTo("This is a test milestone"),
			DueDate:     pointerTo(teamwork.LegacyDate(time.Now().Add(48 * time.Hour))),
			TasklistIDs: []int64{resourceIDs.tasklistID},
			TagIDs:      []int64{resourceIDs.tagID},
			Assignees: &teamwork.LegacyUserGroups{
				UserIDs: []int64{resourceIDs.userID},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update milestone: %v", err)
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

func pointerTo[T any](t T) *T {
	return &t
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
