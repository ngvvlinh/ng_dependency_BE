package api

import (
	"context"

	"o.o/api/services/affiliate"
	api "o.o/api/top/services/affiliate"
	"o.o/backend/pkg/etop/authorize/session"
)

type UserService struct {
	session.Session

	AffiliateAggr affiliate.CommandBus
}

func (s *UserService) Clone() api.UserService { res := *s; return &res }

func (s *UserService) UpdateReferral(ctx context.Context, q *api.UpdateReferralRequest) (*api.UserReferral, error) {
	cmd := &affiliate.CreateOrUpdateUserReferralCommand{
		UserID:           s.SS.Claim().UserID,
		ReferralCode:     q.ReferralCode,
		SaleReferralCode: q.SaleReferralCode,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &api.UserReferral{
		UserId:           cmd.Result.UserID,
		ReferralCode:     cmd.Result.ReferralCode,
		SaleReferralCode: cmd.Result.SaleReferralCode,
	}
	return result, nil
}
