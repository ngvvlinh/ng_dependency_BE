package ws_banner_show_style

// +enum
// +enum:zero=null
type WsBannerShowStyle int

type NullWsBannerShowStyle struct {
	Enum  WsBannerShowStyle
	Valid bool
}

const (
	// +enum=slider
	Slider WsBannerShowStyle = 0

	// +enum=grid
	Grind WsBannerShowStyle = 1
)
