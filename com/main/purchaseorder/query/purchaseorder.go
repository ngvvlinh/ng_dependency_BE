package query

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/receipting"
	"o.o/api/shopping"
	"o.o/api/shopping/suppliering"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/purchaseorder/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ purchaseorder.QueryService = &PurchaseOrderQuery{}

type PurchaseOrderQuery struct {
	db                    *cmsql.Database
	store                 sqlstore.PurchaseOrderStoreFactory
	eventBus              capi.EventBus
	supplierQuery         suppliering.QueryBus
	inventoryVoucherQuery inventory.QueryBus
	receiptQuery          receipting.QueryBus
}

func NewPurchaseOrderQuery(
	database com.MainDB, eventBus capi.EventBus,
	supplierQ suppliering.QueryBus, inventoryVoucherQ inventory.QueryBus,
	receiptQ receipting.QueryBus,
) *PurchaseOrderQuery {
	return &PurchaseOrderQuery{
		db:                    database,
		store:                 sqlstore.NewPurchaseOrderStore(database),
		eventBus:              eventBus,
		supplierQuery:         supplierQ,
		inventoryVoucherQuery: inventoryVoucherQ,
		receiptQuery:          receiptQ,
	}
}

func PurchaseOrderQueryMessageBus(q *PurchaseOrderQuery) purchaseorder.QueryBus {
	b := bus.New()
	return purchaseorder.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *PurchaseOrderQuery) GetPurchaseOrderByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*purchaseorder.PurchaseOrder, error) {
	purchaseOrder, err := q.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return nil, err
	}
	query := &suppliering.GetSupplierByIDQuery{
		ID:     purchaseOrder.SupplierID,
		ShopID: args.ShopID,
	}
	err = q.supplierQuery.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		purchaseOrder.Supplier.Deleted = true
	case cm.NoError:
		// no-op
	default:
		return nil, err
	}

	getInventoryVouchersQuery := &inventory.GetInventoryVouchersByRefIDsQuery{
		RefIDs: []dot.ID{purchaseOrder.ID},
		ShopID: args.ShopID,
	}
	if err := q.inventoryVoucherQuery.Dispatch(ctx, getInventoryVouchersQuery); err != nil {
		return nil, err
	}

	if len(getInventoryVouchersQuery.Result.InventoryVoucher) != 0 {
		purchaseOrder.InventoryVoucher = getInventoryVouchersQuery.Result.InventoryVoucher[0]
	}

	if err := q.addPaidAmount(ctx, args.ShopID, []*purchaseorder.PurchaseOrder{purchaseOrder}); err != nil {
		return nil, err
	}
	return purchaseOrder, nil
}

func (q *PurchaseOrderQuery) ListPurchaseOrders(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*purchaseorder.PurchaseOrdersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Filters(args.Filters)
	purchaseOrders, err := query.WithPaging(args.Paging).ListPurchaseOrders()
	if err != nil {
		return nil, err
	}

	var supplierIDs, purchaseOrderIDs []dot.ID
	mapPurchaseOrderIDAndInventoryVoucher := make(map[dot.ID]*inventory.InventoryVoucher)
	mapSupplier := make(map[dot.ID]*suppliering.ShopSupplier)
	for _, purchaseOrder := range purchaseOrders {
		supplierIDs = append(supplierIDs, purchaseOrder.SupplierID)
		purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
	}

	if len(supplierIDs) != 0 {
		listSuppliersQuery := &suppliering.ListSuppliersByIDsQuery{
			IDs:    supplierIDs,
			ShopID: args.ShopID,
		}
		if err := q.supplierQuery.Dispatch(ctx, listSuppliersQuery); err != nil {
			return nil, err
		}

		for _, supplier := range listSuppliersQuery.Result.Suppliers {
			if _, ok := mapSupplier[supplier.ID]; !ok {
				mapSupplier[supplier.ID] = supplier
			}
		}

		for _, purchaseOrder := range purchaseOrders {
			if _, ok := mapSupplier[purchaseOrder.SupplierID]; !ok {
				purchaseOrder.Supplier.Deleted = true
			}
		}
	}

	if len(purchaseOrderIDs) != 0 {
		getInventoryVouchersQuery := &inventory.GetInventoryVouchersByRefIDsQuery{
			RefIDs: purchaseOrderIDs,
			ShopID: args.ShopID,
		}
		if err := q.inventoryVoucherQuery.Dispatch(ctx, getInventoryVouchersQuery); err != nil {
			return nil, err
		}

		for _, inventoryVoucher := range getInventoryVouchersQuery.Result.InventoryVoucher {
			if _, ok := mapPurchaseOrderIDAndInventoryVoucher[inventoryVoucher.RefID]; !ok {
				mapPurchaseOrderIDAndInventoryVoucher[inventoryVoucher.RefID] = inventoryVoucher
			}
		}

		for _, purchaseOrder := range purchaseOrders {
			if _, ok := mapPurchaseOrderIDAndInventoryVoucher[purchaseOrder.ID]; ok {
				purchaseOrder.InventoryVoucher = mapPurchaseOrderIDAndInventoryVoucher[purchaseOrder.ID]
			} else {
				purchaseOrder.InventoryVoucher = nil
			}
		}
	}

	if err := q.addPaidAmount(ctx, args.ShopID, purchaseOrders); err != nil {
		return nil, err
	}

	return &purchaseorder.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
	}, nil
}

