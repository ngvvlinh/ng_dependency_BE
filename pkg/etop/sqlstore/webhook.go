package sqlstore

import (
	"context"
	"time"

	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type WebhookStore struct {
	ctx   context.Context
	ft    WebhookFilters
	preds []interface{}

	includeDeleted
	multiplelity
}

func Webhook(ctx context.Context) *WebhookStore {
	return &WebhookStore{ctx: ctx}
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

func (s *WebhookStore) Get() (*model.Webhook, error) {
	var item model.Webhook
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *WebhookStore) List() ([]*model.Webhook, error) {
	var items model.Webhooks
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).Find(&items)
	return items, err
}

func (s *WebhookStore) Count() (int, error) {
	return x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).Count((*model.Webhook)(nil))
}

func (s *WebhookStore) Create(item *model.Webhook) error {
	if err := item.BeforeInsert(); err != nil {
		return err
	}
	return x.ShouldInsert(item)
}

func (s *WebhookStore) SoftDelete() error {
	if err := s.ensureMultiplelity(s); err != nil {
		return err
	}
	return x.Where(s.preds...).ShouldUpdate(&model.Webhook{
		DeletedAt: time.Now(),
	})
}
