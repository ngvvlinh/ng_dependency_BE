package shop

import (
	"context"

	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		historyService.GetFulfillmentHistory,
		accountService.GetBalanceShop,
		authorizeService.AuthorizePartner,
		authorizeService.GetAvailablePartners,
		authorizeService.GetAuthorizedPartners,
	)
}

func (s *HistoryService) GetFulfillmentHistory(ctx context.Context, r *GetFulfillmentHistoryEndpoint) error {

	filters := map[string]interface{}{
		"shop_id": r.Context.Shop.ID,
	}
	count := 0
	if r.All {
		count++
	}
	if r.OrderId != 0 {
		count++
		filters["order_id"] = r.OrderId
	}
	if r.Id != 0 {
		count++
		filters["id"] = r.Id
	}
	if count != 1 {
		return cm.Error(cm.InvalidArgument, "Must provide either all, id or order_id", nil)
	}

	paging := cmapi.CMPaging(r.Paging, "-rid")
	query := &model.GetHistoryQuery{
		Paging:  paging,
		Table:   "fulfillment",
		Filters: filters,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &etop.HistoryResponse{
		Paging: cmapi.PbPageInfo(paging, 0),
		Data:   cmapi.RawJSONObjectMsg(query.Result.Data),
	}
	return nil
}

func (s *AccountService) GetBalanceShop(ctx context.Context, q *GetBalanceShopEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &model.GetBalanceShopCommand{
		ShopID: shopID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.GetBalanceShopResponse{
		Amount: cmd.Result.Amount,
	}
	return nil
}

func (s *AuthorizeService) AuthorizePartner(ctx context.Context, q *AuthorizePartnerEndpoint) error {
	shopID := q.Context.Shop.ID
	partnerID := q.PartnerId

	queryPartner := &model.GetPartner{
		PartnerID: partnerID,
	}
	if err := bus.Dispatch(ctx, queryPartner); err != nil {
		return err
	}
	partner := queryPartner.Result.Partner
	if partner.AvailableFromEtopConfig == nil || partner.AvailableFromEtopConfig.RedirectUrl == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin partner (redirect_url). Vui lòng liên hệ admin để biết thêm chi tiết.")
	}

	relQuery := &model.GetPartnerRelationQuery{
		PartnerID: partnerID,
		AccountID: shopID,
	}
	err := bus.Dispatch(ctx, relQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// Authorize already
		return cm.Errorf(cm.AlreadyExists, nil, "Shop đã được xác thực bởi '%v'", queryPartner.Result.Partner.Name)
	case cm.NotFound:
		cmd := &model.CreatePartnerRelationCommand{
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
	query := &model.GetPartnersQuery{
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
	query := &model.GetPartnersFromRelationQuery{
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
