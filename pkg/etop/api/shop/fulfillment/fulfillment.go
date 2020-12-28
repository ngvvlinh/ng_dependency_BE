package fulfillment

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/top/int/shop"
	inttypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
)

type FulfillmentService struct {
	session.Session

	ShipmentManager *shippingcarrier.ShipmentManager
	ShippingQuery   shipping.QueryBus

	OrderStore sqlstore.OrderStoreInterface
}

func (s *FulfillmentService) Clone() shop.FulfillmentService { res := *s; return &res }

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *pbcm.IDRequest) (*inttypes.Fulfillment, error) {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		ShopID:        s.SS.Shop().ID,
		PartnerID:     s.SS.CtxPartner().GetID(),
		FulfillmentID: q.Id,
	}
	if err := s.OrderStore.GetFulfillmentExtended(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbFulfillment(query.Result.Fulfillment, account_tag.TagShop, query.Result.Shop, query.Result.Order)
	return result, nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *shop.GetFulfillmentsRequest) (*inttypes.FulfillmentsResponse, error) {
	shopIDs, err := apiroot.MixAccount(s.SS.Claim(), q.Mixed)
	if err != nil {
		return nil, err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs:       shopIDs,
		PartnerID:     s.SS.CtxPartner().GetID(),
		OrderID:       q.OrderId,
		Status:        q.Status,
		ConnectionIDs: q.ConnectionIDs,
		Paging:        paging,
		Filters:       cmapi.ToFilters(q.Filters),
	}
	if err := s.OrderStore.GetFulfillmentExtendeds(ctx, query); err != nil {
		return nil, err
	}
	result := &inttypes.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, account_tag.TagShop),
		Paging:       cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *FulfillmentService) GetFulfillmentsByIDs(ctx context.Context, q *shop.GetFulfillmentsByIDsRequest) (*inttypes.FulfillmentsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &shipping.ListFulfillmentsByIDsQuery{
		IDs:    q.IDs,
		ShopID: shopID,
	}
	if err := s.ShippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &inttypes.FulfillmentsResponse{
		Fulfillments: convertpb.Convert_core_Fulfillments_To_api_Fulfillments(query.Result, account_tag.TagShop),
	}
	return result, nil
}

func (s *FulfillmentService) GetExternalShippingServices(ctx context.Context, q *inttypes.GetExternalShippingServicesRequest) (*inttypes.GetExternalShippingServicesResponse, error) {
	shopID := s.SS.Shop().ID
	args, err := s.ShipmentManager.PrepareDataGetShippingServices(ctx, q.ToGetShippingServicesRequest())
	if err != nil {
		return nil, err
	}
	args.AccountID = shopID
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return nil, err
	}
	result := &inttypes.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, nil
}

func (s *FulfillmentService) GetPublicExternalShippingServices(ctx context.Context, q *inttypes.GetExternalShippingServicesRequest) (*inttypes.GetExternalShippingServicesResponse, error) {
	args, err := s.ShipmentManager.PrepareDataGetShippingServices(ctx, q.ToGetShippingServicesRequest())
	if err != nil {
		return nil, err
	}
	args.AccountID = idutil.EtopAccountID
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return nil, err
	}
	result := &inttypes.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, nil
}

func (s *FulfillmentService) GetPublicFulfillment(ctx context.Context, q *shop.GetPublicFulfillmentRequest) (*inttypes.PublicFulfillment, error) {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: q.Code,
	}
	if err := s.OrderStore.GetFulfillment(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbPublicFulfillment(query.Result)
	return result, nil
}

func (s *FulfillmentService) UpdateFulfillmentsShippingState(ctx context.Context, q *shop.UpdateFulfillmentsShippingStateRequest) (*pbcm.UpdatedResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &shipmodelx.UpdateFulfillmentsShippingStateCommand{
		ShopID:        shopID,
		IDs:           q.Ids,
		ShippingState: q.ShippingState,
	}
	if err := s.OrderStore.UpdateFulfillmentsShippingState(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}
