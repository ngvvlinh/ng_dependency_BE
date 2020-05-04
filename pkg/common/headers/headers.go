package headers

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"o.o/backend/pkg/common/cmenv"
)

type authKey struct{}
type debugKey struct{}
type headerKey struct{}
type CookieKey struct{}

func ForwardHeaders(next http.Handler, configs ...Config) http.HandlerFunc {
	var cfg Config
	if len(configs) > 0 {
		cfg = configs[0]
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, headerKey{}, filterHeader(r.Header))

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
		if cmenv.IsDev() {
			debug := r.Header.Get("debug")
			ctx = context.WithValue(ctx, debugKey{}, debug)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func filterHeader(header http.Header) http.Header {
	result := make(http.Header)
	for key, vals := range header {
		if key == "Cookie" || key == "Authorization" {
			continue
		}
		result[key] = vals
	}
	return result
}

func CtxDebug(ctx context.Context) string {
	v := ctx.Value(debugKey{})
	if v == nil {
		return ""
	}
	return v.(string)
}

func GetHeader(ctx context.Context) http.Header {
	v := ctx.Value(headerKey{})
	if v == nil {
		return nil
	}
	return v.(http.Header)
}

// GetBearerTokenFromCtx ...
func GetBearerTokenFromCtx(ctx context.Context) string {
	authHeader, ok := ctx.Value(authKey{}).(string)
	if !ok {
		return ""
	}
	token, _ := authFromHeaderString(authHeader)
	return token
}

type Config struct {
	AllowQueryAuthorization bool
}

type HeaderItem struct {
	Key    string
	Values []string
}

func GetSortedHeaders(ctx context.Context) []HeaderItem {
	header := GetHeader(ctx)
	if header == nil {
		return nil
	}
	result := make([]HeaderItem, 0, len(header))
	for key, vals := range header {
		result = append(result, HeaderItem{Key: key, Values: vals})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

// FromHeaderString ...
func authFromHeaderString(s string) (string, error) {
	splits := strings.SplitN(s, " ", 2)
	if len(splits) < 2 {
		return "", errors.New("bad authorization string")
	}
	if splits[0] != "Bearer" && splits[0] != "bearer" {
		return "", errors.New("request unauthenticated with Bearer")
	}
	return splits[1], nil
}
