package integration

import (
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"etop.vn/api/top/int/integration"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/user_source"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/idemp"
	cmservice "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/common/code/gencode"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	apipartner "etop.vn/backend/pkg/etop/apix/partner"
	"etop.vn/backend/pkg/etop/authorize/authkey"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/logic/usering"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup
var authStore auth.Generator

func init() {
	bus.AddHandlers("api",
		miscService.VersionInfo,
		integrationService.Init,
		integrationService.RequestLogin,
		integrationService.LoginUsingToken,
		integrationService.Register,
		integrationService.GrantAccess,
		integrationService.SessionInfo,
	)
}

func Init(sd cmservice.Shutdowner, rd redis.Store, s auth.Generator) {
	authStore = s
	idempgroup = idemp.NewRedisGroup(rd, api.PrefixIdempUser, 0)
	sd.Register(idempgroup.Shutdown)
}

type MiscService struct{}
type IntegrationService struct{}

var integrationService = &IntegrationService{}
var miscService = &MiscService{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Integration",
		Version: "0.1",
	}
	return nil
}

func (s *IntegrationService) Init(ctx context.Context, q *InitEndpoint) error {
	authToken := q.AuthToken
	if authToken == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing token")
	}

	var partnerID dot.ID
	var relationError error
	var relationQuery *identitymodelx.GetPartnerRelationQuery
	var requestInfo apipartner.PartnerShopToken

	switch {
	case strings.HasPrefix(q.AuthToken, "shop"):
		_, ok := authkey.ValidateAuthKeyWithType(authkey.TypePartnerShopKey, q.AuthToken)
		if !ok {
			return cm.Errorf(cm.Unauthenticated, nil, "Mã xác thực không hợp lệ")
		}
		relationQuery = &identitymodelx.GetPartnerRelationQuery{
			AuthKey: q.AuthToken,
		}
		relationError = bus.Dispatch(ctx, relationQuery)
		if relationError != nil {
			return cm.MapError(relationError).
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
		if _, err := authStore.Validate(auth.UsagePartnerIntegration, authToken, &requestInfo); err != nil {
			return cm.MapError(err).
				Map(cm.NotFound, cm.Unauthenticated, "Mã xác thực không hợp lệ").
				Throw()
		}
		partnerID = requestInfo.PartnerID
		relationQuery = &identitymodelx.GetPartnerRelationQuery{
			PartnerID: requestInfo.PartnerID,
			AccountID: requestInfo.ShopID,
		}
		relationError = bus.Dispatch(ctx, relationQuery)
		// the error will be handled later

	default:
		return cm.Errorf(cm.Unauthenticated, nil, "Mã xác thực không hợp lệ")
	}

	partner, err := s.validatePartner(ctx, partnerID)
	if err != nil {
		return err
	}
	if requestInfo.ShopID == 0 {
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		q.Result = resp
		return err
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
			resp, err := s.generateNewSession(ctx, nil, partner, shop, requestInfo)
			q.Result = resp
			return err
		}
		if shop.Status != status3.P || !shop.DeletedAt.IsZero() ||
			user.Status != status3.P {
			return cm.Errorf(cm.AccountClosed, nil, "")
		}
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		q.Result = resp
		return err

	case cm.NotFound:
		resp, err := s.actionRequestLogin(ctx, partner, requestInfo)
		q.Result = resp
		return err

	default:
		return cm.Errorf(cm.Internal, err, "")
	}
}

