package lead

import (
	"context"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
)

type Task struct {
	// ID is the unique id for the board.
	ID int64
	// Name is the name of task in question.
	Name string
	// Start is when the task was started
	Start time.Time
	// End is when the task was completed
	End time.Time
}

type CreateParams struct {
	// Name is the name of the task.
	Name string
}

// Create creates a new task.
// encore:api public
func Create(ctx context.Context, params *CreateParams) (*Task, error) {
	b := &Task{Name: params.Name, Start: time.Now()}
	err := sqldb.QueryRow(ctx, `
		INSERT INTO task (name, start_time)
		VALUES ($1, $2)
		RETURNING id
	`, params.Name, b.Start).Scan(&b.ID)
	if err != nil {
		return nil, fmt.Errorf("could not create task: %v", err)
	}

	return b, nil
}

type CompleteParams struct {
	// ID is the id of the task to complete
	ID int64
}

// Complete marks the task as finished by setting an end time.
// encore:api public
func Compete(ctx context.Context, params *CompleteParams) (*Task, error) {
	b := &Task{ID: params.ID, End: time.Now()}

	err := sqldb.QueryRow(ctx, `
		UPDATE task
		SET end_time = $1
		WHERE id = $2
		RETURNING id, name, start_time, end_time;
	`, b.End, b.ID).Scan(&b.ID, &b.Name, &b.Start, &b.End)
	if err != nil {
		return nil, fmt.Errorf("could not complete task: %w", err)
	}

	return b, nil
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
			return nil, fmt.Errorf("could not scan: %v", err)
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
