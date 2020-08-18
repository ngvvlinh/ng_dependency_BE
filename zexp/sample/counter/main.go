package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	health "o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/servedoc"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/zexp/sample/counter/config"
	"o.o/backend/zexp/sample/counter/service"
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
	cmenv.SetEnvironment("counter", cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sdCtx.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})
	mux.HandleFunc(health.DefaultRoute, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	mux.Handle("/", http.RedirectHandler("/doc/sample/counter", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/sample/counter", http.StatusTemporaryRedirect))

	docPath := "sample/counter"
	swaggerPath := "/doc/" + docPath + "/swagger.json"
	mux.Handle("/doc/"+docPath, servedoc.RedocHandler())
	mux.Handle(swaggerPath, servedoc.SwaggerHandler("zext/"+docPath+"/swagger.json"))

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Error(err.Error())
		os.Exit(2)
	}

	sv := service.NewCounterService(db)
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
		Name:   "simple counter",
		Server: s,
	}
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, server)
	sdCtx.Register(cancelHTTP)
}
