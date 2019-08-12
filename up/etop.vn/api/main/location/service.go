package location

import (
	"context"

	locationv1 "etop.vn/api/main/location/v1"
)

type LocationQueryService interface {
	GetAllLocations(ctx context.Context, _ *GetAllLocationsQueryArgs) (*locationv1.GetAllLocationsQueryResult, error)

	GetLocation(ctx context.Context, _ *GetLocationQueryArgs) (*LocationQueryResult, error)

	FindLocation(ctx context.Context, _ *FindLocationQueryArgs) (*LocationQueryResult, error)

	FindOrGetLocation(ctx context.Context, _ *FindOrGetLocationQueryArgs) (*LocationQueryResult, error)
}

//-- queries --//

type GetAllLocationsQueryArgs = locationv1.GetAllLocationsQueryArgs
type GetAllLocationsQueryResult = locationv1.GetAllLocationsQueryResult
type LocationQueryResult = locationv1.LocationQueryResult
type GetLocationQueryArgs = locationv1.GetLocationQueryArgs
type FindLocationQueryArgs = locationv1.FindLocationQueryArgs
type FindOrGetLocationQueryArgs = locationv1.FindOrGetLocationQueryArgs
