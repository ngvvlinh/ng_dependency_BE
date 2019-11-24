package query

import (
	"context"

	"etop.vn/api/shopping/addressing"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
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

func (a *AddressQuery) GetAddressByID(ctx context.Context, ID dot.ID, ShopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := a.store(ctx).ID(ID).ShopID(ShopID).GetAddress()
	return addr, err
}

func (a *AddressQuery) GetAddressByTraderID(ctx context.Context, traderID, shopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := a.store(ctx).ShopTraderID(shopID, traderID).GetAddress()
	return addr, err
}

func (a *AddressQuery) GetAddressActiveByTraderID(ctx context.Context, traderID, ShopID dot.ID) (*addressing.ShopTraderAddress, error) {
	addr, err := a.store(ctx).ShopTraderID(ShopID, traderID).IsDefault(true).GetAddress()
	return addr, err
}

func (a *AddressQuery) ListAddressesByTraderID(ctx context.Context, ShopID dot.ID, TraderID dot.ID) ([]*addressing.ShopTraderAddress, error) {
	addrs, err := a.store(ctx).ShopTraderID(ShopID, TraderID).ListAddresses()
	return addrs, err
}
