package aggregate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/addressing"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/scheme"
)

var _ addressing.Aggregate = &AddressAggregate{}

type AddressAggregate struct {
	store sqlstore.AddressStoreFactory
}

func NewAddressAggregate(db *cmsql.Database) *AddressAggregate {
	return &AddressAggregate{
		store: sqlstore.NewAddressStore(db),
	}
}

func (q *AddressAggregate) MessageBus() addressing.CommandBus {
	b := bus.New()
	return addressing.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a *AddressAggregate) CreateAddress(ctx context.Context, args *addressing.CreateAddressArgs) (*addressing.ShopTraderAddress, error) {
	err := ValidateCreateShopTraderAddress(args)
	if err != nil {
		return nil, err
	}
	if err := a.store(ctx).UpdateStatusAddresses(args.ShopID, args.TraderID, false); err != nil {
		return nil, err
	}

	addr := &addressing.ShopTraderAddress{}
	if err = scheme.Convert(args, addr); err != nil {
		return nil, err
	}
	err = a.store(ctx).CreateAddress(addr)
	return addr, err
}

func (a *AddressAggregate) UpdateAddress(ctx context.Context, ID int64, ShopID int64, args *addressing.UpdateAddressArgs) (*addressing.ShopTraderAddress, error) {
	addr, err := a.store(ctx).ID(ID).ShopID(ShopID).GetAddress()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(args, addr); err != nil {
		return nil, err
	}
	addrDB := &model.ShopTraderAddress{}
	if err = scheme.Convert(addr, addrDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateAddressDB(addrDB)
	return addr, err
}

func (a *AddressAggregate) DeleteAddress(ctx context.Context, ID int64, ShopID int64) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(ID).ShopID(ShopID).SoftDelete()
	return deleted, err
}

func ValidateCreateShopTraderAddress(args *addressing.CreateAddressArgs) error {
	if args.FullName == "" {
		return EditErrorMsg("Tên")
	}
	if args.DistrictCode == "" {
		return EditErrorMsg("Quận/Huyện")
	}
	if args.WardCode == "" {
		return EditErrorMsg("Phường/Xã")
	}
	if args.Address1 == "" {
		return EditErrorMsg("Địa chỉ cụ thể")
	}
	if args.Phone == "" {
		return EditErrorMsg("Số điện thoại")
	}
	return nil
}

func EditErrorMsg(str string) error {
	return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập thông tin bắt buộc, thiếu %v", str)
}

func (a *AddressAggregate) SetDefaultAddress(
	ctx context.Context, ID, traderID, shopID int64,
) (*meta.UpdatedResponse, error) {
	if err := a.store(ctx).UpdateStatusAddresses(shopID, traderID, false); err != nil {
		return nil, err
	}

	updated, err := a.store(ctx).SetDefaultAddress(ID, shopID, traderID)
	return &meta.UpdatedResponse{Updated: int32(updated)}, err
}
