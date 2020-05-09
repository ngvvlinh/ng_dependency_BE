package model

type CommentsByPostIDsResponse struct {
	Data map[string]*Comment
}

type CommentsResponse struct {
	Comments *Comments `json:"comments"`
	PostID   string    `json:"id"`
}

type Comments struct {
	CommentData []*Comment              `json:"data"`
	Paging      *FacebookPagingResponse `json:"paging"`
}

type Comment struct {
	ID           string             `json:"id"`
	IsHidden     bool               `json:"is_hidden"`
	Message      string             `json:"message"`
	From         *ObjectFrom        `json:"from"`
	CommentCount int                `json:"comment_count"`
	Attachment   *CommentAttachment `json:"attachment"`
	Parent       *ObjectParent      `json:"parent"`
	CreatedTime  *FacebookTime      `json:"created_time"`
}

type CommentFrom struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type CommentAttachment struct {
	Media  *CommentAttachmentMedia  `json:"media"`
	Target *CommentAttachmentTarget `json:"target"`
	Title  string                   `json:"title"`
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
