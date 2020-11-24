package types

import (
	exttypes "o.o/api/top/external/types"
	"o.o/capi/dot"
)

type FaboEvent struct {
	PgEventComment              *PgEventComment
	PgEventConversation         *PgEventConversation
	PgEventMessage              *PgEventMessage
	PgEventCustomerConversation *PgEventCustomerConversation
	Timestamp                   int64 `json:"t"`
}

type PgEventComment struct {
	Op             string                      `json:"op"`
	FbPageID       dot.ID                      `json:"fb_page_id"`
	ShopID         dot.ID                      `json:"shop_id"`
	FbEventComment *exttypes.FbExternalComment `json:"fb_comment"`
}

type PgEventConversation struct {
	Op                  string                           `json:"op"`
	FbPageID            dot.ID                           `json:"fb_page_id"`
	ShopID              dot.ID                           `json:"shop_id"`
	FbEventConversation *exttypes.FbExternalConversation `json:"fb_conversation"`
}

type PgEventMessage struct {
	Op             string                      `json:"op"`
	FbPageID       dot.ID                      `json:"fb_page_id"`
	ShopID         dot.ID                      `json:"shop_id"`
	FbEventMessage *exttypes.FbExternalMessage `json:"fb_message"`
}

type PgEventCustomerConversation struct {
	Op                          string                           `json:"op"`
	FbPageID                    dot.ID                           `json:"fb_page_id"`
	ShopID                      dot.ID                           `json:"shop_id"`
	FbEventCustomerConversation *exttypes.FbCustomerConversation `json:"fb_customer_conversation"`
}
