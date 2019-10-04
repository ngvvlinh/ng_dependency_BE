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
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("api",
		CreateReceipt,
		UpdateReceipt,
		DeleteReceipt,
		GetReceipt,
		GetReceipts)
}

func CreateReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (_err error) {
	var receiptLinesTemp []*receipting.ReceiptLine
	for _, receiptLine := range q.Lines {
		receiptLinesTemp = append(receiptLinesTemp, &receipting.ReceiptLine{
			OrderID: receiptLine.OrderId,
			Title:   receiptLine.Title,
			Amount:  receiptLine.Amount,
		})
	}

	receipt := &receipting.Receipt{
		TraderID:    q.TraderId,
		UserID:      q.UserId,
		Title:       q.Title,
		Type:        q.Type,
		Description: q.Description,
		Code:        q.Code,
		Amount:      q.Amount,
		Lines:       receiptLinesTemp,
	}

	key := fmt.Sprintf("Create receipt %v-%v-%v-%v-%v-%v-%v-%v",
		q.Context.Shop.ID, q.UserId, q.TraderId, q.Title, q.Description, q.Amount, q.Code, q.Type)
	result, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			if err := validateForCreateAndUpdateReceipt(ctx, q.Context.Shop.ID, true, receipt); err != nil {
				return nil, err
			}

			code := strings.TrimSpace(q.Code)

			if code == "" {
				// TODO: generate code
				code = "@vnngmq"
			}

			cmd := &receipting.CreateReceiptCommand{
				ShopID:      q.Context.Shop.ID,
				TraderID:    receipt.TraderID,
				UserID:      receipt.UserID,
				Code:        code,
				Title:       receipt.Title,
				Description: receipt.Description,
				Amount:      receipt.Amount,
				Type:        q.Type,
				Lines:       receipt.Lines,
			}
			if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
				return nil, err
			}

			return cmd, nil

		}, "Create receipt")

	if err != nil {
		return err
	}

	q.Result = pbshop.PbReceipt(result.(*receipting.CreateReceiptCommand).Result)

	return nil
}

func UpdateReceipt(ctx context.Context, q *wrapshop.UpdateReceiptEndpoint) error {
	query := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return cm.Errorf(cm.NotFound, nil, "Receipt not found")
	}

	var receiptLinesTemp []*receipting.ReceiptLine
	for _, receiptLine := range q.Lines {
		receiptLinesTemp = append(receiptLinesTemp, &receipting.ReceiptLine{
			OrderID: receiptLine.OrderId,
			Title:   receiptLine.Title,
			Amount:  receiptLine.Amount,
		})
	}

	receipt := &receipting.Receipt{
		ID:          q.Id,
		TraderID:    q.TraderId,
		UserID:      q.UserId,
		Title:       q.Title,
		Description: q.Description,
		Code:        q.Code,
		Amount:      q.Amount,
		Lines:       receiptLinesTemp,
	}
	if err := validateForCreateAndUpdateReceipt(ctx, q.Context.Shop.ID, false, receipt); err != nil {
		return err
	}

	receipt.Code = strings.TrimSpace(receipt.Code)

	cmd := &receipting.UpdateReceiptCommand{
		ID:          q.Id,
		ShopID:      q.Context.Shop.ID,
		TraderID:    PInt64(&receipt.TraderID),
		UserID:      PInt64(&receipt.UserID),
		Code:        PString(&receipt.Code),
		Title:       PString(&receipt.Title),
		Description: PString(&receipt.Description),
		Amount:      PInt32(&receipt.Amount),
		Lines:       receipt.Lines,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = pbshop.PbReceipt(cmd.Result)

	return nil
}

