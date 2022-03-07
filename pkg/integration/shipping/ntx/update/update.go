package update

import (
	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping"
	ntxclient "o.o/backend/pkg/integration/shipping/ntx/client"
	"o.o/common/jsonx"
	"time"
)

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ntxclient.CallbackOrder) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	now := time.Now()
	state := ntxclient.State(msg.StatusName)
	data, _ := jsonx.Marshal(msg)

	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     msg.StatusName,
		ExternalShippingStatus:    state.ToStatus5(),
		ExternalShippingData:      data,
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

	return update, nil
}
