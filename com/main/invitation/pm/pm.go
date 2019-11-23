package pm

import (
	"context"

	"etop.vn/capi/dot"

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
	if event.AutoAcceptInvitation {
		cmd := &invitation.AcceptInvitationCommand{
			UserID: event.UserID,
			Token:  event.InvitationToken,
		}
		if err := m.invitationAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}

	query := &invitation.ListInvitationsAcceptedByEmailQuery{
		Email: event.Email,
	}
	if err := m.invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	for _, invitation := range query.Result.Invitations {
		query := &model.GetAccountUserQuery{
			UserID:    event.UserID,
			AccountID: invitation.AccountID,
		}
		err := bus.Dispatch(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			var roles []string
			for _, role := range invitation.Roles {
				roles = append(roles, string(role))
			}

			cmd := &model.CreateAccountUserCommand{
				AccountUser: &model.AccountUser{
					AccountID: invitation.AccountID,
					UserID:    event.UserID,
					Status:    model.S3Zero,
					Permission: model.Permission{
						Roles: roles,
					},
					FullName:  event.FullName,
					ShortName: event.ShortName,
				},
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return err
			}
		case cm.NoError:
			accountUser := query.Result
			mapRole := make(map[string]bool)
			for _, role := range accountUser.Permission.Roles {
				mapRole[role] = true
			}
			for _, role := range invitation.Roles {
				mapRole[string(role)] = true
			}

			var roles []string
			for key := range mapRole {
				roles = append(roles, key)
			}

			cmd := &model.UpdateRoleCommand{
				AccountID: invitation.AccountID,
				UserID:    event.UserID,
				Permission: model.Permission{
					Roles:       roles,
					Permissions: accountUser.Permission.Permissions,
				},
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				return err
			}

			updateInfosCmd := &model.UpdateInfosCommand{
				AccountID: invitation.AccountID,
				UserID:    event.UserID,
				FullName:  dot.String(event.FullName),
				ShortName: dot.String(event.ShortName),
				Position:  dot.String(event.Position),
			}
			if err := bus.Dispatch(ctx, updateInfosCmd); err != nil {
				return err
			}
		default:
			return err
		}
	}

	return nil
}
