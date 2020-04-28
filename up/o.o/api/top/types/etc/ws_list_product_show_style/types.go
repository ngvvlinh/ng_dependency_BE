package ws_list_product_show_style

// +enum
// +enum:zero=null
type WsListProductShowStyle int

type NullWsListProductShowStyle struct {
	Enum  WsListProductShowStyle
	Valid bool
}

const (
	// +enum=first
	First WsListProductShowStyle = 0

	// +enum=second
	Second WsListProductShowStyle = 1

	// +enum=third
	Third WsListProductShowStyle = 2

	// +enum=fourth
	Fourth WsListProductShowStyle = 3
)
