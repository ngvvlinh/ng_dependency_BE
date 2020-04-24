package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbMessage struct {
	ID                  dot.ID
	ExternalID          string
	FbConversationID    dot.ID
	ExternalMessage     string
	ExternalTo          *FbMessageTo
	ExternalFrom        *FbMessageFrom
	ExternalAttachments []*FbMessageAttachment
	ExternalCreatedTime time.Time
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
}

type FbMessageTo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FbMessageFrom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FbMessageAttachment struct {
	ID        string                        `json:"id"`
	ImageData *FbMessageAttachmentImageData `json:"image_data"`
}

type FbMessageAttachmentImageData struct {
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	MaxWidth        int    `json:"max_width"`
	MaxHeight       int    `json:"max_height"`
	URL             string `json:"url"`
	PreviewURL      string `json:"preview_url"`
	ImageType       int    `json:"image_type"`
	RenderAsSticker bool   `json:"render_as_sticker"`
}
