package fbmessagetemplate

import (
	"time"

	"o.o/capi/dot"
)

type FbMessageTemplate struct {
	ID        dot.ID
	ShopID    dot.ID
	Template  string
	ShortCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// +convert:create=FbMessageTemplate
type CreateFbMessageTemplate struct {
	ShopID    dot.ID
	Template  string
	ShortCode string
}

// type:update=FbMessageTemplate
type UpdateFbMessageTemplate struct {
	ID        dot.ID
	ShopID    dot.ID
	Template  dot.NullString
	ShortCode dot.NullString
}

type DeleteFbMessageTemplate struct {
	ShopID dot.ID
	ID     dot.ID
}
