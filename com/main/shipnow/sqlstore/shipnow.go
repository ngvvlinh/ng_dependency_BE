package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	carriertypes "etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	etopconvert "etop.vn/backend/com/main/etop/convert"
	"etop.vn/backend/com/main/shipnow/convert"
	"etop.vn/backend/com/main/shipnow/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
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
	filters []meta.Filter
}

func (s *ShipnowStore) Filters(filters []*meta.Filter) *ShipnowStore {
	for _, filter := range filters {
		s.filters = append(s.filters, *filter)
	}
	return s
}

func (s *ShipnowStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShipnowFulfillment)(nil))
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

func (s *ShipnowStore) ShopIDs(ids ...int64) *ShipnowStore {
	s.preds = append(s.preds, sq.In("shop_id", ids))
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

func (s *ShipnowStore) GetShipnowDB() (*model.ShipnowFulfillment, error) {
	var ffm model.ShipnowFulfillment
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

func (s *ShipnowStore) ListShipnowsDB(paging *meta.Paging) ([]*model.ShipnowFulfillment, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, paging, SortShipnow)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShipnowWhitelist)
	if err != nil {
		return nil, err
	}

	var ffms model.ShipnowFulfillments
	err = query.Find(&ffms)
	return ffms, err
}

func (s *ShipnowStore) ListShipnows(paging *meta.Paging) ([]*shipnow.ShipnowFulfillment, error) {
	ffms, err := s.ListShipnowsDB(paging)
	return convert.Shipnows(ffms), err
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
	Carrier             carriertypes.Carrier
	ShippingServiceCode string
	ShippingServiceFee  int32
	ShippingNote        string
	RequestPickupAt     time.Time
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
	updateFfm := &model.ShipnowFulfillment{
		SyncStatus:     etopconvert.Status4ToModel(args.SyncStatus),
		Status:         etopconvert.Status5ToModel(args.Status),
		ConfirmStatus:  etopconvert.Status3ToModel(args.ConfirmStatus),
		ShippingState:  shipnowtypes.StateToString(args.ShippingState),
		SyncStates:     convert.SyncStateToModel(args.SyncStates),
		ShippingStatus: etopconvert.Status5ToModel(args.ShippingStatus),
	}
	err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateCancelArgs struct {
	ID            int64
	ConfirmStatus etop.Status3
	Status        etop.Status5
	ShippingState shipnowtypes.State
	CancelReason  string
}

func (s *ShipnowStore) UpdateCancelled(args UpdateCancelArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &model.ShipnowFulfillment{
		ConfirmStatus: etopconvert.Status3ToModel(args.ConfirmStatus),
		ShippingState: shipnowtypes.StateToString(args.ShippingState),
		Status:        etopconvert.Status5ToModel(args.Status),
		CancelReason:  args.CancelReason,
	}
	err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(updateFfm)
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
	CODEtopTransferedAt time.Time
	Status              etop.Status5

	ShippingPickingAt          time.Time
	ShippingDeliveringAt       time.Time
	ShippingDeliveredAt        time.Time
	ShippingCancelledAt        time.Time
	ShippingServiceName        string
	ShippingServiceDescription string
	CancelReason               string
	ShippingSharedLink         string
}

func (s *ShipnowStore) UpdateCarrierInfo(args UpdateCarrierInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &model.ShipnowFulfillment{
		ShippingState:       shipnowtypes.StateToString(args.ShippingState),
		ShippingStatus:      etopconvert.Status5ToModel(args.ShippingStatus),
		ShippingCreatedAt:   args.ShippingCreatedAt,
		ShippingCode:        args.ShippingCode,
		FeeLines:            convert.FeelinesToModel(args.FeeLines),
		CarrierFeeLines:     convert.FeelinesToModel(args.CarrierFeeLines),
		TotalFee:            args.TotalFee,
		EtopPaymentStatus:   etopconvert.Status4ToModel(args.EtopPaymentStatus),
		CODEtopTransferedAt: args.CODEtopTransferedAt,
		Status:              etopconvert.Status5ToModel(args.Status),

		ShippingPickingAt:          args.ShippingPickingAt,
		ShippingDeliveringAt:       args.ShippingDeliveringAt,
		ShippingDeliveredAt:        args.ShippingDeliveredAt,
		ShippingCancelledAt:        args.ShippingCancelledAt,
		ShippingServiceName:        args.ShippingServiceName,
		ShippingServiceDescription: args.ShippingServiceDescription,
		CancelReason:               args.CancelReason,
		ShippingSharedLink:         args.ShippingSharedLink,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(updateFfm); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}
