package model

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type ObjectTo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Picture   *Picture `json:"picture"`
}

type ObjectsTo struct {
	Data []*ObjectTo `json:"data"`
}

type ObjectFrom struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Picture   *Picture `json:"picture"`
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
