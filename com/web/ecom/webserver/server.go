package webserver

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"etop.vn/backend/com/web/ecom/middlewares"
	cm "etop.vn/backend/pkg/common"
)

type Config struct {
	MainSite string
	RootPath string
}

type Server struct {
	echo *echo.Echo
	tpl  *Templates
	cfg  Config
}

func New(cfg Config) (*Server, error) {
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
