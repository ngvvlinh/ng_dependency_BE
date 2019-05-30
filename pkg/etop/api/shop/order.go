package shop

import (
	"context"
	"fmt"
	"strings"
	"time"

	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	ordermodelx "etop.vn/backend/pkg/services/ordering/modelx"
	shipmodelx "etop.vn/backend/pkg/services/shipping/modelx"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandlers("api",
		CancelFulfillment,
		CancelOrder,
		ConfirmOrderAndCreateFulfillments,
		ConfirmOrdersAndCreateFulfillments,
		CreateFulfillmentsForOrder,
		CreateOrder,
		GetExternalShippingServices,
		GetPublicExternalShippingServices,
		GetFulfillment,
		GetOrder,
		GetOrders,
		GetOrdersByIDs,
		UpdateOrder,
		UpdateOrdersStatus,
		GetFulfillments,
		GetPublicFulfillment,
		UpdateFulfillmentsShippingState,
		UpdateOrderPaymentStatus,
	)
}

func GetOrder(ctx context.Context, q *wrapshop.GetOrderEndpoint) error {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             q.Context.Shop.ID,
		PartnerID:          q.CtxPartner.GetID(),
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = pborder.PbOrder(query.Result.Order, nil, model.TagShop)
	q.Result.ShopName = q.Context.Shop.Name
	q.Result.Fulfillments = pborder.PbFulfillments(query.Result.Fulfillments, model.TagShop)
	return nil
}

func GetOrders(ctx context.Context, q *wrapshop.GetOrdersEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	paging := q.Paging.CMPaging()
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		Paging:    paging,
		Filters:   pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.OrdersResponse{
		Paging: pbcm.PbPageInfo(paging, int32(query.Result.Total)),
		Orders: pborder.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	return nil
}

func GetOrdersByIDs(ctx context.Context, q *wrapshop.GetOrdersByIDsEndpoint) error {
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
	q.Result = &pborder.OrdersResponse{
		Orders: pborder.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	return nil
}

func UpdateOrdersStatus(ctx context.Context, q *wrapshop.UpdateOrdersStatusEndpoint) error {
	cmd := &ordermodelx.UpdateOrdersStatusCommand{
		ShopID:       q.Context.Shop.ID,
		PartnerID:    q.CtxPartner.GetID(),
		OrderIDs:     q.Ids,
		ShopConfirm:  q.Confirm.ToModel(),
		CancelReason: q.CancelReason,
		Status:       q.Status.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result.Updated)}
	return nil
}

func CreateOrder(ctx context.Context, q *wrapshop.CreateOrderEndpoint) error {
	if q.ShippingAddress == nil || q.ShippingAddress.FullName == "" {
		return cm.Error(cm.InvalidArgument, "Thiếu thông tin tên khách hàng, vui lòng kiểm tra lại.", nil)
	}
	customerKey := q.ShippingAddress.FullName
	if phone := strings.TrimSpace(q.ShippingAddress.Phone); phone != "" {
		customerKey = phone
	}
	key := fmt.Sprintf("CreateOrder %v-%v", q.Context.Shop.ID, customerKey)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			return createOrder(ctx, q)
		}, "tạo đơn hàng")

	if err != nil {
		return err
	}
	q.Result = res.(*wrapshop.CreateOrderEndpoint).Result
	return nil
}

func createOrder(ctx context.Context, q *wrapshop.CreateOrderEndpoint) (*wrapshop.CreateOrderEndpoint, error) {
	result, err := logicorder.CreateOrder(ctx, &q.Context, q.CtxPartner, q.CreateOrderRequest)
	q.Result = result
	return q, err
}

func UpdateOrder(ctx context.Context, q *wrapshop.UpdateOrderEndpoint) error {
	result, err := logicorder.UpdateOrder(ctx, &q.Context, q.CtxPartner, q.UpdateOrderRequest)
	q.Result = result
	return err
}

func CancelFulfillment(ctx context.Context, q *wrapshop.CancelFulfillmentEndpoint) error {
	return cm.ErrTODO
}

func CancelOrder(ctx context.Context, q *wrapshop.CancelOrderEndpoint) error {
	key := fmt.Sprintf("cancelOrder %v-%v", q.Context.Shop.ID, q.OrderId)
	res, err := idempgroup.DoAndWrap(key, 5*time.Second,
		func() (interface{}, error) {
			return cancelOrder(ctx, q)
		}, "hủy đơn hàng")

	if err != nil {
		return err
	}
	q.Result = res.(*wrapshop.CancelOrderEndpoint).Result
	return nil
}

