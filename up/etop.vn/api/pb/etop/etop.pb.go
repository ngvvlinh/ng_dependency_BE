// source: etop/etop.proto

package etop

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	address_type "etop.vn/api/pb/etop/etc/address_type"
	status3 "etop.vn/api/pb/etop/etc/status3"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	user_source "etop.vn/api/pb/etop/etc/user_source"
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

// Indicates whether given account is **etop**, **shop**, **partner** or **sale**.
type AccountType int32

const (
	AccountType_unknown   AccountType = 0
	AccountType_partner   AccountType = 21
	AccountType_shop      AccountType = 33
	AccountType_affiliate AccountType = 35
	AccountType_etop      AccountType = 101
)

var AccountType_name = map[int32]string{
	0:   "unknown",
	21:  "partner",
	33:  "shop",
	35:  "affiliate",
	101: "etop",
}

var AccountType_value = map[string]int32{
	"unknown":   0,
	"partner":   21,
	"shop":      33,
	"affiliate": 35,
	"etop":      101,
}

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return proto.EnumName(AccountType_name, int32(x))
}

func (x *AccountType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AccountType_value, data, "AccountType")
	if err != nil {
		return err
	}
	*x = AccountType(value)
	return nil
}

func (AccountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{0}
}

type GetInvitationsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetInvitationsRequest) Reset()         { *m = GetInvitationsRequest{} }
func (m *GetInvitationsRequest) String() string { return proto.CompactTextString(m) }
func (*GetInvitationsRequest) ProtoMessage()    {}
func (*GetInvitationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{0}
}

var xxx_messageInfo_GetInvitationsRequest proto.InternalMessageInfo

func (m *GetInvitationsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetInvitationsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetInvitationByTokenRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token" json:"token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInvitationByTokenRequest) Reset()         { *m = GetInvitationByTokenRequest{} }
func (m *GetInvitationByTokenRequest) String() string { return proto.CompactTextString(m) }
func (*GetInvitationByTokenRequest) ProtoMessage()    {}
func (*GetInvitationByTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{1}
}

var xxx_messageInfo_GetInvitationByTokenRequest proto.InternalMessageInfo

func (m *GetInvitationByTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Invitation struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID         `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Email                string         `protobuf:"bytes,3,opt,name=email" json:"email"`
	Roles                []string       `protobuf:"bytes,4,rep,name=roles" json:"roles"`
	Token                string         `protobuf:"bytes,5,opt,name=token" json:"token"`
	Status               status3.Status `protobuf:"varint,6,opt,name=status,enum=status3.Status" json:"status"`
	InvitedBy            dot.ID         `protobuf:"varint,7,opt,name=invited_by,json=invitedBy" json:"invited_by"`
	AcceptedAt           dot.Time       `protobuf:"bytes,8,opt,name=accepted_at,json=acceptedAt" json:"accepted_at"`
	DeclinedAt           dot.Time       `protobuf:"bytes,9,opt,name=declined_at,json=declinedAt" json:"declined_at"`
	ExpiredAt            dot.Time       `protobuf:"bytes,10,opt,name=expired_at,json=expiredAt" json:"expired_at"`
	CreatedAt            dot.Time       `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time       `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Invitation) Reset()         { *m = Invitation{} }
func (m *Invitation) String() string { return proto.CompactTextString(m) }
func (*Invitation) ProtoMessage()    {}
func (*Invitation) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{2}
}

var xxx_messageInfo_Invitation proto.InternalMessageInfo

func (m *Invitation) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Invitation) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Invitation) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Invitation) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *Invitation) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Invitation) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Invitation) GetInvitedBy() dot.ID {
	if m != nil {
		return m.InvitedBy
	}
	return 0
}

type AcceptInvitationRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token" json:"token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptInvitationRequest) Reset()         { *m = AcceptInvitationRequest{} }
func (m *AcceptInvitationRequest) String() string { return proto.CompactTextString(m) }
func (*AcceptInvitationRequest) ProtoMessage()    {}
func (*AcceptInvitationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{3}
}

var xxx_messageInfo_AcceptInvitationRequest proto.InternalMessageInfo

func (m *AcceptInvitationRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type RejectInvitationRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token" json:"token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RejectInvitationRequest) Reset()         { *m = RejectInvitationRequest{} }
func (m *RejectInvitationRequest) String() string { return proto.CompactTextString(m) }
func (*RejectInvitationRequest) ProtoMessage()    {}
func (*RejectInvitationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{4}
}

var xxx_messageInfo_RejectInvitationRequest proto.InternalMessageInfo

func (m *RejectInvitationRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

// Represents a user in eTop system. The user may or may not have associated accounts.
type User struct {
	// @required
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	FullName string `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	// @required
	ShortName string `protobuf:"bytes,3,opt,name=short_name,json=shortName" json:"short_name"`
	// @required
	Phone string `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	// @required
	Email string `protobuf:"bytes,5,opt,name=email" json:"email"`
	// @required
	CreatedAt dot.Time `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	// @required
	UpdatedAt               dot.Time               `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	EmailVerifiedAt         dot.Time               `protobuf:"bytes,8,opt,name=email_verified_at,json=emailVerifiedAt" json:"email_verified_at"`
	PhoneVerifiedAt         dot.Time               `protobuf:"bytes,9,opt,name=phone_verified_at,json=phoneVerifiedAt" json:"phone_verified_at"`
	EmailVerificationSentAt dot.Time               `protobuf:"bytes,10,opt,name=email_verification_sent_at,json=emailVerificationSentAt" json:"email_verification_sent_at"`
	PhoneVerificationSentAt dot.Time               `protobuf:"bytes,11,opt,name=phone_verification_sent_at,json=phoneVerificationSentAt" json:"phone_verification_sent_at"`
	Source                  user_source.UserSource `protobuf:"varint,12,opt,name=source,enum=user_source.UserSource" json:"source"`
	XXX_NoUnkeyedLiteral    struct{}               `json:"-"`
	XXX_sizecache           int32                  `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{5}
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *User) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

func (m *User) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetSource() user_source.UserSource {
	if m != nil {
		return m.Source
	}
	return user_source.UserSource_unknown
}

type IDsRequest struct {
	// @required
	Ids                  []dot.ID      `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Mixed                *MixedAccount `protobuf:"bytes,2,opt,name=mixed" json:"mixed"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *IDsRequest) Reset()         { *m = IDsRequest{} }
func (m *IDsRequest) String() string { return proto.CompactTextString(m) }
func (*IDsRequest) ProtoMessage()    {}
func (*IDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{6}
}

var xxx_messageInfo_IDsRequest proto.InternalMessageInfo

func (m *IDsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *IDsRequest) GetMixed() *MixedAccount {
	if m != nil {
		return m.Mixed
	}
	return nil
}

type MixedAccount struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	All                  bool     `protobuf:"varint,2,opt,name=all" json:"all"`
	AllShops             bool     `protobuf:"varint,3,opt,name=all_shops,json=allShops" json:"all_shops"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MixedAccount) Reset()         { *m = MixedAccount{} }
func (m *MixedAccount) String() string { return proto.CompactTextString(m) }
func (*MixedAccount) ProtoMessage()    {}
func (*MixedAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{7}
}

var xxx_messageInfo_MixedAccount proto.InternalMessageInfo

func (m *MixedAccount) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *MixedAccount) GetAll() bool {
	if m != nil {
		return m.All
	}
	return false
}

func (m *MixedAccount) GetAllShops() bool {
	if m != nil {
		return m.AllShops
	}
	return false
}

type Partner struct {
	Id                   dot.ID           `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string           `protobuf:"bytes,2,opt,name=name" json:"name"`
	PublicName           string           `protobuf:"bytes,24,opt,name=public_name,json=publicName" json:"public_name"`
	Status               status3.Status   `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	IsTest               bool             `protobuf:"varint,4,opt,name=is_test,json=isTest" json:"is_test"`
	ContactPersons       []*ContactPerson `protobuf:"bytes,8,rep,name=contact_persons,json=contactPersons" json:"contact_persons"`
	Phone                string           `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	WebsiteUrl           string           `protobuf:"bytes,14,opt,name=website_url,json=websiteUrl" json:"website_url"`
	ImageUrl             string           `protobuf:"bytes,15,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email                string           `protobuf:"bytes,16,opt,name=email" json:"email"`
	OwnerId              dot.ID           `protobuf:"varint,22,opt,name=owner_id,json=ownerId" json:"owner_id"`
	User                 *User            `protobuf:"bytes,23,opt,name=user" json:"user"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Partner) Reset()         { *m = Partner{} }
func (m *Partner) String() string { return proto.CompactTextString(m) }
func (*Partner) ProtoMessage()    {}
func (*Partner) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{8}
}

var xxx_messageInfo_Partner proto.InternalMessageInfo

func (m *Partner) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Partner) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Partner) GetPublicName() string {
	if m != nil {
		return m.PublicName
	}
	return ""
}

func (m *Partner) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Partner) GetIsTest() bool {
	if m != nil {
		return m.IsTest
	}
	return false
}

func (m *Partner) GetContactPersons() []*ContactPerson {
	if m != nil {
		return m.ContactPersons
	}
	return nil
}

func (m *Partner) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Partner) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *Partner) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *Partner) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Partner) GetOwnerId() dot.ID {
	if m != nil {
		return m.OwnerId
	}
	return 0
}

func (m *Partner) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type PublicAccountInfo struct {
	Id                   dot.ID      `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string      `protobuf:"bytes,2,opt,name=name" json:"name"`
	Type                 AccountType `protobuf:"varint,3,opt,name=type,enum=etop.AccountType" json:"type"`
	ImageUrl             string      `protobuf:"bytes,4,opt,name=image_url,json=imageUrl" json:"image_url"`
	Website              string      `protobuf:"bytes,5,opt,name=website" json:"website"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PublicAccountInfo) Reset()         { *m = PublicAccountInfo{} }
func (m *PublicAccountInfo) String() string { return proto.CompactTextString(m) }
func (*PublicAccountInfo) ProtoMessage()    {}
func (*PublicAccountInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{9}
}

var xxx_messageInfo_PublicAccountInfo proto.InternalMessageInfo

func (m *PublicAccountInfo) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PublicAccountInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PublicAccountInfo) GetType() AccountType {
	if m != nil {
		return m.Type
	}
	return AccountType_unknown
}

func (m *PublicAccountInfo) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *PublicAccountInfo) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

type PublicAuthorizedPartnerInfo struct {
	Id                   dot.ID      `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string      `protobuf:"bytes,2,opt,name=name" json:"name"`
	Type                 AccountType `protobuf:"varint,3,opt,name=type,enum=etop.AccountType" json:"type"`
	ImageUrl             string      `protobuf:"bytes,4,opt,name=image_url,json=imageUrl" json:"image_url"`
	Website              string      `protobuf:"bytes,5,opt,name=website" json:"website"`
	RedirectUrl          string      `protobuf:"bytes,6,opt,name=redirect_url,json=redirectUrl" json:"redirect_url"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PublicAuthorizedPartnerInfo) Reset()         { *m = PublicAuthorizedPartnerInfo{} }
func (m *PublicAuthorizedPartnerInfo) String() string { return proto.CompactTextString(m) }
func (*PublicAuthorizedPartnerInfo) ProtoMessage()    {}
func (*PublicAuthorizedPartnerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{10}
}

var xxx_messageInfo_PublicAuthorizedPartnerInfo proto.InternalMessageInfo

func (m *PublicAuthorizedPartnerInfo) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PublicAuthorizedPartnerInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PublicAuthorizedPartnerInfo) GetType() AccountType {
	if m != nil {
		return m.Type
	}
	return AccountType_unknown
}

func (m *PublicAuthorizedPartnerInfo) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *PublicAuthorizedPartnerInfo) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *PublicAuthorizedPartnerInfo) GetRedirectUrl() string {
	if m != nil {
		return m.RedirectUrl
	}
	return ""
}

type PublicAccountsResponse struct {
	Accounts             []*PublicAccountInfo `protobuf:"bytes,1,rep,name=accounts" json:"accounts"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PublicAccountsResponse) Reset()         { *m = PublicAccountsResponse{} }
func (m *PublicAccountsResponse) String() string { return proto.CompactTextString(m) }
func (*PublicAccountsResponse) ProtoMessage()    {}
func (*PublicAccountsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{11}
}

var xxx_messageInfo_PublicAccountsResponse proto.InternalMessageInfo

func (m *PublicAccountsResponse) GetAccounts() []*PublicAccountInfo {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type Shop struct {
	ExportedFields     []string       `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	InventoryOverstock bool           `protobuf:"varint,30,opt,name=inventory_overstock,json=inventoryOverstock" json:"inventory_overstock"`
	Id                 dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Name               string         `protobuf:"bytes,2,opt,name=name" json:"name"`
	Status             status3.Status `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	IsTest             bool           `protobuf:"varint,4,opt,name=is_test,json=isTest" json:"is_test"`
	Address            *Address       `protobuf:"bytes,5,opt,name=address" json:"address"`
	Phone              string         `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	BankAccount        *BankAccount   `protobuf:"bytes,10,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	AutoCreateFfm      bool           `protobuf:"varint,20,opt,name=auto_create_ffm,json=autoCreateFfm" json:"auto_create_ffm"`
	WebsiteUrl         string         `protobuf:"bytes,14,opt,name=website_url,json=websiteUrl" json:"website_url"`
	ImageUrl           string         `protobuf:"bytes,15,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email              string         `protobuf:"bytes,16,opt,name=email" json:"email"`
	ProductSourceId    dot.ID         `protobuf:"varint,17,opt,name=product_source_id,json=productSourceId" json:"product_source_id"`
	ShipToAddressId    dot.ID         `protobuf:"varint,18,opt,name=ship_to_address_id,json=shipToAddressId" json:"ship_to_address_id"`
	ShipFromAddressId  dot.ID         `protobuf:"varint,19,opt,name=ship_from_address_id,json=shipFromAddressId" json:"ship_from_address_id"`
	// @deprecated use try_on instead
	GhnNoteCode string           `protobuf:"bytes,21,opt,name=ghn_note_code,json=ghnNoteCode" json:"ghn_note_code"`
	TryOn       try_on.TryOnCode `protobuf:"varint,24,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	OwnerId     dot.ID           `protobuf:"varint,22,opt,name=owner_id,json=ownerId" json:"owner_id"`
	User        *User            `protobuf:"bytes,23,opt,name=user" json:"user"`
	CompanyInfo *CompanyInfo     `protobuf:"bytes,25,opt,name=company_info,json=companyInfo" json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                               `protobuf:"bytes,26,opt,name=money_transaction_rrule,json=moneyTransactionRrule" json:"money_transaction_rrule"`
	SurveyInfo                    []*SurveyInfo                        `protobuf:"bytes,27,rep,name=survey_info,json=surveyInfo" json:"survey_info"`
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem `protobuf:"bytes,28,rep,name=shipping_service_select_strategy,json=shippingServiceSelectStrategy" json:"shipping_service_select_strategy"`
	Code                          string                               `protobuf:"bytes,29,opt,name=code" json:"code"`
	XXX_NoUnkeyedLiteral          struct{}                             `json:"-"`
	XXX_sizecache                 int32                                `json:"-"`
}

func (m *Shop) Reset()         { *m = Shop{} }
func (m *Shop) String() string { return proto.CompactTextString(m) }
func (*Shop) ProtoMessage()    {}
func (*Shop) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{12}
}

var xxx_messageInfo_Shop proto.InternalMessageInfo

func (m *Shop) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *Shop) GetInventoryOverstock() bool {
	if m != nil {
		return m.InventoryOverstock
	}
	return false
}

func (m *Shop) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Shop) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Shop) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Shop) GetIsTest() bool {
	if m != nil {
		return m.IsTest
	}
	return false
}

func (m *Shop) GetAddress() *Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Shop) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Shop) GetBankAccount() *BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *Shop) GetAutoCreateFfm() bool {
	if m != nil {
		return m.AutoCreateFfm
	}
	return false
}

func (m *Shop) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *Shop) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *Shop) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Shop) GetProductSourceId() dot.ID {
	if m != nil {
		return m.ProductSourceId
	}
	return 0
}

func (m *Shop) GetShipToAddressId() dot.ID {
	if m != nil {
		return m.ShipToAddressId
	}
	return 0
}

func (m *Shop) GetShipFromAddressId() dot.ID {
	if m != nil {
		return m.ShipFromAddressId
	}
	return 0
}

func (m *Shop) GetGhnNoteCode() string {
	if m != nil {
		return m.GhnNoteCode
	}
	return ""
}

func (m *Shop) GetTryOn() try_on.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *Shop) GetOwnerId() dot.ID {
	if m != nil {
		return m.OwnerId
	}
	return 0
}

func (m *Shop) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Shop) GetCompanyInfo() *CompanyInfo {
	if m != nil {
		return m.CompanyInfo
	}
	return nil
}

func (m *Shop) GetMoneyTransactionRrule() string {
	if m != nil {
		return m.MoneyTransactionRrule
	}
	return ""
}

func (m *Shop) GetSurveyInfo() []*SurveyInfo {
	if m != nil {
		return m.SurveyInfo
	}
	return nil
}

func (m *Shop) GetShippingServiceSelectStrategy() []*ShippingServiceSelectStrategyItem {
	if m != nil {
		return m.ShippingServiceSelectStrategy
	}
	return nil
}

func (m *Shop) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type ShippingServiceSelectStrategyItem struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key"`
	Value                string   `protobuf:"bytes,2,opt,name=value" json:"value"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShippingServiceSelectStrategyItem) Reset()         { *m = ShippingServiceSelectStrategyItem{} }
func (m *ShippingServiceSelectStrategyItem) String() string { return proto.CompactTextString(m) }
func (*ShippingServiceSelectStrategyItem) ProtoMessage()    {}
func (*ShippingServiceSelectStrategyItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{13}
}

var xxx_messageInfo_ShippingServiceSelectStrategyItem proto.InternalMessageInfo

func (m *ShippingServiceSelectStrategyItem) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ShippingServiceSelectStrategyItem) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SurveyInfo struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key"`
	Question             string   `protobuf:"bytes,2,opt,name=question" json:"question"`
	Answer               string   `protobuf:"bytes,3,opt,name=answer" json:"answer"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SurveyInfo) Reset()         { *m = SurveyInfo{} }
func (m *SurveyInfo) String() string { return proto.CompactTextString(m) }
func (*SurveyInfo) ProtoMessage()    {}
func (*SurveyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{14}
}

var xxx_messageInfo_SurveyInfo proto.InternalMessageInfo

func (m *SurveyInfo) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SurveyInfo) GetQuestion() string {
	if m != nil {
		return m.Question
	}
	return ""
}

func (m *SurveyInfo) GetAnswer() string {
	if m != nil {
		return m.Answer
	}
	return ""
}

type CreateUserRequest struct {
	// @required
	FullName string `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	// Can be automatically deduce from full_name.
	ShortName string `protobuf:"bytes,2,opt,name=short_name,json=shortName" json:"short_name"`
	// @required
	Phone string `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	// It's not required if the user provides register_token
	Email string `protobuf:"bytes,4,opt,name=email" json:"email"`
	// @required
	Password string `protobuf:"bytes,5,opt,name=password" json:"password"`
	// @required
	AgreeTos bool `protobuf:"varint,6,opt,name=agree_tos,json=agreeTos" json:"agree_tos"`
	// @required
	AgreeEmailInfo *bool `protobuf:"varint,7,opt,name=agree_email_info,json=agreeEmailInfo" json:"agree_email_info"`
	// This field must be set if the user uses generated password to register.
	// Automatically set phone_verified if it's sent within a specific time.
	RegisterToken        string                 `protobuf:"bytes,8,opt,name=register_token,json=registerToken" json:"register_token"`
	Source               user_source.UserSource `protobuf:"varint,9,opt,name=source,enum=user_source.UserSource" json:"source"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *CreateUserRequest) Reset()         { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string { return proto.CompactTextString(m) }
func (*CreateUserRequest) ProtoMessage()    {}
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{15}
}

