package shop

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/purchaserefund"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/top/int/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/auth"
	"etop.vn/capi/dot"
)

func (s *PurchaseRefundService) CreatePurchaseRefund(ctx context.Context, q *CreatePurchaseRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var lines []*purchaserefund.PurchaseRefundLine
	for _, v := range q.Lines {
		lines = append(lines, &purchaserefund.PurchaseRefundLine{
			VariantID:    v.VariantID,
			Quantity:     v.Quantity,
			PaymentPrice: v.PaymentPrice,
			Adjustment:   v.Adjustment,
		})
	}
	cmd := purchaserefund.CreatePurchaseRefundCommand{
		Lines:           lines,
		PurchaseOrderID: q.PurchaseOrderID,
		AdjustmentLines: q.AdjustmentLines,
		TotalAdjustment: q.TotalAdjustment,
		TotalAmount:     q.TotalAmount,
		BasketValue:     q.BasketValue,
		ShopID:          shopID,
		CreatedBy:       userID,
		Note:            q.Note,
	}
	err := PurchaseRefundAggr.Dispatch(ctx, &cmd)
	if err != nil {
		return err
	}
	result := PbPurchaseRefund(cmd.Result)
	result, err = populatePurchaseRefundWithSupplier(ctx, result)
	if err != nil {
		return err
	}
	result, err = populatePurchaseRefundWithInventoryVoucher(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *PurchaseRefundService) UpdatePurchaseRefund(ctx context.Context, q *UpdatePurchaseRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var lines []*purchaserefund.PurchaseRefundLine
	for _, v := range q.Lines {
		lines = append(lines, &purchaserefund.PurchaseRefundLine{
			VariantID: v.VariantID,
			Quantity:  v.Quantity,
		})
	}
	cmd := purchaserefund.UpdatePurchaseRefundCommand{
		Lines:           lines,
		ID:              q.ID,
		ShopID:          shopID,
		AdjustmentLines: q.AdjustmentLines,
		TotalAdjustment: q.TotalAdjustment,
		TotalAmount:     q.TotalAmount,
		UpdateBy:        userID,
		BasketValue:     q.BasketValue,
		Note:            q.Note,
	}
	if err := PurchaseRefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbPurchaseRefund(cmd.Result)
	result, err := populatePurchaseRefundWithSupplier(ctx, result)
	if err != nil {
		return err
	}
	result, err = populatePurchaseRefundWithInventoryVoucher(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *PurchaseRefundService) ConfirmPurchaseRefund(ctx context.Context, q *ConfirmPurchaseRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	inventoryOverStock := q.Context.Shop.InventoryOverstock
	roles := auth.Roles(q.Context.Roles)
	cmd := purchaserefund.ConfirmPurchaseRefundCommand{
		ShopID:               shopID,
		ID:                   q.ID,
		UpdatedBy:            userID,
		AutoInventoryVoucher: checkRoleAutoInventoryVoucher(roles, q.AutoInventoryVoucher),
		InventoryOverStock:   inventoryOverStock.Apply(true),
	}
	if err := PurchaseRefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbPurchaseRefund(cmd.Result)
	result, err := populatePurchaseRefundWithSupplier(ctx, result)
	if err != nil {
		return err
	}
	result, err = populatePurchaseRefundWithInventoryVoucher(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *PurchaseRefundService) CancelPurchaseRefund(ctx context.Context, q *CancelPurchaseRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	roles := auth.Roles(q.Context.Roles)
	cmd := purchaserefund.CancelPurchaseRefundCommand{
		ShopID:               shopID,
		ID:                   q.ID,
		UpdatedBy:            userID,
		CancelReason:         q.CancelReason,
		InventoryOverStock:   q.Context.Shop.InventoryOverstock.Apply(true),
		AutoInventoryVoucher: checkRoleAutoInventoryVoucher(roles, q.AutoInventoryVoucher),
	}
	if err := PurchaseRefundAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	result := PbPurchaseRefund(cmd.Result)
	result, err := populatePurchaseRefundWithSupplier(ctx, result)
	if err != nil {
		return err
	}
	result, err = populatePurchaseRefundWithInventoryVoucher(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

func (s *PurchaseRefundService) GetPurchaseRefund(ctx context.Context, q *GetPurchaseRefundEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &purchaserefund.GetPurchaseRefundByIDQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := PurchaseRefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	queryPurchaseOrder := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     query.Result.PurchaseOrderID,
		ShopID: q.Context.Shop.ID,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, queryPurchaseOrder); err != nil {
		return err
	}
	result := PbPurchaseRefund(query.Result)
	result, err := populatePurchaseRefundWithSupplier(ctx, result)
	if err != nil {
		return err
	}
	result, err = populatePurchaseRefundWithInventoryVoucher(ctx, result)
	if err != nil {
		return err
	}
	result.Supplier = convertpb.PbPurchaseOrderSupplier(queryPurchaseOrder.Result.Supplier)
	q.Result = result
	return nil
}

func (s *PurchaseRefundService) GetPurchaseRefundsByIDs(ctx context.Context, q *GetPurchaseRefundsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &purchaserefund.GetPurchaseRefundsByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := PurchaseRefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbPurchaseRefunds(query.Result)
	var err error
	if len(result) > 0 {
		result, err = populatePurchaseRefundsWithSupplier(ctx, result)
		if err != nil {
			return err
		}
		result, err = populatePurchaseRefundsWithInventoryVouchers(ctx, result)
		if err != nil {
			return err
		}
	}
	q.Result = &shop.GetPurchaseRefundsByIDsResponse{
		PurchaseRefund: result,
	}
	return nil
}

func (s *PurchaseRefundService) GetPurchaseRefunds(ctx context.Context, q *GetPurchaseRefundsEndpoint) error {
	shopID := q.Context.Shop.ID
	paging := cmapi.CMPaging(q.Paging)
	query := &purchaserefund.ListPurchaseRefundsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := PurchaseRefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbPurchaseRefunds(query.Result.PurchaseRefunds)
	var err error
	if len(result) > 0 {
		result, err = populatePurchaseRefundsWithSupplier(ctx, result)
		if err != nil {
			return err
		}
		result, err = populatePurchaseRefundsWithInventoryVouchers(ctx, result)
		if err != nil {
			return err
		}
	}
	q.Result = &shop.GetPurchaseRefundsResponse{
		PurchaseRefunds: result,
		Paging:          cmapi.PbPageInfo(paging),
	}
	return nil
}

func populatePurchaseRefundsWithSupplier(ctx context.Context, purchaseRefunds []*shop.PurchaseRefund) ([]*shop.PurchaseRefund, error) {
	if len(purchaseRefunds) == 0 {
		return purchaseRefunds, nil
	}
	var purchaseOrderIDs []dot.ID
	for _, value := range purchaseRefunds {
		purchaseOrderIDs = append(purchaseOrderIDs, value.PurchaseOrderID)
	}
	// Get informations about purchase_order
	queryPurchaseOrder := &purchaseorder.GetPurchaseOrdersByIDsQuery{
		IDs:    purchaseOrderIDs,
		ShopID: purchaseRefunds[0].ShopID,
		Result: nil,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, queryPurchaseOrder); err != nil {
		return nil, err
	}
	// make a map [ PurchaseOrderID ] PurchaseOrderID
	var purchaseOrderMap = make(map[dot.ID]*purchaseorder.PurchaseOrder, len(queryPurchaseOrder.Result.PurchaseOrders))
	for _, value := range queryPurchaseOrder.Result.PurchaseOrders {
		purchaseOrderMap[value.ID] = value
	}
	var supplierIDs []dot.ID
	for key, value := range purchaseRefunds {
		// Add supplier's information to purchaseRefunds
		purchaseRefunds[key].Supplier = convertpb.PbPurchaseOrderSupplier(purchaseOrderMap[value.PurchaseOrderID].Supplier)
		supplierID := purchaseOrderMap[value.PurchaseOrderID].SupplierID
		if supplierID != 0 {
			purchaseRefunds[key].SupplierID = supplierID
		}
		supplierIDs = append(supplierIDs, supplierID)
	}

	// Get all suppliers in list
	querySupplier := &suppliering.ListSuppliersByIDsQuery{
		IDs:    supplierIDs,
		ShopID: purchaseRefunds[0].ShopID,
	}
	err := supplierQuery.Dispatch(ctx, querySupplier)
	if err != nil {
		return nil, err
	}
	// make a map [ supplierID ] Suppliers
	var mapSuppliers = make(map[dot.ID]bool, len(querySupplier.Result.Suppliers))
	for _, v := range querySupplier.Result.Suppliers {
		mapSuppliers[v.ID] = true
	}
	for key, value := range purchaseRefunds {
		purchaseRefunds[key].Supplier.Deleted = true
		// Check supplier have been deleted
		if value.SupplierID != 0 && mapSuppliers[value.SupplierID] {
			purchaseRefunds[key].Supplier.Deleted = false
		}
	}
	return purchaseRefunds, nil
}

func populatePurchaseRefundWithSupplier(ctx context.Context, purchaseRefundArg *shop.PurchaseRefund) (*shop.PurchaseRefund, error) {
	// Get information about supplier from pruchase_order
	queryPurchaseOrder := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     purchaseRefundArg.PurchaseOrderID,
		ShopID: purchaseRefundArg.ShopID,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, queryPurchaseOrder); err != nil {
		return nil, err
	}
	// Add supplier's information to purchase_refund
	purchaseRefundArg.Supplier = convertpb.PbPurchaseOrderSupplier(queryPurchaseOrder.Result.Supplier)
	purchaseRefundArg.SupplierID = queryPurchaseOrder.Result.SupplierID
	if queryPurchaseOrder.Result.SupplierID != 0 {
		purchaseRefundArg.Supplier.Deleted = false
		querySupplier := &suppliering.GetSupplierByIDQuery{
			ID:     queryPurchaseOrder.Result.SupplierID,
			ShopID: purchaseRefundArg.ShopID,
		}
		// Check supplier have been deleted
		err := supplierQuery.Dispatch(ctx, querySupplier)
		if err != nil {
			switch cm.ErrorCode(err) {
			case cm.NotFound:
				purchaseRefundArg.Supplier.Deleted = true
			default:
				return nil, err
			}
		}
	}
	return purchaseRefundArg, nil
}

func populatePurchaseRefundWithInventoryVoucher(ctx context.Context, refundArg *shop.PurchaseRefund) (*shop.PurchaseRefund, error) {
	// Get inventory voucher
	queryInventoryVoucher := &inventory.GetInventoryVoucherQuery{
		ShopID: refundArg.ShopID,
		ID:     refundArg.ID,
	}
	if err := inventoryQuery.Dispatch(ctx, queryInventoryVoucher); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return refundArg, nil
		}
		return nil, err
	}
	// Add inventoryvoucher to refund
	refundArg.InventoryVoucher = PbShopInventoryVoucher(queryInventoryVoucher.Result)
	return refundArg, nil
}

func populatePurchaseRefundsWithInventoryVouchers(ctx context.Context, refundsArgs []*shop.PurchaseRefund) ([]*shop.PurchaseRefund, error) {
	if len(refundsArgs) == 0 {
		return refundsArgs, nil
	}
	var refundIDs []dot.ID
	for _, value := range refundsArgs {
		refundIDs = append(refundIDs, value.ID)
	}
	// Get inventoryVoucher
	queryInventoryVoucher := &inventory.GetInventoryVouchersByRefIDsQuery{
		RefIDs: refundIDs,
		ShopID: refundsArgs[0].ShopID,
	}
	if err := inventoryQuery.Dispatch(ctx, queryInventoryVoucher); err != nil {
		return nil, err
	}
	// make map[ref_id]inventoryVoucher
	var mapInventoryVoucher = make(map[dot.ID]*inventory.InventoryVoucher)
	for _, value := range queryInventoryVoucher.Result.InventoryVoucher {
		mapInventoryVoucher[value.RefID] = value
	}
	for key, value := range refundsArgs {
		refundsArgs[key].InventoryVoucher = PbShopInventoryVoucher(mapInventoryVoucher[value.ID])
	}
	return refundsArgs, nil
}
