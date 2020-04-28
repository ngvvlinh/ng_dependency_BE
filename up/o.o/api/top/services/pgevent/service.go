package pgevent

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=services/pgevent

// +apix:path=/pgevent.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/pgevent.Event
type EventService interface {
	GenerateEvents(context.Context, *GenerateEventsRequest) (*cm.Empty, error)
}
