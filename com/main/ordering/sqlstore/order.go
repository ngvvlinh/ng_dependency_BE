package sqlstore

import (
	"context"
	"time"

	"github.com/lib/pq"

	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/main/ordering/convert"
	"o.o/backend/com/main/ordering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/capi/dot"
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

func (s *OrderStore) OptionalShopID(id dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
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

func (s *OrderStore) PaymentStatus(status status4.Status) *OrderStore {
	s.preds = append(s.preds, s.ft.ByPaymentStatus(status))
	return s
}

func (s *OrderStore) ConfirmStatus(status status3.Status) *OrderStore {
	s.preds = append(s.preds, s.ft.ByConfirmStatus(status))
	return s
}

func (s *OrderStore) CreatedBy(createdBy dot.ID) *OrderStore {
	s.preds = append(s.preds, s.ft.ByCreatedBy(createdBy))
	return s
}

func (s *OrderStore) CreatedAtFromAndTo(createdAtFrom, createdAtTo time.Time) *OrderStore {
	s.preds = append(s.preds, sq.NewExpr("created_at >= ? AND created_at < ?", createdAtFrom, createdAtTo))
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
	Fulfill    ordertypes.ShippingType
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
		FulfillmentType: args.Fulfill,
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
		"fulfillment_shipping_codes":  nil,
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
	FulfillmentStatuses        []int
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
		FulfillmentStatuses:        args.FulfillmentStatuses,
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

type UpdateOrderStatus struct {
	ID     dot.ID
	ShopID dot.ID
	Status status5.Status
}

func (s *OrderStore) UpdateOrderStatus(args *UpdateOrderStatus) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	return s.query().Table("order").Where(s.ft.ByID(args.ID)).Where(s.ft.ByShopID(args.ShopID).Optional()).ShouldUpdateMap(map[string]interface{}{
		"status": args.Status,
	})
}

func (s *OrderStore) UpdateOrderPaymentStatus(args *ordering.UpdateOrderPaymentStatusArgs) error {
	return s.query().Table("order").Where(s.ft.ByID(args.OrderID)).Where(s.ft.ByShopID(args.ShopID).Optional()).ShouldUpdateMap(map[string]interface{}{
		"payment_status": args.PaymentStatus.Apply(status4.Z),
	})
}

func (s *OrderStore) UpdateOrderCustomerInfo(args *ordering.UpdateOrderCustomerInfoArgs, oldCustomer *model.OrderCustomer) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	update := &model.Order{
		Customer: oldCustomer,
	}
	count := 0
	if args.FullName.Valid {
		update.Customer.FullName = args.FullName.Apply(oldCustomer.FullName)
		count++
	}
	if args.Phone.Valid {
		update.Customer.Phone = args.Phone.Apply(oldCustomer.Phone)
		count++
	}
	if count > 0 {
		return s.query().Table("order").Where(s.preds).ShouldUpdate(update)
	}
	return nil
}

func (s *OrderStore) UpdateFulfillmentShippingCodes(fulfillmentShippingCodes []string) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	err := s.query().Table("order").Where(s.preds).ShouldUpdateMap(map[string]interface{}{
		"fulfillment_shipping_codes": pq.StringArray(fulfillmentShippingCodes),
	})
	return err
}
