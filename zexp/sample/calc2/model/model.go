package model

import "o.o/capi/dot"

// +sqlgen
type Equation struct {
	ID dot.ID

	Equation string

	Result string

	CreatedAt dot.Time

	UpdatedAt dot.Time
}
