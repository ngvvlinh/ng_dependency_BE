package update

import (
	"time"

	typesshipping "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipping2 "o.o/backend/pkg/integration/shipping"
	dhlclient "o.o/backend/pkg/integration/shipping/dhl/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *dhlclient.ShipmentItemTrackResp) (*shipmodel.Fulfillment, error) {
	if !shipping2.CanUpdateFulfillment(ffm) {
		return nil, nil
	}

	now := time.Now()

	// get latestEvent (event has latest time)
	latestEvent := msg.GetLatestEvent()
	state := dhlclient.ToState(latestEvent.Status.String())
	data, _ := jsonx.Marshal(latestEvent)
	secondaryStatus := latestEvent.SecondaryStatus.String()

	// TH: khách muốn gửi hàng tại bưu cục thay vì shipper lấy hàng
	// DHL sẽ trả về 3 trạng thái:
	// 	- (1): 134 - Lấy hàng không thành công - Trạng thái của shipper
	// 	- (2): 116 - Đơn hàng đã đến điểm dịch vụ
	//	- (3): 130 - Đã lấy hàng thành công - Trạng thái của bưu cục
	// Tính huống thực tế sẽ gồm 3 trường hợp:
	// TH1: 1 -> 2 -> 3
	// TH2: 2 -> 1 -> 3
	// TH3: 3 -> 1 -> 2
	// 2 trường hợp đầu ko cần quan tâm vì sẽ tự đúng
	// với trường hợp 3 sẽ giải quyết: nếu DHL trả về mã 134 mà ffm hiện tại là 130 (hoặc holding) thì ignore
	if state == dhlclient.StateShipmentPickedUpFailed &&
		ffm.ShippingState == typesshipping.Holding {
		return nil, nil
	}

	// ignore when event have the same state and weight
	externalShippingStateCode := latestEvent.Status.String()
	if externalShippingStateCode == ffm.ExternalShippingStateCode &&
		msg.GetWeight() == ffm.ChargeableWeight {

		return nil, nil
	}

	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     dhlclient.StatusMapping[externalShippingStateCode],
		ExternalShippingStateCode: externalShippingStateCode,
		ExternalShippingSubState:  dot.String(secondaryStatus),
		ExternalShippingStatus:    state.ToStatus5(),
		ExternalShippingData:      data,
		ShippingState:             state.ToModel(),
		ShippingSubstate:          state.ToSubstateModel().Wrap(),
		ShippingStatus:            state.ToStatus5(),
		ExternalShippingLogs:      ffm.ExternalShippingLogs,
		ShippingCode:              ffm.ShippingCode,
	}
	if secondaryStatus != "" {
		update.ExternalShippingNote = dot.String(dhlclient.SecondaryStatusMapping[secondaryStatus])
	}

	// Only update status4 if the current status is not ending status
	newStatus := state.ToStatus5()

	// UpdateInfo ClosedAt
	if newStatus == status5.N || newStatus == status5.NS || newStatus == status5.P {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ExternalShippingClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update, nil
}
