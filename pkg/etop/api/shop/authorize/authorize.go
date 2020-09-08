package authorize

import (
	"context"

	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type AuthorizeService struct {
	session.Session

	PartnerStore sqlstore.PartnerStoreInterface
}

func (s *AuthorizeService) Clone() api.AuthorizeService { res := *s; return &res }

func (s *AuthorizeService) AuthorizePartner(ctx context.Context, q *api.AuthorizePartnerRequest) (*api.AuthorizedPartnerResponse, error) {
	shopID := s.SS.Shop().ID
	partnerID := q.PartnerId

	queryPartner := &identitymodelx.GetPartner{
		PartnerID: partnerID,
	}
	if err := s.PartnerStore.GetPartner(ctx, queryPartner); err != nil {
		return nil, err
	}
	partner := queryPartner.Result.Partner
	if partner.AvailableFromEtopConfig == nil || partner.AvailableFromEtopConfig.RedirectUrl == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin partner (redirect_url). Vui lòng liên hệ admin để biết thêm chi tiết.")
	}

	relQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID: partnerID,
		AccountID: shopID,
	}
	err := s.PartnerStore.GetPartnerRelationQuery(ctx, relQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// Authorize already
		return nil, cm.Errorf(cm.AlreadyExists, nil, "Shop đã được xác thực bởi '%v'", queryPartner.Result.Partner.Name)
	case cm.NotFound:
		cmd := &identitymodelx.CreatePartnerRelationCommand{
			PartnerID: partnerID,
			AccountID: shopID,
		}
		if err := s.PartnerStore.CreatePartnerRelation(ctx, cmd); err != nil {
			return nil, err
		}
		result := convertpb.PbAuthorizedPartner(partner, s.SS.Shop())
		return result, nil
	default:
		return nil, err
	}
}

func (s *AuthorizeService) GetAvailablePartners(ctx context.Context, q *pbcm.Empty) (*api.GetPartnersResponse, error) {
	query := &identitymodelx.GetPartnersQuery{
		AvailableFromEtop: true,
	}
	if err := s.PartnerStore.GetPartners(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetPartnersResponse{
		Partners: convertpb.PbPublicPartners(query.Result.Partners),
	}
	return result, nil
}

func (s *AuthorizeService) GetAuthorizedPartners(ctx context.Context, q *pbcm.Empty) (*api.GetAuthorizedPartnersResponse, error) {
	query := &identitymodelx.GetPartnersFromRelationQuery{
		AccountIDs: []dot.ID{s.SS.Shop().ID},
	}
	if err := s.PartnerStore.GetPartnersFromRelation(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetAuthorizedPartnersResponse{
		Partners: convertpb.PbAuthorizedPartners(query.Result.Partners, s.SS.Shop()),
	}
	return result, nil
}
