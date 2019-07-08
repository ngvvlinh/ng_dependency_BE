package sqlstore

import (
	"context"
	"fmt"
	"strconv"

	"etop.vn/api/main/shipnow"

	"github.com/lib/pq"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
	ordermodelx "etop.vn/backend/pkg/services/ordering/modelx"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
	shipmodelx "etop.vn/backend/pkg/services/shipping/modelx"
	shipmodely "etop.vn/backend/pkg/services/shipping/modely"
)

func init() {
	bus.AddHandlers("sql",
		CreateFulfillments,
		CreateOrder,
		CreateOrders,
		GetFulfillment,
		GetOrder,
		GetOrders,
		SimpleGetOrdersByIDs,
		UpdateFulfillment,
		UpdateFulfillments,
		UpdateFulfillmentsStatus,
		UpdateOrder,
		UpdateOrdersStatus,
		GetFulfillments,
		GetFulfillmentExtended,
		GetFulfillmentExtendeds,
		GetFulfillmentsCallbackLogs,
		SyncUpdateFulfillments,
		VerifyOrdersByEdCode,
		UpdateFulfillmentsShippingState,
		UpdateOrderPaymentStatus,
		GetUnCompleteFulfillments,
		UpdateFulfillmentsWithoutTransaction,
		AdminUpdateFulfillment,
	)
}

var filterOrderWhitelist = FilterWhitelist{
	Arrays:   []string{"fulfillment.shipping_code", "fulfillment.shipping_state"},
	Contains: []string{"customer.name", "product.name"},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"shop.id", "code", "source", "external_code", "external_id", "customer.phone"},
	Numbers:  []string{"total_amount", "chargeable_weight"},
	Status:   []string{"status", "confirm_status", "fulfillment.shipping_status", "etop_payment_status"},
	PrefixOrRename: map[string]string{
		"fulfillment.shipping_code":   `"order".fulfillment_shipping_codes`,
		"fulfillment.shipping_state":  `"order".fulfillment_shipping_states`,
		"fulfillment.shipping_status": `"order".fulfillment_shipping_status`,

		"source":            `"order".order_source_type`,
		"shop.id":           `"order".shop_id`,
		"customer.name":     `"order".customer_name_norm`,
		"customer.phone":    `"order".customer_phone`,
		"product.name":      `"order".product_name_norm`,
		"external_code":     `"order".ed_code`,
		"external_id":       `"order".external_order_id`,
		"chargeable_weight": `"order".total_weight`,
	},
}

var filterFulfillmentWhitelist = FilterWhitelist{
	Arrays:   nil,
	Bools:    []string{"include_insurance"},
	Contains: []string{"customer.name"},
	Dates:    []string{"created_at", "updated_at"},
	Equals: []string{
		"shipping_code", "shop.id", "carrier",
		"order.code", "order.external_code", "order.external_id",
		"shipping_state", "customer.phone", "money_transaction.id",
		"address_to.province_code", "address_to.district_code", "address_to.ward_code",
	},
	Numbers: []string{"total_cod_amount", "cod_amount", "shipping_fee_shop", "shipping_service_fee", "basket_value", "chargeable_weight"},
	Status:  []string{"shipping_status", "etop_payment_status"},
	PrefixOrRename: map[string]string{
		"carrier":                  "f.shipping_provider",
		"shipping_code":            "f",
		"total_cod_amount":         "f", // @deprecated: use cod_amount
		"cod_amount":               "f.total_cod_amount",
		"shipping_fee_shop":        "f", // @deprecated: use shipping_service_fee
		"shipping_service_fee":     "f.shipping_fee_shop",
		"shop.id":                  "f.shop_id",
		"shipping_state":           "f",
		"money_transaction.id":     "f.money_transaction_id",
		"address_to.province_code": "f.address_to_province_code",
		"address_to.district_code": "f.address_to_district_code",
		"address_to.ward_code":     "f.address_to_ward_code",
		"created_at":               "f",
		"updated_at":               "f",
		"basket_value":             "f",
		"chargeable_weight":        "f.total_weight",
		"include_insurance":        "f",
		"etop_payment_status":      "f",

		"customer.name":       "o.customer_name_norm",
		"customer.phone":      "o.customer_phone",
		"order.code":          "o.code",
		"order.external_code": "o.ed_code",
		"order.external_id":   "o.external_order_id",
	},
}

