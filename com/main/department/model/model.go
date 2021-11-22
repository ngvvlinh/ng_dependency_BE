package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Department struct {
	ID          dot.ID `paging:"id"`
	AccountID   dot.ID
	Name        string
	Description string
	CreatedAt   time.Time `sq:"create" paging:"created_at"`
	UpdatedAt   time.Time `sq:"update" paging:"updated_at"`
	DeletedAt   time.Time
}
