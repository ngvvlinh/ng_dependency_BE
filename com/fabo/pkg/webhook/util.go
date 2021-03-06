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
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/webhook_type"
	fblog "o.o/backend/com/fabo/main/fblog/model"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/mq"
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
					subAttachment := &fbmessaging.SubAttachment{
						Type:        v.Type,
						URL:         v.URL,
						Description: v.Description,
					}
					if v.Target != nil {
						subAttachment.Target = &fbmessaging.TargetDataSubAttachment{
							ID:  v.Target.ID,
							URL: v.Target.URL,
						}
					}

					if v.Media != nil && v.Media.Image != nil {
						subAttachment.Media = &fbmessaging.MediaDataSubAttachment{
							Height: v.Media.Image.Height,
							Width:  v.Media.Image.Width,
							Src:    v.Media.Image.Src,
						}
					}
					pAtt.SubAttachments = append(pAtt.SubAttachments, subAttachment)
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
		StatusType:          fb_status_type.ParseFbStatusTypeWithDefault(post.StatusType, fb_status_type.Unknown),
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
			var externalAttachments []*fbmessaging.PostAttachment
			externalAttachments = append(externalAttachments, &fbmessaging.PostAttachment{
				MediaType: subAtt.Type,
				Media: &fbmessaging.MediaPostAttachment{
					Image: &fbmessaging.ImageMediaPostAttachment{
						Height: subAtt.Media.Height,
						Width:  subAtt.Media.Width,
						Src:    subAtt.Media.Src,
					},
				},
				Type: subAtt.Type,
			})
			childPost := &fbmessaging.CreateFbExternalPostArgs{
				ID:                  cm.NewID(),
				ExternalPageID:      post.ExternalPageID,
				ExternalID:          childPostID,
				ExternalParentID:    post.ExternalID,
				ExternalFrom:        post.ExternalFrom,
				ExternalPicture:     subAtt.Media.Src,
				ExternalIcon:        post.ExternalIcon,
				ExternalMessage:     subAtt.Description,
				ExternalAttachments: externalAttachments,
				ExternalCreatedTime: post.ExternalCreatedTime,
				ExternalUpdatedTime: post.ExternalUpdatedTime,
				FeedType:            fb_feed_type.Post,
				StatusType:          post.StatusType,
				Type:                fb_post_type.Page,
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
		InternalSource:       fb_internal_source.Facebook, // comment comes from webhook, by default create it with `Facebook`
		PostType:             fb_post_type.Page,
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

type CallbackOptions struct {
	UserID string // when object == user
	PageID string // when object == page
	Object string // user or page
}

func convertCallbackOptions(options map[string]string) (res CallbackOptions) {
	if userID, ok := options["user_id"]; ok {
		res.UserID = userID
	}
	if pageID, ok := options["page_id"]; ok {
		res.PageID = pageID
	}
	if object, ok := options["object"]; ok {
		res.Object = object
	}
	return
}

func (wh *Webhook) forwardWebhook(c *httpx.Context, message WebhookMessages) {
	callbacks, err := wh.webhookCallbackService.GetWebhookCallbacks(webhook_type.Fabo.String())
	if err != nil {
		return
	}

	client := http.Client{
		Timeout: 60 * time.Second,
	}
	for _, callback := range callbacks {
		callbackURL, header := callback.URL, c.Req.Header // closure
		callbackOptions := convertCallbackOptions(callback.Options)

		// filter message by options
		if callbackOptions.Object != "" && message.Object != callbackOptions.Object {
			continue
		}
		if callbackOptions.Object == "user" && callbackOptions.UserID != "" &&
			len(message.Entry) != 0 && message.Entry[0].ID != callbackOptions.UserID {
			continue
		}
		if callbackOptions.Object == "page" && callbackOptions.PageID != "" &&
			len(message.Entry) != 0 && message.Entry[0].ID != callbackOptions.PageID {
			continue
		}

		go func() (_err error) { // ignore the error
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
func (wh *WebhookHandler) getPageAccessToken(ctx context.Context, extPageID string) (string, mq.Code, error) {
	getAccessTokenQuery := &fbpaging.GetPageAccessTokenQuery{
		ExternalID: extPageID,
	}

	if err := wh.fbPageQuery.Dispatch(ctx, getAccessTokenQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return "", mq.CodeIgnore, cm.Errorf(cm.FacebookWebhookIgnored, nil, "external page with id %v not found", extPageID)
		}
		return "", mq.CodeStop, nil
	}
	return getAccessTokenQuery.Result, mq.CodeOK, nil
}

func (wh *Webhook) saveLogsWebhookPage(msg WebhookMessages, err error) {
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
		ID:             cm.NewID(),
		ExternalPageID: pageID,
		ExternalID:     externalID,
		Data:           nil,
		Error:          nil,
		Type:           string(msg.MessageType()),
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

func (wh *Webhook) saveLogsWebhookUser(msg WebhookUser, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	logData := &fblog.FbWebhookLog{
		ID:             cm.NewID(),
		ExternalUserID: msg.ExternalUserID(),
		ExternalID:     msg.ExternalID(),
		Data:           nil,
		Error:          nil,
		Type:           string(msg.Type()),
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

func (wh *WebhookHandler) getProfile(accessToken, externalPageID, PSID string, profileDefault *fbclientmodel.Profile) (*fbclientmodel.Profile, error) {
	profile, err := wh.faboRedis.LoadProfilePSID(externalPageID, PSID)
	switch err {
	// If profile not in redis then call api getProfileByPSID
	case redis.ErrNil:
		_profile, _err := wh.fbClient.CallAPIGetProfileByPSID(&fbclient.GetProfileRequest{
			AccessToken:    accessToken,
			PSID:           PSID,
			PageID:         externalPageID,
			ProfileDefault: profileDefault,
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

func (wh *WebhookHandler) IsTestPage(ctx context.Context, externalPageID string) (bool, error) {
	getExternalPageQuery := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := wh.fbPageQuery.Dispatch(ctx, getExternalPageQuery); err != nil {
		return false, err
	}

	return strings.HasPrefix(getExternalPageQuery.Result.ExternalName, fbclient.PrefixFanPageNameTest), nil
}
