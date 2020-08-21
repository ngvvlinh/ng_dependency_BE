package _uploader

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/etop/upload"
)

const (
	PurposeImportShopOrder       = storage.Purpose("import_shop_order")
	PurposeImportShopProduct     = storage.Purpose("import_shop_product")
	PurposeImportShopFulfillment = storage.Purpose("import_shop_fulfillment")
)

func SupportedPurposes() []storage.Purpose {
	return []storage.Purpose{
		PurposeImportShopOrder,
		PurposeImportShopProduct,
		PurposeImportShopFulfillment,
	}
}

func DefaultConfig() storage.DirConfigs {
	cfg := storage.DirConfigs{}
	cfg[PurposeImportShopOrder] = storage.DirConfig{
		Path: "import_shop_order",
	}
	cfg[PurposeImportShopProduct] = storage.DirConfig{
		Path: "import_shop_product",
	}
	cfg[PurposeImportShopFulfillment] = storage.DirConfig{
		Path: "import_shop_fulfillment",
	}
	return cfg
}

func NewUploader(ctx context.Context, cfgDirs storage.DirConfigs, driver storage.Bucket) (*upload.Uploader, error) {
	mapDirs := map[string]string{}
	for _, purpose := range SupportedPurposes() {
		cfg := cfgDirs[purpose]
		if err := cfg.Validate(); err != nil {
			return nil, cm.Errorf(cm.Internal, err, "invalid config for purpose `%v`", purpose)
		}
		mapDirs[string(purpose)] = cfg.Path
	}
	return upload.NewUploader(driver, mapDirs)
}
