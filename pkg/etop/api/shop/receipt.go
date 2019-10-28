package shop

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/shopping/vendoring"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		receiptService.CreateReceipt,
		receiptService.UpdateReceipt,
		receiptService.DeleteReceipt,
		receiptService.ConfirmReceipt,
		receiptService.CancelReceipt,
		receiptService.GetReceipt,
		receiptService.GetReceipts)
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (_err error) {
	key := fmt.Sprintf("Create receipt %v-%v-%v-%v-%v-%v-%v",
		q.Context.Shop.ID, q.Context.UserID, q.TraderId, q.Title, q.Description, q.Amount, q.Type)
	result, err := idempgroup.DoAndWrap(
		key, 15*time.Second,
		func() (interface{}, error) { return s.createReceipt(ctx, q) },
		"Create receipt")
	if err != nil {
		return err
	}
	q.Result = pbshop.PbReceipt(result.(*receipting.CreateReceiptCommand).Result)
	return nil
}

func (s *ReceiptService) createReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (*receipting.CreateReceiptCommand, error) {
	if q.PaidAt == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ngày tạo phiếu không hợp lệ")
	}
	if q.TraderId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Đối tác không hợp lệ")
	}
	if q.LedgerId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Sổ quỹ không hợp lệ")
	}
	if q.Type != string(receipting.ReceiptTypeReceipt) && q.Type != string(receipting.ReceiptTypePayment) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Loại đối tác không hợp lệ")
	}

	paidAt, err := ptypes.Timestamp(q.PaidAt)
	if err != nil {
		return nil, err
	}
	receipt := &receipting.Receipt{
		TraderID:    q.TraderId,
		ShopID:      q.Context.Shop.ID,
		CreatedBy:   q.Context.UserID,
		Title:       q.Title,
		Type:        q.Type,
		Description: q.Description,
		Amount:      q.Amount,
		LedgerID:    q.LedgerId,
		PaidAt:      paidAt,
		Lines:       pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
	}

	//if receipt.PaidAt == nil {
	//	return nil, cm.Errorf(cm.InvalidArgument, nil, "Ngày tạo phiếu không hợp lệ")
	//}
	if err := s.validateReceiptForCreateOrUpdate(ctx, q.Context.Shop.ID, receipt); err != nil {
		return nil, err
	}

	cmd := &receipting.CreateReceiptCommand{
		ShopID:      q.Context.Shop.ID,
		CreatedBy:   q.Context.UserID,
		TraderID:    q.TraderId,
		Title:       q.Title,
		Description: q.Description,
		Amount:      q.Amount,
		LedgerID:    q.LedgerId,
		PaidAt:      paidAt,
		Type:        receipting.ReceiptType(q.Type),
		CreatedType: string(receipting.ReceiptCreatedTypeManual),
		Status:      int32(model.S3Zero),
		Lines:       receipt.Lines,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (s *ReceiptService) UpdateReceipt(ctx context.Context, q *wrapshop.UpdateReceiptEndpoint) error {
	query := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}

	if query.Result.Status == int32(model.S3Negative) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể thay đổi phiếu đã hủy.")
	}

	paidAt, err := ptypes.Timestamp(&q.PaidAt)
	if err != nil {
		return err
	}

	lines := pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines)
	receipt := &receipting.Receipt{
		ID:          q.Id,
		Title:       PString(q.Title).Apply(""),
		Description: PString(q.Description).Apply(""),
		LedgerID:    PInt64(q.LedgerId).Apply(0),
	}
	if query.Result.Status == int32(model.S3Zero) {
		receipt.TraderID = PInt64(q.TraderId).Apply(0)
		receipt.Amount = PInt32(q.Amount).Apply(0)
		receipt.Lines = lines
		receipt.PaidAt = paidAt
	}

	if err := s.validateReceiptForCreateOrUpdate(ctx, q.Context.Shop.ID, receipt); err != nil {
		return err
	}

	cmd := &receipting.UpdateReceiptCommand{
		ID:          q.Id,
		ShopID:      q.Context.Shop.ID,
		Title:       PString(q.Title),
		Description: PString(q.Description),
		LedgerID:    PInt64(q.LedgerId),
	}
	if query.Result.Status == int32(model.S3Zero) {
		cmd.TraderID = PInt64(q.TraderId)
		cmd.Amount = PInt32(q.Amount)
		cmd.Lines = lines
		cmd.PaidAt = paidAt
	}
	err = receiptAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}

	q.Result = pbshop.PbReceipt(cmd.Result)
	return nil
}

