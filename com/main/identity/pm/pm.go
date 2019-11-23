package pm

import (
	"context"
	"time"

	"etop.vn/backend/com/main/authorization/convert"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

type ProcessManager struct {
	identityQuery   identity.QueryBus
	invitationQuery *invitation.QueryBus
}

func New(
	identityQ identity.QueryBus, invitationQ *invitation.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		identityQuery:   identityQ,
		invitationQuery: invitationQ,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvitationAccepted)
}

func (m *ProcessManager) InvitationAccepted(ctx context.Context, event *invitation.InvitationAcceptedEvent) error {
	getInvitationQuery := &invitation.GetInvitationQuery{
		ID: event.ID,
	}
	if err := m.invitationQuery.Dispatch(ctx, getInvitationQuery); err != nil {
		return err
	}
	currInvitation := getInvitationQuery.Result

	getUserQuery := &identity.GetUserByEmailQuery{
		Email: currInvitation.Email,
	}
	err := m.identityQuery.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// no-op
	case cm.NotFound:
		return nil
	default:
		return err
	}

	currUser := getUserQuery.Result
	if currUser.EmailVerifiedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thao tác thất bại, email chưa được xác nhận")
	}

	getAccountUserQuery := &model.GetAccountUserQuery{
		UserID:    currUser.ID,
		AccountID: currInvitation.AccountID,
	}
	err = bus.Dispatch(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		err := m.createAccountUserWithRoles(ctx, currInvitation, currUser)
		if err != nil {
			return err
		}
	case cm.NoError:
		err := m.updateAccountUserWithRoles(ctx, getAccountUserQuery.Result, currInvitation)
		if err != nil {
			return err
		}
	default:
		return err
	}

	return nil
}

func (m *ProcessManager) updateAccountUserWithRoles(
	ctx context.Context, accountUser *model.AccountUser, currInvitation *invitation.Invitation,
) error {
	mapRole := make(map[string]bool)
	roles := []string{}
	for _, role := range accountUser.Roles {
		mapRole[role] = true
	}
	for _, role := range currInvitation.Roles {
		mapRole[string(role)] = true
	}
	for role := range mapRole {
		roles = append(roles, role)
	}
	accountUser.Roles = roles
	accountUser.Permission.Roles = roles
	updateAccountUserCmd := &model.UpdateAccountUserCommand{
		AccountUser: accountUser,
	}
	if err := bus.Dispatch(ctx, updateAccountUserCmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) createAccountUserWithRoles(
	ctx context.Context, currInvitation *invitation.Invitation, currUser *identity.User,
) error {
	createAccountUserCmd := &model.CreateAccountUserCommand{
		AccountUser: &model.AccountUser{
			AccountID: currInvitation.AccountID,
			UserID:    currUser.ID,
			Status:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Permission: model.Permission{
				Roles: convert.ConvertRolesToStrings(currInvitation.Roles),
			},
			FullName:  currUser.FullName,
			ShortName: currUser.ShortName,
		},
	}
	if err := bus.Dispatch(ctx, createAccountUserCmd); err != nil {
		return err
	}
	return nil
}
