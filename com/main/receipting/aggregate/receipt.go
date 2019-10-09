package aggregate

import (
	"context"

	"etop.vn/backend/com/main/receipting/convert"

	"etop.vn/backend/com/main/receipting/model"
	"etop.vn/backend/pkg/common/scheme"

	"etop.vn/api/main/receipting"
	"etop.vn/backend/com/main/receipting/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ receipting.Aggregate = &ReceiptAggregate{}

type ReceiptAggregate struct {
	store sqlstore.ReceiptStoreFactory
}

func NewReceiptAggregate(db cmsql.Database) *ReceiptAggregate {
	return &ReceiptAggregate{
		store: sqlstore.NewReceiptStore(db),
	}
}

func (a *ReceiptAggregate) MessageBus() receipting.CommandBus {
	b := bus.New()
	return receipting.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ReceiptAggregate) CreateReceipt(
	ctx context.Context, args *receipting.CreateReceiptArgs,
) (*receipting.Receipt, error) {
	receipt := new(receipting.Receipt)
	if err := scheme.Convert(args, receipt); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateReceipt(receipt)
	return receipt, err
}

func (a *ReceiptAggregate) UpdateReceipt(
	ctx context.Context, args *receipting.UpdateReceiptArgs,
) (*receipting.Receipt, error) {
	receipt, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetReceipt()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(args, receipt); err != nil {
		return nil, err
	}
	receiptDB := new(model.Receipt)
	if err := scheme.Convert(receipt, receiptDB); err != nil {
		return nil, err
	}
	receiptDB.Lines = convert.Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(receipt.Lines)
	err = a.store(ctx).UpdateReceiptDB(receiptDB)
	return receipt, err
}

func (a *ReceiptAggregate) DeleteReceipt(
	ctx context.Context, id int64, shopID int64,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(id).ShopID(shopID).SoftDelete()
	return deleted, err
}
