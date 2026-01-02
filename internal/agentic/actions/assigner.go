package actions

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/webhook"
	twapi "github.com/teamwork/twapi-go-sdk"
	"github.com/teamwork/twapi-go-sdk/projects"
)

var processing sync.Map

// AutoAssignTaskOptions contains the options for the AutoAssignTask function.
type AutoAssignTaskOptions struct {
	skipRates      bool
	skipWorkload   bool
	skipAssignment bool
	skipComment    bool
}

// AutoAssignTaskOption is a function that sets an option for the AutoAssignTask
// function.
type AutoAssignTaskOption func(*AutoAssignTaskOptions)

// WithAutoAssignTaskSkipRates sets the skipRates option for the AutoAssignTask
// function. If set to true, the function will not consider the rates of the
// users when assigning the task.
func WithAutoAssignTaskSkipRates() AutoAssignTaskOption {
	return func(o *AutoAssignTaskOptions) {
		o.skipRates = true
	}
}

// WithAutoAssignTaskSkipWorkload sets the skipWorkload option for the
// AutoAssignTask function. If set to true, the function will not consider the
// workload of the users when assigning the task.
func WithAutoAssignTaskSkipWorkload() AutoAssignTaskOption {
	return func(o *AutoAssignTaskOptions) {
		o.skipWorkload = true
	}
}

// WithAutoAssignTaskSkipAssignment sets the skipAssignment option for the
// AutoAssignTask function. If set to true, the function will not assign the
// task to the users.
func WithAutoAssignTaskSkipAssignment() AutoAssignTaskOption {
	return func(o *AutoAssignTaskOptions) {
		o.skipAssignment = true
	}
}

// WithAutoAssignTaskSkipComment sets the skipComment option for the
// AutoAssignTask function. If set to true, the function will not create a
// comment on the task after assigning it.
func WithAutoAssignTaskSkipComment() AutoAssignTaskOption {
	return func(o *AutoAssignTaskOptions) {
		o.skipComment = true
	}
}

