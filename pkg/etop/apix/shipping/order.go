package shipping

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ordertypes "etop.vn/api/main/ordering/types"
	shippingcore "etop.vn/api/main/shipping"
	shippingtypes "etop.vn/api/main/shipping/types"
	exttypes "etop.vn/api/top/external/types"
	apishop "etop.vn/api/top/int/shop"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/inventory_auto"
	pbsource "etop.vn/api/top/types/etc/order_source"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	"etop.vn/backend/com/main/ordering/modelx"
	ordersqlstore "etop.vn/backend/com/main/ordering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	convertpbint "etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

func CreateOrder(ctx context.Context, shopClaim *claims.ShopClaim, r *exttypes.CreateOrderRequest) (_ *exttypes.Order, _err error) {
	shipping := r.Shipping
	if shipping == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping")
	}
	lines, err := convertpb.OrderLinesToCreateOrderLines(r.Lines)
	if err != nil {
		return nil, err
	}
	if !shipping.CodAmount.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.cod_amount")
	}
	if !shipping.IncludeInsurance.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.include_insurance")
	}
	if shipping.TryOn == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.try_on")
	}

	var partner *identitymodel.Partner
	if shopClaim.AuthPartnerID != 0 {
		queryPartner := &identitymodelx.GetPartner{
			PartnerID: shopClaim.AuthPartnerID,
		}
		if err := bus.Dispatch(ctx, queryPartner); err != nil {
			return nil, err
		}
		partner = queryPartner.Result.Partner
	}

	externalCode := validate.NormalizeExternalCode(r.ExternalCode)
	if r.ExternalCode != "" && externalCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã đơn hàng external_code không hợp lệ")
	}

	conn, serviceCode, err := parseServiceCode(ctx, shipping.ShippingServiceCode.Apply(""))
	if err != nil {
		return nil, err
	}

	req := &types.CreateOrderRequest{
		Source:          pbsource.API,
		ExternalId:      r.ExternalId,
		ExternalCode:    externalCode,
		ExternalMeta:    r.ExternalMeta,
		ExternalUrl:     r.ExternalUrl,
		PaymentMethod:   0, // will be set automatically
		Customer:        convertpb.OrderAddressToPbCustomer(r.CustomerAddress),
		CustomerAddress: convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		BillingAddress:  convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		ShippingAddress: convertpb.OrderAddressToPbOrder(r.ShippingAddress),
		ShopAddress:     convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
		Lines:           lines,
		Discounts:       nil,
		TotalItems:      r.TotalItems,
		BasketValue:     r.BasketValue,
		TotalWeight:     shipping.ChargeableWeight.Apply(0),
		OrderDiscount:   r.OrderDiscount,
		TotalFee:        r.TotalFee,
		FeeLines:        r.FeeLines,
		TotalDiscount:   dot.Int(r.TotalDiscount),
		TotalAmount:     r.TotalAmount,
		OrderNote:       r.OrderNote,
		ShippingNote:    shipping.ShippingNote.Apply(""),
		ShopShippingFee: 0, // deprecated
		ShopCod:         shipping.CodAmount.Apply(0),
		ReferenceUrl:    "",
		ShopShipping:    nil, // deprecated
		Shipping: &types.OrderShipping{
			ExportedFields:      nil,
			ShAddress:           nil,
			XServiceId:          "",
			XShippingFee:        0,
			XServiceName:        "",
			PickupAddress:       convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
			ReturnAddress:       convertpb.OrderAddressToPbOrder(shipping.ReturnAddress),
			ShippingServiceName: "", // TODO: be filled when confirm
			ShippingServiceCode: serviceCode,
			ShippingServiceFee:  shipping.ShippingServiceFee.Apply(0),
			ShippingProvider:    0,
			Carrier:             shipping.Carrier,
			IncludeInsurance:    shipping.IncludeInsurance.Apply(false),
			TryOn:               shipping.TryOn,
			ShippingNote:        shipping.ShippingNote.Apply(""),
			CodAmount:           shipping.CodAmount,
			GrossWeight:         shipping.GrossWeight,
			Length:              shipping.Length,
			Width:               shipping.Width,
			Height:              shipping.Height,
			ChargeableWeight:    shipping.ChargeableWeight,
		},
		GhnNoteCode: 0, // will be over-written by try_on
	}

	if err := validateAddress(req.CustomerAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	if err := validateAddress(req.ShippingAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ người nhận không hợp lệ: %v", err)
	}
	if err := validateAddress(req.Shipping.PickupAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}
	resp, err := logicorder.CreateOrder(ctx, shopClaim, partner, req, nil, 0)
	if err != nil {
		return nil, err
	}
	orderQuery := &modelx.GetOrderQuery{
		OrderID:            resp.Id,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, orderQuery); err != nil {
		return nil, cm.MapError(err).
			Map(cm.NotFound, cm.Internal, "").
			Throw()
	}

	return convertpb.PbOrder(orderQuery.Result.Order), nil
}

