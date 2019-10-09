package sqlstore

import (
	"context"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/xerrors"
)

func init() {
	bus.AddHandlers("sql",
		GetUserByID,
		GetUserByLogin,
		CreateUser,
		SetPassword,
		GetSignedInUser,
		UpdateUserVerification,
		UpdateUserIdentifier,
	)
}

type UserStore struct {
	ctx   context.Context
	ft    UserFilters
	preds []interface{}
}

func User(ctx context.Context) *UserStore {
	return &UserStore{ctx: ctx}
}

func (s *UserStore) ID(id int64) *UserStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *UserStore) Get() (*model.User, error) {
	var user model.User
	err := x.Where(s.preds...).ShouldGet(&user)
	return &user, err
}

func (s *UserStore) List() ([]*model.User, error) {
	var users model.Users
	err := x.Where(s.preds).Find(&users)
	return users, err
}

func GetSignedInUser(ctx context.Context, query *model.GetSignedInUserQuery) error {
	userQuery := &model.GetUserByIDQuery{UserID: query.UserID}
	if err := GetUserByID(ctx, userQuery); err != nil {
		return err
	}
	query.Result = &model.SignedInUser{
		User: userQuery.Result,
	}
	return nil
}

func GetUserByID(ctx context.Context, query *model.GetUserByIDQuery) error {
	if query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	query.Result = new(model.User)
	return x.Where("id = ?", query.UserID).
		ShouldGet(query.Result)
}

func GetUserByLogin(ctx context.Context, query *model.GetUserByLoginQuery) error {
	if query.UserID == 0 && query.PhoneOrEmail == "" {
		return cm.Error(cm.InvalidArgument, "Missing required fields", nil)
	}

	s := x.Table("user")
	if query.UserID != 0 {
		s = s.Where("id = ?", query.UserID)
	}
	if query.PhoneOrEmail != "" {
		if validate.IsEmail(query.PhoneOrEmail) {
			email, ok := validate.NormalizeEmail(query.PhoneOrEmail)
			if !ok {
				return cm.Error(cm.NotFound, "", nil)
			}
			s = s.Where("email = ?", email)

		} else {
			phone, ok := validate.NormalizePhone(query.PhoneOrEmail)
			if !ok {
				return cm.Error(cm.NotFound, "", nil)
			}
			s = s.Where("phone = ?", phone)
		}
	}

	user := new(model.User)
	if err := s.ShouldGet(user); err != nil {
		return err
	}

	userInternal := new(model.UserInternal)
	if err := x.Where("id = ?", user.ID).ShouldGet(userInternal); err != nil {
		return err
	}

	query.Result.User = user
	query.Result.UserInternal = userInternal
	return nil
}

func CreateUser(ctx context.Context, cmd *model.CreateUserCommand) error {
	if cmd.IsStub {
		if cmd.Email == "" && cmd.Phone == "" {
			return cm.Error(cm.InvalidArgument, "Missing email and phone", nil)
		}
	} else {
		if cmd.Email == "" {
			return cm.Error(cm.InvalidArgument, "Vui lòng nhập email", nil)
		}
		if cmd.Phone == "" {
			return cm.Error(cm.InvalidArgument, "Vui lòng nhập số điện thoại", nil)
		}
		if cmd.FullName == "" {
			return cm.Error(cm.InvalidArgument, "Vui lòng nhập tên", nil)
		}
		if cmd.Password == "" {
			return cm.Error(cm.InvalidArgument, "Vui lòng nhập mật khẩu", nil)
		}
		if len(cmd.Password) < 8 {
			return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
		}
		if !cmd.AgreeTOS {
			return cm.Error(cm.InvalidArgument, "Must agree tos", nil)
		}
	}

	var ok bool
	var err error
	var emailNorm model.NormalizedEmail
	var phoneNorm model.NormalizedPhone
	if emailNorm, ok = validate.NormalizeEmail(cmd.Email); cmd.Email != "" && !ok {
		return cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
	}
	cmd.Email = emailNorm.String()
	if phoneNorm, ok = validate.NormalizePhone(cmd.Phone); cmd.Phone != "" && !ok {
		return cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}
	cmd.Phone = phoneNorm.String()
	cmd.FullName, cmd.ShortName, err = NormalizeFullName(cmd.FullName, cmd.ShortName)
	if err != nil {
		return err
	}

	// IsTest
	{
		_, _, emailTest := validate.TrimTest(cmd.Email)
		_, _, phoneTest := validate.TrimTest(cmd.Phone)
		switch {
		case !emailTest && !phoneTest:
			// continue
		case emailTest && phoneTest:
			cmd.IsTest = true
		default:
			return cm.Error(cm.InvalidArgument, "Both phone and email must end with -test", nil)
		}
	}

	return inTransaction(func(s Qx) error {
		return createUser(ctx, s, cmd)
	})
}

func NormalizeFullName(fullName, shortName string) (string, string, error) {
	var ok bool
	if fullName, ok = validate.NormalizeName(fullName); fullName != "" && !ok {
		return "", "", cm.Error(cm.InvalidArgument, "Tên không hợp lệ", nil)
	}
	if shortName == "" {
		ss := strings.Split(fullName, " ")
		shortName = ss[len(ss)-1]
	} else if shortName, ok = validate.NormalizeName(shortName); shortName != "" && !ok {
		return "", "", cm.Error(cm.InvalidArgument, "Tên không hợp lệ (short name)", nil)
	}
	return fullName, shortName, nil
}

