package try_on

// +enum
// +enum:zero=null
type TryOnCode int

type NullTryOnCode struct {
	Enum  TryOnCode
	Valid bool
}

const (
	// +enum=unknown
	Unknown TryOnCode = 0

	// +enum=none
	// +enum:RefName:Không cho xem hàng
	None TryOnCode = 1

	// +enum=open
	// +enum:RefName:Cho xem hàng không thử
	Open TryOnCode = 2

	// +enum=try
	// +enum:RefName:Cho thử hàng
	Try TryOnCode = 3
)

const NoteChoThuHang = "cho thu hang"
const NoteChoXemHang = "cho xem hang khong thu"
const NoteKhongXemHang = "khong cho xem hang"