func ConfirmOrder(ctx context.Context, accountID dot.ID, shopClaim *claims.ShopClaim, orderID dot.ID) (_ *exttypes.OrderAndFulfillments, _err error) {
	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			_, err := logicorder.CancelOrder(ctx, shopClaim.Shop.ID, shopClaim.AuthPartnerID, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", _err), inventory_auto.Unknown)
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()

	resp, err := logicorder.ConfirmOrder(ctx, shopClaim.Shop, &apishop.ConfirmOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, err
	}

	conn, _, err := parseServiceCode(ctx, resp.Shipping.ShippingServiceCode)
	if err != nil {
		return nil, err
	}

	createFfmArgs := &shippingcore.CreateFulfillmentsCommand{
		ShopID:              accountID,
		OrderID:             orderID,
		PickupAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(resp.Shipping.PickupAddress),
		ShippingAddress:     convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(resp.ShippingAddress),
		ReturnAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(resp.Shipping.ReturnAddress),
		ShippingType:        ordertypes.ShippingTypeShipment,
		ShippingServiceCode: resp.Shipping.ShippingServiceCode,
		ShippingServiceFee:  resp.Shipping.ShippingServiceFee,
		ShippingServiceName: resp.Shipping.ShippingServiceName,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      resp.Shipping.GrossWeight.Apply(0),
			ChargeableWeight: resp.Shipping.ChargeableWeight.Apply(0),
			Length:           resp.Shipping.Length.Apply(0),
			Width:            resp.Shipping.Width.Apply(0),
			Height:           resp.Shipping.Height.Apply(0),
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      resp.BasketValue,
			CODAmount:        resp.Shipping.CodAmount.Apply(0),
			IncludeInsurance: resp.Shipping.IncludeInsurance,
		},
		TryOn:        resp.Shipping.TryOn,
		ShippingNote: resp.Shipping.ShippingNote,
		ConnectionID: conn.ID,
	}
	if err := shippingAggr.Dispatch(ctx, createFfmArgs); err != nil {
		return nil, err
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
	return convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments), nil
}

