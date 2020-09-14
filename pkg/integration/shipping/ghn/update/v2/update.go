package v2

import (
	"strconv"
	"time"

	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ghnclient.CallbackOrder) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	now := time.Now()
	state := ghnclient.State(msg.Status)
	data, _ := jsonx.Marshal(msg)

	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     msg.Status.String(),
		ExternalShippingStatus:    state.ToStatus5(),
		ExternalShippingData:      data,
		ProviderShippingFeeLines:  msg.Fee.ToOrderFee().ToFeeLines(),
		ShippingState:             state.ToModel(),
		ShippingSubstate:          state.ToSubstateModel().Wrap(),
		ShippingStatus:            state.ToStatus5(),
		ExternalShippingLogs:      ffm.ExternalShippingLogs,
		ShippingCode:              ffm.ShippingCode,
		ExternalShippingNote:      ffm.ExternalShippingNote,
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

	// update note
	if msg.Reason.String() != "" {
		note, _ := strconv.Unquote("\"" + msg.Reason.String() + "\"")
		update.ExternalShippingNote = dot.String(note)
	}
	return update, nil
}

func CalcRefreshFulfillmentInfo(ffm *shipmodel.Fulfillment, orderGHN *ghnclient.GetOrderInfoResponse) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	state := ghnclient.State(orderGHN.Status)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     orderGHN.Status.String(),
		ExternalShippingStatus:    state.ToStatus5(),
		ShippingState:             state.ToModel(),
		ShippingSubstate:          state.ToSubstateModel().Wrap(),
		ShippingStatus:            state.ToStatus5(),
		ExternalShippingLogs:      ffm.ExternalShippingLogs,
		ShippingCode:              ffm.ShippingCode,
	}
	update.AddressTo = ffm.AddressTo.UpdateAddress(orderGHN.ToPhone.String(), orderGHN.ToName.String())
	update.TotalCODAmount = orderGHN.CODAmount.Int()
	return update, nil
}
