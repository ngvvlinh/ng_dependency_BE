package service

import (
	"context"
	"fmt"
	"strconv"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/calc1/api"
	"o.o/backend/zexp/sample/calc1/sqlstore"
)

var _ api.CalcService = &CalcService{}

type CalcService struct {
	Store sqlstore.EquationStoreFactory
}

func (c *CalcService) Clone() api.CalcService {
	cl := *c
	return &cl
}

func NewCalcService(db *cmsql.Database) *CalcService {
	return &CalcService{
		Store: sqlstore.NewEquationStore(db),
	}
}

func (c *CalcService) Calc(ctx context.Context, req *api.Request) (*api.Response, error) {
	num1 := req.A.Apply("")
	num2 := req.B.Apply("")
	op := req.Op.Apply("")
	equation := num1 + " " + op + " " + num2
	result := "unknown"

	fNum1, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return nil, err
	}

	fNum2, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return nil, err
	}

	if op == "+" {
		result = fmt.Sprintf("%f", fNum1+fNum2)
	}

	if op == "-" {
		result = fmt.Sprintf("%f", fNum1-fNum2)
	}

	if op == "*" {
		result = fmt.Sprintf("%f", fNum1*fNum2)
	}

	if op == "/" {
		result = fmt.Sprintf("%f", fNum1/fNum2)
	}

	eq := &api.Equation{
		ID:       cm.NewID(),
		Equation: equation,
		Result:   result,
	}
	err = c.Store(ctx).CreateEquation(eq)

	if err != nil {
		return nil, err
	}
	return &api.Response{
		ID:     eq.ID,
		Result: eq.Result,
	}, nil
}

func (c *CalcService) Get(ctx context.Context, request *api.GetRequest) (*api.EquationResponse, error) {
	id := request.ID
	eq, err := c.Store(ctx).ID(id).GetEquation()
	if err != nil {
		fmt.Println(err.Error())
	}
	return &api.EquationResponse{Equation: eq}, nil
}
