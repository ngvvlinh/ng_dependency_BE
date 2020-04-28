package usering

import (
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
)

type CreateUserCommand identitymodelx.CreateUserCommand

func (m *CreateUserCommand) Validate() []error {
	return nil
}

func init() {
	bus.AddHandlers("usering",
		CreateUser)
}
