package metrics

import (
	"context"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
)

type Work struct {
	// ID is the unique id for the board
	ID int64
	// Name is the name of task in question
	Branch string
	// PullRequest is the id of the associated pull request in github
	PullRequest int
	// MergeCommit is the sha of main branch commit associated with this work
	MergeCommit string
	// Start is when the task was started
	Start time.Time
	// Merged is when the task merged into the main branch
	Merged time.Time
	// Deployed is when the task completed its deployment to production
	Deployed time.Time
}

type Push struct {
	// Name is the name of the task.
	Ref string `json:"ref"`
}

type WorkID struct {
	ID int64
}

// GitPush Accepts Github Webhooks of the "push" type
// encore:api public
func GitPush(ctx context.Context, push *Push) (*WorkID, error) {
	w := &Work{Branch: push.Ref, Start: time.Now()}

	err := sqldb.QueryRow(ctx, `
		INSERT INTO work (branch, start_time)
		VALUES ($1, $2)
		RETURNING id;
	`, &w.Branch, &w.Start).Scan(&w.ID)
	if err != nil {
		return nil, fmt.Errorf("could not create work: %w", err)
	}

	return &WorkID{w.ID}, nil
}

// Get retreives a work item with a specific ID
// encore:api public
func Get(ctx context.Context, params *WorkID) (*Work, error) {
	w := &Work{}

	// id, branch, pull_request, merge_commit, start_time, merged_time, deployed_time
	err := sqldb.QueryRow(ctx, `
		SELECT
			id, branch
		FROM work
		WHERE id = $1
		LIMIT 1;
	`, params.ID).Scan(&w.ID, &w.Branch)
	// `, params.ID).Scan(&w.ID, &w.Branch, &w.PullRequest, &w.MergeCommit, &w.Start, &w.Merged, &w.Deployed)
	if err != nil {
		return nil, fmt.Errorf("could not get work: %w", err)
	}

	return w, nil
}

type AverageParams struct {
	Since time.Time
}

type AverageResponse struct {
	Time float64
}

// Average returns the avearge lead time since a certain time
// encore:api public
func Average(ctx context.Context, params *AverageParams) (*AverageResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT start_time, end_time
		FROM task
		WHERE start_time IS NOT NULL and end_time IS NOT NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("could not get tasks: %w", err)
	}
	defer rows.Close()

	var sumTaskMinutes float64
	var numTasks float64

	for rows.Next() {
		var start time.Time
		var end time.Time

		if err := rows.Scan(&start, &end); err != nil {
			return nil, fmt.Errorf("could not scan: %w", err)
		}

		taskDuration := end.Sub(start)
		sumTaskMinutes += taskDuration.Minutes()
		numTasks++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %w", err)
	}

	return &AverageResponse{Time: sumTaskMinutes / numTasks}, nil
}
