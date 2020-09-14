package sqlstore

import (
	"context"
	"strings"
	"time"

	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	identitysqlstore "o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

type UserStoreInterface interface {
	CreateUser(ctx context.Context, cmd *identitymodelx.CreateUserCommand) error

	GetSignedInUser(ctx context.Context, query *identitymodelx.GetSignedInUserQuery) error

	GetUserByEmail(ctx context.Context, query *identitymodelx.GetUserByEmailOrPhoneQuery) error

	GetUserByID(ctx context.Context, query *identitymodelx.GetUserByIDQuery) error

	GetUserByLogin(ctx context.Context, query *identitymodelx.GetUserByLoginQuery) error

	SetPassword(ctx context.Context, cmd *identitymodelx.SetPasswordCommand) error

	UpdateUserVerification(ctx context.Context, cmd *identitymodelx.UpdateUserVerificationCommand) error
}

type UserStore struct {
	db    *cmsql.Database
	query cmsql.QueryFactory
	ft    identitysqlstore.UserFilters
	preds []interface{}
}

type UserStoreFactory func(ctx context.Context) *UserStore

func NewUserStore(db com.MainDB) UserStoreFactory {
	return func(ctx context.Context) *UserStore {
		s := &UserStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
		return s
	}
}

func BuildUserStore(db com.MainDB) *UserStore {
	return NewUserStore(db)(context.Background())
}

func BindUserStore(s *UserStore) (to UserStoreInterface) {
	return s
}

func (st *UserStore) ID(id dot.ID) *UserStore {
	st.preds = append(st.preds, st.ft.ByID(id))
	return st
}

func (st *UserStore) IDs(ids ...dot.ID) *UserStore {
	st.preds = append(st.preds, sq.In("id", ids))
	return st
}

func (st *UserStore) Get() (*identitymodel.User, error) {
	var user identitymodel.User
	err := st.query().Where(st.preds...).ShouldGet(&user)
	return &user, err
}

func (st *UserStore) List() ([]*identitymodel.User, error) {
	var users identitymodel.Users
	err := st.query().Where(st.preds).Find(&users)
	return users, err
}

func (st *UserStore) GetSignedInUser(ctx context.Context, query *identitymodelx.GetSignedInUserQuery) error {
	userQuery := &identitymodelx.GetUserByIDQuery{
		UserID: query.UserID,
	}
	if err := st.GetUserByID(ctx, userQuery); err != nil {
		return err
	}
	query.Result = &identitymodelx.SignedInUser{
		User: userQuery.Result,
	}
	return nil
}

func (st *UserStore) GetUserByID(ctx context.Context, query *identitymodelx.GetUserByIDQuery) error {
	if query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	query.Result = new(identitymodel.User)
	s := st.query().Where("id = ?", query.UserID)
	s = FilterByWhiteLabelPartner(s, wl.GetWLPartnerID(ctx))
	return s.ShouldGet(query.Result)
}

func (st *UserStore) GetUserByEmail(ctx context.Context, query *identitymodelx.GetUserByEmailOrPhoneQuery) error {
	count := 0
	q := st.query().Table("user")
	if query.Email != "" {
		q = q.Where("email = ?", query.Email)
		count++
	}
	if query.Phone != "" {
		q = q.Where("phone = ?", query.Phone)
		count++
	}
	q = FilterByWhiteLabelPartner(q, wl.GetWLPartnerID(ctx))
	if count != 1 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}
	query.Result = new(identitymodel.User)
	return q.ShouldGet(query.Result)
}

func (st *UserStore) GetUserByLogin(ctx context.Context, query *identitymodelx.GetUserByLoginQuery) error {
	if query.UserID == 0 && query.PhoneOrEmail == "" {
		return cm.Error(cm.InvalidArgument, "Missing required fields", nil)
	}

	s := st.query().Table("user")
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
	s = FilterByWhiteLabelPartner(s, wl.GetWLPartnerID(ctx))

	user := new(identitymodel.User)
	if err := s.ShouldGet(user); err != nil {
		return err
	}

	userInternal := new(identitymodel.UserInternal)
	if err := st.query().Where("id = ?", user.ID).ShouldGet(userInternal); err != nil {
		return err
	}
	query.Result.User = user
	query.Result.UserInternal = userInternal
	return nil
}

