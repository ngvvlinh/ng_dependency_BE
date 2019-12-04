package aggregate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/addressing"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/capi/dot"
)

var _ addressing.Aggregate = &AddressAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

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

func (q *AddressAggregate) CreateAddress(ctx context.Context, args *addressing.CreateAddressArgs) (*addressing.ShopTraderAddress, error) {
	err := ValidateCreateShopTraderAddress(args)
	if err != nil {
		return nil, err
	}
	if err := q.store(ctx).UpdateStatusAddresses(args.ShopID, args.TraderID, false); err != nil {
		return nil, err
	}

	addr := &addressing.ShopTraderAddress{}
	if err = scheme.Convert(args, addr); err != nil {
		return nil, err
	}
	err = q.store(ctx).CreateAddress(addr)
	return addr, err
}

func (q *AddressAggregate) UpdateAddress(ctx context.Context, ID dot.ID, ShopID dot.ID, args *addressing.UpdateAddressArgs) (*addressing.ShopTraderAddress, error) {
	addr, err := q.store(ctx).ID(ID).ShopID(ShopID).GetAddress()
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
	err = q.store(ctx).UpdateAddressDB(addrDB)
	return addr, err
}

func (q *AddressAggregate) DeleteAddress(ctx context.Context, ID dot.ID, ShopID dot.ID) (deleted int, _ error) {
	deleted, err := q.store(ctx).ID(ID).ShopID(ShopID).SoftDelete()
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

func (q *AddressAggregate) SetDefaultAddress(
	ctx context.Context, ID, traderID, shopID dot.ID,
) (*meta.UpdatedResponse, error) {
	if err := q.store(ctx).UpdateStatusAddresses(shopID, traderID, false); err != nil {
		return nil, err
	}

	updated, err := q.store(ctx).SetDefaultAddress(ID, shopID, traderID)
	return &meta.UpdatedResponse{Updated: updated}, err
}
