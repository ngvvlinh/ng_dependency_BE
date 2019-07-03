package gateway

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"

	"etop.vn/api/external/haravan"
	"etop.vn/api/external/haravan/gateway"
	"etop.vn/api/main/identity"
	pbsp "etop.vn/backend/pb/etop/etc/shipping_provider"
	pbtryon "etop.vn/backend/pb/etop/etc/try_on"
	pborder "etop.vn/backend/pb/etop/order"
	pbsource "etop.vn/backend/pb/etop/order/source"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/external/haravan/gateway/convert"
	identityconvert "etop.vn/backend/pkg/services/identity/convert"
	"etop.vn/backend/pkg/services/ordering/modelx"
	shipmodelx "etop.vn/backend/pkg/services/shipping/modelx"
)

var (
	_  gateway.Aggregate = &Aggregate{}
	ll                   = l.New()
)

type Aggregate struct {
	db           cmsql.Transactioner
	ShippingCtrl *shipping_provider.ProviderManager
	identityQS   identity.QueryBus
}

func NewAggregate(db cmsql.Database, providerManager *shipping_provider.ProviderManager, identityQuery identity.QueryBus) *Aggregate {
	return &Aggregate{
		db:           db,
		ShippingCtrl: providerManager,
		identityQS:   identityQuery,
	}
}

func (a *Aggregate) MessageBus() gateway.CommandBus {
	b := bus.New()
	return gateway.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) GetShippingRate(ctx context.Context, args *gateway.GetShippingRateRequestArgs) (*gateway.GetShippingRateResponse, error) {
	if args.Origin == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ lấy hàng")
	}
	if args.Destination == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ giao hàng")
	}

	// Haravan: default includeInsurance is false
	req := &pborder.GetExternalShippingServicesRequest{
		Provider:         0,
		Carrier:          0,
		FromProvince:     args.Origin.Province,
		FromDistrict:     args.Origin.District,
		ToProvince:       args.Destination.Province,
		ToDistrict:       args.Destination.District,
		Weight:           int32(args.TotalGrams),
		GrossWeight:      0,
		ChargeableWeight: 0,
		Value:            0,
		TotalCodAmount:   int32(args.CodAmount),
		CodAmount:        int32(args.CodAmount),
		BasketValue:      0,
		IncludeInsurance: nil,
	}

	services, err := a.ShippingCtrl.GetExternalShippingServices(ctx, args.EtopShopID, req)
	if err != nil {
		return nil, err
	}
	shippingRates := make([]*haravan.ShippingRate, len(services))
	for i, s := range services {
		serviceID := SHA256StringToInt32(s.ProviderServiceID)
		shippingRates[i] = &haravan.ShippingRate{
			ServiceID:       serviceID,
			ServiceName:     s.Name,
			ServiceCode:     s.ProviderServiceID,
			Currency:        "vnd",
			TotalPrice:      s.ServiceFee,
			PhoneRequired:   true,
			MinDeliveryDate: s.ExpectedDeliveryAt,
			MaxDeliveryDate: s.ExpectedDeliveryAt,
			Description:     "",
		}
	}

	return &gateway.GetShippingRateResponse{
		ShippingRates: shippingRates,
	}, nil
}

