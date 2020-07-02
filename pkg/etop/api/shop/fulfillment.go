package shop

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/top/int/shop"
	inttypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
)

type FulfillmentService struct {
	session.Session

	ShippingQuery shipping.QueryBus
	ShippingCtrl  *shipping_provider.CarrierManager
}

func (s *FulfillmentService) Clone() shop.FulfillmentService { res := *s; return &res }

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *pbcm.IDRequest) (*inttypes.Fulfillment, error) {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		ShopID:        s.SS.Shop().ID,
		PartnerID:     s.SS.CtxPartner().GetID(),
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbFulfillment(query.Result.Fulfillment, model.TagShop, query.Result.Shop, query.Result.Order)
	return result, nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *shop.GetFulfillmentsRequest) (*inttypes.FulfillmentsResponse, error) {
	shopIDs, err := api.MixAccount(s.SS.Claim(), q.Mixed)
	if err != nil {
		return nil, err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs:   shopIDs,
		PartnerID: s.SS.CtxPartner().GetID(),
		OrderID:   q.OrderId,
		Status:    q.Status,
		Paging:    paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &inttypes.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagShop),
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
		Fulfillments: convertpb.Convert_core_Fulfillments_To_api_Fulfillments(query.Result, model.TagShop),
	}
	return result, nil
}

func (s *FulfillmentService) GetExternalShippingServices(ctx context.Context, q *inttypes.GetExternalShippingServicesRequest) (*inttypes.GetExternalShippingServicesResponse, error) {
	resp, err := s.ShippingCtrl.GetExternalShippingServices(ctx, s.SS.Shop().ID, q)
	result := &inttypes.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, err
}

func (s *FulfillmentService) GetPublicExternalShippingServices(ctx context.Context, q *inttypes.GetExternalShippingServicesRequest) (*inttypes.GetExternalShippingServicesResponse, error) {
	resp, err := s.ShippingCtrl.GetExternalShippingServices(ctx, model.EtopAccountID, q)
	result := &inttypes.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, err
}

func (s *FulfillmentService) GetPublicFulfillment(ctx context.Context, q *shop.GetPublicFulfillmentRequest) (*inttypes.PublicFulfillment, error) {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: q.Code,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
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
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}