var xxx_messageInfo_CreateUserRequest proto.InternalMessageInfo

func (m *CreateUserRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateUserRequest) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

func (m *CreateUserRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CreateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *CreateUserRequest) GetAgreeTos() bool {
	if m != nil {
		return m.AgreeTos
	}
	return false
}

func (m *CreateUserRequest) GetAgreeEmailInfo() bool {
	if m != nil && m.AgreeEmailInfo != nil {
		return *m.AgreeEmailInfo
	}
	return false
}

func (m *CreateUserRequest) GetRegisterToken() string {
	if m != nil {
		return m.RegisterToken
	}
	return ""
}

func (m *CreateUserRequest) GetSource() user_source.UserSource {
	if m != nil {
		return m.Source
	}
	return user_source.UserSource_unknown
}

type RegisterResponse struct {
	// @required
	User                 *User    `protobuf:"bytes,1,opt,name=user" json:"user"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{16}
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

// Exactly one of phone or email must be provided.
type ResetPasswordRequest struct {
	// Phone number to send reset password instruction.
	Phone string `protobuf:"bytes,1,opt,name=phone" json:"phone"`
	// Email address to send reset password instruction.
	Email                string   `protobuf:"bytes,2,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResetPasswordRequest) Reset()         { *m = ResetPasswordRequest{} }
func (m *ResetPasswordRequest) String() string { return proto.CompactTextString(m) }
func (*ResetPasswordRequest) ProtoMessage()    {}
func (*ResetPasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{17}
}

var xxx_messageInfo_ResetPasswordRequest proto.InternalMessageInfo

func (m *ResetPasswordRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ResetPasswordRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

// Exactly one of current_password or reset_password_token must be provided.
type ChangePasswordRequest struct {
	// @required
	CurrentPassword string `protobuf:"bytes,2,opt,name=current_password,json=currentPassword" json:"current_password"`
	// @required
	NewPassword string `protobuf:"bytes,3,opt,name=new_password,json=newPassword" json:"new_password"`
	// @required
	ConfirmPassword      string   `protobuf:"bytes,4,opt,name=confirm_password,json=confirmPassword" json:"confirm_password"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordRequest) Reset()         { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordRequest) ProtoMessage()    {}
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{18}
}

var xxx_messageInfo_ChangePasswordRequest proto.InternalMessageInfo

func (m *ChangePasswordRequest) GetCurrentPassword() string {
	if m != nil {
		return m.CurrentPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetConfirmPassword() string {
	if m != nil {
		return m.ConfirmPassword
	}
	return ""
}

// Exactly one of email or phone must be provided.
type ChangePasswordUsingTokenRequest struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email"`
	Phone string `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	// @required
	ResetPasswordToken string `protobuf:"bytes,3,opt,name=reset_password_token,json=resetPasswordToken" json:"reset_password_token"`
	// @required
	NewPassword string `protobuf:"bytes,4,opt,name=new_password,json=newPassword" json:"new_password"`
	// @required
	ConfirmPassword      string   `protobuf:"bytes,5,opt,name=confirm_password,json=confirmPassword" json:"confirm_password"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordUsingTokenRequest) Reset()         { *m = ChangePasswordUsingTokenRequest{} }
func (m *ChangePasswordUsingTokenRequest) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordUsingTokenRequest) ProtoMessage()    {}
func (*ChangePasswordUsingTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{19}
}

var xxx_messageInfo_ChangePasswordUsingTokenRequest proto.InternalMessageInfo

func (m *ChangePasswordUsingTokenRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ChangePasswordUsingTokenRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ChangePasswordUsingTokenRequest) GetResetPasswordToken() string {
	if m != nil {
		return m.ResetPasswordToken
	}
	return ""
}

func (m *ChangePasswordUsingTokenRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

func (m *ChangePasswordUsingTokenRequest) GetConfirmPassword() string {
	if m != nil {
		return m.ConfirmPassword
	}
	return ""
}

// Represents permission of the current user relation with an account.
type Permission struct {
	Roles                []string `protobuf:"bytes,2,rep,name=roles" json:"roles"`
	Permissions          []string `protobuf:"bytes,3,rep,name=permissions" json:"permissions"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Permission) Reset()         { *m = Permission{} }
func (m *Permission) String() string { return proto.CompactTextString(m) }
func (*Permission) ProtoMessage()    {}
func (*Permission) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{20}
}

var xxx_messageInfo_Permission proto.InternalMessageInfo

func (m *Permission) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *Permission) GetPermissions() []string {
	if m != nil {
		return m.Permissions
	}
	return nil
}

type LoginRequest struct {
	// @required Phone or email
	Login string `protobuf:"bytes,1,opt,name=login" json:"login"`
	// @required
	Password string `protobuf:"bytes,2,opt,name=password" json:"password"`
	// Automatically switch to this account if available.
	//
	// It's *ignored* if the user *does not* have permission to this account.
	AccountId dot.ID `protobuf:"varint,3,opt,name=account_id,json=accountId" json:"account_id"`
	// Automatically switch to the only account of this account type if available.
	//
	// It's *ignored* if the user *does not* have any account of this account type, or if the user *has more than one* account of this account type.
	AccountType AccountType `protobuf:"varint,4,opt,name=account_type,json=accountType,enum=etop.AccountType" json:"account_type"`
	// Not implemented.
	AccountKey           string   `protobuf:"bytes,5,opt,name=account_key,json=accountKey" json:"account_key"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{21}
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *LoginRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *LoginRequest) GetAccountType() AccountType {
	if m != nil {
		return m.AccountType
	}
	return AccountType_unknown
}

func (m *LoginRequest) GetAccountKey() string {
	if m != nil {
		return m.AccountKey
	}
	return ""
}

// Represents an account associated with the current user. It has extra fields to represents relation with the user.
type LoginAccount struct {
	ExportedFields []string `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	// @required
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	Name string `protobuf:"bytes,2,opt,name=name" json:"name"`
	// @required
	Type AccountType `protobuf:"varint,3,opt,name=type,enum=etop.AccountType" json:"type"`
	// Associated token for the account. It's returned when calling Login or
	// SwitchAccount with regenerate_tokens set to true.
	AccessToken string `protobuf:"bytes,5,opt,name=access_token,json=accessToken" json:"access_token"`
	// The same as access_token.
	ExpiresIn            int32            `protobuf:"varint,6,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	ImageUrl             string           `protobuf:"bytes,7,opt,name=image_url,json=imageUrl" json:"image_url"`
	UrlSlug              string           `protobuf:"bytes,8,opt,name=url_slug,json=urlSlug" json:"url_slug"`
	UserAccount          *UserAccountInfo `protobuf:"bytes,9,opt,name=user_account,json=userAccount" json:"user_account"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *LoginAccount) Reset()         { *m = LoginAccount{} }
func (m *LoginAccount) String() string { return proto.CompactTextString(m) }
func (*LoginAccount) ProtoMessage()    {}
func (*LoginAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{22}
}

var xxx_messageInfo_LoginAccount proto.InternalMessageInfo

func (m *LoginAccount) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *LoginAccount) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *LoginAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LoginAccount) GetType() AccountType {
	if m != nil {
		return m.Type
	}
	return AccountType_unknown
}

func (m *LoginAccount) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *LoginAccount) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func (m *LoginAccount) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *LoginAccount) GetUrlSlug() string {
	if m != nil {
		return m.UrlSlug
	}
	return ""
}

func (m *LoginAccount) GetUserAccount() *UserAccountInfo {
	if m != nil {
		return m.UserAccount
	}
	return nil
}

type LoginResponse struct {
	// @required
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token"`
	// @required
	ExpiresIn int32 `protobuf:"varint,2,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	// @required
	User      *User         `protobuf:"bytes,3,opt,name=user" json:"user"`
	Account   *LoginAccount `protobuf:"bytes,4,opt,name=account" json:"account"`
	Shop      *Shop         `protobuf:"bytes,5,opt,name=shop" json:"shop"`
	Affiliate *Affiliate    `protobuf:"bytes,6,opt,name=affiliate" json:"affiliate"`
	// @required
	AvailableAccounts    []*LoginAccount    `protobuf:"bytes,7,rep,name=available_accounts,json=availableAccounts" json:"available_accounts"`
	InvitationAccounts   []*UserAccountInfo `protobuf:"bytes,8,rep,name=invitation_accounts,json=invitationAccounts" json:"invitation_accounts"`
	Stoken               bool               `protobuf:"varint,9,opt,name=stoken" json:"stoken"`
	StokenExpiresAt      dot.Time           `protobuf:"bytes,10,opt,name=stoken_expires_at,json=stokenExpiresAt" json:"stoken_expires_at"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{23}
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

func (m *LoginResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *LoginResponse) GetAccount() *LoginAccount {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *LoginResponse) GetShop() *Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *LoginResponse) GetAffiliate() *Affiliate {
	if m != nil {
		return m.Affiliate
	}
	return nil
}

func (m *LoginResponse) GetAvailableAccounts() []*LoginAccount {
	if m != nil {
		return m.AvailableAccounts
	}
	return nil
}

func (m *LoginResponse) GetInvitationAccounts() []*UserAccountInfo {
	if m != nil {
		return m.InvitationAccounts
	}
	return nil
}

func (m *LoginResponse) GetStoken() bool {
	if m != nil {
		return m.Stoken
	}
	return false
}

type SwitchAccountRequest struct {
	// @required
	AccountId dot.ID `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	// This field should only be used after creating new accounts. If it is set,
	// account_id can be left empty.
	RegenerateTokens     bool     `protobuf:"varint,2,opt,name=regenerate_tokens,json=regenerateTokens" json:"regenerate_tokens"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchAccountRequest) Reset()         { *m = SwitchAccountRequest{} }
func (m *SwitchAccountRequest) String() string { return proto.CompactTextString(m) }
func (*SwitchAccountRequest) ProtoMessage()    {}
func (*SwitchAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{24}
}

var xxx_messageInfo_SwitchAccountRequest proto.InternalMessageInfo

func (m *SwitchAccountRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *SwitchAccountRequest) GetRegenerateTokens() bool {
	if m != nil {
		return m.RegenerateTokens
	}
	return false
}

type UpgradeAccessTokenRequest struct {
	// @required
	Stoken               string   `protobuf:"bytes,1,opt,name=stoken" json:"stoken"`
	Password             string   `protobuf:"bytes,2,opt,name=password" json:"password"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpgradeAccessTokenRequest) Reset()         { *m = UpgradeAccessTokenRequest{} }
func (m *UpgradeAccessTokenRequest) String() string { return proto.CompactTextString(m) }
func (*UpgradeAccessTokenRequest) ProtoMessage()    {}
func (*UpgradeAccessTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{25}
}

var xxx_messageInfo_UpgradeAccessTokenRequest proto.InternalMessageInfo

func (m *UpgradeAccessTokenRequest) GetStoken() string {
	if m != nil {
		return m.Stoken
	}
	return ""
}

func (m *UpgradeAccessTokenRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type AccessTokenResponse struct {
	// @required
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token"`
	// @required
	ExpiresIn int32 `protobuf:"varint,2,opt,name=expires_in,json=expiresIn" json:"expires_in"`
	// @required
	User                 *User         `protobuf:"bytes,3,opt,name=user" json:"user"`
	Account              *LoginAccount `protobuf:"bytes,4,opt,name=account" json:"account"`
	Shop                 *Shop         `protobuf:"bytes,5,opt,name=shop" json:"shop"`
	Affiliate            *Affiliate    `protobuf:"bytes,6,opt,name=affiliate" json:"affiliate"`
	Stoken               bool          `protobuf:"varint,8,opt,name=stoken" json:"stoken"`
	StokenExpiresAt      dot.Time      `protobuf:"bytes,9,opt,name=stoken_expires_at,json=stokenExpiresAt" json:"stoken_expires_at"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AccessTokenResponse) Reset()         { *m = AccessTokenResponse{} }
func (m *AccessTokenResponse) String() string { return proto.CompactTextString(m) }
func (*AccessTokenResponse) ProtoMessage()    {}
func (*AccessTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{26}
}

var xxx_messageInfo_AccessTokenResponse proto.InternalMessageInfo

func (m *AccessTokenResponse) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *AccessTokenResponse) GetExpiresIn() int32 {
	if m != nil {
		return m.ExpiresIn
	}
	return 0
}

func (m *AccessTokenResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *AccessTokenResponse) GetAccount() *LoginAccount {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *AccessTokenResponse) GetShop() *Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *AccessTokenResponse) GetAffiliate() *Affiliate {
	if m != nil {
		return m.Affiliate
	}
	return nil
}

func (m *AccessTokenResponse) GetStoken() bool {
	if m != nil {
		return m.Stoken
	}
	return false
}

type SendSTokenEmailRequest struct {
	// @required
	Email string `protobuf:"bytes,1,opt,name=email" json:"email"`
	// @required
	AccountId            dot.ID   `protobuf:"varint,2,opt,name=account_id,json=accountId" json:"account_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendSTokenEmailRequest) Reset()         { *m = SendSTokenEmailRequest{} }
func (m *SendSTokenEmailRequest) String() string { return proto.CompactTextString(m) }
func (*SendSTokenEmailRequest) ProtoMessage()    {}
func (*SendSTokenEmailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{27}
}

var xxx_messageInfo_SendSTokenEmailRequest proto.InternalMessageInfo

func (m *SendSTokenEmailRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *SendSTokenEmailRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

type SendEmailVerificationRequest struct {
	// @required
	Email                string   `protobuf:"bytes,1,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendEmailVerificationRequest) Reset()         { *m = SendEmailVerificationRequest{} }
func (m *SendEmailVerificationRequest) String() string { return proto.CompactTextString(m) }
func (*SendEmailVerificationRequest) ProtoMessage()    {}
func (*SendEmailVerificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{28}
}

var xxx_messageInfo_SendEmailVerificationRequest proto.InternalMessageInfo

func (m *SendEmailVerificationRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type SendPhoneVerificationRequest struct {
	// @required
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendPhoneVerificationRequest) Reset()         { *m = SendPhoneVerificationRequest{} }
func (m *SendPhoneVerificationRequest) String() string { return proto.CompactTextString(m) }
func (*SendPhoneVerificationRequest) ProtoMessage()    {}
func (*SendPhoneVerificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{29}
}

var xxx_messageInfo_SendPhoneVerificationRequest proto.InternalMessageInfo

func (m *SendPhoneVerificationRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type VerifyEmailUsingTokenRequest struct {
	// @required
	VerificationToken    string   `protobuf:"bytes,1,opt,name=verification_token,json=verificationToken" json:"verification_token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyEmailUsingTokenRequest) Reset()         { *m = VerifyEmailUsingTokenRequest{} }
func (m *VerifyEmailUsingTokenRequest) String() string { return proto.CompactTextString(m) }
func (*VerifyEmailUsingTokenRequest) ProtoMessage()    {}
func (*VerifyEmailUsingTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{30}
}

var xxx_messageInfo_VerifyEmailUsingTokenRequest proto.InternalMessageInfo

func (m *VerifyEmailUsingTokenRequest) GetVerificationToken() string {
	if m != nil {
		return m.VerificationToken
	}
	return ""
}

type VerifyPhoneUsingTokenRequest struct {
	// @required
	VerificationToken    string   `protobuf:"bytes,1,opt,name=verification_token,json=verificationToken" json:"verification_token"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyPhoneUsingTokenRequest) Reset()         { *m = VerifyPhoneUsingTokenRequest{} }
func (m *VerifyPhoneUsingTokenRequest) String() string { return proto.CompactTextString(m) }
func (*VerifyPhoneUsingTokenRequest) ProtoMessage()    {}
func (*VerifyPhoneUsingTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{31}
}

var xxx_messageInfo_VerifyPhoneUsingTokenRequest proto.InternalMessageInfo

func (m *VerifyPhoneUsingTokenRequest) GetVerificationToken() string {
	if m != nil {
		return m.VerificationToken
	}
	return ""
}

type InviteUserToAccountRequest struct {
	// @required account to manage, must be owner of the account
	AccountId dot.ID `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	// @required phone or email
	InviteeIdentifier    string      `protobuf:"bytes,2,opt,name=invitee_identifier,json=inviteeIdentifier" json:"invitee_identifier"`
	Permission           *Permission `protobuf:"bytes,4,opt,name=permission" json:"permission"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *InviteUserToAccountRequest) Reset()         { *m = InviteUserToAccountRequest{} }
func (m *InviteUserToAccountRequest) String() string { return proto.CompactTextString(m) }
func (*InviteUserToAccountRequest) ProtoMessage()    {}
func (*InviteUserToAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{32}
}

var xxx_messageInfo_InviteUserToAccountRequest proto.InternalMessageInfo

func (m *InviteUserToAccountRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *InviteUserToAccountRequest) GetInviteeIdentifier() string {
	if m != nil {
		return m.InviteeIdentifier
	}
	return ""
}

func (m *InviteUserToAccountRequest) GetPermission() *Permission {
	if m != nil {
		return m.Permission
	}
	return nil
}

type AnswerInvitationRequest struct {
	AccountId            dot.ID          `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	Response             *status3.Status `protobuf:"varint,2,opt,name=response,enum=status3.Status" json:"response"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AnswerInvitationRequest) Reset()         { *m = AnswerInvitationRequest{} }
func (m *AnswerInvitationRequest) String() string { return proto.CompactTextString(m) }
func (*AnswerInvitationRequest) ProtoMessage()    {}
func (*AnswerInvitationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{33}
}

var xxx_messageInfo_AnswerInvitationRequest proto.InternalMessageInfo

func (m *AnswerInvitationRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *AnswerInvitationRequest) GetResponse() status3.Status {
	if m != nil && m.Response != nil {
		return *m.Response
	}
	return status3.Status_Z
}

type GetUsersInCurrentAccountsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	Mixed                *MixedAccount    `protobuf:"bytes,3,opt,name=mixed" json:"mixed"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetUsersInCurrentAccountsRequest) Reset()         { *m = GetUsersInCurrentAccountsRequest{} }
func (m *GetUsersInCurrentAccountsRequest) String() string { return proto.CompactTextString(m) }
func (*GetUsersInCurrentAccountsRequest) ProtoMessage()    {}
func (*GetUsersInCurrentAccountsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{34}
}

var xxx_messageInfo_GetUsersInCurrentAccountsRequest proto.InternalMessageInfo

func (m *GetUsersInCurrentAccountsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetUsersInCurrentAccountsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *GetUsersInCurrentAccountsRequest) GetMixed() *MixedAccount {
	if m != nil {
		return m.Mixed
	}
	return nil
}

type PublicUserInfo struct {
	// @required
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	FullName string `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	// @required
	ShortName            string   `protobuf:"bytes,3,opt,name=short_name,json=shortName" json:"short_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicUserInfo) Reset()         { *m = PublicUserInfo{} }
func (m *PublicUserInfo) String() string { return proto.CompactTextString(m) }
func (*PublicUserInfo) ProtoMessage()    {}
func (*PublicUserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{35}
}

var xxx_messageInfo_PublicUserInfo proto.InternalMessageInfo

func (m *PublicUserInfo) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PublicUserInfo) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *PublicUserInfo) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

// Presents user information inside an account
type UserAccountInfo struct {
	// @required
	UserId dot.ID `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id"`
	// @required
	UserFullName string `protobuf:"bytes,2,opt,name=user_full_name,json=userFullName" json:"user_full_name"`
	// @required
	UserShortName string `protobuf:"bytes,3,opt,name=user_short_name,json=userShortName" json:"user_short_name"`
	// @required
	AccountId dot.ID `protobuf:"varint,4,opt,name=account_id,json=accountId" json:"account_id"`
	// @required
	AccountName string `protobuf:"bytes,5,opt,name=account_name,json=accountName" json:"account_name"`
	// @required
	AccountType AccountType `protobuf:"varint,6,opt,name=account_type,json=accountType,enum=etop.AccountType" json:"account_type"`
	Position    string      `protobuf:"bytes,7,opt,name=position" json:"position"`
	// @required
	Permission           *Permission    `protobuf:"bytes,8,opt,name=permission" json:"permission"`
	Status               status3.Status `protobuf:"varint,9,opt,name=status,enum=status3.Status" json:"status"`
	ResponseStatus       status3.Status `protobuf:"varint,10,opt,name=response_status,json=responseStatus,enum=status3.Status" json:"response_status"`
	InvitationSentBy     dot.ID         `protobuf:"varint,11,opt,name=invitation_sent_by,json=invitationSentBy" json:"invitation_sent_by"`
	InvitationSentAt     dot.Time       `protobuf:"bytes,12,opt,name=invitation_sent_at,json=invitationSentAt" json:"invitation_sent_at"`
	InvitationAcceptedAt dot.Time       `protobuf:"bytes,13,opt,name=invitation_accepted_at,json=invitationAcceptedAt" json:"invitation_accepted_at"`
	InvitationRejectedAt dot.Time       `protobuf:"bytes,14,opt,name=invitation_rejected_at,json=invitationRejectedAt" json:"invitation_rejected_at"`
	DisabledAt           dot.Time       `protobuf:"bytes,15,opt,name=disabled_at,json=disabledAt" json:"disabled_at"`
	DisabledBy           dot.ID         `protobuf:"varint,16,opt,name=disabled_by,json=disabledBy" json:"disabled_by"`
	DisableReason        string         `protobuf:"bytes,17,opt,name=disable_reason,json=disableReason" json:"disable_reason"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UserAccountInfo) Reset()         { *m = UserAccountInfo{} }
func (m *UserAccountInfo) String() string { return proto.CompactTextString(m) }
func (*UserAccountInfo) ProtoMessage()    {}
func (*UserAccountInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{36}
}

var xxx_messageInfo_UserAccountInfo proto.InternalMessageInfo

func (m *UserAccountInfo) GetUserId() dot.ID {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *UserAccountInfo) GetUserFullName() string {
	if m != nil {
		return m.UserFullName
	}
	return ""
}

func (m *UserAccountInfo) GetUserShortName() string {
	if m != nil {
		return m.UserShortName
	}
	return ""
}

func (m *UserAccountInfo) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *UserAccountInfo) GetAccountName() string {
	if m != nil {
		return m.AccountName
	}
	return ""
}

func (m *UserAccountInfo) GetAccountType() AccountType {
	if m != nil {
		return m.AccountType
	}
	return AccountType_unknown
}

func (m *UserAccountInfo) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *UserAccountInfo) GetPermission() *Permission {
	if m != nil {
		return m.Permission
	}
	return nil
}

func (m *UserAccountInfo) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *UserAccountInfo) GetResponseStatus() status3.Status {
	if m != nil {
		return m.ResponseStatus
	}
	return status3.Status_Z
}

func (m *UserAccountInfo) GetInvitationSentBy() dot.ID {
	if m != nil {
		return m.InvitationSentBy
	}
	return 0
}

func (m *UserAccountInfo) GetDisabledBy() dot.ID {
	if m != nil {
		return m.DisabledBy
	}
	return 0
}

func (m *UserAccountInfo) GetDisableReason() string {
	if m != nil {
		return m.DisableReason
	}
	return ""
}

// Prepsents users in current account
type ProtectedUsersResponse struct {
	Paging               *common.PageInfo   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Users                []*UserAccountInfo `protobuf:"bytes,2,rep,name=users" json:"users"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ProtectedUsersResponse) Reset()         { *m = ProtectedUsersResponse{} }
func (m *ProtectedUsersResponse) String() string { return proto.CompactTextString(m) }
func (*ProtectedUsersResponse) ProtoMessage()    {}
func (*ProtectedUsersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{37}
}

var xxx_messageInfo_ProtectedUsersResponse proto.InternalMessageInfo

func (m *ProtectedUsersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ProtectedUsersResponse) GetUsers() []*UserAccountInfo {
	if m != nil {
		return m.Users
	}
	return nil
}

type LeaveAccountRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaveAccountRequest) Reset()         { *m = LeaveAccountRequest{} }
func (m *LeaveAccountRequest) String() string { return proto.CompactTextString(m) }
func (*LeaveAccountRequest) ProtoMessage()    {}
func (*LeaveAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{38}
}

var xxx_messageInfo_LeaveAccountRequest proto.InternalMessageInfo

type RemoveUserFromCurrentAccountRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveUserFromCurrentAccountRequest) Reset()         { *m = RemoveUserFromCurrentAccountRequest{} }
func (m *RemoveUserFromCurrentAccountRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveUserFromCurrentAccountRequest) ProtoMessage()    {}
func (*RemoveUserFromCurrentAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{39}
}

var xxx_messageInfo_RemoveUserFromCurrentAccountRequest proto.InternalMessageInfo

type Province struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Region               string   `protobuf:"bytes,3,opt,name=region" json:"region"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Province) Reset()         { *m = Province{} }
func (m *Province) String() string { return proto.CompactTextString(m) }
func (*Province) ProtoMessage()    {}
func (*Province) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{40}
}

var xxx_messageInfo_Province proto.InternalMessageInfo

func (m *Province) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Province) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Province) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

type GetProvincesResponse struct {
	Provinces            []*Province `protobuf:"bytes,1,rep,name=provinces" json:"provinces"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetProvincesResponse) Reset()         { *m = GetProvincesResponse{} }
func (m *GetProvincesResponse) String() string { return proto.CompactTextString(m) }
func (*GetProvincesResponse) ProtoMessage()    {}
func (*GetProvincesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{41}
}

var xxx_messageInfo_GetProvincesResponse proto.InternalMessageInfo

func (m *GetProvincesResponse) GetProvinces() []*Province {
	if m != nil {
		return m.Provinces
	}
	return nil
}

type District struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	ProvinceCode         string   `protobuf:"bytes,2,opt,name=province_code,json=provinceCode" json:"province_code"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	IsFreeship           bool     `protobuf:"varint,4,opt,name=is_freeship,json=isFreeship" json:"is_freeship"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *District) Reset()         { *m = District{} }
func (m *District) String() string { return proto.CompactTextString(m) }
func (*District) ProtoMessage()    {}
func (*District) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{42}
}

