package aggregate

import (
	"context"

	"etop.vn/api/shopping/customering"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/main/receipting/convert"
	"etop.vn/backend/com/main/receipting/model"
	"etop.vn/backend/com/main/receipting/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/capi"
	. "etop.vn/capi/dot"
)

var _ receipting.Aggregate = &ReceiptAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ReceiptAggregate struct {
	store         sqlstore.ReceiptStoreFactory
	eventBus      capi.EventBus
	traderQuery   tradering.QueryBus
	ledgerQuery   ledgering.QueryBus
	orderQuery    ordering.QueryBus
	customerQuery customering.QueryBus
}

func NewReceiptAggregate(
	db *cmsql.Database, eventBus capi.EventBus,
	traderQuery tradering.QueryBus, ledgerQuery ledgering.QueryBus,
	orderQuery ordering.QueryBus, customerQuery customering.QueryBus,
) *ReceiptAggregate {
	return &ReceiptAggregate{
		store:         sqlstore.NewReceiptStore(db),
		eventBus:      eventBus,
		traderQuery:   traderQuery,
		ledgerQuery:   ledgerQuery,
		orderQuery:    orderQuery,
		customerQuery: customerQuery,
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
	if args.Type != receipting.ReceiptTypeReceipt && args.Type != receipting.ReceiptTypePayment {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Kiểu phiếu không hợp lệ")
	}

	receiptNeedValidate := &receipting.Receipt{
		TraderID:    args.TraderID,
		ShopID:      args.ShopID,
		CreatedBy:   args.CreatedBy,
		Title:       args.Title,
		Type:        args.Type,
		Description: args.Description,
		Amount:      args.Amount,
		LedgerID:    args.LedgerID,
		PaidAt:      args.PaidAt,
		Lines:       args.Lines,
	}
	if err := a.validateReceiptForCreateOrUpdate(ctx, args.ShopID, receiptNeedValidate); err != nil {
		return nil, err
	}

	receipt := new(receipting.Receipt)
	if err := scheme.Convert(args, receipt); err != nil {
		return nil, err
	}

	var maxCodeNorm int32
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
	receipt.Code = convert.GenerateCode(int(codeNorm))
	receipt.CodeNorm = codeNorm

	err = a.store(ctx).CreateReceipt(receipt)
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

	if receipt.Status == int32(etopmodel.S3Negative) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể thay đổi phiếu đã hủy.")
	}

	receiptNeedValidate := &receipting.Receipt{
		ID:          args.ID,
		Title:       args.Title.String,
		Description: args.Description.String,
		LedgerID:    args.LedgerID.Int64,
	}
	if receipt.Status == int32(etopmodel.S3Zero) {
		if args.TraderID.Int64 != receipt.TraderID {
			receiptNeedValidate.TraderID = args.TraderID.Int64
		}
		receiptNeedValidate.Amount = args.Amount.Int32
		receiptNeedValidate.Lines = args.Lines
		receiptNeedValidate.PaidAt = args.PaidAt
	}
	if err := a.validateReceiptForCreateOrUpdate(ctx, args.ShopID, receiptNeedValidate); err != nil {
		return nil, err
	}

	if receipt.Status != int32(etopmodel.S3Zero) {
		args.TraderID = PInt64(&receipt.TraderID)
		args.Amount = PInt32(&receipt.Amount)
		args.Lines = receipt.Lines
		args.PaidAt = receipt.PaidAt
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

func (a *ReceiptAggregate) validateReceiptForCreateOrUpdate(ctx context.Context, shopID int64, receipt *receipting.Receipt) error {
	if receipt.ID == 0 && receipt.Title == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiêu đề không hợp lệ")
	}
	if receipt.ID == 0 && receipt.Amount <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị phiếu phải lớn hơn 0")
	}

	if receipt.Amount > 0 && len(receipt.Lines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị phiếu không hợp lệ")
	}

	var traderType string
	// validate TraderID
	if receipt.TraderID != 0 {
		switch receipt.TraderID {
		case etopmodel.TopShipID:
			traderType = tradering.CarrierType
		case etopmodel.IndependentCustomerID:
			traderType = tradering.CustomerType
		default:
			query := &tradering.GetTraderByIDQuery{
				ID:     receipt.TraderID,
				ShopID: shopID,
			}
			if err := a.traderQuery.Dispatch(ctx, query); err != nil {
				return cm.MapError(err).
					Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
					Throw()
			}
			traderType = query.Result.Type

			switch traderType {
			case tradering.CustomerType:
				query := &customering.GetCustomerByIDQuery{
					ID:     receipt.TraderID,
					ShopID: shopID,
				}
				err := a.customerQuery.Dispatch(ctx, query)
				if err != nil {
					return cm.MapError(err).
						Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
						Throw()
				}
			case tradering.CarrierType, tradering.VendorType:
				// TODO:
			}
		}
	}

	// Validate ledger
	if receipt.LedgerID != 0 {
		query := &ledgering.GetLedgerByIDQuery{
			ID:     receipt.LedgerID,
			ShopID: shopID,
		}
		if err := a.ledgerQuery.Dispatch(ctx, query); err != nil {
			return cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Tài khoản thanh toán không tồn tại").
				Throw()
		}
	}

	// validate receipt lines
	if receipt.Lines != nil && len(receipt.Lines) > 0 {
		if err := a.validateReceiptLines(ctx, traderType, receipt); err != nil {
			return err
		}
	}

	return nil
}

