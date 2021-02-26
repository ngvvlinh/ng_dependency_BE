package fb_live_video_status

// +enum
// +enum:zero=null
type FbLiveVideoStatus int

type NullFbLiveVideoStatus struct {
	Enum  FbLiveVideoStatus
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbLiveVideoStatus = 0

	// +enum=created
	Created FbLiveVideoStatus = 54

	// +enum=live
	Live FbLiveVideoStatus = 63

	// +enum=live_stopped
	LiveStopped FbLiveVideoStatus = 97

	// +enum=cancelled
	Cancelled FbLiveVideoStatus = 100
)

func ConvertToFbLiveVideoStatus(status string) FbLiveVideoStatus {
	//SCHEDULED_UNPUBLISHED, SCHEDULED_LIVE, SCHEDULED_EXPIRED, SCHEDULED_CANCELED
	switch status {
	case "UNPUBLISHED", "SCHEDULED_UNPUBLISHED", "SCHEDULED_LIVE":
		return Created
	case "LIVE":
		return Live
	case "PROCESSING", "LIVE_STOPPED", "VOD":
		return LiveStopped
	case "SCHEDULED_EXPIRED", "SCHEDULED_CANCELED":
		return Cancelled
	}
	return Unknown
}
