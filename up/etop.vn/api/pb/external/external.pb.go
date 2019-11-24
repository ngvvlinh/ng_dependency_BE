// source: external/external.proto

package external

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	gender "etop.vn/api/pb/etop/etc/gender"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	shipping_provider "etop.vn/api/pb/etop/etc/shipping_provider"
	status3 "etop.vn/api/pb/etop/etc/status3"
	status4 "etop.vn/api/pb/etop/etc/status4"
	status5 "etop.vn/api/pb/etop/etc/status5"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	order "etop.vn/api/pb/etop/order"
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

type Partner struct {
	Id         dot.ID           `protobuf:"varint,1,opt,name=id" json:"id"`
	Name       string           `protobuf:"bytes,2,opt,name=name" json:"name"`
	PublicName string           `protobuf:"bytes,24,opt,name=public_name,json=publicName" json:"public_name"`
	Type       etop.AccountType `protobuf:"varint,3,opt,name=type,enum=etop.AccountType" json:"type"`
	Phone      string           `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	// only domain, no scheme
	Website              string   `protobuf:"bytes,5,opt,name=website" json:"website"`
	WebsiteUrl           string   `protobuf:"bytes,14,opt,name=website_url,json=websiteUrl" json:"website_url"`
	ImageUrl             string   `protobuf:"bytes,15,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email                string   `protobuf:"bytes,16,opt,name=email" json:"email"`
	RecognizedHosts      []string `protobuf:"bytes,17,rep,name=recognized_hosts,json=recognizedHosts" json:"recognized_hosts"`
	RedirectUrls         []string `protobuf:"bytes,18,rep,name=redirect_urls,json=redirectUrls" json:"redirect_urls"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Partner) Reset()         { *m = Partner{} }
func (m *Partner) String() string { return proto.CompactTextString(m) }
func (*Partner) ProtoMessage()    {}
func (*Partner) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{0}
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

func (m *Partner) GetType() etop.AccountType {
	if m != nil {
		return m.Type
	}
	return etop.AccountType_unknown
}

func (m *Partner) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Partner) GetWebsite() string {
	if m != nil {
		return m.Website
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

func (m *Partner) GetRecognizedHosts() []string {
	if m != nil {
		return m.RecognizedHosts
	}
	return nil
}

func (m *Partner) GetRedirectUrls() []string {
	if m != nil {
		return m.RedirectUrls
	}
	return nil
}

type CreateWebhookRequest struct {
	Entities             []string `protobuf:"bytes,2,rep,name=entities" json:"entities"`
	Fields               []string `protobuf:"bytes,3,rep,name=fields" json:"fields"`
	Url                  string   `protobuf:"bytes,4,opt,name=url" json:"url"`
	Metadata             string   `protobuf:"bytes,5,opt,name=metadata" json:"metadata"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateWebhookRequest) Reset()         { *m = CreateWebhookRequest{} }
func (m *CreateWebhookRequest) String() string { return proto.CompactTextString(m) }
func (*CreateWebhookRequest) ProtoMessage()    {}
func (*CreateWebhookRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{1}
}

var xxx_messageInfo_CreateWebhookRequest proto.InternalMessageInfo

func (m *CreateWebhookRequest) GetEntities() []string {
	if m != nil {
		return m.Entities
	}
	return nil
}

func (m *CreateWebhookRequest) GetFields() []string {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *CreateWebhookRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *CreateWebhookRequest) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

type DeleteWebhookRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteWebhookRequest) Reset()         { *m = DeleteWebhookRequest{} }
func (m *DeleteWebhookRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteWebhookRequest) ProtoMessage()    {}
func (*DeleteWebhookRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{2}
}

var xxx_messageInfo_DeleteWebhookRequest proto.InternalMessageInfo

func (m *DeleteWebhookRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

type WebhooksResponse struct {
	Webhooks             []*Webhook `protobuf:"bytes,1,rep,name=webhooks" json:"webhooks"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *WebhooksResponse) Reset()         { *m = WebhooksResponse{} }
func (m *WebhooksResponse) String() string { return proto.CompactTextString(m) }
func (*WebhooksResponse) ProtoMessage()    {}
func (*WebhooksResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{3}
}

var xxx_messageInfo_WebhooksResponse proto.InternalMessageInfo

func (m *WebhooksResponse) GetWebhooks() []*Webhook {
	if m != nil {
		return m.Webhooks
	}
	return nil
}

type Webhook struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Entities             []string       `protobuf:"bytes,2,rep,name=entities" json:"entities"`
	Fields               []string       `protobuf:"bytes,3,rep,name=fields" json:"fields"`
	Url                  string         `protobuf:"bytes,4,opt,name=url" json:"url"`
	Metadata             string         `protobuf:"bytes,5,opt,name=metadata" json:"metadata"`
	CreatedAt            dot.Time       `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	States               *WebhookStates `protobuf:"bytes,7,opt,name=states" json:"states"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Webhook) Reset()         { *m = Webhook{} }
func (m *Webhook) String() string { return proto.CompactTextString(m) }
func (*Webhook) ProtoMessage()    {}
func (*Webhook) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{4}
}

var xxx_messageInfo_Webhook proto.InternalMessageInfo

func (m *Webhook) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Webhook) GetEntities() []string {
	if m != nil {
		return m.Entities
	}
	return nil
}

func (m *Webhook) GetFields() []string {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *Webhook) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Webhook) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

func (m *Webhook) GetStates() *WebhookStates {
	if m != nil {
		return m.States
	}
	return nil
}

type WebhookStates struct {
	State                string        `protobuf:"bytes,1,opt,name=state" json:"state"`
	LastSentAt           dot.Time      `protobuf:"bytes,2,opt,name=last_sent_at,json=lastSentAt" json:"last_sent_at"`
	LastError            *WebhookError `protobuf:"bytes,3,opt,name=last_error,json=lastError" json:"last_error"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *WebhookStates) Reset()         { *m = WebhookStates{} }
func (m *WebhookStates) String() string { return proto.CompactTextString(m) }
func (*WebhookStates) ProtoMessage()    {}
func (*WebhookStates) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{5}
}

var xxx_messageInfo_WebhookStates proto.InternalMessageInfo

func (m *WebhookStates) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *WebhookStates) GetLastError() *WebhookError {
	if m != nil {
		return m.LastError
	}
	return nil
}

type WebhookError struct {
	Error                string   `protobuf:"bytes,1,opt,name=error" json:"error"`
	RespStatus           int32    `protobuf:"varint,2,opt,name=resp_status,json=respStatus" json:"resp_status"`
	RespBody             string   `protobuf:"bytes,3,opt,name=resp_body,json=respBody" json:"resp_body"`
	Retried              int32    `protobuf:"varint,4,opt,name=retried" json:"retried"`
	SentAt               dot.Time `protobuf:"bytes,5,opt,name=sent_at,json=sentAt" json:"sent_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WebhookError) Reset()         { *m = WebhookError{} }
func (m *WebhookError) String() string { return proto.CompactTextString(m) }
func (*WebhookError) ProtoMessage()    {}
func (*WebhookError) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{6}
}

var xxx_messageInfo_WebhookError proto.InternalMessageInfo

func (m *WebhookError) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *WebhookError) GetRespStatus() int32 {
	if m != nil {
		return m.RespStatus
	}
	return 0
}

func (m *WebhookError) GetRespBody() string {
	if m != nil {
		return m.RespBody
	}
	return ""
}

func (m *WebhookError) GetRetried() int32 {
	if m != nil {
		return m.Retried
	}
	return 0
}

type GetChangesRequest struct {
	Paging               *common.ForwardPaging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Entity               *string               `protobuf:"bytes,2,opt,name=entity" json:"entity"`
	EntityId             *string               `protobuf:"bytes,3,opt,name=entity_id,json=entityId" json:"entity_id"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *GetChangesRequest) Reset()         { *m = GetChangesRequest{} }
func (m *GetChangesRequest) String() string { return proto.CompactTextString(m) }
func (*GetChangesRequest) ProtoMessage()    {}
func (*GetChangesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{7}
}

var xxx_messageInfo_GetChangesRequest proto.InternalMessageInfo

func (m *GetChangesRequest) GetPaging() *common.ForwardPaging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetChangesRequest) GetEntity() string {
	if m != nil && m.Entity != nil {
		return *m.Entity
	}
	return ""
}

func (m *GetChangesRequest) GetEntityId() string {
	if m != nil && m.EntityId != nil {
		return *m.EntityId
	}
	return ""
}

type Callback struct {
	Id                   dot.ID    `protobuf:"varint,1,opt,name=id" json:"id"`
	Changes              []*Change `protobuf:"bytes,2,rep,name=changes" json:"changes"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Callback) Reset()         { *m = Callback{} }
func (m *Callback) String() string { return proto.CompactTextString(m) }
func (*Callback) ProtoMessage()    {}
func (*Callback) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{8}
}

var xxx_messageInfo_Callback proto.InternalMessageInfo

func (m *Callback) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Callback) GetChanges() []*Change {
	if m != nil {
		return m.Changes
	}
	return nil
}

// ChangesData serialize changes data for storing in MongoDB
type ChangesData struct {
	// for using with mongodb
	XId                  dot.ID    `protobuf:"varint,1,opt,name=_id,json=Id" json:"_id"`
	WebhookId            dot.ID    `protobuf:"varint,3,opt,name=webhook_id,json=webhookId" json:"webhook_id"`
	AccountId            dot.ID    `protobuf:"varint,4,opt,name=account_id,json=accountId" json:"account_id"`
	CreatedAt            dot.Time  `protobuf:"bytes,5,opt,name=created_at,json=createdAt" json:"created_at"`
	Changes              []*Change `protobuf:"bytes,6,rep,name=changes" json:"changes"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ChangesData) Reset()         { *m = ChangesData{} }
func (m *ChangesData) String() string { return proto.CompactTextString(m) }
func (*ChangesData) ProtoMessage()    {}
func (*ChangesData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{9}
}

var xxx_messageInfo_ChangesData proto.InternalMessageInfo

func (m *ChangesData) GetXId() dot.ID {
	if m != nil {
		return m.XId
	}
	return 0
}

func (m *ChangesData) GetWebhookId() dot.ID {
	if m != nil {
		return m.WebhookId
	}
	return 0
}

func (m *ChangesData) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *ChangesData) GetChanges() []*Change {
	if m != nil {
		return m.Changes
	}
	return nil
}

