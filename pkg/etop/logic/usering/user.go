package usering

import (
	"context"

	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
)

func CreateUser(ctx context.Context, cmd *CreateUserCommand) error {

	createUserCmd := (*identitymodelx.CreateUserCommand)(cmd)
	if err := bus.Dispatch(ctx, createUserCmd); err != nil {
		return err
	}
	return nil
}
