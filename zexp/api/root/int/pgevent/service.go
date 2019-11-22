package pgevent

import (
	"context"

	cm "etop.vn/backend/pb/common"
	event "etop.vn/backend/pb/services/pgevent"
)

// +gen:apix
// +gen:apix:doc-path=services/pgevent

// +apix:path=/pgevent.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/pgevent.Event
type EventAPI interface {
	GenerateEvents(context.Context, *event.GenerateEventsRequest) (*cm.Empty, error)
}
