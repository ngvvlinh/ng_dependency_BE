package shop

import (
	"context"
	"fmt"
	"strings"
	"time"

	"etop.vn/api/top/types/etc/status3"

	"etop.vn/api/main/receipting"
	"etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		orderService.CancelOrder,
		orderService.ConfirmOrderAndCreateFulfillments,
		orderService.CreateOrder,
		fulfillmentService.GetExternalShippingServices,
		fulfillmentService.GetPublicExternalShippingServices,
		fulfillmentService.GetFulfillment,
		orderService.GetOrder,
		orderService.GetOrders,
		orderService.GetOrdersByIDs,
		orderService.UpdateOrder,
		orderService.UpdateOrdersStatus,
		fulfillmentService.GetFulfillments,
		fulfillmentService.GetPublicFulfillment,
		fulfillmentService.UpdateFulfillmentsShippingState,
		orderService.UpdateOrderPaymentStatus,
		orderService.GetOrdersByReceiptID,
	)
}

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
		Paging: cmapi.PbPageInfo(paging, query.Result.Total),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}

	if err := s.addReceivedAmountToOrders(ctx, q.Context.Shop.ID, q.Result.Orders); err != nil {
		return err
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
	if err := receiptQuery.Dispatch(ctx, queryReceipt); err != nil {
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
	customerKey := ""
	if q.Customer != nil {
		customerKey = q.Customer.FullName
		if phone := strings.TrimSpace(q.Customer.Phone); phone != "" {
			customerKey = phone
		}
	} else {
		customerKey = q.CustomerId.String()
	}
	key := fmt.Sprintf("CreateOrder %v-%v", q.Context.Shop.ID, customerKey)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			return s.createOrder(ctx, q)
		}, "tạo đơn hàng")

	if err != nil {
		return err
	}

	q.Result = res.(*CreateOrderEndpoint).Result
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
	res, err := idempgroup.DoAndWrap(key, 5*time.Second,
		func() (interface{}, error) {
			return s.cancelOrder(ctx, q)
		}, "hủy đơn hàng")

	if err != nil {
		return err
	}
	q.Result = res.(*CancelOrderEndpoint).Result
	return nil
}

func (s *OrderService) cancelOrder(ctx context.Context, q *CancelOrderEndpoint) (*CancelOrderEndpoint, error) {
	resp, err := logicorder.CancelOrder(ctx, q.Context.Shop.ID, q.Context.AuthPartnerID, q.OrderId, q.CancelReason, q.AutoInventoryVoucher)
	q.Result = resp
	return q, err
}

func (s *OrderService) ConfirmOrder(ctx context.Context, q *ConfirmOrderEndpoint) error {
	key := fmt.Sprintf("ConfirmOrder %v-%v", q.Context.Shop.ID, q.OrderId)
	res, err := idempgroup.DoAndWrap(key, 10*time.Second,
		func() (interface{}, error) {
			return s.confirmOrder(ctx, q)
		}, "Xác nhận đơn hàng")
	if err != nil {
		return err
	}
	q.Result = res.(*ConfirmOrderEndpoint).Result
	return nil
}

func (s *OrderService) confirmOrder(ctx context.Context, q *ConfirmOrderEndpoint) (_ *ConfirmOrderEndpoint, _err error) {
	resp, err := logicorder.ConfirmOrder(ctx, q.Context.Shop, q.ConfirmOrderRequest)
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
	res, err := idempgroup.DoAndWrap(key, 10*time.Second,
		func() (interface{}, error) {
			return s.confirmOrderAndCreateFulfillments(ctx, q)
		}, "xác nhận đơn hàng")

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
		RefType: receipting.ReceiptRefTypeOrder,
		Status:  int(status3.P),
	}
	if err := receiptQuery.Dispatch(ctx, listReceiptsByRefIDsAndStatusQuery); err != nil {
		return err
	}

	for _, receipt := range listReceiptsByRefIDsAndStatusQuery.Result.Receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mOrderIDsAndReceivedAmounts[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipting.ReceiptTypeReceipt:
					mOrderIDsAndReceivedAmounts[receiptLine.RefID] += receiptLine.Amount
				case receipting.ReceiptTypePayment:
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
	resp, err := logicorder.ConfirmOrderAndCreateFulfillments(ctx, q.Context.Shop, q.Context.AuthPartnerID, q.OrderIDRequest)
	if err != nil {
		return q, err
	}
	q.Result = resp

	return q, nil
}

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		ShopID:        q.Context.Shop.ID,
		PartnerID:     q.CtxPartner.GetID(),
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbFulfillment(query.Result.Fulfillment, model.TagShop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *GetFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		OrderID:   q.OrderId,
		Status:    q.Status,
		Paging:    paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagShop),
		Paging:       cmapi.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}

func (s *FulfillmentService) GetExternalShippingServices(ctx context.Context, q *GetExternalShippingServicesEndpoint) error {
	resp, err := shippingCtrl.GetExternalShippingServices(ctx, q.Context.Shop.ID, q.GetExternalShippingServicesRequest)
	q.Result = &types.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return err
}

func (s *FulfillmentService) GetPublicExternalShippingServices(ctx context.Context, q *GetPublicExternalShippingServicesEndpoint) error {
	resp, err := shippingCtrl.GetExternalShippingServices(ctx, model.EtopAccountID, q.GetExternalShippingServicesRequest)
	q.Result = &types.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return err
}

func (s *FulfillmentService) GetPublicFulfillment(ctx context.Context, q *GetPublicFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: q.Code,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbPublicFulfillment(query.Result)
	return nil
}

func (s *FulfillmentService) UpdateFulfillmentsShippingState(ctx context.Context, q *UpdateFulfillmentsShippingStateEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &shipmodelx.UpdateFulfillmentsShippingStateCommand{
		ShopID:        shopID,
		IDs:           q.Ids,
		ShippingState: convertpb.ShippingStateToModel(&q.ShippingState),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *OrderService) UpdateOrderPaymentStatus(ctx context.Context, q *UpdateOrderPaymentStatusEndpoint) error {
	cmd := &ordermodelx.UpdateOrderPaymentStatusCommand{
		ShopID:  q.Context.Shop.ID,
		OrderID: q.OrderId,
		Status:  q.Status,
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
	if err := convertpb.OrderShippingToModel(q.Shipping, order); err != nil {
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
