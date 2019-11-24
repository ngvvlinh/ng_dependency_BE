// source: services/crm/crm.proto

package crm

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	notifier_entity "etop.vn/api/pb/etop/etc/notifier_entity"
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

type RefreshFulfillmentFromCarrierRequest struct {
	// @required
	ShippingCode         string   `protobuf:"bytes,1,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RefreshFulfillmentFromCarrierRequest) Reset()         { *m = RefreshFulfillmentFromCarrierRequest{} }
func (m *RefreshFulfillmentFromCarrierRequest) String() string { return proto.CompactTextString(m) }
func (*RefreshFulfillmentFromCarrierRequest) ProtoMessage()    {}
func (*RefreshFulfillmentFromCarrierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{0}
}

var xxx_messageInfo_RefreshFulfillmentFromCarrierRequest proto.InternalMessageInfo

func (m *RefreshFulfillmentFromCarrierRequest) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

type SendNotificationRequest struct {
	AccountId            dot.ID                         `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	Title                string                         `protobuf:"bytes,2,opt,name=title" json:"title"`
	Message              string                         `protobuf:"bytes,3,opt,name=message" json:"message"`
	MetaData             common.RawJSONObject           `protobuf:"bytes,4,opt,name=meta_data,json=metaData" json:"meta_data"`
	Entity               notifier_entity.NotifierEntity `protobuf:"varint,5,opt,name=entity,enum=notifier_entity.NotifierEntity" json:"entity"`
	EntityId             dot.ID                         `protobuf:"varint,6,opt,name=entity_id,json=entityId" json:"entity_id"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *SendNotificationRequest) Reset()         { *m = SendNotificationRequest{} }
func (m *SendNotificationRequest) String() string { return proto.CompactTextString(m) }
func (*SendNotificationRequest) ProtoMessage()    {}
func (*SendNotificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{1}
}

var xxx_messageInfo_SendNotificationRequest proto.InternalMessageInfo

func (m *SendNotificationRequest) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *SendNotificationRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SendNotificationRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *SendNotificationRequest) GetMetaData() common.RawJSONObject {
	if m != nil {
		return m.MetaData
	}
	return common.RawJSONObject{}
}

func (m *SendNotificationRequest) GetEntity() notifier_entity.NotifierEntity {
	if m != nil {
		return m.Entity
	}
	return notifier_entity.NotifierEntity_unknown
}

func (m *SendNotificationRequest) GetEntityId() dot.ID {
	if m != nil {
		return m.EntityId
	}
	return 0
}

type GetCallHistoriesRequest struct {
	Paging               *common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	TextSearch           string         `protobuf:"bytes,2,opt,name=text_search,json=textSearch" json:"text_search"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetCallHistoriesRequest) Reset()         { *m = GetCallHistoriesRequest{} }
func (m *GetCallHistoriesRequest) String() string { return proto.CompactTextString(m) }
func (*GetCallHistoriesRequest) ProtoMessage()    {}
func (*GetCallHistoriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{2}
}

var xxx_messageInfo_GetCallHistoriesRequest proto.InternalMessageInfo

func (m *GetCallHistoriesRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCallHistoriesRequest) GetTextSearch() string {
	if m != nil {
		return m.TextSearch
	}
	return ""
}

type GetCallHistoriesResponse struct {
	VhtCallLog           []*VHTCallLog `protobuf:"bytes,1,rep,name=vht_call_log,json=vhtCallLog" json:"vht_call_log"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetCallHistoriesResponse) Reset()         { *m = GetCallHistoriesResponse{} }
func (m *GetCallHistoriesResponse) String() string { return proto.CompactTextString(m) }
func (*GetCallHistoriesResponse) ProtoMessage()    {}
func (*GetCallHistoriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{3}
}

var xxx_messageInfo_GetCallHistoriesResponse proto.InternalMessageInfo

func (m *GetCallHistoriesResponse) GetVhtCallLog() []*VHTCallLog {
	if m != nil {
		return m.VhtCallLog
	}
	return nil
}

type VHTCallLog struct {
	CdrId                string   `protobuf:"bytes,1,opt,name=cdr_id,json=cdrId" json:"cdr_id"`
	CallId               string   `protobuf:"bytes,2,opt,name=call_id,json=callId" json:"call_id"`
	SipCallId            string   `protobuf:"bytes,3,opt,name=sip_call_id,json=sipCallId" json:"sip_call_id"`
	SdkCallId            string   `protobuf:"bytes,4,opt,name=sdk_call_id,json=sdkCallId" json:"sdk_call_id"`
	Cause                string   `protobuf:"bytes,5,opt,name=cause" json:"cause"`
	Q850Cause            string   `protobuf:"bytes,6,opt,name=q850_cause,json=q850Cause" json:"q850_cause"`
	FromExtension        string   `protobuf:"bytes,7,opt,name=from_extension,json=fromExtension" json:"from_extension"`
	ToExtension          string   `protobuf:"bytes,8,opt,name=to_extension,json=toExtension" json:"to_extension"`
	FromNumber           string   `protobuf:"bytes,9,opt,name=from_number,json=fromNumber" json:"from_number"`
	ToNumber             string   `protobuf:"bytes,10,opt,name=to_number,json=toNumber" json:"to_number"`
	Duration             int32    `protobuf:"varint,11,opt,name=duration" json:"duration"`
	Direction            int32    `protobuf:"varint,12,opt,name=direction" json:"direction"`
	TimeStarted          dot.Time `protobuf:"bytes,13,opt,name=time_started,json=timeStarted" json:"time_started"`
	TimeConnected        dot.Time `protobuf:"bytes,14,opt,name=time_connected,json=timeConnected" json:"time_connected"`
	TimeEnded            dot.Time `protobuf:"bytes,15,opt,name=time_ended,json=timeEnded" json:"time_ended"`
	RecordingPath        string   `protobuf:"bytes,16,opt,name=recording_path,json=recordingPath" json:"recording_path"`
	RecordingUrl         string   `protobuf:"bytes,17,opt,name=recording_url,json=recordingUrl" json:"recording_url"`
	RecordFileSize       int32    `protobuf:"varint,18,opt,name=record_file_size,json=recordFileSize" json:"record_file_size"`
	EtopAccountId        dot.ID   `protobuf:"varint,19,opt,name=etop_account_id,json=etopAccountId" json:"etop_account_id"`
	VtigerAccountId      string   `protobuf:"bytes,21,opt,name=vtiger_account_id,json=vtigerAccountId" json:"vtiger_account_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VHTCallLog) Reset()         { *m = VHTCallLog{} }
func (m *VHTCallLog) String() string { return proto.CompactTextString(m) }
func (*VHTCallLog) ProtoMessage()    {}
func (*VHTCallLog) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{4}
}

var xxx_messageInfo_VHTCallLog proto.InternalMessageInfo

func (m *VHTCallLog) GetCdrId() string {
	if m != nil {
		return m.CdrId
	}
	return ""
}

func (m *VHTCallLog) GetCallId() string {
	if m != nil {
		return m.CallId
	}
	return ""
}

func (m *VHTCallLog) GetSipCallId() string {
	if m != nil {
		return m.SipCallId
	}
	return ""
}

func (m *VHTCallLog) GetSdkCallId() string {
	if m != nil {
		return m.SdkCallId
	}
	return ""
}

func (m *VHTCallLog) GetCause() string {
	if m != nil {
		return m.Cause
	}
	return ""
}

func (m *VHTCallLog) GetQ850Cause() string {
	if m != nil {
		return m.Q850Cause
	}
	return ""
}

func (m *VHTCallLog) GetFromExtension() string {
	if m != nil {
		return m.FromExtension
	}
	return ""
}

func (m *VHTCallLog) GetToExtension() string {
	if m != nil {
		return m.ToExtension
	}
	return ""
}

func (m *VHTCallLog) GetFromNumber() string {
	if m != nil {
		return m.FromNumber
	}
	return ""
}

func (m *VHTCallLog) GetToNumber() string {
	if m != nil {
		return m.ToNumber
	}
	return ""
}

func (m *VHTCallLog) GetDuration() int32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *VHTCallLog) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

func (m *VHTCallLog) GetRecordingPath() string {
	if m != nil {
		return m.RecordingPath
	}
	return ""
}

func (m *VHTCallLog) GetRecordingUrl() string {
	if m != nil {
		return m.RecordingUrl
	}
	return ""
}

func (m *VHTCallLog) GetRecordFileSize() int32 {
	if m != nil {
		return m.RecordFileSize
	}
	return 0
}

func (m *VHTCallLog) GetEtopAccountId() dot.ID {
	if m != nil {
		return m.EtopAccountId
	}
	return 0
}

func (m *VHTCallLog) GetVtigerAccountId() string {
	if m != nil {
		return m.VtigerAccountId
	}
	return ""
}

type CountTicketByStatusRequest struct {
	Status               string   `protobuf:"bytes,1,opt,name=status" json:"status"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountTicketByStatusRequest) Reset()         { *m = CountTicketByStatusRequest{} }
func (m *CountTicketByStatusRequest) String() string { return proto.CompactTextString(m) }
func (*CountTicketByStatusRequest) ProtoMessage()    {}
func (*CountTicketByStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{5}
}

var xxx_messageInfo_CountTicketByStatusRequest proto.InternalMessageInfo

func (m *CountTicketByStatusRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type GetTicketStatusCountResponse struct {
	StatusCount          []*CountTicketByStatusResponse `protobuf:"bytes,1,rep,name=status_count,json=statusCount" json:"status_count"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *GetTicketStatusCountResponse) Reset()         { *m = GetTicketStatusCountResponse{} }
func (m *GetTicketStatusCountResponse) String() string { return proto.CompactTextString(m) }
func (*GetTicketStatusCountResponse) ProtoMessage()    {}
func (*GetTicketStatusCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{6}
}

