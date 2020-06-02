package model

type Profile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
	Locale     string `json:"locale"`
	Timezone   int    `json:"timezone"`
	Gender     string `json:"gender"`
}