func (s *IntegrationService) validatePartner(ctx context.Context, partnerID dot.ID) (*identitymodel.Partner, error) {
	partnerQuery := &identitymodelx.GetPartner{PartnerID: partnerID}
	if err := bus.Dispatch(ctx, partnerQuery); err != nil {
		return nil, cm.MapError(err).Map(cm.NotFound, cm.PermissionDenied, "Mã xác thực không hợp lệ").Throw()
	}
	partner := partnerQuery.Result.Partner
	if partner.Status != status3.P || !partner.DeletedAt.IsZero() {
		return nil, cm.Errorf(cm.Unavailable, nil, "Tài khoản đối tác đã khoá hoặc không còn được sử dụng")
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
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}

	meta := map[string]string{}
	action := &integration.Action{
		Name:  "request_login",
		Label: fmt.Sprintf(`Kết nối với %v`, partner.PublicName),
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

	action.Msg = fmt.Sprintf("Nhập số điện thoại hoặc email để kết nối với %v", partner.PublicName)
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
			AccountID:     shop.ID,
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": jsonx.MustMarshalToString(info),
			},
		},
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	actions := getActionsFromConfig(info.Config)
	if len(actions) == 0 {
		actions = append(actions, &integration.Action{
			Name: "create_order",
		})
	}
	resp := generateShopLoginResponse(
		tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn,
		user, partner, shop, actions, info.RedirectURL,
	)
	return resp, nil
}

func (s *IntegrationService) RequestLogin(ctx context.Context, r *RequestLoginEndpoint) error {
	key := fmt.Sprintf("RequestLogin %v", r.Login)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			return s.requestLogin(ctx, r)
		}, "gửi mã đăng nhập")

	if err != nil {
		return err
	}
	r.Result = res.(*RequestLoginEndpoint).Result
	return err
}

// - Verify the token
// - Check whether the user exists in our database
// - Generate verification code and send to email/phone
// - Response action login_using_token
func (s *IntegrationService) requestLogin(ctx context.Context, r *RequestLoginEndpoint) (*RequestLoginEndpoint, error) {
	partner := r.CtxPartner
	if partner == nil {
		return r, cm.Errorf(cm.Internal, nil, "")
	}

	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(r.Context.Extra["request_login"]), &requestInfo); err != nil {
		return r, cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	emailNorm, phoneNorm, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return r, cm.Errorf(cm.InvalidArgument, nil, "Email hoặc số điện thoại không hợp lệ")
	}

	userQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: r.Login,
	}
	err := bus.Dispatch(ctx, userQuery)
	var exists bool
	switch cm.ErrorCode(err) {
	case cm.NoError:
		exists = true
	case cm.NotFound:
		exists = false
	default:
		return r, err
	}

	_ = exists
	user := userQuery.Result.User

	extra := map[string]string{
		"requested_email": emailNorm,
		"requested_phone": phoneNorm,
	}
	_, generatedCode, _, err := generateTokenWithVerificationCode(partner.ID, cm.Coalesce(emailNorm, phoneNorm), extra)
	if err != nil {
		return r, err
	}

	updateSessionCmd := &tokens.UpdateSessionCommand{
		Token:  r.Context.Claim.Token,
		Values: extra,
	}
	if err := bus.Dispatch(ctx, updateSessionCmd); err != nil {
		return r, err
	}

	var msg, notice string
	notice = "Bằng việc nhập mã xác nhận, bạn đồng ý cấp quyền cho đối tác tạo đơn hàng với tư cách của bạn."
	switch {
	case emailNorm != "":
		hello := "Xin chào"
		if user != nil {
			hello = "Gửi " + user.FullName
		}
		var extraMessage template.HTML
		if cmenv.NotProd() {
			extraMessage = "<br><br><i>Đây là email được gửi từ hệ thống của đối tác thông qua eTop.vn nhằm mục đích thử nghiệm. Nếu bạn cho rằng đây là sự nhầm lẫn, xin vui lòng thông báo cho chúng tôi.<i>"
		} else {
			extraMessage = "<br><br><i>Đây là email được gửi từ hệ thống của đối tác thông qua eTop.vn. Nếu bạn cho rằng đây là sự nhầm lẫn, xin vui lòng thông báo cho chúng tôi.<i>"
		}

		var b strings.Builder
		if err := api.RequestLoginEmailTpl.Execute(&b, map[string]interface{}{
			"Code":              formatVerificationCode(generatedCode),
			"Hello":             hello,
			"Email":             emailNorm,
			"PartnerPublicName": partner.PublicName,
			"PartnerWebsite":    validate.DomainFromURL(partner.WebsiteURL),
			"Notice":            notice,
			"Extra":             extraMessage,
		}); err != nil {
			return r, cm.Errorf(cm.Internal, err, "Không thể gửi mã đăng nhập").WithMeta("reason", "can not generate email content")
		}
		address := emailNorm
		cmd := &email.SendEmailCommand{
			FromName:    "eTop.vn (no-reply)",
			ToAddresses: []string{address},
			Subject:     fmt.Sprintf("Đăng nhập vào eTop.vn thông qua hệ thống %v", partner.PublicName),
			Content:     b.String(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return r, err
		}
		msg = fmt.Sprintf("Đã gửi email kèm mã xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", emailNorm)

	case phoneNorm != "":
		var b strings.Builder
		if err := api.RequestLoginSmsTpl.Execute(&b, map[string]interface{}{
			"Code":           formatVerificationCode(generatedCode),
			"PartnerWebsite": validate.DomainFromURL(partner.WebsiteURL),
			"Notice":         notice,
		}); err != nil {
			return r, cm.Errorf(cm.Internal, err, "Không thể gửi mã đăng nhập").WithMeta("reason", "can not generate sms content")
		}
		phone := phoneNorm
		cmd := &sms.SendSMSCommand{
			Phone:   phone,
			Content: b.String(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return r, err
		}
		msg = fmt.Sprintf("Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", phoneNorm)

	default:
		panic("unexpected")
	}

	r.Result = &integration.RequestLoginResponse{
		Code: "ok",
		Msg:  msg,
		Actions: []*integration.Action{
			{
				Name:  "login_using_token",
				Label: "Nhập mã xác nhận",
				Msg:   notice,
			},
		},
	}
	return r, nil
}

func formatVerificationCode(code string) string {
	if len(code) != 8 {
		panic("unexpected length")
	}
	return code[:4] + " " + code[4:]
}

func getToken(partnerID dot.ID, login string) (*auth.Token, string, map[string]string) {
	tokStr := fmt.Sprintf("%v-%v", partnerID, login)
	var v map[string]string
	var code string

	tok, err := authStore.Validate(auth.UsageRequestLogin, tokStr, &v)
	if err == nil && v != nil && len(v["code"]) == 8 {
		code = v["code"]
		return tok, code, v
	}

	_ = authStore.Revoke(auth.UsageRequestLogin, tokStr)
	return nil, "", nil
}

func generateTokenWithVerificationCode(partnerID dot.ID, login string, extra map[string]string) (*auth.Token, string, map[string]string, error) {
	tokStr := fmt.Sprintf("%v-%v", partnerID, login)
	tok, code, v := getToken(partnerID, login)
	if code != "" {
		return tok, code, v, nil
	}

	code, err := gencode.Random8Digits()
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
		return nil, "", nil, cm.Errorf(cm.Internal, err, "Không thể gửi mã đăng nhập").WithMeta("reason", "can not generate token")
	}
	return tok, code, v, nil
}

