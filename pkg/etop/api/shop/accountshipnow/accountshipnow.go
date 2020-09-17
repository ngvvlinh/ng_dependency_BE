package accountshipnow

import (
	"context"

	"o.o/api/main/accountshipnow"
	api "o.o/api/top/int/shop"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type AccountShipnowService struct {
	session.Session

	AccountshipnowQuery accountshipnow.QueryBus
}

func (s *AccountShipnowService) Clone() api.AccountShipnowService {
	res := *s
	return &res
}

func (s *AccountShipnowService) GetAccountShipnow(ctx context.Context, r *api.GetAccountShipnowRequest) (*api.ExternalAccountAhamove, error) {
	query := &accountshipnow.GetAccountShipnowQuery{
		Phone:        r.Identity,
		OwnerID:      s.SS.Shop().OwnerID,
		ConnectionID: r.ConnectionID,
	}
	if err := s.AccountshipnowQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(query.Result, false)
	return result, nil
}