var xxx_messageInfo_GetTicketStatusCountResponse proto.InternalMessageInfo

func (m *GetTicketStatusCountResponse) GetStatusCount() []*CountTicketByStatusResponse {
	if m != nil {
		return m.StatusCount
	}
	return nil
}

type CountTicketByStatusResponse struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Count                int32    `protobuf:"varint,2,opt,name=count" json:"count"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountTicketByStatusResponse) Reset()         { *m = CountTicketByStatusResponse{} }
func (m *CountTicketByStatusResponse) String() string { return proto.CompactTextString(m) }
func (*CountTicketByStatusResponse) ProtoMessage()    {}
func (*CountTicketByStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{7}
}

var xxx_messageInfo_CountTicketByStatusResponse proto.InternalMessageInfo

func (m *CountTicketByStatusResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CountTicketByStatusResponse) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type GetContactsResponse struct {
	Contacts             []*ContactResponse `protobuf:"bytes,1,rep,name=contacts" json:"contacts"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetContactsResponse) Reset()         { *m = GetContactsResponse{} }
func (m *GetContactsResponse) String() string { return proto.CompactTextString(m) }
func (*GetContactsResponse) ProtoMessage()    {}
func (*GetContactsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{8}
}

var xxx_messageInfo_GetContactsResponse proto.InternalMessageInfo

func (m *GetContactsResponse) GetContacts() []*ContactResponse {
	if m != nil {
		return m.Contacts
	}
	return nil
}

type GetContactsRequest struct {
	TextSearch           string         `protobuf:"bytes,1,opt,name=text_search,json=textSearch" json:"text_search"`
	Paging               *common.Paging `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetContactsRequest) Reset()         { *m = GetContactsRequest{} }
func (m *GetContactsRequest) String() string { return proto.CompactTextString(m) }
func (*GetContactsRequest) ProtoMessage()    {}
func (*GetContactsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{9}
}

var xxx_messageInfo_GetContactsRequest proto.InternalMessageInfo

func (m *GetContactsRequest) GetTextSearch() string {
	if m != nil {
		return m.TextSearch
	}
	return ""
}

func (m *GetContactsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

type ContactRequest struct {
	ContactNo            string   `protobuf:"bytes,1,opt,name=contact_no,json=contactNo" json:"contact_no"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Lastname             string   `protobuf:"bytes,3,opt,name=lastname" json:"lastname"`
	Mobile               string   `protobuf:"bytes,4,opt,name=mobile" json:"mobile"`
	Leadsource           string   `protobuf:"bytes,5,opt,name=leadsource" json:"leadsource"`
	Email                string   `protobuf:"bytes,6,opt,name=email" json:"email"`
	Description          string   `protobuf:"bytes,7,opt,name=description" json:"description"`
	Secondaryemail       string   `protobuf:"bytes,8,opt,name=secondaryemail" json:"secondaryemail"`
	Modifiedby           string   `protobuf:"bytes,10,opt,name=modifiedby" json:"modifiedby"`
	Source               string   `protobuf:"bytes,11,opt,name=source" json:"source"`
	EtopUserId           dot.ID   `protobuf:"varint,12,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Website              string   `protobuf:"bytes,14,opt,name=website" json:"website"`
	Lane                 string   `protobuf:"bytes,15,opt,name=lane" json:"lane"`
	City                 string   `protobuf:"bytes,16,opt,name=city" json:"city"`
	State                string   `protobuf:"bytes,17,opt,name=state" json:"state"`
	Country              string   `protobuf:"bytes,18,opt,name=country" json:"country"`
	OrdersPerDay         string   `protobuf:"bytes,19,opt,name=orders_per_day,json=ordersPerDay" json:"orders_per_day"`
	UsedShippingProvider string   `protobuf:"bytes,20,opt,name=used_shipping_provider,json=usedShippingProvider" json:"used_shipping_provider"`
	Id                   string   `protobuf:"bytes,21,opt,name=id" json:"id"`
	Firstname            string   `protobuf:"bytes,22,opt,name=firstname" json:"firstname"`
	Createdtime          dot.Time `protobuf:"bytes,23,opt,name=createdtime" json:"createdtime"`
	Modifiedtime         dot.Time `protobuf:"bytes,24,opt,name=modifiedtime" json:"modifiedtime"`
	AssignedUserId       string   `protobuf:"bytes,25,opt,name=assigned_user_id,json=assignedUserId" json:"assigned_user_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContactRequest) Reset()         { *m = ContactRequest{} }
func (m *ContactRequest) String() string { return proto.CompactTextString(m) }
func (*ContactRequest) ProtoMessage()    {}
func (*ContactRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{10}
}

var xxx_messageInfo_ContactRequest proto.InternalMessageInfo

func (m *ContactRequest) GetContactNo() string {
	if m != nil {
		return m.ContactNo
	}
	return ""
}

func (m *ContactRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ContactRequest) GetLastname() string {
	if m != nil {
		return m.Lastname
	}
	return ""
}

func (m *ContactRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *ContactRequest) GetLeadsource() string {
	if m != nil {
		return m.Leadsource
	}
	return ""
}

func (m *ContactRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ContactRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ContactRequest) GetSecondaryemail() string {
	if m != nil {
		return m.Secondaryemail
	}
	return ""
}

func (m *ContactRequest) GetModifiedby() string {
	if m != nil {
		return m.Modifiedby
	}
	return ""
}

func (m *ContactRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *ContactRequest) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *ContactRequest) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *ContactRequest) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *ContactRequest) GetLane() string {
	if m != nil {
		return m.Lane
	}
	return ""
}

func (m *ContactRequest) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *ContactRequest) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *ContactRequest) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *ContactRequest) GetOrdersPerDay() string {
	if m != nil {
		return m.OrdersPerDay
	}
	return ""
}

func (m *ContactRequest) GetUsedShippingProvider() string {
	if m != nil {
		return m.UsedShippingProvider
	}
	return ""
}

func (m *ContactRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ContactRequest) GetFirstname() string {
	if m != nil {
		return m.Firstname
	}
	return ""
}

func (m *ContactRequest) GetAssignedUserId() string {
	if m != nil {
		return m.AssignedUserId
	}
	return ""
}

type ContactResponse struct {
	ContactNo            string   `protobuf:"bytes,1,opt,name=contact_no,json=contactNo" json:"contact_no"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Lastname             string   `protobuf:"bytes,3,opt,name=lastname" json:"lastname"`
	Mobile               string   `protobuf:"bytes,4,opt,name=mobile" json:"mobile"`
	Leadsource           string   `protobuf:"bytes,5,opt,name=leadsource" json:"leadsource"`
	Email                string   `protobuf:"bytes,6,opt,name=email" json:"email"`
	Description          string   `protobuf:"bytes,7,opt,name=description" json:"description"`
	Secondaryemail       string   `protobuf:"bytes,8,opt,name=secondaryemail" json:"secondaryemail"`
	Modifiedby           string   `protobuf:"bytes,10,opt,name=modifiedby" json:"modifiedby"`
	Source               string   `protobuf:"bytes,11,opt,name=source" json:"source"`
	EtopUserId           dot.ID   `protobuf:"varint,12,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Website              string   `protobuf:"bytes,14,opt,name=website" json:"website"`
	Lane                 string   `protobuf:"bytes,15,opt,name=lane" json:"lane"`
	City                 string   `protobuf:"bytes,16,opt,name=city" json:"city"`
	State                string   `protobuf:"bytes,17,opt,name=state" json:"state"`
	Country              string   `protobuf:"bytes,18,opt,name=country" json:"country"`
	OrdersPerDay         string   `protobuf:"bytes,19,opt,name=orders_per_day,json=ordersPerDay" json:"orders_per_day"`
	UsedShippingProvider string   `protobuf:"bytes,20,opt,name=used_shipping_provider,json=usedShippingProvider" json:"used_shipping_provider"`
	Id                   string   `protobuf:"bytes,21,opt,name=id" json:"id"`
	Firstname            string   `protobuf:"bytes,22,opt,name=firstname" json:"firstname"`
	Createdtime          dot.Time `protobuf:"bytes,23,opt,name=createdtime" json:"createdtime"`
	Modifiedtime         dot.Time `protobuf:"bytes,24,opt,name=modifiedtime" json:"modifiedtime"`
	AssignedUserId       string   `protobuf:"bytes,25,opt,name=assigned_user_id,json=assignedUserId" json:"assigned_user_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContactResponse) Reset()         { *m = ContactResponse{} }
func (m *ContactResponse) String() string { return proto.CompactTextString(m) }
func (*ContactResponse) ProtoMessage()    {}
func (*ContactResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{11}
}

var xxx_messageInfo_ContactResponse proto.InternalMessageInfo

func (m *ContactResponse) GetContactNo() string {
	if m != nil {
		return m.ContactNo
	}
	return ""
}

func (m *ContactResponse) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ContactResponse) GetLastname() string {
	if m != nil {
		return m.Lastname
	}
	return ""
}

func (m *ContactResponse) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *ContactResponse) GetLeadsource() string {
	if m != nil {
		return m.Leadsource
	}
	return ""
}

func (m *ContactResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ContactResponse) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ContactResponse) GetSecondaryemail() string {
	if m != nil {
		return m.Secondaryemail
	}
	return ""
}

func (m *ContactResponse) GetModifiedby() string {
	if m != nil {
		return m.Modifiedby
	}
	return ""
}

func (m *ContactResponse) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *ContactResponse) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *ContactResponse) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *ContactResponse) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *ContactResponse) GetLane() string {
	if m != nil {
		return m.Lane
	}
	return ""
}

func (m *ContactResponse) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *ContactResponse) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *ContactResponse) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *ContactResponse) GetOrdersPerDay() string {
	if m != nil {
		return m.OrdersPerDay
	}
	return ""
}

func (m *ContactResponse) GetUsedShippingProvider() string {
	if m != nil {
		return m.UsedShippingProvider
	}
	return ""
}

func (m *ContactResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ContactResponse) GetFirstname() string {
	if m != nil {
		return m.Firstname
	}
	return ""
}

func (m *ContactResponse) GetAssignedUserId() string {
	if m != nil {
		return m.AssignedUserId
	}
	return ""
}

type LeadRequest struct {
	ContactNo            string   `protobuf:"bytes,1,opt,name=contact_no,json=contactNo" json:"contact_no"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Lastname             string   `protobuf:"bytes,3,opt,name=lastname" json:"lastname"`
	Mobile               string   `protobuf:"bytes,4,opt,name=mobile" json:"mobile"`
	Leadsource           string   `protobuf:"bytes,5,opt,name=leadsource" json:"leadsource"`
	Email                string   `protobuf:"bytes,6,opt,name=email" json:"email"`
	Secondaryemail       string   `protobuf:"bytes,7,opt,name=secondaryemail" json:"secondaryemail"`
	AssignedUserId       string   `protobuf:"bytes,8,opt,name=assigned_user_id,json=assignedUserId" json:"assigned_user_id"`
	Description          string   `protobuf:"bytes,9,opt,name=description" json:"description"`
	Modifiedby           string   `protobuf:"bytes,10,opt,name=modifiedby" json:"modifiedby"`
	Source               string   `protobuf:"bytes,11,opt,name=source" json:"source"`
	EtopUserId           dot.ID   `protobuf:"varint,12,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Website              string   `protobuf:"bytes,14,opt,name=website" json:"website"`
	Lane                 string   `protobuf:"bytes,15,opt,name=lane" json:"lane"`
	City                 string   `protobuf:"bytes,16,opt,name=city" json:"city"`
	State                string   `protobuf:"bytes,17,opt,name=state" json:"state"`
	Country              string   `protobuf:"bytes,18,opt,name=country" json:"country"`
	OrdersPerDay         string   `protobuf:"bytes,19,opt,name=orders_per_day,json=ordersPerDay" json:"orders_per_day"`
	UsedShippingProvider string   `protobuf:"bytes,20,opt,name=used_shipping_provider,json=usedShippingProvider" json:"used_shipping_provider"`
	Id                   string   `protobuf:"bytes,21,opt,name=id" json:"id"`
	Firstname            string   `protobuf:"bytes,22,opt,name=firstname" json:"firstname"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeadRequest) Reset()         { *m = LeadRequest{} }