const (
	UserEmailKey = "user_email_key"
	UserPhoneKey = "user_phone_key"

	MsgCreateUserDuplicatedPhone = `Số điện thoại đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng số điện thoại khác. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.`
	MsgCreateUserDuplicatedEmail = `Email đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng email khác. Nếu cần thêm thông tin, vui lòng liên hệ hotro@etop.vn.`
)

func createUser(ctx context.Context, s Qx, cmd *model.CreateUserCommand) error {
	switch cmd.Status {
	case model.S3Positive:
	case model.S3Zero:
	default:
		return cm.Error(cm.InvalidArgument, "Invalid status", nil)
	}
	if cmd.IsStub != (cmd.Status == model.S3Zero) {
		return cm.Error(cm.InvalidArgument, "Mismatch status", nil)
	}

	now := time.Now()
	userID := cm.NewIDWithTag(model.TagUser)
	user := &model.User{
		ID:        userID,
		UserInner: cmd.UserInner,
		Status:    cmd.Status,
		Source:    cmd.Source,
	}
	if cmd.IsStub {
		user.Identifying = model.UserIdentifyingStub
	} else {
		user.AgreedTOSAt = now
		if user.Email != "" && user.Phone != "" {
			user.Identifying = model.UserIdentifyingFull
		} else {
			user.Identifying = model.UserIdentifyingHalf
		}
	}
	if cmd.AgreeEmailInfo {
		user.AgreedEmailInfoAt = now
	}
	if cmd.IsTest {
		user.IsTest = 1
	}

	userInternal := &model.UserInternal{
		ID: userID,
	}
	if cmd.Password != "" {
		userInternal.Hashpwd = login.EncodePassword(cmd.Password)
	}

	_, err := s.Insert(user, userInternal)
	if xerr, ok := err.(*xerrors.APIError); ok && xerr.Err != nil {
		msg := xerr.Err.Error()
		switch {
		case strings.Contains(msg, UserEmailKey):
			err = cm.Error(cm.FailedPrecondition, MsgCreateUserDuplicatedEmail, nil)
		case strings.Contains(msg, UserPhoneKey):
			err = cm.Error(cm.FailedPrecondition, MsgCreateUserDuplicatedPhone, nil)
		}
	}
	cmd.Result.User = user
	cmd.Result.UserInternal = userInternal
	return err
}

func SetPassword(ctx context.Context, cmd *model.SetPasswordCommand) error {
	if cmd.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	userInternal := &model.UserInternal{
		Hashpwd: login.EncodePassword(cmd.Password),
	}
	if err := x.Where("id = ?", cmd.UserID).
		ShouldUpdate(userInternal); err != nil {
		return err
	}
	return nil
}

func UpdateUserVerification(ctx context.Context, cmd *model.UpdateUserVerificationCommand) error {
	if cmd.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	count := 0
	var user model.User

	s := x.Where("id = ?", cmd.UserID)
	if !cmd.EmailVerifiedAt.IsZero() {
		count++
		user.EmailVerifiedAt = cmd.EmailVerifiedAt
	}
	if !cmd.PhoneVerifiedAt.IsZero() {
		count++
		user.PhoneVerifiedAt = cmd.PhoneVerifiedAt
	}
	if !cmd.EmailVerificationSentAt.IsZero() {
		count++
		user.EmailVerificationSentAt = cmd.EmailVerificationSentAt
	}
	if !cmd.PhoneVerificationSentAt.IsZero() {
		count++
		user.PhoneVerificationSentAt = cmd.PhoneVerificationSentAt
	}
	if count != 1 {
		return cm.Error(cm.InvalidArgument, "Invalid params", nil)
	}

	return s.ShouldUpdate(&user)
}

func UpdateUserIdentifier(ctx context.Context, cmd *model.UpdateUserIdentifierCommand) error {
	if cmd.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}
	if cmd.Phone != "" && cmd.Email != "" {
		return cm.Error(cm.InvalidArgument, "Can not update both identifier", nil)
	}

	if cmd.Phone != "" {
		phoneNorm, ok := validate.NormalizePhone(cmd.Phone)
		if !ok || phoneNorm.String() != cmd.Phone {
			return cm.Error(cm.InvalidArgument, "Invalid phone", nil)
		}
	}
	if cmd.Email != "" {
		emailNorm, ok := validate.NormalizeEmail(cmd.Email)
		if !ok || emailNorm.String() != cmd.Email {
			return cm.Error(cm.InvalidArgument, "Invalid email", nil)
		}
	}
	if cmd.Password == "" {
		return cm.Error(cm.InvalidArgument, "Vui lòng nhập mật khẩu", nil)
	}
	if len(cmd.Password) < 8 {
		return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	// MUSTDO: Merge update user identifier and create user
	// Handle password, agree tos, email info and other information properly

	var user model.User
	user.UserInner = cmd.UserInner
	user.Identifying = cmd.Identifying
	user.Status = cmd.Status
	user.CreatedAt = cmd.CreatedAt

	var userInternal model.UserInternal

	user.PhoneVerifiedAt = cmd.PhoneVerifiedAt
	if err := x.Where("id = ?", cmd.UserID).ShouldUpdate(&user, &userInternal); err != nil {
		return err
	}

	cmd.Result.User = &user
	return x.Where("id = ?", cmd.UserID).ShouldGet(&user)
}
