package api

import "etop.vn/backend/com/fabo/main/fbpage/model"

type Accounts struct {
	Accounts Account `json:"accounts"`
	Id       string  `json:"id"`
}

type Account struct {
	Data   []AccountData  `json:"data"`
	Paging FacebookPaging `json:"paging"`
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
