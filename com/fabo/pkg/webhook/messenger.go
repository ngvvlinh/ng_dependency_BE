package webhook

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/common/xerrors"
)

func (wh *Webhook) handleMessenger(ctx context.Context, webhookMessages WebhookMessages) error {
	entries := webhookMessages.Entry
	for _, entry := range entries {
		externalPageID := entry.ID
		for _, messge := range entry.Messaging {
			if messge.Message != nil {
				PSID := messge.Sender.ID
				if PSID == externalPageID {
					PSID = messge.Recipient.ID
				}
				// message ID
				mid := messge.Message.Mid
				err := wh.handleMessageReturned(ctx, externalPageID, PSID, mid)
				if err == nil {
					continue
				}
				facebookError := err.(*xerrors.APIError)
				code := facebookError.Meta["code"]
				if code == fbclient.AccessTokenHasExpired.String() {
					continue
				} else {
					ll.SendMessage(err.Error())
					return err
				}
			}
		}
	}
	return nil
}

// TODO(ngoc): refactor
func (wh *Webhook) handleMessageReturned(ctx context.Context, externalPageID, PSID, mid string) error {
	isTestPage, _err := wh.IsTestPage(ctx, externalPageID)
	if _err != nil {
		return _err
	}
	// ignore test page
	if cmenv.IsProd() && isTestPage {
		return nil
	}

	accessToken, err := wh.getPageAccessToken(ctx, externalPageID)
	if err != nil {
		return err
	}

	// Get message
	messageResp, err := wh.fbClient.CallAPIGetMessage(&fbclient.GetMessageRequest{
		AccessToken: accessToken,
		MessageID:   mid,
		PageID:      externalPageID,
	})
	if err != nil {
		return err
	}

	externalUserID, err := wh.faboRedis.LoadPSID(externalPageID, PSID)
	switch err {

	// PSID not in redis then save
	case redis.ErrNil:
		{
			externalUserID = messageResp.From.ID
			if externalUserID == externalPageID {
				externalUserID = messageResp.To.Data[0].ID
			}

			if err := wh.faboRedis.SavePSID(externalPageID, PSID, externalUserID); err != nil {
				return err
			}
		}

	// PSID in redis then do nothing
	case nil:
		// no-op
	default:
		return err
	}

	var profileDefault *fbclientmodel.Profile
	// Get externalConversationID (externalPageID, externalUserID) from redis
	externalConversationID, err := wh.faboRedis.LoadExternalConversationID(externalPageID, externalUserID)
	switch err {
	// If externalConversationID not in redis then query database and save it into redis
	case redis.ErrNil:
		getExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery{
			ExternalPageID: externalPageID,
			ExternalUserID: externalUserID,
		}
		_err := wh.fbmessagingQuery.Dispatch(ctx, getExternalConversationQuery)
		switch cm.ErrorCode(_err) {
		case cm.NoError:
			externalConversation := getExternalConversationQuery.Result
			externalConversationID = externalConversation.ExternalID
			profileDefault = &fbclientmodel.Profile{
				ID:   externalConversation.ExternalUserID,
				Name: externalConversation.ExternalUserName,
			}
		case cm.NotFound:
			// if conversation not found then get externalConversation through call api
			// and create new externalConversation
			conversations, err := wh.fbClient.CallAPIGetConversationByUserID(&fbclient.GetConversationByUserIDRequest{
				AccessToken: accessToken,
				PageID:      externalPageID,
				UserID:      PSID,
			})
			if err != nil {
				return err
			}
			if len(conversations.ConversationsData) == 0 {
				return cm.Errorf(cm.Internal, nil, fmt.Sprintf("Wrong PSID %s", PSID))
			}
			externalConversation := conversations.ConversationsData[0]
			externalConversationID = externalConversation.ID

			var externalUserName string
			for _, sender := range externalConversation.Senders.Data {
				if sender.ID != externalPageID {
					externalUserName = sender.Name
					break
				}
			}

			if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalConversationsCommand{
				FbExternalConversations: []*fbmessaging.CreateFbExternalConversationArgs{
					{
						ID:                   cm.NewID(),
						ExternalPageID:       externalPageID,
						ExternalID:           externalConversation.ID,
						PSID:                 PSID,
						ExternalUserID:       externalUserID,
						ExternalUserName:     externalUserName,
						ExternalLink:         externalConversation.Link,
						ExternalUpdatedTime:  externalConversation.UpdatedTime.ToTime(),
						ExternalMessageCount: externalConversation.MessageCount,
					},
				},
			}); err != nil {
				return err
			}

			profileDefault = &fbclientmodel.Profile{
				ID:   externalUserID,
				Name: externalUserName,
			}
		default:
			return _err
		}

		if _err := wh.faboRedis.SaveExternalConversationID(externalPageID, externalUserID, externalConversationID); _err != nil {
			return _err
		}
	case nil:
		getExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalIDAndExternalPageIDQuery{
			ExternalPageID: externalPageID,
			ExternalID:     externalConversationID,
		}
		if _err := wh.fbmessagingQuery.Dispatch(ctx, getExternalConversationQuery); _err != nil {
			return _err
		}

		externalConversation := getExternalConversationQuery.Result
		profileDefault = &fbclientmodel.Profile{
			ID:   externalUserID,
			Name: externalConversation.ExternalUserName,
		}

	// no-op
	default:
		return err
	}

	profile, err := wh.getProfile(accessToken, externalPageID, PSID, profileDefault)
	if err != nil {
		return err
	}
	if profile.ProfilePic == "" {
		profile.ProfilePic = fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", PSID)
	}

	if messageResp.From.ID == PSID {
		messageResp.From.Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: profile.ProfilePic,
			},
		}

		messageResp.To.Data[0].Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", messageResp.To.Data[0].ID),
			},
		}
	} else {
		messageResp.From.Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", messageResp.From.ID),
			},
		}

		messageResp.To.Data[0].Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: profile.ProfilePic,
			},
		}
	}

	// Create new message
	var externalAttachments []*fbmessaging.FbMessageAttachment
	if messageResp.Attachments != nil {
		externalAttachments = fbclientconvert.ConvertMessageDataAttachments(messageResp.Attachments.Data)
	}
	if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalMessagesCommand{
		FbExternalMessages: []*fbmessaging.CreateFbExternalMessageArgs{
			{
				ID:                     cm.NewID(),
				ExternalConversationID: externalConversationID,
				ExternalPageID:         externalPageID,
				ExternalID:             messageResp.ID,
				ExternalMessage:        messageResp.Message,
				ExternalSticker:        messageResp.Sticker,
				ExternalTo:             fbclientconvert.ConvertObjectsTo(messageResp.To),
				ExternalFrom:           fbclientconvert.ConvertObjectFrom(messageResp.From),
				ExternalAttachments:    externalAttachments,
				ExternalCreatedTime:    messageResp.CreatedTime.ToTime(),
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
