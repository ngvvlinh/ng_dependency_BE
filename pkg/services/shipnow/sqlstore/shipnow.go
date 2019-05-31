package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/meta"

	"etop.vn/api/main/shipnow/carrier"

	"etop.vn/api/main/etop"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	metav1 "etop.vn/api/meta/v1"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	etopconvert "etop.vn/backend/pkg/services/etop/convert"
	"etop.vn/backend/pkg/services/shipnow/convert"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
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
	query   func() cmsql.QueryInterface
	ft      ShipnowFulfillmentFilters
	preds   []interface{}
	filters meta.Filters
}

func (s *ShipnowStore) Filters(filters meta.Filters) *ShipnowStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShipnowStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*shipnowmodel.ShipnowFulfillment)(nil))
}

func (s *ShipnowStore) ID(id int64) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipnowStore) IDs(ids ...int64) *ShipnowStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShipnowStore) ShopID(id int64) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShipnowStore) ShippingCode(code string) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code))
	return s
}

func (s *ShipnowStore) PartnerID(id int64) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *ShipnowStore) GetShipnowDB() (*shipnowmodel.ShipnowFulfillment, error) {
	var ffm shipnowmodel.ShipnowFulfillment
	err := s.query().Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *ShipnowStore) GetShipnow() (*shipnow.ShipnowFulfillment, error) {
	shipnowFfm, err := s.GetShipnowDB()
	if err != nil {
		return nil, err
	}
	return convert.Shipnow(shipnowFfm), nil
}

func (s *ShipnowStore) ListShipnowsDB(paging *meta.Paging) ([]*shipnowmodel.ShipnowFulfillment, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, paging, SortShipnow)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShipnowWhitelist)
	if err != nil {
		return nil, err
	}

	var ffms shipnowmodel.ShipnowFulfillments
	err = query.Find(&ffms)
	return ffms, err
}

func (s *ShipnowStore) ListShipnows(paging *meta.Paging) ([]*shipnow.ShipnowFulfillment, error) {
	ffms, err := s.ListShipnowsDB(paging)
	return shipnowconvert.Shipnows(ffms), err
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
	Carrier             carrier.Carrier
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
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(modelShipnowFfm); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateStateArgs struct {
	ID             int64
	SyncStatus     etop.Status4
	Status         etop.Status5
	ShippingState  shipnowtypes.State
	SyncStates     *shipnow.SyncStates
	ConfirmStatus  etop.Status3
	ShippingStatus etop.Status5
}

func (s *ShipnowStore) UpdateSyncState(args UpdateStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &shipnowmodel.ShipnowFulfillment{
		SyncStatus:     etopconvert.Status4ToModel(args.SyncStatus),
		Status:         etopconvert.Status5ToModel(args.Status),
		ConfirmStatus:  etopconvert.Status3ToModel(args.ConfirmStatus),
		ShippingState:  shipnowtypes.StateToString(args.ShippingState),
		SyncStates:     shipnowconvert.SyncStateToModel(args.SyncStates),
		ShippingStatus: etopconvert.Status5ToModel(args.ShippingStatus),
	}
	var ft ShipnowFulfillmentFilters
	err := s.query().Where(ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateCarrierInfoArgs struct {
	ID                  int64
	FeeLines            []*shippingtypes.FeeLine
	CarrierFeeLines     []*shippingtypes.FeeLine
	TotalFee            int
	ShippingCode        string
	ShippingCreatedAt   time.Time
	ShippingState       shipnowtypes.State
	ShippingStatus      etop.Status5
	EtopPaymentStatus   etop.Status4
	CODEtopTransferedAt *metav1.Timestamp
	Status              etop.Status5

	ShippingPickingAt    *metav1.Timestamp
	ShippingDeliveringAt *metav1.Timestamp
	ShippingDeliveredAt  *metav1.Timestamp
	ShippingCancelledAt  *metav1.Timestamp
}

func (s *ShipnowStore) UpdateCarrierInfo(args UpdateCarrierInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &shipnowmodel.ShipnowFulfillment{
		ShippingState:       shipnowtypes.StateToString(args.ShippingState),
		ShippingStatus:      etopconvert.Status5ToModel(args.ShippingStatus),
		ShippingCreatedAt:   args.ShippingCreatedAt,
		ShippingCode:        args.ShippingCode,
		FeeLines:            convert.FeelinesToModel(args.FeeLines),
		CarrierFeeLines:     convert.FeelinesToModel(args.CarrierFeeLines),
		TotalFee:            args.TotalFee,
		EtopPaymentStatus:   etopconvert.Status4ToModel(args.EtopPaymentStatus),
		CODEtopTransferedAt: args.CODEtopTransferedAt.ToTime(),
		Status:              etopconvert.Status5ToModel(args.Status),

		ShippingPickingAt:    args.ShippingPickingAt.ToTime(),
		ShippingDeliveringAt: args.ShippingDeliveringAt.ToTime(),
		ShippingDeliveredAt:  args.ShippingDeliveredAt.ToTime(),
		ShippingCancelledAt:  args.ShippingCancelledAt.ToTime(),
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(updateFfm); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}
