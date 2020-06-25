package shop

import (
	"context"

	"github.com/asaskevich/govalidator"

	"o.o/api/main/address"
	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	apietop "o.o/api/top/int/etop"
	"o.o/api/top/int/shop"
	apishop "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	etop "o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/tools/pkg/acl"
)

type AccountService struct {
	IdentityAggr  identity.CommandBus
	IdentityQuery identity.QueryBus
	AddressQuery  address.QueryBus
}

func (s *AccountService) Clone() *AccountService { res := *s; return &res }

func (s *AccountService) RegisterShop(ctx context.Context, q *RegisterShopEndpoint) error {
	if q.UrlSlug != "" && !validate.URLSlug(q.UrlSlug) {
		return cm.Error(cm.InvalidArgument, "Thông tin url_slug không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}
	addr, err := convertpb.AddressToModel(q.Address)
	if err != nil {
		return err
	}
	cmd := &identitymodelx.CreateShopCommand{
		Name:                        q.Name,
		OwnerID:                     q.Context.UserID,
		Phone:                       q.Phone,
		BankAccount:                 convertpb.BankAccountToModel(q.BankAccount),
		WebsiteURL:                  q.WebsiteUrl,
		ImageURL:                    q.ImageUrl,
		Email:                       q.Email,
		Address:                     addr,
		AutoCreateFFM:               true,
		URLSlug:                     q.UrlSlug,
		IsTest:                      q.Context.User.IsTest != 0,
		CompanyInfo:                 convertpb.CompanyInfoToModel(q.CompanyInfo),
		MoneyTransactionRRule:       q.MoneyTransactionRrule,
		SurveyInfo:                  convertpb.SurveyInfosToModel(q.SurveyInfo),
		ShippingServicePickStrategy: convertpb.ShippingServiceSelectStrategyToModel(q.ShippingServiceSelectStrategy),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &apishop.RegisterShopResponse{
		Shop: convertpb.PbShopExtended(cmd.Result),
	}
	return nil
}

func (s *AccountService) UpdateShop(ctx context.Context, q *UpdateShopEndpoint) error {
	shop := q.Context.Shop
	if q.BankAccount != nil {
		user, err := sqlstore.User(ctx).ID(shop.OwnerID).Get()
		if err != nil {
			return cm.Errorf(cm.Internal, err, "Không thể gửi mã xác nhận thay đổi tài khoản ngân hàng")
		}

		if !q.Context.Claim.SToken {
			stokenCmd := &etop.SendSTokenEmailEndpoint{
				SendSTokenEmailRequest: &apietop.SendSTokenEmailRequest{
					Email:     user.Email,
					AccountId: q.Context.Shop.ID,
				},
			}
			stokenCmd.Context.Claim = q.Context.Claim
			stokenCmd.Context.Admin = q.Context.Admin
			stokenCmd.Context.User = &identitymodelx.SignedInUser{user} // TODO: remove this hack
			if err := etop.UserServiceImpl.SendSTokenEmail(ctx, stokenCmd); err != nil {
				return err
			}

			return cm.Errorf(cm.STokenRequired, nil, "Cần được xác nhận trước khi thay đổi tài khoản ngân hàng. "+stokenCmd.Result.Msg)
		}
	}

	address, err := convertpb.AddressToModel(q.Address)
	if err != nil {
		return err
	}
	cmd := &identitymodelx.UpdateShopCommand{
		Shop: &identitymodel.Shop{
			ID:                            q.Context.Shop.ID,
			InventoryOverstock:            q.InventoryOverstock,
			Name:                          q.Name,
			Phone:                         q.Phone,
			BankAccount:                   convertpb.BankAccountToModel(q.BankAccount),
			WebsiteURL:                    q.WebsiteUrl,
			ImageURL:                      q.ImageUrl,
			Email:                         q.Email,
			Address:                       address,
			TryOn:                         q.TryOn.Apply(0),
			GhnNoteCode:                   q.GhnNoteCode.Apply(0),
			CompanyInfo:                   convertpb.CompanyInfoToModel(q.CompanyInfo),
			MoneyTransactionRRule:         q.MoneyTransactionRrule,
			SurveyInfo:                    convertpb.SurveyInfosToModel(q.SurveyInfo),
			ShippingServiceSelectStrategy: convertpb.ShippingServiceSelectStrategyToModel(q.ShippingServiceSelectStrategy),
		},
		AutoCreateFFM: q.AutoCreateFfm,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &apishop.UpdateShopResponse{
		Shop: convertpb.PbShopExtended(cmd.Result),
	}
	return nil
}

func (s *AccountService) DeleteShop(ctx context.Context, q *DeleteShopEndpoint) error {
	cmd := &identitymodelx.DeleteShopCommand{
		ID:      q.Id,
		OwnerID: q.Context.UserID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.Empty{}
	return nil
}

func (s *AccountService) SetDefaultAddress(ctx context.Context, q *SetDefaultAddressEndpoint) error {
	cmd := &identitymodelx.SetDefaultAddressShopCommand{
		ShopID:    q.Context.Shop.ID,
		Type:      q.Type.String(),
		AddressID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}

	return nil
}

func (s *AccountService) CreateExternalAccountAhamove(ctx context.Context, q *CreateExternalAccountAhamoveEndpoint) error {
	query := &identity.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	queryAddress := &address.GetAddressByIDQuery{
		ID: q.Context.Shop.AddressID,
	}
	if err := s.AddressQuery.Dispatch(ctx, queryAddress); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Thiếu thông tin địa chỉ cửa hàng")
	}
	addr := queryAddress.Result
	cmd := &identity.CreateExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
		Name:    user.FullName,
		Address: addr.GetFullAddress(),
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(cmd.Result, false)
	return nil
}

func (s *AccountService) GetExternalAccountAhamove(ctx context.Context, q *GetExternalAccountAhamoveEndpoint) error {
	queryUser := &identity.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryUser); err != nil {
		return err
	}
	user := queryUser.Result
	phone := user.Phone

	query := &identity.GetExternalAccountAhamoveQuery{
		Phone:   phone,
		OwnerID: user.ID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	account := query.Result
	if !account.ExternalVerified && account.ExternalTicketID != "" {
		cmd := &identity.UpdateVerifiedExternalAccountAhamoveCommand{
			OwnerID: user.ID,
			Phone:   phone,
		}
		if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
		account = cmd.Result
	}

	var hideInfo bool
	if !authorization.IsContainsActionString(auth.ListActionsByRoles(q.Context.Roles), string(acl.ShopExternalAccountManage)) {
		hideInfo = true
	}
	q.Result = convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(account, hideInfo)
	return nil
}

func (s *AccountService) RequestVerifyExternalAccountAhamove(ctx context.Context, q *RequestVerifyExternalAccountAhamoveEndpoint) error {
	query := &identitymodelx.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.RequestVerifyExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *AccountService) UpdateExternalAccountAhamoveVerification(ctx context.Context, r *UpdateExternalAccountAhamoveVerificationEndpoint) error {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return err
	}

	query := &identitymodelx.GetUserByIDQuery{
		UserID: r.Context.Shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.UpdateExternalAccountAhamoveVerificationCommand{
		OwnerID:             user.ID,
		Phone:               phone,
		IDCardFrontImg:      r.IdCardFrontImg,
		IDCardBackImg:       r.IdCardBackImg,
		PortraitImg:         r.PortraitImg,
		WebsiteURL:          r.WebsiteUrl,
		FanpageURL:          r.FanpageUrl,
		CompanyImgs:         r.CompanyImgs,
		BusinessLicenseImgs: r.BusinessLicenseImgs,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *AccountService) GetBalanceShop(ctx context.Context, q *GetBalanceShopEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &model.GetBalanceShopCommand{
		ShopID: shopID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.GetBalanceShopResponse{
		Amount: cmd.Result.Amount,
	}
	return nil
}

func validateUrl(imgsUrl ...string) error {
	for _, url := range imgsUrl {
		if url == "" {
			continue
		}
		if !govalidator.IsURL(url) {
			return cm.Errorf(cm.InvalidArgument, nil, "Invalid url: %v", url)
		}
	}
	return nil
}
