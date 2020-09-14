package sqlstore

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"

	"o.o/api/main/location"
	"o.o/api/main/shipnow"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	ordermodel "o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	ordermodely "o.o/backend/com/main/ordering/modely"
	shipnowconvert "o.o/backend/com/main/shipnow/convert"
	shipnowmodel "o.o/backend/com/main/shipnow/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shipmodely "o.o/backend/com/main/shipping/modely"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etc/typeutil"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/capi/util"
)

var filterOrderWhitelist = sqlstore.FilterWhitelist{
	Arrays:   []string{"fulfillment.shipping_code", "fulfillment.shipping_state", "fulfillment.ids"},
	Contains: []string{"customer.name", "product.name"},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"shop.id", "code", "source", "external_code", "external_id", "customer.phone", "customer.id", "pre_order"},
	Numbers:  []string{"total_amount", "chargeable_weight"},
	Status:   []string{"status", "confirm_status", "fulfillment.shipping_status", "etop_payment_status"},
	PrefixOrRename: map[string]string{
		"fulfillment.shipping_code":   `"order".fulfillment_shipping_codes`,
		"fulfillment.shipping_state":  `"order".fulfillment_shipping_states`,
		"fulfillment.shipping_status": `"order".fulfillment_shipping_status`,

		"source":            `"order".order_source_type`,
		"shop.id":           `"order".shop_id`,
		"customer.id":       `"order".customer_id`,
		"customer.name":     `"order".customer_name_norm`,
		"customer.phone":    `"order".customer_phone`,
		"product.name":      `"order".product_name_norm`,
		"external_code":     `"order".ed_code`,
		"external_id":       `"order".external_order_id`,
		"chargeable_weight": `"order".total_weight`,
		"status":            `"order".status`,
		"fulfillment.ids":   `"order".fulfillment_ids`,
	},
}

var filterFulfillmentWhitelist = sqlstore.FilterWhitelist{
	Arrays:   nil,
	Bools:    []string{"include_insurance"},
	Contains: []string{"customer.name"},
	Dates:    []string{"created_at", "updated_at"},
	Equals: []string{
		"shipping_code", "shop.id", "carrier",
		"order.code", "order.external_code", "order.external_id",
		"shipping_state", "customer.phone", "money_transaction.id",
		"address_to.province_code", "address_to.district_code", "address_to.ward_code", "money_transaction_shipping_external_id", "money_transaction_id", "connection_id",
	},
	Numbers: []string{"total_cod_amount", "cod_amount", "shipping_fee_shop", "shipping_service_fee", "basket_value", "chargeable_weight"},
	Status:  []string{"shipping_status", "etop_payment_status"},
	PrefixOrRename: map[string]string{
		"carrier":                                "f.shipping_provider",
		"shipping_code":                          "f",
		"total_cod_amount":                       "f", // @deprecated: use cod_amount
		"cod_amount":                             "f.total_cod_amount",
		"shipping_fee_shop":                      "f", // @deprecated: use shipping_service_fee
		"shipping_service_fee":                   "f.shipping_fee_shop",
		"shop.id":                                "f.shop_id",
		"shipping_state":                         "f",
		"money_transaction.id":                   "f.money_transaction_id",
		"address_to.province_code":               "f.address_to_province_code",
		"address_to.district_code":               "f.address_to_district_code",
		"address_to.ward_code":                   "f.address_to_ward_code",
		"created_at":                             "f",
		"updated_at":                             "f",
		"basket_value":                           "f",
		"chargeable_weight":                      "f.total_weight",
		"include_insurance":                      "f",
		"etop_payment_status":                    "f",
		"money_transaction_shipping_external_id": "f",
		"money_transaction_id":                   "f",

		"customer.name":       "o.customer_name_norm",
		"customer.phone":      "o.customer_phone",
		"order.code":          "o.code",
		"order.external_code": "o.ed_code",
		"order.external_id":   "o.external_order_id",
	},
}

