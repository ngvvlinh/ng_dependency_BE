package sqlstore

import (
	"context"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/main/shipnow/convert"
	"o.o/backend/com/main/shipnow/model"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShipnowStoreFactory func(context.Context) *ShipnowStore

func NewShipnowStore(db *cmsql.Database) ShipnowStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShipnowStore {
		return &ShipnowStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShipnowStore struct {
	query   cmsql.QueryFactory
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

func (s *ShipnowStore) ID(id dot.ID) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipnowStore) IDs(ids ...dot.ID) *ShipnowStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShipnowStore) ShopID(id dot.ID) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShipnowStore) ShopIDs(ids ...dot.ID) *ShipnowStore {
	s.preds = append(s.preds, sq.In("shop_id", ids))
	return s
}

func (s *ShipnowStore) ShippingCode(code string) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code))
	return s
}

func (s *ShipnowStore) OptionalShippingCode(code string) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code).Optional())
	return s
}

func (s *ShipnowStore) OptionalID(id dot.ID) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByID(id).Optional())
	return s
}

func (s *ShipnowStore) OptionalExternalID(externalID string) *ShipnowStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID).Optional(), sq.NotIn("status", -1))
	return s
}

func (s *ShipnowStore) PartnerID(id dot.ID) *ShipnowStore {
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

func (s *ShipnowStore) ListShipnowsDB(paging *sqlstore.Paging) ([]*model.ShipnowFulfillment, error) {
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
	ffms, err := s.ListShipnowsDB(sqlstore.ConvertPaging(paging))
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
	ID                  dot.ID
	PickupAddress       *ordertypes.Address
	Carrier             carriertypes.ShipnowCarrier
	ShippingServiceCode string
	ShippingServiceFee  int
	ShippingNote        string
	RequestPickupAt     time.Time
	DeliveryPoints      []*shipnow.DeliveryPoint
	WeightInfo          shippingtypes.WeightInfo
	ValueInfo           shippingtypes.ValueInfo
	Coupon              string
}

func (s *ShipnowStore) UpdateInfo(args UpdateInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	update := &shipnow.ShipnowFulfillment{
		ID:                  args.ID,
		PickupAddress:       args.PickupAddress,
		DeliveryPoints:      args.DeliveryPoints,
		Carrier:             args.Carrier,
		ShippingServiceCode: args.ShippingServiceCode,
		ShippingServiceFee:  args.ShippingServiceFee,
		WeightInfo:          args.WeightInfo,
		ValueInfo:           args.ValueInfo,
		ShippingNote:        args.ShippingNote,
		RequestPickupAt:     args.RequestPickupAt,
		Coupon:              args.Coupon,
	}

	modelShipnowFfm := convert.ShipnowToModel(update)
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(modelShipnowFfm); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateStateArgs struct {
	ID             dot.ID
	SyncStatus     status4.Status
	Status         status5.Status
	ShippingState  shipnow_state.State
	SyncStates     *shipnow.SyncStates
	ConfirmStatus  status3.Status
	ShippingStatus status5.Status
}

func (s *ShipnowStore) UpdateSyncState(args UpdateStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &model.ShipnowFulfillment{
		SyncStatus:     args.SyncStatus,
		Status:         args.Status,
		ConfirmStatus:  args.ConfirmStatus,
		ShippingState:  args.ShippingState,
		SyncStates:     convert.SyncStateToModel(args.SyncStates),
		ShippingStatus: args.ShippingStatus,
	}
	err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateCancelArgs struct {
	ID            dot.ID
	ConfirmStatus status3.Status
	Status        status5.Status
	ShippingState shipnow_state.State
	CancelReason  string
}

func (s *ShipnowStore) UpdateCancelled(args UpdateCancelArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &model.ShipnowFulfillment{
		ConfirmStatus: args.ConfirmStatus,
		ShippingState: args.ShippingState,
		Status:        args.Status,
		CancelReason:  args.CancelReason,
	}
	err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(updateFfm)
	if err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}

type UpdateCarrierInfoArgs struct {
	ID                  dot.ID
	FeeLines            []*shippingtypes.ShippingFeeLine
	CarrierFeeLines     []*shippingtypes.ShippingFeeLine
	TotalFee            int
	ShippingCode        string
	ShippingCreatedAt   time.Time
	ShippingState       shipnow_state.State
	ShippingStatus      status5.Status
	EtopPaymentStatus   status4.Status
	CODEtopTransferedAt time.Time
	Status              status5.Status

	ShippingPickingAt          time.Time
	ShippingDeliveringAt       time.Time
	ShippingDeliveredAt        time.Time
	ShippingCancelledAt        time.Time
	ShippingServiceName        string
	ShippingServiceDescription string
	CancelReason               string
	ShippingSharedLink         string
	DeliveryPoints             []*shipnow.DeliveryPoint
	DriverPhone                string
	DriverName                 string
}

func (s *ShipnowStore) UpdateCarrierInfo(args UpdateCarrierInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	updateFfm := &model.ShipnowFulfillment{
		ShippingState:       args.ShippingState,
		ShippingStatus:      args.ShippingStatus,
		ShippingCreatedAt:   args.ShippingCreatedAt,
		ShippingCode:        args.ShippingCode,
		FeeLines:            shippingconvert.Convert_shippingtypes_ShippingFeeLines_sharemodel_ShippingFeeLines(args.FeeLines),
		CarrierFeeLines:     shippingconvert.Convert_shippingtypes_ShippingFeeLines_sharemodel_ShippingFeeLines(args.CarrierFeeLines),
		TotalFee:            args.TotalFee,
		EtopPaymentStatus:   args.EtopPaymentStatus,
		CODEtopTransferedAt: args.CODEtopTransferedAt,
		Status:              args.Status,

		ShippingPickingAt:          args.ShippingPickingAt,
		ShippingDeliveringAt:       args.ShippingDeliveringAt,
		ShippingDeliveredAt:        args.ShippingDeliveredAt,
		ShippingCancelledAt:        args.ShippingCancelledAt,
		ShippingServiceName:        args.ShippingServiceName,
		ShippingServiceDescription: args.ShippingServiceDescription,
		CancelReason:               args.CancelReason,
		ShippingSharedLink:         args.ShippingSharedLink,
		DeliveryPoints:             convert.Convert_shipnowtypes_DeliveryPoints_shipnowmodel_DeliveryPoints(args.DeliveryPoints),
		DriverName:                 args.DriverName,
		DriverPhone:                args.DriverPhone,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(updateFfm); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetShipnow()
}
