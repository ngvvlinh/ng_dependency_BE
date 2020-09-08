package integration

import (
	"context"
	"fmt"
	"strings"

	api "o.o/api/top/int/integration"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/validate"
	apipartner "o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/common/jsonx"
)

func (s *IntegrationService) LoginUsingTokenWL(ctx context.Context, r *api.LoginUsingTokenRequest) (_ *api.LoginResponse, _err error) {
	if r.Login == "" || r.VerificationCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin")
	}
	r.VerificationCode = strings.Replace(r.VerificationCode, " ", "", -1)
	if _, err := s.validateWhiteLabel(ctx); err != nil {
		return nil, err
	}
	claim := s.SS.Claim()
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	verifiedEmail, verifiedPhone, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email hoặc số điện thoại không hợp lệ")
	}

	partner := s.SS.CtxPartner()
	key := fmt.Sprintf("%v-%v", partner.ID, cm.Coalesce(verifiedEmail, verifiedPhone))

	var v map[string]string
	tok, err := s.AuthStore.Validate(auth.UsageRequestLogin, key, &v)
	if err != nil {
		return nil, cm.Errorf(cm.Unauthenticated, nil, "Mã xác nhận không hợp lệ")
	}
	if v == nil || v["code"] != r.VerificationCode {
		return nil, cm.Errorf(cm.Unauthenticated, nil, "Mã xác thực không hợp lệ")
	}
	// delete the token after 5 minutes if login successfully
	defer func() {
		if _err != nil {
			_ = s.AuthStore.SetTTL(tok, 5*60)
		}
	}()

	actionUser := ""
	var user *identitymodel.User
	if requestInfo.ExternalUserID != "" {
		relationUserQuery := &identitymodelx.GetPartnerRelationQuery{
			PartnerID:      partner.ID,
			ExternalUserID: requestInfo.ExternalUserID,
		}
		err = s.PartnerStore.GetPartnerRelationQuery(ctx, relationUserQuery)
		switch cm.ErrorCode(err) {
		case cm.OK:
			actionUser = "ok"
			user = relationUserQuery.Result.User
		case cm.NotFound:
			actionUser = "create"
		default:
			return nil, err
		}
	} else {
		userQuery := &identitymodelx.GetUserByLoginQuery{
			PhoneOrEmail: r.Login,
		}
		err = s.UserStore.GetUserByLogin(ctx, userQuery)
		switch cm.ErrorCode(err) {
		case cm.OK:
			actionUser = "ok"
			user = userQuery.Result.User
		case cm.NotFound:
			actionUser = "create"
		default:
			return nil, err
		}
	}

	switch actionUser {
	case "ok":
		// continue
	case "create":
		// --- register user ---
		// trust thông tin từ đối tác gửi qua
		// xem như email, phone đều đã xác thực
		newUser, err := s.registerUser(ctx, false, partner.ID, requestInfo.ShopOwnerName, requestInfo.ShopOwnerEmail, requestInfo.ShopOwnerPhone, true, true, true, true)
		if err != nil {
			return nil, err
		}
		user = newUser
		if requestInfo.ExternalUserID != "" {
			// create partner relation with user
			relationCmd := &identitymodelx.CreatePartnerRelationCommand{
				UserID:     user.ID,
				PartnerID:  partner.ID,
				ExternalID: requestInfo.ExternalUserID,
			}
			if err := s.PartnerStore.CreatePartnerRelation(ctx, relationCmd); err != nil {
				return nil, err
			}
		}
	}
	if user == nil {
		panic("unexpected")
	}

	userTokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:          user.ID,
			AuthPartnerID:   partner.ID,
			SToken:          false,
			AccountIDs:      nil,
			STokenExpiresAt: nil,
			CAS:             0,
			Extra: map[string]string{
				"request_login": claim.Extra["request_login"],
			},
		},
	}
	if err := s.TokenStore.GenerateToken(ctx, userTokenCmd); err != nil {
		return nil, err
	}
	availableAccounts, err := s.getAvailableAccounts(ctx, user.ID, requestInfo)
	if err != nil {
		return nil, err
	}
	if requestInfo.ShopID != 0 && len(availableAccounts) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Bạn đã từng liên kết với đối tác này, nhưng tài khoản cũ không còn hiệu lực (mã tài khoản: &v)", requestInfo.ShopID).WithMeta("reason", "shop_id not found")
	}

	if requestInfo.ExternalUserID != "" && isExtraTokenInvitation(requestInfo.ExtraToken) {
		// Trường hợp invite user đối tác whitelabel
		// Đối tác chỉ cung cấp ExternalUserID và ExtraToken
		resp, err := s.generateNewSession(ctx, user, partner, nil, requestInfo)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	var shop *identitymodel.Shop
	var isGrantAccess bool
	for _, acc := range availableAccounts {
		if acc.ExternalId == requestInfo.ExternalShopID {
			query := &identitymodelx.GetShopQuery{
				ShopID: acc.Id,
			}
			if err := s.ShopStore.GetShop(ctx, query); err != nil {
				return nil, err
			}
			shop = query.Result
			isGrantAccess = true
		}
	}

	if shop == nil {
		// --- register shop --- //
		cmd := &identitymodelx.CreateShopCommand{
			Name:    requestInfo.ShopName,
			OwnerID: user.ID,
			Phone:   requestInfo.ShopOwnerPhone,
			Email:   requestInfo.ShopOwnerEmail,
		}
		if err := s.AccountStore.CreateShop(ctx, cmd); err != nil {
			return nil, err
		}
		shop = cmd.Result.Shop
	}
	if !isGrantAccess {
		// --- Grant access --- //
		cmd := &identitymodelx.CreatePartnerRelationCommand{
			PartnerID:  partner.ID,
			AccountID:  shop.ID,
			ExternalID: requestInfo.ExternalShopID,
		}
		if err := s.PartnerStore.CreatePartnerRelation(ctx, cmd); err != nil {
			return nil, err
		}
	}

	resp, err := s.generateNewSession(ctx, user, partner, shop, requestInfo)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *IntegrationService) validateWhiteLabel(ctx context.Context) (*whitelabel.WL, error) {
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return nil, cm.Errorf(cm.NotFound, nil, "Không tìm thấy thông tin whitelabel của partner")
	}
	return wlPartner, nil
}

func isExtraTokenInvitation(extraToken string) bool {
	if strings.HasPrefix(extraToken, auth.UsageInviteUser+":") {
		return true
	}
	return false
}
