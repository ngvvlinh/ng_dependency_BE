package aggregate

import (
	"context"
	"time"

	"o.o/api/main/ledgering"
	"o.o/api/main/ordering"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/receipting"
	"o.o/api/shopping/carrying"
	"o.o/api/shopping/customering"
	"o.o/api/shopping/suppliering"
	"o.o/api/shopping/tradering"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/receipting/convert"
	"o.o/backend/com/main/receipting/model"
	"o.o/backend/com/main/receipting/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
	. "o.o/capi/dot"
	"o.o/common/l"
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
	database com.MainDB, eventBus capi.EventBus,
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

func ReceiptAggregateMessageBus(a *ReceiptAggregate) receipting.CommandBus {
	b := bus.New()
	return receipting.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ReceiptAggregate) CreateReceipt(
	ctx context.Context, args *receipting.CreateReceiptArgs,
) (*receipting.Receipt, error) {
	if args.PaidAt.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ng??y t???o phi???u kh??ng h???p l???")
	}
	if args.TraderID == 0 && args.RefType != receipt_ref.Order {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "?????i t??c kh??ng h???p l???")
	}
	if args.LedgerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "T??i kho???n thanh to??n kh??ng h???p l???")
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui l??ng nh???p m??")
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
		if err = a.eventBus.Publish(ctx, receiptConfirmEvent); err != nil {
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
			Wrap(cm.NotFound, "kh??ng t??m th???y phi???u").
			Throw()
	}

	if receipt.Status == status3.N {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Kh??ng th??? thay ?????i phi???u ???? h???y.")
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
		receiptNeedValidate.RefType = args.RefType.Apply(receipt.RefType)
		receiptNeedValidate.Amount = args.Amount.Int
		receiptNeedValidate.Lines = args.Lines
		receiptNeedValidate.PaidAt = args.PaidAt
	}
	if err = a.validateReceiptForCreateOrUpdate(ctx, args.ShopID, receiptNeedValidate); err != nil {
		return nil, err
	}

	if receipt.Status != status3.Z {
		args.TraderID = WrapID(receipt.TraderID)
		args.Amount = Int(receipt.Amount)
		args.Lines = receipt.Lines
		args.Trader = receipt.Trader
		args.PaidAt = receipt.PaidAt
	} else {
		if !args.TraderID.Valid || args.TraderID.ID == receipt.TraderID {
			args.TraderID = PID(&receipt.TraderID)
			args.Trader = receipt.Trader
		} else {
			args.Trader = receiptNeedValidate.Trader
		}
	}

	if err = scheme.Convert(args, receipt); err != nil {
		return nil, err
	}

	receiptDB := new(model.Receipt)
	if err = scheme.Convert(receipt, receiptDB); err != nil {
		return nil, err
	}
	receiptDB.Lines = convert.Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(receipt.Lines)
	err = a.store(ctx).UpdateReceiptDB(receiptDB)
	return receipt, err
}

func (a *ReceiptAggregate) validateReceiptForCreateOrUpdate(ctx context.Context, shopID dot.ID, receipt *receipting.Receipt) error {
	if receipt.ID == 0 && receipt.Title == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Ti??u ????? kh??ng h???p l???")
	}
	if receipt.ID == 0 && receipt.Amount <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Gi?? tr??? phi???u ph???i l???n h??n 0")
	}
	if receipt.Amount > 0 && len(receipt.Lines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Gi?? tr??? phi???u kh??ng h???p l???")
	}

	// Validate type and refID
	if err := validateTypeAndRefType(receipt.Type, receipt.RefType); err != nil {
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

func validateTypeAndRefType(receiptType receipt_type.ReceiptType, receiptRefType receipt_ref.ReceiptRef) error {
	// receipt ref can be empty (""), in this case, we don't care receipt type
	if receiptRefType == receipt_ref.None {
		return nil
	}
	type tuple struct {
		receiptType    receipt_type.ReceiptType
		receiptRefType receipt_ref.ReceiptRef
	}
	t := tuple{receiptType, receiptRefType}
	switch t {
	case
		tuple{receipt_type.Payment, receipt_ref.PurchaseOrder},
		tuple{receipt_type.Payment, receipt_ref.Refund},
		tuple{receipt_type.Receipt, receipt_ref.Order},
		tuple{receipt_type.Payment, receipt_ref.Order},
		tuple{receipt_type.Payment, receipt_ref.Fulfillment}:
		return nil
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Lo???i phi???u kh??ng h???p l???")
	}
}

func (a *ReceiptAggregate) validateLedger(ctx context.Context, ledgerID, shopID dot.ID) error {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     ledgerID,
		ShopID: shopID,
	}
	if err := a.ledgerQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Map(cm.NotFound, cm.FailedPrecondition, "T??i kho???n thanh to??n kh??ng t???n t???i").
			Throw()
	}
	return nil
}

