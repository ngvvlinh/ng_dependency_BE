package integration

import (
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/int/integration"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/authorize_shop_config"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/user_source"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/apifw/whitelabel/templatemessages"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	apipartner "o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/authorize/authkey"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/sms"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

type IntegrationService struct {
	session.Session

	AuthStore   auth.Generator
	TokenStore  tokens.TokenStore
	SMSClient   *sms.Client
	EmailClient *email.Client

	UserStore        sqlstore.UserStoreInterface
	AccountStore     sqlstore.AccountStoreInterface
	AccountUserStore sqlstore.AccountUserStoreInterface
	PartnerStore     sqlstore.PartnerStoreInterface
	ShopStore        sqlstore.ShopStoreInterface
	IdentityAggr     identity.CommandBus
}

func (s *IntegrationService) Clone() integration.IntegrationService {
	res := *s
	return &res
}

func (s *IntegrationService) Init(ctx context.Context, q *integration.InitRequest) (*integration.LoginResponse, error) {
	authToken := q.AuthToken
	if authToken == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing token")
	}

	var partnerID dot.ID
	var relationError error
	var relationQuery *identitymodelx.GetPartnerRelationQuery
	var requestInfo apipartner.PartnerShopToken

	switch {
	case strings.HasPrefix(q.AuthToken, "shop"):
		_, ok := authkey.ValidateAuthKeyWithType(authkey.TypePartnerShopKey, q.AuthToken)
		if !ok {
			return nil, cm.Errorf(cm.Unauthenticated, nil, "M?? x??c th???c kh??ng h???p l???")
		}
		relationQuery = &identitymodelx.GetPartnerRelationQuery{
			AuthKey: q.AuthToken,
		}
		relationError = s.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
		if relationError != nil {
			return nil, cm.MapError(relationError).
				Map(cm.NotFound, cm.PermissionDenied, "").
				Throw()
		}
		partnerID = relationQuery.Result.PartnerID
		shop := relationQuery.Result.Shop
		user := relationQuery.Result.User
		requestInfo = apipartner.PartnerShopToken{
			PartnerID:         partnerID,
			ShopID:            shop.ID,
			ShopName:          shop.Name,
			ShopOwnerEmail:    user.Email,
			ShopOwnerPhone:    user.Phone,
			AuthMode:          apipartner.AuthModeManual,
			RetainCurrentInfo: true,
		}

	case strings.HasPrefix(q.AuthToken, "request:"):
		if _, err := s.AuthStore.Validate(auth.UsagePartnerIntegration, authToken, &requestInfo); err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.Unauthenticated, "M?? x??c th???c kh??ng h???p l???").
				Throw()
		}
		partnerID = requestInfo.PartnerID
		relationQuery = &identitymodelx.GetPartnerRelationQuery{
			PartnerID:         requestInfo.PartnerID,
			AccountID:         requestInfo.ShopID,
			ExternalAccountID: requestInfo.ExternalShopID,
		}
		relationError = s.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
		// the error will be handled later

	default:
		return nil, cm.Errorf(cm.Unauthenticated, nil, "M?? x??c th???c kh??ng h???p l???")
	}

	partner, err := s.validatePartner(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	if requestInfo.AuthType == apipartner.AuthTypeUserKey {
		// Tr?????ng h???p ???? x??c ?????nh c???n tr??? v??? user_key
		// C??c tr?????ng h???p:
		// - Invitation
		// - Auth with external_user_id (tr?????ng h???p ????ng nh???p v??o account do invitation
		relationQuery := &identitymodelx.GetPartnerRelationQuery{
			PartnerID:      partner.ID,
			ExternalUserID: requestInfo.ExternalUserID,
		}
		err := s.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
		if cm.ErrorCode(err) == cm.OK {
			// User ???? c?? s???n t??i kho???n
			user := relationQuery.Result.User
			if user.Status != status3.P {
				return nil, cm.Errorf(cm.AccountClosed, nil, "")
			}
			resp, err := s.generateNewSession(ctx, user, partner, nil, requestInfo)
			return resp, err
		}
	}

	if requestInfo.ShopID == 0 {
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		return resp, err
	}

	// TODO: refactor with partner.AuthorizeShop
	switch cm.ErrorCode(relationError) {
	case cm.OK:
		rel := relationQuery.Result.PartnerRelation
		shop := relationQuery.Result.Shop
		user := relationQuery.Result.User
		if rel.Status == status3.P && rel.DeletedAt.IsZero() &&
			shop.Status == status3.P && shop.DeletedAt.IsZero() &&
			user.Status == status3.P {
			// everything looks good
			resp, err := s.generateNewSession(ctx, user, partner, shop, requestInfo)
			return resp, err
		}
		if shop.Status != status3.P || !shop.DeletedAt.IsZero() ||
			user.Status != status3.P {
			return nil, cm.Errorf(cm.AccountClosed, nil, "")
		}
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		return resp, err

	case cm.NotFound:
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		return resp, err

	default:
		return nil, cm.Errorf(cm.Internal, err, "")
	}
}

