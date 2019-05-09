package v1

import (
	"context"

	"etop.vn/api/main/location"
	locationv1 "etop.vn/api/main/location/v1"
)

var _ locationv1.LocationQueryService = &Server{}

type Server struct {
	s location.LocationQueryService
}

func (im *Server) GetAllLocations(ctx context.Context, args *locationv1.GetAllLocationsQueryArgs) (*locationv1.GetAllLocationsQueryResult, error) {
	return im.s.GetAllLocations(ctx, args)
}

func (im *Server) GetLocation(ctx context.Context, args *locationv1.GetLocationQueryArgs) (*locationv1.LocationQueryResult, error) {
	return im.s.GetLocation(ctx, args)
}

func (im *Server) FindLocation(ctx context.Context, args *locationv1.FindLocationQueryArgs) (*locationv1.LocationQueryResult, error) {
	return im.s.FindLocation(ctx, args)
}

func (im *Server) FindOrGetLocation(ctx context.Context, args *locationv1.FindOrGetLocationQueryArgs) (*locationv1.LocationQueryResult, error) {
	return im.s.FindOrGetLocation(ctx, args)
}