func (a *ReceiptAggregate) validateAndFillTrader(ctx context.Context, shopID dot.ID, receipt *receipting.Receipt) error {
	if receipt.RefType == receipt_ref.Order && receipt.TraderID == 0 {
		return nil
	}
	query := &tradering.GetTraderInfoByIDQuery{
		ID:     receipt.TraderID,
		ShopID: shopID,
	}
	if err := a.traderQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Map(cm.NotFound, cm.FailedPrecondition, "?????i t??c kh??ng h???p l???").
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
	case receipt_ref.Order:
		if traderType != tradering.CustomerType {
			return cm.Errorf(cm.FailedPrecondition, nil, "?????i t??c kh??ng h???p l???")
		}
	case receipt_ref.Fulfillment:
		if traderType != tradering.CarrierType {
			return cm.Errorf(cm.FailedPrecondition, nil, "?????i t??c kh??ng h???p l???")
		}
	case receipt_ref.PurchaseOrder:
		if traderType != tradering.SupplierType {
			return cm.Errorf(cm.FailedPrecondition, nil, "?????i t??c kh??ng h???p l???")
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
			err = cm.Errorf(cm.FailedPrecondition, nil, "Gi?? tr??? m???i h??ng ph???i l???n h??n 0").WithMetap("receipt_line_ref_id", receiptLine.RefID)
			return
		}
		totalAmount += receiptLine.Amount

		if receiptLine.RefID == 0 {
			continue
		}

		if _, ok := mapRefIDAmount[receiptLine.RefID]; ok {
			err = cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d tr??ng nhau trong phi???u", receiptLine.RefID)
			return
		}

		mapRefIDAmount[receiptLine.RefID] = receiptLine.Amount
		refIDs = append(refIDs, receiptLine.RefID)
	}
	return
}

func (a *ReceiptAggregate) validateReceiptLines(
	ctx context.Context, refType receipt_ref.ReceiptRef, receipt *receipting.Receipt,
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
	if err = a.eventBus.Publish(ctx, event); err != nil {
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
			Wrap(cm.NotFound, "Kh??ng t??m th???y phi???u").
			Throw()
	}

	if receipt.Status == status3.N {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phi???u ???? h???y")
	}
	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).CancelReceipt(args.CancelReason)
	if err != nil {
		return 0, err
	}

	receiptCancelledEvent := &receipting.ReceiptCancelledEvent{
		ShopID:    args.ShopID,
		ReceiptID: args.ID,
	}
	if err = a.eventBus.Publish(ctx, receiptCancelledEvent); err != nil {
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
			Wrap(cm.NotFound, "Kh??ng t??m th???y phi???u").
			Throw()
	}
	//receipt confirming
	receiptConfirmingEvent := &receipting.ReceiptConfirmingEvent{
		ShopID:      args.ShopID,
		ReceiptID:   args.ID,
		ReceiptType: receipt.Type,
		RefType:     receipt.RefType,
	}
	if err = a.eventBus.Publish(ctx, receiptConfirmingEvent); err != nil {
		return 0, err
	}

	if err = a.validateReceiptLines(ctx, receipt.RefType, receipt); err != nil {
		return 0, err
	}

	switch receipt.Status {
	case status3.P:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phi???u ???? x??c nh???n")
	case status3.N:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phi???u ???? h???y")
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
	if err = a.eventBus.Publish(ctx, receiptConfirmedEvent); err != nil {
		ll.Error("receiptConfirmedEvent", l.Error(err))
	}

	return updated, err
}