func (st *UserStore) CreateUser(ctx context.Context, cmd *identitymodelx.CreateUserCommand) error {
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

	return inTransaction(st.db, func(tx Qx) error {
		return st.createUser(ctx, tx, cmd)
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
	UserEmailKey            = "user_email_key"
	UserPhoneKey            = "user_phone_key"
	UserPhoneWLPartnerIDKey = "user_phone_wl_partner_id_idx"
	UserEmailWLPartnerIDKey = "user_email_wl_partner_id_idx"

	MsgCreateUserDuplicatedPhone = `Số điện thoại đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng số điện thoại khác. Nếu cần thêm thông tin, vui lòng liên hệ %v.`
	MsgCreateUserDuplicatedEmail = `Email đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng email khác. Nếu cần thêm thông tin, vui lòng liên hệ %v.`
)

var mapUserError = map[string]string{
	UserEmailKey:            MsgCreateUserDuplicatedEmail,
	UserEmailWLPartnerIDKey: MsgCreateUserDuplicatedEmail,
	UserPhoneKey:            MsgCreateUserDuplicatedPhone,
	UserPhoneWLPartnerIDKey: MsgCreateUserDuplicatedPhone,
}

func (st *UserStore) createUser(ctx context.Context, tx Qx, cmd *identitymodelx.CreateUserCommand) error {
	switch cmd.Status {
	case status3.P:
	default:
		return cm.Error(cm.InvalidArgument, "Invalid status", nil)
	}

	now := time.Now()
	userID := cm.NewIDWithTag(account_tag.TagUser)
	user := &identitymodel.User{
		ID:          userID,
		UserInner:   cmd.UserInner,
		Status:      cmd.Status,
		Source:      cmd.Source,
		AgreedTOSAt: now,
		WLPartnerID: wl.GetWLPartnerID(ctx),
	}
	user.FullNameNorm = validate.NormalizeSearch(user.FullName)
	if cmd.AgreeEmailInfo {
		user.AgreedEmailInfoAt = now
	}
	if cmd.IsTest {
		user.IsTest = 1
	}

	userInternal := &identitymodel.UserInternal{
		ID: userID,
	}
	if cmd.Password != "" {
		userInternal.Hashpwd = EncodePassword(cmd.Password)
	}

	objs := []core.IInsert{user, userInternal}
	if cmd.RefAff != "" || cmd.RefSale != "" {
		userRef := &identitymodel.UserRefSaff{
			UserID:  userID,
			RefAff:  cmd.RefAff,
			RefSale: cmd.RefSale,
		}
		objs = append(objs, userRef)
	}

	_, err := tx.Insert(objs...)
	if xerr, ok := err.(*xerrors.APIError); ok && xerr.Err != nil {
		msg := xerr.Err.Error()
		for errKey, errMsg := range mapUserError {
			if strings.Contains(msg, errKey) {
				err = cm.Errorf(cm.FailedPrecondition, nil, errMsg, wl.X(ctx).CSEmail)
			}
		}
	}
	cmd.Result.User = user
	cmd.Result.UserInternal = userInternal
	return err
}

func (st *UserStore) SetPassword(ctx context.Context, cmd *identitymodelx.SetPasswordCommand) error {
	if cmd.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	userInternal := &identitymodel.UserInternal{
		Hashpwd: EncodePassword(cmd.Password),
	}
	if err := st.query().Where("id = ?", cmd.UserID).
		ShouldUpdate(userInternal); err != nil {
		return err
	}
	return nil
}

func (st *UserStore) UpdateUserVerification(ctx context.Context, cmd *identitymodelx.UpdateUserVerificationCommand) error {
	if cmd.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	count := 0
	var user identitymodel.User

	s := st.query().Where("id = ?", cmd.UserID)
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
	s = FilterByWhiteLabelPartner(s, wl.GetWLPartnerID(ctx))

	return s.ShouldUpdate(&user)
}

func (st *UserStore) UpdateUserIdentifier(ctx context.Context, cmd *identitymodelx.UpdateUserIdentifierCommand) error {
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

	var user identitymodel.User
	user.UserInner = cmd.UserInner
	user.FullNameNorm = validate.NormalizeSearch(cmd.FullName)
	user.Status = cmd.Status
	user.CreatedAt = cmd.CreatedAt

	var userInternal identitymodel.UserInternal

	user.PhoneVerifiedAt = cmd.PhoneVerifiedAt
	s := st.query().Where("id = ?", cmd.UserID)
	s = FilterByWhiteLabelPartner(s, wl.GetWLPartnerID(ctx))
	if err := st.query().Where("id = ?", cmd.UserID).ShouldUpdate(&user, &userInternal); err != nil {
		return err
	}

	cmd.Result.User = &user
	return st.query().Where("id = ?", cmd.UserID).ShouldGet(&user)
}

var userFt = identitysqlstore.UserFilters{}

func FilterByWhiteLabelPartner(query cmsql.Query, wlPartnerID dot.ID) cmsql.Query {
	if wlPartnerID != 0 {
		return query.Where(userFt.ByWLPartnerID(wlPartnerID))
	}
	return query.Where(userFt.NotBelongWLPartner())
}
