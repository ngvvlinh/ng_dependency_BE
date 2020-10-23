package webhook

import (
	"fmt"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/com/fabo/pkg/fbclient"
	faboredis "o.o/backend/com/fabo/pkg/redis"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/common/l"
)

const oneHour = 1 * time.Hour

var ll = l.New().WithChannel("webhook")

type Webhook struct {
	db                     *cmsql.Database
	dbLog                  *cmsql.Database
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
	dbLog com.LogDB,
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
		dbLog:                  dbLog,
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

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	var webhookMessages WebhookMessages
	if err := c.DecodeJson(&webhookMessages); err != nil {
		return err
	}
	defer c.Req.Body.Close()

	ctx := c.Context()
	go func() { defer cm.RecoverAndLog(); wh.forwardWebhook(c, webhookMessages) }()

	defer func() {
		writer := c.SetResultRaw()
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("oke"))
		writer.WriteHeader(200)

		wh.saveLogsWebhook(webhookMessages, _err)
	}()

	if webhookMessages.Object != "page" {
		return nil
	}

	switch webhookMessages.MessageType() {
	case WebhookMessage:
		return wh.handleMessenger(ctx, webhookMessages)
	case WebhookFeed:
		return wh.handleFeed(ctx, webhookMessages)
	case WebhookInvalidMessage:
		return nil
	default:
		return nil
	}
}
