// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: common/common.proto

package common

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type VersionInfoResponse struct {
	Service              string               `protobuf:"bytes,1,opt,name=service" json:"service"`
	Version              string               `protobuf:"bytes,2,opt,name=version" json:"version"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *VersionInfoResponse) Reset()         { *m = VersionInfoResponse{} }
func (m *VersionInfoResponse) String() string { return proto.CompactTextString(m) }
func (*VersionInfoResponse) ProtoMessage()    {}
func (*VersionInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{1}
}
func (m *VersionInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionInfoResponse.Unmarshal(m, b)
}
func (m *VersionInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionInfoResponse.Marshal(b, m, deterministic)
}
func (m *VersionInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionInfoResponse.Merge(m, src)
}
func (m *VersionInfoResponse) XXX_Size() int {
	return xxx_messageInfo_VersionInfoResponse.Size(m)
}
func (m *VersionInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionInfoResponse proto.InternalMessageInfo

func (m *VersionInfoResponse) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *VersionInfoResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *VersionInfoResponse) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type IDRequest struct {
	// @required
	Id                   int64    `protobuf:"varint,1,opt,name=id" json:"id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDRequest) Reset()         { *m = IDRequest{} }
func (m *IDRequest) String() string { return proto.CompactTextString(m) }
func (*IDRequest) ProtoMessage()    {}
func (*IDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{2}
}
func (m *IDRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDRequest.Unmarshal(m, b)
}
func (m *IDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDRequest.Marshal(b, m, deterministic)
}
func (m *IDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDRequest.Merge(m, src)
}
func (m *IDRequest) XXX_Size() int {
	return xxx_messageInfo_IDRequest.Size(m)
}
func (m *IDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDRequest proto.InternalMessageInfo

func (m *IDRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type CodeRequest struct {
	// @required
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CodeRequest) Reset()         { *m = CodeRequest{} }
func (m *CodeRequest) String() string { return proto.CompactTextString(m) }
func (*CodeRequest) ProtoMessage()    {}
func (*CodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{3}
}
func (m *CodeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CodeRequest.Unmarshal(m, b)
}
func (m *CodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CodeRequest.Marshal(b, m, deterministic)
}
func (m *CodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CodeRequest.Merge(m, src)
}
func (m *CodeRequest) XXX_Size() int {
	return xxx_messageInfo_CodeRequest.Size(m)
}
func (m *CodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CodeRequest proto.InternalMessageInfo

func (m *CodeRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type NameRequest struct {
	// @required
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NameRequest) Reset()         { *m = NameRequest{} }
func (m *NameRequest) String() string { return proto.CompactTextString(m) }
func (*NameRequest) ProtoMessage()    {}
func (*NameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{4}
}
func (m *NameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameRequest.Unmarshal(m, b)
}
func (m *NameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameRequest.Marshal(b, m, deterministic)
}
func (m *NameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameRequest.Merge(m, src)
}
func (m *NameRequest) XXX_Size() int {
	return xxx_messageInfo_NameRequest.Size(m)
}
func (m *NameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NameRequest proto.InternalMessageInfo

func (m *NameRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type IDsRequest struct {
	// @required
	Ids                  []int64  `protobuf:"varint,1,rep,name=ids" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDsRequest) Reset()         { *m = IDsRequest{} }
func (m *IDsRequest) String() string { return proto.CompactTextString(m) }
func (*IDsRequest) ProtoMessage()    {}
func (*IDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{5}
}
func (m *IDsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDsRequest.Unmarshal(m, b)
}
func (m *IDsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDsRequest.Marshal(b, m, deterministic)
}
func (m *IDsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDsRequest.Merge(m, src)
}
func (m *IDsRequest) XXX_Size() int {
	return xxx_messageInfo_IDsRequest.Size(m)
}
func (m *IDsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDsRequest proto.InternalMessageInfo

func (m *IDsRequest) GetIds() []int64 {
	if m != nil {
		return m.Ids
	}
	return nil
}

type StatusResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status" json:"status"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{6}
}
func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (m *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(m, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type IDMRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id" json:"id"`
	Paging               *Paging  `protobuf:"bytes,2,opt,name=paging" json:"paging,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDMRequest) Reset()         { *m = IDMRequest{} }
func (m *IDMRequest) String() string { return proto.CompactTextString(m) }
func (*IDMRequest) ProtoMessage()    {}
func (*IDMRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{7}
}
func (m *IDMRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDMRequest.Unmarshal(m, b)
}
func (m *IDMRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDMRequest.Marshal(b, m, deterministic)
}
func (m *IDMRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDMRequest.Merge(m, src)
}
func (m *IDMRequest) XXX_Size() int {
	return xxx_messageInfo_IDMRequest.Size(m)
}
func (m *IDMRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDMRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDMRequest proto.InternalMessageInfo

func (m *IDMRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *IDMRequest) GetPaging() *Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

type Paging struct {
	Offset               int32    `protobuf:"varint,5,opt,name=offset" json:"offset"`
	Limit                int32    `protobuf:"varint,6,opt,name=limit" json:"limit"`
	Sort                 string   `protobuf:"bytes,7,opt,name=sort" json:"sort"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Paging) Reset()         { *m = Paging{} }
func (m *Paging) String() string { return proto.CompactTextString(m) }
func (*Paging) ProtoMessage()    {}
func (*Paging) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{8}
}
func (m *Paging) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Paging.Unmarshal(m, b)
}
func (m *Paging) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Paging.Marshal(b, m, deterministic)
}
func (m *Paging) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Paging.Merge(m, src)
}
func (m *Paging) XXX_Size() int {
	return xxx_messageInfo_Paging.Size(m)
}
func (m *Paging) XXX_DiscardUnknown() {
	xxx_messageInfo_Paging.DiscardUnknown(m)
}

var xxx_messageInfo_Paging proto.InternalMessageInfo

func (m *Paging) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *Paging) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *Paging) GetSort() string {
	if m != nil {
		return m.Sort
	}
	return ""
}

type ForwardPaging struct {
	Since                string   `protobuf:"bytes,1,opt,name=since" json:"since"`
	Limit                int32    `protobuf:"varint,2,opt,name=limit" json:"limit"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardPaging) Reset()         { *m = ForwardPaging{} }
func (m *ForwardPaging) String() string { return proto.CompactTextString(m) }
func (*ForwardPaging) ProtoMessage()    {}
func (*ForwardPaging) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{9}
}
func (m *ForwardPaging) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardPaging.Unmarshal(m, b)
}
func (m *ForwardPaging) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardPaging.Marshal(b, m, deterministic)
}
func (m *ForwardPaging) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardPaging.Merge(m, src)
}
func (m *ForwardPaging) XXX_Size() int {
	return xxx_messageInfo_ForwardPaging.Size(m)
}
func (m *ForwardPaging) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardPaging.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardPaging proto.InternalMessageInfo

func (m *ForwardPaging) GetSince() string {
	if m != nil {
		return m.Since
	}
	return ""
}

func (m *ForwardPaging) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type PageInfo struct {
	Total                int32    `protobuf:"varint,1,opt,name=total" json:"total"`
	Limit                int32    `protobuf:"varint,2,opt,name=limit" json:"limit"`
	Sort                 []string `protobuf:"bytes,3,rep,name=sort" json:"sort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PageInfo) Reset()         { *m = PageInfo{} }
func (m *PageInfo) String() string { return proto.CompactTextString(m) }
func (*PageInfo) ProtoMessage()    {}
func (*PageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{10}
}
func (m *PageInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageInfo.Unmarshal(m, b)
}
func (m *PageInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageInfo.Marshal(b, m, deterministic)
}
func (m *PageInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageInfo.Merge(m, src)
}
func (m *PageInfo) XXX_Size() int {
	return xxx_messageInfo_PageInfo.Size(m)
}
func (m *PageInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PageInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PageInfo proto.InternalMessageInfo

func (m *PageInfo) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *PageInfo) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *PageInfo) GetSort() []string {
	if m != nil {
		return m.Sort
	}
	return nil
}

type ForwardPageInfo struct {
	Since                string   `protobuf:"bytes,1,opt,name=since" json:"since"`
	Limit                int32    `protobuf:"varint,2,opt,name=limit" json:"limit"`
	HasNext              bool     `protobuf:"varint,4,opt,name=has_next,json=hasNext" json:"has_next"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardPageInfo) Reset()         { *m = ForwardPageInfo{} }
func (m *ForwardPageInfo) String() string { return proto.CompactTextString(m) }
func (*ForwardPageInfo) ProtoMessage()    {}
func (*ForwardPageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{11}
}
func (m *ForwardPageInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardPageInfo.Unmarshal(m, b)
}
func (m *ForwardPageInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardPageInfo.Marshal(b, m, deterministic)
}
func (m *ForwardPageInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardPageInfo.Merge(m, src)
}
func (m *ForwardPageInfo) XXX_Size() int {
	return xxx_messageInfo_ForwardPageInfo.Size(m)
}
func (m *ForwardPageInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardPageInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardPageInfo proto.InternalMessageInfo

func (m *ForwardPageInfo) GetSince() string {
	if m != nil {
		return m.Since
	}
	return ""
}

func (m *ForwardPageInfo) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ForwardPageInfo) GetHasNext() bool {
	if m != nil {
		return m.HasNext
	}
	return false
}

type RawJSONObject struct {
	Data                 []byte   `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawJSONObject) Reset()         { *m = RawJSONObject{} }
func (m *RawJSONObject) String() string { return proto.CompactTextString(m) }
func (*RawJSONObject) ProtoMessage()    {}
func (*RawJSONObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{12}
}
func (m *RawJSONObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawJSONObject.Unmarshal(m, b)
}
func (m *RawJSONObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawJSONObject.Marshal(b, m, deterministic)
}
func (m *RawJSONObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawJSONObject.Merge(m, src)
}
func (m *RawJSONObject) XXX_Size() int {
	return xxx_messageInfo_RawJSONObject.Size(m)
}
func (m *RawJSONObject) XXX_DiscardUnknown() {
	xxx_messageInfo_RawJSONObject.DiscardUnknown(m)
}

var xxx_messageInfo_RawJSONObject proto.InternalMessageInfo

func (m *RawJSONObject) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Error struct {
	Code                 string            `protobuf:"bytes,1,opt,name=code" json:"code"`
	Msg                  string            `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Meta                 map[string]string `protobuf:"bytes,3,rep,name=meta" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{13}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Error) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *Error) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

