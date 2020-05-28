package partnercarrier

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/top/external/partnercarrier"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ShipmentService struct {
	session.Session

	ConnectionQuery connectioning.QueryBus
	ShippingAggr    shipping.CommandBus
	ShippingQuery   shipping.QueryBus
}

func (s *ShipmentService) Clone() partnercarrier.ShipmentService {
	res := *s
	return &res
}

func (s *ShipmentService) UpdateFulfillment(ctx context.Context, r *partnercarrier.UpdateFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	if r.ShippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_code")
	}
	if !r.ShippingState.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_state")
	}

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
	ffm := ffmQuery.Result

	cmd := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
		FulfillmentID:        ffm.ID,
		ShippingState:        r.ShippingState.Enum,
		ExternalShippingNote: r.Note,
		Weight:               r.Weight.Int(),
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	// update shippingFeeLines
	cmd2 := &shippingcore.UpdateFulfillmentShippingFeesCommand{
		FulfillmentID:            ffm.ID,
		ProviderShippingFeeLines: partnercarrier.Convert_api_ShippingFeeLines_To_core_ShippingFeeLines(r.ShippingFeeLines),
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd2); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: cmd.Result}, nil
}
