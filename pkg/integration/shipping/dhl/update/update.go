package update

import (
	"time"

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