func (s *ReceiptService) validateReceiptForCreateOrUpdate(ctx context.Context, shopID int64, receipt *receipting.Receipt) error {
	if receipt.ID == 0 && receipt.Title == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiêu đề không hợp lệ")
	}
	if receipt.ID == 0 && receipt.Amount <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Amount must rather than 0")
	}

	if receipt.Amount > 0 && len(receipt.Lines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Amount is invalid")
	}

	// validate that the updated code does not exist
	if receipt.Code != "" {
		query := &receipting.GetReceiptByCodeQuery{
			Code:   receipt.Code,
			ShopID: shopID,
		}

		err := receiptQuery.Dispatch(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.NoError:
			if query.Result.ID != receipt.ID {
				return cm.Errorf(cm.FailedPrecondition, nil, "Mã phiếu %v đã tồn tại. Vui lòng chọn mã khác.", receipt.Code)
			}
		case cm.NotFound:
			// no-op
		default:
			return err
		}
	}

	var traderType string
	// validate TraderID
	if receipt.TraderID != 0 {
		query := &tradering.GetTraderByIDQuery{
			ID:     receipt.TraderID,
			ShopID: shopID,
		}
		if err := traderQuery.Dispatch(ctx, query); err != nil {
			return cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
				Throw()
		}
		traderType = query.Result.Type
	}

	// Validate ledger
	if receipt.LedgerID != 0 {
		query := &ledgering.GetLedgerByIDQuery{
			ID:     receipt.LedgerID,
			ShopID: shopID,
		}
		if err := ledgerQuery.Dispatch(ctx, query); err != nil {
			return cm.MapError(err).
				Map(cm.NotFound, cm.FailedPrecondition, "Sổ quỹ không hợp lệ").
				Throw()
		}
	}

	// validate receipt lines
	if receipt.Lines != nil && len(receipt.Lines) > 0 {
		if err := s.validateReceiptLines(ctx, traderType, receipt); err != nil {
			return err
		}
	}

	return nil
}

func (s *ReceiptService) validateReceiptLines(ctx context.Context, traderType string, receipt *receipting.Receipt) error {
	totalAmountOfReceiptLines, orderIDs, mapOrdersAmount, err := calcReceiptLinesTotalAmount(receipt)
	if err != nil {
		return err
	}
	if totalAmountOfReceiptLines != receipt.Amount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Amount of receipt must be equal to total amount of receiptLines")
	}

	if len(orderIDs) == 0 {
		return nil
	}

	// List all orders in orderIDs
	mOrders := make(map[int64]*ordering.Order)
	var ordersTemp []*ordering.Order

	switch traderType {
	case tradering.CustomerType:
		listOrdersQuery := &ordering.GetOrdersByIDsAndCustomerIDQuery{
			ShopID:     receipt.ShopID,
			IDs:        orderIDs,
			CustomerID: receipt.TraderID,
		}
		if err := orderQuery.Dispatch(ctx, listOrdersQuery); err != nil {
			return err
		}
		ordersTemp = listOrdersQuery.Result.Orders
	case tradering.CarrierType, tradering.VendorType:
		listOrdersQuery := &ordering.GetOrdersQuery{
			ShopID: receipt.ShopID,
			IDs:    orderIDs,
		}
		if err := orderQuery.Dispatch(ctx, listOrdersQuery); err != nil {
			return err
		}
		ordersTemp = listOrdersQuery.Result.Orders
	}
	for _, order := range ordersTemp {
		mOrders[order.ID] = order
	}

	// Check orderIds with orderIds of listOrdersQuery.Result
	// When different len
	if len(orderIDs) != len(mOrders) {
		for _, v := range orderIDs {
			if _, ok := mOrders[v]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "OrderID %d not found", v)
			}
		}
	}

	// List all receipts IN orderIDs
	listReceiptsQuery := &receipting.ListReceiptsByRefIDsQuery{
		IDs:    orderIDs,
		ShopID: receipt.ShopID,
	}
	if err := receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}

	// Get total amount each of orderId in orderIDs
	// Map of [ orderId ] amount of receiptLines (old receipts)
	mapOrdersAmountOld := make(map[int64]int32)
	for _, receiptElem := range listReceiptsQuery.Result.Receipts {
		// Ignore current receipt when updating
		if receiptElem.ID == receipt.ID {
			continue
		}
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, has := mapOrdersAmount[receiptLine.RefID]; has {
				mapOrdersAmountOld[receiptLine.RefID] += receiptLine.Amount
			}
		}
	}

	// Check each of amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapOrdersAmount {
		if value >= int32(mOrders[key].TotalAmount)-mapOrdersAmountOld[key] {
			return cm.Errorf(cm.InvalidArgument, nil, "Amount of order_id %d is valid", key)
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
			err = cm.Errorf(cm.FailedPrecondition, nil, "Amount of receiptLine must be greater than 0")
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
			err = cm.Errorf(cm.FailedPrecondition, nil, "Duplicated OrderId %d in receipt", receiptLine.RefID)
			return
		}

		mapOrdersAmount[receiptLine.RefID] = receiptLine.Amount
		orderIDs = append(orderIDs, receiptLine.RefID)
	}
	return
}

