package query

import (
	"context"
	"strconv"

	"o.o/api/etelecom"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/code/gencode"
)

func (q *QueryService) GetExtension(ctx context.Context, args *etelecom.GetExtensionArgs) (*etelecom.Extension, error) {
	query := q.extensionStore(ctx)
	if args.ID != 0 {
		return query.ID(args.ID).GetExtension()
	}
	query = query.UserID(args.UserID).AccountID(args.AccountID).OptionalHotlineID(args.HotlineID).OptionalSubscriptionID(args.SubscriptionID)
	return query.GetExtension()
}

func (q *QueryService) ListExtensions(ctx context.Context, args *etelecom.ListExtensionsArgs) ([]*etelecom.Extension, error) {
	query := q.extensionStore(ctx)
	if len(args.HotlineIDs) > 0 {
		query = query.HotlineIDs(args.HotlineIDs...)
	}
	if len(args.AccountIDs) > 0 {
		query = query.AccountIDs(args.AccountIDs...)
	}
	if len(args.ExtensionNumbers) > 0 {
		query = query.ExtensionNumbers(args.ExtensionNumbers...)
	}
	return query.ListExtensions()
}

func (q *QueryService) GetPrivateExtensionNumber(ctx context.Context, _ *pbcm.Empty) (string, error) {
	var code int
	if err := q.db.SQL(`SELECT nextval('extension_number')`).Scan(&code); err != nil {
		return "", err
	}
	checksumDigit := gencode.CheckSumDigitUPC(strconv.Itoa(code))
	return checksumDigit, nil
}