type Change struct {
	Time                 dot.Time     `protobuf:"bytes,2,opt,name=time" json:"time"`
	ChangeType           string       `protobuf:"bytes,3,opt,name=change_type,json=changeType" json:"change_type"`
	Entity               string       `protobuf:"bytes,4,opt,name=entity" json:"entity"`
	Latest               *LatestOneOf `protobuf:"bytes,5,opt,name=latest" json:"latest"`
	Changed              *ChangeOneOf `protobuf:"bytes,6,opt,name=changed" json:"changed"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Change) Reset()         { *m = Change{} }
func (m *Change) String() string { return proto.CompactTextString(m) }
func (*Change) ProtoMessage()    {}
func (*Change) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{10}
}

var xxx_messageInfo_Change proto.InternalMessageInfo

func (m *Change) GetChangeType() string {
	if m != nil {
		return m.ChangeType
	}
	return ""
}

func (m *Change) GetEntity() string {
	if m != nil {
		return m.Entity
	}
	return ""
}

func (m *Change) GetLatest() *LatestOneOf {
	if m != nil {
		return m.Latest
	}
	return nil
}

func (m *Change) GetChanged() *ChangeOneOf {
	if m != nil {
		return m.Changed
	}
	return nil
}

type LatestOneOf struct {
	Order                *Order       `protobuf:"bytes,7,opt,name=order" json:"order"`
	Fulfillment          *Fulfillment `protobuf:"bytes,8,opt,name=fulfillment" json:"fulfillment"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *LatestOneOf) Reset()         { *m = LatestOneOf{} }
func (m *LatestOneOf) String() string { return proto.CompactTextString(m) }
func (*LatestOneOf) ProtoMessage()    {}
func (*LatestOneOf) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{11}
}

var xxx_messageInfo_LatestOneOf proto.InternalMessageInfo

func (m *LatestOneOf) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *LatestOneOf) GetFulfillment() *Fulfillment {
	if m != nil {
		return m.Fulfillment
	}
	return nil
}

