package fbmessaging

import (
	"time"

	"o.o/api/fabo/fbmessaging/fb_comment_source"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_status_type"
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
	ExternalMessageShares  []*FbMessageShare
	InternalSource         fb_internal_source.FbInternalSource
	ExternalCreatedTime    time.Time
	ExternalTimestamp      int64
	CreatedAt              time.Time `compare:"ignore"`
	UpdatedAt              time.Time `compare:"ignore"`
	CreatedBy              dot.ID
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

type FbMessageShare struct {
	ID          string
	Description string
	Name        string
	Link        string
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
	Source               fb_comment_source.FbCommentSource
	IsLiked              bool
	IsHidden             bool
	IsPrivateReplied     bool
	CreatedAt            time.Time `compare:"ignore"`
	UpdatedAt            time.Time `compare:"ignore"`
	DeletedAt            time.Time
	InternalSource       fb_internal_source.FbInternalSource
	CreatedBy            dot.ID
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
	ID                   dot.ID    `json:"id"`
	ExternalPageID       string    `json:"external_page_id"`
	ExternalID           string    `json:"external_id"`
	PSID                 string    `json:"psid"`
	ExternalUserID       string    `json:"external_user_id"`
	ExternalUserName     string    `json:"external_user_name"`
	ExternalLink         string    `json:"external_link"`
	ExternalUpdatedTime  time.Time `json:"external_updated_time"`
	ExternalMessageCount int       `json:"external_message_count"`
	CreatedAt            time.Time `compare:"ignore" json:"created_at"`
	UpdatedAt            time.Time `compare:"ignore" json:"updated_at"`
}

type FbExternalPosts []*FbExternalPost

type FbExternalPost struct {
	ID                      dot.ID
	ExternalPageID          string
	ExternalID              string
	ExternalParentID        string
	ExternalFrom            *FbObjectFrom
	ExternalPicture         string
	ExternalIcon            string
	ExternalMessage         string
	ExternalAttachments     []*PostAttachment `compare:"ignore"`
	ExternalCreatedTime     time.Time
	CreatedAt               time.Time `compare:"ignore"`
	UpdatedAt               time.Time `compare:"ignore"`
	DeletedAt               time.Time `compare:"ignore"`
	ExternalParent          *FbExternalPost
	FeedType                fb_feed_type.FbFeedType
	StatusType              fb_status_type.FbStatusType
	TotalComments           int
	TotalReactions          int
	IsLiveVideo             bool
	ExternalLiveVideoStatus string
	LiveVideoStatus         fb_live_video_status.FbLiveVideoStatus
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
	LastCustomerMessageAt      time.Time
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
	Description string
	Media       *MediaDataSubAttachment
	Target      *TargetDataSubAttachment
	Type        string
	URL         string
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

type FbExternalCommentCreatedOrUpdatedEvent struct {
	FbExternalComment *FbExternalComment
}

type FbExternalConversationsCreatedEvent struct {
	FbExternalConversations []*FbExternalConversation
}

type FbExternalConversationsUpdatedEvent struct {
	FbExternalConversations []*FbExternalConversation
}
