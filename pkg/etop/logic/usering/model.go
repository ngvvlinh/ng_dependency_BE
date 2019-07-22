package usering

import (
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/bus"
)

type CreateUserCommand model.CreateUserCommand

func (m *CreateUserCommand) Validate() []error {
	return nil
}

func init() {
	bus.AddHandlers("usering",
		CreateUser)
}
