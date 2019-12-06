package change_type

// +enum
type ChangeType int

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
