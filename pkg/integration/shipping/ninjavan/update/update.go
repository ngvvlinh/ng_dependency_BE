package update

import (
	"strconv"
	"time"

	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping"
	ninjavanclient "o.o/backend/pkg/integration/shipping/ninjavan/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ninjavanclient.CallbackOrder) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	now := time.Now()
	state := ninjavanclient.State(msg.Status)
	data, _ := jsonx.Marshal(msg)

	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     msg.Status.String(),
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

	// update note
	if msg.Comments.String() != "" {
		comments := ninjavanclient.MapReasons[state][msg.Comments.String()]
		if comments != "" {
			comments = msg.Comments.String()
		}
		note, _ := strconv.Unquote("\"" + comments + "\"")
		update.ExternalShippingNote = dot.String(note)
	}
	return update, nil
}
