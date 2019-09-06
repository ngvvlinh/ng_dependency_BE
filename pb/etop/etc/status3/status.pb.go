// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/etc/status3/status.proto

package status3

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

type Status int32

const (
	Status_Z Status = 0
	Status_P Status = 1
	Status_N Status = 127
)

var Status_name = map[int32]string{
	0:   "Z",
	1:   "P",
	127: "N",
}

var Status_value = map[string]int32{
	"Z": 0,
	"P": 1,
	"N": 127,
}

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (x *Status) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Status_value, data, "Status")
	if err != nil {
		return err
	}
	*x = Status(value)
	return nil
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d94e5fa78a06ad34, []int{0}
}

func init() {
	proto.RegisterEnum("status3.Status", Status_name, Status_value)
}

func init() { proto.RegisterFile("etop/etc/status3/status.proto", fileDescriptor_d94e5fa78a06ad34) }

var fileDescriptor_d94e5fa78a06ad34 = []byte{
	// 117 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0x2d, 0xc9, 0x2f,
	0xd0, 0x4f, 0x2d, 0x49, 0xd6, 0x2f, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0x36, 0x86, 0xd2, 0x7a, 0x05,
	0x45, 0xf9, 0x25, 0xf9, 0x42, 0xec, 0x50, 0x51, 0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0, 0x98,
	0x3e, 0x88, 0x05, 0x91, 0xd6, 0x92, 0xe5, 0x62, 0x0b, 0x06, 0x2b, 0x10, 0x62, 0xe5, 0x62, 0x8c,
	0x12, 0x60, 0x00, 0x51, 0x01, 0x02, 0x8c, 0x20, 0xca, 0x4f, 0xa0, 0xde, 0x49, 0x73, 0xc6, 0x63,
	0x39, 0x86, 0x28, 0x65, 0x90, 0x15, 0x7a, 0x65, 0x79, 0xfa, 0x49, 0x89, 0xc9, 0xd9, 0xa9, 0x79,
	0x29, 0xfa, 0x05, 0x49, 0xfa, 0xe8, 0xb6, 0x02, 0x02, 0x00, 0x00, 0xff, 0xff, 0x93, 0xfd, 0x95,
	0xb2, 0x88, 0x00, 0x00, 0x00,
}
