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
	CreatedAt              time.Time `compare:"ignore"`
	UpdatedAt              time.Time `compare:"ignore"`
}

type FbObjectTo struct {
	ID        string
	Name      string
	Email     string
	FirstName string
	LastName  string
	ImageURL  string
}

type FbObjectFrom struct {
	ID        string
	Name      string
	Email     string
	FirstName string
	LastName  string
	ImageURL  string `compare:"ignore"`
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
	Width      int
	Height     int
	Length     int
	VideoType  int
	URL        string
	PreviewURL string
	Rotation   int
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
	ExternalCommentCount int `compare:"ignore"`
	ExternalParent       *FbObjectParent
	ExternalFrom         *FbObjectFrom
	ExternalAttachment   *CommentAttachment `compare:"ignore"`
	ExternalCreatedTime  time.Time
	CreatedAt            time.Time `compare:"ignore"`
	UpdatedAt            time.Time `compare:"ignore"`
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
	CreatedAt            time.Time `compare:"ignore"`
	UpdatedAt            time.Time `compare:"ignore"`
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
	ExternalAttachments []*PostAttachment `compare:"ignore"`
	ExternalCreatedTime time.Time
	CreatedAt           time.Time `compare:"ignore"`
	UpdatedAt           time.Time `compare:"ignore"`
	ExternalParent      *FbExternalPost
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
	ExternalUserPictureURL     string
	ExternalPostAttachments    []*PostAttachment
	ExternalCommentAttachment  *CommentAttachment
	ExternalMessageAttachments []*FbMessageAttachment
	LastMessage                string
	LastMessageAt              time.Time
	LastMessageExternalID      string
	CreatedAt                  time.Time `compare:"ignore"`
	UpdatedAt                  time.Time `compare:"ignore"`
}

type PostAttachment struct {
	MediaType      string
	Media          *MediaPostAttachment
	Type           string
	SubAttachments []*SubAttachment
}

type MediaPostAttachment struct {
	Image *ImageMediaPostAttachment
}

type ImageMediaPostAttachment struct {
	Height int
	Width  int
	Src    string
}

type SubAttachment struct {
	Media  *MediaDataSubAttachment
	Target *TargetDataSubAttachment
	Type   string
	URL    string
}

type ImageMediaDataSubAttachment struct {
	Image *MediaDataSubAttachment
}

type CommentAttachment struct {
	Media  *ImageMediaDataSubAttachment
	Target *TargetDataSubAttachment
	Title  string
	Type   string
	URL    string
}

type FbCustomerConversationState struct {
	ID             dot.ID
	IsRead         bool
	ExternalPageID string
	UpdatedAt      time.Time
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

type FbCreatePostArgs struct {
	ExternalPageID string
	AccessToken    string
	Message        string
}
