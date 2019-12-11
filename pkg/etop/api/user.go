package api

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	apisms "etop.vn/api/etc/logging/smslog"
	"etop.vn/api/main/authorization"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/api/top/int/etop"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/api/convertpb"
	authservice "etop.vn/backend/pkg/etop/authorize/auth"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/logic/usering"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	eventBus               capi.EventBus
	idempgroup             *idemp.RedisGroup
	authStore              auth.Generator
	enabledEmail           bool
	enabledSMS             bool
	smsAggr                apisms.CommandBus
	cfgEmail               EmailConfig
	identityAggr           identity.CommandBus
	identityQS             identity.QueryBus
	invitationAggregate    invitation.CommandBus
	invitationQuery        invitation.QueryBus
	authorizationQuery     authorization.QueryBus
	authorizationAggregate authorization.CommandBus
)

const PrefixIdempUser = "IdempUser"

const (
	keyRequestVerifyCode                string = "request_phone_verification_code"
	keyRequestVerifyPhone               string = "request_phone_verification"
	keyRequestPhoneVerificationVerified string = "request_phone_verification_verified"
)

func init() {
	bus.AddHandlers("api",
		userService.UpdatePermission,
		userService.ChangePassword,
		userService.ChangePasswordUsingToken,
		userService.Login,
		userService.Register,
		userService.ResetPassword,
		userService.SendEmailVerification,
		userService.SendPhoneVerification,
		userService.SendSTokenEmail,
		userService.SessionInfo,
		userService.SwitchAccount,
		userService.UpgradeAccessToken,
		userService.VerifyEmailUsingToken,
		userService.VerifyPhoneUsingToken,
		userService.UpdateReferenceUser,
		userService.UpdateReferenceSale,
		userService.CheckUserRegistration,
		userService.RegisterUsingToken,
	)
}

type UserService struct{}

var userService = &UserService{}

type EmailConfig = config.EmailConfig

func Init(
	bus capi.EventBus,
	smsCommandBus apisms.CommandBus,
	identityCommandBus identity.CommandBus,
	identityQueryBus identity.QueryBus,
	invitationAggr invitation.CommandBus,
	invitationQS invitation.QueryBus,
	authorizationQ authorization.QueryBus,
	authorizationA authorization.CommandBus,
	sd cmservice.Shutdowner,
	rd redis.Store,
	s auth.Generator,
	_cfgEmail EmailConfig,
	_cfgSMS sms.Config,
) {
	eventBus = bus
	smsAggr = smsCommandBus
	identityAggr = identityCommandBus
	identityQS = identityQueryBus
	invitationAggregate = invitationAggr
	invitationQuery = invitationQS
	authorizationQuery = authorizationQ
	authorizationAggregate = authorizationA
	authStore = s
	enabledEmail = _cfgEmail.Enabled
	enabledSMS = _cfgSMS.Enabled
	cfgEmail = _cfgEmail
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempUser, 0)
	sd.Register(idempgroup.Shutdown)

	if enabledEmail {
		if _, err := validate.ValidateStruct(cfgEmail); err != nil {
			ll.Fatal("Can not validate config", l.Error(err))
		}
	}
}

type authKey struct{}

// Register
// 1a. If the user does not exist in the system, create a new user.
// 1b. If any email or phone is activated -> AlreadyExists.
//   - If both email and phone exist (but not activated) -> Merge them.
//   - Otherwise, update existing user with the other identifier.
func (s *UserService) Register(ctx context.Context, r *RegisterEndpoint) error {
	if err := validateRegister(r.CreateUserRequest); err != nil {
		return err
	}
	resp, err := s.register(ctx, r.Context, r.CreateUserRequest, false)
	r.Result = resp
	return err
}

func (s *UserService) RegisterUsingToken(ctx context.Context, r *RegisterUsingTokenEndpoint) error {
	if r.Context.Extra[keyRequestVerifyPhone] == "" || r.Context.Extra[keyRequestPhoneVerificationVerified] == "" {
		return cm.Error(cm.InvalidArgument, "Bạn vui lòng xác nhận só điện thoại trước khi đăng kí", nil)
	}
	if err := validateRegister(r.CreateUserRequest); err != nil {
		return err
	}

	phoneNorm, _ := validate.NormalizePhone(r.Phone)
	if r.Context.Extra[keyRequestVerifyPhone] != phoneNorm.String() {
		return cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại %v không hợp lệ bởi vì bạn đã xác nhận số điện thoại: %v", phoneNorm.String(), r.Context.Extra[keyRequestVerifyPhone])
	}

	resp, err := s.register(ctx, r.Context, r.CreateUserRequest, true)
	r.Result = resp
	return err
}

