package aggregate

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"

	"etop.vn/api/shopping/suppliering"
	"etop.vn/backend/com/shopping/suppliering/convert"
	"etop.vn/backend/com/shopping/suppliering/model"
	"etop.vn/backend/com/shopping/suppliering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
)

var _ suppliering.Aggregate = &SupplierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SupplierAggregate struct {
	store sqlstore.SupplierStoreFactory
}

func NewSupplierAggregate(db *cmsql.Database) *SupplierAggregate {
	return &SupplierAggregate{
		store: sqlstore.NewSupplierStore(db),
	}
}

func (a *SupplierAggregate) MessageBus() suppliering.CommandBus {
	b := bus.New()
	return suppliering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SupplierAggregate) CreateSupplier(
	ctx context.Context, args *suppliering.CreateSupplierArgs,
) (*suppliering.ShopSupplier, error) {
	if args.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập tên đầy đủ")
	}
	if args.Phone != "" {
		phone, isPhone := validate.NormalizePhone(args.Phone)
		if isPhone != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
		}
		args.Phone = phone.String()
		_, err := a.store(ctx).Phone(args.Phone).ShopID(args.ShopID).GetSupplier()
		if err == nil {
			return nil, cm.Error(cm.InvalidArgument, "số điện thoại đã tồn tại", nil)
		}
	}
	if args.Email != "" {
		email, isEmail := validate.NormalizeEmail(args.Email)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email = email.String()
		_, err := a.store(ctx).Email(args.Email).ShopID(args.ShopID).GetSupplier()
		if err == nil {
			return nil, cm.Error(cm.InvalidArgument, "Email đã tồn tại", nil)
		}
	}
	supplier := new(suppliering.ShopSupplier)

	if err := scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateSupplier(supplier)
	return supplier, err
}

func (a *SupplierAggregate) UpdateSupplier(
	ctx context.Context, args *suppliering.UpdateSupplierArgs,
) (*suppliering.ShopSupplier, error) {
	supplier, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetSupplier()
	if err != nil {
		return nil, err
	}
	if args.Phone.Valid {
		if args.Phone.String == "" {
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không thể rỗng", err)
		} else {
			phone, isPhone := validate.NormalizePhone(args.Phone.String)
			if isPhone != true {
				return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
			}
			args.Phone.String = phone.String()
			supplier, err := a.store(ctx).Phone(args.Phone.String).ShopID(args.ShopID).GetSupplier()
			if err == nil && args.ID != supplier.ID {
				return nil, cm.Error(cm.InvalidArgument, "số điện thoại đã tồn tại", nil)
			}
		}
	}
	if args.Email.Valid && args.Email.String != "" {
		email, isEmail := validate.NormalizeEmail(args.Email.String)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email.String = email.String()
		supplier, err := a.store(ctx).Email(args.Email.String).ShopID(args.ShopID).GetSupplier()
		if err == nil && args.ID != supplier.ID {
			return nil, cm.Error(cm.InvalidArgument, "Email đã tồn tại", nil)
		}
	}
	if err := scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	supplierDB := new(model.ShopSupplier)
	if err := scheme.Convert(supplier, supplierDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateSupplierDB(supplierDB)
	return supplier, err
}

func (a *SupplierAggregate) DeleteSupplier(
	ctx context.Context, ID int64, shopID int64,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return deleted, err
}
