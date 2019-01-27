package epic

import (
	"context"
	"log"

	"github.com/JohnTitor/frau/operation"
	"github.com/google/go-github/github"
)

func AssignReviewer(ctx context.Context, client *github.Client, ev *github.IssueCommentEvent, reviewers []string) (bool, error) {
	log.Printf("info: Start: assign the reviewer by %v\n", *ev.Comment.ID)
	defer log.Printf("info: End: assign the reviewer by %v\n", *ev.Comment.ID)

	issueSvc := client.Issues
	pullReqSvc := client.PullRequests

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

	log.Printf("debug: assignees is %v\n", reviewers)

	if containsPullReqOwner, index := Contains(reviewers, *ev.Issue.User.Login); containsPullReqOwner {
		pullReqOwner := []string{reviewers[index]}
		removeSliceElement(reviewers, reviewers[index])
		_, _, err := issueSvc.AddAssignees(ctx, repoOwner, repo, issueNum, pullReqOwner)
		if err != nil {
			log.Println("info: could not change asignees.")
			return false, err
		}
	}

	requestedReviewers := github.ReviewersRequest{Reviewers: reviewers}

	_, _, err := pullReqSvc.RequestReviewers(ctx, repoOwner, repo, issueNum, requestedReviewers)
	if err != nil {
		log.Println("info: could not change reviewers.")
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

func AssignReviewerFromPR(ctx context.Context, client *github.Client, ev *github.PullRequestEvent, reviewers []string) (bool, error) {
	log.Printf("info: Start: assign the reviewer by %v\n", *ev.Number)
	defer log.Printf("info: End: assign the reviewer by %v\n", *ev.Number)

	issueSvc := client.Issues
	pullReqSvc := client.PullRequests

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

	log.Printf("debug: reviewers is %v\n", reviewers)

	if containsPullReqOwner, index := Contains(reviewers, *ev.PullRequest.User.Login); containsPullReqOwner {
		pullReqOwner := []string{reviewers[index]}
		removeSliceElement(reviewers, reviewers[index])
		_, _, err := issueSvc.AddAssignees(ctx, repoOwner, repo, pullReqNum, pullReqOwner)
		if err != nil {
			log.Println("info: could not change asignees.")
			return false, err
		}
	}

	requestedReviewers := github.ReviewersRequest{Reviewers: reviewers}

	_, _, err := pullReqSvc.RequestReviewers(ctx, repoOwner, repo, pullReqNum, requestedReviewers)
	if err != nil {
		log.Println("info: could not change reviewers.")
		return false, err
	}

	labels := operation.AddAwaitingReviewLabel(currentLabels)
	_, _, err = issueSvc.ReplaceLabelsForIssue(ctx, repoOwner, repo, pullReqNum, labels)
	if err != nil {
		log.Println("info: could not change labels.")
		return false, err
	}

	log.Println("info: Complete assign the reviewer with no errors.")

	return true, nil
}

func removeSliceElement(strings []string, search string) []string {
	result := []string{}
	for _, v := range strings {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