func (s *UserService) register(
	ctx context.Context,
	claim claims.EmptyClaim,
	r *etop.CreateUserRequest,
	usingToken bool,
) (*etop.RegisterResponse, error) {
	phoneNorm, _ := validate.NormalizePhone(r.Phone)
	emailNorm, _ := validate.NormalizeEmail(r.Email)
	userByPhone, userByEmail, err := s.getUserByPhoneAndByEmail(ctx, phoneNorm.String(), emailNorm.String())
	if err != nil {
		return nil, err
	}

	switch {
	case userByPhone.User != nil:
		return nil, cm.Errorf(cm.AlreadyExists, nil, "Số điện thoại đã được đăng ký. Vui lòng đăng nhập hoặc sử dụng số điện thoại khác")
	case userByEmail.User != nil:
		return nil, cm.Errorf(cm.AlreadyExists, nil, "Email đã được đăng ký. Vui lòng đăng nhập hoặc sử dụng số điện thoại khác")
	}

	var invitationTemp *invitation.Invitation
	if strings.HasPrefix(r.RegisterToken, "iv:") {
		invitationTemp, err = getInvitation(ctx, r)
		if err != nil {
			return nil, err
		}
	}

	user, err := createUser(ctx, r)
	if err != nil {
		return nil, err
	}
	{
		now := time.Now()
		shouldUpdate := false
		updateCmd := &model.UpdateUserVerificationCommand{
			UserID: user.ID,
		}
		// auto verify email when accept invitation from email
		if invitationTemp != nil {
			if invitationTemp.Email != "" && user.EmailVerifiedAt.IsZero() {
				updateCmd.EmailVerifiedAt = now
				shouldUpdate = true
			}
		}
		// auto verify phone when register using token
		if usingToken {
			if user.PhoneVerificationSentAt.IsZero() {
				updateCmd.PhoneVerifiedAt = now
				shouldUpdate = true
			}
		}
		if shouldUpdate {
			if err := bus.Dispatch(ctx, updateCmd); err != nil {
				return nil, err
			}
		}
	}
	{
		event := &identity.UserCreatedEvent{
			UserID:    user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			ShortName: user.ShortName,
		}
		if invitationTemp != nil {
			event.Invitation = &identity.UserInvitation{
				Token:      r.RegisterToken,
				AutoAccept: r.AutoAcceptInvitation,
				FullName:   invitationTemp.FullName,
				ShortName:  invitationTemp.ShortName,
				Position:   invitationTemp.Position,
			}
		}
		if err := eventBus.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	return &etop.RegisterResponse{
		User: convertpb.PbUser(user),
	}, nil
}

func createUser(ctx context.Context, r *etop.CreateUserRequest) (*model.User, error) {
	info := r
	cmd := &usering.CreateUserCommand{
		UserInner: model.UserInner{
			FullName:  info.FullName,
			ShortName: info.ShortName,
			Phone:     info.Phone,
			Email:     info.Email,
		},
		Password:       info.Password,
		Status:         model.StatusActive,
		AgreeTOS:       r.AgreeTos,
		AgreeEmailInfo: r.AgreeEmailInfo.Bool,
		Source:         r.Source,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	query := &model.GetUserByIDQuery{
		UserID: cmd.Result.User.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func validateRegister(r *etop.CreateUserRequest) error {
	if !r.AgreeTos {
		return cm.Error(cm.InvalidArgument, "Bạn cần đồng ý với điều khoản sử dụng dịch vụ để tiếp tục. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !r.AgreeEmailInfo.Valid {
		return cm.Error(cm.InvalidArgument, "Missing agree_email_info", nil)
	}
	if r.Phone == "" {
		return cm.Error(cm.InvalidArgument, "Vui lòng cung cấp số điện thoại", nil)
	}
	if r.Email == "" && r.RegisterToken == "" {
		return cm.Error(cm.InvalidArgument, "Vui lòng cung cấp địa chỉ email", nil)
	}
	_, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}
	if r.Email != "" {
		_, ok := validate.NormalizeEmail(r.Email)
		if !ok {
			return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
	}
	return nil
}

func getInvitation(ctx context.Context, r *etop.CreateUserRequest) (*invitation.Invitation, error) {
	getInvitationByToken := &invitation.GetInvitationByTokenQuery{
		Token:  r.RegisterToken,
		Result: nil,
	}
	if err := invitationQuery.Dispatch(ctx, getInvitationByToken); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Token không hợp lệ").
			Throw()
	}
	invitationTemp := getInvitationByToken.Result
	if r.Email != invitationTemp.Email {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email gửi lên và email trong token không khớp nhau")
	}
	return invitationTemp, nil
}

func (s *UserService) CheckUserRegistration(ctx context.Context, q *CheckUserRegistrationEndpoint) error {
	_, ok := validate.NormalizePhone(q.Phone)
	if !ok {
		q.Result = &etop.GetUserByPhoneResponse{Exists: false}
		return nil
	}

	userByPhoneQuery := &model.GetUserByLoginQuery{
		PhoneOrEmail: q.Phone,
	}
	err := bus.Dispatch(ctx, userByPhoneQuery)
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return err
	}
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		q.Result = &etop.GetUserByPhoneResponse{Exists: false}
		return nil
	}
	q.Result = &etop.GetUserByPhoneResponse{Exists: true}
	return nil
}

func (s *UserService) getUserByPhoneAndByEmail(ctx context.Context, phone, email string) (userByPhone, userByEmail model.UserExtended, err error) {
	userByPhoneQuery := &model.GetUserByLoginQuery{
		PhoneOrEmail: phone,
	}
	if err := bus.Dispatch(ctx, userByPhoneQuery); err != nil &&
		cm.ErrorCode(err) != cm.NotFound {
		return model.UserExtended{}, model.UserExtended{}, err
	}
	userByPhone = userByPhoneQuery.Result

	if email != "" {
		userByEmailQuery := &model.GetUserByLoginQuery{
			PhoneOrEmail: email,
		}
		if err := bus.Dispatch(ctx, userByEmailQuery); err != nil &&
			cm.ErrorCode(err) != cm.NotFound {
			return model.UserExtended{}, model.UserExtended{}, err
		}
		userByEmail = userByEmailQuery.Result
	}
	return
}

func (s *UserService) Login(ctx context.Context, r *LoginEndpoint) error {
	query := &login.LoginUserQuery{
		PhoneOrEmail: r.Login,
		Password:     r.Password,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	user := query.Result.User
	resp, err := s.CreateLoginResponse(
		ctx, nil, "", user.ID, user,
		r.AccountId, int(r.AccountType),
		true, // Generate tokens for all accounts
		0,
	)
	r.Result = resp
	return err
}

func (s *UserService) ResetPassword(ctx context.Context, r *ResetPasswordEndpoint) error {
	key := fmt.Sprintf("ResetPassword %v", r.Email)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			return s.resetPassword(ctx, r)
		}, "gửi email khôi phục mật khẩu")

	if err != nil {
		return err
	}
	r.Result = res.(*ResetPasswordEndpoint).Result
	return err
}

func (s *UserService) resetPassword(ctx context.Context, r *ResetPasswordEndpoint) (*ResetPasswordEndpoint, error) {
	if !enabledEmail {
		return r, cm.Error(cm.FailedPrecondition, "Không thể gửi email khôi phục mật khẩu. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return r, cm.Error(cm.FailedPrecondition, "Thiếu thông tin email. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !strings.Contains(r.Email, "@") {
		return r, cm.Error(cm.FailedPrecondition, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}

	query := &model.GetUserByLoginQuery{
		PhoneOrEmail: r.Email,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return r, cm.MapError(err).
			Wrap(cm.NotFound, "Người dùng chưa đăng ký. Vui lòng kiểm tra lại thông tin (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.").
			Throw()
	}
	user := query.Result.User

	tok := &auth.Token{
		TokenStr: "",
		Usage:    auth.UsageResetPassword,
		UserID:   user.ID,
		Value: map[string]string{
			"email": user.Email,
		},
	}
	if _, err := authStore.GenerateWithValue(tok, 24*60*60); err != nil {
		return r, cm.Errorf(cm.Internal, err, "Không thể khôi phục mật khẩu").WithMeta("reason", "can not generate token")
	}

	resetUrl, err := url.Parse(cfgEmail.ResetPasswordURL)
	if err != nil {
		return r, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := resetUrl.Query()
	urlQuery.Set("t", tok.TokenStr)
	resetUrl.RawQuery = urlQuery.Encode()

	var b strings.Builder
	if err := resetPasswordTpl.Execute(&b, map[string]interface{}{
		"FullName": user.FullName,
		"URL":      resetUrl.String(),
		"Email":    user.Email,
	}); err != nil {
		return r, cm.Errorf(cm.Internal, err, "Không thể khôi phục mật khẩu").WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    "eTop.vn (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Khôi phục mật khẩu eTop",
		Content:     b.String(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}
	r.Result = cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email khôi phục mật khẩu đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", address))
	return r, nil
}

func (s *UserService) ChangePassword(ctx context.Context, r *ChangePasswordEndpoint) error {
	if r.CurrentPassword == "" {
		return cm.Error(cm.InvalidArgument, "Missing current_password", nil)
	}

	query := &login.LoginUserQuery{
		UserID:   r.Context.User.ID,
		Password: r.CurrentPassword,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.Unauthenticated, "Mật khẩu không đúng. Vui lòng kiểm tra lại thông tin đăng nhập. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.").
			Throw()
	}

	if len(r.NewPassword) < 8 {
		return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.NewPassword != r.ConfirmPassword {
		return cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &model.SetPasswordCommand{
		UserID:   r.Context.User.ID,
		Password: r.NewPassword,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.Empty{}
	return nil
}

func (s *UserService) ChangePasswordUsingToken(ctx context.Context, r *ChangePasswordUsingTokenEndpoint) error {
	key := fmt.Sprintf("ChangePasswordUsingToken %v-%v-%v", r.ResetPasswordToken, r.NewPassword, r.ConfirmPassword)
	res, err := idempgroup.DoAndWrap(key, 30*time.Second,
		func() (interface{}, error) {
			return s.changePasswordUsingToken(ctx, r)
		}, "khôi phục mật khẩu")

	if err != nil {
		return err
	}
	r.Result = res.(*ChangePasswordUsingTokenEndpoint).Result
	return err
}

func (s *UserService) changePasswordUsingToken(ctx context.Context, r *ChangePasswordUsingTokenEndpoint) (*ChangePasswordUsingTokenEndpoint, error) {
	if r.ResetPasswordToken == "" {
		return r, cm.Error(cm.InvalidArgument, "Missing reset_password_token", nil)
	}

	var v map[string]string
	tok, err := authStore.Validate(auth.UsageResetPassword, r.ResetPasswordToken, &v)
	if err != nil {
		return r, cm.Error(cm.InvalidArgument, "Không thể khôi phục mật khẩu (token không hợp lệ). Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", err)
	}

	query := &model.GetUserByLoginQuery{
		PhoneOrEmail: v["email"],
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return r, err
	}
	user := query.Result.User

	if tok.UserID != user.ID {
		return r, cm.Error(cm.InvalidArgument, "Không thể khôi phục mật khẩu (token không hợp lệ). Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil).WithMeta("reason", "user is not correct")
	}

	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok {
		return r, cm.Error(cm.InvalidArgument, "Email không hợp lệ. Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}
	if emailNorm.String() != user.Email {
		return r, cm.Error(cm.InvalidArgument, "Email không đúng. Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}

	if len(r.NewPassword) < 8 {
		return r, cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.NewPassword != r.ConfirmPassword {
		return r, cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &model.SetPasswordCommand{
		UserID:   user.ID,
		Password: r.NewPassword,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}

	authStore.Revoke(auth.UsageResetPassword, r.ResetPasswordToken)
	r.Result = &pbcm.Empty{}
	return r, nil
}

func (s *UserService) SessionInfo(ctx context.Context, r *SessionInfoEndpoint) error {
	resp, err := s.CreateLoginResponse(
		ctx,
		&r.Context.ClaimInfo,
		r.Context.Token,
		r.Context.UserID,
		r.Context.User.User,
		r.Context.AccountID,
		0,
		false,
		0,
	)
	r.Result = resp
	return err
}

func (s *UserService) SwitchAccount(ctx context.Context, r *SwitchAccountEndpoint) error {
	if r.AccountId == 0 && !r.RegenerateTokens {
		return cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}
	resp, err := s.CreateSessionResponse(
		ctx,
		nil, // Do not forward claim data
		"",  // Empty to generate new token
		r.Context.UserID,
		r.Context.User.User,
		r.AccountId,
		0,
		0,
	)
	if err != nil {
		return err
	}
	if resp.Account == nil {
		return cm.Error(cm.PermissionDenied, "Tài khoản không hợp lệ.", nil)
	}
	r.Result = resp
	return err
}

func (s *UserService) CreateSessionResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, adminID dot.ID) (*etop.AccessTokenResponse, error) {
	resp, err := s.CreateLoginResponse(ctx, claim, token, userID, user, preferAccountID, preferAccountType, false, adminID)
	if err != nil {
		return nil, err
	}
	return &etop.AccessTokenResponse{
		AccessToken:     resp.AccessToken,
		ExpiresIn:       resp.ExpiresIn,
		User:            resp.User,
		Account:         resp.Account,
		Shop:            resp.Shop,
		Affiliate:       resp.Affiliate,
		Stoken:          resp.Stoken,
		StokenExpiresAt: resp.StokenExpiresAt,
	}, nil
}

func (s *UserService) CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (*etop.LoginResponse, error) {
	resp, _, err := s.CreateLoginResponse2(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
	return resp, err
}

func (s *UserService) CreateLoginResponse2(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (_ *etop.LoginResponse, respShop *model.Shop, _ error) {

	// Retrieve user info
	if user != nil && user.ID != userID {
		return nil, nil, cm.Error(cm.Internal, "Invalid user", nil)
	}
	if user == nil {
		userQuery := &model.GetUserByIDQuery{UserID: userID}
		if err := bus.Dispatch(ctx, userQuery); err != nil {
			return nil, nil, err
		}
		user = userQuery.Result
	}

	// Retrieve list of accounts
	accQuery := &model.GetAllAccountRolesQuery{UserID: userID}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return nil, nil, err
	}

	if preferAccountID != 0 && preferAccountType != 0 {
		return nil, nil, cm.Error(cm.InvalidArgument, "Can not set both account_id and account_type", nil)
	}

	var currentAccount *etop.LoginAccount
	var currentAccountID dot.ID
	availableAccounts := make([]*etop.LoginAccount, len(accQuery.Result))
	for i, accUserX := range accQuery.Result {
		availableAccounts[i] = convertpb.PbLoginAccount(accUserX)
		account := accUserX.Account
		switch {
		case preferAccountID == account.ID,
			preferAccountType == model.TagShop && account.Type == account_type.Shop,
			preferAccountType == model.TagEtop && account.Type == account_type.Etop,
			preferAccountType == model.TagAffiliate && account.Type == account_type.Affiliate:
			currentAccount = availableAccounts[i]
			currentAccountID = currentAccount.Id
		}
	}

	resp := &etop.LoginResponse{
		User:              convertpb.PbUser(user),
		Account:           currentAccount,
		AvailableAccounts: availableAccounts,
	}

	// Retrieve shop account
	if currentAccount != nil {
		switch {
		case model.IsShopID(currentAccountID):
			query := &model.GetShopExtendedQuery{ShopID: currentAccountID}
			if err := bus.Dispatch(ctx, query); err != nil {
				return nil, nil, cm.ErrorTracef(cm.Internal, err, "")
			}
			resp.Shop = convertpb.PbShopExtended(query.Result)
			respShop = query.Result.Shop

		case model.IsAffiliateID(currentAccountID):
			query := &identity.GetAffiliateByIDQuery{ID: currentAccountID}
			if err := identityQS.Dispatch(ctx, query); err != nil {
				return nil, nil, cm.ErrorTracef(cm.Internal, err, "Account affiliate not found")
			}
			resp.Affiliate = convertpb.Convert_core_Affiliate_To_api_Affiliate(query.Result)
		case model.IsEtopAccountID(currentAccountID):
			// nothing
		default:
			return nil, nil, cm.ErrorTracef(cm.Internal, nil, "Invalid account")
		}
	}

	// MixedAccount
	accountIDs := make(map[dot.ID]int)
	for _, acc := range availableAccounts {
		accountIDs[acc.Id] = int(acc.Type)
	}

	// Generate new access token.
	//
	// TODO: Invalidate / refresh the token
	if token == "" {
		tokenCmd := &tokens.GenerateTokenCommand{
			ClaimInfo: claims.ClaimInfo{
				UserID:     userID,
				AdminID:    adminID,
				AccountIDs: accountIDs,
			},
		}
		if currentAccount != nil {
			tokenCmd.AccountID = currentAccount.Id
		}
		if claim != nil && claim.STokenExpiresAt != nil {
			tokenCmd.ClaimInfo.SToken = claim.SToken
			tokenCmd.STokenExpiresAt = claim.STokenExpiresAt
		}
		if err := bus.Dispatch(ctx, tokenCmd); err != nil {
			return nil, nil, err
		}
		token = tokenCmd.Result.TokenStr
	}

	const ttl = tokens.DefaultAccessTokenTTL // TODO: use tokenCmd.Result.ExpiresIn
	if generateAllTokens {
		for _, acc := range availableAccounts {
			if acc.Id == currentAccountID {
				acc.AccessToken = token
				acc.ExpiresIn = ttl
				continue
			}
			tokenCmd := &tokens.GenerateTokenCommand{
				ClaimInfo: claims.ClaimInfo{
					UserID:     userID,
					AccountID:  acc.Id,
					AccountIDs: accountIDs,
				},
			}
			if err := bus.Dispatch(ctx, tokenCmd); err != nil {
				return nil, nil, err
			}
			acc.AccessToken = tokenCmd.Result.TokenStr
			acc.ExpiresIn = tokenCmd.Result.ExpiresIn
		}
	}

	resp.AccessToken = token
	resp.ExpiresIn = ttl

	// UpdateInfo response from claim
	//
	// TODO: refactor due to duplicated with token generation above
	if claim != nil && claim.STokenExpiresAt != nil {
		resp.Stoken = claim.SToken
		resp.StokenExpiresAt = cmapi.PbTime(*claim.STokenExpiresAt)
	}

	// Add actions into permission
	if resp.Account != nil {
		resp.Account.UserAccount.Permission.Permissions = authservice.ListActionsByRoles(resp.Account.UserAccount.Permission.Roles)
	}
	for _, account := range resp.AvailableAccounts {
		account.UserAccount.Permission.Permissions = authservice.ListActionsByRoles(account.UserAccount.Permission.Roles)
	}

	return resp, respShop, nil
}

func (s *UserService) SendEmailVerification(ctx context.Context, r *SendEmailVerificationEndpoint) error {
	key := fmt.Sprintf("SendEmailVerification %v-%v", r.Context.User.ID, r.Email)
	res, err := idempgroup.DoAndWrap(key, 30*time.Second,
		func() (interface{}, error) {
			return s.sendEmailVerification(ctx, r)
		}, "gửi email xác nhận tài khoản")

	if err != nil {
		return err
	}
	r.Result = res.(*SendEmailVerificationEndpoint).Result
	return err
}

func (s *UserService) sendEmailVerification(ctx context.Context, r *SendEmailVerificationEndpoint) (*SendEmailVerificationEndpoint, error) {
	if !enabledEmail {
		return r, cm.Error(cm.FailedPrecondition, "Không thể gửi email xác nhận tài khoản. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return r, cm.Error(cm.FailedPrecondition, "Thiếu thông tin địa chỉ email. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !strings.Contains(r.Email, "@") {
		return r, cm.Error(cm.FailedPrecondition, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}

	user := r.Context.User.User
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok || user.Email != emailNorm.String() {
		return r, cm.Error(cm.FailedPrecondition, "Địa chỉ email không đúng. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !user.EmailVerifiedAt.IsZero() {
		r.Result = cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
		return r, nil
	}

	tok := &auth.Token{
		TokenStr: "",
		Usage:    auth.UsageEmailVerification,
		UserID:   user.ID,
		Value: map[string]string{
			"email": user.Email,
		},
	}
	if _, err := authStore.GenerateWithValue(tok, 24*60*60); err != nil {
		return r, cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate token")
	}

	verificationUrl, err := url.Parse(cfgEmail.EmailVerificationURL)
	if err != nil {
		return r, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := verificationUrl.Query()
	urlQuery.Set("t", tok.TokenStr)
	verificationUrl.RawQuery = urlQuery.Encode()

	var b strings.Builder
	if err := emailVerificationTpl.Execute(&b, map[string]interface{}{
		"FullName": user.FullName,
		"URL":      verificationUrl.String(),
		"Email":    user.Email,
	}); err != nil {
		return r, cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    "eTop.vn (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận địa chỉ email",
		Content:     b.String(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}
	r.Result = cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", address))

	updateCmd := &model.UpdateUserVerificationCommand{
		UserID:                  user.ID,
		EmailVerificationSentAt: time.Now(),
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return r, err
	}
	return r, nil
}

func (s *UserService) SendPhoneVerification(ctx context.Context, r *SendPhoneVerificationEndpoint) error {
	key := fmt.Sprintf("SendPhoneVerification %v-%v", r.Context.Token, r.Phone)
	res, err := idempgroup.DoAndWrap(key, 2*time.Hour,
		func() (interface{}, error) {
			return s.sendPhoneVerification(ctx, r)
		}, "gửi tin nhắn xác nhận số điện thoại")

	if err != nil {
		return err
	}
	r.Result = res.(*SendPhoneVerificationEndpoint).Result
	return err
}

func (s *UserService) sendPhoneVerification(ctx context.Context, r *SendPhoneVerificationEndpoint) (*SendPhoneVerificationEndpoint, error) {
	if !enabledSMS {
		return r, cm.Error(cm.FailedPrecondition, "Không thể gửi tin nhắn xác nhận tài khoản. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil).WithMeta("reason", "not configured")
	}
	if r.Phone == "" {
		return r, cm.Error(cm.FailedPrecondition, "Thiếu thông tin số điện thoại. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	// update token when user not exists
	if r.Context.UserID == 0 {
		return sendPhoneVerificationForRegister(ctx, r)
	}
	getUserByID := &model.GetUserByIDQuery{
		UserID: r.Context.UserID,
	}
	if err := bus.Dispatch(ctx, getUserByID); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return r, err
	}
	user := getUserByID.Result
	phoneNorm, ok := validate.NormalizePhone(r.Phone)

	if !ok || user.Phone != phoneNorm.String() {
		return r, cm.Error(cm.FailedPrecondition, "Số điện thoại không đúng. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !user.PhoneVerifiedAt.IsZero() {
		r.Result = cmapi.Message("ok", "Số điện thoại đã được xác nhận thành công.")
		return r, nil
	}
	_, code, _, err := generateToken(auth.UsagePhoneVerification, user.ID, true, 2*60*60, user.Phone)
	if err != nil {
		return r, err
	}
	msg := fmt.Sprintf(smsVerificationTpl, code, user.Phone)
	phone := user.Phone
	cmd := &sms.SendSMSCommand{
		Phone:   phone,
		Content: msg,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}
	r.Result = cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", phone))

	updateCmd := &model.UpdateUserVerificationCommand{
		UserID:                  user.ID,
		PhoneVerificationSentAt: time.Now(),
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return r, err
	}
	return r, nil
}
func (s *UserService) VerifyEmailUsingToken(ctx context.Context, r *VerifyEmailUsingTokenEndpoint) error {
	key := fmt.Sprintf("VerifyEmailUsingToken %v-%v", r.Context.User.ID, r.VerificationToken)
	res, err := idempgroup.DoAndWrap(key, 30*time.Second,
		func() (interface{}, error) {
			return s.verifyEmailUsingToken(ctx, r)
		}, "xác nhận địa chỉ email")

	if err != nil {
		return err
	}
	r.Result = res.(*VerifyEmailUsingTokenEndpoint).Result
	return err
}

func (s *UserService) verifyEmailUsingToken(ctx context.Context, r *VerifyEmailUsingTokenEndpoint) (*VerifyEmailUsingTokenEndpoint, error) {
	if r.VerificationToken == "" {
		return r, cm.Error(cm.InvalidArgument, "Missing verification_token", nil)
	}

	var v map[string]string
	tok, err := authStore.Validate(auth.UsageEmailVerification, r.VerificationToken, &v)
	if err != nil {
		return r, cm.Error(cm.InvalidArgument, "Không thể xác nhận địa chỉ email (token không hợp lệ). Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}

	user := r.Context.User.User
	if user.ID != tok.UserID || user.Email != v["email"] {
		return r, cm.Error(cm.InvalidArgument, "Không thể xác nhận địa chỉ email (địa chỉ email không đúng). Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}

	if user.EmailVerifiedAt.IsZero() {
		cmd := &model.UpdateUserVerificationCommand{
			UserID: user.ID,

			EmailVerifiedAt: time.Now(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return r, err
		}
	}

	authStore.Revoke(auth.UsageEmailVerification, r.VerificationToken)
	r.Result = cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
	return r, nil
}

func (s *UserService) VerifyPhoneUsingToken(ctx context.Context, r *VerifyPhoneUsingTokenEndpoint) error {
	key := fmt.Sprintf("VerifyPhoneUsingToken %v-%v", r.Context.Token, r.VerificationToken)
	res, err := idempgroup.DoAndWrap(key, 30*time.Second,
		func() (interface{}, error) {
			return s.verifyPhoneUsingToken(ctx, r)
		}, "xác nhận số điện thoại")

	if err != nil {
		return err
	}
	r.Result = res.(*VerifyPhoneUsingTokenEndpoint).Result
	return err
}

func (s *UserService) verifyPhoneUsingToken(ctx context.Context, r *VerifyPhoneUsingTokenEndpoint) (*VerifyPhoneUsingTokenEndpoint, error) {
	if r.VerificationToken == "" {
		return r, cm.Error(cm.InvalidArgument, "Missing code", nil)
	}
	if r.Context.Extra[keyRequestVerifyCode] != "" && r.Context.Extra[keyRequestVerifyCode] != r.VerificationToken {
		r.Result = cmapi.Message("fail", "Mã xác thực không chính xác.")
		return r, nil
	}
	if r.Context.UserID == 0 && r.Context.Extra != nil {
		extra := r.Context.Extra
		extra[keyRequestPhoneVerificationVerified] = "1"
		updateSessionCmd := &tokens.UpdateSessionCommand{
			Token:  r.Context.Claim.Token,
			Values: extra,
		}
		if err := bus.Dispatch(ctx, updateSessionCmd); err != nil {
			return r, err
		}
		r.Result = cmapi.Message("ok", "Số điện thoại đã được xác nhận thành công.")
		return r, nil
	}
	getUserByID := &model.GetUserByIDQuery{
		UserID: r.Context.UserID,
	}
	if err := bus.Dispatch(ctx, getUserByID); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return r, err
	}
	var v map[string]string
	tok, _ := authStore.Validate(auth.UsagePhoneVerification, r.VerificationToken, &v)
	user := getUserByID.Result
	tok, code, v := getToken(auth.UsagePhoneVerification, user.ID)
	if tok == nil || code == "" || v == nil {
		return r, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}
	if v["extra"] != user.Phone {
		authStore.Revoke(auth.UsageSToken, tok.TokenStr)
		return r, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}

	var err error
	if code != r.VerificationToken {
		// Delete token after 3 times
		if len(v["tries"]) >= 2 {
			if err = authStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
				ll.Error("Can not revoke token", l.Error(err))
			}
		} else {
			v["tries"] += "."
			authStore.SetTTL(tok, 60*60)
		}
		return r, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}

	if user.PhoneVerifiedAt.IsZero() {
		cmd := &model.UpdateUserVerificationCommand{
			UserID: user.ID,

			PhoneVerifiedAt: time.Now(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return r, err
		}
	}

	authStore.Revoke(auth.UsagePhoneVerification, tok.TokenStr)
	r.Result = cmapi.Message("ok", "Số điện thoại đã được xác nhận thành công.")
	return r, nil
}

func (s *UserService) UpgradeAccessToken(ctx context.Context, r *UpgradeAccessTokenEndpoint) error {
	key := fmt.Sprintf("UpgradeAccessToken %v-%v", r.Context.User.ID, r.Stoken)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
		func() (interface{}, error) {
			return s.upgradeAccessToken(ctx, r)
		}, "cập nhật thông tin")

	if err != nil {
		return err
	}
	r.Result = res.(*UpgradeAccessTokenEndpoint).Result
	return err
}

func (s *UserService) upgradeAccessToken(ctx context.Context, r *UpgradeAccessTokenEndpoint) (*UpgradeAccessTokenEndpoint, error) {
	if r.Stoken == "" {
		return r, cm.Error(cm.InvalidArgument, "Missing code", nil)
	}

	user := r.Context.User.User
	tok, code, v := getToken(auth.UsageSToken, user.ID)
	if tok == nil || code == "" || v == nil {
		return r, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}
	if v["extra"] != user.Email {
		authStore.Revoke(auth.UsageSToken, tok.TokenStr)
		return r, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}

	var err error
	if code == r.Stoken {
		expiresAt := time.Now().Add(30 * time.Minute)
		r.Result, err = s.CreateSessionResponse(
			ctx,
			&claims.ClaimInfo{
				SToken:          true,
				STokenExpiresAt: &expiresAt,
			},
			"", // Generate new token
			user.ID,
			user,
			r.Context.AccountID,
			0,
			0,
		)
		if err != nil {
			return r, err
		}
		if err = authStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
			ll.Error("Can not revoke token", l.Error(err))
		}
		return r, nil
	}

	// Delete token after 3 times
	if len(v["tries"]) >= 2 {
		if err = authStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
			ll.Error("Can not revoke token", l.Error(err))
		}
	} else {
		v["tries"] += "."
		authStore.SetTTL(tok, 60*60)
	}
	return r, cm.Errorf(cm.PermissionDenied, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
}

func (s *UserService) SendSTokenEmail(ctx context.Context, r *SendSTokenEmailEndpoint) error {
	key := fmt.Sprintf("SendSTokenEmail %v-%v-%v", r.Context.User.ID, r.Email, r.AccountId)
	res, err := idempgroup.DoAndWrap(key, 30*time.Second,
		func() (interface{}, error) {
			return s.sendSTokenEmail(ctx, r)
		}, "gửi email xác nhận")

	if err != nil {
		return err
	}
	r.Result = res.(*SendSTokenEmailEndpoint).Result
	return err
}

func (s *UserService) sendSTokenEmail(ctx context.Context, r *SendSTokenEmailEndpoint) (*SendSTokenEmailEndpoint, error) {
	if !enabledEmail {
		return r, cm.Error(cm.FailedPrecondition, "Không thể gửi email xác nhận. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return r, cm.Error(cm.FailedPrecondition, "Thiếu thông tin email. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if !strings.Contains(r.Email, "@") {
		return r, cm.Error(cm.FailedPrecondition, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	if r.AccountId == 0 {
		return r, cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}

	user := r.Context.User.User
	accQuery := &model.GetAccountRolesQuery{
		AccountID: r.AccountId,
		UserID:    user.ID,
	}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return r, err
	}
	account := accQuery.Result.Account
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	userEmail, _ := validate.NormalizeEmail(user.Email)
	if !ok {
		return r, cm.Error(cm.InvalidArgument, "Email không hợp lệ. Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}
	if emailNorm != userEmail {
		return r, cm.Error(cm.InvalidArgument, "Email không đúng. Vui lòng thử lại hoặc liên hệ hotro@etop.vn.", nil)
	}

	_, code, _, err := generateToken(auth.UsageSToken, user.ID, true, 2*60*60, user.Email)
	if err != nil {
		return r, err
	}

	emailData := map[string]interface{}{
		"FullName":    user.FullName,
		"Code":        code,
		"Email":       user.Email,
		"AccountName": account.Name,
	}
	switch account.Type {
	case account_type.Shop:
		emailData["AccountType"] = model.AccountTypeLabel(account.Type)
	case account_type.Etop:
		return r, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email đến tài khoản %v", account.Name, account).WithMeta("type", string(account.Type))
	default:
		return r, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email đến tài khoản %v", account.Name).WithMeta("type", string(account.Type))
	}

	var b strings.Builder
	if err := emailSTokenTpl.Execute(&b, emailData); err != nil {
		return r, cm.Errorf(cm.Internal, err, "Không thể gửi email đến tài khoản %v", account.Name).WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    "eTop.vn (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận thay đổi thông tin tài khoản",
		Content:     b.String(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}
	r.Result = cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email kèm mã xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", address))
	return r, nil
}

func getToken(usage string, userID dot.ID) (*auth.Token, string, map[string]string) {
	tokStr := userID.String()
	var v map[string]string
	var code string

	tok, err := authStore.Validate(usage, tokStr, &v)
	if err == nil && v != nil && len(v["code"]) == 6 {
		code = v["code"]
		return tok, code, v
	}

	_ = authStore.Revoke(usage, tokStr)
	return nil, "", nil
}

func generateToken(usage string, userID dot.ID, generate bool, ttl int, extra string) (*auth.Token, string, map[string]string, error) {
	tokStr := userID.String()
	tok, code, v := getToken(usage, userID)
	if code != "" {
		return tok, code, v, nil
	}

	code, err := gencode.Random6Digits()
	if err != nil {
		return nil, "", nil, cm.Error(cm.Internal, "", err)
	}

	v = map[string]string{
		"code":  code,
		"tries": "",
		"extra": extra,
	}
	tok = &auth.Token{
		TokenStr: tokStr,
		Usage:    usage,
		UserID:   userID,
		Value:    v,
	}
	if ttl == 0 {
		return nil, "", nil, cm.Error(cm.Internal, "Invalid ttl", nil)
	}
	tok, err = authStore.GenerateWithValue(tok, ttl)
	if err != nil {
		return nil, "", nil, cm.Error(cm.Internal, "", err)
	}
	return tok, code, v, nil
}

func (s *UserService) UpdateReferenceUser(ctx context.Context, r *UpdateReferenceUserEndpoint) error {
	cmd := &identity.UpdateUserReferenceUserIDCommand{
		UserID:       r.Context.UserID,
		RefUserPhone: r.Phone,
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *UserService) UpdateReferenceSale(ctx context.Context, r *UpdateReferenceSaleEndpoint) error {
	cmd := &identity.UpdateUserReferenceSaleIDCommand{
		UserID:       r.Context.UserID,
		RefSalePhone: r.Phone,
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *UserService) InitSession(ctx context.Context, r *InitSessionEndpoint) error {
	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{},
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		return cm.Errorf(cm.Internal, err, "")
	}

	r.Result = generateShopLoginResponse(
		tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn,
	)
	return nil
}

func generateShopLoginResponse(accessToken string, expiresIn int) *etop.LoginResponse {
	resp := &etop.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}
	return resp
}

func getTokenValid(ctx context.Context, token string) string {
	if token != "" {
		return token
	}
	return getBearerTokenFromCtx(ctx)
}

func getBearerTokenFromCtx(ctx context.Context) string {
	authHeader, ok := ctx.Value(authKey{}).(string)
	if !ok {
		return ""
	}
	token, _ := auth.FromHeaderString(authHeader)
	return token
}

func sendPhoneVerificationForRegister(ctx context.Context, r *SendPhoneVerificationEndpoint) (*SendPhoneVerificationEndpoint, error) {
	token := getTokenValid(ctx, r.Context.Token)
	claim, err := tokens.Store.Validate(token)

	if err != nil {
		return r, nil
	}
	phoneNorm, ok := validate.NormalizePhone(r.Phone)

	if !ok {
		return r, cm.Error(cm.FailedPrecondition, "Số điện thoại không đúng. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
	}
	_, code, _, err := generateToken(auth.UsagePhoneVerification, cm.NewID(), true, 2*60*60, r.Phone)
	if err != nil {
		return r, err
	}
	msg := fmt.Sprintf(smsVerificationTpl, code, r.Phone)
	phone := r.Phone
	cmd := &sms.SendSMSCommand{
		Phone:   phone,
		Content: msg,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return r, err
	}
	if claim.Extra == nil || claim.Extra[keyRequestPhoneVerificationVerified] != "" {
		claim.Extra = map[string]string{}
	}
	claim.Extra[keyRequestVerifyPhone] = phoneNorm.String()
	claim.Extra[keyRequestVerifyCode] = code

	updateSessionCmd := &tokens.UpdateSessionCommand{
		Token:  r.Context.Claim.Token,
		Values: claim.Extra,
	}
	if err := bus.Dispatch(ctx, updateSessionCmd); err != nil {
		return r, err
	}
	r.Result = cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.", phone))
	return r, nil
}
