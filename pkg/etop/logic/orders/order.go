package orderS

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"

	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
)

var ll = l.New()
var shippingManager *shipping_provider.ProviderManager

func CreateOrder(ctx context.Context, claim *claims.ShopClaim, authPartner *model.Partner, r *pborder.CreateOrderRequest) (*pborder.Order, error) {
	shipping := r.ShopShipping
	if r.Shipping != nil {
		shipping = r.Shipping
	}
	if shipping == nil {
		shipping = &pborder.OrderShipping{}
	}

	if shipping.Carrier == 0 {
		shipping.Carrier = shipping.ShippingProvider
	}
	if (r.Shipping != nil || r.ShopShipping != nil) &&
		!model.VerifyShippingProvider(shipping.Carrier.ToModel()) {
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

	// if r.GhnNoteCode.ToModel() == "" && shipping.TryOn == 0 {
	// 	return nil, cm.Error(cm.InvalidArgument, "Vui lòng chọn ghi chú xem hàng cho đơn hàng.", nil)
	// }
	// if shipping.TryOn == 0 {
	// 	shipping.TryOn = r.GhnNoteCode.ToTryOn()
	// } else if r.GhnNoteCode == 0 {
	// 	r.GhnNoteCode = ghn_note_code.FromTryOn(shipping.TryOn)
	// }

	src := r.Source.ToModel()
	if !model.VerifyOrderSource(src) {
		return nil, cm.Error(cm.InvalidArgument, "Invalid source", nil)
	}

	shop := claim.Shop
	lines, err := PrepareOrderLines(ctx, shop.ID, r.Lines)
	if err != nil {
		return nil, err
	}
	order, err := PrepareOrder(r, lines)
	if err != nil {
		return nil, err
	}
	if authPartner != nil {
		order.PartnerID = authPartner.ID
	}

	order.ShopID = shop.ID
	order.OrderSourceType = src
	cmd := &model.CreateOrderCommand{
		Order: order,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		// TODO: refactor
		if xerr, ok := err.(*cm.APIError); ok && xerr.Err != nil {
			msg := xerr.Err.Error()
			switch {
			case strings.Contains(msg, "order_shop_external_id_idx"):
				newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn hàng external_id đã tồn tại. Vui lòng kiểm tra lại.").
					WithMeta("duplicated", "external_id")
				orderQuery := &model.GetOrderQuery{
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
				orderQuery := &model.GetOrderQuery{
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
	result := pborder.PbOrder(order, nil, model.TagShop)
	result.ShopName = claim.Shop.Name
	return result, nil
}

func PrepareOrderLines(ctx context.Context, shopID int64, lines []*pborder.CreateOrderLine) ([]*model.OrderLine, error) {
	variantIDs := make([]int64, len(lines))
	if len(lines) > 40 {
		return nil, cm.Error(cm.InvalidArgument, "Đơn hàng có quá nhiều sản phẩm", nil)
	}
	for i, line := range lines {
		if line == nil {
			return nil, cm.Error(cm.InvalidArgument, "Invalid order line", nil)
		}
		if line.VariantId == 0 {
			continue
		}
		variantIDs[i] = line.VariantId

		for j := 0; j < i; j++ {
			if line.VariantId == lines[j].VariantId {
				return nil, cm.Error(cm.InvalidArgument,
					cm.F(`Sản phẩm "%v" đã được nhập nhiều lần. Vui lòng kiểm tra lại.`, line.ProductName), nil)
			}
		}
	}

	shopQuery := &model.GetShopQuery{
		ShopID: shopID,
	}
	if err := bus.Dispatch(ctx, shopQuery); err != nil {
		return nil, err
	}
	shop := shopQuery.Result

	var variants []*model.ShopVariantExtended
	if len(variantIDs) > 0 {
		variantsQuery := &model.GetAllShopVariantsQuery{
			ShopID:          shop.ID,
			VariantIDs:      variantIDs,
			ProductSourceID: shop.ProductSourceID,
		}
		if err := bus.Dispatch(ctx, variantsQuery); err != nil {
			return nil, err
		}
		variants = variantsQuery.Result.Variants
	}

	res := make([]*model.OrderLine, len(lines))
	for i, line := range lines {
		if line.VariantId == 0 {
			item, err := prepareOrderLine(line, shopID, nil, nil)
			if err != nil {
				return nil, err
			}
			res[i] = item
			continue
		}

		var prod *model.ShopVariantExtended
		for _, p := range variants {
			if line.VariantId == p.VariantID {
				prod = p
				break
			}
		}
		if prod == nil {
			return nil, cm.Error(cm.InvalidArgument,
				cm.F(`Sản phẩm "%v" không được đăng bán. Vui lòng kiểm tra lại`,
					line.ProductName), nil)
		}

		item, err := PrepareOrderLine(line, &model.VariantExtended{
			Variant: prod.Variant,
			Product: prod.Product,
		}, prod)
		if err != nil {
			return nil, err
		}
		res[i] = item
	}
	return res, nil
}

func UpdateOrder(ctx context.Context, claim *claims.ShopClaim, authPartner *model.Partner, q *pborder.UpdateOrderRequest) (*pborder.Order, error) {
	query := &model.GetOrderQuery{
		OrderID: q.Id,
		ShopID:  claim.Shop.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	oldOrder := query.Result.Order

	// make sure update always has Lines and FeeLines
	lines, err := PrepareOrderLines(ctx, claim.Shop.ID, q.Lines)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		lines = oldOrder.Lines
	}
	orderDiscount := 0
	if q.OrderDiscount != nil {
		orderDiscount = int(*q.OrderDiscount)
	} else {
		orderDiscount = oldOrder.OrderDiscount
	}
	feeLines := pborder.PbOrderFeeLinesToModel(q.FeeLines)
	if len(feeLines) == 0 {
		feeLines = oldOrder.FeeLines
	} else {
		// calculate fee lines from shop_shipping_fee
		feeLines = model.GetFeeLinesWithFallback(feeLines, nil, q.ShopShippingFee)
	}

	var basketValue, totalDiscount, totalAmount, totalItems int
	for _, line := range lines {
		basketValue += line.LineAmount
		totalItems += line.Quantity
	}
	totalLineDiscount := model.SumOrderLineDiscount(lines)
	totalDiscount = totalLineDiscount + orderDiscount
	totalFee := model.CalcTotalFee(feeLines)

	// calculate shop_cod back from fee_lines
	shopShippingFee := model.GetShippingFeeFromFeeLines(feeLines)
	if q.ShopShippingFee != nil {
		if int(*q.ShopShippingFee) != shopShippingFee {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Phí giao hàng không đúng").
				WithMetap("expected shop_shipping_cod (= SUM(fee_lines.amount) WHERE (type=shipping))", totalFee)
		}
	}
	if q.TotalFee != nil {
		if int(*q.TotalFee) != totalFee {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng phí không đúng").
				WithMetap("expected total_fee (= SUM(fee_lines.amount))", totalFee)
		}
	}
	totalAmount = basketValue - totalDiscount + totalFee

	if basketValue != int(q.BasketValue) {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không đúng", nil).
			WithMetap("expected basket_value (= sum(lines.retail_price))", basketValue)
	}
	if totalAmount != int(q.TotalAmount) {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không đúng", nil).
			WithMetap("expected total_amount (= basket_value + shop_shipping_fee - total_discount)", totalAmount)
	}
	if totalItems != int(q.TotalItems) {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số lượng sản phẩm không đúng", nil).
			WithMetap("expected total_items", totalItems)
	}

	customerAddress, err := q.CustomerAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	shippingAddress, err := q.ShippingAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	billingAddress, err := q.BillingAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ thanh toán không hợp lệ: %v", err)
	}

	shipping := q.ShopShipping
	var shopCod = q.ShopCod
	if q.Shipping != nil {
		shipping = q.Shipping
		shopCod = shipping.CodAmount
	}
	fakeOrder := &model.Order{}
	if err := shipping.ToModel(fakeOrder); err != nil {
		return nil, err
	}

	cmd := &model.UpdateOrderCommand{
		ID:              q.Id,
		ShopID:          claim.Shop.ID,
		Customer:        q.Customer.ToModel(),
		CustomerAddress: customerAddress,
		BillingAddress:  billingAddress,
		ShippingAddress: shippingAddress,
		OrderNote:       q.OrderNote,
		ShippingNote:    cm.Coalesce(q.ShippingNote, fakeOrder.ShippingNote),
		ShopShippingFee: cm.PInt(shopShippingFee),
		TryOn:           fakeOrder.TryOn,
		TotalWeight:     cm.CoalesceInt(int(q.TotalWeight), fakeOrder.TotalWeight),
		ShopShipping:    fakeOrder.ShopShipping,
		Lines:           lines,
		FeeLines:        feeLines,
		TotalFee:        cm.PInt(totalFee),
		BasketValue:     int(q.BasketValue),
		TotalAmount:     int(q.TotalAmount),
		TotalItems:      int(q.TotalItems),
		OrderDiscount:   cm.PInt(orderDiscount),
		TotalDiscount:   totalDiscount,
		ShopCOD:         cm.PInt32(shopCod),
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
	result := pborder.PbOrder(query.Result.Order, nil, model.TagShop)
	result.ShopName = claim.Shop.Name

	return result, nil
}

func PrepareOrderLine(
	m *pborder.CreateOrderLine,
	v *model.VariantExtended, sp *model.ShopVariantExtended,
) (*model.OrderLine, error) {
	if int(m.RetailPrice) != sp.ShopVariant.RetailPrice {
		return nil, cm.Error(cm.FailedPrecondition, cm.F(
			`Có sự khác biệt về giá của sản phẩm "%v". Vui lòng kiểm tra lại. Giá đăng bán %v, giá đơn hàng %v`,
			v.GetFullName(), sp.ShopVariant.RetailPrice, m.RetailPrice), nil)
	}
	if m.PaymentPrice > m.RetailPrice {
		return nil, cm.Error(cm.InvalidArgument, cm.F(
			`Giá phải trả của sản phẩm "%v" không được lớn hơn giá đăng bán. Vui lòng kiểm tra lại.`,
			m.ProductName), nil)
	}
	return prepareOrderLine(m, sp.ShopVariant.ShopID, v, sp)
}

func prepareOrderLine(m *pborder.CreateOrderLine, shopID int64, v *model.VariantExtended, sp *model.ShopVariantExtended) (*model.OrderLine, error) {
	productName, ok := validate.NormalizeGenericName(m.ProductName)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Tên sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}
	if m.PaymentPrice < 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Giá phải trả của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}
	if m.Quantity <= 0 || m.Quantity >= 1000 {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Số lượng của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`, m.ProductName)
	}

	if m.VariantId == 0 && productName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp product_name hoặc variant_id")
	}

	line := &model.OrderLine{
		ShopID:          shopID,
		IsOutsideEtop:   m.VariantId == 0,
		Quantity:        int(m.Quantity),
		ListPrice:       int(m.ListPrice),
		RetailPrice:     int(m.RetailPrice), // will be over-written if a variant is provided
		PaymentPrice:    int(m.PaymentPrice),
		LineAmount:      int(m.Quantity) * int(m.RetailPrice),
		ProductName:     productName,
		Attributes:      pborder.PbAttributesToModel(m.Attributes),
		TotalDiscount:   0, // will be filled later
		TotalLineAmount: 0, // will be filled later
	}

	originalPrice := int(m.RetailPrice)
	if v != nil && sp != nil {
		var externalVariantID string
		if v.VariantExternal != nil {
			externalVariantID = v.VariantExternal.ExternalID
		}

		line.VariantID = m.VariantId
		line.SupplierID = v.SupplierID
		line.ExternalVariantID = externalVariantID
		line.ProductID = sp.Product.ID
		line.ProductName = model.CoalesceString2(sp.ShopProduct.Name, sp.Product.Name)

		line.WholesalePrice = int(v.WholesalePrice)
		line.WholesalePrice0 = int(v.WholesalePrice0)
		line.ListPrice = int(v.ListPrice)

		if len(sp.ShopVariant.ImageURLs) > 0 {
			line.ImageURL = sp.ShopProduct.ImageURLs[0]
		} else if sp.ShopProduct != nil && len(sp.ShopProduct.ImageURLs) > 0 {
			line.ImageURL = sp.ShopProduct.ImageURLs[0]
		} else if sp.Product != nil && len(sp.Product.ImageURLs) > 0 {
			line.ImageURL = sp.Product.ImageURLs[0]
		}

		if sp.ShopVariant != nil {
			line.RetailPrice = int(sp.ShopVariant.RetailPrice)
			originalPrice = sp.ShopVariant.RetailPrice
		} else {
			originalPrice = v.WholesalePrice
		}
		line.Attributes = v.Attributes
	}
	if line.RetailPrice <= 0 {
		return nil, cm.Error(cm.InvalidArgument, cm.F(
			`Giá bán lẻ của sản phẩm "%v" không hợp lệ. Vui lòng kiểm tra lại.`,
			m.ProductName), nil)
	}
	line.TotalDiscount = int(m.Quantity) * (originalPrice - int(m.PaymentPrice))
	line.TotalLineAmount = int(m.Quantity) * int(m.PaymentPrice)
	return line, nil
}

func PrepareOrder(m *pborder.CreateOrderRequest, lines []*model.OrderLine) (*model.Order, error) {
	if m.Customer == nil {
		return nil, cm.Error(cm.InvalidArgument, "Missing Customer", nil)
	}
	if m.CustomerAddress == nil {
		return nil, cm.Error(cm.InvalidArgument, "Missing CustomerAddress", nil)
	}
	if m.BillingAddress == nil {
		return nil, cm.Error(cm.InvalidArgument, "Missing BillingAddress", nil)
	}
	if m.ShippingAddress == nil {
		return nil, cm.Error(cm.InvalidArgument, "Missing ShippingAddress", nil)
	}
	if m.BasketValue <= 0 {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không hợp lệ", nil).
			WithMeta("reason", "basket_value <= 0")
	}
	if m.TotalAmount < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không hợp lệ", nil).
			WithMeta("reason", "total_amount < 0")
	}

	productIDs := make([]int64, len(lines))
	variantIDs := make([]int64, len(lines))

	// {0} and duplicated ids are allowed
	for i, line := range lines {
		productIDs[i] = line.ProductID
		variantIDs[i] = line.VariantID
	}

	// calculate fee lines from shop_shipping_fee
	feeLines := pborder.PbOrderFeeLinesToModel(m.FeeLines)
	feeLines = model.GetFeeLinesWithFallback(feeLines, &m.TotalFee, &m.ShopShippingFee)
	totalFee := 0
	for _, line := range feeLines {
		totalFee += line.Amount
	}
	if m.TotalFee != 0 && int(m.TotalFee) != totalFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng phí không đúng").
			WithMetap("expected total_fee (= SUM(fee_lines.amount))", totalFee)
	}

	// calculate shop_cod back from fee_lines
	shopShippingFee := 0
	for _, line := range feeLines {
		if line.Type == model.OrderFeeShipping {
			shopShippingFee += line.Amount
			if line.Amount < 0 {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Phí không được nhỏ hơn 0")
			}
		}
	}
	if m.ShopShippingFee != 0 && int(m.ShopShippingFee) != shopShippingFee {
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
		basketValue = int(m.BasketValue)
		totalItems = int(m.TotalItems)
	}
	totalLineDiscount := model.SumOrderLineDiscount(lines)
	orderDiscount := int(m.OrderDiscount)
	totalDiscount = totalLineDiscount + orderDiscount
	if m.TotalDiscount != nil {
		if int32(totalDiscount) != *m.TotalDiscount {
			return nil, cm.Error(cm.InvalidArgument, "Tổng giá trị giảm không đúng", nil).
				WithMetap("expected total_discount (= order_discount + sum(lines.total_discount))", totalDiscount)
		}
	}
	if len(lines) != 0 && basketValue != int(m.BasketValue) {
		return nil, cm.Error(cm.InvalidArgument, "Giá trị đơn hàng không đúng", nil).
			WithMetap("expected basket_value (= sum(lines.retail_price))", basketValue)
	}

	basketValue = int(m.BasketValue)
	totalAmount = basketValue - totalDiscount + int(totalFee)

	// if totalDiscount != int(m.TotalDiscount) {
	// 	return nil, cm.Error(cm.InvalidArgument, "Invalid TotalDiscount", nil)
	// }
	if totalAmount != int(m.TotalAmount) {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số tiền không đúng", nil).
			WithMetap("expected total_amount (= basket_value + total_fee - total_discount)", totalAmount)
	}
	if totalItems != int(m.TotalItems) {
		return nil, cm.Error(cm.InvalidArgument, "Tổng số lượng sản phẩm không đúng", nil).
			WithMetap("expected total_items", totalItems)
	}

	var confirm model.Status3 = 0
	if s := m.ShConfirm.ToModel(); s != nil {
		confirm = *s
	}

	customerAddress, err := m.CustomerAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ khách hàng không hợp lệ: %v", err)
	}
	shippingAddress, err := m.ShippingAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	billingAddress, err := m.BillingAddress.ToModel()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ thanh toán không hợp lệ: %v", err)
	}

	shipping := m.Shipping
	if shipping == nil {
		shipping = m.ShopShipping
	}
	paymentMethod := m.PaymentMethod
	var tryOn model.TryOn
	if shipping != nil {
		// return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping/shop_shipping")
		if m.PaymentMethod == "" {
			paymentMethod = model.PaymentMethodCOD
			if m.ShopCod == 0 && (shipping.CodAmount == nil || *shipping.CodAmount == 0) {
				paymentMethod = model.PaymentMethodOther
			}
		}
		tryOn = shipping.TryOn.ToModel()
	}
	if !model.VerifyPaymentMethod(paymentMethod) {
		return nil, cm.Error(cm.InvalidArgument, "Phương thức thanh toán không hợp lệ", nil)
	}
	// TODO: note that m.ExternalCode is validated at ext/partner, not here
	if m.ExternalId != "" && !validate.ExternalCode(m.ExternalId) {
		return nil, cm.Error(cm.InvalidArgument, "Mã đơn hàng external_id không hợp lệ", nil)
	}

	order := &model.Order{
		ID:          0,
		ShopID:      0,
		Code:        "", // will be filled by sqlstore
		EdCode:      m.ExternalCode,
		ProductIDs:  productIDs,
		VariantIDs:  variantIDs,
		SupplierIDs: nil,
		PartnerID:   0,
		Currency:    "",
		// Source:          m.Source.ToModel(),
		PaymentMethod:              paymentMethod,
		Customer:                   m.Customer.ToModel(),
		CustomerAddress:            customerAddress,
		BillingAddress:             billingAddress,
		ShippingAddress:            shippingAddress,
		CustomerName:               "",
		CustomerPhone:              m.ShippingAddress.Phone,
		CustomerEmail:              m.Customer.Email,
		CreatedAt:                  time.Now(),
		ProcessedAt:                time.Time{},
		UpdatedAt:                  time.Time{},
		ClosedAt:                   time.Time{},
		ConfirmedAt:                time.Time{},
		CancelledAt:                time.Time{},
		CancelReason:               "",
		CustomerConfirm:            0,
		ExternalConfirm:            0,
		ShopConfirm:                confirm,
		ConfirmStatus:              0,
		FulfillmentShippingStatus:  0,
		CustomerPaymentStatus:      0,
		EtopPaymentStatus:          0,
		Status:                     0,
		FulfillmentShippingStates:  nil,
		FulfillmentPaymentStatuses: nil,
		Lines:                      lines,
		Discounts:                  pborder.PbOrderDiscountsToModel(m.Discounts),
		TotalItems:                 int(m.TotalItems),
		BasketValue:                int(m.BasketValue),
		TotalWeight:                int(m.TotalWeight),
		TotalTax:                   0,
		OrderDiscount:              orderDiscount,
		TotalDiscount:              int(totalDiscount),
		ShopShippingFee:            shopShippingFee,
		TotalFee:                   totalFee,
		FeeLines:                   feeLines,
		ShopCOD:                    int(m.ShopCod),
		TotalAmount:                int(m.TotalAmount),
		OrderNote:                  m.OrderNote,
		ShopNote:                   "",
		ShippingNote:               m.ShippingNote,
		OrderSourceType:            "",
		OrderSourceID:              0,
		ExternalOrderID:            m.ExternalId,
		ReferenceURL:               m.ReferenceUrl,
		ExternalURL:                m.ExternalUrl,
		ShopShipping:               nil, // will be filled later
		IsOutsideEtop:              false,
		Fulfillments:               nil,
		ExternalData:               nil,
		GhnNoteCode:                m.GhnNoteCode.ToModel(),
		TryOn:                      tryOn,
		CustomerNameNorm:           "",
		ProductNameNorm:            "",
	}
	if err := shipping.ToModel(order); err != nil {
		return nil, err
	}

	if order.ShopShipping != nil {
		shippingServiceCode := order.ShopShipping.GetShippingServiceCode()
		carrierName := order.ShopShipping.GetShippingProvider()

		// handle etop custom service code here
		// TODO: refactor, move to shipping_provider
		// check ETOP service
		shippingServiceName, ok := etop_shipping_price.ParseEtopServiceCode(shippingServiceCode)
		if !ok {
			shippingServiceName, ok = ctrl.ParseServiceCode(carrierName, shippingServiceCode)
		}
		if carrierName != "" && !ok {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Mã dịch vụ không hợp lệ. Vui lòng F5 thử lại hoặc liên hệ hotro@etop.vn")
		}
		order.ShopShipping.ExternalServiceName = shippingServiceName
	}

	return order, nil
}

func CancelOrder(ctx context.Context, shopID int64, authPartnerID int64, orderID int64, cancelReason string) (*pborder.OrderWithErrorsResponse, error) {
	getOrderQuery := &model.GetOrderQuery{
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
	case model.S5Negative:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil)
	case model.S5Positive:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã hoàn thành.", nil)
	case model.S5NegSuper:
		return nil, cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil)
	}

	// MUSTDO: Handle confirm status

	//if order.ConfirmStatus == model.S3Negative ||
	//	order.ShopConfirm == model.S3Negative {
	//	return cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil)
	//}

	updateOrderCmd := &model.UpdateOrdersStatusCommand{
		ShopID:        shopID,
		PartnerID:     authPartnerID,
		OrderIDs:      []int64{orderID},
		ShopConfirm:   model.S3Negative.P(),
		ConfirmStatus: model.S3Negative.P(),
		CancelReason:  cancelReason,
	}
	if err := bus.Dispatch(ctx, updateOrderCmd); err != nil {
		return nil, err
	}

	// fulfillment errors when canceling order, it will appear in response
	var errs []error
	fulfillments := getOrderQuery.Result.Fulfillments
	if len(fulfillments) > 0 {
		err, _errs := TryCancellingFulfillments(ctx, order, fulfillments)
		if err != nil {
			return nil, err
		}
		errs = _errs
	}

	// Get the order again
	if err := bus.Dispatch(ctx, getOrderQuery); err != nil {
		return nil, err
	}

	resp := &pborder.OrderWithErrorsResponse{
		Order:  pborder.PbOrder(getOrderQuery.Result.Order, getOrderQuery.Result.Fulfillments, model.TagShop),
		Errors: pbcm.PbErrors(errs),

		FulfillmentErrors: pbcm.PbErrors(errs),
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
