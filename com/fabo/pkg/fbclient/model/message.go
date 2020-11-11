package model

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
	Shares      *MessageDataShares      `json:"shares"`
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

type MessageDataShares struct {
	Data   []*MessageDataShare     `json:"data"`
	Paging *FacebookPagingResponse `json:"paging"`
}

type MessageDataShare struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
}

type SubscribedAppResponse struct {
	Success bool `json:"success"`
}

type SendMessageArgs struct {
	Recipient *RecipientSendMessageRequest `json:"recipient"`
	Message   *MessageSendMessageRequest   `json:"message"`
	Tag       string                       `json:"tag"`
}

type RecipientSendMessageRequest struct {
	ID        string `json:"id,omitempty"`
	CommentID string `json:"comment_id,omitempty"`
}

type MessageSendMessageRequest struct {
	Text       string                        `json:"text,omitempty"`
	Attachment *AttachmentSendMessageRequest `json:"attachment,omitempty"`
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

type SendCommentArgs struct {
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

type CommonResponse struct {
	Success bool `json:"success"`
}
