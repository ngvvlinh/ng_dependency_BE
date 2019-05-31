package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandlers("sql",
		CreateAddress,
		GetAddress,
		UpdateAddress,
		DeleteAddress,
		GetAddresses,
	)
}

func GetAddress(ctx context.Context, query *model.GetAddressQuery) error {
	if query.AddressID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	query.Result = new(model.Address)
	return x.Table("address").
		Where("id = ?", query.AddressID).
		ShouldGet(query.Result)
}

func GetAddresses(ctx context.Context, query *model.GetAddressesQuery) error {
	s := x.Table("address")
	if query.AccountID != 0 {
		s = s.Where("account_id = ?", query.AccountID)
	}
	if err := s.Find((*model.Addresses)(&query.Result.Addresses)); err != nil {
		return err
	}
	return nil
}

func CreateAddress(ctx context.Context, cmd *model.CreateAddressCommand) error {
	return createAddress(ctx, x, cmd)
}

func createAddress(ctx context.Context, x Qx, cmd *model.CreateAddressCommand) error {
	address := cmd.Address

	if address.Province == "" || address.ProvinceCode == "" {
		return cm.Error(cm.InvalidArgument, "Missing province information", nil)
	}
	if address.District == "" || address.DistrictCode == "" {
		return cm.Error(cm.InvalidArgument, "Missing district information", nil)
	}
	if address.Ward == "" || address.WardCode == "" {
		return cm.Error(cm.InvalidArgument, "Missing ward information", nil)
	}
	if address.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	// if err := location.CheckValidLocation(address.ProvinceCode, "province"); err != nil {
	// 	return err
	// }
	//
	// if err := location.CheckValidLocation(address.DistrictCode, "district"); err != nil {
	// 	return err
	// }
	//
	// if err := location.CheckValidLocation(address.WardCode, "ward"); err != nil {
	// 	return err
	// }

	address.ID = cm.NewID()
	if _, err := x.Table("address").Insert(address); err != nil {
		return err
	}
	cmd.Result = address
	return nil
}

func UpdateAddress(ctx context.Context, cmd *model.UpdateAddressCommand) error {
	return updateAddress(ctx, x, cmd)
}

func updateAddress(ctx context.Context, x Qx, cmd *model.UpdateAddressCommand) error {
	address := cmd.Address
	if address.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	// if address.ProvinceCode != "" {
	// 	if err := location.CheckValidLocation(address.ProvinceCode, "province"); err != nil {
	// 		return err
	// 	}
	// }
	// if address.DistrictCode != "" {
	// 	if err := location.CheckValidLocation(address.DistrictCode, "district"); err != nil {
	// 		return err
	// 	}
	// }
	// if address.WardCode != "" {
	// 	if err := location.CheckValidLocation(address.WardCode, "ward"); err != nil {
	// 		return err
	// 	}
	// }

	if err := x.Table("address").
		Where("id = ?", address.ID).
		ShouldUpdate(address); err != nil {
		return err
	}
	cmd.Result = address
	return nil
}

func DeleteAddress(ctx context.Context, cmd *model.DeleteAddressCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	if cmd.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	s := x.Table("address").Where("id = ? AND account_id = ?", cmd.ID, cmd.AccountID)
	if deleted, err := s.Delete(&model.Address{}); err != nil {
		return err
	} else if deleted == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	return nil
}
