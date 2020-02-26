package integration

import (
	"context"
	"fmt"
	"strings"

	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	apipartner "etop.vn/backend/pkg/etop/apix/partner"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/common/jsonx"
)

func (s *IntegrationService) LoginUsingTokenWL(ctx context.Context, r *LoginUsingTokenWLEndpoint) (_err error) {
	if r.Login == "" || r.VerificationCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin")
	}
	r.VerificationCode = strings.Replace(r.VerificationCode, " ", "", -1)
	if _, err := s.validateWhiteLabel(ctx); err != nil {
		return err
	}
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(r.Context.Extra["request_login"]), &requestInfo); err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	verifiedEmail, verifiedPhone, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "Email hoặc số điện thoại không hợp lệ")
	}

	partner := r.CtxPartner
	key := fmt.Sprintf("%v-%v", partner.ID, cm.Coalesce(verifiedEmail, verifiedPhone))

	var v map[string]string
	tok, err := authStore.Validate(auth.UsageRequestLogin, key, &v)
	if err != nil {
		return cm.Errorf(cm.Unauthenticated, nil, "Mã xác nhận không hợp lệ")
	}
	if v == nil || v["code"] != r.VerificationCode {
		return cm.Errorf(cm.Unauthenticated, nil, "Mã xác thực không hợp lệ")
	}
	// delete the token after 5 minutes if login successfully
	defer func() {
		if _err != nil {
			_ = authStore.SetTTL(tok, 5*60)
		}
	}()
	userQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: r.Login,
	}
	err = bus.Dispatch(ctx, userQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// continue
	case cm.NotFound:
		// --- register user ---
		// trust thông tin từ đối tác gửi qua
		// xem như email, phone đều đã xác thực
		user, err := s.registerUser(ctx, false, partner.ID, requestInfo.ShopOwnerName, requestInfo.ShopOwnerEmail, requestInfo.ShopOwnerPhone, true, true, true, true)
		if err != nil {
			return err
		}
		userQuery.Result.User = user
		if requestInfo.ExternalUserID != "" {
			// create partner relation with user
			relationCmd := &identitymodelx.CreatePartnerRelationCommand{
				UserID:     user.ID,
				PartnerID:  partner.ID,
				ExternalID: requestInfo.ExternalUserID,
			}
			if err := bus.Dispatch(ctx, relationCmd); err != nil {
				return err
			}
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

	if requestInfo.ExternalShopID == "" {
		// Trường hợp invite user đối tác whitelabel
		// Đối tác chỉ cung cấp ExternalUserID và ExtraToken
		resp, err := s.generateNewSession(ctx, user, partner, nil, requestInfo)
		if err != nil {
			return err
		}
		r.Result = resp
		return nil
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
