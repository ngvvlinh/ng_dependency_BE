package shop

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		s.CreateReceipt,
		s.UpdateReceipt,
		s.DeleteReceipt,
		s.GetReceipt,
		s.GetReceipts)
}

func (s *Service) CreateReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (_err error) {
	key := fmt.Sprintf("Create receipt %v-%v-%v-%v-%v-%v-%v-%v",
		q.Context.Shop.ID, q.Context.UserID, q.TraderId, q.Title, q.Description, q.Amount, q.Code, q.Type)
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

func (s *Service) createReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (*receipting.CreateReceiptCommand, error) {
	receipt := &receipting.Receipt{
		TraderID:    q.TraderId,
		ShopID:      q.Context.Shop.ID,
		CreatedBy:   q.Context.UserID,
		Title:       q.Title,
		Type:        q.Type,
		Description: q.Description,
		Code:        q.Code,
		Amount:      q.Amount,
		Lines:       pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
	}
	if receipt.TraderID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Đối tác không hợp lệ")
	}
	if receipt.Type != receipting.ReceiptType && receipt.Type != receipting.PaymentType {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Loại đối tác không hợp lệ")
	}
	if err := s.validateReceiptForCreateOrUpdate(ctx, q.Context.Shop.ID, receipt); err != nil {
		return nil, err
	}

	code := strings.TrimSpace(q.Code)
	if code == "" {
		// TODO: generate code
		code = "TODO"
	}

	cmd := &receipting.CreateReceiptCommand{
		ShopID:      q.Context.Shop.ID,
		CreatedBy:   q.Context.UserID,
		TraderID:    q.TraderId,
		Code:        code,
		Title:       q.Title,
		Description: q.Description,
		Amount:      q.Amount,
		Type:        q.Type,
		Lines:       receipt.Lines,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (s *Service) UpdateReceipt(ctx context.Context, q *wrapshop.UpdateReceiptEndpoint) error {
	query := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu thu").
			Throw()
	}

	lines := pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines)
	receipt := &receipting.Receipt{
		ID:          q.Id,
		TraderID:    PInt64(q.TraderId).Apply(0),
		Title:       PString(q.Title).Apply(""),
		Description: PString(q.Description).Apply(""),
		Code:        PString(q.Code).Apply(""),
		Amount:      PInt32(q.Amount).Apply(0),
		Lines:       lines,
	}
	if err := s.validateReceiptForCreateOrUpdate(ctx, q.Context.Shop.ID, receipt); err != nil {
		return err
	}

	cmd := &receipting.UpdateReceiptCommand{
		ID:          q.Id,
		ShopID:      q.Context.Shop.ID,
		TraderID:    PInt64(q.TraderId),
		Code:        PString(q.Code),
		Title:       PString(q.Title),
		Description: PString(q.Description),
		Amount:      PInt32(q.Amount),
		Lines:       lines,
	}
	err := receiptAggr.Dispatch(ctx, cmd)
	if err != nil {
		errMgs := err.Error()
		switch {
		case strings.Contains(errMgs, "receipt_shop_id_code_idx"):
			err = cm.Errorf(cm.FailedPrecondition, nil, "Mã phiếu %v đã tồn tại. Vui lòng chọn mã khác.", *q.Code)
		}
		return err
	}

	q.Result = pbshop.PbReceipt(cmd.Result)
	return nil
}

