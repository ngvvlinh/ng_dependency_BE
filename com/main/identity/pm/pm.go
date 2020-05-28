package pm

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/backend/com/main/authorization/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	identityQuery   identity.QueryBus
	invitationQuery invitation.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	identityQ identity.QueryBus,
	invitationQ invitation.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		identityQuery:   identityQ,
		invitationQuery: invitationQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
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

	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Email: currInvitation.Email,
		Phone: currInvitation.Phone,
	}
	if err := m.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return err
	}
	currUser := getUserQuery.Result
	userID := currUser.ID

	if currInvitation.Email != "" {
		if currUser.EmailVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Thao tác thất bại, email chưa được xác nhận")
		}
	} else {
		if currUser.PhoneVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Thao tác thất bại, phone chưa được xác nhận")
		}
	}

	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    userID,
		AccountID: currInvitation.AccountID,
	}
	err := bus.Dispatch(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		err := m.createAccountUserWithRoles(ctx, currInvitation, currUser)
		if err != nil {
			return err
		}
	case cm.NoError:
		return cm.Errorf(cm.Internal, nil, "unexpected (account_user exists)")
	default:
		return err
	}
	return nil
}

func (m *ProcessManager) createAccountUserWithRoles(
	ctx context.Context, currInvitation *invitation.Invitation, currUser *identity.User,
) error {
	createAccountUserCmd := &identitymodelx.CreateAccountUserCommand{
		AccountUser: &identitymodel.AccountUser{
			AccountID: currInvitation.AccountID,
			UserID:    currUser.ID,
			Status:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Permission: identitymodel.Permission{
				Roles: convert.ConvertRolesToStrings(currInvitation.Roles),
			},
			FullName:  cm.Coalesce(currInvitation.FullName, currUser.FullName),
			ShortName: cm.Coalesce(currInvitation.ShortName, currUser.ShortName),
			Position:  currInvitation.Position,
		},
	}
	if err := bus.Dispatch(ctx, createAccountUserCmd); err != nil {
		return err
	}
	return nil
}