func (s *IntegrationService) validatePartner(ctx context.Context, partnerID dot.ID) (*identitymodel.Partner, error) {
	partnerQuery := &identitymodelx.GetPartner{PartnerID: partnerID}
	if err := s.PartnerStore.GetPartner(ctx, partnerQuery); err != nil {
		return nil, cm.MapError(err).Map(cm.NotFound, cm.PermissionDenied, "M?? x??c th???c kh??ng h???p l???").Throw()
	}
	partner := partnerQuery.Result.Partner
	if partner.Status != status3.P || !partner.DeletedAt.IsZero() {
		return nil, cm.Errorf(cm.Unavailable, nil, "T??i kho???n ?????i t??c ???? kho?? ho???c kh??ng c??n ???????c s??? d???ng")
	}
	return partner, nil
}

func (s *IntegrationService) actionRequestLogin(ctx context.Context, partner *identitymodel.Partner, info apipartner.PartnerShopToken) (*integration.LoginResponse, error) {
	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			Token:         "",
			AccountID:     0, // shop has not login yet
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": jsonx.MustMarshalToString(info),
			},
		},
		TTL: 2 * 60 * 60,
	}
	if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}

	meta := map[string]string{}
	action := &integration.Action{
		Name:  "request_login",
		Label: fmt.Sprintf(`K???t n???i v???i %v`, partner.PublicName),
		Meta:  meta,
	}
	if info.ShopOwnerEmail != "" {
		meta["email"] = info.ShopOwnerEmail
	}
	if info.ShopOwnerEmail != "" {
		meta["phone"] = info.ShopOwnerPhone
	}
	if info.ShopName != "" {
		meta["name"] = info.ShopName
	}
	if info.RetainCurrentInfo {
		meta["retain_info"] = "name,email,phone"
	} else {
		meta["retain_info"] = ""
	}

	action.Msg = fmt.Sprintf("Nh???p s??? ??i???n tho???i ho???c email ????? k???t n???i v???i %v", partner.PublicName)
	resp := &integration.LoginResponse{
		AccessToken: tokenCmd.Result.TokenStr,
		ExpiresIn:   tokenCmd.Result.ExpiresIn,
		Actions:     []*integration.Action{action},
		AuthPartner: convertpb.PbPublicAccountInfo(partner),
		RedirectUrl: info.RedirectURL,
	}
	return resp, nil
}

func (s *IntegrationService) generateNewSession(ctx context.Context, user *identitymodel.User, partner *identitymodel.Partner, shop *identitymodel.Shop, info apipartner.PartnerShopToken) (*integration.LoginResponse, error) {
	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			Token:         "",
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": jsonx.MustMarshalToString(info),
			},
		},
	}
	if shop != nil {
		tokenCmd.ClaimInfo.AccountID = shop.ID
	}
	if user != nil {
		tokenCmd.ClaimInfo.UserID = user.ID
	}
	if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	actions := getActionsFromConfig(info.Config)
	if len(actions) == 0 {
		actions = append(actions, &integration.Action{
			Name: "using",
		})
	}
	resp := generateShopLoginResponse(
		tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn,
		user, partner, shop, actions, info.RedirectURL,
	)
	return resp, nil
}

func (s *IntegrationService) RequestLogin(ctx context.Context, r *integration.RequestLoginRequest) (*integration.RequestLoginResponse, error) {
	key := fmt.Sprintf("RequestLogin %v", r.Login)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 15*time.Second, "g???i m?? ????ng nh???p",
		func() (interface{}, error) { return s.requestLogin(ctx, r) })

	if err != nil {
		return nil, err
	}
	result := res.(*integration.RequestLoginResponse)
	return result, err
}