func (a *ReceiptAggregate) validateReceiptLines(ctx context.Context, traderType string, receipt *receipting.Receipt) error {
	totalAmountOfReceiptLines, refIDs, mapOrdersAmount, err := calcReceiptLinesTotalAmount(receipt)
	if err != nil {
		return err
	}
	if totalAmountOfReceiptLines != receipt.Amount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Amount of receipt must be equal to total amount of receiptLines")
	}

	if len(refIDs) == 0 {
		return nil
	}

	// List all erf in orderIDs
	mOrders := make(map[int64]*ordering.Order)
	var ordersTemp []*ordering.Order

	switch traderType {
	case tradering.CustomerType:
		query := &ordering.GetOrdersQuery{
			ShopID: receipt.ShopID,
			IDs:    refIDs,
		}
		if err := a.orderQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		ordersTemp = query.Result.Orders
	case tradering.CarrierType, tradering.VendorType:
		query := &ordering.GetOrdersQuery{
			ShopID: receipt.ShopID,
			IDs:    refIDs,
		}
		if err := a.orderQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		ordersTemp = query.Result.Orders
	}
	for _, order := range ordersTemp {
		mOrders[order.ID] = order
	}

	// Check orderIds with orderIds of listOrdersQuery.Result
	// When different len
	if len(refIDs) != len(mOrders) {
		for _, v := range refIDs {
			if _, ok := mOrders[v]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d không tìm thấy", v)
			}
		}
	}

	// Don't check received_amount of receipt type payment
	if receipt.Type == receipting.ReceiptTypePayment {
		return nil
	}

	// List all receipts IN orderIDs
	receipts, err := a.store(ctx).ShopID(receipt.ShopID).RefIDs(refIDs...).Status(etopmodel.S3Positive).ListReceipts()
	if err != nil {
		return err
	}

	// Get total amount each of orderId in orderIDs
	// Map of [ orderId ] amount of receiptLines (old receipts)
	mapOrdersAmountOld := make(map[int64]int32)
	for _, receiptElem := range receipts {
		// Ignore current receipt when updating
		if receiptElem.ID == receipt.ID {
			continue
		}
		for _, receiptLine := range receiptElem.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, has := mapOrdersAmount[receiptLine.RefID]; has {
				switch receipt.Type {
				case receipting.ReceiptTypeReceipt:
					mapOrdersAmountOld[receiptLine.RefID] += receiptLine.Amount
				case receipting.ReceiptTypePayment:
					mapOrdersAmountOld[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}

	// Check each of amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapOrdersAmount {
		if value > int32(mOrders[key].TotalAmount)-mapOrdersAmountOld[key] {
			return cm.Errorf(cm.InvalidArgument, nil, "Giá trị của đơn hàng không hợp lệ, Vui lòng tải lại trang và thử lại")
		}
	}

	return nil
}

func calcReceiptLinesTotalAmount(receipt *receipting.Receipt) (totalAmount int32, orderIDs []int64, mapOrdersAmount map[int64]int32, err error) {
	// Map of [ orderId ] amount of receiptLines (params)
	mapOrdersAmount = make(map[int64]int32)
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

		// Check has key in map
		// hasKey = true -> duplicate orderId in receipt
		// hasKey = false -> add orderId in map
		if _, has := mapOrdersAmount[receiptLine.RefID]; has {
			err = cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d trùng nhau trong phiếu", receiptLine.RefID)
			return
		}

		mapOrdersAmount[receiptLine.RefID] = receiptLine.Amount
		orderIDs = append(orderIDs, receiptLine.RefID)
	}
	return
}

func (a *ReceiptAggregate) DeleteReceipt(
	ctx context.Context, id int64, shopID int64,
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

	if receipt.Status == int32(etopmodel.S3Negative) {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	}

	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).CancelReceipt(args.Reason)
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

	getTraderByIDQuery := &tradering.GetTraderByIDQuery{
		ShopID: args.ShopID,
		ID:     receipt.TraderID,
	}
	if err := a.traderQuery.Dispatch(ctx, getTraderByIDQuery); err != nil {
		return 0, err
	}

	if err := a.validateReceiptLines(ctx, getTraderByIDQuery.Result.Type, receipt); err != nil {
		return 0, err
	}

	switch receipt.Status {
	case int32(etopmodel.S3Positive):
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã xác nhận")
	case int32(etopmodel.S3Negative):
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	default:
		//no-op
	}

	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).ConfirmReceipt()
	return updated, err
}