func (s *ReceiptService) DeleteReceipt(ctx context.Context, q *wrapshop.DeleteReceiptEndpoint) error {
	cmd := &receipting.DeleteReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Receipt not found").
			Throw()
	}
	q.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func (s *ReceiptService) ConfirmReceipt(ctx context.Context, q *wrapshop.ConfirmReceiptEndpoint) error {
	query := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}
	switch query.Result.Status {
	case int32(model.S3Positive):
		return cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã xác nhận")
	case int32(model.S3Negative):
		return cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	default:
		//no-op
	}

	cmd := &receipting.ConfirmReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}

	return nil
}

func (s *ReceiptService) CancelReceipt(ctx context.Context, q *wrapshop.CancelReceiptEndpoint) error {
	query := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}
	switch query.Result.Status {
	case int32(model.S3Positive):
		return cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã xác nhận")
	case int32(model.S3Negative):
		return cm.Errorf(cm.FailedPrecondition, nil, "Phiếu đã hủy")
	default:
		// no-op
	}

	cmd := &receipting.CancelReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
		Reason: q.Reason,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}

	return nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, q *wrapshop.GetReceiptEndpoint) error {
	// Check receipt is exist
	getReceiptQuery := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, getReceiptQuery); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Receipt not found").
			Throw()
	}

	if receipts, err := s.getInfosForReceipts(ctx, q.Context.Shop.ID, []*receipting.Receipt{getReceiptQuery.Result}); err != nil {
		return err
	} else {
		q.Result = receipts[0]
	}

	return nil
}

func (s *ReceiptService) GetReceipts(ctx context.Context, q *wrapshop.GetReceiptsEndpoint) error {
	paging := q.Paging.CMPaging()
	listReceiptsQuery := &receipting.ListReceiptsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}

	if receipts, err := s.getInfosForReceipts(ctx, q.Context.Shop.ID, listReceiptsQuery.Result.Receipts); err != nil {
		return err
	} else {
		q.Result = &pbshop.ReceiptsResponse{
			Receipts: receipts,
			Paging:   pbcm.PbPageInfo(paging, listReceiptsQuery.Result.Count),
		}
	}

	return nil
}

