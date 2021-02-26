package fb_customer_conversation_type

// +enum
// +enum:sql=int
// +enum:zero=null
type FbCustomerConversationType int

type NullFbCustomerConversationType struct {
	Enum  FbCustomerConversationType
	Valid bool
}

const (
	// +enum=unknown
	Unknown FbCustomerConversationType = 0

	// +enum=message
	Message FbCustomerConversationType = 872

	// +enum=comment
	Comment FbCustomerConversationType = 90

	// +enum=live_video
	LiveVideo FbCustomerConversationType = 176

	// +enum=all
	All FbCustomerConversationType = 585
)
