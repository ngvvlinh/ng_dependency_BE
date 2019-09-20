// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop/affiliate/affiliate.proto

package affiliate

import (
	_ "etop.vn/backend/pb/common"
	etop "etop.vn/backend/pb/etop"
	_ "etop.vn/backend/pb/etop/etc/status3"
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

type RegisterAffiliateRequest struct {
	// @required
	Name                 string            `protobuf:"bytes,1,opt,name=name" json:"name"`
	Phone                string            `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Email                string            `protobuf:"bytes,3,opt,name=email" json:"email"`
	BankAccount          *etop.BankAccount `protobuf:"bytes,4,opt,name=bank_account,json=bankAccount" json:"bank_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RegisterAffiliateRequest) Reset()         { *m = RegisterAffiliateRequest{} }
func (m *RegisterAffiliateRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterAffiliateRequest) ProtoMessage()    {}
func (*RegisterAffiliateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{0}
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

func (m *RegisterAffiliateRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
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
	return fileDescriptor_7d5a21d966c22074, []int{1}
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

type UpdateAffiliateBankAccountRequest struct {
	BankAccount          *etop.BankAccount `protobuf:"bytes,1,opt,name=bank_account,json=bankAccount" json:"bank_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateAffiliateBankAccountRequest) Reset()         { *m = UpdateAffiliateBankAccountRequest{} }
func (m *UpdateAffiliateBankAccountRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateAffiliateBankAccountRequest) ProtoMessage()    {}
func (*UpdateAffiliateBankAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{2}
}
func (m *UpdateAffiliateBankAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateAffiliateBankAccountRequest.Unmarshal(m, b)
}
func (m *UpdateAffiliateBankAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateAffiliateBankAccountRequest.Marshal(b, m, deterministic)
}
func (m *UpdateAffiliateBankAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateAffiliateBankAccountRequest.Merge(m, src)
}
func (m *UpdateAffiliateBankAccountRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateAffiliateBankAccountRequest.Size(m)
}
func (m *UpdateAffiliateBankAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateAffiliateBankAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateAffiliateBankAccountRequest proto.InternalMessageInfo

func (m *UpdateAffiliateBankAccountRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func init() {
	proto.RegisterType((*RegisterAffiliateRequest)(nil), "affiliate.RegisterAffiliateRequest")
	proto.RegisterType((*UpdateAffiliateRequest)(nil), "affiliate.UpdateAffiliateRequest")
	proto.RegisterType((*UpdateAffiliateBankAccountRequest)(nil), "affiliate.UpdateAffiliateBankAccountRequest")
}

func init() { proto.RegisterFile("etop/affiliate/affiliate.proto", fileDescriptor_7d5a21d966c22074) }

var fileDescriptor_7d5a21d966c22074 = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x4f, 0x6e, 0xd3, 0x40,
	0x14, 0xc6, 0xc7, 0x69, 0x58, 0x64, 0x02, 0x44, 0x9d, 0x22, 0x30, 0x46, 0x0c, 0x49, 0xd8, 0x54,
	0xa2, 0xb5, 0xd5, 0xc0, 0x8a, 0x15, 0x89, 0xda, 0x45, 0x16, 0x48, 0xc8, 0xfc, 0x91, 0x80, 0x05,
	0x9a, 0xb8, 0x2f, 0xc6, 0x24, 0x9e, 0x99, 0x7a, 0x26, 0x41, 0xdc, 0x80, 0x25, 0x62, 0xc5, 0x9a,
	0x15, 0x0b, 0x0e, 0xc0, 0x11, 0xb2, 0xec, 0x09, 0x50, 0xe3, 0x5e, 0x80, 0x23, 0xa0, 0xd8, 0x4e,
	0x6d, 0x62, 0x8c, 0x58, 0xcd, 0xcb, 0xf7, 0xe9, 0x7d, 0xf3, 0x9b, 0xf7, 0x62, 0x4c, 0x41, 0x0b,
	0xe9, 0xb0, 0xf1, 0x38, 0x98, 0x06, 0x4c, 0x43, 0x5e, 0xd9, 0x32, 0x12, 0x5a, 0x90, 0xc6, 0x85,
	0x60, 0xed, 0x25, 0x8a, 0xb7, 0xef, 0x03, 0xdf, 0x57, 0xef, 0x99, 0xef, 0x43, 0xe4, 0x08, 0xa9,
	0x03, 0xc1, 0x95, 0xc3, 0x38, 0x17, 0x9a, 0x25, 0x75, 0xda, 0x68, 0x5d, 0xf3, 0x85, 0x2f, 0x92,
	0xd2, 0x59, 0x55, 0x99, 0xba, 0xe3, 0x89, 0x30, 0x14, 0xdc, 0x49, 0x8f, 0x4c, 0xbc, 0x9d, 0x30,
	0x80, 0xf6, 0x1c, 0xa5, 0x99, 0x9e, 0xa9, 0xfb, 0xd9, 0x99, 0xd9, 0xad, 0xcc, 0x16, 0x32, 0x15,
	0xba, 0x5f, 0x0d, 0x6c, 0xba, 0xe0, 0x07, 0x4a, 0x43, 0xd4, 0x5f, 0xe3, 0xb9, 0x70, 0x32, 0x03,
	0xa5, 0x89, 0x89, 0xeb, 0x9c, 0x85, 0x60, 0x1a, 0x6d, 0x63, 0xb7, 0x31, 0xa8, 0x2f, 0x7e, 0xde,
	0x41, 0x6e, 0xa2, 0x10, 0x0b, 0x5f, 0x92, 0x6f, 0x05, 0x07, 0xb3, 0x56, 0xb0, 0x52, 0x69, 0xe5,
	0x41, 0xc8, 0x82, 0xa9, 0xb9, 0x55, 0xf4, 0x12, 0x89, 0x3c, 0xc0, 0x97, 0x47, 0x8c, 0x4f, 0xde,
	0x30, 0xcf, 0x13, 0x33, 0xae, 0xcd, 0x7a, 0xdb, 0xd8, 0x6d, 0xf6, 0xb6, 0xed, 0x84, 0x68, 0xc0,
	0xf8, 0xa4, 0x9f, 0x1a, 0x6e, 0x73, 0x94, 0xff, 0xe8, 0xbe, 0xc3, 0xd7, 0x9f, 0xcb, 0x63, 0xa6,
	0xa1, 0x92, 0xb0, 0x56, 0x4d, 0xb8, 0xf5, 0x0f, 0xc2, 0x7a, 0x89, 0xb0, 0xfb, 0x12, 0x77, 0x36,
	0xee, 0x2a, 0x62, 0x65, 0xd7, 0x6e, 0x3e, 0xc3, 0xf8, 0x9f, 0x67, 0xf4, 0x1e, 0xe1, 0xe6, 0xe3,
	0x40, 0x79, 0x4f, 0x21, 0x9a, 0x07, 0x1e, 0x90, 0x03, 0xdc, 0x7c, 0x01, 0x91, 0x0a, 0x04, 0x1f,
	0xf2, 0xb1, 0x20, 0x0d, 0xdb, 0x0b, 0xed, 0xa3, 0x50, 0xea, 0x0f, 0xd6, 0x8d, 0x55, 0x59, 0xf0,
	0x5c, 0x50, 0x52, 0x70, 0x05, 0xbd, 0xef, 0x35, 0x7c, 0x35, 0x4b, 0x5b, 0xa7, 0x0c, 0xf1, 0x76,
	0x69, 0x7f, 0xe4, 0xae, 0x9d, 0xff, 0xf7, 0xaa, 0xb6, 0x6b, 0xb5, 0x52, 0xdc, 0xbc, 0xeb, 0x08,
	0xb7, 0x36, 0x9e, 0x4e, 0x3a, 0x85, 0xa0, 0xbf, 0xaf, 0xa0, 0x1c, 0xf3, 0x1a, 0x5b, 0xd5, 0x13,
	0x24, 0x7b, 0xd5, 0x89, 0xe5, 0x41, 0x97, 0xc3, 0xef, 0xe1, 0xd6, 0x21, 0x4c, 0xa1, 0xc8, 0x78,
	0x65, 0x35, 0xad, 0xe1, 0xe1, 0xba, 0x25, 0x9f, 0xe3, 0xe0, 0xe4, 0x73, 0xff, 0x16, 0xb9, 0x89,
	0x77, 0xe0, 0x99, 0x90, 0xed, 0xfe, 0x93, 0xe1, 0xc3, 0xf6, 0x45, 0x57, 0xaf, 0x36, 0x3f, 0x58,
	0x2c, 0xa9, 0x71, 0xba, 0xa4, 0xc6, 0xd9, 0x92, 0xa2, 0x5f, 0x4b, 0x8a, 0x3e, 0xc6, 0x14, 0x7d,
	0x8b, 0x29, 0xfa, 0x11, 0x53, 0xb4, 0x88, 0x29, 0x3a, 0x8d, 0x29, 0x3a, 0x8b, 0x29, 0xfa, 0x74,
	0x4e, 0xd1, 0x97, 0x73, 0x8a, 0x5e, 0x75, 0x12, 0x98, 0x39, 0x77, 0x46, 0xcc, 0x9b, 0x00, 0x3f,
	0x76, 0xe4, 0xc8, 0xf9, 0xf3, 0x9b, 0xff, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x24, 0x4e, 0x35, 0x24,
	0x04, 0x04, 0x00, 0x00,
}
