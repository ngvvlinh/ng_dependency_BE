package location

import (
	"context"

	locationv1 "etop.vn/api/main/location/v1"
	"etop.vn/api/meta"
)

type Bus struct {
	meta.Bus
}

type LocationQueryService interface {
	GetAllLocations(ctx context.Context, args *GetAllLocationsQueryArgs) (*locationv1.GetAllLocationsQueryResult, error)

	GetLocation(ctx context.Context, args *GetLocationQueryArgs) (*LocationQueryResult, error)

	FindLocation(ctx context.Context, args *FindLocationQueryArgs) (*LocationQueryResult, error)

	FindOrGetLocation(ctx context.Context, args *FindOrGetLocationQueryArgs) (*LocationQueryResult, error)
}

//-- queries --//

type GetAllLocationsQueryArgs = locationv1.GetAllLocationsQueryArgs
type GetAllLocationsQueryResult = locationv1.GetAllLocationsQueryResult
type LocationQueryResult = locationv1.LocationQueryResult
type GetLocationQueryArgs = locationv1.GetLocationQueryArgs
type FindLocationQueryArgs = locationv1.FindLocationQueryArgs
type FindOrGetLocationQueryArgs = locationv1.FindOrGetLocationQueryArgs
