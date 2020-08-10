package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"o.o/api/main/connectioning"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/external/partnercarrier"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/com/etc/logging/shippingwebhook"
	logmodel "o.o/backend/com/etc/logging/shippingwebhook/model"
	com "o.o/backend/com/main"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

type Webhook struct {
	db                     *cmsql.Database
	shippingAggr           shippingcore.CommandBus
	shippingQS             shippingcore.QueryBus
	connectionQS           connectioning.QueryBus
	shipmentWebhookLogAggr *shippingwebhook.Aggregate
}

func New(db com.MainDB,
	shippingQS shippingcore.QueryBus,
	shippingAggr shippingcore.CommandBus,
	connectionQS connectioning.QueryBus,
	shipmentWebhookLogAggr *shippingwebhook.Aggregate,
) *Webhook {
	return &Webhook{
		db:                     db,
		shippingQS:             shippingQS,
		shippingAggr:           shippingAggr,
		connectionQS:           connectionQS,
		shipmentWebhookLogAggr: shipmentWebhookLogAggr,
	}
}

func (wh *Webhook) Callback(ctx context.Context, args *partnercarrier.UpdateFulfillmentRequest, partnerID dot.ID) (_err error) {
	var ffm *shipmodel.Fulfillment
	var err error
	var conn *connectioning.Connection
	defer func() {
		wh.saveLogsFfmUpdate(ctx, args, ffm, _err)
	}()
	ffm, conn, err = wh.validateDataAndGetFfm(ctx, args.ShippingCode, partnerID)
	if err != nil {
		return err
	}
	data, _ := jsonx.Marshal(args)

	return wh.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		t0 := time.Now()
		updateFfm := &shipmodel.Fulfillment{
			ID:                        ffm.ID,
			ShippingState:             args.ShippingState.Enum,
			ExternalShippingData:      data,
			ExternalShippingUpdatedAt: t0,
			ExternalShippingState:     args.ShippingState.Enum.String(),
			ExternalShippingCode:      args.ShippingCode,
			LastSyncAt:                t0,
		}
		updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)

		// update shipping fee lines
		providerFeeLines := partnercarrier.Convert_api_ShippingFeeLines_To_core_ShippingFeeLines(args.ShippingFeeLines)
		updateFeeLinesCmd := &shippingcore.UpdateFulfillmentShippingFeesFromWebhookCommand{
			FulfillmentID:    ffm.ID,
			NewWeight:        cm.CoalesceInt(args.Weight.Int(), ffm.ChargeableWeight),
			NewState:         args.ShippingState.Apply(ffm.ShippingState),
			ProviderFeeLines: providerFeeLines,
		}
		if err := wh.shippingAggr.Dispatch(ctx, updateFeeLinesCmd); err != nil {
			msg := "–––\n👹 %v: đơn %v có thay đổi cước phí. Không thể cập nhật. Vui lòng kiểm tra lại. 👹\n- Weight: %v\n- State: %v\n- Lỗi: %v\n–––"
			ll.SendMessage(fmt.Sprintf(msg, conn.Name, ffm.ShippingCode, updateFeeLinesCmd.NewWeight, updateFeeLinesCmd.NewState, err.Error()))
		}

		if args.CODAmount.Valid {
			updateCODAmountArgs := &shipping.UpdateFfmCODAmountArgs{
				NewCODAmount: args.CODAmount.Int,
				Ffm:          ffm,
				CarrierName:  conn.Name,
			}
			// update COD Amount
			shipping.ValidateAndUpdateFulfillmentCOD(ctx, wh.shippingAggr, updateCODAmountArgs)
		}

		// update info
		update := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
			FulfillmentID:         ffm.ID,
			ShippingState:         args.ShippingState.Enum,
			ExternalShippingData:  updateFfm.ExternalShippingData,
			ExternalShippingState: updateFfm.ExternalShippingState,
			ExternalShippingNote:  args.Note,
			Weight:                args.Weight.Int(),
			ClosedAt:              updateFfm.ClosedAt,
			LastSyncAt:            updateFfm.LastSyncAt,
			ShippingCreatedAt:     updateFfm.ShippingCreatedAt,
			ShippingPickingAt:     updateFfm.ShippingPickingAt,
			ShippingHoldingAt:     updateFfm.ShippingHoldingAt,
			ShippingDeliveringAt:  updateFfm.ShippingDeliveringAt,
			ShippingDeliveredAt:   updateFfm.ShippingDeliveredAt,
			ShippingReturningAt:   updateFfm.ShippingReturningAt,
			ShippingReturnedAt:    updateFfm.ShippingReturnedAt,
			ShippingCancelledAt:   updateFfm.ShippingCancelledAt,
		}
		if err := wh.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}
		return nil
	})
}

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, shippingCode string, partnerID dot.ID) (*shipmodel.Fulfillment, *connectioning.Connection, error) {
	if shippingCode == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_code")
	}
	query := &connectioning.ListConnectionsQuery{
		PartnerID: partnerID,
	}
	if err := wh.connectionQS.Dispatch(ctx, query); err != nil {
		return nil, nil, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy đối tác vận chuyển (partner_id = %v)", partnerID)
	}
	connIDs := []dot.ID{}
	var mapConnections = make(map[dot.ID]*connectioning.Connection)
	for _, conn := range query.Result {
		connIDs = append(connIDs, conn.ID)
		mapConnections[conn.ID] = conn
	}
	if len(connIDs) == 0 {
		return nil, nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể sử dụng api này. Vui lòng liên hệ %v để biết thêm chi tiết.", wl.X(ctx).CSEmail)
	}

	ffmQuery := &shippingcore.GetFulfillmentByIDOrShippingCodeQuery{
		ShippingCode:  shippingCode,
		ConnectionIDs: connIDs,
	}
	if err := wh.shippingQS.Dispatch(ctx, ffmQuery); err != nil {
		return nil, nil, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy ffm (shipping_code = %v)", shippingCode)
	}
	ffm := ffmQuery.Result
	res := shippingconvert.Convert_shipping_Fulfillment_shippingmodel_Fulfillment(ffm, nil)
	return res, mapConnections[ffm.ConnectionID], nil
}

func (wh *Webhook) saveLogsFfmUpdate(ctx context.Context, data *partnercarrier.UpdateFulfillmentRequest, ffm *shipmodel.Fulfillment, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	webhookData := &logmodel.ShippingProviderWebhook{
		ID:               cm.NewID(),
		ShippingProvider: shipping_provider.Partner.String(),
		ShippingCode:     data.ShippingCode,
		Error:            model.ToError(err),
	}
	if data.ShippingState.Valid {
		webhookData.ShippingState = data.ShippingState.Enum.String()
		webhookData.ExternalShippingState = data.ShippingState.Enum.String()
	}
	if err := enc.Encode(data); err == nil {
		webhookData.Data = buf.Bytes()
	}
	if ffm != nil {
		webhookData.ConnectionID = ffm.ConnectionID
	}
	if err := wh.shipmentWebhookLogAggr.CreateShippingWebhookLog(ctx, webhookData); err != nil {
		ll.Error("Partner carrier insert db webhook log error", l.Error(err))
	}
}
