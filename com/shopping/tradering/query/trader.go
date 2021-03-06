package query

import (
	"context"

	"o.o/api/shopping"
	"o.o/api/shopping/carrying"
	"o.o/api/shopping/customering"
	"o.o/api/shopping/suppliering"
	"o.o/api/shopping/tradering"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/tradering/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ tradering.QueryService = &TraderQuery{}

type TraderQuery struct {
	store         sqlstore.TraderStoreFactory
	customerQuery customering.QueryBus
	carrierQuery  carrying.QueryBus
	supplierQuery suppliering.QueryBus
}

func NewTraderQuery(
	db com.MainDB, customerQ customering.QueryBus,
	carrierQ carrying.QueryBus, supplierQ suppliering.QueryBus,
) *TraderQuery {
	return &TraderQuery{
		store:         sqlstore.NewTraderStore(db),
		customerQuery: customerQ,
		carrierQuery:  carrierQ,
		supplierQuery: supplierQ,
	}
}

func TraderQueryMessageBus(q *TraderQuery) tradering.QueryBus {
	b := bus.New()
	return tradering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *TraderQuery) GetTraderByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*tradering.ShopTrader, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetTrader()
}

func (q *TraderQuery) GetTraderInfoByID(
	ctx context.Context, id, shopID dot.ID,
) (*tradering.ShopTrader, error) {
	var fullName, phone string
	trader, err := q.store(ctx).IncludeDeleted().ID(id).GetTrader()
	if err != nil {
		return nil, err
	}
	var _shopID dot.ID
	switch trader.Type {
	case tradering.CustomerType:
		query := &customering.GetCustomerByIDQuery{
			ID:             id,
			IncludeDeleted: true,
		}
		if err := q.customerQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		_shopID = query.Result.ShopID
		fullName = query.Result.FullName
		phone = query.Result.Phone

	case tradering.CarrierType:
		query := &carrying.GetCarrierByIDQuery{
			ID: id,
		}
		if err := q.carrierQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		_shopID = query.Result.ShopID
		fullName = query.Result.FullName
	case tradering.SupplierType:
		query := &suppliering.GetSupplierByIDQuery{
			ID: id,
		}
		if err := q.supplierQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		_shopID = query.Result.ShopID
		fullName = query.Result.FullName
		phone = query.Result.Phone
	}
	if _shopID != 0 {
		if _shopID != shopID {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Đối tác không thuộc cửa hàng").WithMetap("traderID", id).WithMetap("shopID", shopID)
		}
	}

	traderResult := &tradering.ShopTrader{
		ID:       id,
		ShopID:   shopID,
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
