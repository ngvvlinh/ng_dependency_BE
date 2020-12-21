package shipping

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ordertypes "o.o/api/main/ordering/types"
	shippingcore "o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	exttypes "o.o/api/top/external/types"
	apishop "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	pbsource "o.o/api/top/types/etc/order_source"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/com/main/ordering/modelx"
	ordersqlstore "o.o/backend/com/main/ordering/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/validate"
	convertpbint "o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

func (s *Shipping) CreateOrder(ctx context.Context, shop *identitymodel.Shop, partner *identitymodel.Partner, r *exttypes.CreateOrderRequest) (_ *exttypes.OrderWithoutShipping, _err error) {
	lines, err := convertpb.OrderLinesToCreateOrderLines(r.Lines)
	if err != nil {
		return nil, err
	}

	externalCode := validate.NormalizeExternalCode(r.ExternalCode)
	if r.ExternalCode != "" && externalCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã đơn hàng external_code không hợp lệ")
	}

	req := &types.CreateOrderRequest{
		Source:          pbsource.API,
		ExternalId:      r.ExternalId,
		ExternalCode:    externalCode,
		ExternalMeta:    r.ExternalMeta,
		ExternalUrl:     r.ExternalUrl,
		PaymentMethod:   r.PaymentMethod,
		Customer:        convertpb.OrderAddressToPbCustomer(r.CustomerAddress),
		CustomerAddress: convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		BillingAddress:  convertpb.OrderAddressToPbOrder(r.CustomerAddress),
		ShippingAddress: convertpb.OrderAddressToPbOrder(r.ShippingAddress),
		Lines:           lines,
		Discounts:       nil,
		TotalItems:      r.TotalItems,
		BasketValue:     r.BasketValue,
		OrderDiscount:   r.OrderDiscount,
		TotalFee:        r.TotalFee,
		FeeLines:        r.FeeLines,
		TotalDiscount:   dot.Int(r.TotalDiscount),
		TotalAmount:     r.TotalAmount,
		OrderNote:       r.OrderNote,
		ReferenceUrl:    "",
		ShopShipping:    nil,
		GhnNoteCode:     0, // will be over-written by try_on
	}

	resp, err := s.OrderLogic.CreateOrder(ctx, shop, partner, req, nil, shop.OwnerID)
	if err != nil {
		return nil, err
	}
	orderQuery := &modelx.GetOrderQuery{
		OrderID:            resp.Id,
		IncludeFulfillment: true,
	}
	if err = s.OrderStoreIface.GetOrder(ctx, orderQuery); err != nil {
		return nil, cm.MapError(err).
			Map(cm.NotFound, cm.Internal, "").
			Throw()
	}
	return convertpb.PbOrderWithoutShipping(orderQuery.Result.Order), nil
}

func (s *Shipping) ConfirmOrder(ctx context.Context, userID dot.ID, shop *identitymodel.Shop, partner *identitymodel.Partner, orderID dot.ID, autoInventoryVoucher inventory_auto.AutoInventoryVoucher) (_err error) {
	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			_, err := s.OrderLogic.CancelOrder(ctx, userID, shop.ID, partner.ID, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", _err), inventory_auto.Unknown)
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()
	_, err := s.OrderLogic.ConfirmOrder(ctx, userID, shop, &apishop.ConfirmOrderRequest{
		OrderId:              orderID,
		AutoInventoryVoucher: autoInventoryVoucher,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Shipping) CreateAndConfirmOrder(ctx context.Context, userID dot.ID, shop *identitymodel.Shop, partner *identitymodel.Partner, r *exttypes.CreateAndConfirmOrderRequest) (_ *exttypes.OrderAndFulfillments, _err error) {
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

	externalCode := validate.NormalizeExternalCode(r.ExternalCode)
	if r.ExternalCode != "" && externalCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã đơn hàng external_code không hợp lệ")
	}

	conn, serviceCode, err := s.parseServiceCode(ctx, shipping.ShippingServiceCode.Apply(""))
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
	}

	if err := validateAddress(req.CustomerAddress); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	if err := validateAddress(req.ShippingAddress); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ người nhận không hợp lệ: %v", err)
	}
	if err := validateAddress(req.Shipping.PickupAddress); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}
	resp, err := s.OrderLogic.CreateOrder(ctx, shop, partner, req, nil, userID)
	if err != nil {
		return nil, err
	}

	orderID := resp.Id
	defer func() {
		if _err != nil {
			// always cancel order if confirm unsuccessfully
			partnerID := dot.ID(0)
			if partner != nil {
				partnerID = partner.ID
			}
			_, err := s.OrderLogic.CancelOrder(ctx, userID, shop.ID, partnerID, orderID, fmt.Sprintf("Tạo đơn không thành công: %v", err), inventory_auto.Unknown)
			if err != nil {
				ll.Error("error cancelling order", l.Error(err))
			}
		}
	}()

	_, err = s.OrderLogic.ConfirmOrder(ctx, userID, shop, &apishop.ConfirmOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, err
	}
	createFfmArgs := &shippingcore.CreateFulfillmentsCommand{
		ShopID:              shop.ID,
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
	if err := s.ShippingAggr.Dispatch(ctx, createFfmArgs); err != nil {
		return nil, err
	}

	orderQuery := &modelx.GetOrderQuery{
		OrderID:            orderID,
		IncludeFulfillment: true,
	}
	if err := s.OrderStoreIface.GetOrder(ctx, orderQuery); err != nil {
		return nil, cm.MapError(err).
			Map(cm.NotFound, cm.Internal, "").
			Throw()
	}
	return convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments), nil
}

