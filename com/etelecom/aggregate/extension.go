package aggregate

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
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
	if args.ConnectionID == 0 {
		args.ConnectionID = connectioning.DefaultBuiltinVHTEtelecomConnectionID
	}

	ext, err := a.extensionStore(ctx).UserID(args.UserID).AccountID(args.AccountID).ConnectionID(args.ConnectionID).GetExtension()
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

	externalExtensionResp, err := a.telecomManager.CreateExtension(ctx, ext)
	if err != nil {
		return nil, err
	}
	updateExt := &etelecom.UpdateExternalExtensionInfoArgs{
		ID:                externalExtensionResp.ExtensionID,
		HotlineID:         externalExtensionResp.HotlineID,
		ExternalID:        externalExtensionResp.ExternalID,
		ExtensionNumber:   externalExtensionResp.ExtensionNumber,
		ExtensionPassword: externalExtensionResp.ExtensionPassword,
	}
	if err := a.UpdateExternalExtensionInfo(ctx, updateExt); err != nil {
		return nil, err
	}

	return a.extensionStore(ctx).ID(ext.ID).GetExtension()
}

func (a *EtelecomAggregate) DeleteExtension(ctx context.Context, id dot.ID) error {
	_, err := a.extensionStore(ctx).ID(id).SoftDelete()
	return err
}

func (a *EtelecomAggregate) UpdateExternalExtensionInfo(ctx context.Context, args *etelecom.UpdateExternalExtensionInfoArgs) error {
	update := &etelecom.Extension{
		HotlineID:         args.HotlineID,
		ExtensionNumber:   args.ExtensionNumber,
		ExtensionPassword: args.ExtensionPassword,
		ExternalData: &etelecom.ExtensionExternalData{
			ID: args.ExternalID,
		},
	}
	return a.extensionStore(ctx).ID(args.ID).UpdateExtension(update)
}