// AutoAssignTask assigns a task to users based on the skills and job roles
// associated with the task.
func AutoAssignTask(
	ctx context.Context,
	resources *config.Resources,
	taskData webhook.TaskData,
	optFuncs ...AutoAssignTaskOption,
) error {
	var options AutoAssignTaskOptions
	for _, optFunc := range optFuncs {
		optFunc(&options)
	}

	logger := resources.Logger.With(
		slog.String("action", "autoAssignTask"),
		slog.Int64("taskID", taskData.Task.ID),
	)

	if _, ok := processing.LoadOrStore(taskData.Task.ID, struct{}{}); ok {
		logger.Info("task already being processed, skipping AI assignment")
		return nil
	}
	defer processing.Delete(taskData.Task.ID)

	// if there's already an assigned user, we don't need to do anything
	if len(taskData.Task.AssignedUserIDs) > 0 {
		logger.Info("task already has assigned users, skipping AI assignment")
		return nil
	}

	skills, err := loadSkills(ctx, resources)
	if err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}
	skillsMap := skills.toMap()

	jobRoles, err := loadJobRoles(ctx, resources)
	if err != nil {
		return fmt.Errorf("failed to load job roles: %w", err)
	}
	jobRolesMap := jobRoles.toMap()

	projectUsers, err := loadProjectUsers(ctx, resources, taskData.Project.ID)
	if err != nil {
		return fmt.Errorf("failed to load project users: %w", err)
	}
	projectUsersMap := projectUsers.toMap()

	skillIDs, jobRoleIDs, reasoning, err := resources.Agentic.FindTaskSkillsAndJobRoles(ctx, taskData, skills, jobRoles)
	if err != nil {
		return fmt.Errorf("failed to find task skills and job roles: %w", err)
	}

	var userIDsWithSkills []int64
	for _, skillID := range skillIDs {
		skill, ok := skillsMap[skillID]
		if !ok {
			logger.Info("skill not found in the loaded skills, AI halucination",
				slog.Int64("skillID", skillID),
			)
			continue
		}
		userIDsWithSkills = append(userIDsWithSkills, extractMappedIDs(skill.Users, projectUsersMap)...)
	}

	var userIDsWithJobRoles []int64
	for _, jobRoleID := range jobRoleIDs {
		jobRole, ok := jobRolesMap[jobRoleID]
		if !ok {
			logger.Info("job role not found in the loaded job roles, AI halucination",
				slog.Int64("jobRoleID", jobRoleID),
			)
			continue
		}
		userIDsWithJobRoles = append(userIDsWithJobRoles, extractMappedIDs(jobRole.PrimaryUsers, projectUsersMap)...)
		if len(jobRole.PrimaryUsers) == 0 {
			userIDsWithJobRoles = append(userIDsWithJobRoles, extractMappedIDs(jobRole.Users, projectUsersMap)...)
		}
	}

	idealUserIDs := intersection(userIDsWithSkills, userIDsWithJobRoles)
	if len(idealUserIDs) == 0 {
		idealUserIDs = append(idealUserIDs, userIDsWithSkills...)
		idealUserIDs = append(idealUserIDs, userIDsWithJobRoles...)
	}

	if reasoning != "" && !strings.HasSuffix(reasoning, ".") {
		reasoning += "."
	}

	var processors []autoAssignTaskProcessor
	if !options.skipRates {
		processors = append(processors, autoAssignTaskProcessRates(projectUsersMap, &reasoning, logger))
	}
	if !options.skipWorkload {
		processors = append(processors, autoAssignTaskProcessWorkload(ctx, taskData, resources, &reasoning, logger))
	}
	userScores := newUserScores(idealUserIDs)
	for _, processor := range processors {
		if userScores, err = processor(userScores); err != nil {
			return fmt.Errorf("failed to process ideal user IDs: %w", err)
		}
	}
	idealUserIDs = userScores.chooseIDs()
	if len(idealUserIDs) == 0 {
		logger.Info("no users found with the AI suggested skills or job roles, skipping task assignment")
		return nil
	}

	if !options.skipAssignment {
		taskUpdate := projects.NewTaskUpdateRequest(taskData.Task.ID)
		taskUpdate.Path.ID = taskData.Task.ID
		taskUpdate.Assignees = &projects.UserGroups{
			UserIDs: idealUserIDs,
		}
		if _, err := projects.TaskUpdate(ctx, resources.TeamworkEngine, taskUpdate); err != nil {
			return fmt.Errorf("failed to assign task to users: %w", err)
		}
		logger.Info("task assigned to users based on AI",
			slog.Int64("id", taskData.Task.ID),
		)
	}

	if !options.skipComment {
		commentCreate := projects.NewCommentCreateRequestInTask(
			taskData.Task.ID,
			"🤖 Assignment of this task was performed by artificial intelligence.\n",
		)
		for _, userID := range idealUserIDs {
			if user, ok := projectUsersMap[userID]; ok {
				commentCreate.Body += fmt.Sprintf("\n  • %s %s", user.FirstName, user.LastName)
			}
		}
		commentCreate.Body += "\n\n" + reasoning
		if _, err := projects.CommentCreate(ctx, resources.TeamworkEngine, commentCreate); err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	return nil
}

type userScore struct {
	ID    int64
	Score int64
}

type userScores []userScore

func newUserScores(userIDs []int64) userScores {
	userScores := make(userScores, len(userIDs))
	for i, userID := range userIDs {
		userScores[i] = userScore{
			ID:    userID,
			Score: 0,
		}
	}
	return userScores
}

func (u userScores) ids() []int64 {
	ids := make([]int64, len(u))
	for i, userScore := range u {
		ids[i] = userScore.ID
	}
	return ids
}

func (u userScores) chooseIDs() []int64 {
	var highestScore int64
	groupedIDs := make(map[int64][]int64)
	for _, userScore := range u {
		groupedIDs[userScore.Score] = append(groupedIDs[userScore.Score], userScore.ID)
		if userScore.Score > highestScore {
			highestScore = userScore.Score
		}
	}
	return groupedIDs[highestScore]
}

type autoAssignTaskProcessor func(userIDs userScores) (userScores, error)

