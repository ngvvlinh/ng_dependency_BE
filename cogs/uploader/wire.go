package _uploader

import (
	"github.com/google/wire"

	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
)

var WireSet = wire.NewSet(
	NewUploader,
)

func DefaultConfig() upload.Config {
	return upload.Config{
		DirImportShopOrder:   "/tmp",
		DirImportShopProduct: "/tmp",
	}
}

func NewUploader(cfg upload.Config) (*upload.Uploader, error) {
	return upload.NewUploader(map[string]string{
		model.ImportTypeShopOrder.String():   cfg.DirImportShopOrder,
		model.ImportTypeShopProduct.String(): cfg.DirImportShopProduct,
	})
}
