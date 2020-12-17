package aggregate

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/com/main/invitation/convert"
	"o.o/backend/com/main/invitation/model"
	invstore "o.o/backend/com/main/invitation/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/templatemessages"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/sms"
	"o.o/backend/pkg/integration/sms/util"
	"o.o/capi"
	"o.o/capi/dot"
)

type FlagEnableNewLinkInvitation bool

var _ invitation.Aggregate = &InvitationAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type InvitationAggregate struct {
	db            *cmsql.Database
	eventBus      capi.EventBus
	jwtKey        string
	store         invstore.InvitationStoreFactory
	customerQuery customering.QueryBus
	identityQuery identity.QueryBus
	smsClient     *sms.Client
	emailClient   *email.Client
	flagNewLink   FlagEnableNewLinkInvitation
	redisStore    redis.Store

	AccountUserStore sqlstore.AccountUserStoreInterface
	ShopStore        sqlstore.ShopStoreInterface
	UserStore        sqlstore.UserStoreInterface
}

func NewInvitationAggregate(
	database com.MainDB,
	cfg invitation.Config,
	customerQ customering.QueryBus,
	identityQ identity.QueryBus,
	eventBus capi.EventBus,
	smsClient *sms.Client,
	emailClient *email.Client,
	flagNewLink FlagEnableNewLinkInvitation,
	AccountUserStore sqlstore.AccountUserStoreInterface,
	ShopStore sqlstore.ShopStoreInterface,
	UserStore sqlstore.UserStoreInterface,
	redis redis.Store,
) *InvitationAggregate {
	return &InvitationAggregate{
		db:               database,
		eventBus:         eventBus,
		store:            invstore.NewInvitationStore(database),
		jwtKey:           cfg.Secret,
		customerQuery:    customerQ,
		identityQuery:    identityQ,
		smsClient:        smsClient,
		emailClient:      emailClient,
		flagNewLink:      flagNewLink,
		AccountUserStore: AccountUserStore,
		ShopStore:        ShopStore,
		UserStore:        UserStore,
		redisStore:       redis,
	}
}

func InvitationAggregateMessageBus(a *InvitationAggregate) invitation.CommandBus {
	b := bus.New()
	return invitation.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *InvitationAggregate) CreateInvitation(
	ctx context.Context, args *invitation.CreateInvitationArgs,
) (*invitation.Invitation, error) {
	var err error
	if args.Email, args.Phone, err = validateEmailOrPhone(args.Email, args.Phone); err != nil {
		return nil, err
	}

	if !a.checkRoles(args.Roles) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "role không hợp lệ")
	}

	_, userIsInvited, err := a.getInvitationByEmailOrPhone(ctx, args.Email, args.Phone, args.AccountID)

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

	URL, err := GetInvitationURL(ctx, invitationItem, a.flagNewLink)
	if err != nil {
		return nil, err
	}
	invitationItem.InvitationURL = URL.String()

	if err := a.db.InTransaction(ctx, func(q cmsql.QueryInterface) error {
		// create invitation
		if err := a.store(ctx).CreateInvitation(invitationItem); err != nil {
			return err
		}
		err := a.sendInvitation(ctx, invitationItem)
		return err
	}); err != nil {
		return nil, err
	}

	return invitationItem, nil
}

func validateEmailOrPhone(email, phone string) (validatedEmail string, validatedPhone string, err error) {
	var emailNorm validate.NormalizedEmail
	var phoneNorm validate.NormalizedPhone

	if email != "" {
		var ok bool
		emailNorm, ok = validate.NormalizeEmail(email)
		if !ok {
			return "", "", cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
		}
		validatedEmail = emailNorm.String()
	}

	if phone != "" {
		var ok bool
		phoneNorm, ok = validate.NormalizePhone(phone)
		if !ok {
			return "", "", cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
		}
		validatedPhone = phoneNorm.String()
	}

	return
}

func (a *InvitationAggregate) getInvitationByEmailOrPhone(ctx context.Context, email string, phone string, accountID dot.ID) (*invitation.Invitation, *identity.User, error) {
	var err error
	if email, phone, err = validateEmailOrPhone(email, phone); err != nil {
		return nil, nil, err
	}

	userIsInvited, err := a.checkUserBelongsToShop(ctx, email, phone, accountID)
	if err != nil {
		return nil, nil, err
	}

	// check have invitation that was sent to email
	query := a.store(ctx).AccountID(accountID).NotExpires().AcceptedAt(nil).RejectedAt(nil)
	if userIsInvited != nil {
		query = query.PhoneOrEmail(userIsInvited.Phone, userIsInvited.Email)
	} else {
		if phone != "" {
			query = query.Phone(phone)
		}
		if email != "" {
			query = query.Email(email)
		}
	}
	result, err := query.GetInvitation()
	return result, userIsInvited, err
}

func (a *InvitationAggregate) ResendInvitation(ctx context.Context, args *invitation.ResendInvitationArgs) (*invitation.Invitation, error) {
	var err error
	if args.Email, args.Phone, err = validateEmailOrPhone(args.Email, args.Phone); err != nil {
		return nil, err
	}
	invitationCore, _, err := a.getInvitationByEmailOrPhone(ctx, args.Email, args.Phone, args.AccountID)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		if args.Phone != "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại %s chưa được gửi lời mời.", args.Phone)
		}
		if args.Email != "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Email %s chưa được gửi lời mời.", args.Email)
		}
	case cm.NoError:
		// no-op
	default:
		return nil, err
	}

	URL, err := GetInvitationURL(ctx, invitationCore, a.flagNewLink)
	if err != nil {
		return nil, err
	}
	invitationCore.InvitationURL = URL.String()
	err = a.sendInvitation(ctx, invitationCore)
	if err != nil {
		return nil, err
	}
	return invitationCore, nil
}

