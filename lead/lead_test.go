package lead

import (
	"context"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	ctx := context.Background()
	task, err := Create(ctx, &CreateParams{
		Name: t.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	_, err = Compete(ctx, &CompleteParams{ID: task.ID})
	if err != nil {
		t.Fatal(err)
	}

	lead, err := Lead(ctx, &LeadParams{})
	if err != nil {
		t.Fatal(err)
	}
	// 100ms = 0,001666667minutes
	if lead.Time < 0.0012 || lead.Time > 0.0020 {
		t.Error("lead time calculation was off")
	}
}
