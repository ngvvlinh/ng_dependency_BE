// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/location/v1/api.proto

package v1

import (
	fmt "fmt"
	math "math"

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

func (m *Province) Reset()         { *m = Province{} }
func (m *Province) String() string { return proto.CompactTextString(m) }
func (*Province) ProtoMessage()    {}
func (*Province) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{0}
}
func (m *Province) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Province.Unmarshal(m, b)
}
func (m *Province) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Province.Marshal(b, m, deterministic)
}
func (m *Province) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Province.Merge(m, src)
}
func (m *Province) XXX_Size() int {
	return xxx_messageInfo_Province.Size(m)
}
func (m *Province) XXX_DiscardUnknown() {
	xxx_messageInfo_Province.DiscardUnknown(m)
}

var xxx_messageInfo_Province proto.InternalMessageInfo

func (m *Province) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Province) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Province) GetRegion() VietnamRegion {
	if m != nil {
		return m.Region
	}
	return VietnamRegion_unknown
}

func (m *District) Reset()         { *m = District{} }
func (m *District) String() string { return proto.CompactTextString(m) }
func (*District) ProtoMessage()    {}
func (*District) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{1}
}
func (m *District) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_District.Unmarshal(m, b)
}
func (m *District) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_District.Marshal(b, m, deterministic)
}
func (m *District) XXX_Merge(src proto.Message) {
	xxx_messageInfo_District.Merge(m, src)
}
func (m *District) XXX_Size() int {
	return xxx_messageInfo_District.Size(m)
}
func (m *District) XXX_DiscardUnknown() {
	xxx_messageInfo_District.DiscardUnknown(m)
}

var xxx_messageInfo_District proto.InternalMessageInfo

func (m *District) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *District) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *District) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *District) GetUrbanType() UrbanType {
	if m != nil {
		return m.UrbanType
	}
	return UrbanType_unknown
}

func (m *Ward) Reset()         { *m = Ward{} }
func (m *Ward) String() string { return proto.CompactTextString(m) }
func (*Ward) ProtoMessage()    {}
func (*Ward) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{2}
}
func (m *Ward) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ward.Unmarshal(m, b)
}
func (m *Ward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ward.Marshal(b, m, deterministic)
}
func (m *Ward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ward.Merge(m, src)
}
func (m *Ward) XXX_Size() int {
	return xxx_messageInfo_Ward.Size(m)
}
func (m *Ward) XXX_DiscardUnknown() {
	xxx_messageInfo_Ward.DiscardUnknown(m)
}

var xxx_messageInfo_Ward proto.InternalMessageInfo

func (m *Ward) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Ward) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Ward) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *Extra) Reset()         { *m = Extra{} }
func (m *Extra) String() string { return proto.CompactTextString(m) }
func (*Extra) ProtoMessage()    {}
func (*Extra) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{3}
}
func (m *Extra) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Extra.Unmarshal(m, b)
}
func (m *Extra) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Extra.Marshal(b, m, deterministic)
}
func (m *Extra) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Extra.Merge(m, src)
}
func (m *Extra) XXX_Size() int {
	return xxx_messageInfo_Extra.Size(m)
}
func (m *Extra) XXX_DiscardUnknown() {
	xxx_messageInfo_Extra.DiscardUnknown(m)
}

var xxx_messageInfo_Extra proto.InternalMessageInfo

func (m *Extra) GetSpecial() bool {
	if m != nil {
		return m.Special
	}
	return false
}

func (m *Extra) GetGhnId() int32 {
	if m != nil {
		return m.GhnId
	}
	return 0
}

func (m *Extra) GetVtpostId() int32 {
	if m != nil {
		return m.VtpostId
	}
	return 0
}

