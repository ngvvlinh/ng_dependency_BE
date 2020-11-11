package shipping

import (
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
)

// +enum
// +enum:zero=null
type State int

type NullState struct {
	Enum  State
	Valid bool
}

const (
	// +enum=unknown
	Unknown State = 0

	// +enum=default
	Default State = 1

	// +enum=created
	Created State = 2

	// +enum=confirmed
	Confirmed State = 3

	// +enum=processing
	Processing State = 4

	// +enum=picking
	// +enum:RefName:Đang lấy hàng
	Picking State = 5

	// +enum=holding
	// +enum:RefName:Đang giữ hàng
	Holding State = 6

	// +enum=returning
	// +enum:RefName:Đang trả hàng
	Returning State = 7

	// +enum=returned
	// +enum:RefName:Đã trả hàng
	Returned State = 8

	// +enum=delivering
	// +enum:RefName:Đang giao hàng
	Delivering State = 9

	// +enum=delivered
	// +enum:RefName:Đã giao hàng
	Delivered State = 10

	// +enum=cancelled
	// +enum:RefName:Đã huỷ đơn hàng
	Cancelled State = -1

	// Trạng thái Bồi hoàn
	// +enum=undeliverable
	Undeliverable State = -2
)

func (s State) Text() string {
	name := s.Name()
	if name == "" {
		return "Không xác định"
	}
	return name
}

func (s State) ToStatus4() status4.Status {
	switch s {
	case Default:
		return status4.Z
	case Cancelled, Returned:
		return status4.N
	case Delivered:
		return status4.P
	}
	return status4.S
}

func (s State) ToStatus5() status5.Status {
	switch s {
	case Default:
		return status5.Z
	case Cancelled:
		return status5.N
	case Returning, Returned:
		return status5.NS
	case Delivered:
		return status5.P
	default:
		return status5.S
	}
}
