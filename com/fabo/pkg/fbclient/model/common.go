package model

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type From struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
