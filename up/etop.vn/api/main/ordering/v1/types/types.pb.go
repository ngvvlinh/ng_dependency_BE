// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/ordering/v1/types/types.proto

package types

import (
	fmt "fmt"
	math "math"

	types "etop.vn/api/main/catalog/v1/types"
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

func (m *ItemLine) Reset()         { *m = ItemLine{} }
func (m *ItemLine) String() string { return proto.CompactTextString(m) }
func (*ItemLine) ProtoMessage()    {}
func (*ItemLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_e54d1d17e4555a6f, []int{0}
}
func (m *ItemLine) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ItemLine.Unmarshal(m, b)
}
func (m *ItemLine) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ItemLine.Marshal(b, m, deterministic)
}
func (m *ItemLine) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemLine.Merge(m, src)
}
func (m *ItemLine) XXX_Size() int {
	return xxx_messageInfo_ItemLine.Size(m)
}
func (m *ItemLine) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemLine.DiscardUnknown(m)
}

var xxx_messageInfo_ItemLine proto.InternalMessageInfo

func (m *ItemLine) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *ItemLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *ItemLine) GetProductId() int64 {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *ItemLine) GetVariantId() int64 {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *ItemLine) GetIsOutside() bool {
	if m != nil {
		return m.IsOutside
	}
	return false
}

func (m *ItemLine) GetProductInfo() ProductInfo {
	if m != nil {
		return m.ProductInfo
	}
	return ProductInfo{}
}

func (m *ItemLine) GetTotalPrice() int32 {
	if m != nil {
		return m.TotalPrice
	}
	return 0
}

func (m *ProductInfo) Reset()         { *m = ProductInfo{} }
func (m *ProductInfo) String() string { return proto.CompactTextString(m) }
func (*ProductInfo) ProtoMessage()    {}
func (*ProductInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_e54d1d17e4555a6f, []int{1}
}
func (m *ProductInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductInfo.Unmarshal(m, b)
}
func (m *ProductInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductInfo.Marshal(b, m, deterministic)
}
func (m *ProductInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductInfo.Merge(m, src)
}
func (m *ProductInfo) XXX_Size() int {
	return xxx_messageInfo_ProductInfo.Size(m)
}
func (m *ProductInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ProductInfo proto.InternalMessageInfo

func (m *ProductInfo) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *ProductInfo) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *ProductInfo) GetAttributes() []*types.Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *ProductInfo) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_e54d1d17e4555a6f, []int{2}
}
func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Address) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Address) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Address) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *Address) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *Address) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *Location) Reset()         { *m = Location{} }
func (m *Location) String() string { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()    {}
func (*Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_e54d1d17e4555a6f, []int{3}
}
func (m *Location) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Location.Unmarshal(m, b)
}
func (m *Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Location.Marshal(b, m, deterministic)
}
func (m *Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Location.Merge(m, src)
}
func (m *Location) XXX_Size() int {
	return xxx_messageInfo_Location.Size(m)
}
func (m *Location) XXX_DiscardUnknown() {
	xxx_messageInfo_Location.DiscardUnknown(m)
}

var xxx_messageInfo_Location proto.InternalMessageInfo

func (m *Location) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *Location) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *Location) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *Location) GetCoordinates() *Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

func (m *Coordinates) Reset()         { *m = Coordinates{} }
func (m *Coordinates) String() string { return proto.CompactTextString(m) }
func (*Coordinates) ProtoMessage()    {}
func (*Coordinates) Descriptor() ([]byte, []int) {
	return fileDescriptor_e54d1d17e4555a6f, []int{4}
}
func (m *Coordinates) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Coordinates.Unmarshal(m, b)
}
func (m *Coordinates) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Coordinates.Marshal(b, m, deterministic)
}
func (m *Coordinates) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Coordinates.Merge(m, src)
}
func (m *Coordinates) XXX_Size() int {
	return xxx_messageInfo_Coordinates.Size(m)
}
func (m *Coordinates) XXX_DiscardUnknown() {
	xxx_messageInfo_Coordinates.DiscardUnknown(m)
}

var xxx_messageInfo_Coordinates proto.InternalMessageInfo

func (m *Coordinates) GetLatitude() float32 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *Coordinates) GetLongitude() float32 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func init() {
	proto.RegisterType((*ItemLine)(nil), "etop.vn.api.main.ordering.v1.types.ItemLine")
	proto.RegisterType((*ProductInfo)(nil), "etop.vn.api.main.ordering.v1.types.ProductInfo")
	proto.RegisterType((*Address)(nil), "etop.vn.api.main.ordering.v1.types.Address")
	proto.RegisterType((*Location)(nil), "etop.vn.api.main.ordering.v1.types.Location")
	proto.RegisterType((*Coordinates)(nil), "etop.vn.api.main.ordering.v1.types.Coordinates")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/ordering/v1/types/types.proto", fileDescriptor_e54d1d17e4555a6f)
}

