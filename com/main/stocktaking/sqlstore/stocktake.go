package sqlstore

import (
	"context"

	st "etop.vn/api/main/stocktaking"
	"etop.vn/backend/com/main/stocktaking/convert"
	"etop.vn/backend/com/main/stocktaking/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

var Sort = map[string]string{
	"id":         "",
	"created_at": "",
}

type ShopStocktakeFactory func(context.Context) *ShopStocktakeStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewStocktakeStore(db *cmsql.Database) ShopStocktakeFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopStocktakeStore {
		return &ShopStocktakeStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopStocktakeStore struct {
	query  cmsql.QueryFactory
	ft     ShopStocktakeFilters
	paging *cm.Paging
	preds  []interface{}
}

func (s *ShopStocktakeStore) Paging(page *cm.Paging) *ShopStocktakeStore {
	s.paging = page
	return s
}

func (s *ShopStocktakeStore) ID(id int64) *ShopStocktakeStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShopStocktakeStore) IDs(ids ...int64) *ShopStocktakeStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *ShopStocktakeStore) ShopID(id int64) *ShopStocktakeStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShopStocktakeStore) Create(stocktake *st.ShopStocktake) error {

	var stocktakeDB = &model.ShopStocktake{}
	err := scheme.Convert(stocktake, stocktakeDB)
	if err != nil {
		return err
	}
	err = s.CreateStocktakeDB(stocktakeDB)
	return err

}

func (s *ShopStocktakeStore) Update(stocktake *st.ShopStocktake) error {
	var stocktakeDB = &model.ShopStocktake{}
	err := scheme.Convert(stocktake, stocktakeDB)
	if err != nil {
		return err
	}
	if err = s.query().Where(s.preds).ShouldUpdate(stocktakeDB); err != nil {
		return err
	}
	return nil
}

func (s *ShopStocktakeStore) UpdateAll(stocktake *st.ShopStocktake) error {
	var stocktakeDB = &model.ShopStocktake{}
	err := scheme.Convert(stocktake, stocktakeDB)
	if err != nil {
		return err
	}
	if err = s.query().Where(s.preds).UpdateAll().ShouldUpdate(stocktakeDB); err != nil {
		return err
	}
	return nil
}

func (s *ShopStocktakeStore) CreateStocktakeDB(stocktake *model.ShopStocktake) error {
	query := s.query().Where(s.preds)
	return query.ShouldInsert(stocktake)
}

func (s *ShopStocktakeStore) GetShopStocktake() (*st.ShopStocktake, error) {
	result, err := s.GetShopStocktakeDB()
	if err != nil {
		return nil, err
	}
	var stocktake = &st.ShopStocktake{}
	err = scheme.Convert(result, stocktake)
	return stocktake, err
}

func (s *ShopStocktakeStore) GetShopStocktakeDB() (*model.ShopStocktake, error) {
	query := s.query().Where(s.preds)
	var stocktake model.ShopStocktake
	err := query.ShouldGet(&stocktake)
	return &stocktake, err
}

func (s *ShopStocktakeStore) ListShopStocktake() ([]*st.ShopStocktake, error) {
	result, err := s.ListShopStocktakeDB()
	if err != nil {
		return nil, err
	}
	var stocktake []*st.ShopStocktake
	stocktake = convert.Convert_stocktakingmodel_ShopStocktakes_stocktaking_ShopStocktakes(result)
	return stocktake, nil
}

func (s *ShopStocktakeStore) ListShopStocktakeDB() ([]*model.ShopStocktake, error) {
	query := s.query().Where(s.preds)
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, s.paging, SortShopStocktake, s.ft.prefix)
	if err != nil {
		return nil, err
	}

	var stocktakes model.ShopStocktakes
	err = query.Find(&stocktakes)
	return stocktakes, err
}

func (s *ShopStocktakeStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShopStocktake)(nil))
}

func (s *ShopStocktakeStore) GetStocktakeMaximumCodeNorm() (*model.ShopStocktake, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var stocktakeModel model.ShopStocktake
	if err := query.ShouldGet(&stocktakeModel); err != nil {
		return nil, err
	}
	return &stocktakeModel, nil
}