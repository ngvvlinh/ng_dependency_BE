package shop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbshop "etop.vn/backend/pb/etop/shop"
	wrapetop "etop.vn/backend/wrapper/etop"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandler("api", RegisterShop)
	bus.AddHandler("api", UpdateShop)
	bus.AddHandler("api", deleteShop)
	bus.AddHandler("api", SetDefaultAddress)
}

func RegisterShop(ctx context.Context, q *wrapshop.RegisterShopEndpoint) error {
	if q.UrlSlug != "" && !validate.URLSlug(q.UrlSlug) {
		return cm.Error(cm.InvalidArgument, "Thông tin url_slug không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}
	address, err := q.Address.ToModel()
	if err != nil {
		return err
	}
	cmd := &model.CreateShopCommand{
		Name:                        q.Name,
		OwnerID:                     q.Context.UserID,
		Phone:                       q.Phone,
		BankAccount:                 q.BankAccount.ToModel(),
		WebsiteURL:                  q.WebsiteUrl,
		ImageURL:                    q.ImageUrl,
		Email:                       q.Email,
		Address:                     address,
		AutoCreateFFM:               true,
		URLSlug:                     q.UrlSlug,
		IsTest:                      q.Context.User.IsTest != 0,
		CompanyInfo:                 q.CompanyInfo.ToModel(),
		MoneyTransactionRRule:       q.MoneyTransactionRrule,
		SurveyInfo:                  pbetop.SurveyInfosToModel(q.SurveyInfo),
		ShippingServicePickStrategy: pbetop.ShippingServiceSelectStrategyToModel(q.ShippingServiceSelectStrategy),
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.RegisterShopResponse{
		Shop: pbetop.PbShopExtended(cmd.Result),
	}
	return nil
}

func UpdateShop(ctx context.Context, q *wrapshop.UpdateShopEndpoint) error {
	shop := q.Context.Shop
	if q.BankAccount != nil {
		user, err := sqlstore.User(ctx).ID(shop.OwnerID).Get()
		if err != nil {
			return cm.Errorf(cm.Internal, err, "Không thể gửi mã xác nhận thay đổi tài khoản ngân hàng")
		}

		if !q.Context.Claim.SToken {
			stokenCmd := &wrapetop.SendSTokenEmailEndpoint{
				SendSTokenEmailRequest: &pbetop.SendSTokenEmailRequest{
					Email:     user.Email,
					AccountId: q.Context.Shop.ID,
				},
			}
			stokenCmd.Context.Claim = q.Context.Claim
			stokenCmd.Context.Admin = q.Context.Admin
			stokenCmd.Context.User = &model.SignedInUser{user} // TODO: remove this hack
			if err := bus.Dispatch(ctx, stokenCmd); err != nil {
				return err
			}

			return cm.Errorf(cm.STokenRequired, nil, "Cần được xác nhận trước khi thay đổi tài khoản ngân hàng. "+stokenCmd.Result.Msg)
		}
	}

	address, err := q.Address.ToModel()
	if err != nil {
		return err
	}
	cmd := &model.UpdateShopCommand{
		Shop: &model.Shop{
			ID:                            q.Context.Shop.ID,
			Name:                          q.Name,
			Phone:                         q.Phone,
			BankAccount:                   q.BankAccount.ToModel(),
			WebsiteURL:                    q.WebsiteUrl,
			ImageURL:                      q.ImageUrl,
			Email:                         q.Email,
			Address:                       address,
			TryOn:                         q.TryOn.ToModel(),
			GhnNoteCode:                   q.GhnNoteCode.ToModel(),
			CompanyInfo:                   q.CompanyInfo.ToModel(),
			MoneyTransactionRRule:         q.MoneyTransactionRrule,
			SurveyInfo:                    pbetop.SurveyInfosToModel(q.SurveyInfo),
			ShippingServiceSelectStrategy: pbetop.ShippingServiceSelectStrategyToModel(q.ShippingServiceSelectStrategy),
		},
		AutoCreateFFM: q.AutoCreateFfm,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateShopResponse{
		Shop: pbetop.PbShopExtended(cmd.Result),
	}
	return nil
}

func deleteShop(ctx context.Context, q *wrapshop.DeleteShopEndpoint) error {
	cmd := &model.DeleteShopCommand{
		ID:      q.Id,
		OwnerID: q.Context.UserID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.Empty{}
	return nil
}

func SetDefaultAddress(ctx context.Context, q *wrapshop.SetDefaultAddressEndpoint) error {
	cmd := &model.SetDefaultAddressShopCommand{
		ShopID:    q.Context.Shop.ID,
		Type:      q.Type.ToModel(),
		AddressID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}

	return nil
}
