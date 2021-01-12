package fbpaging

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/fbpaging

type FbExternalPage struct {
	ID                   dot.ID              `json:"id"`
	ExternalID           string              `json:"external_id"`
	ShopID               dot.ID              `json:"shop_id"`
	ExternalName         string              `json:"external_name"`
	ExternalCategory     string              `json:"external_category"`
	ExternalCategoryList []*ExternalCategory `json:"external_category_list"`
	ExternalTasks        []string            `json:"external_tasks"`
	ExternalPermissions  []string            `json:"external_permissions"`
	ExternalImageURL     string              `json:"external_image_url"`
	Status               status3.Status      `json:"status"`
	ConnectionStatus     status3.Status      `json:"connection_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

type ExternalCategory struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

type FbExternalPageInternal struct {
	ID         dot.ID    `json:"id"`
	ExternalID string    `json:"external_id"`
	Token      string    `json:"token"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FbExternalPageCombineds []*FbExternalPageCombined

type FbExternalPageCombined struct {
	FbExternalPage         *FbExternalPage
	FbExternalPageInternal *FbExternalPageInternal
}

type FbExternalPagesCreatedOrUpdatedEvent struct {
	ExternalPageIDs []string
}
