package project_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/company"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/project"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/tag"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/user"
)

const timeout = 5 * time.Second

var (
	engine      *twapi.Engine
	resourceIDs struct {
		tagID     int64
		companyID int64
		userID    int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := project.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var projectID int64
	projectIDSetter := twapi.WithIDCallback("id", func(i int64) {
		projectID = i
	})
	if err := engine.Do(ctx, &create, projectIDSetter); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var projectDelete project.Delete
		projectDelete.Request.Path.ID = projectID
		if err := engine.Do(ctx, &projectDelete); err != nil {
			t.Logf("⚠️  failed to delete project: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single project.Single
	single.ID = projectID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get project: %v", err)
	}
	if single.ID != projectID {
		t.Errorf("expected project ID %d, got %d", projectID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := project.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var projectID int64
	projectIDSetter := twapi.WithIDCallback("id", func(i int64) {
		projectID = i
	})
	if err := engine.Do(ctx, &create, projectIDSetter); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var projectDelete project.Delete
		projectDelete.Request.Path.ID = projectID
		if err := engine.Do(ctx, &projectDelete); err != nil {
			t.Logf("⚠️  failed to delete project: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple project.Multiple
	}{{
		name: "all projects",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get projects: %v", err)

			} else if len(tt.multiple.Response.Projects) == 0 {
				t.Error("expected at least one project, got none")
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
		create project.Create
	}{{
		name: "only required fields",
		create: project.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields",
		create: project.Create{
			Name:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Description: twapi.Ref("This is a test project"),
			StartAt:     twapi.Ref(twapi.LegacyDate(time.Now().Add(24 * time.Hour))),
			EndAt:       twapi.Ref(twapi.LegacyDate(time.Now().Add(48 * time.Hour))),
			CompanyID:   resourceIDs.companyID,
			OwnerID:     &resourceIDs.userID,
			Tags:        []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var projectID int64
			projectIDSetter := twapi.WithIDCallback("id", func(id int64) {
				projectID = id
			})

			if err := engine.Do(ctx, &tt.create, projectIDSetter); err != nil {
				t.Errorf("failed to create project: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var projectDelete project.Delete
					projectDelete.Request.Path.ID = projectID
					if err := engine.Do(ctx, &projectDelete); err != nil {
						t.Logf("⚠️  failed to delete project: %v", err)
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

	create := project.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var projectID int64
	projectIDSetter := twapi.WithIDCallback("id", func(i int64) {
		projectID = i
	})
	if err := engine.Do(ctx, &create, projectIDSetter); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var projectDelete project.Delete
		projectDelete.Request.Path.ID = projectID
		if err := engine.Do(ctx, &projectDelete); err != nil {
			t.Logf("⚠️  failed to delete project: %v", err)
		}
	})

	tests := []struct {
		name   string
		create project.Update
	}{{
		name: "all fields",
		create: project.Update{
			Name:        twapi.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description: twapi.Ref("This is a test project"),
			StartAt:     twapi.Ref(twapi.LegacyDate(time.Now().Add(24 * time.Hour))),
			EndAt:       twapi.Ref(twapi.LegacyDate(time.Now().Add(48 * time.Hour))),
			CompanyID:   &resourceIDs.companyID,
			OwnerID:     &resourceIDs.userID,
			Tags:        []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			tt.create.ID = projectID
			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update project: %v", err)
			}
		})
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

func createCompany(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	companyCreate := company.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	companyIDSetter := twapi.WithIDCallback("id", func(id int64) {
		resourceIDs.companyID = id
	})

	logger.Info("⚙️  Creating company")
	if err := engine.Do(ctx, &companyCreate, companyIDSetter); err != nil {
		logger.Error("failed to create company",
			slog.String("error", err.Error()),
		)
		return func() {}
	}
	logger.Info("✅ Created company",
		slog.Int64("id", resourceIDs.companyID),
		slog.String("name", companyCreate.Name),
	)

	return func() {
		logger.Info("🗑️  Cleaning up company",
			slog.Int64("id", resourceIDs.companyID),
		)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var companyDelete company.Delete
		companyDelete.Request.Path.ID = resourceIDs.companyID
		if err := engine.Do(ctx, &companyDelete); err != nil {
			logger.Warn("⚠️  failed to delete company",
				slog.Int64("id", resourceIDs.companyID),
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
		CompanyID: &resourceIDs.companyID,
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

	deleteTag := createTag(logger)
	if resourceIDs.tagID == 0 {
		exitCode = 1
		return
	}
	defer deleteTag()

	deleteCompany := createCompany(logger)
	if resourceIDs.companyID == 0 {
		exitCode = 1
		return
	}
	defer deleteCompany()

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
