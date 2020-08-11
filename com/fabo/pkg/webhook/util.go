package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/webhook_type"
	fblog "o.o/backend/com/fabo/main/fblog/model"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/redis"
	emodel "o.o/backend/pkg/etop/model"
	"o.o/common/jsonx"
)

func convertModelPostToCreatePostArgs(pageID string, externalCreatedTime time.Time, post *model.Post) *fbmessaging.CreateFbExternalPostArgs {
	var extAttachments []*fbmessaging.PostAttachment
	if post.Attachments != nil {
		for _, att := range post.Attachments.Data {
			pAtt := &fbmessaging.PostAttachment{
				MediaType: att.MediaType,
				Type:      att.Type,
			}

			if att.Media != nil {
				pAtt.Media = &fbmessaging.MediaPostAttachment{
					Image: &fbmessaging.ImageMediaPostAttachment{
						Height: att.Media.Image.Height,
						Width:  att.Media.Image.Width,
						Src:    att.Media.Image.Src,
					},
				}
			}

			if att.SubAttachments != nil {
				pAtt.SubAttachments = []*fbmessaging.SubAttachment{}
				for _, v := range att.SubAttachments.Data {
					pAtt.SubAttachments = append(pAtt.SubAttachments, &fbmessaging.SubAttachment{
						Media: &fbmessaging.MediaDataSubAttachment{
							Height: v.Media.Image.Height,
							Width:  v.Media.Image.Width,
							Src:    v.Media.Image.Src,
						},
						Target: &fbmessaging.TargetDataSubAttachment{
							ID:  v.Target.ID,
							URL: v.Target.URL,
						},
						Type:        v.Type,
						URL:         v.URL,
						Description: v.Description,
					})
				}
			}
			extAttachments = append(extAttachments, pAtt)
		}

	}
	res := &fbmessaging.CreateFbExternalPostArgs{
		ExternalPageID:      pageID,
		ExternalID:          post.ID,
		ExternalParentID:    "",
		ExternalFrom:        nil,
		ExternalPicture:     post.Picture,
		ExternalIcon:        post.Icon,
		ExternalMessage:     post.Message,
		ExternalAttachments: extAttachments,
		ExternalCreatedTime: externalCreatedTime,
		ExternalUpdatedTime: time.Time{},
	}
	if post.From != nil {
		res.ExternalFrom = &fbmessaging.FbObjectFrom{
			ID:        post.From.ID,
			Name:      post.From.Name,
			Email:     post.From.Email,
			FirstName: post.From.FirstName,
			LastName:  post.From.LastName,
		}
	}

	return res
}

func buildAllChildPost(post *fbmessaging.CreateFbExternalPostArgs) []*fbmessaging.CreateFbExternalPostArgs {
	if post.ExternalAttachments != nil && len(post.ExternalAttachments) > 0 {
		subAttachments := post.ExternalAttachments[0].SubAttachments
		var res []*fbmessaging.CreateFbExternalPostArgs
		for _, subAtt := range subAttachments {
			childPostID := fmt.Sprintf("%v_%v", post.ExternalPageID, subAtt.Target.ID)
			childPost := &fbmessaging.CreateFbExternalPostArgs{
				ID:                  cm.NewID(),
				ExternalPageID:      post.ExternalPageID,
				ExternalID:          childPostID,
				ExternalParentID:    post.ExternalID,
				ExternalFrom:        post.ExternalFrom,
				ExternalPicture:     subAtt.Media.Src,
				ExternalIcon:        post.ExternalIcon,
				ExternalMessage:     subAtt.Description,
				ExternalCreatedTime: post.ExternalCreatedTime,
				ExternalUpdatedTime: post.ExternalUpdatedTime,
				FeedType:            fb_feed_type.Post,
			}
			res = append(res, childPost)
		}

		return res
	}
	return []*fbmessaging.CreateFbExternalPostArgs{}
}

