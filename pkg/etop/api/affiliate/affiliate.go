package affiliate

import (
	"context"

	"etop.vn/api/main/identity"
	pbcm "etop.vn/api/pb/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

func init() {
	bus.AddHandlers("api",
		miscService.VersionInfo,
		accountService.RegisterAffiliate,
		accountService.UpdateAffiliate,
		accountService.UpdateAffiliateBankAccount,
		accountService.DeleteAffiliate,
	)
}

var (
	identityAggr identity.CommandBus
)

func Init(identityA identity.CommandBus) {
	identityAggr = identityA
}

type MiscService struct{}
type AccountService struct{}

var miscService = &MiscService{}
var accountService = &AccountService{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.affiliate",
		Version: "0.1",
	}
	return nil
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
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.Empty{}
	return nil
}
