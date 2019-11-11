package query

import (
	"context"

	"etop.vn/api/main/inventory"

	"etop.vn/api/shopping/suppliering"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/capi"

	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/purchaseorder/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ purchaseorder.QueryService = &PurchaseOrderQuery{}

type PurchaseOrderQuery struct {
	db                    *cmsql.Database
	store                 sqlstore.PurchaseOrderStoreFactory
	eventBus              capi.EventBus
	supplierQuery         suppliering.QueryBus
	inventoryVoucherQuery inventory.QueryBus
}

func NewPurchaseOrderQuery(
	database *cmsql.Database, eventBus capi.EventBus,
	supplierQ suppliering.QueryBus, inventoryVoucherQ inventory.QueryBus,
) *PurchaseOrderQuery {
	return &PurchaseOrderQuery{
		db:                    database,
		store:                 sqlstore.NewPurchaseOrderStore(database),
		eventBus:              eventBus,
		supplierQuery:         supplierQ,
		inventoryVoucherQuery: inventoryVoucherQ,
	}
}

func (q *PurchaseOrderQuery) MessageBus() purchaseorder.QueryBus {
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
		RefIDs: []int64{purchaseOrder.ID},
		ShopID: args.ShopID,
	}
	if err := q.inventoryVoucherQuery.Dispatch(ctx, getInventoryVouchersQuery); err != nil {
		return nil, err
	}

	if len(getInventoryVouchersQuery.Result.InventoryVoucher) != 0 {
		purchaseOrder.InventoryVoucher = getInventoryVouchersQuery.Result.InventoryVoucher[0]
	}

	return purchaseOrder, nil
}

func (q *PurchaseOrderQuery) ListPurchaseOrders(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*purchaseorder.PurchaseOrdersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Filters(args.Filters)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	purchaseOrders, err := query.Paging(args.Paging).ListPurchaseOrders()
	if err != nil {
		return nil, err
	}

	var supplierIDs, purchaseOrderIDs []int64
	mapPurchaseOrderIDAndInventoryVoucher := make(map[int64]*inventory.InventoryVoucher)
	mapSupplier := make(map[int64]*suppliering.ShopSupplier)
	for _, purchaseOrder := range purchaseOrders {
		supplierIDs = append(supplierIDs, purchaseOrder.SupplierID)
		purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
		mapPurchaseOrderIDAndInventoryVoucher[purchaseOrder.ID] = &inventory.InventoryVoucher{}
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
			if _, ok := mapPurchaseOrderIDAndInventoryVoucher[inventoryVoucher.RefID]; ok {
				mapPurchaseOrderIDAndInventoryVoucher[inventoryVoucher.RefID] = inventoryVoucher
			}
		}

		for _, purchaseOrder := range purchaseOrders {
			if _, ok := mapPurchaseOrderIDAndInventoryVoucher[purchaseOrder.ID]; ok {
				purchaseOrder.InventoryVoucher = mapPurchaseOrderIDAndInventoryVoucher[purchaseOrder.ID]
			}
		}
	}

	return &purchaseorder.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
		Count:          int32(count),
	}, nil
}
