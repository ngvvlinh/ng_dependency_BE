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

type ShipnowStoreFactory func(context.Context) *ShipnowStore

func NewShipnowStore(db cmsql.Database) ShipnowStoreFactory {
	return func(ctx context.Context) *ShipnowStore {
		return &ShipnowStore{query: db.WithContext(ctx)}
	}
}

type ShipnowStore struct {
	query cmsql.Query
}

func (s *ShipnowStore) GetByID(args shipnowmodel.GetByIDArgs) (*shipnow.ShipnowFulfillment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	q := s.query.Where("id = ?", args.ID)
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

func (s *ShipnowStore) GetShipnowFulfillments(args *shipnowmodel.GetShipnowFulfillmentsArgs) (results []*shipnowmodel.ShipnowFulfillment, err error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	err = s.query.Where("shop_id = ?", args.ShopID).Find((*shipnowmodel.ShipnowFulfillments)(&results))
	return
}

func (s *ShipnowStore) Create(shipnowFfm *shipnow.ShipnowFulfillment) error {
	modelShipnowFfm := convert.ShipnowToModel(shipnowFfm)
	if err := modelShipnowFfm.Validate(); err != nil {
		return err
	}
	if modelShipnowFfm.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	err := s.query.ShouldInsert(modelShipnowFfm)
	return err
}

func (s *ShipnowStore) Update(shipnowFfm *shipnow.ShipnowFulfillment) (*shipnow.ShipnowFulfillment, error) {
	if shipnowFfm.Id == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	modelShipnowFfm := convert.ShipnowToModel(shipnowFfm)
	err := s.query.Where("id = ?", shipnowFfm.Id).ShouldUpdate(modelShipnowFfm)

	var result shipnowmodel.ShipnowFulfillment
	err = s.query.Where("id = ?", shipnowFfm.Id).ShouldGet(&result)
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

func (s *ShipnowStore) UpdateSyncState(args UpdateSyncStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &shipnowmodel.ShipnowFulfillment{
		SyncStatus:    model.Status4(args.SyncStatus),
		ShippingState: args.State.String(),
		SyncStates:    args.SyncStates,
	}
	var ft ShipnowFulfillmentFilters
	err := s.query.Where(ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	var result shipnowmodel.ShipnowFulfillment
	err = s.query.Where(ft.ByID(updateFfm.ID)).ShouldGet(&result)
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(&result), nil
}

func (s *ShipnowStore) GetActiveShipnowFulfillmentsByOrderID(args *shipnowmodel.GetActiveShipnowFulfillmentsByOrderIDArgs) ([]*shipnow.ShipnowFulfillment, error) {
	var ffms []*shipnowmodel.ShipnowFulfillment
	x := s.query.Where("? = ANY(order_ids) AND status in (?, ?)", args.OrderID, model.S5Zero, model.S5SuperPos)
	if args.ExcludeShipnowFulfillmentID != 0 {
		x = x.Where("id != ?", args.ExcludeShipnowFulfillmentID)
	}
	if err := x.Find((*shipnowmodel.ShipnowFulfillments)(&ffms)); err != nil {
		return nil, err
	}

	return convert.Shipnows(ffms), nil
}
