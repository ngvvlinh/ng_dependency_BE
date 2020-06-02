package httpx

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func Compose(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		handler := next
		for _, m := range middlewares {
			handler = m(next)
		}
		return handler
	}
}

type Server interface {
	http.Handler
	PathPrefix() string
}

type server struct {
	http.Handler
	path string
}

func MakeServer(path string, handler http.Handler) Server {
	return &server{Handler: handler, path: path}
}

func (s *server) PathPrefix() string {
	return s.path
}