var fileDescriptor_e54d1d17e4555a6f = []byte{
	// 660 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x9d, 0xfe, 0xa1, 0xa9, 0x33, 0x2e, 0x16, 0x87, 0x6a, 0x12, 0x5e, 0x57, 0x84, 0x28,
	0x12, 0x24, 0x5a, 0xbf, 0x00, 0xda, 0x76, 0x9a, 0x34, 0xc1, 0x28, 0x42, 0x42, 0x5c, 0x2a, 0x13,
	0x7b, 0xc5, 0x52, 0x62, 0x07, 0xc7, 0x29, 0xda, 0x37, 0xe0, 0xc8, 0x91, 0x03, 0x1f, 0x80, 0x8f,
	0xc0, 0x85, 0xfb, 0x0e, 0x1c, 0x7a, 0x42, 0x9c, 0xa6, 0x35, 0xfd, 0x02, 0x7c, 0x04, 0x64, 0x27,
	0x69, 0xd3, 0x75, 0x12, 0x5c, 0xa2, 0xe8, 0x79, 0x7f, 0x7e, 0xfd, 0x3c, 0xaf, 0x5f, 0xe8, 0x33,
	0x2d, 0x13, 0x7f, 0x26, 0x02, 0x92, 0xf0, 0x20, 0x26, 0x5c, 0x04, 0x52, 0x51, 0xa6, 0xb8, 0x98,
	0x06, 0xb3, 0x83, 0x40, 0x5f, 0x24, 0x2c, 0x2d, 0xbe, 0x7e, 0xa2, 0xa4, 0x96, 0x68, 0x50, 0xf2,
	0x3e, 0x49, 0xb8, 0x6f, 0x78, 0xbf, 0xe2, 0xfd, 0xd9, 0x81, 0x6f, 0xc9, 0xdd, 0x7b, 0x53, 0x39,
	0x95, 0x16, 0x0f, 0xcc, 0x5f, 0x71, 0x72, 0xf7, 0xfe, 0xc6, 0x4d, 0x4c, 0x13, 0x73, 0x81, 0xe9,
	0x52, 0x94, 0x9f, 0x6e, 0x19, 0x09, 0x89, 0x26, 0x91, 0xbc, 0xdd, 0xc7, 0xe0, 0x47, 0x03, 0xba,
	0x27, 0x9a, 0xc5, 0xa7, 0x5c, 0x30, 0xb4, 0x07, 0x5d, 0xeb, 0x62, 0xc2, 0x69, 0xcf, 0xe9, 0x3b,
	0xc3, 0xe6, 0x51, 0xeb, 0xf2, 0x6a, 0x0f, 0x8c, 0x3b, 0x56, 0x3d, 0xa1, 0xa8, 0x0f, 0xdd, 0x0f,
	0x19, 0x11, 0x9a, 0xeb, 0x8b, 0x5e, 0xa3, 0xef, 0x0c, 0xdb, 0x25, 0xb0, 0x52, 0xd1, 0x03, 0x08,
	0x13, 0x25, 0x69, 0x16, 0x6a, 0xd3, 0xa4, 0x59, 0x6b, 0xd2, 0x2d, 0xf5, 0x13, 0x6a, 0xa0, 0x19,
	0x51, 0x9c, 0x08, 0x0b, 0xb5, 0xea, 0x50, 0xa9, 0x17, 0x10, 0x4f, 0x27, 0x32, 0xd3, 0x29, 0xa7,
	0xac, 0xd7, 0xee, 0x3b, 0x43, 0xb7, 0x82, 0x78, 0xfa, 0xa2, 0x90, 0xd1, 0x1b, 0xb8, 0xb3, 0xba,
	0x4e, 0x9c, 0xcb, 0xde, 0x9d, 0xbe, 0x33, 0xf4, 0x46, 0x81, 0xff, 0xef, 0xe9, 0xfa, 0x67, 0xa5,
	0x1d, 0x71, 0x2e, 0xcb, 0xbe, 0x5e, 0xb2, 0x96, 0xd0, 0x43, 0xe8, 0x69, 0xa9, 0x49, 0x34, 0x49,
	0x14, 0x0f, 0x59, 0xaf, 0x53, 0x4b, 0x0b, 0x6d, 0xe1, 0xcc, 0xe8, 0x83, 0x9f, 0x0e, 0xf4, 0x6a,
	0x9d, 0xd0, 0xa3, 0xb5, 0x21, 0x41, 0x62, 0x66, 0xc7, 0xd8, 0xbd, 0xd1, 0xff, 0x39, 0x89, 0x19,
	0xda, 0x87, 0x5d, 0x1e, 0x93, 0x29, 0x9b, 0x64, 0x2a, 0xb2, 0xb3, 0xac, 0x28, 0xd7, 0xca, 0xaf,
	0x55, 0x84, 0x4e, 0x21, 0x24, 0x5a, 0x2b, 0xfe, 0x2e, 0xd3, 0x2c, 0xed, 0x35, 0xfb, 0xcd, 0xa1,
	0x37, 0x7a, 0xb2, 0x1d, 0xad, 0x7c, 0xdf, 0x75, 0xb2, 0xc3, 0xea, 0xd0, 0xb8, 0x76, 0xde, 0xcc,
	0x33, 0xe2, 0xa9, 0x2e, 0xf3, 0xb4, 0x6a, 0x79, 0xba, 0x46, 0x2f, 0xe2, 0x7c, 0x6d, 0xc0, 0xce,
	0x21, 0xa5, 0x8a, 0xa5, 0xa9, 0x71, 0x78, 0x9e, 0x45, 0xd1, 0x76, 0x0e, 0xd7, 0xc8, 0x36, 0xc4,
	0x2e, 0x6c, 0x27, 0xef, 0xa5, 0x60, 0x1b, 0x01, 0x0a, 0xc9, 0xd4, 0x58, 0x4c, 0x78, 0x64, 0x97,
	0x60, 0x55, 0xb3, 0x12, 0xc2, 0xb0, 0x13, 0xca, 0x38, 0x21, 0xe2, 0xc2, 0x1a, 0xa9, 0xaa, 0x95,
	0x68, 0xf6, 0x8c, 0x14, 0x2e, 0x0e, 0xec, 0xcb, 0xaf, 0x6e, 0xae, 0xd4, 0x1a, 0x31, 0xb2, 0x8f,
	0x7e, 0x93, 0x18, 0xa1, 0x31, 0x74, 0x23, 0x19, 0x12, 0xcd, 0xa5, 0xb0, 0xaf, 0x77, 0xeb, 0xec,
	0xb6, 0xd7, 0xe2, 0xb4, 0x3c, 0x73, 0xe4, 0x9a, 0x7e, 0xf3, 0xab, 0x3d, 0x67, 0xbc, 0xea, 0x33,
	0xf8, 0xe5, 0x40, 0xb7, 0x02, 0xd0, 0x63, 0x78, 0x37, 0x51, 0x72, 0xc6, 0x45, 0xc8, 0x26, 0xa1,
	0xa4, 0x9b, 0x33, 0xda, 0xa9, 0x4a, 0xc7, 0x92, 0x32, 0x83, 0x52, 0x9e, 0x6a, 0xc5, 0x43, 0x5d,
	0xa0, 0xf5, 0x79, 0xed, 0x54, 0x25, 0x8b, 0xee, 0xc3, 0xee, 0x47, 0xa2, 0x68, 0x81, 0xd5, 0x47,
	0xe7, 0x1a, 0xd9, 0x22, 0x2f, 0xa1, 0x17, 0x4a, 0xa9, 0x28, 0x17, 0xc4, 0x2c, 0x46, 0xeb, 0xff,
	0x77, 0xfe, 0x78, 0x7d, 0x6c, 0x5c, 0xef, 0x31, 0x78, 0x05, 0xbd, 0x5a, 0xcd, 0x4c, 0x37, 0x22,
	0x9a, 0xeb, 0xac, 0x4c, 0xd5, 0xa8, 0x3c, 0x54, 0x2a, 0x1a, 0xc0, 0x6e, 0x24, 0xc5, 0xb4, 0x40,
	0x1a, 0x35, 0x64, 0x2d, 0x1f, 0x3d, 0xbb, 0x5c, 0x60, 0x67, 0xbe, 0xc0, 0xe0, 0x7a, 0x81, 0xc1,
	0x9f, 0x05, 0x06, 0x9f, 0x72, 0x0c, 0xbe, 0xe5, 0x18, 0x7c, 0xcf, 0x31, 0xb8, 0xcc, 0x31, 0x98,
	0xe7, 0x18, 0x5c, 0xe7, 0x18, 0x7c, 0x5e, 0x62, 0xf0, 0x65, 0x89, 0xc1, 0x7c, 0x89, 0xc1, 0xef,
	0x25, 0x06, 0x6f, 0xdb, 0xd6, 0xee, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x10, 0x98, 0xf8, 0x02,
	0x55, 0x05, 0x00, 0x00,
}
