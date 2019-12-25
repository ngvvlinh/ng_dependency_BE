package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ tradering.QueryService = &TraderQuery{}

type TraderQuery struct {
	store         sqlstore.TraderStoreFactory
	customerQuery customering.QueryBus
	carrierQuery  carrying.QueryBus
	supplierQuery suppliering.QueryBus
}

func NewTraderQuery(
	db *cmsql.Database, customerQ customering.QueryBus,
	carrierQ carrying.QueryBus, supplierQ suppliering.QueryBus,
) *TraderQuery {
	return &TraderQuery{
		store:         sqlstore.NewTraderStore(db),
		customerQuery: customerQ,
		carrierQuery:  carrierQ,
		supplierQuery: supplierQ,
	}
}

func (q *TraderQuery) MessageBus() tradering.QueryBus {
	b := bus.New()
	return tradering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *TraderQuery) GetTraderByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*tradering.ShopTrader, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetTrader()
}

func (q *TraderQuery) GetTraderInfoByID(
	ctx context.Context, ID, ShopID dot.ID,
) (*tradering.ShopTrader, error) {
	var fullName, phone string
	trader, err := q.store(ctx).ID(ID).ShopID(ShopID).GetTrader()
	if err != nil {
		return nil, err
	}
	switch trader.Type {
	case tradering.CustomerType:
		query := &customering.GetCustomerByIDQuery{
			ID:     ID,
			ShopID: ShopID,
		}
		if err := q.customerQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		fullName = query.Result.FullName
		phone = query.Result.Phone

	case tradering.CarrierType:
		query := &carrying.GetCarrierByIDQuery{
			ID:     ID,
			ShopID: ShopID,
		}
		if err := q.carrierQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		fullName = query.Result.FullName
	case tradering.SupplierType:
		query := &suppliering.GetSupplierByIDQuery{
			ID:     ID,
			ShopID: ShopID,
		}
		if err := q.supplierQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		fullName = query.Result.FullName
		phone = query.Result.Phone
	}

	traderResult := &tradering.ShopTrader{
		ID:       ID,
		ShopID:   ShopID,
		Type:     trader.Type,
		FullName: fullName,
		Phone:    phone,
	}
	return traderResult, nil
}

func (q *TraderQuery) ListTradersByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*tradering.TradersResponse, error) {
	traders, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListTraders()
	if err != nil {
		return nil, err
	}
	return &tradering.TradersResponse{Traders: traders}, nil
}
