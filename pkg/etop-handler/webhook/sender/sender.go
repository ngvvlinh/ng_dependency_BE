package sender

import (
	"bytes"
	"context"
	"crypto/tls"
	"math/rand"
	"net/http"
	"sync"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop-handler/webhook/storage"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var ll = l.New()

var client = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

var redisStore redis.Store
var changesStore *storage.ChangesStore

type WebhookSender struct {
	db       cmsql.Database
	ssenders map[int64][]*SingleSender
	running  bool

	wg sync.WaitGroup
	m  sync.RWMutex
}

func New(db cmsql.Database, redis redis.Store, cs *storage.ChangesStore) *WebhookSender {
	redisStore = redis
	changesStore = cs
	return &WebhookSender{db: db}
}

func (s *WebhookSender) Load() error {
	if s.ssenders != nil {
		ll.Fatal("already init")
	}

	var items model.Webhooks
	if err := s.db.
		Where("deleted_at IS NULL").
		OrderBy("account_id").
		Find(&items); err != nil {
		return err
	}

	webhooks := make(map[int64][]*SingleSender)
	for _, item := range items {
		ss := NewSingleSender(item)
		webhooks[item.AccountID] = append(webhooks[item.AccountID], ss)
	}
	s.ssenders = webhooks

	ll.Info("Loaded all ssenders from database")
	return nil
}

func (s *WebhookSender) Reload(ctx context.Context, accountID int64) error {
	webhooks, err := sqlstore.Webhook(ctx).AccountID(accountID).List()
	if err != nil {
		ll.Error("webhook/reload", l.Error(err))
		return err
	}

	s.m.Lock()
	defer s.m.Unlock()

	// compare and reload
	news := NewSingleSenders(webhooks)
	ssenders := s.ssenders[accountID]

	for _, ss := range ssenders {
		if findWebhook(news, ss.webhook.ID) == nil {
			// stop the old ones if not found in news
			close(ss.stop)
		}
	}

	newSenders := make([]*SingleSender, len(news))
	for i, newOne := range news {
		if ss := findWebhook(ssenders, newOne.webhook.ID); ss != nil {
			// keep the running ones
			newSenders[i] = ss

		} else {
			// and start the new ones
			newSenders[i] = newOne

			// TODO: use main context instead of this ctx
			s.startOne(ctx, newOne)
		}
	}
	s.ssenders[accountID] = newSenders

	ll.Info("reloaded ssenders", l.Int("n", len(webhooks)), l.Int64("account_id", accountID))
	return nil
}

func findWebhook(items []*SingleSender, id int64) *SingleSender {
	for _, item := range items {
		if item.webhook.ID == id {
			return item
		}
	}
	return nil
}

func (s *WebhookSender) Start(ctx context.Context) {
	if s.running {
		ll.Fatal("already start")
	}
	s.running = true

	s.m.Lock()
	defer s.m.Unlock()
	for _, whs := range s.ssenders {
		for _, ss := range whs {
			s.startOne(ctx, ss)
		}
	}
}

func (s *WebhookSender) startOne(ctx context.Context, ss *SingleSender) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// wait for a random interval in 5s before starting
		after := time.Duration(rand.Intn(5000)) * time.Millisecond
		ss.Run(ctx, after)
	}()
}

func (s *WebhookSender) Wait() {
	s.wg.Wait()
}

var marshaler = jsonpb.Marshaler{OrigName: true, EmitDefaults: false}

func (s *WebhookSender) CollectPb(ctx context.Context, entity string, entityID int64, accountIDs []int64, pb proto.Message) (mq.Code, error) {
	var b bytes.Buffer
	if err := marshaler.Marshal(&b, pb); err != nil {
		ll.Error("error marshalling json", l.Error(err))
		return mq.CodeStop, err
	}

	s.Collect(ctx, entity, entityID, accountIDs, b.Bytes())
	return mq.CodeOK, nil
}

func (s *WebhookSender) Collect(ctx context.Context, entity string, entityID int64, accountIDs []int64, msg []byte) {
	ll.Debug("Collect items for accounts", l.Any("ids", accountIDs))
	for _, accountID := range accountIDs {
		if accountID == 0 {
			continue
		}

		s.m.RLock()
		whs := s.ssenders[accountID]
		s.m.RUnlock()
		for _, wh := range whs {
			if cm.StringsContain(wh.webhook.Entities, entity) {
				wh.Collect(entity, entityID, msg)
			}
		}
	}
}

func (s *WebhookSender) ResetState(accountID int64) error {
	s.m.RLock()
	defer s.m.RUnlock()

	ss := s.ssenders[accountID]
	if ss == nil {
		return cm.Error(cm.NotFound, "no active account", nil)
	}

	for _, sender := range ss {
		sender.ResetState()
	}
	return nil
}
