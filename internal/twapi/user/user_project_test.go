package user_test

import (
	"context"
	"testing"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi/user"
)

func TestAddProject(t *testing.T) {
	if engine == nil {
		t.Skip("Skipping test because the engine is not initialized")
	}

	tests := []struct {
		name       string
		addProject user.AddProject
	}{{
		name: "only required fields",
		addProject: func() user.AddProject {
			var addProject user.AddProject
			addProject.Request.Path.ProjectID = resourceIDs.projectID
			addProject.Request.Users.IDs = []int64{resourceIDs.userID}
			return addProject
		}(),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			if err := engine.Do(ctx, &tt.addProject); err != nil {
				t.Errorf("failed to add user to project: %v", err)
			}
		})
	}
}
