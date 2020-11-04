package vht

import "o.o/common/jsonx"

type VHTRegisterUser struct {
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

func (m *VHTRegisterUser) String() string { return jsonx.MustMarshalToString(m) }
