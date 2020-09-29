package sender

import (
	"bytes"
	"context"
	"crypto/tls"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/webhook/storage"
	com "o.o/backend/com/main"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	callbackmodel "o.o/backend/pkg/etc/xmodel/callback/model"
	callbackstore "o.o/backend/pkg/etc/xmodel/callback/sqlstore"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

var client = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

const (
	PrefixGetPartnerIDs = "partnerIDs"
	TTL                 = 5 * 60
)

var redisStore redis.Store
var changesStore *storage.ChangesStore

type WebhookSender struct {
	db       *cmsql.Database
	ssenders map[dot.ID][]*SingleSender
	running  bool

	wg sync.WaitGroup
	m  sync.RWMutex

	PartnerStore sqlstore.PartnerStoreFactory
}

func New(db com.MainDB, redis redis.Store, cs *storage.ChangesStore, partnerStore sqlstore.PartnerStoreInterface) *WebhookSender {
	redisStore = redis
	changesStore = cs
	return &WebhookSender{
		db:           db,
		PartnerStore: sqlstore.NewPartnerStore(db),
	}
}

func (s *WebhookSender) Load() error {
	if s.ssenders != nil {
		ll.Fatal("already init")
	}

	var items callbackmodel.Webhooks
	if err := s.db.
		Where("deleted_at IS NULL").
		OrderBy("account_id").
		Find(&items); err != nil {
		return err
	}

	webhooks := make(map[dot.ID][]*SingleSender)
	for _, item := range items {
		ss := NewSingleSender(item)
		webhooks[item.AccountID] = append(webhooks[item.AccountID], ss)
	}
	s.ssenders = webhooks

	ll.Info("Loaded all ssenders from database")
	return nil
}

func (s *WebhookSender) Reload(ctx context.Context, accountID dot.ID) error {
	webhooks, err := callbackstore.Webhook(ctx, s.db).AccountID(accountID).List()
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

	ll.Info("reloaded ssenders", l.Int("n", len(webhooks)), l.ID("account_id", accountID))
	return nil
}

func findWebhook(items []*SingleSender, id dot.ID) *SingleSender {
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

// TODO: entityID: handle composition primary key (e.g. product collection relationship)
func (s *WebhookSender) CollectPb(ctx context.Context, entity entity_type.EntityType, entityID dot.ID, shopID dot.ID, accountIDs []dot.ID, pb capi.Message) (mq.Code, error) {
	var b bytes.Buffer
	if err := jsonx.MarshalTo(&b, pb); err != nil {
		ll.Error("error marshalling json", l.Error(err))
		return mq.CodeStop, err
	}

	if err := s.Collect(ctx, entity, entityID, shopID, accountIDs, b.Bytes()); err != nil {
		return mq.CodeStop, err
	}
	return mq.CodeOK, nil
}

func (s *WebhookSender) Collect(ctx context.Context, entity entity_type.EntityType, entityID dot.ID, shopID dot.ID, accountIDs []dot.ID, msg []byte) error {
	mapAccountIDs := make(map[dot.ID]bool)
	for _, accountID := range accountIDs {
		mapAccountIDs[accountID] = true
	}
	var partnerIDs []dot.ID
	value, err := redisStore.GetString(PrefixGetPartnerIDs + "-" + shopID.String())
	if err != nil && err != redis.ErrNil {
		return err
	}
	if value != "" { // use cache
		stringIDs := strings.Split(value, ",")
		for _, stringID := range stringIDs {
			id, _err := dot.ParseID(stringID)
			if _err != nil {
				panic("parse ID unexpected")
			}
			partnerIDs = append(partnerIDs, id)
		}
	} else { // reload from database
		query := &identitymodelx.GetPartnersFromRelationQuery{
			AccountIDs: []dot.ID{shopID},
		}
		if err := s.PartnerStore(ctx).GetPartnersFromRelation(ctx, query); err != nil {
			return err
		}

		if len(query.Result.Partners) > 0 {
			for _, partner := range query.Result.Partners {
				partnerIDs = append(partnerIDs, partner.ID)
			}
			value := dot.JoinIDs(partnerIDs)

			if err := redisStore.SetStringWithTTL(PrefixGetPartnerIDs+"-"+shopID.String(), value, TTL); err != nil {
				return err
			}
		}
	}

	for _, partnerID := range partnerIDs {
		if ok := mapAccountIDs[partnerID]; !ok {
			accountIDs = append(accountIDs, partnerID)
			mapAccountIDs[partnerID] = true
		}
	}

	ll.Debug("Collect items for accounts", l.Any("ids", accountIDs))

	for _, accountID := range accountIDs {
		if accountID == 0 {
			continue
		}

		s.m.RLock()
		whs := s.ssenders[accountID]
		s.m.RUnlock()
		for _, wh := range whs {
			if entity_type.Contain(wh.webhook.Entities, entity) {
				wh.Collect(entity, entityID, msg)
			}
		}
	}
	return nil
}

func (s *WebhookSender) ResetState(accountID dot.ID) error {
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
