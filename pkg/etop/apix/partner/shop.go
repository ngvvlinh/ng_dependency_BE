package partner

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	api "o.o/api/top/external/partner"
	etop "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/authorize_shop_config"
	"o.o/api/top/types/etc/status3"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/validate"
	apiconvertpb "o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
	"o.o/common/l"
)

type ShopService struct {
	session.Session

	PartnerStore sqlstore.PartnerStoreInterface
}

func (s *ShopService) Clone() api.ShopService { res := *s; return &res }

func (s *ShopService) CurrentShop(ctx context.Context, q *pbcm.Empty) (*etop.PublicAccountInfo, error) {
	if s.SS.Shop() == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}
	result := apiconvertpb.PbPublicAccountInfo(s.SS.Shop())
	return result, nil
}

func getAuthorizeShopConfig(configs []authorize_shop_config.AuthorizeShopConfig) string {
	var res []string
	for _, c := range configs {
		res = append(res, c.String())
	}
	return strings.Join(res, ",")
}

func (s *ShopService) AuthorizeShop(ctx context.Context, r *api.AuthorizeShopRequest) (*api.AuthorizeShopResponse, error) {
	partner := s.SS.Partner()
	if r.RedirectUrl != "" {
		if err := validateRedirectURL(partner.RedirectURLs, r.RedirectUrl, true); err != nil {
			return nil, err
		}
	}

	whiteLabelData, err := validateWhiteLabel(ctx, r)
	if err != nil {
		return nil, err
	}

	// always verify email and phone
	var emailNorm, phoneNorm string
	if r.Email != "" {
		email, ok := validate.NormalizeEmail(r.Email)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, `Giá trị email không hợp lệ`)
		}
		emailNorm = email.String()
	}
	if r.Phone != "" {
		phone, ok := validate.NormalizePhone(r.Phone)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, `Giá trị phone không hợp lệ`)
		}
		phoneNorm = phone.String()
	}
	_, _ = emailNorm, phoneNorm // temporary ignore error here

	// verify external code
	if r.ExternalShopID != "" && !validate.ExternalCode(r.ExternalShopID) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giá trị external_shop_id không hợp lệ")
	}
	if r.ExternalUserID != "" && !validate.ExternalCode(r.ExternalUserID) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giá trị external_user_id không hợp lệ")
	}

	// case 1: if the shop has linked to the partner
	if r.ShopId != 0 && r.ExternalShopID != "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Chỉ cần cung cấp shop_id hoặc external_shop_id")
	}

	if r.ExternalUserID != "" {
		if !wl.X(ctx).IsWhiteLabel() || whiteLabelData == nil {
			// only whitelabel partner can provide this field,
			// we fake the error message here
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_user_id is removed")
		}
		if r.ExtraToken != "" || r.ExternalShopID == "" {
			return s.handleWLAuthorizeShopByExternalUserID(ctx, r, s.SS.Partner().ID)
		}
	}

	if r.ShopId != 0 || r.ExternalShopID != "" {
		if r.ShopId != 0 && (r.Email != "" || r.Phone != "") {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Nếu cung cấp shop_id thì không cần kèm theo email hoặc phone. Nếu cung cấp email hoặc phone thì không cần kèm theo shop_id.")
		}

		relationQuery := &identitymodelx.GetPartnerRelationQuery{
			PartnerID:         partner.ID,
			AccountID:         r.ShopId,
			ExternalAccountID: r.ExternalShopID,
		}
		err := s.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
		switch {
		case err == nil:
			rel := relationQuery.Result.PartnerRelation
			shop := relationQuery.Result.Shop
			user := relationQuery.Result.User
			if r.ExternalUserID != "" {
				relationUserQuery := &identitymodelx.GetPartnerRelationQuery{
					PartnerID:      partner.ID,
					ExternalUserID: r.ExternalUserID,
				}
				if err2 := s.PartnerStore.GetPartnerRelationQuery(ctx, relationUserQuery); err2 == nil {
					userID := relationUserQuery.Result.SubjectID
					if userID != 0 && userID != user.ID {
						return nil, cm.Errorf(cm.FailedPrecondition, nil, "external_shop_id (id=%v) does not belong to external_user_id (id=%v).", r.ExternalShopID, r.ExternalUserID)
					}
				}
			}

			info := PartnerShopToken{
				PartnerID:      partner.ID,
				ShopID:         shop.ID,
				ShopName:       shop.Name,
				ShopOwnerEmail: user.Email,
				ShopOwnerPhone: user.Phone,
				ShopOwnerName:  user.FullName,
				ExternalShopID: r.ExternalShopID,
				ExternalUserID: r.ExternalUserID,

				// client must keep the current email/phone when calling
				// request_login
				RetainCurrentInfo: true,
				Config:            getAuthorizeShopConfig(r.Config),
				RedirectURL:       r.RedirectUrl,
				AuthType:          AuthTypeShopKey,
			}
			if whiteLabelData != nil {
				// Trường hợp whiteLabel
				// Đã có sẵn tài khoản sẽ redirect về trang whiteLabel
				info.RedirectURL = whiteLabelData.Config.RootURL
			}
			token, err := generateAuthToken(info)
			if err != nil {
				return nil, err
			}

			if rel.Status == status3.P && rel.DeletedAt.IsZero() &&
				shop.Status == status3.P && shop.DeletedAt.IsZero() &&
				user.Status == status3.P {
				// TODO: handle config: "shipment"
				result := &api.AuthorizeShopResponse{
					Code:      "ok",
					Msg:       msgShopKey,
					Type:      "shop_key",
					AuthToken: rel.AuthKey,
					ExpiresIn: -1,
					AuthUrl:   generateAuthURL(whiteLabelData, authURL, token.TokenStr),
				}
				return result, nil
			}
			if shop.Status != status3.P || !shop.DeletedAt.IsZero() ||
				user.Status != status3.P {
				return nil, cm.Errorf(cm.AccountClosed, nil, "")
			}
			if rel.Status != status3.P || !rel.DeletedAt.IsZero() {
				result := &api.AuthorizeShopResponse{
					Code:      "ok",
					Msg:       msgShopRequest,
					Type:      "shop_request",
					AuthToken: token.TokenStr,
					ExpiresIn: token.ExpiresIn,
				}
				if r.RedirectUrl != "" {
					result.AuthUrl = generateAuthURL(whiteLabelData, authURL, result.AuthToken)
				}
				return result, nil
			}

		case cm.ErrorCode(err) == cm.NotFound:
			if r.ShopId != 0 {
				return nil, cm.Errorf(cm.PermissionDenied, nil, "").
					WithMeta("reason", "Chỉ có thể sử dụng shop_id nếu shop đã từng đăng nhập qua hệ thống của đối tác")
			}
			return generateAuthTokenWithRequestLogin(ctx, r, s.SS.Partner().ID, 0)

		default:
			return nil, cm.Errorf(cm.Internal, err, "")
		}

		// prevent unexpected condition
		return nil, cm.Errorf(cm.Internal, nil, "").
			WithMeta("reason", "unexpected condition")
	}

	// case 2: if the user/shop has not linked to the partner
	return generateAuthTokenWithRequestLogin(ctx, r, s.SS.Partner().ID, 0)
}

