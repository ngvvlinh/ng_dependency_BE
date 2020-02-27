package orderS

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/backend/com/main/catalog/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

func CreateOrder(
	ctx context.Context, claim *claims.ShopClaim,
	authPartner *identitymodel.Partner, r *types.CreateOrderRequest,
	tradingShopID *dot.ID, userID dot.ID) (*types.Order, error) {
	shipping := r.ShopShipping
	if r.Shipping != nil {
		shipping = r.Shipping
	}
	if shipping == nil {
		shipping = &types.OrderShipping{}
	}

	if shipping.Carrier == 0 {
		shipping.Carrier = shipping.ShippingProvider
	}
	if (r.Shipping != nil || r.ShopShipping != nil) &&
		!model.VerifyShippingProvider(shipping.Carrier) {
		return nil, cm.Error(cm.InvalidArgument, "Nhà vận chuyển không hợp lệ", nil)
	}
	if r.ExternalUrl != "" {
		recognizedHosts := claim.Shop.RecognizedHosts
		if authPartner != nil {
			recognizedHosts = authPartner.RecognizedHosts
		}
		if err := validateExternalURL(recognizedHosts, r.ExternalUrl); err != nil {
			return nil, err
		}
	}

	if !model.VerifyOrderSource(r.Source) {
		return nil, cm.Error(cm.InvalidArgument, "Invalid source", nil)
	}

	if r.CustomerId != 0 && r.Customer != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "customer_id và customer không được gửi cùng 1 lúc")
	}

	if r.CustomerId != 0 && r.CustomerAddress != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "customer_id và customerAddress không được gửi cùng 1 lúc")
	}

	if r.CustomerId != 0 && r.Customer == nil {
		getCustomerQuery := &customering.GetCustomerByIDQuery{
			ID:     r.CustomerId,
			ShopID: claim.Shop.ID,
		}
		if err := customerQuery.Dispatch(ctx, getCustomerQuery); err != nil {
			return nil, err
		}
		r.Customer = &types.OrderCustomer{
			FullName: getCustomerQuery.Result.FullName,
			Email:    getCustomerQuery.Result.Email,
			Phone:    getCustomerQuery.Result.Phone,
		}
	}
	if r.CustomerId != 0 && r.CustomerAddress == nil {
		isHaveCustomerAddress := true
		getAddressQuery := &addressing.GetAddressActiveByTraderIDQuery{
			TraderID: r.CustomerId,
			ShopID:   claim.Shop.ID,
		}
		if err := traderAddressQuery.Dispatch(ctx, getAddressQuery); err != nil {
			switch cm.ErrorCode(err) {
			case cm.NotFound:
				isHaveCustomerAddress = false
			default:
				return nil, err
			}
		}
		if isHaveCustomerAddress {
			customerAddress, err := convertpb.PbShopAddress(ctx, getAddressQuery.Result, locationQuery)
			if err != nil {
				return nil, err
			}
			r.CustomerAddress = &types.OrderAddress{
				FullName:     customerAddress.FullName,
				Phone:        customerAddress.Phone,
				Email:        customerAddress.Email,
				Country:      customerAddress.Country,
				Province:     customerAddress.Province,
				District:     customerAddress.District,
				Ward:         customerAddress.Ward,
				Company:      customerAddress.Company,
				Address1:     customerAddress.Address1,
				Address2:     customerAddress.Address2,
				ProvinceCode: customerAddress.ProvinceCode,
				DistrictCode: customerAddress.DistrictCode,
				WardCode:     customerAddress.WardCode,
				Coordinates:  customerAddress.Coordinates,
			}
		}
	}

	shippingAddress := r.ShippingAddress
	if r.CustomerId == 0 && shippingAddress != nil &&
		shippingAddress.Phone != "" && shippingAddress.FullName != "" {
		getCustomerByPhone := &customering.GetCustomerByPhoneQuery{
			Phone:  shippingAddress.Phone,
			ShopID: claim.Shop.ID,
		}
		if err := customerQuery.Dispatch(ctx, getCustomerByPhone); err != nil && cm.ErrorCode(err) != cm.NotFound {
			return nil, err
		}
		phoneCustomer := getCustomerByPhone.Result

		var emailCustomer *customering.ShopCustomer
		if shippingAddress.Email != "" {
			getCustomerByEmail := &customering.GetCustomerByEmailQuery{
				Email:  shippingAddress.Email,
				ShopID: claim.Shop.ID,
			}
			if err := customerQuery.Dispatch(ctx, getCustomerByEmail); err != nil && cm.ErrorCode(err) != cm.NotFound {
				return nil, err
			}
			emailCustomer = getCustomerByEmail.Result
		}

		if phoneCustomer == nil {
			if emailCustomer == nil {
				createCustomerCmd := &customering.CreateCustomerCommand{
					ShopID:   claim.Shop.ID,
					FullName: shippingAddress.FullName,
					Type:     customer_type.Individual,
					Phone:    shippingAddress.Phone,
					Email:    shippingAddress.Email,
				}
				if err := customerAggr.Dispatch(ctx, createCustomerCmd); err != nil {
					return nil, err
				}
				r.CustomerId = createCustomerCmd.Result.ID
			} else {
				createCustomerCmd := &customering.CreateCustomerCommand{
					ShopID:   claim.Shop.ID,
					FullName: shippingAddress.FullName,
					Type:     customer_type.Individual,
					Phone:    shippingAddress.Phone,
				}
				if err := customerAggr.Dispatch(ctx, createCustomerCmd); err != nil {
					return nil, err
				}
				r.CustomerId = createCustomerCmd.Result.ID
			}
		} else {
			if emailCustomer == nil {
				cmd := &customering.UpdateCustomerCommand{
					ID:       phoneCustomer.ID,
					ShopID:   claim.Shop.ID,
					FullName: dot.String(shippingAddress.FullName),
					Email:    dot.String(shippingAddress.Email),
				}
				if err := customerAggr.Dispatch(ctx, cmd); err != nil {
					return nil, err
				}
			} else {
				cmd := &customering.UpdateCustomerCommand{
					ID:       phoneCustomer.ID,
					ShopID:   claim.Shop.ID,
					FullName: dot.String(shippingAddress.FullName),
				}
				if err := customerAggr.Dispatch(ctx, cmd); err != nil {
					return nil, err
				}
			}
			r.CustomerId = phoneCustomer.ID
		}
		// ignore err
		if _err := updateOrCreateCustomerAddress(ctx, claim.Shop.ID, r.CustomerId, shippingAddress); _err != nil {
			ll.Error("Auto cập nhật Customer Address lỗi", l.Error(_err))
		}
		r.Customer = getCustomerByID(ctx, claim.Shop.ID, r.CustomerId)
	}
	if r.CustomerId == 0 && r.ShippingAddress == nil {
		cmd := &customering.GetCustomerIndependentQuery{}
		if err := customerQuery.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
		r.CustomerId = cmd.Result.ID
		r.Customer = &types.OrderCustomer{
			FullName: cmd.Result.FullName,
			Type:     cmd.Result.Type,
		}
	}
	shop := claim.Shop
	lines, err := PrepareOrderLines(ctx, shop.ID, r.Lines)
	if err != nil {
		return nil, err
	}
	order, err := PrepareOrder(ctx, shop.ID, r, lines, userID)
	if err != nil {
		return nil, err
	}
	if authPartner != nil {
		order.PartnerID = authPartner.ID
	}

	order.ShopID = shop.ID
	order.OrderSourceType = r.Source
	// fulfillment_type will be filled after create fulfillment
	order.FulfillmentType = ordertypes.ShippingTypeNone
	if tradingShopID != nil {
		order.TradingShopID = *tradingShopID
	}

	cmd := &ordermodelx.CreateOrderCommand{
		Order: order,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		// TODO: refactor
		if xerr, ok := err.(*xerrors.APIError); ok && xerr.Err != nil {
			msg := xerr.Err.Error()
			switch {
			case strings.Contains(msg, "order_shop_external_id_idx"):
				newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn hàng external_id đã tồn tại. Vui lòng kiểm tra lại.").
					WithMeta("duplicated", "external_id")
				orderQuery := &ordermodelx.GetOrderQuery{
					ShopID:     shop.ID,
					ExternalID: r.ExternalId, // TODO: external id may be normalized, this won't work
				}
				_ = bus.Dispatch(ctx, orderQuery)
				if orderQuery.Result.Order != nil {
					newErr = newErr.WithMetap("order_id", orderQuery.Result.Order.ID)
				}
				return nil, newErr

			case strings.Contains(msg, "order_partner_external_id_idx"):
				newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn hàng external_id đã tồn tại. Vui lòng kiểm tra lại.").
					WithMeta("duplicated", "external_id")
				orderQuery := &ordermodelx.GetOrderQuery{
					PartnerID:  shop.ID,
					ExternalID: r.ExternalId,
				}
				_ = bus.Dispatch(ctx, orderQuery)
				if orderQuery.Result.Order != nil {
					newErr = newErr.WithMetap("order_id", orderQuery.Result.Order.ID)
				}
				return nil, newErr

			case strings.Contains(msg, "order_shop_id_ed_code_idx"):
				newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn hàng external_code đã tồn tại. Vui lòng kiểm tra lại.").
					WithMeta("duplicated", "external_code")
				// TODO: include order_id
				return nil, newErr

			case strings.Contains(msg, "order_partner_shop_id_external_code_idx"):
				newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn hàng external_code đã tồn tại. Vui lòng kiểm tra lại.").
					WithMeta("duplicated", "external_code")
				// TODO: include order_id
				return nil, newErr
			}
		}
		return nil, err
	}
	result := convertpb.PbOrder(order, nil, model.TagShop)
	result.ShopName = claim.Shop.Name
	return result, nil
}

