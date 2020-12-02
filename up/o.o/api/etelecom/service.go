package etelecom

import (
	"context"

	cm "o.o/api/top/types/common"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateExtension(context.Context, *CreateExtensionArgs) (*Extension, error)
	DeleteExtension(ctx context.Context, id dot.ID) error
	UpdateExternalExtensionInfo(context.Context, *UpdateExternalExtensionInfoArgs) error
}

type QueryService interface {
	GetHotline(context.Context, *GetHotlineArgs) (*Hotline, error)
	ListHotlines(context.Context, *ListHotlinesArgs) ([]*Hotline, error)

	GetExtension(context.Context, *GetExtensionArgs) (*Extension, error)
	ListExtensions(context.Context, *ListExtensionsArgs) ([]*Extension, error)

	// generate extension number => then use it for create external extension
	GetPrivateExtensionNumber(context.Context, *cm.Empty) (extensionNumber string, _ error)
}

// +convert:create=Extension
type CreateExtensionArgs struct {
	ExtensionNumber   string
	UserID            dot.ID
	AccountID         dot.ID
	ExtensionPassword string
	ConnectionID      dot.ID
}

func (args *CreateExtensionArgs) Validate() error {
	if args.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing user ID")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	return nil
}

type GetHotlineArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

type ListHotlinesArgs struct {
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type GetExtensionArgs struct {
	ID           dot.ID
	UserID       dot.ID
	AccountID    dot.ID
	ConnectionID dot.ID
}

type ListExtensionsArgs struct {
	UserID       dot.ID
	ConnectionID dot.ID
}

type UpdateExternalExtensionInfoArgs struct {
	ID                dot.ID
	HotlineID         dot.ID
	ExternalID        string
	ExtensionNumber   string
	ExtensionPassword string
}
