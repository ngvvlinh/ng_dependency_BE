// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: external/etc/auth_mode/mode.proto

package auth_mode

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AuthMode int32

const (
	AuthMode_default AuthMode = 0
	AuthMode_manual  AuthMode = 1
)

var AuthMode_name = map[int32]string{
	0: "default",
	1: "manual",
}

var AuthMode_value = map[string]int32{
	"default": 0,
	"manual":  1,
}

func (x AuthMode) Enum() *AuthMode {
	p := new(AuthMode)
	*p = x
	return p
}

func (x AuthMode) String() string {
	return proto.EnumName(AuthMode_name, int32(x))
}

func (x *AuthMode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AuthMode_value, data, "AuthMode")
	if err != nil {
		return err
	}
	*x = AuthMode(value)
	return nil
}

func (AuthMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_26ae1217a6caae4c, []int{0}
}

func init() {
	proto.RegisterEnum("auth_mode.AuthMode", AuthMode_name, AuthMode_value)
}

func init() { proto.RegisterFile("external/etc/auth_mode/mode.proto", fileDescriptor_26ae1217a6caae4c) }

var fileDescriptor_26ae1217a6caae4c = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0xad, 0x28, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x4f, 0x2d, 0x49, 0xd6, 0x4f, 0x2c, 0x2d, 0xc9, 0x88, 0xcf, 0xcd,
	0x4f, 0x49, 0xd5, 0x07, 0x11, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x9c, 0x70, 0x51, 0x29,
	0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0, 0xa8, 0x3e, 0x88, 0x05, 0x51, 0xa0, 0xa5, 0xcc, 0xc5, 0xe1,
	0x58, 0x5a, 0x92, 0xe1, 0x9b, 0x9f, 0x92, 0x2a, 0xc4, 0xcd, 0xc5, 0x9e, 0x92, 0x9a, 0x96, 0x58,
	0x9a, 0x53, 0x22, 0xc0, 0x20, 0xc4, 0xc5, 0xc5, 0x96, 0x9b, 0x98, 0x57, 0x9a, 0x98, 0x23, 0xc0,
	0xe8, 0xa4, 0x3f, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x66, 0x6a, 0x49, 0x7e, 0x81, 0x5e, 0x59, 0x9e,
	0x7e, 0x52, 0x62, 0x72, 0x76, 0x6a, 0x5e, 0x8a, 0x7e, 0x41, 0x92, 0x3e, 0x76, 0x17, 0x00, 0x02,
	0x00, 0x00, 0xff, 0xff, 0xab, 0x27, 0xd1, 0x05, 0x9a, 0x00, 0x00, 0x00,
}