type ChangeOneOf struct {
	Order                *Order       `protobuf:"bytes,5,opt,name=order" json:"order"`
	Fulfillment          *Fulfillment `protobuf:"bytes,6,opt,name=fulfillment" json:"fulfillment"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ChangeOneOf) Reset()         { *m = ChangeOneOf{} }
func (m *ChangeOneOf) String() string { return proto.CompactTextString(m) }
func (*ChangeOneOf) ProtoMessage()    {}
func (*ChangeOneOf) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{12}
}

var xxx_messageInfo_ChangeOneOf proto.InternalMessageInfo

func (m *ChangeOneOf) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *ChangeOneOf) GetFulfillment() *Fulfillment {
	if m != nil {
		return m.Fulfillment
	}
	return nil
}

type CancelOrderRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	ExternalId           string   `protobuf:"bytes,3,opt,name=external_id,json=externalId" json:"external_id"`
	CancelReason         string   `protobuf:"bytes,4,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelOrderRequest) Reset()         { *m = CancelOrderRequest{} }
func (m *CancelOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CancelOrderRequest) ProtoMessage()    {}
func (*CancelOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{13}
}

var xxx_messageInfo_CancelOrderRequest proto.InternalMessageInfo

func (m *CancelOrderRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelOrderRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CancelOrderRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *CancelOrderRequest) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

type OrderIDRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	ExternalId           string   `protobuf:"bytes,3,opt,name=external_id,json=externalId" json:"external_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderIDRequest) Reset()         { *m = OrderIDRequest{} }
func (m *OrderIDRequest) String() string { return proto.CompactTextString(m) }
func (*OrderIDRequest) ProtoMessage()    {}
func (*OrderIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{14}
}

var xxx_messageInfo_OrderIDRequest proto.InternalMessageInfo

func (m *OrderIDRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *OrderIDRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *OrderIDRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

type FulfillmentIDRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShippingCode         string   `protobuf:"bytes,2,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FulfillmentIDRequest) Reset()         { *m = FulfillmentIDRequest{} }
func (m *FulfillmentIDRequest) String() string { return proto.CompactTextString(m) }
func (*FulfillmentIDRequest) ProtoMessage()    {}
func (*FulfillmentIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{15}
}

var xxx_messageInfo_FulfillmentIDRequest proto.InternalMessageInfo

func (m *FulfillmentIDRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *FulfillmentIDRequest) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

type OrderAndFulfillments struct {
	Order                *Order          `protobuf:"bytes,1,opt,name=order" json:"order"`
	Fulfillments         []*Fulfillment  `protobuf:"bytes,2,rep,name=fulfillments" json:"fulfillments"`
	FulfillmentErrors    []*common.Error `protobuf:"bytes,3,rep,name=fulfillment_errors,json=fulfillmentErrors" json:"fulfillment_errors"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *OrderAndFulfillments) Reset()         { *m = OrderAndFulfillments{} }
func (m *OrderAndFulfillments) String() string { return proto.CompactTextString(m) }
func (*OrderAndFulfillments) ProtoMessage()    {}
func (*OrderAndFulfillments) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{16}
}

var xxx_messageInfo_OrderAndFulfillments proto.InternalMessageInfo

func (m *OrderAndFulfillments) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *OrderAndFulfillments) GetFulfillments() []*Fulfillment {
	if m != nil {
		return m.Fulfillments
	}
	return nil
}

func (m *OrderAndFulfillments) GetFulfillmentErrors() []*common.Error {
	if m != nil {
		return m.FulfillmentErrors
	}
	return nil
}

type Order struct {
	Id                        dot.ID                `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId                    dot.ID                `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Code                      *string               `protobuf:"bytes,3,opt,name=code" json:"code"`
	ExternalId                *string               `protobuf:"bytes,4,opt,name=external_id,json=externalId" json:"external_id"`
	ExternalCode              *string               `protobuf:"bytes,5,opt,name=external_code,json=externalCode" json:"external_code"`
	ExternalUrl               *string               `protobuf:"bytes,6,opt,name=external_url,json=externalUrl" json:"external_url"`
	SelfUrl                   *string               `protobuf:"bytes,7,opt,name=self_url,json=selfUrl" json:"self_url"`
	CustomerAddress           *OrderAddress         `protobuf:"bytes,10,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	ShippingAddress           *OrderAddress         `protobuf:"bytes,12,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	CreatedAt                 dot.Time              `protobuf:"bytes,13,opt,name=created_at,json=createdAt" json:"created_at"`
	ProcessedAt               dot.Time              `protobuf:"bytes,14,opt,name=processed_at,json=processedAt" json:"processed_at"`
	UpdatedAt                 dot.Time              `protobuf:"bytes,15,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt                  dot.Time              `protobuf:"bytes,16,opt,name=closed_at,json=closedAt" json:"closed_at"`
	ConfirmedAt               dot.Time              `protobuf:"bytes,17,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelledAt               dot.Time              `protobuf:"bytes,18,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CancelReason              *string               `protobuf:"bytes,19,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	ConfirmStatus             *status3.Status       `protobuf:"varint,70,opt,name=confirm_status,json=confirmStatus,enum=status3.Status" json:"confirm_status"`
	Status                    *status5.Status       `protobuf:"varint,34,opt,name=status,enum=status5.Status" json:"status"`
	FulfillmentShippingStatus *status5.Status       `protobuf:"varint,36,opt,name=fulfillment_shipping_status,json=fulfillmentShippingStatus,enum=status5.Status" json:"fulfillment_shipping_status"`
	EtopPaymentStatus         *status4.Status       `protobuf:"varint,69,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	Lines                     []*OrderLine          `protobuf:"bytes,41,rep,name=lines" json:"lines"`
	TotalItems                *int32                `protobuf:"varint,43,opt,name=total_items,json=totalItems" json:"total_items"`
	BasketValue               *int32                `protobuf:"varint,44,opt,name=basket_value,json=basketValue" json:"basket_value"`
	OrderDiscount             *int32                `protobuf:"varint,46,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	TotalDiscount             *int32                `protobuf:"varint,47,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalFee                  *int32                `protobuf:"varint,49,opt,name=total_fee,json=totalFee" json:"total_fee"`
	FeeLines                  []*order.OrderFeeLine `protobuf:"bytes,50,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	TotalAmount               *int32                `protobuf:"varint,51,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	OrderNote                 *string               `protobuf:"bytes,59,opt,name=order_note,json=orderNote" json:"order_note"`
	Shipping                  *OrderShipping        `protobuf:"bytes,60,opt,name=shipping" json:"shipping"`
	XXX_NoUnkeyedLiteral      struct{}              `json:"-"`
	XXX_sizecache             int32                 `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{17}
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Order) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Order) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *Order) GetExternalId() string {
	if m != nil && m.ExternalId != nil {
		return *m.ExternalId
	}
	return ""
}

func (m *Order) GetExternalCode() string {
	if m != nil && m.ExternalCode != nil {
		return *m.ExternalCode
	}
	return ""
}

func (m *Order) GetExternalUrl() string {
	if m != nil && m.ExternalUrl != nil {
		return *m.ExternalUrl
	}
	return ""
}

func (m *Order) GetSelfUrl() string {
	if m != nil && m.SelfUrl != nil {
		return *m.SelfUrl
	}
	return ""
}

func (m *Order) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *Order) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *Order) GetCancelReason() string {
	if m != nil && m.CancelReason != nil {
		return *m.CancelReason
	}
	return ""
}

func (m *Order) GetConfirmStatus() status3.Status {
	if m != nil && m.ConfirmStatus != nil {
		return *m.ConfirmStatus
	}
	return status3.Status_Z
}

func (m *Order) GetStatus() status5.Status {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return status5.Status_Z
}

func (m *Order) GetFulfillmentShippingStatus() status5.Status {
	if m != nil && m.FulfillmentShippingStatus != nil {
		return *m.FulfillmentShippingStatus
	}
	return status5.Status_Z
}

func (m *Order) GetEtopPaymentStatus() status4.Status {
	if m != nil && m.EtopPaymentStatus != nil {
		return *m.EtopPaymentStatus
	}
	return status4.Status_Z
}

func (m *Order) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *Order) GetTotalItems() int32 {
	if m != nil && m.TotalItems != nil {
		return *m.TotalItems
	}
	return 0
}

func (m *Order) GetBasketValue() int32 {
	if m != nil && m.BasketValue != nil {
		return *m.BasketValue
	}
	return 0
}

func (m *Order) GetOrderDiscount() int32 {
	if m != nil && m.OrderDiscount != nil {
		return *m.OrderDiscount
	}
	return 0
}

func (m *Order) GetTotalDiscount() int32 {
	if m != nil && m.TotalDiscount != nil {
		return *m.TotalDiscount
	}
	return 0
}

func (m *Order) GetTotalFee() int32 {
	if m != nil && m.TotalFee != nil {
		return *m.TotalFee
	}
	return 0
}

func (m *Order) GetFeeLines() []*order.OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *Order) GetTotalAmount() int32 {
	if m != nil && m.TotalAmount != nil {
		return *m.TotalAmount
	}
	return 0
}

func (m *Order) GetOrderNote() string {
	if m != nil && m.OrderNote != nil {
		return *m.OrderNote
	}
	return ""
}

func (m *Order) GetShipping() *OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

type OrderShipping struct {
	PickupAddress        *OrderAddress                       `protobuf:"bytes,1,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ReturnAddress        *OrderAddress                       `protobuf:"bytes,17,opt,name=return_address,json=returnAddress" json:"return_address"`
	ShippingServiceName  *string                             `protobuf:"bytes,16,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceCode  *string                             `protobuf:"bytes,2,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee   *int32                              `protobuf:"varint,3,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	Carrier              *shipping_provider.ShippingProvider `protobuf:"varint,5,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	IncludeInsurance     *bool                               `protobuf:"varint,6,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	TryOn                *try_on.TryOnCode                   `protobuf:"varint,7,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	ShippingNote         *string                             `protobuf:"bytes,8,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	CodAmount            *int32                              `protobuf:"varint,9,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	GrossWeight          *int32                              `protobuf:"varint,10,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	Length               *int32                              `protobuf:"varint,11,opt,name=length" json:"length"`
	Width                *int32                              `protobuf:"varint,12,opt,name=width" json:"width"`
	Height               *int32                              `protobuf:"varint,13,opt,name=height" json:"height"`
	ChargeableWeight     *int32                              `protobuf:"varint,15,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *OrderShipping) Reset()         { *m = OrderShipping{} }
func (m *OrderShipping) String() string { return proto.CompactTextString(m) }
func (*OrderShipping) ProtoMessage()    {}
func (*OrderShipping) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{18}
}

var xxx_messageInfo_OrderShipping proto.InternalMessageInfo

func (m *OrderShipping) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *OrderShipping) GetReturnAddress() *OrderAddress {
	if m != nil {
		return m.ReturnAddress
	}
	return nil
}

func (m *OrderShipping) GetShippingServiceName() string {
	if m != nil && m.ShippingServiceName != nil {
		return *m.ShippingServiceName
	}
	return ""
}

func (m *OrderShipping) GetShippingServiceCode() string {
	if m != nil && m.ShippingServiceCode != nil {
		return *m.ShippingServiceCode
	}
	return ""
}

func (m *OrderShipping) GetShippingServiceFee() int32 {
	if m != nil && m.ShippingServiceFee != nil {
		return *m.ShippingServiceFee
	}
	return 0
}

func (m *OrderShipping) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil && m.Carrier != nil {
		return *m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *OrderShipping) GetIncludeInsurance() bool {
	if m != nil && m.IncludeInsurance != nil {
		return *m.IncludeInsurance
	}
	return false
}

func (m *OrderShipping) GetTryOn() try_on.TryOnCode {
	if m != nil && m.TryOn != nil {
		return *m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *OrderShipping) GetShippingNote() string {
	if m != nil && m.ShippingNote != nil {
		return *m.ShippingNote
	}
	return ""
}

func (m *OrderShipping) GetCodAmount() int32 {
	if m != nil && m.CodAmount != nil {
		return *m.CodAmount
	}
	return 0
}

func (m *OrderShipping) GetGrossWeight() int32 {
	if m != nil && m.GrossWeight != nil {
		return *m.GrossWeight
	}
	return 0
}

func (m *OrderShipping) GetLength() int32 {
	if m != nil && m.Length != nil {
		return *m.Length
	}
	return 0
}

func (m *OrderShipping) GetWidth() int32 {
	if m != nil && m.Width != nil {
		return *m.Width
	}
	return 0
}

func (m *OrderShipping) GetHeight() int32 {
	if m != nil && m.Height != nil {
		return *m.Height
	}
	return 0
}

func (m *OrderShipping) GetChargeableWeight() int32 {
	if m != nil && m.ChargeableWeight != nil {
		return *m.ChargeableWeight
	}
	return 0
}

type CreateOrderRequest struct {
	ExternalId      string        `protobuf:"bytes,4,opt,name=external_id,json=externalId" json:"external_id"`
	ExternalCode    string        `protobuf:"bytes,5,opt,name=external_code,json=externalCode" json:"external_code"`
	ExternalUrl     string        `protobuf:"bytes,6,opt,name=external_url,json=externalUrl" json:"external_url"`
	CustomerAddress *OrderAddress `protobuf:"bytes,10,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	ShippingAddress *OrderAddress `protobuf:"bytes,12,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	Lines           []*OrderLine  `protobuf:"bytes,41,rep,name=lines" json:"lines"`
	TotalItems      int32         `protobuf:"varint,43,opt,name=total_items,json=totalItems" json:"total_items"`
	// basket_value = SUM(lines.retail_price)
	BasketValue          int32                 `protobuf:"varint,44,opt,name=basket_value,json=basketValue" json:"basket_value"`
	OrderDiscount        int32                 `protobuf:"varint,46,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	TotalDiscount        int32                 `protobuf:"varint,47,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalFee             *int32                `protobuf:"varint,49,opt,name=total_fee,json=totalFee" json:"total_fee"`
	FeeLines             []*order.OrderFeeLine `protobuf:"bytes,48,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	TotalAmount          int32                 `protobuf:"varint,50,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	OrderNote            string                `protobuf:"bytes,51,opt,name=order_note,json=orderNote" json:"order_note"`
	Shipping             *OrderShipping        `protobuf:"bytes,60,opt,name=shipping" json:"shipping"`
	ExternalMeta         map[string]string     `protobuf:"bytes,62,rep,name=external_meta,json=externalMeta" json:"external_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *CreateOrderRequest) Reset()         { *m = CreateOrderRequest{} }
func (m *CreateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrderRequest) ProtoMessage()    {}
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{19}
}

var xxx_messageInfo_CreateOrderRequest proto.InternalMessageInfo

func (m *CreateOrderRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *CreateOrderRequest) GetExternalCode() string {
	if m != nil {
		return m.ExternalCode
	}
	return ""
}

func (m *CreateOrderRequest) GetExternalUrl() string {
	if m != nil {
		return m.ExternalUrl
	}
	return ""
}

func (m *CreateOrderRequest) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *CreateOrderRequest) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *CreateOrderRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *CreateOrderRequest) GetOrderDiscount() int32 {
	if m != nil {
		return m.OrderDiscount
	}
	return 0
}

func (m *CreateOrderRequest) GetTotalDiscount() int32 {
	if m != nil {
		return m.TotalDiscount
	}
	return 0
}

func (m *CreateOrderRequest) GetTotalFee() int32 {
	if m != nil && m.TotalFee != nil {
		return *m.TotalFee
	}
	return 0
}

func (m *CreateOrderRequest) GetFeeLines() []*order.OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *CreateOrderRequest) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *CreateOrderRequest) GetOrderNote() string {
	if m != nil {
		return m.OrderNote
	}
	return ""
}

func (m *CreateOrderRequest) GetShipping() *OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

func (m *CreateOrderRequest) GetExternalMeta() map[string]string {
	if m != nil {
		return m.ExternalMeta
	}
	return nil
}

type OrderLine struct {
	VariantId   dot.ID `protobuf:"varint,3,opt,name=variant_id,json=variantId" json:"variant_id"`
	ProductId   dot.ID `protobuf:"varint,23,opt,name=product_id,json=productId" json:"product_id"`
	ProductName string `protobuf:"bytes,7,opt,name=product_name,json=productName" json:"product_name"`
	Quantity    int32  `protobuf:"varint,17,opt,name=quantity" json:"quantity"`
	ListPrice   int32  `protobuf:"varint,18,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice int32  `protobuf:"varint,19,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	// payment_price = retail_price - discount_per_item
	PaymentPrice         *int32             `protobuf:"varint,20,opt,name=payment_price,json=paymentPrice" json:"payment_price"`
	ImageUrl             string             `protobuf:"bytes,21,opt,name=image_url,json=imageUrl" json:"image_url"`
	Attributes           []*order.Attribute `protobuf:"bytes,22,rep,name=attributes" json:"attributes"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *OrderLine) Reset()         { *m = OrderLine{} }
func (m *OrderLine) String() string { return proto.CompactTextString(m) }
func (*OrderLine) ProtoMessage()    {}
func (*OrderLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{20}
}

var xxx_messageInfo_OrderLine proto.InternalMessageInfo

func (m *OrderLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *OrderLine) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *OrderLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *OrderLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *OrderLine) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *OrderLine) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *OrderLine) GetPaymentPrice() int32 {
	if m != nil && m.PaymentPrice != nil {
		return *m.PaymentPrice
	}
	return 0
}

func (m *OrderLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *OrderLine) GetAttributes() []*order.Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type Fulfillment struct {
	Id                       dot.ID                              `protobuf:"varint,1,opt,name=id" json:"id"`
	OrderId                  dot.ID                              `protobuf:"varint,2,opt,name=order_id,json=orderId" json:"order_id"`
	ShopId                   dot.ID                              `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	SelfUrl                  *string                             `protobuf:"bytes,4,opt,name=self_url,json=selfUrl" json:"self_url"`
	TotalItems               *int32                              `protobuf:"varint,6,opt,name=total_items,json=totalItems" json:"total_items"`
	BasketValue              *int32                              `protobuf:"varint,10,opt,name=basket_value,json=basketValue" json:"basket_value"`
	CreatedAt                dot.Time                            `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                dot.Time                            `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt                 dot.Time                            `protobuf:"bytes,14,opt,name=closed_at,json=closedAt" json:"closed_at"`
	CancelledAt              dot.Time                            `protobuf:"bytes,16,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CancelReason             *string                             `protobuf:"bytes,17,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	Carrier                  *shipping_provider.ShippingProvider `protobuf:"varint,18,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	ShippingServiceName      *string                             `protobuf:"bytes,27,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceFee       *int32                              `protobuf:"varint,28,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ActualShippingServiceFee *int32                              `protobuf:"varint,33,opt,name=actual_shipping_service_fee,json=actualShippingServiceFee" json:"actual_shipping_service_fee"`
	ShippingServiceCode      *string                             `protobuf:"bytes,29,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingCode             *string                             `protobuf:"bytes,22,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	ShippingNote             *string                             `protobuf:"bytes,19,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	TryOn                    *try_on.TryOnCode                   `protobuf:"varint,20,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	IncludeInsurance         *bool                               `protobuf:"varint,21,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	ConfirmStatus            *status3.Status                     `protobuf:"varint,24,opt,name=confirm_status,json=confirmStatus,enum=status3.Status" json:"confirm_status"`
	ShippingState            *shipping.State                     `protobuf:"varint,25,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	ShippingStatus           *status5.Status                     `protobuf:"varint,50,opt,name=shipping_status,json=shippingStatus,enum=status5.Status" json:"shipping_status"`
	Status                   *status5.Status                     `protobuf:"varint,26,opt,name=status,enum=status5.Status" json:"status"`
	CodAmount                *int32                              `protobuf:"varint,9,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	ActualCodAmount          *int32                              `protobuf:"varint,40,opt,name=actual_cod_amount,json=actualCodAmount" json:"actual_cod_amount"`
	ChargeableWeight         *int32                              `protobuf:"varint,23,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	PickupAddress            *OrderAddress                       `protobuf:"bytes,30,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ReturnAddress            *OrderAddress                       `protobuf:"bytes,31,opt,name=return_address,json=returnAddress" json:"return_address"`
	ShippingAddress          *OrderAddress                       `protobuf:"bytes,32,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	EtopPaymentStatus        *status4.Status                     `protobuf:"varint,69,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	EstimatedDeliveryAt      dot.Time                            `protobuf:"bytes,74,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	EstimatedPickupAt        dot.Time                            `protobuf:"bytes,75,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	XXX_NoUnkeyedLiteral     struct{}                            `json:"-"`
	XXX_sizecache            int32                               `json:"-"`
}

func (m *Fulfillment) Reset()         { *m = Fulfillment{} }
func (m *Fulfillment) String() string { return proto.CompactTextString(m) }
func (*Fulfillment) ProtoMessage()    {}
func (*Fulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{21}
}

var xxx_messageInfo_Fulfillment proto.InternalMessageInfo

func (m *Fulfillment) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Fulfillment) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *Fulfillment) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Fulfillment) GetSelfUrl() string {
	if m != nil && m.SelfUrl != nil {
		return *m.SelfUrl
	}
	return ""
}

func (m *Fulfillment) GetTotalItems() int32 {
	if m != nil && m.TotalItems != nil {
		return *m.TotalItems
	}
	return 0
}

func (m *Fulfillment) GetBasketValue() int32 {
	if m != nil && m.BasketValue != nil {
		return *m.BasketValue
	}
	return 0
}

func (m *Fulfillment) GetCancelReason() string {
	if m != nil && m.CancelReason != nil {
		return *m.CancelReason
	}
	return ""
}

func (m *Fulfillment) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil && m.Carrier != nil {
		return *m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *Fulfillment) GetShippingServiceName() string {
	if m != nil && m.ShippingServiceName != nil {
		return *m.ShippingServiceName
	}
	return ""
}

func (m *Fulfillment) GetShippingServiceFee() int32 {
	if m != nil && m.ShippingServiceFee != nil {
		return *m.ShippingServiceFee
	}
	return 0
}

func (m *Fulfillment) GetActualShippingServiceFee() int32 {
	if m != nil && m.ActualShippingServiceFee != nil {
		return *m.ActualShippingServiceFee
	}
	return 0
}

func (m *Fulfillment) GetShippingServiceCode() string {
	if m != nil && m.ShippingServiceCode != nil {
		return *m.ShippingServiceCode
	}
	return ""
}

func (m *Fulfillment) GetShippingCode() string {
	if m != nil && m.ShippingCode != nil {
		return *m.ShippingCode
	}
	return ""
}

func (m *Fulfillment) GetShippingNote() string {
	if m != nil && m.ShippingNote != nil {
		return *m.ShippingNote
	}
	return ""
}

func (m *Fulfillment) GetTryOn() try_on.TryOnCode {
	if m != nil && m.TryOn != nil {
		return *m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *Fulfillment) GetIncludeInsurance() bool {
	if m != nil && m.IncludeInsurance != nil {
		return *m.IncludeInsurance
	}
	return false
}

func (m *Fulfillment) GetConfirmStatus() status3.Status {
	if m != nil && m.ConfirmStatus != nil {
		return *m.ConfirmStatus
	}
	return status3.Status_Z
}

func (m *Fulfillment) GetShippingState() shipping.State {
	if m != nil && m.ShippingState != nil {
		return *m.ShippingState
	}
	return shipping.State_default
}

func (m *Fulfillment) GetShippingStatus() status5.Status {
	if m != nil && m.ShippingStatus != nil {
		return *m.ShippingStatus
	}
	return status5.Status_Z
}

func (m *Fulfillment) GetStatus() status5.Status {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return status5.Status_Z
}

func (m *Fulfillment) GetCodAmount() int32 {
	if m != nil && m.CodAmount != nil {
		return *m.CodAmount
	}
	return 0
}

func (m *Fulfillment) GetActualCodAmount() int32 {
	if m != nil && m.ActualCodAmount != nil {
		return *m.ActualCodAmount
	}
	return 0
}

func (m *Fulfillment) GetChargeableWeight() int32 {
	if m != nil && m.ChargeableWeight != nil {
		return *m.ChargeableWeight
	}
	return 0
}

func (m *Fulfillment) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *Fulfillment) GetReturnAddress() *OrderAddress {
	if m != nil {
		return m.ReturnAddress
	}
	return nil
}

func (m *Fulfillment) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *Fulfillment) GetEtopPaymentStatus() status4.Status {
	if m != nil && m.EtopPaymentStatus != nil {
		return *m.EtopPaymentStatus
	}
	return status4.Status_Z
}

type Ward struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ward) Reset()         { *m = Ward{} }
func (m *Ward) String() string { return proto.CompactTextString(m) }
func (*Ward) ProtoMessage()    {}
func (*Ward) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{22}
}

var xxx_messageInfo_Ward proto.InternalMessageInfo

func (m *Ward) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type District struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Wards                []Ward   `protobuf:"bytes,2,rep,name=wards" json:"wards"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *District) Reset()         { *m = District{} }
func (m *District) String() string { return proto.CompactTextString(m) }
func (*District) ProtoMessage()    {}
func (*District) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{23}
}

