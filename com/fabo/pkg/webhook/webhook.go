package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/webhook_type"
	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	faboredis "o.o/backend/com/fabo/pkg/redis"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/common/l"
	"o.o/common/xerrors"
)

const oneHour = 1 * time.Hour

var ll = l.New()

type Webhook struct {
	db                     *cmsql.Database
	webhookCallbackService *sadmin.WebhookCallbackService
	verifyToken            string
	faboRedis              *faboredis.FaboRedis
	fbClient               *fbclient.FbClient
	fbmessagingQuery       fbmessaging.QueryBus
	fbmessagingAggr        fbmessaging.CommandBus
	fbPageQuery            fbpaging.QueryBus

	mu                               sync.RWMutex
	mapCallbackURLAndLatestTimeError map[string]time.Time
}

func New(
	db com.MainDB,
	rd redis.Store,
	cfg config.WebhookConfig,
	faboRedis *faboredis.FaboRedis,
	fbClient *fbclient.FbClient,
	fbmessagingQuery fbmessaging.QueryBus,
	fbmessagingAggregate fbmessaging.CommandBus,
	fbPageQuery fbpaging.QueryBus,
) *Webhook {
	wh := &Webhook{
		db:                     db,
		webhookCallbackService: sadmin.NewWebhookCallbackService(rd),
		verifyToken:            cfg.VerifyToken,
		faboRedis:              faboRedis,
		fbClient:               fbClient,
		fbmessagingQuery:       fbmessagingQuery,
		fbmessagingAggr:        fbmessagingAggregate,
		fbPageQuery:            fbPageQuery,

		mapCallbackURLAndLatestTimeError: make(map[string]time.Time),
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.GET("/webhook/fbmessenger/:id", wh.HandleWebhookVerification)
	rt.POST("/webhook/fbmessenger/:id", wh.Callback)

	// backward-compatible
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

	go func() { defer cm.RecoverAndLog(); wh.forwardWebhook(c, body) }()

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
					err := wh.handleMessageReturned(ctx, externalPageID, PSID, mid)
					if err == nil {
						continue
					}
					facebookError := err.(*xerrors.APIError)
					code := facebookError.Meta["code"]
					if code == fbclient.AccessTokenHasExpired.String() {
						continue
					} else {
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

func (wh *Webhook) forwardWebhook(c *httpx.Context, body []byte) {
	callbackURLs, err := wh.webhookCallbackService.GetWebhookCallbackURLs(webhook_type.Fabo.String())
	if err != nil {
		return
	}

	client := http.Client{
		Timeout: 3 * time.Second,
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

			req, err2 := http.NewRequest("POST", callbackURL, bytes.NewReader(body))
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

// TODO: Ngoc
func (wh *Webhook) handleMessageReturned(ctx context.Context, externalPageID, PSID, mid string) error {
	// Get page (active) by externalPageID
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

	// Get token page to get message
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

	// Get message
	messageResp, err := wh.fbClient.CallAPIGetMessage(accessToken, mid)
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
	// PSID not in redis then do nothing
	case nil:
		// no-op
	default:
		return err
	}

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
			externalConversationID = getExternalConversationQuery.Result.ExternalID
		case cm.NotFound:
			// if conversation not found then get externalConversation through call api
			// and create new externalConversation
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

		if _err := wh.faboRedis.SaveExternalConversationID(externalPageID, externalUserID, externalConversationID); _err != nil {
			return _err
		}
	case nil:
	// no-op
	default:
		return err
	}

	profile, err := wh.fbClient.CallAPIGetProfileByPSID(accessToken, PSID)
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

func (wh *Webhook) getProfileByPSID(accessToken, pageID, PSID string) (profile *fbclientmodel.Profile, err error) {
	profile, err = wh.faboRedis.LoadProfilePSID(pageID, PSID)
	switch err {
	case redis.ErrNil:
		profile, err = wh.fbClient.CallAPIGetProfileByPSID(accessToken, PSID)
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