func autoAssignTaskProcessRates(
	projectUsersMap map[int64]projects.User,
	reasoning *string,
	logger *slog.Logger,
) autoAssignTaskProcessor {
	type userCost struct {
		ID   int64
		Cost twapi.Money
	}
	logger = logger.With(
		slog.String("subAction", "processRates"),
	)
	return func(userScores userScores) (userScores, error) {
		var userCosts []userCost
		distinctCosts := make(map[twapi.Money]struct{})
		for _, userScore := range userScores {
			user, ok := projectUsersMap[userScore.ID]
			if !ok {
				continue
			}
			if user.Cost == nil || *user.Cost == 0 || len(userCosts) == 0 {
				var cost twapi.Money
				if user.Cost != nil {
					cost = *user.Cost
				}
				userCosts = append(userCosts, userCost{
					ID:   user.ID,
					Cost: cost,
				})
				distinctCosts[cost] = struct{}{}
				continue
			}
			for i := range userCosts {
				if userCosts[i].Cost > *user.Cost {
					userCosts = slices.Insert(userCosts, i, userCost{
						ID:   user.ID,
						Cost: *user.Cost,
					})
					distinctCosts[*user.Cost] = struct{}{}
					break
				}
			}
		}
		weight := len(distinctCosts) + 1
		userCostsWeights := make(map[int64]int, len(userCosts))
		for i, userCost := range userCosts {
			if i > 0 && userCosts[i-1].Cost == userCost.Cost {
				userCostsWeights[userCost.ID] = weight
			} else {
				weight--
				userCostsWeights[userCost.ID] = weight
			}
		}
		var changed bool
		for i, userScore := range userScores {
			weight, ok := userCostsWeights[userScore.ID]
			if !ok {
				continue
			}
			userScore.Score += int64(weight)
			userScores[i] = userScore
			changed = true
			logger.Debug("user score changed",
				slog.Int64("userID", userScore.ID),
				slog.Int("delta", weight),
				slog.Int64("score", userScore.Score),
			)
		}
		if changed && reasoning != nil {
			if *reasoning != "" {
				*reasoning += " "
			}
			*reasoning += "Concerns over user cost significantly impacted the decision."
		}
		return userScores, nil
	}
}

func autoAssignTaskProcessWorkload(
	ctx context.Context,
	taskData webhook.TaskData,
	resources *config.Resources,
	reasoning *string,
	logger *slog.Logger,
) autoAssignTaskProcessor {
	logger = logger.With(
		slog.String("subAction", "processWorkload"),
	)
	return func(userScores userScores) (userScores, error) {
		if taskData.Task.StartDate == nil || taskData.Task.DueDate == nil {
			// without a window period, we can't calculate the workload
			return userScores, nil
		}

		var workloadRequest projects.WorkloadRequest
		workloadRequest.Filters.StartDate = *taskData.Task.StartDate
		workloadRequest.Filters.EndDate = *taskData.Task.DueDate
		workloadRequest.Filters.UserIDs = userScores.ids()
		workloadRequest.Filters.PageSize = int64(len(workloadRequest.Filters.UserIDs))
		workloadRequest.Filters.Include = []projects.WorkloadGetRequestSideload{
			projects.WorkloadGetRequestSideloadWorkingHourEntries,
		}

		workloadResponse, err := projects.WorkloadGet(ctx, resources.TeamworkEngine, workloadRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to load workload: %w", err)
		}

		availableUserIDs := make(map[int64]struct{})
		for _, user := range workloadResponse.Workload.Users {
			userIDStr := strconv.FormatInt(user.ID, 10)
			var workingHoursID int64
			if relationship := workloadResponse.Included.Users[userIDStr].WorkingHour; relationship != nil {
				workingHoursID = relationship.ID
			}

			var availableHours float64
			for date, dateData := range user.Dates {
				var workingHours *float64
				for _, entry := range workloadResponse.Included.WorkingHoursEntries {
					if entry.WorkingHour.ID != workingHoursID {
						continue
					}
					if weekday := strings.ToLower(time.Time(date).Weekday().String()); entry.Weekday == weekday {
						workingHours = &entry.TaskHours
						break
					}
				}
				if workingHours == nil {
					workingHours = func() *float64 {
						var v float64
						if workloadResponse.Included.Users != nil {
							//nolint:staticcheck
							v = workloadResponse.Included.Users[userIDStr].LengthOfDay
						}
						if v == 0 {
							// last resort to a default value
							v = 8 // hours
						}
						return &v
					}()
				}
				if !dateData.UnavailableDay {
					availableHours += *workingHours - (float64(dateData.CapacityMinutes) / 60)
				}
			}

			if availableHours > float64(taskData.Task.EstimatedMinutes)/60 {
				availableUserIDs[user.ID] = struct{}{}
			}
		}
		var changed bool
		for i, userScore := range userScores {
			if _, ok := availableUserIDs[userScore.ID]; !ok {
				continue
			}
			userScore.Score += int64(len(userScores))
			userScores[i] = userScore
			changed = true
			logger.Debug("user score changed",
				slog.Int64("userID", userScore.ID),
				slog.Int("delta", len(userScores)),
				slog.Int64("score", userScore.Score),
			)
		}
		if changed && reasoning != nil {
			if *reasoning != "" {
				*reasoning += " "
			}
			*reasoning += "Workload was a key consideration in the decision-making process."
		}
		return userScores, nil
	}
}

