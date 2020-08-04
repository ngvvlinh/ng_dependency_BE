package _uploader

import (
	"context"

	"github.com/google/wire"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/etop/upload"
)

var WireSet = wire.NewSet(
	NewUploader,
)

const (
	PurposeImportShopOrder   = storage.Purpose("import_shop_order")
	PurposeImportShopProduct = storage.Purpose("import_shop_product")
)

func SupportedPurposes() []storage.Purpose {
	return []storage.Purpose{
		PurposeImportShopOrder,
		PurposeImportShopProduct,
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
