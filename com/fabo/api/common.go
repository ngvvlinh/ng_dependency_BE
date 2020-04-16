package api

type FacebookPaging struct {
	Cursors  Cursors `json:"cursors"`
	Previous string  `json:"previous"`
	Next     string  `json:"next"`
}

type Cursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