type ErrorsResponse struct {
	Errors               []*Error `protobuf:"bytes,1,rep,name=errors" json:"errors,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorsResponse) Reset()         { *m = ErrorsResponse{} }
func (m *ErrorsResponse) String() string { return proto.CompactTextString(m) }
func (*ErrorsResponse) ProtoMessage()    {}
func (*ErrorsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{14}
}
func (m *ErrorsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorsResponse.Unmarshal(m, b)
}
func (m *ErrorsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorsResponse.Marshal(b, m, deterministic)
}
func (m *ErrorsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorsResponse.Merge(m, src)
}
func (m *ErrorsResponse) XXX_Size() int {
	return xxx_messageInfo_ErrorsResponse.Size(m)
}
func (m *ErrorsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorsResponse proto.InternalMessageInfo

func (m *ErrorsResponse) GetErrors() []*Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

type UpdatedResponse struct {
	Updated              int32    `protobuf:"varint,1,opt,name=updated" json:"updated"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatedResponse) Reset()         { *m = UpdatedResponse{} }
func (m *UpdatedResponse) String() string { return proto.CompactTextString(m) }
func (*UpdatedResponse) ProtoMessage()    {}
func (*UpdatedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{15}
}
func (m *UpdatedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatedResponse.Unmarshal(m, b)
}
func (m *UpdatedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatedResponse.Marshal(b, m, deterministic)
}
func (m *UpdatedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatedResponse.Merge(m, src)
}
func (m *UpdatedResponse) XXX_Size() int {
	return xxx_messageInfo_UpdatedResponse.Size(m)
}
func (m *UpdatedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatedResponse proto.InternalMessageInfo

func (m *UpdatedResponse) GetUpdated() int32 {
	if m != nil {
		return m.Updated
	}
	return 0
}

type RemovedResponse struct {
	Removed              int32    `protobuf:"varint,1,opt,name=removed" json:"removed"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemovedResponse) Reset()         { *m = RemovedResponse{} }
func (m *RemovedResponse) String() string { return proto.CompactTextString(m) }
func (*RemovedResponse) ProtoMessage()    {}
func (*RemovedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{16}
}
func (m *RemovedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemovedResponse.Unmarshal(m, b)
}
func (m *RemovedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemovedResponse.Marshal(b, m, deterministic)
}
func (m *RemovedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemovedResponse.Merge(m, src)
}
func (m *RemovedResponse) XXX_Size() int {
	return xxx_messageInfo_RemovedResponse.Size(m)
}
func (m *RemovedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemovedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemovedResponse proto.InternalMessageInfo

func (m *RemovedResponse) GetRemoved() int32 {
	if m != nil {
		return m.Removed
	}
	return 0
}

type DeletedResponse struct {
	Deleted              int32    `protobuf:"varint,1,opt,name=deleted" json:"deleted"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletedResponse) Reset()         { *m = DeletedResponse{} }
func (m *DeletedResponse) String() string { return proto.CompactTextString(m) }
func (*DeletedResponse) ProtoMessage()    {}
func (*DeletedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{17}
}
func (m *DeletedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletedResponse.Unmarshal(m, b)
}
func (m *DeletedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletedResponse.Marshal(b, m, deterministic)
}
func (m *DeletedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletedResponse.Merge(m, src)
}
func (m *DeletedResponse) XXX_Size() int {
	return xxx_messageInfo_DeletedResponse.Size(m)
}
func (m *DeletedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeletedResponse proto.InternalMessageInfo

func (m *DeletedResponse) GetDeleted() int32 {
	if m != nil {
		return m.Deleted
	}
	return 0
}

type MessageResponse struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageResponse) Reset()         { *m = MessageResponse{} }
func (m *MessageResponse) String() string { return proto.CompactTextString(m) }
func (*MessageResponse) ProtoMessage()    {}
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{18}
}
func (m *MessageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageResponse.Unmarshal(m, b)
}
func (m *MessageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageResponse.Marshal(b, m, deterministic)
}
func (m *MessageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageResponse.Merge(m, src)
}
func (m *MessageResponse) XXX_Size() int {
	return xxx_messageInfo_MessageResponse.Size(m)
}
func (m *MessageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MessageResponse proto.InternalMessageInfo

func (m *MessageResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *MessageResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type Filter struct {
	// Comma separated properties: "name,s_name"
	Name string `protobuf:"bytes,1,opt,name=name" json:"name"`
	// Can be = ≠ (!=) < ≤ (<=) > ≥ (>=) ⊃ (c) ∈ (in) ∩ (n)
	//
	// - Text or set: ⊃ ∩
	// - Exactly: = ≠ ∈
	// - Numeric: = ≠ ∈ < ≤ > ≥
	Op string `protobuf:"bytes,2,opt,name=op" json:"op"`
	// Must always be string
	Value                string   `protobuf:"bytes,3,opt,name=value" json:"value"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{19}
}
func (m *Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Filter.Unmarshal(m, b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
}
func (m *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(m, src)
}
func (m *Filter) XXX_Size() int {
	return xxx_messageInfo_Filter.Size(m)
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (m *Filter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Filter) GetOp() string {
	if m != nil {
		return m.Op
	}
	return ""
}

func (m *Filter) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*Empty)(nil), "cm.Empty")
	proto.RegisterType((*VersionInfoResponse)(nil), "cm.VersionInfoResponse")
	proto.RegisterType((*IDRequest)(nil), "cm.IDRequest")
	proto.RegisterType((*CodeRequest)(nil), "cm.CodeRequest")
	proto.RegisterType((*NameRequest)(nil), "cm.NameRequest")
	proto.RegisterType((*IDsRequest)(nil), "cm.IDsRequest")
	proto.RegisterType((*StatusResponse)(nil), "cm.StatusResponse")
	proto.RegisterType((*IDMRequest)(nil), "cm.IDMRequest")
	proto.RegisterType((*Paging)(nil), "cm.Paging")
	proto.RegisterType((*ForwardPaging)(nil), "cm.ForwardPaging")
	proto.RegisterType((*PageInfo)(nil), "cm.PageInfo")
	proto.RegisterType((*ForwardPageInfo)(nil), "cm.ForwardPageInfo")
	proto.RegisterType((*RawJSONObject)(nil), "cm.RawJSONObject")
	proto.RegisterType((*Error)(nil), "cm.Error")
	proto.RegisterMapType((map[string]string)(nil), "cm.Error.MetaEntry")
	proto.RegisterType((*ErrorsResponse)(nil), "cm.ErrorsResponse")
	proto.RegisterType((*UpdatedResponse)(nil), "cm.UpdatedResponse")
	proto.RegisterType((*RemovedResponse)(nil), "cm.RemovedResponse")
	proto.RegisterType((*DeletedResponse)(nil), "cm.DeletedResponse")
	proto.RegisterType((*MessageResponse)(nil), "cm.MessageResponse")
	proto.RegisterType((*Filter)(nil), "cm.Filter")
}

func init() { proto.RegisterFile("common/common.proto", fileDescriptor_8f954d82c0b891f6) }

var fileDescriptor_8f954d82c0b891f6 = []byte{
	// 789 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x95, 0x4d, 0x6e, 0xdb, 0x46,
	0x14, 0xc7, 0xf9, 0x61, 0xc9, 0xd6, 0x53, 0x13, 0x07, 0xb4, 0x11, 0xb0, 0x42, 0x41, 0x33, 0xd3,
	0x45, 0xb4, 0x68, 0x48, 0x44, 0x5d, 0xb4, 0xcd, 0xce, 0xf1, 0x47, 0xe1, 0x02, 0x76, 0x0c, 0xc6,
	0xcd, 0xa2, 0x28, 0x10, 0x8c, 0xc9, 0x27, 0x86, 0x89, 0xc8, 0x61, 0x39, 0x23, 0x39, 0xbe, 0x41,
	0x97, 0x41, 0x57, 0x05, 0x7a, 0x81, 0x1e, 0xa1, 0x47, 0xf0, 0x32, 0x27, 0x28, 0x22, 0xe5, 0x02,
	0x3d, 0x42, 0x31, 0x33, 0xa4, 0x2c, 0xba, 0x48, 0x80, 0x66, 0xe5, 0xc7, 0xdf, 0xfb, 0xfa, 0xcf,
	0xcc, 0x7b, 0x16, 0x6c, 0xc5, 0x2c, 0xcf, 0x59, 0x11, 0xea, 0x3f, 0x41, 0x59, 0x31, 0xc1, 0x1c,
	0x2b, 0xce, 0x07, 0x5f, 0x29, 0x33, 0x7e, 0x90, 0x62, 0xf1, 0x80, 0x5f, 0xd0, 0x34, 0xc5, 0x2a,
	0x64, 0xa5, 0xc8, 0x58, 0xc1, 0x43, 0x5a, 0x14, 0x4c, 0x50, 0x65, 0xeb, 0x8c, 0xc1, 0x76, 0xca,
	0x52, 0xa6, 0xcc, 0x50, 0x5a, 0x35, 0xdd, 0x49, 0x19, 0x4b, 0x27, 0x18, 0xaa, 0xaf, 0xf3, 0xe9,
	0x38, 0x14, 0x59, 0x8e, 0x5c, 0xd0, 0xbc, 0xd4, 0x01, 0x64, 0x1d, 0x3a, 0x07, 0x79, 0x29, 0x2e,
	0xc9, 0x1b, 0x13, 0xb6, 0x9e, 0x61, 0xc5, 0x33, 0x56, 0x1c, 0x15, 0x63, 0x16, 0x21, 0x2f, 0x59,
	0xc1, 0xd1, 0xf1, 0x60, 0x9d, 0x63, 0x35, 0xcb, 0x62, 0x74, 0x4d, 0xdf, 0x1c, 0xf6, 0x1e, 0xaf,
	0x5d, 0xfd, 0xbd, 0x63, 0x44, 0x0d, 0x94, 0xfe, 0x99, 0x4e, 0x73, 0xad, 0x55, 0x7f, 0x0d, 0x9d,
	0xef, 0x00, 0xa6, 0x65, 0x42, 0x05, 0x26, 0xcf, 0xa9, 0x70, 0x6d, 0xdf, 0x1c, 0xf6, 0x47, 0x83,
	0x40, 0xcb, 0x0a, 0x1a, 0x59, 0xc1, 0x59, 0x23, 0x2b, 0xea, 0xd5, 0xd1, 0xbb, 0x82, 0xdc, 0x83,
	0xde, 0xd1, 0x7e, 0x84, 0xbf, 0x4c, 0x91, 0x0b, 0x67, 0x1b, 0xac, 0x2c, 0x51, 0x12, 0xec, 0xba,
	0x85, 0x95, 0x25, 0xe4, 0x3e, 0xf4, 0xf7, 0x58, 0x82, 0x4d, 0x90, 0x0b, 0x6b, 0x31, 0x4b, 0xda,
	0x4a, 0x15, 0x91, 0x81, 0x27, 0x34, 0x5f, 0x0d, 0x2c, 0x68, 0x7e, 0x23, 0x50, 0x12, 0xe2, 0x01,
	0x1c, 0xed, 0xf3, 0x26, 0xee, 0x0e, 0xd8, 0x59, 0xc2, 0x5d, 0xd3, 0xb7, 0x87, 0x76, 0x24, 0x4d,
	0x12, 0xc0, 0xed, 0xa7, 0x82, 0x8a, 0x29, 0x5f, 0xde, 0xd0, 0x17, 0xd0, 0xe5, 0x8a, 0xb4, 0xaa,
	0xd5, 0x8c, 0x1c, 0xca, 0x7a, 0xc7, 0x1f, 0x3d, 0x85, 0x43, 0xa0, 0x5b, 0xd2, 0x34, 0x2b, 0x52,
	0x75, 0x85, 0xfd, 0x11, 0x04, 0x71, 0x1e, 0x9c, 0x2a, 0x12, 0xd5, 0x1e, 0xf2, 0x33, 0x74, 0x35,
	0x91, 0xfd, 0xd8, 0x78, 0xcc, 0x51, 0xb8, 0x1d, 0xdf, 0x1c, 0x76, 0x9a, 0x7e, 0x9a, 0x39, 0x03,
	0xe8, 0x4c, 0xb2, 0x3c, 0x13, 0x6e, 0x77, 0xc5, 0xa9, 0x91, 0x3c, 0x35, 0x67, 0x95, 0x70, 0xd7,
	0x57, 0x4f, 0x2d, 0x09, 0xf9, 0x1e, 0x6e, 0x1d, 0xb2, 0xea, 0x82, 0x56, 0x49, 0xdd, 0x64, 0x00,
	0x1d, 0x9e, 0x15, 0x37, 0x1e, 0x5d, 0xa3, 0xeb, 0x16, 0xd6, 0x7f, 0x5a, 0x90, 0x67, 0xb0, 0x71,
	0x4a, 0x53, 0x94, 0x23, 0x24, 0xe3, 0x04, 0x13, 0x74, 0xa2, 0x6a, 0x2c, 0xe3, 0x14, 0xfa, 0x58,
	0x0d, 0xc7, 0xa9, 0x65, 0xda, 0xbe, 0x3d, 0xec, 0xd5, 0x02, 0x5f, 0xc2, 0xe6, 0xb5, 0xc0, 0x65,
	0xf9, 0x4f, 0x91, 0xe8, 0xec, 0xc0, 0xc6, 0x0b, 0xca, 0x9f, 0x17, 0xf8, 0x5a, 0xb8, 0x6b, 0xbe,
	0x39, 0xdc, 0x68, 0x46, 0xf6, 0x05, 0xe5, 0x27, 0xf8, 0x5a, 0x90, 0x2f, 0xe1, 0x56, 0x44, 0x2f,
	0x7e, 0x78, 0xfa, 0xe4, 0xe4, 0xc9, 0xf9, 0x4b, 0x8c, 0x95, 0xa0, 0x84, 0x0a, 0xea, 0x9a, 0xbe,
	0x35, 0xfc, 0x2c, 0x52, 0x36, 0xf9, 0xc3, 0x84, 0xce, 0x41, 0x55, 0xb1, 0xea, 0xc3, 0x43, 0xe7,
	0xdc, 0x05, 0x3b, 0xe7, 0x69, 0x6b, 0x2f, 0x24, 0x70, 0xee, 0xc3, 0x5a, 0x8e, 0x82, 0xaa, 0x03,
	0xf6, 0x47, 0x5b, 0xf2, 0xb5, 0x55, 0xa9, 0xe0, 0x18, 0x05, 0x3d, 0x28, 0x44, 0x75, 0x19, 0xa9,
	0x80, 0xc1, 0x37, 0xd0, 0x5b, 0x22, 0x39, 0x8b, 0xaf, 0xf0, 0x52, 0xb7, 0x89, 0xa4, 0xe9, 0x6c,
	0x43, 0x67, 0x46, 0x27, 0x53, 0xd4, 0x1d, 0x22, 0xfd, 0xf1, 0xc8, 0xfa, 0xd6, 0x24, 0x5f, 0xc3,
	0x6d, 0x55, 0xf1, 0x7a, 0x4a, 0xef, 0x41, 0x17, 0x15, 0x51, 0xc3, 0xdc, 0x1f, 0xf5, 0x96, 0x5d,
	0xa3, 0xda, 0x41, 0x1e, 0xc2, 0xe6, 0x8f, 0x7a, 0xf9, 0x56, 0xb7, 0xbf, 0xde, 0xc7, 0xd6, 0x23,
	0x36, 0x50, 0xa6, 0x44, 0x98, 0xb3, 0x59, 0x3b, 0xa5, 0xd2, 0xa8, 0x9d, 0x52, 0x43, 0x99, 0xb2,
	0x8f, 0x13, 0xbc, 0xd1, 0x25, 0xd1, 0xa8, 0x9d, 0x52, 0x43, 0xb2, 0x07, 0x9b, 0xc7, 0xc8, 0x39,
	0x4d, 0x71, 0x99, 0xf2, 0xbf, 0x2f, 0x9d, 0x9c, 0x41, 0xf7, 0x30, 0x9b, 0x08, 0xac, 0x3e, 0xbc,
	0xfc, 0x72, 0x3d, 0x59, 0xd9, 0x4a, 0xb5, 0x58, 0x29, 0x87, 0x49, 0x5f, 0xb3, 0xbd, 0x3a, 0x68,
	0x0a, 0x3d, 0x1e, 0xff, 0xb6, 0xeb, 0x3a, 0x77, 0xe1, 0x0e, 0x9e, 0xb1, 0xd2, 0xdf, 0x3d, 0x3d,
	0x7a, 0xe4, 0xef, 0xa9, 0xff, 0xe3, 0x23, 0x6b, 0xf6, 0xf0, 0x6a, 0xee, 0x99, 0x6f, 0xe7, 0x9e,
	0xf1, 0x6e, 0xee, 0x19, 0xff, 0xcc, 0x3d, 0xe3, 0xd7, 0x85, 0x67, 0xfc, 0xb9, 0xf0, 0x8c, 0xbf,
	0x16, 0x9e, 0x71, 0xb5, 0xf0, 0x8c, 0xb7, 0x0b, 0xcf, 0x78, 0xb7, 0xf0, 0x8c, 0x37, 0xef, 0x3d,
	0xe3, 0xf7, 0xf7, 0x9e, 0xf1, 0xd3, 0xe7, 0x28, 0x58, 0x19, 0xcc, 0x8a, 0xf0, 0x9c, 0xc6, 0xaf,
	0xb0, 0x48, 0xc2, 0xf2, 0xbc, 0xfe, 0x59, 0xf8, 0x37, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x85, 0x77,
	0xd4, 0x26, 0x06, 0x00, 0x00,
}