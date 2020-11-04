package vht

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/mc/vht
// +gen:swagger:title: VHT - VNPost
// +gen:swagger:version=v1

// +apix:path=/vht.User
type UserService interface {
	// Chỉ dùng để VHT gọi tạo User
	// Sử dụng token của VNPost wl partner
	// Dùng để đăng ký user, sử dụng cho webphone (gọi điện, tạo được ticket trên eTop)
	//
	// +apix:path=RegisterUser
	RegisterUser(context.Context, *VHTRegisterUser) (*cm.Empty, error)
}
