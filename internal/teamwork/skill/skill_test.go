package skill_test

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/skill"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		userID int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := skill.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var skillID int64
	skillIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		skillID = i
	})
	if err := engine.Do(ctx, &create, skillIDSetter); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var skillDelete skill.Delete
		skillDelete.Request.Path.ID = skillID
		if err := engine.Do(ctx, &skillDelete); err != nil {
			t.Logf("⚠️  failed to delete skill: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single skill.Single
	single.ID = skillID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get skill: %v", err)
	}
	if single.ID != skillID {
		t.Errorf("expected skill ID %d, got %d", skillID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := skill.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var skillID int64
	skillIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		skillID = i
	})
	if err := engine.Do(ctx, &create, skillIDSetter); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var skillDelete skill.Delete
		skillDelete.Request.Path.ID = skillID
		if err := engine.Do(ctx, &skillDelete); err != nil {
			t.Logf("⚠️  failed to delete skill: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple skill.Multiple
	}{{
		name: "all skills",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get skills: %v", err)

			} else if len(tt.multiple.Response.Skills) == 0 {
				t.Error("expected at least one skill, got none")
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
		create skill.Create
	}{{
		name: "only required fields",
		create: skill.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields",
		create: skill.Create{
			Name:    fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			UserIDs: []int64{resourceIDs.userID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var skillID int64
			skillIDSetter := teamwork.WithIDCallback("skillId", func(id int64) {
				skillID = id
			})

			if err := engine.Do(ctx, &tt.create, skillIDSetter); err != nil {
				t.Errorf("failed to create skill: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var skillDelete skill.Delete
					skillDelete.Request.Path.ID = skillID
					if err := engine.Do(ctx, &skillDelete); err != nil {
						t.Logf("⚠️  failed to delete skill: %v", err)
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

	create := skill.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var skillID int64
	skillIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		skillID = i
	})
	if err := engine.Do(ctx, &create, skillIDSetter); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var skillDelete skill.Delete
		skillDelete.Request.Path.ID = skillID
		if err := engine.Do(ctx, &skillDelete); err != nil {
			t.Logf("⚠️  failed to delete skill: %v", err)
		}
	})

	tests := []struct {
		name   string
		create skill.Update
	}{{
		name: "all fields",
		create: skill.Update{
			Name:    teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			UserIDs: []int64{resourceIDs.userID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			tt.create.ID = skillID
			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update skill: %v", err)
			}
		})
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
