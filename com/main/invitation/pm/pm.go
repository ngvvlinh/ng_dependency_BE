package pm

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
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
	return nil
}
