package model

type ConversationsResponse struct {
	Conversations *Conversations `json:"conversations"`
	ID            string         `json:"id"`
}

type Conversations struct {
	ConversationsData []*Conversation         `json:"data"`
	Paging            *FacebookPagingResponse `json:"paging"`
}

type Conversation struct {
	ID           string        `json:"id"`
	Link         string        `json:"link"`
	MessageCount int           `json:"message_count"`
	UpdatedTime  *FacebookTime `json:"updated_time"`
}
