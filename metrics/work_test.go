package metrics

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"
	"time"
)

//go:embed testdata/first_commit_new_branch.json
var newBrachJSON []byte

func TestCreateBranch(t *testing.T) {
	ctx := context.Background()

	var push Push

	err := json.Unmarshal(newBrachJSON, &push)
	if err != nil {
		t.Fatal(err)
	}

	pushResp, err := GitPush(ctx, &push)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Get(ctx, &WorkID{pushResp.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = time.Parse(time.RFC3339, "2021-05-13T09:09:18+02:00")
	if err != nil {
		t.Fatal(err)
	}

	// if work.Start != shouldStart {
	// 	t.Fatal("work started at the wrong time")
	// }
}