func (m *LeadRequest) String() string { return proto.CompactTextString(m) }
func (*LeadRequest) ProtoMessage()    {}
func (*LeadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{12}
}

var xxx_messageInfo_LeadRequest proto.InternalMessageInfo

func (m *LeadRequest) GetContactNo() string {
	if m != nil {
		return m.ContactNo
	}
	return ""
}

func (m *LeadRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *LeadRequest) GetLastname() string {
	if m != nil {
		return m.Lastname
	}
	return ""
}

func (m *LeadRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *LeadRequest) GetLeadsource() string {
	if m != nil {
		return m.Leadsource
	}
	return ""
}

func (m *LeadRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *LeadRequest) GetSecondaryemail() string {
	if m != nil {
		return m.Secondaryemail
	}
	return ""
}

func (m *LeadRequest) GetAssignedUserId() string {
	if m != nil {
		return m.AssignedUserId
	}
	return ""
}

func (m *LeadRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *LeadRequest) GetModifiedby() string {
	if m != nil {
		return m.Modifiedby
	}
	return ""
}

func (m *LeadRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *LeadRequest) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *LeadRequest) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *LeadRequest) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *LeadRequest) GetLane() string {
	if m != nil {
		return m.Lane
	}
	return ""
}

func (m *LeadRequest) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *LeadRequest) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *LeadRequest) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *LeadRequest) GetOrdersPerDay() string {
	if m != nil {
		return m.OrdersPerDay
	}
	return ""
}

func (m *LeadRequest) GetUsedShippingProvider() string {
	if m != nil {
		return m.UsedShippingProvider
	}
	return ""
}

func (m *LeadRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *LeadRequest) GetFirstname() string {
	if m != nil {
		return m.Firstname
	}
	return ""
}

type LeadResponse struct {
	ContactNo            string   `protobuf:"bytes,1,opt,name=contact_no,json=contactNo" json:"contact_no"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Lastname             string   `protobuf:"bytes,3,opt,name=lastname" json:"lastname"`
	Mobile               string   `protobuf:"bytes,4,opt,name=mobile" json:"mobile"`
	Leadsource           string   `protobuf:"bytes,5,opt,name=leadsource" json:"leadsource"`
	Email                string   `protobuf:"bytes,6,opt,name=email" json:"email"`
	Secondaryemail       string   `protobuf:"bytes,7,opt,name=secondaryemail" json:"secondaryemail"`
	AssignedUserId       string   `protobuf:"bytes,8,opt,name=assigned_user_id,json=assignedUserId" json:"assigned_user_id"`
	Description          string   `protobuf:"bytes,9,opt,name=description" json:"description"`
	Modifiedby           string   `protobuf:"bytes,10,opt,name=modifiedby" json:"modifiedby"`
	Source               string   `protobuf:"bytes,11,opt,name=source" json:"source"`
	EtopUserId           dot.ID   `protobuf:"varint,12,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Website              string   `protobuf:"bytes,14,opt,name=website" json:"website"`
	Lane                 string   `protobuf:"bytes,15,opt,name=lane" json:"lane"`
	City                 string   `protobuf:"bytes,16,opt,name=city" json:"city"`
	State                string   `protobuf:"bytes,17,opt,name=state" json:"state"`
	Country              string   `protobuf:"bytes,18,opt,name=country" json:"country"`
	OrdersPerDay         string   `protobuf:"bytes,19,opt,name=orders_per_day,json=ordersPerDay" json:"orders_per_day"`
	UsedShippingProvider string   `protobuf:"bytes,20,opt,name=used_shipping_provider,json=usedShippingProvider" json:"used_shipping_provider"`
	Id                   string   `protobuf:"bytes,21,opt,name=id" json:"id"`
	Firstname            string   `protobuf:"bytes,22,opt,name=firstname" json:"firstname"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeadResponse) Reset()         { *m = LeadResponse{} }
func (m *LeadResponse) String() string { return proto.CompactTextString(m) }
func (*LeadResponse) ProtoMessage()    {}
func (*LeadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{13}
}

var xxx_messageInfo_LeadResponse proto.InternalMessageInfo

func (m *LeadResponse) GetContactNo() string {
	if m != nil {
		return m.ContactNo
	}
	return ""
}

func (m *LeadResponse) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *LeadResponse) GetLastname() string {
	if m != nil {
		return m.Lastname
	}
	return ""
}

func (m *LeadResponse) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *LeadResponse) GetLeadsource() string {
	if m != nil {
		return m.Leadsource
	}
	return ""
}

func (m *LeadResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *LeadResponse) GetSecondaryemail() string {
	if m != nil {
		return m.Secondaryemail
	}
	return ""
}

func (m *LeadResponse) GetAssignedUserId() string {
	if m != nil {
		return m.AssignedUserId
	}
	return ""
}

func (m *LeadResponse) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *LeadResponse) GetModifiedby() string {
	if m != nil {
		return m.Modifiedby
	}
	return ""
}

func (m *LeadResponse) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *LeadResponse) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *LeadResponse) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *LeadResponse) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *LeadResponse) GetLane() string {
	if m != nil {
		return m.Lane
	}
	return ""
}

func (m *LeadResponse) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *LeadResponse) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *LeadResponse) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *LeadResponse) GetOrdersPerDay() string {
	if m != nil {
		return m.OrdersPerDay
	}
	return ""
}

func (m *LeadResponse) GetUsedShippingProvider() string {
	if m != nil {
		return m.UsedShippingProvider
	}
	return ""
}

func (m *LeadResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *LeadResponse) GetFirstname() string {
	if m != nil {
		return m.Firstname
	}
	return ""
}

type GetTicketsRequest struct {
	Paging               *common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Ticket               TicketRequest  `protobuf:"bytes,2,opt,name=ticket" json:"ticket"`
	Orderby              OrderBy        `protobuf:"bytes,3,opt,name=orderby" json:"orderby"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetTicketsRequest) Reset()         { *m = GetTicketsRequest{} }
func (m *GetTicketsRequest) String() string { return proto.CompactTextString(m) }
func (*GetTicketsRequest) ProtoMessage()    {}
func (*GetTicketsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{14}
}

var xxx_messageInfo_GetTicketsRequest proto.InternalMessageInfo

func (m *GetTicketsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTicketsRequest) GetTicket() TicketRequest {
	if m != nil {
		return m.Ticket
	}
	return TicketRequest{}
}

func (m *GetTicketsRequest) GetOrderby() OrderBy {
	if m != nil {
		return m.Orderby
	}
	return OrderBy{}
}

type GetTicketsResponse struct {
	Tickets              []*Ticket `protobuf:"bytes,1,rep,name=tickets" json:"tickets"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetTicketsResponse) Reset()         { *m = GetTicketsResponse{} }
func (m *GetTicketsResponse) String() string { return proto.CompactTextString(m) }
func (*GetTicketsResponse) ProtoMessage()    {}
func (*GetTicketsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{15}
}

var xxx_messageInfo_GetTicketsResponse proto.InternalMessageInfo

func (m *GetTicketsResponse) GetTickets() []*Ticket {
	if m != nil {
		return m.Tickets
	}
	return nil
}

type OrderBy struct {
	Field                string   `protobuf:"bytes,1,opt,name=field" json:"field"`
	Sort                 string   `protobuf:"bytes,2,opt,name=sort" json:"sort"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderBy) Reset()         { *m = OrderBy{} }
