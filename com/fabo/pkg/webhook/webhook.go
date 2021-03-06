package webhook

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/com/fabo/pkg/fbclient"
	faboredis "o.o/backend/com/fabo/pkg/redis"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/common/jsonx"
	"o.o/common/l"
)

const oneHour = 1 * time.Hour

var ll = l.New().WithChannel("webhook")

type WebhookHandler struct {
	db               com.MainDB
	faboRedis        *faboredis.FaboRedis
	fbClient         *fbclient.FbClient
	fbmessagingQuery fbmessaging.QueryBus
	fbmessagingAggr  fbmessaging.CommandBus
	fbPageQuery      fbpaging.QueryBus
	fbUserQuery      fbusering.QueryBus
	jobKeeper        *JobKeeper
}

func NewWebhookHandler(
	db com.MainDB,
	faboRedis *faboredis.FaboRedis,
	fbClient *fbclient.FbClient,
	fbmessagingQuery fbmessaging.QueryBus,
	fbmessagingAggregate fbmessaging.CommandBus,
	fbPageQuery fbpaging.QueryBus,
	fbUserQuery fbusering.QueryBus,
) *WebhookHandler {
	wh := &WebhookHandler{
		db:               db,
		faboRedis:        faboRedis,
		fbClient:         fbClient,
		fbmessagingQuery: fbmessagingQuery,
		fbmessagingAggr:  fbmessagingAggregate,
		fbPageQuery:      fbPageQuery,
		fbUserQuery:      fbUserQuery,
		jobKeeper:        NewJobKeeper(),
	}
	return wh
}

type Webhook struct {
	//db                     *cmsql.Database
	dbLog                  *cmsql.Database
	webhookCallbackService *sadmin.WebhookCallbackService
	verifyToken            string
	faboRedis              *faboredis.FaboRedis
	producer               *mq.KafkaProducer
	prefix                 string

	mu                               sync.RWMutex
	mapCallbackURLAndLatestTimeError map[string]time.Time
}

func New(
	dbLog com.LogDB,
	rd redis.Store,
	cfg config.WebhookConfig,
	faboRedis *faboredis.FaboRedis,
	producer *mq.KafkaProducer,
	kafka cc.Kafka,
) *Webhook {
	wh := &Webhook{
		dbLog:                  dbLog,
		webhookCallbackService: sadmin.NewWebhookCallbackService(rd),
		verifyToken:            cfg.VerifyToken,
		faboRedis:              faboRedis,
		producer:               producer,
		prefix:                 kafka.TopicPrefix + "_pgrid_",

		mapCallbackURLAndLatestTimeError: make(map[string]time.Time),
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.GET("/webhook/fbmessenger/:id", wh.HandleWebhookVerification)
	rt.POST("/webhook/fbmessenger/:id", wh.Callback)

	rt.GET("/webhook/fbuser/:id", wh.HandleWebhookVerification)
	rt.POST("/webhook/fbuser/:id", wh.CallbackForUser)

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
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	ll.Info("->"+c.Req.URL.Path, l.String("data", string(body)))

	err = jsonx.Unmarshal(body, &webhookMessages)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}

	go func() { defer cm.RecoverAndLog(); wh.forwardWebhook(c, webhookMessages) }()

	defer func() {
		writer := c.SetResultRaw()
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("oke"))
		writer.WriteHeader(200)

		wh.saveLogsWebhookPage(webhookMessages, _err)
	}()

	if webhookMessages.Object != "page" {
		return nil
	}

	switch webhookMessages.MessageType() {
	case WebhookMessage:
		topic := wh.prefix + "facebook_webhook_message"
		wh.produceMessage(topic, webhookMessages.GetKey(), 64, body)
		return nil
	case WebhookFeed:
		topic := wh.prefix + "facebook_webhook_feed"
		wh.produceMessage(topic, webhookMessages.GetKey(), 64, body)
		return nil
	case WebhookInvalidMessage:
		return nil
	default:
		return nil
	}
}

func (wh *Webhook) CallbackForUser(c *httpx.Context) (_err error) {
	var webhookUser WebhookUser
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	ll.Info("->"+c.Req.URL.Path, l.String("data", string(body)))

	if err := jsonx.Unmarshal(body, &webhookUser); err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}

	ll.SendMessagef("fbUser: %v", string(body))

	defer func() {
		writer := c.SetResultRaw()
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("oke"))
		writer.WriteHeader(200)

		wh.saveLogsWebhookUser(webhookUser, _err)
	}()

	if webhookUser.Object != "user" {
		return nil
	}

	switch webhookUser.Type() {
	case WebhookUserLiveVideos:
		topic := wh.prefix + "facebook_webhook_user_live_video"
		wh.produceMessage(topic, webhookUser.GetKey(), 64, body)
		return nil
	}
	return nil
}

func (wh *Webhook) produceMessage(topic, key string, partitions int, message []byte) {
	partition := hash(key, partitions)
	wh.producer.Send(topic, int(partition), key, message)
}

func hash(s string, modular int) int {
	algorithm := fnv.New32()
	algorithm.Write([]byte(s))
	return int(algorithm.Sum32()) % modular
}
