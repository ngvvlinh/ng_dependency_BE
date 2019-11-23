// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: external/partner/partner.proto

package partner

import (
	_ "etop.vn/api/pb/common"
	_ "etop.vn/api/pb/etop"
	_ "etop.vn/api/pb/external"
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

type AuthorizeShopRequest struct {
	ShopId               int64    `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	ExternalShopId       string   `protobuf:"bytes,3,opt,name=external_shop_id,json=externalShopId" json:"external_shop_id"`
	Name                 string   `protobuf:"bytes,4,opt,name=name" json:"name"`
	Phone                string   `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,6,opt,name=email" json:"email"`
	RedirectUrl          string   `protobuf:"bytes,8,opt,name=redirect_url,json=redirectUrl" json:"redirect_url"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizeShopRequest) Reset()         { *m = AuthorizeShopRequest{} }
func (m *AuthorizeShopRequest) String() string { return proto.CompactTextString(m) }
func (*AuthorizeShopRequest) ProtoMessage()    {}
func (*AuthorizeShopRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_19642bc1754e1e01, []int{0}
}
func (m *AuthorizeShopRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthorizeShopRequest.Unmarshal(m, b)
}
func (m *AuthorizeShopRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthorizeShopRequest.Marshal(b, m, deterministic)
}
func (m *AuthorizeShopRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizeShopRequest.Merge(m, src)
}
func (m *AuthorizeShopRequest) XXX_Size() int {
	return xxx_messageInfo_AuthorizeShopRequest.Size(m)
}
func (m *AuthorizeShopRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizeShopRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizeShopRequest proto.InternalMessageInfo

func (m *AuthorizeShopRequest) GetShopId() int64 {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *AuthorizeShopRequest) GetExternalShopId() string {
	if m != nil {
		return m.ExternalShopId
	}
	return ""
}

func (m *AuthorizeShopRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AuthorizeShopRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *AuthorizeShopRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *AuthorizeShopRequest) GetRedirectUrl() string {
	if m != nil {
		return m.RedirectUrl
	}
	return ""
}

type AuthorizeShopResponse struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Type                 string   `protobuf:"bytes,3,opt,name=type" json:"type"`
	AuthToken            string   `protobuf:"bytes,4,opt,name=auth_token,json=authToken" json:"auth_token"`
	ExpiresIn            int32    `protobuf:"varint,5,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	AuthUrl              string   `protobuf:"bytes,6,opt,name=auth_url,json=authUrl" json:"auth_url"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizeShopResponse) Reset()         { *m = AuthorizeShopResponse{} }
func (m *AuthorizeShopResponse) String() string { return proto.CompactTextString(m) }
func (*AuthorizeShopResponse) ProtoMessage()    {}
func (*AuthorizeShopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_19642bc1754e1e01, []int{1}
}
func (m *AuthorizeShopResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthorizeShopResponse.Unmarshal(m, b)
}
func (m *AuthorizeShopResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthorizeShopResponse.Marshal(b, m, deterministic)
}
func (m *AuthorizeShopResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizeShopResponse.Merge(m, src)
}
func (m *AuthorizeShopResponse) XXX_Size() int {
	return xxx_messageInfo_AuthorizeShopResponse.Size(m)
}
func (m *AuthorizeShopResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizeShopResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizeShopResponse proto.InternalMessageInfo

func (m *AuthorizeShopResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *AuthorizeShopResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *AuthorizeShopResponse) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *AuthorizeShopResponse) GetAuthToken() string {
	if m != nil {
		return m.AuthToken
	}
	return ""
}

func (m *AuthorizeShopResponse) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func (m *AuthorizeShopResponse) GetAuthUrl() string {
	if m != nil {
		return m.AuthUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthorizeShopRequest)(nil), "partner.AuthorizeShopRequest")
	proto.RegisterType((*AuthorizeShopResponse)(nil), "partner.AuthorizeShopResponse")
}

func init() { proto.RegisterFile("external/partner/partner.proto", fileDescriptor_19642bc1754e1e01) }

var fileDescriptor_19642bc1754e1e01 = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x4d, 0x8e, 0xd3, 0x30,
	0x18, 0x86, 0x1d, 0xfa, 0x6f, 0x10, 0x20, 0x53, 0x20, 0xaa, 0x84, 0x5b, 0x95, 0x05, 0x5d, 0x35,
	0x67, 0xa0, 0xbb, 0x2e, 0x69, 0xe9, 0x86, 0x4d, 0x14, 0x9a, 0x4f, 0x4d, 0x44, 0x62, 0x1b, 0xc7,
	0x41, 0x85, 0x13, 0xb0, 0x64, 0xc9, 0x11, 0xe6, 0x08, 0x73, 0x84, 0xae, 0x46, 0x3d, 0xc1, 0xa8,
	0x49, 0x2f, 0x30, 0x47, 0x18, 0xd9, 0x89, 0x47, 0x69, 0x37, 0xb1, 0xf3, 0x3c, 0x6f, 0x1c, 0xbf,
	0xd2, 0x87, 0x29, 0xec, 0x15, 0x48, 0x16, 0x24, 0x9e, 0x08, 0xa4, 0x62, 0x20, 0xed, 0x3a, 0x17,
	0x92, 0x2b, 0x4e, 0x7a, 0xf5, 0xeb, 0x68, 0xb8, 0xe3, 0x3b, 0x6e, 0x98, 0xa7, 0x77, 0x95, 0x1e,
	0xbd, 0xd9, 0xf2, 0x34, 0xe5, 0xcc, 0xab, 0x96, 0x1a, 0xbe, 0x02, 0xc5, 0x85, 0xa7, 0x1f, 0x35,
	0x78, 0xff, 0xf4, 0x13, 0xbb, 0xa9, 0xc4, 0xf4, 0xe4, 0xe0, 0xe1, 0xe7, 0x5c, 0x45, 0x5c, 0xc6,
	0x7f, 0x60, 0x1d, 0x71, 0xb1, 0x82, 0x9f, 0x39, 0x64, 0x8a, 0x7c, 0xc0, 0xbd, 0x2c, 0xe2, 0xc2,
	0x8f, 0x43, 0xf7, 0xd9, 0xc4, 0x99, 0xb5, 0x16, 0xed, 0xc3, 0xfd, 0x18, 0xad, 0xba, 0x1a, 0x2e,
	0x43, 0x32, 0xc7, 0xaf, 0xed, 0x49, 0xbe, 0xcd, 0xb5, 0x26, 0xce, 0x6c, 0x50, 0xe7, 0x5e, 0x5a,
	0xbb, 0xae, 0xf2, 0x2e, 0x6e, 0xb3, 0x20, 0x05, 0xb7, 0xdd, 0xc8, 0x18, 0x42, 0x46, 0xb8, 0x23,
	0x22, 0xce, 0xc0, 0xed, 0x34, 0x54, 0x85, 0xb4, 0x83, 0x34, 0x88, 0x13, 0xb7, 0xdb, 0x74, 0x06,
	0x91, 0x4f, 0xf8, 0x85, 0x84, 0x30, 0x96, 0xb0, 0x55, 0x7e, 0x2e, 0x13, 0xb7, 0xdf, 0x88, 0x3c,
	0xb7, 0x66, 0x23, 0x93, 0xe9, 0x9d, 0x83, 0xdf, 0x5e, 0x55, 0xcc, 0x04, 0x67, 0x19, 0xe8, 0x4b,
	0x6d, 0x79, 0x08, 0xae, 0xd3, 0xbc, 0x94, 0x26, 0xe4, 0x1d, 0x6e, 0xa5, 0xd9, 0xce, 0x34, 0xb7,
	0x42, 0x03, 0xfd, 0x85, 0xfa, 0x2d, 0xe0, 0xa2, 0xaa, 0x21, 0xe4, 0x23, 0xc6, 0x41, 0xae, 0x22,
	0x5f, 0xf1, 0x1f, 0xc0, 0x2e, 0x6a, 0x0e, 0x34, 0xff, 0xaa, 0xb1, 0x0e, 0xc1, 0x5e, 0xc4, 0x12,
	0x32, 0x3f, 0x66, 0xa6, 0x70, 0xc7, 0x86, 0x6a, 0xbe, 0x64, 0x64, 0x8c, 0xfb, 0xe6, 0x24, 0x5d,
	0xaa, 0xd9, 0xbb, 0xa7, 0xe9, 0x46, 0x26, 0x8b, 0x2f, 0x87, 0x82, 0x3a, 0xc7, 0x82, 0x3a, 0xa7,
	0x82, 0xa2, 0x87, 0x82, 0xa2, 0xbf, 0x25, 0x45, 0x37, 0x25, 0x45, 0xb7, 0x25, 0x45, 0x87, 0x92,
	0xa2, 0x63, 0x49, 0xd1, 0xa9, 0xa4, 0xe8, 0xdf, 0x99, 0xa2, 0xff, 0x67, 0x8a, 0xbe, 0x8d, 0xcd,
	0x28, 0xfc, 0x62, 0x5e, 0x20, 0x62, 0x4f, 0x7c, 0xf7, 0xae, 0x47, 0xee, 0x31, 0x00, 0x00, 0xff,
	0xff, 0x90, 0xfe, 0xc9, 0xaa, 0x85, 0x02, 0x00, 0x00,
}