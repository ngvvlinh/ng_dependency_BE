package model

import (
	"encoding/json"
	"time"
)

//-- Kiotviet

type SyncUpdateCategoriesCommand struct {
	SourceID   int64
	Data       []*ProductSourceCategoryExternal
	DeletedIDs []string
	LastSyncAt time.Time
	SyncState  json.RawMessage
}

type SyncUpdateProductsCommand struct {
	SourceID   int64
	Data       []*VariantExternalWithQuantity
	DeletedIDs []string
	LastSyncAt time.Time
	SyncState  json.RawMessage
}

type SyncGetOrCreateProductCommand struct {
	SourceID   int64
	Variant    *VariantExternal
	LastSyncAt time.Time

	Result struct {
		ProductID int64
	}
}

type SyncUpdateProductsQuantityCommand struct {
	SourceID int64

	Updates []*ExternalQuantity
}
