package model

type LiveVideosResponse struct {
	LiveVideos *LiveVideos `json:"live_videos"`
}

type LiveVideos struct {
	Data []*LiveVideo `json:"data"`
}

type LiveVideo struct {
	ID           string             `json:"id"`
	Video        *LiveVideoVideo    `json:"video"`
	PermalinkURL string             `json:"permalink_url"`
	From         *ObjectFrom        `json:"from"`
	EmbedHTML    string             `json:"embed_html"`
	Status       string             `json:"status"`
	Comments     *LiveVideoComments `json:"comments"`
	CreationTime *FacebookTime      `json:"creation_time"`
}

type LiveVideoVideo struct {
	ID      string `json:"id"`
	Picture string `json:"picture"`
	Source  string `json:"source"`
}

type LiveVideoComments struct {
	Data    []*Comment               `json:"data"`
	Summary *LiveVideoCommentSummary `json:"summary"`
}

type LiveVideoCommentSummary struct {
	TotalCount int64 `json:"total_count"`
}