func validateForCreateAndUpdateReceipt(
	ctx context.Context, shopID int64,
	isCreateReceipt bool, receiptParam *receipting.Receipt) error {

	// Validate data
	if receiptParam.Title == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Title not empty")
	}
	if receiptParam.Amount <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Amount must rather than 0")
	}
	if isCreateReceipt {
		if receiptParam.UserID == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "UserId not empty")
		}
		if receiptParam.TraderID == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "ShopTraderId not empty")
		}
		if receiptParam.Type != receipting.ReceiptType && receiptParam.Type != receipting.PaymentType {
			return cm.Errorf(cm.InvalidArgument, nil, "Invalid type")
		}
	}
	if strings.TrimSpace(receiptParam.Code) != "" {
		query := &receipting.GetReceiptByCodeQuery{
			Code:   strings.TrimSpace(receiptParam.Code),
			ShopID: shopID,
		}

		if err := receiptQuery.Dispatch(ctx, query); err != nil {
			if cm.ErrorCode(err) != cm.NotFound {
				return err
			}
		}

		if query.Result != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Code is exist")
		}
	}

	if receiptParam.UserID != 0 {
		query := &model.GetAccountUserQuery{
			UserID:    receiptParam.UserID,
			AccountID: shopID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "User not found")
		}
	}
	if receiptParam.TraderID != 0 {
		query := &tradering.GetTraderByIDQuery{
			ID:     receiptParam.TraderID,
			ShopID: shopID,
		}
		if err := traderQuery.Dispatch(ctx, query); err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "ShopTrader not found")
		}
	}

	// List orderId of receiptLines
	var orderIds []int64

	// Check receiptLines
	if receiptParam.Lines != nil && len(receiptParam.Lines) > 0 {
		var totalAmountOfReceiptLines int32 = 0

		// Map of [ orderId ] amount of receiptLines (params)
		mOrdersAmountParam := make(map[int64]int32)

		for _, receiptLine := range receiptParam.Lines {
			totalAmountOfReceiptLines += receiptLine.Amount
			if receiptLine.OrderID != 0 {
				// Check has key in map
				// hasKey = true -> duplicate orderId in receipt
				// hasKey = false -> add orderId in map
				if _, hasKey := mOrdersAmountParam[receiptLine.OrderID]; !hasKey {
					mOrdersAmountParam[receiptLine.OrderID] = receiptLine.Amount
					orderIds = append(orderIds, receiptLine.OrderID)
				} else {
					return cm.Errorf(cm.InvalidArgument, nil, "Duplicate OrderId %d is exist in receipt", receiptLine.OrderID)
				}

			}

			// check amount < 0
			if receiptLine.Amount <= 0 {
				return cm.Errorf(cm.InvalidArgument, nil, "Amount of receiptLine must be rather than 0")
			}
		}

		if totalAmountOfReceiptLines != receiptParam.Amount {
			return cm.Errorf(cm.InvalidArgument, nil, "Amount of receipt must equal to total amount of receiptLines")
		}

		if orderIds != nil && len(orderIds) > 0 {
			mOrders := make(map[int64]int32)
			// List all orders in orderIDs
			listOrdersQuery := &ordering.GetOrdersQuery{
				ShopID: shopID,
				IDs:    orderIds,
			}
			if err := orderQuery.Dispatch(ctx, listOrdersQuery); err != nil {
				return err
			}

			for _, order := range listOrdersQuery.Result.Orders {
				mOrders[order.ID] = int32(order.TotalAmount)
			}

			// Check orderIds with orderIds of listOrdersQuery.Result
			// When different len
			if len(orderIds) != len(mOrders) {
				for _, v := range orderIds {
					if _, ok := mOrders[v]; !ok {
						return cm.Errorf(cm.InvalidArgument, nil, "OrderId %d not found", v)
					}
				}
			}
		}
	}
	return nil
}