func (s *ReceiptService) getInfosForReceipts(ctx context.Context, shopID int64, receipts []*receipting.Receipt) (receiptsResult []*pbshop.Receipt, _ error) {
	mapOrderIDAndReceivedAmount := make(map[int64]int32)
	mapLedger := make(map[int64]*ledgering.ShopLedger)
	var refIDs, userIDs, traderIDs, ledgerIDs []int64

	receiptsResult = pbshop.PbReceipts(receipts)

	// Get ref_ids into receiptLines
	for _, receipt := range receipts {
		userIDs = append(userIDs, receipt.CreatedBy)
		traderIDs = append(traderIDs, receipt.TraderID)
		ledgerIDs = append(ledgerIDs, receipt.LedgerID)

		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.RefID]; !ok {
				refIDs = append(refIDs, receiptLine.RefID)
				mapOrderIDAndReceivedAmount[receiptLine.RefID] = 0
			}
		}
	}

	// Get all users into receipts
	getUsersOfCurrAccount := &model.GetAccountUserExtendedsQuery{
		AccountIDs: []int64{shopID},
	}
	if err := bus.Dispatch(ctx, getUsersOfCurrAccount); err != nil {
		return nil, err
	}
	mapUserIDAndUser := make(map[int64]*model.User)
	for _, accountUser := range getUsersOfCurrAccount.Result.AccountUsers {
		mapUserIDAndUser[accountUser.User.ID] = accountUser.User
	}
	for _, receipt := range receiptsResult {
		receipt.User = pbetop.PbUser(mapUserIDAndUser[receipt.CreatedBy])
	}

	// List traders
	if err := listTraders(ctx, shopID, traderIDs, receiptsResult); err != nil {
		return nil, err
	}

	// List ledgers
	getLedgersByIDs := &ledgering.ListLedgersByIDsQuery{
		ShopID: shopID,
		IDs:    ledgerIDs,
	}
	if err := ledgerQuery.Dispatch(ctx, getLedgersByIDs); err != nil {
		return nil, err
	}
	for _, ledger := range getLedgersByIDs.Result.Ledger {
		mapLedger[ledger.ID] = ledger
	}
	for _, receipt := range receiptsResult {
		receipt.Ledger = pbshop.PbLedger(mapLedger[receipt.LedgerId])
	}
	return receiptsResult, nil
}

func listTraders(
	ctx context.Context, shopID int64,
	traderIDs []int64, receiptsResult []*pbshop.Receipt,
) error {
	mapVendor := make(map[int64]*vendoring.ShopVendor)
	var vendorIDs, customerIDs, carrierIDs []int64
	mapCustomer := make(map[int64]*customering.ShopCustomer)
	mapCarrier := make(map[int64]*carrying.ShopCarrier)
	getTradersByIDsQuery := &tradering.ListTradersByIDsQuery{
		ShopID: shopID,
		IDs:    traderIDs,
	}
	if err := traderQuery.Dispatch(ctx, getTradersByIDsQuery); err != nil {
		return err
	}
	for _, trader := range getTradersByIDsQuery.Result.Traders {
		switch trader.Type {
		case tradering.CarrierType:
			carrierIDs = append(carrierIDs, trader.ID)
		case tradering.CustomerType:
			customerIDs = append(customerIDs, trader.ID)
		case tradering.VendorType:
			vendorIDs = append(vendorIDs, trader.ID)
		}
	}
	// Get elements for each of type
	if vendorIDs != nil && len(vendorIDs) > 0 {
		query := &vendoring.ListVendorsByIDsQuery{
			ShopID: shopID,
			IDs:    vendorIDs,
		}
		if err := vendorQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, vendor := range query.Result.Vendors {
			mapVendor[vendor.ID] = vendor
		}
	}
	if customerIDs != nil && len(customerIDs) > 0 {
		query := &customering.ListCustomersByIDsQuery{
			ShopID: shopID,
			IDs:    customerIDs,
		}
		if err := customerQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, customer := range query.Result.Customers {
			mapCustomer[customer.ID] = customer
		}
	}
	if carrierIDs != nil && len(carrierIDs) > 0 {
		query := &carrying.ListCarriersByIDsQuery{
			ShopID: shopID,
			IDs:    carrierIDs,
		}
		if err := carrierQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, carrier := range query.Result.Carriers {
			mapCarrier[carrier.ID] = carrier
		}
	}
	for _, receipt := range receiptsResult {
		if vendor, ok := mapVendor[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       vendor.ID,
				Type:     tradering.VendorType,
				FullName: vendor.FullName,
			}
		}
		if customer, ok := mapCustomer[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       customer.ID,
				Type:     tradering.CustomerType,
				FullName: customer.FullName,
				Phone:    customer.Phone,
			}
		}
		if carrier, ok := mapCarrier[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       carrier.ID,
				Type:     tradering.CarrierType,
				FullName: carrier.FullName,
			}
		}
	}
	return nil
}
