// source: etop/etc/ghn_note_code/code.proto

package ghn_note_code

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

type GHNNoteCode int32

const (
	GHNNoteCode_unknown            GHNNoteCode = 0
	GHNNoteCode_CHOTHUHANG         GHNNoteCode = 1
	GHNNoteCode_CHOXEMHANGKHONGTHU GHNNoteCode = 2
	GHNNoteCode_KHONGCHOXEMHANG    GHNNoteCode = 3
)

var GHNNoteCode_name = map[int32]string{
	0: "unknown",
	1: "CHOTHUHANG",
	2: "CHOXEMHANGKHONGTHU",
	3: "KHONGCHOXEMHANG",
}

var GHNNoteCode_value = map[string]int32{
	"unknown":            0,
	"CHOTHUHANG":         1,
	"CHOXEMHANGKHONGTHU": 2,
	"KHONGCHOXEMHANG":    3,
}

func (x GHNNoteCode) Enum() *GHNNoteCode {
	p := new(GHNNoteCode)
	*p = x
	return p
}

func (x GHNNoteCode) String() string {
	return proto.EnumName(GHNNoteCode_name, int32(x))
}

func (x *GHNNoteCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(GHNNoteCode_value, data, "GHNNoteCode")
	if err != nil {
		return err
	}
	*x = GHNNoteCode(value)
	return nil
}

func (GHNNoteCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fdc157f02903a8b5, []int{0}
}

func init() {
	proto.RegisterEnum("ghn_note_code.GHNNoteCode", GHNNoteCode_name, GHNNoteCode_value)
}

func init() { proto.RegisterFile("etop/etc/ghn_note_code/code.proto", fileDescriptor_fdc157f02903a8b5) }

var fileDescriptor_fdc157f02903a8b5 = []byte{
	// 173 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x2d, 0xc9, 0x2f,
	0xd0, 0x4f, 0x2d, 0x49, 0xd6, 0x4f, 0xcf, 0xc8, 0x8b, 0xcf, 0xcb, 0x2f, 0x49, 0x8d, 0x4f, 0xce,
	0x4f, 0x49, 0xd5, 0x07, 0x11, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xbc, 0x28, 0x32, 0x52,
	0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x19, 0x7d, 0x10, 0x0b, 0xa2, 0x48, 0x2b, 0x9c, 0x8b, 0xdb,
	0xdd, 0xc3, 0xcf, 0x2f, 0xbf, 0x24, 0xd5, 0x39, 0x3f, 0x25, 0x55, 0x88, 0x9b, 0x8b, 0xbd, 0x34,
	0x2f, 0x3b, 0x2f, 0xbf, 0x3c, 0x4f, 0x80, 0x41, 0x88, 0x8f, 0x8b, 0xcb, 0xd9, 0xc3, 0x3f, 0xc4,
	0x23, 0xd4, 0xc3, 0xd1, 0xcf, 0x5d, 0x80, 0x51, 0x48, 0x8c, 0x4b, 0xc8, 0xd9, 0xc3, 0x3f, 0xc2,
	0xd5, 0x17, 0xc4, 0xf7, 0xf6, 0xf0, 0xf7, 0x73, 0x0f, 0xf1, 0x08, 0x15, 0x60, 0x12, 0x12, 0xe6,
	0xe2, 0x07, 0xf3, 0x10, 0x92, 0x02, 0xcc, 0x4e, 0xda, 0x33, 0x1e, 0xcb, 0x31, 0x44, 0xa9, 0x82,
	0x9c, 0xa9, 0x57, 0x96, 0xa7, 0x9f, 0x58, 0x90, 0xa9, 0x5f, 0x90, 0xa4, 0x8f, 0xdd, 0xd5, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xe9, 0xd3, 0xb4, 0xce, 0x00, 0x00, 0x00,
}
