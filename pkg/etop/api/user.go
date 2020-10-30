package api

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/api/main/notify"
	"o.o/api/top/int/etop"
	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/authentication_method"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/apifw/whitelabel/templatemessages"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/code/gencode"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	authservice "o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/sms"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/timex"
)

var (
	idempgroup   *idemp.RedisGroup
	enabledEmail bool
	enabledSMS   bool
	cfgEmail     cc.EmailConfig
)

const PrefixIdempUser = "IdempUser"

type SignalUpdate string

const (
	keyRequestVerifyCode                string       = "request_phone_verification_code"
	keyRequestVerifyPhone               string       = "request_phone_verification"
	keyRequestEmailVerifyCode           string       = "request_email_verification_code"
	keyRequestPhoneVerificationVerified string       = "request_phone_verification_verified"
	keyRequestAuthUsage                 string       = "request_auth_usage"
	keyRedisFirstCodeUpdateUser         string       = "update-first-code"
	keyRedisSecondCodeUpdateUser        string       = "update-second-code"
	signalUpdateUserEmail               SignalUpdate = "update-email"
	signalUpdateUserPhone               SignalUpdate = "update-phone"
	keyRequestRegisterSimplify          string       = "register-simplify"
)

type UserService struct {
	session.Session

	IdentityAggr     identity.CommandBus
	IdentityQuery    identity.QueryBus
	InvitationQuery  invitation.QueryBus
	NotifyQuery      notify.QueryBus
	NotifyAggr       notify.CommandBus
	EventBus         capi.EventBus
	AuthStore        auth.Generator
	TokenStore       tokens.TokenStore
	RedisStore       redis.Store
	SMSClient        *sms.Client
	EmailClient      *email.Client
	UserStore        sqlstore.UserStoreFactory
	UserStoreIface   sqlstore.UserStoreInterface
	ShopStore        sqlstore.ShopStoreInterface
	AccountUserStore sqlstore.AccountUserStoreInterface
	LoginIface       sqlstore.LoginInterface
}

var UserServiceImpl = &UserService{} // MUSTDO: fix it

func (s *UserService) Clone() api.UserService {
	res := *s
	return &res
}

func (s *UserService) UpdatePermission(ctx context.Context, q *api.UpdatePermissionRequest) (*api.UpdatePermissionResponse, error) {
	return nil, cm.ErrTODO
}

func (s *UserService) UpdateUserEmail(ctx context.Context, r *api.UpdateUserEmailRequest) (*api.UpdateUserEmailResponse, error) {
	key := fmt.Sprintf("UpdateUserEmail %v-%v-%v-%v", r.Email, r.FirstCode, r.SecondCode, r.AuthenticationMethod)
	result, _, err := idempgroup.DoAndWrap(
		ctx, key, 60*time.Second, "thay đổi email",
		func() (interface{}, error) { return s.updateUserEmail(ctx, r) })

	if err != nil {
		return nil, err
	}
	return result.(*api.UpdateUserEmailResponse), nil
}

func (s *UserService) UpdateUserPhone(ctx context.Context, r *api.UpdateUserPhoneRequest) (*api.UpdateUserPhoneResponse, error) {
	key := fmt.Sprintf("UpdateUserPhone %v-%v-%v-%v", r.Phone, r.FirstCode, r.SecondCode, r.AuthenticationMethod)
	result, _, err := idempgroup.DoAndWrap(
		ctx, key, 60*time.Second, "thay đổi số điện thoại",
		func() (interface{}, error) { return s.updateUserPhone(ctx, r) })
	if err != nil {
		return nil, err
	}
	return result.(*api.UpdateUserPhoneResponse), nil
}

func (s *UserService) updateUserPhone(ctx context.Context, r *api.UpdateUserPhoneRequest) (*api.UpdateUserPhoneResponse, error) {
	code, count, err := s.getRedisCode(s.SS.User().ID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserPhone)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	var failCount int
	failCount, err = s.checkFailCount(s.SS.Claim().UserID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserPhone)
	if err != nil {
		return nil, err
	}
	user, err := s.checkUserInfo(ctx, s.SS.User().ID, r.AuthenticationMethod)
	if err != nil {
		return nil, err
	}
	if r.AuthenticationMethod == authentication_method.Phone && timex.IsZeroTime(user.PhoneVerifiedAt) {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Số điện thoại chưa được xác thực.")
	}
	if r.AuthenticationMethod == authentication_method.Email && timex.IsZeroTime(user.EmailVerifiedAt) {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Email chưa được xác thực.")
	}
	switch r.FirstCode {
	case "":
		if r.AuthenticationMethod == authentication_method.Phone {
			msg, err := s.sendPhoneUserCode(ctx, user, user.Phone, keyRedisFirstCodeUpdateUser, signalUpdateUserPhone, code, count, r.AuthenticationMethod)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserPhoneResponse{Msg: msg}
			return result, nil
		}
		if r.AuthenticationMethod == authentication_method.Email {
			msg, err := s.sendEmailUserCode(ctx, user, user.Email, keyRedisFirstCodeUpdateUser, signalUpdateUserPhone, code, count, r.AuthenticationMethod)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserPhoneResponse{Msg: msg}
			return result, nil
		}
		return nil, cm.Errorf(cm.STokenRequired, nil, "Cần chọn phương thức xác nhận. Vui lòng chọn email hoặc số điện thoại.")

	case code:
		return s.updatePhoneVerifySecondCode(ctx, r, user)

	default:
		failCount++
		err = s.setFailCountToRedis(s.SS.User().ID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserPhone, failCount)
		if err != nil {
			return nil, err
		}
		return nil, cm.Errorf(cm.STokenRequired, nil, "Mã xác thực không tồn tại vui lòng thử lại.")
	}
}

func (s *UserService) updatePhoneVerifySecondCode(ctx context.Context, r *api.UpdateUserPhoneRequest, user *identitymodel.User) (*api.UpdateUserPhoneResponse, error) {
	if r.Phone == user.Phone {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại mới đang trùng với số điện thoại hiện tại.")
	}
	userByPhoneQuery := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: r.Phone,
	}
	err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, userByPhoneQuery)
	if err == nil || cm.ErrorCode(err) != cm.NotFound {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại đã tồn tại vui lòng kiểm tra lại.")
	}
	var codeSecond string
	codeSecond, count, err := s.getRedisCode(user.ID, keyRedisSecondCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserPhone)
	switch err {
	case redis.ErrNil:
		msg, err2 := s.sendPhoneUserCode(ctx, user, r.Phone, keyRedisSecondCodeUpdateUser, signalUpdateUserPhone, codeSecond, count, r.AuthenticationMethod)
		if err2 != nil {
			return nil, err2
		}
		result := &api.UpdateUserPhoneResponse{Msg: msg}
		return result, nil
	case nil:
		// continue
	default:
		return nil, err
	}

	switch r.SecondCode {
	case "":
		msg, err2 := s.sendPhoneUserCode(ctx, user, r.Phone, keyRedisSecondCodeUpdateUser, signalUpdateUserPhone, codeSecond, count, r.AuthenticationMethod)
		if err2 != nil {
			return nil, err2
		}
		result := &api.UpdateUserPhoneResponse{Msg: msg}
		return result, nil

	case codeSecond:
		cmd := &identity.UpdateUserPhoneCommand{
			UserID: user.ID,
			Phone:  r.Phone,
		}
		err = s.IdentityAggr.Dispatch(ctx, cmd)
		if err != nil {
			return nil, err
		}
		err = s.clearRedisUpdateUser(user.ID, r.AuthenticationMethod, signalUpdateUserPhone)
		if err != nil {
			return nil, err
		}
		result := &api.UpdateUserPhoneResponse{
			Msg: "Cập nhật số điện thoại thành công",
		}
		ll.SendMessage(fmt.Sprintf("–– User: %v (%v) \n Update: thay đổi số điện thoại từ %v thành %v", user.FullName, user.ID, user.Phone, r.Phone))
		return result, nil

	default:
		return nil, cm.Errorf(cm.STokenRequired, nil, "Mã xác thực không tồn tại vui lòng thử lại.")
	}
}

