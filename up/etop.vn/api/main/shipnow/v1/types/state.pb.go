// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/shipnow/v1/types/state.proto

package types

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
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

var State_name = map[int32]string{
	0:   "default",
	1:   "created",
	2:   "assigning",
	3:   "picking",
	4:   "delivering",
	5:   "delivered",
	6:   "returning",
	7:   "returned",
	101: "unknown",
	126: "undeliverable",
	127: "cancelled",
}

var State_value = map[string]int32{
	"default":       0,
	"created":       1,
	"assigning":     2,
	"picking":       3,
	"delivering":    4,
	"delivered":     5,
	"returning":     6,
	"returned":      7,
	"unknown":       101,
	"undeliverable": 126,
	"cancelled":     127,
}

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}

func (x *State) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(State_value, data, "State")
	if err != nil {
		return err
	}
	*x = State(value)
	return nil
}

func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_42995b7c28920429, []int{0}
}

func init() {
	proto.RegisterEnum("etop.vn.api.main.shipnow.v1.state.State", State_name, State_value)
}

func init() {
	proto.RegisterFile("etop.vn/api/main/shipnow/v1/types/state.proto", fileDescriptor_42995b7c28920429)
}

var fileDescriptor_42995b7c28920429 = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x2c, 0x8f, 0x31, 0x4e, 0x03, 0x31,
	0x10, 0x45, 0x13, 0x60, 0x09, 0x38, 0x04, 0x19, 0x8b, 0x8a, 0x02, 0x89, 0x16, 0x09, 0x5b, 0xb9,
	0x02, 0x57, 0xa0, 0xa3, 0x1b, 0xd6, 0xc3, 0x62, 0xc5, 0x8c, 0x2d, 0x7b, 0x76, 0x23, 0x1a, 0xb8,
	0x06, 0x47, 0xe0, 0x5e, 0x5c, 0x24, 0x9a, 0x8d, 0x3b, 0x3f, 0xfd, 0xf7, 0x2c, 0x8d, 0x7a, 0x42,
	0x4e, 0xd9, 0x4e, 0xe4, 0x20, 0x07, 0xf7, 0x09, 0x81, 0x5c, 0xfd, 0x08, 0x99, 0xd2, 0xde, 0x4d,
	0x5b, 0xc7, 0x5f, 0x19, 0xab, 0xab, 0x0c, 0x8c, 0x36, 0x97, 0xc4, 0xc9, 0x3c, 0x34, 0xdd, 0x42,
	0x0e, 0x56, 0x74, 0xdb, 0x74, 0x3b, 0x6d, 0xed, 0x2c, 0xde, 0xdd, 0x0e, 0x69, 0x48, 0xb3, 0xed,
	0xe4, 0x75, 0x0c, 0x1f, 0xff, 0x96, 0xaa, 0x7b, 0x91, 0xdd, 0xac, 0xd5, 0xca, 0xe3, 0x3b, 0x8c,
	0x91, 0xf5, 0x42, 0xa0, 0x2f, 0x08, 0x8c, 0x5e, 0x2f, 0xcd, 0x46, 0x5d, 0x42, 0xad, 0x61, 0xa0,
	0x40, 0x83, 0x3e, 0x91, 0x2d, 0x87, 0x7e, 0x27, 0x70, 0x6a, 0xae, 0x95, 0xf2, 0x18, 0xc3, 0x84,
	0x45, 0xf8, 0x4c, 0xdc, 0xc6, 0xe8, 0x75, 0x27, 0x58, 0x90, 0xc7, 0x32, 0xa7, 0xe7, 0xe6, 0x4a,
	0x5d, 0x1c, 0x11, 0xbd, 0x5e, 0xc9, 0x47, 0x23, 0xed, 0x28, 0xed, 0x49, 0xa3, 0xb9, 0x51, 0x9b,
	0x91, 0x5a, 0x0a, 0x6f, 0x11, 0xf5, 0xb7, 0xc4, 0x3d, 0x50, 0x8f, 0x31, 0xa2, 0xd7, 0x3f, 0xcf,
	0xeb, 0xdf, 0xff, 0xfb, 0xc5, 0x6b, 0x37, 0x1f, 0x7f, 0x08, 0x00, 0x00, 0xff, 0xff, 0x15, 0xc5,
	0x98, 0xfd, 0x20, 0x01, 0x00, 0x00,
}
