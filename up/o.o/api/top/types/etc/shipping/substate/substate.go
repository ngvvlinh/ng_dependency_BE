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
	// +enum:RefName:Lấy hàng thất bại
	PickFail Substate = 17

	// Giao hàng thất bại, NVC sẽ giao lại tối đa 3 lần
	//
	// +enum=delivery_fail
	// +enum:RefName:Giao hàng thất bại
	DeliveryFail Substate = 25

	// Chờ trả hàng, NVC sẽ giữ hàng, kích hoạt giao lại nếu shop yêu cầu, nếu không tiến hành trả hàng
	//
	// +enum=devivery_giveup
	// +enum:RefName:Giao hàng thất bại. Chờ trả hàng
	DeliveryGiveup Substate = 37

	// +enum=return_fail
	// +enum:RefName:Trả hàng thất bại
	ReturnFail Substate = 43

	// +enum=cancelling
	Cancelling Substate = 51
)
