package try_on

// +enum
// +enum:zero=null
type TryOnCode int

const (
	// +enum=unknown
	Unknown TryOnCode = 0

	// +enum=none
	None TryOnCode = 1

	// +enum=open
	Open TryOnCode = 2

	// +enum=try
	Try TryOnCode = 3
)

const NoteChoThuHang = "cho thu hang"
const NoteChoXemHang = "cho xem hang khong thu"
const NoteKhongXemHang = "khong cho xem hang"