// - Verify the token
// - Check whether the user exists in our database
// - Generate verification code and send to email/phone
// - Response action login_using_token
func (s *IntegrationService) requestLogin(ctx context.Context, r *integration.RequestLoginRequest) (*integration.RequestLoginResponse, error) {
	partner := s.SS.CtxPartner()
	if partner == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}

	claim := s.SS.Claim()
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Y??u c???u ????ng nh???p kh??ng c??n hi???u l???c")
	}

	emailNorm, phoneNorm, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email ho???c s??? ??i???n tho???i kh??ng h???p l???")
	}

	userQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: r.Login,
	}

	err := s.UserStore.GetUserByLogin(ctx, userQuery)
	var exists bool
	switch cm.ErrorCode(err) {
	case cm.NoError:
		exists = true
	case cm.NotFound:
		exists = false
	default:
		return nil, err
	}

	_ = exists
	user := userQuery.Result.User

	extra := map[string]string{
		"requested_email": emailNorm,
		"requested_phone": phoneNorm,
	}
	_, generatedCode, _, err := generateTokenWithVerificationCode(s.AuthStore, partner.ID, cm.Coalesce(emailNorm, phoneNorm), extra)
	if err != nil {
		return nil, err
	}

	if err := s.TokenStore.UpdateSession(ctx, claim.Token, extra); err != nil {
		return nil, err
	}

	var msg, notice string
	notice = "B???ng vi???c nh???p m?? x??c nh???n, b???n ?????ng ?? c???p quy???n cho ?????i t??c t???o ????n h??ng v???i t?? c??ch c???a b???n."
	switch {
	case emailNorm != "":
		hello := "Xin ch??o"
		if user != nil {
			hello = "G???i " + user.FullName
		}
		var extraMessage template.HTML
		if cmenv.NotProd() {
			extraMessage = "<br><br><i>????y l?? email ???????c g???i t??? h??? th???ng c???a ?????i t??c th??ng qua eTop.vn nh???m m???c ????ch th??? nghi???m. N???u b???n cho r???ng ????y l?? s??? nh???m l???n, xin vui l??ng th??ng b??o cho ch??ng t??i.<i>"
		} else {
			extraMessage = "<br><br><i>????y l?? email ???????c g???i t??? h??? th???ng c???a ?????i t??c th??ng qua eTop.vn. N???u b???n cho r???ng ????y l?? s??? nh???m l???n, xin vui l??ng th??ng b??o cho ch??ng t??i.<i>"
		}

		var b strings.Builder
		if err := templatemessages.RequestLoginEmailTpl.Execute(&b, map[string]interface{}{
			"Code":              generatedCode,
			"Hello":             hello,
			"Email":             emailNorm,
			"PartnerPublicName": partner.PublicName,
			"PartnerWebsite":    validate.DomainFromURL(partner.WebsiteURL),
			"Notice":            notice,
			"Extra":             extraMessage,
			"WlName":            wl.X(ctx).Name,
		}); err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Kh??ng th??? g???i m?? ????ng nh???p").WithMeta("reason", "can not generate email content")
		}
		address := emailNorm
		cmd := &email.SendEmailCommand{
			FromName:    wl.X(ctx).CompanyName + " (no-reply)",
			ToAddresses: []string{address},
			Subject:     fmt.Sprintf("????ng nh???p v??o eTop.vn th??ng qua h??? th???ng %v", partner.PublicName),
			Content:     b.String(),
		}
		if err := s.EmailClient.SendMail(ctx, cmd); err != nil {
			return nil, err
		}
		msg = fmt.Sprintf("???? g???i email k??m m?? x??c nh???n ?????n ?????a ch??? %v. Vui l??ng ki???m tra email (k??? c??? trong h???p th?? spam). N???u c???n th??m th??ng tin, vui l??ng li??n h??? %v.", emailNorm, wl.X(ctx).CSEmail)

	case phoneNorm != "":
		tpl := wl.X(ctx).Templates.RequestLoginSmsTpl
		var b strings.Builder
		if err := tpl.Execute(&b, map[string]interface{}{
			"Code":           generatedCode,
			"PartnerWebsite": validate.DomainFromURL(partner.WebsiteURL),
			"Notice":         "",
		}); err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Kh??ng th??? g???i m?? ????ng nh???p").WithMeta("reason", "can not generate sms content")
		}
		phone := phoneNorm
		cmd := &sms.SendSMSCommand{
			Phone:   phone,
			Content: b.String(),
		}
		if err := s.SMSClient.SendSMS(ctx, cmd); err != nil {
			return nil, err
		}
		msg = fmt.Sprintf("???? g???i tin nh???n k??m m?? x??c nh???n ?????n s??? ??i???n tho???i %v. Vui l??ng ki???m tra tin nh???n. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", phoneNorm, wl.X(ctx).CSEmail)

	default:
		panic("unexpected")
	}

	result := &integration.RequestLoginResponse{
		Code: "ok",
		Msg:  msg,
		Actions: []*integration.Action{
			{
				Name:  "login_using_token",
				Label: "Nh???p m?? x??c nh???n",
				Msg:   notice,
			},
		},
	}
	return result, nil
}

