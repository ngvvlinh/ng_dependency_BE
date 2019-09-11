// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/etc/payment_provider/payment_provider.proto

package payment_provider

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

type PaymentProvider int32

const (
	PaymentProvider_unknown PaymentProvider = 0
	PaymentProvider_vtpay   PaymentProvider = 1
)

var PaymentProvider_name = map[int32]string{
	0: "unknown",
	1: "vtpay",
}

var PaymentProvider_value = map[string]int32{
	"unknown": 0,
	"vtpay":   1,
}

func (x PaymentProvider) Enum() *PaymentProvider {
	p := new(PaymentProvider)
	*p = x
	return p
}

func (x PaymentProvider) String() string {
	return proto.EnumName(PaymentProvider_name, int32(x))
}

func (x *PaymentProvider) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PaymentProvider_value, data, "PaymentProvider")
	if err != nil {
		return err
	}
	*x = PaymentProvider(value)
	return nil
}

func (PaymentProvider) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_89ff7f3474d4522a, []int{0}
}

func init() {
	proto.RegisterEnum("payment_provider.PaymentProvider", PaymentProvider_name, PaymentProvider_value)
}

func init() {
	proto.RegisterFile("etop/etc/payment_provider/payment_provider.proto", fileDescriptor_89ff7f3474d4522a)
}

var fileDescriptor_89ff7f3474d4522a = []byte{
	// 138 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x48, 0x2d, 0xc9, 0x2f,
	0xd0, 0x4f, 0x2d, 0x49, 0xd6, 0x2f, 0x48, 0xac, 0xcc, 0x4d, 0xcd, 0x2b, 0x89, 0x2f, 0x28, 0xca,
	0x2f, 0xcb, 0x4c, 0x49, 0x2d, 0xc2, 0x10, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x40,
	0x17, 0x97, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x4b, 0xea, 0x83, 0x58, 0x10, 0x75, 0x5a, 0x9a,
	0x5c, 0xfc, 0x01, 0x10, 0x95, 0x01, 0x50, 0x85, 0x42, 0xdc, 0x5c, 0xec, 0xa5, 0x79, 0xd9, 0x79,
	0xf9, 0xe5, 0x79, 0x02, 0x0c, 0x42, 0x9c, 0x5c, 0xac, 0x65, 0x25, 0x05, 0x89, 0x95, 0x02, 0x8c,
	0x4e, 0x46, 0x33, 0x1e, 0xcb, 0x31, 0x44, 0xe9, 0x80, 0x9c, 0xa2, 0x57, 0x96, 0xa7, 0x9f, 0x94,
	0x98, 0x9c, 0x9d, 0x9a, 0x97, 0xa2, 0x5f, 0x90, 0xa4, 0x8f, 0xd3, 0x75, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xd9, 0xf0, 0xa8, 0xdf, 0xb9, 0x00, 0x00, 0x00,
}
