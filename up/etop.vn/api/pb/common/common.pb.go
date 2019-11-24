// source: common/common.proto

package common

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

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

var xxx_messageInfo_Empty proto.InternalMessageInfo

type VersionInfoResponse struct {
	Service              string   `protobuf:"bytes,1,opt,name=service" json:"service"`
	Version              string   `protobuf:"bytes,2,opt,name=version" json:"version"`
	UpdatedAt            dot.Time `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionInfoResponse) Reset()         { *m = VersionInfoResponse{} }
func (m *VersionInfoResponse) String() string { return proto.CompactTextString(m) }
func (*VersionInfoResponse) ProtoMessage()    {}
func (*VersionInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{1}
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

type IDRequest struct {
	// @required
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDRequest) Reset()         { *m = IDRequest{} }
func (m *IDRequest) String() string { return proto.CompactTextString(m) }
func (*IDRequest) ProtoMessage()    {}
func (*IDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{2}
}

var xxx_messageInfo_IDRequest proto.InternalMessageInfo

func (m *IDRequest) GetId() dot.ID {
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

var xxx_messageInfo_NameRequest proto.InternalMessageInfo

func (m *NameRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type IDsRequest struct {
	// @required
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDsRequest) Reset()         { *m = IDsRequest{} }
func (m *IDsRequest) String() string { return proto.CompactTextString(m) }
func (*IDsRequest) ProtoMessage()    {}
func (*IDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{5}
}

var xxx_messageInfo_IDsRequest proto.InternalMessageInfo

func (m *IDsRequest) GetIds() []dot.ID {
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

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type IDMRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Paging               *Paging  `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDMRequest) Reset()         { *m = IDMRequest{} }
func (m *IDMRequest) String() string { return proto.CompactTextString(m) }
func (*IDMRequest) ProtoMessage()    {}
func (*IDMRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{7}
}

var xxx_messageInfo_IDMRequest proto.InternalMessageInfo

func (m *IDMRequest) GetId() dot.ID {
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
	Sort                 []string `protobuf:"bytes,3,rep,name=sort" json:"sort"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PageInfo) Reset()         { *m = PageInfo{} }
func (m *PageInfo) String() string { return proto.CompactTextString(m) }
func (*PageInfo) ProtoMessage()    {}
func (*PageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{10}
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
	Data                 []byte   `protobuf:"bytes,1,req,name=data" json:"data"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawJSONObject) Reset()         { *m = RawJSONObject{} }
func (m *RawJSONObject) String() string { return proto.CompactTextString(m) }
func (*RawJSONObject) ProtoMessage()    {}
func (*RawJSONObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{12}
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
	Meta                 map[string]string `protobuf:"bytes,3,rep,name=meta" json:"meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{13}
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
	Errors               []*Error `protobuf:"bytes,1,rep,name=errors" json:"errors"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorsResponse) Reset()         { *m = ErrorsResponse{} }
