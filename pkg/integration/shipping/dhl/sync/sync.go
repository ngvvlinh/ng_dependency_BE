package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"o.o/api/main/shipping"
	"o.o/api/meta"
	shippingstate "o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	"o.o/backend/com/main/shipping/convert"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/scheduler"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	etopmodel "o.o/backend/pkg/etop/model"
	shipping2 "o.o/backend/pkg/integration/shipping"
	dhlclient "o.o/backend/pkg/integration/shipping/dhl/client"
	dhldriver "o.o/backend/pkg/integration/shipping/dhl/driver"
	"o.o/backend/pkg/integration/shipping/dhl/update"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ carriertypes.ShipmentSync = &DHLSync{}

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

const (
	defaultNumWorkers            = 16
	defaultRecurrent             = 5 * time.Minute
	defaultRandomTime            = 5 * time.Minute
	defaultTimeCancelOrder       = 5 * time.Minute
	defaultRandomTimeCancelOrder = 5 * time.Minute
	defaultNumFfmsInRequest      = 5
)

type status string

const (
	pending status = "pending"
	running status = "running"
	fail    status = "fail"
	success status = "success"
)

type TaskStatus struct {
	Retry  int
	Status status
	Err    error
}

type TaskTrackingArguments struct {
	ffmIDs       []dot.ID
	shopID       dot.ID
	connectionID dot.ID
}

type TaskCancelOrderArguments struct {
	ffm *shipmodel.Fulfillment
}

type DHLSync struct {
	db                     *cmsql.Database
	schedulerTrackingOrder *scheduler.Scheduler
	schedulerCancelOrder   *scheduler.Scheduler
	shipmentManager        carrier.ShipmentManager
	shippingQS             shipping.QueryBus
	shippingAggr           shipping.CommandBus

	mapTaskTrackingArguments map[dot.ID]*TaskTrackingArguments
	ffmInProgress            map[dot.ID]*shipmodel.Fulfillment

	// key: ffmID
	mapTaskCancelOrdersStatus map[dot.ID]*TaskStatus
	// key: taskID
	mapTaskCancelOrderArgs map[dot.ID]*TaskCancelOrderArguments

	trackingOrderTicker *time.Ticker
	cancelOrderTicker   *time.Ticker
	mutex               sync.Mutex
}

func New(
	db com.MainDB, shipmentManager *carrier.ShipmentManager,
	shippingQS shipping.QueryBus, shippingAggr shipping.CommandBus,
) *DHLSync {
	sched := scheduler.New(defaultNumWorkers)
	return &DHLSync{
		schedulerTrackingOrder:    sched,
		schedulerCancelOrder:      scheduler.New(5),
		shipmentManager:           *shipmentManager,
		shippingQS:                shippingQS,
		shippingAggr:              shippingAggr,
		db:                        db,
		trackingOrderTicker:       time.NewTicker(defaultRecurrent),
		cancelOrderTicker:         time.NewTicker(defaultTimeCancelOrder),
		mapTaskTrackingArguments:  make(map[dot.ID]*TaskTrackingArguments),
		ffmInProgress:             make(map[dot.ID]*shipmodel.Fulfillment),
		mapTaskCancelOrdersStatus: make(map[dot.ID]*TaskStatus),
		mapTaskCancelOrderArgs:    make(map[dot.ID]*TaskCancelOrderArguments),
	}
}

func (d *DHLSync) listFulfillments() (ffms []*shipmodel.Fulfillment, err error) {
	fromID := dot.ID(0)

	for {
		var _ffms shipmodel.Fulfillments

		// l???y t???t c??? ffm ch??a ho??n th??nh
		// shipping_status not in (1, -1, -2)
		// ri??ng tr?????ng h???p shipping_status = -2 c???n x??? l?? tr?????ng h???p returning (ch??a ph???i tr???ng th??i cu???i)
		err = d.db.
			Where("id > ?", fromID.Int64()).
			Where("shipping_provider = ?", shipping_provider.DHL.Name()).
			Where("shipping_code IS NOT NULL").
			Where("status = ? AND (shipping_status not in (?, ?, ?) OR shipping_state = ?)", status5.S, status5.N, status5.P, status5.NS, shippingstate.Returning.String()).
			OrderBy("id asc").
			Limit(1000).
			Find(&_ffms)
		if err != nil {
			return nil, err
		}

		if len(_ffms) == 0 {
			break
		}

		fromID = _ffms[len(_ffms)-1].ID
		ffms = append(ffms, _ffms...)
	}

	return
}