func GetRequestLoginToken(authStore auth.Generator, partnerID dot.ID, login string) (*auth.Token, string, map[string]string) {
	tokStr := fmt.Sprintf("%v-%v", partnerID, login)
	var v map[string]string
	var code string

	tok, err := authStore.Validate(auth.UsageRequestLogin, tokStr, &v)
	if err == nil && v != nil && len(v["code"]) == 6 {
		code = v["code"]
		return tok, code, v
	}

	_ = authStore.Revoke(auth.UsageRequestLogin, tokStr)
	return nil, "", nil
}

func generateTokenWithVerificationCode(authStore auth.Generator, partnerID dot.ID, login string, extra map[string]string) (*auth.Token, string, map[string]string, error) {
	tokStr := fmt.Sprintf("%v-%v", partnerID, login)
	tok, code, v := GetRequestLoginToken(authStore, partnerID, login)
	if code != "" {
		return tok, code, v, nil
	}

	code, err := gencode.Random6Digits()
	if err != nil {
		return nil, "", nil, cm.Errorf(cm.Internal, nil, "")
	}
	v = map[string]string{
		"code":  code,
		"tries": "",
	}
	for key, val := range extra {
		v[key] = val
	}
	tok = &auth.Token{
		TokenStr: tokStr,
		Usage:    auth.UsageRequestLogin,
		Value:    v,
	}
	if _, err := authStore.GenerateWithValue(tok, 2*60*60); err != nil {
		return nil, "", nil, cm.Errorf(cm.Internal, err, "Kh??ng th??? g???i m?? ????ng nh???p").WithMeta("reason", "can not generate token")
	}
	return tok, code, v, nil
}

