package telecom

import (
	"context"
	"database/sql"

	"o.o/api/etelecom"
	etelecommodel "o.o/backend/com/etelecom/model"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq/core"
)

type TelecomStoreInterface interface {
	GetCallLogsExport(ctx context.Context, args *etelecom.ListCallLogsExportArgs) (*ListCallLogsExportResponse, error)
}

type ListCallLogsExportResponse struct {
	Total int
	// only for ResultAsRows
	Rows *sql.Rows
	Opts core.Opts
}

type TelecomStore struct {
	TelecomDB com.EtelecomDB
}

func BindTelecomStore(s *TelecomStore) TelecomStoreInterface {
	return s
}

func (st *TelecomStore) GetCallLogsExport(ctx context.Context, args *etelecom.ListCallLogsExportArgs) (*ListCallLogsExportResponse, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	s := (*cmsql.Database)(st.TelecomDB).Table("call_log")

	if args.OwnerID != 0 {
		s = s.Where("account_id = ? OR (owner_id = ? AND account_id IS NULL)", args.AccountID, args.OwnerID)
	} else {
		s = s.Where("account_id = ?", args.AccountID)
	}
	if len(args.ExtensionIDs) > 0 {
		s = s.In("extension_id", args.ExtensionIDs)
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !args.DateFrom.IsZero() {
		s = s.Where("started_at >= ? AND started_at < ?", args.DateFrom, args.DateTo)
	}

	res := &ListCallLogsExportResponse{}
	{
		s2 := s.Clone()
		total, err := s2.Count(&etelecommodel.CallLogs{})
		if err != nil {
			return nil, err
		}
		res.Total = total
	}
	{
		s = s.OrderBy("started_at DESC")
		opts, rows, err := s.FindRows((*etelecommodel.CallLogs)(nil))
		if err != nil {
			return nil, err
		}
		res.Opts = opts
		res.Rows = rows
	}
	return res, nil
}
