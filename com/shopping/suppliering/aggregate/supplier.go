package aggregate

import (
	"context"

	"o.o/api/shopping/suppliering"
	"o.o/api/shopping/tradering"
	"o.o/backend/com/shopping/suppliering/convert"
	"o.o/backend/com/shopping/suppliering/model"
	"o.o/backend/com/shopping/suppliering/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ suppliering.Aggregate = &SupplierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SupplierAggregate struct {
	store    sqlstore.SupplierStoreFactory
	db       *cmsql.Database
	eventBus capi.EventBus
}

func NewSupplierAggregate(eventBus capi.EventBus, db *cmsql.Database) *SupplierAggregate {
	return &SupplierAggregate{
		store:    sqlstore.NewSupplierStore(db),
		db:       db,
		eventBus: eventBus,
	}
}

func SupplierAggregateMessageBus(a *SupplierAggregate) suppliering.CommandBus {
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
	if args.Email != "" {
		email, isEmail := validate.NormalizeEmail(args.Email)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email = email.String()
	}
	supplier := new(suppliering.ShopSupplier)

	if err := scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	var maxCodeNorm int
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
	supplier.Code = convert.GenerateCode(codeNorm)
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
			_, isPhone := validate.NormalizePhone(args.Phone.String)
			if isPhone != true {
				return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
			}
		}
	}
	if args.Email.Valid && args.Email.String != "" {
		_, isEmail := validate.NormalizeEmail(args.Email.String)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
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
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		deletedTrader, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
		if err != nil {
			return err
		}
		event := &tradering.TraderDeletedEvent{
			ShopID:   shopID,
			TraderID: ID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		deleted = deletedTrader
		return nil
	})
	return deleted, err
}
