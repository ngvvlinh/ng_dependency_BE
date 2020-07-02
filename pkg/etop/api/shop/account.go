package shop

import (
	"context"

	"github.com/asaskevich/govalidator"

	"o.o/api/main/address"
	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	apietop "o.o/api/top/int/etop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	etop "o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/tools/pkg/acl"
)

type AccountService struct {
	session.Session

	IdentityAggr  identity.CommandBus
	IdentityQuery identity.QueryBus
	AddressQuery  address.QueryBus
}

func (s *AccountService) Clone() api.AccountService { res := *s; return &res }

func (s *AccountService) RegisterShop(ctx context.Context, q *api.RegisterShopRequest) (*api.RegisterShopResponse, error) {
	if q.UrlSlug != "" && !validate.URLSlug(q.UrlSlug) {
		return nil, cm.Error(cm.InvalidArgument, "Thông tin url_slug không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}
	addr, err := convertpb.AddressToModel(q.Address)
	if err != nil {
		return nil, err
	}
	cmd := &identitymodelx.CreateShopCommand{
		Name:                        q.Name,
		OwnerID:                     s.SS.Claim().UserID,
		Phone:                       q.Phone,
		BankAccount:                 convertpb.BankAccountToModel(q.BankAccount),
		WebsiteURL:                  q.WebsiteUrl,
		ImageURL:                    q.ImageUrl,
		Email:                       q.Email,
		Address:                     addr,
		AutoCreateFFM:               true,
		URLSlug:                     q.UrlSlug,
		IsTest:                      s.SS.User().IsTest != 0,
		CompanyInfo:                 convertpb.CompanyInfoToModel(q.CompanyInfo),
		MoneyTransactionRRule:       q.MoneyTransactionRrule,
		SurveyInfo:                  convertpb.SurveyInfosToModel(q.SurveyInfo),
		ShippingServicePickStrategy: convertpb.ShippingServiceSelectStrategyToModel(q.ShippingServiceSelectStrategy),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &api.RegisterShopResponse{
		Shop: convertpb.PbShopExtended(cmd.Result),
	}
	return result, nil
}

func (s *AccountService) UpdateShop(ctx context.Context, q *api.UpdateShopRequest) (*api.UpdateShopResponse, error) {
	shop := s.SS.Shop()
	if q.BankAccount != nil {
		user, err := sqlstore.User(ctx).ID(shop.OwnerID).Get()
		if err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Không thể gửi mã xác nhận thay đổi tài khoản ngân hàng")
		}

		if !s.SS.Claim().SToken {
			req := &apietop.SendSTokenEmailRequest{
				Email:     user.Email,
				AccountId: s.SS.Shop().ID,
			}
			userService := etop.UserServiceImpl.Clone().(*etop.UserService)
			userService.Session = s.Session
			result, err := userService.SendSTokenEmail(ctx, req)
			if err != nil {
				return nil, err
			}
			return nil, cm.Errorf(cm.STokenRequired, nil, "Cần được xác nhận trước khi thay đổi tài khoản ngân hàng. "+result.Msg)
		}
	}

	address, err := convertpb.AddressToModel(q.Address)
	if err != nil {
		return nil, err
	}
	cmd := &identitymodelx.UpdateShopCommand{
		Shop: &identitymodel.Shop{
			ID:                            s.SS.Shop().ID,
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
		return nil, err
	}
	result := &api.UpdateShopResponse{
		Shop: convertpb.PbShopExtended(cmd.Result),
	}
	return result, nil
}

func (s *AccountService) DeleteShop(ctx context.Context, q *pbcm.IDRequest) (*pbcm.Empty, error) {
	cmd := &identitymodelx.DeleteShopCommand{
		ID:      q.Id,
		OwnerID: s.SS.Claim().UserID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.Empty{}
	return result, nil
}

func (s *AccountService) SetDefaultAddress(ctx context.Context, q *apietop.SetDefaultAddressRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &identitymodelx.SetDefaultAddressShopCommand{
		ShopID:    s.SS.Shop().ID,
		Type:      q.Type.String(),
		AddressID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}

	return result, nil
}

func (s *AccountService) CreateExternalAccountAhamove(ctx context.Context, q *pbcm.Empty) (*api.ExternalAccountAhamove, error) {
	query := &identity.GetUserByIDQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	user := query.Result
	phone := user.Phone

	queryAddress := &address.GetAddressByIDQuery{
		ID: s.SS.Shop().AddressID,
	}
	if err := s.AddressQuery.Dispatch(ctx, queryAddress); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Thiếu thông tin địa chỉ cửa hàng")
	}
	addr := queryAddress.Result
	cmd := &identity.CreateExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
		Name:    user.FullName,
		Address: addr.GetFullAddress(),
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(cmd.Result, false)
	return result, nil
}

func (s *AccountService) GetExternalAccountAhamove(ctx context.Context, q *pbcm.Empty) (*api.ExternalAccountAhamove, error) {
	queryUser := &identity.GetUserByIDQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryUser); err != nil {
		return nil, err
	}
	user := queryUser.Result
	phone := user.Phone

	query := &identity.GetExternalAccountAhamoveQuery{
		Phone:   phone,
		OwnerID: user.ID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	account := query.Result
	if !account.ExternalVerified && account.ExternalTicketID != "" {
		cmd := &identity.UpdateVerifiedExternalAccountAhamoveCommand{
			OwnerID: user.ID,
			Phone:   phone,
		}
		if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
		account = cmd.Result
	}

	var hideInfo bool
	if !authorization.IsContainsActionString(auth.ListActionsByRoles(s.SS.Permission().Roles), string(acl.ShopExternalAccountManage)) {
		hideInfo = true
	}
	result := convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(account, hideInfo)
	return result, nil
}

func (s *AccountService) RequestVerifyExternalAccountAhamove(ctx context.Context, q *pbcm.Empty) (*pbcm.UpdatedResponse, error) {
	query := &identitymodelx.GetUserByIDQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.RequestVerifyExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *AccountService) UpdateExternalAccountAhamoveVerification(ctx context.Context, r *api.UpdateXAccountAhamoveVerificationRequest) (*pbcm.UpdatedResponse, error) {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return nil, err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return nil, err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return nil, err
	}

	query := &identitymodelx.GetUserByIDQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
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
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *AccountService) GetBalanceShop(ctx context.Context, q *pbcm.Empty) (*api.GetBalanceShopResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &model.GetBalanceShopCommand{
		ShopID: shopID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.GetBalanceShopResponse{
		Amount: cmd.Result.Amount,
	}
	return result, nil
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
