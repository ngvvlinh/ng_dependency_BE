package affiliate

import (
	"context"

	"o.o/api/main/identity"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
)

type AccountService struct {
	IdentityAggr identity.CommandBus
}

func (s *AccountService) Clone() *AccountService {
	res := *s
	return &res
}

func (s *AccountService) RegisterAffiliate(ctx context.Context, r *RegisterAffiliateEndpoint) error {
	cmd := &identity.CreateAffiliateCommand{
		Name:        r.Name,
		OwnerID:     r.Context.UserID,
		Phone:       r.Phone,
		Email:       r.Email,
		BankAccount: convertpb.BankAccountToCoreBankAccount(r.BankAccount),
		IsTest:      r.Context.User.IsTest != 0,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *AccountService) UpdateAffiliate(ctx context.Context, r *UpdateAffiliateEndpoint) error {
	affiliate := r.Context.Affiliate
	cmd := &identity.UpdateAffiliateInfoCommand{
		ID:      affiliate.ID,
		OwnerID: affiliate.OwnerID,
		Phone:   r.Phone,
		Email:   r.Email,
		Name:    r.Name,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *AccountService) UpdateAffiliateBankAccount(ctx context.Context, r *UpdateAffiliateBankAccountEndpoint) error {
	cmd := &identity.UpdateAffiliateBankAccountCommand{
		ID:          r.Context.Affiliate.ID,
		OwnerID:     r.Context.Affiliate.OwnerID,
		BankAccount: convertpb.BankAccountToCoreBankAccount(r.BankAccount),
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *AccountService) DeleteAffiliate(ctx context.Context, r *DeleteAffiliateEndpoint) error {
	cmd := &identity.DeleteAffiliateCommand{
		ID:      r.Id,
		OwnerID: r.Context.Affiliate.OwnerID,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.Empty{}
	return nil
}
