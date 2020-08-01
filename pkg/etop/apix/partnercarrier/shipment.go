package partnercarrier

import (
	"bytes"
	"context"
	"encoding/json"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/top/external/partnercarrier"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/com/etc/logging/shippingwebhook"
	logmodel "o.o/backend/com/etc/logging/shippingwebhook/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type ShipmentService struct {
	session.Session

	ConnectionQuery        connectioning.QueryBus
	ShippingAggr           shipping.CommandBus
	ShippingQuery          shipping.QueryBus
	ShipmentWebhookLogAggr *shippingwebhook.Aggregate
}

func (s *ShipmentService) Clone() partnercarrier.ShipmentService {
	res := *s
	return &res
}

func (s *ShipmentService) UpdateFulfillment(ctx context.Context, r *partnercarrier.UpdateFulfillmentRequest) (_ *pbcm.UpdatedResponse, _err error) {
	if r.ShippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_code")
	}

	var ffm *shipping.Fulfillment
	defer func() {
		s.saveLogsFfmUpdate(ctx, r, ffm, _err)
	}()

	query := &connectioning.ListConnectionsQuery{
		PartnerID: s.SS.Partner().ID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	connIDs := []dot.ID{}
	for _, conn := range query.Result {
		connIDs = append(connIDs, conn.ID)
	}
	if len(connIDs) == 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể sử dụng api này. Vui lòng liên hệ %v để biết thêm chi tiết.", wl.X(ctx).CSEmail)
	}

	ffmQuery := &shippingcore.GetFulfillmentByIDOrShippingCodeQuery{
		ShippingCode:  r.ShippingCode,
		ConnectionIDs: connIDs,
	}
	if err := s.ShippingQuery.Dispatch(ctx, ffmQuery); err != nil {
		return nil, err
	}
	ffm = ffmQuery.Result

	if r.ShippingState.Valid || r.Note.Valid || r.Weight != 0 {
		cmd := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
			FulfillmentID:        ffm.ID,
			ShippingState:        r.ShippingState.Enum,
			ExternalShippingNote: r.Note,
			Weight:               r.Weight.Int(),
		}

		if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
	}

	if r.ShippingFeeLines != nil || r.CODAmount.Valid {
		// update shippingFeeLines
		cmd2 := &shippingcore.UpdateFulfillmentShippingFeesCommand{
			FulfillmentID:            ffm.ID,
			ProviderShippingFeeLines: partnercarrier.Convert_api_ShippingFeeLines_To_core_ShippingFeeLines(r.ShippingFeeLines),
			TotalCODAmount:           r.CODAmount,
		}
		if err := s.ShippingAggr.Dispatch(ctx, cmd2); err != nil {
			return nil, err
		}
	}

	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *ShipmentService) saveLogsFfmUpdate(ctx context.Context, data *partnercarrier.UpdateFulfillmentRequest, ffm *shipping.Fulfillment, err error) {
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
	}
	if err := enc.Encode(data); err == nil {
		webhookData.Data = buf.Bytes()
	}
	if ffm != nil {
		webhookData.ConnectionID = ffm.ConnectionID
	}
	if err := s.ShipmentWebhookLogAggr.CreateShippingWebhookLog(ctx, webhookData); err != nil {
		ll.Error("Partner carrier insert db webhook log error", l.Error(err))
	}
}
