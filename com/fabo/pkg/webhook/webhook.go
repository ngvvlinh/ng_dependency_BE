package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

const (
	PREFIX_PSID                  = "psid"
	PREFIX_EXTERNAL_CONVERSATION = "external_conversation"
)

var (
	ll = l.New()
)

type Webhook struct {
	db               *cmsql.Database
	bot              *telebot.Channel
	verifyToken      string
	redisStore       redis.Store
	fbClient         *fbclient.FbClient
	fbmessagingQuery *fbmessaging.QueryBus
	fbmessagingAggr  *fbmessaging.CommandBus
	fbPageQuery      *fbpaging.QueryBus
}

func New(
	db *cmsql.Database, bot *telebot.Channel, verifyToken string,
	redisStore redis.Store, fbClient *fbclient.FbClient,
	fbmessagingQuery *fbmessaging.QueryBus, fbmessagingAggregate *fbmessaging.CommandBus,
	fbPageQuery *fbpaging.QueryBus,
) *Webhook {
	wh := &Webhook{
		db:               db,
		bot:              bot,
		verifyToken:      verifyToken,
		redisStore:       redisStore,
		fbClient:         fbClient,
		fbmessagingQuery: fbmessagingQuery,
		fbmessagingAggr:  fbmessagingAggregate,
		fbPageQuery:      fbPageQuery,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.GET("/webhook/fbmessager/:id", wh.HandleWebhookVerification)
	rt.POST("/webhook/fbmessager/:id", wh.Callback)
}

func (wh *Webhook) HandleWebhookVerification(c *httpx.Context) error {
	mode := c.Req.URL.Query().Get("hub.mode")
	token := c.Req.URL.Query().Get("hub.verify_token")
	challenge := c.Req.URL.Query().Get("hub.challenge")

	writer := c.SetResultRaw()

	if mode != "" && token != "" {
		if mode == "subscribe" && token == wh.verifyToken {
			fmt.Println("WEBHOOK_VERIFIED")

			writer.Write([]byte(challenge))
			writer.WriteHeader(200)
		} else {
			writer.WriteHeader(403)
		}
	}
	return nil
}

func (wh *Webhook) Callback(c *httpx.Context) error {
	body, err := ioutil.ReadAll(c.Req.Body)
	defer c.Req.Body.Close()
	if err != nil {
		return err
	}
	ctx := c.Context()

	var webhookMessages WebhookMessages

	if err := json.Unmarshal(body, &webhookMessages); err != nil {
		return err
	}

	if webhookMessages.Object == "page" && len(webhookMessages.Entry) > 0 {
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
					if err := wh.handleMessageReturned(ctx, externalPageID, PSID, mid); err != nil {
						return err
					}
				}
			}
		}
	}

	writer := c.SetResultRaw()
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("EVENT_RECEIVED"))
	writer.WriteHeader(200)

	return nil
}

// TODO: Ngoc
func (wh *Webhook) handleMessageReturned(ctx context.Context, externalPageID, PSID, mid string) error {
	getFbExternalPageActiveByExternalIDQuery := &fbpaging.GetFbExternalPageActiveByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := wh.fbPageQuery.Dispatch(ctx, getFbExternalPageActiveByExternalIDQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil
		}
		return err
	}
	fbPageID := getFbExternalPageActiveByExternalIDQuery.Result.ID

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByIDQuery{
		ID: fbPageID,
	}
	if err := wh.fbPageQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil
		}
		return err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	messageResp, err := wh.fbClient.CallAPIGetMessage(accessToken, mid)
	if err != nil {
		return err
	}

	externalUserID, err := wh.loadPSID(externalPageID, PSID)
	switch err {
	case redis.ErrNil:
		{
			externalUserID = messageResp.From.ID

			if err := wh.savePSID(externalPageID, PSID, externalUserID); err != nil {
				return err
			}
		}
	case nil:
		// no-op
	default:
		return err
	}

	externalConversationID, err := wh.loadExternalConversationID(externalPageID, externalUserID)
	switch err {
	case redis.ErrNil:
		getExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery{
			ExternalPageID: externalPageID,
			ExternalUserID: externalUserID,
		}
		_err := wh.fbmessagingQuery.Dispatch(ctx, getExternalConversationQuery)

		switch cm.ErrorCode(_err) {
		case cm.NoError:
			externalConversationID = getExternalConversationQuery.Result.ExternalID
		case cm.NotFound:
			conversations, err := wh.fbClient.CallAPIListConversations(accessToken, externalPageID, PSID, nil)
			if err != nil {
				return err
			}
			if len(conversations.Conversations.ConversationsData) == 0 {
				return cm.Errorf(cm.Internal, nil, fmt.Sprintf("Wrong PSID %s", PSID))
			}
			externalConversation := conversations.Conversations.ConversationsData[0]
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
		default:
			return _err
		}

		if _err := wh.saveExternalConversationID(externalPageID, externalUserID, externalConversationID); _err != nil {
			return _err
		}
	case nil:
	// no-op
	default:
		return err
	}

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

func (wh *Webhook) loadPSID(pageID, PSID string) (string, error) {
	return wh.redisStore.GetString(wh.generatePSIDKey(pageID, PSID))
}

func (wh *Webhook) savePSID(pageID, PSID, externalUserID string) error {
	return wh.redisStore.SetString(wh.generatePSIDKey(pageID, PSID), externalUserID)
}

func (wh *Webhook) generatePSIDKey(externalPageID, psid string) string {
	return fmt.Sprintf("%s:%s_%s", PREFIX_PSID, externalPageID, psid)
}

func (wh *Webhook) loadExternalConversationID(externalPageID, externalUserID string) (string, error) {
	return wh.redisStore.GetString(wh.generateExternalConversationKey(externalPageID, externalUserID))
}

func (wh *Webhook) saveExternalConversationID(externalPageID, externalUserID, externalConversationID string) error {
	return wh.redisStore.SetString(wh.generateExternalConversationKey(externalPageID, externalUserID), externalConversationID)
}

func (wh *Webhook) generateExternalConversationKey(externalPageID, externalUserID string) string {
	return fmt.Sprintf("%s:%s_%s", PREFIX_EXTERNAL_CONVERSATION, externalPageID, externalUserID)
}

type WebhookMessages struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

type Entry struct {
	ID        string                     `json:"id"`
	Time      fbclientmodel.FacebookTime `json:"time"`
	Messaging []*Messaging               `json:"messaging"`
}

type Messaging struct {
	Sender    *Sender                    `json:"sender"`
	Recipient *Recipient                 `json:"recipient"`
	Timestamp fbclientmodel.FacebookTime `json:"timestamp"`
	Message   *Message                   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Mid         string               `json:"mid"`
	Text        string               `json:"text"`
	StickerID   string               `json:"sticher_id"`
	Attachments []*MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	Title   string                 `json:"title"`
	Url     string                 `json:"url"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}
