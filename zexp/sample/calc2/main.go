package main

import (
	"context"
	"net/http"
	"os"
	"time"

	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/calc2/config"
	"o.o/backend/zexp/sample/calc2/service"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	// load config
	cfg, err := config.Load()
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment(cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sdCtx.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))

	docPath := "sample/calc"
	swaggerPath := "/doc/" + docPath + "/swagger.json"
	mux.Handle("/doc/"+docPath, cmservice.RedocHandler())
	mux.Handle(swaggerPath, cmservice.SwaggerHandler(docPath+"/swagger.json"))

	// connect db
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Error(err.Error())
		os.Exit(2)
	}

	sv := service.NewCalcService(db)
	ser, err := httprpc.NewServer(sv.Clone)
	if err != nil {
		ll.Error(err.Error())
		os.Exit(2)
	}

	mux.Handle(ser.PathPrefix(), ser)
	mux.Handle("/doc/sample/calc/", cmservice.RedocHandler())

	s := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	server := lifecycle.HTTPServer{
		Name:   "simple calc",
		Server: s,
	}
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, server)
	sdCtx.Register(cancelHTTP)
}