var xxx_messageInfo_District proto.InternalMessageInfo

func (m *District) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *District) GetWards() []Ward {
	if m != nil {
		return m.Wards
	}
	return nil
}

type Province struct {
	Name                 string     `protobuf:"bytes,1,opt,name=name" json:"name"`
	Districts            []District `protobuf:"bytes,2,rep,name=districts" json:"districts"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Province) Reset()         { *m = Province{} }
func (m *Province) String() string { return proto.CompactTextString(m) }
func (*Province) ProtoMessage()    {}
func (*Province) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{24}
}

var xxx_messageInfo_Province proto.InternalMessageInfo

func (m *Province) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Province) GetDistricts() []District {
	if m != nil {
		return m.Districts
	}
	return nil
}

type LocationResponse struct {
	Provinces            []Province `protobuf:"bytes,1,rep,name=provinces" json:"provinces"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *LocationResponse) Reset()         { *m = LocationResponse{} }
func (m *LocationResponse) String() string { return proto.CompactTextString(m) }
func (*LocationResponse) ProtoMessage()    {}
func (*LocationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{25}
}

var xxx_messageInfo_LocationResponse proto.InternalMessageInfo

func (m *LocationResponse) GetProvinces() []Province {
	if m != nil {
		return m.Provinces
	}
	return nil
}

type GetShippingServicesRequest struct {
	PickupAddress   *LocationAddress `protobuf:"bytes,1,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ShippingAddress *LocationAddress `protobuf:"bytes,2,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	// in gram (g)
	GrossWeight int32 `protobuf:"varint,12,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	// in gram (g)
	ChargeableWeight int32 `protobuf:"varint,5,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	// in centimetre (cm)
	Length int32 `protobuf:"varint,6,opt,name=length" json:"length"`
	// in centimetre (cm)
	Width int32 `protobuf:"varint,7,opt,name=width" json:"width"`
	// in centimetre (cm)
	Height               int32    `protobuf:"varint,8,opt,name=height" json:"height"`
	BasketValue          int32    `protobuf:"varint,10,opt,name=basket_value,json=basketValue" json:"basket_value"`
	CodAmount            int32    `protobuf:"varint,11,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	IncludeInsurance     *bool    `protobuf:"varint,13,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetShippingServicesRequest) Reset()         { *m = GetShippingServicesRequest{} }
func (m *GetShippingServicesRequest) String() string { return proto.CompactTextString(m) }
func (*GetShippingServicesRequest) ProtoMessage()    {}
func (*GetShippingServicesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{26}
}

var xxx_messageInfo_GetShippingServicesRequest proto.InternalMessageInfo

func (m *GetShippingServicesRequest) GetPickupAddress() *LocationAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *GetShippingServicesRequest) GetShippingAddress() *LocationAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *GetShippingServicesRequest) GetGrossWeight() int32 {
	if m != nil {
		return m.GrossWeight
	}
	return 0
}

func (m *GetShippingServicesRequest) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *GetShippingServicesRequest) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *GetShippingServicesRequest) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *GetShippingServicesRequest) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *GetShippingServicesRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *GetShippingServicesRequest) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *GetShippingServicesRequest) GetIncludeInsurance() bool {
	if m != nil && m.IncludeInsurance != nil {
		return *m.IncludeInsurance
	}
	return false
}

type LocationAddress struct {
	Province             string   `protobuf:"bytes,1,opt,name=province" json:"province"`
	District             string   `protobuf:"bytes,2,opt,name=district" json:"district"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LocationAddress) Reset()         { *m = LocationAddress{} }
func (m *LocationAddress) String() string { return proto.CompactTextString(m) }
func (*LocationAddress) ProtoMessage()    {}
func (*LocationAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{27}
}

var xxx_messageInfo_LocationAddress proto.InternalMessageInfo

func (m *LocationAddress) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *LocationAddress) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

type GetShippingServicesResponse struct {
	Services             []*ShippingService `protobuf:"bytes,1,rep,name=services" json:"services"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetShippingServicesResponse) Reset()         { *m = GetShippingServicesResponse{} }
func (m *GetShippingServicesResponse) String() string { return proto.CompactTextString(m) }
func (*GetShippingServicesResponse) ProtoMessage()    {}
func (*GetShippingServicesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{28}
}

var xxx_messageInfo_GetShippingServicesResponse proto.InternalMessageInfo

func (m *GetShippingServicesResponse) GetServices() []*ShippingService {
	if m != nil {
		return m.Services
	}
	return nil
}

type ShippingService struct {
	Code                 string                             `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string                             `protobuf:"bytes,2,opt,name=name" json:"name"`
	Fee                  int32                              `protobuf:"varint,3,opt,name=fee" json:"fee"`
	Carrier              shipping_provider.ShippingProvider `protobuf:"varint,4,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	EstimatedPickupAt    dot.Time                           `protobuf:"bytes,5,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	EstimatedDeliveryAt  dot.Time                           `protobuf:"bytes,6,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *ShippingService) Reset()         { *m = ShippingService{} }
func (m *ShippingService) String() string { return proto.CompactTextString(m) }
func (*ShippingService) ProtoMessage()    {}
func (*ShippingService) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{29}
}

var xxx_messageInfo_ShippingService proto.InternalMessageInfo

func (m *ShippingService) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ShippingService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShippingService) GetFee() int32 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func (m *ShippingService) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

type OrderCustomer struct {
	FullName             string        `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Email                string        `protobuf:"bytes,2,opt,name=email" json:"email"`
	Phone                string        `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Gender               gender.Gender `protobuf:"varint,5,opt,name=gender,enum=gender.Gender" json:"gender"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *OrderCustomer) Reset()         { *m = OrderCustomer{} }
func (m *OrderCustomer) String() string { return proto.CompactTextString(m) }
func (*OrderCustomer) ProtoMessage()    {}
func (*OrderCustomer) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{30}
}

var xxx_messageInfo_OrderCustomer proto.InternalMessageInfo

func (m *OrderCustomer) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *OrderCustomer) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OrderCustomer) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *OrderCustomer) GetGender() gender.Gender {
	if m != nil {
		return m.Gender
	}
	return gender.Gender_unknown
}

type OrderAddress struct {
	FullName             string   `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	Province             string   `protobuf:"bytes,7,opt,name=province" json:"province"`
	District             string   `protobuf:"bytes,8,opt,name=district" json:"district"`
	Ward                 string   `protobuf:"bytes,9,opt,name=ward" json:"ward"`
	Company              string   `protobuf:"bytes,11,opt,name=company" json:"company"`
	Address1             string   `protobuf:"bytes,12,opt,name=address1" json:"address1"`
	Address2             string   `protobuf:"bytes,13,opt,name=address2" json:"address2"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderAddress) Reset()         { *m = OrderAddress{} }
