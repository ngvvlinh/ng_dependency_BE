package sharemodel

type BankAccount struct {
	Name          string `json:"name"`
	Province      string `json:"province"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

type CompanyInfo struct {
	Name                string         `json:"name"`
	TaxCode             string         `json:"tax_code"`
	Address             string         `json:"address"`
	Website             string         `json:"website"`
	LegalRepresentative *ContactPerson `json:"legal_representative"`
}

type ContactPerson struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type AdjustmentLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}

type FeeLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}

type DiscountLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}