func (d *DHLSync) Init(ctx context.Context) error {
	return nil
}

func (d *DHLSync) Start(ctx context.Context) error {
	d.schedulerTrackingOrder.Start()
	d.schedulerCancelOrder.Start()

	d.addTrackingTasks(ctx)
	d.addCancelOrderTasks(ctx)
	for {
		select {
		case <-d.trackingOrderTicker.C:
			d.addTrackingTasks(ctx)

		case <-d.cancelOrderTicker.C:
			d.addCancelOrderTasks(ctx)
		}
	}
}

func (d *DHLSync) Stop(ctx context.Context) error {
	d.schedulerTrackingOrder.Stop()
	d.schedulerCancelOrder.Stop()
	d.trackingOrderTicker.Stop()
	d.cancelOrderTicker.Stop()
	return nil
}

func (d *DHLSync) addTrackingTasks(ctx context.Context) error {
	// list fulfillments
	ffms, err := d.listFulfillments()
	if err != nil {
		ll.SendMessagef("func listFulfillments error %s", err.Error())
	}

	// key: shopID
	// value:
	// 	 key: connectionID
	//	 value: []*fulfillment
	mapFulfillment := make(map[dot.ID]map[dot.ID][]*shipmodel.Fulfillment)
	for _, ffm := range ffms {
		shopID := ffm.ShopID
		connectionID := ffm.ConnectionID

		if _, ok := mapFulfillment[shopID]; !ok {
			mapFulfillment[shopID] = make(map[dot.ID][]*shipmodel.Fulfillment)
		}
		mapFulfillment[shopID][connectionID] = append(mapFulfillment[shopID][connectionID], ffm)
	}

	// create task arguments
	var taskIDs []dot.ID

	d.mutex.Lock()
	for shopID, mConnectionIDAndFfms := range mapFulfillment {
		for connectionID, _ffms := range mConnectionIDAndFfms {
			var ffmIDs []dot.ID
			var ffms []*shipmodel.Fulfillment
			for _, ffm := range _ffms {
				// ignore ffmInProgress
				if _, ok := d.ffmInProgress[ffm.ID]; ok {
					continue
				}
				ffmIDs = append(ffmIDs, ffm.ID)
				ffms = append(ffms, ffm)
			}

			start := 0
			for start < len(ffmIDs) {
				end := minInt(start+defaultNumFfmsInRequest, len(ffmIDs))

				// add ffm into ffmInProgress
				for i := start; i < end; i++ {
					d.ffmInProgress[ffmIDs[i]] = ffms[i]
				}

				// add task arguments
				taskID := cm.NewID()
				d.mapTaskTrackingArguments[taskID] = &TaskTrackingArguments{
					ffmIDs:       ffmIDs[start:end],
					shopID:       shopID,
					connectionID: connectionID,
				}
				start = end
				taskIDs = append(taskIDs, taskID)
			}
		}
	}
	d.mutex.Unlock()

	// add tasks
	for _, taskID := range taskIDs {
		t := rand.Intn(int(defaultRandomTime))
		d.schedulerTrackingOrder.AddAfter(taskID, time.Duration(t), d.trackingOrder)
	}

	return nil
}

