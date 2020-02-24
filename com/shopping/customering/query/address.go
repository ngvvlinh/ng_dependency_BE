package query

import (
	"context"

	"etop.vn/api/shopping/addressing"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ addressing.QueryService = &AddressQuery{}

type AddressQuery struct {
	store sqlstore.AddressStoreFactory
}

func NewAddressQuery(db *cmsql.Database) *AddressQuery {
	return &AddressQuery{
		store: sqlstore.NewAddressStore(db),
	}
}

func (q *AddressQuery) MessageBus() addressing.QueryBus {
	b := bus.New()
	return addressing.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *AddressQuery) GetAddressByID(ctx context.Context, ID dot.ID, ShopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := q.store(ctx).ID(ID).ShopID(ShopID).GetAddress()
	return addr, err
}

func (q *AddressQuery) GetAddressByTraderID(ctx context.Context, traderID, shopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := q.store(ctx).ShopTraderID(shopID, traderID).GetAddress()
	return addr, err
}

func (q *AddressQuery) GetAddressActiveByTraderID(ctx context.Context, traderID, ShopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := q.store(ctx).ShopTraderID(ShopID, traderID).IsDefault(true).GetAddress()
	return addr, err
}

func (q *AddressQuery) ListAddressesByTraderID(ctx context.Context, args *addressing.ListAddressesByTraderIDArgs) (*addressing.ShopTraderAddressesResponse, error) {
	query := q.store(ctx).ShopTraderID(args.ShopID, args.TraderID)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	addrs, err := query.WithPaging(args.Paging).ListAddresses()
	if err != nil {
		return nil, err
	}
	return &addressing.ShopTraderAddressesResponse{
		ShopTraderAddresses: addrs,
		Count:               count,
		Paging:              query.GetPaging(),
	}, err
}
func (q *AddressQuery) ListAddressesByTraderIDs(ctx context.Context, args *addressing.ListAddressesByTraderIDsArgs) (*addressing.ShopTraderAddressesResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID)
	if len(args.TraderIDs) != 0 {
		query = query.IDs(args.TraderIDs...)
	}
	if args.IncludeDeleted {
		query = query.IncludeDeleted()
	}
	addrs, err := query.WithPaging(args.Paging).ListAddresses()
	if err != nil {
		return nil, err
	}
	return &addressing.ShopTraderAddressesResponse{
		ShopTraderAddresses: addrs,
		Paging:              query.GetPaging(),
	}, err
}
