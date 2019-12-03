package api

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	pbcm "etop.vn/api/pb/common"
	pbetop "etop.vn/api/pb/etop"
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
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/logic/usering"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	eventBus            capi.EventBus
	idempgroup          *idemp.RedisGroup
	authStore           auth.Generator
	enabledEmail        bool
	enabledSMS          bool
	cfgEmail            EmailConfig
	identityAggr        identity.CommandBus
	identityQS          identity.QueryBus
	invitationAggregate invitation.CommandBus
	invitationQuery     invitation.QueryBus
)

const PrefixIdempUser = "IdempUser"

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
	)
}

type UserService struct{}

var userService = &UserService{}

type EmailConfig = config.EmailConfig

func Init(
	bus capi.EventBus,
	identityCommandBus identity.CommandBus,
	identityQueryBus identity.QueryBus,
	invitationAggr invitation.CommandBus,
	invitationQS invitation.QueryBus,
	sd cmservice.Shutdowner,
	rd redis.Store,
	s auth.Generator,
	_cfgEmail EmailConfig,
	_cfgSMS sms.Config,
) {
	eventBus = bus
	identityAggr = identityCommandBus
	identityQS = identityQueryBus
	invitationAggregate = invitationAggr
	invitationQuery = invitationQS
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

// Register
// 1a. If the user does not exist in the system, create a new user.
// 1b. If any email or phone is activated -> AlreadyExists.
//   - If both email and phone exist (but not activated) -> Merge them.
//   - Otherwise, update existing user with the other identifier.
func (s *UserService) Register(ctx context.Context, r *RegisterEndpoint) error {
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

	now := time.Now()
	var phoneNorm model.NormalizedPhone
	var emailNorm model.NormalizedEmail
	var ok bool
	phoneNorm, ok = validate.NormalizePhone(r.Phone)
	if !ok {
		return cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}
	if r.Email != "" {
		emailNorm, ok = validate.NormalizeEmail(r.Email)
		if !ok {
			return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
	}

	if r.RegisterToken != "" {
		getInvitationByToken := &invitation.GetInvitationByTokenQuery{
			Token:  r.RegisterToken,
			Result: nil,
		}
		if err := invitationQuery.Dispatch(ctx, getInvitationByToken); err != nil {
			return cm.MapError(err).
				Wrap(cm.NotFound, "Token không hợp lệ").
				Throw()
		}
		invitation := getInvitationByToken.Result
		if r.Email != invitation.Email {
			return cm.Errorf(cm.InvalidArgument, nil, "Email gửi lên và email trong token không khớp nhau")
		}
	}

	userByPhone, userByEmail, err := s.getUserByPhoneAndByEmail(ctx, phoneNorm.String(), emailNorm.String())
	if err != nil {
		return err
	}

	// If no identifier exists, create new user
	if userByPhone.User == nil && userByEmail.User == nil {
		info := r.CreateUserRequest
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
			Source:         convertpb.UserSourceToModel(&r.Source),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}

		if err := s.publishUserCreatedEvent(ctx, cmd.Result.User); err != nil {
			return err
		}

		r.Result = &pbetop.RegisterResponse{
			User: convertpb.PbUser(cmd.Result.User),
		}
		return nil
	}

	// If any identifier is activated -> AlreadyExists
	if userByPhone.User != nil && userByPhone.User.Status != model.S3Zero {
		return cm.Error(cm.FailedPrecondition, sqlstore.MsgCreateUserDuplicatedPhone, nil)
	}
	if userByEmail.User != nil && userByEmail.User.Status != model.S3Zero {
		return cm.Error(cm.FailedPrecondition, sqlstore.MsgCreateUserDuplicatedEmail, nil)
	}

	// If both one identifier exists, merge them and activate the user
	if userByPhone.User != nil && userByEmail.User != nil {
		resp, err := s.mergeUserByPhoneAndEmail(ctx, userByPhone, userByEmail)
		r.Result = resp
		return err
	}

	// Lastly, update one identifier with the other and activate the user
	updateUserCmd := &model.UpdateUserIdentifierCommand{}
	var user *model.User
	var inner model.UserInner
	var hashpwd string
	switch {
	case userByPhone.User != nil && userByPhone.User.Status == model.S3Zero:
		user = userByPhone.User
		if user.Email != "" {
			// unexpected: if both identifier exist, it must not be a stub!
			return cm.Error(cm.Internal, "", nil).
				Log("unexpected condition (1)", l.ID("user_id", user.ID))
		}
		hashpwd = userByPhone.UserInternal.Hashpwd

		inner = user.UserInner
		inner.Email = emailNorm.String()
		updateUserCmd.UserID = user.ID
		updateUserCmd.Status = model.S3Positive
		updateUserCmd.CreatedAt = now
		updateUserCmd.Identifying = "full"
		if emailNorm == "" {
			updateUserCmd.Identifying = "half"
		}

	case userByEmail.User != nil && userByEmail.User.Status == model.S3Zero:
		user = userByEmail.User
		if user.Phone != "" {
			// unexpected: if both identifier exist, it must not be a stub!
			return cm.Error(cm.Internal, "", nil).
				Log("unexpected condition (2)", l.ID("user_id", user.ID))
		}
		hashpwd = userByPhone.UserInternal.Hashpwd

		inner = user.UserInner
		inner.Phone = phoneNorm.String()
		updateUserCmd.UserID = user.ID
		updateUserCmd.Status = model.S3Positive
		updateUserCmd.CreatedAt = now
		updateUserCmd.Identifying = "full"

	default:
		return cm.Error(cm.Internal, "", nil).
			Log("unexpected condition (3)")
	}

	inner.FullName, inner.ShortName, err = sqlstore.NormalizeFullName(r.FullName, r.ShortName)
	if err != nil {
		return err
	}
	updateUserCmd.UserInner = inner
	updateUserCmd.Password = r.Password

	// Automatically verify phone number
	if r.RegisterToken != "" && hashpwd != "" {
		if !login.VerifyPassword(r.RegisterToken, hashpwd) {
			return cm.Error(cm.Internal, "Mã số không hợp lệ. Vui lòng liên hệ hotro@etop.vn.", nil)
		}
		if user.PhoneVerifiedAt.IsZero() &&
			!user.PhoneVerificationSentAt.IsZero() &&
			now.Add(-24*time.Hour).Before(user.PhoneVerificationSentAt) {
			updateUserCmd.PhoneVerifiedAt = now
		}
	}

	if err := bus.Dispatch(ctx, updateUserCmd); err != nil {
		return cm.Error(cm.Internal, "", err)
	}

	if err := s.publishUserCreatedEvent(ctx, updateUserCmd.Result.User); err != nil {
		return err
	}

	r.Result = createRegisterResponse(updateUserCmd.Result.User)
	return nil
}

func (s *UserService) publishUserCreatedEvent(ctx context.Context, user *model.User) error {
	event := &identity.UserCreatedEvent{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		ShortName: user.ShortName,
	}
	if err := eventBus.Publish(ctx, event); err != nil {
		return err
	}
	return nil
}

func createRegisterResponse(user *model.User) *pbetop.RegisterResponse {
	return &pbetop.RegisterResponse{
		User: convertpb.PbUser(user),
	}
}

func (s *UserService) CheckUserRegistration(ctx context.Context, q *CheckUserRegistrationEndpoint) error {
	_, ok := validate.NormalizePhone(q.Phone)
	if !ok {
		q.Result = &pbetop.GetUserByPhoneResponse{Exists: false}
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
		q.Result = &pbetop.GetUserByPhoneResponse{Exists: false}
		return nil
	}
	q.Result = &pbetop.GetUserByPhoneResponse{Exists: true}
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

func (s *UserService) mergeUserByPhoneAndEmail(ctx context.Context, userByPhone, userByEmail model.UserExtended) (*pbetop.RegisterResponse, error) {
	return nil, cm.ErrTODO
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

func (s *UserService) CreateSessionResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, adminID dot.ID) (*pbetop.AccessTokenResponse, error) {
	resp, err := s.CreateLoginResponse(ctx, claim, token, userID, user, preferAccountID, preferAccountType, false, adminID)
	if err != nil {
		return nil, err
	}
	return &pbetop.AccessTokenResponse{
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

func (s *UserService) CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (*pbetop.LoginResponse, error) {
	resp, _, err := s.CreateLoginResponse2(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
	return resp, err
}

func (s *UserService) CreateLoginResponse2(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *model.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (_ *pbetop.LoginResponse, respShop *model.Shop, _ error) {

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

	var currentAccount *pbetop.LoginAccount
	var currentAccountID dot.ID
	availableAccounts := make([]*pbetop.LoginAccount, len(accQuery.Result))
	for i, accUserX := range accQuery.Result {
		availableAccounts[i] = convertpb.PbLoginAccount(accUserX)
		account := accUserX.Account
		switch {
		case preferAccountID == account.ID,
			preferAccountType == model.TagShop && account.Type == model.TypeShop,
			preferAccountType == model.TagEtop && account.Type == model.TypeEtop,
			preferAccountType == model.TagAffiliate && account.Type == model.TypeAffiliate:
			currentAccount = availableAccounts[i]
			currentAccountID = currentAccount.Id
		}
	}

	resp := &pbetop.LoginResponse{
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
			acc.ExpiresIn = int(tokenCmd.Result.ExpiresIn)
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
	key := fmt.Sprintf("SendPhoneVerification %v-%v", r.Context.User.ID, r.Phone)
	res, err := idempgroup.DoAndWrap(key, 60*time.Second,
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
	user := r.Context.User.User
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
	key := fmt.Sprintf("VerifyPhoneUsingToken %v-%v", r.Context.User.ID, r.VerificationToken)
	res, err := idempgroup.DoAndWrap(key, 15*time.Second,
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

	user := r.Context.User
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
	res, err := idempgroup.DoAndWrap(key, 60*time.Second,
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
	case model.TypeShop:
		emailData["AccountType"] = account.Type.Label()
	case model.TypeEtop:
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