func GetOrder(ctx context.Context, query *ordermodelx.GetOrderQuery) error {
	if query.OrderID == 0 && query.ExternalID == "" && query.Code == "" {
		return cm.Error(cm.InvalidArgument, "Missing id or code", nil)
	}

	s := x.Table("order")
	if query.OrderID != 0 {
		s = s.Where("id = ?", query.OrderID)
	}
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}
	if query.ExternalID != "" {
		s = s.Where("external_order_id = ?", query.ExternalID)
	}
	if query.Code != "" {
		s = s.Where("code = ?", query.Code)
	}

	order := new(ordermodel.Order)
	if err := s.ShouldGet(order); err != nil {
		return err
	}

	if query.IncludeFulfillment {
		var ffms []*ordermodelx.Fulfillment

		s := x.Table("fulfillment").
			Where("order_id = ?", order.ID).
			OrderBy("id")
		if query.ShopID != 0 {
			s = s.Where("shop_id = ?", query.ShopID)
		}
		var shipments []*shipmodel.Fulfillment
		if err := s.Find((*shipmodel.Fulfillments)(&shipments)); err != nil {
			return err
		}
		for _, sm := range shipments {
			ffms = append(ffms, &ordermodelx.Fulfillment{Shipment: sm})
		}
		query.Result.Fulfillments = shipments

		var shipnows []*shipnowmodel.ShipnowFulfillment
		if err := x.Table("shipnow_fulfillment").In("id", order.FulfillmentIDs).Find((*shipnowmodel.ShipnowFulfillments)(&shipnows)); err != nil {
			return err
		}
		for _, sn := range shipnows {
			_snCore := shipnowconvert.Shipnow(sn)
			ffms = append(ffms, &ordermodelx.Fulfillment{Shipnow: _snCore})
		}
		query.Result.XFulfillments = ffms
	}

	query.Result.Order = order
	return nil
}

func GetOrders(ctx context.Context, query *ordermodelx.GetOrdersQuery) error {
	s := x.Table("order")
	if query.ShopIDs != nil {
		s = s.InOrEqIDs("shop_id", query.ShopIDs)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}

	s, _, err := Filters(s, query.Filters, filterOrderWhitelist)
	if err != nil {
		return err
	}
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}

	var orders ordermodel.Orders
	{

		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "", "created_at": "", "updated_at": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("id", query.IDs)
		}
		if err := s2.Find(&orders); err != nil {
			return err
		}
		query.Result.Orders = make([]ordermodelx.OrderWithFulfillments, len(orders))
		for i, order := range orders {
			query.Result.Orders[i] = ordermodelx.OrderWithFulfillments{Order: order}
		}
	}
	if len(query.Filters) == 0 {
		total, err := s.Count(&ordermodel.Order{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}

	orderIds := make([]int64, len(query.Result.Orders))
	shopIdsMap := make(map[int64]int64)
	for i, o := range query.Result.Orders {
		orderIds[i] = o.ID
		shopIdsMap[o.ShopID] = o.ShopID
	}
	var fulfillments []*shipmodel.Fulfillment
	if err := x.Table("fulfillment").
		In("order_id", orderIds).
		Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return err
	}

	var shipnows []*shipnowmodel.ShipnowFulfillment
	if err := x.Table("shipnow_fulfillment").
		Where("status != ?", model.S5Negative).
		Where("order_ids && ?", pq.Int64Array(orderIds)).
		Find((*shipnowmodel.ShipnowFulfillments)(&shipnows)); err != nil {
		return err
	}

	orderShipments := make(map[int64][]*shipmodel.Fulfillment)
	for _, ffm := range fulfillments {
		orderShipments[ffm.OrderID] = append(orderShipments[ffm.OrderID], ffm)
	}
	orderShipnows := make(map[int64][]*shipnow.ShipnowFulfillment)
	for _, ffm := range shipnows {
		for _, orderID := range ffm.OrderIDs {
			sn := shipnowconvert.Shipnow(ffm)
			orderShipnows[orderID] = append(orderShipnows[orderID], sn)
		}
	}

	// getShop
	shopIds := make([]int64, 0, len(shopIdsMap))
	for _, shopId := range shopIdsMap {
		shopIds = append(shopIds, shopId)
	}
	shopQuery := &model.GetShopsQuery{
		ShopIDs: shopIds,
	}
	if err := bus.Dispatch(ctx, shopQuery); err != nil {
		return err
	}
	query.Result.Shops = shopQuery.Result.Shops

	for i := range query.Result.Orders {
		order := &query.Result.Orders[i] // it's not a pointer

		shipnows := orderShipnows[order.ID]
		for _, sn := range shipnows {
			order.Fulfillments = append(order.Fulfillments, &ordermodelx.Fulfillment{Shipnow: sn})
		}
		shipments := orderShipments[order.ID]
		for _, sm := range shipments {
			order.Fulfillments = append(order.Fulfillments, &ordermodelx.Fulfillment{Shipment: sm})
		}
	}

	return nil
}

func VerifyOrdersByEdCode(ctx context.Context, query *ordermodelx.VerifyOrdersByEdCodeQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if len(query.EdCodes) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing codes", nil)
	}

	// Currently we only support getting active orders
	if !query.OnlyActiveOrders {
		return cm.Error(cm.InvalidArgument, "Unexpected", nil)
	}

	s := x.SQL(`SELECT DISTINCT ed_code FROM "order"`).
		Where("shop_id = ?", query.ShopID).
		Where("ed_code = ANY(?)", pq.StringArray(query.EdCodes)).
		Where("shop_confirm != -1")
	sql, args, err := s.Build()
	if err != nil {
		return err
	}

	sql2 := fmt.Sprintf(
		"SELECT array_agg(ed_code) FROM (%v) AS s",
		sql,
	)
	return x.QueryRow(sql2, args...).Scan((*pq.StringArray)(&query.Result.EdCodes))
}