func (m *GetAllLocationsQueryArgs) Reset()         { *m = GetAllLocationsQueryArgs{} }
func (m *GetAllLocationsQueryArgs) String() string { return proto.CompactTextString(m) }
func (*GetAllLocationsQueryArgs) ProtoMessage()    {}
func (*GetAllLocationsQueryArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{4}
}
func (m *GetAllLocationsQueryArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAllLocationsQueryArgs.Unmarshal(m, b)
}
func (m *GetAllLocationsQueryArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAllLocationsQueryArgs.Marshal(b, m, deterministic)
}
func (m *GetAllLocationsQueryArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAllLocationsQueryArgs.Merge(m, src)
}
func (m *GetAllLocationsQueryArgs) XXX_Size() int {
	return xxx_messageInfo_GetAllLocationsQueryArgs.Size(m)
}
func (m *GetAllLocationsQueryArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAllLocationsQueryArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GetAllLocationsQueryArgs proto.InternalMessageInfo

func (m *GetAllLocationsQueryArgs) GetAll() bool {
	if m != nil {
		return m.All
	}
	return false
}

func (m *GetAllLocationsQueryArgs) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *GetAllLocationsQueryArgs) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *GetAllLocationsQueryResult) Reset()         { *m = GetAllLocationsQueryResult{} }
func (m *GetAllLocationsQueryResult) String() string { return proto.CompactTextString(m) }
func (*GetAllLocationsQueryResult) ProtoMessage()    {}
func (*GetAllLocationsQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{5}
}
func (m *GetAllLocationsQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAllLocationsQueryResult.Unmarshal(m, b)
}
func (m *GetAllLocationsQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAllLocationsQueryResult.Marshal(b, m, deterministic)
}
func (m *GetAllLocationsQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAllLocationsQueryResult.Merge(m, src)
}
func (m *GetAllLocationsQueryResult) XXX_Size() int {
	return xxx_messageInfo_GetAllLocationsQueryResult.Size(m)
}
func (m *GetAllLocationsQueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAllLocationsQueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_GetAllLocationsQueryResult proto.InternalMessageInfo

func (m *GetAllLocationsQueryResult) GetProvinces() []*Province {
	if m != nil {
		return m.Provinces
	}
	return nil
}

func (m *GetAllLocationsQueryResult) GetDistricts() []*District {
	if m != nil {
		return m.Districts
	}
	return nil
}

func (m *GetAllLocationsQueryResult) GetWards() []*Ward {
	if m != nil {
		return m.Wards
	}
	return nil
}

func (m *LocationQueryResult) Reset()         { *m = LocationQueryResult{} }
func (m *LocationQueryResult) String() string { return proto.CompactTextString(m) }
func (*LocationQueryResult) ProtoMessage()    {}
func (*LocationQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{6}
}
func (m *LocationQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocationQueryResult.Unmarshal(m, b)
}
func (m *LocationQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocationQueryResult.Marshal(b, m, deterministic)
}
func (m *LocationQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocationQueryResult.Merge(m, src)
}
func (m *LocationQueryResult) XXX_Size() int {
	return xxx_messageInfo_LocationQueryResult.Size(m)
}
func (m *LocationQueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_LocationQueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_LocationQueryResult proto.InternalMessageInfo

func (m *LocationQueryResult) GetProvince() *Province {
	if m != nil {
		return m.Province
	}
	return nil
}

func (m *LocationQueryResult) GetDistrict() *District {
	if m != nil {
		return m.District
	}
	return nil
}

func (m *LocationQueryResult) GetWard() *Ward {
	if m != nil {
		return m.Ward
	}
	return nil
}

func (m *GetLocationQueryArgs) Reset()         { *m = GetLocationQueryArgs{} }
func (m *GetLocationQueryArgs) String() string { return proto.CompactTextString(m) }
func (*GetLocationQueryArgs) ProtoMessage()    {}
func (*GetLocationQueryArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{7}
}
func (m *GetLocationQueryArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLocationQueryArgs.Unmarshal(m, b)
}
func (m *GetLocationQueryArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLocationQueryArgs.Marshal(b, m, deterministic)
}
func (m *GetLocationQueryArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLocationQueryArgs.Merge(m, src)
}
func (m *GetLocationQueryArgs) XXX_Size() int {
	return xxx_messageInfo_GetLocationQueryArgs.Size(m)
}
func (m *GetLocationQueryArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLocationQueryArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GetLocationQueryArgs proto.InternalMessageInfo

func (m *GetLocationQueryArgs) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *GetLocationQueryArgs) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *GetLocationQueryArgs) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *FindLocationQueryArgs) Reset()         { *m = FindLocationQueryArgs{} }
func (m *FindLocationQueryArgs) String() string { return proto.CompactTextString(m) }
func (*FindLocationQueryArgs) ProtoMessage()    {}
func (*FindLocationQueryArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{8}
}
func (m *FindLocationQueryArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindLocationQueryArgs.Unmarshal(m, b)
}
func (m *FindLocationQueryArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindLocationQueryArgs.Marshal(b, m, deterministic)
}
func (m *FindLocationQueryArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindLocationQueryArgs.Merge(m, src)
}
func (m *FindLocationQueryArgs) XXX_Size() int {
	return xxx_messageInfo_FindLocationQueryArgs.Size(m)
}
func (m *FindLocationQueryArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_FindLocationQueryArgs.DiscardUnknown(m)
}