func (a *Aggregate) CreateOrder(ctx context.Context, args *gateway.CreateOrderRequestArgs) (_resp *gateway.CreateOrderResponse, _err error) {
	query := &identity.GetShopByIDQuery{
		ID: args.EtopShopID,
	}
	if err := a.identityQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shop := identityconvert.ShopToModel(query.Result.Shop)

	req := &pborder.GetExternalShippingServicesRequest{
		Provider:         0,
		Carrier:          0,
		FromProvince:     args.Origin.Province,
		FromDistrict:     args.Origin.District,
		ToProvince:       args.Destination.Province,
		ToDistrict:       args.Destination.District,
		Weight:           int32(args.TotalGrams),
		GrossWeight:      0,
		ChargeableWeight: 0,
		Value:            0,
		TotalCodAmount:   int32(args.CodAmount),
		CodAmount:        int32(args.CodAmount),
		BasketValue:      0,
		IncludeInsurance: nil,
	}
	services, err := a.ShippingCtrl.GetExternalShippingServices(ctx, args.EtopShopID, req)
	if err != nil {
		return nil, err
	}

	var service *model.AvailableShippingService
	for _, s := range services {
		serviceID := SHA256StringToInt32(s.ProviderServiceID)
		if serviceID == args.ShippingRateID {
			service = s
			break
		}
	}
	if service == nil {
		return nil, cm.Errorf(cm.NotFound, nil, "Không có gói vận chuyển phù hợp")
	}
	externalID := strconv.FormatInt(int64(args.ExternalOrderID), 10)
	totalValue := int32(getOrderValue(args.Items))
	req2 := &pborder.CreateOrderRequest{
		Source:          pbsource.Source_etop_pos,
		ExternalId:      externalID,
		ExternalCode:    args.ExternalCode,
		ExternalUrl:     "",
		PaymentMethod:   "",
		Customer:        convert.ToPbOrderCustomer(args.Origin),
		CustomerAddress: convert.ToPbOrderAddress(args.Origin),
		BillingAddress:  convert.ToPbOrderAddress(args.Origin),
		ShippingAddress: convert.ToPbOrderAddress(args.Destination),
		ShopAddress:     convert.ToPbOrderAddress(args.Origin),
		ShConfirm:       nil,
		Lines:           convert.ToPbCreateOrderLines(args.Items),
		Discounts:       nil,
		TotalItems:      int32(len(args.Items)),
		BasketValue:     totalValue,
		TotalWeight:     int32(args.TotalGrams),
		OrderDiscount:   0,
		TotalFee:        0,
		FeeLines:        nil,
		TotalDiscount:   nil,
		TotalAmount:     totalValue,
		OrderNote:       args.Note,
		ShippingNote:    args.Note,
		ShopCod:         0,
		ReferenceUrl:    "",
		ShopShipping:    nil,
		Shipping: &pborder.OrderShipping{
			ShAddress:           nil,
			XServiceId:          "",
			XShippingFee:        0,
			XServiceName:        "",
			PickupAddress:       convert.ToPbOrderAddress(args.Origin),
			ReturnAddress:       nil,
			ShippingServiceName: "",
			ShippingServiceCode: service.ProviderServiceID,
			ShippingServiceFee:  int32(service.ServiceFee),
			ShippingProvider:    pbsp.PbShippingProviderType(service.Provider),
			Carrier:             pbsp.PbShippingProviderType(service.Provider),
			IncludeInsurance:    false,
			TryOn:               pbtryon.TryOnCode_none,
			ShippingNote:        args.Note,
			CodAmount:           cm.PIntToInt32(args.CodAmount),
			Weight:              cm.PIntToInt32(args.TotalGrams),
			GrossWeight:         nil,
			ChargeableWeight:    nil,
		},
		GhnNoteCode: 0,
	}
	shopClaim := &claims.ShopClaim{
		Shop: shop,
	}
	resp, err := logicorder.CreateOrder(ctx, shopClaim, nil, req2)
	if err != nil {
		return nil, err
	}
	orderID := resp.Id

	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			_, err := logicorder.CancelOrder(ctx, args.EtopShopID, 0, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", err))
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()

	cfmResp, err := logicorder.ConfirmOrderAndCreateFulfillments(ctx, shop, 0, &pbshop.OrderIDRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, err
	}
	for _, err := range cfmResp.FulfillmentErrors {
		if err.Code != "ok" {
			return nil, err
		}
	}

	orderQuery := &modelx.GetOrderQuery{
		OrderID:            orderID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, orderQuery); err != nil {
		return nil, cm.MapError(err).
			Map(cm.NotFound, cm.Internal, "").
			Throw()
	}
	ffm := orderQuery.Result.Fulfillments[0]
	return &gateway.CreateOrderResponse{
		TrackingNumber: ffm.ShippingCode,
		ShippingFee:    int32(ffm.ShippingFeeShop),
		TrackingURL:    ffm.SelfURL(cm.MainSiteBaseURL(), model.TagShop),
		CodAmount:      int32(ffm.TotalCODAmount),
	}, nil
}

func getOrderValue(items []*haravan.Item) (total int) {
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return
}

func (a *Aggregate) GetOrder(ctx context.Context, args *gateway.GetOrderRequestArgs) (*gateway.GetOrderResponse, error) {
	ffmQuery := &shipmodelx.GetFulfillmentQuery{
		ShopID:       args.EtopShopID,
		ShippingCode: args.TrackingNumber,
	}
	if err := bus.Dispatch(ctx, ffmQuery); err != nil {
		return nil, err
	}
	ffm := ffmQuery.Result
	return &gateway.GetOrderResponse{
		TrackingNumber: ffm.ShippingCode,
		ShippingFee:    int32(ffm.ShippingFeeShop),
		TrackingURL:    ffm.SelfURL(cm.MainSiteBaseURL(), model.TagShop),
		CodAmount:      int32(ffm.TotalCODAmount),
		Status:         convert.ToFulfillmentStatus(ffm.ShippingState),
		CodStatus:      convert.ToCODStatus(ffm),
	}, nil
}

func (a *Aggregate) CancelOrder(ctx context.Context, args *gateway.CancelOrderRequestArgs) (*gateway.GetOrderResponse, error) {
	ffmQuery := &shipmodelx.GetFulfillmentQuery{
		ShopID:       args.EtopShopID,
		ShippingCode: args.TrackingNumber,
	}
	if err := bus.Dispatch(ctx, ffmQuery); err != nil {
		return nil, err
	}
	ffm := ffmQuery.Result
	orderID := ffm.OrderID

	_, err := logicorder.CancelOrder(ctx, args.EtopShopID, 0, orderID, "Yêu cầu hủy đơn hàng")
	if err != nil {
		return nil, err
	}

	return a.GetOrder(ctx, &gateway.GetOrderRequestArgs{
		EtopShopID:     args.EtopShopID,
		TrackingNumber: args.TrackingNumber,
	})
}

func SHA256StringToInt32(s string) int32 {
	sha256Bytes := sha256.Sum256([]byte(s))
	num := binary.LittleEndian.Uint32(sha256Bytes[:4])

	// convert to int32 and make sure it always positive
	res := int32(num)
	if res < 0 {
		return ^(res - 1)
	}
	return res
}
