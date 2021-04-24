package model

type LiveVideosWithCommentsResponse struct {
	LiveVideos *LiveVideosWithComments `json:"live_videos"`
}

type LiveVideosWithComments struct {
	Data []*LiveVideoWithComments `json:"data"`
}

type LiveVideoWithComments struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Video        *LiveVideoVideo    `json:"video"`
	PermalinkURL string             `json:"permalink_url"`
	From         *ObjectFrom        `json:"from"`
	EmbedHTML    string             `json:"embed_html"`
	Status       string             `json:"status"`
	Comments     *LiveVideoComments `json:"comments"`
	CreationTime *FacebookTime      `json:"creation_time"`
}

type LiveVideosResponse struct {
	Data   []*LiveVideo            `json:"data"`
	Paging *FacebookPagingResponse `json:"paging"`
}

type LiveVideo struct {
	ID           string                     `json:"id"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Video        *LiveVideoVideo            `json:"video"`
	From         *ObjectFrom                `json:"from"`
	PermalinkURL string                     `json:"permalink_url"`
	Status       string                     `json:"status"`
	Comments     *LiveVideoCommentsSummary  `json:"comments"`
	Reactions    *LiveVideoReactionsSummary `json:"reactions"`
	CreationTime *FacebookTime              `json:"creation_time"`
}

func (s LiveVideo) GetExternalPostID() string {
	if s.Video == nil || s.From == nil {
		return ""
	}

	return s.From.ID + "_" + s.Video.ID
}

type LiveVideoVideo struct {
	ID      string `json:"id"`
	Picture string `json:"picture"`
	Source  string `json:"source"`
}

type LiveVideoComments struct {
	Data    []*Comment         `json:"data"`
	Summary *SummaryTotalCount `json:"summary"`
}

type LiveVideoCommentsSummary struct {
	Summary *SummaryTotalCount `json:"summary"`
}

type LiveVideoReactionsSummary struct {
	Summary *SummaryTotalCount `json:"summary"`
}

type SummaryTotalCount struct {
	TotalCount int `json:"total_count"`
}
