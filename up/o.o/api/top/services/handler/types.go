package handler

import (
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type ResetStateRequest struct {
	WebhookId dot.ID `json:"webhook_id"`
	AccountId dot.ID `json:"account_id"`
}

func (m *ResetStateRequest) String() string { return jsonx.MustMarshalToString(m) }
