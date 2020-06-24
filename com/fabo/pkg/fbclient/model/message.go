package model

import (
	"encoding/json"
	"fmt"
)

type MessagesResponse struct {
	ID       string    `json:"id"`
	Messages *Messages `json:"messages"`
}

type Messages struct {
	MessagesData []*MessageData          `json:"data"`
	Paging       *FacebookPagingResponse `json:"paging"`
}

type MessageData struct {
	ID          string                  `json:"id"`
	CreatedTime *FacebookTime           `json:"created_time"`
	Message     string                  `json:"message"`
	To          *ObjectsTo              `json:"to"`
	From        *ObjectFrom             `json:"from"`
	Sticker     string                  `json:"sticker"`
	Attachments *MessageDataAttachments `json:"attachments"`
}

type MessageDataAttachments struct {
	Data   []*MessageDataAttachment `json:"data"`
	Paging *FacebookPagingResponse  `json:"paging"`
}

type MessageDataAttachment struct {
	ID        string                          `json:"id"`
	ImageData *MessageDataAttachmentImage     `json:"image_data"`
	MimeType  string                          `json:"mime_type"`
	Name      string                          `json:"name"`
	Size      int                             `json:"size"`
	VideoData *MessageDataAttachmentVideoData `json:"video_data"`
	FileURL   string                          `json:"file_url"`
}

type MessageDataAttachmentImage struct {
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	MaxWidth        int    `json:"max_width"`
	MaxHeight       int    `json:"max_height"`
	URL             string `json:"url"`
	PreviewURL      string `json:"preview_url"`
	ImageType       int    `json:"image_type"`
	RenderAsSticker bool   `json:"render_as_sticker"`
}

type MessageDataAttachmentVideoData struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Length     int    `json:"length"`
	VideoType  int    `json:"video_type"`
	URL        string `json:"url"`
	PreviewURL string `json:"preview_url"`
	Rotation   int    `json:"rotation"`
}

type SubcribedAppResponse struct {
	Success bool `json:"success"`
}

type SendMessageRequest struct {
	Recipient *RecipientSendMessageRequest `json:"recipient"`
	Message   *MessageSendMessageRequest   `json:"message"`
}

type RecipientSendMessageRequest struct {
	ID string `json:"id"`
}

type MessageSendMessageRequest struct {
	Text       string                        `json:"text"`
	Attachment *AttachmentSendMessageRequest `json:"attachment"`
}

func (m MessageSendMessageRequest) MarshalJSON() ([]byte, error) {
	if m.Text != "" {
		return []byte(fmt.Sprintf(`{"text": "%s"}`, m.Text)), nil
	}
	attachmentJSON, err := json.Marshal(m.Attachment)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(`{"attachment": %s}`, string(attachmentJSON))), nil
}

type AttachmentSendMessageRequest struct {
	Type    string                              `json:"type"`
	Payload PayloadAttachmentSendMessageRequest `json:"payload"`
}

type PayloadAttachmentSendMessageRequest struct {
	Url        string `json:"url"`
	IsReusable bool   `json:"is_reusable"`
}

type SendMessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type SendCommentRequest struct {
	ID            string `json:"id"`
	Message       string `json:"message"`
	AttachmentURL string `json:"attachment_url"`
}

type SendCommentResponse struct {
	ID string `json:"id"`
}

type CreatePostRequest struct {
	Message string `json:"message"`
}

type CreatePostResponse struct {
	ID string `json:"id"`
}
