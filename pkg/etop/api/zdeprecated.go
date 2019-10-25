package api

import (
	"context"

	pbetop "etop.vn/backend/pb/etop"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/model"
)

func CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID int64, user *model.User, preferAccountID int64, preferAccountType int, generateAllTokens bool, adminID int64) (*pbetop.LoginResponse, error) {
	return s.CreateLoginResponse(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
}
