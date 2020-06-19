package api

import (
	"context"

	"o.o/capi/dot"
	"o.o/common/jsonx"
)

// +gen:apix
// +gen:swagger:doc-path=sample/calc

// +apix:path=/calc.Calc
type CalcService interface {
	Calc(context.Context, *Request) (*Response, error)

	Get(context.Context, *GetRequest) (*EquationResponse, error)
}

type Request struct {
	// we only support "add" at the only operator
	Op dot.NullString `json:"op"`

	// the first param
	A dot.NullString `json:"a"`

	// the second param
	B dot.NullString `json:"b"`
}

func (m *Request) String() string { return jsonx.MustMarshalToString(m) }

type Response struct {
	ID dot.ID `json:"id"`

	Result string `json:"result"`
}

func (m *Response) String() string { return jsonx.MustMarshalToString(m) }

type GetRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetRequest) String() string { return jsonx.MustMarshalToString(m) }

type EquationResponse struct {
	Equation *Equation `json:"equation"`
}

type Equation struct {
	ID dot.ID `json:"id"`

	Equation string `json:"equation"`

	Result string `json:"result"`

	CreatedAt dot.Time `json:"created_at"`

	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *EquationResponse) String() string { return jsonx.MustMarshalToString(m) }
