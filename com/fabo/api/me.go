package api

type Me struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	ShortName string  `json:"short_name"`
	Picture   Picture `json:"picture"`
}

type Picture struct {
	Data PictureData `json:"data"`
}

type PictureData struct {
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	IsSilhouette bool   `json:"is_silhouette"`
	Url          string `json:"url"`
}
