// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/etc/gender/gender.proto

package gender

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

type Gender int32

const (
	Gender_unknown Gender = 0
	Gender_male    Gender = 1
	Gender_female  Gender = 2
	Gender_other   Gender = 3
)

var Gender_name = map[int32]string{
	0: "unknown",
	1: "male",
	2: "female",
	3: "other",
}

var Gender_value = map[string]int32{
	"unknown": 0,
	"male":    1,
	"female":  2,
	"other":   3,
}

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) String() string {
	return proto.EnumName(Gender_name, int32(x))
}

func (x *Gender) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Gender_value, data, "Gender")
	if err != nil {
		return err
	}
	*x = Gender(value)
	return nil
}

func (Gender) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8dbfda7bf0b11a5a, []int{0}
}

func init() {
	proto.RegisterEnum("gender.Gender", Gender_name, Gender_value)
}

func init() { proto.RegisterFile("etop/etc/gender/gender.proto", fileDescriptor_8dbfda7bf0b11a5a) }

var fileDescriptor_8dbfda7bf0b11a5a = []byte{
	// 138 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x49, 0x2d, 0xc9, 0x2f,
	0xd0, 0x4f, 0x2d, 0x49, 0xd6, 0x4f, 0x4f, 0xcd, 0x4b, 0x49, 0x2d, 0x82, 0x52, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58, 0x48, 0x1f,
	0xc4, 0x82, 0xc8, 0x6a, 0x99, 0x71, 0xb1, 0xb9, 0x83, 0xe5, 0x85, 0xb8, 0xb9, 0xd8, 0x4b, 0xf3,
	0xb2, 0xf3, 0xf2, 0xcb, 0xf3, 0x04, 0x18, 0x84, 0x38, 0xb8, 0x58, 0x72, 0x13, 0x73, 0x52, 0x05,
	0x18, 0x85, 0xb8, 0xb8, 0xd8, 0xd2, 0x52, 0xc1, 0x6c, 0x26, 0x21, 0x4e, 0x2e, 0xd6, 0xfc, 0x92,
	0x8c, 0xd4, 0x22, 0x01, 0x66, 0x27, 0x95, 0x19, 0x8f, 0xe5, 0x18, 0xa2, 0xe4, 0x40, 0x36, 0xeb,
	0x95, 0xe5, 0xe9, 0x27, 0x16, 0x64, 0xea, 0x17, 0x24, 0xe9, 0xa3, 0x39, 0x04, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x36, 0x03, 0x44, 0xba, 0x9a, 0x00, 0x00, 0x00,
}