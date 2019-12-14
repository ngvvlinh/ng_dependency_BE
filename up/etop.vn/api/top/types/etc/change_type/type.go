package change_type

// +enum
// +enum:zero=null
type ChangeType int

type NullChangeType struct {
	Enum  ChangeType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ChangeType = 0

	// +enum=update
	Update ChangeType = 1

	// +enum=create
	Create ChangeType = 2

	// +enum=delete
	Delete ChangeType = 3
)
