package sqlstore

import (
	"context"
	"strings"
	"time"

	"etop.vn/api/main/authorization"
	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/try_on"
	addressmodel "etop.vn/backend/com/main/address/model"
	addressmodelx "etop.vn/backend/com/main/address/modelx"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		CreateShop,
		UpdateShop,
		DeleteShop,
		SetDefaultAddressShop,
		UpdateAccountURLSlug,
		GetAccountAuth,
	)
}

func CreateShop(ctx context.Context, cmd *identitymodelx.CreateShopCommand) error {
	if cmd.OwnerID == 0 {
		return cm.Error(cm.Internal, "Missing OwnerID", nil)
	}

	var ok bool
	var emailNorm model.NormalizedEmail
	var phoneNorm model.NormalizedPhone
	if cmd.Name, ok = validate.NormalizeName(cmd.Name); !ok {
		return cm.Error(cm.InvalidArgument, "Invalid name", nil)
	}
	if cmd.Email != "" {
		if emailNorm, ok = validate.NormalizeEmail(cmd.Email); !ok {
			return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
	}
	if phoneNorm, ok = validate.NormalizePhone(cmd.Phone); !ok {
		return cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}

	ownerQuery := &identity.GetUserByIDQuery{UserID: cmd.OwnerID}
	if err := bus.Dispatch(ctx, ownerQuery); err != nil {
		return cm.Error(cm.Internal, "invalid owner_id", nil)
	}

	id := model.NewShopID()
	return x.InTransaction(ctx, func(s cmsql.QueryInterface) error {
		account := &identitymodel.Account{
			ID:       id,
			Name:     cmd.Name,
			Type:     account_type.Shop,
			ImageURL: cmd.ImageURL,
			URLSlug:  cmd.URLSlug,
		}
		permission := &identitymodel.AccountUser{
			AccountID: id,
			UserID:    cmd.OwnerID,
			Status:    status3.P,
			Permission: identitymodel.Permission{
				Roles: []string{string(authorization.RoleShopOwner)},
			},
		}
		if _, err := x.Insert(account); err != nil {
			return err
		}
		if _, err := x.Insert(permission); err != nil {
			return err
		}
		if cmd.Address != nil {
			addressID, err := updateOrCreateAddress(ctx, x, cmd.Address, account.ID, model.AddressTypeGeneral)
			if err != nil {
				return err
			}
			cmd.AddressID = addressID
		}

		code, errCode := GenerateCode(ctx, x, model.CodeTypeShop, "")
		if errCode != nil {
			return errCode
		}

		shop := &identitymodel.Shop{
			ID:                            id,
			Name:                          cmd.Name,
			OwnerID:                       cmd.OwnerID,
			WLPartnerID:                   ownerQuery.Result.WLPartnerID,
			Status:                        status3.P,
			AddressID:                     cmd.AddressID,
			Phone:                         phoneNorm.String(),
			Email:                         emailNorm.String(),
			BankAccount:                   cmd.BankAccount,
			WebsiteURL:                    cmd.WebsiteURL,
			ImageURL:                      cmd.ImageURL,
			Code:                          code,
			TryOn:                         try_on.Open,
			CompanyInfo:                   cmd.CompanyInfo,
			MoneyTransactionRRule:         cmd.MoneyTransactionRRule,
			SurveyInfo:                    cmd.SurveyInfo,
			ShippingServiceSelectStrategy: cmd.ShippingServicePickStrategy,
			AutoCreateFFM:                 cmd.AutoCreateFFM,
		}
		if cmd.MoneyTransactionRRule == "" {
			// set shop MoneyTransactionRRule default value: FREQ=WEEKLY;BYDAY=MO,WE,FR
			shop.MoneyTransactionRRule = "FREQ=WEEKLY;BYDAY=MO,WE,FR"
		}
		if err := shop.CheckInfo(); err != nil {
			return err
		}
		if cmd.IsTest {
			shop.IsTest = 1
		}
		if _, err := x.Insert(shop); err != nil {
			return err
		}

		cmd.Result = new(identitymodel.ShopExtended)
		if has, err := x.
			Table("shop").
			Where("s.id = ?", shop.ID).
			Get(cmd.Result); err != nil || !has {
			return cm.Error(cm.Internal, "", err)
		}
		event := &identity.AccountCreatedEvent{
			ShopID: id,
			UserID: cmd.OwnerID,
		}
		if err := eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
}

func UpdateShop(ctx context.Context, cmd *identitymodelx.UpdateShopCommand) error {
	shop := cmd.Shop
	if shop.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	var ok bool
	var emailNorm model.NormalizedEmail
	var phoneNorm model.NormalizedPhone
	if shop.Name != "" {
		if shop.Name, ok = validate.NormalizeName(shop.Name); !ok {
			return cm.Error(cm.InvalidArgument, "Invalid name", nil)
		}
	}
	if shop.Email != "" {
		if emailNorm, ok = validate.NormalizeEmail(shop.Email); !ok {
			return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
		shop.Email = emailNorm.String()
	}
	if shop.Phone != "" {
		if phoneNorm, ok = validate.NormalizePhone(shop.Phone); !ok {
			return cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
		}
		shop.Phone = phoneNorm.String()
	}
	if err := shop.CheckInfo(); err != nil {
		return err
	}

	return inTransaction(func(x Qx) error {
		if shop.Address != nil {
			addressID, err := updateOrCreateAddress(ctx, x, shop.Address, shop.ID, model.AddressTypeGeneral)
			if err != nil {
				return err
			}
			shop.AddressID = addressID
		}

		if err := x.Table("shop").
			Where("id = ? AND deleted_at is NULL", shop.ID).
			ShouldUpdate(shop); err != nil {
			return err
		}
		updateMapValue := make(map[string]interface{})
		if shop.InventoryOverstock.Valid {
			updateMapValue["inventory_overstock"] = shop.InventoryOverstock
		}
		if cmd.AutoCreateFFM.Valid {
			updateMapValue["auto_create_ffm"] = cmd.AutoCreateFFM
		}
		if len(updateMapValue) != 0 {
			if err := x.Table("shop").Where("id= ?", shop.ID).ShouldUpdateMap(updateMapValue); err != nil {
				return err
			}
		}

		cmd.Result = new(identitymodel.ShopExtended)
		if has, err := x.
			Table("shop").
			Where("s.id = ?", shop.ID).
			Get(cmd.Result); err != nil || !has {
			return cm.Error(cm.Internal, "", err)
		}
		return nil
	})
}

func DeleteShop(ctx context.Context, cmd *identitymodelx.DeleteShopCommand) error {
	return inTransaction(func(s Qx) error {
		if cmd.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing ID", nil)
		}
		now := time.Now()
		{
			if updated, err := s.Table("shop").
				Where("id = ? AND owner_id = ?", cmd.ID, cmd.OwnerID).
				Update(&identitymodel.ShopDelete{
					DeletedAt: now,
				}); err != nil {
				return err
			} else if updated == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
		}
		if _, err := s.Table("account_user").
			Where("account_id = ? AND user_id = ?", cmd.ID, cmd.OwnerID).
			Update(
				&identitymodel.AccountUserDelete{
					DeletedAt: now,
				}); err != nil {
			return err
		}
		return nil
	})
}

func UpdateOrCreateAddress(ctx context.Context, address *addressmodel.Address, accountID dot.ID, AddressType string) (dot.ID, error) {
	return updateOrCreateAddress(ctx, x, address, accountID, AddressType)
}

func updateOrCreateAddress(ctx context.Context, x Qx, address *addressmodel.Address, accountID dot.ID, AddressType string) (dot.ID, error) {
	addressObj := &addressmodel.Address{
		Province:     address.Province,
		ProvinceCode: address.ProvinceCode,
		District:     address.District,
		DistrictCode: address.DistrictCode,
		Ward:         address.Ward,
		WardCode:     address.WardCode,
		Address1:     address.Address1,
		Address2:     address.Address2,
		FullName:     address.FullName,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		Phone:        address.Phone,
		Position:     address.Position,
		Email:        address.Email,
	}

	if address.ID != 0 {
		// update warehouse address
		addressObj.ID = address.ID
		addressCmd := &addressmodelx.UpdateAddressCommand{
			Address: addressObj,
		}

		if err := updateAddress(ctx, x, addressCmd); err != nil {
			return 0, err
		}
		return addressCmd.Result.ID, nil
	} else {
		// create new warehouse address
		addressObj.AccountID = accountID
		addressObj.Type = AddressType
		addressCmd := &addressmodelx.CreateAddressCommand{
			Address: addressObj,
		}
		if err := createAddress(ctx, x, addressCmd); err != nil {
			return 0, err
		}
		return addressCmd.Result.ID, nil
	}
}

func SetDefaultAddressShop(ctx context.Context, cmd *identitymodelx.SetDefaultAddressShopCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if cmd.Type == "" {
		return cm.Error(cm.InvalidArgument, "Missing Address Type", nil)
	}
	if cmd.AddressID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	var address = new(addressmodel.Address)
	if err := x.Table("address").Where("id = ? AND account_id = ? and type = ?", cmd.AddressID, cmd.ShopID, cmd.Type).
		ShouldGet(address); err != nil {
		return err
	}
	shopObj := &identitymodel.Shop{
		ID: cmd.ShopID,
	}

	switch cmd.Type {
	case model.AddressTypeShipTo:
		shopObj.ShipToAddressID = cmd.AddressID
	case model.AddressTypeShipFrom:
		shopObj.ShipFromAddressID = cmd.AddressID
	default:
		return cm.Error(cm.Unimplemented, "Address type does not valid", nil)
	}

	cmdUpdateShop := &identitymodelx.UpdateShopCommand{
		Shop: shopObj,
	}
	if err := bus.Dispatch(ctx, cmdUpdateShop); err != nil {
		return err
	}
	cmd.Result.Updated = 1
	return nil
}

func UpdateAccountURLSlug(ctx context.Context, cmd *identitymodelx.UpdateAccountURLSlugCommand) error {
	if cmd.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}

	s := strings.TrimSpace(cmd.URLSlug)
	if ok := validate.URLSlug(s); !ok {
		return cm.Error(cm.InvalidArgument, "Thông tin truyền vào không hợp lệ.", nil)
	}

	return x.Table("account").
		Where("id = ?", cmd.AccountID).
		ShouldUpdateMap(map[string]interface{}{
			"url_slug": cmd.URLSlug,
		})
}

func GetAccountAuth(ctx context.Context, query *identitymodelx.GetAccountAuthQuery) error {
	if query.AuthKey == "" || query.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing key")
	}

	s := x.Where("auth_key = ?", query.AuthKey).
		Where("account_id = ?", query.AccountID).
		Where("aa.status = 1 AND aa.deleted_at IS NULL")

	switch query.AccountType {
	case account_type.Partner:
		if cm.GetTag(query.AccountID) != model.TagPartner {
			return cm.Errorf(cm.NotFound, nil, "")
		}

		var res identitymodel.AccountAuthFtPartner
		if err := s.Where("p.status = 1 AND p.deleted_at IS NULL").
			ShouldGet(&res); err != nil {
			return err
		}
		query.Result.AccountAuth = res.AccountAuth
		query.Result.Account = res.Partner
		return nil

	case account_type.Shop:
		if cm.GetTag(query.AccountID) != model.TagShop {
			return cm.Errorf(cm.NotFound, nil, "")
		}

		var res identitymodel.AccountAuthFtShop
		if err := s.Where("s.status = 1 AND s.deleted_at IS NULL").
			ShouldGet(&res); err != nil {
			return err
		}
		query.Result.AccountAuth = res.AccountAuth
		query.Result.Account = res.Shop
		return nil

	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Missing AccountType")
	}
}
