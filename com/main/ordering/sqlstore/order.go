package sqlstore

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	etopconvert "etop.vn/backend/com/main/etop/convert"
	"etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/main/ordering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
)

type OrderStoreFactory func(context.Context) *OrderStore

func NewOrderStore(db cmsql.Database) OrderStoreFactory {
	return func(ctx context.Context) *OrderStore {
		return &OrderStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type M map[string]interface{}

type OrderStore struct {
	query func() cmsql.QueryInterface
	ft    OrderFilters
	preds []interface{}
}

func (s *OrderStore) ID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderStore) ShopID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *OrderStore) PartnerID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
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

func (s *OrderStore) ExternalShopID(shopID int64, externalID string) *OrderStore {
	s.preds = append(s.preds,
		s.ft.ByShopID(shopID),
		s.ft.ByExternalOrderID(externalID),
	)
	return s
}

func (s *OrderStore) ExternalPartnerID(partnerID int64, externalID string) *OrderStore {
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
	err = x.Find((*model.Orders)(&results))
	return convert.Orders(results), err
}

type UpdateOrdersForReserveOrdersFfmArgs struct {
	OrderIDs   []int64
	Fulfill    ordertypes.Fulfill
	FulfillIDs []int64
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

type UpdateOrdersForReleaseOrderFfmArgs struct {
	OrderIDs []int64
}

func (s *OrderStore) UpdateOrdersForReleaseOrdersFfm(args UpdateOrdersForReleaseOrderFfmArgs) error {
	if len(args.OrderIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	if err := s.query().Table("order").In("id", args.OrderIDs).ShouldUpdateMap(M{
		"fulfillment_type":            nil,
		"fulfillment_ids":             nil,
		"fulfillment_shipping_states": nil,
		"fulfillment_shipping_status": etop.S5Zero,
		"etop_payment_status":         etop.S5Zero,
		"confirm_status":              etop.S5Zero,
		"shop_confirm":                etop.S5Zero,
		"status":                      etop.S5Zero,
	}); err != nil {
		return err
	}
	return nil
}

type UpdateOrderShippingStatusArgs struct {
	ID                        int64
	FulfillmentShippingStatus etop.Status5
	EtopPaymentStatus         etop.Status4

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int
}

func (s *OrderStore) UpdateOrderShippingStatus(args UpdateOrderShippingStatusArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Mising Order ID")
	}
	update := &model.Order{
		FulfillmentShippingStatus:  etopconvert.Status5ToModel(args.FulfillmentShippingStatus),
		EtopPaymentStatus:          etopconvert.Status4ToModel(args.EtopPaymentStatus),
		FulfillmentShippingStates:  args.FulfillmentShippingStates,
		FulfillmentPaymentStatuses: args.FulfillmentPaymentStatuses,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(update); err != nil {
		return err
	}
	return nil
}

type UpdateOrdersConfirmStatusArgs struct {
	IDs           []int64
	ShopConfirm   etop.Status3
	ConfirmStatus etop.Status3
}

func (s *OrderStore) UpdateOrdersConfirmStatus(args UpdateOrdersConfirmStatusArgs) error {
	if len(args.IDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	update := &model.Order{
		ShopConfirm:   etopconvert.Status3ToModel(args.ShopConfirm),
		ConfirmStatus: etopconvert.Status3ToModel(args.ConfirmStatus),
	}
	if _, err := s.query().Table("order").In("id", args.IDs).Update(update); err != nil {
		return err
	}
	return nil
}
