package team_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/company"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/project"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/team"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		companyID int64
		projectID int64
		userID    int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := team.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var teamID int64
	teamIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		teamID = i
	})
	if err := engine.Do(ctx, &create, teamIDSetter); err != nil {
		t.Fatalf("failed to create team: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var teamDelete team.Delete
		teamDelete.Request.Path.ID = teamID
		if err := engine.Do(ctx, &teamDelete); err != nil {
			t.Logf("⚠️  failed to delete team: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single team.Single
	single.ID = teamwork.LegacyNumber(teamID)

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get team: %v", err)
	}
	if single.ID != teamwork.LegacyNumber(teamID) {
		t.Errorf("expected team ID %d, got %d", teamID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := team.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var teamID int64
	teamIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		teamID = i
	})
	if err := engine.Do(ctx, &create, teamIDSetter); err != nil {
		t.Fatalf("failed to create team: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var teamDelete team.Delete
		teamDelete.Request.Path.ID = teamID
		if err := engine.Do(ctx, &teamDelete); err != nil {
			t.Logf("⚠️  failed to delete team: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple team.Multiple
	}{{
		name: "all teams",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get teams: %v", err)

			} else if len(tt.multiple.Response.Teams) == 0 {
				t.Error("expected at least one team, got none")
			}
		})
	}
}

func TestCreate(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := team.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var parentTeamID int64
	teamIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		parentTeamID = i
	})
	if err := engine.Do(ctx, &create, teamIDSetter); err != nil {
		t.Fatalf("failed to create team: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var teamDelete team.Delete
		teamDelete.Request.Path.ID = parentTeamID
		if err := engine.Do(ctx, &teamDelete); err != nil {
			t.Logf("⚠️  failed to delete team: %v", err)
		}
	})

	tests := []struct {
		name   string
		create team.Create
	}{{
		name: "only required fields",
		create: team.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields for company",
		create: team.Create{
			Name:         fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Handle:       teamwork.Ref(fmt.Sprintf("testhandle%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description:  teamwork.Ref("This is a test team."),
			ParentTeamID: &parentTeamID,
			CompanyID:    &resourceIDs.companyID,
			UserIDs:      []int64{resourceIDs.userID},
		},
	}, {
		name: "all fields for project",
		create: team.Create{
			Name:         fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Handle:       teamwork.Ref(fmt.Sprintf("testhandle%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description:  teamwork.Ref("This is a test team."),
			ParentTeamID: &parentTeamID,
			ProjectID:    &resourceIDs.projectID,
			UserIDs:      []int64{resourceIDs.userID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var teamID int64
			teamIDSetter := teamwork.WithIDCallback("id", func(id int64) {
				teamID = id
			})

			if err := engine.Do(ctx, &tt.create, teamIDSetter); err != nil {
				t.Errorf("failed to create team: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var teamDelete team.Delete
					teamDelete.Request.Path.ID = teamID
					if err := engine.Do(ctx, &teamDelete); err != nil {
						t.Logf("⚠️  failed to delete team: %v", err)
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

	create := team.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var teamID int64
	teamIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		teamID = i
	})
	if err := engine.Do(ctx, &create, teamIDSetter); err != nil {
		t.Fatalf("failed to create team: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var teamDelete team.Delete
		teamDelete.Request.Path.ID = teamID
		if err := engine.Do(ctx, &teamDelete); err != nil {
			t.Logf("⚠️  failed to delete team: %v", err)
		}
	})

	tests := []struct {
		name   string
		create team.Update
	}{{
		name: "all fields for company",
		create: team.Update{
			Name:        teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Handle:      teamwork.Ref(fmt.Sprintf("testhandle%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description: teamwork.Ref("This is a test team."),
			CompanyID:   &resourceIDs.companyID,
			UserIDs:     []int64{resourceIDs.userID},
		},
	}, {
		name: "all fields for project",
		create: team.Update{
			Name:        teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Handle:      teamwork.Ref(fmt.Sprintf("testhandle%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Description: teamwork.Ref("This is a test team."),
			ProjectID:   &resourceIDs.projectID,
			UserIDs:     []int64{resourceIDs.userID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			tt.create.ID = teamID
			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update team: %v", err)
			}
		})
	}
}

func createCompany(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	companyCreate := company.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	companyIDSetter := teamwork.WithIDCallback("id", func(id int64) {
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

func createProject(logger *slog.Logger) func() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	projectCreate := project.Create{
		Name:      fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		CompanyID: resourceIDs.companyID,
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

	deleteCompany := createCompany(logger)
	if resourceIDs.companyID == 0 {
		exitCode = 1
		return
	}
	defer deleteCompany()

	deleteProject := createProject(logger)
	if resourceIDs.projectID == 0 {
		exitCode = 1
		return
	}
	defer deleteProject()

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