func (a *ReceiptAggregate) CancelReceiptByRefID(ctx context.Context, args *receipting.CancelReceiptByRefIDRequest) error {
	switch args.RefType {
	case receipt_ref.Order:
		return a.CancelReceiptByOrderID(ctx, args)
	default:
		panic("implement me")
	}
}

func (a *ReceiptAggregate) CancelReceiptByOrderID(ctx context.Context, args *receipting.CancelReceiptByRefIDRequest) error {
	// - link notion: https://www.notion.so/Cancel-Task-Receipt-70f76208f27c437b9ae650037e08e283
	// - M???c ????ch: ph???c v??? t???ng h???p, th???ng k??
	// - M?? t??? s??:
	// 	+ Tr?????ng h???p receipt ch??a confirm => b??? qua, cho kh??ch h??ng t??? ch???nh s???a
	// 	+ Tr?????ng h???p receipt ???? confirm:
	// 		* Tr?????ng h???p receipt ch??? ch???a 1 ????n h??ng => h???y receipt
	// 		* Tr?????ng h???p receipt g???m nhi???u ????n => t???o m???t phi???u ?????i l???p
	query := a.store(ctx).ShopID(args.ShopID).RefsID(args.RefID).RefType(args.RefType).Status(status3.P)
	receipts, err := query.ListReceipts()
	if err != nil {
		return err
	}
	queryCustomerDefault := &customering.GetCustomerIndependentQuery{}
	err = a.customerQuery.Dispatch(ctx, queryCustomerDefault)
	if err != nil {
		return err
	}
	// check order was canceled
	for _, receipt := range receipts {
		if receipt.Type == receipt_type.Payment {
			return cm.Errorf(cm.InvalidArgument, nil, "Phi???u thu c???a order ???? ???????c h???y b???")
		}
	}

	for _, receipt := range receipts {
		if receipt.Status != status3.P {
			continue
		}
		count := 0
		// check receipt have only one order or not
		for _, refID := range receipt.RefIDs {
			if refID != 0 {
				count++
			}
			if count > 1 {
				break
			}
		}
		// receipt only have one line
		if count == 1 {
			_, err = a.CancelReceipt(ctx, &receipting.CancelReceiptArgs{
				ID:           receipt.ID,
				ShopID:       receipt.ShopID,
				CancelReason: "Cancel order",
			})
			if err != nil {
				return err
			}
			continue
		}

		var customer *customering.ShopCustomer
		queryCustomer := &customering.GetCustomerByIDQuery{
			ID:     receipt.TraderID,
			ShopID: receipt.ShopID,
		}
		err = a.customerQuery.Dispatch(ctx, queryCustomer)
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}
		if err != nil && cm.ErrorCode(err) == cm.NotFound {
			customer = queryCustomerDefault.Result
		} else {
			customer = queryCustomer.Result
		}

		var newLines = []*receipting.ReceiptLine{}
		var lineOrderTarget = &receipting.ReceiptLine{}
		for _, line := range receipt.Lines {
			if line.RefID == args.RefID {
				lineOrderTarget = line
				break
			}
		}
		newLines = append(newLines, lineOrderTarget)

		receipCreate, err := a.CreateReceipt(ctx, &receipting.CreateReceiptArgs{
			ShopID:      receipt.ShopID,
			TraderID:    customer.ID,
			Title:       "Ho??n tr??? khi h???y ????n h??ng",
			Type:        receipt_type.Payment,
			Status:      int(status3.Z),
			Description: "T??? ?????ng ",
			Amount:      lineOrderTarget.Amount,
			LedgerID:    receipt.LedgerID,
			RefIDs:      []dot.ID{args.RefID},
			RefType:     receipt_ref.Order,
			Lines:       newLines,
			PaidAt:      time.Now(),
			Trader: &receipting.Trader{
				ID:       customer.ID,
				Type:     "customer",
				FullName: customer.FullName,
				Phone:    customer.Phone,
			},
			CreatedBy: args.UpdatedBy,
			Mode:      receipt_mode.Auto,
		})
		if err != nil {
			return err
		}
		_, err = a.ConfirmReceipt(ctx, &receipting.ConfirmReceiptArgs{
			ID:     receipCreate.ID,
			ShopID: receipCreate.ShopID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
