package middleware

import (
	"net/http"
	"strings"

	"o.o/backend/pkg/common/cmenv"
)

func CORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		switch {
		case
			origin == "ionic://localhost",
			origin == "capacitor://localhost",
			origin == "http://localhost",
			origin == "http://localhost:8080",
			origin == "http://localhost:8100",
			strings.HasSuffix(origin, ".localhost:8100"),
			strings.HasSuffix(origin, ".ecomify.vn"),
			strings.HasSuffix(origin, ".ecom.d.etop.vn"),
			strings.HasSuffix(origin, ".d.etop.vn"):
			w.Header().Set("Access-Control-Allow-Origin", origin)

		case cmenv.IsSandBox(), cmenv.IsDevOrStag():
			w.Header().Set("Access-Control-Allow-Origin", "*")

		case cmenv.IsProd():

		default:
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Add("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	}
}