func (s *IntegrationService) LoginUsingToken(ctx context.Context, r *LoginUsingTokenEndpoint) (_err error) {
	if r.Login == "" || r.VerificationCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin")
	}
	r.VerificationCode = strings.Replace(r.VerificationCode, " ", "", -1)

	// verify shop_id and external_shop_id in the request info:
	//
	//     - shop_id == 0 && external_shop_id == "": New shop
	//     - shop_id != 0 && external_shop_id == "": Must use existing shop
	//     - shop_id == 0 && external_shop_id != "": New shop and update existing relation
	//     - shop_id != 0 && external_shop_id != "": Must validate that they match
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(r.Context.Extra["request_login"]), &requestInfo); err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	emailNorm, phoneNorm, ok := validate.NormalizeEmailOrPhone(r.Login)
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "Email hoặc số điện thoại không hợp lệ")
	}

	partner := r.CtxPartner
	getToken(partner.ID, cm.Coalesce(emailNorm, phoneNorm))
	key := fmt.Sprintf("%v-%v", partner.ID, cm.Coalesce(emailNorm, phoneNorm))

	var v map[string]string
	tok, err := authStore.Validate(auth.UsageRequestLogin, key, &v)
	if err != nil {
		return cm.Errorf(cm.Unauthenticated, nil, "Mã xác nhận không hợp lệ")
	}
	if v == nil || v["code"] != r.VerificationCode {
		return cm.Errorf(cm.Unauthenticated, nil, "Mã xác nhận không hợp lệ")
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
		tokenCmd := &tokens.GenerateTokenCommand{
			ClaimInfo: claims.ClaimInfo{
				AuthPartnerID: partner.ID,
				Extra: map[string]string{
					"action:register": "1",
					"verified_phone":  phoneNorm,
					"verified_email":  emailNorm,
					"request_login":   r.Context.Extra["request_login"],
				},
			},
		}
		if err := bus.Dispatch(ctx, tokenCmd); err != nil {
			return err
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
		r.Result = &integration.LoginResponse{
			AccessToken:       tokenCmd.Result.TokenStr,
			ExpiresIn:         tokenCmd.Result.ExpiresIn,
			User:              nil,
			Account:           nil,
			Shop:              nil,
			AvailableAccounts: nil,
			AuthPartner:       convertpb.PbPublicAccountInfo(partner),
			Actions:           actions,
		}
		return nil

	default:
		return err
	}

	user := userQuery.Result.User
	relationQuery := &identitymodelx.GetPartnerRelationsQuery{
		PartnerID: partner.ID,
		OwnerID:   user.ID,
	}
	err = bus.Dispatch(ctx, relationQuery)
	if err != nil {
		return err
	}

	accQuery := &identitymodelx.GetAllAccountRolesQuery{UserID: user.ID}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return err
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
				"request_login": r.Context.Extra["request_login"],
			},
		},
	}
	if err := bus.Dispatch(ctx, userTokenCmd); err != nil {
		return err
	}

	// we map from all accounts to partner relations and generate tokens for each one
	availableAccounts := make([]*integration.PartnerShopLoginAccount, 0, len(accQuery.Result))
	for _, acc := range accQuery.Result {
		if acc.Account.Type != account_type.Shop {
			continue
		}

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
			if rel.SubjectType == identitymodel.SubjectTypeAccount &&
				rel.SubjectID == acc.Account.ID {

				// the partner has permission to access this account
				tokenCmd := &tokens.GenerateTokenCommand{
					ClaimInfo: claims.ClaimInfo{
						UserID:        user.ID,
						AccountID:     acc.Account.ID,
						AuthPartnerID: partner.ID,
					},
				}
				if err := bus.Dispatch(ctx, tokenCmd); err != nil {
					return err
				}

				availAcc.ExternalId = rel.ExternalSubjectID
				availAcc.AccessToken = tokenCmd.Result.TokenStr
				availAcc.ExpiresIn = tokenCmd.Result.ExpiresIn
				break
			}
		}
		availableAccounts = append(availableAccounts, availAcc)
	}
	if requestInfo.ShopID != 0 && len(availableAccounts) == 0 {
		return cm.Errorf(cm.NotFound, nil, "Bạn đã từng liên kết với đối tác này, nhưng tài khoản cũ không còn hiệu lực (mã tài khoản %v). Liên hệ với hotro@etop.vn để được hướng dẫn.", requestInfo.ShopID).
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
				return cm.Errorf(cm.NotFound, nil, "Bạn đã từng liên kết với đối tác này, hãy sử dụng đúng tài khoản cũ của bạn trên trang web đối tác (mã tài khoản đối tác \"%v\"). Liên hệ với đối tác và cung cấp mã tài khoản đối tác ở trên để tìm lại tài khoản cũ của bạn.", strings.Join(externalIDs, `", "`)).
					WithMeta("reason", "external_shop_id not found")
			}
		}
	}

	r.Result = &integration.LoginResponse{
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
	return nil
}