var xxx_messageInfo_District proto.InternalMessageInfo

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

func (m *District) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *District) GetIsFreeship() bool {
	if m != nil {
		return m.IsFreeship
	}
	return false
}

type GetDistrictsResponse struct {
	Districts            []*District `protobuf:"bytes,1,rep,name=districts" json:"districts"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetDistrictsResponse) Reset()         { *m = GetDistrictsResponse{} }
func (m *GetDistrictsResponse) String() string { return proto.CompactTextString(m) }
func (*GetDistrictsResponse) ProtoMessage()    {}
func (*GetDistrictsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{43}
}

var xxx_messageInfo_GetDistrictsResponse proto.InternalMessageInfo

func (m *GetDistrictsResponse) GetDistricts() []*District {
	if m != nil {
		return m.Districts
	}
	return nil
}

type GetDistrictsByProvinceRequest struct {
	ProvinceCode         string   `protobuf:"bytes,1,opt,name=province_code,json=provinceCode" json:"province_code"`
	ProvinceName         string   `protobuf:"bytes,2,opt,name=province_name,json=provinceName" json:"province_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDistrictsByProvinceRequest) Reset()         { *m = GetDistrictsByProvinceRequest{} }
func (m *GetDistrictsByProvinceRequest) String() string { return proto.CompactTextString(m) }
func (*GetDistrictsByProvinceRequest) ProtoMessage()    {}
func (*GetDistrictsByProvinceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{44}
}

var xxx_messageInfo_GetDistrictsByProvinceRequest proto.InternalMessageInfo

func (m *GetDistrictsByProvinceRequest) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *GetDistrictsByProvinceRequest) GetProvinceName() string {
	if m != nil {
		return m.ProvinceName
	}
	return ""
}

type Ward struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	DistrictCode         string   `protobuf:"bytes,2,opt,name=district_code,json=districtCode" json:"district_code"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ward) Reset()         { *m = Ward{} }
func (m *Ward) String() string { return proto.CompactTextString(m) }
func (*Ward) ProtoMessage()    {}
func (*Ward) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{45}
}

var xxx_messageInfo_Ward proto.InternalMessageInfo

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

func (m *Ward) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetWardsResponse struct {
	Wards                []*Ward  `protobuf:"bytes,1,rep,name=wards" json:"wards"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWardsResponse) Reset()         { *m = GetWardsResponse{} }
func (m *GetWardsResponse) String() string { return proto.CompactTextString(m) }
func (*GetWardsResponse) ProtoMessage()    {}
func (*GetWardsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{46}
}

var xxx_messageInfo_GetWardsResponse proto.InternalMessageInfo

func (m *GetWardsResponse) GetWards() []*Ward {
	if m != nil {
		return m.Wards
	}
	return nil
}

type GetWardsByDistrictRequest struct {
	DistrictCode         string   `protobuf:"bytes,1,opt,name=district_code,json=districtCode" json:"district_code"`
	DistrictName         string   `protobuf:"bytes,2,opt,name=district_name,json=districtName" json:"district_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWardsByDistrictRequest) Reset()         { *m = GetWardsByDistrictRequest{} }
func (m *GetWardsByDistrictRequest) String() string { return proto.CompactTextString(m) }
func (*GetWardsByDistrictRequest) ProtoMessage()    {}
func (*GetWardsByDistrictRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{47}
}

var xxx_messageInfo_GetWardsByDistrictRequest proto.InternalMessageInfo

func (m *GetWardsByDistrictRequest) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *GetWardsByDistrictRequest) GetDistrictName() string {
	if m != nil {
		return m.DistrictName
	}
	return ""
}

type ParseLocationRequest struct {
	ProvinceName         string   `protobuf:"bytes,1,opt,name=province_name,json=provinceName" json:"province_name"`
	DistrictName         string   `protobuf:"bytes,2,opt,name=district_name,json=districtName" json:"district_name"`
	WardName             string   `protobuf:"bytes,3,opt,name=ward_name,json=wardName" json:"ward_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ParseLocationRequest) Reset()         { *m = ParseLocationRequest{} }
func (m *ParseLocationRequest) String() string { return proto.CompactTextString(m) }
func (*ParseLocationRequest) ProtoMessage()    {}
func (*ParseLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{48}
}

var xxx_messageInfo_ParseLocationRequest proto.InternalMessageInfo

func (m *ParseLocationRequest) GetProvinceName() string {
	if m != nil {
		return m.ProvinceName
	}
	return ""
}

func (m *ParseLocationRequest) GetDistrictName() string {
	if m != nil {
		return m.DistrictName
	}
	return ""
}

func (m *ParseLocationRequest) GetWardName() string {
	if m != nil {
		return m.WardName
	}
	return ""
}

type ParseLocationResponse struct {
	Province             *Province `protobuf:"bytes,1,opt,name=province" json:"province"`
	District             *District `protobuf:"bytes,2,opt,name=district" json:"district"`
	Ward                 *Ward     `protobuf:"bytes,3,opt,name=ward" json:"ward"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ParseLocationResponse) Reset()         { *m = ParseLocationResponse{} }
func (m *ParseLocationResponse) String() string { return proto.CompactTextString(m) }
func (*ParseLocationResponse) ProtoMessage()    {}
func (*ParseLocationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{49}
}

var xxx_messageInfo_ParseLocationResponse proto.InternalMessageInfo

func (m *ParseLocationResponse) GetProvince() *Province {
	if m != nil {
		return m.Province
	}
	return nil
}

func (m *ParseLocationResponse) GetDistrict() *District {
	if m != nil {
		return m.District
	}
	return nil
}

func (m *ParseLocationResponse) GetWard() *Ward {
	if m != nil {
		return m.Ward
	}
	return nil
}

type Bank struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Type                 string   `protobuf:"bytes,3,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bank) Reset()         { *m = Bank{} }
func (m *Bank) String() string { return proto.CompactTextString(m) }
func (*Bank) ProtoMessage()    {}
func (*Bank) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{50}
}

var xxx_messageInfo_Bank proto.InternalMessageInfo

func (m *Bank) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Bank) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Bank) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type GetBanksResponse struct {
	Banks                []*Bank  `protobuf:"bytes,1,rep,name=banks" json:"banks"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBanksResponse) Reset()         { *m = GetBanksResponse{} }
func (m *GetBanksResponse) String() string { return proto.CompactTextString(m) }
func (*GetBanksResponse) ProtoMessage()    {}
func (*GetBanksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{51}
}

var xxx_messageInfo_GetBanksResponse proto.InternalMessageInfo

func (m *GetBanksResponse) GetBanks() []*Bank {
	if m != nil {
		return m.Banks
	}
	return nil
}

type BankProvince struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	BankCode             string   `protobuf:"bytes,3,opt,name=bank_code,json=bankCode" json:"bank_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BankProvince) Reset()         { *m = BankProvince{} }
func (m *BankProvince) String() string { return proto.CompactTextString(m) }
func (*BankProvince) ProtoMessage()    {}
func (*BankProvince) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{52}
}

var xxx_messageInfo_BankProvince proto.InternalMessageInfo

func (m *BankProvince) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *BankProvince) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BankProvince) GetBankCode() string {
	if m != nil {
		return m.BankCode
	}
	return ""
}

type GetBankProvincesResponse struct {
	Provinces            []*BankProvince `protobuf:"bytes,1,rep,name=provinces" json:"provinces"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetBankProvincesResponse) Reset()         { *m = GetBankProvincesResponse{} }
func (m *GetBankProvincesResponse) String() string { return proto.CompactTextString(m) }
func (*GetBankProvincesResponse) ProtoMessage()    {}
func (*GetBankProvincesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{53}
}

var xxx_messageInfo_GetBankProvincesResponse proto.InternalMessageInfo

func (m *GetBankProvincesResponse) GetProvinces() []*BankProvince {
	if m != nil {
		return m.Provinces
	}
	return nil
}

type GetProvincesByBankResquest struct {
	BankCode             string   `protobuf:"bytes,1,opt,name=bank_code,json=bankCode" json:"bank_code"`
	BankName             string   `protobuf:"bytes,2,opt,name=bank_name,json=bankName" json:"bank_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProvincesByBankResquest) Reset()         { *m = GetProvincesByBankResquest{} }
func (m *GetProvincesByBankResquest) String() string { return proto.CompactTextString(m) }
func (*GetProvincesByBankResquest) ProtoMessage()    {}
func (*GetProvincesByBankResquest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{54}
}

var xxx_messageInfo_GetProvincesByBankResquest proto.InternalMessageInfo

func (m *GetProvincesByBankResquest) GetBankCode() string {
	if m != nil {
		return m.BankCode
	}
	return ""
}

func (m *GetProvincesByBankResquest) GetBankName() string {
	if m != nil {
		return m.BankName
	}
	return ""
}

type BankBranch struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	BankCode             string   `protobuf:"bytes,3,opt,name=bank_code,json=bankCode" json:"bank_code"`
	ProvinceCode         string   `protobuf:"bytes,4,opt,name=province_code,json=provinceCode" json:"province_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BankBranch) Reset()         { *m = BankBranch{} }
func (m *BankBranch) String() string { return proto.CompactTextString(m) }
func (*BankBranch) ProtoMessage()    {}
func (*BankBranch) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{55}
}

var xxx_messageInfo_BankBranch proto.InternalMessageInfo

func (m *BankBranch) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *BankBranch) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BankBranch) GetBankCode() string {
	if m != nil {
		return m.BankCode
	}
	return ""
}

func (m *BankBranch) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

type GetBranchesByBankProvinceResponse struct {
	Branches             []*BankBranch `protobuf:"bytes,1,rep,name=branches" json:"branches"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetBranchesByBankProvinceResponse) Reset()         { *m = GetBranchesByBankProvinceResponse{} }
func (m *GetBranchesByBankProvinceResponse) String() string { return proto.CompactTextString(m) }
func (*GetBranchesByBankProvinceResponse) ProtoMessage()    {}
func (*GetBranchesByBankProvinceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{56}
}

var xxx_messageInfo_GetBranchesByBankProvinceResponse proto.InternalMessageInfo

func (m *GetBranchesByBankProvinceResponse) GetBranches() []*BankBranch {
	if m != nil {
		return m.Branches
	}
	return nil
}

