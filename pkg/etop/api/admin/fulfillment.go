package admin

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/shipping"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
)

type FulfillmentService struct {
	session.Session

	EventBus      capi.EventBus
	IdentityQuery identity.QueryBus
	RedisStore    redis.Store
	ShippingAggr  shipping.CommandBus
	ShippingQuery shipping.QueryBus
}

func (s *FulfillmentService) Clone() admin.FulfillmentService {
	res := *s
	return &res
}

func (s *FulfillmentService) UpdateFulfillmentInfo(ctx context.Context, q *admin.UpdateFulfillmentInfoRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.UpdateFulfillmentInfoCommand{
		FulfillmentID: q.ID,
		ShippingCode:  q.ShippingCode,
		FullName:      q.FullName,
		Phone:         q.Phone,
		AdminNote:     q.AdminNote,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}, nil
}

func (s *FulfillmentService) UpdateFulfillmentCODAmount(ctx context.Context, q *admin.UpdateFulfillmentCODAmountRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.UpdateFulfillmentCODAmountCommand{
		FulfillmentID:     q.ID,
		ShippingCode:      q.ShippingCode,
		TotalCODAmount:    q.TotalCODAmount,
		IsPartialDelivery: q.IsPartialDelivery,
		AdminNote:         q.AdminNote,
		UpdatedBy:         s.SS.Claim().UserID,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *pbcm.IDRequest) (*types.Fulfillment, error) {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return result, nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *admin.GetFulfillmentsRequest) (*types.FulfillmentsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		OrderID:       q.OrderId,
		Status:        q.Status,
		ConnectionIDs: q.ConnectionIDs,
		Paging:        paging,
		Filters:       cmapi.ToFilters(q.Filters),
	}
	if q.ShopId != 0 {
		query.ShopIDs = []dot.ID{q.ShopId}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagEtop),
		Paging:       cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *FulfillmentService) UpdateFulfillmentShippingState(ctx context.Context, r *admin.UpdateFulfillmentShippingStateRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.UpdateFulfillmentShippingStateCommand{
		FulfillmentID:            r.ID,
		ShippingCode:             r.ShippingCode,
		ShippingState:            r.ShippingState,
		ActualCompensationAmount: r.ActualCompensationAmount,
		UpdatedBy:                s.SS.Claim().UserID,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return result, nil
}

func (s *FulfillmentService) UpdateFulfillmentShippingFees(ctx context.Context, r *admin.UpdateFulfillmentShippingFeesRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.UpdateFulfillmentShippingFeesCommand{
		FulfillmentID:    r.ID,
		ShippingCode:     r.ShippingCode,
		ShippingFeeLines: convertpb.Convert_api_ShippingFeeLines_To_core_ShippingFeeLines(r.ShippingFeeLines),
		TotalCODAmount:   r.TotalCODAmount,
		UpdatedBy:        s.SS.Claim().UserID,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return result, nil
}

func (s *FulfillmentService) addShippingFee(ctx context.Context, r *admin.AddShippingFeeRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.AddFulfillmentShippingFeeCommand{
		FulfillmentID:   r.ID,
		ShippingCode:    r.ShippingCode,
		ShippingFeeType: r.ShippingFeeType,
		UpdatedBy:       s.SS.Claim().UserID,
	}
	if err := s.ShippingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	resp := &pbcm.UpdatedResponse{Updated: 1}
	return resp, nil
}

func (s *FulfillmentService) AddShippingFee(ctx context.Context, r *admin.AddShippingFeeRequest) (*pbcm.UpdatedResponse, error) {
	key := fmt.Sprintf("addShippingFee %v-%v", r.ID, r.ShippingFeeType.String())
	res, _, err := idempgroup.DoAndWrap(ctx, key, 15*time.Second, "Thêm cước phí cho đơn vận chuyển", func() (interface{}, error) { return s.addShippingFee(ctx, r) })
	if err != nil {
		return nil, err
	}
	result := res.(*pbcm.UpdatedResponse)
	return result, nil
}