func convertModelCommentToCreateCommentArgs(pageId string, postID string, createdTime time.Time, comment *model.Comment) *fbmessaging.CreateFbExternalCommentArgs {
	res := &fbmessaging.CreateFbExternalCommentArgs{
		ID:                   cm.NewID(),
		ExternalPostID:       postID,
		ExternalPageID:       pageId,
		ExternalID:           comment.ID,
		ExternalMessage:      comment.Message,
		ExternalCommentCount: comment.CommentCount,
		ExternalCreatedTime:  createdTime,
	}

	if comment.From != nil {
		res.ExternalFrom = &fbmessaging.FbObjectFrom{
			ID:        comment.From.ID,
			Name:      comment.From.Name,
			Email:     comment.From.Email,
			FirstName: comment.From.FirstName,
			LastName:  comment.From.LastName,
		}
		if comment.From.Picture != nil {
			res.ExternalFrom.ImageURL = comment.From.Picture.Data.Url
		}

		res.ExternalUserID = comment.From.ID
	}

	if comment.Attachment != nil {
		res.ExternalAttachment = &fbmessaging.CommentAttachment{
			Title: comment.Attachment.Title,
			Type:  comment.Attachment.Type,
			URL:   comment.Attachment.URL,
		}
		if comment.Attachment.Media != nil && comment.Attachment.Media.Image != nil {
			res.ExternalAttachment.Media = &fbmessaging.ImageMediaDataSubAttachment{
				Image: &fbmessaging.MediaDataSubAttachment{
					Height: comment.Attachment.Media.Image.Height,
					Width:  comment.Attachment.Media.Image.Width,
					Src:    comment.Attachment.Media.Image.Src,
				},
			}
		}

		if res.ExternalAttachment.Target != nil {
			res.ExternalAttachment.Target = &fbmessaging.TargetDataSubAttachment{
				ID:  comment.Attachment.Target.ID,
				URL: comment.Attachment.Target.URL,
			}
		}
	}

	if comment.Parent != nil {
		res.ExternalParent = &fbmessaging.FbObjectParent{
			CreatedTime: comment.Parent.CreatedTime.ToTime(),
			Message:     comment.Parent.Message,
			ID:          comment.Parent.ID,
		}
		if comment.Parent.From != nil {
			res.ExternalParent.From = &fbmessaging.FbObjectFrom{
				ID:        comment.Parent.From.ID,
				Name:      comment.Parent.From.Name,
				Email:     comment.Parent.From.Email,
				FirstName: comment.Parent.From.FirstName,
				LastName:  comment.Parent.From.LastName,
			}
		}
		res.ExternalParentID = comment.Parent.ID
		res.ExternalParentUserID = comment.Parent.From.ID
	}
	return res
}

func (wh *Webhook) getProfileByPSID(accessToken, pageID, PSID string) (profile *fbclientmodel.Profile, err error) {
	profile, err = wh.faboRedis.LoadProfilePSID(pageID, PSID)
	switch err {
	case redis.ErrNil:
		profile, err = wh.fbClient.CallAPIGetProfileByPSID(&fbclient.GetProfileRequest{
			AccessToken: accessToken,
			PSID:        PSID,
			PageID:      pageID,
		})
		if err != nil {
			return nil, err
		}
		if err := wh.faboRedis.SaveProfilePSID(pageID, PSID, profile); err != nil {
			return nil, err
		}
	case nil:
	// no-op
	default:
		return nil, err
	}
	return profile, nil
}

