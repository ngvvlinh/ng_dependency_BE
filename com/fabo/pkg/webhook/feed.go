package webhook

import (
	"context"
	"fmt"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/mq"
)

// Facebook feed is any action on page (create or update a post, make comment,
// any reaction ....)
func (wh *WebhookHandler) HandleFeed(
	ctx context.Context, feed WebhookMessages,
) (mq.Code, error) {
	// Ignore Create or Update action from page owner.
	// But keep other actions like Remove, Delete ...
	if feed.IsCreateOrEditCommentFromPageOwner() {
		return mq.CodeIgnore, nil
	}

	for _, entry := range feed.Entry {
		externalPageID := entry.ID
		createdTime := time.Unix(int64(entry.Time), 0)

		isTestPage, _err := wh.IsTestPage(ctx, externalPageID)
		if _err != nil {
			if cm.ErrorCode(_err) == cm.NotFound {
				return mq.CodeIgnore, nil
			}
			return mq.CodeRetry, _err
		}
		// ignore test page
		if cmenv.IsProd() && isTestPage {
			return mq.CodeIgnore, nil
		}

		accessToken, returnCode, err := wh.getPageAccessToken(ctx, externalPageID)
		if returnCode != mq.CodeOK {
			return returnCode, err
		}

		for _, change := range entry.Changes {
			switch {
			case change.IsEvent():
				return wh.handleFeedEvent(ctx, externalPageID, change, createdTime)
			case change.IsAdminPost(externalPageID):
				return wh.handleFeedPost(ctx, change, createdTime, externalPageID, accessToken)
			case change.IsComment():
				return wh.handleFeedComment(ctx, change, createdTime, externalPageID, accessToken)
			}
		}
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) handleFeedEvent(
	ctx context.Context, extPageID string,
	feedChange FeedChange, createdTime time.Time,
) (mq.Code, error) {
	saveEvent := &fbmessaging.SaveFbExternalPostCommand{
		ExternalPageID:      extPageID,
		ExternalID:          feedChange.Value.PostID,
		ExternalCreatedTime: createdTime,
		FeedType:            fb_feed_type.Event,
		StatusType:          fb_status_type.CreatedEvent,
		Type:                fb_post_type.User,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, saveEvent); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) handleFeedPost(
	ctx context.Context, feedChange FeedChange,
	createdTime time.Time, extPageID, accessToken string,
) (mq.Code, error) {
	postID := feedChange.Value.PostID
	fromID := feedChange.Value.From.ID
	if err := wh.lockFeedPost(extPageID, postID, fromID); err != nil {
		return mq.CodeRetry, err
	}

	if feedChange.IsRemove() {
		return wh.handleRemovePost(ctx, extPageID, postID)
	}

	externalPost, err := wh.getExternalPost(ctx, postID)
	if err != nil {
		return mq.CodeRetry, err
	}

	post, err := wh.fbClient.CallAPIGetPost(&fbclient.GetPostRequest{
		AccessToken: accessToken,
		PostID:      postID,
		PageID:      extPageID,
	})
	if err != nil {
		return mq.CodeIgnore, err
	}

	// if post does not exist in db, create it
	if externalPost == nil {
		return wh.createParentAndChildPosts(ctx, extPageID, createdTime, post)
	}

	if err := wh.updateParentAndChildPost(ctx, extPageID, post); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) lockFeedPost(pageID, postID, fromID string) error {
	// Sometimes facebook may send many requests for one feed. E.g: If a
	// user posts a post containing 4 images, facebook may calls more than 5
	// requests at the same time. In this case, we can build full post and
	// all child posts from one request, so just hold one and ignore the
	// rest.
	actionKey := fmt.Sprintf("FEED_%v_%v_%v", pageID, postID, fromID)
	if wh.faboRedis.IsExist(actionKey) {
		return nil
	}
	if err := wh.faboRedis.SetWithTTL(actionKey, true, 2); err != nil {
		return err
	}
	return nil
}

func (wh *WebhookHandler) handleRemovePost(ctx context.Context, pageID, postID string) (mq.Code, error) {
	removeCmd := &fbmessaging.RemovePostCommand{
		ExternalPostID: postID,
		ExternalPageID: pageID,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, removeCmd); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) updateFeedPost(ctx context.Context, postID string, message string, externalPicture string) error {
	cmdUpdate := &fbmessaging.UpdateFbPostMessageAndPictureCommand{
		ExternalPostID:  postID,
		Message:         message,
		ExternalPicture: externalPicture,
	}
	return wh.fbmessagingAggr.Dispatch(ctx, cmdUpdate)
}

func (wh *WebhookHandler) updateParentAndChildPost(ctx context.Context, extPageID string, extPost *model.Post) error {
	createdTime := time.Unix(int64(extPost.CreatedTime), 0)
	parentPost := convertModelPostToCreatePostArgs(extPageID, createdTime, extPost)
	allPosts := []*fbmessaging.CreateFbExternalPostArgs{parentPost}

	// If all attachments is not from other page, build all child posts.
	// Behaviour of share po
	if extPost.IsResourceFromCurrentPage() {
		childPosts := buildAllChildPost(parentPost)
		allPosts = append(childPosts, parentPost)
	}

	for _, post := range allPosts {
		err := wh.updateFeedPost(ctx, post.ExternalID, post.ExternalMessage, post.ExternalPicture)
		if err != nil {
			if cm.ErrorCode(err) == cm.NotFound {
				createPostCmd := &fbmessaging.SaveFbExternalPostCommand{
					ExternalPageID:      post.ExternalPageID,
					ExternalID:          post.ExternalID,
					ExternalFrom:        post.ExternalFrom,
					ExternalPicture:     post.ExternalPicture,
					ExternalIcon:        post.ExternalIcon,
					ExternalMessage:     post.ExternalMessage,
					ExternalAttachments: post.ExternalAttachments,
					ExternalCreatedTime: post.ExternalCreatedTime,
					ExternalParentID:    parentPost.ExternalID,
					StatusType:          post.StatusType,
					Type:                fb_post_type.Page,
				}
				if err := wh.fbmessagingAggr.Dispatch(ctx, createPostCmd); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func (wh *WebhookHandler) handleFeedComment(
	ctx context.Context, feedChange FeedChange,
	createdTime time.Time, extPageID, accessToken string,
) (mq.Code, error) {
	if feedChange.IsEventComment() {
		return mq.CodeOK, nil
	}

	postID := feedChange.Value.PostID
	commentID := feedChange.Value.CommentID
	if feedChange.IsRemove() {
		return wh.handleRemoveComment(ctx, commentID)
	}

	externalPost, err := wh.getExternalPost(ctx, postID)
	if err != nil {
		return mq.CodeRetry, err
	}

	if externalPost == nil {
		post, err := wh.fbClient.CallAPIGetPost(&fbclient.GetPostRequest{
			AccessToken: accessToken,
			PostID:      postID,
			PageID:      extPageID,
		})
		if err != nil {
			return mq.CodeIgnore, err
		}

		if code, err := wh.createParentAndChildPosts(ctx, extPageID, createdTime, post); err != nil {
			return code, err
		}
	}

	externalCmt, err := wh.getExternalComment(ctx, commentID)
	if err != nil {
		return mq.CodeRetry, err
	}

	if externalCmt == nil {
		comment, err := wh.fbClient.CallAPICommentByID(&fbclient.GetCommentByIDRequest{
			AccessToken: accessToken,
			CommentID:   commentID,
			PageID:      extPageID,
		})
		if err != nil {
			return mq.CodeRetry, err
		}

		createCmtCmd := convertModelCommentToCreateCommentArgs(extPageID, postID, createdTime, comment)
		if feedChange.IsPageLikeComment(extPageID) {
			createCmtCmd.IsLiked = true
		}
		if feedChange.IsPageHideComment(extPageID) {
			createCmtCmd.IsHidden = true
		}

		if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
			FbExternalComments: []*fbmessaging.CreateFbExternalCommentArgs{createCmtCmd},
		}); err != nil {
			return mq.CodeRetry, err
		}
		return mq.CodeOK, nil
	}

	// handle like and unlike
	if (feedChange.IsPageUnLikeComment(extPageID) && externalCmt.IsLiked) ||
		(feedChange.IsPageLikeComment(extPageID) && !externalCmt.IsLiked) {
		if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.LikeOrUnLikeCommentCommand{
			ExternalCommentID: externalCmt.ExternalID,
			IsLiked:           !externalCmt.IsLiked,
		}); err != nil {
			return mq.CodeRetry, err
		}
	}

	// handle hide and unhide
	if (feedChange.IsPageUnHideComment(extPageID) && externalCmt.IsHidden) ||
		(feedChange.IsPageHideComment(extPageID) && !externalCmt.IsHidden) {
		if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.HideOrUnHideCommentCommand{
			ExternalCommentID: externalCmt.ExternalID,
			IsHidden:          !externalCmt.IsHidden,
		}); err != nil {
			return mq.CodeRetry, err
		}
	}

	if feedChange.IsEdited() {
		updateCommentMsgCmd := &fbmessaging.UpdateFbCommentMessageCommand{
			ExternalCommentID: commentID,
			Message:           feedChange.Value.Message,
		}
		if err := wh.fbmessagingAggr.Dispatch(ctx, updateCommentMsgCmd); err != nil {
			return mq.CodeRetry, err
		}
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) handleRemoveComment(ctx context.Context, commentID string) (mq.Code, error) {
	removeCommentArgs := &fbmessaging.RemoveCommentCommand{
		ExternalCommentID: commentID,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, removeCommentArgs); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func (wh *WebhookHandler) getExternalPost(ctx context.Context, extPostID string) (*fbmessaging.FbExternalPost, error) {
	getFbExternalPostQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: extPostID,
	}
	if err := wh.fbmessagingQuery.Dispatch(ctx, getFbExternalPostQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	return getFbExternalPostQuery.Result, nil
}

func (wh *WebhookHandler) getExternalComment(ctx context.Context, commentID string) (*fbmessaging.FbExternalComment, error) {
	getCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: commentID,
	}
	if err := wh.fbmessagingQuery.Dispatch(ctx, getCommentQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	return getCommentQuery.Result, nil
}

func (wh *WebhookHandler) createParentAndChildPosts(
	ctx context.Context, externalPageID string,
	createdTime time.Time, post *model.Post,
) (mq.Code, error) {
	parentPost := convertModelPostToCreatePostArgs(externalPageID, createdTime, post)
	createParentCmd := &fbmessaging.SaveFbExternalPostCommand{
		ExternalPageID:      parentPost.ExternalPageID,
		ExternalID:          parentPost.ExternalID,
		ExternalPicture:     parentPost.ExternalPicture,
		ExternalIcon:        parentPost.ExternalIcon,
		ExternalMessage:     parentPost.ExternalMessage,
		ExternalCreatedTime: parentPost.ExternalCreatedTime,
		ExternalAttachments: parentPost.ExternalAttachments,
		ExternalFrom:        parentPost.ExternalFrom,
		ExternalParentID:    "",
		FeedType:            fb_feed_type.Post,
		Type:                fb_post_type.Page,
		StatusType:          parentPost.StatusType,
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, createParentCmd); err != nil {
		return mq.CodeRetry, err
	}

	// If all attachments is not from other build all child posts.
	if post.IsResourceFromCurrentPage() {
		createChildPostCmd := &fbmessaging.CreateFbExternalPostsCommand{
			FbExternalPosts: buildAllChildPost(parentPost),
		}
		if err := wh.fbmessagingAggr.Dispatch(ctx, createChildPostCmd); err != nil {
			return mq.CodeRetry, err
		}
	}
	return mq.CodeOK, nil
}
