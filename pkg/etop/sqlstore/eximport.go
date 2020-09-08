package sqlstore

import (
	"context"
	"time"

	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type ExportAttemptStoreInterface interface {
	CreateImportAttempt(ctx context.Context, cmd *model.CreateImportAttemptCommand) error
}

type ExportAttemptStore struct {
	query func() cmsql.QueryInterface
	ft    ExportAttemptFilters
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type ExportAttemptStoreFactory func(ctx context.Context) *ExportAttemptStore

func NewExportAttemptStore(db com.MainDB) ExportAttemptStoreFactory {
	return func(ctx context.Context) *ExportAttemptStore {
		return &ExportAttemptStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func BuildExportAttempStore(db com.MainDB) *ExportAttemptStore {
	return NewExportAttemptStore(db)(context.Background())
}

func (s *ExportAttemptStore) IncludeDeleted() *ExportAttemptStore {
	s.includeDeleted = true
	return s
}

func (s *ExportAttemptStore) AccountID(id dot.ID) *ExportAttemptStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *ExportAttemptStore) NotYetExpired() *ExportAttemptStore {
	s.preds = append(s.preds, sq.NewExpr("expires_at > ?", time.Now()))
	return s
}

func (s *ExportAttemptStore) List() ([]*model.ExportAttempt, error) {
	var items model.ExportAttempts
	err := s.query().
		Where(s.preds...).Where(s.includeDeleted.FilterDeleted(&s.ft)).
		OrderBy("created_at DESC").Limit(100).
		Find(&items)
	return items, err
}

func (s *ExportAttemptStore) Create(exportAttempt *model.ExportAttempt) error {
	return s.query().ShouldInsert(exportAttempt)
}

func (s *ExportAttemptStore) UpdateByID(id string, exportAttempt *model.ExportAttempt) error {
	return s.query().
		Where(s.ft.ByID(id)).
		ShouldUpdate(exportAttempt)
}

func (s *ExportAttemptStore) CreateImportAttempt(ctx context.Context, cmd *model.CreateImportAttemptCommand) error {
	return s.query().Table("import_attempt").
		ShouldInsert(cmd.ImportAttempt)
}
