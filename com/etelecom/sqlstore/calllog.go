package sqlstore

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/meta"
	"o.o/backend/com/etelecom/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type CallLogStore struct {
	ft    CallLogFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

type CallLogStoreFactory func(ctx context.Context) *CallLogStore

func NewCallLogStore(db *cmsql.Database) CallLogStoreFactory {
	return func(ctx context.Context) *CallLogStore {
		return &CallLogStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *CallLogStore) ID(id dot.ID) *CallLogStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CallLogStore) WithPaging(
	paging meta.Paging) *CallLogStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *CallLogStore) ExternalID(externalID string) *CallLogStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *CallLogStore) ExtensionIDs(extIDs ...dot.ID) *CallLogStore {
	s.preds = append(s.preds, sq.In("extension_id", extIDs))
	return s
}

func (s *CallLogStore) HotlineIDs(hotlineIDs ...dot.ID) *CallLogStore {
	s.preds = append(s.preds, sq.In("hotline_id", hotlineIDs))
	return s
}

func (s *CallLogStore) GetCallLogDB() (*model.CallLog, error) {
	query := s.query().Where(s.preds)
	var callLog model.CallLog
	err := query.ShouldGet(&callLog)
	return &callLog, err
}

func (s *CallLogStore) GetCallLog() (*etelecom.CallLog, error) {
	ext, err := s.GetCallLogDB()
	if err != nil {
		return nil, err
	}
	var res etelecom.CallLog
	if err := scheme.Convert(ext, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *CallLogStore) ListCallLogsDB() (res []*model.CallLog, err error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-started_at"}
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortCallLog)
	if err != nil {
		return nil, err
	}
	if err = query.Find((*model.CallLogs)(&res)); err != nil {
		return nil, err
	}
	s.Paging.Apply(res)
	return
}

func (s *CallLogStore) ListCallLogs() (res []*etelecom.CallLog, _ error) {
	callLogsDB, err := s.ListCallLogsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(callLogsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *CallLogStore) CreateCallLog(callLog *etelecom.CallLog) (*etelecom.CallLog, error) {
	var callLogDB model.CallLog
	if err := scheme.Convert(callLog, &callLogDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&callLogDB); err != nil {
		return nil, err
	}
	return s.ID(callLog.ID).GetCallLog()
}

func (s *CallLogStore) UpdateCallLog(callLog *etelecom.CallLog) error {
	var callLogDB model.CallLog
	if err := scheme.Convert(callLog, &callLogDB); err != nil {
		return err
	}
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(&callLogDB)
}