var xxx_messageInfo_FindLocationQueryArgs proto.InternalMessageInfo

func (m *FindLocationQueryArgs) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *FindLocationQueryArgs) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *FindLocationQueryArgs) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) Reset()         { *m = FindOrGetLocationQueryArgs{} }
func (m *FindOrGetLocationQueryArgs) String() string { return proto.CompactTextString(m) }
func (*FindOrGetLocationQueryArgs) ProtoMessage()    {}
func (*FindOrGetLocationQueryArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba177a3af9a276cd, []int{9}
}
func (m *FindOrGetLocationQueryArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindOrGetLocationQueryArgs.Unmarshal(m, b)
}
func (m *FindOrGetLocationQueryArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindOrGetLocationQueryArgs.Marshal(b, m, deterministic)
}
func (m *FindOrGetLocationQueryArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindOrGetLocationQueryArgs.Merge(m, src)
}
func (m *FindOrGetLocationQueryArgs) XXX_Size() int {
	return xxx_messageInfo_FindOrGetLocationQueryArgs.Size(m)
}
func (m *FindOrGetLocationQueryArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_FindOrGetLocationQueryArgs.DiscardUnknown(m)
}

var xxx_messageInfo_FindOrGetLocationQueryArgs proto.InternalMessageInfo

func (m *FindOrGetLocationQueryArgs) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *FindOrGetLocationQueryArgs) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func init() {
	proto.RegisterType((*Province)(nil), "etop.vn.api.main.location.v1.Province")
	proto.RegisterType((*District)(nil), "etop.vn.api.main.location.v1.District")
	proto.RegisterType((*Ward)(nil), "etop.vn.api.main.location.v1.Ward")
	proto.RegisterType((*Extra)(nil), "etop.vn.api.main.location.v1.Extra")
	proto.RegisterType((*GetAllLocationsQueryArgs)(nil), "etop.vn.api.main.location.v1.GetAllLocationsQueryArgs")
	proto.RegisterType((*GetAllLocationsQueryResult)(nil), "etop.vn.api.main.location.v1.GetAllLocationsQueryResult")
	proto.RegisterType((*LocationQueryResult)(nil), "etop.vn.api.main.location.v1.LocationQueryResult")
	proto.RegisterType((*GetLocationQueryArgs)(nil), "etop.vn.api.main.location.v1.GetLocationQueryArgs")
	proto.RegisterType((*FindLocationQueryArgs)(nil), "etop.vn.api.main.location.v1.FindLocationQueryArgs")
	proto.RegisterType((*FindOrGetLocationQueryArgs)(nil), "etop.vn.api.main.location.v1.FindOrGetLocationQueryArgs")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/location/v1/api.proto", fileDescriptor_ba177a3af9a276cd)
}

