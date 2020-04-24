package model

import (
	"time"

	"o.o/backend/com/fabo/main/fbpost/model"
	"o.o/capi/dot"
)

// +sqlgen
type FbComment struct {
	ID                   dot.ID
	ExternalID           string
	ParentID             dot.ID
	ExternalParentID     string
	FbPostID             dot.ID
	ExternalMessage      string
	ExternalCommentCount int
	ExternalFrom         *FbCommentFrom
	ExternalAttachment   *model.Attachment
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
}

type FbCommentFrom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
