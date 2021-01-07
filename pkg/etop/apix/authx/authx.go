package authx

import (
	"strings"

	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/middleware"
)

type AuthxService struct {
	St middleware.SessionStarter
}

func (s *AuthxService) AuthUser(c *httpx.Context) error {
	ctx := c.Context()
	tokenStr := headers.GetBearerTokenFromCtx(ctx)

	if tokenStr == "" {
		return cm.Errorf(cm.Unauthenticated, nil, "")
	}

	var claim *claims.Claim
	var err error
	var account identitymodel.AccountInterface
	switch {
	case !strings.Contains(tokenStr, ":"):
		// token internal: shop
		claim, err = s.St.TokenStore.Validate(tokenStr)
		if err != nil {
			return cm.ErrUnauthenticated
		}

		if !idutil.IsShopID(claim.AccountID) {
			return cm.ErrUnauthenticated
		}
		query := &identitymodelx.GetShopQuery{
			ShopID: claim.AccountID,
		}
		if err := s.St.ShopStore.GetShop(ctx, query); err != nil {
			return cm.ErrUnauthenticated
		}
		account = query.Result

	case strings.HasPrefix(tokenStr, "shop"):
		// token partner shop (APIPartnerShopKey)
		claim, account, err = s.St.VerifyAPIPartnerShopKey(ctx, tokenStr)
		if err != nil {
			return err
		}

	default:
		// token partner/shop (APIKey)
		claim, account, err = s.St.VerifyAPIKey(ctx, tokenStr, account_type.Shop)
		if err != nil {
			_, account, err = s.St.VerifyAPIKey(ctx, tokenStr, account_type.Partner)
			if err != nil {
				return err
			}
		}
	}

	acc := account.GetAccount()
	res := &AuthxResponse{
		User: &AuthxUser{
			ID: acc.OwnerID,
		},
		Account: &AuthxAccount{
			ID:      acc.ID,
			Name:    acc.Name,
			OwnerID: acc.OwnerID,
			Type:    acc.Type,
		},
	}
	c.SetResult(res)
	return nil
}
