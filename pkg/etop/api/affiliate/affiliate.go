package affiliate

import (
	"context"

	"etop.vn/api/main/identity"
	pbcm "etop.vn/backend/pb/common"
	pbaffiliate "etop.vn/backend/pb/etop/affiliate"
	"etop.vn/backend/pkg/common/bus"
	wrapaffiliate "etop.vn/backend/wrapper/etop/affiliate"
)

func init() {
	bus.AddHandlers("api",
		s.VersionInfo,
		s.RegisterAffiliate,
		s.UpdateAffiliate,
		s.UpdateAffiliateBankAccount,
		s.DeleteAffiliate,
	)
}

var (
	identityAggr identity.CommandBus

	s = &Service{}
)

func Init(identityA identity.CommandBus) {
	identityAggr = identityA
}

type Service struct{}

func (s *Service) VersionInfo(ctx context.Context, q *wrapaffiliate.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.affiliate",
		Version: "0.1",
	}
	return nil
}

func (s *Service) RegisterAffiliate(ctx context.Context, r *wrapaffiliate.RegisterAffiliateEndpoint) error {
	cmd := &identity.CreateAffiliateCommand{
		Name:        r.Name,
		OwnerID:     r.Context.UserID,
		Phone:       r.Phone,
		Email:       r.Email,
		BankAccount: r.BankAccount.ToCoreBankAccount(),
		IsTest:      r.Context.User.IsTest != 0,
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbaffiliate.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *Service) UpdateAffiliate(ctx context.Context, r *wrapaffiliate.UpdateAffiliateEndpoint) error {
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
	r.Result = pbaffiliate.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *Service) UpdateAffiliateBankAccount(ctx context.Context, r *wrapaffiliate.UpdateAffiliateBankAccountEndpoint) error {
	cmd := &identity.UpdateAffiliateBankAccountCommand{
		ID:          r.Context.Affiliate.ID,
		OwnerID:     r.Context.Affiliate.OwnerID,
		BankAccount: r.BankAccount.ToCoreBankAccount(),
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbaffiliate.Convert_core_Affiliate_To_api_Affiliate(cmd.Result)
	return nil
}

func (s *Service) DeleteAffiliate(ctx context.Context, r *wrapaffiliate.DeleteAffiliateEndpoint) error {
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
