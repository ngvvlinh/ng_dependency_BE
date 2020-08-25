package sqlstore

import (
	"context"
	"strings"

	"o.o/api/main/address"
	"o.o/backend/com/main/address/convert"
	"o.o/backend/com/main/address/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	sq "o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (ft *AddressFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

const tableName = "address"

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type AddressFactory func(context.Context) *AddressStore

func NewAddressStore(db *cmsql.Database) AddressFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AddressStore {
		return &AddressStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AddressStore struct {
	query cmsql.QueryFactory
	ft    AddressFilters
	sqlstore.Paging
	preds []interface{}
}

func (s *AddressStore) ID(id dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *AddressStore) Type(typeAddress string) *AddressStore {
	s.preds = append(s.preds, s.ft.ByType(typeAddress))
	return s
}

func (s *AddressStore) AccountID(id dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AddressStore) IsDefault(isDefault bool) *AddressStore {
	s.preds = append(s.preds, s.ft.ByIsDefault(isDefault))
	return s
}

func (s *AddressStore) CreateAddress(address *address.Address) error {
	sqlstore.MustNoPreds(s.preds)
	addressDB := convert.Convert_address_Address_addressmodel_Address(address, nil)
	if _, err := s.query().Insert(addressDB); err != nil {
		return err
	}
	return nil
}

func (s *AddressStore) UpdateAddress(address *address.Address) (*address.Address, error) {
	addressDB := convert.Convert_address_Address_addressmodel_Address(address, nil)

	if err := s.query().Where(s.ft.ByID(address.ID)).ShouldUpdate(addressDB); err != nil {
		return nil, err
	}

	return s.ID(address.ID).Get()
}

func (s *AddressStore) Update(address *address.Address) error {
	addressDB := convert.Convert_address_Address_addressmodel_Address(address, nil)
	query := s.query().Where(s.preds)

	return query.ShouldUpdate(addressDB)
}

func (s *AddressStore) UpdateDefault(isDefault bool) (*address.Address, error) {
	if _, err := s.query().Where(s.preds).Table(tableName).UpdateMap(map[string]interface{}{
		"is_default": isDefault,
	}); err != nil {
		return nil, err
	}

	return s.Get()
}

func (s *AddressStore) Delete() error {
	query := s.query().Where(s.preds)
	if deleted, err := query.Delete(&model.Address{}); err != nil {
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

func (s *AddressStore) GetDB() (*model.Address, error) {
	var addressModel model.Address
	err := s.query().Where(s.preds).ShouldGet(&addressModel)
	return &addressModel, err
}

func (s *AddressStore) Get() (res *address.Address, _ error) {
	addressDB, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	return convert.Address(addressDB), nil
}

func (s *AddressStore) ListAddressDBs() ([]*model.Address, error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-updated_at"}
	}

	var addrs model.Addresses
	err := query.Find(&addrs)
	return addrs, err
}

func (s *AddressStore) ListAddresses() (res []*address.Address, _ error) {
	addressLists, err := s.ListAddressDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(addressLists, &res); err != nil {
		return nil, err
	}
	return
}