func (m *OrderAddress) String() string { return proto.CompactTextString(m) }
func (*OrderAddress) ProtoMessage()    {}
func (*OrderAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_c7af707b80d96213, []int{31}
}

var xxx_messageInfo_OrderAddress proto.InternalMessageInfo

func (m *OrderAddress) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *OrderAddress) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *OrderAddress) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OrderAddress) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *OrderAddress) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *OrderAddress) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *OrderAddress) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *OrderAddress) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *OrderAddress) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func init() {
	proto.RegisterType((*Partner)(nil), "external.Partner")
	proto.RegisterType((*CreateWebhookRequest)(nil), "external.CreateWebhookRequest")
	proto.RegisterType((*DeleteWebhookRequest)(nil), "external.DeleteWebhookRequest")
	proto.RegisterType((*WebhooksResponse)(nil), "external.WebhooksResponse")
	proto.RegisterType((*Webhook)(nil), "external.Webhook")
	proto.RegisterType((*WebhookStates)(nil), "external.WebhookStates")
	proto.RegisterType((*WebhookError)(nil), "external.WebhookError")
	proto.RegisterType((*GetChangesRequest)(nil), "external.GetChangesRequest")
	proto.RegisterType((*Callback)(nil), "external.Callback")
	proto.RegisterType((*ChangesData)(nil), "external.ChangesData")
	proto.RegisterType((*Change)(nil), "external.Change")
	proto.RegisterType((*LatestOneOf)(nil), "external.LatestOneOf")
	proto.RegisterType((*ChangeOneOf)(nil), "external.ChangeOneOf")
	proto.RegisterType((*CancelOrderRequest)(nil), "external.CancelOrderRequest")
	proto.RegisterType((*OrderIDRequest)(nil), "external.OrderIDRequest")
	proto.RegisterType((*FulfillmentIDRequest)(nil), "external.FulfillmentIDRequest")
	proto.RegisterType((*OrderAndFulfillments)(nil), "external.OrderAndFulfillments")
	proto.RegisterType((*Order)(nil), "external.Order")
	proto.RegisterType((*OrderShipping)(nil), "external.OrderShipping")
	proto.RegisterType((*CreateOrderRequest)(nil), "external.CreateOrderRequest")
	proto.RegisterMapType((map[string]string)(nil), "external.CreateOrderRequest.ExternalMetaEntry")
	proto.RegisterType((*OrderLine)(nil), "external.OrderLine")
	proto.RegisterType((*Fulfillment)(nil), "external.Fulfillment")
	proto.RegisterType((*Ward)(nil), "external.Ward")
	proto.RegisterType((*District)(nil), "external.District")
	proto.RegisterType((*Province)(nil), "external.Province")
	proto.RegisterType((*LocationResponse)(nil), "external.LocationResponse")
	proto.RegisterType((*GetShippingServicesRequest)(nil), "external.GetShippingServicesRequest")
	proto.RegisterType((*LocationAddress)(nil), "external.LocationAddress")
	proto.RegisterType((*GetShippingServicesResponse)(nil), "external.GetShippingServicesResponse")
	proto.RegisterType((*ShippingService)(nil), "external.ShippingService")
	proto.RegisterType((*OrderCustomer)(nil), "external.OrderCustomer")
	proto.RegisterType((*OrderAddress)(nil), "external.OrderAddress")
}

func init() { proto.RegisterFile("external/external.proto", fileDescriptor_c7af707b80d96213) }

