package substate

// +enum
// +enum:zero=null
type Substate int

type NullSubstate struct {
	Enum  Substate
	Valid bool
}

const (
	// +enum=default
	Default Substate = 1

	// +enum=pick_fail
	PickFail Substate = 17

	// Giao hàng thất bại, NVC sẽ giao lại tối đa 3 lần
	//
	// +enum=delivery_fail
	DeliveryFail Substate = 25

	// Chờ trả hàng, NVC sẽ giữ hàng, kích hoạt giao lại nếu shop yêu cầu, nếu không tiến hành trả hàng
	//
	// +enum=devivery_giveup
	DeliveryGiveup Substate = 37

	// +enum=return_fail
	ReturnFail Substate = 43
)
