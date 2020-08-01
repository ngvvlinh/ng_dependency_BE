package sqlstore

import (
	"context"

	"o.o/backend/com/etc/logging/shippingwebhook/model"
	"o.o/backend/pkg/common/sql/cmsql"
)

type ShippingWebhookStoreFactory func(context.Context) *ShippingWebhookStore

func NewShippingWebhookStore(db *cmsql.Database) ShippingWebhookStoreFactory {
	return func(ctx context.Context) *ShippingWebhookStore {
		return &ShippingWebhookStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShippingWebhookStore struct {
	query cmsql.QueryFactory
	ft    ShippingProviderWebhookFilters
}

func (s *ShippingWebhookStore) CreateShippingWebhookLog(log *model.ShippingProviderWebhook) error {
	return s.query().ShouldInsert(log)
}
