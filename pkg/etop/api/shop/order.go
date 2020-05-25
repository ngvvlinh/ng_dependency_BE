package shop

import (
	"context"
	"fmt"
	"strings"
	"time"

	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	ordermodel "o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type OrderService struct {
	OrderAggr     ordering.CommandBus
	CustomerQuery customering.QueryBus
	OrderQuery    ordering.QueryBus
	ReceiptQuery  receipting.QueryBus
}

func (s *OrderService) Clone() *OrderService { res := *s; return &res }

func (s *OrderService) GetOrder(ctx context.Context, q *GetOrderEndpoint) error {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             q.Context.Shop.ID,
		PartnerID:          q.CtxPartner.GetID(),
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = convertpb.PbOrder(query.Result.Order, nil, model.TagShop)
	q.Result.ShopName = q.Context.Shop.Name
	q.Result.Fulfillments = convertpb.XPbFulfillments(query.Result.XFulfillments, model.TagShop)
	if err := s.checkValidateCustomer(ctx, []*types.Order{q.Result}); err != nil {
		return err
	}

	if err := s.addReceivedAmountToOrders(ctx, q.Context.Shop.ID, []*types.Order{q.Result}); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, q *GetOrdersEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		Paging:    paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	if err := s.checkValidateCustomer(ctx, q.Result.Orders); err != nil {
		return err
	}
	if err := s.addReceivedAmountToOrders(ctx, q.Context.Shop.ID, q.Result.Orders); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) checkValidateCustomer(ctx context.Context, orders []*types.Order) error {
	if orders == nil {
		return nil
	}
	customerIDs := make([]dot.ID, 0, len(orders))
	shopIDs := make([]dot.ID, 0, len(orders))
	for _, order := range orders {
		customerIDs = append(customerIDs, order.CustomerId)
		shopIDs = append(shopIDs, order.ShopId)
	}
	queryCustomers := &customering.ListCustomersByIDsQuery{
		IDs:     customerIDs,
		ShopIDs: shopIDs,
	}
	if err := s.CustomerQuery.Dispatch(ctx, queryCustomers); err != nil {
		return err
	}
	customers := queryCustomers.Result.Customers
	var mapCustomerValidate = make(map[dot.ID]bool)
	for _, customer := range customers {
		mapCustomerValidate[customer.ID] = true
	}
	for _, order := range orders {
		if order.CustomerId == 0 || order.CustomerId == customering.CustomerAnonymous {
			continue
		}
		if !mapCustomerValidate[order.CustomerId] {
			order.Customer.Deleted = true
		}
	}
	return nil
}

