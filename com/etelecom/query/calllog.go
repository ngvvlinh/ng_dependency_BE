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
	if args.DateTo.Before(args.DateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	query := q.callLogStore(ctx).WithPaging(args.Paging)
	if args.OwnerID != 0 {
		query = query.AccountIDOrOwnerID(args.AccountID, args.OwnerID)
	} else {
		query = query.AccountID(args.AccountID)
	}
	if len(args.HotlineIDs) > 0 {
		query = query.HotlineIDs(args.HotlineIDs...)
	}
	if len(args.ExtensionIDs) > 0 {
		query = query.ExtensionIDs(args.ExtensionIDs...)
	}
	if args.UserID != 0 {
		query = query.UserID(args.UserID)
	}
	if args.CallerOrCallee != "" {
		query = query.CallerOrCallee(args.CallerOrCallee)
	}
	if args.CallState != 0 {
		query = query.CallState(args.CallState)
	}
	if !args.DateFrom.IsZero() {
		query = query.BetweenDateFromAndDateTo(args.DateFrom, args.DateTo)
	}
	if args.Direction != 0 {
		query = query.Direction(args.Direction)
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
