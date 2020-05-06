package model

type CommentSummariesResponse struct {
	MapCommentSummary map[string]*CommentSummary
}

type CommentSummary struct {
	Data    []interface{}          `json:"data"`
	Summary *SummaryCommentSummary `json:"summary"`
}

type SummaryCommentSummary struct {
	Order      string `json:"order"`
	TotalCount int    `json:"total_count"`
	CanComment bool   `json:"can_comment"`
}
