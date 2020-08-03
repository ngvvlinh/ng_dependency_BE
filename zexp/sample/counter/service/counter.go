package service

import (
	"context"
	"fmt"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/counter/api"
	"o.o/backend/zexp/sample/counter/sqlstore"
)

var _ api.CounterService = &CounterService{}

type CounterService struct {
	Store sqlstore.CounterStoreFactory
}

func (c *CounterService) Clone() api.CounterService {
	cl := *c
	return &cl
}

func NewCounterService(db *cmsql.Database) *CounterService {
	return &CounterService{
		Store: sqlstore.NewCounterStore(db),
	}
}

func (c *CounterService) Counter(ctx context.Context, req *api.CounterRequest) (*api.CounterResponse, error) {
	counter, err := c.Store(ctx).Label(req.Label).GetCounter()
	if err != nil {
		co := &api.Counter{
			ID:       cm.NewID(),
			Label:    req.Label,
			ValueOne: req.Value,
		}

		if err := c.Store(ctx).CreateCounter(co); err != nil {
			return nil, err
		}
		return &api.CounterResponse{
			Value: req.Value,
		}, nil
	}
	if _, err := c.Store(ctx).Label(req.Label).UpdateLabelvalue(counter.ValueOne + req.Value); err != nil {
		return nil, err
	}

	return &api.CounterResponse{
		Value: counter.ValueOne + req.Value,
	}, nil
}

func (c *CounterService) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	label := request.Label
	co, err := c.Store(ctx).Label(label).GetCounter()
	if err != nil {
		fmt.Println(err.Error())
	}
	return &api.GetResponse{Value: co.ValueOne}, nil
}
