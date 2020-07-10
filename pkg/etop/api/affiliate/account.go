package affiliate

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/affiliate"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type AccountService struct {
	session.Session

	IdentityAggr identity.CommandBus
}

func (s *AccountService) Clone() affiliate.AccountService {
	res := *s
	return &res
}

func (s *AccountService) RegisterAffiliate(ctx context.Context, r *affiliate.RegisterAffiliateRequest) (*etop.Affiliate, error) {
	cmd := &identity.CreateAffiliateCommand{
		Name:        r.Name,
		OwnerID:     s.SS.Claim().UserID,
		Phone:       r.Phone,
		Email:       r.Email,
		BankAccount: convertpb.BankAccountToCoreBankAccount(r.BankAccount),
		IsTest:      s.SS.User().IsTest != 0,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return result, nil
}

func (s *AccountService) UpdateAffiliate(ctx context.Context, r *affiliate.UpdateAffiliateRequest) (*etop.Affiliate, error) {
	aff := s.SS.Affiliate()
	cmd := &identity.UpdateAffiliateInfoCommand{
		ID:      aff.ID,
		OwnerID: aff.OwnerID,
		Phone:   r.Phone,
		Email:   r.Email,
		Name:    r.Name,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return result, nil
}

func (s *AccountService) UpdateAffiliateBankAccount(ctx context.Context, r *affiliate.UpdateAffiliateBankAccountRequest) (*etop.Affiliate, error) {
	cmd := &identity.UpdateAffiliateBankAccountCommand{
		ID:          s.SS.Affiliate().ID,
		OwnerID:     s.SS.Affiliate().OwnerID,
		BankAccount: convertpb.BankAccountToCoreBankAccount(r.BankAccount),
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return result, nil
}

func (s *AccountService) DeleteAffiliate(ctx context.Context, r *pbcm.IDRequest) (*pbcm.Empty, error) {
	cmd := &identity.DeleteAffiliateCommand{
		ID:      r.Id,
		OwnerID: s.SS.Affiliate().OwnerID,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.Empty{}
	return result, nil
}
