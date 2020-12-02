package sqlstore

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/backend/com/etelecom/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ExtensionStore struct {
	ft    ExtensionFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type ExtensionStoreFactory func(ctx context.Context) *ExtensionStore

func NewExtensionStore(db *cmsql.Database) ExtensionStoreFactory {
	return func(ctx context.Context) *ExtensionStore {
		return &ExtensionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *ExtensionStore) ID(id dot.ID) *ExtensionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ExtensionStore) UserID(id dot.ID) *ExtensionStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *ExtensionStore) OptionalUserID(id dot.ID) *ExtensionStore {
	s.preds = append(s.preds, s.ft.ByUserID(id).Optional())
	return s
}

func (s *ExtensionStore) AccountID(id dot.ID) *ExtensionStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *ExtensionStore) OptionalConnectionID(connID dot.ID) *ExtensionStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID).Optional())
	return s
}

func (s *ExtensionStore) GetExtensionDB() (*model.Extension, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var extension model.Extension
	err := query.ShouldGet(&extension)
	return &extension, err
}

func (s *ExtensionStore) GetExtension() (*etelecom.Extension, error) {
	ext, err := s.GetExtensionDB()
	if err != nil {
		return nil, err
	}
	var res etelecom.Extension
	if err := scheme.Convert(ext, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ExtensionStore) ListExtensionsDB() (res []*model.Extension, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	err = query.Find((*model.Extensions)(&res))
	return
}

func (s *ExtensionStore) ListExtensions() (res []*etelecom.Extension, _ error) {
	extensionsDB, err := s.ListExtensionsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(extensionsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ExtensionStore) CreateExtension(ext *etelecom.Extension) (*etelecom.Extension, error) {
	var extDB model.Extension
	if err := scheme.Convert(ext, &extDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&extDB); err != nil {
		return nil, err
	}
	return s.ID(ext.ID).GetExtension()
}

func (s *ExtensionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("extension").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
