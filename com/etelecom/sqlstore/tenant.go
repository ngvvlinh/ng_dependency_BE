package sqlstore

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/meta"
	"o.o/backend/com/etelecom/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TenantStore struct {
	ft    TenantFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

type TenantStoreFactory func(ctx context.Context) *TenantStore

func NewTenantStore(db *cmsql.Database) TenantStoreFactory {
	return func(ctx context.Context) *TenantStore {
		return &TenantStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *TenantStore) WithPaging(paging meta.Paging) *TenantStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *TenantStore) ID(id dot.ID) *TenantStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TenantStore) OwnerID(id dot.ID) *TenantStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(id))
	return s
}

func (s *TenantStore) ConnectionID(id dot.ID) *TenantStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(id))
	return s
}

func (s *TenantStore) GetTenanteDB() (*model.Tenant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var tenant model.Tenant
	err := query.ShouldGet(&tenant)
	return &tenant, err
}

func (s *TenantStore) GetTenant() (*etelecom.Tenant, error) {
	tenant, err := s.GetTenanteDB()
	if err != nil {
		return nil, err
	}
	var res etelecom.Tenant
	if err = scheme.Convert(tenant, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *TenantStore) ListTenantsDB() (res []*model.Tenant, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortTenant)
	if err != nil {
		return nil, err
	}
	err = query.Find((*model.Tenants)(&res))
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(res)
	return
}

func (s *TenantStore) ListTenants() (res []*etelecom.Tenant, _ error) {
	tenantsDB, err := s.ListTenantsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(tenantsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *TenantStore) CreateTenant(tenant *etelecom.Tenant) (*etelecom.Tenant, error) {
	var tenantDB model.Tenant
	if err := scheme.Convert(tenant, &tenantDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&tenantDB); err != nil {
		return nil, err
	}
	return s.ID(tenant.ID).GetTenant()
}

func (s *TenantStore) UpdateTenant(tenant *etelecom.Tenant) error {
	var tenantDB model.Tenant
	if err := scheme.Convert(tenant, &tenantDB); err != nil {
		return err
	}
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(&tenantDB)
}

func (s *TenantStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("tenant").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
