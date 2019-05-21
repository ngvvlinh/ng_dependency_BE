package modelx

import (
	"encoding/json"
	"time"

	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	"etop.vn/backend/pkg/zdeprecated/supplier/modelx"
)

type SyncUpdateCategoriesCommand struct {
	SourceID   int64
	Data       []*model.ProductSourceCategoryExternal
	DeletedIDs []string
	LastSyncAt time.Time
	SyncState  json.RawMessage
}

type SyncUpdateProductsCommand struct {
	SourceID   int64
	Data       []*modelx.VariantExternalWithQuantity
	DeletedIDs []string
	LastSyncAt time.Time
	SyncState  json.RawMessage
}

type SyncGetOrCreateProductCommand struct {
	SourceID   int64
	Variant    *catalogmodel.VariantExternal
	LastSyncAt time.Time

	Result struct {
		ProductID int64
	}
}

type SyncUpdateProductsQuantityCommand struct {
	SourceID int64

	Updates []*modelx.ExternalQuantity
}
