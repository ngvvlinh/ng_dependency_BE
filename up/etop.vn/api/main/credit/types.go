package credit

import (
	creditv1 "etop.vn/api/main/credit/v1"
	"etop.vn/api/meta"
)

type CreditEvent = creditv1.CreditEvent
type EventType = creditv1.EventDataEnum
type IsEventData = creditv1.IsEventData

func NewCreditEvent(
	correlationID meta.UUID,
	accountID int64,
	walletType string,
	data IsEventData,
) *CreditEvent {
	meta.AutoFill(&correlationID)
	return &CreditEvent{
		Id:            0,
		Uuid:          meta.NewUUID(),
		CorrelationId: correlationID,
		Type:          int32(data.GetEnumTag()),
		Data:          &creditv1.EventData{Data: data},
	}
}
