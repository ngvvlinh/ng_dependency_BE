package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/etop"
	"o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
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
	shop2 "o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/session"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	"o.o/capi/dot"
)

type OrderService struct {
	session.Session

	OrderAggr     ordering.CommandBus
	CustomerQuery customering.QueryBus
	OrderQuery    ordering.QueryBus
	ReceiptQuery  receipting.QueryBus
	OrderLogic    *logicorder.OrderLogic
}

func (s *OrderService) Clone() shop.OrderService { res := *s; return &res }

func (s *OrderService) GetOrder(ctx context.Context, q *pbcm.IDRequest) (*types.Order, error) {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             s.SS.Shop().ID,
		PartnerID:          s.SS.CtxPartner().GetID(),
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := convertpb.PbOrder(query.Result.Order, nil, account_tag.TagShop)
	result.ShopName = s.SS.Shop().Name
	result.Fulfillments = convertpb.XPbFulfillments(query.Result.XFulfillments, account_tag.TagShop)
	if err := s.checkValidateCustomer(ctx, []*types.Order{result}); err != nil {
		return nil, err
	}
	if err := s.addReceivedAmountToOrders(ctx, s.SS.Shop().ID, []*types.Order{result}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderService) GetOrders(ctx context.Context, q *shop.GetOrdersRequest) (*types.OrdersResponse, error) {
	claim, partner := s.SS.Claim(), s.SS.CtxPartner()
	shopIDs, err := api.MixAccount(claim, q.Mixed)
	if err != nil {
		return nil, err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   shopIDs,
		PartnerID: partner.GetID(),
		Paging:    paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, account_tag.TagShop, query.Result.Shops),
	}
	if err := s.checkValidateCustomer(ctx, result.Orders); err != nil {
		return nil, err
	}
	if err := s.addReceivedAmountToOrders(ctx, s.SS.Shop().ID, result.Orders); err != nil {
		return nil, err
	}

	return result, nil
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

func (s *OrderService) GetOrdersByIDs(ctx context.Context, q *etop.IDsRequest) (*types.OrdersResponse, error) {
	shopIDs, err := api.MixAccount(s.SS.Claim(), q.Mixed)
	if err != nil {
		return nil, err
	}

	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   shopIDs,
		PartnerID: s.SS.CtxPartner().GetID(),
		IDs:       q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, account_tag.TagShop, query.Result.Shops),
	}
	if err := s.addReceivedAmountToOrders(ctx, s.SS.Shop().ID, result.Orders); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *OrderService) GetOrdersByReceiptID(ctx context.Context, q *shop.GetOrdersByReceiptIDRequest) (*types.OrdersResponse, error) {
	shopID := s.SS.Shop().ID
	queryReceipt := &receipting.GetReceiptByIDQuery{
		ID:     q.ReceiptId,
		ShopID: shopID,
	}
	if err := s.ReceiptQuery.Dispatch(ctx, queryReceipt); err != nil {
		return nil, err
	}
	var arrOrderID []dot.ID
	for _, value := range queryReceipt.Result.Lines {
		arrOrderID = append(arrOrderID, value.RefID)
	}

	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   []dot.ID{shopID},
		PartnerID: s.SS.CtxPartner().GetID(),
		IDs:       arrOrderID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, account_tag.TagShop, query.Result.Shops),
	}

	if err := s.addReceivedAmountToOrders(ctx, s.SS.Shop().ID, result.Orders); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *OrderService) UpdateOrdersStatus(ctx context.Context, q *shop.UpdateOrdersStatusRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &ordermodelx.UpdateOrdersStatusCommand{
		ShopID:       s.SS.Shop().ID,
		PartnerID:    s.SS.CtxPartner().GetID(),
		OrderIDs:     q.Ids,
		ShopConfirm:  q.Confirm,
		CancelReason: q.CancelReason,
		Status:       q.Status.Wrap(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result.Updated}
	return result, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, q *types.CreateOrderRequest) (result *types.Order, err error) {
	if q.CustomerId == 0 && q.Customer == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin tên khách hàng, vui lòng kiểm tra lại.")
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
		s.SS.Shop().ID, customerKey,
		q.TotalAmount, q.BasketValue, q.ShopCod, dot.JoinIDs(variantIDs))

	res, cached, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "tạo đơn hàng",
		func() (interface{}, error) { return s.createOrder(ctx, q) })
	if err != nil {
		return nil, err
	}

	defer func() { result = res.(*types.Order) }()
	if !cached {
		return result, nil
	}

	// FIX(vu): https://github.com/etopvn/one/issues/1910
	//
	// User may want to retry when an order is cancelled by external events. We
	// allow recreating the order when the response is cached and the previous
	// order has status = N

	key2 := key + "/retry"
	res, _, err = shop2.Idempgroup.DoAndWrap(
		ctx, key2, 5*time.Second, "tạo đơn hàng",
		func() (interface{}, error) {
			query := &ordermodelx.GetOrderQuery{
				OrderID:            res.(*types.Order).Id,
				ShopID:             s.SS.Shop().ID,
				PartnerID:          s.SS.CtxPartner().GetID(),
				IncludeFulfillment: true,
			}
			if _err := bus.Dispatch(ctx, query); _err != nil {
				return res, _err // keep the response
			}
			if query.Result.Order.Status != status5.N {
				return res, nil // keep the response
			}
			// release the old key and retry
			shop2.Idempgroup.ReleaseKey(key, "")
			_res, _, _err := shop2.Idempgroup.DoAndWrap(
				ctx, key, 30*time.Second, "tạo đơn hàng",
				func() (interface{}, error) { return s.createOrder(ctx, q) })
			return _res, _err // keep the response
		})
	return result, nil
}

