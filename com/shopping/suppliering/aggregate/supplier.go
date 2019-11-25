package aggregate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/suppliering/convert"
	"etop.vn/backend/com/shopping/suppliering/model"
	"etop.vn/backend/com/shopping/suppliering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ suppliering.Aggregate = &SupplierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SupplierAggregate struct {
	store    sqlstore.SupplierStoreFactory
	eventBus capi.EventBus
}

func NewSupplierAggregate(eventBus capi.EventBus, db *cmsql.Database) *SupplierAggregate {
	return &SupplierAggregate{
		store:    sqlstore.NewSupplierStore(db),
		eventBus: eventBus,
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
	if args.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập số điện thoại đầy đủ")
	}
	phone, isPhone := validate.NormalizePhone(args.Phone)
	if isPhone != true {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
	}
	args.Phone = phone.String()
	sup, err := a.store(ctx).Phone(args.Phone).ShopID(args.ShopID).GetSupplier()
	if err == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại %v đã tồn tại", sup.Phone)
	}
	if args.Email != "" {
		email, isEmail := validate.NormalizeEmail(args.Email)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email = email.String()
		sup, err := a.store(ctx).Email(args.Email).ShopID(args.ShopID).GetSupplier()
		if err == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Nhà cung cấp với email: %v đã tồn tại", sup.Email)
		}
	}
	supplier := new(suppliering.ShopSupplier)

	if err = scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	var maxCodeNorm int32
	supplierTemp, err := a.store(ctx).ShopID(args.ShopID).IncludeDeleted().GetSupplierByMaximumCodeNorm()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = supplierTemp.CodeNorm
	case cm.NotFound:
		// no-op
	default:
		return nil, err
	}

	if maxCodeNorm >= convert.MaxCodeNorm {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
	}
	codeNorm := maxCodeNorm + 1
	supplier.Code = convert.GenerateCode(int(codeNorm))
	supplier.CodeNorm = codeNorm

	err = a.store(ctx).CreateSupplier(supplier)
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
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không thể để trống", nil)
		} else {
			phone, isPhone := validate.NormalizePhone(args.Phone.String)
			if isPhone != true {
				return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
			}
			args.Phone.String = phone.String()
			supplierByPhone, err := a.store(ctx).Phone(args.Phone.String).ShopID(args.ShopID).GetSupplier()
			if err == nil && args.ID != supplierByPhone.ID {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "số điện thoại %v đã tồn tại", supplierByPhone.Phone)
			}
		}
	}
	if args.Email.Valid && args.Email.String != "" {
		email, isEmail := validate.NormalizeEmail(args.Email.String)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email.String = email.String()
		supplierByEmail, err := a.store(ctx).Email(args.Email.String).ShopID(args.ShopID).GetSupplier()
		if err == nil && args.ID != supplierByEmail.ID {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Nhà cung cấp với email: %v đã tồn tại", supplierByEmail.Email)
		}
	}
	if err = scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	supplierDB := new(model.ShopSupplier)
	if err = scheme.Convert(supplier, supplierDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateSupplierDB(supplierDB)
	return supplier, err
}

func (a *SupplierAggregate) DeleteSupplier(
	ctx context.Context, ID dot.ID, shopID dot.ID,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	event := &tradering.TraderDeletedEvent{
		EventMeta: meta.NewEvent(),
		ShopID:    shopID,
		TraderID:  ID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return 0, err
	}
	eventDeleVariantSupplier := &suppliering.VariantSupplierDeletedEvent{
		ShopID:     shopID,
		SupplierID: ID,
	}
	if err := a.eventBus.Publish(ctx, eventDeleVariantSupplier); err != nil {
		return 0, err
	}
	return deleted, err
}
