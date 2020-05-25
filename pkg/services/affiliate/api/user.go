package api

import (
	"context"

	"o.o/api/services/affiliate"
	apiaffiliate "o.o/api/top/services/affiliate"
)

type UserService struct {
	AffiliateAggr affiliate.CommandBus
}

func (s *UserService) Clone() *UserService { res := *s; return &res }

func (s *UserService) UpdateReferral(ctx context.Context, q *UpdateReferralEndpoint) error {
	cmd := &affiliate.CreateOrUpdateUserReferralCommand{
		UserID:           q.Context.UserID,
		ReferralCode:     q.ReferralCode,
		SaleReferralCode: q.SaleReferralCode,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &apiaffiliate.UserReferral{
		UserId:           cmd.Result.UserID,
		ReferralCode:     cmd.Result.ReferralCode,
		SaleReferralCode: cmd.Result.SaleReferralCode,
	}
	return nil
}
