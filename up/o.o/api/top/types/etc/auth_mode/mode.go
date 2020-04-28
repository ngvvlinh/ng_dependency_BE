package auth_mode

// +enum
type AuthMode int

type NullAuthMode struct {
	Enum  AuthMode
	Valid bool
}

const (
	// +enum=default
	Default AuthMode = 0

	// +enum=manual
	Manual AuthMode = 1
)