func SimpleGetOrdersByIDs(ctx context.Context, query *ordermodelx.SimpleGetOrdersByExternalIDsQuery) error {
	if query.SourceType == "" {
		return cm.Error(cm.InvalidArgument, "Missing ExternalProvider", nil)
	}

	s := x.Table("order").
		Where("order_source_type = ?", query.SourceType)
	if query.SourceID != 0 {
		s.Where("order_source_id = ?", query.SourceID)
	}

	s = s.In("external_order_id", query.ExternalIDs)
	return s.Find((*ordermodel.Orders)(&query.Result.Orders))
}

func UpdateOrdersStatus(ctx context.Context, cmd *ordermodelx.UpdateOrdersStatusCommand) error {
	if cmd.ShopConfirm != nil && *cmd.ShopConfirm != -1 && cmd.CancelReason != "" {
		return cm.Error(cm.InvalidArgument, "Cancel reason provided but confirm status is not cancel", nil)
	}
	if cmd.ShopConfirm != nil && *cmd.ShopConfirm == -1 && cmd.CancelReason == "" {
		return cm.Error(cm.InvalidArgument, "Cancel orders must provide cancel reason", nil)
	}

	s := x.Table("order").
		Where("status = 0 OR status = 2 OR status IS NULL"). // Only update orders in 'processing'
		InOrEqIDs("id", cmd.OrderIDs)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	if cmd.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.PartnerID)
	}

	m := M{}
	if cmd.ShopConfirm != nil {
		m["shop_confirm"] = cmd.ShopConfirm
	}
	if cmd.ConfirmStatus != nil {
		m["confirm_status"] = cmd.ConfirmStatus
	}
	if cmd.Status != nil {
		m["status"] = cmd.Status
	}
	if len(m) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	if cmd.CancelReason != "" {
		m["cancel_reason"] = cmd.CancelReason
	}

	if updated, err := s.UpdateMap(m); err != nil {
		return err
	} else if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	} else {
		cmd.Result.Updated = int(updated)
	}
	return nil
}

func CreateOrder(ctx context.Context, cmd *ordermodelx.CreateOrderCommand) error {
	order := cmd.Order
	if order.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	shop, err := generateShopCode(ctx, order.ShopID)
	if err != nil {
		return err
	}

	return inTransaction(func(x Qx) error {
		order.ID = cm.NewID()
		// generate order code
		code, errCode := GenerateCode(ctx, x, model.CodeTypeOrder, shop.Code)
		if errCode != nil {
			return errCode
		}
		order.Code = code
		if err := order.BeforeInsert(); err != nil {
			return err
		}
		if err := x.Table("order").ShouldInsert(order); err != nil {
			return err
		}

		// there is no line
		if len(order.Lines) == 0 {
			return nil
		}
		fn := gencode.GenerateLineCode(order.Code, len(order.Lines))
		for i, line := range order.Lines {
			line.OrderID = order.ID
			line.Code = fn(i)
			if err := x.Table("order_line").
				ShouldInsert(line); err != nil {
				return err
			}
		}
		return nil
	})
}

func CreateOrders(ctx context.Context, cmd *ordermodelx.CreateOrdersCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if len(cmd.Orders) == 0 {
		return cm.Error(cm.InvalidArgument, "Nothing to create", nil)
	}
	for _, order := range cmd.Orders {
		if order.ShopID != cmd.ShopID {
			return cm.Error(cm.InvalidArgument, "Invalid Name", nil)
		}
	}

	shop, err := generateShopCode(ctx, cmd.ShopID)
	if err != nil {
		return err
	}

	errs := make([]error, len(cmd.Orders))
	for i, order := range cmd.Orders {
		errs[i] = inTransaction(func(x Qx) error {
			order.ID = cm.NewID()
			code, errCode := GenerateCode(ctx, x, model.CodeTypeOrder, shop.Code)
			if errCode != nil {
				return errCode
			}
			order.Code = code
			if err := order.BeforeInsert(); err != nil {
				return err
			}
			if err := x.Table("order").ShouldInsert(order); err != nil {
				return err
			}

			fn := gencode.GenerateLineCode(order.Code, len(order.Lines))
			for i, line := range order.Lines {
				line.OrderID = order.ID
				line.Code = fn(i)
				if err := x.Table("order_line").
					ShouldInsert(line); err != nil {
					return err
				}
			}
			return nil
		})
	}
	cmd.Result.Errors = errs
	return nil
}