type OrderStoreInterface interface {
	CreateFulfillments(ctx context.Context, cmd *shipmodelx.CreateFulfillmentsCommand) error

	CreateOrder(ctx context.Context, cmd *ordermodelx.CreateOrderCommand) error

	GetFulfillment(ctx context.Context, query *shipmodelx.GetFulfillmentQuery) error

	GetFulfillmentExtended(ctx context.Context, cmd *shipmodelx.GetFulfillmentExtendedQuery) error

	GetFulfillmentExtendeds(ctx context.Context, query *shipmodelx.GetFulfillmentExtendedsQuery) error

	GetOrder(ctx context.Context, query *ordermodelx.GetOrderQuery) error

	GetOrderExtends(ctx context.Context, query *ordermodelx.GetOrderExtendedsQuery) error

	GetOrders(ctx context.Context, query *ordermodelx.GetOrdersQuery) error

	UpdateFulfillment(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentCommand) error

	UpdateFulfillments(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsCommand) error

	UpdateFulfillmentsShippingState(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsShippingStateCommand) error

	UpdateOrder(ctx context.Context, cmd *ordermodelx.UpdateOrderCommand) error

	UpdateOrderPaymentStatus(ctx context.Context, cmd *ordermodelx.UpdateOrderPaymentStatusCommand) error

	UpdateOrderShippingInfo(ctx context.Context, cmd *ordermodelx.UpdateOrderShippingInfoCommand) error

	UpdateOrdersStatus(ctx context.Context, cmd *ordermodelx.UpdateOrdersStatusCommand) error

	VerifyOrdersByEdCode(ctx context.Context, query *ordermodelx.VerifyOrdersByEdCodeQuery) error
}

type OrderStore struct {
	DB com.MainDB
	db *cmsql.Database `wire:"-"`

	LocationBus  location.QueryBus
	AccountStore AccountStoreInterface
	ShopStore    ShopStoreInterface
}

func BindOrderStore(s *OrderStore) (to OrderStoreInterface) {
	s.db = s.DB
	return s
}

