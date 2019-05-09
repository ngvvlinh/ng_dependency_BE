// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/shipping/v1/api.proto

package v1

import (
	fmt "fmt"
	math "math"

	types "etop.vn/api/main/order/v1/types"
	_ "etop.vn/api/main/shipping/v1/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
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

func (m *Fulfillment) Reset()         { *m = Fulfillment{} }
func (m *Fulfillment) String() string { return proto.CompactTextString(m) }
func (*Fulfillment) ProtoMessage()    {}
func (*Fulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_461a23e35d98b2d6, []int{0}
}
func (m *Fulfillment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Fulfillment.Unmarshal(m, b)
}
func (m *Fulfillment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Fulfillment.Marshal(b, m, deterministic)
}
func (m *Fulfillment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Fulfillment.Merge(m, src)
}
func (m *Fulfillment) XXX_Size() int {
	return xxx_messageInfo_Fulfillment.Size(m)
}
func (m *Fulfillment) XXX_DiscardUnknown() {
	xxx_messageInfo_Fulfillment.DiscardUnknown(m)
}

var xxx_messageInfo_Fulfillment proto.InternalMessageInfo

func (m *Fulfillment) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Fulfillment) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *Fulfillment) GetShopId() int64 {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Fulfillment) GetSupplierId() int64 {
	if m != nil {
		return m.SupplierId
	}
	return 0
}

func (m *Fulfillment) GetPartnerId() int64 {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *Fulfillment) GetLines() types.ItemLine {
	if m != nil {
		return m.Lines
	}
	return types.ItemLine{}
}

func init() {
	proto.RegisterType((*Fulfillment)(nil), "etop.vn.api.main.shipping.v1.Fulfillment")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/shipping/v1/api.proto", fileDescriptor_461a23e35d98b2d6)
}

var fileDescriptor_461a23e35d98b2d6 = []byte{
	// 370 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xcd, 0x4a, 0xeb, 0x40,
	0x18, 0x86, 0x27, 0xe9, 0xef, 0x99, 0xc2, 0xe1, 0x30, 0x74, 0x51, 0xca, 0x71, 0x5a, 0x14, 0xa5,
	0xa2, 0x9d, 0x90, 0xde, 0x81, 0x05, 0x85, 0x80, 0x20, 0xe8, 0xce, 0x8d, 0x04, 0x67, 0x4c, 0x07,
	0xd2, 0x99, 0x21, 0x99, 0x46, 0xbc, 0x03, 0x97, 0xe2, 0xca, 0x4b, 0xf0, 0x12, 0xbc, 0x84, 0x2e,
	0xbb, 0x74, 0x25, 0x4d, 0x72, 0x03, 0xae, 0x5c, 0x4b, 0x7e, 0x0a, 0x81, 0x8a, 0x9b, 0xf0, 0xf1,
	0xbc, 0x4f, 0x5e, 0xbe, 0x99, 0x81, 0x07, 0x4c, 0x4b, 0x45, 0x22, 0x61, 0xb9, 0x8a, 0x5b, 0x73,
	0x97, 0x0b, 0x2b, 0x9c, 0x71, 0xa5, 0xb8, 0xf0, 0xac, 0xc8, 0xce, 0x20, 0x51, 0x81, 0xd4, 0x12,
	0xfd, 0x2f, 0x3d, 0x92, 0xa1, 0xcc, 0x23, 0x1b, 0x8f, 0x44, 0x76, 0xff, 0x38, 0x97, 0x6e, 0xc7,
	0x1e, 0x13, 0xe3, 0xf0, 0xde, 0xf5, 0x3c, 0x16, 0x58, 0x52, 0x69, 0x2e, 0x45, 0x68, 0xb9, 0x42,
	0x48, 0xed, 0xe6, 0x73, 0xd1, 0xd5, 0xef, 0x7a, 0xd2, 0x93, 0xf9, 0x68, 0x65, 0x53, 0x49, 0xc9,
	0xaf, 0x9b, 0xe8, 0x07, 0xc5, 0xc2, 0xe2, 0x5b, 0xfa, 0x47, 0x5b, 0xbe, 0x0c, 0x28, 0x0b, 0x7e,
	0x94, 0x77, 0xbf, 0x0c, 0xd8, 0x39, 0x5b, 0xf8, 0x77, 0xdc, 0xf7, 0xe7, 0x4c, 0x68, 0xd4, 0x85,
	0x26, 0xa7, 0x3d, 0x63, 0x68, 0x8c, 0x6a, 0xd3, 0xfa, 0xf2, 0x63, 0x00, 0x2e, 0x4d, 0x4e, 0xd1,
	0x00, 0xb6, 0xf3, 0x8e, 0x1b, 0x4e, 0x7b, 0x66, 0x25, 0x6b, 0xe5, 0xd4, 0xa1, 0x68, 0x07, 0xb6,
	0xc2, 0x99, 0x54, 0x59, 0x5e, 0xab, 0xe4, 0xcd, 0x0c, 0x3a, 0x14, 0xed, 0xc3, 0x4e, 0xb8, 0x50,
	0xca, 0xe7, 0x45, 0x45, 0xbd, 0xa2, 0xc0, 0x4d, 0xe0, 0x50, 0xb4, 0x07, 0xa1, 0x72, 0x03, 0x2d,
	0x0a, 0xab, 0x51, 0xb1, 0xfe, 0x94, 0xdc, 0xa1, 0xe8, 0x14, 0x36, 0x7c, 0x2e, 0x58, 0xd8, 0x6b,
	0x0e, 0x8d, 0x51, 0x67, 0x72, 0x48, 0xb6, 0x1e, 0x20, 0x5f, 0x8a, 0x44, 0x36, 0x29, 0x0e, 0xea,
	0x68, 0x36, 0x3f, 0xe7, 0x82, 0x95, 0x55, 0xc5, 0xdf, 0xd3, 0x8b, 0xe7, 0x93, 0x7f, 0xe8, 0x2f,
	0x6c, 0x5f, 0x95, 0x77, 0x39, 0x31, 0x23, 0x7b, 0x19, 0x63, 0x63, 0x15, 0x63, 0xb0, 0x8e, 0x31,
	0xf8, 0x8c, 0x31, 0x78, 0x4c, 0x30, 0x78, 0x4d, 0x30, 0x78, 0x4b, 0x30, 0x58, 0x26, 0x18, 0xac,
	0x12, 0x0c, 0xd6, 0x09, 0x06, 0x4f, 0x29, 0x06, 0x2f, 0x29, 0x06, 0xab, 0x14, 0x83, 0xf7, 0x14,
	0x83, 0x6b, 0x33, 0xb2, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x84, 0x80, 0xd3, 0xd8, 0x31, 0x02,
	0x00, 0x00,
}