func (m *OrderBy) String() string { return proto.CompactTextString(m) }
func (*OrderBy) ProtoMessage()    {}
func (*OrderBy) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{16}
}

var xxx_messageInfo_OrderBy proto.InternalMessageInfo

func (m *OrderBy) GetField() string {
	if m != nil {
		return m.Field
	}
	return ""
}

func (m *OrderBy) GetSort() string {
	if m != nil {
		return m.Sort
	}
	return ""
}

type CreateOrUpdateTicketRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id"`
	Code                 string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Title                string   `protobuf:"bytes,3,opt,name=title" json:"title"`
	Value                string   `protobuf:"bytes,4,opt,name=value" json:"value"`
	OldValue             string   `protobuf:"bytes,5,opt,name=old_value,json=oldValue" json:"old_value"`
	Reason               string   `protobuf:"bytes,6,opt,name=reason" json:"reason"`
	EtopUserId           dot.ID   `protobuf:"varint,7,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	OrderId              dot.ID   `protobuf:"varint,8,opt,name=order_id,json=orderId" json:"order_id"`
	OrderCode            string   `protobuf:"bytes,9,opt,name=order_code,json=orderCode" json:"order_code"`
	FfmCode              string   `protobuf:"bytes,10,opt,name=ffm_code,json=ffmCode" json:"ffm_code"`
	FfmUrl               string   `protobuf:"bytes,11,opt,name=ffm_url,json=ffmUrl" json:"ffm_url"`
	FfmId                dot.ID   `protobuf:"varint,12,opt,name=ffm_id,json=ffmId" json:"ffm_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Provider             string   `protobuf:"bytes,15,opt,name=provider" json:"provider"`
	Note                 string   `protobuf:"bytes,16,opt,name=note" json:"note"`
	Environment          string   `protobuf:"bytes,17,opt,name=environment" json:"environment"`
	FromApp              string   `protobuf:"bytes,18,opt,name=from_app,json=fromApp" json:"from_app"`
	Account              Account  `protobuf:"bytes,19,opt,name=account" json:"account"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrUpdateTicketRequest) Reset()         { *m = CreateOrUpdateTicketRequest{} }
func (m *CreateOrUpdateTicketRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrUpdateTicketRequest) ProtoMessage()    {}
func (*CreateOrUpdateTicketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{17}
}

var xxx_messageInfo_CreateOrUpdateTicketRequest proto.InternalMessageInfo

func (m *CreateOrUpdateTicketRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetOldValue() string {
	if m != nil {
		return m.OldValue
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *CreateOrUpdateTicketRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *CreateOrUpdateTicketRequest) GetOrderCode() string {
	if m != nil {
		return m.OrderCode
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetFfmCode() string {
	if m != nil {
		return m.FfmCode
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetFfmUrl() string {
	if m != nil {
		return m.FfmUrl
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetFfmId() dot.ID {
	if m != nil {
		return m.FfmId
	}
	return 0
}

func (m *CreateOrUpdateTicketRequest) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetFromApp() string {
	if m != nil {
		return m.FromApp
	}
	return ""
}

func (m *CreateOrUpdateTicketRequest) GetAccount() Account {
	if m != nil {
		return m.Account
	}
	return Account{}
}

type TicketRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id"`
	Code                 string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Title                string   `protobuf:"bytes,3,opt,name=title" json:"title"`
	Value                string   `protobuf:"bytes,4,opt,name=value" json:"value"`
	OldValue             string   `protobuf:"bytes,5,opt,name=old_value,json=oldValue" json:"old_value"`
	Reason               string   `protobuf:"bytes,6,opt,name=reason" json:"reason"`
	EtopUserId           dot.ID   `protobuf:"varint,7,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	OrderId              dot.ID   `protobuf:"varint,8,opt,name=order_id,json=orderId" json:"order_id"`
	OrderCode            string   `protobuf:"bytes,9,opt,name=order_code,json=orderCode" json:"order_code"`
	FfmCode              string   `protobuf:"bytes,10,opt,name=ffm_code,json=ffmCode" json:"ffm_code"`
	FfmUrl               string   `protobuf:"bytes,11,opt,name=ffm_url,json=ffmUrl" json:"ffm_url"`
	FfmId                dot.ID   `protobuf:"varint,12,opt,name=ffm_id,json=ffmId" json:"ffm_id"`
	Company              string   `protobuf:"bytes,13,opt,name=company" json:"company"`
	Provider             string   `protobuf:"bytes,15,opt,name=provider" json:"provider"`
	Note                 string   `protobuf:"bytes,16,opt,name=note" json:"note"`
	Environment          string   `protobuf:"bytes,17,opt,name=environment" json:"environment"`
	FromApp              string   `protobuf:"bytes,18,opt,name=from_app,json=fromApp" json:"from_app"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TicketRequest) Reset()         { *m = TicketRequest{} }
func (m *TicketRequest) String() string { return proto.CompactTextString(m) }
func (*TicketRequest) ProtoMessage()    {}
func (*TicketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{18}
}

var xxx_messageInfo_TicketRequest proto.InternalMessageInfo

func (m *TicketRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TicketRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *TicketRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *TicketRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *TicketRequest) GetOldValue() string {
	if m != nil {
		return m.OldValue
	}
	return ""
}

func (m *TicketRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *TicketRequest) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *TicketRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *TicketRequest) GetOrderCode() string {
	if m != nil {
		return m.OrderCode
	}
	return ""
}

func (m *TicketRequest) GetFfmCode() string {
	if m != nil {
		return m.FfmCode
	}
	return ""
}

func (m *TicketRequest) GetFfmUrl() string {
	if m != nil {
		return m.FfmUrl
	}
	return ""
}

func (m *TicketRequest) GetFfmId() dot.ID {
	if m != nil {
		return m.FfmId
	}
	return 0
}

func (m *TicketRequest) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *TicketRequest) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *TicketRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *TicketRequest) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *TicketRequest) GetFromApp() string {
	if m != nil {
		return m.FromApp
	}
	return ""
}

type Account struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	FullName             string   `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	ShortName            string   `protobuf:"bytes,3,opt,name=short_name,json=shortName" json:"short_name"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,5,opt,name=email" json:"email"`
	Company              string   `protobuf:"bytes,6,opt,name=company" json:"company"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{19}
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Account) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Account) GetShortName() string {
	if m != nil {
		return m.ShortName
	}
	return ""
}

func (m *Account) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Account) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Account) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

type Ticket struct {
	TicketNo             string   `protobuf:"bytes,1,opt,name=ticket_no,json=ticketNo" json:"ticket_no"`
	AssignedUserId       string   `protobuf:"bytes,2,opt,name=assigned_user_id,json=assignedUserId" json:"assigned_user_id"`
	ParentId             dot.ID   `protobuf:"varint,3,opt,name=parent_id,json=parentId" json:"parent_id"`
	Ticketpriorities     string   `protobuf:"bytes,4,opt,name=ticketpriorities" json:"ticketpriorities"`
	ProductId            dot.ID   `protobuf:"varint,5,opt,name=product_id,json=productId" json:"product_id"`
	Ticketseverities     string   `protobuf:"bytes,6,opt,name=ticketseverities" json:"ticketseverities"`
	Ticketstatus         string   `protobuf:"bytes,7,opt,name=ticketstatus" json:"ticketstatus"`
	Ticketcategories     string   `protobuf:"bytes,8,opt,name=ticketcategories" json:"ticketcategories"`
	UpdateLog            string   `protobuf:"bytes,9,opt,name=update_log,json=updateLog" json:"update_log"`
	Hours                string   `protobuf:"bytes,10,opt,name=hours" json:"hours"`
	Days                 string   `protobuf:"bytes,11,opt,name=days" json:"days"`
	Createdtime          dot.Time `protobuf:"bytes,12,opt,name=createdtime" json:"createdtime"`
	Modifiedtime         dot.Time `protobuf:"bytes,13,opt,name=modifiedtime" json:"modifiedtime"`
	FromPortal           string   `protobuf:"bytes,14,opt,name=from_portal,json=fromPortal" json:"from_portal"`
	Modifiedby           string   `protobuf:"bytes,15,opt,name=modifiedby" json:"modifiedby"`
	TicketTitle          string   `protobuf:"bytes,16,opt,name=ticket_title,json=ticketTitle" json:"ticket_title"`
	Description          string   `protobuf:"bytes,17,opt,name=description" json:"description"`
	Solution             string   `protobuf:"bytes,18,opt,name=solution" json:"solution"`
	ContactId            string   `protobuf:"bytes,19,opt,name=contact_id,json=contactId" json:"contact_id"`
	Source               string   `protobuf:"bytes,20,opt,name=source" json:"source"`
	Starred              string   `protobuf:"bytes,21,opt,name=starred" json:"starred"`
	Tags                 string   `protobuf:"bytes,22,opt,name=tags" json:"tags"`
	Note                 string   `protobuf:"bytes,24,opt,name=note" json:"note"`
	FfmCode              string   `protobuf:"bytes,25,opt,name=ffm_code,json=ffmCode" json:"ffm_code"`
	FfmUrl               string   `protobuf:"bytes,26,opt,name=ffm_url,json=ffmUrl" json:"ffm_url"`
	FfmId                dot.ID   `protobuf:"varint,27,opt,name=ffm_id,json=ffmId" json:"ffm_id"`
	EtopUserId           dot.ID   `protobuf:"varint,28,opt,name=etop_user_id,json=etopUserId" json:"etop_user_id"`
	OrderId              dot.ID   `protobuf:"varint,29,opt,name=order_id,json=orderId" json:"order_id"`
	OrderCode            string   `protobuf:"bytes,30,opt,name=order_code,json=orderCode" json:"order_code"`
	Company              string   `protobuf:"bytes,31,opt,name=company" json:"company"`
	Provider             string   `protobuf:"bytes,32,opt,name=provider" json:"provider"`
	FromApp              string   `protobuf:"bytes,33,opt,name=from_app,json=fromApp" json:"from_app"`
	Environment          string   `protobuf:"bytes,34,opt,name=environment" json:"environment"`
	Code                 string   `protobuf:"bytes,35,opt,name=code" json:"code"`
	OldValue             string   `protobuf:"bytes,36,opt,name=old_value,json=oldValue" json:"old_value"`
	NewValue             string   `protobuf:"bytes,37,opt,name=new_value,json=newValue" json:"new_value"`
	Substatus            string   `protobuf:"bytes,38,opt,name=substatus" json:"substatus"`
	EtopNote             string   `protobuf:"bytes,39,opt,name=etop_note,json=etopNote" json:"etop_note"`
	Reason               string   `protobuf:"bytes,40,opt,name=reason" json:"reason"`
	Id                   string   `protobuf:"bytes,41,opt,name=id" json:"id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ticket) Reset()         { *m = Ticket{} }
