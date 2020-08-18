package api

import (
	"context"
	fmt "fmt"
	"strconv"

	"o.o/api/top/types/common"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

// +gen:apix
// +gen:swagger:doc-path=zext/sample/calc

// +apix:path=/calc.Calc
type CalcService interface {
	Calc(context.Context, *Request) (*Equation, error)

	Get(context.Context, *GetRequest) (*Equation, error)

	List(context.Context, *ListEquationRequest) (*Equations, error)

	Update(context.Context, *UpdateEquationRequest) (*Equation, error)
}

type ListEquationRequest struct {
	Filters []*common.Filter `json:"filters"`
}

func (m *ListEquationRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

// +convert:Create=Equation
type Request struct {
	// we only support "add" at the only operator
	Op EquationOperator `json:"op"`

	// the first param
	A dot.NullString `json:"a"`

	// the second param
	B dot.NullString `json:"b"`
}

func (m *Request) String() string { return jsonx.MustMarshalToString(m) }

type GetRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetRequest) String() string { return jsonx.MustMarshalToString(m) }

type Equations struct {
	Equations []*Equation `json:"equations"`
}

type Equation struct {
	ID dot.ID `json:"id"`

	Equation string `json:"equation"`

	Result string `json:"result"`

	CreatedAt dot.Time `json:"created_at"`

	UpdatedAt dot.Time `json:"updated_at"`
}

func (eq *Equation) ProcessCalc(num1 dot.NullString, num2 dot.NullString, op EquationOperator) error {

	var opStr string
	switch op {
	case Divide:
		opStr = "/"
	case Multiply:
		opStr = "*"
	case Minus:
		opStr = "-"
	case Plus:
		opStr = "+"
	}

	equation := num1.Apply("") + " " + opStr + " " + num2.Apply("")
	result := "unknown"

	fNum1, err := strconv.ParseFloat(num1.Apply(""), 64)
	if err != nil {
		return err
	}

	fNum2, err := strconv.ParseFloat(num2.Apply(""), 64)
	if err != nil {
		return err
	}

	if op == Plus {
		result = fmt.Sprintf("%f", fNum1+fNum2)
	}

	if op == Minus {
		result = fmt.Sprintf("%f", fNum1-fNum2)
	}

	if op == Multiply {
		result = fmt.Sprintf("%f", fNum1*fNum2)
	}

	if op == Divide {
		result = fmt.Sprintf("%f", fNum1/fNum2)
	}

	eq.Result = result
	eq.Equation = equation
	return nil
}

func (m *Equation) String() string { return jsonx.MustMarshalToString(m) }

func (m *Equations) String() string { return jsonx.MustMarshalToString(m) }

// +convert:update=Equation
type UpdateEquationRequest struct {
	// id of equation
	ID dot.ID `json:"id"`

	// the operator (+, -, *, /)
	Op EquationOperator `json:"op"`

	// the first param
	A dot.NullString `json:"a"`

	// the second param
	B dot.NullString `json:"b"`
}

func (m *UpdateEquationRequest) String() string { return jsonx.MustMarshalToString(m) }
