package shop

import (
	"context"

	"etop.vn/api/main/receipting"
	"etop.vn/api/main/refund"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/top/int/shop"
	"etop.vn/api/top/types/etc/receipt_ref"
	"etop.vn/api/top/types/etc/status3"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/auth"
	"etop.vn/capi/dot"
)

func (s *RefundService) CreateRefund(ctx context.Context, q *CreateRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var lines []*refund.RefundLine
	for _, v := range q.Lines {
		lines = append(lines, &refund.RefundLine{
			VariantID:  v.VariantID,
			Quantity:   v.Quantity,
			Adjustment: v.Adjustment,
		})
	}
	cmd := refund.CreateRefundCommand{
		Lines:           lines,
		OrderID:         q.OrderID,
		TotalAdjustment: q.TotalAjustment,
		AdjustmentLines: q.AdjustmentLines,
		TotalAmount:     q.TotalAmount,
		BasketValue:     q.BasketValue,
		ShopID:          shopID,
		CreatedBy:       userID,
		Note:            q.Note,
	}
	err := RefundAggr.Dispatch(ctx, &cmd)
	if err != nil {
		return err
	}
	result := PbRefund(cmd.Result)
	result, err = populateRefundWithCustomer(ctx, result)
	if err != nil {
		return err
	}
	result, err = populateRefundWithReceiptPaidAmount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *RefundService) UpdateRefund(ctx context.Context, q *UpdateRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var lines []*refund.RefundLine
	for _, v := range q.Lines {
		lines = append(lines, &refund.RefundLine{
			VariantID:  v.VariantID,
			Quantity:   v.Quantity,
			Adjustment: v.Adjustment,
		})
	}
	cmd := refund.UpdateRefundCommand{
		Lines:           lines,
		ID:              q.ID,
		ShopID:          shopID,
		UpdateBy:        userID,
		Note:            q.Note,
		TotalAmount:     q.TotalAmount,
		BasketValue:     q.BasketValue,
		AdjustmentLines: q.AdjustmentLines,
		TotalAdjustment: q.TotalAjustment,
	}
	if err := RefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbRefund(cmd.Result)
	result, err := populateRefundWithCustomer(ctx, result)
	if err != nil {
		return err
	}
	result, err = populateRefundWithReceiptPaidAmount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *RefundService) ConfirmRefund(ctx context.Context, q *ConfirmRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	roles := auth.Roles(q.Context.Roles)
	cmd := refund.ConfirmRefundCommand{
		ShopID:               shopID,
		ID:                   q.ID,
		UpdatedBy:            userID,
		AutoInventoryVoucher: checkRoleAutoInventoryVoucher(roles, q.AutoInventoryVoucher),
	}
	if err := RefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbRefund(cmd.Result)
	result, err := populateRefundWithCustomer(ctx, result)
	if err != nil {
		return err
	}
	result, err = populateRefundWithReceiptPaidAmount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *RefundService) CancelRefund(ctx context.Context, q *CancelRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	roles := auth.Roles(q.Context.Roles)
	cmd := refund.CancelRefundCommand{
		ShopID:               shopID,
		ID:                   q.ID,
		UpdatedBy:            userID,
		CancelReason:         q.CancelReason,
		AutoInventoryVoucher: checkRoleAutoInventoryVoucher(roles, q.AutoInventoryVoucher),
	}
	if err := RefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbRefund(cmd.Result)
	result, err := populateRefundWithCustomer(ctx, result)
	if err != nil {
		return err
	}
	result, err = populateRefundWithReceiptPaidAmount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *RefundService) GetRefund(ctx context.Context, q *GetRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &refund.GetRefundByIDQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := RefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID:            query.Result.OrderID,
		IncludeFulfillment: false,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return err
	}
	result := PbRefund(query.Result)
	result, err := populateRefundWithCustomer(ctx, result)
	if err != nil {
		return err
	}
	result, err = populateRefundWithReceiptPaidAmount(ctx, result)
	if err != nil {
		return err
	}
	result.Customer = convertpb.PbOrderCustomer(queryOrder.Result.Order.Customer)
	q.Result = result
	return nil
}

func (s *RefundService) GetRefundsByIDs(ctx context.Context, q *GetRefundsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &refund.GetRefundsByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := RefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbRefunds(query.Result)
	var err error
	if len(result) > 0 {
		result, err = populateRefundsWithCustomer(ctx, result)
		if err != nil {
			return err
		}
		result, err = populateRefundsWithReceiptPaidAmount(ctx, result)
		if err != nil {
			return err
		}
	}
	q.Result = &shop.GetRefundsByIDsResponse{
		Refund: result,
	}
	return nil
}