func getCustomerByID(ctx context.Context, shopID, customerID dot.ID) *types.OrderCustomer {
	getCustomer := &customering.GetCustomerByIDQuery{
		ID:     customerID,
		ShopID: shopID,
	}
	err := customerQuery.Dispatch(ctx, getCustomer)
	if err != nil {
		return nil
	}
	customer := &types.OrderCustomer{
		FullName: getCustomer.Result.FullName,
		Email:    getCustomer.Result.Email,
		Phone:    getCustomer.Result.Phone,
		Type:     getCustomer.Result.Type,
		Gender:   getCustomer.Result.Gender,
	}
	return customer
}

func updateOrCreateCustomerAddress(ctx context.Context, shopID, customerID dot.ID, orderAddress *types.OrderAddress) error {
	address, err := convertpb.OrderAddressToModel(orderAddress)
	if err != nil {
		return err
	}

	getAddressQuery := &addressing.GetAddressActiveByTraderIDQuery{
		TraderID: customerID,
		ShopID:   shopID,
	}
	err = traderAddressQuery.Dispatch(ctx, getAddressQuery)
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return err
	}
	addressDB := getAddressQuery.Result
	if err == nil && addressDB != nil {
		updateCustomerAddressCmd := &addressing.UpdateAddressCommand{
			ID:           addressDB.ID,
			ShopID:       shopID,
			FullName:     dot.String(address.FullName),
			Phone:        dot.String(address.Phone),
			Email:        dot.String(address.Email),
			Company:      dot.String(address.Company),
			Address1:     dot.String(address.Address1),
			Address2:     dot.String(address.Address2),
			DistrictCode: dot.String(address.DistrictCode),
			WardCode:     dot.String(address.WardCode),
		}
		if err := traderAddressAggr.Dispatch(ctx, updateCustomerAddressCmd); err != nil {
			return err
		}
	} else {
		return createCustomerAddress(ctx, shopID, customerID, address)
	}
	return nil
}

