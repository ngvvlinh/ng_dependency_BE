package aggregate

import (
	"context"
	"time"

	"etop.vn/capi"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/main/invitation/convert"
	"etop.vn/backend/com/main/invitation/model"
	"etop.vn/backend/com/main/invitation/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/email"
)

var _ invitation.Aggregate = &InvitationAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type InvitationAggregate struct {
	db            *cmsql.Database
	eventBus      capi.EventBus
	jwtKey        string
	store         sqlstore.InvitationStoreFactory
	customerQuery customering.QueryBus
	identityQuery identity.QueryBus
}

func NewInvitationAggregate(
	database *cmsql.Database, secretKey string,
	customerQ customering.QueryBus, identityQ identity.QueryBus,
	eventBus capi.EventBus,
) *InvitationAggregate {
	return &InvitationAggregate{
		db:            database,
		eventBus:      eventBus,
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
	if !a.checkRoles(args.Roles) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "role không hợp lệ")
	}

	if err := a.havePermissionToInvite(ctx, args); err != nil {
		return nil, err
	}

	if err := a.checkUserBelongsToShop(ctx, args.Email, args.AccountID); err != nil {
		return nil, err
	}

	// check have invitation that was sent to email
	_, err := a.store(ctx).AccountID(args.AccountID).Email(args.Email).NotExpires(time.Now()).
		AcceptedAt(nil).RejectedAt(nil).GetInvitation()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	// no-op
	case cm.NoError:
		return nil, cm.Errorf(cm.NotFound, nil, "Email %s đã được gửi lời mời.", args.Email)
	default:
		return nil, err
	}

	invitation := new(invitation.Invitation)
	if err := scheme.Convert(args, invitation); err != nil {
		return nil, err
	}

	// generate token
	expiresAt := time.Now().Add(convert.ExpiresIn)
	token, err := convert.GenerateToken(a.jwtKey, args.Email, args.AccountID, args.Roles, expiresAt.Unix())
	if err != nil {
		return nil, err
	}
	invitation.ExpiresAt = expiresAt
	invitation.Token = token

	if err := a.store(ctx).CreateInvitation(invitation); err != nil {
		return nil, err
	}

	// send mail
	// TODO: change content and subject
	cmd := &email.SendEmailCommand{
		FromName:    "eTop.vn (no-reply)",
		ToAddresses: []string{args.Email},
		Subject:     "Invitation",
		Content:     token,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (a *InvitationAggregate) checkStatusInvitation(invitation *model.Invitation) error {
	if !invitation.ExpiresAt.After(time.Now()) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã hết hạn")
	}
	if !invitation.AcceptedAt.IsZero() || invitation.Status == etop.S3Positive {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã được chấp nhận")
	}
	if !invitation.RejectedAt.IsZero() || invitation.Status == etop.S3Negative {
		return cm.Errorf(cm.FailedPrecondition, nil, "Lời mời đã được từ chối")
	}
	return nil
}

func (a *InvitationAggregate) checkUserBelongsToShop(ctx context.Context, email string, shopID int64) error {
	getUserQuery := &identity.GetUserByEmailQuery{
		Email: email,
	}
	err := a.identityQuery.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return nil
	case cm.NoError:
	// no-op
	default:
		return err
	}
	userIsInvited := getUserQuery.Result

	getAccountUserQuery := &etopmodel.GetAccountUserQuery{
		UserID:    userIsInvited.ID,
		AccountID: shopID,
	}
	err = bus.Dispatch(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return nil
	case cm.NoError:
		return cm.Errorf(cm.FailedPrecondition, nil, "Người dùng đã là thành viên của shop")
	default:
		return err
	}
}

func (a *InvitationAggregate) havePermissionToInvite(ctx context.Context, args *invitation.CreateInvitationArgs) error {
	getAccountUserQuery := &etopmodel.GetAccountUserQuery{
		UserID:          args.InvitedBy,
		AccountID:       args.AccountID,
		FindByAccountID: true,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return err
	}
	customerRoles := convert.ConvertStringsToRoles(getAccountUserQuery.Result.Roles)
	if !invitation.IsContainsRole(customerRoles, invitation.RoleShopOwner) &&
		!invitation.IsContainsRole(customerRoles, invitation.RoleStaffManagement) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Người dùng hiện tại không có quyền mời người khác vào shop")
	}
	return nil
}

func (a *InvitationAggregate) checkRoles(roles []invitation.Role) bool {
	for _, role := range roles {
		if string(role) == string(invitation.RoleShopOwner) {
			return false
		}
		if !invitation.IsRole(role) {
			return false
		}
	}
	return true
}

func (a *InvitationAggregate) AcceptInvitation(
	ctx context.Context, userID int64, token string,
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

	if err := a.checkTokenBelongsToUser(ctx, userID, invitationDB.Email); err != nil {
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
	ctx context.Context, userID int64, token string,
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

	if err := a.checkTokenBelongsToUser(ctx, userID, invitationDB.Email); err != nil {
		return 0, err
	}

	updated, err := a.store(ctx).Token(token).Reject()
	if err != nil {
		return 0, err
	}
	return updated, err
}

func (a *InvitationAggregate) checkTokenBelongsToUser(ctx context.Context, userID int64, email string) error {
	getUserQuery := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := a.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "Không tìm thấy user").
			Throw()
	}
	userCurr := getUserQuery.Result
	if userCurr.EmailVerifiedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Email bạn chưa được xác nhận nên, không thể thực hiện thao tác này.")
	}
	if userCurr.Email != email {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không tìm thấy lời mời")
	}
	return nil
}
