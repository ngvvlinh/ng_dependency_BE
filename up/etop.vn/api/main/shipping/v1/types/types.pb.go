// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/shipping/v1/types/types.proto

package types

import (
	fmt "fmt"
	math "math"

	v1 "etop.vn/api/meta/v1"
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

func (m *ShippingService) Reset()         { *m = ShippingService{} }
func (m *ShippingService) String() string { return proto.CompactTextString(m) }
func (*ShippingService) ProtoMessage()    {}
func (*ShippingService) Descriptor() ([]byte, []int) {
	return fileDescriptor_e033cb1613705d25, []int{0}
}
func (m *ShippingService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShippingService.Unmarshal(m, b)
}
func (m *ShippingService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShippingService.Marshal(b, m, deterministic)
}
func (m *ShippingService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShippingService.Merge(m, src)
}
func (m *ShippingService) XXX_Size() int {
	return xxx_messageInfo_ShippingService.Size(m)
}
func (m *ShippingService) XXX_DiscardUnknown() {
	xxx_messageInfo_ShippingService.DiscardUnknown(m)
}

var xxx_messageInfo_ShippingService proto.InternalMessageInfo

func (m *ShippingService) GetCarrier() string {
	if m != nil {
		return m.Carrier
	}
	return ""
}

func (m *ShippingService) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ShippingService) GetFee() int32 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func (m *ShippingService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShippingService) GetEstimatedPickupAt() *v1.Timestamp {
	if m != nil {
		return m.EstimatedPickupAt
	}
	return nil
}

func (m *ShippingService) GetEstimatedDeliveryAt() *v1.Timestamp {
	if m != nil {
		return m.EstimatedDeliveryAt
	}
	return nil
}

func (m *WeightInfo) Reset()         { *m = WeightInfo{} }
func (m *WeightInfo) String() string { return proto.CompactTextString(m) }
func (*WeightInfo) ProtoMessage()    {}
func (*WeightInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_e033cb1613705d25, []int{1}
}
func (m *WeightInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeightInfo.Unmarshal(m, b)
}
func (m *WeightInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeightInfo.Marshal(b, m, deterministic)
}
func (m *WeightInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightInfo.Merge(m, src)
}
func (m *WeightInfo) XXX_Size() int {
	return xxx_messageInfo_WeightInfo.Size(m)
}
func (m *WeightInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightInfo.DiscardUnknown(m)
}

var xxx_messageInfo_WeightInfo proto.InternalMessageInfo

func (m *WeightInfo) GetGrossWeight() int32 {
	if m != nil {
		return m.GrossWeight
	}
	return 0
}

func (m *WeightInfo) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *WeightInfo) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *WeightInfo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *WeightInfo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ValueInfo) Reset()         { *m = ValueInfo{} }
func (m *ValueInfo) String() string { return proto.CompactTextString(m) }
func (*ValueInfo) ProtoMessage()    {}
func (*ValueInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_e033cb1613705d25, []int{2}
}
func (m *ValueInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValueInfo.Unmarshal(m, b)
}
func (m *ValueInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValueInfo.Marshal(b, m, deterministic)
}
func (m *ValueInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValueInfo.Merge(m, src)
}
func (m *ValueInfo) XXX_Size() int {
	return xxx_messageInfo_ValueInfo.Size(m)
}
func (m *ValueInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ValueInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ValueInfo proto.InternalMessageInfo

func (m *ValueInfo) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *ValueInfo) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *ValueInfo) GetIncludeInsurance() bool {
	if m != nil {
		return m.IncludeInsurance
	}
	return false
}

func init() {
	proto.RegisterType((*ShippingService)(nil), "etop.vn.api.main.shipping.v1.types.ShippingService")
	proto.RegisterType((*WeightInfo)(nil), "etop.vn.api.main.shipping.v1.types.WeightInfo")
	proto.RegisterType((*ValueInfo)(nil), "etop.vn.api.main.shipping.v1.types.ValueInfo")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/shipping/v1/types/types.proto", fileDescriptor_e033cb1613705d25)
}

