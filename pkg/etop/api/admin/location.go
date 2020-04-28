package admin

import (
	"context"

	"o.o/api/main/location"
	"o.o/api/top/int/admin"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
)

func (s *LocationService) GetCustomRegion(ctx context.Context, r *GetCustomRegionEndpoint) error {
	query := &location.GetCustomRegionQuery{
		ID: r.Id,
	}
	if err := locationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomRegion(query.Result)
	return nil
}

func (s *LocationService) GetCustomRegions(ctx context.Context, r *GetCustomRegionsEndpoint) error {
	query := &location.ListCustomRegionsQuery{}
	if err := locationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetCustomRegionsResponse{
		CustomRegions: convertpb.PbCustomRegions(query.Result),
	}
	return nil
}

func (s *LocationService) CreateCustomRegion(ctx context.Context, r *CreateCustomRegionEndpoint) error {
	cmd := &location.CreateCustomRegionCommand{
		Name:          r.Name,
		Description:   r.Description,
		ProvinceCodes: r.ProvinceCodes,
	}
	if err := locationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomRegion(cmd.Result)
	return nil
}

func (s *LocationService) UpdateCustomRegion(ctx context.Context, r *UpdateCustomRegionEndpoint) error {
	cmd := &location.UpdateCustomRegionCommand{
		ID:            r.ID,
		Name:          r.Name,
		ProvinceCodes: r.ProvinceCodes,
		Description:   r.Description,
	}
	if err := locationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *LocationService) DeleteCustomRegion(ctx context.Context, r *DeleteCustomRegionEndpoint) error {
	cmd := &location.DeleteCustomRegionCommand{
		ID: r.Id,
	}
	if err := locationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}
