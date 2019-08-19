// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/affiliate/affiliate.proto

package affiliate

import (
	_ "etop.vn/backend/pb/common"
	status3 "etop.vn/backend/pb/etop/etc/status3"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
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

type Affiliate struct {
	Id                   int64          `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string         `protobuf:"bytes,2,opt,name=name" json:"name"`
	Status               status3.Status `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	IsTest               bool           `protobuf:"varint,4,opt,name=is_test,json=isTest" json:"is_test"`
	Phone                string         `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Email                string         `protobuf:"bytes,6,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Affiliate) Reset()         { *m = Affiliate{} }
func (m *Affiliate) String() string { return proto.CompactTextString(m) }
func (*Affiliate) ProtoMessage()    {}
func (*Affiliate) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{0}
}
func (m *Affiliate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Affiliate.Unmarshal(m, b)
}
func (m *Affiliate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Affiliate.Marshal(b, m, deterministic)
}
func (m *Affiliate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Affiliate.Merge(m, src)
}
func (m *Affiliate) XXX_Size() int {
	return xxx_messageInfo_Affiliate.Size(m)
}
func (m *Affiliate) XXX_DiscardUnknown() {
	xxx_messageInfo_Affiliate.DiscardUnknown(m)
}

var xxx_messageInfo_Affiliate proto.InternalMessageInfo

func (m *Affiliate) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Affiliate) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Affiliate) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Affiliate) GetIsTest() bool {
	if m != nil {
		return m.IsTest
	}
	return false
}

