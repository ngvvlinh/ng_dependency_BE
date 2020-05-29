package types

import (
	"time"

	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type FbExternalPost struct {
	Id                  dot.ID            `json:"id"`
	ExternalID          dot.NullString    `json:"external_id"`
	ExternalParentID    dot.NullString    `json:"external_parent_id"`
	ExternalFrom        *FbObjectFrom     `json:"external_from"`
	ExternalPicture     dot.NullString    `json:"external_picture"`
	ExternalIcon        dot.NullString    `json:"external_icon"`
	ExternalMessage     dot.NullString    `json:"external_message"`
	ExternalAttachments []*PostAttachment `json:"external_attachment"`
	ExternalCreatedTime time.Time         `json:"external_created_time"`
	CreatedAt           dot.Time          `json:"created_at"`
	UpdatedAt           dot.Time          `json:"updated_at"`
}

func (m *FbExternalPost) String() string { return jsonx.MustMarshalToString(m) }

type FbExternalConversation struct {
	ID                   dot.ID         `json:"id"`
	ExternalPageID       dot.NullString `json:"external_page_id"`
	ExternalID           dot.NullString `json:"external_id"`
	PSID                 dot.NullString `json:"psid"`
	ExternalUserID       dot.NullString `json:"external_user_id"`
	ExternalUserName     dot.NullString `json:"external_user_name"`
	ExternalLink         dot.NullString `json:"external_link"`
	ExternalUpdatedTime  time.Time      `json:"external_updated_time"`
	ExternalMessageCount dot.NullInt    `json:"external_message_count"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (m *FbExternalConversation) String() string { return jsonx.MustMarshalToString(m) }

type FbExternalComment struct {
	ID                   dot.ID             `json:"id"`
	ExternalUserID       dot.NullString     `json:"external_user_id"`
	ExternalParentID     dot.NullString     `json:"external_parent_id"`
	ExternalParentUserID dot.NullString     `json:"external_parent_user_id"`
	ExternalMessage      dot.NullString     `json:"external_message"`
	ExternalCommentCount dot.NullInt        `json:"external_comment_count"`
	ExternalFrom         *FbObjectFrom      `json:"external_from"`
	ExternalAttachment   *CommentAttachment `json:"external_attachment"`
	ExternalCreatedTime  time.Time          `json:"external_created_time"`
	ExternalID           dot.NullString     `json:"external_id"`
	CreatedAt            dot.Time           `json:"created_at"`
	UpdatedAt            dot.Time           `json:"updated_at"`
}

func (m *FbExternalComment) String() string { return jsonx.MustMarshalToString(m) }

type FbObjectFrom struct {
	ID    dot.NullString `json:"id"`
	Name  dot.NullString `json:"name"`
	Email dot.NullString `json:"email"`
}

func (m *FbObjectFrom) String() string { return jsonx.MustMarshalToString(m) }

type PostAttachment struct {
	MediaType      dot.NullString   `json:"media_type"`
	Type           dot.NullString   `json:"type"`
	SubAttachments []*SubAttachment `json:"sub_attachments"`
}

func (m *PostAttachment) String() string { return jsonx.MustMarshalToString(m) }

type CommentAttachment struct {
	Media  *ImageMediaDataSubAttachment `json:"media"`
	Target *TargetDataSubAttachment     `json:"target"`
	Title  dot.NullString               `json:"title"`
	Type   dot.NullString               `json:"type"`
	URL    dot.NullString               `json:"url"`
}

func (m *CommentAttachment) String() string { return jsonx.MustMarshalToString(m) }

type ImageMediaDataSubAttachment struct {
	Image *MediaDataSubAttachment `json:"image"`
}

func (m *ImageMediaDataSubAttachment) String() string { return jsonx.MustMarshalToString(m) }

type SubAttachment struct {
	Media  *MediaDataSubAttachment  `json:"media"`
	Target *TargetDataSubAttachment `json:"target"`
	Type   dot.NullString           `json:"type"`
	URL    dot.NullString           `json:"url"`
}

func (m *SubAttachment) String() string { return jsonx.MustMarshalToString(m) }

type MediaDataSubAttachment struct {
	Height dot.NullInt    `json:"height"`
	Width  dot.NullInt    `json:"width"`
	Src    dot.NullString `json:"src"`
}

func (m *MediaDataSubAttachment) String() string { return jsonx.MustMarshalToString(m) }

type TargetDataSubAttachment struct {
	ID  dot.NullString `json:"id"`
	URL dot.NullString `json:"url"`
}

func (m *TargetDataSubAttachment) String() string { return jsonx.MustMarshalToString(m) }

type FbExternalMessage struct {
	ID                     dot.ID                 `json:"id"`
	ExternalConversationID dot.NullString         `json:"external_conversation_id"`
	ExternalID             dot.NullString         `json:"external_id"`
	ExternalMessage        dot.NullString         `json:"external_message"`
	ExternalTo             []*FbObjectTo          `json:"external_to"`
	ExternalFrom           *FbObjectFrom          `json:"external_from"`
	ExternalAttachments    []*FbMessageAttachment `json:"external_attachments"`
	ExternalCreatedTime    time.Time              `json:"external_created_time"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

func (m *FbExternalMessage) String() string { return jsonx.MustMarshalToString(m) }

type FbMessageAttachment struct {
	ID        dot.NullString                    `json:"id"`
	ImageData *FbMessageAttachmentImageData     `json:"image_data"`
	MimeType  dot.NullString                    `json:"mime_type"`
	Name      dot.NullString                    `json:"name"`
	Size      dot.NullInt                       `json:"size"`
	VideoData *FbMessageDataAttachmentVideoData `json:"video_data"`
	FileURL   dot.NullString                    `json:"size"`
}

func (m *FbMessageAttachment) String() string { return jsonx.MustMarshalToString(m) }

type FbMessageDataAttachmentVideoData struct {
	Width      dot.NullInt    `json:"width"`
	Height     dot.NullInt    `json:"height"`
	Length     dot.NullInt    `json:"length"`
	VideoType  dot.NullInt    `json:"video_type"`
	URL        dot.NullString `json:"url"`
	PreviewURL dot.NullString `json:"preview_url"`
	Rotation   dot.NullInt    `json:"rotation"`
}

func (m *FbMessageDataAttachmentVideoData) String() string { return jsonx.MustMarshalToString(m) }

type FbMessageAttachmentImageData struct {
	Width           dot.NullInt    `json:"width"`
	Height          dot.NullInt    `json:"height"`
	MaxWidth        dot.NullInt    `json:"max_width"`
	MaxHeight       dot.NullInt    `json:"max_height"`
	URL             dot.NullString `json:"url"`
	PreviewURL      dot.NullString `json:"preview_url"`
	ImageType       dot.NullInt    `json:"image_type"`
	RenderAsSticker dot.NullBool   `json:"render_as_sticker"`
}

func (m *FbMessageAttachmentImageData) String() string { return jsonx.MustMarshalToString(m) }

type FbObjectTo struct {
	ID    dot.NullString `json:"id"`
	Name  dot.NullString `json:"name"`
	Email dot.NullString `json:"email"`
}

func (m *FbObjectTo) String() string { return jsonx.MustMarshalToString(m) }