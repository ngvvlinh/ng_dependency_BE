package pm

import (
	"context"

	"etop.vn/api/top/types/etc/status3"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	invitationQuery invitation.QueryBus
	invitationAggr  invitation.CommandBus
}

func New(
	eventBus capi.EventBus,
	invitationQ invitation.QueryBus,
	invitationA invitation.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:        eventBus,
		invitationQuery: invitationQ,
		invitationAggr:  invitationA,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.UserCreated)
}

func (m *ProcessManager) UserCreated(ctx context.Context, event *identity.UserCreatedEvent) error {
	if event.Invitation == nil || !event.Invitation.AutoAccept {
		return nil
	}

	cmd := &invitation.AcceptInvitationCommand{
		UserID: event.UserID,
		Token:  event.Invitation.Token,
	}
	if err := m.invitationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	query := &invitation.ListInvitationsAcceptedByEmailQuery{
		Email: event.Email,
	}
	if err := m.invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	for _, invitationItem := range query.Result.Invitations {
		query := &model.GetAccountUserQuery{
			UserID:    event.UserID,
			AccountID: invitationItem.AccountID,
		}
		err := bus.Dispatch(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			var roles []string
			for _, role := range invitationItem.Roles {
				roles = append(roles, string(role))
			}
			cmd := &model.CreateAccountUserCommand{
				AccountUser: &model.AccountUser{
					AccountID: invitationItem.AccountID,
					UserID:    event.UserID,
					Status:    status3.Z,
					Permission: model.Permission{
						Roles: roles,
					},
					FullName:  event.Invitation.FullName,
					ShortName: event.Invitation.ShortName,
					Position:  event.Invitation.Position,
				},
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return err
			}
		case cm.NoError:
			return cm.Errorf(cm.Internal, nil, "unexpected (invitation exists)")
		default:
			return err
		}
	}
	return nil
}
