package ghn_note_code

import (
	try_on2 "etop.vn/api/pb/etop/etc/try_on"
)

func PbGHNNoteCode(s string) GHNNoteCode {
	return GHNNoteCode(GHNNoteCode_value[s])
}

func PbGHNNoteCodeFromInt(s int) GHNNoteCode {
	return GHNNoteCode(s)
}

func (s *GHNNoteCode) ToModel() string {
	if s == nil || *s == 0 {
		return ""
	}
	return GHNNoteCode_name[int32(*s)]
}

func (s GHNNoteCode) ToTryOn() try_on2.TryOnCode {
	switch s {
	case GHNNoteCode_KHONGCHOXEMHANG:
		return try_on2.TryOnCode_none
	case GHNNoteCode_CHOXEMHANGKHONGTHU:
		return try_on2.TryOnCode_open
	case GHNNoteCode_CHOTHUHANG:
		return try_on2.TryOnCode_try
	default:
		return 0
	}
}

func FromTryOn(code try_on2.TryOnCode) GHNNoteCode {
	switch code {
	case try_on2.TryOnCode_none:
		return GHNNoteCode_KHONGCHOXEMHANG
	case try_on2.TryOnCode_open:
		return GHNNoteCode_CHOXEMHANGKHONGTHU
	case try_on2.TryOnCode_try:
		return GHNNoteCode_CHOTHUHANG
	default:
		return 0
	}
}

func (x GHNNoteCode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
