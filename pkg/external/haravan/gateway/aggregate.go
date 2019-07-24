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
	"etop.vn/api/main/location"
	pbsp "etop.vn/backend/pb/etop/etc/shipping_provider"
	pbtryon "etop.vn/backend/pb/etop/etc/try_on"
	pborder "etop.vn/backend/pb/etop/order"
	pbsource "etop.vn/backend/pb/etop/order/source"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/external/haravan/gateway/convert"
	identityconvert "etop.vn/backend/pkg/services/identity/convert"
	"etop.vn/backend/pkg/services/ordering/modelx"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
	shipmodelx "etop.vn/backend/pkg/services/shipping/modelx"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var (
	_  gateway.Aggregate = &Aggregate{}
	ll                   = l.New()
)

type Aggregate struct {
	db           cmsql.Transactioner
	locationQS   location.QueryBus
	ShippingCtrl *shipping_provider.ProviderManager
	identityQS   identity.QueryBus
}

func NewAggregate(db cmsql.Database, providerManager *shipping_provider.ProviderManager, locationQuery location.QueryBus, identityQuery identity.QueryBus) *Aggregate {
	return &Aggregate{
		db:           db,
		locationQS:   locationQuery,
		ShippingCtrl: providerManager,
		identityQS:   identityQuery,
	}
}

func (a *Aggregate) MessageBus() gateway.CommandBus {
	b := bus.New()
	return gateway.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) GetShippingRate(ctx context.Context, args *gateway.GetShippingRateRequestArgs) (*gateway.GetShippingRateResponse, error) {
	from, err := a.GetLocation(ctx, args.Origin)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Lỗi địa chỉ gửi")
	}
	to, err := a.GetLocation(ctx, args.Destination)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Lỗi địa chỉ nhận")
	}

	// Haravan: default includeInsurance is false
	req := &pborder.GetExternalShippingServicesRequest{
		Provider:         pbsp.ShippingProvider_ghn,
		Carrier:          pbsp.ShippingProvider_ghn,
		FromProvince:     from.Province.Name,
		FromDistrict:     from.District.Name,
		ToProvince:       to.Province.Name,
		ToDistrict:       to.District.Name,
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

	from, err := a.GetLocation(ctx, args.Origin)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Lỗi địa chỉ gửi")
	}
	to, err := a.GetLocation(ctx, args.Destination)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Lỗi địa chỉ nhận")
	}

	req := &pborder.GetExternalShippingServicesRequest{
		Provider:         pbsp.ShippingProvider_ghn,
		Carrier:          pbsp.ShippingProvider_ghn,
		FromProvince:     from.Province.Name,
		FromDistrict:     from.District.Name,
		ToProvince:       to.Province.Name,
		ToDistrict:       to.District.Name,
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
	codAmount := int(args.CodAmount)
	weight := int(args.TotalGrams)

	req2 := &pborder.CreateOrderRequest{
		Source:          pbsource.Source_api,
		ExternalId:      externalID,
		ExternalCode:    externalID,
		ExternalUrl:     "",
		PaymentMethod:   "",
		Customer:        convert.ToPbOrderCustomer(args.Origin),
		CustomerAddress: convert.ToPbOrderAddress(args.Origin, from),
		BillingAddress:  convert.ToPbOrderAddress(args.Origin, from),
		ShippingAddress: convert.ToPbOrderAddress(args.Destination, to),
		ShopAddress:     convert.ToPbOrderAddress(args.Origin, from),
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
			PickupAddress:       convert.ToPbOrderAddress(args.Origin, from),
			ReturnAddress:       nil,
			ShippingServiceName: "",
			ShippingServiceCode: service.ProviderServiceID,
			ShippingServiceFee:  int32(service.ServiceFee),
			ShippingProvider:    pbsp.PbShippingProviderType(service.Provider),
			Carrier:             pbsp.PbShippingProviderType(service.Provider),
			IncludeInsurance:    false,
			TryOn:               pbtryon.TryOnCode_none,
			ShippingNote:        args.Note,
			CodAmount:           cm.PIntToInt32(codAmount),
			Weight:              cm.PIntToInt32(weight),
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
		TrackingURL:    generateTrackingUrl(ffm),
		CodAmount:      int32(ffm.TotalCODAmount),
	}, nil
}

func getOrderValue(items []*haravan.Item) (total int) {
	for _, item := range items {
		total += int(item.Price) * item.Quantity
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
		TrackingURL:    generateTrackingUrl(ffm),
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

func (a *Aggregate) GetLocation(ctx context.Context, addr *haravan.Address) (*location.LocationQueryResult, error) {
	if addr == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ không được để trống")
	}
	query := &location.GetLocationQuery{
		ProvinceCode:     addr.ProvinceCode,
		DistrictCode:     addr.DistrictCode,
		WardCode:         addr.WardCode,
		LocationCodeType: location.LocCodeTypeHaravan,
	}
	if err := a.locationQS.Dispatch(ctx, query); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ gửi không hợp lệ: %v", err)
	}
	return query.Result, nil
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

func generateTrackingUrl(ffm *shipmodel.Fulfillment) string {
	// Haravan will concat shipping_code to this link
	baseURL := cm.MainSiteBaseURL()
	return fmt.Sprintf("%v/s/%v/fulfillment?code=", baseURL, ffm.ShopID)
}