func (s *OrderService) GetOrdersByIDs(ctx context.Context, q *GetOrdersByIDsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		IDs:       q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}

	if err := s.addReceivedAmountToOrders(ctx, q.Context.Shop.ID, q.Result.Orders); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrdersByReceiptID(ctx context.Context, q *GetOrdersByReceiptIDEndpoint) error {
	shopID := q.Context.Shop.ID
	queryReceipt := &receipting.GetReceiptByIDQuery{
		ID:     q.ReceiptId,
		ShopID: shopID,
	}
	if err := s.ReceiptQuery.Dispatch(ctx, queryReceipt); err != nil {
		return err
	}
	var arrOrderID []dot.ID
	for _, value := range queryReceipt.Result.Lines {
		arrOrderID = append(arrOrderID, value.RefID)
	}

	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   []dot.ID{shopID},
		PartnerID: q.CtxPartner.GetID(),
		IDs:       arrOrderID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}

	if err := s.addReceivedAmountToOrders(ctx, q.Context.Shop.ID, q.Result.Orders); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrdersStatus(ctx context.Context, q *UpdateOrdersStatusEndpoint) error {
	cmd := &ordermodelx.UpdateOrdersStatusCommand{
		ShopID:       q.Context.Shop.ID,
		PartnerID:    q.CtxPartner.GetID(),
		OrderIDs:     q.Ids,
		ShopConfirm:  q.Confirm,
		CancelReason: q.CancelReason,
		Status:       q.Status.Wrap(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result.Updated}
	return nil
}

func (s *OrderService) CreateOrder(ctx context.Context, q *CreateOrderEndpoint) error {
	if q.CustomerId == 0 && q.Customer == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin tên khách hàng, vui lòng kiểm tra lại.")
	}
	customerKey := q.CustomerId.String()
	if q.Customer != nil {
		customerKey = q.Customer.FullName
		if phone := strings.TrimSpace(q.Customer.Phone); phone != "" {
			customerKey = phone
		}
	}

	var variantIDs []dot.ID
	for _, line := range q.Lines {
		variantIDs = append(variantIDs, line.VariantId)
	}
	key := fmt.Sprintf("CreateOrder %v-%v-%v-%v-%v-%v",
		q.Context.Shop.ID, customerKey,
		q.TotalAmount, q.BasketValue, q.ShopCod, dot.JoinIDs(variantIDs))

	res, cached, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "tạo đơn hàng",
		func() (interface{}, error) { return s.createOrder(ctx, q) })
	if err != nil {
		return err
	}

	defer func() { q.Result = res.(*CreateOrderEndpoint).Result }()
	if !cached {
		return nil
	}

	// FIX(vu): https://github.com/etopvn/one/issues/1910
	//
	// User may want to retry when an order is cancelled by external events. We
	// allow recreating the order when the response is cached and the previous
	// order has status = N

	key2 := key + "/retry"
	res, _, err = idempgroup.DoAndWrap(
		ctx, key2, 5*time.Second, "tạo đơn hàng",
		func() (interface{}, error) {
			query := &ordermodelx.GetOrderQuery{
				OrderID:            res.(*CreateOrderEndpoint).Result.Id,
				ShopID:             q.Context.Shop.ID,
				PartnerID:          q.CtxPartner.GetID(),
				IncludeFulfillment: true,
			}
			if _err := bus.Dispatch(ctx, query); _err != nil {
				return res, _err // keep the response
			}
			if query.Result.Order.Status != status5.N {
				return res, nil // keep the response
			}

			// release the old key and retry
			idempgroup.ReleaseKey(key, "")
			_res, _, _err := idempgroup.DoAndWrap(
				ctx, key, 30*time.Second, "tạo đơn hàng",
				func() (interface{}, error) { return s.createOrder(ctx, q) })
			return _res, _err // keep the response
		})
	return nil
}

func (s *OrderService) createOrder(ctx context.Context, q *CreateOrderEndpoint) (*CreateOrderEndpoint, error) {
	result, err := logicorder.CreateOrder(ctx, &q.Context, q.CtxPartner, q.CreateOrderRequest, nil, q.Context.UserID)
	q.Result = result
	return q, err
}

func (s *OrderService) UpdateOrder(ctx context.Context, q *UpdateOrderEndpoint) error {
	result, err := logicorder.UpdateOrder(ctx, &q.Context, q.CtxPartner, q.UpdateOrderRequest)
	q.Result = result
	return err
}

func (s *OrderService) CancelOrder(ctx context.Context, q *CancelOrderEndpoint) error {
	key := fmt.Sprintf("cancelOrder %v-%v", q.Context.Shop.ID, q.OrderId)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 5*time.Second, "hủy đơn hàng",
		func() (interface{}, error) { return s.cancelOrder(ctx, q) })

	if err != nil {
		return err
	}
	q.Result = res.(*CancelOrderEndpoint).Result
	return nil
}

func (s *OrderService) cancelOrder(ctx context.Context, q *CancelOrderEndpoint) (*CancelOrderEndpoint, error) {
	resp, err := logicorder.CancelOrder(ctx, q.Context.UserID, q.Context.Shop.ID, q.Context.AuthPartnerID, q.OrderId, q.CancelReason, q.AutoInventoryVoucher)
	q.Result = resp
	return q, err
}

func (s *OrderService) CompleteOrder(ctx context.Context, q *CompleteOrderEndpoint) error {
	cmd := &ordering.CompleteOrderCommand{
		OrderID: q.OrderId,
		ShopID:  q.Context.Shop.ID,
	}
	if err := s.OrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil

}

func (s *OrderService) ConfirmOrder(ctx context.Context, q *ConfirmOrderEndpoint) error {
	key := fmt.Sprintf("ConfirmOrder %v-%v", q.Context.Shop.ID, q.OrderId)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "Xác nhận đơn hàng",
		func() (interface{}, error) { return s.confirmOrder(ctx, q) })
	if err != nil {
		return err
	}
	q.Result = res.(*ConfirmOrderEndpoint).Result
	return nil
}

