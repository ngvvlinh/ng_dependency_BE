package sqlstore

import (
	"context"

	"o.o/api/shopping/setting"
	"o.o/backend/com/shopping/setting/convert"
	"o.o/backend/com/shopping/setting/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

type ShopSettingStoreFactory func(context.Context) *ShopSettingStore

func NewShopSettingStore(db *cmsql.Database) ShopSettingStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopSettingStore {
		return &ShopSettingStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopSettingStore struct {
	ft ShopSettingFilters

	query cmsql.QueryFactory
	preds []interface{}
}

func (s *ShopSettingStore) ShopID(shopID dot.ID) *ShopSettingStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *ShopSettingStore) CreateShopSettingDB(shopSetting *model.ShopSetting) error {
	sqlstore.MustNoPreds(s.preds)
	if _, err := s.query().Insert(shopSetting); err != nil {
		return err
	}

	var shopSettingResult model.ShopSetting
	if err := s.query().Where(s.ft.ByShopID(shopSetting.ShopID)).ShouldGet(&shopSettingResult); err != nil {
		return err
	}

	shopSetting.CreatedAt = shopSettingResult.CreatedAt
	shopSetting.UpdatedAt = shopSettingResult.UpdatedAt
	return nil
}

func (s *ShopSettingStore) CreateShopSetting(shopSetting *setting.ShopSetting) error {
	sqlstore.MustNoPreds(s.preds)

	shopSettingDB := new(model.ShopSetting)
	convert.Convert_setting_ShopSetting_settingmodel_ShopSetting(shopSetting, shopSettingDB)

	if err := s.CreateShopSettingDB(shopSettingDB); err != nil {
		return err
	}

	shopSettingResult, err := s.ShopID(shopSettingDB.ShopID).GetShopSetting()
	if err != nil {
		return err
	}

	*shopSetting = *shopSettingResult
	return nil
}

func (s *ShopSettingStore) UpdateShopSettingDB(shopSetting *model.ShopSetting) error {
	sqlstore.MustNoPreds(s.preds)
	if err := s.query().Where(s.ft.ByShopID(shopSetting.ShopID)).ShouldUpdate(shopSetting); err != nil {
		return err
	}

	var shopSettingResult model.ShopSetting
	if err := s.query().Where(s.ft.ByShopID(shopSetting.ShopID)).ShouldGet(&shopSettingResult); err != nil {
		return err
	}

	*shopSetting = shopSettingResult
	return nil
}

func (s *ShopSettingStore) UpdateShopSetting(shopSetting *setting.ShopSetting) error {
	shopSettingDB := new(model.ShopSetting)
	if err := scheme.Convert(shopSetting, shopSettingDB); err != nil {
		return err
	}

	if err := s.UpdateShopSettingDB(shopSettingDB); err != nil {
		return err
	}

	shopSettingResult, err := s.ShopID(shopSettingDB.ShopID).GetShopSetting()
	if err != nil {
		return err
	}

	*shopSetting = *shopSettingResult
	return nil
}

func (s *ShopSettingStore) GetShopSettingDB() (*model.ShopSetting, error) {
	query := s.query().Where(s.preds)

	var shopSetting model.ShopSetting
	err := query.ShouldGet(&shopSetting)
	return &shopSetting, err
}

func (s *ShopSettingStore) GetShopSetting() (*setting.ShopSetting, error) {
	shopSettingDB, err := s.GetShopSettingDB()
	if err != nil {
		return nil, err
	}

	result := &setting.ShopSetting{}
	if err = scheme.Convert(shopSettingDB, result); err != nil {
		return nil, err
	}
	return result, err
}

func (s *ShopSettingStore) UpdateDirectShipmentShopSettingDB(shopSetting *setting.ShopSetting) error {
	var update = make(map[string]interface{})
	if shopSetting != nil {
		update["allow_connect_direct_shipment"] = shopSetting.AllowConnectDirectShipment
	}
	if err := s.query().Table("shop_setting").Where(s.ft.ByShopID(shopSetting.ShopID)).ShouldUpdateMap(update); err != nil {
		return err
	}
	return nil
}
