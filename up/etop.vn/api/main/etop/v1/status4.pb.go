// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/etop/v1/status4.proto

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

var Status4_name = map[int32]string{
	0:   "Z",
	1:   "P",
	2:   "S",
	127: "N",
}

var Status4_value = map[string]int32{
	"Z": 0,
	"P": 1,
	"S": 2,
	"N": 127,
}

func (x Status4) Enum() *Status4 {
	p := new(Status4)
	*p = x
	return p
}

func (x Status4) String() string {
	return proto.EnumName(Status4_name, int32(x))
}

func (x *Status4) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Status4_value, data, "Status4")
	if err != nil {
		return err
	}
	*x = Status4(value)
	return nil
}

func (Status4) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_185d5e4044f7cedc, []int{0}
}

func init() {
	proto.RegisterEnum("etop.vn.api.main.etop.v1.status4.Status4", Status4_name, Status4_value)
}

func init() {
	proto.RegisterFile("etop.vn/api/main/etop/v1/status4.proto", fileDescriptor_185d5e4044f7cedc)
}

var fileDescriptor_185d5e4044f7cedc = []byte{
	// 134 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0x2d, 0xc9, 0x2f,
	0xd0, 0x2b, 0xcb, 0xd3, 0x4f, 0x2c, 0xc8, 0xd4, 0xcf, 0x4d, 0xcc, 0xcc, 0xd3, 0x07, 0x09, 0xe8,
	0x97, 0x19, 0xea, 0x17, 0x97, 0x24, 0x96, 0x94, 0x16, 0x9b, 0xe8, 0x15, 0x14, 0xe5, 0x97, 0xe4,
	0x0b, 0x29, 0x40, 0xd5, 0xe9, 0x25, 0x16, 0x64, 0xea, 0x81, 0xd4, 0xe9, 0x41, 0x04, 0x0c, 0xf5,
	0xa0, 0xea, 0xa4, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x8a, 0xf5, 0x41, 0x2c, 0x88, 0x3e, 0x2d,
	0x55, 0x2e, 0xf6, 0x60, 0x88, 0x02, 0x21, 0x56, 0x2e, 0xc6, 0x28, 0x01, 0x06, 0x10, 0x15, 0x20,
	0xc0, 0x08, 0xa2, 0x82, 0x05, 0x98, 0x40, 0x94, 0x9f, 0x40, 0xbd, 0x13, 0xf7, 0x8c, 0xc7, 0x72,
	0x0c, 0x51, 0xac, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x78,
	0xcd, 0x70, 0x94, 0x00, 0x00, 0x00,
}