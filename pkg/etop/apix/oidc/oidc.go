package oidc

import (
	api "o.o/api/top/int/etop"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/authorize/middleware"
)

type OIDCService struct {
	SS          middleware.SessionStarter
	Userservice root.UserService
}

func (s *OIDCService) Callback(c *httpx.Context) error {
	ctx := c.Context()
	code := c.Req.URL.Query().Get("code")
	req := api.VerifyTokenUsingCodeRequest{
		Code: code,
	}

	res, err := s.Userservice.VerifyTokenUsingCode(ctx, &req)
	if err != nil {
		return err
	}

	c.SetResult(res)
	return nil
}
