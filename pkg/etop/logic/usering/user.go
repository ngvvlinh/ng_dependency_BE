package usering

import (
	"context"

	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	"etop.vn/backend/pkg/common/bus"
)

func CreateUser(ctx context.Context, cmd *CreateUserCommand) error {

	createUserCmd := (*identitymodelx.CreateUserCommand)(cmd)
	if err := bus.Dispatch(ctx, createUserCmd); err != nil {
		return err
	}
	return nil
}
