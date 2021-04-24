package webhook

import (
	"fmt"
	"strings"
	"time"

	"github.com/r3labs/sse"
	"golang.org/x/net/context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/convert"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
)

func (wh *WebhookHandler) HandleUserLiveVideo(
	ctx context.Context, webhookUser WebhookUser,
) (mq.Code, error) {
	if len(webhookUser.Entry) == 0 {
		return mq.CodeIgnore, nil
	}

	for _, entry := range webhookUser.Entry {
		extUserID := entry.ID
		for _, change := range entry.Changes {
			if change.Field != WebhookUserLiveVideos {
				continue
			}
			if change.Value == nil {
				continue
			}

			extLiveVideoID := change.Value.ID
			extStatus := change.Value.Status

			jobID := fmt.Sprintf("%s-%s", extUserID, extLiveVideoID)

			getFbUserQuery := &fbusering.GetFbExternalUserInternalByExternalIDQuery{
				ExternalID: extUserID,
			}
			if err := wh.fbUserQuery.Dispatch(ctx, getFbUserQuery); err != nil {
				return mq.CodeIgnore, err
			}
			fbUserToken := getFbUserQuery.Result.Token
			extPostID := fmt.Sprintf("%s-%s", extUserID, extLiveVideoID)

			getLiveVideoReq := &fbclient.GetLiveVideoRequest{
				AccessToken: fbUserToken,
				LiveVideoID: extLiveVideoID,
				UserID:      extUserID,
			}
			liveVideoResp, err := wh.fbClient.CallAPIGetLiveVideo(getLiveVideoReq)
			if err != nil {
				return mq.CodeIgnore, err
			}

			getFbExtPostQuery := &fbmessaging.GetFbExternalPostByExternalIDAndExternalUserIDQuery{
				ExternalID:     extPostID,
				ExternalUserID: extUserID,
			}
			_err := wh.fbmessagingQuery.Dispatch(ctx, getFbExtPostQuery)
			switch cm.ErrorCode(_err) {
			case cm.NoError:
			// no-op
			case cm.NotFound:
				saveFbExtPostCmd := &fbmessaging.CreateFbExternalPostsCommand{
					FbExternalPosts: []*fbmessaging.CreateFbExternalPostArgs{
						{
							ID:                  cm.NewID(),
							ExternalUserID:      extUserID,
							ExternalID:          extPostID,
							ExternalFrom:        convert.ConvertObjectFrom(liveVideoResp.From),
							ExternalMessage:     cm.Coalesce(liveVideoResp.Title, liveVideoResp.Description),
							ExternalPicture:     liveVideoResp.Video.Picture,
							ExternalCreatedTime: liveVideoResp.CreationTime.ToTime(),
							FeedType:            fb_feed_type.Post,
							StatusType:          fb_status_type.AddedVideo,
							Type:                fb_post_type.User,
						},
					},
				}
				if err := wh.fbmessagingAggr.Dispatch(ctx, saveFbExtPostCmd); err != nil {
					return mq.CodeIgnore, err
				}
			default:
				return mq.CodeIgnore, err
			}

			extStatusStr := strings.ToUpper(string(extStatus))
			updateLiveVideoStatusCmd := &fbmessaging.UpdateLiveVideoStatusFromSyncCommand{
				ExternalID:              extPostID,
				ExternalLiveVideoStatus: extStatusStr,
				LiveVideoStatus:         fb_live_video_status.ConvertToFbLiveVideoStatus(extStatusStr),
			}
			if err := wh.fbmessagingAggr.Dispatch(ctx, updateLiveVideoStatusCmd); err != nil {
				return mq.CodeIgnore, err
			}

			switch extStatus {
			case LiveStatus:
				args := &LiveVideoArguments{
					ExtUserID:      extUserID,
					ExtPostID:      extPostID,
					ExtLiveVideoID: extLiveVideoID,
					Token:          fbUserToken,
				}
				wh.jobKeeper.AddJob(jobID, wh.handleLiveVideo, args)
			case LiveStoppedStatus:
				wh.jobKeeper.StopJob(jobID)
			default:
				// no-op
			}
		}
	}

	return mq.CodeOK, nil
}

type LiveVideoArguments struct {
	ExtUserID      string
	ExtPostID      string
	ExtLiveVideoID string
	Token          string
}

func (wh *WebhookHandler) handleLiveVideo(ctx context.Context, _args interface{}) error {
	liveVideoArguments := _args.(*LiveVideoArguments)
	extLiveVideoID := liveVideoArguments.ExtLiveVideoID
	extPostID := liveVideoArguments.ExtPostID
	extUserID := liveVideoArguments.ExtUserID
	token := liveVideoArguments.Token

	url := fmt.Sprintf("https://streaming-graph.facebook.com/%s/live_comments?access_token=%s&comment_rate=ten_per_second&fields=from{name,id,email,first_name,last_name,picture},message,attachment", extLiveVideoID, token)
	client := sse.NewClient(url)

	go func(_ctx context.Context, _extPostID, _extUserID string) {
		client.SubscribeRawWithContext(_ctx, func(msg *sse.Event) {
			// Got some data!
			fmt.Println(string(msg.Data))
			var comment LiveVideoComment
			if err := jsonx.Unmarshal(msg.Data, &comment); err != nil {
				return
			}

			fbExternalComment := &fbmessaging.CreateFbExternalCommentArgs{
				ID:                  cm.NewID(),
				ExternalPostID:      _extPostID,
				ExternalID:          comment.ID,
				ExternalMessage:     comment.Message,
				ExternalCreatedTime: time.Time{},
				InternalSource:      fb_internal_source.Facebook,
				ExternalOwnerPostID: _extUserID,
				PostType:            fb_post_type.User,
			}
			if comment.From != nil {
				fbExternalComment.ExternalUserID = comment.From.ID
				fbExternalComment.ExternalFrom = convert.ConvertObjectFrom(comment.From)
			} else {
				extUserIDFrom := cm.NewID().String()
				from := &fbmessaging.FbObjectFrom{
					ID: extUserIDFrom,
				}
				fbExternalComment.ExternalUserID = extUserIDFrom
				fbExternalComment.ExternalFrom = from
			}

			createFbExtCommentCmd := &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
				FbExternalComments: []*fbmessaging.CreateFbExternalCommentArgs{fbExternalComment},
			}

			if err := wh.fbmessagingAggr.Dispatch(bus.Ctx(), createFbExtCommentCmd); err != nil {
				return
			}
		})
	}(ctx, extPostID, extUserID)

	return nil
}
