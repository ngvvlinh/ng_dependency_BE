package types

type BankAccount struct {
	Name          string
	Province      string
	Branch        string
	AccountNumber string
	AccountName   string
}

type CompanyInfo struct {
	Name                string
	TaxCode             string
	Address             string
	Website             string
	LegalRepresentative *ContactPerson
}

type ContactPerson struct {
	Name     string
	Position string
	Phone    string
	Email    string
}
