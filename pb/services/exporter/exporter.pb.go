// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: services/exporter/exporter.proto

package exporter

import (
	common "etop.vn/backend/pb/common"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

type RequestExportRequest struct {
	Filters              []*common.Filter `protobuf:"bytes,1,rep,name=filters" json:"filters,omitempty"`
	DateFrom             string           `protobuf:"bytes,2,opt,name=date_from,json=dateFrom" json:"date_from"`
	DateTo               string           `protobuf:"bytes,3,opt,name=date_to,json=dateTo" json:"date_to"`
	ShopId               int64            `protobuf:"varint,4,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *RequestExportRequest) Reset()         { *m = RequestExportRequest{} }
func (m *RequestExportRequest) String() string { return proto.CompactTextString(m) }
func (*RequestExportRequest) ProtoMessage()    {}
func (*RequestExportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fcd1644d6206654f, []int{0}
}
func (m *RequestExportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestExportRequest.Unmarshal(m, b)
}
func (m *RequestExportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestExportRequest.Marshal(b, m, deterministic)
}
func (m *RequestExportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestExportRequest.Merge(m, src)
}
func (m *RequestExportRequest) XXX_Size() int {
	return xxx_messageInfo_RequestExportRequest.Size(m)
}
func (m *RequestExportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestExportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RequestExportRequest proto.InternalMessageInfo

func (m *RequestExportRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *RequestExportRequest) GetDateFrom() string {
	if m != nil {
		return m.DateFrom
	}
	return ""
}

func (m *RequestExportRequest) GetDateTo() string {
	if m != nil {
		return m.DateTo
	}
	return ""
}

func (m *RequestExportRequest) GetShopId() int64 {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func init() {
	proto.RegisterType((*RequestExportRequest)(nil), "exporter.RequestExportRequest")
}

func init() { proto.RegisterFile("services/exporter/exporter.proto", fileDescriptor_fcd1644d6206654f) }

var fileDescriptor_fcd1644d6206654f = []byte{
	// 363 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x41, 0xee, 0xd2, 0x40,
	0x14, 0xc6, 0xa7, 0x42, 0x04, 0x86, 0x5d, 0x25, 0xb1, 0x21, 0x71, 0xa8, 0x86, 0x05, 0x0b, 0xe9,
	0x44, 0x4e, 0x60, 0x48, 0x20, 0x61, 0xe1, 0xa6, 0x1a, 0x16, 0x6e, 0x48, 0x69, 0x5f, 0x6b, 0x23,
	0x33, 0x6f, 0x9c, 0x19, 0x50, 0x6f, 0xe0, 0xd2, 0xa5, 0x1b, 0xf7, 0x1e, 0xc1, 0x23, 0xb0, 0xe4,
	0x04, 0x86, 0x96, 0x0b, 0x78, 0x04, 0xd3, 0x16, 0x0c, 0xc9, 0x7f, 0x35, 0xbf, 0xf7, 0xfb, 0xbe,
	0x49, 0x5e, 0x66, 0xa8, 0x6f, 0x40, 0x1f, 0xf2, 0x18, 0x0c, 0x87, 0x2f, 0x0a, 0xb5, 0x05, 0xfd,
	0x1f, 0x02, 0xa5, 0xd1, 0xa2, 0xdb, 0xbd, 0xcd, 0xc3, 0x97, 0xb5, 0x88, 0xa7, 0x19, 0xc8, 0xa9,
	0xf9, 0x1c, 0x65, 0x19, 0x68, 0x8e, 0xca, 0xe6, 0x28, 0x0d, 0x8f, 0xa4, 0x44, 0x1b, 0xd5, 0xdc,
	0xdc, 0x1b, 0x0e, 0x32, 0xcc, 0xb0, 0x46, 0x5e, 0xd1, 0xd5, 0x8e, 0x32, 0xc4, 0x6c, 0x07, 0xbc,
	0x9e, 0xb6, 0xfb, 0x94, 0xdb, 0x5c, 0x80, 0xb1, 0x91, 0x50, 0xd7, 0xc2, 0x93, 0x18, 0x85, 0x40,
	0xc9, 0x9b, 0xa3, 0x91, 0x2f, 0x7e, 0x3a, 0x74, 0x10, 0xc2, 0xa7, 0x3d, 0x18, 0xbb, 0xa8, 0xb7,
	0xb9, 0x0e, 0xee, 0x98, 0x76, 0xd2, 0x7c, 0x67, 0x41, 0x1b, 0xcf, 0xf1, 0x5b, 0x93, 0xfe, 0x8c,
	0x06, 0xb1, 0x08, 0x96, 0xb5, 0x0a, 0x6f, 0x91, 0xfb, 0x9c, 0xf6, 0x92, 0xc8, 0xc2, 0x26, 0xd5,
	0x28, 0xbc, 0x47, 0xbe, 0x33, 0xe9, 0xcd, 0xdb, 0xc7, 0x3f, 0x23, 0x12, 0x76, 0x2b, 0xbd, 0xd4,
	0x28, 0xdc, 0x67, 0xb4, 0x53, 0x57, 0x2c, 0x7a, 0xad, 0xbb, 0xc2, 0xe3, 0x4a, 0xbe, 0xc3, 0x2a,
	0x36, 0x1f, 0x50, 0x6d, 0xf2, 0xc4, 0x6b, 0xfb, 0xce, 0xa4, 0x75, 0x8b, 0x2b, 0xb9, 0x4a, 0x66,
	0xaf, 0x69, 0xff, 0x4d, 0x6e, 0xe2, 0xb7, 0xcd, 0x5b, 0xba, 0xaf, 0x68, 0x7f, 0x0d, 0xda, 0xe4,
	0x28, 0x57, 0x32, 0x45, 0xb7, 0x57, 0xed, 0xb4, 0x10, 0xca, 0x7e, 0x1d, 0x3e, 0xad, 0xf0, 0x2e,
	0x0b, 0xc1, 0x28, 0x94, 0x06, 0xe6, 0xeb, 0x63, 0xc1, 0x9c, 0x53, 0xc1, 0x9c, 0x73, 0xc1, 0xc8,
	0xdf, 0x82, 0x91, 0x6f, 0x25, 0x23, 0xbf, 0x4a, 0x46, 0x7e, 0x97, 0x8c, 0x1c, 0x4b, 0x46, 0x4e,
	0x25, 0x23, 0xe7, 0x92, 0x91, 0xef, 0x17, 0x46, 0x7e, 0x5c, 0x18, 0x79, 0x3f, 0x06, 0x8b, 0x2a,
	0x38, 0x48, 0xbe, 0x8d, 0xe2, 0x8f, 0x20, 0x13, 0xae, 0xb6, 0xfc, 0xc1, 0x67, 0xfe, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0xd5, 0xe0, 0x0c, 0x0f, 0xe0, 0x01, 0x00, 0x00,
}