func (m *ErrorsResponse) String() string { return proto.CompactTextString(m) }
func (*ErrorsResponse) ProtoMessage()    {}
func (*ErrorsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{14}
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
	// - Array: = ≠ (only with value is {})
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

type CommonListRequest struct {
	Paging               *Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CommonListRequest) Reset()         { *m = CommonListRequest{} }
func (m *CommonListRequest) String() string { return proto.CompactTextString(m) }
func (*CommonListRequest) ProtoMessage()    {}
func (*CommonListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{20}
}

var xxx_messageInfo_CommonListRequest proto.InternalMessageInfo

func (m *CommonListRequest) GetPaging() *Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *CommonListRequest) GetFilters() []*Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type MetaField struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key"`
	Value                string   `protobuf:"bytes,2,opt,name=value" json:"value"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetaField) Reset()         { *m = MetaField{} }
func (m *MetaField) String() string { return proto.CompactTextString(m) }
func (*MetaField) ProtoMessage()    {}
func (*MetaField) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f954d82c0b891f6, []int{21}
}

var xxx_messageInfo_MetaField proto.InternalMessageInfo

func (m *MetaField) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *MetaField) GetValue() string {
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
	proto.RegisterType((*CommonListRequest)(nil), "cm.CommonListRequest")
	proto.RegisterType((*MetaField)(nil), "cm.MetaField")
}

func init() { proto.RegisterFile("common/common.proto", fileDescriptor_8f954d82c0b891f6) }

var fileDescriptor_8f954d82c0b891f6 = []byte{
	// 792 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xdd, 0x4e, 0xe3, 0x46,
	0x14, 0xf6, 0x4f, 0x7e, 0xc8, 0x49, 0x81, 0xd6, 0xd0, 0xca, 0x8a, 0x2a, 0x13, 0xa6, 0x95, 0xc8,
	0x95, 0xa3, 0xa6, 0x17, 0xfd, 0xb9, 0xa9, 0xca, 0x4f, 0x10, 0x55, 0xf9, 0x91, 0xa1, 0x5c, 0x54,
	0xad, 0xd0, 0x10, 0x4f, 0x8c, 0x69, 0xec, 0x71, 0x3d, 0x93, 0x00, 0x6f, 0xd0, 0x4b, 0x2e, 0x2b,
	0xed, 0x0b, 0xec, 0x23, 0xec, 0x23, 0x70, 0xc9, 0x13, 0xac, 0x48, 0x78, 0x81, 0x7d, 0x84, 0xd5,
	0xcc, 0xd8, 0x49, 0xcc, 0x0a, 0xa4, 0xdd, 0x2b, 0xc6, 0xdf, 0xf9, 0xce, 0x39, 0xdf, 0x99, 0xf9,
	0x0e, 0x81, 0x95, 0x1e, 0x8d, 0x22, 0x1a, 0xb7, 0xd5, 0x1f, 0x37, 0x49, 0x29, 0xa7, 0x96, 0xd1,
	0x8b, 0x1a, 0xab, 0x01, 0x0d, 0xa8, 0xfc, 0x6c, 0x8b, 0x93, 0x8a, 0x34, 0xd6, 0x02, 0x4a, 0x83,
	0x01, 0x69, 0xcb, 0xaf, 0xf3, 0x61, 0xbf, 0xcd, 0xc3, 0x88, 0x30, 0x8e, 0xa3, 0x44, 0x11, 0x50,
	0x15, 0xca, 0x3b, 0x51, 0xc2, 0x6f, 0xd0, 0xad, 0x0e, 0x2b, 0xa7, 0x24, 0x65, 0x21, 0x8d, 0xf7,
	0xe2, 0x3e, 0xf5, 0x08, 0x4b, 0x68, 0xcc, 0x88, 0xe5, 0x40, 0x95, 0x91, 0x74, 0x14, 0xf6, 0x88,
	0xad, 0x37, 0xf5, 0x56, 0x6d, 0xb3, 0x74, 0xf7, 0x76, 0x4d, 0xf3, 0x72, 0x50, 0xc4, 0x47, 0x2a,
	0xcd, 0x36, 0xe6, 0xe3, 0x19, 0x68, 0xfd, 0x04, 0x30, 0x4c, 0x7c, 0xcc, 0x89, 0x7f, 0x86, 0xb9,
	0x6d, 0x36, 0xf5, 0x56, 0xbd, 0xd3, 0x70, 0x95, 0x2c, 0x37, 0x97, 0xe5, 0x9e, 0xe4, 0xb2, 0xbc,
	0x5a, 0xc6, 0xfe, 0x95, 0xa3, 0x75, 0xa8, 0xed, 0x6d, 0x7b, 0xe4, 0xdf, 0x21, 0x61, 0xdc, 0x5a,
	0x05, 0x23, 0xf4, 0xa5, 0x04, 0x33, 0x6b, 0x61, 0x84, 0x3e, 0xda, 0x80, 0xfa, 0x16, 0xf5, 0x49,
	0x4e, 0xb2, 0xa1, 0xd4, 0xa3, 0x7e, 0x51, 0xa9, 0x44, 0x04, 0xf1, 0x00, 0x47, 0xf3, 0xc4, 0x18,
	0x47, 0x4f, 0x88, 0x02, 0x41, 0x0e, 0xc0, 0xde, 0x36, 0xcb, 0x79, 0x9f, 0x83, 0x19, 0xfa, 0xcc,
	0xd6, 0x9b, 0x66, 0xcb, 0xf4, 0xc4, 0x11, 0xb9, 0xb0, 0x74, 0xcc, 0x31, 0x1f, 0xb2, 0xe9, 0x0d,
	0x7d, 0x0d, 0x15, 0x26, 0x91, 0x42, 0xb5, 0x0c, 0x43, 0x5d, 0x51, 0x6f, 0xff, 0xc5, 0x29, 0x2c,
	0x04, 0x95, 0x04, 0x07, 0x61, 0x1c, 0xc8, 0x2b, 0xac, 0x77, 0xc0, 0xed, 0x45, 0xee, 0x91, 0x44,
	0xbc, 0x2c, 0x82, 0xfe, 0x82, 0x8a, 0x42, 0x44, 0x3f, 0xda, 0xef, 0x33, 0xc2, 0xed, 0x72, 0x53,
	0x6f, 0x95, 0xf3, 0x7e, 0x0a, 0xb3, 0x1a, 0x50, 0x1e, 0x84, 0x51, 0xc8, 0xed, 0xca, 0x5c, 0x50,
	0x41, 0x62, 0x6a, 0x46, 0x53, 0x6e, 0x57, 0xe7, 0xa7, 0x16, 0x08, 0xda, 0x85, 0xc5, 0x2e, 0x4d,
	0xaf, 0x70, 0xea, 0x67, 0x4d, 0x1a, 0x50, 0x66, 0x61, 0xfc, 0xe4, 0xd1, 0x15, 0x34, 0x6b, 0x61,
	0x7c, 0xd0, 0x02, 0x9d, 0xc2, 0xc2, 0x11, 0x0e, 0x88, 0xb0, 0x90, 0xe0, 0x71, 0xca, 0xf1, 0x40,
	0xd6, 0x98, 0xf2, 0x24, 0xf4, 0x52, 0x0d, 0xcb, 0xca, 0x64, 0x9a, 0x4d, 0xb3, 0x55, 0xcb, 0x04,
	0x5e, 0xc2, 0xf2, 0x4c, 0xe0, 0xb4, 0xfc, 0xa7, 0x48, 0xb4, 0xd6, 0x60, 0xe1, 0x02, 0xb3, 0xb3,
	0x98, 0x5c, 0x73, 0xbb, 0xd4, 0xd4, 0x5b, 0x0b, 0xb9, 0x65, 0x2f, 0x30, 0x3b, 0x20, 0xd7, 0x1c,
	0x7d, 0x03, 0x8b, 0x1e, 0xbe, 0xfa, 0xed, 0xf8, 0xf0, 0xe0, 0xf0, 0xfc, 0x92, 0xf4, 0xa4, 0x20,
	0x1f, 0x73, 0x6c, 0xeb, 0x4d, 0xa3, 0xf5, 0x99, 0x27, 0xcf, 0xe8, 0x95, 0x0e, 0xe5, 0x9d, 0x34,
	0xa5, 0xe9, 0xf3, 0xa6, 0xb3, 0xbe, 0x02, 0x33, 0x62, 0x41, 0x61, 0x2f, 0x04, 0x60, 0x6d, 0x40,
	0x29, 0x22, 0x1c, 0xcb, 0x01, 0xeb, 0x9d, 0x15, 0xf1, 0xda, 0xb2, 0x94, 0xbb, 0x4f, 0x38, 0xde,
	0x89, 0x79, 0x7a, 0xe3, 0x49, 0x42, 0xe3, 0x07, 0xa8, 0x4d, 0x21, 0xe1, 0xc5, 0x7f, 0xc8, 0x8d,
	0x6a, 0xe3, 0x89, 0xa3, 0xb5, 0x0a, 0xe5, 0x11, 0x1e, 0x0c, 0x89, 0xea, 0xe0, 0xa9, 0x8f, 0x9f,
	0x8d, 0x1f, 0x75, 0xf4, 0x3d, 0x2c, 0xc9, 0x8a, 0x33, 0x97, 0xae, 0x43, 0x85, 0x48, 0x44, 0x9a,
	0xb9, 0xde, 0xa9, 0x4d, 0xbb, 0x7a, 0x59, 0x00, 0x7d, 0x07, 0xcb, 0x7f, 0xa8, 0xe5, 0x9b, 0xdf,
	0xfe, 0x6c, 0x1f, 0x0b, 0x8f, 0x98, 0x83, 0x22, 0xc5, 0x23, 0x11, 0x1d, 0x15, 0x53, 0x52, 0x05,
	0x15, 0x53, 0x32, 0x50, 0xa4, 0x6c, 0x93, 0x01, 0x79, 0xd2, 0xc5, 0x57, 0x50, 0x31, 0x25, 0x03,
	0xd1, 0x16, 0x2c, 0xef, 0x13, 0xc6, 0x70, 0x40, 0xa6, 0x29, 0x1f, 0x7d, 0xe9, 0xe8, 0x04, 0x2a,
	0xdd, 0x70, 0xc0, 0x49, 0xfa, 0xfc, 0xf2, 0x8b, 0xf5, 0xa4, 0x49, 0x21, 0xd5, 0xa0, 0x89, 0x30,
	0x93, 0xba, 0x66, 0x73, 0xde, 0x68, 0x12, 0x42, 0x7f, 0xc3, 0x17, 0x5b, 0xf2, 0x5f, 0xf1, 0xef,
	0x21, 0xe3, 0xf9, 0x96, 0xcf, 0xf6, 0x59, 0x7f, 0x6e, 0x9f, 0xad, 0x6f, 0xa1, 0xda, 0x97, 0x72,
	0x98, 0x6d, 0xc8, 0x07, 0x91, 0x24, 0xa5, 0xd0, 0xcb, 0x43, 0xe8, 0x17, 0x65, 0x80, 0x6e, 0x48,
	0x06, 0xbe, 0x98, 0x6c, 0x6a, 0x80, 0x7c, 0x32, 0x61, 0x83, 0x46, 0xc1, 0x06, 0x05, 0x7d, 0x9b,
	0xbb, 0x77, 0x63, 0x47, 0xbf, 0x1f, 0x3b, 0xda, 0xc3, 0xd8, 0xd1, 0xde, 0x8d, 0x1d, 0xed, 0xbf,
	0x89, 0xa3, 0xbd, 0x9e, 0x38, 0xda, 0x9b, 0x89, 0xa3, 0xdd, 0x4d, 0x1c, 0xed, 0x7e, 0xe2, 0x68,
	0x0f, 0x13, 0x47, 0xbb, 0x7d, 0x74, 0xb4, 0xff, 0x1f, 0x1d, 0xed, 0xcf, 0x2f, 0x09, 0xa7, 0x89,
	0x3b, 0x8a, 0xdb, 0x38, 0x09, 0xdb, 0xc9, 0x79, 0xf6, 0x4b, 0xf3, 0x3e, 0x00, 0x00, 0xff, 0xff,
	0x17, 0xd5, 0x79, 0xb6, 0x79, 0x06, 0x00, 0x00,
}