func (st *OrderStore) GetOrder(ctx context.Context, query *ordermodelx.GetOrderQuery) error {
	if query.OrderID == 0 && query.ExternalID == "" && query.Code == "" {
		return cm.Error(cm.InvalidArgument, "Missing id or code", nil)
	}

	s := st.db.Table("order")
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
	if query.TradingShopID != 0 {
		s = s.Where("trading_shop_id = ?", query.TradingShopID)
	}

	order := new(ordermodel.Order)
	if err := s.ShouldGet(order); err != nil {
		return err
	}

	if query.IncludeFulfillment {
		var ffms []*ordermodelx.Fulfillment

		s := st.db.Table("fulfillment").
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
		if err := st.db.Table("shipnow_fulfillment").In("id", order.FulfillmentIDs).Find((*shipnowmodel.ShipnowFulfillments)(&shipnows)); err != nil {
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

func (st *OrderStore) GetOrders(ctx context.Context, query *ordermodelx.GetOrdersQuery) error {
	s := st.db.Table("order")
	if query.ShopIDs != nil {
		s = s.InOrEqIDs("shop_id", query.ShopIDs)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}
	if query.TradingShopID != 0 {
		s = s.Where("trading_shop_id = ?", query.TradingShopID)
	}

	s, _, err := sqlstore.Filters(s, query.Filters, filterOrderWhitelist)
	if err != nil {
		return err
	}
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}

	var orders ordermodel.Orders
	{

		s2 := s.Clone()
		s2, err := sqlstore.LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"id": "", "created_at": "", "updated_at": ""})
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

	orderIds := make([]dot.ID, len(query.Result.Orders))
	shopIdsMap := make(map[dot.ID]dot.ID)
	for i, o := range query.Result.Orders {
		orderIds[i] = o.ID
		shopIdsMap[o.ShopID] = o.ShopID
	}
	var fulfillments []*shipmodel.Fulfillment
	if err := st.db.Table("fulfillment").
		In("order_id", orderIds).OrderBy("created_at desc").
		Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return err
	}

	var shipnows []*shipnowmodel.ShipnowFulfillment
	if err := st.db.Table("shipnow_fulfillment").
		Where("status != ?", status5.N).OrderBy("created_at desc").
		Where("order_ids && ?", pq.Int64Array(util.IDsToInt64(orderIds))).
		Find((*shipnowmodel.ShipnowFulfillments)(&shipnows)); err != nil {
		return err
	}

	orderShipments := make(map[dot.ID][]*shipmodel.Fulfillment)
	for _, ffm := range fulfillments {
		orderShipments[ffm.OrderID] = append(orderShipments[ffm.OrderID], ffm)
	}
	orderShipnows := make(map[dot.ID][]*shipnow.ShipnowFulfillment)
	for _, ffm := range shipnows {
		for _, orderID := range ffm.OrderIDs {
			sn := shipnowconvert.Shipnow(ffm)
			orderShipnows[orderID] = append(orderShipnows[orderID], sn)
		}
	}

	// getShop
	shopIds := make([]dot.ID, 0, len(shopIdsMap))
	for _, shopId := range shopIdsMap {
		shopIds = append(shopIds, shopId)
	}
	shopQuery := &identitymodelx.GetShopsQuery{
		ShopIDs: shopIds,
	}
	if err := st.ShopStore.GetShops(ctx, shopQuery); err != nil {
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

func (st *OrderStore) GetOrderExtends(ctx context.Context, query *ordermodelx.GetOrderExtendedsQuery) error {
	s := st.db.Table("order")

	if query.ShopIDs != nil {
		s = s.InOrEqIDs(`"order".shop_id`, query.ShopIDs)
	}
	if query.PartnerID != 0 {
		s = s.Where(`"order".partner_id = ?`, query.PartnerID)
	}
	if query.TradingShopID != 0 {
		s = s.Where(`"order".trading_shop_id = ?`, query.TradingShopID)
	}
	if query.DateFrom.IsZero() != query.DateTo.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !query.DateFrom.IsZero() {
		s = s.Where(`"order".created_at BETWEEN ? AND ?`, query.DateFrom, query.DateTo)
	}

	if query.IDs != nil && len(query.IDs) > 0 {
		s = s.In(`"order".id`, query.IDs)
	} else {
		query, _, err := sqlstore.Filters(s, query.Filters, filterOrderWhitelist)
		if err != nil {
			return err
		}
		s = query
	}

	// for exporting data
	if query.ResultAsRows {
		{
			s2 := s.Clone()
			total, err := s2.Count(&ordermodely.OrderExtendeds{})
			if err != nil {
				return err
			}
			query.Result.Total = total
		}
		{
			if query.Paging != nil && len(query.Paging.Sort) != 0 {
				s = s.OrderBy(query.Paging.Sort...)
			} else {
				s = s.OrderBy(`"order".created_at`)
			}

			opts, rows, err := s.FindRows((*ordermodely.OrderExtendeds)(nil))
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
		s2, err := sqlstore.LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"updated_at": `"order".updated_at`, "created_at": `"order".created_at`, "id": `"order".id`})
		if err != nil {
			return err
		}
		if err := s2.Find((*ordermodely.OrderExtendeds)(&query.Result.Orders)); err != nil {
			return err
		}
	}
	if len(query.Filters) == 0 {
		total, err := s.Count(&ordermodely.OrderExtendeds{})
		if err != nil {
			return err
		}
		query.Result.Total = total
	}

	return nil
}

func (st *OrderStore) VerifyOrdersByEdCode(ctx context.Context, query *ordermodelx.VerifyOrdersByEdCodeQuery) error {
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

	s := st.db.SQL(`SELECT DISTINCT ed_code FROM "order"`).
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
	return st.db.QueryRow(sql2, args...).Scan((*pq.StringArray)(&query.Result.EdCodes))
}

func (st *OrderStore) SimpleGetOrdersByIDs(ctx context.Context, query *ordermodelx.SimpleGetOrdersByExternalIDsQuery) error {
	if query.SourceType == "" {
		return cm.Error(cm.InvalidArgument, "Missing ExternalProvider", nil)
	}

	s := st.db.Table("order").
		Where("order_source_type = ?", query.SourceType)
	if query.SourceID != 0 {
		s.Where("order_source_id = ?", query.SourceID)
	}

	s = s.In("external_order_id", query.ExternalIDs)
	return s.Find((*ordermodel.Orders)(&query.Result.Orders))
}

func (st *OrderStore) UpdateOrdersStatus(ctx context.Context, cmd *ordermodelx.UpdateOrdersStatusCommand) error {
	if cmd.ShopConfirm.Apply(0) != status3.N && cmd.CancelReason != "" {
		return cm.Error(cm.InvalidArgument, "Cancel reason provided but confirm status is not cancel", nil)
	}
	if cmd.ShopConfirm.Apply(0) == status3.N && cmd.CancelReason == "" {
		return cm.Error(cm.InvalidArgument, "Cancel orders must provide cancel reason", nil)
	}

	s := st.db.Table("order").
		// Where("status = 0 OR status = 2 OR status IS NULL"). // Only update orders in 'processing'
		InOrEqIDs("id", cmd.OrderIDs)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	if cmd.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.PartnerID)
	}

	m := M{}
	if cmd.ShopConfirm.Valid {
		m["shop_confirm"] = cmd.ShopConfirm
	}
	if cmd.ConfirmStatus.Valid {
		m["confirm_status"] = cmd.ConfirmStatus
		if cmd.ConfirmStatus.Enum == status3.P {
			m["confirmed_at"] = time.Now()
		}
	}
	if cmd.Status.Valid {
		m["status"] = cmd.Status
		if cmd.Status.Enum == status5.N {
			m["cancelled_at"] = time.Now()
		}
	}
	if cmd.PaymentStatus.Valid {
		m["payment_status"] = cmd.PaymentStatus
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
		cmd.Result.Updated = updated
	}
	return nil
}

