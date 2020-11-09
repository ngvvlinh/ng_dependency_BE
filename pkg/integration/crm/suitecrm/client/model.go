package client

import "o.o/backend/pkg/common/apifw/httpreq"

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type SuiteCRMCfg struct {
	Token string `yaml:"token"`
}

type InsertCaseRequest struct {
	PhoneMobile string `json:"phoneMobile"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	RefType     string `json:"refType"`
	RefID       string `json:"refId"`
	RefCompany  string `json:"refCompany"`
}

type InsertCaseResponse struct {
	Message String `json:"message"`
	Error   string `json:"error"`
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
