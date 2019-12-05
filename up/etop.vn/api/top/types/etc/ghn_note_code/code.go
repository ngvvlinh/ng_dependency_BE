package ghn_note_code

import (
	"etop.vn/api/top/types/etc/try_on"
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
	return GHNNoteCode_name[int(*s)]
}

func (s GHNNoteCode) ToTryOn() try_on.TryOnCode {
	switch s {
	case GHNNoteCode_KHONGCHOXEMHANG:
		return try_on.TryOnCode_none
	case GHNNoteCode_CHOXEMHANGKHONGTHU:
		return try_on.TryOnCode_open
	case GHNNoteCode_CHOTHUHANG:
		return try_on.TryOnCode_try
	default:
		return 0
	}
}

func FromTryOn(code try_on.TryOnCode) GHNNoteCode {
	switch code {
	case try_on.TryOnCode_none:
		return GHNNoteCode_KHONGCHOXEMHANG
	case try_on.TryOnCode_open:
		return GHNNoteCode_CHOXEMHANGKHONGTHU
	case try_on.TryOnCode_try:
		return GHNNoteCode_CHOTHUHANG
	default:
		return 0
	}
}

func (s GHNNoteCode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}
