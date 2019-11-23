// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: external/etc/change_type/type.proto

package change_type

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ChangeType int32

const (
	ChangeType_unknown ChangeType = 0
	ChangeType_update  ChangeType = 1
	ChangeType_create  ChangeType = 2
	ChangeType_delete  ChangeType = 3
)

var ChangeType_name = map[int32]string{
	0: "unknown",
	1: "update",
	2: "create",
	3: "delete",
}

var ChangeType_value = map[string]int32{
	"unknown": 0,
	"update":  1,
	"create":  2,
	"delete":  3,
}

func (x ChangeType) Enum() *ChangeType {
	p := new(ChangeType)
	*p = x
	return p
}

func (x ChangeType) String() string {
	return proto.EnumName(ChangeType_name, int32(x))
}

func (x *ChangeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ChangeType_value, data, "ChangeType")
	if err != nil {
		return err
	}
	*x = ChangeType(value)
	return nil
}

func (ChangeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d58df26b033b6104, []int{0}
}

func init() {
	proto.RegisterEnum("change_type.ChangeType", ChangeType_name, ChangeType_value)
}

func init() {
	proto.RegisterFile("external/etc/change_type/type.proto", fileDescriptor_d58df26b033b6104)
}

var fileDescriptor_d58df26b033b6104 = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4e, 0xad, 0x28, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x4f, 0x2d, 0x49, 0xd6, 0x4f, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0x8d,
	0x2f, 0xa9, 0x2c, 0x48, 0xd5, 0x07, 0x11, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xdc, 0x48,
	0xe2, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x71, 0x7d, 0x10, 0x0b, 0xa2, 0x44, 0xcb, 0x96,
	0x8b, 0xcb, 0x19, 0xac, 0x28, 0xa4, 0xb2, 0x20, 0x55, 0x88, 0x9b, 0x8b, 0xbd, 0x34, 0x2f, 0x3b,
	0x2f, 0xbf, 0x3c, 0x4f, 0x80, 0x41, 0x88, 0x8b, 0x8b, 0xad, 0xb4, 0x20, 0x25, 0xb1, 0x24, 0x55,
	0x80, 0x11, 0xc4, 0x4e, 0x2e, 0x4a, 0x05, 0xb1, 0x99, 0x40, 0xec, 0x94, 0xd4, 0x9c, 0xd4, 0x92,
	0x54, 0x01, 0x66, 0x27, 0xdd, 0x19, 0x8f, 0xe5, 0x18, 0xa2, 0xd4, 0x53, 0x4b, 0xf2, 0x0b, 0xf4,
	0xca, 0xf2, 0xf4, 0x13, 0x0b, 0x32, 0xf5, 0x0b, 0x92, 0xf4, 0x71, 0xb9, 0x0d, 0x10, 0x00, 0x00,
	0xff, 0xff, 0xec, 0x4d, 0xbc, 0x15, 0xb6, 0x00, 0x00, 0x00,
}