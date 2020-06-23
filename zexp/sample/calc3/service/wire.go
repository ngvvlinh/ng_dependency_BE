package service

import (
	cc "o.o/backend/pkg/common/config"

	"github.com/google/wire"

	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

var WireSet = wire.NewSet(
	BuildCalcHandlers,
	GetDBConnection,
	NewCalcService,
)

type CalcHandler httpx.Server

func GetDBConnection(pgCfg cc.Postgres) (*cmsql.Database, error) {
	db, err := cmsql.Connect(pgCfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BuildCalcHandlers(service *CalcService) (CalcHandler, error) {
	ser, err := httprpc.NewServer(service.Clone)
	if err != nil {
		return nil, err
	}
	return ser, nil
}
