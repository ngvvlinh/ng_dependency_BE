package model

import (
	"encoding/json"
	"time"

	cm "etop.vn/backend/pkg/common"
)

type GetHistoryQuery struct {
	*cm.Paging
	Table   string
	Filters map[string]interface{}
	KeepRaw bool

	Result struct {
		Len  int
		Data json.RawMessage
		Raws []json.RawMessage
	}
}

type GetBalanceShopCommand struct {
	ShopID int64

	Result struct {
		Amount int
	}
}

type GetShippingSources struct {
	Type   ShippingProvider
	Names  []string
	Result []*ShippingSource
}

type GetShippingSource struct {
	ID       int64
	Name     string
	Username string
	Type     ShippingProvider

	Result struct {
		ShippingSource         *ShippingSource
		ShippingSourceInternal *ShippingSourceInternal
	}
}

type CreateShippingSource struct {
	Name     string
	Type     ShippingProvider
	Username string
	Password string

	Result struct {
		ShippingSource         *ShippingSource
		ShippingSourceInternal *ShippingSourceInternal
	}
}

type GetShippingSourceInternal struct {
	ID     int64
	Result *ShippingSourceInternal
}

type UpdateOrCreateShippingSourceInternal struct {
	ID          int64
	LastSyncAt  time.Time
	AccessToken string
	ExpiresAt   time.Time
	Secret      *ShippingSourceSecret

	Result struct {
		Updated int
	}
}

type VTPostRequestInternalInfo struct {
	ShippingSourceID int64
	Username         string
	Password         string
}