func createCustomerAddress(
	ctx context.Context, shopID, traderID dot.ID, orderAddress *ordermodel.OrderAddress) error {
	createAddressCmd := &addressing.CreateAddressCommand{
		ShopID:       shopID,
		TraderID:     traderID,
		FullName:     orderAddress.FullName,
		Phone:        orderAddress.Phone,
		Email:        orderAddress.Email,
		Company:      orderAddress.Company,
		Address1:     orderAddress.Address1,
		Address2:     orderAddress.Address2,
		DistrictCode: orderAddress.DistrictCode,
		WardCode:     orderAddress.WardCode,
		IsDefault:    true,
	}
	if err := traderAddressAggr.Dispatch(ctx, createAddressCmd); err != nil {
		return err
	}
	return nil
}

func PrepareOrderLines(
	ctx context.Context,
	shopID dot.ID,
	lines []*types.CreateOrderLine,
) ([]*ordermodel.OrderLine, error) {
	variantIDs := make([]dot.ID, 0, len(lines))
	if len(lines) > 40 {
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng có quá nhiều sản phẩm", nil)
	}
	for i, line := range lines {
		if line == nil {
			return nil, cm.Error(cm.InvalidArgument, "Invalid order line", nil)
		}
		productName := strings.ReplaceAll(line.ProductName, " ", "")
		if len(productName) < 2 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tên sản phẩm phải có ít nhất hai kí tự trở lên: '%v'", line.ProductName)
		}
		if line.VariantId == 0 {
			continue
		}
		variantIDs = append(variantIDs, line.VariantId)

		for j := 0; j < i; j++ {
			if line.VariantId == lines[j].VariantId {
				return nil, cm.Errorf(cm.InvalidArgument, nil,
					`Sản phẩm "%v" đã được nhập nhiều lần. Vui lòng kiểm tra lại.`, line.ProductName)
			}
		}
	}

	var variants []*catalog.ShopVariantWithProduct
	if len(variantIDs) > 0 {
		variantsQuery := &catalog.ListShopVariantsWithProductByIDsQuery{
			IDs:    variantIDs,
			ShopID: shopID,
		}
		if err := catalogQuery.Dispatch(ctx, variantsQuery); err != nil {
			return nil, err
		}
		variants = variantsQuery.Result.Variants
	}

	res := make([]*ordermodel.OrderLine, len(lines))
	for i, line := range lines {
		if line.VariantId == 0 {
			item, err := prepareOrderLine(line, shopID, nil)
			if err != nil {
				return nil, err
			}
			res[i] = item
			continue
		}

		var variant *catalog.ShopVariantWithProduct
		for _, v := range variants {
			if line.VariantId == v.ShopVariant.VariantID {
				variant = v

				// Check meta_fields
				mapMetaField := make(map[string]string)

				for _, metaField := range v.ShopProduct.MetaFields {
					mapMetaField[metaField.Key] = metaField.Value
				}

				if len(mapMetaField) != len(line.MetaFields) {
					return nil, cm.Errorf(cm.InvalidArgument, nil, "meta_fields không hợp lệ")
				}

				mapMetaFieldArg := make(map[string]bool)
				for _, metaField := range line.MetaFields {
					mapMetaFieldArg[metaField.Key] = true
					if _, ok := mapMetaField[metaField.Key]; !ok {
						return nil, cm.Errorf(cm.InvalidArgument, nil, "meta_field %v không tồn tại", metaField.Key)
					}
					if metaField.Value == "" || len(strings.TrimSpace(metaField.Value)) == 0 {
						return nil, cm.Errorf(cm.InvalidArgument, nil, "meta_field %v không được rỗng", metaField.Key)
					}

					metaField.Name = mapMetaField[metaField.Key]
				}

				// goal: check duplicate key in metaFields
				if len(mapMetaField) != len(mapMetaFieldArg) {
					return nil, cm.Errorf(cm.InvalidArgument, nil, "meta_fields không hợp lệ")
				}

				break
			}
		}
		if variant == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil,
				`Sản phẩm "%v" không được đăng bán. Vui lòng kiểm tra lại`,
				line.ProductName)
		}

		item, err := PrepareOrderLine(shopID, line, variant)
		if err != nil {
			return nil, err
		}
		res[i] = item
	}
	return res, nil
}

