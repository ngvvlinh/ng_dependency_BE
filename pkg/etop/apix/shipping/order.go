package shipping

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pbcm "etop.vn/api/pb/common"
	pborder "etop.vn/api/pb/etop/order"
	pbsource "etop.vn/api/pb/etop/order/source"
	"etop.vn/api/pb/etop/shop"
	pbexternal "etop.vn/api/pb/external"
	"etop.vn/backend/com/main/ordering/modelx"
	ordersqlstore "etop.vn/backend/com/main/ordering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/common/validate"
	convertpbint "etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

func CreateAndConfirmOrder(ctx context.Context, accountID dot.ID, shopClaim *claims.ShopClaim, r *pbexternal.CreateOrderRequest) (_ *pbexternal.OrderAndFulfillments, _err error) {
	shipping := r.Shipping
	if shipping == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping")
	}
	lines, err := convertpb.OrderLinesToCreateOrderLines(r.Lines)
	if err != nil {
		return nil, err
	}
	if shipping.CodAmount == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.cod_amount")
	}
	if shipping.Carrier == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.carrier")
	}
	if shipping.IncludeInsurance == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.include_insurance")
	}
	if shipping.TryOn == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp mục shipping.try_on")
	}

	var partner *model.Partner
	if shopClaim.AuthPartnerID != 0 {
		queryPartner := &model.GetPartner{
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

	req := &pborder.CreateOrderRequest{
		Source:          pbsource.Source_api,
		ExternalId:      r.ExternalId,
		ExternalCode:    externalCode,
		ExternalMeta:    r.ExternalMeta,
		ExternalUrl:     r.ExternalUrl,
		PaymentMethod:   "", // will be set automatically
		Customer:        convertpb.OrderAddressToPbCustomer(r.CustomerAddress),
		CustomerAddress: convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		BillingAddress:  convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		ShippingAddress: convertpb.OrderAddressToPbOrder(r.ShippingAddress),
		ShopAddress:     convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
		ShConfirm:       nil,
		Lines:           lines,
		Discounts:       nil,
		TotalItems:      r.TotalItems,
		BasketValue:     r.BasketValue,
		TotalWeight:     cmapi.BareInt32(shipping.ChargeableWeight),
		OrderDiscount:   r.OrderDiscount,
		TotalFee:        cmapi.BareInt32(r.TotalFee),
		FeeLines:        r.FeeLines,
		TotalDiscount:   &r.TotalDiscount,
		TotalAmount:     r.TotalAmount,
		OrderNote:       r.OrderNote,
		ShippingNote:    cmapi.BareString(shipping.ShippingNote),
		ShopShippingFee: 0, // deprecated
		ShopCod:         cmapi.BareInt32(shipping.CodAmount),
		ReferenceUrl:    "",
		ShopShipping:    nil, // deprecated
		Shipping: &pborder.OrderShipping{
			ExportedFields:      nil,
			ShAddress:           nil,
			XServiceId:          "",
			XShippingFee:        0,
			XServiceName:        "",
			PickupAddress:       convertpb.OrderAddressToPbOrder(shipping.PickupAddress),
			ReturnAddress:       convertpb.OrderAddressToPbOrder(shipping.ReturnAddress),
			ShippingServiceName: "", // TODO: be filled when confirm
			ShippingServiceCode: cmapi.BareString(shipping.ShippingServiceCode),
			ShippingServiceFee:  cmapi.BareInt32(shipping.ShippingServiceFee),
			ShippingProvider:    0,
			Carrier:             *shipping.Carrier,
			IncludeInsurance:    *shipping.IncludeInsurance,
			TryOn:               *shipping.TryOn,
			ShippingNote:        cmapi.BareString(shipping.ShippingNote),
			CodAmount:           shipping.CodAmount,
			Weight:              nil,
			GrossWeight:         shipping.GrossWeight,
			Length:              shipping.Length,
			Width:               shipping.Width,
			Height:              shipping.Height,
			ChargeableWeight:    shipping.ChargeableWeight,
		},
		GhnNoteCode: 0, // will be over-written by try_on
	}
	if err := validateAddress(req.CustomerAddress, convertpbint.ShippingProviderToModel(shipping.Carrier)); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	if err := validateAddress(req.ShippingAddress, convertpbint.ShippingProviderToModel(shipping.Carrier)); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ người nhận không hợp lệ: %v", err)
	}
	if err := validateAddress(req.Shipping.PickupAddress, convertpbint.ShippingProviderToModel(shipping.Carrier)); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}

	resp, err := logicorder.CreateOrder(ctx, shopClaim, partner, req, nil)
	if err != nil {
		return nil, err
	}

	orderID := resp.Id
	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			_, err := logicorder.CancelOrder(ctx, shopClaim.Shop.ID, shopClaim.AuthPartnerID, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", err), "")
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()

	cfmResp, err := logicorder.ConfirmOrderAndCreateFulfillments(ctx, shopClaim.Shop, shopClaim.AuthPartnerID,
		&shop.OrderIDRequest{
			OrderId: orderID,
		})
	if err != nil {
		return nil, err
	}
	if len(cfmResp.FulfillmentErrors) > 0 {
		// TODO: refactor
		var anErr *pbcm.Error
		for _, err := range cfmResp.FulfillmentErrors {
			if err.Code != "ok" {
				anErr = err
				break
			}
		}
		return nil, anErr
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

func CancelOrder(ctx context.Context, shopID dot.ID, r *pbexternal.CancelOrderRequest) (*pbexternal.OrderAndFulfillments, error) {
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

	resp, err := logicorder.CancelOrder(ctx, shopID, 0, orderID, r.CancelReason, "")
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

func GetOrder(ctx context.Context, shopID dot.ID, r *pbexternal.OrderIDRequest) (*pbexternal.OrderAndFulfillments, error) {
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

func GetFulfillment(ctx context.Context, shopID dot.ID, r *pbexternal.FulfillmentIDRequest) (*pbexternal.Fulfillment, error) {
	s := fulfillmentStore(ctx).ShopID(shopID)
	if r.Id != 0 {
		s = s.ID(r.Id)
	} else if r.ShippingCode != "" {
		s = s.ShippingCode(r.ShippingCode)
	}
	ffm, err := s.Get()
	if err != nil {
		return nil, err
	}
	return convertpb.PbFulfillment(ffm), nil
}

func validateAddress(address *pborder.OrderAddress, shippingProvider model.ShippingProvider) error {
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

	if shippingProvider == model.TypeVTPost {
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
