package sqlstore

import (
	"context"
	"strings"

	com "o.o/backend/com/main"
	addressmodel "o.o/backend/com/main/address/model"
	addressmodelx "o.o/backend/com/main/address/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
)

type AddressStoreInterface interface {
	GetAddress(ctx context.Context, query *addressmodelx.GetAddressQuery) error
}

type AddressStore struct {
	DB com.MainDB
	db *cmsql.Database `wire:"-"`
}

func BindAddressStore(s *AddressStore) (to AddressStoreInterface) {
	s.db = s.DB
	return s
}

func (st *AddressStore) GetAddress(ctx context.Context, query *addressmodelx.GetAddressQuery) error {
	if query.AddressID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	query.Result = new(addressmodel.Address)
	return st.db.Table("address").
		Where("id = ?", query.AddressID).
		ShouldGet(query.Result)
}

func (st *AddressStore) GetAddresses(ctx context.Context, query *addressmodelx.GetAddressesQuery) error {
	s := st.db.Table("address")
	if query.AccountID != 0 {
		s = s.Where("account_id = ?", query.AccountID)
	}
	if err := s.Find((*addressmodel.Addresses)(&query.Result.Addresses)); err != nil {
		return err
	}
	return nil
}

func (st *AddressStore) CreateAddress(ctx context.Context, cmd *addressmodelx.CreateAddressCommand) error {
	return st.createAddress(ctx, st.db, cmd)
}

func (st *AddressStore) createAddress(ctx context.Context, x Qx, cmd *addressmodelx.CreateAddressCommand) error {
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

func (st *AddressStore) UpdateAddress(ctx context.Context, cmd *addressmodelx.UpdateAddressCommand) error {
	return st.updateAddress(ctx, st.db, cmd)
}

func (st *AddressStore) updateAddress(ctx context.Context, tx Qx, cmd *addressmodelx.UpdateAddressCommand) error {
	address := cmd.Address
	if address.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	if err := tx.Table("address").
		Where("id = ?", address.ID).
		ShouldUpdate(address); err != nil {
		return err
	}
	cmd.Result = address
	return nil
}

func (st *AddressStore) DeleteAddress(ctx context.Context, cmd *addressmodelx.DeleteAddressCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AddressID", nil)
	}

	if cmd.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	s := st.db.Table("address").Where("id = ? AND account_id = ?", cmd.ID, cmd.AccountID)
	if deleted, err := s.Delete(&addressmodel.Address{}); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "shop_address_id_fkey") || strings.Contains(errMsg, "shop_ship_from_address_id_fkey") {
			err = cm.Errorf(cm.FailedPrecondition, nil, "không thể xóa địa chỉ mặc định")
		}
		return err
	} else if deleted == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	return nil
}
