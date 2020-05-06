package model

type CommentsByPostIDsResponse struct {
	Data map[string]*Comment
}

type CommentsResponse struct {
	CommentData []*Comment              `json:"data"`
	Paging      *FacebookPagingResponse `json:"paging"`
}

type Comment struct {
	ID           string             `json:"id"`
	Message      string             `json:"message"`
	CommentCount int                `json:"comment_count"`
	From         *ObjectFrom        `json:"from"`
	Attachment   *CommentAttachment `json:"attachment"`
	Parent       *CommentParent     `json:"parent"`
}

type CommentFrom struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type CommentAttachment struct {
	Media  *CommentAttachmentMedia  `json:"media"`
	Target *CommentAttachmentTarget `json:"target"`
	Type   string                   `json:"type"`
	URL    string                   `json:"url"`
}

type CommentAttachmentMedia struct {
	Image *CommentAttachmentMediaImage `json:"image"`
}

type CommentAttachmentMediaImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type CommentAttachmentTarget struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type CommentParent struct {
	CreatedTime *FacebookTime `json:"created_time"`
	From        *ObjectFrom   `json:"from"`
	Message     string        `json:"message"`
	ID          string        `json:"id"`
}
