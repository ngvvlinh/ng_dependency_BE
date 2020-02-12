package middleware

import (
	"context"
	"net/http"

	"etop.vn/backend/pkg/common/apifw/captcha"
	"etop.vn/backend/pkg/common/cmenv"
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
			origin == "http://localhost:8100":
			w.Header().Add("Access-Control-Allow-Origin", origin)

		case cmenv.IsDev():
			w.Header().Add("Access-Control-Allow-Origin", "*")

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

func VerifyCaptcha(ctx context.Context, token string) error {
	return captcha.Verify(token)
}