func (s *IntegrationService) LoginUsingToken(ctx context.Context, r *integration.LoginUsingTokenRequest) (_ *integration.LoginResponse, _err error) {
	if r.Login == "" || r.VerificationCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thi???u th??ng tin")
	}
	r.VerificationCode = strings.Replace(r.VerificationCode, " ", "", -1)
	claim := s.SS.Claim()

	// verify shop_id and external_shop_id in the request info:
	//
	//     - shop_id == 0 && external_shop_id == "": New shop
	//     - shop_id != 0 && external_shop_id == "": Must use existing shop
	//     - shop_id == 0 && external_shop_id != "": New shop and update existing relation
	//     - shop_id != 0 && external_shop_id != "": Must validate that they match
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Y??u c???u ????ng nh???p kh??ng c??n hi???u l???c")
	}

	emailNorm, phoneNorm, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email ho???c s??? ??i???n tho???i kh??ng h???p l???")
	}

	partner := s.SS.CtxPartner()
	key := fmt.Sprintf("%v-%v", partner.ID, cm.Coalesce(emailNorm, phoneNorm))

	var v map[string]string
	tok, err := s.AuthStore.Validate(auth.UsageRequestLogin, key, &v)
	if err != nil {
		return nil, cm.Errorf(cm.Unauthenticated, nil, "M?? x??c nh???n kh??ng h???p l???")
	}
	if v == nil || v["code"] != r.VerificationCode {
		return nil, cm.Errorf(cm.Unauthenticated, nil, "M?? x??c nh???n kh??ng h???p l???")
	}
	// delete the token after 5 minutes if login successfully
	defer func() {
		if _err != nil {
			_ = s.AuthStore.SetTTL(tok, 5*60)
		}
	}()

	userQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: r.Login,
	}
	err = s.UserStore.GetUserByLogin(ctx, userQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// continue

	case cm.NotFound:
		tokenCmd := &tokens.GenerateTokenCommand{
			ClaimInfo: claims.ClaimInfo{
				AuthPartnerID: partner.ID,
				Extra: map[string]string{
					"action:register": "1",
					"verified_phone":  phoneNorm,
					"verified_email":  emailNorm,
					"request_login":   claim.Extra["request_login"],
				},
			},
		}
		if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
			return nil, err
		}
		meta := map[string]string{}
		if phoneNorm != "" {
			meta["retain_info"] = "phone"
			meta["phone"] = phoneNorm
			meta["email"] = requestInfo.ShopOwnerEmail
		} else {
			meta["retain_info"] = "email"
			meta["phone"] = requestInfo.ShopOwnerPhone
			meta["email"] = emailNorm
		}
		// user can create new account
		actions := []*integration.Action{
			{
				Name: "register",
				Meta: meta,
			},
		}
		result := &integration.LoginResponse{
			AccessToken:       tokenCmd.Result.TokenStr,
			ExpiresIn:         tokenCmd.Result.ExpiresIn,
			User:              nil,
			Account:           nil,
			Shop:              nil,
			AvailableAccounts: nil,
			AuthPartner:       convertpb.PbPublicAccountInfo(partner),
			Actions:           actions,
		}
		return result, nil

	default:
		return nil, err
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
				"request_login": claim.Extra["request_login"],
			},
		},
	}
	if err := s.TokenStore.GenerateToken(ctx, userTokenCmd); err != nil {
		return nil, err
	}

	// we map from all accounts to partner relations and generate tokens for each one
	availableAccounts, err := s.getAvailableAccounts(ctx, user.ID, requestInfo)
	if err != nil {
		return nil, err
	}
	if requestInfo.ShopID != 0 && len(availableAccounts) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "B???n ???? t???ng li??n k???t v???i ?????i t??c n??y, nh??ng t??i kho???n c?? kh??ng c??n hi???u l???c (m?? t??i kho???n %v). Li??n h??? v???i %v ????? ???????c h?????ng d???n.", requestInfo.ShopID, wl.X(ctx).CSEmail).
			WithMeta("reason", "shop_id not found")
	}
	if requestInfo.ExternalShopID != "" && len(availableAccounts) != 0 {
		// automatically select external_shop_id:
		// 1. if there is account with that, use it
		// 2. otherwise, clear all token and wait for the user to grant access
		var aa *integration.PartnerShopLoginAccount
		for _, acc := range availableAccounts {
			if acc.ExternalId == requestInfo.ExternalShopID {
				aa = acc
				break
			}
		}

		if aa != nil {
			// found, respond the only account that match
			availableAccounts = []*integration.PartnerShopLoginAccount{aa}

		} else {
			var externalIDs []string
			avails := make([]*integration.PartnerShopLoginAccount, 0, len(availableAccounts))
			// not found, only list the accounts with external_shop_id empty,
			// clear all tokens and wait for the user to grant access because we
			// don't know which shop will map to external_shop_id
			for _, acc := range availableAccounts {
				if acc.ExternalId == "" {
					avails = append(avails, acc)
					acc.AccessToken = ""
					acc.ExpiresIn = 0
				} else {
					externalIDs = append(externalIDs, acc.ExternalId)
				}
			}
			if len(avails) == 0 {
				return nil, cm.Errorf(cm.NotFound, nil, "B???n ???? t???ng li??n k???t v???i ?????i t??c n??y, h??y s??? d???ng ????ng t??i kho???n c?? c???a b???n tr??n trang web ?????i t??c (m?? t??i kho???n ?????i t??c \"%v\"). Li??n h??? v???i ?????i t??c v?? cung c???p m?? t??i kho???n ?????i t??c ??? tr??n ????? t??m l???i t??i kho???n c?? c???a b???n.", strings.Join(externalIDs, `", "`)).
					WithMeta("reason", "external_shop_id not found")
			}
		}
	}

	result := &integration.LoginResponse{
		AccessToken:       userTokenCmd.Result.TokenStr,
		ExpiresIn:         userTokenCmd.Result.ExpiresIn,
		User:              convertpb.PbPartnerUserInfo(userQuery.Result.User),
		Account:           nil,
		Shop:              nil,
		AvailableAccounts: availableAccounts,
		AuthPartner:       convertpb.PbPublicAccountInfo(partner),
		Actions:           getActionsFromConfig(requestInfo.Config),
		RedirectUrl:       requestInfo.RedirectURL,
	}
	return result, nil
}

