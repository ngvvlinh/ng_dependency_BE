package fbpaging

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

type FbPage struct {
	ID                   dot.ID
	ExternalID           string
	FbUserID             dot.ID
	ShopID               dot.ID
	UserID               dot.ID
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
	ID   dot.ID
	Name string
}

type FbPageInternal struct {
	ID        dot.ID
	Token     string
	UpdatedAt time.Time
}

type FbPageCombined struct {
	FbPage         *FbPage
	FbPageInternal *FbPageInternal
}
