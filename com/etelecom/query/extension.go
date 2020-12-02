package query

import (
	"context"

	"o.o/api/etelecom"
)

func (q *QueryService) GetExtension(ctx context.Context, args *etelecom.GetExtensionArgs) (*etelecom.Extension, error) {
	query := q.extensionStore(ctx)
	if args.ID != 0 {
		return query.ID(args.ID).GetExtension()
	}
	query = query.UserID(args.UserID).AccountID(args.AccountID).OptionalConnectionID(args.ConnectionID)
	return query.GetExtension()
}

func (q *QueryService) ListExtensions(ctx context.Context, args *etelecom.ListExtensionsArgs) ([]*etelecom.Extension, error) {
	return q.extensionStore(ctx).OptionalUserID(args.UserID).OptionalConnectionID(args.ConnectionID).ListExtensions()
}
