// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: external/shop/shop.proto

package shop

import (
	_ "etop.vn/backend/pb/common"
	_ "etop.vn/backend/pb/etop"
	_ "etop.vn/backend/pb/etop/etc/address_type"
	_ "etop.vn/backend/pb/etop/etc/status3"
	_ "etop.vn/backend/pb/external"
	_ "etop.vn/backend/pb/external/etc/auth_mode"
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

func init() { proto.RegisterFile("external/shop/shop.proto", fileDescriptor_f1f5c2bb46255d6c) }

var fileDescriptor_f1f5c2bb46255d6c = []byte{
	// 624 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0x76, 0x7e, 0x3f, 0x0e, 0xb0, 0x51, 0x53, 0xe1, 0x52, 0xb5, 0x0a, 0x74, 0x1b, 0x24, 0xb8,
	0xd1, 0xac, 0x1a, 0x7a, 0x40, 0x48, 0x1c, 0xd2, 0xb4, 0x84, 0xaa, 0x45, 0xad, 0x5a, 0xfe, 0x48,
	0x5c, 0x2a, 0xc7, 0x9e, 0x38, 0x56, 0xed, 0x1d, 0xb3, 0x3b, 0x2e, 0xf4, 0x0d, 0x38, 0x22, 0x4e,
	0x3c, 0x42, 0x8f, 0x1c, 0x79, 0x84, 0x1e, 0x7b, 0xe3, 0xda, 0xa4, 0x2f, 0xc0, 0x23, 0x20, 0x3b,
	0xeb, 0x38, 0x09, 0x05, 0x7a, 0x19, 0xcf, 0x7e, 0xdf, 0x37, 0xdf, 0xec, 0x8c, 0x2d, 0xb3, 0x45,
	0xf8, 0x48, 0xa0, 0xa4, 0x13, 0x0a, 0xdd, 0xc3, 0x38, 0x0b, 0xf5, 0x58, 0x21, 0xa1, 0x7d, 0x23,
	0xcd, 0xab, 0x8f, 0xb2, 0x83, 0xbb, 0xe2, 0x83, 0x5c, 0xd1, 0x1f, 0x1c, 0xdf, 0x07, 0x25, 0x30,
	0xa6, 0x00, 0xa5, 0x16, 0x8e, 0x94, 0x48, 0x4e, 0x96, 0x0f, 0x6b, 0xaa, 0x77, 0x7c, 0xf4, 0x31,
	0x4b, 0x45, 0x9a, 0x19, 0x74, 0xd9, 0x47, 0xf4, 0x43, 0x10, 0xd9, 0xa9, 0x93, 0x74, 0x05, 0x05,
	0x11, 0x68, 0x72, 0x22, 0xd3, 0xaa, 0x3a, 0xe7, 0x62, 0x14, 0xa1, 0x14, 0xc3, 0x87, 0x01, 0x67,
	0x81, 0x30, 0x16, 0x69, 0x30, 0xc0, 0xc2, 0xe8, 0xaa, 0x79, 0x62, 0x88, 0x25, 0xa3, 0x74, 0x85,
	0x26, 0x87, 0x12, 0xfd, 0xd8, 0x3c, 0x0d, 0x5d, 0x1b, 0xd1, 0x8e, 0xe7, 0x29, 0xd0, 0xfa, 0x90,
	0x4e, 0x62, 0x10, 0x69, 0x30, 0x8a, 0xfb, 0x85, 0x73, 0xaa, 0x4a, 0xa8, 0x77, 0x18, 0xa1, 0x07,
	0x22, 0x0d, 0x43, 0x49, 0xe3, 0x5b, 0x89, 0x95, 0x5f, 0x06, 0xda, 0x3d, 0x00, 0x75, 0x1c, 0xb8,
	0x60, 0xaf, 0xb2, 0xf2, 0x1b, 0x50, 0x3a, 0x40, 0xb9, 0x25, 0xbb, 0x68, 0xdf, 0xaa, 0xbb, 0x51,
	0x7d, 0x33, 0x8a, 0xe9, 0xa4, 0xba, 0x90, 0xa6, 0x63, 0xdc, 0x3e, 0xe8, 0x18, 0xa5, 0x06, 0x7b,
	0x8d, 0x55, 0x5a, 0x89, 0x52, 0x20, 0xa9, 0xe9, 0xba, 0x98, 0x48, 0x9a, 0xac, 0xca, 0x26, 0xdd,
	0x4b, 0x3a, 0x61, 0xe0, 0x1a, 0x3e, 0x73, 0x7e, 0xc2, 0x66, 0xdb, 0x40, 0x3b, 0xe8, 0x66, 0x8b,
	0xde, 0x09, 0xf4, 0x44, 0x59, 0xb5, 0x3e, 0xda, 0x45, 0x2e, 0xc9, 0xfb, 0x35, 0x7e, 0x94, 0x58,
	0xe5, 0x2d, 0x74, 0x7a, 0x88, 0x47, 0xf9, 0xad, 0xd7, 0xd9, 0x4c, 0x4b, 0x81, 0x43, 0x60, 0x70,
	0x9b, 0x17, 0xf5, 0x13, 0xc4, 0x3e, 0xbc, 0x4f, 0x40, 0x53, 0xf5, 0x76, 0xc1, 0xe7, 0x25, 0x6b,
	0xac, 0xdc, 0x06, 0x32, 0x27, 0xfd, 0x87, 0xcb, 0xe4, 0xf4, 0x68, 0xf8, 0x6d, 0x36, 0xb3, 0x01,
	0x21, 0x5c, 0xd9, 0x79, 0x82, 0xc8, 0x3b, 0xff, 0xc5, 0xac, 0xb1, 0xcb, 0x2a, 0x2f, 0x02, 0x4d,
	0xa8, 0x4e, 0xf2, 0xc1, 0x9e, 0x31, 0xd6, 0x06, 0x6a, 0xf5, 0x1c, 0xe9, 0x83, 0xb6, 0xef, 0x16,
	0xb5, 0x05, 0x9a, 0x1b, 0xdb, 0x63, 0x23, 0x3b, 0x61, 0xd8, 0x71, 0xdc, 0xa3, 0xc6, 0xe9, 0xff,
	0x6c, 0xf6, 0xa0, 0x17, 0xc4, 0x71, 0x20, 0xfd, 0xdc, 0xb2, 0xc3, 0xe6, 0xda, 0x40, 0x53, 0xa8,
	0xb6, 0x1f, 0x4c, 0x78, 0x4f, 0xd3, 0x79, 0x93, 0x87, 0xff, 0x50, 0x99, 0xad, 0xbc, 0x66, 0xf3,
	0xc3, 0xb5, 0x37, 0xa5, 0xd7, 0x42, 0xd9, 0x0d, 0x54, 0xb4, 0xab, 0x3c, 0x50, 0xf6, 0xbd, 0xe9,
	0xf7, 0x92, 0xc1, 0xb9, 0xfb, 0xd8, 0xee, 0x32, 0xbc, 0x29, 0xbd, 0xe7, 0x49, 0xd8, 0x0d, 0xc2,
	0x30, 0x02, 0x49, 0xda, 0xde, 0x66, 0xe5, 0x96, 0x23, 0x5d, 0x08, 0x7f, 0x37, 0x2b, 0xe0, 0xeb,
	0x9a, 0x6d, 0xb0, 0x9b, 0x6d, 0xa0, 0xa1, 0xd3, 0xe2, 0x94, 0x76, 0x6b, 0xe3, 0xba, 0x2e, 0x6d,
	0x56, 0x69, 0x03, 0x8d, 0x41, 0xe3, 0x1f, 0xc0, 0x18, 0x5c, 0x38, 0xce, 0x5f, 0xc9, 0xaf, 0xd3,
	0x97, 0xe6, 0xb2, 0xbd, 0xc4, 0x16, 0xe0, 0x15, 0xc6, 0xb5, 0x4d, 0xa3, 0xa8, 0x35, 0xf7, 0xb6,
	0x9e, 0xd6, 0x0e, 0x7a, 0x18, 0x37, 0xfe, 0x3b, 0x5e, 0x3d, 0xeb, 0xf3, 0xd2, 0x79, 0x9f, 0x97,
	0x2e, 0xfa, 0xdc, 0xfa, 0xd9, 0xe7, 0xd6, 0xa7, 0x01, 0xb7, 0x4e, 0x07, 0xdc, 0xfa, 0x3e, 0xe0,
	0xd6, 0xd9, 0x80, 0x5b, 0xe7, 0x03, 0x6e, 0x5d, 0x0c, 0xb8, 0xf5, 0xf9, 0x92, 0x5b, 0x5f, 0x2f,
	0xb9, 0xf5, 0x2e, 0xfb, 0x43, 0xd4, 0x8f, 0xa5, 0x48, 0x3f, 0x06, 0x90, 0x9e, 0x88, 0x3b, 0x62,
	0xe2, 0xbf, 0xf8, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x48, 0x00, 0xb6, 0xb7, 0x27, 0x05, 0x00, 0x00,
}
