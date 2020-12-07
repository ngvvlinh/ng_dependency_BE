package call_log_direction

// +enum
// +enum:sql=int
// +enum:zero=null
type CallLogDirection int

type NullCallLogDirection struct {
	Enum  CallLogDirection
	Valid bool
}

const (
	// +enum=unknown
	Unknown CallLogDirection = 0

	// +enum=in
	In CallLogDirection = 3

	// +enum=out
	Out CallLogDirection = 9

	// +enum=ext
	Ext CallLogDirection = 15
)
