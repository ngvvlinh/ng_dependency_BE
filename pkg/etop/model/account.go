package model

type UpdateAccountURLSlugCommand struct {
	AccountID int64
	URLSlug   string
}

type GetAccountAuthQuery struct {
	AuthKey     string
	AccountType AccountType
	AccountID   int64

	Result struct {
		AccountAuth *AccountAuth
		Account     AccountInterface
	}
}
