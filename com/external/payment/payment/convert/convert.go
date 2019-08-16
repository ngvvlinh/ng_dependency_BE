package convert

import (
	"etop.vn/api/external/payment"
	etoptypes "etop.vn/api/main/etop"
	paymentmodel "etop.vn/backend/com/external/payment/payment/model"
)

func Payment(in *paymentmodel.Payment) (out *payment.Payment) {
	out = &payment.Payment{
		ID:              in.ID,
		Amount:          in.Amount,
		Status:          etoptypes.Status4FromInt(int(in.Status)),
		State:           payment.PaymentState(in.State),
		PaymentProvider: payment.PaymentProvider(in.PaymentProvider),
		ExternalTransID: in.ExternalTransID,
		ExternalData:    in.ExternalData,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
	return
}