type GetBranchesByBankProvinceResquest struct {
	BankCode             string   `protobuf:"bytes,1,opt,name=bank_code,json=bankCode" json:"bank_code"`
	BankName             string   `protobuf:"bytes,2,opt,name=bank_name,json=bankName" json:"bank_name"`
	ProvinceCode         string   `protobuf:"bytes,3,opt,name=province_code,json=provinceCode" json:"province_code"`
	ProvinceName         string   `protobuf:"bytes,4,opt,name=province_name,json=provinceName" json:"province_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBranchesByBankProvinceResquest) Reset()         { *m = GetBranchesByBankProvinceResquest{} }
func (m *GetBranchesByBankProvinceResquest) String() string { return proto.CompactTextString(m) }
func (*GetBranchesByBankProvinceResquest) ProtoMessage()    {}
func (*GetBranchesByBankProvinceResquest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{57}
}

var xxx_messageInfo_GetBranchesByBankProvinceResquest proto.InternalMessageInfo

func (m *GetBranchesByBankProvinceResquest) GetBankCode() string {
	if m != nil {
		return m.BankCode
	}
	return ""
}

func (m *GetBranchesByBankProvinceResquest) GetBankName() string {
	if m != nil {
		return m.BankName
	}
	return ""
}

func (m *GetBranchesByBankProvinceResquest) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *GetBranchesByBankProvinceResquest) GetProvinceName() string {
	if m != nil {
		return m.ProvinceName
	}
	return ""
}

type Address struct {
	ExportedFields []string `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	Id             dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Province       string   `protobuf:"bytes,2,opt,name=province" json:"province"`
	ProvinceCode   string   `protobuf:"bytes,6,opt,name=province_code,json=provinceCode" json:"province_code"`
	District       string   `protobuf:"bytes,3,opt,name=district" json:"district"`
	DistrictCode   string   `protobuf:"bytes,7,opt,name=district_code,json=districtCode" json:"district_code"`
	Ward           string   `protobuf:"bytes,4,opt,name=ward" json:"ward"`
	WardCode       string   `protobuf:"bytes,8,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1       string   `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2       string   `protobuf:"bytes,9,opt,name=address2" json:"address2"`
	Zip            string   `protobuf:"bytes,10,opt,name=zip" json:"zip"`
	Country        string   `protobuf:"bytes,11,opt,name=country" json:"country"`
	FullName       string   `protobuf:"bytes,12,opt,name=full_name,json=fullName" json:"full_name"`
	// deprecated: use full_name instead
	FirstName string `protobuf:"bytes,13,opt,name=first_name,json=firstName" json:"first_name"`
	// deprecated: use full_name instead
	LastName             string                   `protobuf:"bytes,14,opt,name=last_name,json=lastName" json:"last_name"`
	Phone                string                   `protobuf:"bytes,15,opt,name=phone" json:"phone"`
	Email                string                   `protobuf:"bytes,16,opt,name=email" json:"email"`
	Position             string                   `protobuf:"bytes,17,opt,name=position" json:"position"`
	Type                 address_type.AddressType `protobuf:"varint,18,opt,name=type,enum=address_type.AddressType" json:"type"`
	Notes                *AddressNote             `protobuf:"bytes,19,opt,name=notes" json:"notes"`
	Coordinates          *Coordinates             `protobuf:"bytes,20,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{58}
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *Address) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Address) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *Address) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *Address) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *Address) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *Address) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *Address) GetWardCode() string {
	if m != nil {
		return m.WardCode
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

func (m *Address) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *Address) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Address) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Address) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *Address) GetLastName() string {
	if m != nil {
		return m.LastName
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

func (m *Address) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *Address) GetType() address_type.AddressType {
	if m != nil {
		return m.Type
	}
	return address_type.AddressType_unknown
}

func (m *Address) GetNotes() *AddressNote {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *Address) GetCoordinates() *Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type Coordinates struct {
	Latitude             float32  `protobuf:"fixed32,1,opt,name=latitude" json:"latitude"`
	Longitude            float32  `protobuf:"fixed32,2,opt,name=longitude" json:"longitude"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Coordinates) Reset()         { *m = Coordinates{} }
func (m *Coordinates) String() string { return proto.CompactTextString(m) }
func (*Coordinates) ProtoMessage()    {}
func (*Coordinates) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{59}
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

type BankAccount struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Province             string   `protobuf:"bytes,2,opt,name=province" json:"province"`
	Branch               string   `protobuf:"bytes,3,opt,name=branch" json:"branch"`
	AccountNumber        string   `protobuf:"bytes,4,opt,name=account_number,json=accountNumber" json:"account_number"`
	AccountName          string   `protobuf:"bytes,5,opt,name=account_name,json=accountName" json:"account_name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BankAccount) Reset()         { *m = BankAccount{} }
func (m *BankAccount) String() string { return proto.CompactTextString(m) }
func (*BankAccount) ProtoMessage()    {}
func (*BankAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{60}
}

var xxx_messageInfo_BankAccount proto.InternalMessageInfo

func (m *BankAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BankAccount) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *BankAccount) GetBranch() string {
	if m != nil {
		return m.Branch
	}
	return ""
}

func (m *BankAccount) GetAccountNumber() string {
	if m != nil {
		return m.AccountNumber
	}
	return ""
}

func (m *BankAccount) GetAccountName() string {
	if m != nil {
		return m.AccountName
	}
	return ""
}

type ContactPerson struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Position             string   `protobuf:"bytes,2,opt,name=position" json:"position"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContactPerson) Reset()         { *m = ContactPerson{} }
func (m *ContactPerson) String() string { return proto.CompactTextString(m) }
func (*ContactPerson) ProtoMessage()    {}
func (*ContactPerson) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{61}
}

var xxx_messageInfo_ContactPerson proto.InternalMessageInfo

func (m *ContactPerson) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ContactPerson) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *ContactPerson) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ContactPerson) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type CompanyInfo struct {
	Name                 string         `protobuf:"bytes,1,opt,name=name" json:"name"`
	TaxCode              string         `protobuf:"bytes,2,opt,name=tax_code,json=taxCode" json:"tax_code"`
	Address              string         `protobuf:"bytes,3,opt,name=address" json:"address"`
	Website              string         `protobuf:"bytes,5,opt,name=website" json:"website"`
	LegalRepresentative  *ContactPerson `protobuf:"bytes,4,opt,name=legal_representative,json=legalRepresentative" json:"legal_representative"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CompanyInfo) Reset()         { *m = CompanyInfo{} }
func (m *CompanyInfo) String() string { return proto.CompactTextString(m) }
func (*CompanyInfo) ProtoMessage()    {}
func (*CompanyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{62}
}

var xxx_messageInfo_CompanyInfo proto.InternalMessageInfo

func (m *CompanyInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CompanyInfo) GetTaxCode() string {
	if m != nil {
		return m.TaxCode
	}
	return ""
}

func (m *CompanyInfo) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *CompanyInfo) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *CompanyInfo) GetLegalRepresentative() *ContactPerson {
	if m != nil {
		return m.LegalRepresentative
	}
	return nil
}

type CreateAddressRequest struct {
	// @required
	Province     string `protobuf:"bytes,2,opt,name=province" json:"province"`
	ProvinceCode string `protobuf:"bytes,6,opt,name=province_code,json=provinceCode" json:"province_code"`
	// @required
	District     string `protobuf:"bytes,3,opt,name=district" json:"district"`
	DistrictCode string `protobuf:"bytes,7,opt,name=district_code,json=districtCode" json:"district_code"`
	// @required
	Ward                 string                   `protobuf:"bytes,4,opt,name=ward" json:"ward"`
	WardCode             string                   `protobuf:"bytes,8,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1             string                   `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2             string                   `protobuf:"bytes,9,opt,name=address2" json:"address2"`
	Zip                  string                   `protobuf:"bytes,10,opt,name=zip" json:"zip"`
	Country              string                   `protobuf:"bytes,11,opt,name=country" json:"country"`
	FullName             string                   `protobuf:"bytes,12,opt,name=full_name,json=fullName" json:"full_name"`
	FirstName            string                   `protobuf:"bytes,13,opt,name=first_name,json=firstName" json:"first_name"`
	LastName             string                   `protobuf:"bytes,14,opt,name=last_name,json=lastName" json:"last_name"`
	Phone                string                   `protobuf:"bytes,15,opt,name=phone" json:"phone"`
	Email                string                   `protobuf:"bytes,16,opt,name=email" json:"email"`
	Position             string                   `protobuf:"bytes,17,opt,name=position" json:"position"`
	Type                 address_type.AddressType `protobuf:"varint,18,opt,name=type,enum=address_type.AddressType" json:"type"`
	Notes                *AddressNote             `protobuf:"bytes,19,opt,name=notes" json:"notes"`
	Coordinates          *Coordinates             `protobuf:"bytes,20,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *CreateAddressRequest) Reset()         { *m = CreateAddressRequest{} }
func (m *CreateAddressRequest) String() string { return proto.CompactTextString(m) }
func (*CreateAddressRequest) ProtoMessage()    {}
func (*CreateAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{63}
}

var xxx_messageInfo_CreateAddressRequest proto.InternalMessageInfo

func (m *CreateAddressRequest) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *CreateAddressRequest) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *CreateAddressRequest) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *CreateAddressRequest) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *CreateAddressRequest) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *CreateAddressRequest) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *CreateAddressRequest) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *CreateAddressRequest) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *CreateAddressRequest) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *CreateAddressRequest) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *CreateAddressRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateAddressRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *CreateAddressRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *CreateAddressRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CreateAddressRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateAddressRequest) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *CreateAddressRequest) GetType() address_type.AddressType {
	if m != nil {
		return m.Type
	}
	return address_type.AddressType_unknown
}

func (m *CreateAddressRequest) GetNotes() *AddressNote {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *CreateAddressRequest) GetCoordinates() *Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type AddressNote struct {
	Note                 string   `protobuf:"bytes,1,opt,name=note" json:"note"`
	OpenTime             string   `protobuf:"bytes,2,opt,name=open_time,json=openTime" json:"open_time"`
	LunchBreak           string   `protobuf:"bytes,3,opt,name=lunch_break,json=lunchBreak" json:"lunch_break"`
	Other                string   `protobuf:"bytes,4,opt,name=other" json:"other"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddressNote) Reset()         { *m = AddressNote{} }
func (m *AddressNote) String() string { return proto.CompactTextString(m) }
func (*AddressNote) ProtoMessage()    {}
func (*AddressNote) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{64}
}

var xxx_messageInfo_AddressNote proto.InternalMessageInfo

func (m *AddressNote) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *AddressNote) GetOpenTime() string {
	if m != nil {
		return m.OpenTime
	}
	return ""
}

func (m *AddressNote) GetLunchBreak() string {
	if m != nil {
		return m.LunchBreak
	}
	return ""
}

func (m *AddressNote) GetOther() string {
	if m != nil {
		return m.Other
	}
	return ""
}

type GetAddressResponse struct {
	Addresses            []*Address `protobuf:"bytes,1,rep,name=addresses" json:"addresses"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GetAddressResponse) Reset()         { *m = GetAddressResponse{} }
func (m *GetAddressResponse) String() string { return proto.CompactTextString(m) }
func (*GetAddressResponse) ProtoMessage()    {}
func (*GetAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{65}
}

var xxx_messageInfo_GetAddressResponse proto.InternalMessageInfo

func (m *GetAddressResponse) GetAddresses() []*Address {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type UpdateAddressRequest struct {
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	Province     string `protobuf:"bytes,2,opt,name=province" json:"province"`
	ProvinceCode string `protobuf:"bytes,6,opt,name=province_code,json=provinceCode" json:"province_code"`
	// @required
	District     string `protobuf:"bytes,3,opt,name=district" json:"district"`
	DistrictCode string `protobuf:"bytes,7,opt,name=district_code,json=districtCode" json:"district_code"`
	// @required
	Ward                 string                   `protobuf:"bytes,4,opt,name=ward" json:"ward"`
	WardCode             string                   `protobuf:"bytes,8,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1             string                   `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2             string                   `protobuf:"bytes,9,opt,name=address2" json:"address2"`
	Zip                  string                   `protobuf:"bytes,10,opt,name=zip" json:"zip"`
	Country              string                   `protobuf:"bytes,11,opt,name=country" json:"country"`
	FullName             string                   `protobuf:"bytes,12,opt,name=full_name,json=fullName" json:"full_name"`
	FirstName            string                   `protobuf:"bytes,13,opt,name=first_name,json=firstName" json:"first_name"`
	LastName             string                   `protobuf:"bytes,14,opt,name=last_name,json=lastName" json:"last_name"`
	Phone                string                   `protobuf:"bytes,15,opt,name=phone" json:"phone"`
	Email                string                   `protobuf:"bytes,16,opt,name=email" json:"email"`
	Position             string                   `protobuf:"bytes,17,opt,name=position" json:"position"`
	Type                 address_type.AddressType `protobuf:"varint,18,opt,name=type,enum=address_type.AddressType" json:"type"`
	Notes                *AddressNote             `protobuf:"bytes,19,opt,name=notes" json:"notes"`
	Coordinates          *Coordinates             `protobuf:"bytes,20,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *UpdateAddressRequest) Reset()         { *m = UpdateAddressRequest{} }
func (m *UpdateAddressRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateAddressRequest) ProtoMessage()    {}
func (*UpdateAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{66}
}

var xxx_messageInfo_UpdateAddressRequest proto.InternalMessageInfo

func (m *UpdateAddressRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateAddressRequest) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *UpdateAddressRequest) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *UpdateAddressRequest) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *UpdateAddressRequest) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *UpdateAddressRequest) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *UpdateAddressRequest) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *UpdateAddressRequest) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *UpdateAddressRequest) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *UpdateAddressRequest) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *UpdateAddressRequest) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *UpdateAddressRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *UpdateAddressRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *UpdateAddressRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *UpdateAddressRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *UpdateAddressRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UpdateAddressRequest) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *UpdateAddressRequest) GetType() address_type.AddressType {
	if m != nil {
		return m.Type
	}
	return address_type.AddressType_unknown
}

func (m *UpdateAddressRequest) GetNotes() *AddressNote {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *UpdateAddressRequest) GetCoordinates() *Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type SetDefaultAddressRequest struct {
	Id                   dot.ID                   `protobuf:"varint,1,opt,name=id" json:"id"`
	Type                 address_type.AddressType `protobuf:"varint,2,opt,name=type,enum=address_type.AddressType" json:"type"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *SetDefaultAddressRequest) Reset()         { *m = SetDefaultAddressRequest{} }
func (m *SetDefaultAddressRequest) String() string { return proto.CompactTextString(m) }
func (*SetDefaultAddressRequest) ProtoMessage()    {}
func (*SetDefaultAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{67}
}

var xxx_messageInfo_SetDefaultAddressRequest proto.InternalMessageInfo

func (m *SetDefaultAddressRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SetDefaultAddressRequest) GetType() address_type.AddressType {
	if m != nil {
		return m.Type
	}
	return address_type.AddressType_unknown
}

type UpdateURLSlugRequest struct {
	AccountId            dot.ID   `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	UrlSlug              *string  `protobuf:"bytes,2,opt,name=url_slug,json=urlSlug" json:"url_slug"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateURLSlugRequest) Reset()         { *m = UpdateURLSlugRequest{} }
func (m *UpdateURLSlugRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateURLSlugRequest) ProtoMessage()    {}
func (*UpdateURLSlugRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{68}
}

var xxx_messageInfo_UpdateURLSlugRequest proto.InternalMessageInfo

func (m *UpdateURLSlugRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *UpdateURLSlugRequest) GetUrlSlug() string {
	if m != nil && m.UrlSlug != nil {
		return *m.UrlSlug
	}
	return ""
}

type HistoryResponse struct {
	Paging               *common.PageInfo      `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Data                 *common.RawJSONObject `protobuf:"bytes,2,opt,name=data" json:"data"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *HistoryResponse) Reset()         { *m = HistoryResponse{} }
func (m *HistoryResponse) String() string { return proto.CompactTextString(m) }
func (*HistoryResponse) ProtoMessage()    {}
func (*HistoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{69}
}

var xxx_messageInfo_HistoryResponse proto.InternalMessageInfo

func (m *HistoryResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *HistoryResponse) GetData() *common.RawJSONObject {
	if m != nil {
		return m.Data
	}
	return nil
}

type Credit struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Amount               int            `protobuf:"varint,2,opt,name=amount" json:"amount"`
	ShopId               dot.ID         `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	Type                 string         `protobuf:"bytes,5,opt,name=type" json:"type"`
	Shop                 *Shop          `protobuf:"bytes,6,opt,name=shop" json:"shop"`
	CreatedAt            dot.Time       `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time       `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	PaidAt               dot.Time       `protobuf:"bytes,9,opt,name=paid_at,json=paidAt" json:"paid_at"`
	Status               status3.Status `protobuf:"varint,10,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Credit) Reset()         { *m = Credit{} }
func (m *Credit) String() string { return proto.CompactTextString(m) }
func (*Credit) ProtoMessage()    {}
func (*Credit) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{70}
}

var xxx_messageInfo_Credit proto.InternalMessageInfo

func (m *Credit) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Credit) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Credit) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Credit) GetShop() *Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *Credit) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

type CreditsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Credits              []*Credit        `protobuf:"bytes,2,rep,name=credits" json:"credits"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CreditsResponse) Reset()         { *m = CreditsResponse{} }
func (m *CreditsResponse) String() string { return proto.CompactTextString(m) }
func (*CreditsResponse) ProtoMessage()    {}
func (*CreditsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{71}
}

var xxx_messageInfo_CreditsResponse proto.InternalMessageInfo

func (m *CreditsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *CreditsResponse) GetCredits() []*Credit {
	if m != nil {
		return m.Credits
	}
	return nil
}

type UpdatePermissionRequest struct {
	Items                []*UpdatePermissionItem `protobuf:"bytes,1,rep,name=items" json:"items"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *UpdatePermissionRequest) Reset()         { *m = UpdatePermissionRequest{} }
func (m *UpdatePermissionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdatePermissionRequest) ProtoMessage()    {}
func (*UpdatePermissionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{72}
}

var xxx_messageInfo_UpdatePermissionRequest proto.InternalMessageInfo

func (m *UpdatePermissionRequest) GetItems() []*UpdatePermissionItem {
	if m != nil {
		return m.Items
	}
	return nil
}

type UpdatePermissionItem struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type" json:"type"`
	Key                  string   `protobuf:"bytes,2,opt,name=key" json:"key"`
	Grants               []string `protobuf:"bytes,3,rep,name=grants" json:"grants"`
	Revokes              []string `protobuf:"bytes,4,rep,name=revokes" json:"revokes"`
	ReplaceAll           []string `protobuf:"bytes,5,rep,name=replace_all,json=replaceAll" json:"replace_all"`
	RevokeAll            bool     `protobuf:"varint,6,opt,name=revoke_all,json=revokeAll" json:"revoke_all"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePermissionItem) Reset()         { *m = UpdatePermissionItem{} }
func (m *UpdatePermissionItem) String() string { return proto.CompactTextString(m) }
func (*UpdatePermissionItem) ProtoMessage()    {}
func (*UpdatePermissionItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{73}
}

var xxx_messageInfo_UpdatePermissionItem proto.InternalMessageInfo

func (m *UpdatePermissionItem) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *UpdatePermissionItem) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *UpdatePermissionItem) GetGrants() []string {
	if m != nil {
		return m.Grants
	}
	return nil
}

func (m *UpdatePermissionItem) GetRevokes() []string {
	if m != nil {
		return m.Revokes
	}
	return nil
}

func (m *UpdatePermissionItem) GetReplaceAll() []string {
	if m != nil {
		return m.ReplaceAll
	}
	return nil
}

func (m *UpdatePermissionItem) GetRevokeAll() bool {
	if m != nil {
		return m.RevokeAll
	}
	return false
}

type UpdatePermissionResponse struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg" json:"msg"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePermissionResponse) Reset()         { *m = UpdatePermissionResponse{} }
func (m *UpdatePermissionResponse) String() string { return proto.CompactTextString(m) }
func (*UpdatePermissionResponse) ProtoMessage()    {}
func (*UpdatePermissionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{74}
}

var xxx_messageInfo_UpdatePermissionResponse proto.InternalMessageInfo

func (m *UpdatePermissionResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type Device struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	DeviceId             string   `protobuf:"bytes,2,opt,name=device_id,json=deviceId" json:"device_id"`
	DeviceName           string   `protobuf:"bytes,3,opt,name=device_name,json=deviceName" json:"device_name"`
	ExternalDeviceId     string   `protobuf:"bytes,4,opt,name=external_device_id,json=externalDeviceId" json:"external_device_id"`
	ExternalServiceId    int32    `protobuf:"varint,8,opt,name=external_service_id,json=externalServiceId" json:"external_service_id"`
	AccountId            dot.ID   `protobuf:"varint,5,opt,name=account_id,json=accountId" json:"account_id"`
	CreatedAt            dot.Time `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Device) Reset()         { *m = Device{} }
