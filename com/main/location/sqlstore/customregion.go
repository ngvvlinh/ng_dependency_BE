package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/backend/com/main/location/convert"
	"etop.vn/backend/com/main/location/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type CustomRegionStore struct {
	ft    CustomRegionFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

var scheme = conversion.Build(convert.RegisterConversions)

type CustomRegionFactory func(ctx context.Context) *CustomRegionStore

func NewCustomRegionStore(db *cmsql.Database) CustomRegionFactory {
	return func(ctx context.Context) *CustomRegionStore {
		return &CustomRegionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *CustomRegionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *CustomRegionFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *CustomRegionStore) ID(id dot.ID) *CustomRegionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CustomRegionStore) ProvinceCode(pCode string) *CustomRegionStore {
	s.preds = append(s.preds, sq.NewExpr("province_codes @> ?", core.Array{V: []string{pCode}}))
	return s
}

func (s *CustomRegionStore) GetCustomRegionDB() (*model.CustomRegion, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = query.OrderBy("created_at DESC")
	var region model.CustomRegion
	err := query.ShouldGet(&region)
	return &region, err
}

func (s *CustomRegionStore) GetCustomRegion() (*location.CustomRegion, error) {
	regionDB, err := s.GetCustomRegionDB()
	if err != nil {
		return nil, err
	}
	var res location.CustomRegion
	if err := scheme.Convert(regionDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *CustomRegionStore) ListCustomRegionDBs() (res []*model.CustomRegion, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	err = query.Find((*model.CustomRegions)(&res))
	return
}

func (s *CustomRegionStore) ListCustomRegions() (res []*location.CustomRegion, _ error) {
	regions, err := s.ListCustomRegionDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(regions, &res); err != nil {
		return nil, err
	}
	return
}

func (s *CustomRegionStore) CreateCustomRegion(region *location.CustomRegion) (*location.CustomRegion, error) {
	sqlstore.MustNoPreds(s.preds)
	if region.ID == 0 {
		region.ID = cm.NewID()
	}
	region.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	var regionDB model.CustomRegion
	if err := scheme.Convert(region, &regionDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&regionDB); err != nil {
		return nil, err
	}
	return s.ID(region.ID).GetCustomRegion()
}

func (s *CustomRegionStore) UpdateCustomRegion(region *location.CustomRegion) error {
	sqlstore.MustNoPreds(s.preds)
	var regionDB model.CustomRegion
	if err := scheme.Convert(region, &regionDB); err != nil {
		return err
	}
	query := s.query().Where(s.ft.ByID(region.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(&regionDB)
}

func (s *CustomRegionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("custom_region").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *CustomRegionStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