func (d *DHLSync) addCancelOrderTasks(ctx context.Context) error {
	ffms, err := d.listFfmsNeedCancel(ctx)
	if err != nil {
		return err
	}
	d.mutex.Lock()
	for _, ffm := range ffms {
		taskID := cm.NewID()
		t := rand.Intn(int(defaultRandomTimeCancelOrder))

		taskStatus, ok := d.mapTaskCancelOrdersStatus[ffm.ID]
		if ok && taskStatus.Retry < 3 && taskStatus.Status == fail {
			d.mapTaskCancelOrderArgs[taskID] = &TaskCancelOrderArguments{ffm}
			d.schedulerCancelOrder.AddAfter(taskID, time.Duration(t), d.cancelOrder)
		}
		if !ok {
			d.mapTaskCancelOrderArgs[taskID] = &TaskCancelOrderArguments{ffm}
			d.mapTaskCancelOrdersStatus[ffm.ID] = &TaskStatus{
				Retry:  0,
				Status: pending,
			}
			d.schedulerCancelOrder.AddAfter(taskID, time.Duration(t), d.cancelOrder)
		}
	}
	d.mutex.Unlock()
	return nil
}

func (d *DHLSync) trackingOrder(id interface{}, p scheduler.Planner) (err error) {
	ctx := bus.Ctx()

	taskArgumentID := id.(dot.ID)

	d.mutex.Lock()
	taskArguments := d.mapTaskTrackingArguments[taskArgumentID]
	d.mutex.Unlock()

	shopID := taskArguments.shopID
	connectionID := taskArguments.connectionID
	ffmIDs := taskArguments.ffmIDs

	mapFfmIDs := make(map[dot.ID]bool)
	{
		for _, ffmID := range ffmIDs {
			mapFfmIDs[ffmID] = true
		}
	}

	defer func() {
		// remove ffm from ffmInProgress
		d.mutex.Lock()
		for _, ffmID := range taskArguments.ffmIDs {
			delete(d.ffmInProgress, ffmID)
		}
		delete(d.mapTaskTrackingArguments, taskArgumentID)
		d.mutex.Unlock()

		if err != nil {
			sendError(shopID, connectionID, ffmIDs, err)
		}
	}()

	shipmentCarrier, err := d.shipmentManager.GetShipmentDriver(ctx, connectionID, shopID)
	if err != nil {
		return err
	}
	dhlDriver := shipmentCarrier.(*dhldriver.DHLDriver)

	// call api trackingOrder
	trackingOrders, err := dhlDriver.GetClient().TrackingOrder(ctx, &dhlclient.TrackingOrdersRequest{
		TrackItemRequest: &dhlclient.TrackItemReq{
			Bd: &dhlclient.BdTrackItemReq{
				TrackingReferenceNumber: convertIDsToStrings(ffmIDs),
			},
		},
	})
	if err != nil {
		return err
	}

	// handle result
	shipmentItems := trackingOrders.TrackItemResponse.Bd.ShipmentItems
	for _, shipmentItem := range shipmentItems {
		ffmID, err := dot.ParseID(shipmentItem.ShipmentID.String())
		if err != nil {
			return cm.Errorf(cm.Internal, err, "Can't parse shipmentID")
		}
		ffmID, _err := d.callback(ctx, shipmentItem)
		if _err != nil {
			sendError(shopID, connectionID, []dot.ID{ffmID}, _err)
		}
		delete(mapFfmIDs, ffmID)
	}

	// send error when can't found shipmentID in DHL system
	for ffmID := range mapFfmIDs {
		sendError(shopID, connectionID, []dot.ID{ffmID}, cm.Errorf(cm.ExternalServiceError, nil, "DHL: Can't found shipmentID %v", ffmID))
	}

	return nil
}

func sendError(shopID, connectionID dot.ID, ffmIDs []dot.ID, err error) {
	ll.SendMessagef("?????????\n???? Shipment-sync-service: DHL ????\n- ShopID: %v\n- FfmIDs: %v \n- Error: %v\n---", shopID, strings.Join(convertIDsToStrings(ffmIDs), ","), err.Error())
}

