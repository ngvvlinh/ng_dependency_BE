package usering

import (
	"context"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

func CreateUser(ctx context.Context, cmd *CreateUserCommand) error {

	createUserCmd := (*model.CreateUserCommand)(cmd)
	if err := bus.Dispatch(ctx, createUserCmd); err != nil {
		return err
	}
	return nil
}
