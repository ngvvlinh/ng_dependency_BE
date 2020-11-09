package sqlstore

import (
	"context"

	"o.o/backend/com/etc/logging/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
)

type TicketLogStoreFactory func(context.Context) *TicketLogStore

func NewTicketLogStore(db *cmsql.Database) TicketLogStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketLogStore {
		return &TicketLogStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketLogStore struct {
	query cmsql.QueryFactory
}

func (s *TicketLogStore) CreateTicketLog(ticketLog *model.TicketProviderWebhook) error {
	_, err := s.query().Insert(ticketLog)
	return err
}
