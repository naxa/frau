package epic

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/go-github/v24/github"
	"github.com/naxa-jp/frau/operation"
)

func AssignReviewer(ctx context.Context, client *github.Client, ev *github.IssueCommentEvent, assignees []string) (bool, error) {
	log.Printf("info: Start: assign the reviewer by %v\n", *ev.Comment.ID)
	defer log.Printf("info: End: assign the reviewer by %v\n", *ev.Comment.ID)

	issueSvc := client.Issues

	repoOwner := *ev.Repo.Owner.Login
	log.Printf("debug: repository owner is %v\n", repoOwner)
	repo := *ev.Repo.Name
	log.Printf("debug: repository name is %v\n", repo)

	issue := *ev.Issue
	issueNum := *ev.Issue.Number
	log.Printf("debug: issue number is %v\n", issueNum)

	// https://godoc.org/github.com/google/go-github/github#Issue
	// 	> If PullRequestLinks is nil, this is an issue, and if PullRequestLinks is not nil, this is a pull request.
	if issue.PullRequestLinks == nil {
		log.Println("info: the issue is pull request")
		return false, nil
	}

	currentLabels := operation.GetLabelsByIssue(ctx, issueSvc, repoOwner, repo, issueNum)
	if currentLabels == nil {
		return false, nil
	}

	log.Printf("debug: assignees is %v\n", assignees)

	presentAssignees := issue.Assignees
	if len(presentAssignees) != 0 {
		var presentAssigneesName []string
		for i := 0; i < len(presentAssignees); i++ {
			presentAssigneesName = append(presentAssigneesName, *presentAssignees[i].Login)
		}
		_, _, err := issueSvc.RemoveAssignees(ctx, repoOwner, repo, issueNum, presentAssigneesName)
		if err != nil {
			log.Println("info: could not remove assignees")
			return false, err
		}
	}

	_, _, err := issueSvc.AddAssignees(ctx, repoOwner, repo, issueNum, assignees)
	if err != nil {
		log.Println("info: could not change assignees.")
		return false, err
	}

	labels := operation.AddAwaitingReviewLabel(currentLabels)
	_, _, err = issueSvc.ReplaceLabelsForIssue(ctx, repoOwner, repo, issueNum, labels)
	if err != nil {
		log.Println("info: could not change labels.")
		return false, err
	}

	log.Println("info: Complete assign the reviewer with no errors.")

	return true, nil
}

func AssignReviewerFromPR(ctx context.Context, client *github.Client, ev *github.PullRequestEvent, assignees []string) (bool, error) {
	log.Printf("info: Start: assign the reviewer by %v\n", *ev.Number)
	defer log.Printf("info: End: assign the reviewer by %v\n", *ev.Number)

	isRandom := false

	issueSvc := client.Issues
	repoSvc := client.Repositories

	repoOwner := *ev.Repo.Owner.Login
	log.Printf("debug: repository owner is %v\n", repoOwner)
	repo := *ev.Repo.Name
	log.Printf("debug: repository name is %v\n", repo)

	pullReqNum := *ev.Number
	log.Printf("debug: pull request number is %v\n", pullReqNum)

	currentLabels := operation.GetLabelsByIssue(ctx, issueSvc, repoOwner, repo, pullReqNum)
	if currentLabels == nil {
		return false, nil
	}

	if assignees == nil {
		log.Println("debug: there are no requested assignees")
		ok, owners := fetchOwnersFile(ctx, repoSvc, repoOwner, repo)
		if !ok {
			log.Println("error: could not handle OWNERS file.")
			return false, nil
		}
		reviewers := owners.ReviewersList()
		if reviewers == nil {
			log.Println("info: could not find any reviewers")
			return false, nil
		}
		sender := *ev.Sender.Login
		_, index := contains(reviewers, sender)
		if index != -1 && len(reviewers) != 1 {
			reviewers = remove(reviewers, index)
		}
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(reviewers))
		assignees = append(assignees, reviewers[i])
		isRandom = true
	}

	log.Printf("debug: assignees is %v\n", assignees)

	_, _, err := issueSvc.AddAssignees(ctx, repoOwner, repo, pullReqNum, assignees)
	if err != nil {
		log.Println("info: could not change assignees.")
		return false, err
	}

	labels := operation.AddAwaitingReviewLabel(currentLabels)
	_, _, err = issueSvc.ReplaceLabelsForIssue(ctx, repoOwner, repo, pullReqNum, labels)
	if err != nil {
		log.Println("info: could not change labels.")
		return false, err
	}

	if isRandom {
		comment := fmt.Sprint(":eggplant: I picked a reviewer randomly, you can use `r?` to overwrite.")
		if ok := operation.AddComment(ctx, issueSvc, repoOwner, repo, pullReqNum, comment); !ok {
			log.Println("info: could not create the comment about assigning a reviewer randomly")
			return false, nil
		}
	}

	log.Println("info: Complete assign the reviewer with no errors.")

	return true, nil
}

func remove(s []string, i int) []string {
	if i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}