func CreateAndConfirmOrder(ctx context.Context, accountID dot.ID, shopClaim *claims.ShopClaim, r *exttypes.CreateOrderRequest) (_ *exttypes.OrderAndFulfillments, _err error) {
	shipping := r.Shipping
	if shipping == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping")
	}
	lines, err := convertpb.OrderLinesToCreateOrderLines(r.Lines)
	if err != nil {
		return nil, err
	}
	if !shipping.CodAmount.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.cod_amount")
	}
	if !shipping.IncludeInsurance.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.include_insurance")
	}
	if shipping.TryOn == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.try_on")
	}

	var partner *identitymodel.Partner
	if shopClaim.AuthPartnerID != 0 {
		queryPartner := &identitymodelx.GetPartner{
			PartnerID: shopClaim.AuthPartnerID,
		}
		if err := bus.Dispatch(ctx, queryPartner); err != nil {
			return nil, err
		}
		partner = queryPartner.Result.Partner
	}

	externalCode := validate.NormalizeExternalCode(r.ExternalCode)
	if r.ExternalCode != "" && externalCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã đơn hàng external_code không hợp lệ")
	}

	conn, serviceCode, err := parseServiceCode(ctx, shipping.ShippingServiceCode.Apply(""))
	if err != nil {
		return nil, err
	}

	req := &types.CreateOrderRequest{
		Source:          pbsource.API,
		ExternalId:      r.ExternalId,
		ExternalCode:    externalCode,
		ExternalMeta:    r.ExternalMeta,
		ExternalUrl:     r.ExternalUrl,
		PaymentMethod:   0, // will be set automatically
		Customer:        convertpb.OrderAddressToPbCustomer(r.CustomerAddress),
		CustomerAddress: convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		BillingAddress:  convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		ShippingAddress: convertpb.OrderAddressToPbOrder(r.ShippingAddress),
		ShopAddress:     convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
		Lines:           lines,
		Discounts:       nil,
		TotalItems:      r.TotalItems,
		BasketValue:     r.BasketValue,
		TotalWeight:     shipping.ChargeableWeight.Apply(0),
		OrderDiscount:   r.OrderDiscount,
		TotalFee:        r.TotalFee,
		FeeLines:        r.FeeLines,
		TotalDiscount:   dot.Int(r.TotalDiscount),
		TotalAmount:     r.TotalAmount,
		OrderNote:       r.OrderNote,
		ShippingNote:    shipping.ShippingNote.Apply(""),
		ShopShippingFee: 0, // deprecated
		ShopCod:         shipping.CodAmount.Apply(0),
		ReferenceUrl:    "",
		ShopShipping:    nil, // deprecated
		Shipping: &types.OrderShipping{
			ExportedFields:      nil,
			ShAddress:           nil,
			XServiceId:          "",
			XShippingFee:        0,
			XServiceName:        "",
			PickupAddress:       convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
			ReturnAddress:       convertpb.OrderAddressToPbOrder(shipping.ReturnAddress),
			ShippingServiceName: "", // TODO: be filled when confirm
			ShippingServiceCode: serviceCode,
			ShippingServiceFee:  shipping.ShippingServiceFee.Apply(0),
			ShippingProvider:    0,
			Carrier:             shipping.Carrier,
			IncludeInsurance:    shipping.IncludeInsurance.Apply(false),
			TryOn:               shipping.TryOn,
			ShippingNote:        shipping.ShippingNote.Apply(""),
			CodAmount:           shipping.CodAmount,
			GrossWeight:         shipping.GrossWeight,
			Length:              shipping.Length,
			Width:               shipping.Width,
			Height:              shipping.Height,
			ChargeableWeight:    shipping.ChargeableWeight,
		},
		GhnNoteCode: 0, // will be over-written by try_on
	}

	if err := validateAddress(req.CustomerAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	if err := validateAddress(req.ShippingAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ người nhận không hợp lệ: %v", err)
	}
	if err := validateAddress(req.Shipping.PickupAddress, conn.ConnectionProvider); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}
	resp, err := logicorder.CreateOrder(ctx, shopClaim, partner, req, nil, 0)
	if err != nil {
		return nil, err
	}

	orderID := resp.Id
	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			_, err := logicorder.CancelOrder(ctx, shopClaim.Shop.ID, shopClaim.AuthPartnerID, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", err), inventory_auto.Unknown)
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()

	_, err = logicorder.ConfirmOrder(ctx, shopClaim.Shop, &apishop.ConfirmOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, err
	}
	createFfmArgs := &shippingcore.CreateFulfillmentsCommand{
		ShopID:              accountID,
		OrderID:             orderID,
		PickupAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(req.Shipping.PickupAddress),
		ShippingAddress:     convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(req.ShippingAddress),
		ReturnAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(req.Shipping.ReturnAddress),
		ShippingType:        ordertypes.ShippingTypeShipment,
		ShippingServiceCode: serviceCode,
		ShippingServiceFee:  shipping.ShippingServiceFee.Apply(0),
		ShippingServiceName: shipping.ShippingServiceName.Apply(""),
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      shipping.GrossWeight.Apply(0),
			ChargeableWeight: shipping.ChargeableWeight.Apply(0),
			Length:           shipping.Length.Apply(0),
			Width:            shipping.Width.Apply(0),
			Height:           shipping.Height.Apply(0),
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      r.BasketValue,
			CODAmount:        shipping.CodAmount.Apply(0),
			IncludeInsurance: shipping.IncludeInsurance.Apply(false),
		},
		TryOn:        shipping.TryOn,
		ShippingNote: shipping.ShippingNote.Apply(""),
		ConnectionID: conn.ID,
	}
	if err := shippingAggr.Dispatch(ctx, createFfmArgs); err != nil {
		return nil, err
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
	return convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments), nil
}

