package build

import (
	"net/http"

	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/calc3/config"
	"o.o/backend/zexp/sample/calc3/service"
	"o.o/common/l"
)

var ll = l.New()

type MainServer *http.Server

func BuildServer(cfg config.Config, calcService service.CalcHandler) lifecycle.HTTPServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)

	mux.Handle("/", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/sample/calc", http.StatusTemporaryRedirect))

	docPath := "sample/calc"
	swaggerPath := "/doc/" + docPath + "/swagger.json"
	mux.Handle("/doc/"+docPath, cmservice.RedocHandler())
	mux.Handle(swaggerPath, cmservice.SwaggerHandler(docPath+"/swagger.json"))

	mux.Handle(calcService.PathPrefix(), calcService)
	mux.Handle("/doc/sample/calc/", cmservice.RedocHandler())

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