func (m *Ticket) String() string { return proto.CompactTextString(m) }
func (*Ticket) ProtoMessage()    {}
func (*Ticket) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{20}
}

var xxx_messageInfo_Ticket proto.InternalMessageInfo

func (m *Ticket) GetTicketNo() string {
	if m != nil {
		return m.TicketNo
	}
	return ""
}

func (m *Ticket) GetAssignedUserId() string {
	if m != nil {
		return m.AssignedUserId
	}
	return ""
}

func (m *Ticket) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

func (m *Ticket) GetTicketpriorities() string {
	if m != nil {
		return m.Ticketpriorities
	}
	return ""
}

func (m *Ticket) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *Ticket) GetTicketseverities() string {
	if m != nil {
		return m.Ticketseverities
	}
	return ""
}

func (m *Ticket) GetTicketstatus() string {
	if m != nil {
		return m.Ticketstatus
	}
	return ""
}

func (m *Ticket) GetTicketcategories() string {
	if m != nil {
		return m.Ticketcategories
	}
	return ""
}

func (m *Ticket) GetUpdateLog() string {
	if m != nil {
		return m.UpdateLog
	}
	return ""
}

func (m *Ticket) GetHours() string {
	if m != nil {
		return m.Hours
	}
	return ""
}

func (m *Ticket) GetDays() string {
	if m != nil {
		return m.Days
	}
	return ""
}

func (m *Ticket) GetFromPortal() string {
	if m != nil {
		return m.FromPortal
	}
	return ""
}

func (m *Ticket) GetModifiedby() string {
	if m != nil {
		return m.Modifiedby
	}
	return ""
}

func (m *Ticket) GetTicketTitle() string {
	if m != nil {
		return m.TicketTitle
	}
	return ""
}

func (m *Ticket) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Ticket) GetSolution() string {
	if m != nil {
		return m.Solution
	}
	return ""
}

func (m *Ticket) GetContactId() string {
	if m != nil {
		return m.ContactId
	}
	return ""
}

func (m *Ticket) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Ticket) GetStarred() string {
	if m != nil {
		return m.Starred
	}
	return ""
}

func (m *Ticket) GetTags() string {
	if m != nil {
		return m.Tags
	}
	return ""
}

func (m *Ticket) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Ticket) GetFfmCode() string {
	if m != nil {
		return m.FfmCode
	}
	return ""
}

func (m *Ticket) GetFfmUrl() string {
	if m != nil {
		return m.FfmUrl
	}
	return ""
}

func (m *Ticket) GetFfmId() dot.ID {
	if m != nil {
		return m.FfmId
	}
	return 0
}

func (m *Ticket) GetEtopUserId() dot.ID {
	if m != nil {
		return m.EtopUserId
	}
	return 0
}

func (m *Ticket) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *Ticket) GetOrderCode() string {
	if m != nil {
		return m.OrderCode
	}
	return ""
}

func (m *Ticket) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *Ticket) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *Ticket) GetFromApp() string {
	if m != nil {
		return m.FromApp
	}
	return ""
}

func (m *Ticket) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *Ticket) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Ticket) GetOldValue() string {
	if m != nil {
		return m.OldValue
	}
	return ""
}

func (m *Ticket) GetNewValue() string {
	if m != nil {
		return m.NewValue
	}
	return ""
}

func (m *Ticket) GetSubstatus() string {
	if m != nil {
		return m.Substatus
	}
	return ""
}

func (m *Ticket) GetEtopNote() string {
	if m != nil {
		return m.EtopNote
	}
	return ""
}

func (m *Ticket) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Ticket) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetCategoriesResponse struct {
	Categories           []*Category `protobuf:"bytes,1,rep,name=categories" json:"categories"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetCategoriesResponse) Reset()         { *m = GetCategoriesResponse{} }
func (m *GetCategoriesResponse) String() string { return proto.CompactTextString(m) }
func (*GetCategoriesResponse) ProtoMessage()    {}
func (*GetCategoriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{21}
}

var xxx_messageInfo_GetCategoriesResponse proto.InternalMessageInfo

func (m *GetCategoriesResponse) GetCategories() []*Category {
	if m != nil {
		return m.Categories
	}
	return nil
}

type Category struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Label                string   `protobuf:"bytes,2,opt,name=label" json:"label"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Category) Reset()         { *m = Category{} }
func (m *Category) String() string { return proto.CompactTextString(m) }
func (*Category) ProtoMessage()    {}
func (*Category) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{22}
}

var xxx_messageInfo_Category proto.InternalMessageInfo

func (m *Category) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Category) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

type GetStatusResponse struct {
	Status               []*Status `protobuf:"bytes,1,rep,name=status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetStatusResponse) Reset()         { *m = GetStatusResponse{} }
func (m *GetStatusResponse) String() string { return proto.CompactTextString(m) }
func (*GetStatusResponse) ProtoMessage()    {}
func (*GetStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{23}
}

var xxx_messageInfo_GetStatusResponse proto.InternalMessageInfo

func (m *GetStatusResponse) GetStatus() []*Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type Status struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Count                string   `protobuf:"bytes,2,opt,name=count" json:"count"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_561df3d6f5fdbaa9, []int{24}
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Status) GetCount() string {
	if m != nil {
		return m.Count
	}
	return ""
}

func init() {
	proto.RegisterType((*RefreshFulfillmentFromCarrierRequest)(nil), "crm.RefreshFulfillmentFromCarrierRequest")
	proto.RegisterType((*SendNotificationRequest)(nil), "crm.SendNotificationRequest")
	proto.RegisterType((*GetCallHistoriesRequest)(nil), "crm.GetCallHistoriesRequest")
	proto.RegisterType((*GetCallHistoriesResponse)(nil), "crm.GetCallHistoriesResponse")
	proto.RegisterType((*VHTCallLog)(nil), "crm.VHTCallLog")
	proto.RegisterType((*CountTicketByStatusRequest)(nil), "crm.CountTicketByStatusRequest")
	proto.RegisterType((*GetTicketStatusCountResponse)(nil), "crm.GetTicketStatusCountResponse")
	proto.RegisterType((*CountTicketByStatusResponse)(nil), "crm.CountTicketByStatusResponse")
	proto.RegisterType((*GetContactsResponse)(nil), "crm.GetContactsResponse")
	proto.RegisterType((*GetContactsRequest)(nil), "crm.GetContactsRequest")
	proto.RegisterType((*ContactRequest)(nil), "crm.ContactRequest")
	proto.RegisterType((*ContactResponse)(nil), "crm.ContactResponse")
	proto.RegisterType((*LeadRequest)(nil), "crm.LeadRequest")
	proto.RegisterType((*LeadResponse)(nil), "crm.LeadResponse")
	proto.RegisterType((*GetTicketsRequest)(nil), "crm.GetTicketsRequest")
	proto.RegisterType((*GetTicketsResponse)(nil), "crm.GetTicketsResponse")
	proto.RegisterType((*OrderBy)(nil), "crm.OrderBy")
	proto.RegisterType((*CreateOrUpdateTicketRequest)(nil), "crm.CreateOrUpdateTicketRequest")
	proto.RegisterType((*TicketRequest)(nil), "crm.TicketRequest")
	proto.RegisterType((*Account)(nil), "crm.Account")
	proto.RegisterType((*Ticket)(nil), "crm.Ticket")
	proto.RegisterType((*GetCategoriesResponse)(nil), "crm.GetCategoriesResponse")
	proto.RegisterType((*Category)(nil), "crm.Category")
	proto.RegisterType((*GetStatusResponse)(nil), "crm.GetStatusResponse")
	proto.RegisterType((*Status)(nil), "crm.Status")
}

func init() { proto.RegisterFile("services/crm/crm.proto", fileDescriptor_561df3d6f5fdbaa9) }

