package sqlstore

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
)

type ShipnowStore struct {
	db cmsql.Database
}

func NewShipnowStore(db cmsql.Database) *ShipnowStore {
	return &ShipnowStore{db: db}
}

func (s *ShipnowStore) WithContext(ctx context.Context) *ShipnowStoreWithContext {
	return &ShipnowStoreWithContext{db: s.db, ctx: ctx}
}

type ShipnowStoreWithContext struct {
	ctx context.Context
	db  cmsql.Database
}

func (s *ShipnowStoreWithContext) query() cmsql.Query {
	return s.db.WithContext(s.ctx)
}

func (s *ShipnowStoreWithContext) GetByID(args shipnowmodel.GetByIDArgs) (*shipnow.ShipnowFulfillment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	q := s.db.WithContext(s.ctx).Where("id = ?", args.ID)
	if args.ShopID != 0 {
		q = q.Where("shop_id = ?", args.ShopID)
	}
	if args.PartnerID != 0 {
		q = q.Where("partner_id = ?", args.PartnerID)
	}

	result := &shipnowmodel.ShipnowFulfillment{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.Shipnow(result), nil
}

func (s *ShipnowStoreWithContext) GetShipnowFulfillments(args *shipnowmodel.GetShipnowFulfillmentsArgs) (results []*shipnowmodel.ShipnowFulfillment, err error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	err = s.db.WithContext(s.ctx).Where("shop_id = ?", args.ShopID).Find((*shipnowmodel.ShipnowFulfillments)(&results))
	return
}

func (s *ShipnowStoreWithContext) Create(shipnowFfm *shipnow.ShipnowFulfillment) error {
	modelShipnowFfm := convert.ShipnowToModel(shipnowFfm)
	if err := modelShipnowFfm.Validate(); err != nil {
		return err
	}
	if modelShipnowFfm.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	err := s.db.WithContext(s.ctx).ShouldInsert(modelShipnowFfm)
	return err
}

func (s *ShipnowStoreWithContext) Update(shipnowFfm *shipnow.ShipnowFulfillment) (*shipnow.ShipnowFulfillment, error) {
	if shipnowFfm.Id == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	modelShipnowFfm := convert.ShipnowToModel(shipnowFfm)
	err := s.db.WithContext(s.ctx).Where("id = ?", shipnowFfm.Id).ShouldUpdate(modelShipnowFfm)

	var result shipnowmodel.ShipnowFulfillment
	err = s.db.Where("id = ?", shipnowFfm.Id).ShouldGet(&result)
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(&result), nil
}

type UpdateSyncStateArgs struct {
	ID         int64
	SyncStatus etop.Status4
	State      shipnowtypes.State
	SyncStates *model.FulfillmentSyncStates
}

func (s *ShipnowStoreWithContext) UpdateSyncState(args UpdateSyncStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &shipnowmodel.ShipnowFulfillment{
		SyncStatus:    model.Status4(args.SyncStatus),
		ShippingState: args.State.String(),
		SyncStates:    args.SyncStates,
	}
	var ft ShipnowFulfillmentFilters
	err := s.query().Where(ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	var result shipnowmodel.ShipnowFulfillment
	err = s.query().Where(ft.ByID(updateFfm.ID)).ShouldGet(&result)
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(&result), nil
}

func (s *ShipnowStoreWithContext) GetActiveShipnowFulfillmentsByOrderID(args *shipnowmodel.GetActiveShipnowFulfillmentsByOrderIDArgs) ([]*shipnow.ShipnowFulfillment, error) {
	var ffms []*shipnowmodel.ShipnowFulfillment
	x := s.db.WithContext(s.ctx).Where("? = ANY(order_ids) AND status in (?, ?)", args.OrderID, model.S5Zero, model.S5SuperPos)
	if args.ExcludeShipnowFulfillmentID != 0 {
		x = x.Where("id != ?", args.ExcludeShipnowFulfillmentID)
	}
	if err := x.Find((*shipnowmodel.ShipnowFulfillments)(&ffms)); err != nil {
		return nil, err
	}

	return convert.Shipnows(ffms), nil
}
