package types

import (
	"o.o/api/top/types/etc/entity_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type CreateWebhookRequest struct {
	Entities []entity_type.EntityType `json:"entities"`
	Fields   []string                 `json:"fields"`
	Url      string                   `json:"url"`
	Metadata string                   `json:"metadata"`
}

func (m *CreateWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteWebhookRequest struct {
	Id dot.ID `json:"id"`
}

func (m *DeleteWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }

type WebhooksResponse struct {
	Webhooks []*Webhook `json:"webhooks"`
}

func (m *WebhooksResponse) String() string { return jsonx.MustMarshalToString(m) }

type Webhook struct {
	Id        dot.ID                   `json:"id"`
	Entities  []entity_type.EntityType `json:"entities"`
	Fields    []string                 `json:"fields"`
	Url       string                   `json:"url"`
	Metadata  string                   `json:"metadata"`
	CreatedAt dot.Time                 `json:"created_at"`
	States    *WebhookStates           `json:"states"`
}

func (m *Webhook) String() string { return jsonx.MustMarshalToString(m) }

type WebhookStates struct {
	State      string        `json:"state"`
	LastSentAt dot.Time      `json:"last_sent_at"`
	LastError  *WebhookError `json:"last_error"`
}

func (m *WebhookStates) String() string { return jsonx.MustMarshalToString(m) }

type WebhookError struct {
	Error      string   `json:"error"`
	RespStatus int      `json:"resp_status"`
	RespBody   string   `json:"resp_body"`
	Retried    int      `json:"retried"`
	SentAt     dot.Time `json:"sent_at"`
}

func (m *WebhookError) String() string { return jsonx.MustMarshalToString(m) }

type Callback struct {
	Id      dot.ID    `json:"id"`
	Changes []*Change `json:"changes"`
}

func (m *Callback) String() string { return jsonx.MustMarshalToString(m) }

// ChangesData serialize changes data for storing in MongoDB
type ChangesData struct {
	// for using with mongodb
	XId       dot.ID    `json:"_id"`
	WebhookId dot.ID    `json:"webhook_id"`
	AccountId dot.ID    `json:"account_id"`
	CreatedAt dot.Time  `json:"created_at"`
	Changes   []*Change `json:"changes"`
}

func (m *ChangesData) String() string { return jsonx.MustMarshalToString(m) }

type Change struct {
	Time       dot.Time     `json:"time"`
	ChangeType string       `json:"change_type"`
	Entity     string       `json:"entity"`
	Latest     *LatestOneOf `json:"latest"`
	Changed    *ChangeOneOf `json:"changed"`
}

func (m *Change) String() string { return jsonx.MustMarshalToString(m) }

type LatestOneOf struct {
	Order                         *Order                         `json:"order,omitempty"`
	Fulfillment                   *Fulfillment                   `json:"fulfillment,omitempty"`
	Variant                       *ShopVariant                   `json:"variant,omitempty"`
	InventoryLevel                *InventoryLevel                `json:"inventory_level,omitempty"`
	CustomerAddress               *CustomerAddress               `json:"customer_address,omitempty"`
	Customer                      *Customer                      `json:"customer,omitempty"`
	CustomerGroup                 *CustomerGroup                 `json:"customer_group,omitempty"`
	CustomerGroupRelationship     *CustomerGroupRelationship     `json:"customer_group_relationship,omitempty"`
	Product                       *ShopProduct                   `json:"product,omitempty"`
	ProductCollection             *ProductCollection             `json:"product_collection,omitempty"`
	ProductCollectionRelationship *ProductCollectionRelationship `json:"product_collection_relationship,omitempty"`
	ShipnowFulfillment            *ShipnowFulfillment            `json:"shipnow_fulfillment,omitempty"`
}

func (m *LatestOneOf) String() string { return jsonx.MustMarshalToString(m) }

type ChangeOneOf struct {
	Order                         *Order                         `json:"order,omitempty"`
	Fulfillment                   *Fulfillment                   `json:"fulfillment,omitempty"`
	Product                       *ShopProduct                   `json:"product,omitempty"`
	Variant                       *ShopVariant                   `json:"variant,omitempty"`
	Customer                      *Customer                      `json:"customer,omitempty"`
	InventoryLevel                *InventoryLevel                `json:"inventory_level,omitempty"`
	CustomerAddress               *CustomerAddress               `json:"customer_address,omitempty"`
	CustomerGroup                 *CustomerGroup                 `json:"customer_group,omitempty"`
	CustomerGroupRelationship     *CustomerGroupRelationship     `json:"customer_group_relationship,omitempty"`
	ProductCollection             *ProductCollection             `json:"product_collection,omitempty"`
	ProductCollectionRelationship *ProductCollectionRelationship `json:"product_collection_relationship,omitempty"`
	ShipnowFulfillment            *ShipnowFulfillment            `json:"shipnow_fulfillment,omitempty"`
}

func (m *ChangeOneOf) String() string { return jsonx.MustMarshalToString(m) }
