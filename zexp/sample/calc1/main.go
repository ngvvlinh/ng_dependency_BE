package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/calc1/config"
	"o.o/backend/zexp/sample/calc1/service"
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
	cmenv.SetEnvironment("calc1", cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sdCtx.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

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