func (st *OrderStore) CreateOrder(ctx context.Context, cmd *ordermodelx.CreateOrderCommand) error {
	order := cmd.Order
	if order.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	shop, err := st.generateShopCode(ctx, order.ShopID)
	if err != nil {
		return err
	}

	return inTransaction(st.db, func(tx Qx) error {
		order.ID = cm.NewID()
		// generate order code
		code, errCode := GenerateCode(ctx, tx, model.CodeTypeOrder, shop.Code)
		if errCode != nil {
			return errCode
		}
		order.Code = code
		if err = order.BeforeInsert(); err != nil {
			return err
		}

		if len(order.Lines) > 0 {
			fn := gencode.GenerateLineCode(order.Code, len(order.Lines))
			for i, _ := range order.Lines {
				order.Lines[i].Code = fn(i)
			}
		}
		if err = tx.Table("order").ShouldInsert(order); err != nil {
			return err
		}
		for _, line := range order.Lines {
			line.OrderID = order.ID
			if err = tx.Table("order_line").
				ShouldInsert(line); err != nil {
				return err
			}
		}
		return nil
	})
}

func (st *OrderStore) CreateOrders(ctx context.Context, cmd *ordermodelx.CreateOrdersCommand) error {
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

	shop, err := st.generateShopCode(ctx, cmd.ShopID)
	if err != nil {
		return err
	}

	errs := make([]error, len(cmd.Orders))
	for i, order := range cmd.Orders {
		errs[i] = inTransaction(st.db, func(tx Qx) error {
			order.ID = cm.NewID()
			code, errCode := GenerateCode(ctx, tx, model.CodeTypeOrder, shop.Code)
			if errCode != nil {
				return errCode
			}
			order.Code = code
			if err := order.BeforeInsert(); err != nil {
				return err
			}
			if err := tx.Table("order").ShouldInsert(order); err != nil {
				return err
			}

			fn := gencode.GenerateLineCode(order.Code, len(order.Lines))
			for i, line := range order.Lines {
				line.OrderID = order.ID
				line.Code = fn(i)
				if err := tx.Table("order_line").
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

func (st *OrderStore) generateShopCode(ctx context.Context, shopID dot.ID) (*identitymodel.Shop, error) {
	queryShop := &identitymodelx.GetShopQuery{
		ShopID: shopID,
	}
	if err := st.ShopStore.GetShop(ctx, queryShop); err != nil {
		return nil, err
	}
	shop := queryShop.Result

	// generate shop code if not existed
	if shop.Code == "" {
		// update shop Code
		var shopUpdate = &identitymodel.Shop{
			ID: shop.ID,
		}
		shopCode := gencode.GenerateShopCode()
		shopUpdate.Code = shopCode

		cmd := &identitymodelx.UpdateShopCommand{
			Shop: shopUpdate,
		}
		if err := st.AccountStore.UpdateShop(ctx, cmd); err != nil {
			return nil, err
		}
		shop.Code = shopCode
	}
	return shop, nil
}

func (st *OrderStore) UpdateOrder(ctx context.Context, cmd *ordermodelx.UpdateOrderCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	query := &ordermodelx.GetOrderQuery{
		OrderID: cmd.ID,
		ShopID:  cmd.ShopID,
	}
	if err := st.GetOrder(ctx, query); err != nil {
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
		GhnNoteCode:     typeutil.GHNNoteCodeFromTryOn(cmd.TryOn),
		TryOn:           cmd.TryOn,
		TotalWeight:     cmd.TotalWeight,
		FeeLines:        cmd.FeeLines,
		TotalItems:      cmd.TotalItems,
		PartnerID:       cmd.PartnerID,
		CustomerID:      cmd.CustomerID,
	}

	if err := order.BeforeUpdate(); err != nil {
		return err
	}

	// only update order_lines if order's fulfillment does not exist
	if len(cmd.Lines) > 0 {
		var ffm = new(shipmodel.Fulfillment)
		has, _ := st.db.Table("fulfillment").Where("order_id = ? AND status != ?", cmd.ID, status5.Z).Get(ffm)
		if has {
			return cm.Error(cm.FailedPrecondition, "Đơn giao hàng đã được tạo. Không thể cập nhật đơn hàng này.", nil)
		}
	}

	m := M{}
	return inTransaction(st.db, func(tx Qx) error {
		if len(cmd.Lines) > 0 {
			// delete old lines + insert new lines
			if _, err := tx.Table("order_line").Where("order_id = ?", cmd.ID).Delete(&ordermodel.OrderLine{}); err != nil {
				return err
			}
			fn := gencode.GenerateLineCode(oldOrder.Code, len(cmd.Lines))
			for i, line := range cmd.Lines {
				line.OrderID = cmd.ID
				line.Code = fn(i)
				if err := tx.Table("order_line").
					ShouldInsert(line); err != nil {
					return err
				}
			}
			order.Lines = cmd.Lines
		}

		// TODO: Handle status
		s2 := tx.Table("order").
			Where("id = ? AND shop_id = ?", order.ID, order.ShopID)
			// Where("status = 0 OR status = 2 OR status IS NULL") // Only update orders in 'processing'
		if err := s2.ShouldUpdate(order); err != nil {
			return err
		}

		// TODO: Handler pointer in common/sql
		if cmd.ShopCOD.Valid {
			m["shop_cod"] = cmd.ShopCOD.Int
		}
		if cmd.ShopShippingFee.Valid {
			m["shop_shipping_fee"] = cmd.ShopShippingFee.Int
		}
		if cmd.OrderDiscount.Valid {
			m["order_discount"] = cmd.OrderDiscount.Int
		}
		if cmd.TotalFee.Valid {
			m["total_fee"] = cmd.TotalFee.Int
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
		if _, err := tx.Table("order").
			Where("id = ? AND shop_id = ?", order.ID, order.ShopID).
			Where("status not in (-1, -2)").
			UpdateMap(m); err != nil {
			return err
		}
		return nil
	})
}

func (st *OrderStore) GetFulfillment(ctx context.Context, query *shipmodelx.GetFulfillmentQuery) error {
	if query.FulfillmentID == 0 && query.ShippingCode == "" && query.ExternalShippingCode == "" {
		return cm.Error(cm.InvalidArgument, "You must provide fulfillment's id or code", nil)
	}

	s := st.db.Table("fulfillment")
	if query.FulfillmentID != 0 {
		s = s.Where("id = ?", query.FulfillmentID)
	}
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.PartnerID != 0 {
		s = s.Where("partner_id = ?", query.PartnerID)
	}
	if query.ShippingProvider != 0 {
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

func (st *OrderStore) GetFulfillmentExtended(ctx context.Context, cmd *shipmodelx.GetFulfillmentExtendedQuery) error {
	if cmd.FulfillmentID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing FulfillmentID", nil)
	}
	s := st.db.Table("fulfillment")
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

func (st *OrderStore) GetFulfillments(ctx context.Context, query *shipmodelx.GetFulfillmentsQuery) error {
	isLimitSort := true
	s := st.db.Table("fulfillment")
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
	if query.Status.Valid {
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

	s, _, err := sqlstore.Filters(s, query.Filters, filterFulfillmentWhitelist)
	if err != nil {
		return err
	}
	{
		s2 := s.Clone()
		if isLimitSort {
			s2, err = sqlstore.LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"updated_at": "", "created_at": "", "id": ""})
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
		query.Result.Total = total
	}
	return nil
}

func (st *OrderStore) GetUnCompleteFulfillments(ctx context.Context, query *shipmodelx.GetUnCompleteFulfillmentsQuery) error {
	s := st.db.Table("fulfillment").Where("status = 2 AND shipping_status not in (1, -2, -1)").OrderBy("created_at DESC")
	if len(query.ShippingProviders) != 0 {
		s = s.In("shipping_provider", query.ShippingProviders)
	}
	var fulfillments []*shipmodel.Fulfillment
	if err := s.Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return err
	}
	query.Result = fulfillments
	return nil
}

func (st *OrderStore) GetFulfillmentsCallbackLogs(ctx context.Context, query *shipmodelx.GetFulfillmentsCallbackLogs) error {
	s := st.db.Table("fulfillment")
	if query.FromID != 0 {
		s.Where("id > ?", query.FromID)
	}
	if len(query.ExcludeShippingStates) > 0 {
		s = s.NotIn("shipping_state", query.ExcludeShippingStates)
	}
	s, err := sqlstore.LimitSort(s, sqlstore.ConvertPaging(query.Paging), Ms{"updated_at": "", "created_at": "", "id": ""})
	if err != nil {
		return err
	}
	if err := s.Find((*shipmodel.Fulfillments)(&query.Result.Fulfillments)); err != nil {
		return err
	}
	return nil
}

func (st *OrderStore) GetFulfillmentExtendeds(ctx context.Context, query *shipmodelx.GetFulfillmentExtendedsQuery) error {
	s := st.db.Table("fulfillment")
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
	if query.Status.Valid {
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
	if len(query.ConnectionIDs) > 0 {
		s = s.In("f.connection_id", query.ConnectionIDs)
	}

	if query.IDs != nil && len(query.IDs) > 0 {
		s = s.In("f.id", query.IDs)
	} else {
		query, _, err := sqlstore.Filters(s, query.Filters, filterFulfillmentWhitelist)
		if err != nil {
			return err
		}
		s = query
	}

	// for exporting data
	if query.ResultAsRows {
		{
			s2 := s.Clone()
			total, err := s2.Table("fulfillment", "f").Count(&shipmodely.FulfillmentExtendeds{})
			if err != nil {
				return err
			}
			query.Result.Total = total
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
		s2, err := sqlstore.LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"updated_at": "f.updated_at", "created_at": "f.created_at", "id": "f.id"})
		if err != nil {
			return err
		}
		if err := s2.Find((*shipmodely.FulfillmentExtendeds)(&query.Result.Fulfillments)); err != nil {
			return err
		}
	}
	return nil
}

func (st *OrderStore) CreateFulfillments(ctx context.Context, cmd *shipmodelx.CreateFulfillmentsCommand) error {
	for _, ffm := range cmd.Fulfillments {
		if ffm.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing FulfillmentID", nil)
		}
	}
	return inTransaction(st.db, func(tx Qx) error {
		for _, ffm := range cmd.Fulfillments {
			deliveryRoute, err := st.getDeliveryRoute(ctx, ffm)
			if err != nil {
				return err
			}
			ffm.DeliveryRoute = deliveryRoute
			if err := ffm.BeforeInsert(); err != nil {
				return err
			}
			if _, err := tx.Insert(ffm); err != nil {
				return err
			}
		}
		return nil
	})
}

func (st *OrderStore) UpdateFulfillment(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentCommand) error {
	if cmd.Fulfillment.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if err := cmd.Fulfillment.BeforeUpdate(); err != nil {
		return err
	}
	s := st.db.Table("fulfillment").
		Where("id = ?", cmd.Fulfillment.ID)
	if cmd.Fulfillment.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.Fulfillment.PartnerID)
	}
	m := M{}
	if cmd.ExternalShippingNote.Valid {
		m["external_shipping_note"] = cmd.ExternalShippingNote.Apply("")
	}
	if cmd.ExternalShippingSubState.Valid {
		m["external_shipping_sub_state"] = cmd.ExternalShippingSubState.Apply("")
	}
	if cmd.Fulfillment.AddressTo != nil && cmd.Fulfillment.AddressFrom != nil {
		deliveryRoute, err := st.getDeliveryRoute(ctx, cmd.Fulfillment)
		if err != nil {
			return err
		}
		cmd.Fulfillment.DeliveryRoute = deliveryRoute
	}
	if err := s.ShouldUpdate(cmd.Fulfillment); err != nil {
		return err
	}
	if len(m) > 0 {
		if err := st.db.Table("fulfillment").Where("id = ?", cmd.Fulfillment.ID).ShouldUpdateMap(m); err != nil {
			return err
		}
	}
	return nil
}

func (st *OrderStore) getDeliveryRoute(ctx context.Context, ffm *shipmodel.Fulfillment) (string, error) {
	deliveryRoute := model.RouteNationWide
	if ffm.AddressTo.ProvinceCode == ffm.AddressFrom.ProvinceCode {
		deliveryRoute = model.RouteSameProvince
	} else {
		queryFrom := location.GetLocationQuery{
			ProvinceCode: ffm.AddressFrom.ProvinceCode,
		}
		err := st.LocationBus.Dispatch(ctx, &queryFrom)
		if err != nil {
			return "", err
		}
		queryTo := location.GetLocationQuery{
			ProvinceCode: ffm.AddressTo.ProvinceCode,
		}
		err = st.LocationBus.Dispatch(ctx, &queryTo)
		if err != nil {
			return "", err
		}
		if queryFrom.Result.Province.Region == queryFrom.Result.Province.Region {
			deliveryRoute = model.RouteSameRegion
		}
	}

	return string(deliveryRoute), nil
}

func (st *OrderStore) UpdateFulfillments(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsCommand) error {
	for _, ffm := range cmd.Fulfillments {
		if ffm.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing ID", nil)
		}
	}

	return inTransaction(st.db, func(tx Qx) error {
		for _, ffm := range cmd.Fulfillments {
			if err := ffm.BeforeUpdate(); err != nil {
				return err
			}
			if err := tx.Table("fulfillment").
				Where("id = ?", ffm.ID).
				Where("status = 0 OR status = 2 OR status IS NULL").
				ShouldUpdate(ffm); err != nil {
				return err
			}
		}
		return nil
	})
}

func (st *OrderStore) UpdateFulfillmentsWithoutTransaction(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsWithoutTransactionCommand) error {
	maxGoroutines := 8
	chUpdate := make(chan error, maxGoroutines)
	guard := make(chan int, maxGoroutines)

	for i, ffm := range cmd.Fulfillments {
		guard <- i
		go func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				<-guard
				chUpdate <- _err
			}()
			if err := ffm.BeforeUpdate(); err != nil {
				return err
			}
			updated, err := st.db.Table("fulfillment").Where("id = ?", ffm.ID).Where("status = 0 OR status = 2 OR status IS NULL").Update(ffm)
			if err != nil {
				return err
			}
			if updated == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
			return nil
		}(ffm)
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

func (st *OrderStore) UpdateFulfillmentsStatus(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsStatusCommand) error {
	if len(cmd.FulfillmentIDs) == 0 || cmd.FulfillmentIDs[0] == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}

	m := map[string]interface{}{}
	if cmd.Status.Valid {
		m["status"] = cmd.Status
	}
	if cmd.ShopConfirm.Valid {
		m["shop_confirm"] = cmd.ShopConfirm
	}
	if cmd.SyncStatus.Valid {
		m["sync_status"] = cmd.SyncStatus
	}
	if cmd.ShippingState != "" {
		m["shipping_state"] = cmd.ShippingState
	}
	return st.db.Table("fulfillment").
		InOrEqIDs("id", cmd.FulfillmentIDs).
		ShouldUpdateMap(m)
}

func (st *OrderStore) UpdateFulfillmentsShippingState(ctx context.Context, cmd *shipmodelx.UpdateFulfillmentsShippingStateCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Fulfillment IDs", nil)
	}

	var ffms []*shipmodel.Fulfillment
	s := st.db.Table("fulfillment").Where("shop_id = ?", cmd.ShopID)
	if cmd.PartnerID != 0 {
		s = s.Where("partner_id = ?", cmd.PartnerID)
	}
	if err := s.In("id", cmd.IDs).
		Find((*shipmodel.Fulfillments)(&ffms)); err != nil {
		return err
	}
	ffmsMap := make(map[dot.ID]*shipmodel.Fulfillment)
	for _, ffm := range ffms {
		ffmsMap[ffm.ID] = ffm
	}
	for _, id := range cmd.IDs {
		ffm := ffmsMap[id]
		if ffm == nil {
			return cm.Errorf(cm.NotFound, nil, "Không tìm thấy đơn giao hàng").WithMetap("id", id)
		}
		var order = new(ordermodel.Order)
		if err := st.db.Table("order").Where("id = ?", ffm.OrderID).ShouldGet(order); err != nil {
			return cm.Errorf(cm.NotFound, nil, "Không tìm thấy đơn hàng của đơn giao hàng").WithMetap("id", id)
		}

		switch order.Status {
		case status5.N:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("ffm ID", id)
		case status5.P:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hoàn thành").WithMetap("fulfillment_id", id)
		case status5.NS:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("fulfillment_id", id)
		}
		if order.ConfirmStatus == status3.N ||
			order.ShopConfirm == status3.N {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("ffm ID", id)
		}
		switch ffm.Status {
		case status5.N:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã huỷ").WithMetap("fulfillment_id", id)
		case status5.P:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hoàn thành").WithMetap("fulfillment_id", id)
		case status5.NS:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("fulfillment_id", id)
		}

		if order.ShopShipping == nil || order.ShopShipping.ShippingProvider != shipping_provider.Manual {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật trạng thái đơn giao hàng này, ID = %v", id)
		}
	}
	update := map[string]interface{}{
		"shipping_state":  cmd.ShippingState,
		"shipping_status": cmd.ShippingState.ToStatus5(),
		"status":          cmd.ShippingState.ToStatus4(),
	}
	if err := st.db.Table("fulfillment").In("id", cmd.IDs).ShouldUpdateMap(update); err != nil {
		return err
	}
	cmd.Result.Updated = len(cmd.IDs)
	return nil
}

func (st *OrderStore) UpdateOrderPaymentStatus(ctx context.Context, cmd *ordermodelx.UpdateOrderPaymentStatusCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShopID", nil)
	}
	if cmd.OrderID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing OrderID", nil)
	}
	if !cmd.PaymentStatus.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing payment status")
	}
	var order = new(ordermodel.Order)
	if err := st.db.Table("order").Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldGet(order); err != nil {
		return err
	}
	if _, err := canUpdateOrder(order); err != nil {
		return err
	}

	update := M{
		"payment_status": cmd.PaymentStatus.Apply(status4.S),
	}
	if err := st.db.Table("order").Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldUpdateMap(update); err != nil {
		return err
	}
	cmd.Result.Updated = 1
	return nil
}

func (st *OrderStore) UpdateOrderShippingInfo(ctx context.Context, cmd *ordermodelx.UpdateOrderShippingInfoCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShopID", nil)
	}
	if cmd.OrderID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing OrderID", nil)
	}
	var order = new(ordermodel.Order)
	if err := st.db.Table("order").Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldGet(order); err != nil {
		return err
	}
	if _, err := canUpdateOrder(order); err != nil {
		return err
	}
	update := &ordermodel.Order{
		ShippingAddress: cmd.ShippingAddress,
		ShopShipping:    cmd.Shipping,
	}
	if err := st.db.Where("shop_id = ? AND id = ?", cmd.ShopID, cmd.OrderID).ShouldUpdate(update); err != nil {
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
	case status5.N:
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("id", order.ID)
	case status5.NS:
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng").WithMetap("id", order.ID)
	}
	if order.ConfirmStatus == status3.N ||
		order.ShopConfirm == status3.N {
		return false, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã huỷ").WithMetap("id", order.ID)
	}
	return true, nil
}

func canUpdateFulfillment(ffm *shipmodel.Fulfillment) (bool, error) {
	if ffm.Status == status5.P {
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
