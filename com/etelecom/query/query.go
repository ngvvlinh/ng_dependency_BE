package query

import (
	"context"

	"o.o/api/etelecom"
	"o.o/backend/com/etelecom/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
)

var _ etelecom.QueryService = &QueryService{}

type QueryService struct {
	hotlineStore   sqlstore.HotlineStoreFactory
	extensionStore sqlstore.ExtensionStoreFactory
}

func NewQueryService(dbEtelecom com.EtelecomDB) *QueryService {
	return &QueryService{
		hotlineStore:   sqlstore.NewHotlineStore(dbEtelecom),
		extensionStore: sqlstore.NewExtensionStore(dbEtelecom),
	}
}

func QueryServiceMessageBus(s *QueryService) etelecom.QueryBus {
	b := bus.New()
	return etelecom.NewQueryServiceHandler(s).RegisterHandlers(b)
}

func (q *QueryService) GetHotline(ctx context.Context, args *etelecom.GetHotlineArgs) (*etelecom.Hotline, error) {
	return q.hotlineStore(ctx).ID(args.ID).GetHotline()
}

func (q *QueryService) ListHotlines(ctx context.Context, args *etelecom.ListHotlinesArgs) ([]*etelecom.Hotline, error) {
	return q.hotlineStore(ctx).OptionalUserID(args.UserID).OptionalConnectionID(args.ConnectionID).ListHotlines()
}
