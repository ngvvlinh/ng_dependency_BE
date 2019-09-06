// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/etc/shipping_fee_type/type.proto

package shipping_fee_type

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

type ShippingFeeType int32

const (
	ShippingFeeType_main           ShippingFeeType = 0
	ShippingFeeType_return         ShippingFeeType = 1
	ShippingFeeType_adjustment     ShippingFeeType = 2
	ShippingFeeType_insurance      ShippingFeeType = 3
	ShippingFeeType_tax            ShippingFeeType = 4
	ShippingFeeType_other          ShippingFeeType = 5
	ShippingFeeType_cods           ShippingFeeType = 6
	ShippingFeeType_address_change ShippingFeeType = 7
	ShippingFeeType_discount       ShippingFeeType = 8
	ShippingFeeType_unknown        ShippingFeeType = 127
)

var ShippingFeeType_name = map[int32]string{
	0:   "main",
	1:   "return",
	2:   "adjustment",
	3:   "insurance",
	4:   "tax",
	5:   "other",
	6:   "cods",
	7:   "address_change",
	8:   "discount",
	127: "unknown",
}

var ShippingFeeType_value = map[string]int32{
	"main":           0,
	"return":         1,
	"adjustment":     2,
	"insurance":      3,
	"tax":            4,
	"other":          5,
	"cods":           6,
	"address_change": 7,
	"discount":       8,
	"unknown":        127,
}

func (x ShippingFeeType) Enum() *ShippingFeeType {
	p := new(ShippingFeeType)
	*p = x
	return p
}

func (x ShippingFeeType) String() string {
	return proto.EnumName(ShippingFeeType_name, int32(x))
}

func (x *ShippingFeeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ShippingFeeType_value, data, "ShippingFeeType")
	if err != nil {
		return err
	}
	*x = ShippingFeeType(value)
	return nil
}

func (ShippingFeeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dc2a878b5abe185c, []int{0}
}

func init() {
	proto.RegisterEnum("shipping_fee_type.ShippingFeeType", ShippingFeeType_name, ShippingFeeType_value)
}

func init() {
	proto.RegisterFile("etop/etc/shipping_fee_type/type.proto", fileDescriptor_dc2a878b5abe185c)
}

var fileDescriptor_dc2a878b5abe185c = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0x3d, 0x4e, 0x03, 0x31,
	0x10, 0x46, 0x13, 0xf2, 0xb3, 0x9b, 0x01, 0xc2, 0x30, 0xa2, 0xa2, 0xa0, 0xa3, 0x41, 0x22, 0x2e,
	0xb8, 0x01, 0x05, 0x17, 0x80, 0x8a, 0x66, 0xe5, 0xd8, 0xc3, 0xee, 0x12, 0x65, 0x6c, 0xd9, 0x63,
	0x20, 0x15, 0x87, 0xa0, 0xe1, 0x7e, 0x5c, 0x04, 0x2d, 0xd0, 0xa1, 0x34, 0xa3, 0x91, 0xde, 0xd3,
	0x27, 0x3d, 0xb8, 0x64, 0x0d, 0xd1, 0xb0, 0x3a, 0x93, 0xbb, 0x3e, 0xc6, 0x5e, 0xda, 0xe6, 0x89,
	0xb9, 0xd1, 0x5d, 0x64, 0x33, 0x9c, 0x55, 0x4c, 0x41, 0x03, 0x9d, 0xfe, 0xa3, 0xe7, 0x67, 0x6d,
	0x68, 0xc3, 0x0f, 0x35, 0xc3, 0xf7, 0x2b, 0x5e, 0x7d, 0x8c, 0xe1, 0xe4, 0xfe, 0xcf, 0xbd, 0x63,
	0x7e, 0xd8, 0x45, 0xa6, 0x1a, 0xa6, 0x5b, 0xdb, 0x0b, 0x8e, 0x08, 0x60, 0x9e, 0x58, 0x4b, 0x12,
	0x1c, 0xd3, 0x12, 0xc0, 0xfa, 0xe7, 0x92, 0x75, 0xcb, 0xa2, 0x78, 0x40, 0xc7, 0xb0, 0xe8, 0x25,
	0x97, 0x64, 0xc5, 0x31, 0x4e, 0xa8, 0x82, 0x89, 0xda, 0x37, 0x9c, 0xd2, 0x02, 0x66, 0x41, 0x3b,
	0x4e, 0x38, 0x1b, 0x86, 0x5c, 0xf0, 0x19, 0xe7, 0x44, 0xb0, 0xb4, 0xde, 0x27, 0xce, 0xb9, 0x71,
	0x9d, 0x95, 0x96, 0xb1, 0xa2, 0x23, 0xa8, 0x7d, 0x9f, 0x5d, 0x28, 0xa2, 0x58, 0xd3, 0x21, 0x54,
	0x45, 0x36, 0x12, 0x5e, 0x05, 0xdf, 0x6f, 0x6f, 0x3e, 0xbf, 0x2e, 0x46, 0x8f, 0xd7, 0x43, 0xeb,
	0xea, 0x45, 0xcc, 0xda, 0xba, 0x0d, 0x8b, 0x37, 0x71, 0x6d, 0xf6, 0xe7, 0x7f, 0x07, 0x00, 0x00,
	0xff, 0xff, 0x28, 0x6e, 0xde, 0x21, 0x1b, 0x01, 0x00, 0x00,
}
