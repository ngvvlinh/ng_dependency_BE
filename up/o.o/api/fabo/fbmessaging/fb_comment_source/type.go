package fb_comment_source

// +enum
// +enum:sql=int
// +enum:zero=null
type FbCommentSource int

type NullFbCommentSource struct {
	Enum  FbCommentSource
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbCommentSource = 0

	// +enum=web
	Web FbCommentSource = 427

	// +enum=app
	App FbCommentSource = 135
)