func DeleteReceipt(ctx context.Context, q *wrapshop.DeleteReceiptEndpoint) error {
	cmd := &receipting.DeleteReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func GetReceipt(ctx context.Context, q *wrapshop.GetReceiptEndpoint) error {
	// Check receipt is exist
	getReceiptQuery := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, getReceiptQuery); err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Receipt not found")
	}
	receiptResult := getReceiptQuery.Result
	receipt := pbshop.PbReceipt(getReceiptQuery.Result)
	receipt.Lines = pbshop.PbReceiptLines(receiptResult.Lines)

	// Get list orderIds
	mOrderIDsAndReceivedAmounts := make(map[int64]int32)
	var orderIds []int64

	for _, receiptLine := range getReceiptQuery.Result.Lines {
		if receiptLine.OrderID != 0 {
			orderIds = append(orderIds, receiptLine.OrderID)
			mOrderIDsAndReceivedAmounts[receiptLine.OrderID] = 0
		}
	}

	// Get list orders
	if orderIds != nil && len(orderIds) > 0 {
		getOrdersQuery := &ordering.GetOrdersQuery{
			ShopID: q.Context.Shop.ID,
			IDs:    orderIds,
		}
		if err := orderQuery.Dispatch(ctx, getOrdersQuery); err != nil {
			return err
		}

		getReceiptsByOrderIDs := &receipting.ListReceiptsByOrderIDsQuery{
			IDs:    orderIds,
			ShopID: q.Context.Shop.ID,
		}
		if err := receiptQuery.Dispatch(ctx, getReceiptsByOrderIDs); err != nil {
			return err
		}
		for _, receipt := range getReceiptsByOrderIDs.Result.Receipts {
			for _, receiptLine := range receipt.Lines {
				if receiptLine.OrderID == 0 {
					continue
				}
				if _, ok := mOrderIDsAndReceivedAmounts[receiptLine.OrderID]; ok {
					switch receipt.Type {
					case receipting.ReceiptType:
						mOrderIDsAndReceivedAmounts[receiptLine.OrderID] += receiptLine.Amount
					case receipting.PaymentType:
						mOrderIDsAndReceivedAmounts[receiptLine.OrderID] -= receiptLine.Amount
					}

				}
			}
		}

		for _, receiptLine := range receipt.Lines {
			if receiptLine.OrderId == 0 {
				continue
			}
			receiptLine.ReceivedAmount = mOrderIDsAndReceivedAmounts[receiptLine.OrderId] - receiptLine.Amount
			for _, order := range getOrdersQuery.Result.Orders {
				if order.ID == receiptLine.OrderId {
					receiptLine.Order = &pbshop.OrderOfReceiptLine{
						Id:          order.ID,
						ShopId:      order.ShopID,
						Code:        order.Code,
						TotalAmount: int32(order.TotalAmount),
					}
					break
				}
			}
		}
	}

	// Get user
	getUserQuery := &model.GetUserByIDQuery{
		UserID: receipt.UserId,
	}
	if err := bus.Dispatch(ctx, getUserQuery); err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "User not found")
	}
	receipt.User = pbetop.PbUser(getUserQuery.Result)

	// Get trader
	getTraderQuery := &tradering.GetTraderByIDQuery{
		ID:     receipt.TraderId,
		ShopID: q.Context.Shop.ID,
	}
	if err := traderQuery.Dispatch(ctx, getTraderQuery); err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Trader not found")
	}

	switch getTraderQuery.Result.Type {
	case tradering.CustomerType:
		query := &customering.GetCustomerByIDQuery{
			ID:     getTraderQuery.Result.ID,
			ShopID: q.Context.Shop.ID,
		}
		if err := customerQuery.Dispatch(ctx, query); err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Customer not found")
		}
		receipt.Partner = &pbshop.Partner{
			Type:     tradering.CustomerType,
			FullName: query.Result.FullName,
		}
	case tradering.VendorType:
		query := &vendoring.GetVendorByIDQuery{
			ID:     getTraderQuery.Result.ID,
			ShopID: q.Context.Shop.ID,
		}
		if err := vendorQuery.Dispatch(ctx, query); err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Vendor not found")
		}
		receipt.Partner = &pbshop.Partner{
			Type:     tradering.VendorType,
			FullName: query.Result.FullName,
		}
	case tradering.CarrierType:
		query := &carrying.GetCarrierByIDQuery{
			ID:     getTraderQuery.Result.ID,
			ShopID: q.Context.Shop.ID,
		}
		if err := carrierQuery.Dispatch(ctx, query); err != nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Carrier not found")
		}
		receipt.Partner = &pbshop.Partner{
			Type:     tradering.CarrierType,
			FullName: query.Result.FullName,
		}
	}
	receipt.Partner.Id = getTraderQuery.Result.ID

	q.Result = receipt

	return nil
}

