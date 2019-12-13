package sqlstore

import (
	"context"

	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/main/ordering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/capi/dot"
)

type OrderStoreFactory func(context.Context) *OrderStore

func NewOrderStore(db *cmsql.Database) OrderStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *OrderStore {
		return &OrderStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type M map[string]interface{}

type OrderStore struct {
	query cmsql.QueryFactory
	ft    OrderFilters
	preds []interface{}
}

func (s *OrderStore) ID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderStore) IDs(ids ...dot.ID) *OrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *OrderStore) ShopID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *OrderStore) PartnerID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *OrderStore) CustomerID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByCustomerID(id))
	return s
}

func (s *OrderStore) Statuses(values []status5.Status) *OrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "status", values))
	return s
}

func (s *OrderStore) CustomerIDs(ids ...dot.ID) *OrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "customer_id", ids))
	return s
}

func (s *OrderStore) ExternalID(id string) *OrderStore {
	s.preds = append(s.preds, s.ft.ByExternalOrderID(id))
	return s
}

func (s *OrderStore) Code(code string) *OrderStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *OrderStore) ExternalShopID(shopID dot.ID, externalID string) *OrderStore {
	s.preds = append(s.preds,
		s.ft.ByShopID(shopID),
		s.ft.ByExternalOrderID(externalID),
	)
	return s
}

func (s *OrderStore) ExternalPartnerID(partnerID dot.ID, externalID string) *OrderStore {
	s.preds = append(s.preds,
		s.ft.ByPartnerID(partnerID),
		s.ft.ByExternalOrderID(externalID),
	)
	return s
}

func (s *OrderStore) GetOrderDB() (*model.Order, error) {
	var order model.Order
	err := s.query().Where(s.preds...).ShouldGet(&order)
	return &order, err
}

func (s *OrderStore) GetOrder() (*ordering.Order, error) {
	order, err := s.GetOrderDB()
	if err != nil {
		return nil, err
	}
	return convert.Order(order), nil
}

func (s *OrderStore) GetOrders(args *ordering.GetOrdersArgs) (orders []*ordering.Order, err error) {
	if len(args.IDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing IDs")
	}
	x := s.query().In("id", args.IDs)
	if args.ShopID != 0 {
		x = x.Where("shop_id = ?", args.ShopID)
	}
	var results model.Orders
	err = x.Find(&results)
	return convert.Orders(results), err
}

type UpdateOrdersForReserveOrdersFfmArgs struct {
	OrderIDs   []dot.ID
	Fulfill    ordertypes.Fulfill
	FulfillIDs []dot.ID
}

func (s *OrderStore) UpdateOrdersForReserveOrdersFfm(args UpdateOrdersForReserveOrdersFfmArgs) ([]*ordering.Order, error) {
	if len(args.OrderIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	if len(args.FulfillIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs")
	}
	update := &model.Order{
		FulfillmentType: model.FulfillType(args.Fulfill),
		FulfillmentIDs:  args.FulfillIDs,
	}
	if err := s.query().In("id", args.OrderIDs).ShouldUpdate(update); err != nil {
		return nil, err
	}

	return s.GetOrders(&ordering.GetOrdersArgs{
		IDs: args.OrderIDs,
	})
}

func (s *OrderStore) ListOrdersDB() ([]*model.Order, error) {
	query := s.query().Where(s.preds)

	var orders model.Orders
	err := query.Find(&orders)
	return orders, err
}

func (s *OrderStore) ListOrders() ([]*model.Order, error) {
	orders, err := s.ListOrdersDB()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type UpdateOrdersForReleaseOrderFfmArgs struct {
	OrderIDs []dot.ID
}

func (s *OrderStore) UpdateOrdersForReleaseOrdersFfm(args UpdateOrdersForReleaseOrderFfmArgs) error {
	if len(args.OrderIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	if err := s.query().Table("order").In("id", args.OrderIDs).ShouldUpdateMap(M{
		"fulfillment_type":            nil,
		"fulfillment_ids":             nil,
		"fulfillment_shipping_states": nil,
		"fulfillment_shipping_status": status5.Z,
		"etop_payment_status":         status5.Z,
		"confirm_status":              status5.Z,
		"shop_confirm":                status5.Z,
		"status":                      status5.Z,
	}); err != nil {
		return err
	}
	return nil
}

type UpdateOrderShippingStatusArgs struct {
	ID                        dot.ID
	FulfillmentShippingStatus status5.Status
	EtopPaymentStatus         status4.Status

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int
}

func (s *OrderStore) UpdateOrderShippingStatus(args UpdateOrderShippingStatusArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Mising Order ID")
	}
	update := &model.Order{
		FulfillmentShippingStatus:  args.FulfillmentShippingStatus,
		EtopPaymentStatus:          args.EtopPaymentStatus,
		FulfillmentShippingStates:  args.FulfillmentShippingStates,
		FulfillmentPaymentStatuses: args.FulfillmentPaymentStatuses,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(update); err != nil {
		return err
	}
	return nil
}

type UpdateOrdersConfirmStatusArgs struct {
	IDs           []dot.ID
	ShopConfirm   status3.Status
	ConfirmStatus status3.Status
}

func (s *OrderStore) UpdateOrdersConfirmStatus(args UpdateOrdersConfirmStatusArgs) error {
	if len(args.IDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	update := &model.Order{
		ShopConfirm:   args.ShopConfirm,
		ConfirmStatus: args.ConfirmStatus,
	}
	if _, err := s.query().Table("order").In("id", args.IDs).Update(update); err != nil {
		return err
	}
	return nil
}

type UpdateOrderPaymentInfoArgs struct {
	ID            dot.ID
	PaymentStatus status4.Status
	PaymentID     dot.ID
}

func (s *OrderStore) UpdateOrderPaymentInfo(args UpdateOrderPaymentInfoArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	update := &model.Order{
		PaymentStatus: args.PaymentStatus,
		PaymentID:     args.PaymentID,
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(update); err != nil {
		return err
	}
	return nil
}