func getActionsFromConfig(config string) (actions []*integration.Action) {
	if config == "" {
		return
	}
	_actions := strings.Split(config, ",")
	for _, a := range _actions {
		actions = append(actions, &integration.Action{
			Name: a,
		})
	}
	return
}

func (s *IntegrationService) Register(ctx context.Context, r *RegisterEndpoint) error {
	partner := r.CtxPartner
	claim := r.Context.ClaimInfo
	if claim.Extra == nil || claim.Extra["action:register"] == "" {
		return cm.Errorf(cm.PermissionDenied, nil, "Cần xác nhận email hoặc số điện thoại trước khi đăng ký")
	}

	if !r.AgreeTos {
		return cm.Error(cm.InvalidArgument, "Bạn cần đồng ý với điều khoản sử dụng dịch vụ để tiếp tục. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !r.AgreeEmailInfo.Valid {
		return cm.Error(cm.InvalidArgument, "Missing agree_email_info", nil)
	}

	phoneNorm, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ")
	}

	verifiedPhone, verifiedEmail := claim.Extra["verified_phone"], claim.Extra["verified_email"]
	if verifiedPhone != "" && verifiedPhone != string(phoneNorm) {
		return cm.Errorf(cm.InvalidArgument, nil, "Cần sử dụng số điện thoại đã được xác nhận khi đăng ký")
	}
	if verifiedEmail != "" && verifiedEmail != string(emailNorm) {
		return cm.Errorf(cm.InvalidArgument, nil, "Cần sử dụng địa chỉ email đã được xác nhận khi đăng ký")
	}

	generatedPassword := gencode.GenerateCode(gencode.Alphabet32, 8)
	// set default source: "partner"
	source := user_source.Partner

	cmd := &usering.CreateUserCommand{
		UserInner: identitymodel.UserInner{
			FullName:  r.FullName,
			ShortName: "",
			Phone:     string(phoneNorm),
			Email:     string(emailNorm),
		},
		Password:       generatedPassword,
		Status:         status3.P,
		AgreeTOS:       r.AgreeTos,
		AgreeEmailInfo: r.AgreeEmailInfo.Bool,
		Source:         source,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	user := cmd.Result.User

	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:        user.ID,
			AuthPartnerID: partner.ID,
			Extra: map[string]string{
				"request_login": r.Context.Extra["request_login"],
			},
		},
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return err
	}

	switch {
	case verifiedEmail != "":
		verifyCmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			EmailVerifiedAt: time.Now(),
		}
		if err := bus.Dispatch(ctx, verifyCmd); err != nil {
			ll.Error("Can not update verification", l.Error(err))
		}

		var b strings.Builder
		if err := api.NewAccountViaPartnerEmailTpl.Execute(&b, map[string]interface{}{
			"FullName":          user.FullName,
			"PartnerPublicName": partner.PublicName,
			"PartnerWebsite":    validate.DomainFromURL(partner.WebsiteURL),
			"LoginLabel":        "Email",
			"Login":             verifiedEmail,
			"Password":          generatedPassword,
		}); err != nil {
			ll.Error("Can not send email", l.Error(err))
			break
		}
		emailCmd := &email.SendEmailCommand{
			FromName:    "eTop.vn (no-reply)",
			ToAddresses: []string{string(emailNorm)},
			Subject:     "Mật khẩu đăng nhập vào tài khoản ở eTop.vn",
			Content:     b.String(),
		}
		if err := bus.Dispatch(ctx, emailCmd); err != nil {
			ll.Error("Can not send email", l.Error(err))
		}

	case verifiedPhone != "":
		verifyCmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			PhoneVerifiedAt: time.Now(),
		}
		if err := bus.Dispatch(ctx, verifyCmd); err != nil {
			ll.Error("Can not update verification", l.Error(err))
		}

		var b strings.Builder
		if err := api.NewAccountViaPartnerSmsTpl.Execute(&b, map[string]interface{}{
			"Password": generatedPassword,
		}); err != nil {
			ll.Error("Can not send email", l.Error(err))
			break
		}
		smsCmd := &sms.SendSMSCommand{
			Phone:   verifiedPhone,
			Content: b.String(),
		}
		if err := bus.Dispatch(ctx, smsCmd); err != nil {
			ll.Error("Can not send sms", l.Error(err))
		}
	}

	r.Result = &integration.RegisterResponse{
		User:        convertpb.PbUser(user),
		AccessToken: tokenCmd.Result.TokenStr,
		ExpiresIn:   tokenCmd.Result.ExpiresIn,
	}
	return nil
}

