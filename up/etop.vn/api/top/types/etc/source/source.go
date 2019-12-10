package source

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
	Api Source = 3

	// +enum=etop_pos
	EtopPos Source = 5

	// +enum=etop_pxs
	EtopPxs Source = 6

	// +enum=etop_cmx
	EtopCmx Source = 7

	// +enum=ts_app
	TsApp Source = 8
)
