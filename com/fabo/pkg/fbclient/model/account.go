package model

import "o.o/backend/com/fabo/main/fbpage/model"

type AccountsResponse struct {
	Accounts    Accounts            `json:"accounts"`
	Permissions AccountsPermissions `json:"permissions"`
	Id          string              `json:"id"`
}

type Accounts struct {
	Data   []AccountData          `json:"data"`
	Paging FacebookPagingResponse `json:"paging"`
}

type AccountData struct {
	AccessToken  string                   `json:"access_token"`
	Category     string                   `json:"category"`
	CategoryList []model.ExternalCategory `json:"category_list"`
	Name         string                   `json:"name"`
	Id           string                   `json:"id"`
	Tasks        []string                 `json:"tasks"`
	FanCount     int                      `json:"fan_count"`
	Picture      Picture                  `json:"picture"`
}

type AccountsPermissions struct {
	Data []AccountPermission `json:"data"`
}

type AccountPermission struct {
	Permission string `json:"permission"`
	Status     string `json:"status"`
}
