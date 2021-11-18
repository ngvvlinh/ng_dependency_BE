package model

import (
	"time"

	"o.o/api/fabo/fbmessaging/fb_comment_source"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/capi/dot"
)

// +sqlgen
type FbExternalMessage struct {
	ID                     dot.ID `paging:"id"`
	ExternalConversationID string
	ExternalPageID         string
	ExternalID             string
	ExternalMessage        string
	ExternalSticker        string
	ExternalTo             []*FbObjectTo
	ExternalFrom           *FbObjectFrom
	ExternalFromID         string
	ExternalAttachments    []*FbMessageAttachment
	ExternalMessageShares  []*FbMessageShare
	ExternalCreatedTime    time.Time `paging:"external_created_time"`
	// Webhook fb trả vê timestamp sẽ phân biệt (sort) tốt hơn external_created_time
	// vì khi parse timestamp sang time.Time sẽ mất phần milisecond trong timestamp
	ExternalTimestamp int64
	InternalSource    fb_internal_source.FbInternalSource

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
	CreatedBy dot.ID
}

type FbObjectTo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageURL  string `json:"image_url"`
}

type FbObjectFrom struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageURL  string `json:"image_url"`
}

type FbMessageAttachment struct {
	ID        string                            `json:"id"`
	ImageData *FbMessageAttachmentImageData     `json:"image_data"`
	MimeType  string                            `json:"mime_type"`
	Name      string                            `json:"name"`
	Size      int                               `json:"size"`
	VideoData *FbMessageDataAttachmentVideoData `json:"video_data"`
	FileURL   string                            `json:"file_url"`
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

type FbMessageDataAttachmentVideoData struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Length     int    `json:"length"`
	VideoType  int    `json:"video_type"`
	URL        string `json:"url"`
	PreviewURL string `json:"preview_url"`
	Rotation   int    `json:"rotation"`
}

type FbMessageShare struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Link        string `json:"link"`
}

// +sqlgen
type FbExternalConversation struct {
	ID                   dot.ID
	ExternalPageID       string
	ExternalID           string
	PSID                 string // page scope ID
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
	ID                   dot.ID `paging:"id"`
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
	ExternalCreatedTime  time.Time `paging:"external_created_time"`
	Source               fb_comment_source.FbCommentSource
	InternalSource       fb_internal_source.FbInternalSource
	IsLiked              bool      // like by current page
	IsHidden             bool      // hide by current page
	IsPrivateReplied     bool      // reply by current page
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
	CreatedBy            dot.ID

	ExternalOwnerPostID string // user create post
	PostType            fb_post_type.FbPostType
}

type CommentAttachment struct {
	Media  *ImageMediaDataSubAttachment `json:"media"`
	Target *TargetDataSubAttachment     `json:"target"`
	Title  string                       `json:"title"`
	Type   string                       `json:"type"`
	URL    string                       `json:"url"`
}

type ImageMediaDataSubAttachment struct {
	Image *MediaDataSubAttachment `json:"image"`
}

type FbObjectParent struct {
	CreatedTime time.Time     `json:"created_time"`
	From        *FbObjectFrom `json:"from"`
	Message     string        `json:"message"`
	ID          string        `json:"id"`
}

// +sqlgen
type FbExternalPost struct {
	ID                  dot.ID
	ExternalPageID      string // external_page_id when post was created by page
	ExternalUserID      string // external_user_id when post was created by user
	ExternalID          string
	ExternalParentID    string
	ExternalFrom        *FbObjectFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*PostAttachment
	ExternalCreatedTime time.Time `paging:"external_created_time"`
	ExternalUpdatedTime time.Time
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
	TotalComments       int
	TotalReactions      int
	Type                fb_post_type.FbPostType
	FeedType            fb_feed_type.FbFeedType
	StatusType          fb_status_type.FbStatusType

	IsLiveVideo             bool
	ExternalLiveVideoStatus string
	LiveVideoStatus         fb_live_video_status.FbLiveVideoStatus
}

// +sqlsel
type FbExternalPostFtTotalComment struct {
	ExternalPostID string `sel:"external_post_id"`
	Count          int    `sel:"count(id)"`
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
	ID                         dot.ID `paging:"id"`
	ExternalPageID             string
	ExternalOwnerPostID        string // user create post
	ExternalID                 string
	ExternalUserID             string
	ExternalUserName           string
	ExternalFrom               *FbObjectFrom
	ExternalPostAttachments    []*PostAttachment
	ExternalCommentAttachment  *CommentAttachment
	ExternalMessageAttachments []*FbMessageAttachment
	Type                       fb_customer_conversation_type.FbCustomerConversationType
	LastMessage                string
	LastMessageAt              time.Time `paging:"last_message_at"`
	LastCustomerMessageAt      time.Time
	LastMessageExternalID      string
	CreatedAt                  time.Time `sq:"create"`
	UpdatedAt                  time.Time `sq:"update"`
	DeletedAt                  time.Time
}

type PostAttachment struct {
	Media          *MediaPostAttachment `json:"media"`
	MediaType      string               `json:"media_type"`
	Type           string               `json:"type"`
	SubAttachments []*SubAttachment     `json:"sub_attachments"`
}

type MediaPostAttachment struct {
	Image *ImageMediaPostAttachment `json:"image"`
}

type ImageMediaPostAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

// +sqlgen
type FbCustomerConversationState struct {
	ID             dot.ID
	IsRead         bool
	ExternalPageID string
	UpdatedAt      time.Time `sq:"update"`
}
