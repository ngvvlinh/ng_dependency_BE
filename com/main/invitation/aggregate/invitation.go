package aggregate

import (
	"context"
	"net/url"
	"strings"
	"time"

	"etop.vn/api/main/authorization"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/cmd/etop-server/config"
	authorizationconvert "etop.vn/backend/com/main/authorization/convert"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	"etop.vn/backend/com/main/invitation/convert"
	"etop.vn/backend/com/main/invitation/model"
	"etop.vn/backend/com/main/invitation/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ invitation.Aggregate = &InvitationAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type InvitationAggregate struct {
	db            *cmsql.Database
	eventBus      capi.EventBus
	cfg           config.Config
	jwtKey        string
	store         sqlstore.InvitationStoreFactory
	customerQuery customering.QueryBus
	identityQuery identity.QueryBus
}

func NewInvitationAggregate(
	database *cmsql.Database, secretKey string,
	customerQ customering.QueryBus, identityQ identity.QueryBus,
	eventBus capi.EventBus, config config.Config,
) *InvitationAggregate {
	return &InvitationAggregate{
		db:            database,
		eventBus:      eventBus,
		cfg:           config,
		store:         sqlstore.NewInvitationStore(database),
		jwtKey:        secretKey,
		customerQuery: customerQ,
		identityQuery: identityQ,
	}
}

func (a *InvitationAggregate) MessageBus() invitation.CommandBus {
	b := bus.New()
	return invitation.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *InvitationAggregate) CreateInvitation(
	ctx context.Context, args *invitation.CreateInvitationArgs,
) (*invitation.Invitation, error) {
	var emailNorm validate.NormalizedEmail
	if args.Email != "" {
		var ok bool
		emailNorm, ok = validate.NormalizeEmail(args.Email)
		if !ok {
			return nil, cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
		args.Email = emailNorm.String()
	}

	if !a.checkRoles(args.Roles) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "role không hợp lệ")
	}
	if err := a.havePermissionToInvite(ctx, args); err != nil {
		return nil, err
	}
	userIsInvited, err := a.checkUserBelongsToShop(ctx, args.Email, args.Phone, args.AccountID)
	if err != nil {
		return nil, err
	}

	// check have invitation that was sent to email
	query := a.store(ctx).AccountID(args.AccountID).NotExpires().AcceptedAt(nil).RejectedAt(nil)
	if userIsInvited != nil {
		query = query.PhoneOrEmail(userIsInvited.Phone, userIsInvited.Email)
	} else {
		if args.Phone != "" {
			query = query.Phone(args.Phone)
		}
		if args.Email != "" {
			query = query.Email(args.Email)
		}
	}
	_, err = query.GetInvitation()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	// no-op
	case cm.NoError:
		if userIsInvited == nil {
			if args.Phone != "" {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại %s đã được gửi lời mời.", args.Phone)
			}
			if args.Email != "" {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Email %s đã được gửi lời mời.", args.Email)
			}
		} else {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Tài khoản đã được gửi lời mời.")
		}
	default:
		return nil, err
	}

	invitationItem := new(invitation.Invitation)
	if err := scheme.Convert(args, invitationItem); err != nil {
		return nil, err
	}

	// generate token
	expiresAt := time.Now().Add(convert.ExpiresIn)

	token := auth.UsageInviteUser + ":" + auth.RandomToken(auth.DefaultTokenLength)
	invitationItem.ExpiresAt = expiresAt
	invitationItem.Token = token

	getUserQuery := &identitymodelx.GetUserByIDQuery{
		UserID: invitationItem.InvitedBy,
	}
	if err := bus.Dispatch(ctx, getUserQuery); err != nil {
		return nil, err
	}

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: invitationItem.AccountID,
	}
	if err := bus.Dispatch(ctx, getAccountQuery); err != nil {
		return nil, err
	}

	var invitationUrl string
	if args.Email != "" || (args.Phone != "" && !a.cfg.FlagEnableNewLinkInvitation) {
		invitationUrl = wl.X(ctx).InviteUserURLByEmail
	} else {
		// format url: https://etop.vn/i/p000000
		invitationUrl = wl.X(ctx).InviteUserURLByPhone + "/p" + args.Phone
	}

	URL, err := url.Parse(invitationUrl)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := URL.Query()
	if args.Email != "" || (args.Phone != "" && !a.cfg.FlagEnableNewLinkInvitation) {
		urlQuery.Set("t", token)
	}
	URL.RawQuery = urlQuery.Encode()

	fullName := "bạn"
	if args.FullName != "" {
		fullName = args.FullName
	}
	shopRoles := strings.Join(authorization.ParseRoleLabels(invitationItem.Roles), ", ")
	shopName := getAccountQuery.Result.Name
	invitedUsername := getUserQuery.Result.FullName

	if err := a.db.InTransaction(ctx, func(q cmsql.QueryInterface) error {
		// create invitation
		if err := a.store(ctx).CreateInvitation(invitationItem); err != nil {
			return err
		}

		var b strings.Builder
		if args.Email != "" {
			if err := api.EmailInvitationTpl.Execute(&b, map[string]interface{}{
				"FullName":        fullName,
				"URL":             URL.String(),
				"ShopRoles":       shopRoles,
				"ShopName":        shopName,
				"InvitedUsername": invitedUsername,
			}); err != nil {
				return cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate email content")
			}
			cmd := &email.SendEmailCommand{
				FromName:    "eTop.vn (no-reply)",
				ToAddresses: []string{string(emailNorm)},
				Subject:     "Invitation",
				Content:     b.String(),
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return err
			}
		} else {
			if err := api.PhoneInvitationTpl.Execute(&b, map[string]interface{}{
				"InvitedUsername": invitedUsername,
				"URL":             URL.String(),
				"ShopRoles":       shopRoles,
				"ShopName":        shopName,
			}); err != nil {
				return cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ phone").WithMeta("reason", "can not generate phone content")
			}
			cmd := &sms.SendSMSCommand{
				Phone:   args.Phone,
				Content: b.String(),
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return invitationItem, nil
}

func (a *InvitationAggregate) checkStatusInvitation(invitation *model.Invitation) error {
	if !invitation.ExpiresAt.After(time.Now()) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã hết hạn")
	}
	if !invitation.AcceptedAt.IsZero() || invitation.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã được chấp nhận")
	}
	if !invitation.RejectedAt.IsZero() || invitation.Status == status3.N {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã được từ chối")
	}
	return nil
}

func (a *InvitationAggregate) checkUserBelongsToShop(ctx context.Context, email, phone string, shopID dot.ID) (user *identity.User, _ error) {
	var userIsInvited *identity.User
	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Phone: phone,
		Email: email,
	}
	err := a.identityQuery.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return nil, nil
	case cm.NoError:
	// no-op
	default:
		return nil, err
	}
	userIsInvited = getUserQuery.Result

	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    userIsInvited.ID,
		AccountID: shopID,
	}
	err = bus.Dispatch(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return userIsInvited, nil
	case cm.NoError:
		return userIsInvited, cm.Errorf(cm.FailedPrecondition, nil, "Người dùng đã là thành viên của shop")
	default:
		return userIsInvited, err
	}
}

