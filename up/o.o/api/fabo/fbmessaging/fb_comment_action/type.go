package fb_comment_action

// +enum
// +enum:sql=int
// +enum:zero=null
type FbCommentAction int

type NullFbCommentAction struct {
	Enum  FbCommentAction
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbCommentAction = 0

	// +enum=like
	Like FbCommentAction = 123

	// +enum=unlike
	UnLike FbCommentAction = 145

	// +enum=hide
	Hide FbCommentAction = 275

	// +enum=unhide
	UnHide FbCommentAction = 842
)

func (s FbCommentAction) Label() string {
	switch s {
	case Like:
		return "thích"
	case UnLike:
		return "bỏ thích"
	case Hide:
		return "ẩn"
	case UnHide:
		return "bỏ ẩn"
	default:
		return ""
	}
}
