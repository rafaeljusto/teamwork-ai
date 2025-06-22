package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/rafaeljusto/teamwork-ai/internal/config"
	"github.com/rafaeljusto/teamwork-ai/internal/twapi/activity"
)

// SummarizeActivitiesOptions contains the options for the SummarizeActivities
// function.
type SummarizeActivitiesOptions struct {
	projectID int64
	startDate time.Time
	endDate   time.Time
}

// SummarizeActivitiesOption is a function that modifies the
// SummarizeActivitiesOptions struct. It allows for optional configuration of
// the SummarizeActivities function.
type SummarizeActivitiesOption func(*SummarizeActivitiesOptions)

// WithSummarizeActivitiesPeriod sets the start and end dates for the
// SummarizeActivities function. The start date is the beginning of the period
// to summarize, and the end date is the end of the period to summarize. The end
// date cannot be before or equal the start date, and the period cannot be over
// 365 days.
func WithSummarizeActivitiesPeriod(startDate, endDate time.Time) SummarizeActivitiesOption {
	return func(o *SummarizeActivitiesOptions) {
		o.startDate = startDate
		o.endDate = endDate
	}
}

// WithSummarizeActivitiesProjectID sets the project ID for the
// SummarizeActivities function. If the project ID is set, the summary will be
// for the specified project. If not set, the summary will be for all activities
// within the specified period.
func WithSummarizeActivitiesProjectID(projectID int64) SummarizeActivitiesOption {
	return func(o *SummarizeActivitiesOptions) {
		o.projectID = projectID
	}
}

// SummarizeActivities summarizes the activities for a given period. It uses the
// start and end dates to filter the activities and generate a summary. The
// summary is returned as a string. It's possible to specify a project ID to
// summarize activities for a specific project. By default, the summary is for
// all activities within the last 365 days.
func SummarizeActivities(
	ctx context.Context,
	resources *config.Resources,
	optFuncs ...SummarizeActivitiesOption,
) (string, error) {
	options := SummarizeActivitiesOptions{
		startDate: time.Now().AddDate(-1, 0, 0),
		endDate:   time.Now(),
	}
	for _, optFunc := range optFuncs {
		optFunc(&options)
	}

	switch {
	case options.startDate.IsZero(), options.endDate.IsZero():
		return "", fmt.Errorf("startDate and endDate are required")
	case !options.endDate.After(options.startDate):
		return "", fmt.Errorf("startDate must be before endDate")
	case options.startDate.After(time.Now()):
		return "", fmt.Errorf("startDate must be before now")
	case options.endDate.Sub(options.startDate) > 365*24*time.Hour:
		return "", fmt.Errorf("startDate and endDate must be within 1 year")
	}

	// TODO(rafaeljusto): add support for pagination
	var multiple activity.Multiple
	multiple.Request.Path.ProjectID = options.projectID
	multiple.Request.Filters.StartDate = options.startDate
	multiple.Request.Filters.EndDate = options.endDate
	if err := resources.TeamworkEngine.Do(ctx, &multiple); err != nil {
		return "", fmt.Errorf("failed to load activities: %w", err)
	}

	if len(multiple.Response.Activities) == 0 {
		return "No activity during this period", nil
	}

	summary, err := resources.Agentic.SummarizeActivities(ctx, multiple.Response.Activities)
	if err != nil {
		return "", fmt.Errorf("failed to summarize activities: %w", err)
	}
	return summary, nil
}
