package query

import (
	"context"

	"o.o/api/etelecom"
	"o.o/backend/com/etelecom/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ etelecom.QueryService = &QueryService{}

type QueryService struct {
	db             *cmsql.Database
	hotlineStore   sqlstore.HotlineStoreFactory
	extensionStore sqlstore.ExtensionStoreFactory
}

func NewQueryService(dbEtelecom com.EtelecomDB) *QueryService {
	return &QueryService{
		db:             dbEtelecom,
		hotlineStore:   sqlstore.NewHotlineStore(dbEtelecom),
		extensionStore: sqlstore.NewExtensionStore(dbEtelecom),
	}
}

func QueryServiceMessageBus(s *QueryService) etelecom.QueryBus {
	b := bus.New()
	return etelecom.NewQueryServiceHandler(s).RegisterHandlers(b)
}

func (q *QueryService) GetHotline(ctx context.Context, args *etelecom.GetHotlineArgs) (*etelecom.Hotline, error) {
	return q.hotlineStore(ctx).ID(args.ID).OptionalOwnerID(args.OwnerID).GetHotline()
}

func (q *QueryService) ListHotlines(ctx context.Context, args *etelecom.ListHotlinesArgs) ([]*etelecom.Hotline, error) {
	return q.hotlineStore(ctx).OptionalOwnerID(args.OwnerID).OptionalConnectionID(args.ConnectionID).ListHotlines()
}
