// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/order/source/source.proto

package source

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

type Source int32

const (
	Source_unknown  Source = 0
	Source_self     Source = 1
	Source_import   Source = 2
	Source_api      Source = 3
	Source_etop_pos Source = 5
	Source_etop_pxs Source = 6
	Source_etop_cmx Source = 7
	Source_ts_app   Source = 8
)

var Source_name = map[int32]string{
	0: "unknown",
	1: "self",
	2: "import",
	3: "api",
	5: "etop_pos",
	6: "etop_pxs",
	7: "etop_cmx",
	8: "ts_app",
}

var Source_value = map[string]int32{
	"unknown":  0,
	"self":     1,
	"import":   2,
	"api":      3,
	"etop_pos": 5,
	"etop_pxs": 6,
	"etop_cmx": 7,
	"ts_app":   8,
}

func (x Source) Enum() *Source {
	p := new(Source)
	*p = x
	return p
}

func (x Source) String() string {
	return proto.EnumName(Source_name, int32(x))
}

func (x *Source) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Source_value, data, "Source")
	if err != nil {
		return err
	}
	*x = Source(value)
	return nil
}

func (Source) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_81452fd3d5f5edaa, []int{0}
}

func init() {
	proto.RegisterEnum("source.Source", Source_name, Source_value)
}

func init() { proto.RegisterFile("etop/order/source/source.proto", fileDescriptor_81452fd3d5f5edaa) }

var fileDescriptor_81452fd3d5f5edaa = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0x2d, 0xc9, 0x2f,
	0xd0, 0xcf, 0x2f, 0x4a, 0x49, 0x2d, 0xd2, 0x2f, 0xce, 0x2f, 0x2d, 0x4a, 0x4e, 0x85, 0x52, 0x7a,
	0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58,
	0x48, 0x1f, 0xc4, 0x82, 0xc8, 0x6a, 0x65, 0x71, 0xb1, 0x05, 0x83, 0xe5, 0x85, 0xb8, 0xb9, 0xd8,
	0x4b, 0xf3, 0xb2, 0xf3, 0xf2, 0xcb, 0xf3, 0x04, 0x18, 0x84, 0x38, 0xb8, 0x58, 0x8a, 0x53, 0x73,
	0xd2, 0x04, 0x18, 0x85, 0xb8, 0xb8, 0xd8, 0x32, 0x73, 0x0b, 0xf2, 0x8b, 0x4a, 0x04, 0x98, 0x84,
	0xd8, 0xb9, 0x98, 0x13, 0x0b, 0x32, 0x05, 0x98, 0x85, 0x78, 0xb8, 0x38, 0x40, 0xb6, 0xc6, 0x17,
	0xe4, 0x17, 0x0b, 0xb0, 0x22, 0x78, 0x15, 0xc5, 0x02, 0x6c, 0x70, 0x5e, 0x72, 0x6e, 0x85, 0x00,
	0x3b, 0x48, 0x7b, 0x49, 0x71, 0x7c, 0x62, 0x41, 0x81, 0x00, 0x87, 0x93, 0xd6, 0x8c, 0xc7, 0x72,
	0x0c, 0x51, 0x2a, 0x20, 0x59, 0xbd, 0xb2, 0x3c, 0xfd, 0xa4, 0xc4, 0xe4, 0xec, 0xd4, 0xbc, 0x14,
	0xfd, 0x82, 0x24, 0x7d, 0x0c, 0x2f, 0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3b, 0x85, 0xb2, 0xd6,
	0xd6, 0x00, 0x00, 0x00,
}
