// source: etop/affiliate/affiliate.proto

package affiliate

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	_ "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	_ "etop.vn/api/pb/etop/etc/status3"
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
	BankAccount          *etop.BankAccount `protobuf:"bytes,4,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RegisterAffiliateRequest) Reset()         { *m = RegisterAffiliateRequest{} }
func (m *RegisterAffiliateRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterAffiliateRequest) ProtoMessage()    {}
func (*RegisterAffiliateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{0}
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
	BankAccount          *etop.BankAccount `protobuf:"bytes,1,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateAffiliateBankAccountRequest) Reset()         { *m = UpdateAffiliateBankAccountRequest{} }
func (m *UpdateAffiliateBankAccountRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateAffiliateBankAccountRequest) ProtoMessage()    {}
func (*UpdateAffiliateBankAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d5a21d966c22074, []int{2}
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
	// 322 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x4e, 0x02, 0x31,
	0x14, 0x86, 0x5b, 0x18, 0x17, 0x14, 0x13, 0xe3, 0x68, 0x4c, 0x43, 0xe2, 0x13, 0x59, 0xb1, 0xa2,
	0x89, 0x7a, 0x01, 0xb8, 0x80, 0xc9, 0x24, 0x2e, 0x74, 0x63, 0xca, 0x58, 0x70, 0x84, 0x69, 0x2b,
	0x53, 0x3c, 0x83, 0x4b, 0x97, 0xae, 0x5d, 0x79, 0x04, 0x8f, 0x30, 0x4b, 0x4e, 0x60, 0x98, 0x72,
	0x01, 0x8f, 0x60, 0xe8, 0x8c, 0x32, 0x62, 0x30, 0xae, 0xde, 0x9b, 0xef, 0x9f, 0xf4, 0x7d, 0x2f,
	0x2d, 0x01, 0x61, 0x94, 0x66, 0x7c, 0x30, 0x88, 0xc6, 0x11, 0x37, 0x62, 0xd5, 0x75, 0xf4, 0x44,
	0x19, 0xe5, 0xd7, 0xbe, 0x41, 0x63, 0x7f, 0xa8, 0x86, 0xca, 0x51, 0xb6, 0xec, 0xf2, 0x1f, 0x1a,
	0x7b, 0xa1, 0x8a, 0x63, 0x25, 0x59, 0x5e, 0x0a, 0x78, 0xe8, 0x4e, 0x15, 0x26, 0x64, 0x89, 0xe1,
	0x66, 0x9a, 0x9c, 0x16, 0xb5, 0x88, 0x77, 0x8a, 0x58, 0xe9, 0x1c, 0xb4, 0x5e, 0x30, 0xa1, 0x81,
	0x18, 0x46, 0x89, 0x11, 0x93, 0xee, 0xd7, 0xc0, 0x40, 0xdc, 0x4f, 0x45, 0x62, 0x7c, 0x4a, 0x3c,
	0xc9, 0x63, 0x41, 0x71, 0x13, 0xb7, 0x6b, 0x3d, 0x2f, 0x7d, 0x3f, 0x42, 0x81, 0x23, 0x7e, 0x83,
	0x6c, 0xe9, 0x5b, 0x25, 0x05, 0xad, 0x94, 0xa2, 0x1c, 0x2d, 0x33, 0x11, 0xf3, 0x68, 0x4c, 0xab,
	0xe5, 0xcc, 0x21, 0xff, 0x8c, 0x6c, 0xf7, 0xb9, 0x1c, 0x5d, 0xf3, 0x30, 0x54, 0x53, 0x69, 0xa8,
	0xd7, 0xc4, 0xed, 0xfa, 0xc9, 0x6e, 0xc7, 0x19, 0xf5, 0xb8, 0x1c, 0x75, 0xf3, 0x20, 0xa8, 0xf7,
	0x57, 0x1f, 0xad, 0x3b, 0x72, 0x70, 0xa1, 0x6f, 0xb8, 0x11, 0x1b, 0x0d, 0x2b, 0x9b, 0x0d, 0xab,
	0x7f, 0x18, 0x7a, 0xbf, 0x0c, 0x5b, 0x97, 0xe4, 0x78, 0x6d, 0x56, 0x59, 0xab, 0x18, 0xbb, 0xbe,
	0x06, 0xfe, 0xcf, 0x1a, 0xbd, 0xf3, 0x34, 0x03, 0x3c, 0xcb, 0x00, 0xcf, 0x33, 0x40, 0x1f, 0x19,
	0xa0, 0x47, 0x0b, 0xe8, 0xd5, 0x02, 0x7a, 0xb3, 0x80, 0x52, 0x0b, 0x68, 0x66, 0x01, 0xcd, 0x2d,
	0xa0, 0xa7, 0x05, 0xa0, 0xe7, 0x05, 0xa0, 0x2b, 0x77, 0x9b, 0x9d, 0x07, 0xc9, 0xb8, 0x8e, 0x98,
	0xee, 0xb3, 0x9f, 0x4f, 0xe6, 0x33, 0x00, 0x00, 0xff, 0xff, 0xf9, 0x2c, 0xdc, 0x8f, 0x43, 0x02,
	0x00, 0x00,
}