func (a *InvitationAggregate) sendInvitation(ctx context.Context, args *invitation.Invitation) error {
	getUserQuery := &identitymodelx.GetUserByIDQuery{
		UserID: args.InvitedBy,
	}
	if err := a.UserStore.GetUserByID(ctx, getUserQuery); err != nil {
		return err
	}

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: args.AccountID,
	}
	if err := a.ShopStore.GetShop(ctx, getAccountQuery); err != nil {
		return err
	}
	fullName := "bạn"
	if args.FullName != "" {
		fullName = args.FullName
	}
	shopRoles := strings.Join(authorization.ParseRoleLabels(args.Roles), ", ")
	shopName := getAccountQuery.Result.Name
	invitedUsername := getUserQuery.Result.FullName
	var b strings.Builder
	if args.Email != "" {
		if err := templatemessages.EmailInvitationTpl.Execute(&b, map[string]interface{}{
			"FullName":        fullName,
			"URL":             args.InvitationURL,
			"ShopRoles":       shopRoles,
			"ShopName":        shopName,
			"InvitedUsername": invitedUsername,
			"WlName":          wl.X(ctx).Name,
		}); err != nil {
			return cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ email").WithMeta("reason", "can not generate email content")
		}
		cmd := &email.SendEmailCommand{
			FromName:    wl.X(ctx).CompanyName + " (no-reply)",
			ToAddresses: []string{args.Email},
			Subject:     "Invitation",
			Content:     b.String(),
		}
		if err := a.emailClient.SendMail(ctx, cmd); err != nil {
			return err
		}
	} else {
		var countSend int
		var redisKey = fmt.Sprintf("Invitation-sendphone-%v-%v", args.AccountID, args.Phone)
		err := a.redisStore.Get(redisKey, &countSend)
		if err != nil && err != redis.ErrNil {
			return err
		}
		if err != nil && err == redis.ErrNil {
			countSend = 1
		}
		if countSend == 1 {
			if err := templatemessages.PhoneInvitationTpl.Execute(&b, map[string]interface{}{
				"URL":      args.InvitationURL,
				"ShopName": util.ModifyMsgPhone(shopName),
			}); err != nil {
				return cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ phone").WithMeta("reason", "can not generate phone content")
			}
		} else {
			if err := templatemessages.PhoneInvitationTplRepeat.Execute(&b, map[string]interface{}{
				"URL":      args.InvitationURL,
				"ShopName": util.ModifyMsgPhone(shopName),
				"SendTime": countSend,
			}); err != nil {
				return cm.Errorf(cm.Internal, err, "Không thể xác nhận địa chỉ phone").WithMeta("reason", "can not generate phone content")
			}
		}
		cmd := &sms.SendSMSCommand{
			Phone:   args.Phone,
			Content: b.String(),
		}
		if err := a.smsClient.SendSMS(ctx, cmd); err != nil {
			return err
		}

		// update countSend
		countSend += 1
		if err := a.redisStore.Set(redisKey, countSend); err != nil {
			return err
		}
	}
	return nil
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
	err = a.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return userIsInvited, nil
	case cm.NoError:
		return userIsInvited, cm.Errorf(cm.FailedPrecondition, nil, "Người dùng đã là thành viên của shop")
	default:
		return userIsInvited, err
	}
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
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getAccountUserQuery); err != nil {
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

func GetInvitationURL(ctx context.Context, args *invitation.Invitation, flag FlagEnableNewLinkInvitation) (*url.URL, error) {
	var invitationUrl string
	if isRoleForTelecom(ctx, args) {
		invitationUrl = getInvitationURLForTelecom(ctx, args, flag)
	} else {
		invitationUrl = getInvitationURL(ctx, args, flag)
	}

	URL, err := url.Parse(invitationUrl)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can not parse url")
	}
	urlQuery := URL.Query()
	if args.Email != "" || (args.Phone != "" && !flag) {
		urlQuery.Set("t", args.Token)
	}
	URL.RawQuery = urlQuery.Encode()
	return URL, nil
}

func getInvitationURLForTelecom(ctx context.Context, args *invitation.Invitation, flag FlagEnableNewLinkInvitation) string {
	wlPartner := wl.X(ctx)
	if args.Email != "" || (args.Phone != "" && !flag) {
		return wlPartner.Config.GetTelecomInviteUserURLByEmail(wlPartner.TelecomURL)
	}
	return wlPartner.Config.GetTelecomInviteUserURLByPhone(wlPartner.TelecomURL, args.Phone)
}

func getInvitationURL(ctx context.Context, args *invitation.Invitation, flag FlagEnableNewLinkInvitation) string {
	wlPartner := wl.X(ctx)
	if args.Email != "" || (args.Phone != "" && !flag) {
		return wlPartner.InviteUserURLByEmail
	}
	return wlPartner.InviteUserURLByPhone + "/p" + args.Phone
}

func isRoleForTelecom(ctx context.Context, args *invitation.Invitation) bool {
	return authorization.IsContainsRole(args.Roles, authorization.RoleTelecomCustomerService)
}
