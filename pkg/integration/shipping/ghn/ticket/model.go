package ticket

import "time"

type CreateTicketReplyRequest struct {
	TicketID    string `url:"ticket_id"`
	UserID      string `url:"user_id"`
	Description string `url:"description"`
}

type CreateTicketReplyResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		FromEmail string    `json:"from_email"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    string    `json:"user_id"`
	} `json:"data"`
}

type CreateTicketRequest struct {
	CEmail      string `url:"c_email"`
	OrderCode   string `url:"order_code"`
	Category    string `url:"category"`
	Description string `url:"description"`
}

type CreateTicketResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Attachments   []string  `json:"attachments"`
		ClientID      string    `json:"client_id"`
		Conversations []string  `json:"conversations"`
		CreatedAt     time.Time `json:"created_at"`
		Description   string    `json:"description"`
		CreatedBy     string    `json:"created_by"`
		ID            string    `json:"id"`
		OrderCode     string    `json:"order_code"`
		Status        string    `json:"status"`
		StatusID      int       `json:"status_id"`
		Type          string    `json:"type"`
		UpdatedAt     time.Time
	} `json:"data"`
}
