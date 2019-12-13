package api

import (
	"context"

	"etop.vn/api/top/int/etop"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (*etop.LoginResponse, error) {
	return userService.CreateLoginResponse(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
}