var fileDescriptor_c7af707b80d96213 = []byte{
	// 2944 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x5a, 0x4d, 0x70, 0x1c, 0x47,
	0xf5, 0xd7, 0x4a, 0xda, 0xaf, 0xb7, 0x1f, 0xd2, 0xb6, 0x3e, 0x3c, 0x96, 0xed, 0xb5, 0x3c, 0xfa,
	0xbb, 0x2c, 0xc5, 0xc9, 0xca, 0x91, 0x63, 0x27, 0xf9, 0x13, 0x13, 0x14, 0xc9, 0x0e, 0x0a, 0x26,
	0x56, 0x8d, 0x1c, 0x5c, 0x45, 0x51, 0xb5, 0x35, 0x9a, 0x69, 0xad, 0xa6, 0x3c, 0x3b, 0xb3, 0x99,
	0xe9, 0xb5, 0x11, 0x67, 0x0e, 0x1c, 0x39, 0x41, 0x8e, 0x1c, 0xa9, 0xe2, 0xc6, 0x89, 0x13, 0xe7,
	0x70, 0x4b, 0x41, 0x71, 0x84, 0x8a, 0x9d, 0x0b, 0x47, 0x6e, 0x14, 0x27, 0xa8, 0x7e, 0xdd, 0x3d,
	0xdb, 0xb3, 0x3b, 0xab, 0x95, 0xec, 0x82, 0xe2, 0x22, 0xed, 0xfc, 0xde, 0xef, 0x75, 0xbf, 0xee,
	0x7e, 0xfd, 0xfa, 0xbd, 0x9e, 0x81, 0x0b, 0xf4, 0xc7, 0x8c, 0x46, 0x81, 0xed, 0x6f, 0xaa, 0x1f,
	0xad, 0x5e, 0x14, 0xb2, 0x90, 0x94, 0xd4, 0xf3, 0xca, 0x62, 0x27, 0xec, 0x84, 0x08, 0x6e, 0xf2,
	0x5f, 0x42, 0xbe, 0xb2, 0xe0, 0x84, 0xdd, 0x6e, 0x18, 0x6c, 0x8a, 0x7f, 0x12, 0xbc, 0x4c, 0x59,
	0xd8, 0xdb, 0xa4, 0xcc, 0xd9, 0xec, 0xd0, 0xc0, 0xa5, 0x91, 0xfc, 0x27, 0xa5, 0x57, 0x12, 0x69,
	0x7c, 0xec, 0xf5, 0x7a, 0x5e, 0xd0, 0xd9, 0x8c, 0x99, 0xcd, 0xa8, 0x14, 0x6f, 0x8c, 0x88, 0xdb,
	0xbd, 0x28, 0x7c, 0xe6, 0xf1, 0x76, 0xd4, 0x8f, 0xd1, 0x96, 0x98, 0xcd, 0xfa, 0xf1, 0x6d, 0xf9,
	0x7f, 0x8c, 0xf8, 0x9d, 0xd3, 0xc5, 0x77, 0xd2, 0xe2, 0x95, 0x44, 0xcc, 0xa2, 0x93, 0x36, 0x0e,
	0xd1, 0x55, 0x36, 0xce, 0x49, 0x59, 0xd8, 0x93, 0xc0, 0x32, 0x02, 0x61, 0xc4, 0x8d, 0xc4, 0xbf,
	0x12, 0xbf, 0xda, 0x09, 0xc3, 0x8e, 0x4f, 0x37, 0xf1, 0xe9, 0xb0, 0x7f, 0xb4, 0xc9, 0xbc, 0x2e,
	0x8d, 0x99, 0xdd, 0x95, 0x8a, 0xe6, 0xbf, 0xa6, 0xa1, 0xb8, 0x6f, 0x47, 0x2c, 0xa0, 0x11, 0x59,
	0x84, 0x69, 0xcf, 0x35, 0x72, 0xab, 0xb9, 0xf5, 0x99, 0x8f, 0x66, 0xbf, 0xfc, 0xeb, 0xd5, 0x29,
	0x6b, 0xda, 0x73, 0x89, 0x01, 0xb3, 0x81, 0xdd, 0xa5, 0xc6, 0xf4, 0x6a, 0x6e, 0xbd, 0x2c, 0x71,
	0x44, 0xc8, 0x75, 0xa8, 0xf4, 0xfa, 0x87, 0xbe, 0xe7, 0xb4, 0x91, 0x60, 0x68, 0x04, 0x10, 0x82,
	0x4f, 0x39, 0xed, 0x26, 0xcc, 0xb2, 0x93, 0x1e, 0x35, 0x66, 0x56, 0x73, 0xeb, 0xf5, 0xad, 0x46,
	0x0b, 0xcd, 0xde, 0x76, 0x9c, 0xb0, 0x1f, 0xb0, 0xc7, 0x27, 0x3d, 0xaa, 0xda, 0xe4, 0x24, 0xb2,
	0x02, 0xf9, 0xde, 0x71, 0x18, 0x50, 0xa3, 0xac, 0xb5, 0x26, 0x20, 0xd2, 0x84, 0xe2, 0x73, 0x7a,
	0x18, 0x7b, 0x8c, 0x1a, 0x79, 0x4d, 0xaa, 0x40, 0x6e, 0x8f, 0xfc, 0xd9, 0xee, 0x47, 0xbe, 0x51,
	0xd7, 0xed, 0x91, 0x82, 0xcf, 0x22, 0x9f, 0x5c, 0x83, 0xb2, 0xd7, 0xb5, 0x3b, 0x82, 0x34, 0xa7,
	0x91, 0x4a, 0x08, 0x73, 0xca, 0x0a, 0xe4, 0x69, 0xd7, 0xf6, 0x7c, 0x63, 0x5e, 0xb7, 0x02, 0x21,
	0xb2, 0x01, 0xf3, 0x11, 0x75, 0xc2, 0x4e, 0xe0, 0xfd, 0x84, 0xba, 0xed, 0xe3, 0x30, 0x66, 0xb1,
	0xd1, 0x58, 0x9d, 0x59, 0x2f, 0x5b, 0x73, 0x03, 0xfc, 0xbb, 0x1c, 0x26, 0x6b, 0x50, 0x8b, 0xa8,
	0xeb, 0x45, 0xd4, 0x61, 0xbc, 0xb3, 0xd8, 0x20, 0xc8, 0xab, 0x2a, 0xf0, 0xb3, 0xc8, 0x8f, 0xcd,
	0x9f, 0xe6, 0x60, 0x71, 0x27, 0xa2, 0x36, 0xa3, 0x4f, 0xe8, 0xe1, 0x71, 0x18, 0x3e, 0xb5, 0xe8,
	0xe7, 0x7d, 0x1a, 0x33, 0xb2, 0x02, 0x25, 0x1a, 0x30, 0x8f, 0x79, 0x34, 0x36, 0xa6, 0x51, 0x31,
	0x79, 0x26, 0xcb, 0x50, 0x38, 0xf2, 0xa8, 0xef, 0xc6, 0xc6, 0x0c, 0x4a, 0xe4, 0x13, 0x59, 0x86,
	0x19, 0x3e, 0xaa, 0x59, 0xcd, 0x6c, 0x0e, 0x90, 0x55, 0x28, 0x75, 0x29, 0xb3, 0x5d, 0x9b, 0xd9,
	0xa9, 0xb9, 0x4b, 0x50, 0xf3, 0x4d, 0x58, 0xdc, 0xa5, 0x3e, 0x1d, 0xb1, 0x22, 0xd3, 0x29, 0xcc,
	0x6d, 0x98, 0x97, 0xbc, 0xd8, 0xa2, 0x71, 0x2f, 0x0c, 0x62, 0x4a, 0xde, 0x82, 0xd2, 0x73, 0x89,
	0x19, 0xb9, 0xd5, 0x99, 0xf5, 0x0a, 0x5f, 0x6b, 0xb5, 0x9b, 0x55, 0xab, 0x09, 0xc5, 0xfc, 0x67,
	0x0e, 0x8a, 0x12, 0x1d, 0xe3, 0x79, 0xff, 0xd5, 0x09, 0x20, 0xef, 0x03, 0x38, 0xb8, 0x0c, 0x6e,
	0xdb, 0x66, 0x46, 0x61, 0x35, 0xb7, 0x5e, 0xd9, 0x5a, 0x69, 0x89, 0xfd, 0xd3, 0x52, 0xfb, 0xa7,
	0xf5, 0x58, 0xed, 0x1f, 0xab, 0x2c, 0xd9, 0xdb, 0x8c, 0x6c, 0x42, 0x01, 0x23, 0x48, 0x6c, 0x14,
	0x51, 0xed, 0xc2, 0xc8, 0xb8, 0x0f, 0x50, 0x6c, 0x49, 0x9a, 0xf9, 0xab, 0x1c, 0xd4, 0x52, 0x12,
	0xee, 0x71, 0x28, 0xc3, 0x49, 0x48, 0x3c, 0x0e, 0x21, 0xf2, 0x01, 0x54, 0x7d, 0x3b, 0x66, 0xed,
	0x98, 0x06, 0x8c, 0xdb, 0x36, 0x3d, 0xd1, 0x36, 0xe0, 0xfc, 0x03, 0x1a, 0xb0, 0x6d, 0x46, 0xee,
	0x00, 0x3e, 0xb5, 0x69, 0x14, 0x85, 0x11, 0x6e, 0xc2, 0xca, 0xd6, 0xf2, 0x88, 0x81, 0xf7, 0xb9,
	0xd4, 0x2a, 0x73, 0x26, 0xfe, 0x34, 0xff, 0x90, 0x83, 0xaa, 0x2e, 0xc3, 0x3d, 0x81, 0x4d, 0xa4,
	0x2c, 0x44, 0x88, 0xef, 0xbc, 0x88, 0xc6, 0xbd, 0xb6, 0x08, 0x60, 0x68, 0x60, 0x5e, 0xed, 0x3c,
	0x2e, 0x38, 0x40, 0x9c, 0xef, 0x3c, 0xa4, 0x1d, 0x86, 0xee, 0x09, 0x5a, 0x92, 0xac, 0x02, 0x87,
	0x3f, 0x0a, 0xdd, 0x13, 0xbe, 0xc7, 0x23, 0xca, 0x22, 0x8f, 0xba, 0xb8, 0x86, 0xaa, 0x15, 0x05,
	0x92, 0xdb, 0x50, 0x54, 0xd3, 0x90, 0x9f, 0x38, 0x0d, 0x85, 0x18, 0xa7, 0xc0, 0x8c, 0xa1, 0xf1,
	0x31, 0x65, 0x3b, 0xc7, 0x76, 0xd0, 0xa1, 0xb1, 0x72, 0xec, 0x0d, 0x28, 0xf4, 0xec, 0x8e, 0x17,
	0x74, 0x70, 0x40, 0xdc, 0x59, 0x9d, 0x6e, 0xeb, 0x41, 0x18, 0x3d, 0xb7, 0x23, 0x77, 0x1f, 0x05,
	0x96, 0x24, 0x70, 0x67, 0x43, 0xc7, 0x3b, 0x11, 0x41, 0xd0, 0x92, 0x4f, 0xe4, 0x12, 0x94, 0xc5,
	0xaf, 0xb6, 0xe7, 0x8a, 0xf1, 0x48, 0x0f, 0x3d, 0xd9, 0x73, 0xcd, 0x87, 0x50, 0xda, 0xb1, 0x7d,
	0xff, 0xd0, 0x76, 0xc6, 0xf9, 0xf7, 0x1b, 0x50, 0x74, 0x84, 0x4d, 0xe8, 0xde, 0x95, 0xad, 0xf9,
	0xc1, 0xb2, 0x08, 0x63, 0x2d, 0x45, 0x30, 0xff, 0x9c, 0x83, 0x8a, 0x1c, 0xc0, 0x2e, 0xf7, 0xd6,
	0x25, 0x98, 0x69, 0x0f, 0x37, 0xb9, 0xe7, 0x92, 0x35, 0x00, 0xb9, 0xc1, 0x94, 0x49, 0x4a, 0x5a,
	0x96, 0xb8, 0x20, 0xd9, 0x22, 0xfc, 0x72, 0xd2, 0xac, 0x4e, 0x92, 0xf8, 0x9e, 0x3b, 0xb4, 0x1d,
	0xf2, 0xe7, 0xd9, 0x0e, 0xda, 0xb8, 0x0a, 0x93, 0xc6, 0xf5, 0x97, 0x1c, 0x14, 0x04, 0x46, 0x5a,
	0x30, 0xcb, 0x4f, 0xa7, 0x33, 0xb8, 0x37, 0xf2, 0xb8, 0xd3, 0x89, 0x56, 0xda, 0xc9, 0xf1, 0x92,
	0x84, 0x7b, 0x21, 0xe0, 0xa7, 0x0b, 0xb9, 0x9c, 0x2c, 0x9e, 0x1e, 0x14, 0xd4, 0x12, 0xbe, 0x05,
	0x05, 0x9f, 0x6f, 0x40, 0x35, 0xc4, 0xa5, 0x81, 0xa9, 0x0f, 0x11, 0x7f, 0x14, 0xd0, 0x47, 0x47,
	0x96, 0x24, 0x91, 0x4d, 0x35, 0x34, 0x57, 0x46, 0x88, 0xa5, 0xe1, 0xa1, 0x09, 0xbe, 0x62, 0x99,
	0x5d, 0xa8, 0x68, 0xed, 0x90, 0xeb, 0x90, 0xc7, 0xe3, 0x59, 0x06, 0x8a, 0xb9, 0x81, 0xf6, 0x23,
	0x0e, 0x5b, 0x42, 0x4a, 0xde, 0x85, 0xca, 0x51, 0xdf, 0x3f, 0xf2, 0x7c, 0xbf, 0x4b, 0x03, 0x66,
	0x94, 0x86, 0xbb, 0x7a, 0x30, 0x10, 0x5a, 0x3a, 0x93, 0x77, 0xa7, 0x99, 0x31, 0xe8, 0x2e, 0x7f,
	0x9e, 0xee, 0x0a, 0x67, 0xee, 0xee, 0x17, 0x39, 0x20, 0x3b, 0x76, 0xe0, 0x50, 0x5f, 0xb4, 0x77,
	0xda, 0x99, 0xc1, 0x13, 0x09, 0x9e, 0xc2, 0xa4, 0x13, 0x09, 0x8e, 0xf0, 0x95, 0x54, 0x7d, 0x25,
	0x3b, 0x49, 0xad, 0xa4, 0x12, 0xec, 0xb9, 0x64, 0x03, 0x6a, 0x0e, 0x76, 0xd6, 0x8e, 0xa8, 0x1d,
	0x87, 0x41, 0x6a, 0x41, 0xab, 0x42, 0x64, 0xa1, 0xc4, 0xec, 0x40, 0x1d, 0x2d, 0xda, 0xdb, 0xfd,
	0xcf, 0xda, 0x64, 0x3e, 0x81, 0x45, 0x6d, 0x76, 0x26, 0x75, 0xb7, 0x01, 0xb5, 0x24, 0xa9, 0x1c,
	0xe9, 0xb7, 0xaa, 0x44, 0x3b, 0xa1, 0x4b, 0xcd, 0xdf, 0xe6, 0x60, 0x11, 0x87, 0xb0, 0x1d, 0xb8,
	0x5a, 0x0f, 0xf1, 0x60, 0x4d, 0x73, 0xa7, 0xae, 0xe9, 0xfb, 0x50, 0xd5, 0x56, 0x4a, 0x45, 0x98,
	0x31, 0x8b, 0x9a, 0xa2, 0x92, 0xf7, 0x80, 0x68, 0xcf, 0xe2, 0xe0, 0x10, 0xe7, 0x6c, 0x65, 0xab,
	0xcc, 0xa3, 0xa4, 0x38, 0x2c, 0x1a, 0x1a, 0x09, 0x91, 0xd8, 0x7c, 0x01, 0x90, 0x47, 0x2b, 0xc6,
	0x8c, 0xff, 0x0a, 0x14, 0xe3, 0xe3, 0xb0, 0xc7, 0x27, 0x74, 0x5a, 0x13, 0x15, 0x38, 0xb8, 0xe7,
	0x12, 0x22, 0x57, 0x43, 0x84, 0x52, 0xb1, 0x0e, 0x57, 0xd3, 0xeb, 0x80, 0x4b, 0x9e, 0xf2, 0x8a,
	0x35, 0xa8, 0x25, 0x04, 0xd4, 0xc6, 0xe3, 0xdd, 0xaa, 0x2a, 0x90, 0xcf, 0x26, 0xb9, 0x06, 0xc9,
	0x33, 0xa6, 0x7d, 0x05, 0xe4, 0x24, 0x2d, 0xf3, 0x9c, 0xef, 0x22, 0x94, 0x62, 0xea, 0x1f, 0xa1,
	0xb8, 0x88, 0xe2, 0x22, 0x7f, 0xe6, 0xa2, 0x6d, 0x98, 0x77, 0xfa, 0x31, 0x0b, 0xbb, 0x34, 0x6a,
	0xdb, 0xae, 0x1b, 0xd1, 0x38, 0x36, 0x60, 0xf8, 0x20, 0x15, 0x8b, 0x25, 0xa4, 0xd6, 0x9c, 0xe2,
	0x4b, 0x80, 0x37, 0x91, 0xac, 0xbc, 0x6a, 0xa2, 0x7a, 0x7a, 0x13, 0x8a, 0xaf, 0x9a, 0x48, 0x47,
	0xe4, 0xda, 0x79, 0x22, 0xf2, 0x3d, 0xa8, 0xf6, 0xa2, 0xd0, 0xa1, 0x71, 0x2c, 0x94, 0xeb, 0x13,
	0x95, 0x2b, 0x09, 0x7f, 0x9b, 0xf1, 0x9e, 0xfb, 0x3d, 0x57, 0xf5, 0x3c, 0x37, 0xb9, 0x67, 0xc9,
	0xde, 0x66, 0xe4, 0x5d, 0x28, 0x3b, 0x7e, 0x28, 0xbb, 0x9d, 0x9f, 0xa8, 0x59, 0x12, 0x64, 0x61,
	0xb2, 0x13, 0x06, 0x47, 0x5e, 0xd4, 0x15, 0xba, 0x8d, 0xc9, 0x26, 0x27, 0x7c, 0xa9, 0x8e, 0x01,
	0xc1, 0x17, 0xea, 0xe4, 0x0c, 0xea, 0x8a, 0xbf, 0xcd, 0xb8, 0x53, 0xa5, 0x43, 0xcd, 0x82, 0x70,
	0x2a, 0x3d, 0xc8, 0x90, 0xbb, 0x50, 0x97, 0x5d, 0xaa, 0xc4, 0xe7, 0x01, 0x96, 0x38, 0x73, 0x2d,
	0x59, 0x0e, 0xb6, 0x44, 0xde, 0x63, 0xd5, 0x24, 0x4d, 0xa6, 0x41, 0x37, 0x44, 0xba, 0xd8, 0x8f,
	0x0d, 0x33, 0xc5, 0xbf, 0xa3, 0xf8, 0x52, 0x4c, 0x1e, 0xc1, 0x25, 0x7d, 0x23, 0x26, 0x0e, 0x24,
	0xb5, 0xff, 0x2f, 0x5b, 0xfb, 0xa2, 0xa6, 0x73, 0x20, 0x55, 0x64, 0xcf, 0x1f, 0xc2, 0x02, 0xaf,
	0xbe, 0xda, 0x3d, 0xfb, 0x44, 0xb4, 0x28, 0x1a, 0xba, 0x9f, 0x6a, 0xe8, 0x1d, 0xd5, 0x50, 0x83,
	0x73, 0xf7, 0x05, 0x55, 0x36, 0xb0, 0x01, 0x79, 0xdf, 0x0b, 0x68, 0x6c, 0x6c, 0x60, 0x34, 0x58,
	0x18, 0xf2, 0xdd, 0x87, 0x5e, 0x40, 0x2d, 0xc1, 0xe0, 0x1b, 0x97, 0x85, 0x8c, 0xef, 0x5a, 0x46,
	0xbb, 0xb1, 0x71, 0x93, 0x67, 0x73, 0x16, 0x20, 0xb4, 0xc7, 0x11, 0xbe, 0x27, 0x0f, 0xed, 0xf8,
	0x29, 0x65, 0xed, 0x67, 0xb6, 0xdf, 0xa7, 0xc6, 0x9b, 0xc8, 0xa8, 0x08, 0xec, 0x07, 0x1c, 0x22,
	0xd7, 0xa1, 0x8e, 0xd1, 0xac, 0xed, 0x7a, 0x31, 0x26, 0x26, 0x46, 0x0b, 0x49, 0x35, 0x44, 0x77,
	0x25, 0xc8, 0x69, 0xa2, 0xab, 0x84, 0xb6, 0x29, 0x68, 0x88, 0x26, 0xb4, 0x4b, 0x50, 0x16, 0xb4,
	0x23, 0x4a, 0x8d, 0xb7, 0x91, 0x51, 0x42, 0xe0, 0x01, 0xa5, 0xe4, 0x16, 0x94, 0x8f, 0x28, 0x6d,
	0x8b, 0xd1, 0x6d, 0xc9, 0xd1, 0x89, 0x52, 0x1a, 0x87, 0xf6, 0x80, 0x52, 0x1c, 0x5d, 0xe9, 0x48,
	0xfc, 0x40, 0xfb, 0x45, 0x73, 0x76, 0x17, 0xfb, 0xbc, 0x2d, 0xec, 0x47, 0x6c, 0x1b, 0x21, 0x72,
	0x05, 0x40, 0xd8, 0x1f, 0x84, 0x8c, 0x1a, 0xdf, 0x42, 0x1f, 0x2a, 0x23, 0xf2, 0x69, 0xc8, 0x28,
	0xb9, 0x0d, 0x25, 0xb5, 0xa6, 0xc6, 0x07, 0xc3, 0x95, 0x03, 0xf6, 0xaa, 0xd6, 0xcf, 0x4a, 0x88,
	0xe6, 0x2f, 0xf3, 0x50, 0x4b, 0xc9, 0xc8, 0x3d, 0xa8, 0xf7, 0x3c, 0xe7, 0x69, 0xbf, 0x97, 0x44,
	0x96, 0xdc, 0xa9, 0x91, 0xa5, 0x26, 0xd8, 0x2a, 0xae, 0xdc, 0x83, 0x7a, 0x44, 0x59, 0x3f, 0x0a,
	0x12, 0xf5, 0xc6, 0xe9, 0xea, 0x82, 0xad, 0xd4, 0xb7, 0x60, 0x69, 0xe0, 0x98, 0x34, 0x7a, 0xe6,
	0x39, 0x54, 0xdc, 0x07, 0x60, 0xed, 0x6c, 0x2d, 0x28, 0xe1, 0x81, 0x90, 0xe1, 0x95, 0x40, 0x96,
	0xce, 0xe0, 0x3c, 0x1c, 0xd1, 0xc1, 0x10, 0x7e, 0x0b, 0x16, 0x47, 0x74, 0xf8, 0x42, 0xce, 0xe0,
	0xb4, 0x93, 0x21, 0x15, 0xbe, 0xa4, 0xf7, 0xa0, 0xe8, 0xd8, 0x51, 0xe4, 0xc9, 0xfc, 0xa7, 0xbe,
	0xb5, 0xd6, 0x1a, 0xb9, 0xd2, 0x69, 0xa9, 0x59, 0xdc, 0x97, 0x80, 0xa5, 0x74, 0xc8, 0x4d, 0x68,
	0x78, 0x81, 0xe3, 0xf7, 0x5d, 0xda, 0xf6, 0x82, 0xb8, 0x1f, 0xf1, 0xad, 0x8f, 0x07, 0x47, 0xc9,
	0x9a, 0x97, 0x82, 0x3d, 0x85, 0x93, 0x75, 0x28, 0x88, 0x6b, 0x1a, 0x3c, 0x3b, 0xea, 0x5b, 0x8d,
	0x96, 0x78, 0x6c, 0x3d, 0x8e, 0x4e, 0x1e, 0x05, 0x7c, 0x00, 0x56, 0x9e, 0xf1, 0x9f, 0x3c, 0xb4,
	0x24, 0x56, 0xa0, 0x5b, 0x94, 0x44, 0x68, 0x51, 0x20, 0x7a, 0xc6, 0x15, 0x00, 0x27, 0x74, 0x95,
	0x67, 0x95, 0x71, 0x88, 0x65, 0x27, 0x74, 0xa5, 0x5f, 0x5d, 0x83, 0x6a, 0x27, 0x0a, 0xe3, 0xb8,
	0xfd, 0x9c, 0x7a, 0x9d, 0x63, 0x86, 0x87, 0x51, 0xde, 0xaa, 0x20, 0xf6, 0x04, 0x21, 0x5e, 0xb3,
	0xf8, 0x34, 0xe8, 0xb0, 0x63, 0xa3, 0x82, 0x42, 0xf9, 0x44, 0x16, 0x21, 0xff, 0xdc, 0x73, 0xd9,
	0x31, 0x9e, 0x3e, 0x79, 0x4b, 0x3c, 0x70, 0xf6, 0xb1, 0x68, 0xaa, 0x26, 0xd8, 0xe2, 0x89, 0xcf,
	0x81, 0x73, 0x6c, 0x47, 0x1d, 0x6a, 0x1f, 0xfa, 0x54, 0xf5, 0x36, 0x87, 0x94, 0xf9, 0x81, 0x40,
	0x74, 0x69, 0xfe, 0xb1, 0x00, 0x44, 0xdc, 0x64, 0xa4, 0xb2, 0xc1, 0xeb, 0x19, 0x27, 0x78, 0x76,
	0x76, 0x97, 0x71, 0x8e, 0xab, 0xdc, 0x28, 0x75, 0x9a, 0xdf, 0xc8, 0x3a, 0xcd, 0x25, 0x33, 0x75,
	0xa6, 0xff, 0x6f, 0x1c, 0xdc, 0xe7, 0x08, 0x9a, 0xd7, 0x33, 0x82, 0xa6, 0x9a, 0x2b, 0x2d, 0x74,
	0xde, 0xc8, 0x0a, 0x9d, 0x6a, 0x02, 0xf4, 0x00, 0x7a, 0x33, 0x3b, 0x80, 0x4a, 0xea, 0x50, 0x18,
	0xbd, 0x99, 0x1d, 0x46, 0x15, 0xf9, 0x55, 0x83, 0xe9, 0xad, 0xb3, 0x04, 0xd3, 0x1b, 0x43, 0xc1,
	0x74, 0x4b, 0x1f, 0x91, 0x1e, 0x52, 0xd7, 0x52, 0x21, 0xf5, 0xb6, 0xb6, 0xf2, 0xaf, 0x19, 0x58,
	0xc9, 0x81, 0xe6, 0x80, 0x5d, 0xca, 0x6c, 0xe3, 0xdb, 0x68, 0x78, 0x4b, 0xab, 0xf0, 0x46, 0x9c,
	0xbb, 0x75, 0x5f, 0xca, 0xbe, 0x4f, 0x99, 0x7d, 0x3f, 0x60, 0xd1, 0xc9, 0xc0, 0x55, 0x39, 0xb4,
	0xf2, 0x21, 0x34, 0x46, 0x28, 0x64, 0x1e, 0x66, 0x9e, 0xd2, 0x13, 0x71, 0x91, 0x62, 0xf1, 0x9f,
	0x7c, 0x57, 0x8a, 0x95, 0x14, 0x01, 0x50, 0x3c, 0xfc, 0xff, 0xf4, 0x7b, 0x39, 0xf3, 0x6f, 0xd3,
	0x50, 0x4e, 0xdc, 0x84, 0x8f, 0xfe, 0x99, 0x1d, 0x79, 0xb6, 0x28, 0xdd, 0x53, 0xf5, 0xbd, 0xc4,
	0x45, 0x7d, 0xdf, 0x8b, 0x42, 0xb7, 0xef, 0x20, 0xe9, 0x82, 0x4e, 0x92, 0xf8, 0x9e, 0xcb, 0x27,
	0x5c, 0x91, 0x30, 0x5a, 0x17, 0xf5, 0x3d, 0x24, 0x25, 0x18, 0xab, 0x57, 0xa1, 0xf4, 0x79, 0xdf,
	0x16, 0x15, 0x74, 0x43, 0x5b, 0x95, 0x04, 0xe5, 0xfd, 0xf9, 0x5e, 0xcc, 0xda, 0xbd, 0xc8, 0x73,
	0x28, 0x66, 0x5a, 0x8a, 0x53, 0xe6, 0xf8, 0x3e, 0x87, 0x79, 0x7f, 0x11, 0x65, 0xb6, 0xe7, 0x4b,
	0xda, 0x82, 0xbe, 0xc0, 0x42, 0x22, 0x88, 0x6b, 0x50, 0x53, 0xe9, 0x89, 0x60, 0x2e, 0xa2, 0x73,
	0x55, 0x25, 0x28, 0x48, 0xa9, 0x3b, 0xdc, 0xa5, 0xcc, 0x3b, 0xdc, 0x5b, 0x00, 0x36, 0x63, 0x91,
	0x77, 0xd8, 0x67, 0x34, 0x36, 0x96, 0xe5, 0x45, 0x84, 0x70, 0xc2, 0x6d, 0x25, 0xb0, 0x34, 0x8e,
	0xf9, 0x8f, 0x2a, 0x54, 0xb4, 0xaa, 0x68, 0x4c, 0x0d, 0x73, 0x15, 0x4a, 0xc2, 0x01, 0x87, 0x8a,
	0x98, 0x22, 0xa2, 0x7b, 0xa9, 0x22, 0x67, 0x26, 0xa3, 0xc8, 0xd1, 0xeb, 0x8c, 0xd9, 0x74, 0x9d,
	0x31, 0x94, 0x32, 0x15, 0x26, 0xa6, 0x4c, 0x30, 0x9a, 0x32, 0xa5, 0xab, 0x84, 0xca, 0x79, 0xaa,
	0x84, 0x74, 0x9a, 0x5f, 0x7b, 0xe5, 0x34, 0xbf, 0x7e, 0xce, 0x34, 0x5f, 0xcf, 0xd3, 0xe7, 0x5f,
	0x33, 0x4f, 0x6f, 0x64, 0xe4, 0xe9, 0x5a, 0x1e, 0x40, 0x5e, 0x21, 0x0f, 0x18, 0x9b, 0xe0, 0x5c,
	0x1a, 0x9f, 0xe0, 0x8c, 0x4b, 0x56, 0x2e, 0x9f, 0x92, 0xac, 0x5c, 0xb2, 0x1d, 0xd6, 0xb7, 0xfd,
	0x76, 0xa6, 0xe2, 0x35, 0x54, 0x34, 0x04, 0xe5, 0x60, 0x54, 0x7d, 0x6c, 0x46, 0x75, 0x65, 0x7c,
	0x46, 0xb5, 0x36, 0x7c, 0x1b, 0xb1, 0x9c, 0xce, 0x44, 0x46, 0x48, 0x18, 0x72, 0x17, 0x32, 0xd2,
	0x95, 0x41, 0xf6, 0xb3, 0x38, 0x21, 0xfb, 0xc9, 0x4c, 0xaa, 0x96, 0xc6, 0x24, 0x55, 0xa3, 0x05,
	0x96, 0x71, 0xa6, 0x02, 0xeb, 0x2e, 0xd4, 0x53, 0xb5, 0x12, 0x35, 0x2e, 0x2a, 0x3d, 0x09, 0xa3,
	0x22, 0xb5, 0x92, 0xa1, 0xe1, 0x23, 0x79, 0x0f, 0xe6, 0x86, 0x6b, 0xac, 0xad, 0xec, 0x1a, 0xab,
	0x1e, 0xa7, 0x0b, 0xab, 0x41, 0x49, 0xb7, 0x72, 0x7a, 0x49, 0x37, 0x21, 0xb1, 0x7b, 0x03, 0x1a,
	0xd2, 0x0b, 0x34, 0xd6, 0x3a, 0xb2, 0xe6, 0x84, 0x60, 0x27, 0xe1, 0x66, 0xe6, 0x66, 0x17, 0xb2,
	0x73, 0xb3, 0x8c, 0x1a, 0xa1, 0xf9, 0x7a, 0x35, 0xc2, 0xd5, 0xf3, 0xd4, 0x08, 0x59, 0x49, 0xd4,
	0xea, 0xf9, 0x92, 0xa8, 0xd7, 0x2e, 0x5d, 0x3f, 0x85, 0x25, 0x1a, 0x33, 0xaf, 0x8b, 0xf1, 0xcd,
	0xa5, 0xbe, 0xf7, 0x8c, 0x46, 0x27, 0x3c, 0xe4, 0x7c, 0x32, 0x31, 0xe4, 0x2c, 0x24, 0x8a, 0xbb,
	0x52, 0x6f, 0x9b, 0x91, 0x4f, 0x60, 0x00, 0xb7, 0xd5, 0xdc, 0x32, 0xe3, 0x7b, 0x13, 0x5b, 0x6b,
	0x24, 0x6a, 0xfb, 0x62, 0x8e, 0x99, 0xb9, 0x0a, 0xb3, 0x4f, 0xec, 0x68, 0xf0, 0xae, 0x35, 0x37,
	0xfc, 0xae, 0xd5, 0xdc, 0x87, 0xd2, 0xae, 0x17, 0xb3, 0xc8, 0x73, 0xd8, 0x78, 0x16, 0x79, 0x03,
	0xf2, 0xcf, 0xed, 0xc8, 0x55, 0xb7, 0x7d, 0x75, 0xed, 0x35, 0x8f, 0x1d, 0xb9, 0xea, 0x9d, 0x0d,
	0x52, 0xcc, 0x1f, 0x41, 0x09, 0x63, 0x1d, 0xdf, 0x68, 0xe3, 0x5b, 0xbc, 0x0b, 0x65, 0x57, 0xf6,
	0xab, 0x5a, 0x25, 0x83, 0x56, 0x95, 0x49, 0xea, 0xb8, 0x4f, 0xa8, 0xe6, 0x27, 0x30, 0xff, 0x30,
	0x74, 0x6c, 0xe6, 0x85, 0x41, 0xf2, 0x82, 0xf0, 0x2e, 0x94, 0x7b, 0xb2, 0x47, 0xf5, 0x86, 0x50,
	0x6b, 0x4b, 0x19, 0xa3, 0xa5, 0x2a, 0x82, 0x6a, 0xfe, 0x69, 0x06, 0x56, 0x3e, 0xa6, 0x6c, 0x28,
	0xea, 0x25, 0x2f, 0x72, 0xbe, 0x33, 0xa6, 0xfc, 0xbd, 0xa8, 0x5d, 0xe5, 0x4b, 0x53, 0xc6, 0x78,
	0xf7, 0x6e, 0x86, 0x7b, 0x4e, 0x4f, 0x6a, 0x63, 0xc4, 0x43, 0x6f, 0x0c, 0x15, 0x65, 0x55, 0x3d,
	0xc3, 0xd1, 0x4b, 0xb3, 0xb7, 0xb3, 0x36, 0x6e, 0x5e, 0x63, 0x8f, 0x6e, 0xdf, 0xcb, 0x49, 0x35,
	0x57, 0xd0, 0x78, 0xaa, 0xa6, 0x5b, 0x51, 0x35, 0x5d, 0x51, 0x13, 0xca, 0xca, 0xee, 0x72, 0x52,
	0xd9, 0x95, 0x74, 0x4d, 0x59, 0xdf, 0xdd, 0xc8, 0x4a, 0x28, 0xb2, 0x0a, 0x89, 0xb5, 0x54, 0xdc,
	0xaa, 0xe8, 0x39, 0x9e, 0xa3, 0x47, 0xa4, 0xd1, 0xe0, 0x5e, 0xcb, 0x0e, 0xee, 0xe6, 0x67, 0x30,
	0x37, 0x34, 0xa5, 0x3c, 0xd5, 0x54, 0xab, 0x9e, 0x72, 0xc5, 0x04, 0xe5, 0x0c, 0xe5, 0x63, 0xa9,
	0xbb, 0xf3, 0x04, 0x35, 0x1f, 0xc3, 0xa5, 0x4c, 0x5f, 0x91, 0x3e, 0x78, 0x87, 0x67, 0x5f, 0x02,
	0x93, 0x2e, 0xa8, 0x2d, 0xf1, 0x90, 0x96, 0x95, 0x50, 0xcd, 0xdf, 0x4f, 0xc3, 0xdc, 0x90, 0x34,
	0x79, 0x77, 0x90, 0x1b, 0x79, 0x77, 0x30, 0xfe, 0x93, 0x89, 0x65, 0x98, 0x49, 0xee, 0x2c, 0xd4,
	0xeb, 0xe9, 0x23, 0x4a, 0xc9, 0xce, 0x20, 0x45, 0x99, 0x3d, 0x73, 0x8a, 0xa2, 0x12, 0x4f, 0x95,
	0xa8, 0x8c, 0x89, 0x48, 0xf9, 0x57, 0x88, 0x48, 0xe3, 0xa3, 0x65, 0xe1, 0x95, 0xa2, 0xa5, 0xf9,
	0x45, 0x4e, 0xde, 0x5a, 0xed, 0xc8, 0xfa, 0x9a, 0xa7, 0xf0, 0x47, 0x7d, 0xdf, 0x6f, 0x8f, 0x04,
	0x9e, 0x12, 0x87, 0x31, 0x8b, 0x4a, 0x3e, 0xc3, 0x98, 0x1e, 0xfd, 0x0c, 0x23, 0xf9, 0x50, 0x64,
	0x66, 0xf4, 0x43, 0x91, 0x37, 0xa1, 0x20, 0xbe, 0xf8, 0x91, 0xf7, 0x3e, 0xf5, 0x96, 0xfc, 0x00,
	0xe8, 0x63, 0xfc, 0xa7, 0xf6, 0x80, 0x00, 0xcd, 0xdf, 0x4c, 0x43, 0x55, 0x3f, 0x7b, 0xce, 0x68,
	0xd9, 0xd8, 0xde, 0x13, 0xab, 0x67, 0x47, 0xad, 0xd6, 0x3d, 0xbc, 0x38, 0xd1, 0xc3, 0x4b, 0x59,
	0x1e, 0xce, 0xbd, 0x8b, 0x47, 0xf0, 0xd4, 0x17, 0x32, 0x88, 0x90, 0x26, 0x14, 0x9d, 0xb0, 0xdb,
	0xb3, 0x83, 0x13, 0xdc, 0xa1, 0xc9, 0x07, 0x32, 0x12, 0xe4, 0x6d, 0xcb, 0xf0, 0xf6, 0x36, 0x46,
	0xa7, 0xa4, 0x6d, 0x85, 0x6a, 0x8c, 0x2d, 0xdc, 0xb8, 0xc3, 0x8c, 0xad, 0x8f, 0xf6, 0xbe, 0x7c,
	0xd1, 0xcc, 0x7d, 0xf5, 0xa2, 0x99, 0xfb, 0xfa, 0x45, 0x73, 0xea, 0xef, 0x2f, 0x9a, 0x53, 0x3f,
	0x7b, 0xd9, 0x9c, 0xfa, 0xf5, 0xcb, 0xe6, 0xd4, 0xef, 0x5e, 0x36, 0xa7, 0xbe, 0x7c, 0xd9, 0x9c,
	0xfa, 0xea, 0x65, 0x73, 0xea, 0xeb, 0x97, 0xcd, 0xa9, 0x9f, 0x7f, 0xd3, 0x9c, 0xfa, 0xe2, 0x9b,
	0xe6, 0xd4, 0x0f, 0x2f, 0xe0, 0x07, 0x3f, 0xcf, 0x82, 0x4d, 0xbb, 0xe7, 0x6d, 0xf6, 0x0e, 0x93,
	0x2f, 0xbc, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x21, 0x87, 0xe2, 0xc3, 0xf5, 0x25, 0x00, 0x00,
}