func (q *PurchaseOrderQuery) ListPurchaseOrdersByReceiptID(
	ctx context.Context, receiptID, shopID dot.ID,
) (*purchaseorder.PurchaseOrdersResponse, error) {
	getReceipt := &receipting.GetReceiptByIDQuery{
		ID:     receiptID,
		ShopID: shopID,
	}
	if err := q.receiptQuery.Dispatch(ctx, getReceipt); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y phi???u").
			Throw()
	}

	purchaseOrders, err := q.store(ctx).IDs(getReceipt.Result.RefIDs...).ShopID(shopID).ListPurchaseOrders()
	if err != nil {
		return nil, err
	}
	return &purchaseorder.PurchaseOrdersResponse{PurchaseOrders: purchaseOrders}, nil
}

func (q *PurchaseOrderQuery) GetPurchaseOrdersByIDs(
	ctx context.Context, IDs []dot.ID, ShopID dot.ID,
) (*purchaseorder.PurchaseOrdersResponse, error) {
	query := q.store(ctx).ShopID(ShopID).IDs(IDs...)
	purchaseOrders, err := query.ListPurchaseOrders()
	if err != nil {
		return nil, err
	}

	if err := q.addPaidAmount(ctx, ShopID, purchaseOrders); err != nil {
		return nil, err
	}

	return &purchaseorder.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
	}, nil
}

func (q *PurchaseOrderQuery) addPaidAmount(ctx context.Context, shopID dot.ID, purchaseOrders []*purchaseorder.PurchaseOrder) error {
	mapPurchaseOrderIDAndPaidAmount := make(map[dot.ID]int)
	var purchaseOrderIDs []dot.ID
	for _, purchaseOrder := range purchaseOrders {
		mapPurchaseOrderIDAndPaidAmount[purchaseOrder.ID] = 0
		purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
	}
	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  shopID,
		RefIDs:  purchaseOrderIDs,
		RefType: receipt_ref.PurchaseOrder,
		Status:  int(status3.P),
	}
	if err := q.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}
	receipts := listReceiptsQuery.Result.Receipts
	for _, receipt := range receipts {
		for _, line := range receipt.Lines {
			if _, ok := mapPurchaseOrderIDAndPaidAmount[line.RefID]; ok {
				mapPurchaseOrderIDAndPaidAmount[line.RefID] += line.Amount
			}
		}
	}
	for _, purchaseOrder := range purchaseOrders {
		purchaseOrder.PaidAmount = mapPurchaseOrderIDAndPaidAmount[purchaseOrder.ID]
	}
	return nil
}

func (q *PurchaseOrderQuery) ListPurchaseOrdersBySupplierIDsAndStatuses(
	ctx context.Context, shopID dot.ID, supplierIDs []dot.ID, statuses []status3.Status,
) (*purchaseorder.PurchaseOrdersResponse, error) {
	query := q.store(ctx).ShopID(shopID).SupplierIDs(supplierIDs...)
	if len(statuses) != 0 {
		query.Statuses(statuses...)
	}
	purchaseOrders, err := query.ListPurchaseOrders()
	if err != nil {
		return nil, err
	}
	return &purchaseorder.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
	}, nil
}
