package ghn_note_code

// +enum
// +enum:zero=null
type GHNNoteCode int

const (
	// +enum=unknown
	Unknown GHNNoteCode = 0

	// +enum=CHOTHUHANG
	CHOTHUHANG GHNNoteCode = 1

	// +enum=CHOXEMHANGKHONGTHU
	CHOXEMHANGKHONGTHU GHNNoteCode = 2

	// +enum=KHONGCHOXEMHANG
	KHONGCHOXEMHANG GHNNoteCode = 3
)
