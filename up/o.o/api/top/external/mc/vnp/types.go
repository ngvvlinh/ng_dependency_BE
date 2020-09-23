package vnp

import (
	vnpentitytype "o.o/api/top/external/mc/vnp/etc/entity_type"
	"o.o/api/top/external/types"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type CreateWebhookRequest struct {
	Entities []vnpentitytype.EntityType `json:"entities"`
	URL      string                     `json:"url"`
}

func (m *CreateWebhookRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type Webhook struct {
	ID        dot.ID                     `json:"id"`
	Entities  []vnpentitytype.EntityType `json:"entities"`
	URL       string                     `json:"url"`
	CreatedAt dot.Time                   `json:"created_at"`
	States    *types.WebhookStates       `json:"states"`
}

func (m *Webhook) String() string {
	return jsonx.MustMarshalToString(m)
}

type WebhooksResponse struct {
	Webhooks []*Webhook `json:"webhooks"`
}

func (m *WebhooksResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type DataCallback struct {
	ID      dot.ID    `json:"id"`
	Changes []*Change `json:"changes"`
}

func (m *DataCallback) String() string { return jsonx.MustMarshalToString(m) }

type Change struct {
	Time       dot.Time                 `json:"time"`
	ChangeType string                   `json:"change_type"`
	Entity     vnpentitytype.EntityType `json:"entity"`
	// Giá trị hiện tại của đối tượng (entity)
	// Chứa đầy đủ data
	Lastest *LastestOneOf `json:"lastest"`
	// Giá trị thay đổi của đối tượng (entity)
	// Chỉ chứa những trường có thay đổi
	Changed *ChangeOneOf `json:"changed"`
}

func (m *Change) String() string { return jsonx.MustMarshalToString(m) }

type LastestOneOf struct {
	ShipnowFulfillment *types.ShipnowFulfillment `json:"shipnow_fulfillment,omitempty"`
}

func (m *LastestOneOf) String() string { return jsonx.MustMarshalToString(m) }

type ChangeOneOf struct {
	ShipnowFulfillment *types.ShipnowFulfillment `json:"shipnow_fulfillment,omitempty"`
}

func (m *ChangeOneOf) String() string { return jsonx.MustMarshalToString(m) }