func UpdateOrder(ctx context.Context, claim *claims.ShopClaim, authPartner *identitymodel.Partner, q *types.UpdateOrderRequest) (*types.Order, error) {
	query := &ordermodelx.GetOrderQuery{
		OrderID: q.Id,
		ShopID:  claim.Shop.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	oldOrder := query.Result.Order

	switch oldOrder.Status {
	case status5.N:
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng đã hủy, không thể cập nhật đơn", nil)
	case status5.NS:
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng đã trả hàng, không thể cập nhật đơn", nil)
	case status5.P:
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng đã hoàn thành, không thể cập nhật đơn", nil)
	case status5.S:
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng đang xử lý, không thể cập nhật đơn", nil)
	}

	customerId := query.Result.Order.CustomerID

	if q.CustomerId != 0 {
		customerId = q.CustomerId
	}

	// make sure update always has Lines and FeeLines
	lines, err := PrepareOrderLines(ctx, claim.Shop.ID, q.Lines)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		lines = oldOrder.Lines
	}
	orderDiscount := q.OrderDiscount.Apply(oldOrder.OrderDiscount)
	feeLines := convertpb.PbOrderFeeLinesToModel(q.FeeLines)
	if len(feeLines) == 0 {
		feeLines = oldOrder.FeeLines
	} else {
		// calculate fee lines from shop_shipping_fee
		feeLines = ordermodel.GetFeeLinesWithFallback(feeLines, dot.NullInt{}, q.ShopShippingFee)
	}

	var basketValue, totalDiscount, totalAmount, totalItems int
	for _, line := range lines {
		basketValue += line.LineAmount
		totalItems += line.Quantity
	}
	totalLineDiscount := ordermodelx.SumOrderLineDiscount(lines)
	totalDiscount = totalLineDiscount + orderDiscount
	totalFee := ordermodel.CalcTotalFee(feeLines)

	// calculate shop_cod back from fee_lines
	shopShippingFee := ordermodel.GetShippingFeeFromFeeLines(feeLines)
	if q.ShopShippingFee.Apply(shopShippingFee) != shopShippingFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phí giao hàng không đúng").
			WithMetap("expected shop_shipping_cod (= SUM(fee_lines.amount) WHERE (type=shipping))", totalFee)
	}
	if q.TotalFee.Apply(totalFee) != totalFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng phí không đúng").
			WithMetap("expected total_fee (= SUM(fee_lines.amount))", totalFee)
	}
	totalAmount = basketValue - totalDiscount + totalFee

	if basketValue != q.BasketValue {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không đúng", nil).
			WithMetap("expected basket_value (= sum(lines.retail_price))", basketValue)
	}
	if totalAmount != q.TotalAmount {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không đúng", nil).
			WithMetap("expected total_amount (= basket_value + shop_shipping_fee - total_discount)", totalAmount)
	}
	if totalItems != q.TotalItems {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số lượng sản phẩm không đúng", nil).
			WithMetap("expected total_items", totalItems)
	}

	customerAddress, err := convertpb.OrderAddressToModel(q.CustomerAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	shippingAddress, err := convertpb.OrderAddressToModel(q.ShippingAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	billingAddress, err := convertpb.OrderAddressToModel(q.BillingAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ thanh toán không hợp lệ: %v", err)
	}

	shipping := q.ShopShipping
	var shopCod = q.ShopCod
	if q.Shipping != nil {
		shipping = q.Shipping
		shopCod = shipping.CodAmount
	}
	fakeOrder := &ordermodel.Order{}
	if err := convertpb.OrderShippingToModel(ctx, shipping, fakeOrder); err != nil {
		return nil, err
	}

	if q.CustomerId != 0 && q.Customer != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "customer_id và customer không được gửi cùng 1 lúc", err)
	}

	if q.CustomerId != 0 && q.CustomerAddress != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "customer_id và customer_address không được gửi cùng 1 lúc", err)
	}

	if q.CustomerId != 0 {
		query := &customering.GetCustomerByIDQuery{
			ID:     q.CustomerId,
			ShopID: claim.Shop.ID,
		}
		if err := customerQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).
				Wrapf(cm.NotFound, "customer_id %v không tồn tại", q.CustomerId).
				Throw()
		}

		q.Customer = &types.OrderCustomer{
			FullName: query.Result.FullName,
			Email:    query.Result.Email,
			Phone:    query.Result.Phone,
			Gender:   query.Result.Gender,
			Type:     query.Result.Type,
		}

		isHaveAddress := true
		getAddressQuery := &addressing.GetAddressActiveByTraderIDQuery{
			ShopID:   claim.Shop.ID,
			TraderID: q.CustomerId,
		}
		if err := traderAddressQuery.Dispatch(ctx, getAddressQuery); err != nil {
			switch cm.ErrorCode(err) {
			case cm.NotFound:
				isHaveAddress = false
			default:
				return nil, err
			}
		}
		if isHaveAddress {
			customerAddressResult, err := convertpb.PbShopAddress(ctx, getAddressQuery.Result, locationQuery)
			if err != nil {
				return nil, err
			}
			q.CustomerAddress = &types.OrderAddress{
				FullName:     customerAddressResult.FullName,
				Phone:        customerAddressResult.Phone,
				Email:        customerAddressResult.Email,
				District:     customerAddressResult.District,
				Ward:         customerAddressResult.Ward,
				Company:      customerAddressResult.Company,
				Address1:     customerAddressResult.Address1,
				Address2:     customerAddressResult.Address2,
				DistrictCode: customerAddressResult.DistrictCode,
				WardCode:     customerAddressResult.WardCode,
				Coordinates:  customerAddressResult.Coordinates,
			}
		} else {
			// TODO: handle when customerAddress is empty (from customer_id)
			// q.CustomerAddress = &pborder.OrderAddress{}
		}
	}

	customerAddress, err = convertpb.OrderAddressToModel(q.CustomerAddress)
	if err != nil {
		return nil, err
	}

	cmd := &ordermodelx.UpdateOrderCommand{
		ID:              q.Id,
		ShopID:          claim.Shop.ID,
		Customer:        convertpb.OrderCustomerToModel(q.Customer),
		CustomerAddress: customerAddress,
		BillingAddress:  billingAddress,
		ShippingAddress: shippingAddress,
		OrderNote:       q.OrderNote,
		ShippingNote:    cm.Coalesce(q.ShippingNote, fakeOrder.ShippingNote),
		ShopShippingFee: dot.Int(shopShippingFee),
		TryOn:           fakeOrder.TryOn,
		TotalWeight:     cm.CoalesceInt(q.TotalWeight, fakeOrder.TotalWeight),
		ShopShipping:    fakeOrder.ShopShipping,
		Lines:           lines,
		FeeLines:        feeLines,
		TotalFee:        dot.Int(totalFee),
		BasketValue:     q.BasketValue,
		TotalAmount:     q.TotalAmount,
		TotalItems:      q.TotalItems,
		OrderDiscount:   dot.Int(orderDiscount),
		TotalDiscount:   totalDiscount,
		ShopCOD:         shopCod,
		CustomerID:      customerId,
	}
	if authPartner != nil {
		cmd.PartnerID = authPartner.ID
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	// re-get order
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbOrder(query.Result.Order, nil, model.TagShop)
	result.ShopName = claim.Shop.Name

	return result, nil
}

