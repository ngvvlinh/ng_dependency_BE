package middleware

import (
	"context"
	"net/http"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/captcha"
)

type authKey struct{}
type debugKey struct{}

type Config struct {
	AllowQueryAuthorization bool
}

func ForwardHeaders(next http.Handler, configs ...Config) http.HandlerFunc {
	var cfg Config
	if len(configs) > 0 {
		cfg = configs[0]
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authorization := r.Header.Get("Authorization")
		if cfg.AllowQueryAuthorization && authorization == "" {
			token := r.URL.Query().Get("__token")
			if token != "" {
				authorization = "Bearer " + token
			}
		}
		if authorization != "" {
			ctx = context.WithValue(ctx, authKey{}, authorization)
		}
		if cm.IsDev() {
			debug := r.Header.Get("debug")
			ctx = context.WithValue(ctx, debugKey{}, debug)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func CtxDebug(ctx context.Context) string {
	v := ctx.Value(debugKey{})
	if v == nil {
		return ""
	}
	return v.(string)
}

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

		case cm.IsDev():
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
