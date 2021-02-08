package fb_post_source

// +enum
// +enum:zero=null
type FbPostSource int

type NullFbPostSource struct {
	Enum  FbPostSource
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbPostSource = 0

	// +enum=page
	Page FbPostSource = 13

	// +enum=user
	User FbPostSource = 42
)