func (s *UserService) updateUserEmail(ctx context.Context, r *api.UpdateUserEmailRequest) (*api.UpdateUserEmailResponse, error) {
	code, count, err := s.getRedisCode(s.SS.User().ID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserEmail)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	var failCount int
	failCount, err = s.checkFailCount(s.SS.Claim().UserID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserEmail)
	if err != nil {
		return nil, err
	}
	user, err := s.checkUserInfo(ctx, s.SS.User().ID, r.AuthenticationMethod)
	if err != nil {
		return nil, err
	}
	switch r.FirstCode {
	case "":
		if r.AuthenticationMethod == authentication_method.Phone {
			msg, err := s.sendPhoneUserCode(ctx, user, user.Phone, keyRedisFirstCodeUpdateUser, signalUpdateUserEmail, code, count, r.AuthenticationMethod)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserEmailResponse{Msg: msg}
			return result, nil
		}
		if r.AuthenticationMethod == authentication_method.Email {
			msg, err := s.sendEmailUserCode(ctx, user, user.Email, keyRedisFirstCodeUpdateUser, signalUpdateUserEmail, code, count, r.AuthenticationMethod)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserEmailResponse{Msg: msg}
			return result, nil
		}
		return nil, cm.Errorf(cm.STokenRequired, nil, "Cần chọn phương thức xác nhận. Vui lòng chọn email hoặc số điện thoại.")

	case code:
		return s.updateEmailVerifySecondCode(ctx, r, user)

	default:
		failCount++
		err = s.setFailCountToRedis(s.SS.User().ID, keyRedisFirstCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserEmail, failCount)
		if err != nil {
			return nil, err
		}
		return nil, cm.Errorf(cm.STokenRequired, nil, "Mã xác thực không tồn tại vui lòng thử lại.")
	}
}

func (s *UserService) getRedisCode(userID dot.ID, keyCode string, method authentication_method.AuthenticationMethod, signal SignalUpdate) (string, int, error) {
	var code string
	var count int
	var err error
	err = s.RedisStore.Get(fmt.Sprintf("Code:%v-%v-%v-%v", userID, keyCode, method.String(), signal), &code)
	if err != nil {
		return "", 0, err
	}
	err = s.RedisStore.Get(fmt.Sprintf("Code-Count:%v-%v-%v-%v", userID, keyCode, method.String(), signal), &count)
	if err != nil && err != redis.ErrNil {
		return "", 0, err
	}
	return code, count, nil
}

func (s *UserService) setRedisCode(userID dot.ID, redisKeyCode string, method string, code string, count int, signal SignalUpdate) error {
	err := s.RedisStore.SetWithTTL(fmt.Sprintf("Code:%v-%v-%v-%v", userID, redisKeyCode, method, signal), code, 2*60*60)
	if err != nil {
		return err
	}
	err = s.RedisStore.SetWithTTL(fmt.Sprintf("Code-Count:%v-%v-%v-%v", userID, redisKeyCode, method, signal), count, 2*60*60)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) checkUserInfo(ctx context.Context, userID dot.ID, method authentication_method.AuthenticationMethod) (*identitymodel.User, error) {
	user, err := s.UserStore(ctx).ID(userID).Get()
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "User không tồn tại")
	}
	if method == authentication_method.Phone && timex.IsZeroTime(user.PhoneVerifiedAt) {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Số điện thoại chưa được xác thực.")
	}
	if method == authentication_method.Email && timex.IsZeroTime(user.EmailVerifiedAt) {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Email chưa được xác thực.")
	}
	return user, nil
}

func (s *UserService) setFailCountToRedis(userID dot.ID, keyRedisCode string, method authentication_method.AuthenticationMethod, signal SignalUpdate, failCount int) error {
	if failCount >= 3 {
		err := s.clearRedisUpdateUser(userID, method, signal)
		if err != nil {
			return err
		}
	}
	err := s.RedisStore.SetWithTTL(fmt.Sprintf("UpdateUserFailCount-%v-%v-%v-%v", keyRedisCode, userID, method.String(), signal), failCount, 30*60)
	return err
}

func (s *UserService) checkFailCount(userID dot.ID, keyRediscode string, method authentication_method.AuthenticationMethod, signal SignalUpdate) (int, error) {
	var failCount int
	redisKey := fmt.Sprintf("UpdateUserFailCount-%v-%v-%v-%v", keyRediscode, userID, method.String(), signal)
	err := s.RedisStore.Get(redisKey, &failCount)
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	if err == redis.ErrNil {
		failCount = 0
	}
	ttl, err := s.RedisStore.GetTTL(redisKey)
	minutes := (ttl - 1) / 60
	minutes = minutes / 10
	if minutes > 0 {
		minutes = minutes*10 + 10
	} else {
		minutes = 5
	}
	if failCount >= 3 {
		return 0, cm.Errorf(cm.FailedPrecondition, err, fmt.Sprintf("Mã sai quá 3 lần. Vui lòng thử lại sau %v phút", minutes))
	}
	return failCount, nil
}

func (s *UserService) clearRedisUpdateUser(userID dot.ID, method authentication_method.AuthenticationMethod, signal SignalUpdate) error {
	err := s.RedisStore.Del(fmt.Sprintf("UpdateUserFailCount-%v-%v-%v", keyRedisFirstCodeUpdateUser, userID, method.String()))
	if err != nil && err != redis.ErrNil {
		return err
	}
	err = s.RedisStore.Del(fmt.Sprintf("Code:%v-%v-%v-%v", userID, keyRedisFirstCodeUpdateUser, method.String(), signal))
	if err != nil && err != redis.ErrNil {
		return err
	}
	err = s.RedisStore.Del(fmt.Sprintf("Code-Count:%v-%v-%v-%v", userID, keyRedisFirstCodeUpdateUser, method.String(), signal))
	if err != nil && err != redis.ErrNil {
		return err
	}
	err = s.RedisStore.Del(fmt.Sprintf("Code:%v-%v-%v-%v", userID, keyRedisSecondCodeUpdateUser, method.String(), signal))
	if err != nil && err != redis.ErrNil {
		return err
	}
	return nil
}

func (s *UserService) updateEmailVerifySecondCode(ctx context.Context, r *api.UpdateUserEmailRequest, user *identitymodel.User) (*api.UpdateUserEmailResponse, error) {
	normalizeEmail, _ := validate.NormalizeEmail(r.Email)
	r.Email = normalizeEmail.String()
	if r.Email == user.Email {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Email mới đang trùng với email hiện tại.")
	}
	userByEmailQuery := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Email: r.Email,
	}
	err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, userByEmailQuery)
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	if err == nil || userByEmailQuery.Result.ID != 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Email đã tồn tại vui lòng kiểm tra lại.")
	}
	codeSecond, count, err := s.getRedisCode(s.SS.User().ID, keyRedisSecondCodeUpdateUser, r.AuthenticationMethod, signalUpdateUserEmail)
	switch err {
	case redis.ErrNil:
		msg, err := s.sendEmailUserCode(ctx, user, r.Email, keyRedisSecondCodeUpdateUser, signalUpdateUserEmail, codeSecond, count, r.AuthenticationMethod)
		if err != nil {
			return nil, err
		}
		result := &api.UpdateUserEmailResponse{Msg: msg}
		return result, nil
	case nil:
		switch r.SecondCode {
		case "":
			msg, err := s.sendEmailUserCode(ctx, user, r.Email, keyRedisSecondCodeUpdateUser, signalUpdateUserEmail, codeSecond, count, r.AuthenticationMethod)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserEmailResponse{Msg: msg}
			return result, nil
		case codeSecond:
			cmd := &identity.UpdateUserEmailCommand{
				Email:  r.Email,
				UserID: user.ID,
			}
			err = s.IdentityAggr.Dispatch(ctx, cmd)
			if err != nil {
				return nil, err
			}
			err = s.clearRedisUpdateUser(user.ID, r.AuthenticationMethod, signalUpdateUserEmail)
			if err != nil {
				return nil, err
			}
			result := &api.UpdateUserEmailResponse{
				Msg: "Cập nhật email thành công",
			}
			ll.SendMessage(fmt.Sprintf("–– User: %v (%v) \n Update: thay đổi email từ %v thành %v", user.FullName, user.ID, user.Email, r.Email))
			return result, nil
		default:
			return nil, cm.Errorf(cm.STokenRequired, nil, "Mã xác thực không tồn tại vui lòng thử lại.")
		}
	default:
		return nil, err
	}
}

func (s *UserService) sendPhoneUserCode(ctx context.Context, user *identitymodel.User, phone string, redisCode string, signal SignalUpdate, code6Digits string, sendCount int, method authentication_method.AuthenticationMethod) (string, error) {
	phoneUse := phone
	var code string
	var err error
	if code6Digits == "" {
		code, err = gencode.Random6Digits()
		if err != nil {
			return "", err
		}
	} else {
		code = code6Digits
	}
	sendCount++
	var msgUser string
	if redisCode == keyRedisFirstCodeUpdateUser {
		if signal == signalUpdateUserEmail {
			if sendCount > 1 {
				msgUser = fmt.Sprintf(templatemessages.SmsChangeEmailTplRepeat, code, wl.X(ctx).Name, sendCount)
			} else {
				msgUser = fmt.Sprintf(templatemessages.SmsChangeEmailTpl, code, wl.X(ctx).Name)
			}
		} else {
			if sendCount > 1 {
				msgUser = fmt.Sprintf(templatemessages.SmsChangePhoneTplRepeat, code, wl.X(ctx).Name, sendCount)
			} else {
				msgUser = fmt.Sprintf(templatemessages.SmsChangePhoneTpl, code, wl.X(ctx).Name)
			}
		}
	} else {
		if sendCount > 1 {
			msgUser = fmt.Sprintf(templatemessages.SmsChangePhoneTplConfirmRepeat, code, wl.X(ctx).Name, sendCount)
		} else {
			msgUser = fmt.Sprintf(templatemessages.SmsChangePhoneTplConfirm, code, wl.X(ctx).Name)
		}
	}

	err = s.setRedisCode(user.ID, redisCode, method.String(), code, sendCount, signal)
	if err != nil {
		return "", err
	}
	cmd := &sms.SendSMSCommand{
		Phone:   phoneUse,
		Content: msgUser,
	}
	if err = s.SMSClient.SendSMS(ctx, cmd); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ %v.", phoneUse, wl.X(ctx).CSEmail), nil
}