func PrepareOrderLine(
	shopID dot.ID,
	m *types.CreateOrderLine,
	v *catalog.ShopVariantWithProduct,
) (*ordermodel.OrderLine, error) {
	retailPrice := v.GetRetailPrice()
	if m.RetailPrice != retailPrice {
		return nil, cm.Errorf(cm.FailedPrecondition, nil,
			`Có sự khác biệt về giá của sản phẩm "%v". Vui lòng kiểm tra lại. Giá đăng bán %v, giá đơn hàng %v`,
			v.ProductWithVariantName(), retailPrice, m.RetailPrice)
	}
	if m.PaymentPrice > m.RetailPrice {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Giá phải trả của sản phẩm "%v" không được lớn hơn giá đăng bán. Vui lòng kiểm tra lại.`,
			m.ProductName)
	}
	return prepareOrderLine(m, shopID, v)
}

func prepareOrderLine(
	m *types.CreateOrderLine,
	shopID dot.ID,
	v *catalog.ShopVariantWithProduct,
) (*ordermodel.OrderLine, error) {
	productName, ok := validate.NormalizeGenericName(m.ProductName)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Tên sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}
	if m.PaymentPrice < 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Giá phải trả của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}
	if m.Quantity <= 0 || m.Quantity >= 1000000 {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Số lượng của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}

	if m.VariantId == 0 && productName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp product_name hoặc variant_id")
	}

	metaFields := []*ordermodel.MetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &ordermodel.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
			Name:  metaField.Name,
		})
	}

	line := &ordermodel.OrderLine{
		ShopID:          shopID,
		IsOutsideEtop:   m.VariantId == 0,
		Quantity:        m.Quantity,
		ListPrice:       m.ListPrice,
		RetailPrice:     m.RetailPrice, // will be over-written if a variant is provided
		PaymentPrice:    m.PaymentPrice,
		LineAmount:      m.Quantity * m.RetailPrice,
		ImageURL:        m.ImageUrl,
		ProductName:     productName,
		Attributes:      convert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(m.Attributes),
		TotalDiscount:   0, // will be filled later
		TotalLineAmount: 0, // will be filled later
		MetaFields:      metaFields,
	}

	originalPrice := m.RetailPrice
	if v != nil {
		line.Code = v.Code
		line.VariantID = m.VariantId
		line.ProductID = v.ShopProduct.ProductID
		line.ProductName = model.CoalesceString2(v.ShopProduct.Name, v.ShopProduct.Name)

		line.ListPrice = v.GetListPrice()

		if len(v.ShopVariant.ImageURLs) > 0 {
			line.ImageURL = v.ShopVariant.ImageURLs[0]
		} else if v.ShopProduct != nil && len(v.ShopProduct.ImageURLs) > 0 {
			line.ImageURL = v.ShopProduct.ImageURLs[0]
		}

		if v.ShopVariant != nil {
			line.RetailPrice = v.GetRetailPrice()
			originalPrice = v.GetRetailPrice()
		}
		line.Attributes = convert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(v.ShopVariant.Attributes)
	}
	if line.RetailPrice <= 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Giá bán lẻ của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`,
			m.ProductName)
	}
	line.TotalDiscount = m.Quantity * (originalPrice - m.PaymentPrice)
	line.TotalLineAmount = m.Quantity * m.PaymentPrice
	return line, nil
}

