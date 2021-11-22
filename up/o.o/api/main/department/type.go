package department

import (
	"time"

	"o.o/capi/dot"
)

// +gen:event:topic=event/department

type Department struct {
	ID          dot.ID
	AccountID   dot.ID
	Name        string
	Description string
	Count       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
