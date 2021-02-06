package fb_post_type

// +enum
// +enum:sql=int
// +enum:zero=null
type FbPostType int

type NullFbPostType struct {
	Enum  FbPostType
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbPostType = 0

	// +enum=page
	Page FbPostType = 101

	// +enum=user
	User FbPostType = 256
)