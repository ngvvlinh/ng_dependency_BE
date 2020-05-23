package fbmessaging

import (
	"time"

	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/capi/dot"
)

// +gen:event:topic=event/fbmessaging

type FbExternalMessage struct {
	ID                     dot.ID
	ExternalConversationID string
	ExternalPageID         string
	ExternalID             string
	ExternalMessage        string
	ExternalSticker        string
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
	VideoData *FbMessageDataAttachmentVideoData
	FileURL   string
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

type FbMessageDataAttachmentVideoData struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Length     int    `json:"length"`
	VideoType  int    `json:"video_type"`
	URL        string `json:"url"`
	PreviewURL string `json:"preview_url"`
	Rotation   int    `json:"rotation"`
}

type FbExternalComment struct {
	ID                   dot.ID
	ExternalPostID       string
	ExternalPageID       string
	ExternalID           string
	ExternalUserID       string
	ExternalParentID     string
	ExternalParentUserID string
	ExternalMessage      string
	ExternalCommentCount int
	ExternalParent       *FbObjectParent
	ExternalFrom         *FbObjectFrom
	ExternalAttachment   *CommentAttachment
	ExternalCreatedTime  time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type FbObjectParent struct {
	CreatedTime time.Time
	From        *FbObjectFrom
	Message     string
	ID          string
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
	ExternalPageID       string
	ExternalID           string
	PSID                 string
	ExternalUserID       string
	ExternalUserName     string
	ExternalLink         string
	ExternalUpdatedTime  time.Time
	ExternalMessageCount int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type FbExternalPosts []*FbExternalPost

type FbExternalPost struct {
	ID                  dot.ID
	ExternalPageID      string
	ExternalID          string
	ExternalParentID    string
	ExternalFrom        *FbObjectFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*PostAttachment
	ExternalCreatedTime time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type FbCustomerConversation struct {
	ID                         dot.ID
	ExternalPageID             string
	ExternalID                 string
	ExternalUserID             string
	ExternalUserName           string
	ExternalFrom               *FbObjectFrom
	IsRead                     bool
	Type                       fb_customer_conversation_type.FbCustomerConversationType
	ExternalPostAttachments    []*PostAttachment
	ExternalCommentAttachment  *CommentAttachment
	ExternalMessageAttachments []*FbMessageAttachment
	LastMessage                string
	LastMessageAt              time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	DeletedAt                  time.Time
}

type PostAttachment struct {
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

type ImageMediaDataSubAttachment struct {
	Image *MediaDataSubAttachment `json:"image"`
}

type CommentAttachment struct {
	Media  *ImageMediaDataSubAttachment `json:"media"`
	Target *TargetDataSubAttachment     `json:"target"`
	Title  string                       `json:"title"`
	Type   string                       `json:"type"`
	URL    string                       `json:"url"`
}

type FbExternalMessagesCreatedEvent struct {
	FbExternalMessages []*FbExternalMessage
}

type FbExternalCommentsCreatedEvent struct {
	FbExternalComments []*FbExternalComment
}

type FbExternalConversationsCreatedEvent struct {
	FbExternalConversations []*FbExternalConversation
}
