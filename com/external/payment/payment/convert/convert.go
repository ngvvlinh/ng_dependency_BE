package convert

import (
	"etop.vn/api/external/payment"
	paymentmodel "etop.vn/backend/com/external/payment/payment/model"
)

func Payment(in *paymentmodel.Payment) (out *payment.Payment) {
	out = &payment.Payment{
		ID:              in.ID,
		Amount:          in.Amount,
		Status:          in.Status,
		State:           in.State,
		PaymentProvider: in.PaymentProvider,
		ExternalTransID: in.ExternalTransID,
		ExternalData:    in.ExternalData,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
	return
}
