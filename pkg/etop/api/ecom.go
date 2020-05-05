package api

import (
	"context"
	"net/http"

	"o.o/api/top/int/etop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/headers"
	authservice "o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
)

func (s *EcomService) SessionInfo(ctx context.Context, r *EcomSessionInfoEndpoint) error {
	cookies := headers.GetCookiesFromCtx(ctx)
	if cookies == nil || len(cookies) == 0 {
		return cm.ErrUnauthenticated
	}
	cookie := getEcomSessionCookie(cookies)
	if cookie == nil {
		return cm.ErrUnauthenticated
	}

	var err error
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSessionWithToken(ctx, cookie.Value, sessionQuery)
	if err != nil {
		return err
	}
	session := sessionQuery.Result
	if session.Shop == nil {
		return cm.ErrUnauthenticated
	}
	r.Result = &etop.EcomSessionInfoResponse{
		AllowAccess: true,
	}
	return nil
}

func getEcomSessionCookie(cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == authservice.EcomAuthorization {
			return c
		}
	}
	return nil
}