func PrepareOrder(ctx context.Context, shopID dot.ID, m *types.CreateOrderRequest, lines []*ordermodel.OrderLine, userID dot.ID) (*ordermodel.Order, error) {
	if m.BasketValue <= 0 {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không hợp lệ", nil).
			WithMeta("reason", "basket_value <= 0")
	}
	if m.TotalAmount < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không hợp lệ", nil).
			WithMeta("reason", "total_amount < 0")
	}

	productIDs := make([]dot.ID, len(lines))
	variantIDs := make([]dot.ID, len(lines))

	// {0} and duplicated ids are allowed
	for i, line := range lines {
		productIDs[i] = line.ProductID
		variantIDs[i] = line.VariantID
	}

	// calculate fee lines from shop_shipping_fee
	feeLines := convertpb.PbOrderFeeLinesToModel(m.FeeLines)
	feeLines = ordermodel.GetFeeLinesWithFallback(feeLines, dot.Int(m.TotalFee), dot.Int(m.ShopShippingFee))
	totalFee := 0
	for _, line := range feeLines {
		totalFee += line.Amount
	}
	if m.TotalFee != 0 && m.TotalFee != totalFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng phí không đúng").
			WithMetap("expected total_fee (= SUM(fee_lines.amount))", totalFee)
	}

	// calculate shop_cod back from fee_lines
	shopShippingFee := 0
	for _, line := range feeLines {
		if line.Type == ordermodel.OrderFeeShipping {
			shopShippingFee += line.Amount
			if line.Amount < 0 {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Phí không được nhỏ hơn 0")
			}
		}
	}
	if m.ShopShippingFee != 0 && m.ShopShippingFee != shopShippingFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phí giao hàng không đúng").
			WithMetap("expected shop_shipping_cod (= SUM(fee_lines.amount) WHERE (type=shipping))", totalFee)
	}

	// caclulate basket_value and total_amount
	var basketValue, totalDiscount, totalAmount, totalItems int
	if len(lines) != 0 {
		for _, line := range lines {
			basketValue += line.LineAmount
			totalItems += line.Quantity
		}
	} else {
		basketValue = m.BasketValue
		totalItems = m.TotalItems
	}
	totalLineDiscount := ordermodelx.SumOrderLineDiscount(lines)
	orderDiscount := m.OrderDiscount
	totalDiscount = totalLineDiscount + orderDiscount
	if m.TotalDiscount.Apply(totalDiscount) != totalDiscount {
		return nil, cm.Error(cm.InvalidArgument, "Tổng giá trị giảm không đúng", nil).
			WithMetap("expected total_discount (= order_discount + sum(lines.total_discount))", totalDiscount)
	}
	if len(lines) != 0 && basketValue != m.BasketValue {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không đúng", nil).
			WithMetap("expected basket_value (= sum(lines.retail_price))", basketValue)
	}

	basketValue = m.BasketValue
	totalAmount = basketValue - totalDiscount + totalFee

	// if totalDiscount != int(m.TotalDiscount) {
	// 	return nil, cm.Error(cm.InvalidArgument, "Invalid TotalDiscount", nil)
	// }
	if totalAmount != m.TotalAmount {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không đúng", nil).
			WithMetap("expected total_amount (= basket_value + total_fee - total_discount)", totalAmount)
	}
	if totalItems != m.TotalItems {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số lượng sản phẩm không đúng", nil).
			WithMetap("expected total_items", totalItems)
	}

	confirm := m.ShConfirm.Apply(0)

	customerAddress, err := convertpb.OrderAddressToModel(m.CustomerAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	shippingAddress, err := convertpb.OrderAddressToModel(m.ShippingAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	billingAddress, err := convertpb.OrderAddressToModel(m.BillingAddress)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ thanh toán không hợp lệ: %v", err)
	}

	shipping := m.Shipping
	if shipping == nil {
		shipping = m.ShopShipping
	}
	paymentMethod := m.PaymentMethod
	var tryOn try_on.TryOnCode
	if shipping != nil {
		// return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping/shop_shipping")
		if m.PaymentMethod == 0 {
			paymentMethod = payment_method.COD
			if m.ShopCod == 0 && shipping.CodAmount.Apply(0) == 0 {
				paymentMethod = payment_method.Other
			}
		}
		tryOn = shipping.TryOn
	}
	if !model.VerifyPaymentMethod(paymentMethod) {
		return nil, cm.Error(cm.InvalidArgument, "Phương thức thanh toán không hợp lệ", nil)
	}
	// TODO: note that m.ExternalCode is validated at ext/partner, not here
	if m.ExternalId != "" && !validate.ExternalCode(m.ExternalId) {
		return nil, cm.Error(cm.InvalidArgument, "Mã đơn hàng external_id không hợp lệ", nil)
	}
	externalMeta, _ := jsonx.Marshal(m.ExternalMeta)
	referralMeta, _ := jsonx.Marshal(m.ReferralMeta)
	order := &ordermodel.Order{
		ID:         0,
		ShopID:     0,
		Code:       "", // will be filled by sqlstore
		EdCode:     m.ExternalCode,
		ProductIDs: productIDs,
		VariantIDs: variantIDs,
		CreatedBy:  userID,
		PartnerID:  0,
		Currency:   "",
		// Source:          m.Source.ToModel(),
		PaymentMethod:              paymentMethod,
		Customer:                   convertpb.OrderCustomerToModel(m.Customer),
		CustomerAddress:            customerAddress,
		BillingAddress:             billingAddress,
		ShippingAddress:            shippingAddress,
		CustomerName:               "",
		CustomerPhone:              m.Customer.Phone,
		CustomerEmail:              m.Customer.Email,
		CreatedAt:                  time.Now(),
		ProcessedAt:                time.Time{},
		UpdatedAt:                  time.Time{},
		ClosedAt:                   time.Time{},
		ConfirmedAt:                time.Time{},
		CancelledAt:                time.Time{},
		CancelReason:               "",
		CustomerConfirm:            0,
		ShopConfirm:                confirm,
		ConfirmStatus:              0,
		FulfillmentShippingStatus:  0,
		EtopPaymentStatus:          0,
		Status:                     0,
		FulfillmentShippingStates:  nil,
		FulfillmentPaymentStatuses: nil,
		Lines:                      lines,
		Discounts:                  convertpb.PbOrderDiscountsToModel(m.Discounts),
		TotalItems:                 m.TotalItems,
		BasketValue:                m.BasketValue,
		TotalWeight:                m.TotalWeight,
		TotalTax:                   0,
		OrderDiscount:              orderDiscount,
		TotalDiscount:              totalDiscount,
		ShopShippingFee:            shopShippingFee,
		TotalFee:                   totalFee,
		FeeLines:                   feeLines,
		ShopCOD:                    m.ShopCod,
		TotalAmount:                m.TotalAmount,
		OrderNote:                  m.OrderNote,
		ShopNote:                   "",
		ShippingNote:               m.ShippingNote,
		OrderSourceType:            0,
		OrderSourceID:              0,
		ExternalOrderID:            m.ExternalId,
		ReferenceURL:               m.ReferenceUrl,
		ExternalURL:                m.ExternalUrl,
		ShopShipping:               nil, // will be filled later
		IsOutsideEtop:              false,
		GhnNoteCode:                m.GhnNoteCode,
		TryOn:                      tryOn,
		CustomerNameNorm:           "",
		ProductNameNorm:            "",
		ExternalMeta:               externalMeta,
		ReferralMeta:               referralMeta,
		CustomerID:                 m.CustomerId,
	}
	if err = convertpb.OrderShippingToModel(ctx, shipping, order); err != nil {
		return nil, err
	}

	if order.ShopShipping != nil {
		shippingServiceCode := order.ShopShipping.GetShippingServiceCode()
		carrier := order.ShopShipping.GetShippingProvider()

		// handle etop custom service code here
		// TODO: refactor, move to shipping_provider
		// check ETOP service
		shippingServiceName, ok := etop_shipping_price.ParseEtopServiceCode(shippingServiceCode)
		if !ok {
			shippingServiceName, ok = ctrl.ParseServiceCode(carrier, shippingServiceCode)
		}
		if carrier != 0 && !ok {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Mã dịch vụ không hợp lệ. Vui lòng F5 thử lại hoặc liên hệ %v", wl.X(ctx).CSEmail)
		}
		order.ShopShipping.ExternalServiceName = shippingServiceName
	}

	return order, nil
}

