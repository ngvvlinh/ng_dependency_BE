package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		CreateImportAttempt)
}

type ExportAttemptStore struct {
	ctx   context.Context
	ft    ExportAttemptFilters
	preds []interface{}

	includeDeleted
}

func ExportAttempt(ctx context.Context) *ExportAttemptStore {
	return &ExportAttemptStore{ctx: ctx}
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
	err := x.NewQuery().WithContext(s.ctx).
		Where(s.preds...).Where(s.filterDeleted(&s.ft)).
		OrderBy("created_at DESC").Limit(100).
		Find(&items)
	return items, err
}

func (s *ExportAttemptStore) Create(exportAttempt *model.ExportAttempt) error {
	return x.NewQuery().WithContext(s.ctx).ShouldInsert(exportAttempt)
}

func (s *ExportAttemptStore) UpdateByID(id string, exportAttempt *model.ExportAttempt) error {
	return x.NewQuery().WithContext(s.ctx).
		Where(s.ft.ByID(id)).
		ShouldUpdate(exportAttempt)
}

func CreateImportAttempt(ctx context.Context, cmd *model.CreateImportAttemptCommand) error {
	return x.Table("import_attempt").
		ShouldInsert(cmd.ImportAttempt)
}
