package model

import (
	"fmt"
	"strings"
)

type PublishedPostsResponse struct {
	Data   []*Post                 `json:"data"`
	Paging *FacebookPagingResponse `json:"paging"`
}

type Post struct {
	ID               string            `json:"id"`
	From             *ObjectFrom       `json:"from"`
	FullPicture      string            `json:"full_picture"`
	Icon             string            `json:"icon"`
	IsExpired        bool              `json:"is_expired"`
	IsHidden         bool              `json:"is_hidden"`
	IsPopular        bool              `json:"is_popular"`
	IsPublished      bool              `json:"is_published"`
	Message          string            `json:"message"`
	Story            string            `json:"story"`
	PermalinkURL     string            `json:"permalink_url"`
	Shares           *Shares           `json:"shares"`
	StatusType       string            `json:"status_type"`
	Picture          string            `json:"picture"`
	Attachments      *Attachments      `json:"attachments"`
	CommentsSummary  *CommentsSummary  `json:"comments"`
	ReactionsSummary *ReactionsSummary `json:"reactions"`
	CreatedTime      FacebookTime      `json:"created_time"`
	UpdatedTime      FacebookTime      `json:"updated_time"`
}

// IsResourceFromCurrentPage check current post is share or not.
func (p *Post) IsResourceFromCurrentPage() bool {
	splited := strings.Split(p.ID, "_")
	if len(splited) < 2 {
		return false
	}
	currentPageID := splited[0]

	if p.Attachments == nil {
		return false
	}

	if p.Attachments.Data == nil || len(p.Attachments.Data) == 0 {
		return false
	}

	data := p.Attachments.Data[0]
	if data.Target != nil {
		prefixUrl := fmt.Sprintf("https://www.facebook.com/%v/posts", currentPageID)
		if strings.HasPrefix(data.Target.Url, prefixUrl) {
			return true
		}
	}
	return false
}

type Shares struct {
	Count int `json:"count"`
}

type Attachments struct {
	Data []*DataAttachment `json:"data"`
}

type DataAttachment struct {
	Media          *MediaPostAttachment `json:"media"`
	MediaType      string               `json:"media_type"`
	Type           string               `json:"type"`
	Title          string               `json:"title"`
	Target         *DataTarget          `json:"target"`
	SubAttachments *SubAttachments      `json:"subattachments"`
}

type DataTarget struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type MediaPostAttachment struct {
	Image *ImageMediaPostAttachment `json:"image"`
}

type ImageMediaPostAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type SubAttachments struct {
	Data []*DataSubAttachment `json:"data"`
}

type DataSubAttachment struct {
	Description string                   `json:"description"`
	Media       *MediaDataSubAttachment  `json:"media"`
	Target      *TargetDataSubAttachment `json:"target"`
	Type        string                   `json:"type"`
	URL         string                   `json:"url"`
}

type MediaDataSubAttachment struct {
	Image *ImageMediaDataSubAttachment `json:"image"`
}

type ImageMediaDataSubAttachment struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Src    string `json:"src"`
}

type TargetDataSubAttachment struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type PostsWithCommentsResponse struct {
	Feeds *PostsWithComments `json:"feed"`
}

type PostsWithComments struct {
	Data []*PostWithComments `json:"data"`
}

type PostWithComments struct {
	Post
	Comments *Comments `json:"comments"`
}
