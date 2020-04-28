package shop

import (
	"context"

	apietop "o.o/api/top/int/etop"
	apishop "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	etop "o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/sqlstore"
)

func init() {
	bus.AddHandler("api", accountService.RegisterShop)
	bus.AddHandler("api", accountService.UpdateShop)
	bus.AddHandler("api", accountService.DeleteShop)
	bus.AddHandler("api", accountService.SetDefaultAddress)
}

func (s *AccountService) RegisterShop(ctx context.Context, q *RegisterShopEndpoint) error {
	if q.UrlSlug != "" && !validate.URLSlug(q.UrlSlug) {
		return cm.Error(cm.InvalidArgument, "Thông tin url_slug không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}
	address, err := convertpb.AddressToModel(q.Address)
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
		Address:                     address,
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
			if err := bus.Dispatch(ctx, stokenCmd); err != nil {
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
