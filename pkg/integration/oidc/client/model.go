package client

type VerifyTokenResponse struct {
	IDTokenClaims *IDTokenClaims `json:"id_token_claims"`
}

type IDTokenClaims struct {
	UniqueName          string   `json:"unique_name"`
	PreferredUsername   string   `json:"preferred_username"`
	GivenName           string   `json:"given_name"`
	FamilyName          string   `json:"family_name"`
	Role                []string `json:"role"`
	PhoneNumber         string   `json:"phone_number"`
	PhoneNumberVerified string   `json:"phone_number_verified"`
	Email               string   `json:"email"`
	EmailVerified       string   `json:"email_verified"`
	Name                string   `json:"name"`
}
