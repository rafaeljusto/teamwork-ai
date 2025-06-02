package agentic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rafaeljusto/teamwork-ai/internal/twapi/jobrole"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/skill"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
)

var registered map[string]Agentic

// Register registers an agentic implementation with the given name. The name is
// used to identify the agentic implementation when initializing it. The agentic
// implementation must implement the Agentic interface.
func Register(name string, agentic Agentic) {
	if registered == nil {
		registered = make(map[string]Agentic)
	}
	registered[name] = agentic
}

// Init initializes the agentic system with the provided name, and DSN. The name
// must be from a pre-registered agentic implementation. The DSN is specific to
// the agentic implementation and is used to configure it.
func Init(name, dsn string, logger *slog.Logger) Agentic {
	if name == "" {
		return nil
	}
	agentic, ok := registered[name]
	if !ok {
		panic(fmt.Errorf("unknown agentic implementation: %s", name))
	}
	if err := agentic.Init(dsn, logger); err != nil {
		panic(fmt.Errorf("failed to initialize agentic implementation: %w", err))
	}
	return agentic
}

// Agentic stores mechanisms to build autonomous systems capable of making
// decisions and performing tasks without constant human intervention.
type Agentic interface {
	// Init initializes the agentic system with the provided DSN.
	Init(dsn string, logger *slog.Logger) error

	// FindTaskSkillsAndJobRoles finds the skills and job roles for a given task.
	// It uses the task data, available skills, and available job roles to
	// determine the most relevant skills and job roles IDs for the task.
	FindTaskSkillsAndJobRoles(
		ctx context.Context,
		taskDate webhook.TaskData,
		availableSkills []skill.Skill,
		availableJobRoles []jobrole.JobRole,
	) (skillIDs, jobRoleIDs []int64, reasoning string, err error)
}
