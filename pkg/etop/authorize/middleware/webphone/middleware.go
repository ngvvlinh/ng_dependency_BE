package webphone

import (
	"net/http"
	"strings"

	config_server "o.o/backend/cogs/config/_server"
)

const (
	XRequestClientHeader = "X-Request-Client"
)

func CORS(webphonePublicKey config_server.WebphonePublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Bypass preflight request (OPTION Method)
			// We need to allow CORS for request from webphone (extension chrome)
			// by identify the header X-Request-Client
			origin := r.Header.Get("origin")
			if r.Method == "OPTIONS" {
				requestHeaders := r.Header.Get("Access-Control-Request-Headers")
				if strings.Contains(requestHeaders, strings.ToLower(XRequestClientHeader)) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				next.ServeHTTP(w, r)
				return
			}

			requestClient := r.Header.Get(XRequestClientHeader)
			if requestClient == "" {
				next.ServeHTTP(w, r)
				return
			}

			if requestClient == string(webphonePublicKey) && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}
