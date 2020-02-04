package types

import (
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type CreateWebhookRequest struct {
	Entities []string `json:"entities"`
	Fields   []string `json:"fields"`
	Url      string   `json:"url"`
	Metadata string   `json:"metadata"`
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
	Id        dot.ID         `json:"id"`
	Entities  []string       `json:"entities"`
	Fields    []string       `json:"fields"`
	Url       string         `json:"url"`
	Metadata  string         `json:"metadata"`
	CreatedAt dot.Time       `json:"created_at"`
	States    *WebhookStates `json:"states"`
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
	Order                         *Order                         `json:"order"`
	Fulfillment                   *Fulfillment                   `json:"fulfillment"`
	Variant                       *ShopVariant                   `json:"variant"`
	InventoryLevel                *InventoryLevel                `json:"inventory_level"`
	CustomerAddress               *CustomerAddress               `json:"customer_address"`
	Customer                      *Customer                      `json:"customer"`
	CustomerGroup                 *CustomerGroup                 `json:"customer_group"`
	CustomerGroupRelationship     *CustomerGroupRelationship     `json:"customer_group_relationship"`
	Product                       *ShopProduct                   `json:"product"`
	ProductCollection             *ProductCollection             `json:"product_collection"`
	ProductCollectionRelationship *ProductCollectionRelationship `json:"product_collection_relationship"`
}

func (m *LatestOneOf) String() string { return jsonx.MustMarshalToString(m) }

type ChangeOneOf struct {
	Order                         *Order                         `json:"order"`
	Fulfillment                   *Fulfillment                   `json:"fulfillment"`
	Product                       *ShopProduct                   `json:"product"`
	Variant                       *ShopVariant                   `json:"variant"`
	Customer                      *Customer                      `json:"customer"`
	InventoryLevel                *InventoryLevel                `json:"inventory_level"`
	CustomerAddress               *CustomerAddress               `json:"customer_address"`
	CustomerGroup                 *CustomerGroup                 `json:"customer_group"`
	CustomerGroupRelationship     *CustomerGroupRelationship     `json:"customer_group_relationship"`
	ProductCollection             *ProductCollection             `json:"product_collection"`
	ProductCollectionRelationship *ProductCollectionRelationship `json:"product_collection_relationship"`
}

func (m *ChangeOneOf) String() string { return jsonx.MustMarshalToString(m) }