func (s *Shipping) CancelOrder(ctx context.Context, userID dot.ID, shopID dot.ID, r *exttypes.CancelOrderRequest) (*exttypes.OrderAndFulfillments, error) {
	var orderID dot.ID
	var sqlQuery *ordersqlstore.OrderStore

	count := 0
	if r.Id != 0 {
		count++
		orderID = r.Id
	}
	if r.ExternalId != "" {
		count++
		sqlQuery = s.OrderStore(ctx).ShopID(shopID).ExternalID(r.ExternalId)
	}
	if r.Code != "" {
		count++
		sqlQuery = s.OrderStore(ctx).ShopID(shopID).Code(r.Code)
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

	resp, err := s.OrderLogic.CancelOrder(ctx, userID, shopID, 0, orderID, r.CancelReason, inventory_auto.Unknown)
	if err != nil {
		return nil, err
	}
	orderQuery := &modelx.GetOrderQuery{
		OrderID:            orderID,
		IncludeFulfillment: true,
	}
	if err := s.OrderStoreIface.GetOrder(ctx, orderQuery); err != nil {
		return nil, err
	}
	resp2 := convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments)
	resp2.FulfillmentErrors = resp.Errors
	return resp2, nil
}

func (s *Shipping) GetOrder(ctx context.Context, shopID dot.ID, r *exttypes.OrderIDRequest) (*exttypes.OrderAndFulfillments, error) {
	orderQuery := &modelx.GetOrderQuery{
		ShopID:             shopID,
		OrderID:            r.Id,
		ExternalID:         r.ExternalId,
		Code:               r.Code,
		IncludeFulfillment: true,
	}
	if err := s.OrderStoreIface.GetOrder(ctx, orderQuery); err != nil {
		return nil, err
	}
	return convertpb.PbOrderAndFulfillments(orderQuery.Result.Order, orderQuery.Result.Fulfillments), nil
}

func (s *Shipping) ListFulfillments(ctx context.Context, shopID dot.ID, r *exttypes.ListFulfillmentsRequest) (*exttypes.FulfillmentsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	q := s.FulfillmentStore(ctx).ShopID(shopID).WithPaging(*paging)
	if len(r.Filter.OrderID) != 0 {
		q = q.OrderIDs(r.Filter.OrderID...)
	}

	ffms, err := q.ListFfmsDB()
	if err != nil {
		return nil, err
	}
	pagingResult := q.GetPaging()
	return &exttypes.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillments(ffms),
		Paging:       convertpb.PbPageInfo(paging, &pagingResult),
	}, nil
}

func (s *Shipping) CreateFulfillment(ctx context.Context, shopID dot.ID, r *exttypes.CreateFulfillmentRequest) (*exttypes.Fulfillment, error) {
	conn, serviceCode, err := s.parseServiceCode(ctx, r.ShippingServiceCode)
	if err != nil {
		return nil, err
	}

	createFfmArgs := &shippingcore.CreateFulfillmentsCommand{
		ShopID:              shopID,
		OrderID:             r.OrderID,
		PickupAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(r.PickupAddress),
		ShippingAddress:     convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(r.ShippingAddress),
		ReturnAddress:       convertpbint.Convert_api_OrderAddress_To_core_OrderAddress(r.ReturnAddress),
		ShippingType:        ordertypes.ShippingTypeShipment,
		ShippingServiceCode: serviceCode,
		ShippingServiceFee:  r.ShippingServiceFee,
		ShippingServiceName: r.ShippingServiceName,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      r.GrossWeight,
			ChargeableWeight: r.ChargeableWeight,
			Length:           r.Length,
			Width:            r.Width,
			Height:           r.Height,
		},
		ValueInfo: shippingtypes.ValueInfo{
			CODAmount:        r.CODAmount,
			IncludeInsurance: r.IncludeInsurance,
		},
		TryOn:         r.TryOn,
		ShippingNote:  r.ShippingNote,
		ConnectionID:  conn.ID,
		ShopCarrierID: r.ShopCarrierID,
	}
	if err := s.ShippingAggr.Dispatch(ctx, createFfmArgs); err != nil {
		return nil, err
	}

	query := &shippingcore.GetFulfillmentByIDOrShippingCodeQuery{
		ID: createFfmArgs.Result[0],
	}
	if err := s.ShippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	ffm, err := s.FulfillmentStore(ctx).ShopID(shopID).ID(createFfmArgs.Result[0]).GetFfmDB()
	if err != nil {
		return nil, err
	}

	return convertpb.PbFulfillment(ffm), nil
}

func (s *Shipping) GetFulfillment(ctx context.Context, shopID dot.ID, r *exttypes.FulfillmentIDRequest) (*exttypes.Fulfillment, error) {
	q := s.FulfillmentStore(ctx).ShopID(shopID)
	if r.Id != 0 {
		q = q.ID(r.Id)
	} else if r.ShippingCode != "" {
		q = q.ShippingCode(r.ShippingCode)
	}
	ffm, err := q.GetFfmDB()
	if err != nil {
		return nil, err
	}
	return convertpb.PbFulfillment(ffm), nil
}

func validateAddress(address *types.OrderAddress) error {
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

	// required Ward
	_address, err := convertpbint.OrderAddressToModel(address)
	if err != nil {
		return err
	}
	if _address.Ward == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin phường/xã (%v, %v)", address.District, address.Province)
	}
	if _address.WardCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Phường/xã không hợp lệ (%v, %v, %v) ", address.Ward, address.District, address.Province)
	}
	return nil
}

func (s *Shipping) CancelFulfillment(ctx context.Context, fulfillmentID dot.ID, cancelReason string) error {
	cmd := &shippingcore.CancelFulfillmentCommand{
		FulfillmentID: fulfillmentID,
		CancelReason:  cancelReason,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
