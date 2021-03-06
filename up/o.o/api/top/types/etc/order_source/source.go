package order_source

// +enum
// +enum:zero=null
type Source int

type NullSource struct {
	Enum  Source
	Valid bool
}

const (
	// +enum=unknown
	Unknown Source = 0

	// +enum=self
	Self Source = 1

	// +enum=import
	Import Source = 2

	// +enum=api
	API Source = 3

	// +enum=etop_pos
	EtopPOS Source = 5

	// +enum=etop_pxs
	EtopPXS Source = 6

	// +enum=etop_cmx
	EtopCMX Source = 7

	// +enum=ts_app
	TSApp Source = 8

	// +enum=etop_app
	EtopApp Source = 9

	// +enum=haravan
	Haravan Source = 10

	// +enum=ecomify
	Ecomify Source = 11
)