func (m *Device) String() string { return proto.CompactTextString(m) }
func (*Device) ProtoMessage()    {}
func (*Device) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{75}
}

var xxx_messageInfo_Device proto.InternalMessageInfo

func (m *Device) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Device) GetDeviceId() string {
	if m != nil {
		return m.DeviceId
	}
	return ""
}

func (m *Device) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *Device) GetExternalDeviceId() string {
	if m != nil {
		return m.ExternalDeviceId
	}
	return ""
}

func (m *Device) GetExternalServiceId() int32 {
	if m != nil {
		return m.ExternalServiceId
	}
	return 0
}

func (m *Device) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

type Notification struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Title                string         `protobuf:"bytes,2,opt,name=title" json:"title"`
	Message              string         `protobuf:"bytes,3,opt,name=message" json:"message"`
	IsRead               bool           `protobuf:"varint,4,opt,name=is_read,json=isRead" json:"is_read"`
	Entity               string         `protobuf:"bytes,5,opt,name=entity" json:"entity"`
	EntityId             dot.ID         `protobuf:"varint,6,opt,name=entity_id,json=entityId" json:"entity_id"`
	AccountId            dot.ID         `protobuf:"varint,7,opt,name=account_id,json=accountId" json:"account_id"`
	SendNotification     bool           `protobuf:"varint,12,opt,name=send_notification,json=sendNotification" json:"send_notification"`
	SeenAt               dot.Time       `protobuf:"bytes,8,opt,name=seen_at,json=seenAt" json:"seen_at"`
	CreatedAt            dot.Time       `protobuf:"bytes,9,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time       `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	SyncStatus           status3.Status `protobuf:"varint,11,opt,name=sync_status,json=syncStatus,enum=status3.Status" json:"sync_status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Notification) Reset()         { *m = Notification{} }
func (m *Notification) String() string { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()    {}
func (*Notification) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{76}
}

var xxx_messageInfo_Notification proto.InternalMessageInfo

func (m *Notification) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Notification) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Notification) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Notification) GetIsRead() bool {
	if m != nil {
		return m.IsRead
	}
	return false
}

func (m *Notification) GetEntity() string {
	if m != nil {
		return m.Entity
	}
	return ""
}

func (m *Notification) GetEntityId() dot.ID {
	if m != nil {
		return m.EntityId
	}
	return 0
}

func (m *Notification) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *Notification) GetSendNotification() bool {
	if m != nil {
		return m.SendNotification
	}
	return false
}

func (m *Notification) GetSyncStatus() status3.Status {
	if m != nil {
		return m.SyncStatus
	}
	return status3.Status_Z
}

type CreateDeviceRequest struct {
	DeviceId             string   `protobuf:"bytes,1,opt,name=device_id,json=deviceId" json:"device_id"`
	DeviceName           string   `protobuf:"bytes,2,opt,name=device_name,json=deviceName" json:"device_name"`
	ExternalDeviceId     string   `protobuf:"bytes,3,opt,name=external_device_id,json=externalDeviceId" json:"external_device_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateDeviceRequest) Reset()         { *m = CreateDeviceRequest{} }
func (m *CreateDeviceRequest) String() string { return proto.CompactTextString(m) }
func (*CreateDeviceRequest) ProtoMessage()    {}
func (*CreateDeviceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{77}
}

var xxx_messageInfo_CreateDeviceRequest proto.InternalMessageInfo

func (m *CreateDeviceRequest) GetDeviceId() string {
	if m != nil {
		return m.DeviceId
	}
	return ""
}

func (m *CreateDeviceRequest) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *CreateDeviceRequest) GetExternalDeviceId() string {
	if m != nil {
		return m.ExternalDeviceId
	}
	return ""
}

type DeleteDeviceRequest struct {
	DeviceId             string   `protobuf:"bytes,1,opt,name=device_id,json=deviceId" json:"device_id"`
	ExternalDeviceId     string   `protobuf:"bytes,2,opt,name=external_device_id,json=externalDeviceId" json:"external_device_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteDeviceRequest) Reset()         { *m = DeleteDeviceRequest{} }
func (m *DeleteDeviceRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteDeviceRequest) ProtoMessage()    {}
func (*DeleteDeviceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{78}
}

var xxx_messageInfo_DeleteDeviceRequest proto.InternalMessageInfo

func (m *DeleteDeviceRequest) GetDeviceId() string {
	if m != nil {
		return m.DeviceId
	}
	return ""
}

func (m *DeleteDeviceRequest) GetExternalDeviceId() string {
	if m != nil {
		return m.ExternalDeviceId
	}
	return ""
}

type GetNotificationsRequest struct {
	Paging               *common.Paging `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetNotificationsRequest) Reset()         { *m = GetNotificationsRequest{} }
func (m *GetNotificationsRequest) String() string { return proto.CompactTextString(m) }
func (*GetNotificationsRequest) ProtoMessage()    {}
func (*GetNotificationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{79}
}

var xxx_messageInfo_GetNotificationsRequest proto.InternalMessageInfo

func (m *GetNotificationsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

type NotificationsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Notifications        []*Notification  `protobuf:"bytes,2,rep,name=notifications" json:"notifications"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NotificationsResponse) Reset()         { *m = NotificationsResponse{} }
func (m *NotificationsResponse) String() string { return proto.CompactTextString(m) }
func (*NotificationsResponse) ProtoMessage()    {}
func (*NotificationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{80}
}

var xxx_messageInfo_NotificationsResponse proto.InternalMessageInfo

func (m *NotificationsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *NotificationsResponse) GetNotifications() []*Notification {
	if m != nil {
		return m.Notifications
	}
	return nil
}

type UpdateNotificationsRequest struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	IsRead               bool     `protobuf:"varint,2,opt,name=is_read,json=isRead" json:"is_read"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateNotificationsRequest) Reset()         { *m = UpdateNotificationsRequest{} }
func (m *UpdateNotificationsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateNotificationsRequest) ProtoMessage()    {}
func (*UpdateNotificationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{81}
}

var xxx_messageInfo_UpdateNotificationsRequest proto.InternalMessageInfo

func (m *UpdateNotificationsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *UpdateNotificationsRequest) GetIsRead() bool {
	if m != nil {
		return m.IsRead
	}
	return false
}

type UpdateReferenceUserRequest struct {
	// @required
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateReferenceUserRequest) Reset()         { *m = UpdateReferenceUserRequest{} }
func (m *UpdateReferenceUserRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateReferenceUserRequest) ProtoMessage()    {}
func (*UpdateReferenceUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{82}
}

var xxx_messageInfo_UpdateReferenceUserRequest proto.InternalMessageInfo

