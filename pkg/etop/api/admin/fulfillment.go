package admin

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type FulfillmentService struct {
	ShippingAgg shipping.CommandBus
	RedisStore  redis.Store
}

func (s *FulfillmentService) Clone() *FulfillmentService {
	res := *s
	return &res
}

func (s *FulfillmentService) UpdateFulfillment(ctx context.Context, q *UpdateFulfillmentEndpoint) error {
	cmd := &shipmodelx.AdminUpdateFulfillmentCommand{
		FulfillmentID:            q.Id,
		FullName:                 q.FullName,
		Phone:                    q.Phone,
		TotalCODAmount:           q.TotalCodAmount,
		IsPartialDelivery:        q.IsPartialDelivery,
		AdminNote:                q.AdminNote,
		ActualCompensationAmount: q.ActualCompensationAmount,
		ShippingState:            q.ShippingState,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *FulfillmentService) UpdateFulfillmentInfo(ctx context.Context, q *UpdateFulfillmentInfoEndpoint) error {
	cmd := &shipping.UpdateFulfillmentInfoCommand{
		ID:        q.Id,
		FullName:  q.FullName,
		Phone:     q.Phone,
		AdminNote: q.AdminNote,
	}
	if err := s.ShippingAgg.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return nil
}

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *GetFulfillmentsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		OrderID: q.OrderId,
		Status:  q.Status,
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if q.ShopId != 0 {
		query.ShopIDs = []dot.ID{q.ShopId}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagEtop),
		Paging:       cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *FulfillmentService) UpdateFulfillmentShippingState(ctx context.Context, r *UpdateFulfillmentShippingStateEndpoint) error {
	return cm.ErrTODO
}

func (s *FulfillmentService) UpdateFulfillmentShippingFee(ctx context.Context, r *UpdateFulfillmentShippingFeeEndpoint) error {
	return cm.ErrTODO
}
