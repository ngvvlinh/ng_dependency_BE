package filter_type

// +enum
// +enum:zero=null
type FilterType int

type NullFilterType struct {
	Enum  FilterType
	Valid bool
}

const (
	// +enum=include
	Include FilterType = 1

	// +enum=exclude
	Exclude FilterType = 2
)
