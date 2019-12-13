package aggregate

import (
	"context"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/receipting/convert"
	"etop.vn/backend/com/main/receipting/model"
	"etop.vn/backend/com/main/receipting/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	. "etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()
var _ receipting.Aggregate = &ReceiptAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ReceiptAggregate struct {
	db                 *cmsql.Database
	store              sqlstore.ReceiptStoreFactory
	eventBus           capi.EventBus
	traderQuery        tradering.QueryBus
	ledgerQuery        ledgering.QueryBus
	orderQuery         ordering.QueryBus
	customerQuery      customering.QueryBus
	carrierQuery       carrying.QueryBus
	supplierQuery      suppliering.QueryBus
	purchaseOrderQuery purchaseorder.QueryBus
}

func NewReceiptAggregate(
	database *cmsql.Database, eventBus capi.EventBus,
	traderQuery tradering.QueryBus, ledgerQuery ledgering.QueryBus,
	orderQuery ordering.QueryBus, customerQuery customering.QueryBus,
	carrierQuery carrying.QueryBus, supplierQuery suppliering.QueryBus,
	purchaseOrderQ purchaseorder.QueryBus,
) *ReceiptAggregate {
	return &ReceiptAggregate{
		db:                 database,
		store:              sqlstore.NewReceiptStore(database),
		eventBus:           eventBus,
		traderQuery:        traderQuery,
		ledgerQuery:        ledgerQuery,
		orderQuery:         orderQuery,
		customerQuery:      customerQuery,
		carrierQuery:       carrierQuery,
		supplierQuery:      supplierQuery,
		purchaseOrderQuery: purchaseOrderQ,
	}
}

func (a *ReceiptAggregate) MessageBus() receipting.CommandBus {
	b := bus.New()
	return receipting.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ReceiptAggregate) CreateReceipt(
	ctx context.Context, args *receipting.CreateReceiptArgs,
) (*receipting.Receipt, error) {
	if args.PaidAt.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ngày tạo phiếu không hợp lệ")
	}
	if args.TraderID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Đối tác không hợp lệ")
	}
	if args.LedgerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tài khoản thanh toán không hợp lệ")
	}

	receiptNeedValidate := &receipting.Receipt{
		TraderID:    args.TraderID,
		ShopID:      args.ShopID,
		CreatedBy:   args.CreatedBy,
		Title:       args.Title,
		Type:        args.Type,
		RefType:     args.RefType,
		Description: args.Description,
		Amount:      args.Amount,
		LedgerID:    args.LedgerID,
		PaidAt:      args.PaidAt,
		Lines:       args.Lines,
	}
	if err := a.validateReceiptForCreateOrUpdate(ctx, args.ShopID, receiptNeedValidate); err != nil {
		return nil, err
	}

	args.Trader = receiptNeedValidate.Trader
	receipt := new(receipting.Receipt)
	if err := scheme.Convert(args, receipt); err != nil {
		return nil, err
	}

	var maxCodeNorm int
	receiptTemp, err := a.store(ctx).ShopID(args.ShopID).GetReceiptByMaximumCodeNorm()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = receiptTemp.CodeNorm
	case cm.NotFound:
	// no-op
	default:
		return nil, err
	}

	if maxCodeNorm >= convert.MaxCodeNorm {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
	}
	codeNorm := maxCodeNorm + 1
	receipt.Code = convert.GenerateCode(codeNorm)
	receipt.CodeNorm = codeNorm

	if err = a.store(ctx).CreateReceipt(receipt); err != nil {
		return nil, err
	}

	if !args.ConfirmedAt.IsZero() {
		receiptConfirmEvent := &receipting.ReceiptConfirmedEvent{
			ShopID:    args.ShopID,
			ReceiptID: receipt.ID,
		}
		if err := a.eventBus.Publish(ctx, receiptConfirmEvent); err != nil {
			ll.Error("receiptConfirmedEvent", l.Error(err))
		}
	}

	return receipt, err
}