func (s *IntegrationService) GrantAccess(ctx context.Context, r *GrantAccessEndpoint) error {
	var requestInfo apipartner.PartnerShopToken
	if err := jsonx.Unmarshal([]byte(r.Context.Extra["request_login"]), &requestInfo); err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu đăng nhập không còn hiệu lực")
	}

	partner := r.CtxPartner
	user := r.Context.User
	if r.ShopId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}

	if requestInfo.ShopID != 0 && r.ShopId != requestInfo.ShopID {
		return cm.Errorf(cm.FailedPrecondition, nil, "Bạn cần sử dụng tài khoản đã từng liên kết với đối tác này").WithMeta("reason", "shop_id does not match")
	}

	shopQuery := &identitymodelx.GetShopQuery{
		ShopID: r.ShopId,
	}
	if err := bus.Dispatch(ctx, shopQuery); err != nil {
		return err
	}
	shop := shopQuery.Result
	if shop.OwnerID != user.ID {
		return cm.Errorf(cm.NotFound, nil, "")
	}

	relQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID: partner.ID,
		AccountID: shop.ID,
	}
	err := bus.Dispatch(ctx, relQuery)
	switch cm.ErrorCode(err) {
	case cm.OK:
		// verify the external_shop_id or automatically fill it
		rel := relQuery.Result
		if rel.ExternalSubjectID != "" && requestInfo.ExternalShopID != "" &&
			rel.ExternalSubjectID != requestInfo.ExternalShopID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Bạn cần sử dụng tài khoản đã từng liên kết với đối tác này").WithMeta("reason", "external_shop_id does not match")
		}
		if rel.ExternalSubjectID == "" && requestInfo.ExternalShopID != "" {
			cmd := &identitymodelx.UpdatePartnerRelationCommand{
				PartnerID:  partner.ID,
				AccountID:  r.ShopId,
				ExternalID: requestInfo.ExternalShopID,
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return cm.Errorf(cm.Internal, err, "Không thể cập nhật thông tin tài khoản. Vui lòng liên hệ hotro@etop.vn để được hỗ trợ.").
					WithMeta("reason", "can not update external_shop_id")
			}
		}

	case cm.NotFound:
		cmd := &identitymodelx.CreatePartnerRelationCommand{
			PartnerID:  partner.ID,
			AccountID:  r.ShopId,
			ExternalID: requestInfo.ExternalShopID,
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}

	default:
		return err
	}

	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{
			UserID:        user.ID,
			AccountID:     shop.ID,
			AuthPartnerID: partner.ID,
		},
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return err
	}

	r.Result = &integration.GrantAccessResponse{
		AccessToken: tokenCmd.Result.TokenStr,
		ExpiresIn:   tokenCmd.Result.ExpiresIn,
	}
	return nil
}

func (s *IntegrationService) SessionInfo(ctx context.Context, q *SessionInfoEndpoint) error {
	var shop *identitymodel.Shop
	if q.Context.Claim.AccountID != 0 {
		query := &identitymodelx.GetShopQuery{
			ShopID: q.Context.Claim.AccountID,
		}
		err := bus.Dispatch(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.OK, cm.NotFound:
			// continue
		default:
			return err
		}
		shop = query.Result
	}
	claim := q.Context.Claim
	var actions []*integration.Action
	redirectURL := ""
	if claim.Extra != nil && claim.Extra["request_login"] != "" {
		var requestInfo apipartner.PartnerShopToken
		jsonx.Unmarshal([]byte(claim.Extra["request_login"]), &requestInfo)
		actions = getActionsFromConfig(requestInfo.Config)
		redirectURL = requestInfo.RedirectURL
	}

	q.Result = generateShopLoginResponse(
		q.Context.Token, tokens.DefaultAccessTokenTTL,
		nil, q.CtxPartner, shop, actions, redirectURL,
	)
	return nil
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
