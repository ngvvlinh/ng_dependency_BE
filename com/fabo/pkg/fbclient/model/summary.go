package model

type CommentsSummary struct {
	Data    []*Comment              `json:"data"`
	Summary *SummaryCommentsSummary `json:"summary"`
}

type SummaryCommentsSummary struct {
	Order      string `json:"order"`
	TotalCount int    `json:"total_count"`
	CanComment bool   `json:"can_comment"`
}

type ReactionsSummary struct {
	Data    []interface{}            `json:"data"`
	Summary *SummaryReactionsSummary `json:"summary"`
}

type SummaryReactionsSummary struct {
	TotalCount int `json:"total_count"`
}