func (s *UserService) sendEmailUserCode(ctx context.Context, user *identitymodel.User, emailVerify string, redisKeyCode string, signal SignalUpdate, code6Digits string, sendCount int, method authentication_method.AuthenticationMethod) (string, error) {
	if !enabledEmail {
		return "", cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email xác nhận. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if !validate.IsEmail(emailVerify) {
		return "", cm.Errorf(cm.FailedPrecondition, nil, "Email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	var code string
	var err error
	code = code6Digits
	if code6Digits == "" {
		code, err = gencode.Random6Digits()
		if err != nil {
			return "", err
		}
	}
	emailData := map[string]interface{}{
		"FullName": user.FullName,
		"Code":     code,
		"Email":    emailVerify,
		"WlName":   wl.X(ctx).Name,
	}
	var b strings.Builder

	if redisKeyCode == keyRedisFirstCodeUpdateUser {
		if signal == signalUpdateUserEmail {
			if err = templatemessages.UpdateEmailTpl.Execute(&b, emailData); err != nil {
				return "", cm.Errorf(cm.Internal, err, "Không thể gửi email đến tài khoản %v", user.FullName).WithMeta("reason", "can not generate email content")
			}
		} else {
			if err = templatemessages.UpdatePhoneTpl.Execute(&b, emailData); err != nil {
				return "", cm.Errorf(cm.Internal, err, "Không thể gửi email đến tài khoản %v", user.FullName).WithMeta("reason", "can not generate email content")
			}
		}
	} else {
		if err = templatemessages.UpdateEmailTplConfirm.Execute(&b, emailData); err != nil {
			return "", cm.Errorf(cm.Internal, err, "Không thể gửi email đến tài khoản %v", user.FullName).WithMeta("reason", "can not generate email content")
		}
	}
	err = s.setRedisCode(user.ID, redisKeyCode, method.String(), code, sendCount, signal)
	if err != nil {
		return "", err
	}
	address := emailVerify
	cmd := &email.SendEmailCommand{
		FromName:    wl.X(ctx).Name + " (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận thay đổi thông tin tài khoản",
		Content:     b.String(),
	}
	if err = s.EmailClient.SendMail(ctx, cmd); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Đã gửi email kèm mã xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ %v.", address, wl.X(ctx).CSEmail), nil
}

// Register
// 1a. If the user does not exist in the system, create a new user.
// 1b. If any email or phone is activated -> AlreadyExists.
//   - If both email and phone exist (but not activated) -> Merge them.
//   - Otherwise, update existing user with the other identifier.
func (s *UserService) Register(ctx context.Context, r *api.CreateUserRequest) (*api.RegisterResponse, error) {
	if err := validateRegister(ctx, r); err != nil {
		return nil, err
	}
	claim := s.SS.Claim()
	if claim.Extra != nil {
		if claim.Extra[keyRequestVerifyPhone] != "" || claim.Extra[keyRequestPhoneVerificationVerified] != "" {
			phoneNorm, _ := validate.NormalizePhone(r.Phone)
			if claim.Extra[keyRequestVerifyPhone] != phoneNorm.String() {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại %v không hợp lệ bởi vì bạn đã xác nhận số điện thoại: %v", phoneNorm.String(), claim.Extra[keyRequestVerifyPhone])
			}

			resp, err := s.register(ctx, claim.Extra, r, true)
			s.AuthStore.Revoke(auth.UsageAccessToken, claim.Token)
			return resp, err
		}
	}
	resp, err := s.register(ctx, nil, r, false)
	return resp, err
}

func (s *UserService) RegisterUsingToken(ctx context.Context, r *api.CreateUserRequest) (*api.RegisterResponse, error) {
	claimExtra := s.SS.Claim().Extra
	if claimExtra[keyRequestVerifyPhone] == "" || claimExtra[keyRequestPhoneVerificationVerified] == "" {
		return nil, cm.Error(cm.InvalidArgument, "Bạn vui lòng xác nhận só điện thoại trước khi đăng kí", nil)
	}
	if err := validateRegister(ctx, r); err != nil {
		return nil, err
	}

	phoneNorm, _ := validate.NormalizePhone(r.Phone)
	if claimExtra[keyRequestVerifyPhone] != phoneNorm.String() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại %v không hợp lệ bởi vì bạn đã xác nhận số điện thoại: %v", phoneNorm.String(), claimExtra[keyRequestVerifyPhone])
	}

	resp, err := s.register(ctx, claimExtra, r, true)
	return resp, err
}

func (s *UserService) register(
	ctx context.Context,
	claimExtra map[string]string,
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
	if strings.HasPrefix(r.RegisterToken, auth.UsageInviteUser+":") {
		invitationTemp, err = s.getInvitation(ctx, r)
		if err != nil {
			return nil, err
		}
	}

	user, err := s.createUser(ctx, r)
	if err != nil {
		return nil, err
	}
	{
		now := time.Now()

		// auto verify email when accept invitation by email
		if invitationTemp != nil && invitationTemp.Email != "" && user.EmailVerifiedAt.IsZero() {
			updateCmd := &identitymodelx.UpdateUserVerificationCommand{
				UserID: user.ID,
			}
			updateCmd.EmailVerifiedAt = now
			if err := s.UserStoreIface.UpdateUserVerification(ctx, updateCmd); err != nil {
				return nil, err
			}
		}

		shouldUpdate := false

		if user.PhoneVerifiedAt.IsZero() {
			if usingToken {
				if invitationTemp == nil ||
					(invitationTemp.Phone == "" || invitationTemp.Phone == claimExtra[keyRequestVerifyPhone]) {
					shouldUpdate = true
				}
			} else {
				if invitationTemp != nil && invitationTemp.Phone != "" {
					shouldUpdate = true
				}
			}
		}

		if shouldUpdate {
			updateCmd := &identitymodelx.UpdateUserVerificationCommand{
				UserID: user.ID,
			}
			updateCmd.PhoneVerifiedAt = now
			if err := s.UserStoreIface.UpdateUserVerification(ctx, updateCmd); err != nil {
				return nil, err
			}
		}
	}
	{
		event := &identity.UserCreatedEvent{
			UserID: user.ID,
		}
		if invitationTemp != nil {
			event.Invitation = &identity.UserInvitation{
				Token:      r.RegisterToken,
				AutoAccept: r.AutoAcceptInvitation,
			}
		}
		if err := s.EventBus.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	return &etop.RegisterResponse{
		User: convertpb.PbUser(user),
	}, nil
}

func (s *UserService) createUser(ctx context.Context, r *etop.CreateUserRequest) (*identitymodel.User, error) {
	info := r
	cmd := &identitymodelx.CreateUserCommand{
		UserInner: identitymodel.UserInner{
			FullName:  info.FullName,
			ShortName: info.ShortName,
			Phone:     info.Phone,
			Email:     info.Email,
		},
		Password:       info.Password,
		Status:         status3.P,
		AgreeTOS:       r.AgreeTos,
		AgreeEmailInfo: r.AgreeEmailInfo.Bool,
		Source:         r.Source,
		RefSale:        r.RefSale,
		RefAff:         r.RefAff,
	}
	if err := s.UserStoreIface.CreateUser(ctx, cmd); err != nil {
		return nil, err
	}
	query := &identitymodelx.GetUserByIDQuery{
		UserID: cmd.Result.User.ID,
	}
	if err := s.UserStoreIface.GetUserByID(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func validateRegister(ctx context.Context, r *etop.CreateUserRequest) error {
	if !r.AgreeTos {
		return cm.Errorf(cm.InvalidArgument, nil, "Bạn cần đồng ý với điều khoản sử dụng dịch vụ để tiếp tục. Nếu cần thêm thông tin, vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
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
		emailNorm, ok := validate.NormalizeEmail(r.Email)
		if !ok {
			return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
		if err := validate.PopularEmailAddressMistake(ctx, emailNorm.String()); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) getInvitation(ctx context.Context, r *etop.CreateUserRequest) (*invitation.Invitation, error) {
	getInvitationByToken := &invitation.GetInvitationByTokenQuery{
		Token:  r.RegisterToken,
		Result: nil,
	}
	if err := s.InvitationQuery.Dispatch(ctx, getInvitationByToken); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Token không hợp lệ").
			Throw()
	}
	invitationTemp := getInvitationByToken.Result
	if invitationTemp.Email != "" && r.Email != invitationTemp.Email {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email gửi lên và email trong token không khớp nhau")
	}
	if invitationTemp.Phone != "" && r.Phone != invitationTemp.Phone {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phone gửi lên và phone trong token không khớp nhau")
	}

	return invitationTemp, nil
}

func (s *UserService) CheckUserRegistration(ctx context.Context, q *api.GetUserByPhoneRequest) (*api.GetUserByPhoneResponse, error) {
	_, ok := validate.NormalizePhone(q.Phone)
	if !ok {
		result := &api.GetUserByPhoneResponse{Exists: false}
		return result, nil
	}

	userByPhoneQuery := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: q.Phone,
	}
	err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, userByPhoneQuery)
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		result := &api.GetUserByPhoneResponse{Exists: false}
		return result, nil
	}
	result := &api.GetUserByPhoneResponse{Exists: true}
	return result, nil
}

func (s *UserService) getUserByPhoneAndByEmail(ctx context.Context, phone, email string) (userByPhone, userByEmail identitymodel.UserExtended, err error) {
	userByPhoneQuery := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: phone,
	}
	if err := s.UserStoreIface.GetUserByLogin(ctx, userByPhoneQuery); err != nil &&
		cm.ErrorCode(err) != cm.NotFound {
		return identitymodel.UserExtended{}, identitymodel.UserExtended{}, err
	}
	userByPhone = userByPhoneQuery.Result

	if email != "" {
		userByEmailQuery := &identitymodelx.GetUserByLoginQuery{
			PhoneOrEmail: email,
		}
		if err := s.UserStoreIface.GetUserByLogin(ctx, userByEmailQuery); err != nil &&
			cm.ErrorCode(err) != cm.NotFound {
			return identitymodel.UserExtended{}, identitymodel.UserExtended{}, err
		}
		userByEmail = userByEmailQuery.Result
	}
	return
}

func (s *UserService) Login(ctx context.Context, r *api.LoginRequest) (*api.LoginResponse, error) {
	query := &sqlstore.LoginUserQuery{
		PhoneOrEmail: r.Login,
		Password:     r.Password,
	}
	if err := s.LoginIface.LoginUser(ctx, query); err != nil {
		return nil, err
	}

	user := query.Result.User
	resp, err := s.CreateLoginResponse(
		ctx, nil, "", user.ID, user,
		r.AccountId, r.AccountType.Enum(),
		true, // Generate tokens for all accounts
		0,
	)
	if err != nil {
		return nil, err
	}
	setCookieForEcomify(ctx, resp.Account)
	return resp, err
}

func (s *UserService) ResetPassword(ctx context.Context, r *api.ResetPasswordRequest) (*api.ResetPasswordResponse, error) {
	key := fmt.Sprintf("ResetPassword %v-%v", r.Email, r.Phone)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 60*time.Second, "gửi email khôi phục mật khẩu",
		func() (interface{}, error) { return s.resetPassword(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*api.ResetPasswordResponse), nil
}

func (s *UserService) resetPassword(ctx context.Context, r *api.ResetPasswordRequest) (*api.ResetPasswordResponse, error) {
	// không thể gửi cùng 1 lúc cả phone và email
	if r.Email == "" && r.Phone == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu khôi phục mật khẩu không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Email != "" && r.Phone != "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu khôi phục mật khẩu không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Phone != "" {
		return s.resetPasswordUsingPhone(ctx, r)
	}
	if r.Email != "" {
		return s.resetPasswordUsingEmail(ctx, r)
	}
	panic("unreachable")
}

func (s *UserService) resetPasswordUsingPhone(ctx context.Context, r *api.ResetPasswordRequest) (*api.ResetPasswordResponse, error) {
	user, err := s.getUserByPhone(ctx, r.Phone)
	if err != nil {
		return nil, err
	}
	token, expiresIn := s.SS.Claim().Token, 0
	if token == "" {
		tokenCmd := &tokens.GenerateTokenCommand{}
		if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
			return nil, cm.Errorf(cm.Internal, err, "")
		}
		token, expiresIn = tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn
	}
	var redisCodeCount = fmt.Sprintf("reset-pasword-phone-%v", user.ID)
	if err = s.verifyPhone(ctx, auth.UsageResetPassword, user, 1*60*60, r.Phone, redisCodeCount, templatemessages.SmsResetPasswordTpl, templatemessages.SmsResetPasswordTplRepeat, false, token); err != nil {
		return nil, err
	}

	result := &api.ResetPasswordResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		Code:        "ok",
		Msg: fmt.Sprintf(
			"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ %v.", r.Phone, wl.X(ctx).CSEmail),
	}
	return result, nil
}

func (s *UserService) resetPasswordUsingEmail(ctx context.Context, r *api.ResetPasswordRequest) (*api.ResetPasswordResponse, error) {
	if !enabledEmail {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email khôi phục mật khẩu. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if !strings.Contains(r.Email, "@") {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	query := &identitymodelx.GetUserByLoginQuery{
		PhoneOrEmail: r.Email,
	}
	if err := s.UserStoreIface.GetUserByLogin(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, fmt.Sprintf("Người dùng chưa đăng ký. Vui lòng kiểm tra lại thông tin (hoặc đăng ký nếu chưa có tài khoản). Nếu cần thêm thông tin, vui lòng liên hệ %v.", wl.X(ctx).CSEmail)).
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
	if _, err := s.AuthStore.GenerateWithValue(tok, 24*60*60); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể khôi phục mật khẩu").WithMeta("reason", "can not generate token")
	}

	resetUrl, err := url.Parse(cfgEmail.ResetPasswordURL)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := resetUrl.Query()
	urlQuery.Set("t", tok.TokenStr)
	resetUrl.RawQuery = urlQuery.Encode()

	var b strings.Builder
	if err = templatemessages.ResetPasswordTpl.Execute(&b, map[string]interface{}{
		"FullName": user.FullName,
		"URL":      resetUrl.String(),
		"Email":    user.Email,
		"WlName":   wl.X(ctx).Name,
	}); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể khôi phục mật khẩu").WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    wl.X(ctx).CompanyName + " (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Khôi phục mật khẩu eTop",
		Content:     b.String(),
	}
	if err := s.EmailClient.SendMail(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.ResetPasswordResponse{
		Code: "ok",
		Msg: fmt.Sprintf(
			"Đã gửi email khôi phục mật khẩu đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ %v.", address, wl.X(ctx).CSEmail),
	}
	return result, nil
}

func (s *UserService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*pbcm.Empty, error) {
	if r.CurrentPassword == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing current_password", nil)
	}

	query := &sqlstore.LoginUserQuery{
		UserID:   s.SS.User().ID,
		Password: r.CurrentPassword,
	}
	if err := s.LoginIface.LoginUser(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.Unauthenticated, fmt.Sprintf("Mật khẩu không đúng. Vui lòng kiểm tra lại thông tin đăng nhập. Nếu cần thêm thông tin, vui lòng liên hệ %v.", wl.X(ctx).CSEmail)).
			Throw()
	}

	if len(r.NewPassword) < 8 {
		return nil, cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.NewPassword != r.ConfirmPassword {
		return nil, cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &identitymodelx.SetPasswordCommand{
		UserID:   s.SS.User().ID,
		Password: r.NewPassword,
	}
	if err := s.UserStoreIface.SetPassword(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.Empty{}
	return result, nil
}

func (s *UserService) ChangePasswordUsingToken(ctx context.Context, r *api.ChangePasswordUsingTokenRequest) (*pbcm.Empty, error) {
	key := fmt.Sprintf("ChangePasswordUsingToken %v-%v-%v", r.ResetPasswordToken, r.NewPassword, r.ConfirmPassword)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "khôi phục mật khẩu",
		func() (interface{}, error) { return s.changePasswordUsingToken(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.Empty), nil
}

func (s *UserService) changePasswordUsingToken(ctx context.Context, r *api.ChangePasswordUsingTokenRequest) (*pbcm.Empty, error) {
	return s.changePasswordUsingTokenForEmail(ctx, r)
}

func (s *UserService) changePasswordUsingTokenForEmail(ctx context.Context, r *api.ChangePasswordUsingTokenRequest) (*pbcm.Empty, error) {
	if r.ResetPasswordToken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing reset_password_token", nil)
	}
	var v map[string]string
	tok, err := s.AuthStore.Validate(auth.UsageResetPassword, r.ResetPasswordToken, &v)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể khôi phục mật khẩu (token không hợp lệ). Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	if err := s.changePassword("", v["email"], tok.UserID, ctx, r.NewPassword, r.ConfirmPassword); err != nil {
		return nil, err
	}
	s.AuthStore.Revoke(auth.UsageResetPassword, r.ResetPasswordToken)
	result := &pbcm.Empty{}
	return result, nil
}

func (s *UserService) changePassword(phone string, email string, tokUserID dot.ID, ctx context.Context, newPassword string, confirmPassword string) error {
	query := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: phone,
		Email: email,
	}
	if err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, query); err != nil {
		return err
	}
	user := query.Result
	if tokUserID != user.ID {
		return cm.Errorf(cm.InvalidArgument, nil, "Không thể khôi phục mật khẩu (token không hợp lệ). Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "user is not correct")
	}
	if len(newPassword) < 8 {
		return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if newPassword != confirmPassword {
		return cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}
	cmd := &identitymodelx.SetPasswordCommand{
		UserID:   user.ID,
		Password: newPassword,
	}
	if err := s.UserStoreIface.SetPassword(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (s *UserService) SessionInfo(ctx context.Context, r *pbcm.Empty) (*api.LoginResponse, error) {
	claim := s.SS.Claim().ClaimInfo
	resp, err := s.CreateLoginResponse(
		ctx,
		&claim,
		s.SS.Claim().Token,
		s.SS.Claim().UserID,
		s.SS.User().User,
		s.SS.Claim().AccountID,
		0,
		false,
		0,
	)
	return resp, err
}

func (s *UserService) SwitchAccount(ctx context.Context, r *api.SwitchAccountRequest) (*api.AccessTokenResponse, error) {
	if r.AccountId == 0 && !r.RegenerateTokens {
		return nil, cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}
	resp, err := s.CreateSessionResponse(
		ctx,
		nil, // Do not forward claim data
		"",  // Empty to generate new token
		s.SS.Claim().UserID,
		s.SS.User().User,
		r.AccountId,
		0,
		0,
	)
	if err != nil {
		return nil, err
	}
	if resp.Account == nil {
		return nil, cm.Error(cm.PermissionDenied, "Tài khoản không hợp lệ.", nil)
	}

	// set cookie for ecomify
	setCookieForEcomify(ctx, resp.Account)
	return resp, err
}

func setCookieForEcomify(ctx context.Context, account *etop.LoginAccount) {
	if account == nil {
		return
	}
	cookie := &http.Cookie{
		Name:     authservice.EcomAuthorization,
		Value:    account.AccessToken,
		Domain:   "",
		Expires:  time.Now().Add(24 * 60 * 60 * time.Second),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	bus.GetContext(ctx).WithValue(headers.CookieKey{}, []*http.Cookie{cookie})
	return
}

func (s *UserService) CreateSessionResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *identitymodel.User, preferAccountID dot.ID, preferAccountType int, adminID dot.ID) (*etop.AccessTokenResponse, error) {
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

func (s *UserService) CreateLoginResponse(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *identitymodel.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (*etop.LoginResponse, error) {
	resp, _, err := s.CreateLoginResponse2(ctx, claim, token, userID, user, preferAccountID, preferAccountType, generateAllTokens, adminID)
	return resp, err
}

func (s *UserService) CreateLoginResponse2(ctx context.Context, claim *claims.ClaimInfo, token string, userID dot.ID, user *identitymodel.User, preferAccountID dot.ID, preferAccountType int, generateAllTokens bool, adminID dot.ID) (_ *etop.LoginResponse, respShop *identitymodel.Shop, _ error) {

	// Retrieve user info
	if user != nil && user.ID != userID {
		return nil, nil, cm.Error(cm.Internal, "Invalid user", nil)
	}
	if user == nil {
		userQuery := &identitymodelx.GetUserByIDQuery{
			UserID: userID,
		}
		if err := s.UserStoreIface.GetUserByID(ctx, userQuery); err != nil {
			return nil, nil, err
		}
		user = userQuery.Result
	}

	// Retrieve list of accounts
	accQuery := &identitymodelx.GetAllAccountRolesQuery{UserID: userID}
	if err := s.AccountUserStore.GetAllAccountRoles(ctx, accQuery); err != nil {
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
			preferAccountType == account_tag.TagShop && account.Type == account_type.Shop,
			preferAccountType == account_tag.TagEtop && account.Type == account_type.Etop,
			preferAccountType == account_tag.TagAffiliate && account.Type == account_type.Affiliate:
			currentAccount = availableAccounts[i]
			currentAccountID = currentAccount.Id
		}
	}

	// Make sure user included `ref_sale` and `ref_aff`
	userFtRefSaffQuery := &identity.GetUserFtRefSaffByIDQuery{UserID: userID}
	if err := s.IdentityQuery.Dispatch(ctx, userFtRefSaffQuery); err != nil {
		return nil, nil, err
	}
	resp := &etop.LoginResponse{
		User:              convertpb.Convert_core_UserFtRefSaff_To_api_User(userFtRefSaffQuery.Result),
		Account:           currentAccount,
		AvailableAccounts: availableAccounts,
	}

	// Retrieve shop account
	if currentAccount != nil {
		switch {
		case idutil.IsShopID(currentAccountID):
			query := &identitymodelx.GetShopExtendedQuery{ShopID: currentAccountID}
			if err := s.ShopStore.GetShopExtended(ctx, query); err != nil {
				return nil, nil, cm.ErrorTracef(cm.Internal, err, "")
			}
			resp.Shop = convertpb.PbShopExtended(query.Result)
			respShop = query.Result.Shop

		case idutil.IsAffiliateID(currentAccountID):
			query := &identity.GetAffiliateByIDQuery{ID: currentAccountID}
			if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
				return nil, nil, cm.ErrorTracef(cm.Internal, err, "Account affiliate not found")
			}
			resp.Affiliate = convertpb.Convert_core_Affiliate_To_api_Affiliate(query.Result)
		case idutil.IsEtopAccountID(currentAccountID):
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
		if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
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
			if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
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
		resp.Account.UserAccount.Permission.Permissions = s.SS.Authorizer().ListActionsByRoles(resp.Account.UserAccount.Permission.Roles)
	}
	for _, account := range resp.AvailableAccounts {
		account.UserAccount.Permission.Permissions = s.SS.Authorizer().ListActionsByRoles(account.UserAccount.Permission.Roles)
	}

	return resp, respShop, nil
}

func (s *UserService) SendEmailVerificationUsingOTP(ctx context.Context, r *api.SendEmailVerificationUsingOTPRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("SendEmailVerificationUsingOTP %s-%s", s.SS.User().ID, r.Email)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "gửi email xác nhận tài khoản",
		func() (interface{}, error) { return s.sendEmailVerificationUsingOTP(ctx, r) })
	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.MessageResponse), nil
}

func (s *UserService) sendEmailVerificationUsingOTP(
	ctx context.Context, r *api.SendEmailVerificationUsingOTPRequest,
) (_ *pbcm.MessageResponse, err error) {
	if !enabledEmail {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email xác nhận tài khoản. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin địa chỉ email. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if !strings.Contains(r.Email, "@") {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	user := s.SS.User()
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok || user.Email != emailNorm.String() {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không đúng. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if !user.EmailVerifiedAt.IsZero() {
		result := cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
		return result, nil
	}

	var code string
	err = s.RedisStore.Get(fmt.Sprintf("Code:%v-%v-%v", user.Email, user.ID, keyRequestEmailVerifyCode), &code)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if code == "" {
		code, err = gencode.Random6Digits()
		if err != nil {
			return nil, err
		}
		err = s.RedisStore.SetWithTTL(fmt.Sprintf("Code:%v-%v-%v", user.Email, user.ID, keyRequestEmailVerifyCode), code, 2*60*60)
		if err != nil {
			return nil, err
		}
	}

	var b strings.Builder
	if err := templatemessages.EmailVerificationByOTPTpl.Execute(&b, map[string]interface{}{
		"Email":  user.Email,
		"Code":   code,
		"WlName": wl.X(ctx).Name,
	}); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    wl.X(ctx).CompanyName + " (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận địa chỉ email",
		Content:     b.String(),
	}
	if err := s.EmailClient.SendMail(ctx, cmd); err != nil {
		return nil, err
	}
	result := cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ %v.", address, wl.X(ctx).CSEmail))

	updateCmd := &identitymodelx.UpdateUserVerificationCommand{
		UserID:                  user.ID,
		EmailVerificationSentAt: time.Now(),
	}
	if err := s.UserStoreIface.UpdateUserVerification(ctx, updateCmd); err != nil {
		return nil, err
	}

	extra := s.SS.Claim().Extra
	if extra == nil {
		extra = map[string]string{}
	}

	extra[keyRequestEmailVerifyCode] = code
	extra[keyRequestAuthUsage] = auth.UsageEmailVerification
	if err := s.TokenStore.UpdateSession(ctx, s.SS.Claim().Token, extra); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) SendEmailVerification(ctx context.Context, r *api.SendEmailVerificationRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("SendEmailVerification %v-%v", s.SS.User().ID, r.Email)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "gửi email xác nhận tài khoản",
		func() (interface{}, error) { return s.sendEmailVerification(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.MessageResponse), nil
}

func (s *UserService) sendEmailVerification(ctx context.Context, r *api.SendEmailVerificationRequest) (*pbcm.MessageResponse, error) {
	if !enabledEmail {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email xác nhận tài khoản. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin địa chỉ email. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if !strings.Contains(r.Email, "@") {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	user := s.SS.User()
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	if !ok || user.Email != emailNorm.String() {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không đúng. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if !user.EmailVerifiedAt.IsZero() {
		result := cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
		return result, nil
	}

	tok := &auth.Token{
		TokenStr: "",
		Usage:    auth.UsageEmailVerification,
		UserID:   user.ID,
		Value: map[string]string{
			"email": user.Email,
		},
	}
	if _, err := s.AuthStore.GenerateWithValue(tok, 24*60*60); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate token")
	}

	verificationUrl, err := url.Parse(cfgEmail.EmailVerificationURL)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := verificationUrl.Query()
	urlQuery.Set("t", tok.TokenStr)
	verificationUrl.RawQuery = urlQuery.Encode()

	var b strings.Builder
	if err := templatemessages.EmailVerificationTpl.Execute(&b, map[string]interface{}{
		"FullName": user.FullName,
		"URL":      verificationUrl.String(),
		"Email":    user.Email,
		"WlName":   wl.X(ctx).Name,
	}); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    wl.X(ctx).CompanyName + " (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận địa chỉ email",
		Content:     b.String(),
	}
	if err := s.EmailClient.SendMail(ctx, cmd); err != nil {
		return nil, err
	}
	result := cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ %v.", address, wl.X(ctx).CSEmail))

	updateCmd := &identitymodelx.UpdateUserVerificationCommand{
		UserID:                  user.ID,
		EmailVerificationSentAt: time.Now(),
	}
	if err := s.UserStoreIface.UpdateUserVerification(ctx, updateCmd); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) SendPhoneVerification(ctx context.Context, r *api.SendPhoneVerificationRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("SendPhoneVerification %v-%v", s.SS.Claim().Token, r.Phone)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 60*time.Second, "gửi tin nhắn xác nhận số điện thoại",
		func() (interface{}, error) { return s.sendPhoneVerification(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.MessageResponse), nil
}

func (s *UserService) sendPhoneVerification(ctx context.Context, r *api.SendPhoneVerificationRequest) (resp *pbcm.MessageResponse, _ error) {
	if !enabledSMS {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi tin nhắn xác nhận tài khoản. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Phone == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin số điện thoại. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	// update token when user not exists
	if s.SS.Claim().UserID == 0 {
		return s.sendPhoneVerificationForRegister(ctx, r)
	}
	getUserByID := &identitymodelx.GetUserByIDQuery{
		UserID: s.SS.Claim().UserID,
	}
	if err := s.UserStoreIface.GetUserByID(ctx, getUserByID); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	user := getUserByID.Result
	_, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại không hợp le. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	var redisCodeCount = fmt.Sprintf("confirm-phone-%v", s.SS.User().ID)
	if err := s.verifyPhone(ctx, auth.UsagePhoneVerification, user, 2*60*60, r.Phone, redisCodeCount, templatemessages.SmsVerificationTpl, templatemessages.SmsVerificationTplRepeat, true, s.SS.Claim().Token); err != nil {
		return nil, err
	}

	return cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ %v.", r.Phone, wl.X(ctx).CSEmail)), nil
}

func (s *UserService) countSendtimeMsg(code, msgFirst, msgSecond string) (string, error) {
	h := sha1.New()
	h.Write([]byte(code))
	bs := h.Sum(nil)
	code = fmt.Sprintf("%x", bs)
	var msg string
	var sendTime int
	err := s.RedisStore.Get(code, &sendTime)
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err != nil && err == redis.ErrNil {
		err = s.RedisStore.SetWithTTL(code, 1, 2*60*60)
		if err != nil {
			return "", err
		}
		msg = msgFirst
	} else {
		sendTime++
		msg = msgSecond
		err = s.RedisStore.SetWithTTL(code, sendTime, 2*60*60)
		if err != nil {
			return "", err
		}
		// Print sendtime in last value of msg
		countMsg := strings.Count(msgSecond, "%v")
		i := 1
		var arrString = []interface{}{}
		for i < countMsg {
			arrString = append(arrString, "%v")
			i++
		}
		arrString = append(arrString, sendTime)
		msg = fmt.Sprintf(msgSecond, arrString...)
	}
	return msg, nil
}

func (s *UserService) VerifyEmailUsingToken(ctx context.Context, r *api.VerifyEmailUsingTokenRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("VerifyEmailUsingToken %v-%v", s.SS.User().ID, r.VerificationToken)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "xác nhận địa chỉ email",
		func() (interface{}, error) { return s.verifyEmailUsingToken(ctx, r) })

	if err != nil {
		return nil, err
	}
	return res.(*pbcm.MessageResponse), nil
}

func (s *UserService) verifyEmailUsingToken(ctx context.Context, r *api.VerifyEmailUsingTokenRequest) (*pbcm.MessageResponse, error) {
	if r.VerificationToken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing verification_token", nil)
	}

	var v map[string]string
	tok, err := s.AuthStore.Validate(auth.UsageEmailVerification, r.VerificationToken, &v)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể xác nhận địa chỉ email (token không hợp lệ). Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	user := s.SS.User()
	if user.ID != tok.UserID || user.Email != v["email"] {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể xác nhận địa chỉ email (địa chỉ email không đúng). Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	if user.EmailVerifiedAt.IsZero() {
		cmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			EmailVerifiedAt: time.Now(),
		}
		if err := s.UserStoreIface.UpdateUserVerification(ctx, cmd); err != nil {
			return nil, err
		}
	}

	s.AuthStore.Revoke(auth.UsageEmailVerification, r.VerificationToken)
	result := cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
	return result, nil
}

func (s *UserService) VerifyEmailUsingOTP(ctx context.Context, r *api.VerifyEmailUsingOTPRequest) (*pbcm.MessageResponse, error) {
	extra := s.SS.Claim().Extra
	user := s.SS.User()

	if r.VerificationToken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing verification_token", nil)
	}

	if extra == nil || extra[keyRequestAuthUsage] != auth.UsageEmailVerification {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể xác nhận địa chỉ email. Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if extra[keyRequestEmailVerifyCode] != r.VerificationToken {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể xác nhận địa chỉ email (mã xác thực không đúng). Vui lòng thử lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	delete(extra, keyRequestAuthUsage)
	delete(extra, keyRequestEmailVerifyCode)

	if err := s.TokenStore.UpdateSession(ctx, s.SS.Claim().Token, extra); err != nil {
		return nil, err
	}

	if err := s.RedisStore.Del(fmt.Sprintf("Code:%v-%v-%v", user.Email, user.ID, keyRequestEmailVerifyCode)); err != nil {
		return nil, err
	}

	if user.EmailVerifiedAt.IsZero() {
		cmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			EmailVerifiedAt: time.Now(),
		}
		if err := s.UserStoreIface.UpdateUserVerification(ctx, cmd); err != nil {
			return nil, err
		}
	}

	result := cmapi.Message("ok", "Địa chỉ email đã được xác nhận thành công.")
	return result, nil
}

func (s *UserService) VerifyPhoneUsingToken(ctx context.Context, r *api.VerifyPhoneUsingTokenRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("VerifyPhoneUsingToken %v-%v", s.SS.Claim().Token, r.VerificationToken)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "xác nhận số điện thoại",
		func() (interface{}, error) { return s.verifyPhoneUsingToken(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.MessageResponse), nil
}

func (s *UserService) verifyPhoneUsingToken(ctx context.Context, r *api.VerifyPhoneUsingTokenRequest) (*pbcm.MessageResponse, error) {
	if r.VerificationToken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing code", nil)
	}
	if s.SS.Claim().Extra[keyRequestVerifyCode] != "" && s.SS.Claim().Extra[keyRequestVerifyCode] != r.VerificationToken {
		result := cmapi.Message("fail", "Mã xác thực không chính xác.")
		return result, nil
	}
	if s.SS.Claim().UserID == 0 && s.SS.Claim().Extra != nil {
		extra := s.SS.Claim().Extra
		extra[keyRequestPhoneVerificationVerified] = "1"

		if err := s.TokenStore.UpdateSession(ctx, s.SS.Claim().Token, extra); err != nil {
			return nil, err
		}
		result := cmapi.Message("ok", "Số điện thoại đã được xác nhận thành công.")
		return result, nil
	}
	getUserByID := &identitymodelx.GetUserByIDQuery{
		UserID: s.SS.Claim().UserID,
	}
	if err := s.UserStoreIface.GetUserByID(ctx, getUserByID); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	var v map[string]string
	user := getUserByID.Result
	tok, code, v := s.getToken(s.SS.Claim().Extra[keyRequestAuthUsage], user.ID, user.Phone)
	if tok == nil || code == "" || v == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if v["extra"] != user.Phone {
		s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr)
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	var err error
	if code != r.VerificationToken {
		// Delete token after 3 times
		if len(v["tries"]) >= 2 {
			if err = s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
				ll.Error("Can not revoke token", l.Error(err))
			}
		} else {
			v["tries"] += "."
			s.AuthStore.SetTTL(tok, 60*60)
		}
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	if user.PhoneVerifiedAt.IsZero() {
		cmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:          user.ID,
			PhoneVerifiedAt: time.Now(),
		}
		if err := s.UserStoreIface.UpdateUserVerification(ctx, cmd); err != nil {
			return nil, err
		}
	}

	s.AuthStore.Revoke(auth.UsagePhoneVerification, tok.TokenStr)
	result := cmapi.Message("ok", "Số điện thoại đã được xác nhận thành công.")
	return result, nil
}

func (s *UserService) VerifyPhoneResetPasswordUsingToken(ctx context.Context, r *api.VerifyPhoneResetPasswordUsingTokenRequest) (*api.VerifyPhoneResetPasswordUsingTokenResponse, error) {
	key := fmt.Sprintf("VerifyPhoneResetPasswordUsingToken %v-%v", s.SS.Claim().Token, r.VerificationToken)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "xác nhận số điện thoại",
		func() (interface{}, error) { return s.verifyPhoneResetPasswordUsingToken(ctx, r) })
	if err != nil {
		return nil, err
	}
	return resp.(*api.VerifyPhoneResetPasswordUsingTokenResponse), nil
}

func (s *UserService) verifyPhoneResetPasswordUsingToken(ctx context.Context, r *api.VerifyPhoneResetPasswordUsingTokenRequest) (*api.VerifyPhoneResetPasswordUsingTokenResponse, error) {
	if r.VerificationToken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing code", nil)
	}
	if s.SS.Claim().Extra[keyRequestVerifyCode] != "" && s.SS.Claim().Extra[keyRequestVerifyCode] != r.VerificationToken {
		return nil, cm.Error(cm.InvalidArgument, "Mã xác nhận không chính xác.", nil)
	}
	getUserByID := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: s.SS.Claim().Extra[keyRequestVerifyPhone],
	}
	if err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, getUserByID); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	var v map[string]string
	user := getUserByID.Result
	tok, code, v := s.getToken(s.SS.Claim().Extra[keyRequestAuthUsage], user.ID, user.Phone)
	if tok == nil || code == "" || v == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if v["extra"] != user.Phone {
		s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr)
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	var err error
	if code != r.VerificationToken {
		// Delete token after 3 times
		if len(v["tries"]) >= 2 {
			if err = s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
				ll.Error("Can not revoke token", l.Error(err))
			}
		} else {
			v["tries"] += "."
			s.AuthStore.SetTTL(tok, 60*60)
		}
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	recaptchaToken := &auth.Token{
		TokenStr: "",
		Usage:    auth.UsageResetPassword,
		UserID:   user.ID,
		Value: map[string]string{
			"email": user.Email,
		},
	}
	if _, err := s.AuthStore.GenerateWithValue(recaptchaToken, 1*60*60); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể khôi phục mật khẩu").WithMeta("reason", "can not generate token")
	}

	s.AuthStore.Revoke(auth.UsageResetPassword, tok.TokenStr)
	result := &api.VerifyPhoneResetPasswordUsingTokenResponse{ResetPasswordToken: recaptchaToken.TokenStr}
	return result, nil
}

func (s *UserService) UpgradeAccessToken(ctx context.Context, r *api.UpgradeAccessTokenRequest) (*api.AccessTokenResponse, error) {
	key := fmt.Sprintf("UpgradeAccessToken %v-%v", s.SS.User().ID, r.Stoken)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 15*time.Second, "cập nhật thông tin",
		func() (interface{}, error) { return s.upgradeAccessToken(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*api.AccessTokenResponse), nil
}

func (s *UserService) upgradeAccessToken(ctx context.Context, r *api.UpgradeAccessTokenRequest) (*api.AccessTokenResponse, error) {
	if r.Stoken == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing code", nil)
	}

	user := s.SS.User()
	tok, code, v := s.getToken(auth.UsageSToken, user.ID, user.Email)
	if tok == nil || code == "" || v == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if v["extra"] != user.Email {
		s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr)
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	if code == r.Stoken {
		expiresAt := time.Now().Add(30 * time.Minute)
		result, err := s.CreateSessionResponse(
			ctx,
			&claims.ClaimInfo{
				SToken:          true,
				STokenExpiresAt: &expiresAt,
			},
			"", // Generate new token
			user.ID,
			user.User,
			s.SS.Claim().AccountID,
			0,
			0,
		)
		if err != nil {
			return nil, err
		}
		if err = s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
			ll.Error("Can not revoke token", l.Error(err))
		}
		return result, nil
	}

	// Delete token after 3 times
	if len(v["tries"]) >= 2 {
		if err := s.AuthStore.Revoke(auth.UsageSToken, tok.TokenStr); err != nil {
			ll.Error("Can not revoke token", l.Error(err))
		}
	} else {
		v["tries"] += "."
		s.AuthStore.SetTTL(tok, 60*60)
	}
	return nil, cm.Errorf(cm.PermissionDenied, nil, "Mã xác nhận không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
}

func (s *UserService) SendSTokenEmail(ctx context.Context, r *api.SendSTokenEmailRequest) (*pbcm.MessageResponse, error) {
	key := fmt.Sprintf("SendSTokenEmail %v-%v-%v", s.SS.User().ID, r.Email, r.AccountId)
	resp, _, err := idempgroup.DoAndWrap(
		ctx, key, 30*time.Second, "gửi email xác nhận",
		func() (interface{}, error) { return s.sendSTokenEmail(ctx, r) })

	if err != nil {
		return nil, err
	}
	return resp.(*pbcm.MessageResponse), nil
}

func (s *UserService) sendSTokenEmail(ctx context.Context, r *api.SendSTokenEmailRequest) (*pbcm.MessageResponse, error) {
	if !enabledEmail {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email xác nhận. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "not configured")
	}
	if r.Email == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin email. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if !strings.Contains(r.Email, "@") {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Địa chỉ email không hợp lệ. Nếu cần thêm thông tin vui lòng liên hệ  %v.", wl.X(ctx).CSEmail)
	}
	if r.AccountId == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}

	user := s.SS.User()
	accQuery := &identitymodelx.GetAccountRolesQuery{
		AccountID: r.AccountId,
		UserID:    user.ID,
	}
	if err := s.AccountUserStore.GetAccountUserExtended(ctx, accQuery); err != nil {
		return nil, err
	}
	account := accQuery.Result.Account
	emailNorm, ok := validate.NormalizeEmail(r.Email)
	userEmail, _ := validate.NormalizeEmail(user.Email)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ. Vui lòng thử lại hoặc liên hệ  %v.", wl.X(ctx).CSEmail)
	}
	if emailNorm != userEmail {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Email không đúng. Vui lòng thử lại hoặc liên hệ  %v.", wl.X(ctx).CSEmail)
	}

	_, code, _, err := s.generateToken(auth.UsageSToken, user.ID, true, 2*60*60, user.Email)
	if err != nil {
		return nil, err
	}

	emailData := map[string]interface{}{
		"FullName":    user.FullName,
		"Code":        code,
		"Email":       user.Email,
		"AccountName": account.Name,
		"WlName":      wl.X(ctx).Name,
	}
	switch account.Type {
	case account_type.Shop:
		emailData["AccountType"] = model.AccountTypeLabel(account.Type)
	case account_type.Etop:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email đến tài khoản %v", account.Name, account).WithMeta("type", account.Type.String())
	default:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi email đến tài khoản %v", account.Name).WithMeta("type", account.Type.String())
	}

	var b strings.Builder
	if err := templatemessages.EmailSTokenTpl.Execute(&b, emailData); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Không thể gửi email đến tài khoản %v", account.Name).WithMeta("reason", "can not generate email content")
	}

	address := user.Email
	cmd := &email.SendEmailCommand{
		FromName:    wl.X(ctx).CompanyName + " (no-reply)",
		ToAddresses: []string{address},
		Subject:     "Xác nhận thay đổi thông tin tài khoản",
		Content:     b.String(),
	}
	if err := s.EmailClient.SendMail(ctx, cmd); err != nil {
		return nil, err
	}
	result := cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi email kèm mã xác nhận đến địa chỉ %v. Vui lòng kiểm tra email (kể cả trong hộp thư spam). Nếu cần thêm thông tin, vui lòng liên hệ %v.", address, wl.X(ctx).CSEmail))
	return result, nil
}

