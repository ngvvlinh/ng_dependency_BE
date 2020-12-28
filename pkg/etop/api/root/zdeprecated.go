package root

import (
	"context"

	"o.o/api/top/int/etop"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/capi/dot"
)

func CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *identitymodel.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (*etop.LoginResponse, error) {
	return UserServiceImpl.CreateLoginResponse(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
}
