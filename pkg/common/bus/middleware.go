package bus

import "net/http"

func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewRootContext(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