func (d *DHLSync) callback(
	ctx context.Context, shipmentItem *dhlclient.ShipmentItemTrackResp,
) (ffmID dot.ID, err error) {
	t0 := time.Now()

	ffmID, err = dot.ParseID(shipmentItem.ShipmentID.String())
	if err != nil {
		return 0, cm.Errorf(cm.Internal, err, "Can't parse shipmentID")
	}

	query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
		ID: ffmID,
	}
	if err := d.shippingQS.Dispatch(ctx, query); err != nil {
		return 0, cm.Errorf(cm.Internal, err, "")
	}
	ffm := query.Result
	var ffmModel shipmodel.Fulfillment
	shippingconvert.Convert_shipping_Fulfillment_To_shippingmodel_Fulfillment(ffm, &ffmModel)

	err = d.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// check to update fulfillment
		updateFfm, err := update.CalcUpdateFulfillment(&ffmModel, shipmentItem)
		if err != nil {
			return cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMeta("result", "ignore")
		}
		if updateFfm == nil {
			return nil
		}

		updateFfm.LastSyncAt = t0
		// UpdateInfo other time
		updateFfm = shipping2.CalcOtherTimeBaseOnState(updateFfm, &ffmModel, t0)

		// update shipping fee lines
		newWeight := shipmentItem.GetWeight()
		updateFeeLinesArgs := &shipping2.UpdateShippingFeeLinesArgs{
			FfmID:  ffmID,
			Weight: newWeight,
			State:  updateFfm.ShippingState,
		}
		if err := shipping2.UpdateShippingFeeLines(ctx, d.shippingAggr, updateFeeLinesArgs); err != nil {
			msg := "?????????\n???? DHL: ????n %v c?? thay ?????i c?????c ph??. Kh??ng th??? c???p nh???t. Vui l??ng ki???m tra l???i. ????\n- Weight: %v\n- State: %v\n- L???i: %v\n?????????"
			ll.SendMessage(fmt.Sprintf(msg, ffm.ShippingCode, updateFeeLinesArgs.Weight, updateFeeLinesArgs.State, err.Error()))
		}

		// update info
		update := &shipping.UpdateFulfillmentExternalShippingInfoCommand{
			FulfillmentID:             ffmID,
			ShippingState:             updateFfm.ShippingState,
			ShippingStatus:            updateFfm.ShippingStatus,
			ShippingSubstate:          updateFfm.ShippingSubstate,
			ExternalShippingData:      updateFfm.ExternalShippingData,
			ExternalShippingState:     updateFfm.ExternalShippingState,
			ExternalShippingStatus:    updateFfm.ExternalShippingStatus,
			ExternalShippingUpdatedAt: updateFfm.ExternalShippingUpdatedAt,
			ExternalShippingLogs:      convert.Convert_shippingmodel_ExternalShippingLogs_shipping_ExternalShippingLogs(updateFfm.ExternalShippingLogs),
			ExternalShippingStateCode: updateFfm.ExternalShippingStateCode,
			Weight:                    newWeight,
			ClosedAt:                  updateFfm.ClosedAt,
			LastSyncAt:                updateFfm.LastSyncAt,
			ShippingCreatedAt:         updateFfm.ShippingCreatedAt,
			ShippingPickingAt:         updateFfm.ShippingPickingAt,
			ShippingHoldingAt:         updateFfm.ShippingHoldingAt,
			ShippingDeliveringAt:      updateFfm.ShippingDeliveringAt,
			ShippingDeliveredAt:       updateFfm.ShippingDeliveredAt,
			ShippingReturningAt:       updateFfm.ShippingReturningAt,
			ShippingReturnedAt:        updateFfm.ShippingReturnedAt,
			ShippingCancelledAt:       updateFfm.ShippingCancelledAt,
		}
		if err := d.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}

		// Tr??? h??ng
		//
		// - Khi tr??? h??ng, DHL s??? sinh ra m???t trackingID m???i (m?? DHL m???i).
		// - shipment_id (ffm_id c???a TopShip) s??? kh??ng ?????i
		// - => Lu??n g???i api get tracking (`rest/v3/Tracking`) b???ng ffm_id c???a TopShip. `trackingID` m???i s??? ???????c tr??? v??? trong k???t qu??? c???a api n??y.
		// - ?????i so??t:
		//  - ?????i so??t COD s??? kh??ng c?? ????n tr??? h??ng.
		//  - ?????i so??t c?????c ph?? (ch??? l??m vi???c v???i k??? to??n) - s??? c?? 2 ????n DHL
		// => Update shipping_code (tracking_id) n???u c???n
		newShippingCode := shipmentItem.TrackingID.String()
		if newShippingCode != "" && ffm.ShippingCode != newShippingCode {
			if !shipping.IsStateReturn(update.ShippingState) {
				return nil
			}
			updateShippingCodeCmd := &shipping.UpdateFulfillmentShippingCodeCommand{
				FulfillmentID: ffmID,
				ShippingCode:  newShippingCode,
			}
			return d.shippingAggr.Dispatch(ctx, updateShippingCodeCmd)
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}

