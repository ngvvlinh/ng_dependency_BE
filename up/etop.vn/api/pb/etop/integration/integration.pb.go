// source: etop/integration/integration.proto

package integration

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	_ "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	status3 "etop.vn/api/pb/etop/etc/status3"
	"etop.vn/capi/dot"
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

type InitRequest struct {
	AuthToken            string   `protobuf:"bytes,1,opt,name=auth_token,json=authToken" json:"auth_token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitRequest) Reset()         { *m = InitRequest{} }
func (m *InitRequest) String() string { return proto.CompactTextString(m) }
func (*InitRequest) ProtoMessage()    {}
func (*InitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{0}
}

var xxx_messageInfo_InitRequest proto.InternalMessageInfo

func (m *InitRequest) GetAuthToken() string {
	if m != nil {
		return m.AuthToken
	}
	return ""
}

type LoginResponse struct {
	AccessToken          string                     `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token"`
	ExpiresIn            int32                      `protobuf:"varint,2,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	User                 *PartnerUserLogin          `protobuf:"bytes,3,opt,name=user" json:"user"`
	Account              *PartnerShopLoginAccount   `protobuf:"bytes,4,opt,name=account" json:"account"`
	Shop                 *PartnerShopInfo           `protobuf:"bytes,5,opt,name=shop" json:"shop"`
	AvailableAccounts    []*PartnerShopLoginAccount `protobuf:"bytes,7,rep,name=available_accounts,json=availableAccounts" json:"available_accounts"`
	AuthPartner          *etop.PublicAccountInfo    `protobuf:"bytes,11,opt,name=auth_partner,json=authPartner" json:"auth_partner"`
	Actions              []*Action                  `protobuf:"bytes,12,rep,name=actions" json:"actions"`
	RedirectUrl          string                     `protobuf:"bytes,13,opt,name=redirect_url,json=redirectUrl" json:"redirect_url"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{1}
}

var xxx_messageInfo_LoginResponse proto.InternalMessageInfo

func (m *LoginResponse) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *LoginResponse) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func (m *LoginResponse) GetUser() *PartnerUserLogin {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *LoginResponse) GetAccount() *PartnerShopLoginAccount {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *LoginResponse) GetShop() *PartnerShopInfo {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *LoginResponse) GetAvailableAccounts() []*PartnerShopLoginAccount {
	if m != nil {
		return m.AvailableAccounts
	}
	return nil
}

func (m *LoginResponse) GetAuthPartner() *etop.PublicAccountInfo {
	if m != nil {
		return m.AuthPartner
	}
	return nil
}

func (m *LoginResponse) GetActions() []*Action {
	if m != nil {
		return m.Actions
	}
	return nil
}

func (m *LoginResponse) GetRedirectUrl() string {
	if m != nil {
		return m.RedirectUrl
	}
	return ""
}

type Action struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name" json:"name"`
	Label                string            `protobuf:"bytes,2,opt,name=label" json:"label"`
	Msg                  string            `protobuf:"bytes,3,opt,name=msg" json:"msg"`
	Meta                 map[string]string `protobuf:"bytes,4,rep,name=meta" json:"meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}
func (*Action) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{2}
}

var xxx_messageInfo_Action proto.InternalMessageInfo

func (m *Action) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Action) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Action) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *Action) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

type PartnerUserLogin struct {
	// @required
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	FullName string `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	// @required
	ShortName string `protobuf:"bytes,3,opt,name=short_name,json=shortName" json:"short_name"`
	// @required
	Phone string `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	// @required
	Email                string   `protobuf:"bytes,5,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PartnerUserLogin) Reset()         { *m = PartnerUserLogin{} }
func (m *PartnerUserLogin) String() string { return proto.CompactTextString(m) }
func (*PartnerUserLogin) ProtoMessage()    {}
func (*PartnerUserLogin) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{3}
}

var xxx_messageInfo_PartnerUserLogin proto.InternalMessageInfo

func (m *PartnerUserLogin) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PartnerUserLogin) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *PartnerUserLogin) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

