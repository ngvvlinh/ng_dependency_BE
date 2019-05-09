// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/product/v1/types/types.proto

package types

import (
	fmt "fmt"
	math "math"

	_ "etop.vn/api/meta/v1"
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

func (m *Attribute) Reset()         { *m = Attribute{} }
func (m *Attribute) String() string { return proto.CompactTextString(m) }
func (*Attribute) ProtoMessage()    {}
func (*Attribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_0864ca9e661cc3cd, []int{0}
}
func (m *Attribute) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Attribute.Unmarshal(m, b)
}
func (m *Attribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Attribute.Marshal(b, m, deterministic)
}
func (m *Attribute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attribute.Merge(m, src)
}
func (m *Attribute) XXX_Size() int {
	return xxx_messageInfo_Attribute.Size(m)
}
func (m *Attribute) XXX_DiscardUnknown() {
	xxx_messageInfo_Attribute.DiscardUnknown(m)
}

var xxx_messageInfo_Attribute proto.InternalMessageInfo

func (m *Attribute) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*Attribute)(nil), "etop.vn.api.main.product.v1.types.Attribute")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/product/v1/types/types.proto", fileDescriptor_0864ca9e661cc3cd)
}

var fileDescriptor_0864ca9e661cc3cd = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4d, 0x2d, 0xc9, 0x2f,
	0xd0, 0x2b, 0xcb, 0xd3, 0x4f, 0x2c, 0xc8, 0xd4, 0xcf, 0x4d, 0xcc, 0xcc, 0xd3, 0x2f, 0x28, 0xca,
	0x4f, 0x29, 0x4d, 0x2e, 0xd1, 0x2f, 0x33, 0xd4, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0x86, 0x90, 0x7a,
	0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x8a, 0x50, 0xe5, 0x7a, 0x89, 0x05, 0x99, 0x7a, 0x20, 0xe5,
	0x7a, 0x50, 0xe5, 0x7a, 0x65, 0x86, 0x7a, 0x60, 0x85, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60,
	0xd5, 0xfa, 0x20, 0x16, 0x44, 0xa3, 0x94, 0x2c, 0x8a, 0x3d, 0xa9, 0x25, 0x89, 0x20, 0xf3, 0x41,
	0x86, 0x80, 0xa5, 0x95, 0x1c, 0xb9, 0x38, 0x1d, 0x4b, 0x4a, 0x8a, 0x32, 0x93, 0x4a, 0x4b, 0x52,
	0x85, 0x24, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x9d, 0x58,
	0x4e, 0xdc, 0x93, 0x67, 0x08, 0x02, 0x8b, 0x08, 0x49, 0x71, 0xb1, 0x96, 0x25, 0xe6, 0x94, 0xa6,
	0x4a, 0x30, 0x21, 0x49, 0x41, 0x84, 0x9c, 0xec, 0x4f, 0x3c, 0x94, 0x63, 0xbc, 0xf0, 0x50, 0x8e,
	0xe1, 0xc1, 0x43, 0x39, 0x86, 0x0f, 0x0f, 0xe5, 0x18, 0x3a, 0x1e, 0xc9, 0x31, 0xac, 0x78, 0x24,
	0xc7, 0xb0, 0xe3, 0x91, 0x1c, 0xc3, 0x89, 0x47, 0x72, 0x0c, 0x17, 0x1e, 0xc9, 0x31, 0x3c, 0x78,
	0x24, 0xc7, 0x30, 0xe1, 0xb1, 0x1c, 0xc3, 0x8c, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc,
	0x78, 0x2c, 0xc7, 0x10, 0xc5, 0x0a, 0x76, 0x38, 0x20, 0x00, 0x00, 0xff, 0xff, 0x16, 0x4b, 0xe4,
	0x45, 0x0b, 0x01, 0x00, 0x00,
}
