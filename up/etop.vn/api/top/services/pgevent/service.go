package pgevent

import (
	"context"

	cm "etop.vn/api/pb/common"
	event "etop.vn/api/pb/services/pgevent"
)

// +gen:apix
// +gen:swagger:doc-path=services/pgevent

// +apix:path=/pgevent.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/pgevent.Event
type EventService interface {
	GenerateEvents(context.Context, *event.GenerateEventsRequest) (*cm.Empty, error)
}
