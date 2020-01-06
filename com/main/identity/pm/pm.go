package pm

import (
	"context"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/backend/com/main/authorization/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
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

	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
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