func (m *Affiliate) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Affiliate) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type RegisterAffiliateRequest struct {
	// @required
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,3,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterAffiliateRequest) Reset()         { *m = RegisterAffiliateRequest{} }
func (m *RegisterAffiliateRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterAffiliateRequest) ProtoMessage()    {}
func (*RegisterAffiliateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{1}
}
func (m *RegisterAffiliateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterAffiliateRequest.Unmarshal(m, b)
}
func (m *RegisterAffiliateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterAffiliateRequest.Marshal(b, m, deterministic)
}
func (m *RegisterAffiliateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterAffiliateRequest.Merge(m, src)
}
func (m *RegisterAffiliateRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterAffiliateRequest.Size(m)
}
func (m *RegisterAffiliateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterAffiliateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterAffiliateRequest proto.InternalMessageInfo

func (m *RegisterAffiliateRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RegisterAffiliateRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterAffiliateRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UpdateAffiliateRequest struct {
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateAffiliateRequest) Reset()         { *m = UpdateAffiliateRequest{} }
func (m *UpdateAffiliateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateAffiliateRequest) ProtoMessage()    {}
func (*UpdateAffiliateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{2}
}
func (m *UpdateAffiliateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateAffiliateRequest.Unmarshal(m, b)
}
func (m *UpdateAffiliateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateAffiliateRequest.Marshal(b, m, deterministic)
}
func (m *UpdateAffiliateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateAffiliateRequest.Merge(m, src)
}
func (m *UpdateAffiliateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateAffiliateRequest.Size(m)
}
func (m *UpdateAffiliateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateAffiliateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateAffiliateRequest proto.InternalMessageInfo

func (m *UpdateAffiliateRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateAffiliateRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *UpdateAffiliateRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func init() {
	proto.RegisterType((*Affiliate)(nil), "affiliate.Affiliate")
	proto.RegisterType((*RegisterAffiliateRequest)(nil), "affiliate.RegisterAffiliateRequest")
	proto.RegisterType((*UpdateAffiliateRequest)(nil), "affiliate.UpdateAffiliateRequest")
}

func init() { proto.RegisterFile("etop/affiliate/affiliate.proto", fileDescriptor_7d5a21d966c22074) }

var fileDescriptor_7d5a21d966c22074 = []byte{
	// 508 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xdd, 0x75, 0xd2, 0x42, 0xb6, 0xa2, 0x11, 0x4b, 0x54, 0xac, 0x48, 0x5d, 0xdc, 0x72, 0x89,
	0x04, 0xb1, 0xd5, 0x70, 0xe3, 0x44, 0xaa, 0x72, 0x08, 0x12, 0x08, 0x25, 0x85, 0x03, 0x17, 0xb4,
	0x75, 0x26, 0x66, 0xc1, 0xde, 0x35, 0xd9, 0x4d, 0x10, 0x7f, 0xc0, 0x11, 0x71, 0xe2, 0x13, 0xf8,
	0x04, 0xae, 0xdc, 0x72, 0xec, 0x95, 0x0b, 0x6a, 0xdc, 0x1f, 0xe0, 0x13, 0x90, 0x1d, 0x3b, 0x31,
	0x84, 0xf4, 0x34, 0xa3, 0xf7, 0x66, 0xde, 0x9b, 0x19, 0xed, 0x12, 0x06, 0x46, 0xc5, 0x1e, 0x1f,
	0x8d, 0x44, 0x28, 0xb8, 0x81, 0x55, 0xe6, 0xc6, 0x63, 0x65, 0x14, 0xad, 0x2d, 0x81, 0xe6, 0xfd,
	0x0c, 0xf1, 0xdb, 0x01, 0xc8, 0xb6, 0xfe, 0xc0, 0x83, 0x00, 0xc6, 0x9e, 0x8a, 0x8d, 0x50, 0x52,
	0x7b, 0x5c, 0x4a, 0x65, 0x78, 0x96, 0x2f, 0x1a, 0x9b, 0x8d, 0x40, 0x05, 0x2a, 0x4b, 0xbd, 0x34,
	0xcb, 0xd1, 0x5b, 0xbe, 0x8a, 0x22, 0x25, 0xbd, 0x45, 0xc8, 0xc1, 0xfd, 0x6c, 0x06, 0x30, 0xbe,
	0xa7, 0x0d, 0x37, 0x13, 0xfd, 0x20, 0x8f, 0x0b, 0xfa, 0xf0, 0x07, 0x26, 0xb5, 0x6e, 0x31, 0x05,
	0x6d, 0x10, 0x4b, 0x0c, 0x6d, 0xec, 0xe0, 0x56, 0xe5, 0xb8, 0x3a, 0xfb, 0x75, 0x07, 0xf5, 0x2d,
	0x31, 0xa4, 0x36, 0xa9, 0x4a, 0x1e, 0x81, 0x6d, 0x39, 0xb8, 0x55, 0xcb, 0xf1, 0x0c, 0xa1, 0x6d,
	0xb2, 0xbd, 0x50, 0xb3, 0x2b, 0x0e, 0x6e, 0xed, 0x76, 0xea, 0x6e, 0x6e, 0xe2, 0x0e, 0xb2, 0x98,
	0x17, 0xe7, 0x45, 0x74, 0x9f, 0x5c, 0x13, 0xfa, 0xb5, 0x01, 0x6d, 0xec, 0xaa, 0x83, 0x5b, 0xd7,
	0x0b, 0x5a, 0xe8, 0x53, 0xd0, 0x86, 0x36, 0xc9, 0x56, 0xfc, 0x46, 0x49, 0xb0, 0xb7, 0x4a, 0x46,
	0x0b, 0x28, 0xe5, 0x20, 0xe2, 0x22, 0xb4, 0xb7, 0xcb, 0x5c, 0x06, 0x1d, 0x86, 0xc4, 0xee, 0x43,
	0x20, 0xb4, 0x81, 0xf1, 0x72, 0x95, 0x3e, 0xbc, 0x9f, 0xa4, 0x9a, 0xc5, 0xec, 0x78, 0x6d, 0xf6,
	0xa5, 0x9b, 0x75, 0x85, 0x5b, 0x65, 0xdd, 0xed, 0x2d, 0xd9, 0x7b, 0x11, 0x0f, 0xb9, 0x81, 0x8d,
	0x5e, 0xd6, 0x66, 0xaf, 0xca, 0x15, 0x5e, 0xd5, 0x35, 0xaf, 0xce, 0x23, 0xb2, 0xf3, 0x54, 0x68,
	0x7f, 0x00, 0xe3, 0xa9, 0xf0, 0x81, 0x1e, 0x91, 0x9d, 0x97, 0x30, 0xd6, 0x42, 0xc9, 0x9e, 0x1c,
	0x29, 0x5a, 0x73, 0xfd, 0xc8, 0x7d, 0x1c, 0xc5, 0xe6, 0x63, 0xf3, 0x76, 0x9a, 0x96, 0xb8, 0x3e,
	0xe8, 0x58, 0x49, 0x0d, 0x9d, 0x9f, 0x98, 0xec, 0x76, 0x7d, 0x5f, 0x4d, 0xa4, 0x29, 0x54, 0x9e,
	0x91, 0x9b, 0x6b, 0xe7, 0xa2, 0x77, 0xdd, 0xd5, 0xe3, 0xdc, 0x74, 0xcc, 0x66, 0xa3, 0x54, 0xb4,
	0x6a, 0x7d, 0x42, 0xea, 0xff, 0x1c, 0x84, 0x1e, 0x94, 0x0a, 0xff, 0x7f, 0xac, 0x0d, 0x5a, 0xf7,
	0x48, 0xfd, 0x04, 0x42, 0x28, 0x6b, 0xdd, 0x48, 0x57, 0xeb, 0x9d, 0x14, 0x7d, 0xab, 0xa5, 0x8f,
	0xc3, 0x2f, 0xdd, 0x3d, 0xda, 0x20, 0xbb, 0x70, 0xaa, 0x62, 0xa7, 0xfb, 0xbc, 0xf7, 0xd0, 0x19,
	0xf0, 0x10, 0x3a, 0xd6, 0xf4, 0x68, 0x36, 0x67, 0xf8, 0x7c, 0xce, 0xf0, 0xc5, 0x9c, 0xa1, 0xdf,
	0x73, 0x86, 0x3e, 0x25, 0x0c, 0x7d, 0x4b, 0x18, 0xfa, 0x9e, 0x30, 0x34, 0x4b, 0x18, 0x3a, 0x4f,
	0x18, 0xba, 0x48, 0x18, 0xfa, 0x7c, 0xc9, 0xd0, 0xd7, 0x4b, 0x86, 0x5e, 0x1d, 0xa4, 0xbf, 0xc4,
	0x9d, 0x4a, 0xef, 0x8c, 0xfb, 0xef, 0x40, 0x0e, 0xbd, 0xf8, 0xcc, 0xfb, 0xfb, 0xf3, 0xfe, 0x09,
	0x00, 0x00, 0xff, 0xff, 0xf4, 0xe0, 0x69, 0xf9, 0xcd, 0x03, 0x00, 0x00,
}
