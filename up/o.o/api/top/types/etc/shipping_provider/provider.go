package shipping_provider

// +enum
// +enum:zero=null
type ShippingProvider int

type NullShippingProvider struct {
	Enum  ShippingProvider
	Valid bool
}

const (
	// +enum=unknown
	Unknown ShippingProvider = 0

	// +enum=all
	All ShippingProvider = 22

	// +enum=manual
	Manual ShippingProvider = 20

	// +enum=ghn
	GHN ShippingProvider = 19

	// +enum=ghtk
	GHTK ShippingProvider = 21

	// +enum=vtpost
	VTPost ShippingProvider = 23

	// +enum=etop
	Etop ShippingProvider = 24

	// +enum=partner
	Partner ShippingProvider = 25

	// +enum=ninjavan
	NinjaVan ShippingProvider = 26

	// +enum=dhl
	DHL ShippingProvider = 27

	// +enum=ntx
	NTX ShippingProvider = 28
)

func (s ShippingProvider) Label() string {
	switch s {
	case GHN:
		return "Giao Hàng Nhanh"
	case GHTK:
		return "Giao Hàng Tiết Kiệm"
	case VTPost:
		return "Viettel Post"
	case Manual:
		return "Tự giao"
	case Partner:
		return "Đối tác vận chuyển"
	default:
		return ""
	}
}
