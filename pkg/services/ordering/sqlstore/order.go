package sqlstore

import (
	"context"

	"etop.vn/api/meta"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/ordering"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/backend/pkg/common/cmsql"
	orderconvert "etop.vn/backend/pkg/services/ordering/convert"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

type OrderStore struct {
	db    cmsql.Database
	query cmsql.QueryInterface
	ctx   context.Context
	ft    OrderFilters
	preds []interface{}
}

func NewOrderStore(db cmsql.Database) *OrderStore {
	ctx := context.Background()
	return &OrderStore{
		db:    db,
		ctx:   ctx,
		query: db.WithContext(ctx),
	}
}

func (s *OrderStore) WithContext(ctx context.Context) *OrderStore {
	store := &OrderStore{
		db:    s.db,
		ctx:   ctx,
		query: s.db.WithContext(ctx),
	}
	tx := ctx.Value(meta.KeyTx{})
	if tx != nil {
		store.query = tx.(cmsql.Tx)
	}
	return store
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

func (s *OrderStore) Get() (*ordermodel.Order, error) {
	var order ordermodel.Order
	err := s.query.Where(s.preds...).ShouldGet(&order)
	return &order, err
}

func (s *OrderStore) GetOrdes(args *ordering.GetOrdersArgs) (orders []*ordering.Order, err error) {
	if len(args.IDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing IDs")
	}
	x := s.query.In("id", args.IDs)
	if args.ShopID != 0 {
		x = x.Where("shop_id = ?", args.ShopID)
	}
	var results ordermodel.Orders
	err = x.Find((*ordermodel.Orders)(&results))
	return orderconvert.Orders(results), err
}

type UpdateOrdersForReserveOrdersArgs struct {
	OrderIDs   []int64
	Fulfill    ordertypes.Fulfill
	FulfillIDs []int64
}

func (s *OrderStore) UpdateOrdersForReverseOrders(args UpdateOrdersForReserveOrdersArgs) ([]*ordering.Order, error) {
	if len(args.OrderIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OrderIDs")
	}
	if len(args.FulfillIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillIDs")
	}
	update := &ordermodel.Order{
		Fulfill:    ordermodel.FulfillType(args.Fulfill),
		FulfillIDs: args.FulfillIDs,
	}
	if err := s.query.In("id", args.OrderIDs).ShouldUpdate(update); err != nil {
		return nil, err
	}

	return s.GetOrdes(&ordering.GetOrdersArgs{
		IDs: args.OrderIDs,
	})
}