func (s *UserService) getToken(usage string, userID dot.ID, extra string) (*auth.Token, string, map[string]string) {
	tokStr := fmt.Sprintf("%v-%v", userID, extra)
	var v map[string]string
	var code string

	tok, err := s.AuthStore.Validate(usage, tokStr, &v)
	if err == nil && v != nil && len(v["code"]) == 6 {
		code = v["code"]
		return tok, code, v
	}

	_ = s.AuthStore.Revoke(usage, tokStr)
	return nil, "", nil
}

func (s *UserService) generateToken(usage string, userID dot.ID, generate bool, ttl int, extra string) (*auth.Token, string, map[string]string, error) {
	tokStr := fmt.Sprintf("%v-%v", userID, extra)
	tok, code, v := s.getToken(usage, userID, extra)
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
	tok, err = s.AuthStore.GenerateWithValue(tok, ttl)
	if err != nil {
		return nil, "", nil, cm.Error(cm.Internal, "", err)
	}
	return tok, code, v, nil
}

func (s *UserService) UpdateReferenceUser(ctx context.Context, r *api.UpdateReferenceUserRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &identity.UpdateUserReferenceUserIDCommand{
		UserID:       s.SS.Claim().UserID,
		RefUserPhone: r.Phone,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *UserService) UpdateReferenceSale(ctx context.Context, r *api.UpdateReferenceSaleRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &identity.UpdateUserReferenceSaleIDCommand{
		UserID:       s.SS.Claim().UserID,
		RefSalePhone: r.Phone,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *UserService) InitSession(ctx context.Context, r *pbcm.Empty) (*api.LoginResponse, error) {
	tokenCmd := &tokens.GenerateTokenCommand{
		ClaimInfo: claims.ClaimInfo{},
	}
	if err := s.TokenStore.GenerateToken(ctx, tokenCmd); err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}

	result := generateShopLoginResponse(
		tokenCmd.Result.TokenStr, tokenCmd.Result.ExpiresIn,
	)
	return result, nil
}

func generateShopLoginResponse(accessToken string, expiresIn int) *etop.LoginResponse {
	resp := &etop.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}
	return resp
}

