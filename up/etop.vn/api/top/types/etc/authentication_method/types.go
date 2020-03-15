package authentication_method

// +enum
type AuthenticationMethod int

type NullAuthenticationMethod struct {
	Enum  AuthenticationMethod
	Valid bool
}

const (
	// +enum=email
	Email AuthenticationMethod = 1

	// +enum=phone
	Phone AuthenticationMethod = 2
)
