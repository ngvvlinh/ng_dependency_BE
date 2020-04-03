package location_type

// +enum
// +enum:zero=null
type RegionType int

type NullRegionType struct {
	Enum  RegionType
	Valid bool
}

const (
	// +enum=north
	NorthRegion RegionType = 1

	// +enum=middle
	MiddleRegion RegionType = 2

	// +enum=south
	SouthRegion RegionType = 3
)

func RegionTypeContains(regionTypes []RegionType, typ RegionType) bool {
	for _, rt := range regionTypes {
		if rt == typ {
			return true
		}
	}
	return false
}

// +enum
// +enum:zero=null
type ShippingLocationType int

type NullShippingLocationType struct {
	Enum  ShippingLocationType
	Valid bool
}

const (
	// +enum=pick
	// +enum:RefName:Địa chỉ lấy hàng
	ShippingLocationPick ShippingLocationType = 1

	// +enum=deliver
	// +enum:RefName:Địa chỉ giao hàng
	ShippingLocationDeliver ShippingLocationType = 2
)
