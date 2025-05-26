package user_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork"
	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

func TestSingle(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := user.Create{
		FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
		Email:     fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var userID int64
	userIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		userID = i
	})
	if err := engine.Do(ctx, &create, userIDSetter); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var userDelete user.Delete
		userDelete.Request.Path.ID = userID
		if err := engine.Do(ctx, &userDelete); err != nil {
			t.Logf("⚠️  failed to delete user: %v", err)
		}
	})

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var single user.Single
	single.ID = userID

	if err := engine.Do(ctx, &single); err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if single.ID != userID {
		t.Errorf("expected user ID %d, got %d", userID, single.ID)
	}
}

func TestMultiple(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	create := user.Create{
		FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
		Email:     fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var userID int64
	userIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		userID = i
	})
	if err := engine.Do(ctx, &create, userIDSetter); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var userDelete user.Delete
		userDelete.Request.Path.ID = userID
		if err := engine.Do(ctx, &userDelete); err != nil {
			t.Logf("⚠️  failed to delete user: %v", err)
		}
	})

	tests := []struct {
		name     string
		multiple user.Multiple
	}{{
		name: "all users",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.multiple); err != nil {
				t.Errorf("failed to get users: %v", err)

			} else if len(tt.multiple.Response.Users) == 0 {
				t.Error("expected at least one user, got none")
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
		create user.Create
	}{{
		name: "only required fields",
		create: user.Create{
			FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Email:     fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100)),
		},
	}, {
		name: "all fields",
		create: user.Create{
			FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
			LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
			Title:     teamwork.Ref("Test User"),
			Email:     fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100)),
			Admin:     teamwork.Ref(true),
			Type:      teamwork.Ref("account"),
			CompanyID: &resourceIDs.companyID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var userID int64
			userIDSetter := teamwork.WithIDCallback("id", func(id int64) {
				userID = id
			})

			if err := engine.Do(ctx, &tt.create, userIDSetter); err != nil {
				t.Errorf("failed to create user: %v", err)

			} else {
				t.Cleanup(func() {
					ctx := context.Background()
					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					var userDelete user.Delete
					userDelete.Request.Path.ID = userID
					if err := engine.Do(ctx, &userDelete); err != nil {
						t.Logf("⚠️  failed to delete user: %v", err)
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

	create := user.Create{
		FirstName: fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100)),
		LastName:  fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100)),
		Email:     fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100)),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var userID int64
	userIDSetter := teamwork.WithIDCallback("id", func(i int64) {
		userID = i
	})
	if err := engine.Do(ctx, &create, userIDSetter); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var userDelete user.Delete
		userDelete.Request.Path.ID = userID
		if err := engine.Do(ctx, &userDelete); err != nil {
			t.Logf("⚠️  failed to delete user: %v", err)
		}
	})

	tests := []struct {
		name   string
		create user.Update
	}{{
		name: "all fields",
		create: user.Update{
			ID:        userID,
			FirstName: teamwork.Ref(fmt.Sprintf("test%d%d", time.Now().UnixNano(), rand.Intn(100))),
			LastName:  teamwork.Ref(fmt.Sprintf("user%d%d", time.Now().UnixNano(), rand.Intn(100))),
			Title:     teamwork.Ref("Test User"),
			Email:     teamwork.Ref(fmt.Sprintf("email%d%d@example.com", time.Now().UnixNano(), rand.Intn(100))),
			Admin:     teamwork.Ref(true),
			Type:      teamwork.Ref("account"),
			CompanyID: &resourceIDs.companyID,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.create); err != nil {
				t.Errorf("failed to update user: %v", err)
			}
		})
	}
}