func generateShopCode(ctx context.Context, shopID int64) (*model.Shop, error) {
	queryShop := &model.GetShopQuery{
		ShopID: shopID,
	}
	if err := GetShop(ctx, queryShop); err != nil {
		return nil, err
	}
	shop := queryShop.Result

	// generate shop code if not existed
	if shop.Code == "" {
		// update shop Code
		var shopUpdate = &model.Shop{
			ID: shop.ID,
		}
		shopCode := gencode.GenerateShopCode()
		shopUpdate.Code = shopCode

		cmd := &model.UpdateShopCommand{
			Shop: shopUpdate,
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
		shop.Code = shopCode
	}
	return shop, nil
}

func UpdateOrder(ctx context.Context, cmd *ordermodelx.UpdateOrderCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	query := &ordermodelx.GetOrderQuery{
		OrderID: cmd.ID,
		ShopID:  cmd.ShopID,
	}
	if err := GetOrder(ctx, query); err != nil {
		return nil
	}
	oldOrder := query.Result.Order

	order := &ordermodel.Order{
		ID:              cmd.ID,
		ShopID:          cmd.ShopID,
		Customer:        cmd.Customer,
		CustomerAddress: cmd.CustomerAddress,
		BillingAddress:  cmd.BillingAddress,
		ShippingAddress: cmd.ShippingAddress,
		OrderNote:       cmd.OrderNote,
		ShippingNote:    cmd.ShippingNote,
		ShopShipping:    cmd.ShopShipping,
		GhnNoteCode:     model.GHNNoteCodeFromTryOn(cmd.TryOn),
		TryOn:           cmd.TryOn,
		TotalWeight:     cmd.TotalWeight,
		FeeLines:        cmd.FeeLines,
		TotalItems:      cmd.TotalItems,
		PartnerID:       cmd.PartnerID,
	}

	if err := order.BeforeUpdate(); err != nil {
		return err
	}

	// only update order_lines if order's fulfillment does not exist
	if len(cmd.Lines) > 0 {
		var ffm = new(shipmodel.Fulfillment)
		has, _ := x.Table("fulfillment").Where("order_id = ? AND status != ?", cmd.ID, model.S5Zero).Get(ffm)
		if has {
			return cm.Error(cm.FailedPrecondition, "Đơn giao hàng đã được tạo. Không thể cập nhật đơn hàng này.", nil)
		}
	}

	m := M{}
	return inTransaction(func(x Qx) error {
		if len(cmd.Lines) > 0 {
			// delete old lines + insert new lines
			if _, err := x.Table("order_line").Where("order_id = ?", cmd.ID).Delete(&ordermodel.OrderLine{}); err != nil {
				return err
			}
			fn := gencode.GenerateLineCode(oldOrder.Code, len(cmd.Lines))
			for i, line := range cmd.Lines {
				line.OrderID = cmd.ID
				line.Code = fn(i)
				if err := x.Table("order_line").
					ShouldInsert(line); err != nil {
					return err
				}
			}
			order.Lines = cmd.Lines
		}

		// TODO: Handle status
		s2 := x.Table("order").
			Where("id = ? AND shop_id = ?", order.ID, order.ShopID).
			Where("status = 0 OR status = 2 OR status IS NULL") // Only update orders in 'processing'
		if err := s2.ShouldUpdate(order); err != nil {
			return err
		}

		// TODO: Handler pointer in common/sql
		if cmd.ShopCOD != nil {
			m["shop_cod"] = *cmd.ShopCOD
		}
		if cmd.ShopShippingFee != nil {
			m["shop_shipping_fee"] = *cmd.ShopShippingFee
		}
		if cmd.OrderDiscount != nil {
			m["order_discount"] = *cmd.OrderDiscount
		}
		if cmd.TotalFee != nil {
			m["total_fee"] = *cmd.TotalFee
		}
		if len(m) == 0 {
			return nil
		}
		// require update basket_value, total_amount, total_discount at the same time because of constraint
		if cmd.BasketValue != 0 {
			m["basket_value"] = cmd.BasketValue
		}
		if cmd.TotalAmount != 0 {
			m["total_amount"] = cmd.TotalAmount
		}
		if cmd.TotalDiscount != 0 {
			m["total_discount"] = cmd.TotalDiscount
		}
		if _, err := x.Table("order").
			Where("id = ? AND shop_id = ?", order.ID, order.ShopID).
			Where("status = 0 OR status = 2 OR status IS NULL").
			UpdateMap(m); err != nil {
			return err
		}
		return nil
	})
}

func GetFulfillment(ctx context.Context, query *shipmodelx.GetFulfillmentQuery) error {
	if query.FulfillmentID == 0 && query.ShippingCode == "" && query.ExternalShippingCode == "" {
		return cm.Error(cm.InvalidArgument, "You must provide fulfillment's id or code", nil)
	}

	s := x.Table("fulfillment")
	if query.FulfillmentID != 0 {
		s = s.Where("id = ?", query.FulfillmentID)
	}
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}
	if query.ShippingProvider != "" {
		s = s.Where("shipping_provider = ?", query.ShippingProvider)
	}
	switch {
	case query.ShippingCode != "":
		s = s.Where("shipping_code = ?", query.ShippingCode).
			OrderBy("created_at DESC")
		// shipping_code may be duplicated (for example, partners reuse old codes)
	case query.ExternalShippingCode != "":
		s = s.Where("external_shipping_code = ?", query.ExternalShippingCode).
			OrderBy("created_at DESC")
	}

	query.Result = new(shipmodel.Fulfillment)
	err := s.ShouldGet(query.Result)
	return err
}