func GetReceipts(ctx context.Context, q *wrapshop.GetReceiptsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &receipting.ListReceiptsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	mOrderIDsAndReceivedAmounts := make(map[int64]int32)
	var orderIds, userIds, traderIds []int64

	receipts := pbshop.PbReceipts(query.Result.Receipts)

	for _, receipt := range receipts {
		userIds = append(userIds, receipt.UserId)
		traderIds = append(traderIds, receipt.TraderId)
		for _, receiptLine := range receipt.Lines {
			if receiptLine.OrderId != 0 {
				if _, ok := mOrderIDsAndReceivedAmounts[receiptLine.OrderId]; ok != false {
					orderIds = append(orderIds, receiptLine.OrderId)
				}
				switch receipt.Type {
				case receipting.ReceiptType:
					mOrderIDsAndReceivedAmounts[receiptLine.OrderId] += receiptLine.Amount
				case receipting.PaymentType:
					mOrderIDsAndReceivedAmounts[receiptLine.OrderId] -= receiptLine.Amount
				}
			}
		}
	}

	// Get users of current account
	getUsersOfCurrAccount := &model.GetAccountUserExtendedsQuery{
		AccountIDs: []int64{q.Context.Shop.ID},
	}
	if err := bus.Dispatch(ctx, getUsersOfCurrAccount); err != nil {
		return err
	}
	mUserIdsAndUser := make(map[int64]*model.User)
	for _, accountUser := range getUsersOfCurrAccount.Result.AccountUsers {
		mUserIdsAndUser[accountUser.User.ID] = accountUser.User
	}
	for _, receipt := range receipts {
		receipt.User = pbetop.PbUser(mUserIdsAndUser[receipt.UserId])
	}

	// Get list orders
	if orderIds != nil && len(orderIds) > 0 {
		getOrdersQuery := &ordering.GetOrdersQuery{
			ShopID: q.Context.Shop.ID,
			IDs:    orderIds,
		}
		if err := orderQuery.Dispatch(ctx, getOrdersQuery); err != nil {
			return err
		}
		for _, receipt := range receipts {
			for _, receiptLine := range receipt.Lines {
				if receiptLine.OrderId == 0 {
					continue
				}
				receiptLine.ReceivedAmount = mOrderIDsAndReceivedAmounts[receiptLine.OrderId] - receiptLine.Amount
				for _, order := range getOrdersQuery.Result.Orders {
					if order.ID == receiptLine.OrderId {
						receiptLine.Order = &pbshop.OrderOfReceiptLine{
							Id:          order.ID,
							ShopId:      order.ShopID,
							Code:        order.Code,
							TotalAmount: int32(order.TotalAmount),
						}
						break
					}
				}
			}
		}
	}

	// List traders
	mVendors := make(map[int64]*vendoring.ShopVendor)
	var vendorIDs, customerIDs, carrierIDs []int64
	mCustomers := make(map[int64]*customering.ShopCustomer)
	mCarriers := make(map[int64]*carrying.ShopCarrier)
	getTradersByIDsQuery := &tradering.ListTradersByIDsQuery{
		ShopID: q.Context.Shop.ID,
		IDs:    traderIds,
	}
	if err := traderQuery.Dispatch(ctx, getTradersByIDsQuery); err != nil {
		return nil
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
			IDs:    vendorIDs,
			ShopID: q.Context.Shop.ID,
		}
		if err := vendorQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, vendor := range query.Result.Vendors {
			mVendors[vendor.ID] = vendor
		}
	}
	if customerIDs != nil && len(customerIDs) > 0 {
		query := &customering.ListCustomersByIDsQuery{
			IDs:    customerIDs,
			ShopID: q.Context.Shop.ID,
		}
		if err := customerQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, customer := range query.Result.Customers {
			mCustomers[customer.ID] = customer
		}
	}
	if carrierIDs != nil && len(carrierIDs) > 0 {
		query := &carrying.ListCarriersByIDsQuery{
			IDs:    carrierIDs,
			ShopID: q.Context.Shop.ID,
		}
		if err := carrierQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, carrier := range query.Result.Carriers {
			mCarriers[carrier.ID] = carrier
		}
	}

	for _, receipt := range receipts {
		if value, ok := mVendors[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       value.ID,
				Type:     tradering.VendorType,
				FullName: value.FullName,
			}
		}
		if value, ok := mCustomers[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       value.ID,
				Type:     tradering.CustomerType,
				FullName: value.FullName,
			}
		}
		if value, ok := mCarriers[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       value.ID,
				Type:     tradering.CarrierType,
				FullName: value.FullName,
			}
		}
	}

	q.Result = &pbshop.ReceiptsResponse{
		Receipts: receipts,
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}