func (a *ReceiptAggregate) UpdateReceipt(
	ctx context.Context, args *receipting.UpdateReceiptArgs,
) (*receipting.Receipt, error) {
	receipt, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetReceipt()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "không tìm thấy phiếu").
			Throw()
	}

	if receipt.Status == status3.N {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể thay đổi phiếu đã hủy.")
	}

	receiptNeedValidate := &receipting.Receipt{
		ID:          args.ID,
		Title:       args.Title.String,
		Description: args.Description.String,
		LedgerID:    args.LedgerID.ID,
		TraderID:    receipt.TraderID,
		Type:        receipt.Type,
		RefType:     receipt.RefType,
		ShopID:      receipt.ShopID,
	}
	if receipt.Status == status3.Z {
		if args.TraderID.Valid && args.TraderID.ID != receipt.TraderID {
			receiptNeedValidate.TraderID = args.TraderID.ID
		}
		receiptNeedValidate.RefType = args.RefType
		receiptNeedValidate.Amount = args.Amount.Int
		receiptNeedValidate.Lines = args.Lines
		receiptNeedValidate.PaidAt = args.PaidAt
	}
	if err := a.validateReceiptForCreateOrUpdate(ctx, args.ShopID, receiptNeedValidate); err != nil {
		return nil, err
	}

	if receipt.Status != status3.Z {
		args.TraderID = WrapID(receipt.TraderID)
		args.Amount = Int(receipt.Amount)
		args.RefType = receipt.RefType
		args.Lines = receipt.Lines
		args.Trader = receipt.Trader
		args.PaidAt = receipt.PaidAt
	} else {
		if !args.TraderID.Valid || args.TraderID.ID == receipt.TraderID {
			args.TraderID = PID(&receipt.TraderID)
			args.RefType = receipt.RefType
			args.Trader = receipt.Trader
		} else {
			args.Trader = receiptNeedValidate.Trader
		}
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

func (a *ReceiptAggregate) validateReceiptForCreateOrUpdate(ctx context.Context, shopID dot.ID, receipt *receipting.Receipt) error {
	if receipt.ID == 0 && receipt.Title == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiêu đề không hợp lệ")
	}
	if receipt.ID == 0 && receipt.Amount <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị phiếu phải lớn hơn 0")
	}
	if receipt.Amount > 0 && len(receipt.Lines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị phiếu không hợp lệ")
	}

	// Validate type and refID
	if err := a.validateTypeAndRefType(receipt.Type, receipt.RefType); err != nil {
		return err
	}

	// Validate TraderID
	if receipt.TraderID != 0 {
		if err := a.validateAndFillTrader(ctx, shopID, receipt); err != nil {
			return err
		}
	}

	// Validate ledger
	if receipt.LedgerID != 0 {
		if err := a.validateLedger(ctx, receipt.LedgerID, shopID); err != nil {
			return err
		}
	}

	// Validate receipt lines
	if receipt.Lines != nil && len(receipt.Lines) > 0 {
		if err := a.validateReceiptLines(ctx, receipt.RefType, receipt); err != nil {
			return err
		}
	}

	return nil
}

func (a *ReceiptAggregate) validateTypeAndRefType(receiptType receipting.ReceiptType, receiptRefType receipting.ReceiptRefType) error {
	if receiptRefType == receipting.ReceiptRefTypePurchaseOrder && receiptType == receipting.ReceiptTypeReceipt {
		return cm.Errorf(cm.InvalidArgument, nil, "Loại phiếu không hợp lệ")
	}
	if receiptType == receipting.ReceiptTypePayment && receiptRefType == receipting.ReceiptRefTypeOrder {
		return cm.Errorf(cm.InvalidArgument, nil, "Loại phiếu không hợp lệ")
	}
	return nil
}

func (a *ReceiptAggregate) validateLedger(ctx context.Context, ledgerID, shopID dot.ID) error {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     ledgerID,
		ShopID: shopID,
	}
	if err := a.ledgerQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Map(cm.NotFound, cm.FailedPrecondition, "Tài khoản thanh toán không tồn tại").
			Throw()
	}
	return nil
}

