package webserver

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"o.o/api/webserver"
	"o.o/backend/com/web/ecom/middlewares"
	cm "o.o/backend/pkg/common"
)

var webserverQueryBus webserver.QueryBus
var config Config

type Config struct {
	MainSite string
	RootPath string
}

type Server struct {
	echo *echo.Echo
	tpl  *Templates
	cfg  Config
}

func New(cfg Config, query webserver.QueryBus) (*Server, error) {
	webserverQueryBus = query
	if cfg.MainSite == "" {
		return nil, cm.Errorf(cm.Internal, nil, "missing main_site")
	}
	if cfg.RootPath == "" {
		return nil, cm.Errorf(cm.Internal, nil, "missing root_path")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.SiteRouter)
	e.Static("/assets", filepath.Join(cfg.RootPath, "com/web/ecom/assets"))
	e.Static("/fonts", filepath.Join(cfg.RootPath, "com/web/ecom/templates/fonts"))
	e.Static("/css", filepath.Join(cfg.RootPath, "com/web/ecom/templates/css"))
	e.Static("/js", filepath.Join(cfg.RootPath, "com/web/ecom/templates/js"))
	e.Static("/images", filepath.Join(cfg.RootPath, "com/web/ecom/templates/images"))
	e.Static("/vendor", filepath.Join(cfg.RootPath, "com/web/ecom/templates/vendor"))

	templates, err := parseTemplates(filepath.Join(cfg.RootPath, "com/web/ecom/templates", "*.html"))
	if err != nil {
		return nil, err
	}
	e.Renderer = templates

	s := &Server{echo: e, cfg: cfg, tpl: templates}
	s.registerHandlers(e)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}

func (s *Server) registerHandlers(e *echo.Echo) {
	e.GET("/", s.Index)
}
