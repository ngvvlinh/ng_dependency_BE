package query

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/etelecom/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ etelecom.QueryService = &QueryService{}

type QueryService struct {
	db             *cmsql.Database
	hotlineStore   sqlstore.HotlineStoreFactory
	extensionStore sqlstore.ExtensionStoreFactory
	callLogStore   sqlstore.CallLogStoreFactory
	tenantStore    sqlstore.TenantStoreFactory
	connectionQS   connectioning.QueryBus
}

func NewQueryService(dbEtelecom com.EtelecomDB, connectionQ connectioning.QueryBus) *QueryService {
	return &QueryService{
		db:             dbEtelecom,
		hotlineStore:   sqlstore.NewHotlineStore(dbEtelecom),
		extensionStore: sqlstore.NewExtensionStore(dbEtelecom),
		callLogStore:   sqlstore.NewCallLogStore(dbEtelecom),
		tenantStore:    sqlstore.NewTenantStore(dbEtelecom),
		connectionQS:   connectionQ,
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
	return q.hotlineStore(ctx).OptionalOwnerID(args.OwnerID).OptionalConnectionID(args.ConnectionID).OptionalTenantID(args.TenantID).ListHotlines()
}

func (q *QueryService) ListBuiltinHotlines(ctx context.Context, _ *cm.Empty) (res []*etelecom.Hotline, _ error) {
	queryConn := &connectioning.ListConnectionsQuery{
		Status:           status3.P.Wrap(),
		ConnectionType:   connection_type.Telecom,
		ConnectionMethod: connection_type.ConnectionMethodBuiltin,
	}
	if err := q.connectionQS.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	connIDs := []dot.ID{}
	for _, conn := range queryConn.Result {
		connIDs = append(connIDs, conn.ID)
	}
	if len(connIDs) == 0 {
		return
	}

	return q.hotlineStore(ctx).ConnectionIDs(connIDs...).OptionalStatus(status3.P).ListHotlines()
}
