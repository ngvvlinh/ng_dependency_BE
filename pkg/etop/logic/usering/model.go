package usering

import (
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

type CreateUserCommand model.CreateUserCommand

func (m *CreateUserCommand) Validate() []error {
	return nil
}

func init() {
	bus.AddHandlers("usering",
		CreateUser)
}
