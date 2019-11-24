// source: etop/etc/user_source/user_source.proto

package user_source

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
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

type UserSource int32

const (
	UserSource_unknown        UserSource = 0
	UserSource_psx            UserSource = 1
	UserSource_etop           UserSource = 2
	UserSource_topship        UserSource = 3
	UserSource_ts_app_android UserSource = 4
	UserSource_ts_app_ios     UserSource = 5
	UserSource_ts_app_web     UserSource = 6
	UserSource_partner        UserSource = 7
)

var UserSource_name = map[int32]string{
	0: "unknown",
	1: "psx",
	2: "etop",
	3: "topship",
	4: "ts_app_android",
	5: "ts_app_ios",
	6: "ts_app_web",
	7: "partner",
}

var UserSource_value = map[string]int32{
	"unknown":        0,
	"psx":            1,
	"etop":           2,
	"topship":        3,
	"ts_app_android": 4,
	"ts_app_ios":     5,
	"ts_app_web":     6,
	"partner":        7,
}

func (x UserSource) Enum() *UserSource {
	p := new(UserSource)
	*p = x
	return p
}

func (x UserSource) String() string {
	return proto.EnumName(UserSource_name, int32(x))
}

func (x *UserSource) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(UserSource_value, data, "UserSource")
	if err != nil {
		return err
	}
	*x = UserSource(value)
	return nil
}

func (UserSource) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_92152a3c9db54c68, []int{0}
}

func init() {
	proto.RegisterEnum("user_source.UserSource", UserSource_name, UserSource_value)
}

func init() {
	proto.RegisterFile("etop/etc/user_source/user_source.proto", fileDescriptor_92152a3c9db54c68)
}

var fileDescriptor_92152a3c9db54c68 = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xc1, 0x8a, 0xc2, 0x30,
	0x10, 0x86, 0xbb, 0xdb, 0xee, 0x76, 0x99, 0x42, 0x19, 0xc2, 0x9e, 0x3c, 0x78, 0x11, 0x04, 0x3d,
	0x34, 0xef, 0xe0, 0x2b, 0x88, 0x17, 0x2f, 0x25, 0x6d, 0x43, 0x0d, 0x42, 0x66, 0x48, 0x52, 0x2b,
	0x3e, 0x89, 0xef, 0xe7, 0x8b, 0x48, 0xaa, 0x87, 0x1e, 0xbc, 0xfd, 0xf3, 0xcd, 0xcf, 0x37, 0x0c,
	0xac, 0x75, 0x20, 0x96, 0x3a, 0xb4, 0x72, 0xf0, 0xda, 0xd5, 0x9e, 0x06, 0xd7, 0xea, 0x79, 0xae,
	0xd8, 0x51, 0x20, 0x51, 0xcc, 0xd0, 0xe2, 0xbf, 0xa7, 0x9e, 0x26, 0x2e, 0x63, 0x7a, 0x55, 0xb6,
	0x37, 0x80, 0x83, 0xd7, 0x6e, 0x3f, 0x75, 0x44, 0x01, 0xf9, 0x60, 0xcf, 0x96, 0x46, 0x8b, 0x89,
	0xc8, 0x21, 0x65, 0x7f, 0xc5, 0x2f, 0xf1, 0x07, 0x59, 0x3c, 0x88, 0xdf, 0x71, 0x1f, 0x88, 0xfd,
	0xc9, 0x30, 0xa6, 0x42, 0x40, 0x19, 0x7c, 0xad, 0x98, 0x6b, 0x65, 0x3b, 0x47, 0xa6, 0xc3, 0x4c,
	0x94, 0x00, 0x6f, 0x66, 0xc8, 0xe3, 0xcf, 0x6c, 0x1e, 0x75, 0x83, 0xbf, 0x51, 0xc0, 0xca, 0x05,
	0xab, 0x1d, 0xe6, 0xbb, 0xcd, 0xfd, 0xb1, 0x4c, 0x8e, 0xab, 0xe8, 0xae, 0x2e, 0x56, 0x2a, 0x36,
	0x92, 0x1b, 0xf9, 0xe9, 0xb7, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0xca, 0x8c, 0xa0, 0x89, 0xf2,
	0x00, 0x00, 0x00,
}