func (s *UserService) sendPhoneVerificationForRegister(ctx context.Context, r *api.SendPhoneVerificationRequest) (*pbcm.MessageResponse, error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Error(cm.FailedPrecondition, "Số điện thoại không hợp lệ", nil)
	}
	var msg string
	var redisCodeCount = fmt.Sprintf("confirm-phone-%v-%v", s.SS.Claim().UserID, r.Phone)
	redisCodeCount = sqlstore.EncodePassword(redisCodeCount)
	msg, err := s.countSendtimeMsg(redisCodeCount, templatemessages.SmsVerificationTpl, templatemessages.SmsVerificationTplRepeat)
	if err != nil {
		return nil, err
	}
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err = s.sendPhoneVerificationImpl(ctx, nil, 2*60*60, auth.UsagePhoneVerification, r.Phone, msg, false, s.SS.Claim().Token); err != nil {
		return nil, err
	}
	return cmapi.Message("ok", fmt.Sprintf(
		"Đã gửi tin nhắn kèm mã xác nhận đến số điện thoại %v. Vui lòng kiểm tra tin nhắn. Nếu cần thêm thông tin, vui lòng liên hệ %v.", phone, wl.X(ctx).CSEmail)), nil
}

func (s *UserService) sendPhoneVerificationImpl(ctx context.Context, user *identitymodel.User, ttl int, usage string,
	phone string, msg string, checkVerifyPhoneForUser bool, token string) error {
	var userIDUse dot.ID
	userIDUse = 0
	if user != nil {
		userIDUse = user.ID
	}
	phoneUse := phone
	_, code, _, err := s.generateToken(usage, userIDUse, true, ttl, phoneUse)
	if err != nil {
		return err
	}
	msgUser := fmt.Sprintf(msg, code)
	cmd := &sms.SendSMSCommand{
		Phone:   phoneUse,
		Content: msgUser,
	}
	if err := s.SMSClient.SendSMS(ctx, cmd); err != nil {
		return err
	}
	extra := s.SS.Claim().Extra
	if extra == nil {
		extra = map[string]string{}
	}
	if !checkVerifyPhoneForUser {
		extra[keyRequestVerifyPhone] = phoneUse
		extra[keyRequestVerifyCode] = code
	}
	extra[keyRequestAuthUsage] = usage

	if err := s.TokenStore.UpdateSession(ctx, token, extra); err != nil {
		return err
	}
	return nil
}

