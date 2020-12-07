package sqlstore

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/etelecom/convert"
	"o.o/backend/com/etelecom/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type HotlineStore struct {
	ft    HotlineFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type HotlineStoreFactory func(ctx context.Context) *HotlineStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewHotlineStore(db *cmsql.Database) HotlineStoreFactory {
	return func(ctx context.Context) *HotlineStore {
		return &HotlineStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *HotlineStore) ID(id dot.ID) *HotlineStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *HotlineStore) OptionalStatus(status status3.Status) *HotlineStore {
	s.preds = append(s.preds, s.ft.ByStatus(status).Optional())
	return s
}

func (s *HotlineStore) OptionalOwnerID(userid dot.ID) *HotlineStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(userid).Optional())
	return s
}

func (s *HotlineStore) OptionalConnectionID(connID dot.ID) *HotlineStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID).Optional())
	return s
}

func (s *HotlineStore) ConnectionID(connID dot.ID) *HotlineStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID))
	return s
}

func (s *HotlineStore) ConnectionIDs(connIDs ...dot.ID) *HotlineStore {
	s.preds = append(s.preds, sq.In("connection_id", connIDs))
	return s
}

func (s *HotlineStore) GetHotlineDB() (*model.Hotline, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var hotline model.Hotline
	err := query.ShouldGet(&hotline)
	return &hotline, err
}

func (s *HotlineStore) GetHotline() (*etelecom.Hotline, error) {
	hotline, err := s.GetHotlineDB()
	if err != nil {
		return nil, err
	}
	var res etelecom.Hotline
	if err := scheme.Convert(hotline, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *HotlineStore) ListHotlinesDB() (res []*model.Hotline, err error) {
	query := s.query().Where(s.preds)
	query = query.OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	err = query.Find((*model.Hotlines)(&res))
	return
}

func (s *HotlineStore) ListHotlines() (res []*etelecom.Hotline, _ error) {
	hotlinesDB, err := s.ListHotlinesDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(hotlinesDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *HotlineStore) CreateHotline(hotline *etelecom.Hotline) (*etelecom.Hotline, error) {
	var hotlineDB model.Hotline
	if err := scheme.Convert(hotline, &hotlineDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&hotlineDB); err != nil {
		return nil, err
	}
	return s.ID(hotline.ID).GetHotline()
}

func (s *HotlineStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("hotline").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
