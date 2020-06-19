package convert

import (
	"strings"

	"o.o/backend/zexp/sample/calc2/api"
	calc2model "o.o/backend/zexp/sample/calc2/model"
)

// +gen:convert: o.o/backend/zexp/sample/calc2/model  -> o.o/backend/zexp/sample/calc2/api
// +gen:convert: o.o/backend/zexp/sample/calc2/api

func equation(arg *calc2model.Equation, out *api.Equation) {
	convert_calc2model_Equation_api_Equation(arg, out)
	out.Equation = strings.TrimSpace(out.Equation)
}
