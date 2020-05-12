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
	Attachments *MessageDataAttachments `json:"attachments"`
}

type MessageDataAttachments struct {
	Data   []*MessageDataAttachment `json:"data"`
	Paging *FacebookPagingResponse  `json:"paging"`
}

type MessageDataAttachment struct {
	ID        string                      `json:"id"`
	ImageData *MessageDataAttachmentImage `json:"image_data"`
	MimeType  string                      `json:"mime_type"`
	Name      string                      `json:"name"`
	Size      int                         `json:"size"`
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