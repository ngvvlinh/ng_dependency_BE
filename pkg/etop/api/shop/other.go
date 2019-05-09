package shop

import (
	"context"

	cmP "etop.vn/backend/pb/common"
	etopP "etop.vn/backend/pb/etop"
	shopP "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	shopW "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandlers("api",
		GetFulfillmentHistory,
		GetBalanceShop,
		AuthorizePartner,
		GetAvailablePartners,
		GetAuthorizedPartners,
	)
}

func GetFulfillmentHistory(ctx context.Context, r *shopW.GetFulfillmentHistoryEndpoint) error {

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

	paging := r.Paging.CMPaging("-rid")
	query := &model.GetHistoryQuery{
		Paging:  paging,
		Table:   "fulfillment",
		Filters: filters,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &etopP.HistoryResponse{
		Paging: cmP.PbPageInfo(paging, 0),
		Data:   cmP.RawJSONObjectMsg(query.Result.Data),
	}
	return nil
}

func GetBalanceShop(ctx context.Context, q *shopW.GetBalanceShopEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &model.GetBalanceShopCommand{
		ShopID: shopID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shopP.GetBalanceShopResponse{
		Amount: int32(cmd.Result.Amount),
	}
	return nil
}

func AuthorizePartner(ctx context.Context, q *shopW.AuthorizePartnerEndpoint) error {
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
		q.Result = shopP.PbAuthorizedPartner(partner, q.Context.Shop)
	default:
		return err
	}
	return nil
}

func GetAvailablePartners(ctx context.Context, q *shopW.GetAvailablePartnersEndpoint) error {
	query := &model.GetPartnersQuery{
		AvailableFromEtop: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.GetPartnersResponse{
		Partners: etopP.PbPublicPartners(query.Result.Partners),
	}
	return nil
}

func GetAuthorizedPartners(ctx context.Context, q *shopW.GetAuthorizedPartnersEndpoint) error {
	query := &model.GetPartnersFromRelationQuery{
		AccountIDs: []int64{q.Context.Shop.ID},
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.GetAuthorizedPartnersResponse{
		Partners: shopP.PbAuthorizedPartners(query.Result.Partners, q.Context.Shop),
	}
	return nil
}
