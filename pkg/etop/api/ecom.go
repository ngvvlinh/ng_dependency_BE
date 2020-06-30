package api

import (
	"context"
	"net/http"

	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
)

type EcomService struct{}

func (s *EcomService) Clone() api.EcomService { res := *s; return &res }

func (s *EcomService) SessionInfo(ctx context.Context, r *pbcm.Empty) (*api.EcomSessionInfoResponse, error) {
	cookies := headers.GetCookiesFromCtx(ctx)
	if cookies == nil || len(cookies) == 0 {
		return nil, cm.ErrUnauthenticated
	}
	cookie := getEcomSessionCookie(cookies)
	if cookie == nil {
		return nil, cm.ErrUnauthenticated
	}

	var err error
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSessionWithToken(ctx, cookie.Value, sessionQuery)
	if err != nil {
		return nil, err
	}
	session := sessionQuery.Result
	if session.Shop == nil {
		return nil, cm.ErrUnauthenticated
	}
	result := &api.EcomSessionInfoResponse{
		AllowAccess: true,
	}
	return result, nil
}

func getEcomSessionCookie(cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == auth.EcomAuthorization {
			return c
		}
	}
	return nil
}
