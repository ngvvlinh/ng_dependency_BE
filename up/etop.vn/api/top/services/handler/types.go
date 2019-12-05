package handler

import (
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type ResetStateRequest struct {
	WebhookId dot.ID `json:"webhook_id"`
	AccountId dot.ID `json:"account_id"`
}

func (m *ResetStateRequest) Reset()         { *m = ResetStateRequest{} }
func (m *ResetStateRequest) String() string { return jsonx.MustMarshalToString(m) }
