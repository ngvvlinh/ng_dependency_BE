package connection_type

// +enum
// +enum:zero=null
type ConnectionType int

type NullConnectionType struct {
	Enum  ConnectionType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ConnectionType = 0

	// +enum=shipping
	Shipping ConnectionType = 1

	// +enum=crm
	CRM ConnectionType = 2
)

// +enum
// +enum:zero=null
type ConnectionSubtype int

type NullConnectionSubtype struct {
	Enum  ConnectionSubtype
	Valid bool
}

const (
	// +enum=unknown
	ConnectionSubtypeUnknown ConnectionSubtype = 0

	// +enum=shipment
	ConnectionSubtypeShipment ConnectionSubtype = 1

	// +enum=manual
	ConnectionSubtypeManual ConnectionSubtype = 2

	// +enum=shipnow
	ConnectionSubtypeShipnow ConnectionSubtype = 3
)

// +enum
// +enum:zero=null
type ConnectionMethod int

type NullConnectionMethod struct {
	Enum  ConnectionMethod
	Valid bool
}

const (
	// +enum=unknown
	ConnectionMethodUnknown ConnectionMethod = 0

	// built-in: dùng tài khoản có sẵn để tạo đơn (tức tạo đơn từ tài khoản của TopShip hoặc tài khoản whitelabel nào đó)
	// topship: backward-compatible (remove later)
	// +enum=builtin,topship
	ConnectionMethodBuiltin ConnectionMethod = 1

	// direct: dùng trực tiếp tài khoản từ NVC để tạo đơn
	// +enum=direct
	ConnectionMethodDirect ConnectionMethod = 2
)

// +enum
// +enum:zero=null
type ConnectionProvider int

type NullConnectionProvider struct {
	Enum  ConnectionProvider
	Valid bool
}

const (
	// +enum=unknown
	ConnectionProviderUnknown ConnectionProvider = 0

	// +enum=ghn
	ConnectionProviderGHN ConnectionProvider = 1

	// +enum=ghtk
	ConnectionProviderGHTK ConnectionProvider = 2

	// +enum=vtpost
	ConnectionProviderVTP ConnectionProvider = 3

	// +enum=partner
	ConnectionProviderPartner ConnectionProvider = 4

	// +enum=ahamove
	ConnectionProviderAhamove ConnectionProvider = 5

	// +enum=ninjavan
	ConnectionProviderNinjaVan ConnectionProvider = 6

	// +enum=dhl
	ConnectionProviderDHL ConnectionProvider = 7

	// +enum=suitecrm
	ConnectionProviderSuiteCRM ConnectionProvider = 8
)
