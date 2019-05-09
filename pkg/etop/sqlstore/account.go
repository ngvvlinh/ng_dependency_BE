package sqlstore

import (
	"context"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandlers("sql",
		CreateSupplierKiotviet,
		CreateShop,
		UpdateShop,
		DeleteShop,
		SetDefaultAddressShop,
		SetDefaultAddressSupplier,
		UpdateAccountURLSlug,
		GetAccountAuth,
	)
}

func CreateSupplierKiotviet(ctx context.Context, cmd *model.CreateSupplierKiotvietCommand) error {
	if cmd.Kiotviet.RetailerID == "" {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if cmd.OwnerID == 0 {
		return cm.Error(cm.Internal, "Missing OwnerID", nil)
	}

	var ps model.ProductSource
	if has, err := x.
		Where("external_key = ?", cmd.Kiotviet.RetailerID).
		Get(&ps); err != nil {
		return err
	} else if has {
		return cm.Error(cm.AlreadyExists,
			"Tài khoản Kiotviet đã được đăng ký. Vui lòng kiểm tra lại.", nil)
	}

	return inTransaction(func(s Qx) error {
		supplierID := model.NewSupplierID()
		sourceID := model.NewID()

		account := &model.Account{
			ID:       supplierID,
			Name:     cmd.Name,
			Type:     model.TypeSupplier,
			ImageURL: cmd.ImageURL,
			URLSlug:  cmd.URLSlug,
		}

		supplier := &model.Supplier{
			ID:              supplierID,
			ProductSourceID: sourceID,
			OwnerID:         cmd.OwnerID,
		}
		if cmd.IsTest {
			supplier.IsTest = 1
		}
		cmd.SupplierInfo.AssignTo(supplier)

		accountUser := &model.AccountUser{
			AccountID: supplierID,
			UserID:    cmd.OwnerID,
			Status:    model.StatusActive,
		}

		productSource := &model.ProductSource{
			ID:         sourceID,
			Type:       model.TypeKiotviet,
			Name:       cmd.Kiotviet.RetailerID,
			SupplierID: supplierID,
		}

		kv := cmd.Kiotviet
		productSourceInternal := &model.ProductSourceInternal{
			ID: sourceID,
			Secret: &model.KiotvietSecret{
				RetailerID:   kv.RetailerID,
				ClientID:     kv.ClientID,
				ClientSecret: kv.ClientSecret,
			},
			AccessToken: kv.ClientToken,
			ExpiresAt:   kv.ExpiresAt,
		}

		cmd.Result.Supplier = supplier
		cmd.Result.ProductSource = productSource
		cmd.Result.ProductSourceInternal = productSourceInternal

		_, err := s.Insert(
			account, supplier, accountUser,
			productSource, productSourceInternal)
		return err
	})
}

func CreateShop(ctx context.Context, cmd *model.CreateShopCommand) error {
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

	id := model.NewShopID()
	return inTransaction(func(x Qx) error {
		account := &model.Account{
			ID:       id,
			Name:     cmd.Name,
			Type:     model.TypeShop,
			ImageURL: cmd.ImageURL,
			URLSlug:  cmd.URLSlug,
		}
		permission := &model.AccountUser{
			AccountID: id,
			UserID:    cmd.OwnerID,
			Status:    model.StatusActive,
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

		shop := &model.Shop{
			ID:                            id,
			Name:                          cmd.Name,
			OwnerID:                       cmd.OwnerID,
			Status:                        model.StatusActive,
			AddressID:                     cmd.AddressID,
			Phone:                         phoneNorm.String(),
			Email:                         emailNorm.String(),
			BankAccount:                   cmd.BankAccount,
			WebsiteURL:                    cmd.WebsiteURL,
			ImageURL:                      cmd.ImageURL,
			Code:                          code,
			TryOn:                         model.TryOnOpen,
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

		cmd.Result = new(model.ShopExtended)
		if has, err := x.
			Table("shop").
			Where("s.id = ?", shop.ID).
			Get(cmd.Result); err != nil || !has {
			return cm.Error(cm.Internal, "", err)
		}
		return nil
	})
}

func UpdateShop(ctx context.Context, cmd *model.UpdateShopCommand) error {
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

		cmd.Result = new(model.ShopExtended)
		if has, err := x.
			Table("shop").
			Where("s.id = ?", shop.ID).
			Get(cmd.Result); err != nil || !has {
			return cm.Error(cm.Internal, "", err)
		}
		return nil
	})
}

func DeleteShop(ctx context.Context, cmd *model.DeleteShopCommand) error {
	return inTransaction(func(s Qx) error {
		if cmd.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing ID", nil)
		}
		now := time.Now()
		{
			if updated, err := s.Table("shop").
				Where("id = ? AND owner_id = ?", cmd.ID, cmd.OwnerID).
				Update(&model.ShopDelete{
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
				&model.AccountUserDelete{
					DeletedAt: now,
				}); err != nil {
			return err
		}
		return nil
	})
}

func UpdateOrCreateAddress(ctx context.Context, address *model.Address, accountID int64, AddressType string) (int64, error) {
	return updateOrCreateAddress(ctx, x, address, accountID, AddressType)
}

func updateOrCreateAddress(ctx context.Context, x Qx, address *model.Address, accountID int64, AddressType string) (int64, error) {
	addressObj := &model.Address{
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
		addressCmd := &model.UpdateAddressCommand{
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
		addressCmd := &model.CreateAddressCommand{
			Address: addressObj,
		}
		if err := createAddress(ctx, x, addressCmd); err != nil {
			return 0, err
		}
		return addressCmd.Result.ID, nil
	}
}

func SetDefaultAddressShop(ctx context.Context, cmd *model.SetDefaultAddressShopCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Type == "" {
		return cm.Error(cm.InvalidArgument, "Missing Address Type", nil)
	}
	if cmd.AddressID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	var address = new(model.Address)
	if err := x.Table("address").Where("id = ? AND account_id = ? and type = ?", cmd.AddressID, cmd.ShopID, cmd.Type).
		ShouldGet(address); err != nil {
		return err
	}
	shopObj := &model.Shop{
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

	cmdUpdateShop := &model.UpdateShopCommand{
		Shop: shopObj,
	}
	if err := bus.Dispatch(ctx, cmdUpdateShop); err != nil {
		return err
	}
	cmd.Result.Updated = 1
	return nil
}

func SetDefaultAddressSupplier(ctx context.Context, cmd *model.SetDefaultAddressSupplierCommand) error {
	if cmd.SupplierID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}
	if cmd.Type == "" {
		return cm.Error(cm.InvalidArgument, "Missing Address Type", nil)
	}
	if cmd.AddressID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	var address = new(model.Address)
	if err := x.Table("address").Where("id = ? AND account_id = ? and type = ?", cmd.AddressID, cmd.SupplierID, cmd.Type).
		ShouldGet(address); err != nil {
		return err
	}
	supplierObj := &model.Supplier{
		ID: cmd.SupplierID,
	}

	switch cmd.Type {
	case model.AddressTypeShipFrom:
		supplierObj.ShipFromAddressID = cmd.AddressID
	default:
		return cm.Error(cm.Unimplemented, "Address type does not valid", nil)
	}

	cmdUpdateSupplier := &model.UpdateSupplierCommand{
		Supplier: supplierObj,
	}
	if err := bus.Dispatch(ctx, cmdUpdateSupplier); err != nil {
		return err
	}
	cmd.Result.Updated = 1
	return nil
}

func UpdateAccountURLSlug(ctx context.Context, cmd *model.UpdateAccountURLSlugCommand) error {
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

func GetAccountAuth(ctx context.Context, query *model.GetAccountAuthQuery) error {
	if query.AuthKey == "" || query.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing key")
	}

	s := x.Where("auth_key = ?", query.AuthKey).
		Where("account_id = ?", query.AccountID).
		Where("aa.status = 1 AND aa.deleted_at IS NULL")

	switch query.AccountType {
	case model.TypePartner:
		if cm.GetTag(query.AccountID) != model.TagPartner {
			return cm.Errorf(cm.NotFound, nil, "")
		}

		var res model.AccountAuthFtPartner
		if err := s.Where("p.status = 1 AND p.deleted_at IS NULL").
			ShouldGet(&res); err != nil {
			return err
		}
		query.Result.AccountAuth = res.AccountAuth
		query.Result.Account = res.Partner
		return nil

	case model.TypeShop:
		if cm.GetTag(query.AccountID) != model.TagShop {
			return cm.Errorf(cm.NotFound, nil, "")
		}

		var res model.AccountAuthFtShop
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
