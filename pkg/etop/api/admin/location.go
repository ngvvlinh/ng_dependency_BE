package admin

import (
	"context"

	"o.o/api/main/location"
	"o.o/api/top/int/admin"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/admin/convert"
	"o.o/backend/pkg/etop/authorize/session"
)

type LocationService struct {
	session.Session

	LocationAggr  location.CommandBus
	LocationQuery location.QueryBus
}

func (s *LocationService) Clone() admin.LocationService {
	res := *s
	return &res
}

func (s *LocationService) GetCustomRegion(ctx context.Context, r *pbcm.IDRequest) (*admin.CustomRegion, error) {
	query := &location.GetCustomRegionQuery{
		ID: r.Id,
	}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convert.PbCustomRegion(query.Result)
	return result, nil
}

func (s *LocationService) GetCustomRegions(ctx context.Context, r *pbcm.Empty) (*admin.GetCustomRegionsResponse, error) {
	query := &location.ListCustomRegionsQuery{}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetCustomRegionsResponse{
		CustomRegions: convert.PbCustomRegions(query.Result),
	}
	return result, nil
}

func (s *LocationService) CreateCustomRegion(ctx context.Context, r *admin.CreateCustomRegionRequest) (*admin.CustomRegion, error) {
	cmd := &location.CreateCustomRegionCommand{
		Name:          r.Name,
		Description:   r.Description,
		ProvinceCodes: r.ProvinceCodes,
	}
	if err := s.LocationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbCustomRegion(cmd.Result)
	return result, nil
}

func (s *LocationService) UpdateCustomRegion(ctx context.Context, r *admin.CustomRegion) (*pbcm.UpdatedResponse, error) {
	cmd := &location.UpdateCustomRegionCommand{
		ID:            r.ID,
		Name:          r.Name,
		ProvinceCodes: r.ProvinceCodes,
		Description:   r.Description,
	}
	if err := s.LocationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *LocationService) DeleteCustomRegion(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &location.DeleteCustomRegionCommand{
		ID: r.Id,
	}
	if err := s.LocationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}
