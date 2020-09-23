package vnp

import (
	"context"

	"o.o/api/top/external/mc/vnp"
	vnpentitytype "o.o/api/top/external/mc/vnp/etc/entity_type"
	"o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/session"
)

var _ vnp.WebhookService = &VNPostWebhookService{}

type VNPostWebhookService struct {
	session.Session
	WebhookInner *webhook.Service
}

func (s *VNPostWebhookService) Clone() vnp.WebhookService {
	res := *s
	return &res
}

func (s *VNPostWebhookService) CreateWebhook(ctx context.Context, r *vnp.CreateWebhookRequest) (*vnp.Webhook, error) {
	accountID := s.SS.Shop().ID
	if len(r.Entities) == 0 {
		r.Entities = append(r.Entities, vnpentitytype.ShipnowFulfillment)
	}
	entities, ok := vnpentitytype.Convert_type_VnpEntities_To_type_Entities(r.Entities)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Entities does not valid")
	}
	req := &types.CreateWebhookRequest{
		Entities: entities,
		Url:      r.URL,
	}
	whResp, err := s.WebhookInner.CreateWebhook(ctx, accountID, req)
	if err != nil {
		return nil, err
	}

	return Convert_apix_Webhook_To_vnp_Webhook(whResp), nil
}

func (s *VNPostWebhookService) GetWebhooks(ctx context.Context, r *pbcm.Empty) (*vnp.WebhooksResponse, error) {
	accountID := s.SS.Shop().ID
	whResp, err := s.WebhookInner.GetWebhooks(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return &vnp.WebhooksResponse{
		Webhooks: Convert_apix_Webhooks_To_vnp_Webhooks(whResp.Webhooks),
	}, nil
}

func (s *VNPostWebhookService) DeleteWebhook(ctx context.Context, r *types.DeleteWebhookRequest) (*vnp.WebhooksResponse, error) {
	accountID := s.SS.Shop().ID
	req := &types.DeleteWebhookRequest{Id: r.Id}

	whResp, err := s.WebhookInner.DeleteWebhook(ctx, accountID, req)
	if err != nil {
		return nil, err
	}
	return &vnp.WebhooksResponse{
		Webhooks: Convert_apix_Webhooks_To_vnp_Webhooks(whResp.Webhooks),
	}, nil
}

func (s *VNPostWebhookService) GetChanges(ctx context.Context, r *pbcm.Empty) (*vnp.DataCallback, error) {
	return nil, cm.ErrTODO
}
