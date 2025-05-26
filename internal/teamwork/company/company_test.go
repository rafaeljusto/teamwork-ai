package company_test

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
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/tag"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

const timeout = 5 * time.Second

var (
	engine      *teamwork.Engine
	resourceIDs struct {
		tagID  int64
		userID int64
	}
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := company.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var companyID int64
	companyIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		companyID = i
	})
	if err := engine.Do(ctx, &create, companyIDSetter); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var companyDelete company.Delete
		companyDelete.Request.Path.ID = companyID
		if err := engine.Do(ctx, &companyDelete); err != nil {
			t.Logf("⚠️  failed to delete company: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single company.Single
	single.ID = companyID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get company: %v", err)
	}
	if single.ID != companyID {
		t.Errorf("expected company ID %d, got %d", companyID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := company.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var companyID int64
	companyIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		companyID = i
	})
	if err := engine.Do(ctx, &create, companyIDSetter); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var companyDelete company.Delete
		companyDelete.Request.Path.ID = companyID
		if err := engine.Do(ctx, &companyDelete); err != nil {
			t.Logf("⚠️  failed to delete company: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple company.Multiple
	}{{
		name: "all companies",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get companies: %v", err)

			} else if len(tt.multiple.Response.Companies) == 0 {
				t.Error("expected at least one company, got none")
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
		create company.Create
	}{{
		name: "only required fields",
		create: company.Create{
			Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields",
		create: company.Create{
			Name:        fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			AddressOne:  teamwork.Ref("123 Main St"),
			AddressTwo:  teamwork.Ref("Apt. 456"),
			City:        teamwork.Ref("Cork"),
			CountryCode: teamwork.Ref("IR"),
			EmailOne:    teamwork.Ref("test1@company.com"),
			EmailTwo:    teamwork.Ref("test2@company.com"),
			EmailThree:  teamwork.Ref("test3@company.com"),
			Fax:         teamwork.Ref("123-456-7890"),
			Phone:       teamwork.Ref("123-456-7890"),
			Profile:     teamwork.Ref("This is a test company profile."),
			State:       teamwork.Ref("Cork"),
			Website:     teamwork.Ref("https://www.example.com"),
			Zip:         teamwork.Ref("12345"),
			ManagerID:   &resourceIDs.userID,
			IndustryID:  teamwork.Ref(int64(1)), // Web Development Agency,
			TagIDs:      []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var companyID int64
			companyIDSetter := teamwork.WithIDCallback("id", func(id int64) {
				companyID = id
			})

			if err := engine.Do(ctx, &tt.create, companyIDSetter); err != nil {
				t.Errorf("failed to create company: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var companyDelete company.Delete
					companyDelete.Request.Path.ID = companyID
					if err := engine.Do(ctx, &companyDelete); err != nil {
						t.Logf("⚠️  failed to delete company: %v", err)
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

	create := company.Create{
		Name: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var companyID int64
	companyIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		companyID = i
	})
	if err := engine.Do(ctx, &create, companyIDSetter); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var companyDelete company.Delete
		companyDelete.Request.Path.ID = companyID
		if err := engine.Do(ctx, &companyDelete); err != nil {
			t.Logf("⚠️  failed to delete company: %v", err)
		}
	})

	tests := []struct {
		name   string
		create company.Update
	}{{
		name: "all fields",
		create: company.Update{
			Name:        teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			AddressOne:  teamwork.Ref("123 Main St"),
			AddressTwo:  teamwork.Ref("Apt. 456"),
			City:        teamwork.Ref("Cork"),
			CountryCode: teamwork.Ref("IR"),
			EmailOne:    teamwork.Ref("test1@company.com"),
			EmailTwo:    teamwork.Ref("test2@company.com"),
			EmailThree:  teamwork.Ref("test3@company.com"),
			Fax:         teamwork.Ref("123-456-7890"),
			Phone:       teamwork.Ref("123-456-7890"),
			Profile:     teamwork.Ref("This is a test company profile."),
			State:       teamwork.Ref("Cork"),
			Website:     teamwork.Ref("https://www.example.com"),
			Zip:         teamwork.Ref("12345"),
			ManagerID:   &resourceIDs.userID,
			IndustryID:  teamwork.Ref(int64(1)), // Web Development Agency,
			TagIDs:      []int64{resourceIDs.tagID},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			tt.create.ID = companyID
			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update company: %v", err)
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