func (m *PartnerUserLogin) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *PartnerUserLogin) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type PartnerShopLoginAccount struct {
	// @required
	Id         dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	ExternalId string `protobuf:"bytes,10,opt,name=external_id,json=externalId" json:"external_id"`
	// @required
	Name string `protobuf:"bytes,2,opt,name=name" json:"name"`
	// @required
	Type etop.AccountType `protobuf:"varint,3,opt,name=type,enum=etop.AccountType" json:"type"`
	// Associated token for the account. It's returned when calling Login or
	// SwitchAccount with regenerate_tokens set to true.
	AccessToken string `protobuf:"bytes,5,opt,name=access_token,json=accessToken" json:"access_token"`
	// The same as access_token.
	ExpiresIn            int32    `protobuf:"varint,6,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	ImageUrl             string   `protobuf:"bytes,7,opt,name=image_url,json=imageUrl" json:"image_url"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PartnerShopLoginAccount) Reset()         { *m = PartnerShopLoginAccount{} }
func (m *PartnerShopLoginAccount) String() string { return proto.CompactTextString(m) }
func (*PartnerShopLoginAccount) ProtoMessage()    {}
func (*PartnerShopLoginAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{4}
}

var xxx_messageInfo_PartnerShopLoginAccount proto.InternalMessageInfo

func (m *PartnerShopLoginAccount) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PartnerShopLoginAccount) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *PartnerShopLoginAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PartnerShopLoginAccount) GetType() etop.AccountType {
	if m != nil {
		return m.Type
	}
	return etop.AccountType_unknown
}

func (m *PartnerShopLoginAccount) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *PartnerShopLoginAccount) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func (m *PartnerShopLoginAccount) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

type PartnerShopInfo struct {
	Id      dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Name    string         `protobuf:"bytes,2,opt,name=name" json:"name"`
	Status  status3.Status `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	IsTest  bool           `protobuf:"varint,4,opt,name=is_test,json=isTest" json:"is_test"`
	Address *etop.Address  `protobuf:"bytes,5,opt,name=address" json:"address"`
	Phone   string         `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	//    optional string website_url = 14 [(gogoproto.nullable)=false];
	ImageUrl string `protobuf:"bytes,15,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email    string `protobuf:"bytes,16,opt,name=email" json:"email"`
	//    optional dot.ID product_source_id = 17 [(gogoproto.nullable)=false];
	//    optional dot.ID ship_to_address_id = 18 [(gogoproto.nullable)=false];
	ShipFromAddressId    dot.ID   `protobuf:"varint,19,opt,name=ship_from_address_id,json=shipFromAddressId" json:"ship_from_address_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PartnerShopInfo) Reset()         { *m = PartnerShopInfo{} }
func (m *PartnerShopInfo) String() string { return proto.CompactTextString(m) }
func (*PartnerShopInfo) ProtoMessage()    {}
func (*PartnerShopInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{5}
}

var xxx_messageInfo_PartnerShopInfo proto.InternalMessageInfo

func (m *PartnerShopInfo) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PartnerShopInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PartnerShopInfo) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *PartnerShopInfo) GetIsTest() bool {
	if m != nil {
		return m.IsTest
	}
	return false
}

func (m *PartnerShopInfo) GetAddress() *etop.Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *PartnerShopInfo) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *PartnerShopInfo) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *PartnerShopInfo) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *PartnerShopInfo) GetShipFromAddressId() dot.ID {
	if m != nil {
		return m.ShipFromAddressId
	}
	return 0
}