func (s *OrderService) createOrder(ctx context.Context, q *types.CreateOrderRequest) (*types.Order, error) {
	result, err := s.OrderLogic.CreateOrder(ctx, s.SS.Shop(), s.SS.CtxPartner(), q, nil, s.SS.Claim().UserID)
	return result, err
}

func (s *OrderService) UpdateOrder(ctx context.Context, q *types.UpdateOrderRequest) (*types.Order, error) {
	result, err := s.OrderLogic.UpdateOrder(ctx, s.SS.Shop(), s.SS.CtxPartner(), q)
	return result, err
}

func (s *OrderService) CancelOrder(ctx context.Context, q *shop.CancelOrderRequest) (*types.OrderWithErrorsResponse, error) {
	key := fmt.Sprintf("cancelOrder %v-%v", s.SS.Shop().ID, q.OrderId)
	res, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 5*time.Second, "hủy đơn hàng",
		func() (interface{}, error) { return s.cancelOrder(ctx, q) })
	if err != nil {
		return nil, err
	}
	result := res.(*types.OrderWithErrorsResponse)
	return result, nil
}

func (s *OrderService) cancelOrder(ctx context.Context, q *shop.CancelOrderRequest) (*types.OrderWithErrorsResponse, error) {
	resp, err := s.OrderLogic.CancelOrder(ctx, s.SS.Claim().UserID, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, q.OrderId, q.CancelReason, q.AutoInventoryVoucher)
	return resp, err
}

func (s *OrderService) CompleteOrder(ctx context.Context, q *shop.OrderIDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &ordering.CompleteOrderCommand{
		OrderID: q.OrderId,
		ShopID:  s.SS.Shop().ID,
	}
	if err := s.OrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil

}

func (s *OrderService) ConfirmOrder(ctx context.Context, q *shop.ConfirmOrderRequest) (*types.Order, error) {
	key := fmt.Sprintf("ConfirmOrder %v-%v", s.SS.Shop().ID, q.OrderId)
	res, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "Xác nhận đơn hàng",
		func() (interface{}, error) { return s.confirmOrder(ctx, q) })
	if err != nil {
		return nil, err
	}
	result := res.(*types.Order)
	return result, nil
}

func (s *OrderService) confirmOrder(ctx context.Context, q *shop.ConfirmOrderRequest) (*types.Order, error) {
	resp, err := s.OrderLogic.ConfirmOrder(ctx, s.SS.Claim().UserID, s.SS.Shop(), q)
	return resp, err
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

func (s *OrderService) ConfirmOrderAndCreateFulfillments(ctx context.Context, q *shop.OrderIDRequest) (resp *types.OrderWithErrorsResponse, _err error) {
	key := fmt.Sprintf("ConfirmOrderAndCreateFulfillments %v-%v", s.SS.Shop().ID, q.OrderId)
	res, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "xác nhận đơn hàng",
		func() (interface{}, error) { return s.confirmOrderAndCreateFulfillments(ctx, q) })

	if err != nil {
		return nil, err
	}
	result := res.(*types.OrderWithErrorsResponse)
	return result, nil
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

func (s *OrderService) confirmOrderAndCreateFulfillments(ctx context.Context, q *shop.OrderIDRequest) (resp *types.OrderWithErrorsResponse, _err error) {
	resp, err := s.OrderLogic.ConfirmOrderAndCreateFulfillments(ctx, s.SS.Claim().UserID, s.SS.Shop(), s.SS.Claim().AuthPartnerID, q)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *OrderService) UpdateOrderPaymentStatus(ctx context.Context, q *shop.UpdateOrderPaymentStatusRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &ordermodelx.UpdateOrderPaymentStatusCommand{
		ShopID:        s.SS.Shop().ID,
		OrderID:       q.OrderId,
		PaymentStatus: q.Status,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}

func (s *OrderService) UpdateOrderShippingInfo(ctx context.Context, q *shop.UpdateOrderShippingInfoRequest) (*pbcm.UpdatedResponse, error) {
	shippingAddressModel, err := convertpb.OrderAddressToModel(q.ShippingAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	var order = new(ordermodel.Order)
	if err := convertpb.OrderShippingToModel(ctx, q.Shipping, order); err != nil {
		return nil, err
	}
	cmd := &ordermodelx.UpdateOrderShippingInfoCommand{
		ShopID:          s.SS.Shop().ID,
		OrderID:         q.OrderId,
		ShippingAddress: shippingAddressModel,
		Shipping:        order.ShopShipping,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}
