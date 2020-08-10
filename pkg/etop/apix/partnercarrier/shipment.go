package partnercarrier

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	"o.o/api/top/external/partnercarrier"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/etc/logging/shippingwebhook"
	"o.o/backend/pkg/etop/authorize/session"
	directwebhook "o.o/backend/pkg/integration/shipping/direct/webhook"
	"o.o/common/l"
)

var ll = l.New()

type ShipmentService struct {
	session.Session

	ConnectionQuery        connectioning.QueryBus
	ShippingAggr           shipping.CommandBus
	ShippingQuery          shipping.QueryBus
	ShipmentWebhookLogAggr *shippingwebhook.Aggregate
	DirectWebhook          *directwebhook.Webhook
}

func (s *ShipmentService) Clone() partnercarrier.ShipmentService {
	res := *s
	return &res
}

func (s *ShipmentService) UpdateFulfillment(ctx context.Context, r *partnercarrier.UpdateFulfillmentRequest) (_ *pbcm.UpdatedResponse, _err error) {
	if err := s.DirectWebhook.Callback(ctx, r, s.SS.Partner().ID); err != nil {
		return nil, err
	}

	return &pbcm.UpdatedResponse{Updated: 1}, nil
}
