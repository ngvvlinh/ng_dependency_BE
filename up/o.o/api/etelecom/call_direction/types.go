package call_direction

// +enum
// +enum:sql=int
// +enum:zero=null
type CallDirection int

type NullCallDirection struct {
	Enum  CallDirection
	Valid bool
}

const (
	// +enum=unknown
	Unknown CallDirection = 0

	// +enum=in
	In CallDirection = 3

	// +enum=out
	Out CallDirection = 9

	// +enum=ext
	Ext CallDirection = 15
)
