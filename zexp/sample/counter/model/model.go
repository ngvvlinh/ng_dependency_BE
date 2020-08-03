package model

import "o.o/capi/dot"

// +sqlgen
type Counter struct {
	ID dot.ID

	Label string

	Value int

	CreatedAt dot.Time

	UpdatedAt dot.Time
}
