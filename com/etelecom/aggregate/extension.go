package aggregate

import (
	"context"

	"o.o/api/etelecom"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (a *EtelecomAggregate) CreateExtension(ctx context.Context, args *etelecom.CreateExtensionArgs) (*etelecom.Extension, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	event := &etelecom.ExtensionCreatingEvent{
		UserID:    args.UserID,
		AccountID: args.AccountID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	ext, err := a.extensionStore(ctx).UserID(args.UserID).AccountID(args.AccountID).GetExtension()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		if ext.ExtensionNumber != "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Extension đã được tạo cho người dùng này.")
		}
	case cm.NotFound:
		// create new one
		var extension etelecom.Extension
		if err := scheme.Convert(args, &extension); err != nil {
			return nil, err
		}

		extension.ID = cm.NewID()
		ext, err = a.extensionStore(ctx).CreateExtension(&extension)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return ext, nil
}

func (a *EtelecomAggregate) DeleteExtension(ctx context.Context, id dot.ID) error {
	_, err := a.extensionStore(ctx).ID(id).SoftDelete()
	return err
}
