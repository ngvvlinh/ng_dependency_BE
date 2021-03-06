package gender

// +enum
// +enum:zero=null
type Gender int

type NullGender struct {
	Enum  Gender
	Valid bool
}

const (
	// +enum=unknown
	Unknown Gender = 0

	// +enum=male
	Male Gender = 1

	// +enum=female
	Female Gender = 2

	// +enum=other
	Other Gender = 3
)
