package gateway

import (
	"o.o/backend/pkg/common/apifw/httpreq"
	kpayclient "o.o/backend/pkg/integration/payment/kpay/client"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type CallbackTransaction struct {
	Transaction        String                   `json:"transaction"`
	PartnerTransaction String                   `json:"partnerTransaction"`
	Amount             Int                      `json:"amount"`
	Fee                Int                      `json:"fee"`
	Total              Int                      `json:"total"`
	State              kpayclient.PaymentStatus `json:"state"`
	Desc               String                   `json:"desc"`
	CreatedAt          Time                     `json:"createdAt"`
	UpdatedAt          Time                     `json:"updatedAt"`
}