func getActionsFromConfig(config string) (actions []*integration.Action) {
	if config == "" {
		return
	}
	_actions := strings.Split(config, ",")
	for _, a := range _actions {
		c, ok := authorize_shop_config.ParseAuthorizeShopConfig(a)
		if !ok {
			continue
		}
		switch c {
		case authorize_shop_config.Shipment:
			actions = append(actions, &integration.Action{
				Name: authorize_shop_config.Shipment.String(),
			})
		default:
			continue
		}
	}
	return
}

func (s *IntegrationService) Register(ctx context.Context, r *integration.RegisterRequest) (*integration.RegisterResponse, error) {
	partner := s.SS.CtxPartner()
	claim := s.SS.Claim()
	if claim.Extra == nil || claim.Extra["action:register"] == "" {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "C???n x??c nh???n email ho???c s??? ??i???n tho???i tr?????c khi ????ng k??")
	}

	if !r.AgreeTos {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "B???n c???n ?????ng ?? v???i ??i???u kho???n s??? d???ng d???ch v??? ????? ti???p t???c. N???u c???n th??m th??ng tin, vui l??ng li??n h??? %.", wl.X(ctx).CSEmail)
	}
	if !r.AgreeEmailInfo.Valid {
		return nil, cm.Error(cm.InvalidArgument, "Missing agree_email_info", nil)
	}

	phoneNorm, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "S??? ??i???n tho???i kh??ng h???p l???")
	}
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email kh??ng h???p l???")
	}

	verifiedPhone, verifiedEmail := claim.Extra["verified_phone"], claim.Extra["verified_email"]
	if verifiedPhone != "" && verifiedPhone != string(phoneNorm) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "C???n s??? d???ng s??? ??i???n tho???i ???? ???????c x??c nh???n khi ????ng k??")
	}
	if verifiedEmail != "" && verifiedEmail != string(emailNorm) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "C???n s??? d???ng ?????a ch??? email ???? ???????c x??c nh???n khi ????ng k??")
	}

	user, err := s.registerUser(ctx, true, partner.ID, r.FullName, string(emailNorm), string(phoneNorm), r.AgreeTos, r.AgreeEmailInfo.Bool, verifiedEmail != "", verifiedPhone != "")
	if err != nil {
		return nil, err
	}

	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:        user.ID,
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": claim.Extra["request_login"],
			},
		},
	}
	if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
		return nil, err
	}

	result := &integration.RegisterResponse{
		User:        convertpb.PbUser(user),
		AccessToken: tokenCmd.Result.TokenStr,
		ExpiresIn:   tokenCmd.Result.ExpiresIn,
	}
	return result, nil
}

func (s *IntegrationService) registerUser(ctx context.Context, sendConfirmInfo bool, partnerID dot.ID, fullName, userEmail, userPhone string, agreeTos, agreeEmailInfo bool, verifiedEmail, verifiedPhone bool) (*identitymodel.User, error) {
	partner, err := s.validatePartner(ctx, partnerID)
	if err != nil {
		return nil, err
	}
	generatedPassword := gencode.GenerateCode(gencode.Alphabet32, 8)
	// set default source: "partner"
	source := user_source.Partner

	cmd := &identitymodelx.CreateUserCommand{
		UserInner: identitymodel.UserInner{
			FullName:  fullName,
			ShortName: "",
			Phone:     userPhone,
			Email:     userEmail,
		},
		Password:       generatedPassword,
		Status:         status3.P,
		AgreeTOS:       agreeTos,
		AgreeEmailInfo: agreeEmailInfo,
		Source:         source,
	}
	if err := s.UserStore.CreateUser(ctx, cmd); err != nil {
		return nil, err
	}
	user := cmd.Result.User

	if verifiedEmail {
		verifyCmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			EmailVerifiedAt: time.Now(),
		}
		if err := s.UserStore.UpdateUserVerification(ctx, verifyCmd); err != nil {
			ll.Error("Can not update verification", l.Error(err))
		}

		if sendConfirmInfo {
			var b strings.Builder
			err := templatemessages.NewAccountViaPartnerEmailTpl.Execute(&b, map[string]interface{}{
				"FullName":          user.FullName,
				"PartnerPublicName": partner.PublicName,
				"PartnerWebsite":    validate.DomainFromURL(partner.WebsiteURL),
				"LoginLabel":        "Email",
				"Login":             userEmail,
				"Password":          generatedPassword,
			})
			if err == nil {
				emailCmd := &email.SendEmailCommand{
					FromName:    "eTop.vn (no-reply)",
					ToAddresses: []string{userEmail},
					Subject:     "M???t kh???u ????ng nh???p v??o t??i kho???n ??? eTop.vn",
					Content:     b.String(),
				}
				if err := s.EmailClient.SendMail(ctx, emailCmd); err != nil {
					ll.Error("Can not send email", l.Error(err))
				}
			}
		}
	}
	if verifiedPhone {
		verifyCmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			PhoneVerifiedAt: time.Now(),
		}
		if err := s.UserStore.UpdateUserVerification(ctx, verifyCmd); err != nil {
			ll.Error("Can not update verification", l.Error(err))
		}

		if sendConfirmInfo {
			tpl := wl.X(ctx).Templates.NewAccountViaPartnerSmsTpl
			var b strings.Builder
			err := tpl.Execute(&b, map[string]interface{}{
				"Password": generatedPassword,
			})
			if err == nil {
				smsCmd := &sms.SendSMSCommand{
					Phone:   userPhone,
					Content: b.String(),
				}
				if err := s.SMSClient.SendSMS(ctx, smsCmd); err != nil {
					ll.Error("Can not send sms", l.Error(err))
				}
			}
		}
	}
	return user, nil
}