func (m *UpdateReferenceUserRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type UpdateReferenceSaleRequest struct {
	// @required
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateReferenceSaleRequest) Reset()         { *m = UpdateReferenceSaleRequest{} }
func (m *UpdateReferenceSaleRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateReferenceSaleRequest) ProtoMessage()    {}
func (*UpdateReferenceSaleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{83}
}

var xxx_messageInfo_UpdateReferenceSaleRequest proto.InternalMessageInfo

func (m *UpdateReferenceSaleRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type Affiliate struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string         `protobuf:"bytes,2,opt,name=name" json:"name"`
	Status               status3.Status `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	IsTest               bool           `protobuf:"varint,4,opt,name=is_test,json=isTest" json:"is_test"`
	Phone                string         `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Email                string         `protobuf:"bytes,6,opt,name=email" json:"email"`
	BankAccount          *BankAccount   `protobuf:"bytes,7,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Affiliate) Reset()         { *m = Affiliate{} }
func (m *Affiliate) String() string { return proto.CompactTextString(m) }
func (*Affiliate) ProtoMessage()    {}
func (*Affiliate) Descriptor() ([]byte, []int) {
	return fileDescriptor_4324bfd8c89695e1, []int{84}
}

var xxx_messageInfo_Affiliate proto.InternalMessageInfo

func (m *Affiliate) GetId() dot.ID {
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

func (m *Affiliate) GetBankAccount() *BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func init() {
	proto.RegisterEnum("etop.AccountType", AccountType_name, AccountType_value)
	proto.RegisterType((*GetInvitationsRequest)(nil), "etop.GetInvitationsRequest")
	proto.RegisterType((*GetInvitationByTokenRequest)(nil), "etop.GetInvitationByTokenRequest")
	proto.RegisterType((*Invitation)(nil), "etop.Invitation")
	proto.RegisterType((*AcceptInvitationRequest)(nil), "etop.AcceptInvitationRequest")
	proto.RegisterType((*RejectInvitationRequest)(nil), "etop.RejectInvitationRequest")
	proto.RegisterType((*User)(nil), "etop.User")
	proto.RegisterType((*IDsRequest)(nil), "etop.IDsRequest")
	proto.RegisterType((*MixedAccount)(nil), "etop.MixedAccount")
	proto.RegisterType((*Partner)(nil), "etop.Partner")
	proto.RegisterType((*PublicAccountInfo)(nil), "etop.PublicAccountInfo")
	proto.RegisterType((*PublicAuthorizedPartnerInfo)(nil), "etop.PublicAuthorizedPartnerInfo")
	proto.RegisterType((*PublicAccountsResponse)(nil), "etop.PublicAccountsResponse")
	proto.RegisterType((*Shop)(nil), "etop.Shop")
	proto.RegisterType((*ShippingServiceSelectStrategyItem)(nil), "etop.ShippingServiceSelectStrategyItem")
	proto.RegisterType((*SurveyInfo)(nil), "etop.SurveyInfo")
	proto.RegisterType((*CreateUserRequest)(nil), "etop.CreateUserRequest")
	proto.RegisterType((*RegisterResponse)(nil), "etop.RegisterResponse")
	proto.RegisterType((*ResetPasswordRequest)(nil), "etop.ResetPasswordRequest")
	proto.RegisterType((*ChangePasswordRequest)(nil), "etop.ChangePasswordRequest")
	proto.RegisterType((*ChangePasswordUsingTokenRequest)(nil), "etop.ChangePasswordUsingTokenRequest")
	proto.RegisterType((*Permission)(nil), "etop.Permission")
	proto.RegisterType((*LoginRequest)(nil), "etop.LoginRequest")
	proto.RegisterType((*LoginAccount)(nil), "etop.LoginAccount")
	proto.RegisterType((*LoginResponse)(nil), "etop.LoginResponse")
	proto.RegisterType((*SwitchAccountRequest)(nil), "etop.SwitchAccountRequest")
	proto.RegisterType((*UpgradeAccessTokenRequest)(nil), "etop.UpgradeAccessTokenRequest")
	proto.RegisterType((*AccessTokenResponse)(nil), "etop.AccessTokenResponse")
	proto.RegisterType((*SendSTokenEmailRequest)(nil), "etop.SendSTokenEmailRequest")
	proto.RegisterType((*SendEmailVerificationRequest)(nil), "etop.SendEmailVerificationRequest")
	proto.RegisterType((*SendPhoneVerificationRequest)(nil), "etop.SendPhoneVerificationRequest")
	proto.RegisterType((*VerifyEmailUsingTokenRequest)(nil), "etop.VerifyEmailUsingTokenRequest")
	proto.RegisterType((*VerifyPhoneUsingTokenRequest)(nil), "etop.VerifyPhoneUsingTokenRequest")
	proto.RegisterType((*InviteUserToAccountRequest)(nil), "etop.InviteUserToAccountRequest")
	proto.RegisterType((*AnswerInvitationRequest)(nil), "etop.AnswerInvitationRequest")
	proto.RegisterType((*GetUsersInCurrentAccountsRequest)(nil), "etop.GetUsersInCurrentAccountsRequest")
	proto.RegisterType((*PublicUserInfo)(nil), "etop.PublicUserInfo")
	proto.RegisterType((*UserAccountInfo)(nil), "etop.UserAccountInfo")
	proto.RegisterType((*ProtectedUsersResponse)(nil), "etop.ProtectedUsersResponse")
	proto.RegisterType((*LeaveAccountRequest)(nil), "etop.LeaveAccountRequest")
	proto.RegisterType((*RemoveUserFromCurrentAccountRequest)(nil), "etop.RemoveUserFromCurrentAccountRequest")
	proto.RegisterType((*Province)(nil), "etop.Province")
	proto.RegisterType((*GetProvincesResponse)(nil), "etop.GetProvincesResponse")
	proto.RegisterType((*District)(nil), "etop.District")
	proto.RegisterType((*GetDistrictsResponse)(nil), "etop.GetDistrictsResponse")
	proto.RegisterType((*GetDistrictsByProvinceRequest)(nil), "etop.GetDistrictsByProvinceRequest")
	proto.RegisterType((*Ward)(nil), "etop.Ward")
	proto.RegisterType((*GetWardsResponse)(nil), "etop.GetWardsResponse")
	proto.RegisterType((*GetWardsByDistrictRequest)(nil), "etop.GetWardsByDistrictRequest")
	proto.RegisterType((*ParseLocationRequest)(nil), "etop.ParseLocationRequest")
	proto.RegisterType((*ParseLocationResponse)(nil), "etop.ParseLocationResponse")
	proto.RegisterType((*Bank)(nil), "etop.Bank")
	proto.RegisterType((*GetBanksResponse)(nil), "etop.GetBanksResponse")
	proto.RegisterType((*BankProvince)(nil), "etop.BankProvince")
	proto.RegisterType((*GetBankProvincesResponse)(nil), "etop.GetBankProvincesResponse")
	proto.RegisterType((*GetProvincesByBankResquest)(nil), "etop.GetProvincesByBankResquest")
	proto.RegisterType((*BankBranch)(nil), "etop.BankBranch")
	proto.RegisterType((*GetBranchesByBankProvinceResponse)(nil), "etop.GetBranchesByBankProvinceResponse")
	proto.RegisterType((*GetBranchesByBankProvinceResquest)(nil), "etop.GetBranchesByBankProvinceResquest")
	proto.RegisterType((*Address)(nil), "etop.Address")
	proto.RegisterType((*Coordinates)(nil), "etop.Coordinates")
	proto.RegisterType((*BankAccount)(nil), "etop.BankAccount")
	proto.RegisterType((*ContactPerson)(nil), "etop.ContactPerson")
	proto.RegisterType((*CompanyInfo)(nil), "etop.CompanyInfo")
	proto.RegisterType((*CreateAddressRequest)(nil), "etop.CreateAddressRequest")
	proto.RegisterType((*AddressNote)(nil), "etop.AddressNote")
	proto.RegisterType((*GetAddressResponse)(nil), "etop.GetAddressResponse")
	proto.RegisterType((*UpdateAddressRequest)(nil), "etop.UpdateAddressRequest")
	proto.RegisterType((*SetDefaultAddressRequest)(nil), "etop.SetDefaultAddressRequest")
	proto.RegisterType((*UpdateURLSlugRequest)(nil), "etop.UpdateURLSlugRequest")
	proto.RegisterType((*HistoryResponse)(nil), "etop.HistoryResponse")
	proto.RegisterType((*Credit)(nil), "etop.Credit")
	proto.RegisterType((*CreditsResponse)(nil), "etop.CreditsResponse")
	proto.RegisterType((*UpdatePermissionRequest)(nil), "etop.UpdatePermissionRequest")
	proto.RegisterType((*UpdatePermissionItem)(nil), "etop.UpdatePermissionItem")
	proto.RegisterType((*UpdatePermissionResponse)(nil), "etop.UpdatePermissionResponse")
	proto.RegisterType((*Device)(nil), "etop.Device")
	proto.RegisterType((*Notification)(nil), "etop.Notification")
	proto.RegisterType((*CreateDeviceRequest)(nil), "etop.CreateDeviceRequest")
	proto.RegisterType((*DeleteDeviceRequest)(nil), "etop.DeleteDeviceRequest")
	proto.RegisterType((*GetNotificationsRequest)(nil), "etop.GetNotificationsRequest")
	proto.RegisterType((*NotificationsResponse)(nil), "etop.NotificationsResponse")
	proto.RegisterType((*UpdateNotificationsRequest)(nil), "etop.UpdateNotificationsRequest")
	proto.RegisterType((*UpdateReferenceUserRequest)(nil), "etop.UpdateReferenceUserRequest")
	proto.RegisterType((*UpdateReferenceSaleRequest)(nil), "etop.UpdateReferenceSaleRequest")
	proto.RegisterType((*Affiliate)(nil), "etop.Affiliate")
}

func init() { proto.RegisterFile("etop/etop.proto", fileDescriptor_4324bfd8c89695e1) }

var fileDescriptor_4324bfd8c89695e1 = []byte{
	// 4282 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x3c, 0x4b, 0x6c, 0x1c, 0x47,
	0x76, 0xec, 0xf9, 0xcf, 0x1b, 0x92, 0x43, 0x36, 0x7f, 0x63, 0x5a, 0xa2, 0xa8, 0xd6, 0x7a, 0xa5,
	0x58, 0x5e, 0xd2, 0xa6, 0xac, 0x85, 0x77, 0xb3, 0x09, 0x40, 0x52, 0xa6, 0xcc, 0x5d, 0x59, 0x66,
	0x7a, 0xa8, 0x35, 0x92, 0x05, 0xd2, 0x68, 0x76, 0xd7, 0x0c, 0x7b, 0xd9, 0xd3, 0xdd, 0x5b, 0x5d,
	0x43, 0x6a, 0x7c, 0xcd, 0x25, 0x48, 0x90, 0x64, 0x81, 0x04, 0x41, 0xce, 0x39, 0xed, 0x2d, 0x87,
	0x5c, 0x72, 0x09, 0x92, 0xdc, 0x7c, 0xf4, 0x21, 0x40, 0x4e, 0x09, 0x6c, 0x2d, 0xb0, 0x97, 0x45,
	0x90, 0x00, 0x39, 0xe4, 0x90, 0x4b, 0x50, 0xbf, 0xee, 0xea, 0xe9, 0x99, 0xe1, 0x8c, 0xac, 0xcd,
	0xc2, 0x88, 0x2f, 0xe2, 0xf4, 0xfb, 0xf5, 0xab, 0x57, 0xef, 0xbd, 0x7a, 0xf5, 0xaa, 0x5a, 0xd0,
	0x44, 0x24, 0x8c, 0x76, 0xe9, 0x3f, 0x3b, 0x11, 0x0e, 0x49, 0xa8, 0x97, 0xe8, 0xef, 0xcd, 0xd5,
	0x6e, 0xd8, 0x0d, 0x19, 0x60, 0x97, 0xfe, 0xe2, 0xb8, 0xcd, 0x5b, 0xdd, 0x30, 0xec, 0xfa, 0x68,
	0x97, 0x3d, 0x9d, 0xf5, 0x3b, 0xbb, 0xc4, 0xeb, 0xa1, 0x98, 0xd8, 0x3d, 0xc1, 0xbc, 0xb9, 0xe2,
	0x84, 0xbd, 0x5e, 0x18, 0xec, 0xf2, 0x3f, 0x02, 0xb8, 0x29, 0x5e, 0xe1, 0xec, 0x12, 0x3c, 0xb0,
	0x18, 0xd6, 0x45, 0x02, 0x77, 0x33, 0xc1, 0xc5, 0xc4, 0x26, 0xfd, 0xf8, 0x81, 0xf8, 0x2b, 0xd0,
	0xdb, 0x09, 0xda, 0x76, 0x5d, 0x8c, 0xe2, 0xd8, 0x22, 0x83, 0x08, 0xed, 0xd2, 0x7f, 0x04, 0xc5,
	0x37, 0x13, 0x8a, 0x7e, 0x8c, 0xb0, 0x15, 0x87, 0x7d, 0xec, 0x20, 0xf5, 0x37, 0xa7, 0x33, 0x6c,
	0x58, 0x7b, 0x8c, 0xc8, 0x71, 0x70, 0xe9, 0x11, 0x9b, 0x78, 0x61, 0x10, 0x9b, 0xe8, 0x27, 0x7d,
	0x14, 0x13, 0xdd, 0x80, 0x4a, 0x64, 0x77, 0xbd, 0xa0, 0xdb, 0xd2, 0xb6, 0xb5, 0x7b, 0x8d, 0x3d,
	0xd8, 0x71, 0x7a, 0x3b, 0x27, 0x0c, 0x62, 0x0a, 0x8c, 0xfe, 0x0d, 0xa8, 0x76, 0x3c, 0x9f, 0x20,
	0x1c, 0xb7, 0x0a, 0xdb, 0x45, 0x49, 0x74, 0xc4, 0x40, 0xa6, 0x44, 0x19, 0xdf, 0x81, 0xd7, 0x33,
	0xaf, 0x38, 0x18, 0x9c, 0x86, 0x17, 0x28, 0x90, 0x2f, 0xda, 0x84, 0x32, 0xa1, 0xcf, 0xec, 0x3d,
	0xf5, 0x83, 0xd2, 0xa7, 0xff, 0x76, 0x6b, 0xce, 0xe4, 0x20, 0xe3, 0x2f, 0x4a, 0x00, 0x29, 0xa3,
	0xbe, 0x0a, 0x05, 0xcf, 0x65, 0x74, 0x45, 0x41, 0x57, 0xf0, 0x5c, 0xfd, 0x26, 0x54, 0xe3, 0xf3,
	0x30, 0xb2, 0x3c, 0xb7, 0x55, 0x50, 0x50, 0x15, 0x0a, 0x3c, 0x76, 0xa9, 0x7c, 0xd4, 0xb3, 0x3d,
	0xbf, 0x55, 0x54, 0xe5, 0x33, 0x10, 0xc5, 0xe1, 0xd0, 0x47, 0x71, 0xab, 0xb4, 0x5d, 0x4c, 0x71,
	0x0c, 0x94, 0xea, 0x55, 0xce, 0xe9, 0xa5, 0x7f, 0x0b, 0x2a, 0x7c, 0x3e, 0x5a, 0x95, 0x6d, 0xed,
	0xde, 0xe2, 0x5e, 0x73, 0x47, 0x4c, 0xd3, 0x4e, 0x9b, 0xfd, 0x4d, 0x54, 0x60, 0x4f, 0xfa, 0x1d,
	0x00, 0x8f, 0x8e, 0x02, 0xb9, 0xd6, 0xd9, 0xa0, 0x55, 0x55, 0x94, 0xac, 0x0b, 0xf8, 0xc1, 0x40,
	0xff, 0x4d, 0x68, 0xd8, 0x8e, 0x83, 0x22, 0x4a, 0x65, 0x93, 0x56, 0x8d, 0x59, 0x7d, 0x73, 0x87,
	0xbb, 0xd6, 0x8e, 0x74, 0xad, 0x9d, 0x53, 0xe9, 0x5a, 0x26, 0x48, 0xf2, 0x7d, 0x42, 0x99, 0x5d,
	0xe4, 0xf8, 0x5e, 0xc0, 0x99, 0xeb, 0xd7, 0x33, 0x4b, 0xf2, 0x7d, 0xa2, 0x7f, 0x07, 0x00, 0x3d,
	0x8f, 0x3c, 0xcc, 0x79, 0xe1, 0x5a, 0xde, 0xba, 0xa0, 0xe6, 0xac, 0x0e, 0x46, 0xb6, 0xd0, 0xb9,
	0x71, 0x3d, 0xab, 0xa0, 0xe6, 0xac, 0xfd, 0xc8, 0x95, 0xac, 0xf3, 0xd7, 0xb3, 0x0a, 0xea, 0x7d,
	0x62, 0x3c, 0x84, 0x8d, 0x7d, 0x36, 0xf6, 0xd4, 0x37, 0xa6, 0xf1, 0xa6, 0x87, 0xb0, 0x61, 0xa2,
	0x1f, 0x23, 0x67, 0x46, 0xb6, 0x3f, 0x2e, 0x43, 0xe9, 0x59, 0x8c, 0xf0, 0x18, 0xf7, 0xbb, 0x0d,
	0xf5, 0x4e, 0xdf, 0xf7, 0xad, 0xc0, 0xee, 0x21, 0xe6, 0x80, 0x92, 0xbd, 0x46, 0xc1, 0x4f, 0xed,
	0x1e, 0xa2, 0xf3, 0x1f, 0x9f, 0x87, 0x98, 0x70, 0x1a, 0xd5, 0x0f, 0xeb, 0x0c, 0xce, 0x88, 0x36,
	0xa1, 0x1c, 0x9d, 0x87, 0x01, 0x6a, 0x95, 0x54, 0x15, 0x18, 0x28, 0xf5, 0xe1, 0x72, 0xde, 0x87,
	0xb3, 0x53, 0x50, 0x79, 0xf9, 0x29, 0xa8, 0xce, 0x30, 0x05, 0xfa, 0x11, 0x2c, 0xb3, 0xd7, 0x5b,
	0x97, 0x08, 0x7b, 0x1d, 0x6f, 0x5a, 0x9f, 0x6d, 0x32, 0xa6, 0x1f, 0x0a, 0x1e, 0x2e, 0x87, 0x0d,
	0x31, 0x23, 0xe7, 0x7a, 0xf7, 0x6d, 0x32, 0x26, 0x45, 0xce, 0xc7, 0xb0, 0xa9, 0xea, 0xe3, 0xb0,
	0xd9, 0xb5, 0x62, 0x14, 0x90, 0xe9, 0x7c, 0x7a, 0x43, 0x51, 0x8c, 0x33, 0xb7, 0x51, 0x40, 0xb8,
	0x60, 0x55, 0xc1, 0x21, 0xc1, 0xd7, 0x7b, 0xfc, 0x86, 0xa2, 0x69, 0x46, 0xf0, 0x43, 0xa8, 0xf0,
	0x4c, 0xcc, 0x7c, 0x7f, 0x71, 0x6f, 0x63, 0x47, 0xcd, 0xce, 0xd4, 0xe1, 0xda, 0xec, 0x67, 0x92,
	0x4b, 0xd8, 0x93, 0xf1, 0x01, 0xc0, 0xf1, 0xa3, 0x24, 0x4b, 0x2f, 0x41, 0xd1, 0x73, 0xe3, 0x96,
	0xb6, 0x5d, 0xbc, 0x57, 0x34, 0xe9, 0x4f, 0xfd, 0x1e, 0x94, 0x7b, 0xde, 0x73, 0xc4, 0x73, 0x61,
	0x63, 0x4f, 0xdf, 0x61, 0x6b, 0xd8, 0x87, 0x14, 0xb4, 0xef, 0x38, 0x61, 0x3f, 0x20, 0x26, 0x27,
	0x30, 0x7e, 0x04, 0xf3, 0x2a, 0x78, 0x84, 0xac, 0x75, 0x28, 0xda, 0xbe, 0xcf, 0x24, 0xd5, 0x84,
	0x1a, 0x14, 0x40, 0x5d, 0xde, 0xf6, 0x7d, 0x8b, 0x26, 0xd8, 0x98, 0xb9, 0xb3, 0xc4, 0xd6, 0x6c,
	0xdf, 0x6f, 0x53, 0xa8, 0xf1, 0x37, 0x45, 0xa8, 0x9e, 0xd8, 0x98, 0x04, 0x63, 0xe3, 0xa6, 0x05,
	0xa5, 0x5c, 0xc8, 0x30, 0x88, 0xfe, 0x06, 0x34, 0xa2, 0xfe, 0x99, 0xef, 0x39, 0x3c, 0x5e, 0x5a,
	0x0a, 0x01, 0x70, 0x04, 0x0b, 0x98, 0x34, 0x09, 0x17, 0xa7, 0x49, 0xc2, 0x37, 0xa1, 0xea, 0xc5,
	0x16, 0x41, 0x31, 0x61, 0x11, 0x26, 0x55, 0xae, 0x78, 0xf1, 0x29, 0xb5, 0xe4, 0xf7, 0xa0, 0xe9,
	0x84, 0x01, 0xb1, 0x1d, 0x62, 0x45, 0x08, 0xc7, 0x61, 0x10, 0xb7, 0x6a, 0x6c, 0x4d, 0x5b, 0xe1,
	0x16, 0x3c, 0xe4, 0xc8, 0x13, 0x86, 0x33, 0x17, 0x1d, 0xf5, 0x31, 0x4e, 0x83, 0xb7, 0x9e, 0x0f,
	0xde, 0x37, 0xa0, 0x71, 0x85, 0xce, 0x62, 0x8f, 0x20, 0xab, 0x8f, 0xfd, 0xd6, 0xa2, 0x3a, 0x1c,
	0x81, 0x78, 0x86, 0x99, 0x51, 0xbd, 0x9e, 0xdd, 0xe5, 0x44, 0x4d, 0x35, 0x8f, 0x30, 0x30, 0x25,
	0x49, 0xd2, 0xc0, 0x52, 0x3e, 0x0d, 0xdc, 0x82, 0x5a, 0x78, 0x15, 0x20, 0x4c, 0x97, 0xc1, 0x75,
	0xc5, 0xd4, 0x55, 0x06, 0x3d, 0x76, 0xf5, 0x2d, 0x28, 0x51, 0x07, 0x6b, 0x6d, 0x88, 0xe5, 0x9c,
	0x8d, 0x8a, 0xba, 0x99, 0xc9, 0xe0, 0xc6, 0xdf, 0x6a, 0xb0, 0x7c, 0xc2, 0xac, 0x2b, 0x1c, 0xe2,
	0x38, 0xe8, 0x84, 0x33, 0xcf, 0xdd, 0x7d, 0x28, 0xd1, 0x2a, 0x44, 0x4c, 0xc9, 0x32, 0x7f, 0x8b,
	0x10, 0x78, 0x3a, 0x88, 0xa4, 0x37, 0x33, 0xa2, 0xec, 0x90, 0x4b, 0x23, 0x87, 0xbc, 0x05, 0x55,
	0x61, 0xa3, 0x4c, 0xee, 0x93, 0x40, 0xe3, 0x17, 0x1a, 0xbc, 0x2e, 0xb4, 0xee, 0x93, 0xf3, 0x10,
	0x7b, 0x9f, 0x20, 0x57, 0xf8, 0xdd, 0x57, 0x41, 0x7f, 0xfd, 0x2e, 0xcc, 0x63, 0xe4, 0x7a, 0x18,
	0x39, 0x84, 0x49, 0xa9, 0x28, 0x44, 0x0d, 0x89, 0x79, 0x86, 0x7d, 0xe3, 0x43, 0x58, 0xcf, 0xcc,
	0x4e, 0x6c, 0xa2, 0x38, 0x0a, 0x83, 0x18, 0xe9, 0x0f, 0xa0, 0x66, 0x0b, 0x18, 0x0b, 0xde, 0xc6,
	0xde, 0x06, 0x57, 0x3b, 0x37, 0x9b, 0x66, 0x42, 0x68, 0xfc, 0x6b, 0x0d, 0x4a, 0x34, 0x52, 0xf5,
	0xbb, 0xd0, 0x44, 0xcf, 0xa3, 0x10, 0xd3, 0x45, 0xa0, 0xe3, 0x21, 0xdf, 0x8d, 0x5b, 0x2e, 0x2d,
	0x86, 0xcc, 0x45, 0x09, 0x3e, 0x62, 0x50, 0xfd, 0x21, 0xac, 0x78, 0xc1, 0x25, 0x0a, 0x48, 0x48,
	0xab, 0xd5, 0x4b, 0x84, 0x63, 0x12, 0x3a, 0x17, 0xad, 0x2d, 0x25, 0x96, 0xf4, 0x84, 0xe0, 0x23,
	0x89, 0x9f, 0x79, 0x02, 0x5e, 0x6d, 0x54, 0xdf, 0x85, 0xaa, 0xa8, 0x90, 0x99, 0xf9, 0x1b, 0x7b,
	0x0b, 0x62, 0x46, 0x39, 0xd0, 0x94, 0xd8, 0x89, 0x01, 0xfc, 0x2e, 0xcc, 0x9f, 0xd9, 0xc1, 0x85,
	0x25, 0x8c, 0x27, 0x56, 0x13, 0xe1, 0x1b, 0x07, 0x76, 0x70, 0x21, 0x13, 0x6b, 0xe3, 0x2c, 0x7d,
	0xd0, 0xdf, 0x82, 0xa6, 0xdd, 0x27, 0xa1, 0xc5, 0x97, 0x5b, 0xab, 0xd3, 0xe9, 0xb5, 0x56, 0x15,
	0x0d, 0x17, 0x28, 0xf2, 0x90, 0xe1, 0x8e, 0x3a, 0xbd, 0xff, 0xa3, 0x24, 0xf1, 0x36, 0x2c, 0x47,
	0x38, 0x74, 0xfb, 0x0e, 0x11, 0xeb, 0x0c, 0xcd, 0x16, 0xcb, 0xca, 0xdc, 0x34, 0x05, 0x9a, 0x2f,
	0x3d, 0xc7, 0xae, 0xfe, 0x0e, 0xe8, 0xf1, 0xb9, 0x17, 0x59, 0x24, 0xb4, 0xe4, 0x56, 0xc3, 0x73,
	0x5b, 0xba, 0xca, 0x42, 0xf1, 0xa7, 0xa1, 0xb0, 0xe8, 0xb1, 0xab, 0x3f, 0x84, 0x55, 0xc6, 0xd2,
	0xc1, 0x61, 0x4f, 0x65, 0x5a, 0x51, 0x98, 0x96, 0x29, 0xc5, 0x11, 0x0e, 0x7b, 0x29, 0xdb, 0x3d,
	0x58, 0xe8, 0x9e, 0x07, 0x56, 0x10, 0x12, 0x64, 0xd1, 0x9d, 0x50, 0x6b, 0x4d, 0x0d, 0x85, 0xee,
	0x79, 0xf0, 0x34, 0x24, 0xe8, 0x30, 0x74, 0x91, 0xbe, 0x03, 0x15, 0xbe, 0x63, 0x62, 0x4b, 0x03,
	0x8d, 0x52, 0xfe, 0xb8, 0x73, 0x8a, 0x07, 0x1f, 0x05, 0x94, 0x24, 0x29, 0xe0, 0x28, 0xe0, 0x4b,
	0xa7, 0x46, 0xea, 0x00, 0x4e, 0xd8, 0x8b, 0xec, 0x60, 0x60, 0x79, 0x41, 0x27, 0x6c, 0xbd, 0xa6,
	0x3a, 0xc0, 0x21, 0xc7, 0xb0, 0xf8, 0x6a, 0x38, 0xe9, 0x83, 0xfe, 0x3d, 0xd8, 0xe8, 0x85, 0x01,
	0x1a, 0x58, 0x04, 0xdb, 0x41, 0x6c, 0x3b, 0xac, 0x70, 0xc0, 0xb8, 0xef, 0xa3, 0xd6, 0xa6, 0x32,
	0xb4, 0x35, 0x46, 0x74, 0x9a, 0xd2, 0x98, 0x94, 0x44, 0x7f, 0x07, 0x1a, 0x71, 0x1f, 0x5f, 0x22,
	0xf1, 0xca, 0xd7, 0x59, 0x60, 0x2f, 0xf1, 0x57, 0xb6, 0x19, 0x82, 0xbd, 0x11, 0xe2, 0xe4, 0xb7,
	0x1e, 0xc1, 0x36, 0x35, 0x6b, 0xe4, 0x05, 0x5d, 0x2b, 0x46, 0xf8, 0xd2, 0x73, 0x90, 0x15, 0x23,
	0x9f, 0xa6, 0x96, 0x98, 0x60, 0x9b, 0xa0, 0xee, 0xa0, 0x75, 0x83, 0xc9, 0xb9, 0x2b, 0xe4, 0x08,
	0xea, 0x36, 0x27, 0x6e, 0x33, 0xda, 0xb6, 0x20, 0x3d, 0x26, 0xa8, 0x67, 0xde, 0x8c, 0x27, 0x91,
	0xd0, 0x30, 0x66, 0x53, 0x75, 0x53, 0x0d, 0x63, 0x0a, 0x31, 0x3e, 0x86, 0xdb, 0xd7, 0x4a, 0xa7,
	0xf5, 0xc5, 0x05, 0x1a, 0x64, 0x6a, 0x6e, 0x0a, 0xa0, 0x2e, 0x7c, 0x69, 0xfb, 0xfd, 0x6c, 0x7a,
	0xe0, 0x20, 0xc3, 0x05, 0x48, 0x87, 0x3f, 0x56, 0xc2, 0x36, 0xd4, 0x58, 0x81, 0xe4, 0x85, 0x41,
	0xb6, 0x26, 0x97, 0x50, 0xfd, 0x06, 0x54, 0xec, 0x20, 0xbe, 0x42, 0x38, 0x53, 0x8f, 0x0b, 0x98,
	0xf1, 0x8b, 0x02, 0x2c, 0xf3, 0xe0, 0x64, 0x6e, 0x20, 0xaa, 0xad, 0x4c, 0xa9, 0xaf, 0x4d, 0x51,
	0xea, 0x17, 0xae, 0x29, 0xf5, 0x8b, 0x13, 0x4a, 0xfd, 0x52, 0x3e, 0x7c, 0xb7, 0xa1, 0x16, 0xd9,
	0x71, 0x7c, 0x15, 0x62, 0x37, 0xb3, 0x9a, 0x24, 0x50, 0x56, 0x99, 0x75, 0x31, 0x42, 0x16, 0x09,
	0xf9, 0xde, 0x34, 0xad, 0xcc, 0x28, 0xf8, 0x34, 0xa4, 0x05, 0xe2, 0x12, 0x27, 0xe1, 0xf5, 0x32,
	0xf3, 0x2e, 0x5a, 0xfa, 0xd7, 0xcc, 0x45, 0x06, 0x7f, 0x9f, 0x82, 0x99, 0x71, 0xef, 0xc3, 0x22,
	0x46, 0x5d, 0x2f, 0x26, 0x08, 0x5b, 0x7c, 0x77, 0x54, 0x53, 0x5e, 0xba, 0x20, 0x71, 0x6c, 0x37,
	0xaf, 0x94, 0xb3, 0xf5, 0x59, 0xca, 0xd9, 0x3d, 0x58, 0x32, 0x85, 0x9c, 0x64, 0x41, 0x93, 0xe1,
	0xa8, 0x8d, 0xa9, 0x54, 0x9e, 0xc2, 0xaa, 0x89, 0x62, 0x44, 0x4e, 0xc4, 0xa8, 0x95, 0x4d, 0x1c,
	0x37, 0xab, 0x36, 0xc1, 0xac, 0x85, 0x9c, 0x59, 0x8d, 0xbf, 0xd6, 0x60, 0xed, 0xf0, 0xdc, 0x0e,
	0xba, 0x68, 0x58, 0xe2, 0x2e, 0x2c, 0x39, 0x7d, 0x8c, 0x69, 0xb1, 0x9f, 0x18, 0x5e, 0x15, 0xd0,
	0x14, 0x58, 0xc9, 0x47, 0x97, 0xf3, 0x00, 0x5d, 0xa5, 0xc4, 0xea, 0x04, 0x37, 0x02, 0x74, 0x95,
	0x10, 0x52, 0xc9, 0x61, 0xd0, 0xf1, 0x70, 0x2f, 0x25, 0x2e, 0x65, 0x24, 0x73, 0xac, 0x64, 0x30,
	0x7e, 0xa9, 0xc1, 0xad, 0xac, 0x92, 0xcf, 0x62, 0x2f, 0xe8, 0x0e, 0xb7, 0x52, 0xf8, 0x20, 0xb5,
	0x91, 0xad, 0x0e, 0x6e, 0x9c, 0x42, 0xde, 0x38, 0xdf, 0x86, 0x55, 0x4c, 0x0d, 0x9a, 0xa8, 0x22,
	0xa6, 0x5b, 0xd5, 0x5e, 0xc7, 0xaa, 0xc9, 0xf9, 0x9c, 0x0f, 0x8f, 0xb6, 0x34, 0xcb, 0x68, 0xcb,
	0x93, 0x46, 0xfb, 0x08, 0xe0, 0x04, 0xe1, 0x9e, 0x17, 0xc7, 0xbc, 0xef, 0x23, 0xda, 0x34, 0x05,
	0x56, 0x99, 0x88, 0x06, 0xcd, 0x36, 0x34, 0xa2, 0x84, 0x86, 0x96, 0x0b, 0x14, 0xa7, 0x82, 0x8c,
	0x7f, 0xd6, 0x60, 0xfe, 0x49, 0xd8, 0xf5, 0x54, 0x03, 0xf9, 0xf4, 0x39, 0x6b, 0x20, 0x06, 0xca,
	0x04, 0x57, 0x61, 0x64, 0x70, 0xdd, 0x01, 0x10, 0x25, 0x00, 0x5d, 0x49, 0x8a, 0x6a, 0x1b, 0x47,
	0xc0, 0x8f, 0x5d, 0xfd, 0xbb, 0x30, 0x2f, 0x89, 0x58, 0x21, 0x59, 0x9a, 0x5c, 0x48, 0x36, 0xec,
	0x14, 0x44, 0x8b, 0x00, 0xc9, 0x4b, 0xb3, 0x9a, 0x6a, 0x21, 0xf9, 0xe6, 0x1f, 0xa0, 0x81, 0xf1,
	0x79, 0x41, 0x0c, 0x4b, 0x96, 0x1a, 0x53, 0xd7, 0x70, 0xbf, 0xd2, 0x6a, 0xf8, 0x2e, 0x1b, 0x39,
	0x6b, 0x46, 0xe6, 0xfa, 0x66, 0x0d, 0x8e, 0xe1, 0x6e, 0x73, 0x47, 0xf6, 0x9b, 0x62, 0xcb, 0x0b,
	0x58, 0x96, 0x2a, 0x4b, 0x3b, 0x0a, 0xf8, 0x71, 0x90, 0xad, 0x74, 0xaa, 0x23, 0x2b, 0x9d, 0x5b,
	0x50, 0xeb, 0x63, 0xdf, 0x8a, 0xfd, 0x7e, 0x37, 0x93, 0x99, 0xaa, 0x7d, 0xec, 0xb7, 0xfd, 0x7e,
	0x57, 0x7f, 0x0f, 0xe6, 0x59, 0x12, 0x92, 0x85, 0x1b, 0xef, 0x2b, 0xac, 0xa5, 0x09, 0x45, 0xad,
	0x8d, 0x1b, 0xfd, 0x14, 0x60, 0xfc, 0xb2, 0x08, 0x0b, 0xc2, 0x73, 0x44, 0x52, 0x1a, 0x1e, 0x9d,
	0x36, 0xdd, 0xe8, 0x0a, 0xa3, 0x47, 0x27, 0x53, 0x5c, 0x71, 0x4c, 0xc5, 0xf1, 0x16, 0x54, 0xa5,
	0xd2, 0x25, 0x75, 0x1f, 0xaf, 0x4e, 0xbb, 0x29, 0x49, 0xa8, 0x34, 0xba, 0x17, 0x17, 0x25, 0x2e,
	0xc8, 0xc5, 0x3d, 0x8c, 0x4c, 0x06, 0xd7, 0xbf, 0x05, 0x75, 0xbb, 0xd3, 0xf1, 0x7c, 0xcf, 0x26,
	0x48, 0x74, 0x88, 0x9a, 0x62, 0x2e, 0x25, 0xd8, 0x4c, 0x29, 0xf4, 0x7d, 0xd0, 0xed, 0x4b, 0xdb,
	0xf3, 0xed, 0x33, 0x1f, 0x59, 0xc9, 0xd6, 0xa2, 0xca, 0x2a, 0x87, 0x51, 0x7a, 0x2c, 0x27, 0xd4,
	0x72, 0x6f, 0xa2, 0x1f, 0xb1, 0xcd, 0x82, 0x68, 0xb2, 0xa5, 0x32, 0xf8, 0x8e, 0x7a, 0xcc, 0x04,
	0xe8, 0x29, 0x47, 0x22, 0xe7, 0x06, 0xdd, 0x0d, 0x30, 0x7b, 0xd7, 0xd5, 0xea, 0x9e, 0xc3, 0xf4,
	0x23, 0x58, 0xe6, 0xbf, 0x2c, 0x69, 0xf1, 0xa9, 0x7a, 0x3d, 0x4d, 0xce, 0xf4, 0x3e, 0xe7, 0xd9,
	0x27, 0x46, 0x00, 0xab, 0xed, 0x2b, 0x8f, 0x38, 0xe7, 0x72, 0x44, 0x22, 0x5d, 0x64, 0x03, 0x5e,
	0x1b, 0x1d, 0xf0, 0xef, 0xc0, 0x32, 0x46, 0x5d, 0x14, 0x20, 0x5a, 0xd8, 0x70, 0xe7, 0x88, 0x33,
	0x2d, 0x93, 0xa5, 0x14, 0xcd, 0x3c, 0x24, 0x36, 0x7e, 0x04, 0xaf, 0x3d, 0x8b, 0xba, 0xd8, 0x76,
	0xa9, 0xc1, 0xa4, 0xe3, 0xc8, 0x97, 0xa6, 0x43, 0x56, 0x5d, 0x4c, 0x0e, 0xf9, 0xda, 0x2c, 0x45,
	0x4b, 0x97, 0x95, 0x8c, 0xd8, 0xff, 0x7f, 0x0e, 0x9c, 0x9a, 0xb0, 0x36, 0xad, 0xd7, 0xd4, 0x67,
	0xf7, 0x9a, 0xdf, 0x85, 0xf5, 0x36, 0x0a, 0xdc, 0x36, 0xb3, 0x1a, 0xab, 0x9a, 0xa6, 0x59, 0x87,
	0xb3, 0x3e, 0x55, 0x18, 0xe9, 0x53, 0xc6, 0x77, 0xe1, 0x06, 0x15, 0xfd, 0xfe, 0x70, 0x4f, 0x72,
	0x8a, 0x17, 0x48, 0xde, 0x93, 0xe1, 0xb6, 0xe3, 0x14, 0x55, 0x92, 0xd1, 0x86, 0x1b, 0x8c, 0x65,
	0xc0, 0xde, 0x9c, 0x2f, 0x30, 0x1e, 0x80, 0x9e, 0x69, 0x83, 0xe6, 0x3d, 0x69, 0x59, 0xc5, 0x33,
	0xde, 0x54, 0x28, 0x53, 0xe9, 0x15, 0x09, 0xfd, 0x99, 0x06, 0x9b, 0xac, 0x8d, 0xcf, 0x0a, 0xf4,
	0xd3, 0xf0, 0x65, 0x22, 0xf7, 0x01, 0xf0, 0x94, 0x83, 0xe8, 0x36, 0x18, 0x05, 0xc4, 0xeb, 0x78,
	0x08, 0x67, 0xa2, 0x6a, 0x59, 0xe0, 0x8f, 0x13, 0xb4, 0xfe, 0x36, 0x40, 0x5a, 0x62, 0x08, 0xdf,
	0x16, 0xdb, 0xb2, 0xb4, 0x62, 0x31, 0x15, 0x1a, 0xe3, 0x02, 0x36, 0xf6, 0xd9, 0xae, 0x22, 0x7f,
	0xec, 0x30, 0x95, 0x9a, 0xf7, 0xa1, 0x86, 0x45, 0x10, 0x33, 0xe5, 0xf2, 0x3d, 0x11, 0x33, 0x21,
	0x30, 0xfe, 0x52, 0x83, 0xed, 0xc7, 0x88, 0x50, 0xa3, 0xc4, 0xc7, 0xc1, 0x21, 0x2f, 0x4f, 0xd3,
	0x96, 0xd1, 0x2b, 0x3e, 0xdb, 0x4b, 0xbb, 0xcd, 0xc5, 0xeb, 0xba, 0xcd, 0x01, 0x2c, 0xf2, 0x7e,
	0x14, 0x55, 0x6d, 0x42, 0x6b, 0xee, 0x15, 0x1d, 0xa7, 0x18, 0x7f, 0x54, 0x85, 0xe6, 0xd0, 0x0a,
	0xa3, 0xdf, 0x84, 0x2a, 0xab, 0x07, 0x86, 0x5e, 0x5b, 0xa1, 0xc0, 0x63, 0x57, 0x7f, 0x13, 0x16,
	0x19, 0x7a, 0xf4, 0xfb, 0x59, 0x29, 0x71, 0x24, 0x75, 0x78, 0x0b, 0x9a, 0x7c, 0x7f, 0x33, 0x5a,
	0x91, 0x05, 0x8a, 0x6c, 0x27, 0x1b, 0xbe, 0xec, 0x3c, 0x97, 0x46, 0xcf, 0xf3, 0xdd, 0xb4, 0x72,
	0x64, 0xf2, 0x86, 0xeb, 0x27, 0x8a, 0x61, 0xd2, 0x86, 0x4b, 0xcc, 0xca, 0x0c, 0x25, 0x26, 0x5d,
	0x3f, 0xc2, 0xd8, 0x63, 0x1b, 0xe3, 0x4c, 0x55, 0x25, 0xa1, 0x43, 0x0e, 0x5e, 0xbb, 0xde, 0xc1,
	0x95, 0x96, 0x5d, 0x7d, 0x9a, 0x96, 0xdd, 0x6f, 0x43, 0x53, 0xba, 0xab, 0x25, 0xf8, 0x60, 0x12,
	0xdf, 0xa2, 0xa4, 0xe6, 0x50, 0x7d, 0x0f, 0x94, 0x4a, 0x81, 0x9f, 0xc4, 0x9c, 0x0d, 0xd8, 0x49,
	0x8c, 0x34, 0xea, 0x52, 0x8a, 0x6f, 0xa3, 0x80, 0x1c, 0x0c, 0xf4, 0x0f, 0xf2, 0x3c, 0x53, 0x1d,
	0x3a, 0x0e, 0x49, 0xda, 0x27, 0xfa, 0x09, 0xac, 0x67, 0x2b, 0x9b, 0xe4, 0xc4, 0x76, 0xe1, 0x5a,
	0x69, 0xab, 0x99, 0x0a, 0x47, 0x9e, 0xdd, 0x66, 0x25, 0x62, 0x76, 0x42, 0xc9, 0x25, 0x2e, 0xce,
	0x22, 0xd1, 0x14, 0x8c, 0xe2, 0x34, 0xd8, 0x8b, 0x69, 0x41, 0xc6, 0xc4, 0x34, 0xa7, 0x38, 0x0d,
	0x16, 0xe4, 0xfb, 0x84, 0x6e, 0x42, 0x12, 0xe6, 0xb3, 0x01, 0xeb, 0x22, 0x4a, 0xbb, 0x26, 0x64,
	0x07, 0x03, 0xfd, 0x3e, 0x2c, 0x8a, 0x27, 0x0b, 0x23, 0x3b, 0x0e, 0x03, 0xd6, 0x47, 0x4c, 0xfc,
	0x5f, 0xe0, 0x4c, 0x86, 0x32, 0x2e, 0x60, 0xfd, 0x04, 0x87, 0x84, 0xe9, 0xc7, 0x52, 0x53, 0x52,
	0x95, 0x7c, 0x63, 0x28, 0x15, 0xcd, 0x8b, 0x54, 0x84, 0x58, 0x49, 0x28, 0x93, 0xd1, 0x7d, 0x28,
	0xd3, 0x80, 0x92, 0xa9, 0x68, 0x4c, 0x01, 0xc9, 0x69, 0x8c, 0x35, 0x58, 0x79, 0x82, 0xec, 0x4b,
	0x94, 0x5d, 0x12, 0x8c, 0x37, 0xe0, 0x8e, 0x89, 0x7a, 0xe1, 0x25, 0x5b, 0x30, 0x8e, 0x70, 0xd8,
	0xcb, 0x26, 0x47, 0x49, 0xf6, 0x07, 0x1a, 0xd4, 0x4e, 0x70, 0x78, 0xe9, 0x05, 0x0e, 0x4a, 0xfa,
	0x5b, 0xda, 0x70, 0x7f, 0x6b, 0xc2, 0x9e, 0xe9, 0x06, 0x54, 0x30, 0xea, 0xd2, 0xd8, 0xc9, 0x34,
	0x96, 0x38, 0x8c, 0x5a, 0x97, 0xff, 0xe2, 0x3d, 0x4e, 0x35, 0x15, 0x00, 0x47, 0x1c, 0x86, 0x2e,
	0x32, 0x1e, 0xc1, 0xea, 0x63, 0x44, 0xa4, 0x1e, 0xa9, 0xb9, 0xde, 0x82, 0x7a, 0x24, 0x81, 0xa2,
	0xd9, 0xbf, 0x28, 0x62, 0x53, 0x80, 0xcd, 0x94, 0xc0, 0xf8, 0xa9, 0x06, 0xb5, 0x47, 0x5e, 0x4c,
	0xb0, 0xe7, 0x90, 0x09, 0x63, 0xf9, 0x0d, 0x58, 0x90, 0x3c, 0x5c, 0xab, 0x4c, 0xda, 0x93, 0xa8,
	0x43, 0x75, 0xd8, 0xc5, 0x51, 0x87, 0x76, 0x5e, 0x6c, 0x75, 0x30, 0x42, 0xf1, 0xb9, 0x17, 0x65,
	0x9a, 0xf1, 0xe0, 0xc5, 0x47, 0x02, 0x2e, 0x06, 0x26, 0x95, 0xca, 0x0c, 0xcc, 0x95, 0xc0, 0xec,
	0xc0, 0x24, 0xad, 0x99, 0x12, 0x18, 0x7d, 0xb8, 0xa9, 0x4a, 0x39, 0x18, 0x24, 0xa3, 0x17, 0x2b,
	0x5c, 0x6e, 0x48, 0xda, 0xd8, 0x21, 0xa9, 0xa4, 0xf9, 0xa4, 0x2f, 0x51, 0x6c, 0x4d, 0x41, 0x50,
	0xfa, 0xd8, 0xc6, 0xee, 0x64, 0x53, 0x4a, 0x2d, 0x47, 0x98, 0x52, 0xa2, 0x26, 0x9b, 0xd2, 0x78,
	0x17, 0x96, 0x1e, 0x23, 0x42, 0xdf, 0x94, 0xda, 0x67, 0x1b, 0xca, 0x57, 0x14, 0x20, 0x6c, 0x23,
	0x4a, 0x64, 0x4a, 0x63, 0x72, 0x84, 0xf1, 0x13, 0x78, 0x4d, 0x72, 0x1d, 0x0c, 0x12, 0xa3, 0xa5,
	0xf6, 0xc8, 0xea, 0xa5, 0x8d, 0xd5, 0x4b, 0x25, 0xcd, 0xdb, 0x43, 0xa2, 0x98, 0x3d, 0xfe, 0x4c,
	0x83, 0xd5, 0x13, 0x1b, 0xc7, 0xe8, 0x49, 0x98, 0xad, 0x31, 0x73, 0x36, 0xd5, 0xc6, 0xd9, 0x74,
	0x86, 0xd7, 0xd1, 0xd2, 0x80, 0x0e, 0x35, 0xbf, 0xda, 0xd6, 0x28, 0x58, 0x6a, 0xb4, 0x36, 0xa4,
	0x91, 0x30, 0xe0, 0x9b, 0x50, 0x93, 0xef, 0x15, 0xa9, 0x66, 0x38, 0x70, 0x12, 0x3c, 0xa5, 0x95,
	0x2f, 0x16, 0xc7, 0xe8, 0xc3, 0xbe, 0x98, 0xe0, 0xe9, 0xd6, 0x85, 0xbe, 0x3d, 0xbb, 0x11, 0x62,
	0xf3, 0xc2, 0xe0, 0xc6, 0x29, 0x94, 0x0e, 0xec, 0xe0, 0xe2, 0xa5, 0x52, 0x49, 0x4b, 0x69, 0xbf,
	0xd4, 0xd5, 0x5e, 0x8b, 0x70, 0x11, 0x2a, 0x38, 0xe3, 0x22, 0x67, 0x14, 0x90, 0x75, 0x11, 0x4a,
	0x63, 0x72, 0x84, 0x81, 0x60, 0x9e, 0x3e, 0x7e, 0xa9, 0xf4, 0x76, 0x1b, 0xea, 0xec, 0x30, 0x8c,
	0x31, 0x66, 0x26, 0x81, 0x82, 0x59, 0xf2, 0x7a, 0x02, 0x2d, 0xa1, 0x5c, 0x3e, 0x81, 0xbd, 0x9d,
	0x4f, 0x60, 0x7a, 0xaa, 0xe8, 0xa8, 0x24, 0x76, 0x06, 0x9b, 0x6a, 0x2a, 0x3c, 0x18, 0xb0, 0x11,
	0xa1, 0x38, 0x69, 0xc9, 0xa7, 0xea, 0x68, 0xa3, 0xd4, 0x49, 0x48, 0xf2, 0x15, 0x25, 0x05, 0x33,
	0xb7, 0xf9, 0x53, 0x0d, 0x80, 0x8a, 0x3d, 0xc0, 0x76, 0xe0, 0x9c, 0xff, 0x8a, 0xec, 0x92, 0x4f,
	0x4a, 0xa5, 0x71, 0x49, 0xc9, 0xf8, 0x1d, 0xb8, 0x4d, 0x4d, 0xc8, 0xd4, 0x91, 0x63, 0x4e, 0x73,
	0x5c, 0x92, 0x33, 0x6b, 0x67, 0x82, 0x42, 0x98, 0x72, 0x29, 0x35, 0x25, 0xe7, 0x35, 0x13, 0x0a,
	0xe3, 0x1f, 0xb4, 0xc9, 0x32, 0x5f, 0xa1, 0x3d, 0xf3, 0x23, 0x2d, 0x4e, 0x9f, 0x7e, 0x4b, 0x63,
	0xd3, 0xef, 0x7f, 0x95, 0xa1, 0x2a, 0xce, 0x0b, 0xbf, 0x6c, 0xcb, 0x73, 0x5b, 0xc9, 0x06, 0xd9,
	0x36, 0x8a, 0x8c, 0x8d, 0xdc, 0x10, 0x2a, 0x63, 0x87, 0xb0, 0xad, 0xa4, 0x8b, 0xcc, 0xcc, 0x27,
	0x49, 0x22, 0x97, 0x7e, 0xab, 0x93, 0x96, 0x05, 0x96, 0x4f, 0x54, 0x33, 0x30, 0x48, 0x92, 0xfe,
	0x98, 0x80, 0xda, 0x70, 0xfa, 0x93, 0x9a, 0x88, 0x03, 0xd7, 0x77, 0xb2, 0x07, 0x44, 0x12, 0xaa,
	0x50, 0xec, 0x65, 0x8e, 0xba, 0x13, 0xa8, 0xbe, 0x0e, 0xc5, 0x4f, 0xbc, 0x88, 0x95, 0xe4, 0xc9,
	0x91, 0xda, 0x27, 0x5e, 0xa4, 0x6f, 0x41, 0x95, 0x15, 0x4a, 0x98, 0xd7, 0xda, 0x49, 0xb3, 0x55,
	0x00, 0xb3, 0x1b, 0xb7, 0xf9, 0x71, 0x1b, 0xb7, 0x8e, 0x87, 0x63, 0x91, 0xe8, 0x17, 0xd4, 0x8d,
	0x1b, 0x83, 0xcb, 0x2c, 0xef, 0xdb, 0x92, 0x46, 0x3d, 0x07, 0xaf, 0x51, 0x70, 0xf6, 0xfc, 0xac,
	0x39, 0xe1, 0xa0, 0x67, 0x69, 0xf4, 0xf9, 0x99, 0xdc, 0xfc, 0x2c, 0x8f, 0xdc, 0xfc, 0x3c, 0x10,
	0x19, 0x57, 0x67, 0x1b, 0x92, 0xd7, 0x76, 0xd4, 0xeb, 0xb5, 0xf2, 0xd2, 0xc0, 0x88, 0xc6, 0x77,
	0x39, 0x08, 0x09, 0x8a, 0xd9, 0x09, 0x77, 0x72, 0x2e, 0x2c, 0xa8, 0x9f, 0x86, 0x04, 0x99, 0x1c,
	0xaf, 0x3f, 0x80, 0x86, 0x13, 0x86, 0xd8, 0xf5, 0x02, 0x9b, 0x92, 0xaf, 0x66, 0x8f, 0x91, 0x13,
	0x84, 0xa9, 0x52, 0x19, 0x6d, 0x68, 0x28, 0x38, 0x3a, 0x06, 0xdf, 0x26, 0x1e, 0xe9, 0x8b, 0xf8,
	0x2c, 0xa4, 0xd6, 0xe1, 0x50, 0xdd, 0x80, 0xba, 0x1f, 0x06, 0x5d, 0x4e, 0x52, 0x50, 0x48, 0x52,
	0xb0, 0xf1, 0xf7, 0x1a, 0x34, 0x94, 0x9b, 0x0b, 0x49, 0x5e, 0xd3, 0x72, 0x79, 0xed, 0xfa, 0x48,
	0xb9, 0x01, 0x15, 0x9e, 0x64, 0xb2, 0x05, 0x2f, 0x87, 0xd1, 0x7d, 0x42, 0xb2, 0xab, 0xed, 0xf7,
	0xce, 0x10, 0xce, 0x78, 0xf6, 0x82, 0xdc, 0xd7, 0x32, 0xd4, 0xd4, 0x5b, 0x60, 0x5a, 0xa5, 0x2f,
	0x64, 0x6e, 0x64, 0x5d, 0x33, 0x02, 0x39, 0xeb, 0x85, 0x91, 0xb3, 0xfe, 0x92, 0xe7, 0xb1, 0xc6,
	0x67, 0x1a, 0x9d, 0x9b, 0xf4, 0xc4, 0x7f, 0xbc, 0x0e, 0xb7, 0xa0, 0x46, 0xec, 0xe7, 0xf9, 0x92,
	0xb0, 0x4a, 0xec, 0xe7, 0x2c, 0x72, 0xb7, 0xd2, 0x8b, 0x2a, 0xaa, 0x12, 0xc9, 0xfd, 0x94, 0xeb,
	0xee, 0x11, 0x1d, 0xc1, 0xaa, 0x8f, 0xba, 0xb6, 0x6f, 0x61, 0x14, 0x61, 0x44, 0x37, 0xb8, 0x36,
	0xf1, 0x2e, 0x91, 0x68, 0x50, 0x8d, 0xbc, 0xc3, 0xb6, 0xc2, 0x18, 0xcc, 0x0c, 0xbd, 0xf1, 0x4f,
	0x65, 0x58, 0xe5, 0x07, 0xdf, 0xf2, 0x8a, 0x8c, 0x28, 0xe9, 0xbe, 0xce, 0x98, 0x5f, 0x67, 0xcc,
	0xaf, 0x48, 0xc6, 0xfc, 0x13, 0x0d, 0x1a, 0x8a, 0x2c, 0x16, 0x96, 0x21, 0x19, 0x0e, 0x4b, 0x8a,
	0xb9, 0x0d, 0xf5, 0x30, 0x42, 0x81, 0x45, 0xbc, 0xe1, 0x52, 0x86, 0x82, 0x4f, 0x3d, 0xbe, 0xaf,
	0xf5, 0xfb, 0x81, 0x73, 0x6e, 0x9d, 0x61, 0x64, 0x5f, 0x64, 0xbc, 0x15, 0x18, 0xe2, 0x80, 0xc2,
	0xa9, 0x11, 0x43, 0x72, 0x3e, 0x94, 0xdd, 0x38, 0xc8, 0xd8, 0x07, 0xfd, 0x31, 0x22, 0x49, 0x3c,
	0x89, 0xea, 0xed, 0x3e, 0xd4, 0x85, 0xad, 0x92, 0xf2, 0x6d, 0xe8, 0x72, 0x5a, 0x8a, 0x37, 0xfe,
	0xa5, 0x0c, 0xab, 0xcf, 0xd8, 0xe5, 0xeb, 0xa1, 0xb0, 0xfc, 0xba, 0xbc, 0xf9, 0x3a, 0x58, 0xbf,
	0xd2, 0xc1, 0x8a, 0xa0, 0xd5, 0x46, 0xe4, 0x11, 0xea, 0xd8, 0x7d, 0x9f, 0x4c, 0xe5, 0xdc, 0x72,
	0x10, 0x85, 0x19, 0x06, 0x61, 0xfc, 0x50, 0xc6, 0xcf, 0x33, 0xf3, 0x49, 0xdb, 0xef, 0x77, 0x67,
	0x3a, 0x81, 0x79, 0x4d, 0xb9, 0x68, 0xc0, 0xc2, 0x29, 0xb9, 0x62, 0x60, 0xfc, 0x3e, 0x34, 0x3f,
	0xf0, 0x62, 0x12, 0xe2, 0xc1, 0x8c, 0x2d, 0xcd, 0x37, 0xa0, 0xe4, 0xda, 0xc4, 0x16, 0xfd, 0x85,
	0x65, 0x4a, 0x63, 0xda, 0x57, 0xdf, 0x6f, 0x7f, 0xf4, 0xf4, 0xa3, 0xb3, 0x1f, 0x23, 0x87, 0x98,
	0x0c, 0x6d, 0xfc, 0x47, 0x01, 0x2a, 0x87, 0x18, 0xb9, 0xde, 0x38, 0x6b, 0xdc, 0x80, 0x8a, 0xdd,
	0x63, 0xe7, 0xac, 0x99, 0x8f, 0x9f, 0x38, 0x4c, 0xfd, 0x36, 0xaa, 0x38, 0xe2, 0xdb, 0x28, 0xd9,
	0x60, 0x28, 0x0f, 0x37, 0x18, 0x92, 0x13, 0xd9, 0xca, 0x98, 0x13, 0xd9, 0xec, 0x57, 0x27, 0xd5,
	0x97, 0xff, 0xea, 0xa4, 0x36, 0xcb, 0x57, 0x27, 0x0f, 0xa0, 0x1a, 0xd9, 0xde, 0x94, 0xdf, 0x88,
	0x54, 0x28, 0xe9, 0x3e, 0x51, 0x8e, 0x27, 0x60, 0x8a, 0xe3, 0x09, 0xc3, 0x82, 0x26, 0x37, 0xf8,
	0xac, 0x4d, 0xea, 0x6f, 0x42, 0xd5, 0xe1, 0x8c, 0xa2, 0x4d, 0x3d, 0x2f, 0x5c, 0x9f, 0x01, 0x4d,
	0x89, 0x34, 0x7e, 0x00, 0x1b, 0xdc, 0x15, 0x95, 0xe3, 0x14, 0xe1, 0x8d, 0x6f, 0x43, 0xd9, 0x23,
	0xa8, 0x27, 0xd7, 0x83, 0x4d, 0xd1, 0xe7, 0x1e, 0xa2, 0x66, 0x37, 0x33, 0x39, 0xa1, 0xf1, 0x8f,
	0x9a, 0x74, 0xec, 0x2c, 0x3e, 0x99, 0x5a, 0x2d, 0x37, 0xb5, 0xe2, 0xce, 0x64, 0x61, 0xf8, 0xce,
	0xe4, 0x3a, 0x54, 0xba, 0xd8, 0x0e, 0x88, 0xbc, 0x4a, 0x25, 0x9e, 0xf4, 0x16, 0x54, 0x31, 0xba,
	0x0c, 0x2f, 0xe4, 0x67, 0x72, 0xa6, 0x7c, 0xd4, 0x6f, 0x41, 0x03, 0xa3, 0xc8, 0xb7, 0x1d, 0x64,
	0xd9, 0xbe, 0xdf, 0x2a, 0x33, 0x2c, 0x08, 0xd0, 0xbe, 0xcf, 0x0e, 0xbb, 0x39, 0x2d, 0xc3, 0xab,
	0xf7, 0x11, 0xeb, 0x1c, 0xbe, 0xef, 0xfb, 0xc6, 0x1e, 0xb4, 0xf2, 0xf6, 0x10, 0x96, 0x5f, 0x87,
	0x62, 0x2f, 0xee, 0x66, 0xef, 0x77, 0xf6, 0xe2, 0xae, 0xf1, 0xef, 0x05, 0xa8, 0x3c, 0x42, 0x97,
	0x9e, 0x83, 0xc6, 0x1f, 0x23, 0xba, 0x0c, 0x2f, 0x4f, 0xd9, 0xd3, 0x55, 0x8b, 0x81, 0x8f, 0x5d,
	0x76, 0xd0, 0xc1, 0x49, 0x72, 0x0d, 0x45, 0xe0, 0x08, 0x96, 0x8d, 0xf7, 0x40, 0x47, 0xcf, 0x09,
	0xc2, 0x81, 0xed, 0x5b, 0xa9, 0x48, 0x75, 0xfd, 0x5a, 0x92, 0xf8, 0x47, 0x52, 0xf4, 0xbb, 0xb0,
	0x92, 0xf0, 0xc8, 0x9b, 0xb8, 0x9e, 0xcb, 0x7c, 0x5d, 0xde, 0xa5, 0x58, 0x96, 0x04, 0xe2, 0x7a,
	0xec, 0xf1, 0xf0, 0xfd, 0xb2, 0xf2, 0xe8, 0x5c, 0xf4, 0x6b, 0xf9, 0xdc, 0xcb, 0xf8, 0x9f, 0x22,
	0xcc, 0x3f, 0x0d, 0x49, 0x72, 0x08, 0x3f, 0xc6, 0xea, 0x9b, 0x50, 0x26, 0x1e, 0xf1, 0x87, 0x2e,
	0x19, 0x32, 0x10, 0x5d, 0x60, 0x7b, 0x28, 0x8e, 0xed, 0x6e, 0xd6, 0xd4, 0x12, 0x28, 0x6e, 0xf2,
	0x63, 0x64, 0xbb, 0xc3, 0x37, 0xf9, 0x4d, 0x64, 0xb3, 0x3c, 0x87, 0x02, 0xe2, 0x91, 0xec, 0xb5,
	0x38, 0x01, 0xa3, 0xd3, 0xcd, 0x7f, 0x51, 0xcb, 0x55, 0x14, 0xad, 0x6a, 0x1c, 0x9c, 0xb3, 0x6e,
	0x75, 0xec, 0x65, 0x9e, 0x18, 0x05, 0xae, 0x15, 0x28, 0x63, 0x65, 0xab, 0x7d, 0x72, 0x99, 0x87,
	0xa2, 0x33, 0x96, 0x78, 0x00, 0xd5, 0x18, 0xa1, 0x60, 0xba, 0x5c, 0x56, 0xa1, 0xa4, 0xb9, 0xef,
	0x26, 0xeb, 0x2f, 0x3f, 0x8b, 0x30, 0x4b, 0xfa, 0xfc, 0x36, 0x34, 0xe2, 0x41, 0xe0, 0xc8, 0x53,
	0xd7, 0xc6, 0xa4, 0x74, 0x08, 0x94, 0x92, 0x43, 0x8c, 0x3f, 0xd7, 0x60, 0x85, 0x6f, 0x0a, 0xb9,
	0x87, 0x2b, 0xf7, 0xa1, 0xd3, 0x88, 0xd0, 0xa6, 0x09, 0xb2, 0xc2, 0x4c, 0x41, 0x56, 0x9c, 0x14,
	0x64, 0x86, 0x0f, 0x2b, 0x8f, 0x90, 0x8f, 0x5e, 0x42, 0xa9, 0xd1, 0x6f, 0x2b, 0x4c, 0x7c, 0xdb,
	0x6f, 0xc1, 0xc6, 0x63, 0x44, 0xd4, 0x99, 0x1f, 0x71, 0x9d, 0xa2, 0x30, 0xee, 0x3a, 0x85, 0x71,
	0x05, 0x6b, 0x43, 0xbc, 0x33, 0xad, 0x2d, 0xef, 0xc1, 0x82, 0xea, 0x92, 0x72, 0x85, 0x11, 0xad,
	0x73, 0x55, 0xb2, 0x99, 0x25, 0x34, 0x3e, 0x84, 0x4d, 0x9e, 0x5d, 0x47, 0xaa, 0x9e, 0xff, 0xe6,
	0x4f, 0x09, 0xc3, 0x42, 0x3e, 0x0c, 0x8d, 0xf7, 0xa4, 0x38, 0x13, 0x75, 0x10, 0x46, 0x81, 0x93,
	0xb9, 0x20, 0x3f, 0xe9, 0x6e, 0x51, 0x9e, 0xb3, 0x6d, 0xfb, 0x68, 0x1a, 0xce, 0xff, 0xd6, 0xa0,
	0x9e, 0xdc, 0xf3, 0xfa, 0x35, 0x7f, 0x50, 0x94, 0xe8, 0x59, 0x9e, 0x50, 0x9b, 0x57, 0xf2, 0xb5,
	0xf9, 0xf0, 0x37, 0x44, 0xd5, 0x69, 0xbe, 0x21, 0x7a, 0xf3, 0xfb, 0xd0, 0x50, 0xee, 0x73, 0xe8,
	0x0d, 0xa8, 0xf6, 0x83, 0x8b, 0x20, 0xbc, 0x0a, 0x96, 0xe6, 0xe8, 0x43, 0xc4, 0x3f, 0x74, 0x5b,
	0x5a, 0xd3, 0x6b, 0xbc, 0x5c, 0x5b, 0xba, 0xad, 0x2f, 0x28, 0x57, 0xe5, 0x96, 0xee, 0x50, 0x04,
	0x7d, 0xc5, 0x12, 0x3a, 0x78, 0xff, 0xd3, 0x2f, 0xb6, 0xb4, 0xcf, 0xbe, 0xd8, 0xd2, 0x3e, 0xff,
	0x62, 0x6b, 0xee, 0x3f, 0xbf, 0xd8, 0x9a, 0xfb, 0xc3, 0x17, 0x5b, 0x73, 0x3f, 0x7b, 0xb1, 0x35,
	0xf7, 0x77, 0x2f, 0xb6, 0xe6, 0x3e, 0x7d, 0xb1, 0x35, 0xf7, 0xd9, 0x8b, 0xad, 0xb9, 0xcf, 0x5f,
	0x6c, 0xcd, 0xfd, 0xf4, 0xe7, 0x5b, 0x73, 0x7f, 0xf5, 0xf3, 0xad, 0xb9, 0xdf, 0x5b, 0x61, 0xca,
	0x5d, 0x06, 0xbb, 0x76, 0xe4, 0xed, 0x46, 0x67, 0xec, 0x7f, 0x43, 0xf8, 0xdf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xee, 0x62, 0xef, 0xc5, 0x19, 0x41, 0x00, 0x00,
}
