// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: services/handler/handler.proto

package handler

import (
	_ "etop.vn/backend/pb/common"
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

type ResetStateRequest struct {
	WebhookId            int64    `protobuf:"varint,1,opt,name=webhook_id,json=webhookId" json:"webhook_id"`
	AccountId            int64    `protobuf:"varint,2,opt,name=account_id,json=accountId" json:"account_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResetStateRequest) Reset()         { *m = ResetStateRequest{} }
func (m *ResetStateRequest) String() string { return proto.CompactTextString(m) }
func (*ResetStateRequest) ProtoMessage()    {}
func (*ResetStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb80502182612866, []int{0}
}
func (m *ResetStateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResetStateRequest.Unmarshal(m, b)
}
func (m *ResetStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResetStateRequest.Marshal(b, m, deterministic)
}
func (m *ResetStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResetStateRequest.Merge(m, src)
}
func (m *ResetStateRequest) XXX_Size() int {
	return xxx_messageInfo_ResetStateRequest.Size(m)
}
func (m *ResetStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResetStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResetStateRequest proto.InternalMessageInfo

func (m *ResetStateRequest) GetWebhookId() int64 {
	if m != nil {
		return m.WebhookId
	}
	return 0
}

func (m *ResetStateRequest) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func init() {
	proto.RegisterType((*ResetStateRequest)(nil), "handler.ResetStateRequest")
}

func init() { proto.RegisterFile("services/handler/handler.proto", fileDescriptor_cb80502182612866) }

var fileDescriptor_cb80502182612866 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2b, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd6, 0xcf, 0x48, 0xcc, 0x4b, 0xc9, 0x49, 0x2d, 0x82, 0xd1, 0x7a, 0x05,
	0x45, 0xf9, 0x25, 0xf9, 0x42, 0xec, 0x50, 0xae, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58, 0x4c,
	0x1f, 0xc4, 0x82, 0x48, 0x4b, 0x09, 0x27, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7, 0xe9, 0x43, 0x28, 0x88,
	0xa0, 0x52, 0x2c, 0x97, 0x60, 0x50, 0x6a, 0x71, 0x6a, 0x49, 0x70, 0x49, 0x62, 0x49, 0x6a, 0x50,
	0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x90, 0x32, 0x17, 0x57, 0x79, 0x6a, 0x52, 0x46, 0x7e, 0x7e,
	0x76, 0x7c, 0x66, 0x8a, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0xb3, 0x13, 0xcb, 0x89, 0x7b, 0xf2, 0x0c,
	0x41, 0x9c, 0x50, 0x71, 0xcf, 0x14, 0x90, 0xa2, 0xc4, 0xe4, 0xe4, 0xfc, 0xd2, 0xbc, 0x12, 0x90,
	0x22, 0x26, 0x64, 0x45, 0x50, 0x71, 0xcf, 0x14, 0xa7, 0xd0, 0x13, 0x0f, 0xe5, 0x18, 0x2f, 0x3c,
	0x94, 0x63, 0x7c, 0xf0, 0x50, 0x8e, 0xe1, 0xc3, 0x43, 0x39, 0x86, 0x8e, 0x47, 0x72, 0x0c, 0x2b,
	0x1e, 0xc9, 0x31, 0xec, 0x78, 0x24, 0xc7, 0x70, 0xe2, 0x91, 0x1c, 0xc3, 0x85, 0x47, 0x72, 0x0c,
	0x0f, 0x1e, 0xc9, 0x31, 0x4c, 0x78, 0x2c, 0xc7, 0x30, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x72, 0x6a,
	0x49, 0x7e, 0x81, 0x5e, 0x59, 0x9e, 0x7e, 0x52, 0x62, 0x72, 0x76, 0x6a, 0x5e, 0x8a, 0x7e, 0x41,
	0x92, 0x3e, 0xba, 0xbf, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x54, 0x88, 0x8b, 0x15, 0x0a, 0x01,
	0x00, 0x00,
}
