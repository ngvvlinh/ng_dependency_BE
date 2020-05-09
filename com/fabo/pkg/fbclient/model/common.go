package model

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type ObjectTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type ObjectsTo struct {
	Data []*ObjectTo `json:"data"`
}

type ObjectFrom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ObjectsFrom struct {
	Data []*ObjectFrom `json:"data"`
}

type ObjectParent struct {
	CreatedTime *FacebookTime `json:"created_time"`
	From        *ObjectFrom   `json:"from"`
	Message     string        `json:"message"`
	ID          string        `json:"id"`
}
