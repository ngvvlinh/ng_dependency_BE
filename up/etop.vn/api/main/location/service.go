package location

import (
	"context"
)

type LocationQueryService interface {
	GetAllLocations(ctx context.Context, _ *GetAllLocationsQueryArgs) (*GetAllLocationsQueryResult, error)

	GetLocation(ctx context.Context, _ *GetLocationQueryArgs) (*LocationQueryResult, error)

	FindLocation(ctx context.Context, _ *FindLocationQueryArgs) (*LocationQueryResult, error)

	FindOrGetLocation(ctx context.Context, _ *FindOrGetLocationQueryArgs) (*LocationQueryResult, error)
}

//-- queries --//

type GetAllLocationsQueryArgs struct {
	All          bool
	ProvinceCode string
	DistrictCode string
}

type GetAllLocationsQueryResult struct {
	Provinces []*Province
	Districts []*District
	Wards     []*Ward
}

type LocationQueryResult struct {
	Province *Province
	District *District
	Ward     *Ward
}

type GetLocationQueryArgs struct {
	ProvinceCode     string
	DistrictCode     string
	WardCode         string
	LocationCodeType LocationCodeType
}

type FindLocationQueryArgs struct {
	Province string
	District string
	Ward     string
}

type FindOrGetLocationQueryArgs struct {
	Province     string
	District     string
	Ward         string
	ProvinceCode string
	DistrictCode string
	WardCode     string
}
