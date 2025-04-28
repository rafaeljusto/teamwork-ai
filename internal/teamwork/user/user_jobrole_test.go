package user_test

import (
	"context"
	"testing"

	"github.com/rafaeljusto/teamwork-ai/internal/teamwork/user"
)

func TestAssignJobRole(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	tests := []struct {
		name            string
		assignJobRole   user.AssignJobRole
		unassignJobRole user.UnassignJobRole
	}{{
		name: "only required fields",
		assignJobRole: func() user.AssignJobRole {
			var assignJobRole user.AssignJobRole
			assignJobRole.Request.Path.JobRoleID = resourceIDs.jobRoleID
			assignJobRole.Request.Users.IDs = []int64{resourceIDs.userID}
			assignJobRole.Request.Users.IsPrimary = true
			return assignJobRole
		}(),
		unassignJobRole: func() user.UnassignJobRole {
			var unassignJobRole user.UnassignJobRole
			unassignJobRole.Request.Path.JobRoleID = resourceIDs.jobRoleID
			unassignJobRole.Request.Users.IDs = []int64{resourceIDs.userID}
			return unassignJobRole
		}(),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.assignJobRole); err != nil {
				t.Errorf("failed to assign a user to a job role: %v", err)
			} else if err := engine.Do(ctx, &tt.unassignJobRole); err != nil {
				t.Errorf("failed to unassign user from a job role: %v", err)
			}
		})
	}
}
