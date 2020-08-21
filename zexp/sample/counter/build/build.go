package build

import (
	"net/http"

	"o.o/backend/pkg/common/apifw/servedoc"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/counter/config"
	"o.o/backend/zexp/sample/counter/service"
	"o.o/common/l"
)

var ll = l.New()

type MainServer *http.Server

func BuildServer(cfg config.Config, counterService service.CounterHandler) lifecycle.HTTPServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)

	mux.Handle("/", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))

	docPath := "sample/calc"
	swaggerPath := "/doc/" + docPath + "/swagger.json"
	mux.Handle("/doc/"+docPath, servedoc.RedocHandler())
	mux.Handle(swaggerPath, servedoc.SwaggerHandler(docPath+"/swagger.json"))

	mux.Handle(counterService.PathPrefix(), counterService)
	mux.Handle("/doc/sample/calc/", servedoc.RedocHandler())

	s := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	server := lifecycle.HTTPServer{
		Name:   "simple calc",
		Server: s,
	}
	return server
}