func cancelOrder(ctx context.Context, q *wrapshop.CancelOrderEndpoint) (*wrapshop.CancelOrderEndpoint, error) {
	resp, err := logicorder.CancelOrder(ctx, q.Context.Shop.ID, q.Context.AuthPartnerID, q.OrderId, q.CancelReason)
	q.Result = resp
	return q, err
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

func ConfirmOrderAndCreateFulfillments(ctx context.Context, q *wrapshop.ConfirmOrderAndCreateFulfillmentsEndpoint) (_err error) {
	key := fmt.Sprintf("ConfirmOrderAndCreateFulfillments %v-%v", q.Context.Shop.ID, q.OrderId)
	res, err := idempgroup.DoAndWrap(key, 10*time.Second,
		func() (interface{}, error) {
			return confirmOrderAndCreateFulfillments(ctx, q)
		}, "xác nhận đơn hàng")

	if err != nil {
		return err
	}
	q.Result = res.(*wrapshop.ConfirmOrderAndCreateFulfillmentsEndpoint).Result
	return err
}

func confirmOrderAndCreateFulfillments(ctx context.Context, q *wrapshop.ConfirmOrderAndCreateFulfillmentsEndpoint) (_ *wrapshop.ConfirmOrderAndCreateFulfillmentsEndpoint, _err error) {
	resp, err := logicorder.ConfirmOrderAndCreateFulfillments(ctx, q.Context.Shop, q.Context.AuthPartnerID, q.OrderIDRequest)
	if err != nil {
		return q, err
	}
	q.Result = resp
	return q, nil
}

func ConfirmOrdersAndCreateFulfillments(ctx context.Context, q *wrapshop.ConfirmOrdersAndCreateFulfillmentsEndpoint) error {
	return cm.ErrTODO
}

func CreateFulfillmentsForOrder(ctx context.Context, q *wrapshop.CreateFulfillmentsForOrderEndpoint) error {
	return cm.ErrTODO
}

func GetFulfillment(ctx context.Context, q *wrapshop.GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		ShopID:        q.Context.Shop.ID,
		PartnerID:     q.CtxPartner.GetID(),
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbFulfillment(query.Result.Fulfillment, model.TagShop, query.Result.Shop, query.Result.Order)
	return nil
}

func GetFulfillments(ctx context.Context, q *wrapshop.GetFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	paging := q.Paging.CMPaging()
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		OrderID:   q.OrderId,
		Status:    q.Status.ToModel(),
		Paging:    paging,
		Filters:   pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.FulfillmentsResponse{
		Fulfillments: pborder.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagShop),
		Paging:       pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func GetExternalShippingServices(ctx context.Context, q *wrapshop.GetExternalShippingServicesEndpoint) error {
	resp, err := shippingCtrl.GetExternalShippingServices(ctx, q.Context.Shop.ID, q.GetExternalShippingServicesRequest)
	q.Result = &pborder.GetExternalShippingServicesResponse{
		Services: pborder.PbAvailableShippingServices(resp),
	}
	return err
}

func GetPublicExternalShippingServices(ctx context.Context, q *wrapshop.GetPublicExternalShippingServicesEndpoint) error {
	resp, err := shippingCtrl.GetExternalShippingServices(ctx, model.EtopAccountID, q.GetExternalShippingServicesRequest)
	q.Result = &pborder.GetExternalShippingServicesResponse{
		Services: pborder.PbAvailableShippingServices(resp),
	}
	return err
}

func GetPublicFulfillment(ctx context.Context, q *wrapshop.GetPublicFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: q.Code,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbPublicFulfillment(query.Result)
	return nil
}

func UpdateFulfillmentsShippingState(ctx context.Context, q *wrapshop.UpdateFulfillmentsShippingStateEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &shipmodelx.UpdateFulfillmentsShippingStateCommand{
		ShopID:        shopID,
		IDs:           q.Ids,
		ShippingState: q.ShippingState.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func UpdateOrderPaymentStatus(ctx context.Context, q *wrapshop.UpdateOrderPaymentStatusEndpoint) error {
	cmd := &ordermodelx.UpdateOrderPaymentStatusCommand{
		ShopID:  q.Context.Shop.ID,
		OrderID: q.OrderId,
		Status:  q.Status.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}
