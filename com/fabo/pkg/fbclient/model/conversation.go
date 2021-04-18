package model

type ConversationsResponse struct {
	ConversationsData []*Conversation         `json:"data"`
	Paging            *FacebookPagingResponse `json:"paging"`
}

type Conversation struct {
	ID           string               `json:"id"`
	Link         string               `json:"link"`
	MessageCount int                  `json:"message_count"`
	Senders      *ConversationSenders `json:"senders"`
	UpdatedTime  *FacebookTime        `json:"updated_time"`
}

type ConversationSenders struct {
	Data []*ConversationSender `json:"data"`
}

type ConversationSender struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