func (s *IntegrationService) GrantAccess(ctx context.Context, r *integration.GrantAccessRequest) (*integration.GrantAccessResponse, error) {
	claim := s.SS.Claim()
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Y??u c???u ????ng nh???p kh??ng c??n hi???u l???c")
	}

	partner := s.SS.CtxPartner()
	user := s.SS.User()
	if r.ShopId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}

	if requestInfo.ShopID != 0 && r.ShopId != requestInfo.ShopID {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "B???n c???n s??? d???ng t??i kho???n ???? t???ng li??n k???t v???i ?????i t??c n??y").WithMeta("reason", "shop_id does not match")
	}

	shopQuery := &identitymodelx.GetShopQuery{
		ShopID: r.ShopId,
	}
	if err := s.ShopStore.GetShop(ctx, shopQuery); err != nil {
		return nil, err
	}
	shop := shopQuery.Result
	if shop.OwnerID != user.ID {
		return nil, cm.Errorf(cm.NotFound, nil, "")
	}

	relQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID: partner.ID,
		AccountID: shop.ID,
	}
	err := s.PartnerStore.GetPartnerRelationQuery(ctx, relQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// verify the external_shop_id or automatically fill it
		rel := relQuery.Result
		if rel.ExternalSubjectID != "" && requestInfo.ExternalShopID != "" &&
			rel.ExternalSubjectID != requestInfo.ExternalShopID {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "B???n c???n s??? d???ng t??i kho???n ???? t???ng li??n k???t v???i ?????i t??c n??y").WithMeta("reason", "external_shop_id does not match")
		}
		if rel.ExternalSubjectID == "" && requestInfo.ExternalShopID != "" {
			cmd := &identitymodelx.UpdatePartnerRelationCommand{
				PartnerID:  partner.ID,
				AccountID:  r.ShopId,
				ExternalID: requestInfo.ExternalShopID,
			}
			if err := s.PartnerStore.UpdatePartnerRelationCommand(ctx, cmd); err != nil {
				return nil, cm.Errorf(cm.Internal, err, "Kh??ng th??? c???p nh???t th??ng tin t??i kho???n. Vui l??ng li??n h??? %v ????? ???????c h??? tr???.", wl.X(ctx).CSEmail).
					WithMeta("reason", "can not update external_shop_id")
			}
		}

	case cm.NotFound:
		cmd := &identitymodelx.CreatePartnerRelationCommand{
			PartnerID:  partner.ID,
			AccountID:  r.ShopId,
			ExternalID: requestInfo.ExternalShopID,
		}
		if err := s.PartnerStore.CreatePartnerRelation(ctx, cmd); err != nil {
			return nil, err
		}

	default:
		return nil, err
	}

	// set shop IsPriorMoneyTransaction = true
	updateShop := &identity.UpdateShopInfoCommand{
		ShopID:                  r.ShopId,
		IsPriorMoneyTransaction: dot.Bool(true),
	}
	_ = s.IdentityAggr.Dispatch(ctx, updateShop)

	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:        user.ID,
			AccountID:     shop.ID,
			AuthPartnerID: partner.ID,
		},
	}
	if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
		return nil, err
	}

	result := &integration.GrantAccessResponse{
		AccessToken: tokenCmd.Result.TokenStr,
		ExpiresIn:   tokenCmd.Result.ExpiresIn,
	}
	return result, nil
}