func CancelOrder(ctx context.Context, shopID dot.ID, authPartnerID dot.ID, orderID dot.ID, cancelReason string, autoInventoryVoucher inventory_auto.AutoInventoryVoucher) (*types.OrderWithErrorsResponse, error) {
	getOrderQuery := &ordermodelx.GetOrderQuery{
		ShopID:             shopID,
		PartnerID:          authPartnerID,
		OrderID:            orderID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, getOrderQuery); err != nil {
		return nil, err
	}
	order := getOrderQuery.Result.Order

	switch order.Status {
	case status5.N:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil)
	case status5.P:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã hoàn thành.", nil)
	case status5.NS:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil)
	}

	// Do not allow cancel order if it had a shipnow fulfillment
	if order.FulfillmentType == ordertypes.ShippingTypeShipnow {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã tạo đơn giao hàng tức thì. Không thể hủy đơn.")
	}

	updateOrderCmd := &ordermodelx.UpdateOrdersStatusCommand{
		ShopID:        shopID,
		PartnerID:     authPartnerID,
		OrderIDs:      []dot.ID{orderID},
		ShopConfirm:   status3.N.Wrap(),
		ConfirmStatus: status3.N.Wrap(),
		CancelReason:  cancelReason,
		Status:        status5.N.Wrap(),
	}
	if err := bus.Dispatch(ctx, updateOrderCmd); err != nil {
		return nil, err
	}
	event := &ordering.OrderCancelledEvent{
		OrderID:              order.ID,
		ShopID:               shopID,
		AutoInventoryVoucher: autoInventoryVoucher,
	}
	if err := eventBus.Publish(ctx, event); err != nil {
		ll.Error("RaiseOrderCancelledEvent", l.Error(err))
	}

	// fulfillment errors when canceling order, it will appear in response
	var errs []error
	fulfillments := getOrderQuery.Result.Fulfillments
	order.CancelReason = cancelReason
	if len(fulfillments) > 0 {
		_errs, err := TryCancellingFulfillments(ctx, order, fulfillments)
		if err != nil {
			return nil, err
		}
		errs = _errs
	}

	// Get the order again
	if err := bus.Dispatch(ctx, getOrderQuery); err != nil {
		return nil, err
	}

	resp := &types.OrderWithErrorsResponse{
		Order:  convertpb.PbOrder(getOrderQuery.Result.Order, getOrderQuery.Result.Fulfillments, model.TagShop),
		Errors: cmapi.PbErrors(errs),

		FulfillmentErrors: cmapi.PbErrors(errs),
	}
	return resp, nil
}

var reSubdomain = regexp.MustCompile("^[a-z0-9]([a-z0-9-]{0,126}[a-z0-9])?$")

func validateExternalURL(recognizedURLs []string, externalURL string) error {
	eURL, err := url.Parse(externalURL)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Thông tin external_url không hợp lệ")
	}
	if eURL.Scheme != "http" && eURL.Scheme != "https" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thông tin external_url không hợp lệ").
			WithMeta("reason", "Chỉ cho phép http và https")
	}

	// Allow domains *.example.com
	for _, recognizedURL := range recognizedURLs {
		// We check *. first because url.Parse() accepts https://*.example.com
		if strings.HasPrefix(recognizedURL, "*.") {
			if strings.HasSuffix(eURL.Host, recognizedURL[1:]) {
				host := strings.TrimSuffix(eURL.Host, recognizedURL[1:])
				if reSubdomain.MatchString(host) {
					return nil
				}
			}
		} else if eURL.Host == recognizedURL {
			return nil
		}
	}
	return cm.Errorf(cm.InvalidArgument, nil, "Thông tin external_url không hợp lệ").
		WithMeta("reason", "Danh sách domain cần được đăng ký trước")
}