func GetFulfillmentExtended(ctx context.Context, cmd *shipmodelx.GetFulfillmentExtendedQuery) error {
	if cmd.FulfillmentID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing FulfillmentID", nil)
	}
	s := x.Table("fulfillment")
	if cmd.ShopID != 0 {
		s = s.Where("f.shop_id = ?", cmd.ShopID)
	}
	if cmd.PartnerID != 0 {
		s = s.Where("f.partner_id = ?", cmd.PartnerID)
	}
	if cmd.ExternalShippingCode != "" {
		s = s.Where("f.external_shipping_code = ?", cmd.ExternalShippingCode)
	}
	ffm := new(shipmodely.FulfillmentExtended)
	err := s.
		Where("f.id = ?", cmd.FulfillmentID).
		ShouldGet(ffm)
	cmd.Result = ffm
	return err
}

func GetFulfillments(ctx context.Context, query *shipmodelx.GetFulfillmentsQuery) error {
	isLimitSort := true
	s := x.Table("fulfillment")
	// ignore failed ffm (missing shipping_code)
	s = s.Where("shipping_code is not null")

	if query.ShopIDs != nil {
		s = s.InOrEqIDs("shop_id", query.ShopIDs)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}
	if query.OrderID != 0 {
		s = s.Where("order_id = ?", query.OrderID)
	}
	if query.Status != nil {
		s = s.Where("status = ?", query.Status)
	}
	if len(query.ShippingCodes) > 0 {
		s = s.In("shipping_code", query.ShippingCodes)
		isLimitSort = false
	}
	if len(query.ExternalShippingCodes) > 0 {
		s = s.In("external_shipping_code", query.ExternalShippingCodes)
		isLimitSort = false
	}
	if query.IDs != nil {
		s = s.In("id", query.IDs)
		isLimitSort = false
	}

	s, _, err := Filters(s, query.Filters, filterFulfillmentWhitelist)
	if err != nil {
		return err
	}
	{
		s2 := s.Clone()
		if isLimitSort {
			s2, err = LimitSort(s2, query.Paging, Ms{"updated_at": "", "created_at": "", "id": ""})
			if err != nil {
				return err
			}
		}
		if err := s2.Find((*shipmodel.Fulfillments)(&query.Result.Fulfillments)); err != nil {
			return err
		}
	}
	if len(query.Filters) == 0 {
		total, err := s.Count(&shipmodel.Fulfillment{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetUnCompleteFulfillments(ctx context.Context, query *shipmodelx.GetUnCompleteFulfillmentsQuery) error {
	s := x.Table("fulfillment").Where("status = 2 AND shipping_status not in (1, -2, -1)").OrderBy("created_at DESC")
	if len(query.ShippingProviders) != 0 {
		s = s.In("shipping_provider ", query.ShippingProviders)
	}
	var fulfillments []*shipmodel.Fulfillment
	if err := s.Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return err
	}
	query.Result = fulfillments
	return nil
}

func GetFulfillmentsCallbackLogs(ctx context.Context, query *shipmodelx.GetFulfillmentsCallbackLogs) error {
	s := x.Table("fulfillment")
	if query.FromID != 0 {
		s.Where("id > ?", query.FromID)
	}
	if len(query.ExcludeShippingStates) > 0 {
		s = s.NotIn("shipping_state", query.ExcludeShippingStates)
	}
	s, err := LimitSort(s, query.Paging, Ms{"updated_at": "", "created_at": "", "id": ""})
	if err != nil {
		return err
	}
	if err := s.Find((*shipmodel.Fulfillments)(&query.Result.Fulfillments)); err != nil {
		return err
	}
	return nil
}

func GetFulfillmentExtendeds(ctx context.Context, query *shipmodelx.GetFulfillmentExtendedsQuery) error {
	s := x.Table("fulfillment")
	// ignore failed ffm (missing shipping_code)
	s = s.Where("f.shipping_code is not null")

	if query.ShopIDs != nil {
		s = s.InOrEqIDs("f.shop_id", query.ShopIDs)
	}
	if query.PartnerID != 0 {
		s = s.Where("f.partner_id = ?", query.PartnerID)
	}
	if query.OrderID != 0 {
		s = s.Where("f.order_id = ?", query.OrderID)
	}
	if query.Status != nil {
		s = s.Where("f.status = ?", query.Status)
	}
	if len(query.ShippingCodes) > 0 {
		s = s.In("f.shipping_code", query.ShippingCodes)
	}
	if query.DateFrom.IsZero() != query.DateTo.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !query.DateFrom.IsZero() {
		s = s.Where("f.created_at BETWEEN ? AND ?", query.DateFrom, query.DateTo)
	}

	s, _, err := Filters(s, query.Filters, filterFulfillmentWhitelist)
	if err != nil {
		return err
	}

	// for exporting data
	if query.ResultAsRows {
		{
			s2 := s.Clone()
			total, err := s2.Count((*shipmodely.FulfillmentExtendeds)(nil))
			if err != nil {
				return err
			}
			query.Result.Total = int(total)
		}
		{
			if query.Paging != nil && len(query.Paging.Sort) != 0 {
				s = s.OrderBy(query.Paging.Sort...)
			} else {
				s = s.OrderBy("f.created_at")
			}

			opts, rows, err := s.FindRows((*shipmodely.FulfillmentExtendeds)(nil))
			if err != nil {
				return err
			}
			query.Result.Opts = opts
			query.Result.Rows = rows
		}
		return nil
	}

	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"updated_at": "f.updated_at", "created_at": "f.created_at", "id": "f.id"})
		if err != nil {
			return err
		}
		if err := s2.Find((*shipmodely.FulfillmentExtendeds)(&query.Result.Fulfillments)); err != nil {
			return err
		}
	}
	if len(query.Filters) == 0 {
		total, err := s.Count(&shipmodely.FulfillmentExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func CreateFulfillments(ctx context.Context, cmd *shipmodelx.CreateFulfillmentsCommand) error {
	for _, ffm := range cmd.Fulfillments {
		if ffm.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing FulfillmentID", nil)
		}
	}

	return inTransaction(func(x Qx) error {
		for _, ffm := range cmd.Fulfillments {
			if err := ffm.BeforeInsert(); err != nil {
				return err
			}
			if _, err := x.Insert(ffm); err != nil {
				return err
			}
		}
		return nil
	})
}

func UpdateFulfillment(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentCommand) error {
	if cmd.Fulfillment.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if err := cmd.Fulfillment.BeforeUpdate(); err != nil {
		return err
	}
	s := x.Table("fulfillment").
		Where("id = ?", cmd.Fulfillment.ID)
	if cmd.Fulfillment.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.Fulfillment.PartnerID)
	}
	m := M{}
	if cmd.ExternalShippingNote != nil {
		m["external_shipping_note"] = *cmd.ExternalShippingNote
	}
	if cmd.ExternalShippingSubState != nil {
		m["external_shipping_sub_state"] = *cmd.ExternalShippingSubState
	}
	if err := s.ShouldUpdate(cmd.Fulfillment); err != nil {
		return err
	}
	if len(m) > 0 {
		if err := x.Table("fulfillment").Where("id = ?", cmd.Fulfillment.ID).ShouldUpdateMap(m); err != nil {
			return err
		}
	}
	return nil
}

