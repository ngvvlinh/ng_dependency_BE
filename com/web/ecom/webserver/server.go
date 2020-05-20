package webserver

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"o.o/api/main/catalog"
	"o.o/api/main/location"
	"o.o/api/subscripting/subscription"
	"o.o/api/webserver"
	"o.o/backend/com/web/ecom/middlewares"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/redis"
)

var webserverQueryBus webserver.QueryBus
var catelogQueryBus catalog.QueryBus
var subscriptionQuery subscription.QueryBus

type Config struct {
	MainSite string
	CoreSite string
	RootPath string
}

var config Config

type Server struct {
	echo *echo.Echo
	tpl  *Templates
	cfg  Config
}

var locationBus location.QueryBus

func New(cfg Config, query webserver.QueryBus, catalogQuery catalog.QueryBus, rd redis.Store, locationQueryBus location.QueryBus, subrQuery subscription.QueryBus) (*Server, error) {
	locationBus = locationQueryBus
	redisStore = rd
	webserverQueryBus = query
	catelogQueryBus = catalogQuery
	subscriptionQuery = subrQuery
	if cfg.MainSite == "" {
		return nil, cm.Errorf(cm.Internal, nil, "missing main_site")
	}
	if cfg.RootPath == "" {
		return nil, cm.Errorf(cm.Internal, nil, "missing root_path")
	}
	if cfg.CoreSite == "" {
		return nil, cm.Errorf(cm.Internal, nil, "missing core_site")
	}
	config = cfg

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.SiteRouter)
	e.Static("/assets", filepath.Join(cfg.RootPath, "com/web/ecom/assets"))
	e.Static("/fonts", filepath.Join(cfg.RootPath, "com/web/ecom/templates/fonts"))
	e.Static("/css", filepath.Join(cfg.RootPath, "com/web/ecom/templates/css"))
	e.Static("/js", filepath.Join(cfg.RootPath, "com/web/ecom/templates/js"))
	e.Static("/images", filepath.Join(cfg.RootPath, "com/web/ecom/templates/images"))
	e.Static("/libs", filepath.Join(cfg.RootPath, "com/web/ecom/templates/libs"))

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
	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	s.echo.ServeHTTP(w, r)
}

func (s *Server) registerHandlers(e *echo.Echo) {
	e.GET("/", s.Index)

	e.GET("/checkout", s.Checkout)
	e.GET("/product/:id", s.Product)
	e.GET("/category/:id/:page", s.Category)
	e.GET("/category/:id/:page/:perpage", s.CategoryWithPagePaging)

	e.GET("/search/:search/:page/:perpage", s.Search)

	e.GET("/about-us", s.AboutUs)
	e.GET("/page/:id", s.Page)

	e.PUT("/cart/remove", s.CartRemoveVariant)
	e.GET("/order", s.CartOrder)
	e.POST("/cart/add-one-product", s.CartQuickAddProduct)
	e.POST("/cart", s.CartAddProduct)
	e.PUT("/cart", s.CartUpdateAllListProduct)
	e.GET("/cart", s.Cart)
	e.POST("/checkout/create-order", s.CreateOrder)
	e.POST("/cart/total-count", s.CartTotalCount)

	e.POST("/provinces", s.Provinces)
	e.POST("/districts", s.Districts)
	e.POST("/wards", s.Wards)

	e.GET("/subscription-outdated", func(c echo.Context) error {
		return c.Render(http.StatusOK, "subscription-outdated.html", nil)
	})
}