func (s *RefundService) GetRefunds(ctx context.Context, q *GetRefundsEndpoint) error {
	shopID := q.Context.Shop.ID
	paging := cmapi.CMPaging(q.Paging)
	query := &refund.GetRefundsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := RefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbRefunds(query.Result.Refunds)
	var err error
	if len(result) > 0 {
		result, err = populateRefundsWithCustomer(ctx, result)
		if err != nil {
			return err
		}
		result, err = populateRefundsWithReceiptPaidAmount(ctx, result)
		if err != nil {
			return err
		}
	}
	q.Result = &shop.GetRefundsResponse{
		Refunds: result,
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}

// Get total paid amount of refund from receipt which have status = P
func populateRefundWithReceiptPaidAmount(ctx context.Context, arg *shop.Refund) (*shop.Refund, error) {
	query := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  arg.ShopID,
		RefIDs:  []dot.ID{arg.ID},
		RefType: receipt_ref.Refund,
		Status:  int(status3.P),
	}
	err := receiptQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	for _, value := range query.Result.Receipts {
		for _, _value := range value.Lines {
			if _value.RefID == arg.ID {
				arg.PaidAmount += _value.Amount
			}
		}
	}
	return arg, nil
}

// Get total paid amount of each refunds  from receipt which have status = P
func populateRefundsWithReceiptPaidAmount(ctx context.Context, refunds []*shop.Refund) ([]*shop.Refund, error) {
	if len(refunds) == 0 {
		return refunds, nil
	}
	var refundIDs []dot.ID
	for _, value := range refunds {
		refundIDs = append(refundIDs, value.ID)
	}
	query := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  refunds[0].ShopID,
		RefIDs:  refundIDs,
		RefType: receipt_ref.Refund,
		Status:  int(status3.P),
	}
	err := receiptQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapPaidAmount = make(map[dot.ID]int)
	for _, value := range query.Result.Receipts {
		for _, _value := range value.Lines {
			mapPaidAmount[_value.RefID] = _value.Amount + mapPaidAmount[_value.RefID]
		}
	}
	for key, value := range refunds {
		refunds[key].PaidAmount = mapPaidAmount[value.ID]
	}
	return refunds, nil
}

func populateRefundsWithCustomer(ctx context.Context, refunds []*shop.Refund) ([]*shop.Refund, error) {
	var orderIDs []dot.ID
	for _, value := range refunds {
		orderIDs = append(orderIDs, value.OrderID)
	}
	// Get informations about customers from orders
	queryOrder := &ordermodelx.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return nil, err
	}
	// make a map [ OrderID ] Order
	var orderCustomerMap = make(map[dot.ID]*ordermodel.Order, len(queryOrder.Result.Orders))
	for _, value := range queryOrder.Result.Orders {
		orderCustomerMap[value.ID] = value.Order
	}
	var cutomerIDs []dot.ID
	for key, value := range refunds {
		// Add customer's information to refunds
		refunds[key].Customer = convertpb.PbOrderCustomer(orderCustomerMap[value.OrderID].Customer)
		customerID := orderCustomerMap[value.OrderID].CustomerID
		if customerID != 0 {
			refunds[key].CustomerID = customerID
		}
		cutomerIDs = append(cutomerIDs, customerID)
	}

	// Get all customers in list
	queryCustomer := &customering.ListCustomersByIDsQuery{
		IDs:    cutomerIDs,
		ShopID: refunds[0].ShopID,
	}
	err := customerQuery.Dispatch(ctx, queryCustomer)
	if err != nil {
		return nil, err
	}
	// make a map [ customerID ] Customer
	var mapCustomers = make(map[dot.ID]bool, len(queryCustomer.Result.Customers))
	for _, v := range queryCustomer.Result.Customers {
		mapCustomers[v.ID] = true
	}
	for key, value := range refunds {
		refunds[key].Customer.Deleted = true
		// Check customer have been deleted
		if value.CustomerID != 0 && mapCustomers[value.CustomerID] {
			refunds[key].Customer.Deleted = false
		}
	}
	return refunds, nil
}

func populateRefundWithCustomer(ctx context.Context, refundArg *shop.Refund) (*shop.Refund, error) {
	// Get information about customer from order
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID:            refundArg.OrderID,
		IncludeFulfillment: false,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return nil, err
	}
	// Add customer's information to refund
	refundArg.Customer = convertpb.PbOrderCustomer(queryOrder.Result.Order.Customer)
	if queryOrder.Result.Order.CustomerID != 0 {
		refundArg.Customer.Deleted = true
		queryCustomer := &customering.ListCustomersByIDsQuery{
			IDs:    []dot.ID{queryOrder.Result.Order.CustomerID},
			ShopID: refundArg.ShopID,
		}
		// Check customer have been deleted
		err := customerQuery.Dispatch(ctx, queryCustomer)
		if err != nil {
			return nil, err
		}
		refundArg.CustomerID = queryOrder.Result.Order.CustomerID
		if len(queryCustomer.Result.Customers) > 0 {
			refundArg.Customer.Deleted = false
		}
	}
	return refundArg, nil
}