func UpdateFulfillments(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsCommand) error {
	for _, ffm := range cmd.Fulfillments {
		if ffm.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing ID", nil)
		}
	}

	return inTransaction(func(s Qx) error {
		for _, ffm := range cmd.Fulfillments {
			if err := ffm.BeforeUpdate(); err != nil {
				return err
			}
			if err := s.Table("fulfillment").
				Where("id = ?", ffm.ID).
				Where("status = 0 OR status = 2 OR status IS NULL").
				ShouldUpdate(ffm); err != nil {
				return err
			}
		}
		return nil
	})
}

func UpdateFulfillmentsWithoutTransaction(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsWithoutTransactionCommand) error {
	maxGoroutines := 8
	chUpdate := make(chan error, maxGoroutines)
	guard := make(chan int, maxGoroutines)

	for i, ffm := range cmd.Fulfillments {
		guard <- i
		go ignoreError(func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				<-guard
				chUpdate <- _err
			}()
			if err := ffm.BeforeUpdate(); err != nil {
				return err
			}
			updated, err := x.Table("fulfillment").Where("id = ?", ffm.ID).Where("status = 0 OR status = 2 OR status IS NULL").Update(ffm)
			if err != nil {
				return err
			}
			if updated == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
			return nil
		}(ffm))
	}

	var updated, errors int
	for i, n := 0, len(cmd.Fulfillments); i < n; i++ {
		err := <-chUpdate
		if err == nil {
			updated++
		} else {
			errors++
		}
	}
	ll.S.Infof("update fulfillment :: updated %v/%v, errors %v/%v",
		updated, len(cmd.Fulfillments),
		errors, len(cmd.Fulfillments))
	cmd.Result.Updated = updated
	cmd.Result.Error = errors
	return nil
}

func UpdateFulfillmentsStatus(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsStatusCommand) error {
	if len(cmd.FulfillmentIDs) == 0 || cmd.FulfillmentIDs[0] == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}

	m := map[string]interface{}{}
	if cmd.Status != nil {
		m["status"] = *cmd.Status
	}
	if cmd.ShopConfirm != nil {
		m["shop_confirm"] = *cmd.ShopConfirm
	}
	if cmd.SyncStatus != nil {
		m["sync_status"] = *cmd.SyncStatus
	}
	if cmd.ShippingState != "" {
		m["shipping_state"] = cmd.ShippingState
	}
	return x.Table("fulfillment").
		InOrEqIDs("id", cmd.FulfillmentIDs).
		ShouldUpdateMap(m)
}

