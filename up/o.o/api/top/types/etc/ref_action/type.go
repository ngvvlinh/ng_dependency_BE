package ref_action

// +enum
// +enum:zero=null
type RefAction int

type NullRefAction struct {
	Enum  RefAction
	Valid bool
}

const (
	// +enum=create
	Create RefAction = 0

	// +enum=cancel
	Cancel RefAction = 1
)
