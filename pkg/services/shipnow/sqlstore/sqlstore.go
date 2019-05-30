package sqlstore

import (
	"context"

	"etop.vn/api/main/etop"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	metav1 "etop.vn/api/meta/v1"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
	etopconvert "etop.vn/backend/pkg/services/etop/convert"
	"etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
	shipnowmodelx "etop.vn/backend/pkg/services/shipnow/modelx"
)

type ShipnowStoreFactory func(context.Context) *ShipnowStore

func NewShipnowStore(db cmsql.Database) ShipnowStoreFactory {
	return func(ctx context.Context) *ShipnowStore {
		return &ShipnowStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShipnowStore struct {
	query func() cmsql.QueryInterface
}

func (s *ShipnowStore) GetByID(args shipnowmodelx.GetByIDArgs) (*shipnow.ShipnowFulfillment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	q := s.query().Where("id = ?", args.ID)
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

func (s *ShipnowStore) GetShipnowFulfillments(args *shipnowmodelx.GetShipnowFulfillmentsArgs) (results []*shipnowmodel.ShipnowFulfillment, err error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	err = s.query().Where("shop_id = ?", args.ShopID).Find((*shipnowmodel.ShipnowFulfillments)(&results))
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
	err := s.query().ShouldInsert(modelShipnowFfm)
	return err
}

type UpdateInfoArgs struct {
	ID                  int64
	PickupAddress       *ordertypes.Address
	Carrier             string
	ShippingServiceCode string
	ShippingServiceFee  int32
	ShippingNote        string
	RequestPickupAt     *metav1.Timestamp
	DeliveryPoints      []*shipnow.DeliveryPoint
	WeightInfo          shippingtypes.WeightInfo
	ValueInfo           shippingtypes.ValueInfo
}

func (s *ShipnowStore) UpdateInfo(args UpdateInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	update := &shipnow.ShipnowFulfillment{
		Id:                  args.ID,
		PickupAddress:       args.PickupAddress,
		DeliveryPoints:      args.DeliveryPoints,
		Carrier:             args.Carrier,
		ShippingServiceCode: args.ShippingServiceCode,
		ShippingServiceFee:  args.ShippingServiceFee,
		WeightInfo:          args.WeightInfo,
		ValueInfo:           args.ValueInfo,
		ShippingNote:        args.ShippingNote,
		RequestPickupAt:     args.RequestPickupAt,
	}

	modelShipnowFfm := convert.ShipnowToModel(update)
	err := s.query().Where("id = ?", args.ID).ShouldUpdate(modelShipnowFfm)

	var result shipnowmodel.ShipnowFulfillment
	err = s.query().Where("id = ?", args.ID).ShouldGet(&result)
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(&result), nil
}

type UpdateStateArgs struct {
	ID            int64
	SyncStatus    etop.Status4
	Status        etop.Status5
	State         shipnowtypes.State
	SyncStates    *model.FulfillmentSyncStates
	ConfirmStatus etop.Status3
}

func (s *ShipnowStore) UpdateSyncState(args UpdateStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &shipnowmodel.ShipnowFulfillment{
		SyncStatus:    etopconvert.Status4ToModel(args.SyncStatus),
		Status:        etopconvert.Status5ToModel(args.Status),
		ConfirmStatus: etopconvert.Status3ToModel(args.ConfirmStatus),
		ShippingState: args.State.String(),
		SyncStates:    args.SyncStates,
	}
	var ft ShipnowFulfillmentFilters
	err := s.query().Where(ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	var result shipnowmodel.ShipnowFulfillment
	err = s.query().Where(ft.ByID(args.ID)).ShouldGet(&result)
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(&result), nil
}
