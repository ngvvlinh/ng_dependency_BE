package _server

import "net/http"

type HTTPServer interface {
	http.Handler
	PathPrefix() string
}