func (a *InvitationAggregate) havePermissionToInvite(ctx context.Context, args *invitation.CreateInvitationArgs) error {
	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    args.InvitedBy,
		AccountID: args.AccountID,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return err
	}
	customerRoles := authorizationconvert.ConvertStringsToRoles(getAccountUserQuery.Result.Roles)
	if !authorization.IsContainsRole(customerRoles, authorization.RoleShopOwner) &&
		!authorization.IsContainsRole(customerRoles, authorization.RoleStaffManagement) &&
		!authorization.IsContainsRole(customerRoles, authorization.RoleAdmin) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Người dùng hiện tại không có quyền mời người khác vào shop")
	}
	return nil
}

func (a *InvitationAggregate) checkRoles(roles []authorization.Role) bool {
	for _, role := range roles {
		if role == authorization.RoleShopOwner || role == authorization.RoleAdmin {
			return false
		}
		if !authorization.IsRole(role) {
			return false
		}
	}
	return true
}

func (a *InvitationAggregate) AcceptInvitation(
	ctx context.Context, userID dot.ID, token string,
) (int, error) {
	invitationDB, err := a.store(ctx).Token(token).GetInvitationDB()
	if err != nil {
		return 0, cm.MapError(err).
			Wrapf(cm.NotFound, "Không tìm thấy lời mời").
			Throw()
	}
	if err := a.checkStatusInvitation(invitationDB); err != nil {
		return 0, err
	}

	if err := a.checkTokenBelongsToUser(ctx, userID, invitationDB.Email, invitationDB.Phone); err != nil {
		return 0, err
	}

	updated, err := a.store(ctx).Token(token).Accept()
	if err != nil {
		return 0, err
	}

	event := &invitation.InvitationAcceptedEvent{
		ID: invitationDB.ID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return 0, err
	}
	return updated, err
}

func (a *InvitationAggregate) RejectInvitation(
	ctx context.Context, userID dot.ID, token string,
) (int, error) {
	invitationDB, err := a.store(ctx).Token(token).GetInvitationDB()
	if err != nil {
		return 0, cm.MapError(err).
			Wrapf(cm.NotFound, "Không tìm thấy lời mời").
			Throw()
	}
	if err := a.checkStatusInvitation(invitationDB); err != nil {
		return 0, err
	}

	if err := a.checkTokenBelongsToUser(ctx, userID, invitationDB.Email, invitationDB.Phone); err != nil {
		return 0, err
	}

	updated, err := a.store(ctx).Token(token).Reject()
	if err != nil {
		return 0, err
	}
	return updated, err
}

func (a *InvitationAggregate) checkTokenBelongsToUser(ctx context.Context, userID dot.ID, email, phone string) error {
	getUserQuery := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := a.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "Không tìm thấy user").
			Throw()
	}
	userCurr := getUserQuery.Result

	if email != "" {
		if userCurr.EmailVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Email bạn chưa được xác nhận nên, không thể thực hiện thao tác này.")
		}
		if userCurr.Email != email {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không tìm thấy lời mời")
		}
	} else {
		if userCurr.PhoneVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Phone bạn chưa được xác nhận nên không thể thực hiện thao tác này.")
		}
		if userCurr.Phone != phone {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không tìm thấy lời mời")
		}
	}
	return nil
}

func (a *InvitationAggregate) DeleteInvitation(
	ctx context.Context, userID, accountID dot.ID, token string,
) (updated int, _ error) {
	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: accountID,
		UserID:    userID,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "tài khoản bạn không thuộc shop này").
			Throw()
	}
	roles := convert.ConvertStringsToRoles(getAccountUserQuery.Result.AccountUser.Permission.Roles)

	if !authorization.IsContainsRole(roles, authorization.RoleShopOwner) &&
		!authorization.IsContainsRole(roles, authorization.RoleStaffManagement) {
		return 0, cm.Error(cm.FailedPrecondition, "Bạn không có quyền thực hiện thao tác này", nil)
	}

	updated, err := a.store(ctx).AccountID(accountID).Token(token).SoftDelete()
	if err != nil {
		return 0, err
	}
	return updated, err
}
