package sqlstore

import (
	"context"

	"o.o/api/main/address"
	"o.o/backend/com/main/address/convert"
	addressmodel "o.o/backend/com/main/address/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

type AddressStore struct {
	ctx context.Context
	db  *cmsql.Database
}

func NewAddressStore(db *cmsql.Database) *AddressStore {
	return &AddressStore{
		ctx: context.Background(),
		db:  db,
	}
}

func (s *AddressStore) WithContext(ctx context.Context) *AddressStore {
	return &AddressStore{
		ctx: ctx,
		db:  s.db,
	}
}

func (s *AddressStore) GetByID(ID dot.ID) (*address.Address, error) {
	if ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing ID")
	}

	var res = new(addressmodel.Address)
	if err := s.db.WithContext(s.ctx).Where("id = ?", ID).ShouldGet(res); err != nil {
		return nil, err
	}
	return convert.Address(res), nil
}
