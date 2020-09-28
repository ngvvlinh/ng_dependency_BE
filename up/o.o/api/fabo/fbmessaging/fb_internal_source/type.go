package fb_internal_source

// +enum
// +enum:sql=int
// +enum:zero=null
type FbInternalSource int

type NullFbInternalSource struct {
	Enum  FbInternalSource
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbInternalSource = 0

	// +enum=fabo
	Fabo FbInternalSource = 101

	// +enum=facebook
	Facebook FbInternalSource = 256
)