var fileDescriptor_ba177a3af9a276cd = []byte{
	// 763 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xde, 0xcd, 0x1f, 0xc9, 0xb4, 0x14, 0xb1, 0x14, 0x64, 0x05, 0xb4, 0x4d, 0x83, 0x54, 0xa5,
	0x07, 0x12, 0x25, 0x48, 0x55, 0x0f, 0x5c, 0x9a, 0x16, 0xaa, 0x4a, 0x48, 0x14, 0xf3, 0x27, 0x71,
	0xa9, 0x4c, 0xbc, 0x4a, 0x2d, 0xa5, 0xb6, 0x65, 0x6f, 0x0c, 0xbd, 0x21, 0xb8, 0xf4, 0x82, 0xe0,
	0xc8, 0x23, 0x20, 0x9e, 0x80, 0x47, 0xe8, 0x31, 0x17, 0x24, 0x4e, 0x55, 0x93, 0xbc, 0x00, 0x12,
	0x2f, 0x80, 0x76, 0x63, 0xa7, 0x69, 0xed, 0xba, 0x26, 0xc0, 0x25, 0x8a, 0x67, 0xbe, 0xf9, 0x66,
	0xe6, 0xdb, 0x6f, 0x17, 0x96, 0x18, 0xb7, 0xec, 0xaa, 0x67, 0xd6, 0x34, 0xdb, 0xa8, 0xed, 0x69,
	0x86, 0x59, 0xeb, 0x58, 0x2d, 0x8d, 0x1b, 0x96, 0x59, 0xf3, 0xea, 0x22, 0x58, 0xb5, 0x1d, 0x8b,
	0x5b, 0xe4, 0x96, 0x8f, 0xab, 0x8a, 0x90, 0xc0, 0x55, 0x03, 0x5c, 0xd5, 0xab, 0x17, 0xe7, 0xdb,
	0x56, 0xdb, 0x92, 0xc0, 0x9a, 0xf8, 0x37, 0xaa, 0x29, 0x2e, 0xc7, 0x72, 0x3b, 0xac, 0x2d, 0xaa,
	0x47, 0xd0, 0x4a, 0x2c, 0xb4, 0xeb, 0xbc, 0xd2, 0x7c, 0x64, 0xf9, 0x3b, 0x86, 0xfc, 0xb6, 0x63,
	0x79, 0x86, 0xd9, 0x62, 0x44, 0x81, 0x8c, 0xa9, 0xed, 0x31, 0x05, 0x97, 0x70, 0xa5, 0xd0, 0xcc,
	0x1c, 0x1e, 0x2d, 0x20, 0x55, 0x46, 0x44, 0xa6, 0x65, 0xe9, 0x4c, 0x49, 0x4d, 0x66, 0x44, 0x84,
	0x6c, 0x43, 0x6e, 0xd4, 0x5a, 0x49, 0x97, 0x70, 0x65, 0xae, 0xd1, 0xa8, 0xc6, 0xad, 0x56, 0xf5,
	0xc7, 0x7c, 0x6e, 0x30, 0x6e, 0x6a, 0x7b, 0xaa, 0xfc, 0xf2, 0xf9, 0x7c, 0x1e, 0xb2, 0x0e, 0x59,
	0xf6, 0x86, 0x3b, 0x9a, 0x92, 0x29, 0xe1, 0xca, 0x4c, 0xe3, 0x76, 0x3c, 0xe1, 0x7d, 0x01, 0x6d,
	0xe6, 0x05, 0x43, 0xef, 0x68, 0x01, 0xab, 0xa3, 0xda, 0xf2, 0xbb, 0x14, 0xe4, 0x37, 0x0c, 0x97,
	0x3b, 0x46, 0x8b, 0x4f, 0xb5, 0xd7, 0x32, 0x5c, 0xb6, 0x7d, 0x5d, 0x76, 0x24, 0x24, 0x3d, 0x01,
	0x99, 0x0d, 0x52, 0xeb, 0x02, 0xaa, 0x02, 0x48, 0x49, 0x77, 0xf8, 0xbe, 0xcd, 0xe4, 0xd4, 0x73,
	0x8d, 0x3b, 0xf1, 0x53, 0x8f, 0x8e, 0xe0, 0x99, 0xf8, 0x7d, 0xba, 0x6f, 0x33, 0x9f, 0xb6, 0xd0,
	0x0d, 0x02, 0x27, 0x22, 0x64, 0xff, 0x42, 0x84, 0xaf, 0x18, 0x32, 0x2f, 0x34, 0x47, 0x9f, 0x56,
	0x00, 0xdd, 0x17, 0x30, 0x42, 0x80, 0x20, 0x25, 0x05, 0xf8, 0x27, 0x27, 0xd6, 0x86, 0xac, 0xcc,
	0x10, 0x0a, 0x97, 0x5c, 0x9b, 0xb5, 0x0c, 0xad, 0x23, 0xe7, 0xcd, 0xfb, 0x2d, 0x83, 0x20, 0xb9,
	0x09, 0xb9, 0xf6, 0xae, 0xb9, 0x63, 0xe8, 0x52, 0x9b, 0xac, 0x9f, 0xce, 0xb6, 0x77, 0xcd, 0x2d,
	0x9d, 0x2c, 0x42, 0xc1, 0xe3, 0xb6, 0xe5, 0x72, 0x91, 0xcf, 0x4d, 0xe4, 0xf3, 0xa3, 0xf0, 0x96,
	0x5e, 0x3e, 0xc0, 0xa0, 0x6c, 0x32, 0xbe, 0xd6, 0xe9, 0x3c, 0xf4, 0xa7, 0x72, 0x1f, 0x77, 0x99,
	0xb3, 0xbf, 0xe6, 0xb4, 0x5d, 0x72, 0x03, 0xd2, 0x5a, 0xe7, 0x74, 0x63, 0x11, 0x08, 0xdb, 0x21,
	0x75, 0xae, 0x1d, 0x92, 0x0b, 0x57, 0x1e, 0x60, 0x28, 0x46, 0x8d, 0xa2, 0x32, 0xb7, 0xdb, 0xe1,
	0x64, 0x03, 0x0a, 0x01, 0xb3, 0xab, 0xe0, 0x52, 0xba, 0x32, 0xd3, 0x58, 0x8a, 0xd7, 0x36, 0xb8,
	0xca, 0xea, 0x49, 0xa1, 0x60, 0x09, 0x9a, 0xba, 0x4a, 0x2a, 0x09, 0x4b, 0x70, 0x71, 0xd4, 0x93,
	0x42, 0xb2, 0x0a, 0xd9, 0xd7, 0x9a, 0xa3, 0xbb, 0x4a, 0x5a, 0x32, 0x94, 0xe3, 0x19, 0x84, 0xeb,
	0xd4, 0x51, 0x81, 0x78, 0x62, 0xae, 0x05, 0xeb, 0x4d, 0x6e, 0xd7, 0x84, 0x7c, 0x30, 0xa4, 0xd4,
	0x3b, 0xf9, 0x72, 0xe3, 0x3a, 0xc1, 0x11, 0x8c, 0x28, 0x4f, 0x24, 0xf9, 0x6a, 0xe3, 0x3a, 0xb2,
	0x02, 0x19, 0x31, 0xa8, 0x3c, 0xa6, 0x64, 0x8b, 0x49, 0x7c, 0xf9, 0x23, 0x86, 0xf9, 0x4d, 0xc6,
	0x4f, 0xad, 0x26, 0x3d, 0x14, 0xf2, 0x0a, 0x4e, 0xee, 0x95, 0xd4, 0xb9, 0x97, 0x6c, 0x11, 0x0a,
	0xa2, 0x6d, 0xd8, 0x52, 0x79, 0x11, 0x96, 0x76, 0xea, 0xc2, 0xf5, 0x07, 0x86, 0xa9, 0x87, 0x27,
	0x2a, 0x9d, 0x91, 0x7a, 0x5c, 0x3a, 0x16, 0xb2, 0x74, 0x46, 0xc8, 0x31, 0x42, 0x9f, 0x78, 0x44,
	0xc7, 0x32, 0x8d, 0x5f, 0x0a, 0x29, 0xc4, 0x2f, 0x0c, 0x45, 0xd1, 0xf7, 0x91, 0x13, 0x29, 0xc7,
	0x7f, 0x6d, 0x1e, 0x16, 0x3b, 0x93, 0x5c, 0xec, 0x6c, 0x32, 0xb1, 0x73, 0x51, 0x62, 0x37, 0x3e,
	0x64, 0x60, 0xfe, 0xd4, 0xb2, 0x4f, 0x98, 0xe3, 0x19, 0x2d, 0x46, 0xde, 0x63, 0xb8, 0x72, 0xe6,
	0x52, 0x93, 0x95, 0x78, 0x57, 0x9d, 0xf7, 0x1c, 0x15, 0x57, 0xff, 0xbc, 0xce, 0xbf, 0x5d, 0x1c,
	0x66, 0x26, 0x4e, 0x83, 0x34, 0x2e, 0x24, 0x0a, 0x1d, 0x5c, 0xb1, 0x1e, 0x5f, 0x13, 0x75, 0xa7,
	0x3d, 0x98, 0x9d, 0x74, 0x20, 0xb9, 0x1b, 0x4f, 0x11, 0xe9, 0xd6, 0x69, 0xfa, 0xbe, 0xc5, 0x70,
	0x35, 0x64, 0x41, 0xb2, 0x7a, 0x71, 0xf7, 0x68, 0xcf, 0x4e, 0x31, 0x42, 0xf3, 0xde, 0x61, 0x9f,
	0xe2, 0x5e, 0x9f, 0xa2, 0xe3, 0x3e, 0x45, 0x3f, 0xfb, 0x14, 0x1d, 0x0c, 0x28, 0xfa, 0x32, 0xa0,
	0xe8, 0xdb, 0x80, 0xa2, 0xc3, 0x01, 0x45, 0xbd, 0x01, 0x45, 0xc7, 0x03, 0x8a, 0x3e, 0x0d, 0x29,
	0xfa, 0x3c, 0xa4, 0xa8, 0x37, 0xa4, 0xe8, 0xc7, 0x90, 0xa2, 0x97, 0x29, 0xaf, 0xfe, 0x3b, 0x00,
	0x00, 0xff, 0xff, 0x8f, 0x7f, 0x89, 0x42, 0x39, 0x0a, 0x00, 0x00,
}