var fileDescriptor_e033cb1613705d25 = []byte{
	// 492 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x3f, 0x8f, 0xd3, 0x30,
	0x18, 0xc6, 0x9d, 0x5e, 0x7b, 0x50, 0x1f, 0x12, 0x9c, 0xf9, 0xa3, 0xa8, 0x02, 0x53, 0x95, 0x81,
	0x4e, 0x89, 0xca, 0x17, 0x40, 0x3d, 0xb1, 0xdc, 0x82, 0x50, 0x0f, 0x81, 0xc4, 0x12, 0xf9, 0x92,
	0xf7, 0x12, 0xeb, 0x1a, 0xdb, 0x4a, 0x9c, 0x9c, 0xee, 0x1b, 0x30, 0x30, 0x30, 0xb2, 0xb2, 0xf1,
	0x11, 0x18, 0x19, 0x3b, 0x76, 0x64, 0x42, 0xd7, 0xf4, 0x0b, 0xf0, 0x11, 0x90, 0x9d, 0x84, 0x36,
	0xdb, 0x2d, 0x51, 0xf4, 0x7b, 0x9e, 0xe7, 0xb5, 0x9f, 0x57, 0xc6, 0x1e, 0x68, 0xa9, 0xbc, 0x52,
	0xf8, 0x4c, 0x71, 0x3f, 0x65, 0x5c, 0xf8, 0x79, 0xc2, 0x95, 0xe2, 0x22, 0xf6, 0xcb, 0x99, 0xaf,
	0xaf, 0x15, 0xe4, 0xf5, 0xd7, 0x53, 0x99, 0xd4, 0x92, 0x4c, 0x1a, 0xbf, 0xc7, 0x14, 0xf7, 0x8c,
	0xdf, 0x6b, 0xfd, 0x5e, 0x39, 0xf3, 0xac, 0x73, 0xf4, 0x28, 0x96, 0xb1, 0xb4, 0x76, 0xdf, 0xfc,
	0xd5, 0xc9, 0xd1, 0xb3, 0xce, 0x49, 0xa0, 0x99, 0x39, 0xc0, 0x4c, 0xb1, 0xf2, 0xe4, 0x7b, 0x0f,
	0xdf, 0x3f, 0x6b, 0x46, 0x9d, 0x41, 0x56, 0xf2, 0x10, 0x08, 0xc5, 0x77, 0x42, 0x96, 0x65, 0x1c,
	0x32, 0xd7, 0x19, 0x3b, 0xd3, 0xe1, 0x49, 0x7f, 0xf5, 0xe7, 0x39, 0x5a, 0xb4, 0x90, 0xb8, 0xb8,
	0x1f, 0xca, 0x08, 0xdc, 0xde, 0x9e, 0x68, 0x09, 0x79, 0x82, 0x0f, 0x2e, 0x00, 0xdc, 0x83, 0xb1,
	0x33, 0x1d, 0x34, 0x82, 0x01, 0x26, 0x21, 0x58, 0x0a, 0x6e, 0x7f, 0x3f, 0x61, 0x08, 0x79, 0x8b,
	0x1f, 0x42, 0xae, 0x79, 0xca, 0x34, 0x44, 0x81, 0xe2, 0xe1, 0x65, 0xa1, 0x02, 0xa6, 0xdd, 0xc1,
	0xd8, 0x99, 0x1e, 0xbd, 0xa2, 0x5e, 0xa7, 0x36, 0x68, 0x66, 0xda, 0xbe, 0xe7, 0x29, 0xe4, 0x9a,
	0xa5, 0x6a, 0x71, 0xfc, 0x3f, 0xfa, 0xce, 0x26, 0xe7, 0x9a, 0x2c, 0xf0, 0xe3, 0xdd, 0xbc, 0x08,
	0x96, 0xbc, 0x84, 0xec, 0xda, 0x4c, 0x3c, 0xbc, 0xd5, 0xc4, 0xdd, 0x65, 0xde, 0x34, 0xd9, 0xb9,
	0x9e, 0xfc, 0x72, 0x30, 0xfe, 0x08, 0x3c, 0x4e, 0xf4, 0xa9, 0xb8, 0x90, 0xe4, 0x25, 0xbe, 0x17,
	0x67, 0x32, 0xcf, 0x83, 0x2b, 0xcb, 0xec, 0x8e, 0xda, 0xb6, 0x47, 0x56, 0xa9, 0xcd, 0x64, 0x86,
	0x8f, 0xc3, 0x84, 0x65, 0x31, 0xb0, 0xf3, 0x25, 0xb4, 0xee, 0xde, 0x9e, 0xfb, 0xc1, 0x4e, 0x6e,
	0x22, 0x4f, 0xf1, 0xe1, 0x12, 0x44, 0xac, 0x93, 0xce, 0x0e, 0x1b, 0x46, 0x46, 0x78, 0x70, 0xc5,
	0x23, 0x9d, 0xd8, 0x3d, 0xb6, 0x62, 0x8d, 0x4c, 0x32, 0xa9, 0x4f, 0x18, 0xec, 0x27, 0x6b, 0x36,
	0xf9, 0xe2, 0xe0, 0xe1, 0x07, 0xb6, 0x2c, 0xa0, 0x6d, 0x70, 0xce, 0xf2, 0x4b, 0xd0, 0x41, 0x69,
	0x58, 0xb7, 0x41, 0xad, 0x58, 0x33, 0x79, 0x81, 0x71, 0x28, 0xa3, 0x80, 0xa5, 0xb2, 0x10, 0xdd,
	0xab, 0x0f, 0x43, 0x19, 0xcd, 0x2d, 0x36, 0x35, 0xb9, 0x08, 0x97, 0x45, 0x04, 0x01, 0x17, 0x79,
	0x91, 0x31, 0x11, 0xd6, 0x4f, 0xe0, 0x6e, 0x5b, 0xb3, 0x91, 0x4f, 0x5b, 0xf5, 0xe4, 0xf5, 0x6a,
	0x43, 0x9d, 0xf5, 0x86, 0xa2, 0x9b, 0x0d, 0x45, 0x7f, 0x37, 0x14, 0x7d, 0xae, 0x28, 0xfa, 0x51,
	0x51, 0xf4, 0xb3, 0xa2, 0x68, 0x55, 0x51, 0xb4, 0xae, 0x28, 0xba, 0xa9, 0x28, 0xfa, 0xba, 0xa5,
	0xe8, 0xdb, 0x96, 0xa2, 0xf5, 0x96, 0xa2, 0xdf, 0x5b, 0x8a, 0x3e, 0x0d, 0xec, 0x5b, 0xff, 0x17,
	0x00, 0x00, 0xff, 0xff, 0x90, 0x77, 0x2b, 0x47, 0x40, 0x03, 0x00, 0x00,
}
