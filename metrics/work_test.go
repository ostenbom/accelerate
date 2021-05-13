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

//go:embed testdata/pr_opened.json
var newPRJSON []byte

func getPush(t *testing.T, bytes []byte) Push {
	var push Push

	err := json.Unmarshal(bytes, &push)
	if err != nil {
		t.Fatal(err)
	}

	return push
}

func getPR(t *testing.T, bytes []byte) PullRequest {
	var pr PullRequest

	err := json.Unmarshal(bytes, &pr)
	if err != nil {
		t.Fatal(err)
	}

	return pr
}

func TestCreateBranch(t *testing.T) {
	ctx := context.Background()
	push := getPush(t, newBrachJSON)

	pushResp, err := GitPush(ctx, &push)
	if err != nil {
		t.Fatal(err)
	}

	work, err := Get(ctx, &WorkID{pushResp.ID})
	if err != nil {
		t.Fatal(err)
	}

	shouldStart, err := time.Parse(time.RFC3339, "2021-05-13T09:09:18+02:00")
	if err != nil {
		t.Fatal(err)
	}

	if work.Start != shouldStart {
		t.Fatalf("work started at the wrong time, %v, %v", work.Start, shouldStart)
	}

	if work.Branch != "lead-test" {
		t.Fatalf("work had the wrong branch name: %s", work.Branch)
	}
}

func TestCreateBranchPR(t *testing.T) {
	ctx := context.Background()
	push := getPush(t, newBrachJSON)
	pr := getPR(t, newPRJSON)

	_, err := GitPush(ctx, &push)
	if err != nil {
		t.Fatal(err)
	}

	prResp, err := GitPullRequest(ctx, &pr)
	if err != nil {
		t.Fatal(err)
	}

	work, err := Get(ctx, &WorkID{prResp.ID})
	if err != nil {
		t.Fatal(err)
	}

	if work.PullRequest != 1 {
		t.Fatal("work item did not get associated with the correct pull request")
	}
}