func (wh *Webhook) forwardWebhook(c *httpx.Context, message WebhookMessages) {
	callbackURLs, err := wh.webhookCallbackService.GetWebhookCallbackURLs(webhook_type.Fabo.String())
	if err != nil {
		return
	}

	client := http.Client{
		Timeout: 60 * time.Second,
	}
	for _, callbackURL := range callbackURLs {
		callbackURL, header := callbackURL, c.Req.Header // closure
		go func() (_err error) {                         // ignore the error
			defer func() {
				wh.mu.RLock()
				latestTimeError, ok := wh.mapCallbackURLAndLatestTimeError[callbackURL]
				wh.mu.RUnlock()
				if _err == nil {
					if ok { // reset the timer when callback is successful
						wh.mu.Lock()
						delete(wh.mapCallbackURLAndLatestTimeError, callbackURL)
						wh.mu.Unlock()
					}
					return
				}

				ll.SendMessagef("[Error] callback_url: %v\n\n%v", callbackURL, _err.Error())
				switch {
				case !ok:
					// first time error, store the time
					wh.mu.Lock()
					wh.mapCallbackURLAndLatestTimeError[callbackURL] = time.Now()
					wh.mu.Unlock()

				case time.Now().Sub(latestTimeError) < oneHour:
					// under one hour, do nothing

				default:
					// error for too long, remove the webhook
					wh.mu.Lock()
					delete(wh.mapCallbackURLAndLatestTimeError, callbackURL)
					_ = wh.webhookCallbackService.RemoveWebhookCallbackURL(webhook_type.Fabo.String(), callbackURL)
					wh.mu.Unlock()
				}
			}()

			data, err := jsonx.Marshal(message)
			if err != nil {
				return err
			}

			req, err2 := http.NewRequest("POST", callbackURL, bytes.NewReader(data))
			if err2 != nil {
				return err2
			}
			req.Header = header
			resp, err2 := client.Do(req)
			if err2 != nil {
				return err2
			}
			if resp.StatusCode != 200 {
				return fmt.Errorf("response status: %v", resp.StatusCode)
			}
			return nil
		}()
	}
}

// TODO(vu): cache, write a query for retrieving only the access token
func (wh *Webhook) getPageAccessToken(ctx context.Context, extPageID string) (string, error) {
	getAccessTokenQuery := &fbpaging.GetPageAccessTokenQuery{
		ExternalID: extPageID,
	}

	if err := wh.fbPageQuery.Dispatch(ctx, getAccessTokenQuery); err != nil {
		return "", cm.MapError(err).
			Mapf(cm.NotFound, cm.FacebookWebhookIgnored, "external page with id %v not found", extPageID).
			Throw()
	}
	return getAccessTokenQuery.Result, nil
}

func (wh *Webhook) saveLogsWebhook(msg WebhookMessages, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	pageID := ""
	externalID := ""

	if len(msg.Entry) > 0 {
		entry := msg.Entry[0]
		pageID = entry.ID
		externalID = entry.ExternalID()
	}

	logData := &fblog.FbWebhookLog{
		ID:         cm.NewID(),
		PageID:     pageID,
		ExternalID: externalID,
		Data:       nil,
		Error:      nil,
		Type:       string(msg.MessageType()),
	}

	if err != nil {
		logData.Error = emodel.ToError(err)
	}

	if err := enc.Encode(msg); err == nil {
		logData.Data = buf.Bytes()
	}

	if _, _err := wh.dbLog.Insert(logData); _err != nil {
		ll.SendMessagef("Insert db etop_log error %v", _err.Error())
	}
}

func (wh *Webhook) getProfile(accessToken, externalPageID, PSID string) (*fbclientmodel.Profile, error) {
	profile, err := wh.faboRedis.LoadProfilePSID(externalPageID, PSID)
	switch err {
	// If profile not in redis then call api getProfileByPSID
	case redis.ErrNil:
		_profile, _err := wh.fbClient.CallAPIGetProfileByPSID(&fbclient.GetProfileRequest{
			AccessToken: accessToken,
			PSID:        PSID,
			PageID:      externalPageID,
		})
		if _err != nil {
			return nil, _err
		}
		if _err := wh.faboRedis.SaveProfilePSID(externalPageID, PSID, _profile); _err != nil {
			return nil, _err
		}
		return _profile, nil
	case nil:
		return profile, nil
	default:
		return nil, err
	}
}

func (wh *Webhook) IsTestPage(ctx context.Context, externalPageID string) (bool, error) {
	getExternalPageQuery := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := wh.fbPageQuery.Dispatch(ctx, getExternalPageQuery); err != nil {
		return false, err
	}

	return strings.HasPrefix(getExternalPageQuery.Result.ExternalName, fbclient.PrefixFanPageNameTest), nil
}
