package query

import (
	"context"

	"o.o/api/etelecom"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (q *QueryService) GetCallLogByExternalID(
	ctx context.Context, args *etelecom.GetCallLogByExternalIDArgs,
) (*etelecom.CallLog, error) {
	if args.ExternalID == "" {
		return nil, cm.Error(cm.InvalidArgument, "external_id must not be empty", nil)
	}
	return q.callLogStore(ctx).ExternalID(args.ExternalID).GetCallLog()
}

func (q *QueryService) ListCallLogs(ctx context.Context, args *etelecom.ListCallLogsArgs) (*etelecom.ListCallLogsResponse, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	query := q.callLogStore(ctx).WithPaging(args.Paging)
	if len(args.HotlineIDs) > 0 {
		query = query.AccountIDAndHotlineIDs(args.AccountID, args.HotlineIDs)
	} else {
		query = query.AccountID(args.AccountID)
	}
	if len(args.ExtensionIDs) > 0 {
		query = query.ExtensionIDs(args.ExtensionIDs...)
	}
	res, err := query.ListCallLogs()
	if err != nil {
		return nil, err
	}
	return &etelecom.ListCallLogsResponse{
		CallLogs: res,
		Paging:   query.GetPaging(),
	}, nil
}

func (q *QueryService) GetCallLog(ctx context.Context, id dot.ID) (*etelecom.CallLog, error) {
	return q.callLogStore(ctx).ID(id).GetCallLog()
}
