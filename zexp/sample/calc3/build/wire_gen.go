// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package build

import (
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/calc3/config"
	"o.o/backend/zexp/sample/calc3/service"
)

// Injectors from wire.go:

func Build(cfg config.Config) (lifecycle.HTTPServer, error) {
	configPostgres := cfg.Postgres
	database, err := service.GetDBConnection(configPostgres)
	if err != nil {
		return lifecycle.HTTPServer{}, err
	}
	calcService := service.NewCalcService(database)
	calcHandler, err := service.BuildCalcHandlers(calcService)
	if err != nil {
		return lifecycle.HTTPServer{}, err
	}
	httpServer := BuildServer(cfg, calcHandler)
	return httpServer, nil
}