func (a *ReceiptAggregate) validateAndFillTrader(ctx context.Context, shopID dot.ID, receipt *receipting.Receipt) error {
	query := &tradering.GetTraderInfoByIDQuery{
		ID:     receipt.TraderID,
		ShopID: shopID,
	}
	if err := a.traderQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
			Throw()
	}
	receipt.Trader = &receipting.Trader{
		ID:       query.Result.ID,
		Type:     query.Result.Type,
		FullName: query.Result.FullName,
		Phone:    query.Result.Phone,
	}
	traderType := query.Result.Type
	switch receipt.RefType {
	case receipting.ReceiptRefTypeOrder:
		if traderType != tradering.CustomerType {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đối tác không hợp lệ")
		}
	case receipting.ReceiptRefTypeFulfillment:
		if traderType != tradering.CarrierType {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đối tác không hợp lệ")
		}
	case receipting.ReceiptRefTypePurchaseOrder:
		if traderType != tradering.SupplierType {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đối tác không hợp lệ")
		}
	}
	return nil
}

func calcReceiptLinesTotalAmount(receipt *receipting.Receipt) (totalAmount int, refIDs []dot.ID, mapRefIDAmount map[dot.ID]int, err error) {
	// Map of [ ref_id ] amount of line
	mapRefIDAmount = make(map[dot.ID]int)
	for _, receiptLine := range receipt.Lines {
		// check amount of a receiptLine < 0
		if receiptLine.Amount <= 0 {
			err = cm.Errorf(cm.FailedPrecondition, nil, "Giá trị mỗi hàng phải lớn hơn 0")
			return
		}
		totalAmount += receiptLine.Amount

		if receiptLine.RefID == 0 {
			continue
		}

		if _, ok := mapRefIDAmount[receiptLine.RefID]; ok {
			err = cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d trùng nhau trong phiếu", receiptLine.RefID)
			return
		}

		mapRefIDAmount[receiptLine.RefID] = receiptLine.Amount
		refIDs = append(refIDs, receiptLine.RefID)
	}
	return
}

func (a *ReceiptAggregate) validateReceiptLines(
	ctx context.Context, refType receipting.ReceiptRefType, receipt *receipting.Receipt,
) error {
	totalAmountOfReceiptLines, refIDs, mapRefIDAmount, err := calcReceiptLinesTotalAmount(receipt)
	if err != nil {
		return err
	}
	if totalAmountOfReceiptLines != receipt.Amount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Amount of receipt must be equal to total amount of receiptLines")
	}
	if len(refIDs) == 0 {
		return nil
	}

	event := &receipting.ReceiptCreatingEvent{
		RefIDs:         refIDs,
		MapRefIDAmount: mapRefIDAmount,
		Receipt:        receipt,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}

func (a *ReceiptAggregate) DeleteReceipt(
	ctx context.Context, id dot.ID, shopID dot.ID,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(id).ShopID(shopID).SoftDelete()
	return deleted, err
}

func (a *ReceiptAggregate) CancelReceipt(
	ctx context.Context, args *receipting.CancelReceiptArgs,
) (updated int, err error) {
	receipt, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetReceipt()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}

	if receipt.Status == status3.N {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	}

	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).CancelReceipt(args.Reason)
	if err != nil {
		return 0, err
	}

	receiptCancelledEvent := &receipting.ReceiptCancelledEvent{
		ShopID:    args.ShopID,
		ReceiptID: args.ID,
	}
	if err := a.eventBus.Publish(ctx, receiptCancelledEvent); err != nil {
		ll.Error("receiptCancelledEvent", l.Error(err))
	}

	return updated, err
}

func (a *ReceiptAggregate) ConfirmReceipt(
	ctx context.Context, args *receipting.ConfirmReceiptArgs,
) (updated int, err error) {
	receipt, err := a.store(ctx).ShopID(args.ShopID).ID(args.ID).GetReceipt()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}

	if err := a.validateReceiptLines(ctx, receipt.RefType, receipt); err != nil {
		return 0, err
	}

	switch receipt.Status {
	case status3.P:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã xác nhận")
	case status3.N:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	default:
		//no-op
	}

	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).ConfirmReceipt()
	if err != nil {
		return 0, err
	}

	receiptConfirmedEvent := &receipting.ReceiptConfirmedEvent{
		ShopID:    args.ShopID,
		ReceiptID: args.ID,
	}
	if err := a.eventBus.Publish(ctx, receiptConfirmedEvent); err != nil {
		ll.Error("receiptConfirmedEvent", l.Error(err))
	}

	return updated, err
}
