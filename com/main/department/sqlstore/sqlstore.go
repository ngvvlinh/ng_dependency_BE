package sqlstore

import (
	"context"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/filter"
	"time"

	"o.o/api/main/department"
	"o.o/api/meta"
	"o.o/backend/com/main/department/convert"
	"o.o/backend/com/main/department/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type DepartmentStoreFactory func(ctx context.Context) *DepartmentStore

type DepartmentStore struct {
	ft DepartmentFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted

	ctx context.Context
}

func NewDepartmentStore(db *cmsql.Database) DepartmentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *DepartmentStore {
		return &DepartmentStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *DepartmentFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (s *DepartmentStore) Filters(filters meta.Filters) *DepartmentStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *DepartmentStore) WithPaging(paging meta.Paging) *DepartmentStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *DepartmentStore) ByID(id dot.ID) *DepartmentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *DepartmentStore) ByIDs(ids []dot.ID) *DepartmentStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *DepartmentStore) ByAccountID(accountID dot.ID) *DepartmentStore {
	s.preds = append(s.preds, s.ft.ByAccountID(accountID))
	return s
}

func (s *DepartmentStore) ByOptionalAccountID(accountID dot.ID) *DepartmentStore {
	s.preds = append(s.preds, s.ft.ByAccountID(accountID).Optional())
	return s
}

func (s *DepartmentStore) ByName(name string) *DepartmentStore {
	s.preds = append(s.preds, s.ft.ByName(name))
	return s
}

func (s *DepartmentStore) ByNameNorm(name filter.FullTextSearch) *DepartmentStore {
	s.preds = append(s.preds, s.ft.Filter(`name @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
}

func (s *DepartmentStore) createDepartmentDB(department *model.Department) error {
	sqlstore.MustNoPreds(s.preds)
	if department.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return s.query().ShouldInsert(department)
}

func (s *DepartmentStore) CreateDepartment(department *department.Department) error {
	var departmentDB model.Department
	if err := scheme.Convert(department, &departmentDB); err != nil {
		return err
	}
	return s.createDepartmentDB(&departmentDB)
}

func (s *DepartmentStore) getDepartmentDB() (*model.Department, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var departmentDB model.Department
	if err := query.ShouldGet(&departmentDB); err != nil {
		return nil, err
	}
	return &departmentDB, nil
}

func (s *DepartmentStore) GetDepartment() (*department.Department, error) {
	inv, err := s.getDepartmentDB()
	if err != nil {
		return nil, err
	}
	var res department.Department
	if err = scheme.Convert(inv, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *DepartmentStore) listDepartmentsDB() ([]*model.Department, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	var err error
	query, err = sqlstore.LimitSort(query, &s.Paging, SortDepartment)
	if err != nil {
		return nil, err
	}
	var departments model.Departments
	err = query.Find(&departments)
	s.Paging.Apply(departments)
	return departments, err
}

func (s *DepartmentStore) ListDepartments() ([]*department.Department, error) {
	departmentsDB, err := s.listDepartmentsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_departmentmodel_Departments_department_Departments(departmentsDB), err
}

func (s *DepartmentStore) updateDepartmentDB(args *model.Department) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *DepartmentStore) UpdateDepartment(args *department.Department) error {
	departmentModel := &model.Department{}
	if err := scheme.Convert(args, departmentModel); err != nil {
		return err
	}
	return s.updateDepartmentDB(departmentModel)
}

func (s *DepartmentStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("department").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
