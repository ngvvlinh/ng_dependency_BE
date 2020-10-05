package aggregate

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping/addressing"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/customering/convert"
	"o.o/backend/com/shopping/customering/model"
	"o.o/backend/com/shopping/customering/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

var _ addressing.Aggregate = &AddressAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type AddressAggregate struct {
	store         sqlstore.AddressStoreFactory
	customerStore sqlstore.CustomerStoreFactory
}

func NewAddressAggregate(db com.MainDB) *AddressAggregate {
	return &AddressAggregate{
		store:         sqlstore.NewAddressStore(db),
		customerStore: sqlstore.NewCustomerStore(db),
	}
}

func AddressAggregateMessageBus(q *AddressAggregate) addressing.CommandBus {
	b := bus.New()
	return addressing.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *AddressAggregate) CreateAddress(ctx context.Context, args *addressing.CreateAddressArgs) (*addressing.ShopTraderAddress, error) {
	if args.Phone == "" {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng cung cấp số điện thoại", nil)
	}
	phoneNorm, ok := validate.NormalizePhone(args.Phone)
	if !ok {
		return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}

	args.Phone = phoneNorm.String()

	// check customer
	{
		customer, err := q.customerStore(ctx).ID(args.TraderID).IncludeDeleted().GetCustomer()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			return nil, cm.Error(cm.InvalidArgument, "Không tìm thấy user", nil)
		case cm.NoError:
			if customer.Deleted {
				return nil, cm.Error(cm.InvalidArgument, "Khách hàng đã bị xóa. Vui lòng kiểm tra lại", nil)
			}
		default:
			return nil, err
		}
	}

	if args.IsDefault {
		if err := q.store(ctx).UpdateStatusAddresses(args.ShopID, args.TraderID, false); err != nil {
			return nil, err
		}
	}

	addr := &addressing.ShopTraderAddress{}
	if err := scheme.Convert(args, addr); err != nil {
		return nil, err
	}
	if err := ValidateCreateShopTraderAddress(addr); err != nil {
		return nil, err
	}
	err := q.store(ctx).CreateAddress(addr)
	return addr, err
}

func (q *AddressAggregate) UpdateAddress(ctx context.Context, ID dot.ID, ShopID dot.ID, args *addressing.UpdateAddressArgs) (*addressing.ShopTraderAddress, error) {
	addr, err := q.store(ctx).ID(ID).ShopID(ShopID).GetAddress()
	if err != nil {
		return nil, err
	}
	if args.Phone.Valid {
		if args.Phone.String == "" {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng cung cấp số điện thoại", nil)
		}
		_, ok := validate.NormalizePhone(args.Phone.String)
		if !ok {
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
		}
	}
	if err = scheme.Convert(args, addr); err != nil {
		return nil, err
	}
	err = ValidateCreateShopTraderAddress(addr)
	if err != nil {
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

func ValidateCreateShopTraderAddress(args *addressing.ShopTraderAddress) error {
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
