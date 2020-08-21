package typeutil

import (
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/try_on"
)

func TryOnFromGHNNoteCode(c ghn_note_code.GHNNoteCode) try_on.TryOnCode {
	switch c {
	case ghn_note_code.KHONGCHOXEMHANG:
		return try_on.None
	case ghn_note_code.CHOXEMHANGKHONGTHU:
		return try_on.Open
	case ghn_note_code.CHOTHUHANG:
		return try_on.Try
	default:
		return 0
	}
}

func GHNNoteCodeFromTryOn(to try_on.TryOnCode) ghn_note_code.GHNNoteCode {
	switch to {
	case try_on.None:
		return ghn_note_code.KHONGCHOXEMHANG
	case try_on.Open:
		return ghn_note_code.CHOXEMHANGKHONGTHU
	case try_on.Try:
		return ghn_note_code.CHOTHUHANG
	default:
		return 0
	}
}