func (s *IntegrationService) SessionInfo(ctx context.Context, q *pbcm.Empty) (*integration.LoginResponse, error) {
	var shop *identitymodel.Shop
	claim := s.SS.Claim()
	if claim.AccountID != 0 {
		query := &identitymodelx.GetShopQuery{
			ShopID: claim.AccountID,
		}
		err := s.ShopStore.GetShop(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.OK, cm.NotFound:
			// continue
		default:
			return nil, err
		}
		shop = query.Result
	}
	var actions []*integration.Action
	redirectURL := ""
	if claim.Extra != nil && claim.Extra["request_login"] != "" {
		var requestInfo apipartner.PartnerShopToken
		_ = jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo)
		actions = getActionsFromConfig(requestInfo.Config)
		redirectURL = requestInfo.RedirectURL
	}

	result := generateShopLoginResponse(
		claim.Token, tokens.DefaultAccessTokenTTL,
		nil, s.SS.CtxPartner(), shop, actions, redirectURL,
	)
	return result, nil
}

func generateShopLoginResponse(accessToken string, expiresIn int, user *identitymodel.User, partner *identitymodel.Partner, shop *identitymodel.Shop, actions []*integration.Action, redirectURL string) *integration.LoginResponse {
	resp := &integration.LoginResponse{
		AccessToken:       accessToken,
		ExpiresIn:         expiresIn,
		Account:           nil,
		AvailableAccounts: nil,
		User:              convertpb.PbPartnerUserInfo(user),
		Shop:              nil,
		AuthPartner:       convertpb.PbPublicAccountInfo(partner),
		Actions:           actions,
		RedirectUrl:       redirectURL,
	}

	if shop != nil {
		account := &integration.PartnerShopLoginAccount{
			Id:          shop.ID,
			Name:        shop.Name,
			Type:        convertpb.PbAccountType(account_type.Shop),
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
			ImageUrl:    shop.ImageURL,
		}
		resp.Account = account
		resp.AvailableAccounts = []*integration.PartnerShopLoginAccount{account}
		resp.Shop = convertpb.PbPartnerShopInfo(shop)
	}
	return resp
}

func (s *IntegrationService) getAvailableAccounts(ctx context.Context, userID dot.ID, requestInfo apipartner.PartnerShopToken) ([]*integration.PartnerShopLoginAccount, error) {
	relationQuery := &identitymodelx.GetPartnerRelationsQuery{
		PartnerID: requestInfo.PartnerID,
		OwnerID:   userID,
	}
	if err := s.PartnerStore.GetPartnerRelations(ctx, relationQuery); err != nil {
		return nil, err
	}

	accQuery := &identitymodelx.GetAllAccountRolesQuery{
		UserID: userID,
		Type:   account_type.Shop.Wrap(),
	}
	if err := s.AccountUserStore.GetAllAccountRoles(ctx, accQuery); err != nil {
		return nil, err
	}
	// we map from all accounts to partner relations and generate tokens for each one
	availableAccounts := make([]*integration.PartnerShopLoginAccount, 0, len(accQuery.Result))
	for _, acc := range accQuery.Result {
		// if there is shop_id, only list the account with that id
		if requestInfo.ShopID != 0 && acc.Account.ID != requestInfo.ShopID {
			continue
		}

		availAcc := &integration.PartnerShopLoginAccount{
			Id:       acc.Account.ID,
			Name:     acc.Account.Name,
			Type:     convertpb.PbAccountType(acc.Account.Type),
			ImageUrl: acc.Account.ImageURL,
		}
		for _, rel := range relationQuery.Result.Relations {
			if rel.SubjectType == identity.SubjectTypeAccount && rel.SubjectID == acc.Account.ID {
				// the partner has permission to access this account
				tokenCmd := &tokens.GenerateTokenCommand{
					ClaimInfo: claims.ClaimInfo{
						UserID:        userID,
						AccountID:     acc.Account.ID,
						AuthPartnerID: requestInfo.PartnerID,
					},
				}
				if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
					return nil, err
				}
				availAcc.ExternalId = rel.ExternalSubjectID
				availAcc.AccessToken = tokenCmd.Result.TokenStr
				availAcc.ExpiresIn = tokenCmd.Result.ExpiresIn
				break
			}
		}
		availableAccounts = append(availableAccounts, availAcc)
	}
	return availableAccounts, nil
}
