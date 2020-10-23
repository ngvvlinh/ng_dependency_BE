package middlewares

import (
	"net/http"

	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/session"
)

func Authorization(_ss session.Session) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string

			// clone session for each request
			ctx := r.Context()
			ss := _ss
			perms := permission.Decl{Type: permission.Shop}

			cookies := r.Cookies()
			for _, cookie := range cookies {
				if cookie.Name == auth.Authorization {
					tokenStr = cookie.Value
				}
			}

			if _, err := ss.StartSession(ctx, perms, tokenStr); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			newCtx := session.NewCtxWithSession(ctx, ss)
			newR := r.WithContext(newCtx)
			next.ServeHTTP(w, newR)
		}
	}
}
