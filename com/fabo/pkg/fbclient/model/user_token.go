package model

type UserToken struct {
	Data *Data `json:"data"`
}

type Data struct {
	AppID               string           `json:"app_id"`
	Type                string           `json:"type"`
	Application         string           `json:"application"`
	DataAccessExpiresAt int              `json:"data_access_expires_at"`
	ExpiresAt           int              `json:"expires_at"`
	IsValid             bool             `json:"is_valid"`
	Scopes              []string         `json:"scopes"`
	GranularScopes      []*GranularScope `json:"granular_scopes"`
	UserID              string           `json:"user_id"`
}

type GranularScope struct {
	Scope     string   `json:"scope"`
	TargetIDs []string `json:"target_ids"`
}