func CancelOrder(ctx context.Context, shopID dot.ID, r *exttypes.CancelOrderRequest) (*exttypes.OrderAndFulfillments, error) {
	var orderID dot.ID
	var sqlQuery *ordersqlstore.OrderStore

	count := 0
	if r.Id != 0 {
		count++
		orderID = r.Id
	}
	if r.ExternalId != "" {
		count++
		sqlQuery = orderStore(ctx).ShopID(shopID).ExternalID(r.ExternalId)
	}
	if r.Code != "" {
		count++
		sqlQuery = orderStore(ctx).ShopID(shopID).Code(r.Code)
	}
	if count != 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp id, code hoặc external_code")
	}
	if sqlQuery != nil {
		order, err := sqlQuery.GetOrderDB()
		if err != nil {
			return nil, err
		}
		orderID = order.ID
	}

	resp, err := logicorder.CancelOrder(ctx, shopID, 0, orderID, r.CancelReason, inventory_auto.Unknown)
	if err != nil {
		return nil, err
	}
	orderQuery := &modelx.GetOrderQuery{
		OrderID:            orderID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, orderQuery); err != nil {
		return nil, err
	}
	resp2 := convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments)
	resp2.FulfillmentErrors = resp.Errors
	return resp2, nil
}

func GetOrder(ctx context.Context, shopID dot.ID, r *exttypes.OrderIDRequest) (*exttypes.OrderAndFulfillments, error) {
	orderQuery := &modelx.GetOrderQuery{
		ShopID:             shopID,
		OrderID:            r.Id,
		ExternalID:         r.ExternalId,
		Code:               r.Code,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, orderQuery); err != nil {
		return nil, err
	}
	return convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments), nil
}

func ListFulfillments(ctx context.Context, shopID dot.ID, r *exttypes.ListFulfillmentsRequest) (*exttypes.FulfillmentsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	s := fulfillmentStore(ctx).ShopID(shopID).WithPaging(*paging)
	if len(r.Filter.OrderID) != 0 {
		s = s.OrderIDs(r.Filter.OrderID...)
	}

	ffms, err := s.ListFfmsDB()
	if err != nil {
		return nil, err
	}
	pagingResult := s.GetPaging()
	return &exttypes.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillments(ffms),
		Paging:       convertpb.PbPageInfo(r.Paging, &pagingResult),
	}, nil
}

func GetFulfillment(ctx context.Context, shopID dot.ID, r *exttypes.FulfillmentIDRequest) (*exttypes.Fulfillment, error) {
	s := fulfillmentStore(ctx).ShopID(shopID)
	if r.Id != 0 {
		s = s.ID(r.Id)
	} else if r.ShippingCode != "" {
		s = s.ShippingCode(r.ShippingCode)
	}
	ffm, err := s.GetFfmDB()
	if err != nil {
		return nil, err
	}
	return convertpb.PbFulfillment(ffm), nil
}

func validateAddress(address *types.OrderAddress, shippingProvider connection_type.ConnectionProvider) error {
	if address == nil {
		return errors.New("Thiếu thông tin địa chỉ")
	}
	if address.Phone == "" {
		return errors.New("Thiếu thông tin số điện thoại")
	}
	phoneNorm, ok := validate.NormalizePhone(address.Phone)
	if !ok {
		return errors.New("Số điện thoại không hợp lệ")
	}
	// số điện thoại bàn bắt đầu bằng 02
	if !strings.HasPrefix(string(phoneNorm), "02") && len(phoneNorm) >= 11 {
		return errors.New("Số điện thoại 11 số không còn được hỗ trợ. Vui lòng sử dụng số điện thoại 10 số.")
	}
	if address.Address1 == "" {
		return errors.New("Thiếu thông tin địa chỉ")
	}
	if address.FullName == "" {
		return errors.New("Thiếu thông tin tên")
	}
	if s, ok := validate.NormalizeGenericName(address.FullName); !ok {
		return errors.New("Tên không hợp lệ")
	} else {
		address.FullName = s
	}

	if shippingProvider == connection_type.ConnectionProviderVTP {
		// required Ward
		_address, err := convertpbint.OrderAddressToModel(address)
		if err != nil {
			return err
		}
		if _address.Ward == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin phường/xã (%v, %v) - ViettelPost yêu cầu thông tin phường/xã đúng để giao hàng.", address.District, address.Province)
		}
		if _address.WardCode == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "Phường/xã không hợp lệ (%v, %v, %v) - ViettelPost yêu cầu thông tin phường/xã đúng để giao hàng.", address.Ward, address.District, address.Province)
		}
	}
	return nil
}
