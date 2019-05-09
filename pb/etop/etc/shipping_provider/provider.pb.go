// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/etc/shipping_provider/provider.proto

package shipping_provider

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

type ShippingProvider int32

const (
	ShippingProvider_unknown ShippingProvider = 0
	ShippingProvider_all     ShippingProvider = 22
	ShippingProvider_manual  ShippingProvider = 20
	ShippingProvider_ghn     ShippingProvider = 19
	ShippingProvider_ghtk    ShippingProvider = 21
	ShippingProvider_vtpost  ShippingProvider = 23
)

var ShippingProvider_name = map[int32]string{
	0:  "unknown",
	22: "all",
	20: "manual",
	19: "ghn",
	21: "ghtk",
	23: "vtpost",
}

var ShippingProvider_value = map[string]int32{
	"unknown": 0,
	"all":     22,
	"manual":  20,
	"ghn":     19,
	"ghtk":    21,
	"vtpost":  23,
}

func (x ShippingProvider) Enum() *ShippingProvider {
	p := new(ShippingProvider)
	*p = x
	return p
}

func (x ShippingProvider) String() string {
	return proto.EnumName(ShippingProvider_name, int32(x))
}

func (x *ShippingProvider) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ShippingProvider_value, data, "ShippingProvider")
	if err != nil {
		return err
	}
	*x = ShippingProvider(value)
	return nil
}

func (ShippingProvider) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_29a71c0b157e22e7, []int{0}
}

func init() {
	proto.RegisterEnum("shipping_provider.ShippingProvider", ShippingProvider_name, ShippingProvider_value)
}

func init() {
	proto.RegisterFile("etop/etc/shipping_provider/provider.proto", fileDescriptor_29a71c0b157e22e7)
}

var fileDescriptor_29a71c0b157e22e7 = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2d, 0xc9, 0x2f,
	0xd0, 0x4f, 0x2d, 0x49, 0xd6, 0x2f, 0xce, 0xc8, 0x2c, 0x28, 0xc8, 0xcc, 0x4b, 0x8f, 0x2f, 0x28,
	0xca, 0x2f, 0xcb, 0x4c, 0x49, 0x2d, 0xd2, 0x87, 0x31, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85,
	0x04, 0x31, 0x54, 0x48, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x65, 0xf5, 0x41, 0x2c, 0x88, 0x42,
	0xad, 0x60, 0x2e, 0x81, 0x60, 0xa8, 0xd2, 0x00, 0xa8, 0x4a, 0x21, 0x6e, 0x2e, 0xf6, 0xd2, 0xbc,
	0xec, 0xbc, 0xfc, 0xf2, 0x3c, 0x01, 0x06, 0x21, 0x76, 0x2e, 0xe6, 0xc4, 0x9c, 0x1c, 0x01, 0x31,
	0x21, 0x2e, 0x2e, 0xb6, 0xdc, 0xc4, 0xbc, 0xd2, 0xc4, 0x1c, 0x01, 0x11, 0x90, 0x60, 0x7a, 0x46,
	0x9e, 0x80, 0xb0, 0x10, 0x07, 0x17, 0x4b, 0x7a, 0x46, 0x49, 0xb6, 0x80, 0x28, 0x48, 0xba, 0xac,
	0xa4, 0x20, 0xbf, 0xb8, 0x44, 0x40, 0xdc, 0xc9, 0x78, 0xc6, 0x63, 0x39, 0x86, 0x28, 0x5d, 0x90,
	0x73, 0xf5, 0xca, 0xf2, 0xf4, 0x93, 0x12, 0x93, 0xb3, 0x53, 0xf3, 0x52, 0xf4, 0x0b, 0x92, 0xf4,
	0x71, 0xfb, 0x00, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe7, 0xd4, 0xa5, 0x97, 0xde, 0x00, 0x00, 0x00,
}
