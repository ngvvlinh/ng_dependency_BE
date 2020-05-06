package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbExternalMessage struct {
	ID                     dot.ID `paging:"id"`
	FbConversationID       dot.ID
	ExternalConversationID string
	FbPageID               dot.ID
	ExternalID             string
	ExternalMessage        string
	ExternalTo             []*FbObjectTo
	ExternalFrom           *FbObjectFrom
	ExternalAttachments    []*FbMessageAttachment
	ExternalCreatedTime    time.Time `paging:"external_created_time"`
	CreatedAt              time.Time `sq:"create"`
	UpdatedAt              time.Time `sq:"update"`
	DeletedAt              time.Time
}

type FbObjectTo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FbObjectFrom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FbMessageAttachment struct {
	ID        string                        `json:"id"`
	ImageData *FbMessageAttachmentImageData `json:"image_data"`
	MimeType  string                        `json:"mime_type"`
	Name      string                        `json:"name"`
	Size      int                           `json:"size"`
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

// +sqlgen
type FbExternalConversation struct {
	ID                   dot.ID
	FbPageID             dot.ID
	ExternalID           string
	ExternalUserID       string
	ExternalUserName     string
	ExternalLink         string
	ExternalUpdatedTime  time.Time
	ExternalMessageCount int
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
}

// +sqlgen
type FbExternalComment struct {
	ID                   dot.ID
	FbPostID             dot.ID
	FbPageID             dot.ID
	ExternalID           string
	ExternalUserID       string
	ExternalParentID     string
	ExternalParentUserID string
	ExternalMessage      string
	ExternalCommentCount int
	ExternalFrom         *FbObjectFrom
	ExternalAttachment   *Attachment
	ExternalCreatedTime  time.Time
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
}

// +sqlgen
type FbExternalPost struct {
	ID                  dot.ID
	FbPageID            dot.ID
	ExternalID          string
	ExternalParentID    string
	ExternalFrom        *FbObjectFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*Attachment
	ExternalCreatedTime time.Time
	ExternalUpdatedTime time.Time
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
}

type Attachment struct {
	MediaType      string           `json:"media_type"`
	Type           string           `json:"type"`
	SubAttachments []*SubAttachment `json:"sub_attachments"`
}

type SubAttachment struct {
	Media  *MediaDataSubAttachment  `json:"media"`
	Target *TargetDataSubAttachment `json:"target"`
	Type   string                   `json:"type"`
	URL    string                   `json:"url"`
}

type MediaDataSubAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type TargetDataSubAttachment struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// +sqlgen
type FbCustomerConversation struct {
	ID               dot.ID `paging:"id"`
	FbPageID         dot.ID
	ExternalID       string
	ExternalUserID   string
	ExternalUserName string
	IsRead           bool
	PostAttachments  []*PostAttachment
	Type             int
	LastMessage      string
	LastMessageAt    time.Time `paging:"last_message_at"`
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
	DeletedAt        time.Time
}

type PostAttachment struct {
	Media *PostAttachmentMedia `json:"media"`
	Type  string               `json:"type"`
}

type PostAttachmentMedia struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}
