package order_source

// +enum
// +enum:zero=null
type Source int

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
)
