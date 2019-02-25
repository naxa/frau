package epic

import (
	"context"
	"log"

	"github.com/google/go-github/github"

	"errors"

	"fmt"

	"github.com/student-kyushu/frau/input"
	"github.com/student-kyushu/frau/operation"
	"github.com/student-kyushu/frau/queue"
)

func (c *AcceptCommand) AcceptChangesetOnReview(ctx context.Context, ev *github.PullRequestReviewEvent, cmd *input.AcceptChangeOnReview) (bool, error) {
	log.Printf("info: Start: merge the pull request by %v\n", *ev.Review.ID)
	defer log.Printf("info: End: merge the pull request by %v\n", *ev.Review.ID)

	if c.BotName != cmd.BotName() {
		log.Println("info: this command works only if target user is actual our bot.")
		return false, nil
	}

	sender := *ev.Sender.Login
	log.Printf("debug: command is sent from %v\n", sender)

	if !c.Info.IsReviewer(sender) {
		log.Printf("info: %v is not an reviewer registered to this bot.\n", sender)
		return false, nil
	}

	return c.acceptChangesetOnReview(ctx, ev, cmd)
}

func (c *AcceptCommand) acceptChangesetOnReview(ctx context.Context, ev *github.PullRequestReviewEvent, cmd input.AcceptChangesetCommand) (bool, error) {
	sender := *ev.Sender.Login

	client := c.Client
	issueSvc := client.Issues

	repoOwner := c.Owner
	repoName := c.Name
	pullrequestNumber := *ev.PullRequest.Number
	log.Printf("debug: issue number is %v\n", pullrequestNumber)

	currentLabels := operation.GetLabelsByIssue(ctx, issueSvc, repoOwner, repoName, pullrequestNumber)
	if currentLabels == nil {
		return false, nil
	}

	doNotMergeLabel, _, err := issueSvc.GetLabel(ctx, repoOwner, repoName, "S-do-not-merge")
	if err != nil {
		log.Println("info: could not find label `S-do-not-merge`")
		return false, err
	}

	if containsDoNotMerge := containsLabel(currentLabels, doNotMergeLabel); containsDoNotMerge {
		comment := fmt.Sprint("warning: forbid merging by label `S-do-not-merge`")
		if ok := operation.AddComment(ctx, issueSvc, repoOwner, repoName, pullrequestNumber, comment); !ok {
			log.Println("info: could not create the comment to declare the head is approved.")
			return false, nil
		}
		log.Println("info: forbid merging by label `S-do-not-merge`")
		return false, nil
	}

	labels := operation.AddAwaitingMergeLabel(currentLabels)

	// https://github.com/nekoya/popuko/blob/master/web.py
	_, _, err = issueSvc.ReplaceLabelsForIssue(ctx, repoOwner, repoName, pullrequestNumber, labels)
	if err != nil {
		log.Println("info: could not change labels by the issue")
		return false, err
	}

	prSvc := client.PullRequests
	pr, _, err := prSvc.Get(ctx, repoOwner, repoName, pullrequestNumber)
	if err != nil {
		log.Println("info: could not fetch the pull request information.")
		return false, err
	}

	headSha := *pr.Head.SHA
	if ok := commentApprovedSha(ctx, cmd, issueSvc, repoOwner, repoName, pullrequestNumber, headSha, sender); !ok {
		log.Println("info: could not create the comment to declare the head is approved.")
		return false, err
	}

	if c.Info.EnableAutoMerge {
		qHandle := c.AutoMergeRepo.Get(repoOwner, repoName)
		if qHandle == nil {
			log.Println("error: cannot get the queue handle")
			return false, errors.New("error: cannot get the queue handle")
		}

		qHandle.Lock()
		defer qHandle.Unlock()

		q := qHandle.Load()

		item := &queue.AutoMergeQueueItem{
			PullRequest: pullrequestNumber,
			PrHead:      headSha,
		}
		ok, mutated := queuePullReq(q, item)
		if !ok {
			return false, errors.New("error: we cannot recover the error")
		}

		if mutated {
			q.Save()
		}

		if q.HasActive() {
			commentAsPostponed(ctx, issueSvc, repoOwner, repoName, pullrequestNumber)
			return true, nil
		}

		if next := q.Front(); next != item {
			commentAsPostponed(ctx, issueSvc, repoOwner, repoName, pullrequestNumber)
		}

		tryNextItem(ctx, client, repoOwner, repoName, q, c.Info.AutoBranchName)
	}

	log.Printf("info: complete merge the pull request %v\n", pullrequestNumber)
	return true, nil
}