func (s *Service) validateReceiptForCreateOrUpdate(ctx context.Context, shopID int64, receipt *receipting.Receipt) error {
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

	// validate receipt lines
	if receipt.Lines != nil && len(receipt.Lines) > 0 {
		if err := s.validateReceiptLines(ctx, traderType, receipt); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) validateReceiptLines(ctx context.Context, traderType string, receipt *receipting.Receipt) error {
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
	listReceiptsQuery := &receipting.ListReceiptsByOrderIDsQuery{
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
			if receiptLine.OrderID == 0 {
				continue
			}
			if _, has := mapOrdersAmount[receiptLine.OrderID]; has {
				mapOrdersAmountOld[receiptLine.OrderID] += receiptLine.Amount
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

		if receiptLine.OrderID == 0 {
			continue
		}

		// Check has key in map
		// hasKey = true -> duplicate orderId in receipt
		// hasKey = false -> add orderId in map
		if _, has := mapOrdersAmount[receiptLine.OrderID]; has {
			err = cm.Errorf(cm.FailedPrecondition, nil, "Duplicated OrderId %d in receipt", receiptLine.OrderID)
			return
		}

		mapOrdersAmount[receiptLine.OrderID] = receiptLine.Amount
		orderIDs = append(orderIDs, receiptLine.OrderID)
	}
	return
}

func (s *Service) DeleteReceipt(ctx context.Context, q *wrapshop.DeleteReceiptEndpoint) error {
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

func (s *Service) GetReceipt(ctx context.Context, q *wrapshop.GetReceiptEndpoint) error {
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

func (s *Service) GetReceipts(ctx context.Context, q *wrapshop.GetReceiptsEndpoint) error {
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

func (s *Service) getInfosForReceipts(ctx context.Context, shopID int64, receipts []*receipting.Receipt) (receiptsResult []*pbshop.Receipt, _ error) {
	mapOrderIDAndReceivedAmount := make(map[int64]int32)
	var orderIDs, userIDs, traderIDs []int64

	receiptsResult = pbshop.PbReceipts(receipts)

	// Get orderIDs into receiptLines
	for _, receipt := range receipts {
		userIDs = append(userIDs, receipt.CreatedBy)
		traderIDs = append(traderIDs, receipt.TraderID)

		for _, receiptLine := range receipt.Lines {
			if receiptLine.OrderID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.OrderID]; !ok {
				orderIDs = append(orderIDs, receiptLine.OrderID)
				mapOrderIDAndReceivedAmount[receiptLine.OrderID] = 0
			}
		}
	}

	// Get receipts have orderID into orderIDs(above)
	// Calculate received amount each of orderID
	if err := calcReceivedAmounts(ctx, shopID, orderIDs, mapOrderIDAndReceivedAmount); err != nil {
		return nil, err
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

	// List orders
	if len(orderIDs) > 0 {
		if err := listOrdersAndRecalcReceivedAmounts(ctx, shopID, orderIDs, mapOrderIDAndReceivedAmount, receiptsResult); err != nil {
			return nil, err
		}
	}

	// List traders
	if err := listTraders(ctx, shopID, traderIDs, receiptsResult); err != nil {
		return nil, err
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

func listOrdersAndRecalcReceivedAmounts(
	ctx context.Context, shopID int64, orderIDs []int64,
	mapOrderIDAndReceivedAmount map[int64]int32, receiptsResult []*pbshop.Receipt,
) error {
	mapOrders := make(map[int64]*ordering.Order)
	getOrdersQuery := &ordering.GetOrdersQuery{
		ShopID: shopID,
		IDs:    orderIDs,
	}
	err := orderQuery.Dispatch(ctx, getOrdersQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		for _, order := range getOrdersQuery.Result.Orders {
			mapOrders[order.ID] = order
		}
	default:
		return err
	}
	for _, receipt := range receiptsResult {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.OrderId == 0 {
				continue
			}
			switch receipt.Type {
			case receipting.ReceiptType:
				receiptLine.ReceivedAmount = mapOrderIDAndReceivedAmount[receiptLine.OrderId] - receiptLine.Amount
			case receipting.PaymentType:
				receiptLine.ReceivedAmount = -(mapOrderIDAndReceivedAmount[receiptLine.OrderId] + receiptLine.Amount)
			}

			orderID := receiptLine.OrderId
			receiptLine.Order = &pbshop.OrderOfReceiptLine{
				Id:          orderID,
				ShopId:      shopID,
				Code:        mapOrders[orderID].Code,
				TotalAmount: int32(mapOrders[orderID].TotalAmount),
			}
		}
	}
	return nil
}

func calcReceivedAmounts(
	ctx context.Context, shopID int64,
	orderIDs []int64, mapOrderIDAndReceivedAmount map[int64]int32,
) error {
	getReceiptsByOrderIDs := &receipting.ListReceiptsByOrderIDsQuery{
		IDs:    orderIDs,
		ShopID: shopID,
	}
	if err := receiptQuery.Dispatch(ctx, getReceiptsByOrderIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByOrderIDs.Result.Receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.OrderID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.OrderID]; !ok {
				continue
			}

			switch receipt.Type {
			case receipting.ReceiptType:
				mapOrderIDAndReceivedAmount[receiptLine.OrderID] += receiptLine.Amount
			case receipting.PaymentType:
				mapOrderIDAndReceivedAmount[receiptLine.OrderID] -= receiptLine.Amount
			}
		}
	}
	return nil
}
