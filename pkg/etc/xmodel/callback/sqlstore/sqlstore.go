package sqlstore

import (
	"context"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etc/xmodel/callback/model"
	callbackmodel "o.o/backend/pkg/etc/xmodel/callback/model"
	"o.o/capi/dot"
)

type WebhookStore struct {
	db    *cmsql.Database
	ctx   context.Context
	ft    WebhookFilters
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
	multiplelity
}

func Webhook(ctx context.Context, db *cmsql.Database) *WebhookStore {
	return &WebhookStore{db: db, ctx: ctx}
}

func (s *WebhookStore) ID(id dot.ID) *WebhookStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *WebhookStore) AccountID(id dot.ID) *WebhookStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *WebhookStore) IncludeDeleted() *WebhookStore {
	s.includeDeleted = true
	return s
}

func (s *WebhookStore) Multiple() *WebhookStore {
	s.multiplelity = true
	return s
}

func (s *WebhookStore) Get() (*callbackmodel.Webhook, error) {
	var item callbackmodel.Webhook
	err := s.db.Where(s.preds...).Where(s.includeDeleted.FilterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *WebhookStore) List() ([]*callbackmodel.Webhook, error) {
	var items model.Webhooks
	err := s.db.Where(s.preds...).Where(s.includeDeleted.FilterDeleted(&s.ft)).Find(&items)
	return items, err
}

func (s *WebhookStore) Count() (int, error) {
	return s.db.Where(s.preds...).Where(s.includeDeleted.FilterDeleted(&s.ft)).Count((*callbackmodel.Webhook)(nil))
}

func (s *WebhookStore) Create(item *callbackmodel.Webhook) error {
	if err := item.BeforeInsert(); err != nil {
		return err
	}
	return s.db.ShouldInsert(item)
}

func (s *WebhookStore) SoftDelete() error {
	if err := s.ensureMultiplelity(s); err != nil {
		return err
	}
	return s.db.Where(s.preds...).ShouldUpdate(&callbackmodel.Webhook{
		DeletedAt: time.Now(),
	})
}

type multiplelity bool

func (m multiplelity) ensureMultiplelity(countable interface{ Count() (int, error) }) error {
	n, err := countable.Count()
	if err != nil {
		return err
	}
	if !m && (n > 1) {
		return cm.Errorf(cm.Internal, nil, "unexpected number of changes")
	}
	return nil
}
