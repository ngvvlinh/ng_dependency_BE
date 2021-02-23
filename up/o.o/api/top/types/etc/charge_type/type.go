package charge_type

// +enum
// +enum:zero=null
type ChargeType int

type NullChargeType struct {
	Enum  ChargeType
	Valid bool
}

const (
	// +enum=prepaid
	Prepaid ChargeType = 1

	// +enum=postpaid
	Postpaid ChargeType = 5

	// +enum=free
	Free ChargeType = 9
)
