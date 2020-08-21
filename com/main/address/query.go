package address

import (
	"context"

	"o.o/api/main/address"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/address/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ address.QueryService = &AddressQueryService{}

type AddressQueryService struct {
	AddressStore sqlstore.AddressFactory
	eventBus     capi.EventBus
}

func NewQueryAddress(db com.MainDB, bus capi.EventBus) *AddressQueryService {
	return &AddressQueryService{
		AddressStore: sqlstore.NewAddressStore(db),
		eventBus:     bus,
	}
}

func QueryServiceMessageBus(q *AddressQueryService) address.QueryBus {
	b := bus.New()
	return address.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (s *AddressQueryService) GetAddressByID(ctx context.Context, q *address.GetAddressByIDQueryArgs) (*address.Address, error) {
	query := s.AddressStore(ctx).ID(q.ID)
	addr, err := query.Get()
	if err != nil {
		return nil, err
	}
	return addr, err
}

func (s *AddressQueryService) ListAddresses(ctx context.Context, AccountID dot.ID) (*address.GetAddressResponse, error) {
	query := s.AddressStore(ctx).AccountID(AccountID)
	addrs, err := query.ListAddresses()
	if err != nil {
		return nil, err
	}

	return &address.GetAddressResponse{
		Addresses: addrs,
	}, nil
}
