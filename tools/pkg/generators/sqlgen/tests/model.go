package tests

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type A struct {
	ID   dot.ID
	Name string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