func SyncUpdateFulfillments(ctx context.Context, cmd *shipmodelx.SyncUpdateFulfillmentsCommand) error {
	if cmd.ShippingSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShippingSourceID", nil)
	}
	if cmd.LastSyncAt.IsZero() {
		return cm.Error(cm.InvalidArgument, "Missing LastSyncAt", nil)
	}
	maxGoroutines := 8
	chUpdate := make(chan error, maxGoroutines)
	guard := make(chan int, maxGoroutines)
	for i, ffm := range cmd.Fulfillments {
		guard <- i
		go ignoreError(func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				<-guard
				chUpdate <- _err
			}()
			if err := ffm.BeforeUpdate(); err != nil {
				return err
			}
			updated, err := x.Table("fulfillment").Where("id = ?", ffm.ID).Update(ffm)
			if err != nil {
				return err
			}
			if updated == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
			return nil
		}(ffm))
	}

	var errs cm.ErrorCollector
	var updated, errors int
	for i, l := 0, len(cmd.Fulfillments); i < l; i++ {
		err := <-chUpdate
		if err == nil {
			updated++
		} else {
			errors++
		}
		errs.Collect(err)
	}
	ll.S.Infof("Sync update fulfillment to db: updated %v/%v, errors %v/%v",
		updated, len(cmd.Fulfillments),
		errors, len(cmd.Fulfillments))
	if errors > 0 {
		ll.Error("Error", l.Error(errs.Any()))
		return errs.Any()
	}
	updateCmd := &model.UpdateOrCreateShippingSourceInternal{
		ID:         cmd.ShippingSourceID,
		LastSyncAt: cmd.LastSyncAt,
	}
	err := bus.Dispatch(ctx, updateCmd)
	return err
}

func UpdateFulfillmentsShippingState(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsShippingStateCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Fulfillment IDs", nil)
	}

	var ffms []*shipmodel.Fulfillment
	s := x.Table("fulfillment").Where("shop_id = ?", cmd.ShopID)
	if cmd.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.PartnerID)
	}
	if err := s.In("id", cmd.IDs).
		Find((*shipmodel.Fulfillments)(&ffms)); err != nil {
		return err
	}
	ffmsMap := make(map[int64]*shipmodel.Fulfillment)
	for _, ffm := range ffms {
		ffmsMap[ffm.ID] = ffm
	}
	for _, id := range cmd.IDs {
		ffm := ffmsMap[id]
		if ffm == nil {
			return cm.Errorf(cm.NotFound, nil, "Không tìm thấy đơn giao hàng").WithMetap("id", id)
		}
		var order = new(ordermodel.Order)
		if err := x.Table("order").Where("id = ?", ffm.OrderID).ShouldGet(order); err != nil {
			return cm.Errorf(cm.NotFound, nil, "Không tìm thấy đơn hàng của đơn giao hàng").WithMetap("id", id)
		}

		switch order.Status {
		case model.S5Negative:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("ffm ID", id)
		case model.S5Positive:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hoàn thành").WithMetap("fulfillment_id", id)
		case model.S5NegSuper:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("fulfillment_id", id)
		}
		if order.ConfirmStatus == model.S3Negative ||
			order.ShopConfirm == model.S3Negative {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("ffm ID", id)
		}
		switch ffm.Status {
		case model.S5Negative:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã huỷ").WithMetap("fulfillment_id", id)
		case model.S5Positive:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hoàn thành").WithMetap("fulfillment_id", id)
		case model.S5NegSuper:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("fulfillment_id", id)
		}

		if order.ShopShipping == nil || order.ShopShipping.ShippingProvider != model.TypeShippingProviderManual {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật trạng thái đơn giao hàng này, ID = %v", id)
		}
	}
	update := map[string]interface{}{
		"shipping_state":  cmd.ShippingState,
		"shipping_status": cmd.ShippingState.ToShippingStatus5(),
		"status":          cmd.ShippingState.ToStatus4(),
	}
	if err := x.Table("fulfillment").In("id", cmd.IDs).ShouldUpdateMap(update); err != nil {
		return err
	}
	cmd.Result.Updated = len(cmd.IDs)
	return nil
}

func UpdateOrderPaymentStatus(ctx context.Context, cmd *ordermodelx.UpdateOrderPaymentStatusCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if cmd.OrderID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing OrderID", nil)
	}
	var order = new(ordermodel.Order)
	if err := x.Table("order").Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldGet(order); err != nil {
		return err
	}
	if _, err := canUpdateOrder(order); err != nil {
		return err
	}
	if order.ShopShipping == nil || order.ShopShipping.ShippingProvider != model.TypeShippingProviderManual {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật trạng thái thanh toán cho đơn hàng này")
	}

	if order.CustomerPaymentStatus == model.S3Positive {
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã được thanh toán", nil)
	}
	if err := x.Table("order").Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldUpdateMap(M{
		"customer_payment_status": cmd.Status,
	}); err != nil {
		return err
	}
	cmd.Result.Updated = 1
	return nil
}

func canUpdateOrder(order *ordermodel.Order) (bool, error) {
	if order == nil {
		return false, cm.Error(cm.FailedPrecondition, "Đơn hàng không tồn tại", nil)
	}
	switch order.Status {
	case model.S5Negative:
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("id", order.ID)
	case model.S5Positive:
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hoàn thành").WithMetap("id", order.ID)
	case model.S5NegSuper:
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("id", order.ID)
	}
	if order.ConfirmStatus == model.S3Negative ||
		order.ShopConfirm == model.S3Negative {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("id", order.ID)
	}
	return true, nil
}

