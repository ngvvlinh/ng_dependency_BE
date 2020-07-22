package webhook

import (
	"context"
	"fmt"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
)

// Facebook feed is any action on page (create or update a post, make comment,
// any reaction ....)
func (wh *Webhook) handleFeed(ctx context.Context, feed WebhookMessages) error {
	for _, entry := range feed.Entry {
		// first we need check page is active then get token of page
		externalPageID := entry.ID
		createdTime := time.Unix(int64(entry.Time), 0)

		accessToken, err := wh.getPageAccessToken(ctx, externalPageID)
		if err != nil {
			return err
		}

		for _, change := range entry.Changes {
			switch {
			case change.IsAdminPost(externalPageID):
				return wh.handleFeedPost(ctx, externalPageID, change, createdTime, accessToken)
			case change.IsComment():
				return wh.handleFeedComment(ctx, externalPageID, change, createdTime, accessToken)
			}
		}
	}
	return nil
}

func (wh *Webhook) handleFeedPost(ctx context.Context, extPageID string, feedChange FeedChange, createdTime time.Time, accessToken string) error {
	postID := feedChange.Value.PostID
	{
		// Sometimes facebook may send many requests for one feed. E.g: If a
		// user posts a post containing 4 images, facebook may calls more than 5
		// requests at the same time. In this case, we can build full post and
		// all child posts from one request, so just hold one and ignore the
		// rest.
		actionKey := fmt.Sprintf("FEED_%v_%v_%v", extPageID, postID, feedChange.Value.From.ID)
		if wh.faboRedis.IsExist(actionKey) {
			return nil
		}
		wh.faboRedis.SetKey(actionKey, true)
		defer wh.faboRedis.DelKey(actionKey)
	}

	externalPost, err := wh.getExternalPost(ctx, postID)
	if err != nil {
		return err
	}

	// if post does not exist in db, create it
	if externalPost == nil {
		post, err := wh.fbClient.CallAPIGetPost(postID, accessToken)
		if err != nil {
			return err
		}

		if err := wh.createParentAndChildPosts(extPageID, createdTime, ctx, post); err != nil {
			return err
		}
	}

	if !feedChange.IsEdited() {
		return nil
	}

	switch {
	case feedChange.IsOnChildPost():
		// At this step, the parent of this post always exists in db. Check
		// if the child post does not exist, create it.
		childPostID := fmt.Sprintf("%v_%v", extPageID, feedChange.Value.PhotoID)
		childPost, err := wh.getExternalPost(ctx, childPostID)
		if err != nil {
			return err
		}

		if childPost == nil {
			createExternalPostCmd := &fbmessaging.CreateFbExternalPostArgs{
				ID:                  cm.NewID(),
				ExternalPageID:      extPageID,
				ExternalID:          childPostID,
				ExternalParentID:    postID,
				ExternalFrom:        externalPost.ExternalFrom,
				ExternalPicture:     feedChange.Value.Link,
				ExternalIcon:        externalPost.ExternalIcon,
				ExternalMessage:     feedChange.Value.Message,
				ExternalCreatedTime: externalPost.CreatedAt,
				ExternalUpdatedTime: externalPost.UpdatedAt,
			}
			return wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateFbExternalPostsCommand{
				FbExternalPosts: []*fbmessaging.CreateFbExternalPostArgs{createExternalPostCmd},
			})
		} else {
			return wh.updateFeedPostMessage(ctx, postID, feedChange.Value.Message)
		}

	case feedChange.IsOnParentPost():
		return wh.updateFeedPostMessage(ctx, postID, feedChange.Value.Message)
	}
	return nil
}

func (wh *Webhook) updateFeedPostMessage(ctx context.Context, postID string, message string) error {
	cmdUpdate := &fbmessaging.UpdateFbPostMessageCommand{
		ExternalPostID: postID,
		Message:        message,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, cmdUpdate); err != nil {
		return err
	}
	return nil
}

func (wh *Webhook) handleFeedComment(ctx context.Context, extPageID string, feedChange FeedChange, createdTime time.Time, accessToken string) error {
	postID := feedChange.Value.PostID
	externalPost, err := wh.getExternalPost(ctx, postID)
	if err != nil {
		return err
	}

	if externalPost == nil {
		// Case post not exists in db and action on this comment is remove,
		// don't do anything prevent for create invalid conversation.
		if feedChange.IsRemove() {
			return nil
		}

		post, err := wh.fbClient.CallAPIGetPost(postID, accessToken)
		if err != nil {
			return err
		}

		if err := wh.createParentAndChildPosts(extPageID, createdTime, ctx, post); err != nil {
			return err
		}
	}

	commentID := feedChange.Value.CommentID
	externalCmt, err := wh.getExternalComment(ctx, commentID)
	if err != nil {
		return err
	}

	if externalCmt == nil {
		if feedChange.IsRemove() {
			return nil
		}

		var createCmtCmd []*fbmessaging.CreateFbExternalCommentArgs
		comment, err := wh.fbClient.CallAPICommentByID(accessToken, commentID)
		if err != nil {
			return err
		}
		createCmtCmd = append(createCmtCmd, convertModelCommentToCreateCommentArgs(extPageID, postID, createdTime, comment))

		if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
			FbExternalComments: createCmtCmd,
		}); err != nil {
			return err
		}
		return nil
	}

	if feedChange.IsEdited() {
		updateCommentMsgCmd := &fbmessaging.UpdateFbCommentMessageCommand{
			ExternalCommentID: commentID,
			Message:           feedChange.Value.Message,
		}
		if err := wh.fbmessagingAggr.Dispatch(ctx, updateCommentMsgCmd); err != nil {
			return err
		}
	}
	return nil
}

func (wh *Webhook) getExternalPost(ctx context.Context, extPostID string) (*fbmessaging.FbExternalPost, error) {
	getFbExternalPostQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: extPostID,
	}
	if err := wh.fbmessagingQuery.Dispatch(ctx, getFbExternalPostQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	return getFbExternalPostQuery.Result, nil
}

func (wh *Webhook) getExternalComment(ctx context.Context, commentID string) (*fbmessaging.FbExternalComment, error) {
	getCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: commentID,
	}
	if err := wh.fbmessagingQuery.Dispatch(ctx, getCommentQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	return getCommentQuery.Result, nil
}

func (wh *Webhook) createParentAndChildPosts(externalPageID string, createdTime time.Time, ctx context.Context, post *model.Post) error {
	parentPost := convertModelPostToCreatePostArgs(externalPageID, createdTime, post)
	saveExternalPostCmd := &fbmessaging.SaveFbExternalPostCommand{
		ExternalPageID:      parentPost.ExternalPageID,
		ExternalID:          parentPost.ExternalID,
		ExternalPicture:     parentPost.ExternalPicture,
		ExternalIcon:        parentPost.ExternalIcon,
		ExternalMessage:     parentPost.ExternalMessage,
		ExternalCreatedTime: parentPost.ExternalCreatedTime,
		ExternalAttachments: parentPost.ExternalAttachments,
		ExternalFrom:        parentPost.ExternalFrom,
		ExternalParent:      nil,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, saveExternalPostCmd); err != nil {
		return err
	}

	createChildPostCmd := &fbmessaging.CreateFbExternalPostsCommand{
		FbExternalPosts: buildAllChildPost(parentPost),
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, createChildPostCmd); err != nil {
		return err
	}
	return nil
}