func (s *OrderService) confirmOrder(ctx context.Context, q *ConfirmOrderEndpoint) (_ *ConfirmOrderEndpoint, _err error) {
	resp, err := logicorder.ConfirmOrder(ctx, q.Context.UserID, q.Context.Shop, q.ConfirmOrderRequest)
	if err != nil {
		return q, err
	}
	q.Result = resp
	return q, nil
}

/*
	0. Check idempotent
	1. Get order
	2. UpdateInfo order to confirm
	3. Check if fulfillment exists
	   no:  Create fulfillments with status sync_status 0
	   yes: Check existing fullfillments sync_status
			if any sync_status 0:
				Create fulfillments with status sync_status 0

	4. UpdateInfo fulfillment information and status from GHN
*/

func (s *OrderService) ConfirmOrderAndCreateFulfillments(ctx context.Context, q *ConfirmOrderAndCreateFulfillmentsEndpoint) (_err error) {
	key := fmt.Sprintf("ConfirmOrderAndCreateFulfillments %v-%v", q.Context.Shop.ID, q.OrderId)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "xác nhận đơn hàng",
		func() (interface{}, error) { return s.confirmOrderAndCreateFulfillments(ctx, q) })

	if err != nil {
		return err
	}
	q.Result = res.(*ConfirmOrderAndCreateFulfillmentsEndpoint).Result
	return err
}

func (s *OrderService) addReceivedAmountToOrders(ctx context.Context, shopID dot.ID, orders []*types.Order) error {
	var orderIDs []dot.ID
	mOrderIDsAndReceivedAmounts := make(map[dot.ID]int)

	for _, order := range orders {
		mOrderIDsAndReceivedAmounts[order.Id] = 0
		orderIDs = append(orderIDs, order.Id)
	}

	listReceiptsByRefIDsAndStatusQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  shopID,
		RefIDs:  orderIDs,
		RefType: receipt_ref.Order,
		Status:  int(status3.P),
	}
	if err := s.ReceiptQuery.Dispatch(ctx, listReceiptsByRefIDsAndStatusQuery); err != nil {
		return err
	}

	for _, receipt := range listReceiptsByRefIDsAndStatusQuery.Result.Receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mOrderIDsAndReceivedAmounts[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipt_type.Receipt:
					mOrderIDsAndReceivedAmounts[receiptLine.RefID] += receiptLine.Amount
				case receipt_type.Payment:
					mOrderIDsAndReceivedAmounts[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}

	for _, order := range orders {
		order.ReceivedAmount = mOrderIDsAndReceivedAmounts[order.Id]
	}

	return nil
}

func (s *OrderService) confirmOrderAndCreateFulfillments(ctx context.Context, q *ConfirmOrderAndCreateFulfillmentsEndpoint) (_ *ConfirmOrderAndCreateFulfillmentsEndpoint, _err error) {
	resp, err := logicorder.ConfirmOrderAndCreateFulfillments(ctx, q.Context.UserID, q.Context.Shop, q.Context.AuthPartnerID, q.OrderIDRequest)
	if err != nil {
		return q, err
	}
	q.Result = resp

	return q, nil
}

func (s *OrderService) UpdateOrderPaymentStatus(ctx context.Context, q *UpdateOrderPaymentStatusEndpoint) error {
	cmd := &ordermodelx.UpdateOrderPaymentStatusCommand{
		ShopID:        q.Context.Shop.ID,
		OrderID:       q.OrderId,
		PaymentStatus: q.Status,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *OrderService) UpdateOrderShippingInfo(ctx context.Context, q *UpdateOrderShippingInfoEndpoint) error {
	shippingAddressModel, err := convertpb.OrderAddressToModel(q.ShippingAddress)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	var order = new(ordermodel.Order)
	if err := convertpb.OrderShippingToModel(ctx, q.Shipping, order); err != nil {
		return err
	}
	cmd := &ordermodelx.UpdateOrderShippingInfoCommand{
		ShopID:          q.Context.Shop.ID,
		OrderID:         q.OrderId,
		ShippingAddress: shippingAddressModel,
		Shipping:        order.ShopShipping,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}
