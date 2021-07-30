package sms_provider

// +enum
// +enum:zero=null
type SmsProvider int

type NullSmsProvider struct {
	Enum  SmsProvider
	Valid bool
}

const (
	// +enum=unknown
	Unknown SmsProvider = 0

	// +enum=mock
	Mock SmsProvider = 1

	// +enum=telegram
	Telegram SmsProvider = 2

	// +enum=vietguys
	Vietguys SmsProvider = 3

	// +enum=incom
	Incom SmsProvider = 4
)
