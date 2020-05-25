package shop

import (
	"context"

	"o.o/api/top/int/shop"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

type AuthorizeService struct{}

func (s *AuthorizeService) Clone() *AuthorizeService { res := *s; return &res }

func (s *AuthorizeService) AuthorizePartner(ctx context.Context, q *AuthorizePartnerEndpoint) error {
	shopID := q.Context.Shop.ID
	partnerID := q.PartnerId

	queryPartner := &identitymodelx.GetPartner{
		PartnerID: partnerID,
	}
	if err := bus.Dispatch(ctx, queryPartner); err != nil {
		return err
	}
	partner := queryPartner.Result.Partner
	if partner.AvailableFromEtopConfig == nil || partner.AvailableFromEtopConfig.RedirectUrl == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin partner (redirect_url). Vui lòng liên hệ admin để biết thêm chi tiết.")
	}

	relQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID: partnerID,
		AccountID: shopID,
	}
	err := bus.Dispatch(ctx, relQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// Authorize already
		return cm.Errorf(cm.AlreadyExists, nil, "Shop đã được xác thực bởi '%v'", queryPartner.Result.Partner.Name)
	case cm.NotFound:
		cmd := &identitymodelx.CreatePartnerRelationCommand{
			PartnerID: partnerID,
			AccountID: shopID,
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}
		q.Result = convertpb.PbAuthorizedPartner(partner, q.Context.Shop)
	default:
		return err
	}
	return nil
}

func (s *AuthorizeService) GetAvailablePartners(ctx context.Context, q *GetAvailablePartnersEndpoint) error {
	query := &identitymodelx.GetPartnersQuery{
		AvailableFromEtop: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetPartnersResponse{
		Partners: convertpb.PbPublicPartners(query.Result.Partners),
	}
	return nil
}

func (s *AuthorizeService) GetAuthorizedPartners(ctx context.Context, q *GetAuthorizedPartnersEndpoint) error {
	query := &identitymodelx.GetPartnersFromRelationQuery{
		AccountIDs: []dot.ID{q.Context.Shop.ID},
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetAuthorizedPartnersResponse{
		Partners: convertpb.PbAuthorizedPartners(query.Result.Partners, q.Context.Shop),
	}
	return nil
}
