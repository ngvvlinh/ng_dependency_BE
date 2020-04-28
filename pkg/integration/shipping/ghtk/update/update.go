package update

import (
	"strconv"
	"time"

	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func CalcUpdateFulfillmentFromWebhook(ffm *shipmodel.Fulfillment, msg *ghtkclient.CallbackOrder, ffmToUpdate *shipmodel.Fulfillment) {
	if ffmToUpdate == nil {
		ffmToUpdate = &shipmodel.Fulfillment{}
	}
	if !shipping.CanUpdateFulfillment(ffm) {
		return
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
	ffm.ShippingStatus = stateID.ToStatus5()

	return
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

	// make sure can not update ffm's shipping fee when it belong to a money transaction
	if shipping.CanUpdateFulfillmentFeelines(ffm) {
		update.ProviderShippingFeeLines = ghtk.CalcAndConvertShippingFeeLines(ghtkOrder)
		shippingFeeShopLines := shippingsharemodel.GetShippingFeeShopLines(update.ProviderShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		shippingFeeShop := 0
		for _, line := range shippingFeeShopLines {
			shippingFeeShop += line.Cost
		}
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.ShippingFeeShop = shipmodel.CalcShopShippingFee(shippingFeeShop, ffm)
	}

	return update, nil
}
