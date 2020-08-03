package api

import (
	"context"

	"o.o/capi/dot"
	"o.o/common/jsonx"
)

// +gen:apix
// +gen:swagger:doc-path=sample/counter

// +apix:path=/counter.Counter
type CounterService interface {
	Counter(context.Context, *CounterRequest) (*CounterResponse, error)

	Get(context.Context, *GetRequest) (*GetResponse, error)
}

type CounterRequest struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

func (m *CounterRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CounterResponse struct {
	Value int `json:"value"`
}

func (m *CounterResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetRequest struct {
	Label string `json:"label"`
}

func (m *GetRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetResponse struct {
	Value int `json:"value"`
}

type Counter struct {
	ID dot.ID

	Label string

	ValueOne int

	CreatedAt dot.Time

	UpdatedAt dot.Time
}

func (m *GetResponse) String() string {
	return jsonx.MustMarshalToString(m)
}