type skills []projects.Skill

func (s skills) toMap() map[int64]projects.Skill {
	m := make(map[int64]projects.Skill, len(s))
	for _, skill := range s {
		m[skill.ID] = skill
	}
	return m
}

func loadSkills(ctx context.Context, resources *config.Resources) (skills, error) {
	skillsNext, err := twapi.Iterate[projects.SkillListRequest, *projects.SkillListResponse](
		ctx,
		resources.TeamworkEngine,
		projects.NewSkillListRequest(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build skills iterator: %w", err)
	}

	var skills skills
	for {
		skillsResponse, hasSkillsNext, err := skillsNext()
		if err != nil {
			return nil, fmt.Errorf("failed to list skills: %w", err)
		}
		if skillsResponse == nil {
			break
		}
		skills = append(skills, skillsResponse.Skills...)
		if !hasSkillsNext {
			break
		}
	}
	return skills, nil
}

type jobRoles []projects.JobRole

func (j jobRoles) toMap() map[int64]projects.JobRole {
	m := make(map[int64]projects.JobRole, len(j))
	for _, jobRole := range j {
		m[jobRole.ID] = jobRole
	}
	return m
}

func loadJobRoles(ctx context.Context, resources *config.Resources) (jobRoles, error) {
	jobRolesNext, err := twapi.Iterate[projects.JobRoleListRequest, *projects.JobRoleListResponse](
		ctx,
		resources.TeamworkEngine,
		projects.NewJobRoleListRequest(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build job roles iterator: %w", err)
	}

	var jobRoles jobRoles
	for {
		jobRolesResponse, hasJobRolesNext, err := jobRolesNext()
		if err != nil {
			return nil, fmt.Errorf("failed to list job roles: %w", err)
		}
		if jobRolesResponse == nil {
			break
		}
		jobRoles = append(jobRoles, jobRolesResponse.JobRoles...)
		if !hasJobRolesNext {
			break
		}
	}
	return jobRoles, nil
}

type projectUsers []projects.User

func (p projectUsers) toMap() map[int64]projects.User {
	m := make(map[int64]projects.User, len(p))
	for _, user := range p {
		m[user.ID] = user
	}
	return m
}

func loadProjectUsers(ctx context.Context, resources *config.Resources, projectID int64) (projectUsers, error) {
	userListRequest := projects.NewUserListRequest()
	userListRequest.Path.ProjectID = projectID

	projectUsersNext, err := twapi.Iterate[projects.UserListRequest, *projects.UserListResponse](
		ctx,
		resources.TeamworkEngine,
		userListRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build project users iterator: %w", err)
	}

	var projectUsers projectUsers
	for {
		projectUsersResponse, hasProjectUsersNext, err := projectUsersNext()
		if err != nil {
			return nil, fmt.Errorf("failed to list project users: %w", err)
		}
		if projectUsersResponse == nil {
			break
		}
		projectUsers = append(projectUsers, projectUsersResponse.Users...)
		if !hasProjectUsersNext {
			break
		}
	}
	return projectUsers, nil
}
