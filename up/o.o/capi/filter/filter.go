package filter

import (
	"o.o/capi/dot"
)

// Filter by id value
//
// +swagger:type=string
type IDs []dot.ID

// Filter by string value
//
// +swagger:type=string
type Strings []string

// Filter by date value
//
// Object:
//
//     - {"from":"2020-01-01T00:00:00Z","to":"2020-02-01T00:00:00Z"}
//
// String:
//
//     - "2020-01-01T00:00:00Z,2010-02-01T00:00:00Z":
//     - "2020-01-01T00:00:00Z,":
//     - ",2020-01-01T00:00:00Z":
//
// +swagger:type=string
type Date date

type date struct {
	From dot.Time `json:"from"`
	To   dot.Time `json:"to"`
}

func (d Date) IsZero() bool {
	return d.From.IsZero() && d.To.IsZero()
}

// +swagger:type=string
type FullTextSearch string

func (f FullTextSearch) String() string { return string(f) }
