package ghn_note_code

import (
	"etop.vn/common/jsonx"
)

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
	return jsonx.EnumName(GHNNoteCode_name, int32(x))
}

func (x *GHNNoteCode) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(GHNNoteCode_value, data, "GHNNoteCode")
	if err != nil {
		return err
	}
	*x = GHNNoteCode(value)
	return nil
}