var fileDescriptor_561df3d6f5fdbaa9 = []byte{
	// 2232 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x5a, 0xdd, 0x72, 0x1b, 0xb7,
	0xf5, 0x27, 0xf5, 0xc1, 0x8f, 0x43, 0x4a, 0xb2, 0xd7, 0x8a, 0xcd, 0x48, 0x36, 0x25, 0xad, 0xbf,
	0x94, 0xff, 0xdf, 0x25, 0x5d, 0x4f, 0x3b, 0x93, 0xa6, 0x4d, 0x5a, 0x5b, 0xb1, 0x1d, 0x75, 0x1c,
	0xd9, 0x95, 0xec, 0x5c, 0xf4, 0x66, 0x07, 0x5a, 0x80, 0x24, 0xea, 0xdd, 0xc5, 0x06, 0x00, 0xe5,
	0x30, 0x4f, 0xd0, 0xcb, 0x4e, 0xa7, 0x69, 0xfb, 0x08, 0x79, 0x84, 0xde, 0xf5, 0xa6, 0x17, 0xbe,
	0xcc, 0x13, 0x74, 0x6c, 0xa5, 0x0f, 0xd0, 0x47, 0xe8, 0xe0, 0x63, 0x97, 0x58, 0x92, 0xb2, 0x9c,
	0xc9, 0x45, 0xa7, 0x19, 0x5d, 0x78, 0xb4, 0xfc, 0x9d, 0x1f, 0xce, 0x62, 0x01, 0x9c, 0xdf, 0x39,
	0x00, 0x0c, 0x17, 0x05, 0xe1, 0x47, 0x34, 0x24, 0xa2, 0x1b, 0xf2, 0x58, 0xfd, 0xeb, 0xa4, 0x9c,
	0x49, 0xe6, 0xcd, 0x87, 0x3c, 0x5e, 0x5b, 0xed, 0xb3, 0x3e, 0xd3, 0xbf, 0xbb, 0xea, 0xc9, 0x98,
	0xd6, 0x2e, 0x84, 0x2c, 0x8e, 0x59, 0xd2, 0x35, 0x7f, 0x2c, 0x78, 0x9d, 0x48, 0x96, 0x76, 0x89,
	0x0c, 0xbb, 0x09, 0x93, 0xb4, 0x47, 0x09, 0x0f, 0x48, 0x22, 0xa9, 0x1c, 0x75, 0xcd, 0x1f, 0x4b,
	0xdb, 0xe8, 0x33, 0xd6, 0x8f, 0x48, 0x57, 0xff, 0x3a, 0x1c, 0xf6, 0xba, 0x92, 0xc6, 0x44, 0x48,
	0x14, 0xa7, 0x86, 0xe0, 0xff, 0x06, 0xae, 0xed, 0x93, 0x1e, 0x27, 0x62, 0xf0, 0x60, 0x18, 0xf5,
	0x68, 0x14, 0xc5, 0x24, 0x91, 0x0f, 0x38, 0x8b, 0x77, 0x10, 0xe7, 0x94, 0xf0, 0x7d, 0xf2, 0xf9,
	0x90, 0x08, 0xe9, 0xbd, 0x07, 0x4b, 0x62, 0x40, 0xd3, 0x94, 0x26, 0xfd, 0x20, 0x64, 0x98, 0xb4,
	0xca, 0x9b, 0xe5, 0xed, 0xfa, 0xbd, 0x85, 0x97, 0xff, 0xdc, 0x28, 0xed, 0x37, 0x33, 0xd3, 0x0e,
	0xc3, 0xc4, 0xff, 0xe3, 0x1c, 0x5c, 0x3a, 0x20, 0x09, 0xde, 0xd3, 0x1d, 0x0b, 0x91, 0xa4, 0x2c,
	0xc9, 0xdc, 0x5c, 0x05, 0x40, 0x61, 0xc8, 0x86, 0x89, 0x0c, 0x28, 0xd6, 0x3e, 0xe6, 0xad, 0x8f,
	0xba, 0xc5, 0x77, 0xb1, 0xb7, 0x06, 0x8b, 0x92, 0xca, 0x88, 0xb4, 0xe6, 0x9c, 0x77, 0x18, 0xc8,
	0x6b, 0x43, 0x35, 0x26, 0x42, 0xa0, 0x3e, 0x69, 0xcd, 0x3b, 0xd6, 0x0c, 0xf4, 0x7e, 0x02, 0xf5,
	0x98, 0x48, 0x14, 0x60, 0x24, 0x51, 0x6b, 0x61, 0xb3, 0xbc, 0xdd, 0xb8, 0x73, 0xbe, 0x13, 0xc6,
	0x9d, 0x7d, 0xf4, 0xe2, 0xd7, 0x07, 0x8f, 0xf7, 0x1e, 0x1f, 0xfe, 0x8e, 0x84, 0xd2, 0x36, 0xaa,
	0x29, 0xe6, 0xc7, 0x48, 0x22, 0xef, 0x43, 0xa8, 0x98, 0x61, 0x6b, 0x2d, 0x6e, 0x96, 0xb7, 0x97,
	0xef, 0x6c, 0x74, 0x26, 0x46, 0xb5, 0xb3, 0x67, 0x7f, 0xdf, 0xd7, 0x3f, 0xad, 0x03, 0xdb, 0xc8,
	0xdb, 0x82, 0xba, 0x79, 0x52, 0x1f, 0x55, 0x71, 0x3e, 0xaa, 0x66, 0xe0, 0x5d, 0xec, 0x63, 0xb8,
	0xf4, 0x90, 0xc8, 0x1d, 0x14, 0x45, 0x9f, 0x50, 0x21, 0x19, 0xa7, 0x44, 0x64, 0x63, 0xe2, 0x43,
	0x25, 0x45, 0x7d, 0x9a, 0xf4, 0xf5, 0x78, 0x34, 0xee, 0x80, 0xea, 0xef, 0x13, 0x8d, 0xec, 0x5b,
	0x8b, 0x77, 0x1d, 0x1a, 0x92, 0x7c, 0x21, 0x03, 0x41, 0x10, 0x0f, 0x07, 0x85, 0x81, 0x01, 0x65,
	0x38, 0xd0, 0xb8, 0xff, 0x29, 0xb4, 0xa6, 0xdf, 0x22, 0x52, 0x96, 0x08, 0xe2, 0xfd, 0x18, 0x9a,
	0x47, 0x03, 0x19, 0x84, 0x28, 0x8a, 0x82, 0x88, 0xa9, 0x97, 0xcd, 0x6f, 0x37, 0xee, 0xac, 0x74,
	0xd4, 0x1a, 0xfc, 0xec, 0x93, 0xa7, 0xaa, 0xd1, 0x23, 0xd6, 0xdf, 0x87, 0xa3, 0x81, 0xb4, 0xcf,
	0xfe, 0xab, 0x0a, 0xc0, 0xd8, 0xe4, 0xad, 0x43, 0x25, 0xc4, 0x3c, 0x9b, 0xb8, 0x7c, 0x62, 0x42,
	0xcc, 0x77, 0xb1, 0x77, 0x05, 0xaa, 0xda, 0x35, 0xc5, 0x85, 0xde, 0x55, 0x14, 0xb8, 0x8b, 0xbd,
	0x6b, 0xd0, 0x10, 0x34, 0x0d, 0x32, 0x8a, 0x3b, 0x77, 0x75, 0x41, 0xd3, 0x9d, 0x31, 0x0b, 0x3f,
	0xcf, 0x59, 0x0b, 0x05, 0x16, 0x7e, 0x6e, 0x59, 0x6b, 0xb0, 0x18, 0xa2, 0xa1, 0x20, 0x7a, 0xb2,
	0xc6, 0xdd, 0x50, 0x90, 0x5a, 0x60, 0x9f, 0xbf, 0xff, 0xd3, 0xdb, 0x81, 0x21, 0x54, 0x5c, 0x07,
	0x0a, 0xdf, 0xd1, 0xa4, 0xff, 0x87, 0xe5, 0x1e, 0x67, 0x71, 0x40, 0xbe, 0x90, 0x24, 0x11, 0x94,
	0x25, 0xad, 0xaa, 0x43, 0x5c, 0x52, 0xb6, 0xfb, 0x99, 0xc9, 0xbb, 0x09, 0x4d, 0xc9, 0x1c, 0x6a,
	0xcd, 0xa1, 0x36, 0x24, 0x1b, 0x13, 0xaf, 0x43, 0x43, 0x7b, 0x4d, 0x86, 0xf1, 0x21, 0xe1, 0xad,
	0xba, 0x3b, 0x47, 0xca, 0xb0, 0xa7, 0x71, 0xb5, 0x58, 0x24, 0xcb, 0x48, 0xe0, 0x90, 0x6a, 0x92,
	0x59, 0xca, 0x26, 0xd4, 0xf0, 0x90, 0xeb, 0xc0, 0x69, 0x35, 0x36, 0xcb, 0xdb, 0x8b, 0x19, 0x23,
	0x43, 0x3d, 0x1f, 0xea, 0x98, 0x72, 0x12, 0x6a, 0x4a, 0xd3, 0xa1, 0x8c, 0x61, 0xef, 0x43, 0x68,
	0xaa, 0x68, 0x0f, 0x84, 0x44, 0x5c, 0x12, 0xdc, 0x5a, 0xd2, 0xab, 0x6b, 0xad, 0x63, 0x24, 0xa1,
	0x93, 0x49, 0x42, 0xe7, 0x69, 0x26, 0x09, 0xfb, 0x0d, 0xc5, 0x3f, 0x30, 0x74, 0xef, 0x2e, 0x2c,
	0xeb, 0xe6, 0x21, 0x4b, 0x12, 0x12, 0x2a, 0x07, 0xcb, 0xa7, 0x3a, 0x58, 0x52, 0x2d, 0x76, 0xb2,
	0x06, 0xde, 0xcf, 0x00, 0xb4, 0x0b, 0x92, 0x60, 0x82, 0x5b, 0x2b, 0xa7, 0x36, 0xaf, 0x2b, 0xf6,
	0x7d, 0x45, 0x56, 0x53, 0xc4, 0x49, 0xc8, 0x38, 0x56, 0x82, 0x93, 0x22, 0x39, 0x68, 0x9d, 0x73,
	0xa7, 0x28, 0xb7, 0x3d, 0x41, 0x72, 0xa0, 0xc4, 0x69, 0x4c, 0x1e, 0xf2, 0xa8, 0x75, 0xde, 0x15,
	0xa7, 0xdc, 0xf4, 0x8c, 0x47, 0x5e, 0x07, 0xce, 0x99, 0xdf, 0x41, 0x8f, 0x46, 0x24, 0x10, 0xf4,
	0x4b, 0xd2, 0xf2, 0x9c, 0xf1, 0xb3, 0x6f, 0x7d, 0x40, 0x23, 0x72, 0x40, 0xbf, 0x24, 0xde, 0x2d,
	0x58, 0x51, 0x4a, 0x1b, 0x38, 0xaa, 0x75, 0xc1, 0x09, 0xf0, 0x25, 0x65, 0xbc, 0x9b, 0x2b, 0xd7,
	0x6d, 0x38, 0x7f, 0x24, 0x69, 0x9f, 0x70, 0x97, 0xff, 0x8e, 0xd3, 0x99, 0x15, 0x63, 0xce, 0x5b,
	0xf8, 0x1f, 0xc0, 0xda, 0x8e, 0x7a, 0x7c, 0x4a, 0xc3, 0xe7, 0x44, 0xde, 0x1b, 0x1d, 0x48, 0x24,
	0x87, 0xb9, 0x34, 0x5c, 0x86, 0x8a, 0xd0, 0x40, 0x21, 0xe2, 0x2c, 0xe6, 0x53, 0xb8, 0xfc, 0x90,
	0xd8, 0x96, 0xa6, 0x9d, 0x76, 0x95, 0x47, 0xfc, 0x2e, 0x34, 0x0d, 0x33, 0xd0, 0x6f, 0xb3, 0x11,
	0xbf, 0xa9, 0x23, 0x7e, 0xe6, 0x4b, 0x4d, 0x3b, 0xfd, 0x96, 0xf2, 0x7e, 0x43, 0x8c, 0x5d, 0xfa,
	0x07, 0xb0, 0xfe, 0x86, 0x16, 0x5e, 0x0b, 0x16, 0xa6, 0x92, 0x82, 0x46, 0x74, 0xac, 0xea, 0x97,
	0xcf, 0x39, 0x83, 0x6c, 0x20, 0xff, 0x21, 0x5c, 0x50, 0x6a, 0xc5, 0x12, 0x89, 0x42, 0x39, 0x76,
	0x76, 0x1b, 0x6a, 0xa1, 0xc5, 0x6c, 0x97, 0x57, 0x6d, 0x97, 0x35, 0x98, 0xf1, 0xf6, 0x73, 0x96,
	0x1f, 0x80, 0x57, 0x70, 0x64, 0x06, 0x6f, 0x42, 0x33, 0xcb, 0xb3, 0x35, 0xd3, 0x91, 0xdf, 0xb9,
	0x93, 0xe4, 0xd7, 0xff, 0xaa, 0x0a, 0xcb, 0xf9, 0xeb, 0xf3, 0x4c, 0x66, 0xdf, 0x1f, 0x24, 0xac,
	0xe0, 0xbc, 0x6e, 0xf1, 0x3d, 0xa6, 0xbe, 0x3e, 0x1d, 0xb0, 0x64, 0x22, 0x93, 0x69, 0x48, 0x05,
	0x79, 0x84, 0x84, 0x4c, 0x50, 0x5c, 0x4c, 0x65, 0x39, 0xaa, 0x66, 0x3f, 0x66, 0x87, 0x34, 0x22,
	0x05, 0x21, 0xb4, 0x98, 0x77, 0x0d, 0x20, 0x22, 0x08, 0x0b, 0x36, 0xe4, 0x61, 0x51, 0x0a, 0x1d,
	0x5c, 0xf5, 0x80, 0xc4, 0x88, 0x46, 0x05, 0x29, 0x34, 0x90, 0x77, 0x03, 0x1a, 0x98, 0x88, 0x90,
	0xd3, 0x54, 0x4e, 0x6a, 0xa0, 0x6b, 0xf0, 0x6e, 0xc1, 0xb2, 0x20, 0x21, 0x4b, 0x30, 0xe2, 0x23,
	0xe3, 0xcc, 0xd5, 0xc0, 0x09, 0x9b, 0xea, 0x57, 0xcc, 0xb0, 0x4a, 0x96, 0xf8, 0x70, 0x54, 0x10,
	0x38, 0x07, 0xd7, 0x2b, 0xdb, 0xf4, 0xbc, 0x51, 0x58, 0xd9, 0xa6, 0xd7, 0x37, 0xa0, 0xa9, 0xa3,
	0x6e, 0x28, 0x88, 0xce, 0x37, 0x4d, 0x27, 0xe4, 0x40, 0x59, 0x9e, 0x09, 0xa2, 0x92, 0x4e, 0x1b,
	0xaa, 0x21, 0x8b, 0x53, 0x94, 0x8c, 0xb4, 0xba, 0xe5, 0xd5, 0x80, 0x05, 0x95, 0xfd, 0x05, 0x39,
	0x14, 0x54, 0x12, 0x2d, 0x5e, 0xb9, 0xdd, 0x82, 0x6a, 0xdd, 0x46, 0x28, 0x21, 0x5a, 0x9a, 0xf2,
	0x75, 0xab, 0x10, 0xbd, 0xa2, 0x55, 0x3d, 0x70, 0xae, 0xb0, 0xa2, 0x55, 0xb2, 0x5f, 0x83, 0x45,
	0x15, 0x19, 0xa4, 0x20, 0x32, 0x06, 0x32, 0xfd, 0x19, 0x26, 0x92, 0x8f, 0xb4, 0xa8, 0x38, 0xfd,
	0xd1, 0xa0, 0xf7, 0x7f, 0xb0, 0xcc, 0x38, 0x26, 0x5c, 0x04, 0x29, 0xe1, 0x01, 0x46, 0x23, 0x2d,
	0x26, 0xb9, 0x52, 0x19, 0xdb, 0x13, 0xc2, 0x3f, 0x46, 0x23, 0xef, 0x03, 0xb8, 0x38, 0x14, 0x04,
	0x07, 0x79, 0xd9, 0x95, 0x72, 0x76, 0x44, 0x31, 0xe1, 0xad, 0x55, 0xa7, 0xcd, 0xaa, 0xe2, 0x1c,
	0x58, 0xca, 0x13, 0xcb, 0xf0, 0x56, 0x61, 0x6e, 0x42, 0x78, 0xe6, 0x28, 0x56, 0x49, 0xa3, 0x47,
	0xb9, 0x5d, 0x72, 0x17, 0xdd, 0x15, 0x9b, 0xc3, 0xde, 0x2f, 0xa0, 0x11, 0x72, 0x82, 0x24, 0xc1,
	0x4a, 0x8b, 0x5b, 0x97, 0x4e, 0xcf, 0x19, 0x0e, 0xdd, 0xfb, 0x08, 0x9a, 0xd9, 0x1c, 0xeb, 0xe6,
	0xad, 0x53, 0x9b, 0x17, 0xf8, 0x4a, 0x9d, 0x91, 0x10, 0xb4, 0x9f, 0x10, 0x9c, 0xcf, 0xfd, 0xbb,
	0xee, 0x5a, 0xcb, 0xac, 0x66, 0xfe, 0xfd, 0x3f, 0x57, 0x61, 0x65, 0x42, 0x16, 0xce, 0x02, 0xf3,
	0x2c, 0x30, 0xcf, 0x02, 0xf3, 0xbf, 0x1d, 0x98, 0x7f, 0xaa, 0x40, 0xe3, 0x11, 0x41, 0xf8, 0x07,
	0x94, 0x2d, 0xa7, 0x83, 0xad, 0xfa, 0x86, 0x60, 0x9b, 0x35, 0x60, 0xb5, 0x93, 0x07, 0x6c, 0x32,
	0xe4, 0xeb, 0x27, 0x85, 0xfc, 0x59, 0x10, 0xff, 0x0f, 0x07, 0xb1, 0xff, 0x55, 0x05, 0x9a, 0x26,
	0x2c, 0x7e, 0x38, 0xc9, 0xea, 0x2c, 0x2e, 0xce, 0xe2, 0xe2, 0xfb, 0xc6, 0xc5, 0x5f, 0xca, 0x70,
	0x3e, 0xdf, 0xca, 0x7e, 0xa7, 0x83, 0xb1, 0xdb, 0x50, 0x91, 0xba, 0x95, 0xdd, 0xbd, 0x79, 0x7a,
	0xab, 0x68, 0x1c, 0x59, 0x3f, 0xd9, 0x2c, 0x1b, 0x9e, 0x77, 0x0b, 0xaa, 0xfa, 0x8b, 0x0f, 0x47,
	0x3a, 0x60, 0x1a, 0x77, 0x9a, 0xba, 0xc9, 0x63, 0x85, 0xdd, 0xcb, 0x4e, 0xf6, 0x32, 0x8a, 0xff,
	0x73, 0xbd, 0xb5, 0xcc, 0x3b, 0x66, 0xc3, 0xf6, 0x3a, 0x54, 0x8d, 0xb7, 0x6c, 0x87, 0xda, 0x70,
	0x5f, 0x9b, 0xd9, 0xfc, 0x5f, 0x42, 0xd5, 0xba, 0x55, 0xf3, 0xd7, 0xa3, 0x24, 0x9a, 0x38, 0x3a,
	0xd3, 0x90, 0x9a, 0x75, 0xc1, 0xb8, 0x2c, 0x84, 0xb7, 0x46, 0xfc, 0x7f, 0x2d, 0xc0, 0xfa, 0x8e,
	0x4e, 0xe3, 0x8f, 0xf9, 0xb3, 0x14, 0x23, 0x49, 0x0a, 0x5f, 0x66, 0x47, 0xbc, 0x3c, 0x31, 0xe2,
	0xd9, 0x6e, 0x7c, 0x6e, 0xd6, 0x6e, 0xdc, 0x9c, 0xac, 0xce, 0x4f, 0x9f, 0xac, 0xae, 0xc1, 0xe2,
	0x11, 0x8a, 0x86, 0x45, 0x99, 0x30, 0x90, 0xb7, 0x05, 0x75, 0x16, 0xe1, 0xc0, 0xd8, 0x5d, 0x91,
	0xa8, 0xb1, 0x08, 0x7f, 0xa6, 0x29, 0x97, 0xa1, 0xc2, 0x09, 0x12, 0x2c, 0x29, 0x68, 0x84, 0xc5,
	0xa6, 0x42, 0xab, 0x7a, 0x42, 0x68, 0x6d, 0x40, 0x4d, 0x8f, 0x7c, 0x26, 0x0b, 0xf3, 0x85, 0xf9,
	0xd8, 0xc5, 0x4a, 0x30, 0x0d, 0x41, 0x7f, 0xa1, 0x2b, 0x07, 0x75, 0x8d, 0xef, 0xa8, 0xcf, 0xdc,
	0x80, 0x5a, 0xaf, 0x17, 0x1b, 0x8a, 0x2b, 0x05, 0xd5, 0x5e, 0x2f, 0xd6, 0x84, 0x2b, 0xa0, 0x1e,
	0xf5, 0x51, 0x51, 0x41, 0x08, 0x7a, 0xbd, 0xf8, 0x19, 0x8f, 0xbc, 0x75, 0x50, 0x4f, 0x93, 0x12,
	0xb0, 0xd8, 0xeb, 0xc5, 0x6f, 0x11, 0xfd, 0x9b, 0x50, 0xcb, 0x63, 0xc6, 0x8d, 0xf0, 0x1c, 0x55,
	0xf3, 0x93, 0x30, 0x49, 0x8a, 0x51, 0xae, 0x10, 0xa5, 0x76, 0x24, 0x39, 0xa2, 0x9c, 0x25, 0x31,
	0x49, 0x64, 0x21, 0xd6, 0x5d, 0x83, 0xfe, 0x40, 0xce, 0xe2, 0x00, 0xa5, 0x69, 0x31, 0xe4, 0x15,
	0x7a, 0x37, 0x4d, 0xd5, 0x22, 0xb7, 0x27, 0x50, 0x3a, 0xd6, 0xb3, 0x45, 0x6e, 0xcf, 0x9d, 0x32,
	0xb6, 0xa5, 0xf8, 0x5f, 0x2f, 0xc0, 0xd2, 0xd9, 0xc2, 0x3a, 0x5b, 0x58, 0xa7, 0x2f, 0x2c, 0xff,
	0xef, 0x65, 0xa8, 0xda, 0x55, 0xe4, 0x2c, 0x92, 0x79, 0x67, 0x91, 0x6c, 0x41, 0xbd, 0x37, 0x8c,
	0xa2, 0x40, 0xeb, 0xbd, 0xbb, 0x52, 0x6a, 0x0a, 0xde, 0x53, 0x25, 0xc9, 0x55, 0x00, 0x31, 0x60,
	0x5c, 0x06, 0x53, 0x65, 0x4b, 0x5d, 0xe3, 0x9a, 0x94, 0x57, 0x3d, 0x0b, 0xd3, 0x55, 0x4f, 0x5e,
	0x8f, 0x2c, 0x4e, 0xd7, 0x23, 0xce, 0x30, 0x56, 0x66, 0x0c, 0xa3, 0xff, 0x8f, 0x06, 0x54, 0xcc,
	0x62, 0xd7, 0x47, 0xf1, 0xfa, 0x69, 0xb2, 0xf8, 0xaa, 0x19, 0x78, 0x8f, 0xcd, 0xac, 0x57, 0xe6,
	0xde, 0x50, 0xaf, 0x6c, 0x41, 0x3d, 0x45, 0x9c, 0x98, 0x93, 0xdf, 0x79, 0xf7, 0x2a, 0xc8, 0xc0,
	0xfa, 0x90, 0xf8, 0x9c, 0x71, 0x9f, 0x72, 0xca, 0x38, 0x95, 0x94, 0x88, 0xc2, 0x37, 0x4e, 0x59,
	0xd5, 0x78, 0xa5, 0x9c, 0xe1, 0x61, 0xa8, 0xbd, 0x2e, 0xba, 0xb7, 0x66, 0x16, 0x77, 0xdd, 0x0a,
	0x72, 0x44, 0xac, 0xdb, 0xca, 0xb4, 0xdb, 0xb1, 0xd5, 0xdb, 0x86, 0xa6, 0xc5, 0xcc, 0x19, 0xb3,
	0x5b, 0xb7, 0x15, 0x2c, 0x63, 0xdf, 0x21, 0x92, 0xa4, 0xaf, 0xef, 0x95, 0x0a, 0x55, 0xdb, 0x94,
	0x55, 0x75, 0x79, 0xa8, 0x13, 0x96, 0xbe, 0x6b, 0x2a, 0x84, 0x93, 0xc1, 0x1f, 0xb1, 0xbe, 0x9a,
	0xc6, 0x01, 0x1b, 0x72, 0x51, 0x88, 0x25, 0x03, 0xa9, 0xb5, 0x8c, 0xd1, 0x48, 0x14, 0xc2, 0x48,
	0x23, 0x93, 0x3b, 0xe1, 0xe6, 0xf7, 0xdb, 0x09, 0x2f, 0x7d, 0xc7, 0x9d, 0x70, 0x76, 0xcb, 0x93,
	0x32, 0x2e, 0x51, 0x54, 0x28, 0xf0, 0xf4, 0x2d, 0xcf, 0x13, 0x8d, 0x4f, 0xd4, 0xa3, 0x2b, 0x27,
	0xd4, 0xa3, 0x37, 0xb3, 0x19, 0x08, 0x8c, 0x7a, 0x9e, 0x2b, 0xdc, 0x2d, 0x69, 0xcb, 0x53, 0xad,
	0xa1, 0x13, 0x65, 0xf0, 0xf9, 0x93, 0xca, 0xe0, 0x4d, 0xa8, 0x09, 0x16, 0x0d, 0x35, 0xc9, 0x8d,
	0xdf, 0x1c, 0x75, 0x77, 0x1c, 0xf6, 0x2e, 0x63, 0x72, 0xc7, 0xb1, 0x8b, 0x9d, 0x3a, 0x79, 0x75,
	0x46, 0x9d, 0xdc, 0x86, 0xaa, 0x90, 0x88, 0x73, 0x52, 0x2c, 0xf6, 0x32, 0x50, 0x4d, 0x9d, 0x44,
	0x7d, 0x51, 0x28, 0xf6, 0x34, 0x92, 0x0b, 0x54, 0x6b, 0x4a, 0xa0, 0x5c, 0x65, 0x7d, 0xf7, 0x14,
	0x65, 0x5d, 0x7b, 0xa3, 0xb2, 0xae, 0x4f, 0x2b, 0xeb, 0x64, 0x92, 0xb8, 0xfc, 0x16, 0x49, 0xe2,
	0xca, 0xe9, 0x49, 0xa2, 0x3d, 0x3b, 0x49, 0x38, 0x02, 0xb4, 0x71, 0x9a, 0x8e, 0x6f, 0xce, 0xd4,
	0x71, 0x57, 0x85, 0xb7, 0x66, 0xa5, 0xf7, 0x09, 0x39, 0xf7, 0x4f, 0x92, 0xf3, 0x2c, 0x61, 0x5f,
	0x9d, 0x4a, 0xd8, 0x85, 0xc4, 0x7b, 0x6d, 0x66, 0xe2, 0xdd, 0x82, 0x7a, 0x42, 0x5e, 0x58, 0xca,
	0x75, 0x97, 0x92, 0x90, 0x17, 0x86, 0xe2, 0x43, 0x5d, 0x0c, 0x0f, 0xad, 0x7c, 0xdc, 0x28, 0xe8,
	0x78, 0x06, 0xeb, 0xcb, 0x71, 0x35, 0xf8, 0x7a, 0xe2, 0x6f, 0xba, 0x6e, 0x14, 0xbc, 0xa7, 0x26,
	0x7f, 0x9c, 0xe2, 0xb7, 0x67, 0xa4, 0x78, 0x93, 0x66, 0xde, 0x2b, 0xd6, 0x22, 0xfe, 0x03, 0x78,
	0x47, 0x5f, 0x75, 0x67, 0x8a, 0x93, 0xd7, 0xe6, 0x3f, 0x02, 0x70, 0x54, 0xca, 0x94, 0xe7, 0x4b,
	0xe6, 0x02, 0xc9, 0xc0, 0xa3, 0x7d, 0x87, 0xe0, 0xff, 0x0a, 0x6a, 0x19, 0xfe, 0xe6, 0x6b, 0xac,
	0x08, 0x1d, 0x92, 0xa8, 0xb8, 0x05, 0xd7, 0x90, 0xff, 0xbe, 0xde, 0xbb, 0x4c, 0xdc, 0x88, 0x5d,
	0x75, 0x6e, 0xee, 0xc6, 0x1b, 0x04, 0x4b, 0xca, 0x2e, 0xf0, 0x3e, 0x82, 0x8a, 0x41, 0xde, 0xf6,
	0x02, 0xad, 0x5e, 0xb8, 0x40, 0xbb, 0xf7, 0xe9, 0xcb, 0xd7, 0xed, 0xf2, 0x37, 0xaf, 0xdb, 0xe5,
	0x57, 0xaf, 0xdb, 0xa5, 0x7f, 0xbf, 0x6e, 0x97, 0x7e, 0x7f, 0xdc, 0x2e, 0x7d, 0x7d, 0xdc, 0x2e,
	0xfd, 0xed, 0xb8, 0x5d, 0x7a, 0x79, 0xdc, 0x2e, 0x7d, 0x73, 0xdc, 0x2e, 0xbd, 0x3a, 0x6e, 0x97,
	0xfe, 0xf0, 0x6d, 0xbb, 0xf4, 0xd7, 0x6f, 0xdb, 0xa5, 0xdf, 0xae, 0xab, 0x31, 0xef, 0x1c, 0x25,
	0x5d, 0x94, 0xd2, 0x6e, 0x7a, 0xd8, 0x75, 0xff, 0x37, 0xca, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x5d, 0x6b, 0xab, 0xe2, 0x9c, 0x22, 0x00, 0x00,
}
