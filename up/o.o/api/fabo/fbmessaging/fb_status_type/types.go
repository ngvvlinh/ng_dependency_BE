package fb_status_type

// +enum
// +enum:sql=int
// +enum:zero=null
type FbStatusType int

type NullFbStatusType struct {
	Enum  FbStatusType
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbStatusType = 0

	// +enum=added_photos
	AddedPhotos FbStatusType = 23

	// +enum=added_video
	AddedVideo FbStatusType = 27

	// +enum=app_created_story
	AppCreatedStory FbStatusType = 28

	// +enum=approved_friend
	ApprovedFriend FbStatusType = 32

	// +enum=created_event
	CreatedEvent FbStatusType = 37

	// +enum=created_group
	CreatedGroup FbStatusType = 41

	// +enum=created_note
	CreatedNote FbStatusType = 45

	// +enum=mobile_status_update
	MobileStatusUpdate FbStatusType = 52

	// +enum=published_story
	PublishedStory FbStatusType = 57

	// +enum=shared_story
	SharedStory FbStatusType = 68

	// +enum=tagged_in_photo
	TaggedInPhoto FbStatusType = 73

	// +enum=wall_post
	WallPost FbStatusType = 78
)
