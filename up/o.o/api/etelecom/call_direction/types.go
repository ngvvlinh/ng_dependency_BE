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
	// +enum:RefName:Gọi vào
	In CallDirection = 3

	// +enum=out
	// +enum:RefName:Gọi ra
	Out CallDirection = 9

	// +enum=ext
	// +enum:RefName:Gọi nội bộ
	Ext CallDirection = 15

	// +enum=ext_in
	// +enum:RefName:Gọi vào nội bộ
	ExtIn CallDirection = 17

	// +enum=ext_out
	// +enum:RefName:Gọi ra nội bộ
	ExtOut CallDirection = 21
)
