package types

import (
	"etop.vn/api/top/types/common"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type ListInventoryLevelsFilter struct {
	VariantID filter.IDs `json:"variant_id"`
}

type ListInventoryLevelsRequest struct {
	Filter ListInventoryLevelsFilter `json:"filter"`
	Paging *common.CursorPaging      `json:"paging"`
}

func (m *ListInventoryLevelsRequest) String() string { return jsonx.MustMarshalToString(m) }

type InventoryLevel struct {
	VariantId         dot.ID   `json:"variant_id"`
	AvailableQuantity int      `json:"available_quantity"`
	ReservedQuantity  int      `json:"reserved_quantity"`
	PickedQuantity    int      `json:"picked_quantity"`
	UpdatedAt         dot.Time `json:"updated_at"`
}

func (m *InventoryLevel) String() string { return jsonx.MustMarshalToString(m) }

type InventoryLevelsResponse struct {
	InventoryLevels []*InventoryLevel      `json:"inventory_levels"`
	Paging          *common.CursorPageInfo `json:"paging"`
}

func (m *InventoryLevelsResponse) String() string { return jsonx.MustMarshalToString(m) }
