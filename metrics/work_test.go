package metrics

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:embed testdata/first_commit_new_branch.json
var newBrachJSON []byte

//go:embed testdata/pr_opened.json
var newPRJSON []byte

//go:embed testdata/pr_closed_rebase.json
var closedPRRebaseJSON []byte

//go:embed testdata/pr_closed_merged.json
var closedPRMergeJSON []byte

var _ = Describe("GitHub Work Flow", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("a commit has been pushed to a branch", func() {
		var push Push
		var pushResp *WorkID
		BeforeEach(func() {
			err := json.Unmarshal(newBrachJSON, &push)
			Expect(err).NotTo(HaveOccurred())

			pushResp, err = GitPush(ctx, &push)
			Expect(err).NotTo(HaveOccurred())
		})

		It("saves branch name and time", func() {
			work, err := Get(ctx, &WorkID{pushResp.ID})
			Expect(err).NotTo(HaveOccurred())

			shouldStart, err := time.Parse(time.RFC3339, "2021-05-13T09:09:18+02:00")
			Expect(err).NotTo(HaveOccurred())

			Expect(work.Start.Unix()).To(Equal(shouldStart.Unix()))
			Expect(work.Branch).To(Equal("lead-test"))
		})

		Context("a pull request has been created", func() {
			var pr PullRequest
			var prResp *WorkID
			BeforeEach(func() {
				ctx := context.Background()

				err := json.Unmarshal(newPRJSON, &pr)
				Expect(err).NotTo(HaveOccurred())

				prResp, err = GitPullRequest(ctx, &pr)
				Expect(err).NotTo(HaveOccurred())
			})

			It("associates the pull request with the correct work", func() {
				work, err := Get(ctx, &WorkID{prResp.ID})
				Expect(err).NotTo(HaveOccurred())

				Expect(work.PullRequest).To(Equal(1))
			})

			Context("the pull request is closed with a rebase", func() {
				var closePR PullRequest
				var closePRResp *WorkID
				BeforeEach(func() {
					ctx := context.Background()

					err := json.Unmarshal(closedPRRebaseJSON, &closePR)
					Expect(err).NotTo(HaveOccurred())

					closePRResp, err = GitPullRequest(ctx, &closePR)
					Expect(err).NotTo(HaveOccurred())
				})

				It("associates the merge commit and time", func() {
					work, err := Get(ctx, &WorkID{closePRResp.ID})
					Expect(err).NotTo(HaveOccurred())

					Expect(work.MergeCommit).To(Equal("9bd73f28b5ed4597123de1d8ecf509078d99bc84"))

					expectedMergeTime, err := time.Parse(time.RFC3339, "2021-05-13T07:26:12Z")
					Expect(err).NotTo(HaveOccurred())
					Expect(work.Merged.Unix()).To(Equal(expectedMergeTime.Unix()))
				})

				It("can associate the correct deployment time", func() {
					deployTime := time.Now()
					_, err := SetDeployed(ctx, &DeployedParams{"9bd73f28b5ed4597123de1d8ecf509078d99bc84", deployTime})

					work, err := Get(ctx, &WorkID{closePRResp.ID})
					Expect(err).NotTo(HaveOccurred())

					Expect(work.Deployed.Unix()).To(Equal(deployTime.Unix()))
				})
			})

			Context("the pull request is closed with a merge commit", func() {
				var closePR PullRequest
				var closePRResp *WorkID
				BeforeEach(func() {
					ctx := context.Background()

					err := json.Unmarshal(closedPRMergeJSON, &closePR)
					Expect(err).NotTo(HaveOccurred())

					closePRResp, err = GitPullRequest(ctx, &closePR)
					Expect(err).NotTo(HaveOccurred())
				})

				It("associates the merge commit and time", func() {
					work, err := Get(ctx, &WorkID{closePRResp.ID})
					Expect(err).NotTo(HaveOccurred())

					Expect(work.MergeCommit).To(Equal("ecc81403853a621bea766bad50d1fb907d1b2689"))

					expectedMergeTime, err := time.Parse(time.RFC3339, "2021-05-13T07:41:12Z")
					Expect(err).NotTo(HaveOccurred())
					Expect(work.Merged.Unix()).To(Equal(expectedMergeTime.Unix()))
				})
			})
		})
	})
})

func TestWork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Work Suite")
}
