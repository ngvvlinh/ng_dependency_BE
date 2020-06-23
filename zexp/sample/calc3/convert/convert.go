package convert

import (
	"strings"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/zexp/sample/calc3/api"
	calc3model "o.o/backend/zexp/sample/calc3/model"
	"o.o/capi/dot"
)

// +gen:convert: o.o/backend/zexp/sample/calc3/model  -> o.o/backend/zexp/sample/calc3/api
// +gen:convert: o.o/backend/zexp/sample/calc3/api

func equation(arg *calc3model.Equation, out *api.Equation) {
	convert_calc3model_Equation_api_Equation(arg, out)
	out.Equation = strings.TrimSpace(out.Equation)
}

func createEquationConvert(arg *api.CreateEquationRequest, out *api.Equation) {
	apply_api_CreateEquationRequest_api_Equation(arg, out)
	out.ID = cm.NewID()
	out.CreatedAt = dot.Time(time.Now())
	out.UpdatedAt = dot.Time(time.Now())
	out.ProcessCalc(arg.A, arg.B, arg.Op)
}

func updateEquationConvert(arg *api.UpdateEquationRequest, out *api.Equation) {
	apply_api_UpdateEquationRequest_api_Equation(arg, out)
	out.UpdatedAt = dot.Time(time.Now())
}
