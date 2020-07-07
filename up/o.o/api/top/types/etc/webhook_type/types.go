package webhook_type

// +enum
// +enum:zero=null
type WebhookType int

type NullWebhookType struct {
	Enum  WebhookType
	Valid bool
}

const (
	// +enum=unknown
	Unknown WebhookType = 0

	// +enum=fabo
	Fabo WebhookType = 1
)