func canUpdateFulfillment(ffm *shipmodel.Fulfillment) (bool, error) {
	if ffm.Status == model.S5Positive {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã hoàn thành")
	}
	if !ffm.CODEtopTransferedAt.IsZero() {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã đối soát").WithMetap("money_transaction_id", ffm.MoneyTransactionID)
	}
	if ffm.MoneyTransactionID != 0 {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã thuộc phiên chuyển tiền").WithMetap("money_transaction_id", ffm.MoneyTransactionID)
	}
	if ffm.MoneyTransactionShippingExternalID != 0 {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã thuộc phiên chuyển tiền nhà vận chuyển").WithMetap("money_transaction_shipping_external_id", ffm.MoneyTransactionShippingExternalID)
	}
	return true, nil
}

func AdminUpdateFulfillment(ctx context.Context, cmd *shipmodelx.AdminUpdateFulfillmentCommand) error {
	if cmd.FulfillmentID == 0 {
		return cm.Error(cm.InvalidArgument, "Thiếu ID đơn vận chuyển", nil)
	}
	if cmd.AdminNote == "" {
		return cm.Error(cm.InvalidArgument, "Ghi chú chỉnh sửa không được để trống", nil)
	}
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: cmd.FulfillmentID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	if ok, err := canUpdateOrder(query.Result.Order); err != nil || !ok {
		return err
	}
	ffm := query.Result.Fulfillment
	if ok, err := canUpdateFulfillment(ffm); err != nil || !ok {
		return err
	}
	updateFfm := &shipmodel.Fulfillment{
		ID:        ffm.ID,
		AdminNote: cmd.AdminNote,
	}

	if cmd.ActualCompensationAmount != 0 {
		if ffm.ShippingState != model.StateUndeliverable && cmd.ShippingState != model.StateUndeliverable {
			return cm.Error(cm.FailedPrecondition, "Chỉ cập nhật ActualCompensationAmount khi đơn vận chuyển không giao được hàng.", nil)
		} else {
			updateFfm.ActualCompensationAmount = cmd.ActualCompensationAmount
		}
	}
	if cmd.ShippingState != "" {
		if ffm.ShippingState != model.StateUndeliverable && cmd.ShippingState != model.StateUndeliverable {
			return cm.Error(cm.PermissionDenied, "Chỉ được cập nhật sang trạng thái không giao được hàng", nil)
		}
		updateFfm.ShippingState = cmd.ShippingState
	}

	order := query.Result.Order
	updateOrder := &ordermodel.Order{
		ID: ffm.OrderID,
	}
	updateFfmMap := M{}
	updateOrderMap := M{}
	if cmd.TotalCODAmount != nil {
		updateFfmMap["total_cod_amount"] = cmd.TotalCODAmount
		if cmd.IsPartialDelivery {
			updateFfmMap["is_partial_delivery"] = true
		} else {
			updateOrderMap["shop_cod"] = cmd.TotalCODAmount
		}
	}

	updateOrder.Customer = order.Customer.UpdateCustomer(cmd.FullName)
	updateOrder.CustomerAddress = order.CustomerAddress.UpdateAddress(cmd.Phone, cmd.FullName)
	updateOrder.BillingAddress = order.BillingAddress.UpdateAddress(cmd.Phone, cmd.FullName)
	updateOrder.ShippingAddress = order.ShippingAddress.UpdateAddress(cmd.Phone, cmd.FullName)
	if cmd.FullName != "" {
		updateOrder.CustomerName = cmd.FullName
	}
	if cmd.Phone != "" {
		updateOrder.CustomerPhone = cmd.Phone
	}
	updateFfm.AddressTo = ffm.AddressTo.UpdateAddress(cmd.Phone, cmd.FullName)

	return inTransaction(func(s Qx) error {
		if err := s.Table("order").Where("id = ?", ffm.OrderID).
			Where("status = 0 OR status = 2 OR status IS NULL").ShouldUpdate(updateOrder); err != nil {
			return err
		}
		if err := s.Table("fulfillment").Where("id = ?", ffm.ID).
			Where("status = 0 OR status = 2 OR status IS NULL").ShouldUpdate(updateFfm); err != nil {
			return err
		}
		if len(updateFfmMap) > 0 {
			if _, err := s.Table("fulfillment").Where("id = ?", ffm.ID).
				Where("status = 0 OR status = 2 OR status IS NULL").UpdateMap(updateFfmMap); err != nil {
				return err
			}
		}
		if len(updateOrderMap) > 0 {
			if _, err := s.Table("order").Where("id = ?", ffm.OrderID).
				Where("status = 0 OR status = 2 OR status IS NULL").UpdateMap(updateOrderMap); err != nil {
				return err
			}
		}
		cmd.Result.Updated = 1
		return nil
	})
}

func GenerateVtpostShippingCode() (string, error) {
	var code int
	if err := x.SQL(`SELECT nextval('shipping_code')`).Scan(&code); err != nil {
		return "", err
	}
	// checksum: avoid input wrong code
	checksumDigit := gencode.CheckSumDigitUPC(strconv.Itoa(code))
	return checksumDigit, nil
}