func (d *DHLSync) cancelOrder(id interface{}, p scheduler.Planner) (err error) {
	ctx := bus.Ctx()

	taskArgumentID := id.(dot.ID)

	d.mutex.Lock()

	taskArgs := d.mapTaskCancelOrderArgs[taskArgumentID]
	ffm := taskArgs.ffm

	taskStatus := d.mapTaskCancelOrdersStatus[ffm.ID]
	taskStatus.Status = running

	d.mutex.Unlock()

	defer func() {
		t0 := time.Now()
		updateFfm := &shipmodel.Fulfillment{
			LastSyncAt: t0,
			SyncStates: &shippingsharemodel.FulfillmentSyncStates{
				SyncAt:    t0,
				TrySyncAt: t0,
				Error:     etopmodel.ToError(err),
			},
		}

		if err == nil {
			updateFfm.ShippingState = shippingstate.Cancelled
			updateFfm.ShippingSubstate = shippingsubstate.WrapSubstate(shippingsubstate.Default)
		}
		_ = d.db.Where("id = ?", ffm.ID).ShouldUpdate(updateFfm)

		if err != nil {
			d.mutex.Lock()
			taskStatus.Err = err
			taskStatus.Status = fail
			taskStatus.Retry += 1
			d.mutex.Unlock()
			if taskStatus.Retry >= 3 {
				ll.SendMessage(fmt.Sprintf("DHL: Kh??ng th??? hu??? ffm (ID: %v). \nError: %v", ffm.ID, err))
			}
		} else {
			d.mutex.Lock()
			taskStatus.Err = nil
			taskStatus.Status = success
			delete(d.mapTaskCancelOrdersStatus, ffm.ID)
			d.mutex.Unlock()
		}

		d.mutex.Lock()
		delete(d.mapTaskCancelOrderArgs, taskArgumentID)
		d.mutex.Unlock()
	}()

	shipmentCarrier, err := d.shipmentManager.GetShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return err
	}
	dhlDriver := shipmentCarrier.(*dhldriver.DHLDriver)

	if err := dhlDriver.CancelFulfillment(ctx, ffm); err != nil {
		return err
	}

	d.mutex.Lock()
	taskStatus.Status = success
	d.mutex.Unlock()

	return nil
}

// because DHL's system behaviour don't allow to cancel order before 10 mins from created
// list ffms that created rather 10 mins and have shippingState = cancelled and shippingSubstate = cancelling
func (d *DHLSync) listFfmsNeedCancel(ctx context.Context) (ffms []*shipmodel.Fulfillment, err error) {
	fromID := dot.ID(0)

	for {
		var _ffms shipmodel.Fulfillments

		err = d.db.
			Where("id > ?", fromID.Int64()).
			Where("shipping_provider = ?", shipping_provider.DHL.Name()).
			Where("shipping_state = ? AND shipping_substate = ?", shippingstate.Cancelled.Name(), shippingsubstate.Cancelling.Name()).
			OrderBy("id asc").
			Limit(1000).
			Find(&_ffms)
		if err != nil {
			return nil, err
		}

		if len(_ffms) == 0 {
			break
		}

		t0 := time.Now()
		fromID = _ffms[len(_ffms)-1].ID
		for _, ffm := range _ffms {
			externalCreatedAt := ffm.ExternalShippingCreatedAt
			// when run in real environment, some ffm can't cancel after 10 mins
			// then get ffm have the diff createdAt and now >= 12 mins.
			if externalCreatedAt.Add(12 * time.Minute).After(t0) {
				continue
			}

			ffms = append(ffms, ffm)
		}
	}

	return
}
