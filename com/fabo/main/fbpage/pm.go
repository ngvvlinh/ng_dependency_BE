package fbpage

import (
	"golang.org/x/net/context"
	"o.o/api/fabo/fbpaging"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	eventBus   bus.EventRegistry
	fbPageUtil *FbPageUtil
}

func NewProcessManager(
	eventBus bus.EventRegistry,
	fbPageUtil *FbPageUtil,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:   eventBus,
		fbPageUtil: fbPageUtil,
	}
	p.RegisterEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.HandleFbExternalPagesCreatedOrUpdatedEvent)
}

func (m *ProcessManager) HandleFbExternalPagesCreatedOrUpdatedEvent(
	ctx context.Context, event *fbpaging.FbExternalPagesCreatedOrUpdatedEvent,
) error {
	if err := m.fbPageUtil.ClearFbPages(event.ExternalPageIDs...); err != nil {
		return err
	}

	return m.fbPageUtil.ClearFbPageInternals(event.ExternalPageIDs...)
}
