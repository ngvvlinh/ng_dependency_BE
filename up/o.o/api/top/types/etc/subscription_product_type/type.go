package subscription_product_type

// +enum
// +enum:zero=null
type ProductSubscriptionType int

type NullProductSubscriptionType struct {
	Enum  ProductSubscriptionType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ProductSubscriptionType = 0

	// +enum=ecomify
	Ecomify ProductSubscriptionType = 1

	// +enum=telecom-extension
	TelecomExtension ProductSubscriptionType = 4
)
