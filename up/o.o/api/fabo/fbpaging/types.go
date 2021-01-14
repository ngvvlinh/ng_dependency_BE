package fbpaging

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/fbpaging

type FbExternalPage struct {
	ID                   dot.ID
	ExternalID           string
	ExternalUserID       string
	ShopID               dot.ID
	ExternalName         string
	ExternalCategory     string
	ExternalCategoryList []*ExternalCategory
	ExternalTasks        []string
	ExternalPermissions  []string
	ExternalImageURL     string
	Status               status3.Status
	ConnectionStatus     status3.Status
	CreatedAt            time.Time
	UpdatedAt            time.Time
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
