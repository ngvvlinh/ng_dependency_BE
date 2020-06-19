package api

// +enum
// +enum:zero=null
type EquationOperator int

const (
	// +enum=add
	Plus EquationOperator = 0

	// +enum=sub
	Minus EquationOperator = 1

	// +enum=div
	Divide EquationOperator = 2

	// +enum=mul
	Multiply EquationOperator = 3
)
