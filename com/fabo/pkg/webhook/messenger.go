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
	"o.o/backend/pkg/common/mq"
	"o.o/common/xerrors"
)

func (wh *WebhookHandler) HandleMessenger(
	ctx context.Context, webhookMessages WebhookMessages,
) (mq.Code, error) {
	entries := webhookMessages.Entry
	for _, entry := range entries {
		externalPageID := entry.ID
		for _, messaging := range entry.Messaging {
			if messaging.Message == nil {
				continue
			}

			externalUserID := messaging.Sender.ID
			if externalUserID == externalPageID {
				externalUserID = messaging.Recipient.ID
			}

			returnCode, err := wh.handleMessageReturned(ctx, messaging, externalPageID, externalUserID)
			if err == nil {
				continue
			}

			// handle error
			facebookError, ok := err.(*xerrors.APIError)
			if ok {
				code := facebookError.Meta["code"]
				if code == fbclient.AccessTokenHasExpired.String() {
					continue
				}
			}

			ll.SendMessage(err.Error())
			return returnCode, err
		}
	}
	return mq.CodeOK, nil
}

// TODO(ngoc): refactor
func (wh *WebhookHandler) handleMessageReturned(
	ctx context.Context, messaging *Messaging,
	externalPageID, externalUserID string,
) (mq.Code, error) {
	if messaging == nil || messaging.Message == nil {
		return mq.CodeIgnore, nil
	}
	externalTimestamp := messaging.Timestamp

	// handle test page
	if code, err := wh.handleTestPage(ctx, externalPageID); code != mq.CodeOK {
		return code, err
	}

	// convert message webhook to message model
	var messageData *fbclientmodel.MessageData
	{
		message := messaging.Message
		messageData = message.ConvertToMessageData(messaging.Sender.ID, messaging.Recipient.ID, externalTimestamp)
		if messageData == nil {
			return mq.CodeIgnore, nil
		}
	}

	pageAccessToken, returnCode, err := wh.getPageAccessToken(ctx, externalPageID)
	if returnCode != mq.CodeOK {
		return returnCode, err
	}

	// Each message belongs to a conversation
	// And conversation is defined by externalConversationID that unique on (externalPageID, externalUserID)
	externalConversation, err := wh.getExternalConversation(ctx, externalPageID, externalUserID, pageAccessToken)
	if err != nil {
		return mq.CodeRetry, err
	}

	// Add information for From and To of messageData
	if err := wh.addInfosForFromAndTo(pageAccessToken, externalConversation, messageData); err != nil {
		return mq.CodeRetry, err
	}

	// Create new message
	var externalAttachments []*fbmessaging.FbMessageAttachment
	var externalShares []*fbmessaging.FbMessageShare
	if messageData.Attachments != nil {
		externalAttachments = fbclientconvert.ConvertMessageDataAttachments(messageData.Attachments.Data)
	}
	if messageData.Shares != nil {
		externalShares = fbclientconvert.ConvertMessageShares(messageData.Shares.Data)
	}

	currentMessage := messageData.Message

	if err := wh.fbmessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalMessagesCommand{
		FbExternalMessages: []*fbmessaging.CreateFbExternalMessageArgs{
			{
				ID:                     cm.NewID(),
				ExternalConversationID: externalConversation.ExternalID,
				ExternalPageID:         externalPageID,
				ExternalID:             messageData.ID,
				ExternalMessage:        currentMessage,
				ExternalSticker:        messageData.Sticker,
				ExternalTo:             fbclientconvert.ConvertObjectsTo(messageData.To),
				ExternalFrom:           fbclientconvert.ConvertObjectFrom(messageData.From),
				ExternalAttachments:    externalAttachments,
				ExternalMessageShares:  externalShares,
				ExternalCreatedTime:    messageData.CreatedTime.WebhookTimeToTime(),
				ExternalTimestamp:      int64(externalTimestamp),
			},
		},
	}); err != nil {
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func (wh *WebhookHandler) addInfosForFromAndTo(
	pageAccessToken string,
	externalConversation *fbmessaging.FbExternalConversation, messageData *fbclientmodel.MessageData,
) error {
	externalPageID := externalConversation.ExternalPageID
	externalUserID := externalConversation.ExternalUserID

	profileUserDefault := &fbclientmodel.Profile{
		ID:   externalConversation.ExternalUserID,
		Name: externalConversation.ExternalUserName,
	}
	profile, err := wh.getProfile(pageAccessToken, externalPageID, externalUserID, profileUserDefault)
	if err != nil {
		return err
	}
	if profile.ProfilePic == "" {
		profile.ProfilePic = defaultAvatar(externalUserID)
	}

	if messageData.From.ID == externalUserID {
		messageData.From.Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: profile.ProfilePic,
			},
		}

		messageData.To.Data[0].Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: defaultAvatar(externalPageID),
			},
		}
	} else {
		messageData.From.Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: defaultAvatar(externalPageID),
			},
		}

		messageData.To.Data[0].Picture = &fbclientmodel.Picture{
			Data: fbclientmodel.PictureData{
				Url: profile.ProfilePic,
			},
		}
	}
	return nil
}

