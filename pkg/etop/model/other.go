package model

import (
	"encoding/json"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
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
	ShopID dot.ID

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
	ID       dot.ID
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
	ID     dot.ID
	Result *ShippingSourceInternal
}

type UpdateOrCreateShippingSourceInternal struct {
	ID          dot.ID
	LastSyncAt  time.Time
	AccessToken string
	ExpiresAt   time.Time
	Secret      *ShippingSourceSecret

	Result struct {
		Updated int
	}
}

type VTPostRequestInternalInfo struct {
	ShippingSourceID dot.ID
	Username         string
	Password         string
}