func (s *UserService) getUserByPhone(ctx context.Context, phone string) (*identitymodel.User, error) {
	_, ok := validate.NormalizePhone(phone)
	if !ok {
		return nil, cm.Error(cm.FailedPrecondition, "Số điện thoại không hợp lệ", nil)
	}
	userByPhone := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: phone,
	}
	if err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, userByPhone); err != nil {
		return nil, cm.Error(cm.FailedPrecondition, "Số điện này không hợp lệ vì chưa được đăng kí", nil)
	}
	return userByPhone.Result, nil
}

func (s *UserService) verifyPhone(ctx context.Context, usage string, user *identitymodel.User, ttl int, phone string, redisCodeCount string, msgFirstime string, msgMultiTime string, checkVerifyPhoneForUser bool, token string) error {
	if user != nil && user.Phone != phone {
		return cm.Error(cm.FailedPrecondition, "Số điện này không hợp lệ vì chưa được đăng kí", nil)
	}
	msg, err := s.countSendtimeMsg(redisCodeCount, msgFirstime, msgMultiTime)
	if err != nil {
		return err
	}
	if checkVerifyPhoneForUser {
		if !user.PhoneVerifiedAt.IsZero() {
			return nil
		}
		if err := s.sendPhoneVerificationImpl(ctx, user, ttl, usage, user.Phone, msg, true, token); err != nil {
			return err
		}
		updateCmd := &identitymodelx.UpdateUserVerificationCommand{
			UserID:                  user.ID,
			PhoneVerificationSentAt: time.Now(),
		}
		if err := s.UserStoreIface.UpdateUserVerification(ctx, updateCmd); err != nil {
			return err
		}
		return nil
	}
	if err := s.sendPhoneVerificationImpl(ctx, user, ttl, usage, user.Phone, msg, false, token); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangeRefAff(ctx context.Context, request *api.ChangeUserRefAffRequest) (*pbcm.Empty, error) {
	cmdUpdate := &identity.UpdateUserRefCommand{
		UserID: s.SS.User().ID,
		RefAff: request.RefAff,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmdUpdate); err != nil {
		return nil, err
	}
	return &pbcm.Empty{}, nil
}

func (s *UserService) RequestRegisterSimplify(ctx context.Context, r *api.RequestRegisterSimplifyRequest) (*pbcm.MessageResponse, error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}

	// send OTP
	key := fmt.Sprintf("%v-%v", keyRequestRegisterSimplify, phone.String())
	code, err := s.RedisStore.GetString(key)
	if err != nil || code == "" {
		code, err = gencode.Random6Digits()
		if err != nil {
			return nil, cm.Error(cm.Internal, "", err)
		}
		_ = s.RedisStore.SetStringWithTTL(key, code, 5*60)
	}

	smsMsg := fmt.Sprintf(templatemessages.SmsRegisterSimplifyTpl, code)
	cmd := &sms.SendSMSCommand{
		Phone:   phone.String(),
		Content: smsMsg,
	}
	if err := s.SMSClient.SendSMS(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.MessageResponse{
		Code: "ok",
		Msg:  "Đã gửi otp tới số điện thoại cung cấp",
	}, nil
}

func (s *UserService) RegisterSimplify(ctx context.Context, r *api.RegisterSimplifyRequest) (*api.LoginResponse, error) {
	phone, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	if r.OTP == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp mã OTP")
	}
	key := fmt.Sprintf("%v-%v", keyRequestRegisterSimplify, phone.String())
	code, err := s.RedisStore.GetString(key)
	if err != nil && err != redis.ErrNil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	if code != r.OTP {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Mã OTP không hợp lệ")
	}

	cmd := &identity.RegisterSimplifyCommand{
		Phone: phone.String(),
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	query := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Phone: phone.String(),
	}
	if err := s.UserStoreIface.GetUserByEmailOrPhone(ctx, query); err != nil {
		return nil, err
	}
	user := query.Result
	resp, err := s.CreateLoginResponse(ctx, nil, "", user.ID, user, 0, account_type.Shop.Enum(), true, 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
