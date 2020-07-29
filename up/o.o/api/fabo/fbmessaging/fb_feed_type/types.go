package fb_feed_type

// +enum
// +enum:sql=int
// +enum:zero=null
type FbFeedType int

type NullFbFeedType struct {
	Enum  FbFeedType
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbFeedType = 0

	// +enum=web
	Post FbFeedType = 233

	// +enum=app
	Event FbFeedType = 278
)
