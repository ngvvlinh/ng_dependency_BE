package location

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCustomRegion(context.Context, *CreateCustomRegionArgs) (*CustomRegion, error)

	UpdateCustomRegion(context.Context, *CustomRegion) error

	DeleteCustomRegion(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetAllLocations(ctx context.Context, _ *GetAllLocationsQueryArgs) (*GetAllLocationsQueryResult, error)

	GetLocation(ctx context.Context, _ *GetLocationQueryArgs) (*LocationQueryResult, error)

	FindLocation(ctx context.Context, _ *FindLocationQueryArgs) (*LocationQueryResult, error)

	FindOrGetLocation(ctx context.Context, _ *FindOrGetLocationQueryArgs) (*LocationQueryResult, error)

	// --- CustomRegion --- //
	GetCustomRegion(ctx context.Context, ID dot.ID) (*CustomRegion, error)

	ListCustomRegionsByCode(ctx context.Context, ProvinceCode string, DistrictCode string) ([]*CustomRegion, error)

	ListCustomRegions(context.Context, *meta.Empty) ([]*CustomRegion, error)
}

func (b QueryBus) DispatchAll(ctx context.Context, msgs ...interface{ query() }) error {
	for _, msg := range msgs {
		if err := b.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
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

// +convert:create=CustomRegion
type CreateCustomRegionArgs struct {
	Name          string
	Description   string
	ProvinceCodes []string
}
