package usering

import (
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	"etop.vn/backend/pkg/common/bus"
)

type CreateUserCommand identitymodelx.CreateUserCommand

func (m *CreateUserCommand) Validate() []error {
	return nil
}

func init() {
	bus.AddHandlers("usering",
		CreateUser)
}
