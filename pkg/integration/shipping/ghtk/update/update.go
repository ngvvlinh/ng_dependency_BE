package update

import (
	"strconv"
	"time"

	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func CalcUpdateFulfillmentFromWebhook(ffm *shipmodel.Fulfillment, msg *ghtkclient.CallbackOrder, ffmToUpdate *shipmodel.Fulfillment) (*shipmodel.Fulfillment, error) {
	if ffmToUpdate == nil {
		ffmToUpdate = &shipmodel.Fulfillment{}
	}
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update Fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	now := time.Now()
	data, _ := jsonx.Marshal(msg)
	var statusID int
	statusID = int(msg.StatusID)
	stateID := ghtkclient.StateID(statusID)
	ffmToUpdate.ID = ffm.ID
	ffmToUpdate.ExternalShippingUpdatedAt = now
	ffmToUpdate.ExternalShippingData = data
	ffmToUpdate.ExternalShippingState = ghtkclient.StateMapping[stateID]
	ffmToUpdate.ShippingState = stateID.ToModel()
	ffmToUpdate.ShippingSubstate = stateID.ToSubstateModel().Wrap()
	ffmToUpdate.ShippingStatus = stateID.ToStatus5()
	ffmToUpdate.ExternalShippingStateCode = strconv.Itoa(statusID)

	subState := ghtkclient.SubStateMapping[stateID]
	ffmToUpdate.ExternalShippingSubState = dot.String(subState)

	// update note
	note, _ := strconv.Unquote("\"" + msg.Reason.String() + "\"")
	ffmToUpdate.ExternalShippingNote = dot.String(note)
	return ffmToUpdate, nil
}

func CalcRefreshFulfillmentInfo(ffm *shipmodel.Fulfillment, ghtkOrder *ghtkclient.OrderInfo) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update this fulfillment").WithMeta("result", "ignore")
	}

	now := time.Now()
	statusID, _ := strconv.Atoi(ghtkOrder.Status.String())
	stateID := ghtkclient.StateID(statusID)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: now,
		ExternalShippingState:     ghtkclient.StateMapping[stateID],
		ExternalShippingStatus:    stateID.ToStatus5(),
		ExternalShippingCode:      ghtkOrder.LabelID.String(),
		ShippingState:             stateID.ToModel(),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            stateID.ToStatus5(),
	}
	update.ProviderShippingFeeLines = ghtk.CalcAndConvertShippingFeeLines(ghtkOrder)

	return update, nil
}
