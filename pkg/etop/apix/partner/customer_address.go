package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type CustomerAddressService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *CustomerAddressService) Clone() api.CustomerAddressService { res := *s; return &res }

func (s *CustomerAddressService) ListAddresses(ctx context.Context, r *externaltypes.ListCustomerAddressesRequest) (*externaltypes.CustomerAddressesResponse, error) {
	resp, err := s.Shopping.ListAddresses(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerAddressService) GetAddress(ctx context.Context, r *externaltypes.OrderIDRequest) (*externaltypes.CustomerAddress, error) {
	resp, err := s.Shopping.GetAddress(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerAddressService) CreateAddress(ctx context.Context, r *externaltypes.CreateCustomerAddressRequest) (*externaltypes.CustomerAddress, error) {
	resp, err := s.Shopping.CreateAddress(ctx, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, r)
	return resp, err
}

func (s *CustomerAddressService) UpdateAddress(ctx context.Context, r *externaltypes.UpdateCustomerAddressRequest) (*externaltypes.CustomerAddress, error) {
	resp, err := s.Shopping.UpdateAddress(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerAddressService) DeleteAddress(ctx context.Context, r *pbcm.IDRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteAddress(ctx, s.SS.Shop().ID, r)
	return resp, err
}
