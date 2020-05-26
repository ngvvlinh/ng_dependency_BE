package shop

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
)

type FulfillmentService struct {
	ShippingQuery shipping.QueryBus
	ShippingCtrl  *shipping_provider.CarrierManager
}

func (s *FulfillmentService) Clone() *FulfillmentService { res := *s; return &res }

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		ShopID:        q.Context.Shop.ID,
		PartnerID:     q.CtxPartner.GetID(),
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbFulfillment(query.Result.Fulfillment, model.TagShop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *GetFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}

	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs:   shopIDs,
		PartnerID: q.CtxPartner.GetID(),
		OrderID:   q.OrderId,
		Status:    q.Status,
		Paging:    paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagShop),
		Paging:       cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *FulfillmentService) GetFulfillmentsByIDs(ctx context.Context, q *GetFulfillmentsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &shipping.ListFulfillmentsByIDsQuery{
		IDs:    q.IDs,
		ShopID: shopID,
	}
	if err := s.ShippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.FulfillmentsResponse{
		Fulfillments: convertpb.Convert_core_Fulfillments_To_api_Fulfillments(query.Result, model.TagShop),
	}
	return nil
}

func (s *FulfillmentService) GetExternalShippingServices(ctx context.Context, q *GetExternalShippingServicesEndpoint) error {
	resp, err := s.ShippingCtrl.GetExternalShippingServices(ctx, q.Context.Shop.ID, q.GetExternalShippingServicesRequest)
	q.Result = &types.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return err
}

func (s *FulfillmentService) GetPublicExternalShippingServices(ctx context.Context, q *GetPublicExternalShippingServicesEndpoint) error {
	resp, err := s.ShippingCtrl.GetExternalShippingServices(ctx, model.EtopAccountID, q.GetExternalShippingServicesRequest)
	q.Result = &types.GetExternalShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return err
}

func (s *FulfillmentService) GetPublicFulfillment(ctx context.Context, q *GetPublicFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: q.Code,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbPublicFulfillment(query.Result)
	return nil
}

func (s *FulfillmentService) UpdateFulfillmentsShippingState(ctx context.Context, q *UpdateFulfillmentsShippingStateEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &shipmodelx.UpdateFulfillmentsShippingStateCommand{
		ShopID:        shopID,
		IDs:           q.Ids,
		ShippingState: q.ShippingState,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}
