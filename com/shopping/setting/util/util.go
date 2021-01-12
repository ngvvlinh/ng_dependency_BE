package util

import (
	"fmt"

	"o.o/api/shopping/setting"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

const (
	versionCached = "v1.0.0"
	ttl           = 1 * 60 * 60
	prefixKey     = "shop_setting"
)

type ShopSettingUtil struct {
	rd redis.Store
}

func NewShopSettingUtil(rd redis.Store) *ShopSettingUtil {
	return &ShopSettingUtil{
		rd: rd,
	}
}

func (s *ShopSettingUtil) generateKey(shopID dot.ID) string {
	return fmt.Sprintf("%s:%s:%d", prefixKey, versionCached, shopID.Int64())
}

func (s *ShopSettingUtil) SetShopSetting(shopID dot.ID, shopSetting setting.ShopSetting) error {
	key := s.generateKey(shopID)

	return s.rd.SetWithTTL(key, shopSetting, ttl)
}

func (s *ShopSettingUtil) GetShopSetting(shopID dot.ID) (*setting.ShopSetting, error) {
	key := s.generateKey(shopID)

	var shopSetting setting.ShopSetting

	err := s.rd.Get(key, &shopSetting)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		return &shopSetting, nil
	default:
		return nil, err
	}
}

func (s *ShopSettingUtil) ClearShopSetting(shopID dot.ID) error {
	key := s.generateKey(shopID)

	return s.rd.Del(key)
}
