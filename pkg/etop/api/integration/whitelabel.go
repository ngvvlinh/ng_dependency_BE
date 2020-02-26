package integration

import (
	"context"

	"etop.vn/api/top/int/integration"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/bus"
	apipartner "etop.vn/backend/pkg/etop/apix/partner"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/common/jsonx"
)

func (s *IntegrationService) LoginUsingTokenWL(ctx context.Context, r *LoginUsingTokenWLEndpoint) (_err error) {
	if _, err := s.validateWhiteLabel(ctx); err != nil {
		return err
	}
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(r.Context.Extra["request_login"]), &requestInfo); err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	partner := r.CtxPartner
	// wl: định danh bằng email
	userQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: requestInfo.ShopOwnerEmail,
	}
	err := bus.Dispatch(ctx, userQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// continue
	case cm.NotFound:
		// --- register user ---
		userQuery.Result.User, err = s.registerUser(ctx, false, partner.ID, requestInfo.ShopOwnerName, requestInfo.ShopOwnerEmail, requestInfo.ShopOwnerPhone, true, true, false, false)
		if err != nil {
			return err
		}
	default:
		return err
	}
	user := userQuery.Result.User

	userTokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:          user.ID,
			AuthPartnerID:   partner.ID,
			SToken:          false,
			AccountIDs:      nil,
			STokenExpiresAt: nil,
			CAS:             0,
			Extra: map[string]string{
				"request_login": r.Context.Extra["request_login"],
			},
		},
	}
	if err := bus.Dispatch(ctx, userTokenCmd); err != nil {
		return err
	}
	availableAccounts, err := getAvailableAccounts(ctx, user.ID, requestInfo)
	if err != nil {
		return err
	}
	if requestInfo.ShopID != 0 && len(availableAccounts) == 0 {
		return cm.Errorf(cm.NotFound, nil, "Bạn đã từng liên kết với đối tác này, nhưng tài khoản cũ không còn hiệu lực (mã tài khoản: &v)", requestInfo.ShopID).WithMeta("reason", "shop_id not found")
	}

	var shop *identitymodel.Shop
	var isGrantAccess bool
	for _, acc := range availableAccounts {
		if acc.ExternalId == requestInfo.ExternalShopID {
			query := &identitymodelx.GetShopQuery{
				ShopID: acc.Id,
			}
			if err := bus.Dispatch(ctx, query); err != nil {
				return err
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
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return nil
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
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}

	resp, err := s.generateNewSession(ctx, user, partner, shop, requestInfo)
	if err != nil {
		return err
	}
	r.Result = resp
	return nil
}

func (s *IntegrationService) validateWhiteLabel(ctx context.Context) (*whitelabel.WL, error) {
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return nil, cm.Errorf(cm.NotFound, nil, "Không tìm thấy thông tin whitelabel của partner")
	}
	return wlPartner, nil
}

func (s *IntegrationService) actionLoginUsingTokenWL(ctx context.Context, partner *identitymodel.Partner, info apipartner.PartnerShopToken) (*integration.LoginResponse, error) {
	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			Token:         "",
			AccountID:     0,
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": jsonx.MustMarshalToString(info),
			},
		},
		TTL: 2 * 60 * 60,
	}

	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	actions := []*integration.Action{
		{
			Name: "login_using_token_wl",
		},
	}
	resp := generateShopLoginResponse(
		tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn,
		nil, partner, nil, actions, info.RedirectURL,
	)
	return resp, nil
}
