package fbmessaging

import (
	"time"

	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/capi/dot"
)

// +gen:event:topic=event/fbmessaging

type FbExternalMessage struct {
	ID                     dot.ID
	FbConversationID       dot.ID
	ExternalConversationID string
	FbPageID               dot.ID
	ExternalID             string
	ExternalMessage        string
	ExternalTo             []*FbObjectTo
	ExternalFrom           *FbObjectFrom
	ExternalAttachments    []*FbMessageAttachment
	ExternalCreatedTime    time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type FbObjectTo struct {
	ID    string
	Name  string
	Email string
}

type FbObjectFrom struct {
	ID    string
	Name  string
	Email string
}

type FbMessageAttachment struct {
	ID        string
	ImageData *FbMessageAttachmentImageData
	MimeType  string
	Name      string
	Size      int
}

type FbMessageAttachmentImageData struct {
	Width           int
	Height          int
	MaxWidth        int
	MaxHeight       int
	URL             string
	PreviewURL      string
	ImageType       int
	RenderAsSticker bool
}

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
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Attachment struct {
	MediaType      string
	Type           string
	SubAttachments []*SubAttachment
}

type SubAttachment struct {
	Media  *MediaDataSubAttachment
	Target *TargetDataSubAttachment
	Type   string
	URL    string
}

type MediaDataSubAttachment struct {
	Height int
	Width  int
	Src    string
}

type TargetDataSubAttachment struct {
	ID  string
	URL string
}

type FbExternalConversation struct {
	ID                   dot.ID
	FbPageID             dot.ID
	ExternalID           string
	ExternalUserID       string
	ExternalUserName     string
	ExternalLink         string
	ExternalUpdatedTime  time.Time
	ExternalMessageCount int
	LastMessage          string
	LastMessageAt        time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

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
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type FbCustomerConversation struct {
	ID               dot.ID
	FbPageID         dot.ID
	ExternalID       string
	ExternalUserID   string
	ExternalUserName string
	IsRead           bool
	Type             fb_customer_conversation_type.FbCustomerConversationType
	PostAttachments  []*PostAttachment
	LastMessage      string
	LastMessageAt    time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type PostAttachment struct {
	Media *PostAttachmentMedia
	Type  string
}

type PostAttachmentMedia struct {
	Height int
	Width  int
	Src    string
}

type FbExternalMessagesCreatedEvent struct {
	FbExternalMessages []*FbExternalMessage
}
