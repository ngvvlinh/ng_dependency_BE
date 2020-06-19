package service

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/calc2/api"
	"o.o/backend/zexp/sample/calc2/convert"
	"o.o/backend/zexp/sample/calc2/sqlstore"
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

func (c *CalcService) Calc(ctx context.Context, req *api.Request) (*api.Equation, error) {
	eq := &api.Equation{
		ID: cm.NewID(),
	}
	err := eq.ProcessCalc(req.A, req.B, req.Op)
	err = c.Store(ctx).CreateEquation(eq)

	if err != nil {
		return nil, err
	}
	return eq, nil
}

func (c *CalcService) Get(ctx context.Context, request *api.GetRequest) (*api.Equation, error) {
	id := request.ID
	eq, err := c.Store(ctx).ID(id).GetEquation()
	if err != nil {
		return nil, err
	}
	return eq, nil
}

func (c *CalcService) Update(ctx context.Context, request *api.UpdateEquationRequest) (*api.Equation, error) {
	id := request.ID
	eq, err := c.Store(ctx).ID(id).GetEquation()
	if err != nil {
		return nil, err
	}

	eq = convert.Apply_api_UpdateEquationRequest_api_Equation(request, eq)
	err = eq.ProcessCalc(request.A, request.B, request.Op)
	if err != nil {
		return nil, err
	}

	_, err = c.Store(ctx).UpdateEquation(eq)
	if err != nil {
		return nil, err
	}

	return eq, nil
}

func (c *CalcService) List(ctx context.Context, request *api.ListEquationRequest) (*api.Equations, error) {
	return c.Store(ctx).Filters(cmapi.ToFilters(request.Filters)).ListEquation()
}
