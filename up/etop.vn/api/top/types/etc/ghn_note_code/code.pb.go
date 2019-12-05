package ghn_note_code

import (
	"etop.vn/common/jsonx"
)

type GHNNoteCode int

const (
	GHNNoteCode_unknown            GHNNoteCode = 0
	GHNNoteCode_CHOTHUHANG         GHNNoteCode = 1
	GHNNoteCode_CHOXEMHANGKHONGTHU GHNNoteCode = 2
	GHNNoteCode_KHONGCHOXEMHANG    GHNNoteCode = 3
)

var GHNNoteCode_name = map[int]string{
	0: "unknown",
	1: "CHOTHUHANG",
	2: "CHOXEMHANGKHONGTHU",
	3: "KHONGCHOXEMHANG",
}

var GHNNoteCode_value = map[string]int{
	"unknown":            0,
	"CHOTHUHANG":         1,
	"CHOXEMHANGKHONGTHU": 2,
	"KHONGCHOXEMHANG":    3,
}

func (s GHNNoteCode) Enum() *GHNNoteCode {
	p := new(GHNNoteCode)
	*p = s
	return p
}

func (s GHNNoteCode) String() string {
	return jsonx.EnumName(GHNNoteCode_name, int(s))
}

func (s *GHNNoteCode) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(GHNNoteCode_value, data, "GHNNoteCode")
	if err != nil {
		return err
	}
	*s = GHNNoteCode(value)
	return nil
}