type RequestLoginRequest struct {
	// @required Phone or email
	Login string `protobuf:"bytes,1,opt,name=login" json:"login"`
	// @required
	RecaptchaToken       string   `protobuf:"bytes,2,opt,name=recaptcha_token,json=recaptchaToken" json:"recaptcha_token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestLoginRequest) Reset()         { *m = RequestLoginRequest{} }
func (m *RequestLoginRequest) String() string { return proto.CompactTextString(m) }
func (*RequestLoginRequest) ProtoMessage()    {}
func (*RequestLoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{6}
}

var xxx_messageInfo_RequestLoginRequest proto.InternalMessageInfo

func (m *RequestLoginRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *RequestLoginRequest) GetRecaptchaToken() string {
	if m != nil {
		return m.RecaptchaToken
	}
	return ""
}

type RequestLoginResponse struct {
	// @required
	Code string `protobuf:"bytes,1,opt,name=code" json:"code"`
	// @required
	Msg string `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	// @required
	Actions              []*Action `protobuf:"bytes,3,rep,name=actions" json:"actions"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *RequestLoginResponse) Reset()         { *m = RequestLoginResponse{} }
func (m *RequestLoginResponse) String() string { return proto.CompactTextString(m) }
func (*RequestLoginResponse) ProtoMessage()    {}
func (*RequestLoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{7}
}

var xxx_messageInfo_RequestLoginResponse proto.InternalMessageInfo

func (m *RequestLoginResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *RequestLoginResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *RequestLoginResponse) GetActions() []*Action {
	if m != nil {
		return m.Actions
	}
	return nil
}

type LoginUsingTokenRequest struct {
	Login                string   `protobuf:"bytes,1,opt,name=login" json:"login"`
	VerificationCode     string   `protobuf:"bytes,2,opt,name=verification_code,json=verificationCode" json:"verification_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginUsingTokenRequest) Reset()         { *m = LoginUsingTokenRequest{} }
func (m *LoginUsingTokenRequest) String() string { return proto.CompactTextString(m) }
func (*LoginUsingTokenRequest) ProtoMessage()    {}
func (*LoginUsingTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{8}
}

var xxx_messageInfo_LoginUsingTokenRequest proto.InternalMessageInfo

func (m *LoginUsingTokenRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *LoginUsingTokenRequest) GetVerificationCode() string {
	if m != nil {
		return m.VerificationCode
	}
	return ""
}

type RegisterRequest struct {
	FullName             string   `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	AgreeTos             bool     `protobuf:"varint,5,opt,name=agree_tos,json=agreeTos" json:"agree_tos"`
	AgreeEmailInfo       *bool    `protobuf:"varint,6,opt,name=agree_email_info,json=agreeEmailInfo" json:"agree_email_info"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{9}
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *RegisterRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *RegisterRequest) GetAgreeTos() bool {
	if m != nil {
		return m.AgreeTos
	}
	return false
}

func (m *RegisterRequest) GetAgreeEmailInfo() bool {
	if m != nil && m.AgreeEmailInfo != nil {
		return *m.AgreeEmailInfo
	}
	return false
}

type RegisterResponse struct {
	User                 *etop.User `protobuf:"bytes,1,opt,name=user" json:"user"`
	AccessToken          string     `protobuf:"bytes,2,opt,name=access_token,json=accessToken" json:"access_token"`
	ExpiresIn            int32      `protobuf:"varint,3,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{10}
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetUser() *etop.User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *RegisterResponse) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *RegisterResponse) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

type GrantAccessRequest struct {
	ShopId               dot.ID   `protobuf:"varint,1,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrantAccessRequest) Reset()         { *m = GrantAccessRequest{} }
func (m *GrantAccessRequest) String() string { return proto.CompactTextString(m) }
func (*GrantAccessRequest) ProtoMessage()    {}
func (*GrantAccessRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{11}
}

var xxx_messageInfo_GrantAccessRequest proto.InternalMessageInfo

func (m *GrantAccessRequest) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

type GrantAccessResponse struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token"`
	ExpiresIn            int32    `protobuf:"varint,2,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrantAccessResponse) Reset()         { *m = GrantAccessResponse{} }
func (m *GrantAccessResponse) String() string { return proto.CompactTextString(m) }
func (*GrantAccessResponse) ProtoMessage()    {}
func (*GrantAccessResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_af618ab6ed9ae35a, []int{12}
}

var xxx_messageInfo_GrantAccessResponse proto.InternalMessageInfo

func (m *GrantAccessResponse) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *GrantAccessResponse) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func init() {
	proto.RegisterType((*InitRequest)(nil), "integration.InitRequest")
	proto.RegisterType((*LoginResponse)(nil), "integration.LoginResponse")
	proto.RegisterType((*Action)(nil), "integration.Action")
	proto.RegisterMapType((map[string]string)(nil), "integration.Action.MetaEntry")
	proto.RegisterType((*PartnerUserLogin)(nil), "integration.PartnerUserLogin")
	proto.RegisterType((*PartnerShopLoginAccount)(nil), "integration.PartnerShopLoginAccount")
	proto.RegisterType((*PartnerShopInfo)(nil), "integration.PartnerShopInfo")
	proto.RegisterType((*RequestLoginRequest)(nil), "integration.RequestLoginRequest")
	proto.RegisterType((*RequestLoginResponse)(nil), "integration.RequestLoginResponse")
	proto.RegisterType((*LoginUsingTokenRequest)(nil), "integration.LoginUsingTokenRequest")
	proto.RegisterType((*RegisterRequest)(nil), "integration.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "integration.RegisterResponse")
	proto.RegisterType((*GrantAccessRequest)(nil), "integration.GrantAccessRequest")
	proto.RegisterType((*GrantAccessResponse)(nil), "integration.GrantAccessResponse")
}

func init() { proto.RegisterFile("etop/integration/integration.proto", fileDescriptor_af618ab6ed9ae35a) }

var fileDescriptor_af618ab6ed9ae35a = []byte{
	// 1049 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0x4f, 0x6f, 0x1b, 0x45,
	0x14, 0xdf, 0xb5, 0x9d, 0x38, 0x7e, 0x4e, 0x62, 0x67, 0x12, 0xb5, 0xab, 0x88, 0x6c, 0x93, 0x2d,
	0xa8, 0x91, 0x50, 0x6c, 0x9a, 0x0a, 0x81, 0x7a, 0x40, 0x4a, 0x51, 0x41, 0x96, 0xa0, 0x2a, 0x9b,
	0xe4, 0xc2, 0x65, 0x99, 0xec, 0x4e, 0xec, 0x51, 0xd7, 0x3b, 0xcb, 0xcc, 0x38, 0x34, 0x37, 0x8e,
	0x9c, 0x10, 0x47, 0x3e, 0x42, 0x3f, 0x02, 0x5c, 0x39, 0x45, 0x9c, 0xfa, 0x09, 0x50, 0xe3, 0x1e,
	0xb9, 0xf0, 0x11, 0xd0, 0xfc, 0x59, 0x67, 0xed, 0x38, 0x90, 0x4b, 0x4f, 0xbb, 0xfe, 0xbd, 0xf7,
	0x66, 0xde, 0x9f, 0xdf, 0xfb, 0xad, 0x21, 0x20, 0x92, 0xe5, 0x5d, 0x9a, 0x49, 0xd2, 0xe7, 0x58,
	0x52, 0x96, 0x95, 0xdf, 0x3b, 0x39, 0x67, 0x92, 0xa1, 0x66, 0x09, 0xda, 0xdc, 0xe8, 0xb3, 0x3e,
	0xd3, 0x78, 0x57, 0xbd, 0x19, 0x97, 0xcd, 0xf5, 0x98, 0x0d, 0x87, 0x2c, 0xeb, 0x9a, 0x87, 0x05,
	0xb7, 0xf4, 0xd9, 0x44, 0xc6, 0x5d, 0x21, 0xb1, 0x1c, 0x89, 0x47, 0xf6, 0x69, 0xcd, 0x2d, 0x6b,
	0x66, 0xb9, 0x01, 0x82, 0x7d, 0x68, 0xf6, 0x32, 0x2a, 0x43, 0xf2, 0xfd, 0x88, 0x08, 0x89, 0xee,
	0x03, 0xe0, 0x91, 0x1c, 0x44, 0x92, 0xbd, 0x20, 0x99, 0xe7, 0x6e, 0xbb, 0xbb, 0x8d, 0x27, 0xb5,
	0x8b, 0xbf, 0xee, 0x39, 0x61, 0x43, 0xe1, 0x47, 0x0a, 0x0e, 0xfe, 0xae, 0xc2, 0xca, 0x57, 0xac,
	0x4f, 0xb3, 0x90, 0x88, 0x9c, 0x65, 0x82, 0xa0, 0x07, 0xb0, 0x8c, 0xe3, 0x98, 0x08, 0x31, 0x27,
	0xb0, 0x69, 0x2c, 0x3a, 0x54, 0x9d, 0x4f, 0x5e, 0xe6, 0x94, 0x13, 0x11, 0xd1, 0xcc, 0xab, 0x6c,
	0xbb, 0xbb, 0x0b, 0xc5, 0xf9, 0x16, 0xef, 0x65, 0xe8, 0x21, 0xd4, 0x46, 0x82, 0x70, 0xaf, 0xba,
	0xed, 0xee, 0x36, 0xf7, 0xb7, 0x3a, 0xe5, 0xee, 0x3c, 0xc7, 0x5c, 0x66, 0x84, 0x1f, 0x0b, 0xc2,
	0x4d, 0x0a, 0xda, 0x15, 0x7d, 0x06, 0x75, 0x1c, 0xc7, 0x6c, 0x94, 0x49, 0xaf, 0xa6, 0xa3, 0xde,
	0x9f, 0x17, 0x75, 0x38, 0x60, 0xb9, 0x8e, 0x3a, 0x30, 0xbe, 0x61, 0x11, 0x84, 0x3e, 0x82, 0x9a,
	0x18, 0xb0, 0xdc, 0x5b, 0xd0, 0xc1, 0xef, 0xdd, 0x14, 0xdc, 0xcb, 0x4e, 0x59, 0xa8, 0x3d, 0xd1,
	0x21, 0x20, 0x7c, 0x86, 0x69, 0x8a, 0x4f, 0x52, 0x12, 0xd9, 0x63, 0x84, 0x57, 0xdf, 0xae, 0xde,
	0xfa, 0xf2, 0xb5, 0x49, 0xbc, 0x45, 0x04, 0x7a, 0x0c, 0xcb, 0xba, 0xfd, 0xb9, 0x09, 0xf1, 0x9a,
	0x3a, 0x9d, 0xbb, 0x1d, 0x3d, 0xb0, 0xe7, 0xa3, 0x93, 0x94, 0xc6, 0xd6, 0x57, 0x67, 0xd2, 0x54,
	0xce, 0xf6, 0x78, 0xb4, 0xa7, 0x5a, 0xa0, 0x2e, 0x14, 0xde, 0xb2, 0xce, 0x62, 0x7d, 0x2a, 0x8b,
	0x03, 0x6d, 0x0b, 0x0b, 0x1f, 0x35, 0x32, 0x4e, 0x12, 0xca, 0x49, 0x2c, 0xa3, 0x11, 0x4f, 0xbd,
	0x95, 0xf2, 0xc8, 0x0a, 0xcb, 0x31, 0x4f, 0x83, 0x3f, 0x5c, 0x58, 0x34, 0xc1, 0xc8, 0x83, 0x5a,
	0x86, 0x87, 0x64, 0x6a, 0xbc, 0x1a, 0x41, 0x9b, 0xb0, 0x90, 0xe2, 0x13, 0x92, 0xea, 0x91, 0x16,
	0x26, 0x03, 0xa1, 0x3b, 0x50, 0x1d, 0x8a, 0xbe, 0x9e, 0x66, 0x61, 0x51, 0x80, 0x1a, 0xf3, 0x90,
	0x48, 0xec, 0xd5, 0x74, 0xb6, 0x5b, 0x73, 0xb2, 0xed, 0x7c, 0x4d, 0x24, 0x7e, 0x9a, 0x49, 0x7e,
	0x1e, 0x6a, 0xd7, 0xcd, 0x4f, 0xa0, 0x31, 0x81, 0x50, 0x1b, 0xaa, 0x2f, 0xc8, 0xb9, 0x49, 0x26,
	0x54, 0xaf, 0x68, 0x03, 0x16, 0xce, 0x70, 0x3a, 0x22, 0x26, 0x8b, 0xd0, 0xfc, 0x78, 0x5c, 0xf9,
	0xd4, 0x0d, 0x5e, 0xb9, 0xd0, 0x9e, 0xa5, 0x0e, 0xda, 0x80, 0x0a, 0x4d, 0x74, 0x7c, 0xd5, 0xe6,
	0x55, 0xa1, 0x09, 0xda, 0x81, 0xc6, 0xe9, 0x28, 0x4d, 0x23, 0x5d, 0x69, 0xb9, 0x9c, 0x25, 0x05,
	0x3f, 0x53, 0xd5, 0xde, 0x07, 0x10, 0x03, 0xc6, 0xa5, 0xf1, 0x29, 0x17, 0xd6, 0xd0, 0xf8, 0x33,
	0xdb, 0x92, 0x7c, 0xc0, 0x32, 0xa2, 0x09, 0x39, 0x69, 0x89, 0x86, 0x94, 0x8d, 0x0c, 0x31, 0x4d,
	0x35, 0xdf, 0x26, 0x36, 0x0d, 0x05, 0x3f, 0x57, 0xe0, 0xee, 0x0d, 0x94, 0xb9, 0x21, 0xe3, 0x0f,
	0xa0, 0x49, 0x5e, 0x4a, 0xc2, 0x33, 0x9c, 0x46, 0x34, 0xf1, 0xa0, 0x74, 0x26, 0x14, 0x86, 0x5e,
	0x32, 0x99, 0x5e, 0xe5, 0xda, 0xf4, 0x3e, 0x84, 0x9a, 0x3c, 0xcf, 0x4d, 0x25, 0xab, 0xfb, 0x6b,
	0x86, 0x6e, 0xf6, 0xce, 0xa3, 0xf3, 0x9c, 0x14, 0xce, 0xca, 0xe9, 0xda, 0xae, 0x2f, 0xdc, 0x6e,
	0xd7, 0x17, 0xe7, 0xef, 0xfa, 0x0e, 0x34, 0xe8, 0x10, 0xf7, 0x89, 0xe6, 0x60, 0xbd, 0xdc, 0x6d,
	0x0d, 0x2b, 0x02, 0xfe, 0x59, 0x81, 0xd6, 0xcc, 0x0e, 0xde, 0xd0, 0x88, 0x9b, 0x2b, 0xdc, 0x83,
	0x45, 0xa3, 0x83, 0xb6, 0xc6, 0x56, 0xc7, 0xca, 0x63, 0xe7, 0x50, 0x3f, 0xad, 0xb3, 0x75, 0x42,
	0x5b, 0x50, 0xa7, 0x22, 0x92, 0x44, 0x18, 0x39, 0x59, 0x2a, 0xcc, 0x54, 0x1c, 0x29, 0x95, 0x7c,
	0x00, 0x75, 0x9c, 0x24, 0x9c, 0x08, 0x61, 0x05, 0x63, 0xc5, 0xb6, 0xcc, 0x80, 0x61, 0x61, 0xbd,
	0xe2, 0x40, 0xe3, 0x3a, 0x07, 0xa6, 0x2a, 0x6f, 0xcd, 0xab, 0xfc, 0x8a, 0x26, 0xed, 0x6b, 0x34,
	0x41, 0x1f, 0xc3, 0x86, 0x18, 0xd0, 0x3c, 0x3a, 0xe5, 0x6c, 0x18, 0xd9, 0xfb, 0xd4, 0xf4, 0xd7,
	0x4b, 0x3d, 0x59, 0x53, 0x1e, 0x5f, 0x70, 0x36, 0xb4, 0x89, 0xf5, 0x92, 0xe0, 0x3b, 0x58, 0xb7,
	0x5a, 0x6f, 0x15, 0xdc, 0xe8, 0xbe, 0xda, 0x5f, 0xf5, 0x7b, 0x6a, 0xb5, 0x0d, 0x84, 0xf6, 0xa0,
	0xc5, 0x49, 0x8c, 0x73, 0x19, 0x0f, 0xb0, 0x9d, 0x79, 0xb9, 0xc1, 0xab, 0x13, 0xa3, 0xf9, 0x3a,
	0xfc, 0x00, 0x1b, 0xd3, 0x37, 0xd8, 0x6f, 0x84, 0x07, 0xb5, 0x98, 0x25, 0x33, 0xe2, 0xa1, 0x90,
	0x42, 0x20, 0x2a, 0xb3, 0x02, 0x51, 0x52, 0xb4, 0xea, 0xff, 0x2b, 0x5a, 0xd0, 0x87, 0x3b, 0xfa,
	0xc6, 0x63, 0x41, 0xb3, 0xbe, 0xce, 0xe5, 0x36, 0xd5, 0x3d, 0x84, 0xb5, 0x33, 0xc2, 0xe9, 0x29,
	0x8d, 0xf5, 0xa9, 0x91, 0xce, 0xb1, 0x9c, 0x4a, 0xbb, 0x6c, 0xfe, 0x9c, 0x25, 0x24, 0xf8, 0xdd,
	0x85, 0x56, 0x48, 0xfa, 0x54, 0x48, 0xc2, 0x8b, 0x2b, 0xa6, 0x54, 0xc3, 0x9d, 0xab, 0x1a, 0x13,
	0x32, 0x54, 0xff, 0x43, 0x10, 0x6a, 0xd7, 0x27, 0xbd, 0x03, 0x0d, 0xdc, 0xe7, 0x84, 0x44, 0x92,
	0x19, 0xbe, 0x15, 0x74, 0x5c, 0xd2, 0xf0, 0x11, 0x13, 0x68, 0x17, 0xda, 0xc6, 0x45, 0x47, 0x44,
	0x34, 0x3b, 0x65, 0x7a, 0xe1, 0x96, 0xc2, 0x55, 0x8d, 0x3f, 0x55, 0xb0, 0x5a, 0x9c, 0xe0, 0x47,
	0x17, 0xda, 0x57, 0xb9, 0xdb, 0xd1, 0xf8, 0xf6, 0x83, 0xeb, 0x6a, 0x32, 0x83, 0x21, 0xb3, 0xd2,
	0x49, 0xfb, 0x75, 0x9d, 0x5d, 0xf9, 0xca, 0xed, 0x56, 0xbe, 0x3a, 0x77, 0xe5, 0x83, 0x47, 0x80,
	0xbe, 0xe4, 0x38, 0x93, 0x07, 0x3a, 0xb0, 0x68, 0xe0, 0x16, 0xd4, 0xd5, 0x77, 0x35, 0x9a, 0x59,
	0xeb, 0x45, 0x05, 0xf6, 0x92, 0x20, 0x86, 0xf5, 0xa9, 0xa0, 0x77, 0xf1, 0xc7, 0xe3, 0xc9, 0x37,
	0x17, 0x97, 0xbe, 0xfb, 0xfa, 0xd2, 0x77, 0xdf, 0x5c, 0xfa, 0xce, 0x3f, 0x97, 0xbe, 0xf3, 0xd3,
	0xd8, 0x77, 0x5e, 0x8d, 0x7d, 0xe7, 0xb7, 0xb1, 0xef, 0x5c, 0x8c, 0x7d, 0xe7, 0xf5, 0xd8, 0x77,
	0xde, 0x8c, 0x7d, 0xe7, 0x97, 0xb7, 0xbe, 0xf3, 0xeb, 0x5b, 0xdf, 0xf9, 0xf6, 0x9e, 0x6e, 0xd5,
	0x59, 0xd6, 0xc5, 0x39, 0xed, 0xe6, 0x27, 0xdd, 0xd9, 0x7f, 0x76, 0xff, 0x06, 0x00, 0x00, 0xff,
	0xff, 0xd4, 0x12, 0xa6, 0x0d, 0xec, 0x09, 0x00, 0x00,
}