func (s *ShopService) handleWLAuthorizeShopByExternalUserID(ctx context.Context, r *api.AuthorizeShopRequest, partnerID dot.ID) (*api.AuthorizeShopResponse, error) {
	if !wl.X(ctx).IsWhiteLabel() {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Only whitelabel partner can use this function")
	}

	wlPartner := wl.X(ctx)
	redirectURL := ""
	switch {
	case r.ExternalUserID != "" && r.ExtraToken != "":
		// TH invitation
		// check if invitation
		if err := checkExtraTokenInvitation(r.ExtraToken); err != nil {
			return nil, err
		}
		redirectURL = fmt.Sprintf("%v?t=%v", wlPartner.InviteUserURLByEmail, r.ExtraToken)
	case r.ExternalUserID != "" && r.ExternalShopID == "":
		// TH login bằng external_user_id
		redirectURL = wlPartner.Config.RootURL
	default:
		return nil, cm.Errorf(cm.Internal, nil, "").WithMeta("reason", "unexpected condition")
	}

	// try retrieving the partner relation
	relationQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID:      wlPartner.ID,
		ExternalUserID: r.ExternalUserID,
	}
	err := s.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// generate user_key and redirect to whiteLabelData.Config.RootURL
		rel := relationQuery.Result.PartnerRelation
		user := relationQuery.Result.User
		info := PartnerShopToken{
			PartnerID:      wlPartner.ID,
			ShopOwnerEmail: user.Email,
			ShopOwnerPhone: user.Phone,
			ShopOwnerName:  user.FullName,
			ExternalUserID: r.ExternalUserID,
			ExtraToken:     r.ExtraToken,
			// client must keep the current email/phone when calling
			// request_login
			RetainCurrentInfo: true,
			Config:            getAuthorizeShopConfig(r.Config),
			RedirectURL:       redirectURL,
			AuthType:          AuthTypeUserKey,
		}
		token, err := generateAuthToken(info)
		if err != nil {
			return nil, err
		}

		if user.Status != status3.P {
			return nil, cm.Errorf(cm.AccountClosed, nil, "")
		}
		// TODO: handle config: "shipment"
		result := &api.AuthorizeShopResponse{
			Code:      "ok",
			Msg:       msgUserKey,
			Type:      "user_key",
			AuthToken: rel.AuthKey,
			ExpiresIn: -1,
			AuthUrl:   generateAuthURL(wlPartner, authURL, token.TokenStr),
		}
		return result, nil

	case cm.NotFound:
		// chỉ cho phép pass bước này trong trường hợp invitation
		if _err := checkExtraTokenInvitation(r.ExtraToken); _err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Request không hợp lệ.")
		}

		// continue to case 2 (the user/shop has not linked to the partner)
		return generateAuthTokenWithRequestLogin(ctx, r, partnerID, 0)
	default:
		return nil, err
	}
}