func (wh *WebhookHandler) getExternalConversation(
	ctx context.Context,
	externalPageID, externalUserID, pageAccessToken string,
) (*fbmessaging.FbExternalConversation, error) {
	// Get externalConversationID (externalPageID, externalUserID) from redis
	externalConversationCached, err := wh.faboRedis.LoadExternalConversation(externalPageID, externalUserID)
	if err != nil {
		return nil, err
	}

	// if ExternalConversation in cache
	if externalConversationCached != nil {
		return externalConversationCached, nil
	}

	// Query FbExternalConversation from DB
	getOldExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery{
		ExternalPageID: externalPageID,
		ExternalUserID: externalUserID,
	}
	_err := wh.fbmessagingQuery.Dispatch(ctx, getOldExternalConversationQuery)

	var externalConversation *fbmessaging.FbExternalConversation
	switch cm.ErrorCode(_err) {
	case cm.NoError:
		externalConversation = getOldExternalConversationQuery.Result
	case cm.NotFound:
		// if conversation not found then get externalConversation through call api
		// and create new externalConversationFromAPI
		conversations, err := wh.fbClient.CallAPIGetConversationByUserID(&fbclient.GetConversationByUserIDRequest{
			AccessToken: pageAccessToken,
			PageID:      externalPageID,
			UserID:      externalUserID,
		})
		if err != nil {
			return nil, err
		}
		if len(conversations.ConversationsData) == 0 {
			return nil, cm.Errorf(cm.Internal, nil, fmt.Sprintf("Wrong externalUserID %s", externalUserID))
		}
		if len(conversations.ConversationsData) > 1 {
			return nil, cm.Error(cm.Internal, "Something wrong", nil)
		}

		externalConversationFromAPI := conversations.ConversationsData[0]
		var externalUserName string
		for _, sender := range externalConversationFromAPI.Senders.Data {
			if sender.ID != externalPageID {
				externalUserName = sender.Name
				break
			}
		}

		// create new FbExternalConversation
		createFbExternalConversationCommand := &fbmessaging.CreateOrUpdateFbExternalConversationCommand{
			ID:                   cm.NewID(),
			ExternalPageID:       externalPageID,
			ExternalID:           externalConversationFromAPI.ID,
			PSID:                 externalUserID,
			ExternalUserID:       externalUserID,
			ExternalUserName:     externalUserName,
			ExternalLink:         externalConversationFromAPI.Link,
			ExternalUpdatedTime:  externalConversationFromAPI.UpdatedTime.ToTime(),
			ExternalMessageCount: externalConversationFromAPI.MessageCount,
		}
		if err := wh.fbmessagingAggr.Dispatch(ctx, createFbExternalConversationCommand); err != nil {
			return nil, err
		}

		externalConversation = createFbExternalConversationCommand.Result
	default:
		return nil, _err
	}
	// cache externalConversation
	if err := wh.faboRedis.SaveExternalConversation(externalPageID, externalUserID, *externalConversation); err != nil {
		return nil, err
	}

	return externalConversation, nil
}

func (wh *WebhookHandler) handleTestPage(ctx context.Context, externalPageID string) (mq.Code, error) {
	// check page is a test page
	isTestPage, err := wh.IsTestPage(ctx, externalPageID)
	if err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return mq.CodeIgnore, nil
		}
		return mq.CodeRetry, err
	}

	// ignore test page on production
	if cmenv.IsProd() && isTestPage {
		return mq.CodeIgnore, nil
	}

	return mq.CodeOK, nil
}

func defaultAvatar(externalUserID string) string {
	return fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", externalUserID)
}