func generateAuthTokenWithRequestLogin(ctx context.Context, r *api.AuthorizeShopRequest, partnerID, shopID dot.ID) (*api.AuthorizeShopResponse, error) {
	whiteLabelData, err := validateWhiteLabel(ctx, r)
	if err != nil {
		return nil, err
	}

	info := PartnerShopToken{
		PartnerID: partnerID,

		// leave this field empty because we don't want to expose our account
		// information to partner
		ShopID:            0,
		ShopName:          r.ShopName,
		ShopOwnerEmail:    r.Email,
		ShopOwnerPhone:    r.Phone,
		ShopOwnerName:     r.Name,
		ExternalShopID:    r.ExternalShopID,
		ExternalUserID:    r.ExternalUserID,
		ExtraToken:        r.ExtraToken,
		RetainCurrentInfo: false,
		RedirectURL:       r.RedirectUrl,
		Config:            getAuthorizeShopConfig(r.Config),
	}
	if shopID != 0 {
		info.ShopID = shopID
		info.RetainCurrentInfo = true
	}
	if whiteLabelData != nil && r.ExternalUserID != "" && r.ExtraToken != "" {
		if err := checkExtraTokenInvitation(r.ExtraToken); err != nil {
			return nil, err
		}
		info.RedirectURL = fmt.Sprintf("%v?t=%v", whiteLabelData.InviteUserURLByEmail, r.ExtraToken)
	}

	token, err := generateAuthToken(info)
	if err != nil {
		return nil, err
	}

	result := &api.AuthorizeShopResponse{
		Code:      "ok",
		Msg:       msgShopRequest,
		Type:      "shop_request",
		AuthToken: token.TokenStr,
		ExpiresIn: token.ExpiresIn,
		Meta: map[string]string{
			"recaptcha_token": "1",
		},
	}
	if r.RedirectUrl != "" {
		result.AuthUrl = generateAuthURL(whiteLabelData, authURL, result.AuthToken)
	}
	return result, nil
}

func generateAuthToken(info PartnerShopToken) (*auth.Token, error) {
	if info.PartnerID == 0 {
		return nil, cm.Errorf(cm.Internal, nil, "Missing PartnerID")
	}

	tokStr := "request:" + auth.RandomToken(auth.DefaultTokenLength)
	tok := &auth.Token{
		TokenStr: tokStr,
		Usage:    auth.UsagePartnerIntegration,
		UserID:   0,
		Value:    info,
	}
	_, err := authStore.GenerateWithValue(tok, ttlShopRequest)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	return tok, nil
}

func generateAuthURL(whitelabelData *whitelabel.WL, authURL string, token string) string {
	if whitelabelData != nil {
		authURL = whitelabelData.Config.AuthURL
	}
	u, err := url.Parse(authURL)
	if err != nil {
		ll.Panic("invalid auth_url", l.Error(err))
	}
	query := u.Query()
	query.Set("token", token)
	u.RawQuery = query.Encode()
	return u.String()
}

var reLoopback = regexp.MustCompile(`^127\.0\.0\.[0-9]{3}$`)

func validateRedirectURL(redirectURLs []string, redirectURL string, skipCheckIfNoURL bool) error {
	rURL, err := url.Parse(redirectURL)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url không hợp lệ")
	}

	// allow localhost for testing
	if rURL.Host == "localhost" || reLoopback.MatchString(rURL.Host) {
		return nil
	}

	if skipCheckIfNoURL && len(redirectURLs) == 0 {
		return nil
	}
	for _, registerURL := range redirectURLs {
		if redirectURL == registerURL {
			return nil
		}
	}
	return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url cần được đăng ký trước")
}

func validateWhiteLabel(ctx context.Context, q *api.AuthorizeShopRequest) (wlPartner *whitelabel.WL, err error) {
	wlPartner = wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return nil, nil
	}

	// validate data
	fields := []cmapi.Field{
		{
			Name:  "Số điện thoại (phone)",
			Value: q.Phone,
		}, {
			Name:  "Email",
			Value: q.Email,
		}, {
			Name:  "Tên (name)",
			Value: q.Name,
		}, {
			Name:  "redirect_url",
			Value: q.RedirectUrl,
		}, {
			Name:  "Mã người dùng (external_user_id)",
			Value: q.ExternalUserID,
		},
		// NOTE(vu): external_shop_id may be empty
	}

	if q.ExternalShopID != "" {
		validateFields := []cmapi.Field{
			{
				Name:  "Tên cửa hàng (shop_name)",
				Value: q.ShopName,
			},
		}
		fields = append(fields, validateFields...)
	}

	if err = cmapi.ValidateEmptyField(fields...); err != nil {
		return nil, err
	}
	return wlPartner, nil
}

func checkExtraTokenInvitation(extraToken string) error {
	if strings.HasPrefix(extraToken, auth.UsageInviteUser+":") {
		return nil
	}
	return cm.Errorf(cm.InvalidArgument, nil, "extra_token does not valid")
}
